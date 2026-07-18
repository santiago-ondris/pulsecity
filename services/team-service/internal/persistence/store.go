package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pulsecity/services/team-service/internal/domain"
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

CREATE TABLE IF NOT EXISTS team_player_match_states (
	game_id TEXT NOT NULL,
	player_id TEXT NOT NULL,
	emotional_state SMALLINT NOT NULL DEFAULT 0,
	source_simulated_date TEXT NOT NULL,
	source_event_id TEXT NOT NULL,
	source_subject TEXT NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	PRIMARY KEY (game_id, player_id)
);

CREATE TABLE IF NOT EXISTS team_injuries (
	injury_id TEXT PRIMARY KEY,
	game_id TEXT NOT NULL,
	player_id TEXT NOT NULL,
	severity TEXT NOT NULL,
	estimated_days_out SMALLINT NOT NULL,
	injured_on TEXT NOT NULL,
	expected_recovery_date TEXT NOT NULL,
	recovered_on TEXT,
	reason TEXT NOT NULL,
	source_match_id TEXT NOT NULL,
	workload_score SMALLINT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_team_injuries_game_recovery
ON team_injuries (game_id, expected_recovery_date)
WHERE recovered_on IS NULL;

ALTER TABLE team_injuries ADD COLUMN IF NOT EXISTS medical_decision TEXT;
ALTER TABLE team_injuries ADD COLUMN IF NOT EXISTS decision_id TEXT;
ALTER TABLE team_injuries ADD COLUMN IF NOT EXISTS forced_return_at TEXT;

CREATE TABLE IF NOT EXISTS team_records (
	game_id TEXT PRIMARY KEY,
	wins SMALLINT NOT NULL DEFAULT 0,
	losses SMALLINT NOT NULL DEFAULT 0,
	points_for INTEGER NOT NULL DEFAULT 0,
	points_against INTEGER NOT NULL DEFAULT 0,
	last_match_id TEXT,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_salary_cap (
	game_id TEXT PRIMARY KEY,
	simulated_date TEXT NOT NULL,
	cap_base INTEGER NOT NULL,
	luxury_tax_line INTEGER NOT NULL,
	committed_salary INTEGER NOT NULL,
	cap_space INTEGER NOT NULL,
	luxury_tax_space INTEGER NOT NULL,
	roster_count SMALLINT NOT NULL,
	status TEXT NOT NULL,
	near_luxury_tax BOOLEAN NOT NULL,
	projected_tax_payment INTEGER NOT NULL,
	source_event_id TEXT NOT NULL,
	source_subject TEXT NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_trades (
	proposal_id TEXT PRIMARY KEY,
	game_id TEXT NOT NULL,
	rival_team_id TEXT NOT NULL,
	offered_player_id TEXT NOT NULL,
	offered_player_name TEXT NOT NULL,
	offered_salary INTEGER NOT NULL,
	requested_position TEXT NOT NULL,
	incoming_salary INTEGER NOT NULL,
	cap_space_after INTEGER NOT NULL,
	status TEXT NOT NULL,
	rejection_reason TEXT,
	rejection_detail TEXT,
	incoming_player_id TEXT,
	incoming_player_name TEXT,
	incoming_rating SMALLINT,
	accepted_additional_asset TEXT,
	accepted_at TEXT,
	source_decision_id TEXT NOT NULL,
	source_event_id TEXT NOT NULL,
	simulated_date TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_team_trades_game_status
ON team_trades (game_id, status);

ALTER TABLE team_trades ADD COLUMN IF NOT EXISTS incoming_player_id TEXT;
ALTER TABLE team_trades ADD COLUMN IF NOT EXISTS incoming_player_name TEXT;
ALTER TABLE team_trades ADD COLUMN IF NOT EXISTS incoming_rating SMALLINT;
ALTER TABLE team_trades ADD COLUMN IF NOT EXISTS accepted_additional_asset TEXT;
ALTER TABLE team_trades ADD COLUMN IF NOT EXISTS accepted_at TEXT;
`

	_, err := s.pool.Exec(ctx, query)
	return err
}

func (s *Store) SaveInitialSeason(ctx context.Context, season domain.InitialSeason) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin save initial season: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
INSERT INTO team_franchises (
	game_id, team_id, name, abbreviation, rating, offense_rating, defense_rating, pace, home_court_advantage, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
ON CONFLICT (game_id) DO UPDATE SET
	team_id = EXCLUDED.team_id,
	name = EXCLUDED.name,
	abbreviation = EXCLUDED.abbreviation,
	rating = EXCLUDED.rating,
	offense_rating = EXCLUDED.offense_rating,
	defense_rating = EXCLUDED.defense_rating,
	pace = EXCLUDED.pace,
	home_court_advantage = EXCLUDED.home_court_advantage,
	updated_at = NOW();
`, season.GameID, season.Team.TeamID, season.Team.Name, season.Team.Abbreviation, int16(season.Team.Rating), int16(season.Team.OffenseRating), int16(season.Team.DefenseRating), int16(season.Team.Pace), int16(season.Team.HomeCourtAdvantage)); err != nil {
		return fmt.Errorf("save franchise: %w", err)
	}

	for _, player := range season.Roster {
		if _, err := tx.Exec(ctx, `
INSERT INTO team_roster_players (
	player_id, game_id, full_name, position, overall_rating, roster_status, contract_years, salary, sort_order, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
ON CONFLICT (player_id) DO UPDATE SET
	full_name = EXCLUDED.full_name,
	position = EXCLUDED.position,
	overall_rating = EXCLUDED.overall_rating,
	roster_status = EXCLUDED.roster_status,
	contract_years = EXCLUDED.contract_years,
	salary = EXCLUDED.salary,
	sort_order = EXCLUDED.sort_order,
	updated_at = NOW();
`, player.PlayerID, player.GameID, player.FullName, player.Position, int16(player.OverallRating), player.RosterStatus, int16(player.ContractYears), player.Salary, player.SortOrder); err != nil {
			return fmt.Errorf("save roster player %s: %w", player.PlayerID, err)
		}
	}

	for _, opponent := range season.Opponents {
		if _, err := tx.Exec(ctx, `
INSERT INTO team_opponents (
	game_id, team_id, name, abbreviation, rating, offense_rating, defense_rating, pace, home_court_advantage, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
ON CONFLICT (game_id, team_id) DO UPDATE SET
	name = EXCLUDED.name,
	abbreviation = EXCLUDED.abbreviation,
	rating = EXCLUDED.rating,
	offense_rating = EXCLUDED.offense_rating,
	defense_rating = EXCLUDED.defense_rating,
	pace = EXCLUDED.pace,
	home_court_advantage = EXCLUDED.home_court_advantage,
	updated_at = NOW();
`, season.GameID, opponent.TeamID, opponent.Name, opponent.Abbreviation, int16(opponent.Rating), int16(opponent.OffenseRating), int16(opponent.DefenseRating), int16(opponent.Pace), int16(opponent.HomeCourtAdvantage)); err != nil {
			return fmt.Errorf("save opponent %s: %w", opponent.TeamID, err)
		}
	}

	for _, match := range season.Schedule {
		if _, err := tx.Exec(ctx, `
INSERT INTO team_schedule (
	match_id, game_id, simulated_date, home_team_id, away_team_id, opponent_team_id, home_game, seed, status, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
ON CONFLICT (match_id) DO UPDATE SET
	simulated_date = EXCLUDED.simulated_date,
	home_team_id = EXCLUDED.home_team_id,
	away_team_id = EXCLUDED.away_team_id,
	opponent_team_id = EXCLUDED.opponent_team_id,
	home_game = EXCLUDED.home_game,
	seed = EXCLUDED.seed,
	status = CASE
		WHEN team_schedule.status IN ('scheduled_dispatched', 'final') THEN team_schedule.status
		ELSE EXCLUDED.status
	END,
	updated_at = NOW();
`, match.MatchID, match.GameID, match.SimulatedDate, match.HomeTeam.TeamID, match.AwayTeam.TeamID, match.OpponentTeam.TeamID, match.HomeGame, int64(match.Seed), match.Status); err != nil {
			return fmt.Errorf("save scheduled match %s: %w", match.MatchID, err)
		}
	}

	if _, err := tx.Exec(ctx, `
INSERT INTO team_records (game_id, wins, losses, points_for, points_against, updated_at)
VALUES ($1, 0, 0, 0, 0, NOW())
ON CONFLICT (game_id) DO NOTHING;
`, season.GameID); err != nil {
		return fmt.Errorf("initialize team record: %w", err)
	}

	salaryCap := domain.CalculateSalaryCap(
		season.GameID,
		season.Roster,
		domain.DefaultSeasonStartDate,
		"",
		"initial-season-"+season.GameID,
		domain.SubjectMapGenerationStarted,
	)
	if err := saveSalaryCap(ctx, tx, salaryCap); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit save initial season: %w", err)
	}

	return nil
}

type queryer interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}
