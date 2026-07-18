package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/pulsecity/services/team-service/internal/domain"
)

func (s *Store) ApplyTradeAcceptance(
	ctx context.Context,
	event domain.GMDecisionRegisteredEvent,
) (*domain.TradeAcceptedEvent, *domain.RosterPatchEnvelope, *domain.SalaryCapCalculatedEvent, *domain.TradeRejectedEvent, error) {
	if event.Kind != "trade_acceptance" {
		return nil, nil, nil, nil, nil
	}

	proposalID := event.Payload["proposal_id"]
	if proposalID == "" {
		return nil, nil, nil, nil, fmt.Errorf("apply trade acceptance: missing proposal_id")
	}
	additionalAsset := event.Payload["accepted_additional_asset"]

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("begin apply trade acceptance: %w", err)
	}
	defer tx.Rollback(ctx)

	trade, err := loadTrade(ctx, tx, proposalID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			rejected := tradeRejectedFromDecision(event, proposalID, "", "proposal_not_found", "La propuesta de trade no existe o ya no esta disponible.")
			return nil, nil, nil, &rejected, tx.Commit(ctx)
		}
		return nil, nil, nil, nil, err
	}
	if trade.Status == "accepted" {
		return nil, nil, nil, nil, tx.Commit(ctx)
	}
	if trade.Status != "proposed" {
		rejected := tradeRejectedFromDecision(event, proposalID, trade.RivalTeamID, "proposal_not_open", "La propuesta ya no esta abierta para aceptacion.")
		return nil, nil, nil, &rejected, tx.Commit(ctx)
	}

	outgoing, err := loadRosterPlayer(ctx, tx, event.GameID, trade.OfferedPlayerID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if outgoing.RosterStatus != "active" {
		rejected := tradeRejectedFromDecision(event, proposalID, trade.RivalTeamID, "player_unavailable", "El jugador ofrecido ya no esta activo.")
		return nil, nil, nil, &rejected, tx.Commit(ctx)
	}

	incoming := domain.MaterializeIncomingTradePlayer(
		event.GameID,
		proposalID,
		trade.RequestedPosition,
		trade.IncomingSalary,
		outgoing.OverallRating,
		outgoing.SortOrder,
	)
	roster, err := loadRoster(ctx, tx, event.GameID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	projected := make([]domain.RosterPlayer, 0, len(roster))
	for _, player := range roster {
		if player.PlayerID == outgoing.PlayerID {
			player.RosterStatus = "traded"
		}
		projected = append(projected, player)
	}
	projected = append(projected, incoming)
	cap := domain.CalculateSalaryCap(event.GameID, projected, event.SimulatedDate, event.OccurredAt, proposalID, domain.SubjectTradeAccepted)
	if cap.CommittedSalary > domain.DefaultLuxuryTaxLine {
		rejected := tradeRejectedFromDecision(event, proposalID, trade.RivalTeamID, "cap_blocked", "Aceptar la propuesta ahora empuja a PulseCity por encima de la linea de luxury tax.")
		return nil, nil, nil, &rejected, tx.Commit(ctx)
	}

	if _, err := tx.Exec(ctx, `
UPDATE team_roster_players
SET roster_status = 'traded', updated_at = NOW()
WHERE game_id = $1 AND player_id = $2 AND roster_status = 'active';
`, event.GameID, outgoing.PlayerID); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("mark outgoing player traded %s: %w", outgoing.PlayerID, err)
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO team_roster_players (
	player_id, game_id, full_name, position, overall_rating, roster_status, contract_years, salary, sort_order, created_at, updated_at
)
VALUES ($1, $2, $3, $4, $5, 'active', $6, $7, $8, NOW(), NOW())
ON CONFLICT (player_id) DO NOTHING;
`, incoming.PlayerID, incoming.GameID, incoming.FullName, incoming.Position, int16(incoming.OverallRating), int16(incoming.ContractYears), incoming.Salary, incoming.SortOrder); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("insert incoming trade player %s: %w", incoming.PlayerID, err)
	}
	if _, err := tx.Exec(ctx, `
UPDATE team_trades
SET status = 'accepted',
	incoming_player_id = $2,
	incoming_player_name = $3,
	incoming_rating = $4,
	accepted_additional_asset = $5,
	accepted_at = $6,
	updated_at = NOW()
WHERE proposal_id = $1 AND status = 'proposed';
`, proposalID, incoming.PlayerID, incoming.FullName, int16(incoming.OverallRating), additionalAsset, event.OccurredAt); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("mark trade accepted %s: %w", proposalID, err)
	}
	if err := saveSalaryCap(ctx, tx, cap); err != nil {
		return nil, nil, nil, nil, err
	}

	accepted := domain.TradeAcceptedEvent{
		EventMeta: domain.EventMeta{
			EventID:       "trade-accepted-" + proposalID,
			GameID:        event.GameID,
			OccurredAt:    event.OccurredAt,
			SchemaVersion: 1,
		},
		ProposalID:              proposalID,
		SimulatedDate:           event.SimulatedDate,
		RivalTeamID:             trade.RivalTeamID,
		OutgoingPlayerID:        outgoing.PlayerID,
		OutgoingPlayerName:      outgoing.FullName,
		IncomingPlayerID:        incoming.PlayerID,
		IncomingPlayerName:      incoming.FullName,
		IncomingPosition:        incoming.Position,
		IncomingRating:          incoming.OverallRating,
		IncomingSalary:          incoming.Salary,
		AcceptedAdditionalAsset: additionalAsset,
	}
	rosterPatch := tradeRosterPatch(accepted)
	salaryCapEvent := cap.SalaryCapCalculatedEvent()

	return &accepted, &rosterPatch, &salaryCapEvent, nil, tx.Commit(ctx)
}

func loadTrade(ctx context.Context, q queryer, proposalID string) (storedTrade, error) {
	var trade storedTrade
	if err := q.QueryRow(ctx, `
SELECT proposal_id, game_id, rival_team_id, offered_player_id, offered_player_name,
	offered_salary, requested_position, incoming_salary, cap_space_after, status,
	simulated_date
FROM team_trades
WHERE proposal_id = $1;
`, proposalID).Scan(
		&trade.ProposalID,
		&trade.GameID,
		&trade.RivalTeamID,
		&trade.OfferedPlayerID,
		&trade.OfferedPlayerName,
		&trade.OfferedSalary,
		&trade.RequestedPosition,
		&trade.IncomingSalary,
		&trade.CapSpaceAfter,
		&trade.Status,
		&trade.SimulatedDate,
	); err != nil {
		return storedTrade{}, fmt.Errorf("load trade %s: %w", proposalID, err)
	}

	return trade, nil
}

type storedTrade struct {
	ProposalID        string
	GameID            string
	RivalTeamID       string
	OfferedPlayerID   string
	OfferedPlayerName string
	OfferedSalary     int
	RequestedPosition string
	IncomingSalary    int
	CapSpaceAfter     int
	Status            string
	SimulatedDate     string
}

func tradeRejectedFromDecision(event domain.GMDecisionRegisteredEvent, proposalID, rivalTeamID, reason, detail string) domain.TradeRejectedEvent {
	return domain.TradeRejectedEvent{
		EventMeta: domain.EventMeta{
			EventID:       "trade-rejected-" + proposalID,
			GameID:        event.GameID,
			OccurredAt:    event.OccurredAt,
			SchemaVersion: 1,
		},
		ProposalID:    proposalID,
		SimulatedDate: event.SimulatedDate,
		RivalTeamID:   rivalTeamID,
		Reason:        reason,
		Detail:        detail,
	}
}

func tradeRosterPatch(event domain.TradeAcceptedEvent) domain.RosterPatchEnvelope {
	return domain.RosterPatchEnvelope{
		Type:    domain.SubjectRosterPatchDelta,
		Subject: domain.SubjectTradeAccepted,
		GameID:  event.GameID,
		Patch: domain.RosterStatePatch{
			SimulatedDate: event.SimulatedDate,
			SourceEventID: event.EventID,
			SourceSubject: domain.SubjectTradeAccepted,
			Players: []domain.PlayerEmotionalPatch{
				{
					PlayerID:       event.OutgoingPlayerID,
					EmotionalState: "traded",
					Summary:        event.OutgoingPlayerName + " fue traspasado.",
					Availability:   "traded",
					FullName:       event.OutgoingPlayerName,
				},
				{
					PlayerID:       event.IncomingPlayerID,
					EmotionalState: "arriving",
					Summary:        event.IncomingPlayerName + " llega a PulseCity via trade.",
					Availability:   "active",
					FullName:       event.IncomingPlayerName,
					Position:       event.IncomingPosition,
					OverallRating:  event.IncomingRating,
					Salary:         event.IncomingSalary,
				},
			},
		},
	}
}
