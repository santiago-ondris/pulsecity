CREATE TABLE IF NOT EXISTS games (
    game_id TEXT PRIMARY KEY,
    city_name TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    current_snapshot JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
