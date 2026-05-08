import { useEffect, useRef, useState } from "react";

import type {
  MapCell,
  MapClientState,
  MapEvent,
  MapSnapshotEnvelope,
} from "./types";

const gatewayBaseUrl = "http://localhost:8080";
const socketBaseUrl = "ws://localhost:8080/ws";

const initialState: MapClientState = {
  game_id: "",
  stage: "idle",
  progress: 0,
  message: "Esperando nueva partida.",
};

export function App() {
  const [cityName, setCityName] = useState("Nueva Aurora");
  const [gameId, setGameId] = useState("");
  const [status, setStatus] = useState("Sincronizado con gateway.");
  const [socketStatus, setSocketStatus] = useState("Conectando...");
  const [mapState, setMapState] = useState<MapClientState>(initialState);
  const [events, setEvents] = useState<MapEvent[]>([]);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    connectSocket("");

    return () => {
      socketRef.current?.close();
    };
  }, []);

  function connectSocket(nextGameId: string) {
    socketRef.current?.close();

    const query = nextGameId ? `?game_id=${encodeURIComponent(nextGameId)}` : "";
    const socket = new WebSocket(`${socketBaseUrl}${query}`);
    socketRef.current = socket;

    socket.addEventListener("open", () => {
      setSocketStatus(nextGameId ? `Socket activo para ${nextGameId}` : "Socket activo");
    });

    socket.addEventListener("close", () => {
      setSocketStatus("Socket desconectado");
    });

    socket.addEventListener("error", () => {
      setSocketStatus("Error en WebSocket");
    });

    socket.addEventListener("message", (event) => {
      const payload = JSON.parse(event.data) as MapEvent;
      applyEvent(payload);
      setEvents((current) => [payload, ...current].slice(0, 10));
    });
  }

  function applyEvent(payload: MapEvent) {
    if (payload.type === "map.snapshot") {
      setMapState(payload.state);
      setGameId(payload.state.game_id);
      return;
    }

    setMapState((current) => ({
      ...current,
      game_id: payload.game_id,
      stage: payload.patch.stage ?? current.stage,
      progress: payload.patch.progress ?? current.progress,
      message: payload.patch.message ?? current.message,
      map_data: payload.patch.map_data ?? current.map_data,
      stadium: payload.patch.stadium ?? current.stadium,
    }));
  }

  async function createGame() {
    setStatus("Creando partida...");
    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ city_name: cityName }),
      });

      const payload = (await response.json()) as { game_id?: string; error?: string };
      if (!response.ok || !payload.game_id) {
        setStatus(payload.error ?? "No se pudo crear la partida.");
        return;
      }

      setGameId(payload.game_id);
      setStatus(`Partida creada: ${payload.game_id}`);
      connectSocket(payload.game_id);
    } catch (error) {
      setStatus(
        error instanceof Error ? error.message : "Fallo de red al crear la partida.",
      );
    }
  }

  async function loadSnapshot() {
    if (!gameId) {
      setStatus("Ingresar o crear un game_id.");
      return;
    }

    setStatus(`Cargando snapshot ${gameId}...`);
    const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${gameId}/snapshot`);
    const payload = (await response.json()) as MapSnapshotEnvelope | { error?: string };
    if (!response.ok || !("type" in payload)) {
      setStatus(("error" in payload && payload.error) || "No se pudo cargar snapshot.");
      return;
    }

    applyEvent(payload);
    setStatus(`Snapshot cargado para ${gameId}`);
  }

  const terrainStats = summarizeTerrain(mapState.map_data?.cells ?? []);

  return (
    <main className="shell">
      <section className="hero">
        <div className="hero-copy">
          <p className="eyebrow">PulseCity / Frontend Slice</p>
          <h1>El mundo nace del backend y toma forma en el cliente.</h1>
          <p className="lede">
            Este primer frontend real ya consume `snapshot` y `patch`, mantiene estado local y
            muestra una escena 2D mínima lista para migrar después a Three.js.
          </p>
        </div>
        <div className="hero-panel card">
          <label className="field">
            <span>Ciudad</span>
            <input value={cityName} onChange={(event) => setCityName(event.target.value)} />
          </label>
          <label className="field">
            <span>Game ID</span>
            <input value={gameId} onChange={(event) => setGameId(event.target.value)} />
          </label>
          <div className="actions">
            <button onClick={createGame}>Nueva partida</button>
            <button className="secondary" onClick={() => connectSocket(gameId)}>
              Reconectar socket
            </button>
            <button className="secondary" onClick={loadSnapshot}>
              Cargar snapshot
            </button>
          </div>
          <div className="status-grid">
            <div className="status-card">
              <span>Gateway</span>
              <strong>{status}</strong>
            </div>
            <div className="status-card">
              <span>WebSocket</span>
              <strong>{socketStatus}</strong>
            </div>
          </div>
        </div>
      </section>

      <section className="layout">
        <article className="map-card card">
          <div className="panel-head">
            <div>
              <p className="eyebrow">Mapa</p>
              <h2>{mapState.message}</h2>
            </div>
            <div className="progress-badge">{mapState.progress}%</div>
          </div>
          <div className="stats-row">
            <span>Agua {terrainStats.water}%</span>
            <span>Bosque {terrainStats.forest}%</span>
            <span>Llano {terrainStats.plain}%</span>
            <span>Colina {terrainStats.hill}%</span>
          </div>
          <div className="map-grid" style={{ gridTemplateColumns: `repeat(${mapState.map_data?.width ?? 1}, 1fr)` }}>
            {(mapState.map_data?.cells ?? []).flatMap((row, y) =>
              row.map((cell, x) => {
                const classes = [
                  "cell",
                  `terrain-${cell.terrain}`,
                  cell.zone ? `zone-${cell.zone}` : "",
                  mapState.stadium?.x === x && mapState.stadium?.y === y ? "stadium" : "",
                ]
                  .filter(Boolean)
                  .join(" ");

                return <div key={`${x}-${y}`} className={classes} />;
              }),
            )}
          </div>
        </article>

        <aside className="side-column">
          <article className="card detail-card">
            <p className="eyebrow">Estado local</p>
            <h2>{mapState.stage}</h2>
            <dl className="detail-list">
              <div>
                <dt>Game ID</dt>
                <dd>{mapState.game_id || "sin partida"}</dd>
              </div>
              <div>
                <dt>Estadio</dt>
                <dd>
                  {mapState.stadium
                    ? `${mapState.stadium.x}, ${mapState.stadium.y}`
                    : "pendiente"}
                </dd>
              </div>
              <div>
                <dt>Resolución</dt>
                <dd>
                  {mapState.map_data
                    ? `${mapState.map_data.width} x ${mapState.map_data.height}`
                    : "sin mapa"}
                </dd>
              </div>
            </dl>
          </article>

          <article className="card detail-card">
            <p className="eyebrow">Eventos recientes</p>
            <ul className="event-list">
              {events.map((event, index) => (
                <li key={`${event.subject}-${index}`}>
                  <strong>{event.subject}</strong>
                  <span>{event.type === "map.snapshot" ? event.state.message : event.patch.message}</span>
                </li>
              ))}
            </ul>
          </article>
        </aside>
      </section>
    </main>
  );
}

function summarizeTerrain(cells: MapCell[][]) {
  const flat = cells.flat();
  if (flat.length === 0) {
    return { water: 0, forest: 0, plain: 0, hill: 0 };
  }

  const counts = flat.reduce(
    (acc, cell) => {
      acc[cell.terrain] += 1;
      return acc;
    },
    { water: 0, forest: 0, plain: 0, hill: 0 },
  );

  return {
    water: Math.round((counts.water / flat.length) * 100),
    forest: Math.round((counts.forest / flat.length) * 100),
    plain: Math.round((counts.plain / flat.length) * 100),
    hill: Math.round((counts.hill / flat.length) * 100),
  };
}
