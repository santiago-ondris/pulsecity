package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/pulsecity/services/gateway/internal/domain"
	natsclient "github.com/pulsecity/services/gateway/internal/nats"
	"github.com/pulsecity/services/gateway/internal/ws"
)

type Dependencies struct {
	Bus *natsclient.Client
	Hub *ws.Hub
}

func RegisterRoutes(mux *http.ServeMux, deps Dependencies) {
	mux.HandleFunc("GET /", debugPage)
	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /ws", deps.Hub.HandleWebSocket)
	mux.HandleFunc("POST /api/v1/games", deps.startGame)
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

func (d Dependencies) startGame(w http.ResponseWriter, r *http.Request) {
	var request domain.StartGameRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&request)
	}

	command := domain.MapGenerationRequest{
		GameID:   uuid.NewString(),
		CityName: request.CityName,
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

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
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
            <button id="create-game" type="button">Crear partida</button>
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
        <h2>Eventos recibidos</h2>
        <ul id="event-log" class="log"></ul>
      </section>
    </main>

    <script>
      const requestStatus = document.getElementById("request-status");
      const socketStatus = document.getElementById("socket-status");
      const latestProgress = document.getElementById("latest-progress");
      const cityNameInput = document.getElementById("city-name");
      const eventLog = document.getElementById("event-log");
      const createGameButton = document.getElementById("create-game");

      const socketProtocol = window.location.protocol === "https:" ? "wss:" : "ws:";
      const socket = new WebSocket(socketProtocol + "//" + window.location.host + "/ws");

      socket.addEventListener("open", () => {
        socketStatus.textContent = "WebSocket conectado.";
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
        latestProgress.textContent =
          message.subject + " -> " + message.payload.progress + "% | " + message.payload.message;
        appendEvent(message);
      });

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

          requestStatus.textContent =
            "Partida creada. game_id: " + payload.game_id;
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
        body.textContent = JSON.stringify(message.payload, null, 2);

        item.prepend(title);
        item.appendChild(body);
        eventLog.prepend(item);
      }
    </script>
  </body>
</html>`
