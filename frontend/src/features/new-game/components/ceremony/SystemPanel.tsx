import type { MapClientState, RealtimeEvent } from "../../../../types";
import { stageMeta, stageSequence } from "../../constants";
import { describeRealtimeEvent } from "../../helpers";

interface SystemPanelProps {
  events: RealtimeEvent[];
  mapState: MapClientState;
}

export function SystemPanel({ events, mapState }: SystemPanelProps) {
  const stageIndex = stageSequence.indexOf(mapState.stage as (typeof stageSequence)[number]);
  const completedSteps = stageIndex >= 0 ? stageIndex : -1;

  return (
    <section className="ceremony-panel ceremony-panel--tab">
      <p className="eyebrow">Sistema</p>

      <ol className="ceremony-pipeline">
        {stageSequence.map((stage, index) => {
          const isActive = mapState.stage === stage;
          const isDone = completedSteps >= index;

          return (
            <li
              key={stage}
              className={[
                "ceremony-pipeline__item",
                isActive ? "active" : "",
                isDone ? "done" : "",
              ]
                .filter(Boolean)
                .join(" ")}
            >
              <span>0{index + 1}</span>
              <div>
                <strong>{stageMeta[stage].label}</strong>
                <small>{stageMeta[stage].title}</small>
              </div>
            </li>
          );
        })}
      </ol>

      <div>
        <p className="eyebrow">Eventos recientes</p>
        <ul className="ceremony-events">
          {events.length === 0 ? (
            <li>Todavia no llegaron eventos para esta partida.</li>
          ) : (
            events.map((event, index) => (
              <li key={`${event.subject}-${index}`}>
                <strong>{event.subject}</strong>
                <span>{describeRealtimeEvent(event)}</span>
              </li>
            ))
          )}
        </ul>
      </div>
    </section>
  );
}
