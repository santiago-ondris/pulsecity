import type { CityClientState, MapClientState } from "../../../../types";
import { CitySnapshotCard } from "./CitySnapshotCard";

interface WorldGenerationPanelProps {
  cityName: string;
  cityState: CityClientState;
  currentStage: {
    label: string;
    title: string;
    description: string;
  };
  mapState: MapClientState;
}

export function WorldGenerationPanel({
  cityName,
  cityState,
  currentStage,
  mapState,
}: WorldGenerationPanelProps) {
  return (
    <main className="world-generation">
      <section className="world-generation__status">
        <p className="eyebrow">Fundación en curso · {currentStage.label}</p>
        <h1>{currentStage.title}</h1>
        <p>{currentStage.description}</p>
        <div className="world-generation__progress" aria-label={`${mapState.progress}% completado`}>
          <span style={{ width: `${Math.max(0, Math.min(mapState.progress, 100))}%` }} />
        </div>
        <div className="world-generation__message">
          <strong>{mapState.progress}%</strong>
          <span>{mapState.message}</span>
        </div>
      </section>
      <CitySnapshotCard cityName={cityName} cityState={cityState} mapState={mapState} />
    </main>
  );
}
