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

CREATE TABLE IF NOT EXISTS agent_processed_matches (
    game_id TEXT NOT NULL,
    match_id TEXT NOT NULL,
    source_event_id TEXT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, match_id)
);
