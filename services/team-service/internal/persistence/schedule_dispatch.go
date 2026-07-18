package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pulsecity/services/team-service/internal/domain"
)

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
	playerStates, err := loadPlayerMatchStates(ctx, tx, scheduled.GameID)
	if err != nil {
		return domain.MatchScheduledEvent{}, false, err
	}
	record, err := loadSeasonRecord(ctx, tx, scheduled.GameID)
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
	event := domain.BuildPreparedMatchScheduledEvent(day, match, roster, domain.MatchPreparation{
		PlayerStates: playerStates,
		Record:       record,
	})

	if err := tx.Commit(ctx); err != nil {
		return domain.MatchScheduledEvent{}, false, fmt.Errorf("commit dispatch scheduled match: %w", err)
	}

	return event, true, nil
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

func loadPlayerMatchStates(ctx context.Context, q queryer, gameID string) (map[string]domain.PlayerMatchState, error) {
	rows, err := q.Query(ctx, `
WITH recent_box_scores AS (
	SELECT player_id, minutes
	FROM team_player_box_scores
	WHERE game_id = $1
	ORDER BY created_at DESC
	LIMIT 45
),
recent_minutes AS (
	SELECT player_id, COALESCE(SUM(minutes), 0) AS recent_minutes
	FROM recent_box_scores
	GROUP BY player_id
)
SELECT
	rp.player_id,
	COALESCE(rm.recent_minutes, 0),
	COALESCE(ms.emotional_state, 0)
FROM team_roster_players rp
LEFT JOIN recent_minutes rm ON rm.player_id = rp.player_id
LEFT JOIN team_player_match_states ms ON ms.game_id = rp.game_id AND ms.player_id = rp.player_id
WHERE rp.game_id = $1;
`, gameID)
	if err != nil {
		return nil, fmt.Errorf("load player match states %s: %w", gameID, err)
	}
	defer rows.Close()

	states := make(map[string]domain.PlayerMatchState)
	for rows.Next() {
		var playerID string
		var recentMinutes int64
		var emotionalState int16
		if err := rows.Scan(&playerID, &recentMinutes, &emotionalState); err != nil {
			return nil, fmt.Errorf("scan player match state %s: %w", gameID, err)
		}
		states[playerID] = domain.PlayerMatchState{
			RecentMinutes:  uint16(recentMinutes),
			EmotionalState: int8(emotionalState),
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate player match states %s: %w", gameID, err)
	}

	return states, nil
}
