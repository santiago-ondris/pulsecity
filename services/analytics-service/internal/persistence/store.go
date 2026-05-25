package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulsecity/services/analytics-service/internal/domain"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, databaseURL string) (*Store, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) EnsureSchema(ctx context.Context) error {
	const query = `
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
`

	if _, err := s.pool.Exec(ctx, query); err != nil {
		return fmt.Errorf("ensure analytics schema: %w", err)
	}

	return nil
}

func (s *Store) IngestMatchFinished(ctx context.Context, event domain.MatchFinishedEvent) error {
	occurredAt := parseEventTime(event.OccurredAt)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin match analytics ingest: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
INSERT INTO analytics_match_results (
	game_id, match_id, occurred_at, simulated_date, home_team_id, away_team_id,
	home_score, away_score, winner_team_id, seed, source_event_id, created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
ON CONFLICT (game_id, match_id) DO NOTHING;
`,
		event.GameID,
		event.MatchID,
		occurredAt,
		event.SimulatedDate,
		event.HomeTeam.TeamID,
		event.AwayTeam.TeamID,
		int16(event.HomeScore),
		int16(event.AwayScore),
		event.WinnerTeamID,
		int64(event.Seed),
		event.EventID,
	); err != nil {
		return fmt.Errorf("insert match result analytics: %w", err)
	}

	for _, line := range event.BoxScore {
		if _, err := tx.Exec(ctx, `
INSERT INTO analytics_player_box_scores (
	game_id, match_id, player_id, team_id, occurred_at, simulated_date,
	minutes, points, rebounds, assists, steals, blocks, turnovers, created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW())
ON CONFLICT (game_id, match_id, player_id) DO NOTHING;
`,
			event.GameID,
			event.MatchID,
			line.PlayerID,
			line.TeamID,
			occurredAt,
			event.SimulatedDate,
			int16(line.Minutes),
			int16(line.Points),
			int16(line.Rebounds),
			int16(line.Assists),
			int16(line.Steals),
			int16(line.Blocks),
			int16(line.Turnovers),
		); err != nil {
			return fmt.Errorf("insert player box score analytics: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit match analytics ingest: %w", err)
	}

	return nil
}

func (s *Store) IngestCityEconomyChange(ctx context.Context, event domain.CityEconomyChangeEvent) error {
	occurredAt := parseEventTime(event.OccurredAt)
	points := []struct {
		metric string
		value  float64
		delta  float64
	}{
		{metric: "fan_sentiment", value: event.FanSentiment, delta: event.FanSentimentDelta},
		{metric: "ticket_sales_index", value: event.TicketSalesIndex, delta: event.TicketSalesDelta},
		{metric: "local_economy_index", value: event.LocalEconomyIndex, delta: event.LocalEconomyDelta},
	}

	for _, point := range points {
		if _, err := s.pool.Exec(ctx, `
INSERT INTO analytics_city_metric_points (
	event_id, metric, game_id, occurred_at, simulated_date, value, delta,
	source_event_id, source_subject, reason, created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
ON CONFLICT (event_id, metric) DO NOTHING;
`,
			event.EventID,
			point.metric,
			event.GameID,
			occurredAt,
			event.SimulatedDate,
			point.value,
			point.delta,
			event.SourceEventID,
			event.SourceSubject,
			event.Reason,
		); err != nil {
			return fmt.Errorf("insert city metric analytics: %w", err)
		}
	}

	return nil
}

func (s *Store) IngestCityLandUpdated(ctx context.Context, event domain.CityLandUpdatedEvent) error {
	_, err := s.pool.Exec(ctx, `
INSERT INTO analytics_land_value_points (
	event_id, game_id, occurred_at, simulated_date, zone_id, value, delta,
	source_event_id, reason, created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
ON CONFLICT (event_id) DO NOTHING;
`,
		event.EventID,
		event.GameID,
		parseEventTime(event.OccurredAt),
		event.SimulatedDate,
		event.ZoneID,
		event.NewLandValue,
		event.LandValueDelta,
		event.SourceEventID,
		event.Reason,
	)
	if err != nil {
		return fmt.Errorf("insert land value analytics: %w", err)
	}

	return nil
}

func (s *Store) IngestAgentStateChanged(ctx context.Context, event domain.AgentStateChangedEvent) error {
	for metric, value := range event.State {
		if _, err := s.pool.Exec(ctx, `
INSERT INTO analytics_agent_state_points (
	event_id, metric, game_id, agent_id, occurred_at, simulated_date, value,
	mood, source_event_id, source_subject, created_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
ON CONFLICT (event_id, metric) DO NOTHING;
`,
			event.EventID,
			metric,
			event.GameID,
			event.AgentID,
			parseEventTime(event.OccurredAt),
			event.SimulatedDate,
			value,
			event.Mood,
			event.SourceEventID,
			event.SourceSubject,
		); err != nil {
			return fmt.Errorf("insert agent state analytics: %w", err)
		}
	}

	return nil
}

func parseEventTime(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Now().UTC()
	}

	return parsed.UTC()
}
