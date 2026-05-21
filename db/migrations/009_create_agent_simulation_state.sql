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
