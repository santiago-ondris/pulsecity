package handlers

import (
	"encoding/json"
	"net/http"

	natsclient "github.com/pulsecity/services/gateway/internal/nats"
	"github.com/pulsecity/services/gateway/internal/persistence"
	"github.com/pulsecity/services/gateway/internal/state"
	"github.com/pulsecity/services/gateway/internal/ws"
)

type Dependencies struct {
	Bus       *natsclient.Client
	Hub       *ws.Hub
	Store     *persistence.Store
	Snapshots *state.MapSnapshots
}

func RegisterRoutes(mux *http.ServeMux, deps Dependencies) {
	mux.HandleFunc("GET /", debugPage)
	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /ws", deps.serveWebSocket)
	mux.HandleFunc("POST /api/v1/auth/register", deps.register)
	mux.HandleFunc("POST /api/v1/auth/login", deps.login)
	mux.HandleFunc("GET /api/v1/auth/session", deps.getCurrentSession)
	mux.HandleFunc("POST /api/v1/auth/upgrade-guest", deps.upgradeGuest)
	mux.HandleFunc("POST /api/v1/guest-sessions", deps.createGuestSession)
	mux.HandleFunc("POST /api/v1/games", deps.startGame)
	mux.HandleFunc("GET /api/v1/games", deps.listGames)
	mux.HandleFunc("GET /api/v1/games/{gameID}", deps.getGame)
	mux.HandleFunc("POST /api/v1/games/{gameID}/owner-intro-response", deps.answerOwnerIntro)
	mux.HandleFunc("POST /api/v1/games/{gameID}/agent-chat", deps.startAgentChat)
	mux.HandleFunc("POST /api/v1/games/{gameID}/medical-decisions", deps.answerMedicalDecision)
	mux.HandleFunc("POST /api/v1/games/{gameID}/trades/proposals", deps.proposeTrade)
	mux.HandleFunc("POST /api/v1/games/{gameID}/trades/acceptances", deps.acceptTrade)
	mux.HandleFunc("GET /api/v1/games/{gameID}/snapshot", deps.getSnapshot)
	mux.HandleFunc("POST /api/v1/games/{gameID}/time-control", deps.updateTimeControl)
}

func debugPage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(debugHTML))
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
