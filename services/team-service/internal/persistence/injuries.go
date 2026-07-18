package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pulsecity/services/team-service/internal/domain"
)

func (s *Store) ApplyMedicalDecision(
	ctx context.Context,
	event domain.GMDecisionRegisteredEvent,
) ([]domain.PlayerRecoveredEvent, error) {
	if event.Kind != "medical_decision" {
		return nil, nil
	}

	injuryID := event.Payload["injury_id"]
	playerID := event.Payload["player_id"]
	choiceID := event.Payload["choice_id"]
	if injuryID == "" || playerID == "" || choiceID == "" {
		return nil, fmt.Errorf("apply medical decision: missing injury_id, player_id or choice_id")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin apply medical decision: %w", err)
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx, `
UPDATE team_injuries
SET medical_decision = $4, decision_id = $5, updated_at = NOW()
WHERE game_id = $1
	AND injury_id = $2
	AND player_id = $3;
`, event.GameID, injuryID, playerID, choiceID, event.DecisionID)
	if err != nil {
		return nil, fmt.Errorf("record medical decision %s: %w", injuryID, err)
	}
	if commandTag.RowsAffected() == 0 {
		return nil, fmt.Errorf("record medical decision %s: injury not found", injuryID)
	}

	var recovered []domain.PlayerRecoveredEvent
	if choiceID == "force_return" {
		commandTag, err := tx.Exec(ctx, `
UPDATE team_injuries
SET recovered_on = $4, forced_return_at = $4, updated_at = NOW()
WHERE game_id = $1
	AND injury_id = $2
	AND player_id = $3
	AND recovered_on IS NULL;
`, event.GameID, injuryID, playerID, event.SimulatedDate)
		if err != nil {
			return nil, fmt.Errorf("force player return %s: %w", injuryID, err)
		}
		if commandTag.RowsAffected() > 0 {
			if _, err := tx.Exec(ctx, `
UPDATE team_roster_players
SET roster_status = 'active', updated_at = NOW()
WHERE game_id = $1
	AND player_id = $2
	AND roster_status = 'injured';
`, event.GameID, playerID); err != nil {
				return nil, fmt.Errorf("activate forced return player %s: %w", playerID, err)
			}
			recovered = append(recovered, domain.PlayerRecoveredEvent{
				EventMeta: domain.EventMeta{
					EventID:       "player-recovered-" + injuryID,
					GameID:        event.GameID,
					OccurredAt:    event.OccurredAt,
					SchemaVersion: 1,
				},
				InjuryID:    injuryID,
				PlayerID:    playerID,
				RecoveredOn: event.SimulatedDate,
			})
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit apply medical decision: %w", err)
	}

	return recovered, nil
}

func (s *Store) RecoverPlayersForDate(
	ctx context.Context,
	day domain.DayAdvancedEvent,
) ([]domain.PlayerRecoveredEvent, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin recover players: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
UPDATE team_injuries
SET recovered_on = $2, updated_at = NOW()
WHERE game_id = $1
	AND recovered_on IS NULL
	AND expected_recovery_date <= $2
RETURNING injury_id, player_id;
`, day.GameID, day.SimulatedDate)
	if err != nil {
		return nil, fmt.Errorf("recover players %s/%s: %w", day.GameID, day.SimulatedDate, err)
	}
	defer rows.Close()

	var recovered []domain.PlayerRecoveredEvent
	for rows.Next() {
		event := domain.PlayerRecoveredEvent{
			EventMeta: domain.EventMeta{
				GameID:        day.GameID,
				OccurredAt:    day.OccurredAt,
				SchemaVersion: 1,
			},
			RecoveredOn: day.SimulatedDate,
		}
		if err := rows.Scan(&event.InjuryID, &event.PlayerID); err != nil {
			return nil, fmt.Errorf("scan recovered player %s: %w", day.GameID, err)
		}
		event.EventID = fmt.Sprintf("player-recovered-%s", event.InjuryID)
		recovered = append(recovered, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recovered players %s: %w", day.GameID, err)
	}

	for _, event := range recovered {
		if _, err := tx.Exec(ctx, `
UPDATE team_roster_players
SET roster_status = 'active', updated_at = NOW()
WHERE game_id = $1
	AND player_id = $2
	AND roster_status = 'injured';
`, day.GameID, event.PlayerID); err != nil {
			return nil, fmt.Errorf("activate recovered player %s: %w", event.PlayerID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit recover players: %w", err)
	}

	return recovered, nil
}

func createInjuriesForMatch(
	ctx context.Context,
	q queryer,
	event domain.MatchFinishedEvent,
) ([]domain.PlayerInjuredEvent, error) {
	playerStates, err := loadPlayerMatchStates(ctx, q, event.GameID)
	if err != nil {
		return nil, err
	}

	var injuries []domain.PlayerInjuredEvent
	for _, line := range event.BoxScore {
		if line.TeamID != domain.OwnTeamID {
			continue
		}

		var rosterStatus string
		if err := q.QueryRow(ctx, `
SELECT roster_status
FROM team_roster_players
WHERE game_id = $1 AND player_id = $2;
`, event.GameID, line.PlayerID).Scan(&rosterStatus); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			return nil, fmt.Errorf("load roster status %s: %w", line.PlayerID, err)
		}
		if rosterStatus != "active" {
			continue
		}

		reaggravation, reaggravated, err := forcedReturnReaggravation(ctx, q, event, line)
		if err != nil {
			return nil, err
		}
		if reaggravated {
			if inserted, err := insertInjury(ctx, q, reaggravation); err != nil {
				return nil, err
			} else if inserted {
				if _, err := q.Exec(ctx, `
UPDATE team_roster_players
SET roster_status = 'injured', updated_at = NOW()
WHERE game_id = $1 AND player_id = $2;
`, reaggravation.GameID, reaggravation.PlayerID); err != nil {
					return nil, fmt.Errorf("mark reaggravated player injured %s: %w", reaggravation.PlayerID, err)
				}
				injuries = append(injuries, reaggravation)
			}
			continue
		}

		state := playerStates[line.PlayerID]
		injury, injured, err := domain.AssessInjuryRisk(domain.InjuryAssessmentInput{
			GameID:         event.GameID,
			MatchID:        event.MatchID,
			OccurredAt:     event.OccurredAt,
			SimulatedDate:  event.SimulatedDate,
			PlayerID:       line.PlayerID,
			TeamID:         line.TeamID,
			Minutes:        line.Minutes,
			RecentMinutes:  state.RecentMinutes,
			EmotionalState: state.EmotionalState,
		})
		if err != nil {
			return nil, err
		}
		if !injured {
			continue
		}

		inserted, err := insertInjury(ctx, q, injury)
		if err != nil {
			return nil, err
		}
		if !inserted {
			continue
		}

		if _, err := q.Exec(ctx, `
UPDATE team_roster_players
SET roster_status = 'injured', updated_at = NOW()
WHERE game_id = $1 AND player_id = $2;
`, injury.GameID, injury.PlayerID); err != nil {
			return nil, fmt.Errorf("mark player injured %s: %w", injury.PlayerID, err)
		}

		injuries = append(injuries, injury)
	}

	return injuries, nil
}

func forcedReturnReaggravation(
	ctx context.Context,
	q queryer,
	event domain.MatchFinishedEvent,
	line domain.PlayerBoxScore,
) (domain.PlayerInjuredEvent, bool, error) {
	var sourceInjuryID, expectedRecoveryDate string
	err := q.QueryRow(ctx, `
SELECT injury_id, expected_recovery_date
FROM team_injuries
WHERE game_id = $1
	AND player_id = $2
	AND forced_return_at IS NOT NULL
	AND expected_recovery_date > $3
ORDER BY forced_return_at DESC
LIMIT 1;
`, event.GameID, line.PlayerID, event.SimulatedDate).Scan(&sourceInjuryID, &expectedRecoveryDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.PlayerInjuredEvent{}, false, nil
		}
		return domain.PlayerInjuredEvent{}, false, fmt.Errorf("load forced return state %s: %w", line.PlayerID, err)
	}

	injuredOn, err := parseDate(event.SimulatedDate)
	if err != nil {
		return domain.PlayerInjuredEvent{}, false, err
	}
	estimatedDays := uint16(14)
	injuryID := fmt.Sprintf("injury-%s-%s-reaggravation", event.MatchID, line.PlayerID)
	return domain.PlayerInjuredEvent{
		EventMeta: domain.EventMeta{
			EventID:       fmt.Sprintf("player-injured-%s-%s-reaggravation", event.MatchID, line.PlayerID),
			GameID:        event.GameID,
			OccurredAt:    event.OccurredAt,
			SchemaVersion: 1,
		},
		InjuryID:             injuryID,
		PlayerID:             line.PlayerID,
		Severity:             "moderate",
		EstimatedDaysOut:     estimatedDays,
		InjuredOn:            event.SimulatedDate,
		ExpectedRecoveryDate: injuredOn.AddDate(0, 0, int(estimatedDays)).Format("2006-01-02"),
		Reason:               "forced_return_reaggravation",
		SourceMatchID:        event.MatchID,
		WorkloadScore:        uint16(line.Minutes),
	}, sourceInjuryID != "", nil
}

func insertInjury(ctx context.Context, q queryer, injury domain.PlayerInjuredEvent) (bool, error) {
	commandTag, err := q.Exec(ctx, `
INSERT INTO team_injuries (
	injury_id, game_id, player_id, severity, estimated_days_out, injured_on,
	expected_recovery_date, reason, source_match_id, workload_score, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
ON CONFLICT (injury_id) DO NOTHING;
`, injury.InjuryID, injury.GameID, injury.PlayerID, injury.Severity, int16(injury.EstimatedDaysOut), injury.InjuredOn, injury.ExpectedRecoveryDate, injury.Reason, injury.SourceMatchID, int16(injury.WorkloadScore))
	if err != nil {
		return false, fmt.Errorf("insert injury %s: %w", injury.InjuryID, err)
	}

	return commandTag.RowsAffected() > 0, nil
}

func parseDate(value string) (time.Time, error) {
	parsed, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse date %q: %w", value, err)
	}

	return parsed, nil
}
