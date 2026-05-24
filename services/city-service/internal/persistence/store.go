package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulsecity/services/city-service/internal/domain"
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
`

	_, err := s.pool.Exec(ctx, query)
	return err
}

func (s *Store) ApplyMatchFinished(ctx context.Context, event domain.MatchFinishedEvent) (domain.CityReaction, bool, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.CityReaction{}, false, fmt.Errorf("begin city reaction: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
INSERT INTO city_processed_matches (game_id, match_id, source_event_id, processed_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (game_id, match_id) DO NOTHING;
`, event.GameID, event.MatchID, event.EventID)
	if err != nil {
		return domain.CityReaction{}, false, fmt.Errorf("mark city match processed: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.CityReaction{}, false, nil
	}

	current, err := loadMetrics(ctx, tx, event.GameID)
	if err != nil {
		return domain.CityReaction{}, false, err
	}

	reaction := domain.ApplyMatchFinished(current, event)
	if err := saveMetrics(ctx, tx, reaction.Metrics); err != nil {
		return domain.CityReaction{}, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.CityReaction{}, false, fmt.Errorf("commit city reaction: %w", err)
	}

	return reaction, true, nil
}

func loadMetrics(ctx context.Context, tx pgx.Tx, gameID string) (domain.CityMetrics, error) {
	var metrics domain.CityMetrics
	err := tx.QueryRow(ctx, `
SELECT game_id, fan_sentiment, ticket_sales_index, local_economy_index,
	stadium_district_land_value, win_streak, loss_streak, COALESCE(last_match_id, '')
FROM city_metrics
WHERE game_id = $1;
`, gameID).Scan(
		&metrics.GameID,
		&metrics.FanSentiment,
		&metrics.TicketSalesIndex,
		&metrics.LocalEconomyIndex,
		&metrics.StadiumDistrictLandValue,
		&metrics.WinStreak,
		&metrics.LossStreak,
		&metrics.LastMatchID,
	)
	if err == nil {
		return metrics, nil
	}
	if err == pgx.ErrNoRows {
		return domain.DefaultCityMetrics(gameID), nil
	}

	return domain.CityMetrics{}, fmt.Errorf("load city metrics: %w", err)
}

func saveMetrics(ctx context.Context, tx pgx.Tx, metrics domain.CityMetrics) error {
	_, err := tx.Exec(ctx, `
INSERT INTO city_metrics (
	game_id, fan_sentiment, ticket_sales_index, local_economy_index, stadium_district_land_value,
	win_streak, loss_streak, last_match_id, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
ON CONFLICT (game_id) DO UPDATE SET
	fan_sentiment = EXCLUDED.fan_sentiment,
	ticket_sales_index = EXCLUDED.ticket_sales_index,
	local_economy_index = EXCLUDED.local_economy_index,
	stadium_district_land_value = EXCLUDED.stadium_district_land_value,
	win_streak = EXCLUDED.win_streak,
	loss_streak = EXCLUDED.loss_streak,
	last_match_id = EXCLUDED.last_match_id,
	updated_at = NOW();
`,
		metrics.GameID,
		metrics.FanSentiment,
		metrics.TicketSalesIndex,
		metrics.LocalEconomyIndex,
		metrics.StadiumDistrictLandValue,
		metrics.WinStreak,
		metrics.LossStreak,
		metrics.LastMatchID,
	)
	if err != nil {
		return fmt.Errorf("save city metrics: %w", err)
	}

	return nil
}
