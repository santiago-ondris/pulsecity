import type { CeremonySharedProps } from "./types";
import { CitySnapshotCard } from "./CitySnapshotCard";
import { formatMatchScore, formatSimulatedDate } from "./helpers";

interface CommandCenterOverviewProps {
  cityName: string;
  data: CeremonySharedProps;
  franchiseName: string;
  rosterCount: number;
  showStartPrompt: boolean;
  startingSeason: boolean;
  onOpenMedicalCenter: () => void;
  onOpenStaff: () => void;
  onOpenTradeCenter: () => void;
  onStartSeason: () => void;
}

export function CommandCenterOverview({
  cityName,
  data,
  franchiseName,
  rosterCount,
  showStartPrompt,
  startingSeason,
  onOpenMedicalCenter,
  onOpenStaff,
  onOpenTradeCenter,
  onStartSeason,
}: CommandCenterOverviewProps) {
  const gamesPlayed = data.seasonState.wins + data.seasonState.losses;
  const injuredPlayers = Object.values(data.rosterStates)
    .filter((player) => player.availability === "injured").length;

  return (
    <section className="command-overview" aria-labelledby="command-overview-title">
      <header className="command-overview__header">
        <div>
          <p className="eyebrow">Pulso de franquicia</p>
          <h1 id="command-overview-title">{franchiseName}</h1>
          <p>{gamesPlayed === 0 ? "La temporada está lista para comenzar." : `Temporada en marcha · ${gamesPlayed}/82 partidos`}</p>
        </div>
        {showStartPrompt ? (
          <button type="button" disabled={startingSeason} onClick={onStartSeason}>
            {startingSeason ? "Iniciando..." : "Comenzar temporada"}
          </button>
        ) : null}
      </header>

      <div className="command-overview__metrics" aria-label="Métricas principales">
        <OverviewMetric label="Récord" value={`${data.seasonState.wins}-${data.seasonState.losses}`} />
        <OverviewMetric
          label="Cap space"
          value={data.financeState.cap_base > 0 ? formatMoney(data.financeState.cap_space) : "—"}
        />
        <OverviewMetric label="Alertas" value={`${data.narrativeInbox.length}`} />
        <OverviewMetric label="Roster" value={rosterCount > 0 ? `${rosterCount}` : "—"} />
      </div>

      <div className="command-overview__body">
        <CitySnapshotCard cityName={cityName} cityState={data.cityState} mapState={data.mapState} />

        <section className="command-operations" aria-labelledby="command-operations-title">
          <div>
            <p className="eyebrow">Operaciones</p>
            <h2 id="command-operations-title">¿Qué necesita tu atención?</h2>
          </div>
          <button type="button" onClick={onOpenTradeCenter}>
            <span>Basketball Ops</span>
            <strong>Trade Center</strong>
            <small>Abrir o seguir negociaciones</small>
          </button>
          <button type="button" onClick={onOpenMedicalCenter}>
            <span>Performance & Health</span>
            <strong>Centro Médico</strong>
            <small>{injuredPlayers > 0 ? `${injuredPlayers} casos activos` : "Plantel sin bajas registradas"}</small>
          </button>
          <button type="button" onClick={onOpenStaff}>
            <span>Consulta directa</span>
            <strong>Hablar con el staff</strong>
            <small>Buscar un agente por nombre o área</small>
          </button>
        </section>
      </div>

      <section className="command-recent" aria-labelledby="command-recent-title">
        <div>
          <p className="eyebrow">Ritmo de temporada</p>
          <h2 id="command-recent-title">Resultados recientes</h2>
        </div>
        {data.recentResults.length === 0 ? (
          <p className="command-recent__empty">Todavía no hay partidos. El calendario avanza cuando corre el tiempo.</p>
        ) : (
          <ul>
            {data.recentResults.slice(0, 4).map((result) => (
              <li key={result.match_id}>
                <span className={result.winner_team_id === "pulsecity" ? "win" : "loss"}>
                  {result.winner_team_id === "pulsecity" ? "W" : "L"}
                </span>
                <strong>{formatMatchScore(result)}</strong>
                <small>{formatSimulatedDate(result.simulated_date)}</small>
              </li>
            ))}
          </ul>
        )}
      </section>
    </section>
  );
}

function OverviewMetric({ label, value }: { label: string; value: string }) {
  return <div><span>{label}</span><strong>{value}</strong></div>;
}

function formatMoney(value: number): string {
  const amount = Math.abs(value) / 1_000_000;
  const sign = value < 0 ? "−" : "";
  return `${sign}$${amount >= 10 ? amount.toFixed(0) : amount.toFixed(1)}M`;
}
