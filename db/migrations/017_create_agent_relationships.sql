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
