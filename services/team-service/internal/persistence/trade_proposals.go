package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pulsecity/services/team-service/internal/domain"
)

func (s *Store) ApplyTradeProposal(
	ctx context.Context,
	event domain.GMDecisionRegisteredEvent,
) (*domain.TradeProposedEvent, *domain.TradeRejectedEvent, error) {
	if event.Kind != "trade_proposal" {
		return nil, nil, nil
	}

	proposalID := event.Payload["proposal_id"]
	rivalTeamID := event.Payload["rival_team_id"]
	offeredPlayerID := event.Payload["offered_player_id"]
	requestedPosition := event.Payload["requested_position"]
	incomingSalary, err := parsePositiveInt(event.Payload["incoming_salary"])
	if proposalID == "" || rivalTeamID == "" || offeredPlayerID == "" || requestedPosition == "" || err != nil {
		return nil, nil, fmt.Errorf("apply trade proposal: invalid proposal payload")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("begin apply trade proposal: %w", err)
	}
	defer tx.Rollback(ctx)

	_, found, err := loadTradeStatus(ctx, tx, proposalID)
	if err != nil {
		return nil, nil, err
	}
	if found {
		return nil, nil, tx.Commit(ctx)
	}

	player, err := loadRosterPlayer(ctx, tx, event.GameID, offeredPlayerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			rejected := tradeRejectedFromDecision(event, proposalID, rivalTeamID, "player_not_found", "El jugador ofrecido no pertenece al roster activo de PulseCity.")
			if err := insertRejectedTrade(ctx, tx, event, rejected, offeredPlayerID, requestedPosition, incomingSalary); err != nil {
				return nil, nil, err
			}
			return nil, &rejected, tx.Commit(ctx)
		}
		return nil, nil, err
	}
	if player.RosterStatus != "active" {
		rejected := tradeRejectedFromDecision(event, proposalID, rivalTeamID, "player_unavailable", "El jugador ofrecido no esta activo para negociar.")
		if err := insertRejectedTrade(ctx, tx, event, rejected, offeredPlayerID, requestedPosition, incomingSalary); err != nil {
			return nil, nil, err
		}
		return nil, &rejected, tx.Commit(ctx)
	}

	roster, err := loadRoster(ctx, tx, event.GameID)
	if err != nil {
		return nil, nil, err
	}
	projected := make([]domain.RosterPlayer, 0, len(roster))
	for _, rosterPlayer := range roster {
		if rosterPlayer.PlayerID == offeredPlayerID {
			continue
		}
		projected = append(projected, rosterPlayer)
	}
	projected = append(projected, domain.RosterPlayer{
		PlayerID:      proposalID + "-incoming",
		GameID:        event.GameID,
		Position:      requestedPosition,
		OverallRating: player.OverallRating,
		RosterStatus:  "active",
		ContractYears: 1,
		Salary:        incomingSalary,
	})
	cap := domain.CalculateSalaryCap(event.GameID, projected, event.SimulatedDate, event.OccurredAt, proposalID, domain.SubjectTradeProposed)
	if cap.CommittedSalary > domain.DefaultLuxuryTaxLine {
		rejected := tradeRejectedFromDecision(event, proposalID, rivalTeamID, "cap_blocked", "La propuesta empuja a PulseCity por encima de la linea de luxury tax.")
		if err := insertRejectedTradeWithPlayer(ctx, tx, event, rejected, player, requestedPosition, incomingSalary, cap.LuxuryTaxSpace); err != nil {
			return nil, nil, err
		}
		return nil, &rejected, tx.Commit(ctx)
	}

	proposed := domain.TradeProposedEvent{
		EventMeta: domain.EventMeta{
			EventID:       "trade-proposed-" + proposalID,
			GameID:        event.GameID,
			OccurredAt:    event.OccurredAt,
			SchemaVersion: 1,
		},
		ProposalID:        proposalID,
		SimulatedDate:     event.SimulatedDate,
		RivalTeamID:       rivalTeamID,
		OfferedPlayerID:   offeredPlayerID,
		OfferedPlayerName: player.FullName,
		OfferedSalary:     player.Salary,
		RequestedPosition: requestedPosition,
		IncomingSalary:    incomingSalary,
		CapSpaceAfter:     cap.CapSpace,
	}
	if err := insertProposedTrade(ctx, tx, event, proposed); err != nil {
		return nil, nil, err
	}

	return &proposed, nil, tx.Commit(ctx)
}

func loadTradeStatus(ctx context.Context, q queryer, proposalID string) (string, bool, error) {
	var status string
	if err := q.QueryRow(ctx, `
SELECT status
FROM team_trades
WHERE proposal_id = $1;
`, proposalID).Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("load trade status %s: %w", proposalID, err)
	}

	return status, true, nil
}

func insertProposedTrade(ctx context.Context, q queryer, event domain.GMDecisionRegisteredEvent, proposed domain.TradeProposedEvent) error {
	if _, err := q.Exec(ctx, `
INSERT INTO team_trades (
	proposal_id, game_id, rival_team_id, offered_player_id, offered_player_name,
	offered_salary, requested_position, incoming_salary, cap_space_after, status,
	source_decision_id, source_event_id, simulated_date, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'proposed', $10, $11, $12, NOW(), NOW());
`, proposed.ProposalID, proposed.GameID, proposed.RivalTeamID, proposed.OfferedPlayerID, proposed.OfferedPlayerName, proposed.OfferedSalary, proposed.RequestedPosition, proposed.IncomingSalary, proposed.CapSpaceAfter, event.DecisionID, event.EventID, event.SimulatedDate); err != nil {
		return fmt.Errorf("insert proposed trade %s: %w", proposed.ProposalID, err)
	}
	return nil
}

func insertRejectedTrade(ctx context.Context, q queryer, event domain.GMDecisionRegisteredEvent, rejected domain.TradeRejectedEvent, offeredPlayerID, requestedPosition string, incomingSalary int) error {
	return insertRejectedTradeWithPlayer(ctx, q, event, rejected, domain.RosterPlayer{
		PlayerID: offeredPlayerID,
		GameID:   event.GameID,
	}, requestedPosition, incomingSalary, 0)
}

func insertRejectedTradeWithPlayer(ctx context.Context, q queryer, event domain.GMDecisionRegisteredEvent, rejected domain.TradeRejectedEvent, player domain.RosterPlayer, requestedPosition string, incomingSalary, capSpaceAfter int) error {
	if _, err := q.Exec(ctx, `
INSERT INTO team_trades (
	proposal_id, game_id, rival_team_id, offered_player_id, offered_player_name,
	offered_salary, requested_position, incoming_salary, cap_space_after, status,
	rejection_reason, rejection_detail, source_decision_id, source_event_id,
	simulated_date, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'rejected', $10, $11, $12, $13, $14, NOW(), NOW());
`, rejected.ProposalID, rejected.GameID, rejected.RivalTeamID, player.PlayerID, player.FullName, player.Salary, requestedPosition, incomingSalary, capSpaceAfter, rejected.Reason, rejected.Detail, event.DecisionID, event.EventID, event.SimulatedDate); err != nil {
		return fmt.Errorf("insert rejected trade %s: %w", rejected.ProposalID, err)
	}
	return nil
}

func parsePositiveInt(value string) (int, error) {
	var parsed int
	if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil {
		return 0, err
	}
	if parsed <= 0 {
		return 0, fmt.Errorf("value must be positive")
	}
	return parsed, nil
}
