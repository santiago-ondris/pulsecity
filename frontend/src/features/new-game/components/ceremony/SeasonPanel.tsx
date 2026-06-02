import type {
  CityClientState,
  MapClientState,
  SeasonClientState,
  SeasonMatchSummary,
} from "../../../../types";
import { formatGameId } from "../../helpers";
import { formatMatchScore, formatPointDifferential, formatSimulatedDate } from "./helpers";

interface SeasonPanelProps {
  cityState: CityClientState;
  gameId: string;
  mapState: MapClientState;
  ownerIntroResponseLabel: string | null;
  recentResults: SeasonMatchSummary[];
  seasonState: SeasonClientState;
  status: string;
}

export function SeasonPanel({
  cityState,
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

function StateItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="ceremony-state-item">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}
