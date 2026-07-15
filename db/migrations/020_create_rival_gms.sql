CREATE TABLE IF NOT EXISTS rival_gms (
    game_id TEXT NOT NULL,
    rival_team_id TEXT NOT NULL,
    gm_agent_id TEXT NOT NULL,
    display_name TEXT NOT NULL,
    team_name TEXT NOT NULL,
    negotiation_style TEXT NOT NULL,
    urgency_current DOUBLE PRECISION NOT NULL,
    build_philosophy TEXT NOT NULL,
    roster_needs_json TEXT NOT NULL,
    relationship_trust DOUBLE PRECISION NOT NULL,
    relationship_history_json TEXT NOT NULL,
    last_interaction_event_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, rival_team_id),
    UNIQUE (game_id, gm_agent_id)
);

CREATE INDEX IF NOT EXISTS idx_rival_gms_game_style
ON rival_gms (game_id, negotiation_style);

CREATE INDEX IF NOT EXISTS idx_rival_gms_game_urgency
ON rival_gms (game_id, urgency_current DESC);
