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
