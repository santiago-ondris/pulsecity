package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulsecity/services/gateway/internal/domain"
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
CREATE TABLE IF NOT EXISTS games (
	game_id TEXT PRIMARY KEY,
	city_name TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	current_snapshot JSONB,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

	_, err := s.pool.Exec(ctx, query)
	return err
}

func (s *Store) CreateGame(ctx context.Context, gameID, cityName string) error {
	const query = `
INSERT INTO games (game_id, city_name, status)
VALUES ($1, $2, 'generation_started')
ON CONFLICT (game_id) DO NOTHING;
`

	_, err := s.pool.Exec(ctx, query, gameID, cityName)
	return err
}

func (s *Store) UpsertSnapshot(ctx context.Context, state domain.MapClientState) error {
	payload, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	const query = `
INSERT INTO games (game_id, city_name, status, current_snapshot, created_at, updated_at)
VALUES ($1, '', $2, $3::jsonb, $4, $4)
ON CONFLICT (game_id) DO UPDATE
SET status = EXCLUDED.status,
	current_snapshot = EXCLUDED.current_snapshot,
	updated_at = EXCLUDED.updated_at;
`

	now := time.Now().UTC()
	_, err = s.pool.Exec(ctx, query, state.GameID, statusFromStage(state.Stage), payload, now)
	return err
}

func (s *Store) GetSnapshot(ctx context.Context, gameID string) (domain.MapClientState, bool, error) {
	const query = `
SELECT current_snapshot
FROM games
WHERE game_id = $1;
`

	var raw []byte
	if err := s.pool.QueryRow(ctx, query, gameID).Scan(&raw); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.MapClientState{}, false, nil
		}
		return domain.MapClientState{}, false, err
	}

	if len(raw) == 0 {
		return domain.MapClientState{}, false, nil
	}

	var state domain.MapClientState
	if err := json.Unmarshal(raw, &state); err != nil {
		return domain.MapClientState{}, false, fmt.Errorf("unmarshal snapshot: %w", err)
	}

	return state, true, nil
}

func statusFromStage(stage string) string {
	switch stage {
	case "complete":
		return "map_generation_complete"
	default:
		return "map_" + stage
	}
}
