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
CREATE TABLE IF NOT EXISTS guest_sessions (
	guest_token TEXT PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS games (
	game_id TEXT PRIMARY KEY,
	guest_token TEXT NOT NULL DEFAULT '',
	city_name TEXT NOT NULL DEFAULT '',
	franchise_name TEXT NOT NULL DEFAULT '',
	abbreviation TEXT NOT NULL DEFAULT '',
	primary_color TEXT NOT NULL DEFAULT '',
	secondary_color TEXT NOT NULL DEFAULT '',
	accent_color TEXT NOT NULL DEFAULT '',
	initial_scenario TEXT NOT NULL DEFAULT 'expansion',
	city_management_mode TEXT NOT NULL DEFAULT 'owner_influence',
	status TEXT NOT NULL,
	owner_intro_event JSONB,
	owner_intro_response JSONB,
	current_snapshot JSONB,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_games_guest_token_updated_at ON games (guest_token, updated_at DESC);

ALTER TABLE games ADD COLUMN IF NOT EXISTS guest_token TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS franchise_name TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS abbreviation TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS primary_color TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS secondary_color TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS accent_color TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS initial_scenario TEXT NOT NULL DEFAULT 'expansion';
ALTER TABLE games ADD COLUMN IF NOT EXISTS city_management_mode TEXT NOT NULL DEFAULT 'owner_influence';
ALTER TABLE games ADD COLUMN IF NOT EXISTS owner_intro_event JSONB;
ALTER TABLE games ADD COLUMN IF NOT EXISTS owner_intro_response JSONB;
`

	_, err := s.pool.Exec(ctx, query)
	return err
}

func (s *Store) CreateGuestSession(ctx context.Context, token string) error {
	const query = `
INSERT INTO guest_sessions (guest_token, created_at, last_seen_at)
VALUES ($1, $2, $2)
ON CONFLICT (guest_token) DO UPDATE
SET last_seen_at = EXCLUDED.last_seen_at;
`

	now := time.Now().UTC()
	_, err := s.pool.Exec(ctx, query, token, now)
	return err
}

func (s *Store) TouchGuestSession(ctx context.Context, token string) (bool, error) {
	const query = `
UPDATE guest_sessions
SET last_seen_at = $2
WHERE guest_token = $1;
`

	commandTag, err := s.pool.Exec(ctx, query, token, time.Now().UTC())
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func (s *Store) CreateGame(ctx context.Context, setup domain.GameSetup) error {
	const query = `
INSERT INTO games (
	game_id,
	guest_token,
	city_name,
	franchise_name,
	abbreviation,
	primary_color,
	secondary_color,
	accent_color,
	initial_scenario,
	city_management_mode,
	owner_intro_event,
	status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NULL, $11)
ON CONFLICT (game_id) DO NOTHING;
`

	_, err := s.pool.Exec(
		ctx,
		query,
		setup.GameID,
		setup.GuestToken,
		setup.CityName,
		setup.FranchiseName,
		setup.Abbreviation,
		setup.PrimaryColor,
		setup.SecondaryColor,
		setup.AccentColor,
		setup.InitialScenario,
		setup.CityManagementMode,
		setup.Status,
	)
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

func (s *Store) GetGame(ctx context.Context, gameID string) (domain.GameSetup, bool, error) {
	const query = `
SELECT
	game_id,
	guest_token,
	city_name,
	franchise_name,
	abbreviation,
	primary_color,
	secondary_color,
	accent_color,
	initial_scenario,
	city_management_mode,
	owner_intro_event,
	owner_intro_response,
	status,
	created_at,
	updated_at
FROM games
WHERE game_id = $1;
`

	var game domain.GameSetup
	var createdAt time.Time
	var updatedAt time.Time
	var ownerIntroRaw []byte
	var ownerIntroResponseRaw []byte
	if err := s.pool.QueryRow(ctx, query, gameID).Scan(
		&game.GameID,
		&game.GuestToken,
		&game.CityName,
		&game.FranchiseName,
		&game.Abbreviation,
		&game.PrimaryColor,
		&game.SecondaryColor,
		&game.AccentColor,
		&game.InitialScenario,
		&game.CityManagementMode,
		&ownerIntroRaw,
		&ownerIntroResponseRaw,
		&game.Status,
		&createdAt,
		&updatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.GameSetup{}, false, nil
		}
		return domain.GameSetup{}, false, err
	}

	game.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	game.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)
	if len(ownerIntroRaw) > 0 {
		var event domain.NarrativeEvent
		if err := json.Unmarshal(ownerIntroRaw, &event); err != nil {
			return domain.GameSetup{}, false, fmt.Errorf("unmarshal owner intro event: %w", err)
		}
		game.OwnerIntroEvent = &event
	}
	if len(ownerIntroResponseRaw) > 0 {
		var choice domain.NarrativeChoice
		if err := json.Unmarshal(ownerIntroResponseRaw, &choice); err != nil {
			return domain.GameSetup{}, false, fmt.Errorf("unmarshal owner intro response: %w", err)
		}
		game.OwnerIntroResponse = &choice
	}

	return game, true, nil
}

func (s *Store) ListGamesByGuest(ctx context.Context, guestToken string) ([]domain.GameSummary, error) {
	const query = `
SELECT
	game_id,
	city_name,
	franchise_name,
	initial_scenario,
	city_management_mode,
	status,
	updated_at
FROM games
WHERE guest_token = $1
ORDER BY updated_at DESC;
`

	rows, err := s.pool.Query(ctx, query, guestToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summaries := make([]domain.GameSummary, 0)
	for rows.Next() {
		var summary domain.GameSummary
		var updatedAt time.Time
		if err := rows.Scan(
			&summary.GameID,
			&summary.CityName,
			&summary.FranchiseName,
			&summary.InitialScenario,
			&summary.CityManagementMode,
			&summary.Status,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		summary.UpdatedAt = updatedAt.UTC().Format(time.RFC3339)
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return summaries, nil
}

func (s *Store) SetOwnerIntroEvent(ctx context.Context, gameID string, event domain.NarrativeEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal owner intro event: %w", err)
	}

	const query = `
UPDATE games
SET owner_intro_event = $2::jsonb,
	updated_at = $3
WHERE game_id = $1;
`

	_, err = s.pool.Exec(ctx, query, gameID, payload, time.Now().UTC())
	return err
}

func (s *Store) SetOwnerIntroResponse(ctx context.Context, gameID string, choice domain.NarrativeChoice) error {
	payload, err := json.Marshal(choice)
	if err != nil {
		return fmt.Errorf("marshal owner intro response: %w", err)
	}

	const query = `
UPDATE games
SET owner_intro_response = $2::jsonb,
	status = 'owner_intro_answered',
	updated_at = $3
WHERE game_id = $1;
`

	_, err = s.pool.Exec(ctx, query, gameID, payload, time.Now().UTC())
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
