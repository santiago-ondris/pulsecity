package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pulsecity/services/gateway/internal/domain"
)

func (d Dependencies) updateTimeControl(w http.ResponseWriter, r *http.Request) {
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

	var request domain.TimeControlRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}
	if request.Speed == nil && request.Paused == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing time control change",
		})
		return
	}
	if request.Speed != nil && !validTimeSpeed(*request.Speed) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid speed",
		})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	if request.Speed != nil {
		if err := d.Bus.PublishJSON(domain.SubjectTimeSpeedChanged, domain.TimeSpeedChangedEvent{
			EventMeta: domain.EventMeta{
				EventID:       uuid.NewString(),
				GameID:        gameID,
				OccurredAt:    now,
				SchemaVersion: 1,
			},
			Speed: *request.Speed,
		}); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{
				"error": "failed to publish speed change",
			})
			return
		}
	}
	if request.Paused != nil {
		if err := d.Bus.PublishJSON(domain.SubjectTimePauseChanged, domain.TimePauseChangedEvent{
			EventMeta: domain.EventMeta{
				EventID:       uuid.NewString(),
				GameID:        gameID,
				OccurredAt:    now,
				SchemaVersion: 1,
			},
			Paused: *request.Paused,
		}); err != nil {
			writeJSON(w, http.StatusBadGateway, map[string]string{
				"error": "failed to publish pause change",
			})
			return
		}
	}

	d.Hub.Broadcast(domain.TimePatchEnvelope{
		Type:    "time.patch",
		Subject: "gateway.time_control",
		GameID:  gameID,
		Patch: domain.TimeStatePatch{
			Speed:  request.Speed,
			Paused: request.Paused,
		},
	})

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status": "time_control_published",
	})
}

func validTimeSpeed(speed uint8) bool {
	return speed == 1 || speed == 5 || speed == 20
}
