import type { CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import type { MapClientState, RealtimeEvent } from "../../../types";
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

          <div className="ceremony-terrain">
            <TerrainStat label="Agua" value={terrainStats.water} />
            <TerrainStat label="Bosque" value={terrainStats.forest} />
            <TerrainStat label="Llano" value={terrainStats.plain} />
            <TerrainStat label="Colina" value={terrainStats.hill} />
          </div>
        </section>

        <aside className="ceremony-control">
          <section className="ceremony-panel">
            <p className="eyebrow">Estado</p>
            <div className="ceremony-state-grid">
              <StateItem label="Partida" value={formatGameId(props.mapState.game_id || props.gameId)} />
              <StateItem label="Sistema" value={props.status} />
              <StateItem
                label="Owner"
                value={props.ownerIntroResponseLabel ?? "Pendiente de llamada"}
              />
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

function StateItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="ceremony-state-item">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function TerrainStat({ label, value }: { label: string; value: number }) {
  return (
    <div className="ceremony-terrain__item">
      <span>{label}</span>
      <strong>{value}%</strong>
    </div>
  );
}
