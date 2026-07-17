import { useCallback, useEffect, useState } from "react";

import type {
  MedicalDecisionInput,
  MedicalDecisionRecord,
  MedicalDecisionRequestResult,
} from "./types";

interface UseMedicalOperationsOptions {
  activeAuthKind: "none" | "guest" | "user";
  gameId: string;
  gatewayBaseUrl: string;
  guestToken: string;
  sessionToken?: string;
  simulatedDate: string;
  onStatusChange: (message: string) => void;
}

interface MedicalOperationsResult {
  decisionsByInjury: Record<string, MedicalDecisionRecord>;
  errorsByInjury: Record<string, string>;
  submittingInjuryIds: Set<string>;
  submitMedicalDecision: (
    input: MedicalDecisionInput,
  ) => Promise<MedicalDecisionRequestResult>;
}

export function useMedicalOperations({
  activeAuthKind,
  gameId,
  gatewayBaseUrl,
  guestToken,
  sessionToken,
  simulatedDate,
  onStatusChange,
}: UseMedicalOperationsOptions): MedicalOperationsResult {
  const [decisionsByInjury, setDecisionsByInjury] =
    useState<Record<string, MedicalDecisionRecord>>({});
  const [errorsByInjury, setErrorsByInjury] = useState<Record<string, string>>({});
  const [submittingInjuryIds, setSubmittingInjuryIds] = useState<Set<string>>(new Set());

  useEffect(() => {
    setDecisionsByInjury({});
    setErrorsByInjury({});
    setSubmittingInjuryIds(new Set());
  }, [gameId]);

  const submitMedicalDecision = useCallback(async (
    input: MedicalDecisionInput,
  ): Promise<MedicalDecisionRequestResult> => {
    if (!gameId || activeAuthKind === "none") {
      const message = "Necesitás una partida activa para responder al equipo médico.";
      setErrorsByInjury((current) => ({ ...current, [input.injuryId]: message }));
      onStatusChange(message);
      return { ok: false };
    }

    setSubmittingInjuryIds((current) => new Set(current).add(input.injuryId));
    setErrorsByInjury((current) => {
      const next = { ...current };
      delete next[input.injuryId];
      return next;
    });
    onStatusChange("Registrando la decisión médica...");

    try {
      const response = await fetch(
        `${gatewayBaseUrl}/api/v1/games/${gameId}/medical-decisions`,
        {
          method: "POST",
          headers: buildAuthHeaders(guestToken, sessionToken),
          body: JSON.stringify({
            injury_id: input.injuryId,
            player_id: input.playerId,
            choice_id: input.choiceId,
            simulated_date: simulatedDate,
          }),
        },
      );
      const payload = (await response.json()) as { decision_id?: string; error?: string };
      if (!response.ok || !payload.decision_id) {
        const message = payload.error ?? "No se pudo registrar la decisión médica.";
        setErrorsByInjury((current) => ({ ...current, [input.injuryId]: message }));
        onStatusChange(message);
        return { ok: false };
      }

      setDecisionsByInjury((current) => ({
        ...current,
        [input.injuryId]: {
          choiceId: input.choiceId,
          decisionId: payload.decision_id ?? "",
        },
      }));
      onStatusChange("Decisión registrada. El staff médico procesará sus consecuencias.");
      return { ok: true, decisionId: payload.decision_id };
    } catch (requestError) {
      const message = requestError instanceof Error
        ? requestError.message
        : "Fallo de red al registrar la decisión médica.";
      setErrorsByInjury((current) => ({ ...current, [input.injuryId]: message }));
      onStatusChange(message);
      return { ok: false };
    } finally {
      setSubmittingInjuryIds((current) => {
        const next = new Set(current);
        next.delete(input.injuryId);
        return next;
      });
    }
  }, [activeAuthKind, gameId, gatewayBaseUrl, guestToken, onStatusChange, sessionToken, simulatedDate]);

  return {
    decisionsByInjury,
    errorsByInjury,
    submittingInjuryIds,
    submitMedicalDecision,
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
