package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pulsecity/services/team-service/internal/domain"
)

func (s *Store) ApplyMatchFinished(
	ctx context.Context,
	event domain.MatchFinishedEvent,
) (domain.SeasonRecord, []domain.PlayerInjuredEvent, bool, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.SeasonRecord{}, nil, false, fmt.Errorf("begin apply match finished: %w", err)
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
		return domain.SeasonRecord{}, nil, false, fmt.Errorf("update finished match %s: %w", event.MatchID, err)
	}
	if commandTag.RowsAffected() == 0 {
		record, err := loadSeasonRecord(ctx, tx, event.GameID)
		if err != nil {
			return domain.SeasonRecord{}, nil, false, err
		}
		return record, nil, false, nil
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
			return domain.SeasonRecord{}, nil, false, fmt.Errorf("save box score %s/%s: %w", event.MatchID, line.PlayerID, err)
		}
	}

	injuries, err := createInjuriesForMatch(ctx, tx, event)
	if err != nil {
		return domain.SeasonRecord{}, nil, false, err
	}

	record, err := recalculateSeasonRecord(ctx, tx, event.GameID)
	if err != nil {
		return domain.SeasonRecord{}, nil, false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.SeasonRecord{}, nil, false, fmt.Errorf("commit apply match finished: %w", err)
	}

	return record, injuries, true, nil
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
