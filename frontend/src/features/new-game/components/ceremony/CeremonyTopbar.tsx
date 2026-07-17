import { formatPointDifferential, formatSimulatedDate } from "./helpers";
import type { CeremonySharedProps } from "./types";

interface CeremonyTopbarProps {
  data: Pick<
    CeremonySharedProps,
    "seasonState" | "timeState"
  >;
  abbreviation: string;
  activeAlerts: number;
  cityName: string;
  franchiseName: string;
  mapProgress: number;
  mode: "generation" | "kickoff" | "active";
  onSetPaused: (paused: boolean) => void;
  onSetSpeed: (speed: 1 | 5 | 20) => void;
}

export function CeremonyTopbar({
  data,
  abbreviation,
  activeAlerts,
  cityName,
  franchiseName,
  mapProgress,
  mode,
  onSetPaused,
  onSetSpeed,
}: CeremonyTopbarProps) {
  return (
    <header className="ceremony-builder__topbar ceremony-command-topbar">
      <div className="ceremony-command-topbar__identity">
        <span>{abbreviation}</span>
        <div>
          <p className="eyebrow">Command Center · {cityName}</p>
          <strong>{franchiseName}</strong>
        </div>
      </div>

      <div className="ceremony-command-topbar__metrics" aria-label="Estado principal de la partida">
        <TopbarMetric label="Fecha" value={formatSimulatedDate(data.timeState.simulated_date)} />
        <TopbarMetric
          label="Record"
          value={`${data.seasonState.wins}-${data.seasonState.losses}`}
          detail={formatPointDifferential(data.seasonState)}
        />
        <TopbarMetric
          label="Alertas"
          value={activeAlerts > 0 ? `${activeAlerts} pendientes` : "Sin urgencias"}
          detail={activeAlerts > 0 ? "Revisar inbox" : "Operación estable"}
        />
      </div>

      {mode === "active" ? (
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
      ) : mode === "kickoff" ? (
        <div className="ceremony-command-topbar__waiting" aria-label="Temporada en espera">
          <span>Tiempo pausado</span>
          <strong>Temporada en espera</strong>
          <small>Usá “Comenzar temporada” cuando estés listo</small>
        </div>
      ) : (
        <div className="ceremony-command-topbar__waiting" aria-label="Generación de ciudad">
          <span>Fundación en curso</span>
          <strong>Construyendo la ciudad</strong>
          <small>{mapProgress}% completado</small>
        </div>
      )}
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
