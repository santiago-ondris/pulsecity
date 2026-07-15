CREATE TABLE IF NOT EXISTS team_trades (
    proposal_id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL,
    rival_team_id TEXT NOT NULL,
    offered_player_id TEXT NOT NULL,
    offered_player_name TEXT NOT NULL,
    offered_salary INTEGER NOT NULL,
    requested_position TEXT NOT NULL,
    incoming_salary INTEGER NOT NULL,
    cap_space_after INTEGER NOT NULL,
    status TEXT NOT NULL,
    rejection_reason TEXT,
    rejection_detail TEXT,
    incoming_player_id TEXT,
    incoming_player_name TEXT,
    incoming_rating SMALLINT,
    accepted_additional_asset TEXT,
    accepted_at TEXT,
    source_decision_id TEXT NOT NULL,
    source_event_id TEXT NOT NULL,
    simulated_date TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_team_trades_game_status
ON team_trades (game_id, status);
