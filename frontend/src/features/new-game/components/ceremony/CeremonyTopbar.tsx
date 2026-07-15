import { formatPointDifferential, formatSimulatedDate } from "./helpers";
import type { CeremonySharedProps } from "./types";

interface CeremonyTopbarProps {
  data: Pick<
    CeremonySharedProps,
    "currentStage" | "mapState" | "seasonState" | "socketStatus" | "timeState"
  >;
  onOpenTradeCenter: () => void;
  onSetPaused: (paused: boolean) => void;
  onSetSpeed: (speed: 1 | 5 | 20) => void;
}

export function CeremonyTopbar({ data, onOpenTradeCenter, onSetPaused, onSetSpeed }: CeremonyTopbarProps) {
  return (
    <header className="ceremony-builder__topbar ceremony-command-topbar">
      <div className="ceremony-command-topbar__stage">
        <p className="eyebrow">Command center</p>
        <strong>{data.currentStage.label}</strong>
        <button type="button" className="ceremony-command-topbar__workspace" onClick={onOpenTradeCenter}>
          Abrir Trade Center
        </button>
      </div>

      <div className="ceremony-command-topbar__metrics" aria-label="Estado principal de la partida">
        <TopbarMetric label="Fecha" value={formatSimulatedDate(data.timeState.simulated_date)} />
        <TopbarMetric
          label="Record"
          value={`${data.seasonState.wins}-${data.seasonState.losses}`}
          detail={formatPointDifferential(data.seasonState)}
        />
        <TopbarMetric label="Socket" value={data.socketStatus} detail={`${data.mapState.progress}% mapa`} />
      </div>

      <div className="ceremony-command-topbar__time" aria-label="Controles de tiempo">
        <button
          type="button"
          className={data.timeState.paused ? "active" : ""}
          onClick={() => onSetPaused(true)}
        >
          Pausa
        </button>
        {[1, 5, 20].map((speed) => (
          <button
            key={speed}
            type="button"
            className={!data.timeState.paused && data.timeState.speed === speed ? "active" : ""}
            onClick={() => onSetSpeed(speed as 1 | 5 | 20)}
          >
            x{speed}
          </button>
        ))}
      </div>
    </header>
  );
}

function TopbarMetric({ label, value, detail }: { label: string; value: string; detail?: string }) {
  return (
    <div className="ceremony-command-topbar__metric">
      <span>{label}</span>
      <strong>{value}</strong>
      {detail ? <small>{detail}</small> : null}
    </div>
  );
}
