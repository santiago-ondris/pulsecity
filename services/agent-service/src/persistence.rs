use anyhow::{Context, Result};
use tokio_postgres::{Client, NoTls};

use crate::{
    agents::{
        AgentStateChange, CoreAgentState, IndividualAgentState, apply_match_finished,
        default_core_agent_state, default_individual_agent_states,
    },
    events::MatchFinishedEvent,
    simulation::SimulationState,
};

pub struct Store {
    client: Client,
}

impl Store {
    pub async fn connect(database_url: &str) -> Result<Self> {
        let (client, connection) = tokio_postgres::connect(database_url, NoTls)
            .await
            .with_context(|| "connect postgres")?;

        tokio::spawn(async move {
            if let Err(err) = connection.await {
                tracing::error!("postgres connection error: {err}");
            }
        });

        Ok(Self { client })
    }

    pub async fn ensure_schema(&self) -> Result<()> {
        self.client
            .batch_execute(
                "
CREATE TABLE IF NOT EXISTS agent_simulation_state (
    game_id TEXT PRIMARY KEY,
    current_simulated_date TEXT NOT NULL,
    speed SMALLINT NOT NULL,
    paused BOOLEAN NOT NULL,
    session_active BOOLEAN NOT NULL,
    last_tick_processed_at TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS agent_core_states (
    game_id TEXT NOT NULL,
    agent_id TEXT NOT NULL,
    mood TEXT NOT NULL,
    state_json TEXT NOT NULL,
    last_match_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, agent_id)
);

CREATE TABLE IF NOT EXISTS agent_individual_states (
    game_id TEXT NOT NULL,
    agent_id TEXT NOT NULL,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL,
    role TEXT NOT NULL,
    domain TEXT NOT NULL,
    emotional_state TEXT NOT NULL,
    confidence DOUBLE PRECISION NOT NULL,
    satisfaction DOUBLE PRECISION NOT NULL,
    loyalty DOUBLE PRECISION NOT NULL,
    role_performance DOUBLE PRECISION NOT NULL,
    state_json TEXT NOT NULL,
    agenda_json TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, agent_id)
);

CREATE INDEX IF NOT EXISTS idx_agent_individual_states_game_category
ON agent_individual_states (game_id, category);

CREATE TABLE IF NOT EXISTS agent_processed_matches (
    game_id TEXT NOT NULL,
    match_id TEXT NOT NULL,
    source_event_id TEXT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, match_id)
);
",
            )
            .await
            .with_context(|| "ensure agent simulation state schema")?;

        Ok(())
    }

    pub async fn load_simulation_state(&self, game_id: &str) -> Result<Option<SimulationState>> {
        let row = self
            .client
            .query_opt(
                "
SELECT
    game_id,
    current_simulated_date,
    speed,
    paused,
    session_active,
    last_tick_processed_at
FROM agent_simulation_state
WHERE game_id = $1;
",
                &[&game_id],
            )
            .await
            .with_context(|| "load simulation state")?;

        let Some(row) = row else {
            return Ok(None);
        };

        let speed: i16 = row.get("speed");
        Ok(Some(SimulationState {
            game_id: row.get("game_id"),
            current_simulated_date: row.get("current_simulated_date"),
            speed: speed as u8,
            paused: row.get("paused"),
            session_active: row.get("session_active"),
            last_tick_processed_at: row.get("last_tick_processed_at"),
        }))
    }

    pub async fn save_simulation_state(&self, state: &SimulationState) -> Result<()> {
        let speed = i16::from(state.speed);
        self.client
            .execute(
                "
INSERT INTO agent_simulation_state (
    game_id,
    current_simulated_date,
    speed,
    paused,
    session_active,
    last_tick_processed_at,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
ON CONFLICT (game_id) DO UPDATE SET
    current_simulated_date = EXCLUDED.current_simulated_date,
    speed = EXCLUDED.speed,
    paused = EXCLUDED.paused,
    session_active = EXCLUDED.session_active,
    last_tick_processed_at = EXCLUDED.last_tick_processed_at,
    updated_at = NOW();
",
                &[
                    &state.game_id,
                    &state.current_simulated_date,
                    &speed,
                    &state.paused,
                    &state.session_active,
                    &state.last_tick_processed_at,
                ],
            )
            .await
            .with_context(|| "save simulation state")?;

        Ok(())
    }

    pub async fn load_or_initialize_simulation_state(
        &self,
        game_id: &str,
    ) -> Result<SimulationState> {
        self.ensure_individual_agents(game_id).await?;
        if let Some(state) = self.load_simulation_state(game_id).await? {
            return Ok(state);
        }

        let state = SimulationState::new(game_id);
        self.save_simulation_state(&state).await?;
        Ok(state)
    }

    pub async fn ensure_individual_agents(&self, game_id: &str) -> Result<u64> {
        let agents = default_individual_agent_states(game_id);
        let mut inserted = 0;

        for agent in agents {
            inserted += self.insert_individual_agent_if_missing(&agent).await?;
        }

        Ok(inserted)
    }

    pub async fn count_individual_agents(&self, game_id: &str) -> Result<i64> {
        let row = self
            .client
            .query_one(
                "
SELECT COUNT(*)
FROM agent_individual_states
WHERE game_id = $1;
",
                &[&game_id],
            )
            .await
            .with_context(|| "count individual agent states")?;

        Ok(row.get(0))
    }

    async fn insert_individual_agent_if_missing(
        &self,
        agent: &IndividualAgentState,
    ) -> Result<u64> {
        let state_json =
            serde_json::to_string(&agent.state).with_context(|| "encode individual state json")?;
        let agenda_json = serde_json::to_string(&agent.agenda)
            .with_context(|| "encode individual agenda json")?;

        self.client
            .execute(
                "
INSERT INTO agent_individual_states (
    game_id,
    agent_id,
    display_name,
    category,
    role,
    domain,
    emotional_state,
    confidence,
    satisfaction,
    loyalty,
    role_performance,
    state_json,
    agenda_json,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
ON CONFLICT (game_id, agent_id) DO NOTHING;
",
                &[
                    &agent.game_id,
                    &agent.agent_id,
                    &agent.display_name,
                    &agent.category,
                    &agent.role,
                    &agent.domain,
                    &agent.emotional_state,
                    &agent.confidence,
                    &agent.satisfaction,
                    &agent.loyalty,
                    &agent.role_performance,
                    &state_json,
                    &agenda_json,
                ],
            )
            .await
            .with_context(|| "insert individual agent state")
    }

    pub async fn apply_match_finished(
        &mut self,
        event: &MatchFinishedEvent,
        occurred_at: String,
    ) -> Result<Option<Vec<AgentStateChange>>> {
        let transaction = self
            .client
            .transaction()
            .await
            .with_context(|| "begin agent match reaction")?;

        let processed = transaction
            .execute(
                "
INSERT INTO agent_processed_matches (game_id, match_id, source_event_id, processed_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (game_id, match_id) DO NOTHING;
",
                &[&event.meta.game_id, &event.match_id, &event.meta.event_id],
            )
            .await
            .with_context(|| "mark agent match processed")?;

        if processed == 0 {
            transaction
                .rollback()
                .await
                .with_context(|| "rollback already processed agent match")?;
            return Ok(None);
        }

        let mut current_states = Vec::with_capacity(crate::agents::CORE_AGENT_IDS.len());
        for agent_id in crate::agents::CORE_AGENT_IDS {
            current_states
                .push(load_core_agent_state(&transaction, &event.meta.game_id, agent_id).await?);
        }

        let changes = apply_match_finished(current_states, event, occurred_at);
        for change in &changes {
            save_core_agent_state(&transaction, &change.state).await?;
        }

        transaction
            .commit()
            .await
            .with_context(|| "commit agent match reaction")?;

        Ok(Some(changes))
    }
}

async fn load_core_agent_state(
    transaction: &tokio_postgres::Transaction<'_>,
    game_id: &str,
    agent_id: &str,
) -> Result<CoreAgentState> {
    let row = transaction
        .query_opt(
            "
SELECT game_id, agent_id, mood, state_json, last_match_id
FROM agent_core_states
WHERE game_id = $1 AND agent_id = $2;
",
            &[&game_id, &agent_id],
        )
        .await
        .with_context(|| "load core agent state")?;

    let Some(row) = row else {
        return Ok(default_core_agent_state(game_id, agent_id));
    };

    let state_json: String = row.get("state_json");
    let state =
        serde_json::from_str(&state_json).with_context(|| "decode core agent state json")?;

    Ok(CoreAgentState {
        game_id: row.get("game_id"),
        agent_id: row.get("agent_id"),
        mood: row.get("mood"),
        state,
        last_match_id: row.get("last_match_id"),
    })
}

async fn save_core_agent_state(
    transaction: &tokio_postgres::Transaction<'_>,
    state: &CoreAgentState,
) -> Result<()> {
    let state_json =
        serde_json::to_string(&state.state).with_context(|| "encode core agent state json")?;
    transaction
        .execute(
            "
INSERT INTO agent_core_states (
    game_id,
    agent_id,
    mood,
    state_json,
    last_match_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT (game_id, agent_id) DO UPDATE SET
    mood = EXCLUDED.mood,
    state_json = EXCLUDED.state_json,
    last_match_id = EXCLUDED.last_match_id,
    updated_at = NOW();
",
            &[
                &state.game_id,
                &state.agent_id,
                &state.mood,
                &state_json,
                &state.last_match_id,
            ],
        )
        .await
        .with_context(|| "save core agent state")?;

    Ok(())
}
