package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pulsecity/services/gateway/internal/domain"
)

func (d Dependencies) answerMedicalDecision(w http.ResponseWriter, r *http.Request) {
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

	var request domain.MedicalDecisionRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}
	request.InjuryID = strings.TrimSpace(request.InjuryID)
	request.PlayerID = strings.TrimSpace(request.PlayerID)
	request.ChoiceID = strings.TrimSpace(request.ChoiceID)
	request.SimulatedDate = strings.TrimSpace(request.SimulatedDate)
	if request.InjuryID == "" || request.PlayerID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "injury_id and player_id are required",
		})
		return
	}
	choiceLabel, ok := medicalDecisionLabel(request.ChoiceID)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid medical decision",
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
	decisionEvent := domain.GMDecisionRegisteredEvent{
		EventMeta: domain.EventMeta{
			EventID:       "decision-medical-" + gameID + "-" + request.InjuryID,
			GameID:        gameID,
			OccurredAt:    now,
			SchemaVersion: 1,
		},
		DecisionID:    "medical-" + request.InjuryID,
		Kind:          "medical_decision",
		SimulatedDate: request.SimulatedDate,
		Payload: map[string]string{
			"injury_id":    request.InjuryID,
			"player_id":    request.PlayerID,
			"choice_id":    request.ChoiceID,
			"choice_label": choiceLabel,
		},
		AgentsAffected: []string{"team_doctor", "strength_conditioning_coach", "head_coach"},
		SourceEventID:  request.InjuryID,
		SourceSubject:  domain.SubjectPlayerInjured,
	}
	if err := d.Bus.PublishJSON(domain.SubjectGMDecisionRegistered, decisionEvent); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": "failed to publish medical decision",
		})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"decision_id": decisionEvent.DecisionID,
		"choice_id":   request.ChoiceID,
		"status":      "queued",
	})
}

func medicalDecisionLabel(choiceID string) (string, bool) {
	switch choiceID {
	case "rest":
		return "Seguir el protocolo medico", true
	case "reduce_minutes":
		return "Reducir carga al volver", true
	case "ignore_doctor":
		return "Ignorar la recomendacion medica", true
	case "force_return":
		return "Forzar alta anticipada", true
	default:
		return "", false
	}
}
