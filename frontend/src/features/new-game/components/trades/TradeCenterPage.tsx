import { useEffect, useMemo, useState, type CSSProperties, type FormEvent } from "react";

import skylineBackdrop from "../../../assets/landing-city-night.svg";
import type {
  FinanceClientState,
  PlayerEmotionalState,
  RosterClientStates,
  TimeClientState,
  TradeClientState,
  TradeClientStates,
} from "../../../../types";
import { formatGameId } from "../../helpers";
import { rivalTeams, requestedPositions } from "../../trades/constants";
import type {
  TradeAcceptanceInput,
  TradeProposalInput,
  TradeRequestResult,
} from "../../trades/types";
import "./tradeCenter.css";

interface TradeCenterPageProps {
  acceptingProposalIds: Set<string>;
  error: string | null;
  financeState: FinanceClientState;
  gameId: string;
  rosterStates: RosterClientStates;
  submittingProposal: boolean;
  timeState: TimeClientState;
  trades: TradeClientStates;
  onAcceptTrade: (input: TradeAcceptanceInput) => Promise<TradeRequestResult>;
  onBack: () => void;
  onProposeTrade: (input: TradeProposalInput) => Promise<TradeRequestResult>;
}

export function TradeCenterPage({
  acceptingProposalIds,
  error,
  financeState,
  gameId,
  rosterStates,
  submittingProposal,
  timeState,
  trades,
  onAcceptTrade,
  onBack,
  onProposeTrade,
}: TradeCenterPageProps) {
  const availablePlayers = useMemo(
    () => Object.values(rosterStates).filter(isTradeEligible),
    [rosterStates],
  );
  const negotiations = useMemo(() => Object.values(trades).reverse(), [trades]);
  const [offeredPlayerId, setOfferedPlayerId] = useState("");
  const [rivalTeamId, setRivalTeamId] = useState("bos");
  const [requestedPosition, setRequestedPosition] = useState("PG");
  const [incomingSalaryMillions, setIncomingSalaryMillions] = useState("12");

  useEffect(() => {
    if (!availablePlayers.some((player) => player.player_id === offeredPlayerId)) {
      setOfferedPlayerId(availablePlayers[0]?.player_id ?? "");
    }
  }, [availablePlayers, offeredPlayerId]);

  async function submitProposal(event: FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    const incomingSalary = Math.round(Number(incomingSalaryMillions) * 1_000_000);
    if (!offeredPlayerId || !Number.isFinite(incomingSalary) || incomingSalary <= 0) {
      return;
    }

    await onProposeTrade({
      rivalTeamId,
      offeredPlayerId,
      requestedPosition,
      incomingSalary,
    });
  }

  return (
    <section
      className="trade-center"
      style={{ "--trade-center-backdrop": `url("${skylineBackdrop}")` } as CSSProperties}
    >
      <div className="trade-center__image" />
      <div className="trade-center__shade" />

      <header className="trade-center__topbar">
        <button type="button" className="trade-center__back" onClick={onBack}>
          Volver al Command Center
        </button>
        <div className="trade-center__identity">
          <p className="eyebrow">Basketball Operations</p>
          <strong>Trade Center</strong>
        </div>
        <div className="trade-center__context" aria-label="Contexto de la partida">
          <ContextValue label="Fecha" value={formatDate(timeState.simulated_date)} />
          <ContextValue label="Cap space" value={formatMoney(financeState.cap_space)} />
          <ContextValue label="Partida" value={formatGameId(gameId)} />
        </div>
      </header>

      <main className="trade-center__workspace">
        <section className="trade-desk" aria-labelledby="trade-desk-title">
          <div className="trade-center__section-heading">
            <p className="eyebrow">Mesa de propuesta</p>
            <h1 id="trade-desk-title">Abrir una negociación</h1>
            <p>Definí el costo propio y el perfil de jugador que querés recibir.</p>
          </div>

          <form className="trade-form" onSubmit={submitProposal}>
            <label>
              <span>Jugador ofrecido</span>
              <select
                value={offeredPlayerId}
                onChange={(event) => setOfferedPlayerId(event.target.value)}
                disabled={availablePlayers.length === 0 || submittingProposal}
              >
                {availablePlayers.length === 0 ? (
                  <option value="">Roster aún no disponible</option>
                ) : (
                  availablePlayers.map((player) => (
                    <option key={player.player_id} value={player.player_id}>
                      {playerLabel(player)}
                    </option>
                  ))
                )}
              </select>
              <small>Solo aparecen jugadores disponibles en el roster activo.</small>
            </label>

            <label>
              <span>Franquicia rival</span>
              <select
                value={rivalTeamId}
                onChange={(event) => setRivalTeamId(event.target.value)}
                disabled={submittingProposal}
              >
                {rivalTeams.map((team) => (
                  <option key={team.id} value={team.id}>
                    {team.name} · GM {team.gm}
                  </option>
                ))}
              </select>
            </label>

            <div className="trade-form__row">
              <label>
                <span>Posición buscada</span>
                <select
                  value={requestedPosition}
                  onChange={(event) => setRequestedPosition(event.target.value)}
                  disabled={submittingProposal}
                >
                  {requestedPositions.map((position) => (
                    <option key={position} value={position}>{position}</option>
                  ))}
                </select>
              </label>

              <label>
                <span>Salario entrante</span>
                <div className="trade-form__salary">
                  <span>$</span>
                  <input
                    type="number"
                    min="0.5"
                    step="0.5"
                    value={incomingSalaryMillions}
                    onChange={(event) => setIncomingSalaryMillions(event.target.value)}
                    disabled={submittingProposal}
                    aria-describedby="incoming-salary-unit"
                  />
                  <span id="incoming-salary-unit">M</span>
                </div>
              </label>
            </div>

            {error ? <p className="trade-form__error" role="alert">{error}</p> : null}

            <button
              type="submit"
              className="trade-form__submit"
              disabled={!offeredPlayerId || submittingProposal}
            >
              {submittingProposal ? "Enviando propuesta..." : "Enviar al GM rival"}
            </button>
          </form>
        </section>

        <section className="trade-ledger" aria-labelledby="trade-ledger-title">
          <div className="trade-ledger__header">
            <div>
              <p className="eyebrow">Negociaciones vivas</p>
              <h2 id="trade-ledger-title">Sala de respuestas</h2>
            </div>
            <span>{negotiations.length} negociaciones registradas</span>
          </div>

          <div className="trade-ledger__list" aria-live="polite">
            {negotiations.length === 0 ? (
              <div className="trade-ledger__empty">
                <strong>Todavía no hay conversaciones.</strong>
                <p>La respuesta del primer GM rival aparecerá acá cuando envíes una propuesta.</p>
              </div>
            ) : (
              negotiations.map((trade) => (
                <NegotiationCard
                  key={trade.proposal_id}
                  trade={trade}
                  accepting={acceptingProposalIds.has(trade.proposal_id)}
                  onAccept={onAcceptTrade}
                />
              ))
            )}
          </div>
        </section>
      </main>
    </section>
  );
}

interface NegotiationCardProps {
  accepting: boolean;
  trade: TradeClientState;
  onAccept: (input: TradeAcceptanceInput) => Promise<TradeRequestResult>;
}

function NegotiationCard({ accepting, trade, onAccept }: NegotiationCardProps) {
  const rival = rivalTeams.find((team) => team.id === trade.rival_team_id);
  const asset = trade.additional_asset_required ?? "";

  return (
    <article className={`trade-negotiation trade-negotiation--${trade.status}`}>
      <div className="trade-negotiation__header">
        <div>
          <span>{rival?.name ?? trade.rival_team_id.toUpperCase()}</span>
          <strong>{rival ? `GM ${rival.gm}` : "Oficina rival"}</strong>
        </div>
        <span className={`trade-status trade-status--${trade.status}`}>
          {statusLabel(trade.status)}
        </span>
      </div>

      <p className="trade-negotiation__detail">{tradeDetail(trade)}</p>

      <dl className="trade-negotiation__terms">
        <Term label="Sale" value={trade.outgoing_player_name ?? trade.offered_player_name ?? trade.offered_player_id ?? "Pendiente"} />
        <Term label="Buscás" value={trade.incoming_player_name ?? trade.requested_position ?? "Pendiente"} />
        <Term label="Salario" value={trade.incoming_salary ? formatMoney(trade.incoming_salary) : "Pendiente"} />
        <Term label="Fecha" value={formatDate(trade.simulated_date)} />
      </dl>

      {trade.status === "countered" ? (
        <div className="trade-negotiation__counter">
          <p>
            Condición adicional
            <strong>{formatAsset(asset)}</strong>
          </p>
          <button
            type="button"
            disabled={accepting}
            onClick={() => void onAccept({
              proposalId: trade.proposal_id,
              acceptedAdditionalAsset: asset,
            })}
          >
            {accepting ? "Confirmando..." : "Aceptar contraoferta"}
          </button>
        </div>
      ) : null}

      {trade.status === "accepted" && trade.incoming_player_name ? (
        <p className="trade-negotiation__result">
          Llega <strong>{trade.incoming_player_name}</strong>
          {trade.incoming_rating ? ` · OVR ${trade.incoming_rating}` : ""}
        </p>
      ) : null}
    </article>
  );
}

function ContextValue({ label, value }: { label: string; value: string }) {
  return <span><small>{label}</small><strong>{value}</strong></span>;
}

function Term({ label, value }: { label: string; value: string }) {
  return <div><dt>{label}</dt><dd>{value}</dd></div>;
}

function isTradeEligible(player: PlayerEmotionalState): boolean {
  return player.availability !== "traded" && player.emotional_state !== "traded";
}

function playerLabel(player: PlayerEmotionalState): string {
  const name = player.full_name ?? player.player_id;
  const position = player.position ? ` · ${player.position}` : "";
  const rating = player.overall_rating ? ` · OVR ${player.overall_rating}` : "";
  return `${name}${position}${rating}`;
}

function tradeDetail(trade: TradeClientState): string {
  if (trade.detail) {
    return trade.detail;
  }
  if (trade.status === "proposed") {
    return "La propuesta superó la validación local y espera respuesta del GM rival.";
  }
  if (trade.status === "accepted") {
    return "El acuerdo quedó persistido y el roster ya fue actualizado.";
  }
  return trade.reason ? formatAsset(trade.reason) : "La negociación fue actualizada.";
}

function statusLabel(status: TradeClientState["status"]): string {
  switch (status) {
    case "countered":
      return "Contraoferta";
    case "rejected":
      return "Rechazada";
    case "accepted":
      return "Aceptada";
    default:
      return "En evaluación";
  }
}

function formatAsset(value: string): string {
  return value ? value.replaceAll("_", " ") : "Sin asset adicional";
}

function formatDate(value: string): string {
  if (!value) {
    return "Sin fecha";
  }
  return new Intl.DateTimeFormat("es-AR", { day: "2-digit", month: "short", year: "numeric" })
    .format(new Date(`${value}T00:00:00Z`));
}

function formatMoney(value: number): string {
  const sign = value < 0 ? "-" : "";
  const amount = Math.abs(value) / 1_000_000;
  return `${sign}$${amount.toFixed(amount >= 10 ? 0 : 1)}M`;
}
