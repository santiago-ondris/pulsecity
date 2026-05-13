package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pulsecity/services/gateway/internal/domain"
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
	mux.HandleFunc("POST /api/v1/guest-sessions", deps.createGuestSession)
	mux.HandleFunc("POST /api/v1/games", deps.startGame)
	mux.HandleFunc("GET /api/v1/games", deps.listGames)
	mux.HandleFunc("GET /api/v1/games/{gameID}", deps.getGame)
	mux.HandleFunc("POST /api/v1/games/{gameID}/owner-intro-response", deps.answerOwnerIntro)
	mux.HandleFunc("GET /api/v1/games/{gameID}/snapshot", deps.getSnapshot)
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

func (d Dependencies) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("game_id")
	guestToken := strings.TrimSpace(r.URL.Query().Get("guest_token"))
	d.Hub.ServeWebSocket(w, r, func(conn *websocket.Conn) error {
		if gameID == "" {
			return nil
		}

		game, found, err := d.Store.GetGame(r.Context(), gameID)
		if err != nil {
			return err
		}
		if !found || !guestOwnsGame(guestToken, game) {
			return nil
		}

		snapshot, ok := d.Snapshots.Get(gameID)
		if !ok {
			rehydrated, found, err := d.Store.GetSnapshot(r.Context(), gameID)
			if err != nil {
				return err
			}
			if !found {
				return nil
			}
			d.Snapshots.Set(rehydrated)
			snapshot = rehydrated
		}

		return conn.WriteJSON(domain.MapSnapshotEnvelope{
			Type:    "map.snapshot",
			Subject: "gateway.snapshot_rehidratado",
			State:   snapshot,
		})
	})
}

func (d Dependencies) startGame(w http.ResponseWriter, r *http.Request) {
	guestToken, ok := d.requireGuestToken(w, r)
	if !ok {
		return
	}

	var request domain.StartGameRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}

	command := domain.MapGenerationRequest{
		GameID:   uuid.NewString(),
		CityName: request.CityName,
	}

	setup := domain.GameSetup{
		GameID:             command.GameID,
		GuestToken:         guestToken,
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

func (d Dependencies) createGuestSession(w http.ResponseWriter, r *http.Request) {
	token := "guest_" + uuid.NewString()
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := d.Store.CreateGuestSession(ctx, token); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to create guest session",
		})
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	writeJSON(w, http.StatusCreated, domain.GuestSession{
		GuestToken: token,
		CreatedAt:  now,
		LastSeenAt: now,
	})
}

func (d Dependencies) listGames(w http.ResponseWriter, r *http.Request) {
	guestToken, ok := d.requireGuestToken(w, r)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	games, err := d.Store.ListGamesByGuest(ctx, guestToken)
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
	guestToken, ok := d.requireGuestToken(w, r)
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
	if !guestOwnsGame(guestToken, game) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "game not found",
		})
		return
	}

	writeJSON(w, http.StatusOK, game)
}

func (d Dependencies) getSnapshot(w http.ResponseWriter, r *http.Request) {
	guestToken, ok := d.requireGuestToken(w, r)
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
	if !found || !guestOwnsGame(guestToken, game) {
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

func (d Dependencies) answerOwnerIntro(w http.ResponseWriter, r *http.Request) {
	guestToken, ok := d.requireGuestToken(w, r)
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
	if !guestOwnsGame(guestToken, game) {
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

	writeJSON(w, http.StatusAccepted, choice)
}

func (d Dependencies) requireGuestToken(w http.ResponseWriter, r *http.Request) (string, bool) {
	token := guestTokenFromRequest(r)
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "missing guest token",
		})
		return "", false
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	ok, err := d.Store.TouchGuestSession(ctx, token)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to validate guest token",
		})
		return "", false
	}
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "invalid guest token",
		})
		return "", false
	}

	return token, true
}

func guestTokenFromRequest(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get("X-Guest-Token"))
}

func guestOwnsGame(guestToken string, game domain.GameSetup) bool {
	return guestToken != "" && game.GuestToken == guestToken
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
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

func findNarrativeChoice(choices []domain.NarrativeChoice, choiceID string) (domain.NarrativeChoice, bool) {
	for _, choice := range choices {
		if choice.ID == choiceID {
			return choice, true
		}
	}

	return domain.NarrativeChoice{}, false
}

const debugHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>PulseCity Gateway Debug</title>
    <style>
      :root {
        color-scheme: light;
        --bg: #f4efe7;
        --panel: #fffdf8;
        --ink: #1e2a2f;
        --muted: #66757d;
        --accent: #0f766e;
        --accent-strong: #115e59;
        --line: #d7cec0;
        --warn: #8a3b12;
      }

      * {
        box-sizing: border-box;
      }

      body {
        margin: 0;
        font-family: Georgia, "Times New Roman", serif;
        background:
          radial-gradient(circle at top left, rgba(15, 118, 110, 0.16), transparent 24%),
          linear-gradient(180deg, #f8f4ec 0%, var(--bg) 100%);
        color: var(--ink);
      }

      main {
        max-width: 980px;
        margin: 0 auto;
        padding: 48px 20px 80px;
      }

      h1 {
        margin: 0 0 8px;
        font-size: clamp(2rem, 4vw, 3.75rem);
        line-height: 0.95;
      }

      p {
        margin: 0;
        color: var(--muted);
      }

      .grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: 18px;
        margin-top: 28px;
      }

      .panel {
        background: color-mix(in srgb, var(--panel) 94%, white 6%);
        border: 1px solid var(--line);
        border-radius: 20px;
        padding: 18px;
        box-shadow: 0 18px 36px rgba(30, 42, 47, 0.08);
      }

      .panel h2 {
        margin: 0 0 14px;
        font-size: 1.15rem;
      }

      label {
        display: block;
        font-size: 0.85rem;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: var(--muted);
        margin-bottom: 8px;
      }

      input,
      button {
        width: 100%;
        border-radius: 14px;
        border: 1px solid var(--line);
        padding: 13px 14px;
        font: inherit;
      }

      input {
        background: #fff;
        color: var(--ink);
      }

      button {
        margin-top: 14px;
        background: var(--accent);
        color: #f4fffd;
        border-color: var(--accent);
        cursor: pointer;
        font-weight: 700;
      }

      button:hover {
        background: var(--accent-strong);
      }

      .status {
        margin-top: 14px;
        padding: 12px 14px;
        border-radius: 14px;
        background: rgba(15, 118, 110, 0.08);
        color: var(--accent-strong);
        min-height: 48px;
      }

      .status.warn {
        background: rgba(138, 59, 18, 0.08);
        color: var(--warn);
      }

      .stack {
        display: grid;
        gap: 10px;
      }

      .stack.two {
        grid-template-columns: repeat(2, minmax(0, 1fr));
      }

      .log {
        margin-top: 12px;
        max-height: 420px;
        overflow: auto;
        padding: 0;
        list-style: none;
        display: grid;
        gap: 10px;
      }

      .log li {
        border: 1px solid var(--line);
        border-radius: 14px;
        background: #fff;
        padding: 12px;
      }

      .log strong {
        display: block;
        margin-bottom: 6px;
      }

      .log pre {
        margin: 0;
        white-space: pre-wrap;
        word-break: break-word;
        color: var(--muted);
        font-size: 0.88rem;
      }

      .pill {
        display: inline-block;
        font-size: 0.8rem;
        padding: 6px 10px;
        border-radius: 999px;
        background: rgba(15, 118, 110, 0.08);
        color: var(--accent-strong);
      }

      .map-shell {
        margin-top: 18px;
        display: grid;
        gap: 10px;
      }

      .map-meta {
        display: grid;
        grid-template-columns: minmax(0, 1fr) auto;
        gap: 12px;
        align-items: start;
      }

      .legend {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        justify-content: flex-end;
      }

      .legend-item {
        display: inline-flex;
        align-items: center;
        gap: 8px;
        padding: 7px 10px;
        border-radius: 999px;
        border: 1px solid var(--line);
        background: rgba(255, 255, 255, 0.76);
        font-size: 0.8rem;
        color: var(--muted);
      }

      .legend-swatch {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        display: inline-block;
      }

      .map-grid {
        display: grid;
        grid-template-columns: repeat(20, minmax(0, 1fr));
        gap: 2px;
        background:
          linear-gradient(180deg, rgba(255,255,255,0.5), rgba(30,42,47,0.02)),
          #d1c7ba;
        border-radius: 22px;
        padding: 12px;
        box-shadow: inset 0 1px 0 rgba(255,255,255,0.7);
      }

      .cell {
        aspect-ratio: 1;
        border-radius: 5px;
        position: relative;
        overflow: hidden;
        transform: translateY(0);
        transition: transform 140ms ease, filter 140ms ease;
      }

      .cell:hover {
        transform: translateY(-1px);
        filter: saturate(1.08);
      }

      .cell[data-terrain="water"] {
        background:
          linear-gradient(135deg, rgba(255,255,255,0.28), transparent 55%),
          linear-gradient(180deg, #7bb1df 0%, #3f74b0 100%);
      }

      .cell[data-terrain="plain"] {
        background:
          linear-gradient(145deg, rgba(255,255,255,0.22), transparent 60%),
          linear-gradient(180deg, #9dca82 0%, #6ea85b 100%);
      }

      .cell[data-terrain="forest"] {
        background:
          radial-gradient(circle at 30% 28%, rgba(255,255,255,0.18), transparent 42%),
          linear-gradient(180deg, #668f52 0%, #365c31 100%);
      }

      .cell[data-terrain="hill"] {
        background:
          linear-gradient(145deg, rgba(255,255,255,0.2), transparent 55%),
          linear-gradient(180deg, #aa9a82 0%, #776751 100%);
      }

      .cell.coast::before {
        content: "";
        position: absolute;
        inset: 0;
        border-radius: inherit;
        box-shadow: inset 0 0 0 1px rgba(244, 239, 231, 0.85);
      }

      .cell[data-zone="residential"] {
        box-shadow: inset 0 0 0 2px rgba(244, 253, 248, 0.52);
      }

      .cell[data-zone="commercial"] {
        box-shadow: inset 0 0 0 2px rgba(247, 196, 72, 0.9);
      }

      .cell[data-zone="industrial"] {
        box-shadow: inset 0 0 0 2px rgba(193, 87, 65, 0.9);
      }

      .cell[data-zone="park"] {
        box-shadow: inset 0 0 0 2px rgba(18, 94, 89, 0.9);
      }

      .cell.stadium::after {
        content: "";
        position: absolute;
        inset: 18%;
        border-radius: 40%;
        background:
          linear-gradient(180deg, rgba(255,255,255,0.95), rgba(220,226,229,0.92));
        border: 2px solid #1e2a2f;
        box-shadow: 0 0 0 2px rgba(255,255,255,0.5);
      }

      .cell.stadium::before {
        content: "";
        position: absolute;
        inset: -18%;
        border-radius: 50%;
        background: radial-gradient(circle, rgba(232, 120, 76, 0.28), transparent 62%);
      }
    </style>
  </head>
  <body>
    <main>
      <span class="pill">Gateway Debug Surface</span>
      <h1>PulseCity<br />Map Slice</h1>
      <p>
        Esta pagina dispara la creacion de partida por HTTP y escucha los deltas
        de progreso por WebSocket.
      </p>

      <section class="grid">
        <article class="panel">
          <h2>Crear partida</h2>
          <div class="stack">
            <div>
              <label for="city-name">Nombre de la ciudad</label>
              <input id="city-name" value="Nueva Aurora" />
            </div>
            <div>
              <label for="game-id">Game ID</label>
              <input id="game-id" placeholder="uuid de partida" />
            </div>
            <button id="create-game" type="button">Crear partida</button>
            <div class="stack two">
              <button id="reconnect-socket" type="button">Reconectar socket</button>
              <button id="load-snapshot" type="button">Cargar snapshot</button>
            </div>
            <div id="request-status" class="status">Esperando accion.</div>
          </div>
        </article>

        <article class="panel">
          <h2>WebSocket</h2>
          <div class="stack">
            <div id="socket-status" class="status">Conectando...</div>
            <div id="latest-progress" class="status">Sin deltas todavia.</div>
          </div>
        </article>
      </section>

      <section class="panel" style="margin-top: 18px;">
        <h2>Mapa</h2>
        <div class="map-shell">
          <div class="map-meta">
            <div id="map-summary" class="status">Esperando datos de mapa.</div>
            <div class="legend">
              <span class="legend-item"><span class="legend-swatch" style="background:#4f86c6;"></span>Agua</span>
              <span class="legend-item"><span class="legend-swatch" style="background:#7aa95f;"></span>Llano</span>
              <span class="legend-item"><span class="legend-swatch" style="background:#436a39;"></span>Bosque</span>
              <span class="legend-item"><span class="legend-swatch" style="background:#8c7d69;"></span>Colina</span>
            </div>
          </div>
          <div id="map-grid" class="map-grid"></div>
        </div>
      </section>

      <section class="panel" style="margin-top: 18px;">
        <h2>Eventos recibidos</h2>
        <ul id="event-log" class="log"></ul>
      </section>
    </main>

    <script>
      const requestStatus = document.getElementById("request-status");
      const socketStatus = document.getElementById("socket-status");
      const latestProgress = document.getElementById("latest-progress");
      const cityNameInput = document.getElementById("city-name");
      const gameIDInput = document.getElementById("game-id");
      const eventLog = document.getElementById("event-log");
      const createGameButton = document.getElementById("create-game");
      const reconnectSocketButton = document.getElementById("reconnect-socket");
      const loadSnapshotButton = document.getElementById("load-snapshot");
      const mapSummary = document.getElementById("map-summary");
      const mapGrid = document.getElementById("map-grid");

      let currentMap = null;
      let currentStadium = null;
      let currentGameID = new URLSearchParams(window.location.search).get("game_id") || "";
      let socket = null;

      const socketProtocol = window.location.protocol === "https:" ? "wss:" : "ws:";
      gameIDInput.value = currentGameID;
      connectSocket();

      createGameButton.addEventListener("click", async () => {
        requestStatus.textContent = "Enviando request...";
        requestStatus.classList.remove("warn");

        try {
          const response = await fetch("/api/v1/games", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ city_name: cityNameInput.value.trim() })
          });

          const payload = await response.json();
          if (!response.ok) {
            throw new Error(payload.error || "request failed");
          }

          currentGameID = payload.game_id;
          gameIDInput.value = currentGameID;
          syncURL();
          connectSocket();
          requestStatus.textContent = "Partida creada. game_id: " + payload.game_id;
        } catch (error) {
          requestStatus.textContent = "Error: " + error.message;
          requestStatus.classList.add("warn");
        }
      });

      reconnectSocketButton.addEventListener("click", () => {
        currentGameID = gameIDInput.value.trim();
        syncURL();
        connectSocket();
      });

      loadSnapshotButton.addEventListener("click", async () => {
        currentGameID = gameIDInput.value.trim();
        syncURL();

        if (!currentGameID) {
          requestStatus.textContent = "Ingresar un game_id para cargar snapshot.";
          requestStatus.classList.add("warn");
          return;
        }

        try {
          const response = await fetch("/api/v1/games/" + currentGameID + "/snapshot");
          const payload = await response.json();
          if (!response.ok) {
            throw new Error(payload.error || "snapshot request failed");
          }

          applyMessage(payload);
          renderMap();
          appendEvent(payload);
          requestStatus.textContent = "Snapshot cargado por HTTP.";
          requestStatus.classList.remove("warn");
        } catch (error) {
          requestStatus.textContent = "Error: " + error.message;
          requestStatus.classList.add("warn");
        }
      });

      function appendEvent(message) {
        const item = document.createElement("li");
        const title = document.createElement("strong");
        const body = document.createElement("pre");

        title.textContent = message.subject;
        body.textContent = JSON.stringify(message.state || message.patch || message.payload, null, 2);

        item.prepend(title);
        item.appendChild(body);
        eventLog.prepend(item);
      }

      function applyMessage(message) {
        if (message.type === "map.snapshot") {
          currentGameID = message.state.game_id || currentGameID;
          gameIDInput.value = currentGameID;
          currentMap = message.state.map_data || null;
          currentStadium = message.state.stadium || null;
          latestProgress.textContent =
            message.subject + " -> " + message.state.progress + "% | " + message.state.message;
          return;
        }

        if (message.type === "map.patch") {
          if (message.patch.map_data) {
            currentMap = message.patch.map_data;
          }
          if (message.patch.stadium) {
            currentStadium = message.patch.stadium;
          }

          const progressText =
            (message.patch.progress ?? "?") + "% | " + (message.patch.message || "Sin mensaje");
          latestProgress.textContent = message.subject + " -> " + progressText;
        }
      }

      function connectSocket() {
        if (socket) {
          socket.close();
        }

        const query = currentGameID ? "?game_id=" + encodeURIComponent(currentGameID) : "";
        socket = new WebSocket(socketProtocol + "//" + window.location.host + "/ws" + query);

        socket.addEventListener("open", () => {
          socketStatus.textContent = currentGameID
            ? "WebSocket conectado para " + currentGameID + "."
            : "WebSocket conectado sin game_id.";
          socketStatus.classList.remove("warn");
        });

        socket.addEventListener("close", () => {
          socketStatus.textContent = "WebSocket desconectado.";
          socketStatus.classList.add("warn");
        });

        socket.addEventListener("error", () => {
          socketStatus.textContent = "Error en WebSocket.";
          socketStatus.classList.add("warn");
        });

        socket.addEventListener("message", (event) => {
          const message = JSON.parse(event.data);
          applyMessage(message);
          renderMap();
          appendEvent(message);
        });
      }

      function syncURL() {
        const url = new URL(window.location.href);
        if (currentGameID) {
          url.searchParams.set("game_id", currentGameID);
        } else {
          url.searchParams.delete("game_id");
        }
        window.history.replaceState({}, "", url);
      }

      function renderMap() {
        if (!currentMap) {
          return;
        }

        mapGrid.innerHTML = "";
        mapGrid.style.gridTemplateColumns = "repeat(" + currentMap.width + ", minmax(0, 1fr))";

        for (let y = 0; y < currentMap.height; y += 1) {
          for (let x = 0; x < currentMap.width; x += 1) {
            const cell = currentMap.cells[y][x];
            const tile = document.createElement("div");

            tile.className = "cell";
            tile.dataset.terrain = cell.terrain;
            if (cell.zone) {
              tile.dataset.zone = cell.zone;
            }
            if (isCoastCell(x, y)) {
              tile.classList.add("coast");
            }

            if (currentStadium && currentStadium.x === x && currentStadium.y === y) {
              tile.classList.add("stadium");
            }

            mapGrid.appendChild(tile);
          }
        }

        const stats = summarizeMap();
        mapSummary.textContent =
          "Mapa " +
          currentMap.width +
          "x" +
          currentMap.height +
          " | agua " + stats.water + "% | bosque " + stats.forest + "% | llano " + stats.plain + "%" +
          (currentStadium
            ? " | estadio en (" + currentStadium.x + ", " + currentStadium.y + ")"
            : " | estadio pendiente");
      }

      function isCoastCell(x, y) {
        const cell = currentMap.cells[y][x];
        if (cell.terrain !== "plain" && cell.terrain !== "forest") {
          return false;
        }

        const neighbors = [
          [x - 1, y],
          [x + 1, y],
          [x, y - 1],
          [x, y + 1]
        ];

        return neighbors.some(([nx, ny]) => {
          if (nx < 0 || ny < 0 || nx >= currentMap.width || ny >= currentMap.height) {
            return false;
          }
          return currentMap.cells[ny][nx].terrain === "water";
        });
      }

      function summarizeMap() {
        let water = 0;
        let forest = 0;
        let plain = 0;
        const total = currentMap.width * currentMap.height;

        for (const row of currentMap.cells) {
          for (const cell of row) {
            if (cell.terrain === "water") water += 1;
            if (cell.terrain === "forest") forest += 1;
            if (cell.terrain === "plain") plain += 1;
          }
        }

        return {
          water: Math.round((water / total) * 100),
          forest: Math.round((forest / total) * 100),
          plain: Math.round((plain / total) * 100)
        };
      }
    </script>
  </body>
</html>`
