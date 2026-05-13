CREATE TABLE IF NOT EXISTS guest_sessions (
    guest_token TEXT PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE games
    ADD COLUMN IF NOT EXISTS guest_token TEXT NOT NULL DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_games_guest_token_updated_at
    ON games (guest_token, updated_at DESC);
