use anyhow::{Context, Result};
use tokio_postgres::{Client, NoTls};

use crate::{
    agents::{
        AgentRelationship, AgentRelationshipChange, AgentStateChange, CoreAgentState,
        IndividualAgentState, MatchAgentReactions, PlayerAgentState, RivalGMProfile,
        TeamRosterPlayer, apply_gm_decision_to_relationships, apply_match_finished,
        apply_match_to_player_agents, apply_match_to_relationships,
        apply_salary_cap_to_core_agents, default_agent_relationships, default_core_agent_state,
        default_individual_agent_states, default_player_agent_state, default_rival_gms,
    },
    events::{GMDecisionRegisteredEvent, MatchFinishedEvent, SalaryCapCalculatedEvent},
    simulation::SimulationState,
};

pub struct Store {
    client: Client,
}

#[derive(Debug, Clone, PartialEq)]
pub struct GMDecisionLogEntry {
    pub event_id: String,
    pub game_id: String,
    pub decision_id: String,
    pub kind: String,
    pub payload: std::collections::BTreeMap<String, String>,
    pub simulated_date: String,
    pub agents_affected: Vec<String>,
    pub source_event_id: Option<String>,
    pub source_subject: Option<String>,
    pub occurred_at: String,
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

CREATE TABLE IF NOT EXISTS agent_player_states (
    game_id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    full_name TEXT NOT NULL,
    position TEXT NOT NULL,
    emotional_state TEXT NOT NULL,
    satisfaction DOUBLE PRECISION NOT NULL,
    loyalty DOUBLE PRECISION NOT NULL,
    ego DOUBLE PRECISION NOT NULL,
    competitive_drive DOUBLE PRECISION NOT NULL,
    city_connection DOUBLE PRECISION NOT NULL,
    last_match_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, player_id)
);

CREATE TABLE IF NOT EXISTS agent_relationships (
    game_id TEXT NOT NULL,
    relationship_key TEXT NOT NULL,
    agent_a_id TEXT NOT NULL,
    agent_b_id TEXT NOT NULL,
    trust DOUBLE PRECISION NOT NULL,
    last_event TEXT NOT NULL,
    trend TEXT NOT NULL,
    short_history_json TEXT NOT NULL,
    last_source_event_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, relationship_key)
);

CREATE TABLE IF NOT EXISTS agent_relationship_event_hashes (
    game_id TEXT NOT NULL,
    relationship_key TEXT NOT NULL,
    source_event_id TEXT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, relationship_key, source_event_id)
);

CREATE TABLE IF NOT EXISTS gm_decisions_log (
    event_id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL,
    decision_id TEXT NOT NULL,
    kind TEXT NOT NULL,
    payload_json TEXT NOT NULL,
    simulated_date TEXT NOT NULL,
    agents_affected_json TEXT NOT NULL,
    source_event_id TEXT,
    source_subject TEXT,
    occurred_at TEXT NOT NULL,
    schema_version SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (game_id, decision_id)
);

CREATE INDEX IF NOT EXISTS idx_gm_decisions_log_game_created
ON gm_decisions_log (game_id, created_at DESC);

CREATE TABLE IF NOT EXISTS rival_gms (
    game_id TEXT NOT NULL,
    rival_team_id TEXT NOT NULL,
    gm_agent_id TEXT NOT NULL,
    display_name TEXT NOT NULL,
    team_name TEXT NOT NULL,
    negotiation_style TEXT NOT NULL,
    urgency_current DOUBLE PRECISION NOT NULL,
    build_philosophy TEXT NOT NULL,
    roster_needs_json TEXT NOT NULL,
    relationship_trust DOUBLE PRECISION NOT NULL,
    relationship_history_json TEXT NOT NULL,
    last_interaction_event_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, rival_team_id),
    UNIQUE (game_id, gm_agent_id)
);

CREATE INDEX IF NOT EXISTS idx_rival_gms_game_style
ON rival_gms (game_id, negotiation_style);

CREATE INDEX IF NOT EXISTS idx_rival_gms_game_urgency
ON rival_gms (game_id, urgency_current DESC);

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
        self.ensure_player_agents(game_id).await?;
        self.ensure_agent_relationships(game_id).await?;
        self.ensure_rival_gms(game_id).await?;
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

    pub async fn ensure_player_agents(&self, game_id: &str) -> Result<u64> {
        let players = self.load_team_roster_players(game_id).await?;
        let mut inserted = 0;

        for player in players {
            let state = default_player_agent_state(&player);
            inserted += self.insert_player_agent_if_missing(&state).await?;
        }

        Ok(inserted)
    }

    pub async fn count_player_agents(&self, game_id: &str) -> Result<i64> {
        let row = self
            .client
            .query_one(
                "
SELECT COUNT(*)
FROM agent_player_states
WHERE game_id = $1;
",
                &[&game_id],
            )
            .await
            .with_context(|| "count player agent states")?;

        Ok(row.get(0))
    }

    pub async fn ensure_agent_relationships(&self, game_id: &str) -> Result<u64> {
        let relationships = default_agent_relationships(game_id);
        let mut inserted = 0;

        for relationship in relationships {
            inserted += self
                .insert_agent_relationship_if_missing(&relationship)
                .await?;
        }

        Ok(inserted)
    }

    pub async fn count_agent_relationships(&self, game_id: &str) -> Result<i64> {
        let row = self
            .client
            .query_one(
                "
SELECT COUNT(*)
FROM agent_relationships
WHERE game_id = $1;
",
                &[&game_id],
            )
            .await
            .with_context(|| "count agent relationships")?;

        Ok(row.get(0))
    }

    pub async fn ensure_rival_gms(&self, game_id: &str) -> Result<u64> {
        let rival_gms = default_rival_gms(game_id);
        let mut inserted = 0;

        for rival_gm in rival_gms {
            inserted += self.insert_rival_gm_if_missing(&rival_gm).await?;
        }

        Ok(inserted)
    }

    pub async fn count_rival_gms(&self, game_id: &str) -> Result<i64> {
        let row = self
            .client
            .query_one(
                "
SELECT COUNT(*)
FROM rival_gms
WHERE game_id = $1;
",
                &[&game_id],
            )
            .await
            .with_context(|| "count rival gms")?;

        Ok(row.get(0))
    }

    pub async fn record_gm_decision(&self, event: &GMDecisionRegisteredEvent) -> Result<bool> {
        let payload_json = serde_json::to_string(&event.payload)
            .with_context(|| "encode gm decision payload json")?;
        let agents_affected_json = serde_json::to_string(&event.agents_affected)
            .with_context(|| "encode gm decision agents affected json")?;
        let schema_version = i16::try_from(event.meta.schema_version).unwrap_or(i16::MAX);

        let inserted = self
            .client
            .execute(
                "
INSERT INTO gm_decisions_log (
    event_id,
    game_id,
    decision_id,
    kind,
    payload_json,
    simulated_date,
    agents_affected_json,
    source_event_id,
    source_subject,
    occurred_at,
    schema_version,
    created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
ON CONFLICT DO NOTHING;
",
                &[
                    &event.meta.event_id,
                    &event.meta.game_id,
                    &event.decision_id,
                    &event.kind,
                    &payload_json,
                    &event.simulated_date,
                    &agents_affected_json,
                    &event.source_event_id,
                    &event.source_subject,
                    &event.meta.occurred_at,
                    &schema_version,
                ],
            )
            .await
            .with_context(|| "record gm decision")?;

        Ok(inserted > 0)
    }

    pub async fn apply_gm_decision_relationships(
        &mut self,
        event: &GMDecisionRegisteredEvent,
        occurred_at: String,
    ) -> Result<Vec<AgentRelationshipChange>> {
        self.ensure_agent_relationships(&event.meta.game_id).await?;

        let transaction = self
            .client
            .transaction()
            .await
            .with_context(|| "begin gm decision relationship reaction")?;

        let relationship_states =
            load_agent_relationships(&transaction, &event.meta.game_id).await?;
        let relationship_changes =
            apply_gm_decision_to_relationships(relationship_states, event, occurred_at);
        let mut persisted_changes = Vec::with_capacity(relationship_changes.len());
        for change in relationship_changes {
            let key = crate::agents::relationship_key(
                &change.relationship.agent_a_id,
                &change.relationship.agent_b_id,
            );
            if mark_relationship_event_processed(
                &transaction,
                &event.meta.game_id,
                &key,
                &event.meta.event_id,
            )
            .await?
            {
                save_agent_relationship(&transaction, &change.relationship).await?;
                persisted_changes.push(change);
            }
        }

        transaction
            .commit()
            .await
            .with_context(|| "commit gm decision relationship reaction")?;

        Ok(persisted_changes)
    }

    pub async fn apply_salary_cap_calculated(
        &mut self,
        event: &SalaryCapCalculatedEvent,
        occurred_at: String,
    ) -> Result<Vec<AgentStateChange>> {
        let transaction = self
            .client
            .transaction()
            .await
            .with_context(|| "begin salary cap agent reaction")?;

        let processed = transaction
            .execute(
                "
INSERT INTO agent_processed_matches (game_id, match_id, source_event_id, processed_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (game_id, match_id) DO NOTHING;
",
                &[
                    &event.meta.game_id,
                    &format!("salary-cap-{}", event.meta.event_id),
                    &event.meta.event_id,
                ],
            )
            .await
            .with_context(|| "mark salary cap agent reaction processed")?;
        if processed == 0 {
            transaction
                .rollback()
                .await
                .with_context(|| "rollback already processed salary cap reaction")?;
            return Ok(Vec::new());
        }

        let current_states = vec![
            load_core_agent_state(&transaction, &event.meta.game_id, "owner").await?,
            load_core_agent_state(&transaction, &event.meta.game_id, "cfo").await?,
        ];
        let changes = apply_salary_cap_to_core_agents(current_states, event, occurred_at);
        for change in &changes {
            save_core_agent_state(&transaction, &change.state).await?;
        }

        transaction
            .commit()
            .await
            .with_context(|| "commit salary cap agent reaction")?;

        Ok(changes)
    }

    pub async fn latest_gm_decisions(
        &self,
        game_id: &str,
        limit: i64,
    ) -> Result<Vec<GMDecisionLogEntry>> {
        let bounded_limit = limit.clamp(1, 50);
        let rows = self
            .client
            .query(
                "
SELECT
    event_id,
    game_id,
    decision_id,
    kind,
    payload_json,
    simulated_date,
    agents_affected_json,
    source_event_id,
    source_subject,
    occurred_at
FROM gm_decisions_log
WHERE game_id = $1
ORDER BY created_at DESC
LIMIT $2;
",
                &[&game_id, &bounded_limit],
            )
            .await
            .with_context(|| "load latest gm decisions")?;

        let mut decisions = Vec::with_capacity(rows.len());
        for row in rows {
            let payload_json: String = row.get("payload_json");
            let agents_affected_json: String = row.get("agents_affected_json");
            decisions.push(GMDecisionLogEntry {
                event_id: row.get("event_id"),
                game_id: row.get("game_id"),
                decision_id: row.get("decision_id"),
                kind: row.get("kind"),
                payload: serde_json::from_str(&payload_json)
                    .with_context(|| "decode gm decision payload json")?,
                simulated_date: row.get("simulated_date"),
                agents_affected: serde_json::from_str(&agents_affected_json)
                    .with_context(|| "decode gm decision agents affected json")?,
                source_event_id: row.get("source_event_id"),
                source_subject: row.get("source_subject"),
                occurred_at: row.get("occurred_at"),
            });
        }

        Ok(decisions)
    }

    async fn load_team_roster_players(&self, game_id: &str) -> Result<Vec<TeamRosterPlayer>> {
        let rows = self
            .client
            .query(
                "
SELECT player_id, game_id, full_name, position, overall_rating, roster_status
FROM team_roster_players
WHERE game_id = $1
ORDER BY sort_order ASC;
",
                &[&game_id],
            )
            .await
            .with_context(|| "load team roster players")?;

        let mut players = Vec::with_capacity(rows.len());
        for row in rows {
            let overall_rating: i16 = row.get("overall_rating");
            players.push(TeamRosterPlayer {
                player_id: row.get("player_id"),
                game_id: row.get("game_id"),
                full_name: row.get("full_name"),
                position: row.get("position"),
                overall_rating: overall_rating as u8,
                roster_status: row.get("roster_status"),
            });
        }

        Ok(players)
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

    async fn insert_player_agent_if_missing(&self, state: &PlayerAgentState) -> Result<u64> {
        self.client
            .execute(
                "
INSERT INTO agent_player_states (
    game_id,
    player_id,
    full_name,
    position,
    emotional_state,
    satisfaction,
    loyalty,
    ego,
    competitive_drive,
    city_connection,
    last_match_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
ON CONFLICT (game_id, player_id) DO NOTHING;
",
                &[
                    &state.game_id,
                    &state.player_id,
                    &state.full_name,
                    &state.position,
                    &state.emotional_state,
                    &state.satisfaction,
                    &state.loyalty,
                    &state.ego,
                    &state.competitive_drive,
                    &state.city_connection,
                    &state.last_match_id,
                ],
            )
            .await
            .with_context(|| "insert player agent state")
    }

    async fn insert_agent_relationship_if_missing(
        &self,
        relationship: &AgentRelationship,
    ) -> Result<u64> {
        let short_history_json = serde_json::to_string(&relationship.short_history)
            .with_context(|| "encode relationship short history json")?;
        let key =
            crate::agents::relationship_key(&relationship.agent_a_id, &relationship.agent_b_id);

        self.client
            .execute(
                "
INSERT INTO agent_relationships (
    game_id,
    relationship_key,
    agent_a_id,
    agent_b_id,
    trust,
    last_event,
    trend,
    short_history_json,
    last_source_event_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
ON CONFLICT (game_id, relationship_key) DO NOTHING;
",
                &[
                    &relationship.game_id,
                    &key,
                    &relationship.agent_a_id,
                    &relationship.agent_b_id,
                    &relationship.trust,
                    &relationship.last_event,
                    &relationship.trend,
                    &short_history_json,
                    &relationship.last_source_event_id,
                ],
            )
            .await
            .with_context(|| "insert agent relationship")
    }

    async fn insert_rival_gm_if_missing(&self, rival_gm: &RivalGMProfile) -> Result<u64> {
        let roster_needs_json = serde_json::to_string(&rival_gm.roster_needs)
            .with_context(|| "encode rival gm roster needs json")?;
        let relationship_history_json = serde_json::to_string(&rival_gm.relationship_history)
            .with_context(|| "encode rival gm relationship history json")?;

        self.client
            .execute(
                "
INSERT INTO rival_gms (
    game_id,
    rival_team_id,
    gm_agent_id,
    display_name,
    team_name,
    negotiation_style,
    urgency_current,
    build_philosophy,
    roster_needs_json,
    relationship_trust,
    relationship_history_json,
    last_interaction_event_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
ON CONFLICT (game_id, rival_team_id) DO NOTHING;
",
                &[
                    &rival_gm.game_id,
                    &rival_gm.rival_team_id,
                    &rival_gm.gm_agent_id,
                    &rival_gm.display_name,
                    &rival_gm.team_name,
                    &rival_gm.negotiation_style,
                    &rival_gm.urgency_current,
                    &rival_gm.build_philosophy,
                    &roster_needs_json,
                    &rival_gm.relationship_trust,
                    &relationship_history_json,
                    &rival_gm.last_interaction_event_id,
                ],
            )
            .await
            .with_context(|| "insert rival gm")
    }

    pub async fn apply_match_finished(
        &mut self,
        event: &MatchFinishedEvent,
        occurred_at: String,
    ) -> Result<Option<MatchAgentReactions>> {
        self.ensure_player_agents(&event.meta.game_id).await?;

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

        let changes = apply_match_finished(current_states, event, occurred_at.clone());
        for change in &changes {
            save_core_agent_state(&transaction, &change.state).await?;
        }
        let player_states = load_player_agent_states(&transaction, &event.meta.game_id).await?;
        let (updated_player_states, roster_patch) =
            apply_match_to_player_agents(player_states, event);
        for state in &updated_player_states {
            save_player_agent_state(&transaction, state).await?;
        }
        let relationship_states =
            load_agent_relationships(&transaction, &event.meta.game_id).await?;
        let relationship_changes =
            apply_match_to_relationships(relationship_states, event, occurred_at.clone());
        let mut persisted_relationship_changes = Vec::with_capacity(relationship_changes.len());
        for change in relationship_changes {
            let key = crate::agents::relationship_key(
                &change.relationship.agent_a_id,
                &change.relationship.agent_b_id,
            );
            if mark_relationship_event_processed(
                &transaction,
                &event.meta.game_id,
                &key,
                &event.meta.event_id,
            )
            .await?
            {
                save_agent_relationship(&transaction, &change.relationship).await?;
                persisted_relationship_changes.push(change);
            }
        }

        transaction
            .commit()
            .await
            .with_context(|| "commit agent match reaction")?;

        Ok(Some(MatchAgentReactions {
            core_agent_changes: changes,
            roster_patch,
            relationship_changes: persisted_relationship_changes,
        }))
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

async fn load_player_agent_states(
    transaction: &tokio_postgres::Transaction<'_>,
    game_id: &str,
) -> Result<Vec<PlayerAgentState>> {
    let rows = transaction
        .query(
            "
SELECT
    game_id,
    player_id,
    full_name,
    position,
    emotional_state,
    satisfaction,
    loyalty,
    ego,
    competitive_drive,
    city_connection,
    last_match_id
FROM agent_player_states
WHERE game_id = $1
ORDER BY player_id ASC;
",
            &[&game_id],
        )
        .await
        .with_context(|| "load player agent states")?;

    let mut states = Vec::with_capacity(rows.len());
    for row in rows {
        states.push(PlayerAgentState {
            game_id: row.get("game_id"),
            player_id: row.get("player_id"),
            full_name: row.get("full_name"),
            position: row.get("position"),
            emotional_state: row.get("emotional_state"),
            satisfaction: row.get("satisfaction"),
            loyalty: row.get("loyalty"),
            ego: row.get("ego"),
            competitive_drive: row.get("competitive_drive"),
            city_connection: row.get("city_connection"),
            last_match_id: row.get("last_match_id"),
        });
    }

    Ok(states)
}

async fn save_player_agent_state(
    transaction: &tokio_postgres::Transaction<'_>,
    state: &PlayerAgentState,
) -> Result<()> {
    transaction
        .execute(
            "
INSERT INTO agent_player_states (
    game_id,
    player_id,
    full_name,
    position,
    emotional_state,
    satisfaction,
    loyalty,
    ego,
    competitive_drive,
    city_connection,
    last_match_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
ON CONFLICT (game_id, player_id) DO UPDATE SET
    full_name = EXCLUDED.full_name,
    position = EXCLUDED.position,
    emotional_state = EXCLUDED.emotional_state,
    satisfaction = EXCLUDED.satisfaction,
    loyalty = EXCLUDED.loyalty,
    ego = EXCLUDED.ego,
    competitive_drive = EXCLUDED.competitive_drive,
    city_connection = EXCLUDED.city_connection,
    last_match_id = EXCLUDED.last_match_id,
    updated_at = NOW();
",
            &[
                &state.game_id,
                &state.player_id,
                &state.full_name,
                &state.position,
                &state.emotional_state,
                &state.satisfaction,
                &state.loyalty,
                &state.ego,
                &state.competitive_drive,
                &state.city_connection,
                &state.last_match_id,
            ],
        )
        .await
        .with_context(|| "save player agent state")?;

    Ok(())
}

async fn load_agent_relationships(
    transaction: &tokio_postgres::Transaction<'_>,
    game_id: &str,
) -> Result<Vec<AgentRelationship>> {
    let rows = transaction
        .query(
            "
SELECT
    game_id,
    agent_a_id,
    agent_b_id,
    trust,
    last_event,
    trend,
    short_history_json,
    last_source_event_id
FROM agent_relationships
WHERE game_id = $1
ORDER BY relationship_key ASC;
",
            &[&game_id],
        )
        .await
        .with_context(|| "load agent relationships")?;

    let mut relationships = Vec::with_capacity(rows.len());
    for row in rows {
        let short_history_json: String = row.get("short_history_json");
        let short_history = serde_json::from_str(&short_history_json)
            .with_context(|| "decode relationship short history json")?;
        relationships.push(AgentRelationship {
            game_id: row.get("game_id"),
            agent_a_id: row.get("agent_a_id"),
            agent_b_id: row.get("agent_b_id"),
            trust: row.get("trust"),
            last_event: row.get("last_event"),
            trend: row.get("trend"),
            short_history,
            last_source_event_id: row.get("last_source_event_id"),
        });
    }

    Ok(relationships)
}

async fn mark_relationship_event_processed(
    transaction: &tokio_postgres::Transaction<'_>,
    game_id: &str,
    relationship_key: &str,
    source_event_id: &str,
) -> Result<bool> {
    let processed = transaction
        .execute(
            "
INSERT INTO agent_relationship_event_hashes (
    game_id,
    relationship_key,
    source_event_id,
    processed_at
)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (game_id, relationship_key, source_event_id) DO NOTHING;
",
            &[&game_id, &relationship_key, &source_event_id],
        )
        .await
        .with_context(|| "mark relationship event processed")?;

    Ok(processed > 0)
}

async fn save_agent_relationship(
    transaction: &tokio_postgres::Transaction<'_>,
    relationship: &AgentRelationship,
) -> Result<()> {
    let short_history_json = serde_json::to_string(&relationship.short_history)
        .with_context(|| "encode relationship short history json")?;
    let key = crate::agents::relationship_key(&relationship.agent_a_id, &relationship.agent_b_id);

    transaction
        .execute(
            "
INSERT INTO agent_relationships (
    game_id,
    relationship_key,
    agent_a_id,
    agent_b_id,
    trust,
    last_event,
    trend,
    short_history_json,
    last_source_event_id,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
ON CONFLICT (game_id, relationship_key) DO UPDATE SET
    agent_a_id = EXCLUDED.agent_a_id,
    agent_b_id = EXCLUDED.agent_b_id,
    trust = EXCLUDED.trust,
    last_event = EXCLUDED.last_event,
    trend = EXCLUDED.trend,
    short_history_json = EXCLUDED.short_history_json,
    last_source_event_id = EXCLUDED.last_source_event_id,
    updated_at = NOW();
",
            &[
                &relationship.game_id,
                &key,
                &relationship.agent_a_id,
                &relationship.agent_b_id,
                &relationship.trust,
                &relationship.last_event,
                &relationship.trend,
                &short_history_json,
                &relationship.last_source_event_id,
            ],
        )
        .await
        .with_context(|| "save agent relationship")?;

    Ok(())
}
