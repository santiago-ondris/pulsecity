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

func (d Dependencies) startGame(w http.ResponseWriter, r *http.Request) {
	currentActor, ok := d.requireActor(w, r)
	if !ok {
		return
	}

	var request domain.StartGameRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}

	command := domain.MapGenerationRequest{
		GameID:        uuid.NewString(),
		CityName:      request.CityName,
		FranchiseName: request.FranchiseName,
		Abbreviation:  request.Abbreviation,
	}

	setup := domain.GameSetup{
		GameID:             command.GameID,
		GuestToken:         currentActor.guestToken,
		UserID:             currentActor.user.UserID,
		CityName:           normalizeText(request.CityName, "Nueva Aurora"),
		FranchiseName:      normalizeText(request.FranchiseName, "Lighthouses"),
		Abbreviation:       normalizeAbbreviation(request.Abbreviation),
		PrimaryColor:       normalizeColor(request.PrimaryColor, "#00C896"),
		SecondaryColor:     normalizeColor(request.SecondaryColor, "#7B8CDE"),
		AccentColor:        normalizeColor(request.AccentColor, "#FFAA00"),
		InitialScenario:    normalizeScenario(request.InitialScenario),
		CityManagementMode: normalizeCityManagementMode(request.CityManagementMode),
		Status:             "generation_started",
	}
	command.CityName = setup.CityName
	command.FranchiseName = setup.FranchiseName
	command.Abbreviation = setup.Abbreviation

	if err := d.Store.CreateGame(r.Context(), setup); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to persist game",
		})
		return
	}

	if err := d.Bus.PublishJSON(domain.SubjectMapGenerationStarted, command); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{
			"error": "failed to publish map generation request",
		})
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"game_id": command.GameID,
		"status":  "map_generation_started",
	})
}

func (d Dependencies) listGames(w http.ResponseWriter, r *http.Request) {
	currentActor, ok := d.requireActor(w, r)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var (
		games []domain.GameSummary
		err   error
	)
	if currentActor.kind == "user" {
		games, err = d.Store.ListGamesByUser(ctx, currentActor.user.UserID)
	} else {
		games, err = d.Store.ListGamesByGuest(ctx, currentActor.guestToken)
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to list games",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"games": games,
	})
}

func (d Dependencies) getGame(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusOK, game)
}

func (d Dependencies) getSnapshot(w http.ResponseWriter, r *http.Request) {
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
			"error": "snapshot not found",
		})
		return
	}

	snapshot, ok := d.Snapshots.Get(gameID)
	if !ok {
		rehydrated, found, err := d.Store.GetSnapshot(r.Context(), gameID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to load snapshot",
			})
			return
		}
		if !found {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "snapshot not found",
			})
			return
		}
		d.Snapshots.Set(rehydrated)
		snapshot = rehydrated
	}

	writeJSON(w, http.StatusOK, domain.MapSnapshotEnvelope{
		Type:    "map.snapshot",
		Subject: "gateway.snapshot_http",
		State:   snapshot,
	})
}

func normalizeText(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}

	return trimmed
}

func normalizeAbbreviation(value string) string {
	trimmed := strings.ToUpper(strings.TrimSpace(value))
	if len(trimmed) == 0 {
		return "NEW"
	}

	if len(trimmed) > 3 {
		return trimmed[:3]
	}

	return trimmed
}

func normalizeColor(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}

	return trimmed
}

func normalizeScenario(value string) string {
	switch value {
	case "rebuild", "contention", "decline", "expansion":
		return value
	default:
		return "expansion"
	}
}

func normalizeCityManagementMode(value string) string {
	switch value {
	case "owner_influence", "dual_figure":
		return value
	default:
		return "owner_influence"
	}
}
