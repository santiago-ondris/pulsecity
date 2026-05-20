package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
CREATE TABLE IF NOT EXISTS users (
	user_id TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	display_name TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_sessions (
	session_token TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS guest_sessions (
	guest_token TEXT PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS games (
	game_id TEXT PRIMARY KEY,
	guest_token TEXT NOT NULL DEFAULT '',
	user_id TEXT,
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

ALTER TABLE games ADD COLUMN IF NOT EXISTS guest_token TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS user_id TEXT;
ALTER TABLE games ADD COLUMN IF NOT EXISTS franchise_name TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS abbreviation TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS primary_color TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS secondary_color TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS accent_color TEXT NOT NULL DEFAULT '';
ALTER TABLE games ADD COLUMN IF NOT EXISTS initial_scenario TEXT NOT NULL DEFAULT 'expansion';
ALTER TABLE games ADD COLUMN IF NOT EXISTS city_management_mode TEXT NOT NULL DEFAULT 'owner_influence';
ALTER TABLE games ADD COLUMN IF NOT EXISTS owner_intro_event JSONB;
ALTER TABLE games ADD COLUMN IF NOT EXISTS owner_intro_response JSONB;

CREATE INDEX IF NOT EXISTS idx_games_guest_token_updated_at ON games (guest_token, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_games_user_id_updated_at ON games (user_id, updated_at DESC);

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'games_exactly_one_owner'
	) THEN
		ALTER TABLE games
		ADD CONSTRAINT games_exactly_one_owner
		CHECK (
			((CASE WHEN guest_token <> '' THEN 1 ELSE 0 END) +
			(CASE WHEN user_id IS NOT NULL AND user_id <> '' THEN 1 ELSE 0 END)) = 1
		) NOT VALID;
	END IF;
END $$;
`

	_, err := s.pool.Exec(ctx, query)
	return err
}

func (s *Store) CreateUser(ctx context.Context, email, displayName, passwordHash string) (domain.User, error) {
	const query = `
INSERT INTO users (user_id, email, display_name, password_hash, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $5)
RETURNING user_id, email, display_name, created_at;
`

	now := time.Now().UTC()
	var user domain.User
	if err := s.pool.QueryRow(
		ctx,
		query,
		uuid.NewString(),
		normalizeEmail(email),
		strings.TrimSpace(displayName),
		passwordHash,
		now,
	).Scan(&user.UserID, &user.Email, &user.DisplayName, &now); err != nil {
		return domain.User{}, err
	}

	user.CreatedAt = now.UTC().Format(time.RFC3339)
	return user, nil
}

func (s *Store) GetUserCredentialsByEmail(ctx context.Context, email string) (domain.User, string, bool, error) {
	const query = `
SELECT user_id, email, display_name, password_hash, created_at
FROM users
WHERE email = $1;
`

	var user domain.User
	var passwordHash string
	var createdAt time.Time
	if err := s.pool.QueryRow(ctx, query, normalizeEmail(email)).Scan(
		&user.UserID,
		&user.Email,
		&user.DisplayName,
		&passwordHash,
		&createdAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, "", false, nil
		}
		return domain.User{}, "", false, err
	}

	user.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	return user, passwordHash, true, nil
}

func (s *Store) CreateUserSession(ctx context.Context, user domain.User) (domain.UserSession, error) {
	const query = `
INSERT INTO user_sessions (session_token, user_id, created_at, last_seen_at)
VALUES ($1, $2, $3, $3);
`

	now := time.Now().UTC()
	sessionToken := "session_" + uuid.NewString()
	if _, err := s.pool.Exec(ctx, query, sessionToken, user.UserID, now); err != nil {
		return domain.UserSession{}, err
	}

	return domain.UserSession{
		SessionToken: sessionToken,
		User:         user,
		CreatedAt:    now.Format(time.RFC3339),
		LastSeenAt:   now.Format(time.RFC3339),
	}, nil
}

func (s *Store) GetUserSession(ctx context.Context, sessionToken string) (domain.UserSession, bool, error) {
	const query = `
SELECT
	s.session_token,
	s.created_at,
	s.last_seen_at,
	u.user_id,
	u.email,
	u.display_name,
	u.created_at
FROM user_sessions s
JOIN users u ON u.user_id = s.user_id
WHERE s.session_token = $1;
`

	var session domain.UserSession
	var sessionCreatedAt time.Time
	var lastSeenAt time.Time
	var userCreatedAt time.Time
	if err := s.pool.QueryRow(ctx, query, strings.TrimSpace(sessionToken)).Scan(
		&session.SessionToken,
		&sessionCreatedAt,
		&lastSeenAt,
		&session.User.UserID,
		&session.User.Email,
		&session.User.DisplayName,
		&userCreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserSession{}, false, nil
		}
		return domain.UserSession{}, false, err
	}

	session.CreatedAt = sessionCreatedAt.UTC().Format(time.RFC3339)
	session.LastSeenAt = lastSeenAt.UTC().Format(time.RFC3339)
	session.User.CreatedAt = userCreatedAt.UTC().Format(time.RFC3339)
	return session, true, nil
}

func (s *Store) TouchUserSession(ctx context.Context, sessionToken string) (bool, error) {
	const query = `
UPDATE user_sessions
SET last_seen_at = $2
WHERE session_token = $1;
`

	commandTag, err := s.pool.Exec(ctx, query, strings.TrimSpace(sessionToken), time.Now().UTC())
	if err != nil {
		return false, err
	}

	return commandTag.RowsAffected() > 0, nil
}

func (s *Store) MigrateGuestGamesToUser(ctx context.Context, guestToken, userID string) (int, error) {
	guestToken = strings.TrimSpace(guestToken)
	userID = strings.TrimSpace(userID)
	if guestToken == "" || userID == "" {
		return 0, fmt.Errorf("guest token and user id are required")
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const query = `
UPDATE games
SET guest_token = '',
	user_id = $2,
	updated_at = $3
WHERE guest_token = $1
	AND (user_id IS NULL OR user_id = '');
`

	commandTag, err := tx.Exec(ctx, query, guestToken, userID, time.Now().UTC())
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return int(commandTag.RowsAffected()), nil
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
	if !hasExclusiveOwner(setup.GuestToken, setup.UserID) {
		return fmt.Errorf("game must have exactly one owner")
	}

	const query = `
INSERT INTO games (
	game_id,
	guest_token,
	user_id,
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NULL, $12)
ON CONFLICT (game_id) DO NOTHING;
`

	_, err := s.pool.Exec(
		ctx,
		query,
		setup.GameID,
		setup.GuestToken,
		nullableText(setup.UserID),
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
UPDATE games
SET status = $2,
	current_snapshot = $3::jsonb,
	updated_at = $4
WHERE game_id = $1;
`

	now := time.Now().UTC()
	commandTag, err := s.pool.Exec(ctx, query, state.GameID, statusFromStage(state.Stage), payload, now)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("game %s not found while persisting snapshot", state.GameID)
	}

	return err
}

func (s *Store) GetGame(ctx context.Context, gameID string) (domain.GameSetup, bool, error) {
	const query = `
SELECT
	game_id,
	guest_token,
	COALESCE(user_id, '') AS user_id,
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
		&game.UserID,
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
	game.OwnerKind = ownerKindForGame(game.GuestToken, game.UserID)
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
	CASE WHEN user_id IS NOT NULL AND user_id <> '' THEN 'user' ELSE 'guest' END AS owner_kind,
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
			&summary.OwnerKind,
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

func (s *Store) ListGamesByUser(ctx context.Context, userID string) ([]domain.GameSummary, error) {
	const query = `
SELECT
	game_id,
	city_name,
	franchise_name,
	'user' AS owner_kind,
	initial_scenario,
	city_management_mode,
	status,
	updated_at
FROM games
WHERE user_id = $1
ORDER BY updated_at DESC;
`

	rows, err := s.pool.Query(ctx, query, userID)
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
			&summary.OwnerKind,
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

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func ownerKindForGame(guestToken, userID string) string {
	if strings.TrimSpace(userID) != "" {
		return domain.OwnerKindUser
	}
	if strings.TrimSpace(guestToken) != "" {
		return domain.OwnerKindGuest
	}

	return ""
}

func hasExclusiveOwner(guestToken, userID string) bool {
	hasGuest := strings.TrimSpace(guestToken) != ""
	hasUser := strings.TrimSpace(userID) != ""
	return (hasGuest || hasUser) && !(hasGuest && hasUser)
}

func nullableText(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return trimmed
}
