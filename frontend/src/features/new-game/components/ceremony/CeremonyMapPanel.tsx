import type { CSSProperties } from "react";

import type { CityClientState, MapClientState } from "../../../../types";
import { stageMeta, stageSequence } from "../../constants";
import { buildCellClassName, gridColumns, summarizeTerrain } from "../../helpers";

interface CeremonyMapPanelProps {
  cityState: CityClientState;
  currentStage: {
    title: string;
    description: string;
  };
  mapState: MapClientState;
}

export function CeremonyMapPanel({ cityState, currentStage, mapState }: CeremonyMapPanelProps) {
  const terrainStats = summarizeTerrain(mapState.map_data?.cells ?? []);
  const stageIndex = stageSequence.indexOf(mapState.stage as (typeof stageSequence)[number]);
  const showZones =
    mapState.stage === "zoning" ||
    mapState.stage === "stadium" ||
    mapState.stage === "complete";
  const showStadium = mapState.stage === "stadium" || mapState.stage === "complete";

  return (
    <section className="ceremony-world ceremony-world--command" aria-label="Mapa generado">
      <div className="ceremony-world__header ceremony-world__header--compact">
        <div>
          <p className="eyebrow">Ciudad activa</p>
          <h1>{currentStage.title}</h1>
          <p>{currentStage.description}</p>
        </div>
        <span>{mapState.message}</span>
      </div>

      <div className="ceremony-world__frame ceremony-world__frame--command">
        <div
          className={["map-grid", cityState.last_match_id ? "map-grid--city-pulse" : ""]
            .filter(Boolean)
            .join(" ")}
          style={gridColumns(mapState.map_data?.width ?? 1) as CSSProperties}
        >
          {(mapState.map_data?.cells ?? []).flatMap((row, y) =>
            row.map((cell, x) => {
              const classes = buildCellClassName({
                cell,
                showZones,
                showStadium:
                  showStadium &&
                  mapState.stadium?.x === x &&
                  mapState.stadium?.y === y,
              });

              return <div key={`${x}-${y}`} className={classes} />;
            }),
          )}
        </div>
      </div>

      <div className="ceremony-map-footer">
        <div className="ceremony-stage-strip" aria-label="Pipeline de generacion">
          {stageSequence.map((stage, index) => (
            <span
              key={stage}
              className={[
                "ceremony-stage-strip__item",
                index <= stageIndex ? "done" : "",
                mapState.stage === stage ? "active" : "",
              ]
                .filter(Boolean)
                .join(" ")}
            >
              {stageMeta[stage].label}
            </span>
          ))}
        </div>

        <div className="ceremony-terrain ceremony-terrain--inline">
          <TerrainStat label="Agua" value={terrainStats.water} />
          <TerrainStat label="Bosque" value={terrainStats.forest} />
          <TerrainStat label="Llano" value={terrainStats.plain} />
          <TerrainStat label="Colina" value={terrainStats.hill} />
        </div>
      </div>
    </section>
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
