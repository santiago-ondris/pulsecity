import type { NarrativeEvent } from "../../../../types";

interface InboxPanelProps {
  events: NarrativeEvent[];
}

export function InboxPanel({ events }: InboxPanelProps) {
  return (
    <section className="ceremony-panel ceremony-panel--tab">
      <div className="ceremony-panel__title">
        <div>
          <p className="eyebrow">Inbox narrativo</p>
          <strong>{events.length} eventos</strong>
        </div>
      </div>
      <ul className="ceremony-inbox ceremony-inbox--tab">
        {events.length === 0 ? (
          <li className="empty">Los reportes post-partido y mensajes importantes van a aparecer aca.</li>
        ) : (
          events.map((event) => (
            <li key={event.event_id}>
              <div>
                <strong>{event.title}</strong>
                <span>{event.emitter}</span>
              </div>
              <p>{event.body}</p>
            </li>
          ))
        )}
      </ul>
    </section>
  );
}
