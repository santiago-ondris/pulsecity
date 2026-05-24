CREATE TABLE IF NOT EXISTS city_metrics (
    game_id TEXT PRIMARY KEY,
    fan_sentiment DOUBLE PRECISION NOT NULL,
    ticket_sales_index DOUBLE PRECISION NOT NULL,
    local_economy_index DOUBLE PRECISION NOT NULL,
    stadium_district_land_value DOUBLE PRECISION NOT NULL,
    win_streak SMALLINT NOT NULL,
    loss_streak SMALLINT NOT NULL,
    last_match_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS city_processed_matches (
    game_id TEXT NOT NULL,
    match_id TEXT NOT NULL,
    source_event_id TEXT NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (game_id, match_id)
);
