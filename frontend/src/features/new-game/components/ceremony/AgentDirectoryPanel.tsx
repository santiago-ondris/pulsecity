import { useMemo, useState, type FormEvent } from "react";

import { canonicalAgents } from "./agents";
import { formatAgentMetric } from "./helpers";
import type {
  AgentChatState,
  AgentDirectoryCategory,
  CeremonySharedProps,
  CoreAgentDefinition,
} from "./types";

interface AgentDirectoryPanelProps {
  chat: AgentChatState;
  data: Pick<CeremonySharedProps, "agentStates" | "rosterStates">;
}

const categoryFilters: { id: AgentDirectoryCategory | "all"; label: string }[] = [
  { id: "all", label: "Todos" },
  { id: "basketball_ops", label: "Basketball" },
  { id: "business_ops", label: "Business" },
  { id: "city", label: "Ciudad" },
  { id: "roster", label: "Roster" },
];

export function AgentDirectoryPanel({ chat, data }: AgentDirectoryPanelProps) {
  const [query, setQuery] = useState("");
  const [category, setCategory] = useState<AgentDirectoryCategory | "all">("all");
  const directory = useMemo(() => buildDirectory(data), [data]);
  const filteredAgents = useMemo(
    () => filterDirectory(directory, query, category),
    [category, directory, query],
  );

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    chat.submitMessage();
  }

  return (
    <section className="ceremony-agents-workspace" aria-label="Directorio y chat de agentes">
      <div className="ceremony-panel ceremony-agent-focus">
        <div>
          <p className="eyebrow">Agente seleccionado</p>
          <h2>{chat.selectedAgent.label}</h2>
          <span>{chat.selectedAgent.role}</span>
        </div>
        <p>{chat.selectedAgent.domain}</p>
      </div>

      <div className="ceremony-agent-search">
        <label>
          <span>Buscar agente</span>
          <input
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Nombre, rol o area..."
            type="search"
          />
        </label>
        <div className="ceremony-agent-filters" aria-label="Filtrar agentes">
          {categoryFilters.map((filter) => (
            <button
              key={filter.id}
              type="button"
              className={category === filter.id ? "active" : ""}
              onClick={() => setCategory(filter.id)}
            >
              {filter.label}
            </button>
          ))}
        </div>
      </div>

      <div className="ceremony-agent-list" aria-label="Agentes disponibles">
        {filteredAgents.length === 0 ? (
          <p className="ceremony-agent-list__empty">No hay agentes para ese filtro.</p>
        ) : (
          filteredAgents.map((agent) => {
            const state = data.agentStates[agent.id];
            const rosterState = data.rosterStates[agent.id];
            return (
              <button
                key={agent.id}
                type="button"
                className={chat.selectedAgentId === agent.id ? "ceremony-agent-row active" : "ceremony-agent-row"}
                onClick={() => chat.setSelectedAgentId(agent.id)}
              >
                <span>
                  <strong>{agent.label}</strong>
                  <small>{state?.summary || rosterState?.summary || agent.domain}</small>
                </span>
                <em className={`ceremony-agent__mood mood-${state?.mood ?? rosterState?.emotional_state ?? "idle"}`}>
                  {state?.mood ?? rosterState?.emotional_state ?? "idle"}
                </em>
              </button>
            );
          })
        )}
      </div>

      {chat.selectedAgent.metrics.length > 0 ? (
        <div className="ceremony-agent-metrics">
          {chat.selectedAgent.metrics.map((metric) => (
            <span key={metric.key}>
              {metric.label}
              <strong>{formatAgentMetric(agentMetricValue(data, chat.selectedAgent.id, metric.key))}</strong>
            </span>
          ))}
        </div>
      ) : null}

      <section className="agent-chat-panel agent-chat-panel--focused">
        <div className="ceremony-panel__title">
          <div>
            <p className="eyebrow">Chat directo</p>
            <strong>{chat.selectedAgent.label}</strong>
          </div>
          <span className="agent-chat-panel__status">LLM real</span>
        </div>

        <div className="agent-chat-log" aria-live="polite">
          {chat.activeMessages.length === 0 ? (
            <p className="agent-chat-log__empty">
              Escribi una consulta para abrir una conversacion directa con este agente.
            </p>
          ) : (
            chat.activeMessages.map((message) => (
              <article
                key={message.message_id}
                className={message.sender === "gm" ? "chat-bubble gm" : "chat-bubble agent"}
              >
                <span>{message.sender === "gm" ? "GM" : chat.selectedAgent.label}</span>
                <p>{message.body}</p>
              </article>
            ))
          )}
        </div>

        <form className="agent-chat-form" onSubmit={submit}>
          <textarea
            value={chat.draftMessage}
            onChange={(event) => chat.setDraftMessage(event.target.value)}
            placeholder={`Preguntale a ${chat.selectedAgent.label} por su area...`}
            rows={3}
          />
          <button type="submit" disabled={!chat.draftMessage.trim()}>
            Enviar mensaje
          </button>
        </form>
      </section>
    </section>
  );
}

function agentMetricValue(
  data: Pick<CeremonySharedProps, "agentStates" | "rosterStates">,
  agentId: string,
  metricKey: string,
) {
  const agentMetric = data.agentStates[agentId]?.state[metricKey];
  if (agentMetric !== undefined) {
    return agentMetric;
  }

  const rosterState = data.rosterStates[agentId];
  if (!rosterState) {
    return undefined;
  }

  if (metricKey === "satisfaction") {
    return rosterState.satisfaction;
  }
  if (metricKey === "loyalty") {
    return rosterState.loyalty;
  }

  return undefined;
}

function buildDirectory(
  data: Pick<CeremonySharedProps, "agentStates" | "rosterStates">,
): CoreAgentDefinition[] {
  const rosterAgents = Object.values(data.rosterStates).map((player) => ({
    id: player.player_id,
    label: player.player_id,
    role: "Jugador del roster",
    domain: "rendimiento, rol, vestuario y experiencia personal dentro del roster",
    category: "roster" as const,
    metrics: [
      { key: "satisfaction", label: "Sat." },
      { key: "loyalty", label: "Leal." },
    ],
  }));

  const knownIds = new Set([...canonicalAgents.map((agent) => agent.id), ...rosterAgents.map((agent) => agent.id)]);
  const liveUnknownAgents = Object.values(data.agentStates)
    .filter((state) => !knownIds.has(state.agent_id))
    .map((state) => ({
      id: state.agent_id,
      label: state.agent_id,
      role: "Agente",
      domain: state.summary || "estado general de la franquicia",
      category: "basketball_ops" as const,
      metrics: [],
    }));

  return [...canonicalAgents, ...rosterAgents, ...liveUnknownAgents];
}

function filterDirectory(
  agents: CoreAgentDefinition[],
  query: string,
  category: AgentDirectoryCategory | "all",
) {
  const normalizedQuery = query.trim().toLowerCase();

  return agents.filter((agent) => {
    const matchesCategory = category === "all" || agent.category === category;
    if (!matchesCategory) {
      return false;
    }
    if (!normalizedQuery) {
      return true;
    }

    return [agent.label, agent.role, agent.domain, agent.id]
      .join(" ")
      .toLowerCase()
      .includes(normalizedQuery);
  });
}
