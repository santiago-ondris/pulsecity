package persistence

import (
	"context"
	"errors"
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

CREATE TABLE IF NOT EXISTS team_records (
	game_id TEXT PRIMARY KEY,
	wins SMALLINT NOT NULL DEFAULT 0,
	losses SMALLINT NOT NULL DEFAULT 0,
	points_for INTEGER NOT NULL DEFAULT 0,
	points_against INTEGER NOT NULL DEFAULT 0,
	last_match_id TEXT,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
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

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit save initial season: %w", err)
	}

	return nil
}

func (s *Store) DispatchScheduledMatchForDate(
	ctx context.Context,
	day domain.DayAdvancedEvent,
) (domain.MatchScheduledEvent, bool, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.MatchScheduledEvent{}, false, fmt.Errorf("begin dispatch scheduled match: %w", err)
	}
	defer tx.Rollback(ctx)

	var scheduled storedScheduleMatch
	if err := tx.QueryRow(ctx, `
UPDATE team_schedule
SET status = 'scheduled_dispatched', updated_at = NOW()
WHERE match_id = (
	SELECT match_id
	FROM team_schedule
	WHERE game_id = $1
		AND simulated_date = $2
		AND status = 'scheduled'
	ORDER BY simulated_date, match_id
	LIMIT 1
)
RETURNING match_id, game_id, simulated_date, home_team_id, away_team_id, opponent_team_id, home_game, seed, status;
`, day.GameID, day.SimulatedDate).Scan(
		&scheduled.MatchID,
		&scheduled.GameID,
		&scheduled.SimulatedDate,
		&scheduled.HomeTeamID,
		&scheduled.AwayTeamID,
		&scheduled.OpponentTeamID,
		&scheduled.HomeGame,
		&scheduled.Seed,
		&scheduled.Status,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.MatchScheduledEvent{}, false, nil
		}
		return domain.MatchScheduledEvent{}, false, fmt.Errorf("claim scheduled match: %w", err)
	}

	franchise, err := loadFranchise(ctx, tx, scheduled.GameID)
	if err != nil {
		return domain.MatchScheduledEvent{}, false, err
	}
	opponent, err := loadOpponent(ctx, tx, scheduled.GameID, scheduled.OpponentTeamID)
	if err != nil {
		return domain.MatchScheduledEvent{}, false, err
	}
	roster, err := loadRoster(ctx, tx, scheduled.GameID)
	if err != nil {
		return domain.MatchScheduledEvent{}, false, err
	}

	homeTeam := franchise
	awayTeam := opponent
	if scheduled.HomeTeamID != domain.OwnTeamID {
		homeTeam = opponent
		awayTeam = franchise
	}

	match := domain.ScheduledMatch{
		MatchID:       scheduled.MatchID,
		GameID:        scheduled.GameID,
		SimulatedDate: scheduled.SimulatedDate,
		HomeTeam:      homeTeam,
		AwayTeam:      awayTeam,
		OpponentTeam:  opponent,
		HomeGame:      scheduled.HomeGame,
		Seed:          uint64(scheduled.Seed),
		Status:        scheduled.Status,
	}
	event := domain.BuildMatchScheduledEvent(day, match, roster)

	if err := tx.Commit(ctx); err != nil {
		return domain.MatchScheduledEvent{}, false, fmt.Errorf("commit dispatch scheduled match: %w", err)
	}

	return event, true, nil
}

func (s *Store) ApplyMatchFinished(
	ctx context.Context,
	event domain.MatchFinishedEvent,
) (domain.SeasonRecord, bool, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.SeasonRecord{}, false, fmt.Errorf("begin apply match finished: %w", err)
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx, `
UPDATE team_schedule
SET
	status = 'final',
	home_score = $3,
	away_score = $4,
	winner_team_id = $5,
	played_at = NOW(),
	updated_at = NOW()
WHERE game_id = $1
	AND match_id = $2
	AND status <> 'final';
`, event.GameID, event.MatchID, int16(event.HomeScore), int16(event.AwayScore), event.WinnerTeamID)
	if err != nil {
		return domain.SeasonRecord{}, false, fmt.Errorf("update finished match %s: %w", event.MatchID, err)
	}
	if commandTag.RowsAffected() == 0 {
		record, err := loadSeasonRecord(ctx, tx, event.GameID)
		if err != nil {
			return domain.SeasonRecord{}, false, err
		}
		return record, false, nil
	}

	for _, line := range event.BoxScore {
		if _, err := tx.Exec(ctx, `
INSERT INTO team_player_box_scores (
	match_id, player_id, game_id, team_id, minutes, points, rebounds, assists, steals, blocks, turnovers, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
ON CONFLICT (match_id, player_id) DO UPDATE SET
	team_id = EXCLUDED.team_id,
	minutes = EXCLUDED.minutes,
	points = EXCLUDED.points,
	rebounds = EXCLUDED.rebounds,
	assists = EXCLUDED.assists,
	steals = EXCLUDED.steals,
	blocks = EXCLUDED.blocks,
	turnovers = EXCLUDED.turnovers,
	updated_at = NOW();
`, event.MatchID, line.PlayerID, event.GameID, line.TeamID, int16(line.Minutes), int16(line.Points), int16(line.Rebounds), int16(line.Assists), int16(line.Steals), int16(line.Blocks), int16(line.Turnovers)); err != nil {
			return domain.SeasonRecord{}, false, fmt.Errorf("save box score %s/%s: %w", event.MatchID, line.PlayerID, err)
		}
	}

	record, err := recalculateSeasonRecord(ctx, tx, event.GameID)
	if err != nil {
		return domain.SeasonRecord{}, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.SeasonRecord{}, false, fmt.Errorf("commit apply match finished: %w", err)
	}

	return record, true, nil
}

type storedScheduleMatch struct {
	MatchID        string
	GameID         string
	SimulatedDate  string
	HomeTeamID     string
	AwayTeamID     string
	OpponentTeamID string
	HomeGame       bool
	Seed           int64
	Status         string
}

type queryer interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

func loadFranchise(ctx context.Context, q queryer, gameID string) (domain.MatchTeam, error) {
	var team domain.MatchTeam
	var rating, offense, defense, pace, homeCourt int16
	if err := q.QueryRow(ctx, `
SELECT team_id, name, abbreviation, rating, offense_rating, defense_rating, pace, home_court_advantage
FROM team_franchises
WHERE game_id = $1;
`, gameID).Scan(
		&team.TeamID,
		&team.Name,
		&team.Abbreviation,
		&rating,
		&offense,
		&defense,
		&pace,
		&homeCourt,
	); err != nil {
		return domain.MatchTeam{}, fmt.Errorf("load franchise %s: %w", gameID, err)
	}

	team.Rating = uint8(rating)
	team.OffenseRating = uint8(offense)
	team.DefenseRating = uint8(defense)
	team.Pace = uint8(pace)
	team.HomeCourtAdvantage = int8(homeCourt)
	return team, nil
}

func loadOpponent(ctx context.Context, q queryer, gameID, teamID string) (domain.MatchTeam, error) {
	var team domain.MatchTeam
	var rating, offense, defense, pace, homeCourt int16
	if err := q.QueryRow(ctx, `
SELECT team_id, name, abbreviation, rating, offense_rating, defense_rating, pace, home_court_advantage
FROM team_opponents
WHERE game_id = $1 AND team_id = $2;
`, gameID, teamID).Scan(
		&team.TeamID,
		&team.Name,
		&team.Abbreviation,
		&rating,
		&offense,
		&defense,
		&pace,
		&homeCourt,
	); err != nil {
		return domain.MatchTeam{}, fmt.Errorf("load opponent %s/%s: %w", gameID, teamID, err)
	}

	team.Rating = uint8(rating)
	team.OffenseRating = uint8(offense)
	team.DefenseRating = uint8(defense)
	team.Pace = uint8(pace)
	team.HomeCourtAdvantage = int8(homeCourt)
	return team, nil
}

func loadRoster(ctx context.Context, q queryer, gameID string) ([]domain.RosterPlayer, error) {
	rows, err := q.Query(ctx, `
SELECT player_id, game_id, full_name, position, overall_rating, roster_status, contract_years, salary, sort_order
FROM team_roster_players
WHERE game_id = $1
ORDER BY sort_order;
`, gameID)
	if err != nil {
		return nil, fmt.Errorf("load roster %s: %w", gameID, err)
	}
	defer rows.Close()

	var roster []domain.RosterPlayer
	for rows.Next() {
		var player domain.RosterPlayer
		var rating, years int16
		if err := rows.Scan(
			&player.PlayerID,
			&player.GameID,
			&player.FullName,
			&player.Position,
			&rating,
			&player.RosterStatus,
			&years,
			&player.Salary,
			&player.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("scan roster %s: %w", gameID, err)
		}
		player.OverallRating = uint8(rating)
		player.ContractYears = uint8(years)
		roster = append(roster, player)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate roster %s: %w", gameID, err)
	}

	return roster, nil
}

func recalculateSeasonRecord(ctx context.Context, q queryer, gameID string) (domain.SeasonRecord, error) {
	rows, err := q.Query(ctx, `
SELECT match_id, simulated_date, home_team_id, away_team_id, home_score, away_score, winner_team_id
FROM team_schedule
WHERE game_id = $1
	AND status = 'final'
ORDER BY simulated_date, match_id;
`, gameID)
	if err != nil {
		return domain.SeasonRecord{}, fmt.Errorf("load final matches %s: %w", gameID, err)
	}
	defer rows.Close()

	record := domain.SeasonRecord{GameID: gameID}
	for rows.Next() {
		var result domain.SeasonMatchSummary
		var homeScore, awayScore int16
		if err := rows.Scan(
			&result.MatchID,
			&result.SimulatedDate,
			&result.HomeTeamID,
			&result.AwayTeamID,
			&homeScore,
			&awayScore,
			&result.WinnerTeamID,
		); err != nil {
			return domain.SeasonRecord{}, fmt.Errorf("scan final match %s: %w", gameID, err)
		}
		result.HomeScore = uint16(homeScore)
		result.AwayScore = uint16(awayScore)
		if result.WinnerTeamID == domain.OwnTeamID {
			record.Wins++
		} else {
			record.Losses++
		}
		if result.HomeTeamID == domain.OwnTeamID {
			record.PointsFor += result.HomeScore
			record.PointsAgainst += result.AwayScore
		} else {
			record.PointsFor += result.AwayScore
			record.PointsAgainst += result.HomeScore
		}
		last := result
		record.LastResult = &last
	}
	if err := rows.Err(); err != nil {
		return domain.SeasonRecord{}, fmt.Errorf("iterate final matches %s: %w", gameID, err)
	}

	var lastMatchID any
	if record.LastResult != nil {
		lastMatchID = record.LastResult.MatchID
	}
	if _, err := q.Exec(ctx, `
INSERT INTO team_records (game_id, wins, losses, points_for, points_against, last_match_id, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
ON CONFLICT (game_id) DO UPDATE SET
	wins = EXCLUDED.wins,
	losses = EXCLUDED.losses,
	points_for = EXCLUDED.points_for,
	points_against = EXCLUDED.points_against,
	last_match_id = EXCLUDED.last_match_id,
	updated_at = NOW();
`, gameID, int16(record.Wins), int16(record.Losses), int(record.PointsFor), int(record.PointsAgainst), lastMatchID); err != nil {
		return domain.SeasonRecord{}, fmt.Errorf("upsert team record %s: %w", gameID, err)
	}

	return record, nil
}

func loadSeasonRecord(ctx context.Context, q queryer, gameID string) (domain.SeasonRecord, error) {
	var record domain.SeasonRecord
	var wins, losses int16
	var pointsFor, pointsAgainst int
	if err := q.QueryRow(ctx, `
SELECT game_id, wins, losses, points_for, points_against
FROM team_records
WHERE game_id = $1;
`, gameID).Scan(
		&record.GameID,
		&wins,
		&losses,
		&pointsFor,
		&pointsAgainst,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SeasonRecord{GameID: gameID}, nil
		}
		return domain.SeasonRecord{}, fmt.Errorf("load team record %s: %w", gameID, err)
	}

	record.Wins = uint16(wins)
	record.Losses = uint16(losses)
	record.PointsFor = uint16(pointsFor)
	record.PointsAgainst = uint16(pointsAgainst)
	return record, nil
}
