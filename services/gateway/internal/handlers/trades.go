package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pulsecity/services/gateway/internal/domain"
)

func (d Dependencies) proposeTrade(w http.ResponseWriter, r *http.Request) {
	currentActor, ok := d.requireActor(w, r)
	if !ok {
		return
	}

	gameID := r.PathValue("gameID")
	if gameID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing game id",
		})
		return
	}

	game, found, err := d.Store.GetGame(r.Context(), gameID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load game",
		})
		return
	}
	if !found || !gameOwnedBy(currentActor, game) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "game not found",
		})
		return
	}

	var request domain.TradeProposalRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}
	request.RivalTeamID = strings.ToLower(strings.TrimSpace(request.RivalTeamID))
	request.OfferedPlayerID = strings.TrimSpace(request.OfferedPlayerID)
	request.RequestedPosition = strings.ToUpper(strings.TrimSpace(request.RequestedPosition))
	request.SimulatedDate = strings.TrimSpace(request.SimulatedDate)
	if request.RivalTeamID == "" || request.OfferedPlayerID == "" || request.RequestedPosition == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "rival_team_id, offered_player_id and requested_position are required",
		})
		return
	}
	if request.IncomingSalary <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "incoming_salary must be positive",
		})
		return
	}
	if request.SimulatedDate == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "simulated_date is required",
		})
		return
	}

	proposalID := "trade-" + uuid.NewString()
	now := time.Now().UTC().Format(time.RFC3339)
	decisionEvent := domain.GMDecisionRegisteredEvent{
		EventMeta: domain.EventMeta{
			EventID:       "decision-" + proposalID,
			GameID:        gameID,
			OccurredAt:    now,
			SchemaVersion: 1,
		},
		DecisionID:    proposalID,
		Kind:          "trade_proposal",
		SimulatedDate: request.SimulatedDate,
		Payload: map[string]string{
			"proposal_id":        proposalID,
			"rival_team_id":      request.RivalTeamID,
			"offered_player_id":  request.OfferedPlayerID,
			"requested_position": request.RequestedPosition,
			"incoming_salary":    strconv.Itoa(request.IncomingSalary),
		},
		AgentsAffected: []string{"director_player_personnel", "cfo", "rival_gm_" + request.RivalTeamID},
		SourceEventID:  proposalID,
		SourceSubject:  "trade.propuesta_gm",
	}
	if err := d.Bus.PublishJSON(domain.SubjectGMDecisionRegistered, decisionEvent); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": "failed to publish trade proposal",
		})
		return
	}

	writeJSON(w, http.StatusAccepted, domain.TradeProposalAcceptedResponse{
		ProposalID: proposalID,
		Status:     "queued",
	})
}

func (d Dependencies) acceptTrade(w http.ResponseWriter, r *http.Request) {
	currentActor, ok := d.requireActor(w, r)
	if !ok {
		return
	}

	gameID := r.PathValue("gameID")
	if gameID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing game id",
		})
		return
	}

	game, found, err := d.Store.GetGame(r.Context(), gameID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load game",
		})
		return
	}
	if !found || !gameOwnedBy(currentActor, game) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "game not found",
		})
		return
	}

	var request domain.TradeAcceptanceRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}
	request.ProposalID = strings.TrimSpace(request.ProposalID)
	request.AcceptedAdditionalAsset = strings.TrimSpace(request.AcceptedAdditionalAsset)
	request.SimulatedDate = strings.TrimSpace(request.SimulatedDate)
	if request.ProposalID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "proposal_id is required",
		})
		return
	}
	if request.SimulatedDate == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "simulated_date is required",
		})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	decisionID := "accept-" + request.ProposalID
	decisionEvent := domain.GMDecisionRegisteredEvent{
		EventMeta: domain.EventMeta{
			EventID:       "decision-" + decisionID,
			GameID:        gameID,
			OccurredAt:    now,
			SchemaVersion: 1,
		},
		DecisionID:    decisionID,
		Kind:          "trade_acceptance",
		SimulatedDate: request.SimulatedDate,
		Payload: map[string]string{
			"proposal_id":               request.ProposalID,
			"accepted_additional_asset": request.AcceptedAdditionalAsset,
		},
		AgentsAffected: []string{"director_player_personnel", "cfo"},
		SourceEventID:  request.ProposalID,
		SourceSubject:  domain.SubjectTradeCountered,
	}
	if err := d.Bus.PublishJSON(domain.SubjectGMDecisionRegistered, decisionEvent); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": "failed to publish trade acceptance",
		})
		return
	}

	writeJSON(w, http.StatusAccepted, domain.TradeProposalAcceptedResponse{
		ProposalID: request.ProposalID,
		Status:     "queued",
	})
}
