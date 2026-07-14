import type {
  CityClientState,
  FinanceClientState,
  MapClientState,
  SeasonClientState,
  SeasonMatchSummary,
} from "../../../../types";
import { formatGameId } from "../../helpers";
import { formatMatchScore, formatPointDifferential, formatSimulatedDate } from "./helpers";

interface SeasonPanelProps {
  cityState: CityClientState;
  financeState: FinanceClientState;
  gameId: string;
  mapState: MapClientState;
  ownerIntroResponseLabel: string | null;
  recentResults: SeasonMatchSummary[];
  seasonState: SeasonClientState;
  status: string;
}

export function SeasonPanel({
  cityState,
  financeState,
  gameId,
  mapState,
  ownerIntroResponseLabel,
  recentResults,
  seasonState,
  status,
}: SeasonPanelProps) {
  return (
    <section className="ceremony-panel ceremony-panel--tab">
      <div className="ceremony-panel__title">
        <div>
          <p className="eyebrow">Temporada viva</p>
          <strong>{seasonState.wins + seasonState.losses}/82</strong>
        </div>
      </div>

      <div className="ceremony-scoreline">
        <strong>{seasonState.wins}-{seasonState.losses}</strong>
        <span>{formatPointDifferential(seasonState)}</span>
      </div>

      <div className="ceremony-state-grid">
        <StateItem label="Partida" value={formatGameId(mapState.game_id || gameId)} />
        <StateItem
          label="Ciudad"
          value={`${Math.round(cityState.fan_sentiment)} animo / ${Math.round(cityState.stadium_district_land_value)} suelo`}
        />
        <StateItem label="Entradas" value={`${Math.round(cityState.ticket_sales_index)} demanda`} />
        <StateItem label="Sistema" value={status} />
        <StateItem label="Owner" value={ownerIntroResponseLabel ?? "Pendiente de llamada"} />
      </div>

      <div>
        <p className="eyebrow">Finanzas</p>
        <div className="ceremony-state-grid">
          <StateItem label="Salarios" value={formatMoney(financeState.committed_salary)} />
          <StateItem label="Cap space" value={formatMoney(financeState.cap_space)} />
          <StateItem label="Luxury tax" value={formatMoney(financeState.luxury_tax_space)} />
          <StateItem label="Estado" value={formatCapStatus(financeState.status)} />
        </div>
      </div>

      <div>
        <p className="eyebrow">Resultados recientes</p>
        <ul className="ceremony-results ceremony-results--tab">
          {recentResults.length === 0 ? (
            <li className="empty">Todavia no hay partidos finalizados.</li>
          ) : (
            recentResults.map((result) => (
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
      </div>
    </section>
  );
}

function formatCapStatus(status: FinanceClientState["status"]) {
  switch (status) {
    case "luxury_tax":
      return "Luxury tax";
    case "over_cap":
      return "Sobre el cap";
    default:
      return "Bajo el cap";
  }
}

function formatMoney(value: number) {
  const sign = value < 0 ? "-" : "";
  const absoluteValue = Math.abs(value);

  if (absoluteValue >= 1_000_000) {
    const amount = absoluteValue / 1_000_000;
    const formatted = amount >= 10 ? amount.toFixed(0) : amount.toFixed(1);
    return `${sign}$${formatted}M`;
  }

  return `${sign}$${absoluteValue.toLocaleString("en-US")}`;
}

function StateItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="ceremony-state-item">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}
