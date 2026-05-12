import type { NarrativeChoice, NarrativeEvent } from "../../../types";
import { StatusBadge } from "./common";

interface OwnerIntroModalProps {
  event: NarrativeEvent;
  submitting: boolean;
  onSelect: (choice: NarrativeChoice) => void;
}

export function OwnerIntroModal(props: OwnerIntroModalProps) {
  return (
    <div className="narrative-overlay" role="dialog" aria-modal="true">
      <article className="narrative-modal">
        <div className="panel-header">
          <p className="eyebrow">Evento obligatorio</p>
          <h2>{props.event.title}</h2>
        </div>

        <StatusBadge label="Owner" tone="urgent" />
        <p className="narrative-body">{props.event.body}</p>

        <div className="narrative-actions">
          {(props.event.choices ?? [{ id: "continue", label: "Entendido" }]).map((choice) => (
            <button
              key={choice.id}
              type="button"
              className="primary-action"
              onClick={() => props.onSelect(choice)}
              disabled={props.submitting}
            >
              {props.submitting ? "Confirmando..." : choice.label}
            </button>
          ))}
        </div>
      </article>
    </div>
  );
}
