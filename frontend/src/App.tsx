import { useEffect, useRef, useState, type CSSProperties } from "react";

import type {
  GameSetup,
  MapCell,
  MapClientState,
  MapSnapshotEnvelope,
  NarrativeEvent,
  RealtimeEvent,
} from "./types";

const gatewayBaseUrl = "http://localhost:8080";
const socketBaseUrl = "ws://localhost:8080/ws";

const initialState: MapClientState = {
  game_id: "",
  stage: "idle",
  progress: 0,
  message: "Esperando la orden de fundacion.",
};

const stageSequence = ["terrain", "zoning", "stadium", "complete"] as const;

const stageMeta: Record<
  string,
  {
    label: string;
    title: string;
    description: string;
  }
> = {
  idle: {
    label: "En espera",
    title: "La ciudad todavia no fue fundada.",
    description: "Defini la identidad de la franquicia y dispará la ceremonia desde el frontend.",
  },
  terrain: {
    label: "Terreno",
    title: "Primero aparece la geografia.",
    description: "Costa, relieve y base territorial emergen antes de cualquier lectura urbana.",
  },
  zoning: {
    label: "Zonificacion",
    title: "Los distritos toman forma.",
    description: "La ciudad empieza a especializar sus barrios y a declarar vocaciones de suelo.",
  },
  stadium: {
    label: "Estadio",
    title: "La franquicia encuentra su centro.",
    description: "El estadio fija el punto emocional y urbano alrededor del cual se organiza el mapa.",
  },
  complete: {
    label: "Completo",
    title: "El mundo quedo listo.",
    description: "La base de la partida ya existe y el cliente solo tiene que seguir sus deltas.",
  },
};

const initialScenarios = [
  {
    id: "rebuild",
    label: "Reconstruccion",
    roster: "Joven, rating bajo, alto potencial.",
    pressure: "Dueño paciente, presión mediática baja.",
    city: "Ciudad pequeña, economía modesta y fanbase leal.",
  },
  {
    id: "contention",
    label: "Ventana de contencion",
    roster: "2-3 estrellas, veteranos y contratos grandes.",
    pressure: "Dueño exigente, playoffs este año.",
    city: "Ciudad desarrollada, estadio lleno y fanbase caliente.",
  },
  {
    id: "decline",
    label: "Historica en declive",
    roster: "Viejas glorias mezcladas con jovenes sin rumbo.",
    pressure: "Dueño nostalgico, comparación constante con el pasado.",
    city: "Ciudad grande, medios encima y frustración alta.",
  },
  {
    id: "expansion",
    label: "Expansion pura",
    roster: "Draft de expansion, sin compromisos heredados.",
    pressure: "Dueño visionario, horizonte largo y paciencia total.",
    city: "Ciudad virgen, todo por construir desde cero.",
  },
] as const;

const colorPresets = [
  { label: "Signal Green", value: "#00C896" },
  { label: "Steel Blue", value: "#7B8CDE" },
  { label: "Arena Gold", value: "#FFAA00" },
  { label: "Burnt Orange", value: "#FF6B2B" },
  { label: "Infra Red", value: "#E05555" },
  { label: "Night Glass", value: "#1A1A1E" },
] as const;

type ScenarioId = (typeof initialScenarios)[number]["id"];
const cityManagementModes = [
  {
    id: "owner_influence",
    label: "Dueño con influencia",
    description:
      "Sos el GM de la franquicia. Tu poder sobre la ciudad es indirecto: lobby, financiamiento de proyectos y presión política.",
    impact:
      "El alcalde tiene agenda propia y la ciudad reacciona como un organismo independiente.",
  },
  {
    id: "dual_figure",
    label: "Figura dual",
    description:
      "Controlás tanto la franquicia como la ciudad directamente. Dos sombreros, control total.",
    impact:
      "La ciudad queda bajo tu manejo directo y el basket convive con ese control urbano.",
  },
] as const;
type CityManagementModeId = (typeof cityManagementModes)[number]["id"];

export function App() {
  const [cityName, setCityName] = useState("Nueva Aurora");
  const [franchiseName, setFranchiseName] = useState("Lighthouses");
  const [abbreviation, setAbbreviation] = useState("NAR");
  const [primaryColor, setPrimaryColor] = useState("#00C896");
  const [secondaryColor, setSecondaryColor] = useState("#7B8CDE");
  const [accentColor, setAccentColor] = useState("#FFAA00");
  const [selectedScenario, setSelectedScenario] = useState<ScenarioId>("expansion");
  const [cityManagementMode, setCityManagementMode] = useState<CityManagementModeId>("owner_influence");
  const [gameId, setGameId] = useState("");
  const [status, setStatus] = useState("Esperando la orden de fundacion.");
  const [socketStatus, setSocketStatus] = useState("Conectando...");
  const [mapState, setMapState] = useState<MapClientState>(initialState);
  const [events, setEvents] = useState<RealtimeEvent[]>([]);
  const [activeNarrativeEvent, setActiveNarrativeEvent] = useState<NarrativeEvent | null>(null);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    connectSocket("");

    return () => {
      socketRef.current?.close();
    };
  }, []);

  useEffect(() => {
    if (!gameId) {
      return;
    }

    void loadGameSetup(gameId);
  }, [gameId]);

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
      const payload = JSON.parse(event.data) as RealtimeEvent;
      applyEvent(payload);
      setEvents((current) => [payload, ...current].slice(0, 12));
    });
  }

  function applyEvent(payload: RealtimeEvent) {
    if (payload.type === "narrative.event") {
      setActiveNarrativeEvent(payload);
      return;
    }

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
    setStatus(`Fundando ${cityName} ${franchiseName}...`);
    setEvents([]);

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          city_name: cityName,
          franchise_name: franchiseName,
          abbreviation,
          primary_color: primaryColor,
          secondary_color: secondaryColor,
          accent_color: accentColor,
          initial_scenario: selectedScenario,
          city_management_mode: cityManagementMode,
        }),
      });

      const payload = (await response.json()) as { game_id?: string; error?: string };
      if (!response.ok || !payload.game_id) {
        setStatus(payload.error ?? "No se pudo crear la partida.");
        return;
      }

      setGameId(payload.game_id);
      setStatus(`Partida creada para ${cityName}. Ceremonia en curso.`);
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

    setStatus(`Rehidratando partida ${gameId}...`);

    const setupLoaded = await loadGameSetup(gameId);
    const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${gameId}/snapshot`);
    const payload = (await response.json()) as MapSnapshotEnvelope | { error?: string };
    if (!response.ok || !("type" in payload)) {
      setStatus(("error" in payload && payload.error) || "No se pudo cargar snapshot.");
      return;
    }

    applyEvent(payload);
    if (setupLoaded) {
      setStatus(`Partida rehidratada para ${gameId}`);
      return;
    }
    setStatus(`Snapshot cargado para ${gameId}`);
  }

  async function loadGameSetup(nextGameId: string) {
    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${nextGameId}`);
      const payload = (await response.json()) as GameSetup | { error?: string };
      if (!response.ok || !("game_id" in payload)) {
        return false;
      }

      applyGameSetup(payload);
      return true;
    } catch {
      return false;
    }
  }

  function applyGameSetup(setup: GameSetup) {
    setCityName(setup.city_name);
    setFranchiseName(setup.franchise_name);
    setAbbreviation(setup.abbreviation);
    setPrimaryColor(setup.primary_color);
    setSecondaryColor(setup.secondary_color);
    setAccentColor(setup.accent_color);
    if (isScenarioId(setup.initial_scenario)) {
      setSelectedScenario(setup.initial_scenario);
    }
    if (isCityManagementModeId(setup.city_management_mode)) {
      setCityManagementMode(setup.city_management_mode);
    }
    if (setup.owner_intro_event) {
      setActiveNarrativeEvent(setup.owner_intro_event);
    }
  }

  const currentStage = stageMeta[mapState.stage] ?? stageMeta.idle;
  const terrainStats = summarizeTerrain(mapState.map_data?.cells ?? []);
  const stageIndex = stageSequence.indexOf(mapState.stage as (typeof stageSequence)[number]);
  const completedSteps = stageIndex >= 0 ? stageIndex : -1;
  const showZones =
    mapState.stage === "zoning" || mapState.stage === "stadium" || mapState.stage === "complete";
  const showStadium = mapState.stage === "stadium" || mapState.stage === "complete";
  const scenario = initialScenarios.find((item) => item.id === selectedScenario) ?? initialScenarios[0];
  const previewStyle = {
    "--franchise-primary": primaryColor,
    "--franchise-secondary": secondaryColor,
    "--franchise-accent": accentColor,
  } as CSSProperties;

  return (
    <main className="shell">
      <header className="topbar">
        <div>
          <p className="eyebrow">PulseCity / Nueva Partida</p>
          <h1>Fundá la franquicia antes de fundar la ciudad.</h1>
        </div>
        <StatusBadge label={socketStatus} tone="primary" />
      </header>

      <section className="setup-layout">
        <article className="panel">
          <div className="panel-header">
            <p className="eyebrow">01 / Identidad visual</p>
            <h2>Definí nombre, sigla y paleta.</h2>
          </div>

          <div className="field-grid">
            <label className="field">
              <span>Ciudad</span>
              <input value={cityName} onChange={(event) => setCityName(event.target.value)} />
            </label>

            <label className="field">
              <span>Franquicia</span>
              <input
                value={franchiseName}
                onChange={(event) => setFranchiseName(event.target.value)}
              />
            </label>

            <label className="field field-short">
              <span>Abreviatura</span>
              <input
                value={abbreviation}
                maxLength={3}
                onChange={(event) => setAbbreviation(event.target.value.toUpperCase())}
              />
            </label>
          </div>

          <div className="panel-section">
            <div className="section-title-row">
              <div>
                <p className="eyebrow">Paleta</p>
                <h3>Colores fundacionales</h3>
              </div>
              <p className="microcopy">
                La identidad ahora tambien se persiste en backend y puede rehidratarse por
                `game_id`.
              </p>
            </div>

            <div className="color-stack">
              <ColorField
                label="Primario"
                value={primaryColor}
                onChange={setPrimaryColor}
                presets={colorPresets}
              />
              <ColorField
                label="Secundario"
                value={secondaryColor}
                onChange={setSecondaryColor}
                presets={colorPresets}
              />
              <ColorField
                label="Acento"
                value={accentColor}
                onChange={setAccentColor}
                presets={colorPresets}
              />
            </div>
          </div>
        </article>

        <article className="panel">
          <div className="panel-header">
            <p className="eyebrow">02 / Estado inicial</p>
            <h2>Elegí desde dónde empieza la historia.</h2>
          </div>

          <div className="scenario-grid">
            {initialScenarios.map((item) => {
              const active = item.id === selectedScenario;

              return (
                <button
                  key={item.id}
                  type="button"
                  className={active ? "scenario-card active" : "scenario-card"}
                  onClick={() => setSelectedScenario(item.id)}
                >
                  <span className="scenario-label">{item.label}</span>
                  <strong>{item.roster}</strong>
                  <p>{item.pressure}</p>
                  <p>{item.city}</p>
                </button>
              );
            })}
          </div>
        </article>

        <article className="panel">
          <div className="panel-header">
            <p className="eyebrow">03 / Gestión de ciudad</p>
            <h2>Definí cómo gobernás la ciudad.</h2>
          </div>

          <div className="management-grid">
            {cityManagementModes.map((mode) => {
              const active = mode.id === cityManagementMode;

              return (
                <button
                  key={mode.id}
                  type="button"
                  className={active ? "management-card active" : "management-card"}
                  onClick={() => setCityManagementMode(mode.id)}
                >
                  <span className="scenario-label">{mode.label}</span>
                  <strong>{mode.description}</strong>
                  <p>{mode.impact}</p>
                </button>
              );
            })}
          </div>
        </article>

        <aside className="preview-rail">
          <article className="panel preview-panel" style={previewStyle}>
            <div className="panel-header">
              <p className="eyebrow">Preview</p>
              <h2>Identidad de franquicia</h2>
            </div>

            <div className="franchise-crest">
              <span className="crest-mark">{safeAbbreviation(abbreviation)}</span>
            </div>

            <div className="franchise-headline">
              <p>{cityName}</p>
              <h3>{franchiseName}</h3>
            </div>

            <div className="swatch-row">
              <Swatch label="Primario" value={primaryColor} />
              <Swatch label="Secundario" value={secondaryColor} />
              <Swatch label="Acento" value={accentColor} />
            </div>

            <div className="preview-meta">
              <div>
                <span>Escenario</span>
                <strong>{scenario.label}</strong>
              </div>
              <div>
                <span>Gestión de ciudad</span>
                <strong>{managementModeLabel(cityManagementMode)}</strong>
              </div>
              <div>
                <span>Lectura inicial</span>
                <strong>{scenario.pressure}</strong>
              </div>
            </div>
          </article>

          <article className="panel launch-panel">
            <div className="panel-header">
              <p className="eyebrow">04 / Confirmacion</p>
              <h2>Cuando confirmás, nace el mundo.</h2>
            </div>

            <p className="copy">
              Esta pantalla ya modela la identidad y el punto de partida narrativo. La ceremonia
              sigue usando el pipeline real de `map-service`.
            </p>

            <div className="action-stack">
              <button type="button" className="primary-action" onClick={createGame}>
                Confirmar franquicia y generar mapa
              </button>
              <button type="button" className="secondary-action" onClick={() => connectSocket(gameId)}>
                Reabrir socket
              </button>
              <button type="button" className="secondary-action" onClick={loadSnapshot}>
                Rehidratar snapshot
              </button>
            </div>

            <div className="status-list">
              <StatusRow label="Gateway" value={status} />
              <StatusRow label="Partida activa" value={formatGameId(mapState.game_id || gameId)} />
              <StatusRow label="Fase actual" value={currentStage.label} />
            </div>
          </article>
        </aside>
      </section>

      <section className="ceremony-panel panel">
        <div className="ceremony-header">
          <div>
            <p className="eyebrow">Ceremonia de generacion</p>
            <h2>{currentStage.title}</h2>
            <p className="copy">{currentStage.description}</p>
          </div>

          <div className="ceremony-stats">
            <Metric label="Progreso" value={`${mapState.progress}%`} />
            <Metric
              label="Resolucion"
              value={
                mapState.map_data
                  ? `${mapState.map_data.width} x ${mapState.map_data.height}`
                  : "Sin mapa"
              }
            />
            <Metric
              label="Estadio"
              value={mapState.stadium ? `${mapState.stadium.x}, ${mapState.stadium.y}` : "Pendiente"}
            />
          </div>
        </div>

        <div className="ceremony-body">
          <article className="world-card">
            <div className="world-header">
              <StatusBadge label={mapState.message} tone="info" />
            </div>

            <div className={showStadium ? "world-frame stage-live has-stadium" : "world-frame stage-live"}>
              <div className="map-grid" style={gridColumns(mapState.map_data?.width ?? 1)}>
                {(mapState.map_data?.cells ?? []).flatMap((row, y) =>
                  row.map((cell, x) => {
                    const classes = buildCellClassName({
                      cell,
                      showZones,
                      showStadium:
                        showStadium && mapState.stadium?.x === x && mapState.stadium?.y === y,
                    });

                    return <div key={`${x}-${y}`} className={classes} />;
                  }),
                )}
              </div>
            </div>

            <div className="terrain-band">
              <Metric label="Agua" value={`${terrainStats.water}%`} />
              <Metric label="Bosque" value={`${terrainStats.forest}%`} />
              <Metric label="Llano" value={`${terrainStats.plain}%`} />
              <Metric label="Colina" value={`${terrainStats.hill}%`} />
            </div>
          </article>

          <aside className="ceremony-sidebar">
            <article className="panel panel-nested">
              <div className="panel-header">
                <p className="eyebrow">Pipeline</p>
                <h2>Etapas del backend</h2>
              </div>

              <ol className="timeline">
                {stageSequence.map((stage, index) => {
                  const meta = stageMeta[stage];
                  const isActive = mapState.stage === stage;
                  const isDone = completedSteps >= index;

                  return (
                    <li
                      key={stage}
                      className={[
                        "timeline-item",
                        isActive ? "active" : "",
                        isDone ? "done" : "",
                      ]
                        .filter(Boolean)
                        .join(" ")}
                    >
                      <span className="timeline-index">0{index + 1}</span>
                      <div>
                        <strong>{meta.label}</strong>
                        <p>{meta.title}</p>
                      </div>
                    </li>
                  );
                })}
              </ol>
            </article>

            <article className="panel panel-nested">
              <div className="panel-header">
                <p className="eyebrow">Eventos</p>
                <h2>Traza reciente</h2>
              </div>

              <ul className="event-list">
                {events.length === 0 ? (
                  <li className="event-empty">Todavia no llegaron eventos para esta partida.</li>
                ) : (
                  events.map((event, index) => (
                    <li key={`${event.subject}-${index}`}>
                      <strong>{event.subject}</strong>
                      <span>{describeRealtimeEvent(event)}</span>
                    </li>
                  ))
                )}
              </ul>
            </article>
          </aside>
        </div>
      </section>

      {activeNarrativeEvent ? (
        <div className="narrative-overlay" role="dialog" aria-modal="true">
          <article className="narrative-modal panel">
            <div className="panel-header">
              <p className="eyebrow">Evento obligatorio</p>
              <h2>{activeNarrativeEvent.title}</h2>
            </div>

            <StatusBadge label="Owner" tone="urgent" />
            <p className="narrative-body">{activeNarrativeEvent.body}</p>

            <div className="narrative-actions">
              {(activeNarrativeEvent.choices ?? [{ id: "continue", label: "Entendido" }]).map(
                (choice) => (
                  <button
                    key={choice.id}
                    type="button"
                    className="primary-action"
                    onClick={() => setActiveNarrativeEvent(null)}
                  >
                    {choice.label}
                  </button>
                ),
              )}
            </div>
          </article>
        </div>
      ) : null}
    </main>
  );
}

function ColorField({
  label,
  value,
  onChange,
  presets,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  presets: readonly { label: string; value: string }[];
}) {
  return (
    <div className="color-field">
      <div className="color-row">
        <label className="field field-inline">
          <span>{label}</span>
          <div className="color-input-row">
            <input
              className="color-picker"
              type="color"
              value={value}
              onChange={(event) => onChange(event.target.value)}
            />
            <input value={value} onChange={(event) => onChange(event.target.value)} />
          </div>
        </label>
      </div>

      <div className="preset-row">
        {presets.map((preset) => (
          <button
            key={`${label}-${preset.value}`}
            type="button"
            className={value.toLowerCase() === preset.value.toLowerCase() ? "preset active" : "preset"}
            onClick={() => onChange(preset.value)}
            aria-label={`${label} ${preset.label}`}
          >
            <span style={{ background: preset.value }} />
          </button>
        ))}
      </div>
    </div>
  );
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="metric">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function StatusBadge({
  label,
  tone,
}: {
  label: string;
  tone: "primary" | "info" | "urgent";
}) {
  return (
    <span
      className={
        tone === "primary"
          ? "badge badge-primary"
          : tone === "urgent"
            ? "badge badge-urgent"
            : "badge badge-info"
      }
    >
      <span className="badge-dot" />
      {label}
    </span>
  );
}

function StatusRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="status-row">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function Swatch({ label, value }: { label: string; value: string }) {
  return (
    <div className="swatch">
      <span className="swatch-chip" style={{ background: value }} />
      <div>
        <span>{label}</span>
        <strong>{value}</strong>
      </div>
    </div>
  );
}

function buildCellClassName({
  cell,
  showZones,
  showStadium,
}: {
  cell: MapCell;
  showZones: boolean;
  showStadium: boolean;
}) {
  return [
    "cell",
    `terrain-${cell.terrain}`,
    showZones && cell.zone ? `zone-${cell.zone}` : "",
    showStadium ? "stadium" : "",
  ]
    .filter(Boolean)
    .join(" ");
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

function formatGameId(value: string) {
  if (!value) {
    return "sin partida";
  }

  if (value.length <= 18) {
    return value;
  }

  return `${value.slice(0, 8)}...${value.slice(-6)}`;
}

function safeAbbreviation(value: string) {
  const trimmed = value.trim().toUpperCase();
  return trimmed.length > 0 ? trimmed : "NEW";
}

function gridColumns(width: number): CSSProperties {
  return {
    gridTemplateColumns: `repeat(${width}, 1fr)`,
  };
}

function isScenarioId(value: string): value is ScenarioId {
  return initialScenarios.some((scenario) => scenario.id === value);
}

function isCityManagementModeId(value: string): value is CityManagementModeId {
  return cityManagementModes.some((mode) => mode.id === value);
}

function managementModeLabel(value: CityManagementModeId) {
  return cityManagementModes.find((mode) => mode.id === value)?.label ?? "Dueño con influencia";
}

function describeRealtimeEvent(event: RealtimeEvent) {
  if (event.type === "narrative.event") {
    return event.title;
  }

  if (event.type === "map.snapshot") {
    return event.state.message;
  }

  return event.patch.message ?? "Patch recibido";
}
