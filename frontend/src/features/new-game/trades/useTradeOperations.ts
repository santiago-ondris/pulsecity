import { useCallback, useEffect, useState } from "react";

import type { TradeClientStates, TradePatchEnvelope } from "../../../types";
import type {
  TradeAcceptanceInput,
  TradeProposalInput,
  TradeRequestResult,
} from "./types";

interface UseTradeOperationsOptions {
  activeAuthKind: "none" | "guest" | "user";
  gameId: string;
  gatewayBaseUrl: string;
  guestToken: string;
  sessionToken?: string;
  simulatedDate: string;
  onStatusChange: (message: string) => void;
}

interface TradeOperationsResult {
  trades: TradeClientStates;
  submittingProposal: boolean;
  acceptingProposalIds: Set<string>;
  error: string | null;
  applyTradePatch: (event: TradePatchEnvelope) => void;
  proposeTrade: (input: TradeProposalInput) => Promise<TradeRequestResult>;
  acceptTrade: (input: TradeAcceptanceInput) => Promise<TradeRequestResult>;
}

export function useTradeOperations({
  activeAuthKind,
  gameId,
  gatewayBaseUrl,
  guestToken,
  sessionToken,
  simulatedDate,
  onStatusChange,
}: UseTradeOperationsOptions): TradeOperationsResult {
  const [trades, setTrades] = useState<TradeClientStates>({});
  const [submittingProposal, setSubmittingProposal] = useState(false);
  const [acceptingProposalIds, setAcceptingProposalIds] = useState<Set<string>>(new Set());
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setTrades({});
    setSubmittingProposal(false);
    setAcceptingProposalIds(new Set());
    setError(null);
  }, [gameId]);

  const applyTradePatch = useCallback((event: TradePatchEnvelope): void => {
    setTrades((current) => ({
      ...current,
      [event.patch.proposal_id]: {
        ...current[event.patch.proposal_id],
        ...event.patch,
      },
    }));
  }, []);

  const proposeTrade = useCallback(async (input: TradeProposalInput): Promise<TradeRequestResult> => {
    if (!gameId || activeAuthKind === "none") {
      const message = "Necesitás una partida activa para proponer un trade.";
      setError(message);
      onStatusChange(message);
      return { ok: false };
    }

    setSubmittingProposal(true);
    setError(null);
    onStatusChange("Enviando propuesta al GM rival...");

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${gameId}/trades/proposals`, {
        method: "POST",
        headers: buildAuthHeaders(guestToken, sessionToken),
        body: JSON.stringify({
          rival_team_id: input.rivalTeamId,
          offered_player_id: input.offeredPlayerId,
          requested_position: input.requestedPosition,
          incoming_salary: input.incomingSalary,
          simulated_date: simulatedDate,
        }),
      });
      const payload = (await response.json()) as { proposal_id?: string; error?: string };
      if (!response.ok || !payload.proposal_id) {
        const message = payload.error ?? "No se pudo enviar la propuesta.";
        setError(message);
        onStatusChange(message);
        return { ok: false };
      }

      onStatusChange("Propuesta enviada. Esperando la evaluación del GM rival.");
      return { ok: true, proposalId: payload.proposal_id };
    } catch (requestError) {
      const message = requestError instanceof Error
        ? requestError.message
        : "Fallo de red al proponer el trade.";
      setError(message);
      onStatusChange(message);
      return { ok: false };
    } finally {
      setSubmittingProposal(false);
    }
  }, [activeAuthKind, gameId, gatewayBaseUrl, guestToken, onStatusChange, sessionToken, simulatedDate]);

  const acceptTrade = useCallback(async (input: TradeAcceptanceInput): Promise<TradeRequestResult> => {
    if (!gameId || activeAuthKind === "none") {
      const message = "Necesitás una partida activa para aceptar una contraoferta.";
      setError(message);
      onStatusChange(message);
      return { ok: false };
    }

    setAcceptingProposalIds((current) => new Set(current).add(input.proposalId));
    setError(null);
    onStatusChange("Confirmando la contraoferta...");

    try {
      const response = await fetch(`${gatewayBaseUrl}/api/v1/games/${gameId}/trades/acceptances`, {
        method: "POST",
        headers: buildAuthHeaders(guestToken, sessionToken),
        body: JSON.stringify({
          proposal_id: input.proposalId,
          accepted_additional_asset: input.acceptedAdditionalAsset,
          simulated_date: simulatedDate,
        }),
      });
      const payload = (await response.json()) as { proposal_id?: string; error?: string };
      if (!response.ok || !payload.proposal_id) {
        const message = payload.error ?? "No se pudo aceptar la contraoferta.";
        setError(message);
        onStatusChange(message);
        return { ok: false };
      }

      onStatusChange("Aceptación enviada. Esperando la confirmación del cierre.");
      return { ok: true, proposalId: payload.proposal_id };
    } catch (requestError) {
      const message = requestError instanceof Error
        ? requestError.message
        : "Fallo de red al aceptar la contraoferta.";
      setError(message);
      onStatusChange(message);
      return { ok: false };
    } finally {
      setAcceptingProposalIds((current) => {
        const next = new Set(current);
        next.delete(input.proposalId);
        return next;
      });
    }
  }, [activeAuthKind, gameId, gatewayBaseUrl, guestToken, onStatusChange, sessionToken, simulatedDate]);

  return {
    trades,
    submittingProposal,
    acceptingProposalIds,
    error,
    applyTradePatch,
    proposeTrade,
    acceptTrade,
  };
}

function buildAuthHeaders(guestToken: string, sessionToken?: string): Record<string, string> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (sessionToken) {
    headers["X-Session-Token"] = sessionToken;
  } else if (guestToken) {
    headers["X-Guest-Token"] = guestToken;
  }
  return headers;
}
