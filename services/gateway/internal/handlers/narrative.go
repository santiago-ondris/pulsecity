package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pulsecity/services/gateway/internal/domain"
)

func (d Dependencies) answerOwnerIntro(w http.ResponseWriter, r *http.Request) {
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
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "game not found",
		})
		return
	}
	if !gameOwnedBy(currentActor, game) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "game not found",
		})
		return
	}
	if game.OwnerIntroEvent == nil {
		writeJSON(w, http.StatusConflict, map[string]string{
			"error": "owner intro event not available yet",
		})
		return
	}
	if game.OwnerIntroResponse != nil {
		writeJSON(w, http.StatusOK, game.OwnerIntroResponse)
		return
	}

	var request domain.OwnerIntroResponseRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}

	choice, ok := findNarrativeChoice(game.OwnerIntroEvent.Choices, request.ChoiceID)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid owner intro choice",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := d.Store.SetOwnerIntroResponse(ctx, gameID, choice); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to persist owner intro response",
		})
		return
	}

	responseEvent := domain.NarrativeResponseEvent{
		Type:    "narrative.response",
		Subject: "narrativa.respuesta_gm",
		GameID:  gameID,
		EventID: game.OwnerIntroEvent.EventID,
		Choice:  choice,
		Emitter: "gm",
		Metadata: map[string]string{
			"city_name":            game.CityName,
			"franchise_name":       game.FranchiseName,
			"initial_scenario":     game.InitialScenario,
			"city_management_mode": game.CityManagementMode,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	if err := d.Bus.PublishJSON("narrativa.respuesta_gm", responseEvent); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": "failed to publish owner intro response",
		})
		return
	}
	decisionEvent := domain.GMDecisionRegisteredEvent{
		EventMeta: domain.EventMeta{
			EventID:       "decision-owner-intro-" + gameID,
			GameID:        gameID,
			OccurredAt:    responseEvent.Timestamp,
			SchemaVersion: 1,
		},
		DecisionID:    "owner-intro-" + gameID,
		Kind:          "owner_intro_response",
		SimulatedDate: "2026-10-01",
		Payload: map[string]string{
			"choice_id":            choice.ID,
			"choice_label":         choice.Label,
			"city_name":            game.CityName,
			"franchise_name":       game.FranchiseName,
			"initial_scenario":     game.InitialScenario,
			"city_management_mode": game.CityManagementMode,
		},
		AgentsAffected: []string{"owner", "president_basketball_ops", "ceo_business_ops"},
		SourceEventID:  game.OwnerIntroEvent.EventID,
		SourceSubject:  "narrativa.respuesta_gm",
	}
	if err := d.Bus.PublishJSON(domain.SubjectGMDecisionRegistered, decisionEvent); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": "failed to publish gm decision",
		})
		return
	}

	writeJSON(w, http.StatusAccepted, choice)
}

func (d Dependencies) startAgentChat(w http.ResponseWriter, r *http.Request) {
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

	var request domain.AgentChatRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}
	request.AgentID = strings.TrimSpace(request.AgentID)
	request.Message = strings.TrimSpace(request.Message)
	request.ConversationID = strings.TrimSpace(request.ConversationID)
	if request.AgentID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing agent id",
		})
		return
	}
	if request.Message == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing message",
		})
		return
	}
	if len(request.Message) > 1200 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "message too long",
		})
		return
	}
	if request.ConversationID == "" {
		request.ConversationID = "chat-" + uuid.NewString()
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	eventID := "agent-consultation-" + uuid.NewString()
	event := domain.AgentConsultationStartedEvent{
		EventMeta: domain.EventMeta{
			EventID:       eventID,
			GameID:        gameID,
			OccurredAt:    now,
			SchemaVersion: 1,
		},
		ConversationID: request.ConversationID,
		AgentID:        request.AgentID,
		Sender:         "gm",
		Message:        request.Message,
	}
	if err := d.Bus.PublishJSON(domain.SubjectAgentConsultationStarted, event); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": "failed to publish agent chat request",
		})
		return
	}

	writeJSON(w, http.StatusAccepted, domain.AgentChatAcceptedResponse{
		ConversationID: request.ConversationID,
		RequestEventID: eventID,
		Status:         "agent_chat_requested",
	})
}

func findNarrativeChoice(choices []domain.NarrativeChoice, choiceID string) (domain.NarrativeChoice, bool) {
	for _, choice := range choices {
		if choice.ID == choiceID {
			return choice, true
		}
	}

	return domain.NarrativeChoice{}, false
}
