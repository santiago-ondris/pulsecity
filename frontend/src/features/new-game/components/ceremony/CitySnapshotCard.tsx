import type { CSSProperties } from "react";

import type { CityClientState, MapClientState } from "../../../../types";
import { buildCellClassName, summarizeTerrain } from "../../helpers";

interface CitySnapshotCardProps {
  cityName: string;
  cityState: CityClientState;
  mapState: MapClientState;
}

export function CitySnapshotCard({ cityName, cityState, mapState }: CitySnapshotCardProps) {
  const cells = mapState.map_data?.cells ?? [];
  const width = Math.max(mapState.map_data?.width ?? 1, 1);
  const height = Math.max(mapState.map_data?.height ?? cells.length ?? 1, 1);
  const terrain = summarizeTerrain(cells);
  const showZones = mapState.stage === "zoning" || mapState.stage === "stadium" || mapState.stage === "complete";
  const showStadium = mapState.stage === "stadium" || mapState.stage === "complete";

  return (
    <section className="city-snapshot" aria-label={`Snapshot de ${cityName}`}>
      <div className="city-snapshot__header">
        <div>
          <p className="eyebrow">Vista ciudad</p>
          <h2>{cityName}</h2>
        </div>
        <span>{Math.round(cityState.fan_sentiment)} ánimo</span>
      </div>

      <div
        className="city-snapshot__viewport"
        style={{ "--city-map-ratio": `${width} / ${height}` } as CSSProperties}
      >
        <div
          className={[
            "city-snapshot__grid",
            cityState.last_match_id ? "map-grid--city-pulse" : "",
          ].filter(Boolean).join(" ")}
          style={{
            gridTemplateColumns: `repeat(${width}, minmax(0, 1fr))`,
            gridTemplateRows: `repeat(${height}, minmax(0, 1fr))`,
          }}
        >
          {cells.flatMap((row, y) => row.map((cell, x) => (
            <div
              key={`${x}-${y}`}
              className={buildCellClassName({
                cell,
                showZones,
                showStadium: showStadium && mapState.stadium?.x === x && mapState.stadium?.y === y,
              })}
            />
          )))}
        </div>
      </div>

      <div className="city-snapshot__footer">
        <span>Agua <strong>{terrain.water}%</strong></span>
        <span>Verde <strong>{terrain.forest}%</strong></span>
        <span>Distrito estadio <strong>{Math.round(cityState.stadium_district_land_value)}</strong></span>
      </div>
    </section>
  );
}
