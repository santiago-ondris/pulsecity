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
