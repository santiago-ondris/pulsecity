package persistence

import (
	"context"
	"fmt"

	"github.com/pulsecity/services/team-service/internal/domain"
)

func (s *Store) ApplyRosterPatch(ctx context.Context, event domain.RosterPatchEnvelope) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin apply roster patch: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, player := range event.Patch.Players {
		emotionalState := domain.EmotionalStateScore(player.EmotionalState)
		if _, err := tx.Exec(ctx, `
INSERT INTO team_player_match_states (
	game_id, player_id, emotional_state, source_simulated_date, source_event_id, source_subject, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
ON CONFLICT (game_id, player_id) DO UPDATE SET
	emotional_state = EXCLUDED.emotional_state,
	source_simulated_date = EXCLUDED.source_simulated_date,
	source_event_id = EXCLUDED.source_event_id,
	source_subject = EXCLUDED.source_subject,
	updated_at = NOW()
WHERE team_player_match_states.source_simulated_date <= EXCLUDED.source_simulated_date;
`, event.GameID, player.PlayerID, int16(emotionalState), event.Patch.SimulatedDate, event.Patch.SourceEventID, event.Patch.SourceSubject); err != nil {
			return fmt.Errorf("upsert roster patch player %s: %w", player.PlayerID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit apply roster patch: %w", err)
	}

	return nil
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

func loadRosterPlayer(ctx context.Context, q queryer, gameID, playerID string) (domain.RosterPlayer, error) {
	var player domain.RosterPlayer
	var overallRating, contractYears int16
	if err := q.QueryRow(ctx, `
SELECT player_id, game_id, full_name, position, overall_rating, roster_status, contract_years, salary, sort_order
FROM team_roster_players
WHERE game_id = $1 AND player_id = $2;
`, gameID, playerID).Scan(
		&player.PlayerID,
		&player.GameID,
		&player.FullName,
		&player.Position,
		&overallRating,
		&player.RosterStatus,
		&contractYears,
		&player.Salary,
		&player.SortOrder,
	); err != nil {
		return domain.RosterPlayer{}, fmt.Errorf("load roster player %s/%s: %w", gameID, playerID, err)
	}

	player.OverallRating = uint8(overallRating)
	player.ContractYears = uint8(contractYears)
	return player, nil
}
