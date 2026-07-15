CREATE TABLE IF NOT EXISTS rival_gm_trade_evaluations (
    proposal_id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL,
    rival_team_id TEXT NOT NULL,
    result_subject TEXT NOT NULL,
    result_event_id TEXT NOT NULL,
    evaluated_at TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
