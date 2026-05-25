import type { CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import type {
  AgentClientStates,
  CityClientState,
  MapClientState,
  NarrativeEvent,
  RealtimeEvent,
  SeasonClientState,
  SeasonMatchSummary,
  TimeClientState,
} from "../../../types";
import { stageMeta, stageSequence } from "../constants";
import {
  buildCellClassName,
  describeRealtimeEvent,
  formatGameId,
  gridColumns,
  summarizeTerrain,
} from "../helpers";

interface CeremonyPageProps {
  currentStage: {
    label: string;
    title: string;
    description: string;
  };
  events: RealtimeEvent[];
  gameId: string;
  agentStates: AgentClientStates;
  cityState: CityClientState;
  mapState: MapClientState;
  narrativeInbox: NarrativeEvent[];
  ownerIntroResponseLabel: string | null;
  recentResults: SeasonMatchSummary[];
  seasonState: SeasonClientState;
  socketStatus: string;
  status: string;
  timeState: TimeClientState;
  onSetPaused: (paused: boolean) => void;
  onSetSpeed: (speed: 1 | 5 | 20) => void;
}

export function CeremonyPage(props: CeremonyPageProps) {
  const terrainStats = summarizeTerrain(props.mapState.map_data?.cells ?? []);
  const stageIndex = stageSequence.indexOf(
    props.mapState.stage as (typeof stageSequence)[number],
  );
  const completedSteps = stageIndex >= 0 ? stageIndex : -1;
  const showZones =
    props.mapState.stage === "zoning" ||
    props.mapState.stage === "stadium" ||
    props.mapState.stage === "complete";
  const showStadium = props.mapState.stage === "stadium" || props.mapState.stage === "complete";

  return (
    <section
      className="ceremony-builder"
      style={{ "--ceremony-backdrop": `url("${skylineBackdrop}")` } as CSSProperties}
    >
      <div className="ceremony-builder__image" />
      <div className="ceremony-builder__shade" />

      <header className="ceremony-builder__topbar">
        <div>
          <p className="eyebrow">Ceremonia de fundacion</p>
          <strong>{props.currentStage.label}</strong>
        </div>
        <div className="ceremony-live-status">
          <span>{props.socketStatus}</span>
          <strong>{props.mapState.progress}%</strong>
        </div>
      </header>

      <main className="ceremony-builder__main">
        <section className="ceremony-world" aria-label="Mapa generado">
          <div className="ceremony-world__header">
            <div>
              <p className="eyebrow">Nacimiento del mundo</p>
              <h1>{props.currentStage.title}</h1>
              <p>{props.currentStage.description}</p>
            </div>
            <span>{props.mapState.message}</span>
          </div>

          <div className="ceremony-world__frame">
            <div
              className={[
                "map-grid",
                props.cityState.last_match_id ? "map-grid--city-pulse" : "",
              ]
                .filter(Boolean)
                .join(" ")}
              style={gridColumns(props.mapState.map_data?.width ?? 1)}
            >
              {(props.mapState.map_data?.cells ?? []).flatMap((row, y) =>
                row.map((cell, x) => {
                  const classes = buildCellClassName({
                    cell,
                    showZones,
                    showStadium:
                      showStadium &&
                      props.mapState.stadium?.x === x &&
                      props.mapState.stadium?.y === y,
                  });

                  return <div key={`${x}-${y}`} className={classes} />;
                }),
              )}
            </div>
          </div>

          <div className="ceremony-terrain">
            <TerrainStat label="Agua" value={terrainStats.water} />
            <TerrainStat label="Bosque" value={terrainStats.forest} />
            <TerrainStat label="Llano" value={terrainStats.plain} />
            <TerrainStat label="Colina" value={terrainStats.hill} />
          </div>
        </section>

        <aside className="ceremony-control">
          <section className="ceremony-panel time-hud-panel">
            <div className="time-hud__header">
              <div>
                <p className="eyebrow">Tiempo simulado</p>
                <strong>{formatSimulatedDate(props.timeState.simulated_date)}</strong>
              </div>
              <span className={props.timeState.paused ? "time-hud__badge paused" : "time-hud__badge"}>
                {props.timeState.paused ? "Pausa" : `x${props.timeState.speed}`}
              </span>
            </div>

            <div className="time-hud__controls" aria-label="Controles de tiempo">
              <button
                type="button"
                className={props.timeState.paused ? "active" : ""}
                onClick={() => props.onSetPaused(true)}
              >
                Pausa
              </button>
              {[1, 5, 20].map((speed) => (
                <button
                  key={speed}
                  type="button"
                  className={!props.timeState.paused && props.timeState.speed === speed ? "active" : ""}
                  onClick={() => props.onSetSpeed(speed as 1 | 5 | 20)}
                >
                  x{speed}
                </button>
              ))}
            </div>
          </section>

          <section className="ceremony-panel">
            <div className="ceremony-panel__title">
              <p className="eyebrow">Temporada viva</p>
              <strong>{props.seasonState.wins + props.seasonState.losses}/82</strong>
            </div>
            <div className="ceremony-scoreline">
              <strong>{props.seasonState.wins}-{props.seasonState.losses}</strong>
              <span>{formatPointDifferential(props.seasonState)}</span>
            </div>
            <div className="ceremony-state-grid">
              <StateItem label="Partida" value={formatGameId(props.mapState.game_id || props.gameId)} />
              <StateItem
                label="Ciudad"
                value={`${Math.round(props.cityState.fan_sentiment)} ánimo / ${Math.round(props.cityState.stadium_district_land_value)} suelo`}
              />
              <StateItem
                label="Entradas"
                value={`${Math.round(props.cityState.ticket_sales_index)} demanda`}
              />
              <StateItem label="Sistema" value={props.status} />
              <StateItem
                label="Owner"
                value={props.ownerIntroResponseLabel ?? "Pendiente de llamada"}
              />
            </div>
          </section>

          <section className="ceremony-panel">
            <div className="ceremony-panel__title">
              <p className="eyebrow">Resultados recientes</p>
              <strong>{props.recentResults.length}</strong>
            </div>
            <ul className="ceremony-results">
              {props.recentResults.length === 0 ? (
                <li className="empty">Todavia no hay partidos finalizados.</li>
              ) : (
                props.recentResults.map((result) => (
                  <li key={result.match_id}>
                    <span className={result.winner_team_id === "pulsecity" ? "win" : "loss"}>
                      {result.winner_team_id === "pulsecity" ? "W" : "L"}
                    </span>
                    <div>
                      <strong>{formatMatchScore(result)}</strong>
                      <small>{formatSimulatedDate(result.simulated_date)}</small>
                    </div>
                  </li>
                ))
              )}
            </ul>
          </section>

          <section className="ceremony-panel">
            <div className="ceremony-panel__title">
              <p className="eyebrow">Inbox narrativo</p>
              <strong>{props.narrativeInbox.length}</strong>
            </div>
            <ul className="ceremony-inbox">
              {props.narrativeInbox.length === 0 ? (
                <li className="empty">Los reportes post-partido van a aparecer aca.</li>
              ) : (
                props.narrativeInbox.map((event) => (
                  <li key={event.event_id}>
                    <div>
                      <strong>{event.title}</strong>
                      <span>{event.emitter}</span>
                    </div>
                    <p>{event.body}</p>
                  </li>
                ))
              )}
            </ul>
          </section>

          <section className="ceremony-panel">
            <p className="eyebrow">Agentes core</p>
            <div className="ceremony-agents">
              {coreAgentOrder.map((agent) => {
                const state = props.agentStates[agent.id];
                return (
                  <article key={agent.id} className="ceremony-agent">
                    <div className="ceremony-agent__header">
                      <strong>{agent.label}</strong>
                      <span className={`ceremony-agent__mood mood-${state?.mood ?? "idle"}`}>
                        {state?.mood ?? "idle"}
                      </span>
                    </div>
                    <p>{state?.summary || "Esperando el primer resultado."}</p>
                    <div className="ceremony-agent__metrics">
                      {agent.metrics.map((metric) => (
                        <span key={metric.key}>
                          {metric.label}
                          <strong>{formatAgentMetric(state?.state[metric.key])}</strong>
                        </span>
                      ))}
                    </div>
                  </article>
                );
              })}
            </div>
          </section>

          <section className="ceremony-panel">
            <p className="eyebrow">Pipeline</p>
            <ol className="ceremony-pipeline">
              {stageSequence.map((stage, index) => {
                const isActive = props.mapState.stage === stage;
                const isDone = completedSteps >= index;

                return (
                  <li
                    key={stage}
                    className={[
                      "ceremony-pipeline__item",
                      isActive ? "active" : "",
                      isDone ? "done" : "",
                    ]
                      .filter(Boolean)
                      .join(" ")}
                  >
                    <span>0{index + 1}</span>
                    <div>
                      <strong>{stageMeta[stage].label}</strong>
                      <small>{stageMeta[stage].title}</small>
                    </div>
                  </li>
                );
              })}
            </ol>
          </section>

          <section className="ceremony-panel">
            <p className="eyebrow">Eventos recientes</p>
            <ul className="ceremony-events">
              {props.events.length === 0 ? (
                <li>Todavia no llegaron eventos para esta partida.</li>
              ) : (
                props.events.map((event, index) => (
                  <li key={`${event.subject}-${index}`}>
                    <strong>{event.subject}</strong>
                    <span>{describeRealtimeEvent(event)}</span>
                  </li>
                ))
              )}
            </ul>
          </section>
        </aside>
      </main>
    </section>
  );
}

const coreAgentOrder = [
  {
    id: "owner",
    label: "Owner",
    metrics: [
      { key: "sporting_trust", label: "Dep." },
      { key: "patience_remaining", label: "Pac." },
    ],
  },
  {
    id: "head_coach",
    label: "Head Coach",
    metrics: [
      { key: "roster_satisfaction", label: "Roster" },
      { key: "results_pressure", label: "Pres." },
    ],
  },
  {
    id: "cfo",
    label: "CFO",
    metrics: [
      { key: "financial_trust", label: "Fin." },
      { key: "budget_alert", label: "Alerta" },
    ],
  },
  {
    id: "scouting_director",
    label: "Dir. Scouting",
    metrics: [
      { key: "criteria_trust", label: "Crit." },
      { key: "motivation", label: "Mot." },
    ],
  },
  {
    id: "sports_psychologist",
    label: "Sports Psych.",
    metrics: [
      { key: "locker_room_climate", label: "Clima" },
      { key: "emotional_alert", label: "Alerta" },
    ],
  },
] as const;

function formatSimulatedDate(value: string) {
  const date = new Date(`${value}T00:00:00Z`);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat("es", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    timeZone: "UTC",
  }).format(date);
}

function StateItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="ceremony-state-item">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function formatPointDifferential(season: SeasonClientState) {
  const games = season.wins + season.losses;
  if (games === 0) {
    return "DIF 0.0";
  }

  const differential = (season.points_for - season.points_against) / games;
  const sign = differential > 0 ? "+" : "";
  return `DIF ${sign}${differential.toFixed(1)}`;
}

function formatMatchScore(result: SeasonMatchSummary) {
  const ownHome = result.home_team_id === "pulsecity";
  const ownScore = ownHome ? result.home_score : result.away_score;
  const opponentScore = ownHome ? result.away_score : result.home_score;
  const venue = ownHome ? "vs" : "@";
  const opponent = ownHome ? result.away_team_id : result.home_team_id;
  return `${ownScore}-${opponentScore} ${venue} ${opponent}`;
}

function formatAgentMetric(value: number | undefined) {
  if (value === undefined) {
    return "--";
  }

  return value.toFixed(2);
}

function TerrainStat({ label, value }: { label: string; value: number }) {
  return (
    <div className="ceremony-terrain__item">
      <span>{label}</span>
      <strong>{value}%</strong>
    </div>
  );
}
