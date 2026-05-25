CREATE TABLE IF NOT EXISTS analytics_match_results (
    game_id TEXT NOT NULL,
    match_id TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    simulated_date TEXT NOT NULL,
    home_team_id TEXT NOT NULL,
    away_team_id TEXT NOT NULL,
    home_score SMALLINT NOT NULL,
    away_score SMALLINT NOT NULL,
    winner_team_id TEXT NOT NULL,
    seed BIGINT NOT NULL,
    source_event_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, match_id)
);

CREATE TABLE IF NOT EXISTS analytics_player_box_scores (
    game_id TEXT NOT NULL,
    match_id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    team_id TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    simulated_date TEXT NOT NULL,
    minutes SMALLINT NOT NULL,
    points SMALLINT NOT NULL,
    rebounds SMALLINT NOT NULL,
    assists SMALLINT NOT NULL,
    steals SMALLINT NOT NULL,
    blocks SMALLINT NOT NULL,
    turnovers SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, match_id, player_id)
);

CREATE INDEX IF NOT EXISTS idx_analytics_player_box_scores_player_time
ON analytics_player_box_scores (game_id, player_id, occurred_at DESC);

CREATE TABLE IF NOT EXISTS analytics_city_metric_points (
    event_id TEXT NOT NULL,
    metric TEXT NOT NULL,
    game_id TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    simulated_date TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    delta DOUBLE PRECISION NOT NULL,
    source_event_id TEXT NOT NULL,
    source_subject TEXT NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (event_id, metric)
);

CREATE INDEX IF NOT EXISTS idx_analytics_city_metric_points_game_metric_time
ON analytics_city_metric_points (game_id, metric, occurred_at DESC);

CREATE TABLE IF NOT EXISTS analytics_land_value_points (
    event_id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    simulated_date TEXT NOT NULL,
    zone_id TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    delta DOUBLE PRECISION NOT NULL,
    source_event_id TEXT NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analytics_land_value_points_game_zone_time
ON analytics_land_value_points (game_id, zone_id, occurred_at DESC);

CREATE TABLE IF NOT EXISTS analytics_agent_state_points (
    event_id TEXT NOT NULL,
    metric TEXT NOT NULL,
    game_id TEXT NOT NULL,
    agent_id TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    simulated_date TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    mood TEXT NOT NULL,
    source_event_id TEXT NOT NULL,
    source_subject TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (event_id, metric)
);

CREATE INDEX IF NOT EXISTS idx_analytics_agent_state_points_game_agent_time
ON analytics_agent_state_points (game_id, agent_id, occurred_at DESC);
