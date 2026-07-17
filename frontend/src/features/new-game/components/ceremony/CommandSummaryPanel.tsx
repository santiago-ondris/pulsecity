import type { CeremonySharedProps } from "./types";

interface CommandSummaryPanelProps {
  data: CeremonySharedProps;
  onOpenStaff: () => void;
}

export function CommandSummaryPanel({ data, onOpenStaff }: CommandSummaryPanelProps) {
  const gamesPlayed = data.seasonState.wins + data.seasonState.losses;
  const urgentEvents = data.narrativeInbox.filter((event) => (
    event.urgency === "urgent" || event.urgency === "critical"
  ));

  return (
    <section className="command-summary-panel">
      <div className="command-summary-panel__state">
        <p className="eyebrow">Ahora</p>
        <strong>{urgentEvents.length > 0 ? `${urgentEvents.length} alertas urgentes` : "Sin urgencias activas"}</strong>
        <p>
          {gamesPlayed === 0
            ? "La temporada todavía no avanzó. Iniciá el tiempo cuando estés listo."
            : `${gamesPlayed}/82 partidos completados. El mundo sigue reaccionando a tus decisiones.`}
        </p>
      </div>

      <div className="command-summary-panel__list">
        <SummaryItem label="Owner" value={data.ownerIntroResponseLabel ?? "Sin mandato registrado"} />
        <SummaryItem label="Inbox" value={`${data.narrativeInbox.length} eventos`} />
        <SummaryItem label="Ciudad" value={`${Math.round(data.cityState.fan_sentiment)} ánimo`} />
        <SummaryItem label="Sistema" value={data.timeState.paused ? "Tiempo pausado" : `Corriendo x${data.timeState.speed}`} />
      </div>

      <button type="button" className="command-summary-panel__staff" onClick={onOpenStaff}>
        Buscar y hablar con un agente
      </button>
    </section>
  );
}

function SummaryItem({ label, value }: { label: string; value: string }) {
  return <div><span>{label}</span><strong>{value}</strong></div>;
}
