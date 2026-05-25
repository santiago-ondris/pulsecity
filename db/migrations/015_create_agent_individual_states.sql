CREATE TABLE IF NOT EXISTS agent_individual_states (
    game_id TEXT NOT NULL,
    agent_id TEXT NOT NULL,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL,
    role TEXT NOT NULL,
    domain TEXT NOT NULL,
    emotional_state TEXT NOT NULL,
    confidence DOUBLE PRECISION NOT NULL,
    satisfaction DOUBLE PRECISION NOT NULL,
    loyalty DOUBLE PRECISION NOT NULL,
    role_performance DOUBLE PRECISION NOT NULL,
    state_json TEXT NOT NULL,
    agenda_json TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, agent_id)
);

CREATE INDEX IF NOT EXISTS idx_agent_individual_states_game_category
ON agent_individual_states (game_id, category);
