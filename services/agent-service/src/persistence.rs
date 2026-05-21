use anyhow::{Context, Result};
use tokio_postgres::{Client, NoTls};

use crate::simulation::SimulationState;

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
        if let Some(state) = self.load_simulation_state(game_id).await? {
            return Ok(state);
        }

        let state = SimulationState::new(game_id);
        self.save_simulation_state(&state).await?;
        Ok(state)
    }
}
