import { useMemo, type CSSProperties } from "react";

import skylineBackdrop from "../../../assets/landing-city-night.svg";
import type { PlayerEmotionalState, RosterClientStates, TimeClientState } from "../../../../types";
import { formatGameId } from "../../helpers";
import type {
  MedicalDecisionChoice,
  MedicalDecisionInput,
  MedicalDecisionRecord,
  MedicalDecisionRequestResult,
} from "../../medical/types";
import "./medicalCenter.css";

interface MedicalCenterPageProps {
  decisionsByInjury: Record<string, MedicalDecisionRecord>;
  errorsByInjury: Record<string, string>;
  gameId: string;
  rosterStates: RosterClientStates;
  submittingInjuryIds: Set<string>;
  timeState: TimeClientState;
  onBack: () => void;
  onSubmitDecision: (
    input: MedicalDecisionInput,
  ) => Promise<MedicalDecisionRequestResult>;
}

interface DecisionOption {
  id: MedicalDecisionChoice;
  label: string;
  detail: string;
  tone: "recommended" | "cautious" | "warning" | "urgent";
}

const decisionOptions: DecisionOption[] = [
  {
    id: "rest",
    label: "Seguir protocolo",
    detail: "Mantener la baja hasta el alta médica prevista.",
    tone: "recommended",
  },
  {
    id: "reduce_minutes",
    label: "Carga reducida al volver",
    detail: "Registrar un regreso protegido para la rotación.",
    tone: "cautious",
  },
  {
    id: "ignore_doctor",
    label: "Ignorar recomendación",
    detail: "Sostener la postura deportiva contra el criterio médico.",
    tone: "warning",
  },
  {
    id: "force_return",
    label: "Forzar alta anticipada",
    detail: "Activarlo ahora, asumiendo riesgo de reagravación.",
    tone: "urgent",
  },
];

const severityOrder = { major: 0, moderate: 1, minor: 2 } as const;

export function MedicalCenterPage({
  decisionsByInjury,
  errorsByInjury,
  gameId,
  rosterStates,
  submittingInjuryIds,
  timeState,
  onBack,
  onSubmitDecision,
}: MedicalCenterPageProps) {
  const players = useMemo(() => Object.values(rosterStates), [rosterStates]);
  const injuredPlayers = useMemo(
    () => [...players]
      .filter((player) => player.availability === "injured" && player.injury_id)
      .sort((left, right) => severityRank(left) - severityRank(right)),
    [players],
  );
  const availablePlayers = useMemo(
    () => players.filter((player) => (
      player.availability !== "injured"
      && player.availability !== "traded"
      && player.emotional_state !== "traded"
    )),
    [players],
  );
  const majorCases = injuredPlayers.filter((player) => player.severity === "major").length;

  return (
    <section
      className="medical-center"
      style={{ "--medical-center-backdrop": `url("${skylineBackdrop}")` } as CSSProperties}
    >
      <div className="medical-center__image" />
      <div className="medical-center__shade" />

      <header className="medical-center__topbar">
        <button type="button" className="medical-center__back" onClick={onBack}>
          Volver al Command Center
        </button>
        <div className="medical-center__identity">
          <p className="eyebrow">Performance & Health</p>
          <strong>Centro Médico</strong>
        </div>
        <div className="medical-center__context" aria-label="Contexto de la partida">
          <ContextValue label="Fecha" value={formatDate(timeState.simulated_date)} />
          <ContextValue label="Disponibles" value={`${availablePlayers.length}`} />
          <ContextValue label="Partida" value={formatGameId(gameId)} />
        </div>
      </header>

      <main className="medical-center__workspace">
        <section className="medical-board" aria-labelledby="medical-board-title">
          <div className="medical-board__heading">
            <div>
              <p className="eyebrow">Parte diario</p>
              <h1 id="medical-board-title">Disponibilidad del roster</h1>
              <p>El staff diagnostica. Vos decidís cuánto riesgo acepta la franquicia.</p>
            </div>
            <div className="medical-board__summary" aria-label="Resumen médico">
              <MedicalMetric label="Casos activos" value={injuredPlayers.length} tone="urgent" />
              <MedicalMetric label="Bajas graves" value={majorCases} tone="negative" />
              <MedicalMetric label="Listos" value={availablePlayers.length} tone="healthy" />
            </div>
          </div>

          <div className="medical-cases" aria-live="polite">
            {players.length === 0 ? (
              <EmptyState
                eyebrow="Sin snapshot clínico"
                title="Esperando el primer parte del roster"
                detail="Esta página se completa con roster.patch. No genera disponibilidad por su cuenta."
              />
            ) : injuredPlayers.length === 0 ? (
              <EmptyState
                eyebrow="Plantel disponible"
                title="No hay lesiones activas"
                detail="El equipo médico seguirá monitoreando carga y minutos después de cada partido."
              />
            ) : (
              injuredPlayers.map((player) => {
                const injuryId = player.injury_id ?? "";
                return (
                  <PatientCard
                    key={injuryId}
                    decision={decisionsByInjury[injuryId]}
                    error={errorsByInjury[injuryId]}
                    player={player}
                    submitting={submittingInjuryIds.has(injuryId)}
                    onSubmitDecision={onSubmitDecision}
                  />
                );
              })
            )}
          </div>
        </section>

        <aside className="medical-sidebar" aria-label="Contexto del staff médico">
          <section className="medical-brief">
            <p className="eyebrow">Criterio del staff</p>
            <h2>Disponibilidad no es salud</h2>
            <p>
              El Médico protege el alta. El Head Coach necesita rotación. Cada decisión también
              modifica la confianza entre ellos y con el GM.
            </p>
          </section>

          <section className="medical-protocol">
            <p className="eyebrow">Lectura de riesgo</p>
            <dl>
              <ProtocolRow label="Protocolo" detail="Conserva la fecha de recuperación" tone="healthy" />
              <ProtocolRow label="Carga reducida" detail="Protege el regreso a la rotación" tone="info" />
              <ProtocolRow label="Ignorar" detail="Deteriora la relación con el Médico" tone="warning" />
              <ProtocolRow label="Alta forzada" detail="Puede causar una nueva lesión" tone="urgent" />
            </dl>
          </section>

          <section className="medical-roster-strip">
            <div>
              <p className="eyebrow">Unidad disponible</p>
              <strong>{availablePlayers.length} jugadores listos</strong>
            </div>
            <ul>
              {availablePlayers.slice(0, 6).map((player) => (
                <li key={player.player_id}>
                  <span>{player.position ?? "—"}</span>
                  <strong>{playerName(player)}</strong>
                </li>
              ))}
            </ul>
            {availablePlayers.length > 6 ? (
              <small>+{availablePlayers.length - 6} disponibles en el roster</small>
            ) : null}
          </section>
        </aside>
      </main>
    </section>
  );
}

interface PatientCardProps {
  decision?: MedicalDecisionRecord;
  error?: string;
  player: PlayerEmotionalState;
  submitting: boolean;
  onSubmitDecision: (
    input: MedicalDecisionInput,
  ) => Promise<MedicalDecisionRequestResult>;
}

function PatientCard({
  decision,
  error,
  player,
  submitting,
  onSubmitDecision,
}: PatientCardProps) {
  const injuryId = player.injury_id ?? "";
  const selectedOption = decisionOptions.find((option) => option.id === decision?.choiceId);

  return (
    <article className={`patient-card patient-card--${player.severity ?? "moderate"}`}>
      <div className="patient-card__header">
        <div className="patient-card__identity">
          <span>{player.position ?? "—"}</span>
          <div>
            <strong>{playerName(player)}</strong>
            <small>{player.overall_rating ? `OVR ${player.overall_rating}` : player.player_id}</small>
          </div>
        </div>
        <span className={`medical-severity medical-severity--${player.severity ?? "moderate"}`}>
          {severityLabel(player.severity)}
        </span>
      </div>

      <dl className="patient-card__diagnosis">
        <DiagnosisValue label="Disponibilidad" value="Baja médica" />
        <DiagnosisValue
          label="Retorno estimado"
          value={formatDate(player.expected_recovery_date)}
        />
        <DiagnosisValue
          label="Días estimados"
          value={player.estimated_days_out ? `${player.estimated_days_out} días` : "Pendiente"}
        />
        <DiagnosisValue label="Último cambio" value={formatDate(player.availability_changed_on)} />
      </dl>

      {decision ? (
        <div className="patient-card__confirmed" role="status">
          <span>Decisión registrada</span>
          <strong>{selectedOption?.label ?? decision.choiceId}</strong>
          <p>
            El estado clínico seguirá cambiando únicamente cuando el backend publique un nuevo delta.
          </p>
        </div>
      ) : (
        <div className="patient-card__decision">
          <div>
            <p className="eyebrow">Decisión del GM</p>
            <strong>¿Cómo manejás este caso?</strong>
          </div>
          <div className="medical-decisions">
            {decisionOptions.map((option) => (
              <button
                key={option.id}
                type="button"
                className={`medical-decision medical-decision--${option.tone}`}
                disabled={submitting}
                onClick={() => void onSubmitDecision({
                  injuryId,
                  playerId: player.player_id,
                  choiceId: option.id,
                })}
              >
                <strong>{option.label}</strong>
                <span>{option.detail}</span>
              </button>
            ))}
          </div>
          {submitting ? <p className="patient-card__sending">Registrando decisión...</p> : null}
          {error ? <p className="patient-card__error" role="alert">{error}</p> : null}
        </div>
      )}
    </article>
  );
}

function EmptyState({ eyebrow, title, detail }: { eyebrow: string; title: string; detail: string }) {
  return (
    <div className="medical-cases__empty">
      <span className="medical-cases__pulse" />
      <div>
        <p className="eyebrow">{eyebrow}</p>
        <strong>{title}</strong>
        <p>{detail}</p>
      </div>
    </div>
  );
}

function MedicalMetric({
  label,
  value,
  tone,
}: {
  label: string;
  value: number;
  tone: "urgent" | "negative" | "healthy";
}) {
  return (
    <div className={`medical-metric medical-metric--${tone}`}>
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function ContextValue({ label, value }: { label: string; value: string }) {
  return <span><small>{label}</small><strong>{value}</strong></span>;
}

function DiagnosisValue({ label, value }: { label: string; value: string }) {
  return <div><dt>{label}</dt><dd>{value}</dd></div>;
}

function ProtocolRow({
  label,
  detail,
  tone,
}: {
  label: string;
  detail: string;
  tone: "healthy" | "info" | "warning" | "urgent";
}) {
  return (
    <div className={`medical-protocol__row medical-protocol__row--${tone}`}>
      <dt>{label}</dt>
      <dd>{detail}</dd>
    </div>
  );
}

function severityRank(player: PlayerEmotionalState): number {
  return player.severity ? severityOrder[player.severity] : severityOrder.moderate;
}

function severityLabel(severity?: PlayerEmotionalState["severity"]): string {
  if (severity === "major") {
    return "Grave";
  }
  if (severity === "minor") {
    return "Leve";
  }
  return "Moderada";
}

function playerName(player: PlayerEmotionalState): string {
  return player.full_name ?? player.player_id;
}

function formatDate(value?: string): string {
  if (!value) {
    return "Pendiente";
  }
  const date = new Date(`${value.slice(0, 10)}T00:00:00Z`);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return new Intl.DateTimeFormat("es-AR", {
    day: "2-digit",
    month: "short",
    timeZone: "UTC",
  }).format(date);
}
