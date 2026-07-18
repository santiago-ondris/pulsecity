package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pulsecity/services/gateway/internal/domain"
)

func (d Dependencies) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	guestToken := strings.TrimSpace(r.URL.Query().Get("guest_token"))
	sessionToken := strings.TrimSpace(r.URL.Query().Get("session_token"))
	d.Hub.ServeWebSocket(w, r, func(conn *websocket.Conn) (func(), error) {
		if gameID == "" {
			return nil, nil
		}
		if sessionToken != "" && guestToken != "" {
			return nil, nil
		}

		game, found, err := d.Store.GetGame(r.Context(), gameID)
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, nil
		}
		if sessionToken != "" {
			ok, err := d.Store.TouchUserSession(r.Context(), sessionToken)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, nil
			}
			session, sessionFound, err := d.Store.GetUserSession(r.Context(), sessionToken)
			if err != nil {
				return nil, err
			}
			if !sessionFound || game.UserID != session.User.UserID {
				return nil, nil
			}
		} else {
			ok, err := d.Store.TouchGuestSession(r.Context(), guestToken)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, nil
			}
			if !guestOwnsGame(guestToken, game) {
				return nil, nil
			}
		}

		snapshot, ok := d.Snapshots.Get(gameID)
		if !ok {
			rehydrated, found, err := d.Store.GetSnapshot(r.Context(), gameID)
			if err != nil {
				return nil, err
			}
			if !found {
				return nil, nil
			}
			d.Snapshots.Set(rehydrated)
			snapshot = rehydrated
		}

		if err := conn.WriteJSON(domain.MapSnapshotEnvelope{
			Type:    "map.snapshot",
			Subject: "gateway.snapshot_rehidratado",
			State:   snapshot,
		}); err != nil {
			return nil, err
		}

		simulationSessionID := uuid.NewString()
		if d.Hub.ActivateGame(gameID, simulationSessionID) {
			now := time.Now().UTC().Format(time.RFC3339Nano)
			if err := d.Bus.PublishJSON(domain.SubjectTimeSessionStarted, domain.TimeSessionStartedEvent{
				EventMeta: domain.EventMeta{
					EventID:       uuid.NewString(),
					GameID:        gameID,
					OccurredAt:    now,
					SchemaVersion: 1,
				},
				SessionID: simulationSessionID,
				ClientID:  clientIDFromRequest(r),
			}); err != nil {
				d.Hub.DeactivateGame(gameID)
				return nil, err
			}
		}

		return func() {
			activeSessionID, isLastConnection := d.Hub.DeactivateGame(gameID)
			if !isLastConnection {
				return
			}
			if err := d.Bus.PublishJSON(domain.SubjectTimeSessionEnded, domain.TimeSessionEndedEvent{
				EventMeta: domain.EventMeta{
					EventID:       uuid.NewString(),
					GameID:        gameID,
					OccurredAt:    time.Now().UTC().Format(time.RFC3339Nano),
					SchemaVersion: 1,
				},
				SessionID: activeSessionID,
				Reason:    "client_closed",
			}); err != nil {
				// No se puede responder al cliente durante el cierre del socket; loguear alcanza.
				// El siguiente connect volvera a activar la sesion.
				_ = err
			}
		}, nil
	})
}

func clientIDFromRequest(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		return forwarded
	}
	if r.RemoteAddr != "" {
		return r.RemoteAddr
	}

	return "browser"
}

func guestOwnsGame(guestToken string, game domain.GameSetup) bool {
	return guestToken != "" && game.GuestToken == guestToken
}
