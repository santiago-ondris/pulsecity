CREATE TABLE IF NOT EXISTS team_franchises (
    game_id TEXT PRIMARY KEY,
    team_id TEXT NOT NULL,
    name TEXT NOT NULL,
    abbreviation TEXT NOT NULL,
    rating SMALLINT NOT NULL,
    offense_rating SMALLINT NOT NULL,
    defense_rating SMALLINT NOT NULL,
    pace SMALLINT NOT NULL,
    home_court_advantage SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_roster_players (
    player_id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL,
    full_name TEXT NOT NULL,
    position TEXT NOT NULL,
    overall_rating SMALLINT NOT NULL,
    roster_status TEXT NOT NULL,
    contract_years SMALLINT NOT NULL,
    salary INTEGER NOT NULL,
    sort_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_team_roster_players_game_id_sort_order
ON team_roster_players (game_id, sort_order);

CREATE TABLE IF NOT EXISTS team_opponents (
    game_id TEXT NOT NULL,
    team_id TEXT NOT NULL,
    name TEXT NOT NULL,
    abbreviation TEXT NOT NULL,
    rating SMALLINT NOT NULL,
    offense_rating SMALLINT NOT NULL,
    defense_rating SMALLINT NOT NULL,
    pace SMALLINT NOT NULL,
    home_court_advantage SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, team_id)
);

CREATE TABLE IF NOT EXISTS team_schedule (
    match_id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL,
    simulated_date TEXT NOT NULL,
    home_team_id TEXT NOT NULL,
    away_team_id TEXT NOT NULL,
    opponent_team_id TEXT NOT NULL,
    home_game BOOLEAN NOT NULL,
    seed BIGINT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE team_schedule ADD COLUMN IF NOT EXISTS home_score SMALLINT;
ALTER TABLE team_schedule ADD COLUMN IF NOT EXISTS away_score SMALLINT;
ALTER TABLE team_schedule ADD COLUMN IF NOT EXISTS winner_team_id TEXT;
ALTER TABLE team_schedule ADD COLUMN IF NOT EXISTS played_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_team_schedule_game_id_date
ON team_schedule (game_id, simulated_date);

CREATE TABLE IF NOT EXISTS team_player_box_scores (
    match_id TEXT NOT NULL,
    player_id TEXT NOT NULL,
    game_id TEXT NOT NULL,
    team_id TEXT NOT NULL,
    minutes SMALLINT NOT NULL,
    points SMALLINT NOT NULL,
    rebounds SMALLINT NOT NULL,
    assists SMALLINT NOT NULL,
    steals SMALLINT NOT NULL,
    blocks SMALLINT NOT NULL,
    turnovers SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (match_id, player_id)
);

CREATE INDEX IF NOT EXISTS idx_team_player_box_scores_game_id_player_id
ON team_player_box_scores (game_id, player_id);

CREATE TABLE IF NOT EXISTS team_records (
    game_id TEXT PRIMARY KEY,
    wins SMALLINT NOT NULL DEFAULT 0,
    losses SMALLINT NOT NULL DEFAULT 0,
    points_for INTEGER NOT NULL DEFAULT 0,
    points_against INTEGER NOT NULL DEFAULT 0,
    last_match_id TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
