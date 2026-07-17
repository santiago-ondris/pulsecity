import type {
  CityClientState,
  FinanceClientState,
  MapClientState,
  NarrativeEvent,
  SeasonClientState,
} from "../../../../types";
import { CitySnapshotCard } from "./CitySnapshotCard";

interface SeasonKickoffPanelProps {
  cityName: string;
  cityState: CityClientState;
  error: string | null;
  financeState: FinanceClientState;
  franchiseName: string;
  mapState: MapClientState;
  narrativeInbox: NarrativeEvent[];
  ownerIntroResponseLabel: string | null;
  rosterCount: number;
  seasonState: SeasonClientState;
  starting: boolean;
  onOpenMedicalCenter: () => void;
  onOpenStaff: () => void;
  onOpenTradeCenter: () => void;
  onStartSeason: () => void;
}

export function SeasonKickoffPanel({
  cityName,
  cityState,
  error,
  financeState,
  franchiseName,
  mapState,
  narrativeInbox,
  ownerIntroResponseLabel,
  rosterCount,
  seasonState,
  starting,
  onOpenMedicalCenter,
  onOpenStaff,
  onOpenTradeCenter,
  onStartSeason,
}: SeasonKickoffPanelProps) {
  return (
    <main className="season-kickoff">
      <section className="season-kickoff__decision" aria-labelledby="season-kickoff-title">
        <div className="season-kickoff__copy">
          <p className="eyebrow">La temporada espera tu orden</p>
          <h1 id="season-kickoff-title">{franchiseName} está lista.</h1>
          <p>
            El mundo permanece pausado. Cuando empieces, cada día moverá partidos,
            relaciones, lesiones y decisiones alrededor de la franquicia.
          </p>
        </div>

        <div className="season-kickoff__brief" aria-label="Situación inicial">
          <KickoffFact
            label="Mandato del Owner"
            value={ownerIntroResponseLabel ?? "Dirección pendiente"}
          />
          <KickoffFact
            label="Temporada"
            value={`${seasonState.wins + seasonState.losses}/82 partidos`}
          />
          <KickoffFact
            label="Cap space"
            value={financeState.cap_base > 0 ? formatMoney(financeState.cap_space) : "Esperando reporte del CFO"}
          />
          <KickoffFact
            label="Plantel observado"
            value={rosterCount > 0 ? `${rosterCount} jugadores` : "Esperando parte del staff"}
          />
          <KickoffFact
            label="Alertas"
            value={narrativeInbox.length > 0 ? `${narrativeInbox.length} pendientes` : "Sin urgencias"}
          />
        </div>

        <div className="season-kickoff__actions">
          <button
            type="button"
            className="season-kickoff__start"
            disabled={starting}
            onClick={onStartSeason}
          >
            <span>{starting ? "Poniendo el mundo en marcha..." : "Comenzar temporada"}</span>
            <small>La simulación avanzará en velocidad x1</small>
          </button>
          {error ? <p className="season-kickoff__error" role="alert">{error}</p> : null}
        </div>

        <div className="season-kickoff__secondary" aria-label="Revisiones opcionales">
          <span>Antes de empezar, si querés:</span>
          <button type="button" onClick={onOpenTradeCenter}>Revisar Trade Center</button>
          <button type="button" onClick={onOpenMedicalCenter}>Revisar Centro Médico</button>
          <button type="button" onClick={onOpenStaff}>Hablar con el staff</button>
        </div>
      </section>

      <CitySnapshotCard cityName={cityName} cityState={cityState} mapState={mapState} />
    </main>
  );
}

function KickoffFact({ label, value }: { label: string; value: string }) {
  return <div><span>{label}</span><strong>{value}</strong></div>;
}

function formatMoney(value: number): string {
  const amount = Math.abs(value) / 1_000_000;
  const sign = value < 0 ? "−" : "";
  return `${sign}$${amount >= 10 ? amount.toFixed(0) : amount.toFixed(1)}M`;
}
