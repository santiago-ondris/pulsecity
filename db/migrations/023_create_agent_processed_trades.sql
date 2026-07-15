CREATE TABLE IF NOT EXISTS agent_processed_trades (
    game_id TEXT NOT NULL,
    proposal_id TEXT NOT NULL,
    source_event_id TEXT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, proposal_id)
);
