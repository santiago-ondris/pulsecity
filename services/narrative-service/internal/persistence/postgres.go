package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulsecity/services/narrative-service/internal/domain"
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

func (s *Store) SetOwnerIntroEventIfEmpty(ctx context.Context, gameID string, event domain.NarrativeEvent) (bool, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return false, fmt.Errorf("marshal owner intro event: %w", err)
	}

	const query = `
UPDATE games
SET owner_intro_event = $2::jsonb,
	updated_at = $3
WHERE game_id = $1
	AND owner_intro_event IS NULL;
`

	commandTag, err := s.pool.Exec(ctx, query, gameID, payload, time.Now().UTC())
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}
