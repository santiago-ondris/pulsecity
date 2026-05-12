import type { MapClientState, RealtimeEvent } from "../../../types";
import { stageMeta, stageSequence } from "../constants";
import {
  buildCellClassName,
  describeRealtimeEvent,
  formatGameId,
  gridColumns,
  summarizeTerrain,
} from "../helpers";
import { Metric, StatusBadge } from "./common";

interface CeremonyPageProps {
  currentStage: {
    label: string;
    title: string;
    description: string;
  };
  events: RealtimeEvent[];
  gameId: string;
  mapState: MapClientState;
  ownerIntroResponseLabel: string | null;
  socketStatus: string;
  status: string;
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
    <section className="screen ceremony-screen">
      <div className="ceremony-topbar">
        <div>
          <p className="eyebrow">Ceremonia de fundación</p>
          <h1>{props.currentStage.title}</h1>
          <p className="copy">{props.currentStage.description}</p>
        </div>
        <div className="ceremony-badges">
          <StatusBadge label={props.socketStatus} tone="primary" />
          <StatusBadge label={props.status} tone="info" />
        </div>
      </div>

      <div className="ceremony-grid">
        <article className="ceremony-map-card">
          <div className="ceremony-stats">
            <Metric label="Partida" value={formatGameId(props.mapState.game_id || props.gameId)} />
            <Metric label="Etapa" value={props.currentStage.label} />
            <Metric label="Progreso" value={`${props.mapState.progress}%`} />
            <Metric
              label="Dirección inicial"
              value={props.ownerIntroResponseLabel ?? "Pendiente de llamada del Owner"}
            />
          </div>

          <div className="world-frame">
            <div className="world-header">
              <StatusBadge label={props.mapState.message} tone="info" />
            </div>
            <div className="map-grid" style={gridColumns(props.mapState.map_data?.width ?? 1)}>
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

          <div className="terrain-band">
            <Metric label="Agua" value={`${terrainStats.water}%`} />
            <Metric label="Bosque" value={`${terrainStats.forest}%`} />
            <Metric label="Llano" value={`${terrainStats.plain}%`} />
            <Metric label="Colina" value={`${terrainStats.hill}%`} />
          </div>
        </article>

        <aside className="ceremony-sidebar">
          <article className="step-card">
            <div className="panel-header">
              <p className="eyebrow">Pipeline</p>
              <h2>Orden del backend</h2>
            </div>
            <ol className="timeline">
              {stageSequence.map((stage, index) => {
                const isActive = props.mapState.stage === stage;
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
                        <strong>{stageMeta[stage].label}</strong>
                        <p>{stageMeta[stage].title}</p>
                      </div>
                    </li>
                  );
              })}
            </ol>
          </article>

          <article className="step-card">
            <div className="panel-header">
              <p className="eyebrow">Eventos</p>
              <h2>Traza reciente</h2>
            </div>
            <ul className="event-list">
              {props.events.length === 0 ? (
                <li className="event-empty">Todavia no llegaron eventos para esta partida.</li>
              ) : (
                props.events.map((event, index) => (
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
  );
}
