import type { AgentDirectoryCategory, CoreAgentDefinition } from "./types";
import type { AgentClientStates, RosterClientStates } from "../../../../types";

export const canonicalAgents: CoreAgentDefinition[] = [
  agent("owner", "Owner", "Propietario", "ownership, franchise vision, pressure, spending", "basketball_ops", [
    ["sporting_trust", "Dep."],
    ["patience_remaining", "Pac."],
  ]),
  agent(
    "president_basketball_ops",
    "President of Basketball Ops",
    "President of Basketball Operations",
    "front office alignment, roster philosophy, owner buffer",
    "basketball_ops",
  ),
  agent("assistant_gm_cap", "Assistant GM, Cap", "Assistant General Manager", "contracts, cap sheets, trade mechanics", "basketball_ops"),
  agent(
    "assistant_gm_personnel",
    "Assistant GM, Personnel",
    "Assistant General Manager",
    "personnel evaluation, market reads, player agents",
    "basketball_ops",
  ),
  agent(
    "assistant_gm_operations",
    "Assistant GM, Ops",
    "Assistant General Manager",
    "front office operations, staff coordination, internal process",
    "basketball_ops",
  ),
  agent("scouting_director", "Director de Scouting", "Director de Scouting", "draft, scouting reports, talent projection", "basketball_ops", [
    ["criteria_trust", "Crit."],
    ["motivation", "Mot."],
  ]),
  agent(
    "player_personnel_director",
    "Director de Player Personnel",
    "Director de Player Personnel",
    "roster balance, player agents, trade opportunities",
    "basketball_ops",
  ),
  agent("head_analytics", "Head of Analytics", "Head of Analytics", "models, lineup data, projections", "basketball_ops"),
  agent("head_coach", "Head Coach", "Entrenadora principal", "rotation, tactics, locker room leadership", "basketball_ops", [
    ["roster_satisfaction", "Roster"],
    ["results_pressure", "Pres."],
  ]),
  agent("assistant_coach_offense", "Assistant Coach, Offense", "Assistant Coach", "offensive sets, shot profile, player communication", "basketball_ops"),
  agent("assistant_coach_defense", "Assistant Coach, Defense", "Assistant Coach", "defensive system, matchup prep, effort standards", "basketball_ops"),
  agent(
    "player_development_director",
    "Director de Player Development",
    "Director de Player Development",
    "young players, skill plans, long-term growth",
    "basketball_ops",
  ),
  agent("team_doctor", "Medico del Equipo", "Medico del Equipo", "injury diagnosis, return protocol, health risk", "basketball_ops"),
  agent(
    "strength_conditioning_coach",
    "Strength & Conditioning",
    "Strength & Conditioning Coach",
    "load management, prevention, conditioning",
    "basketball_ops",
  ),
  agent("sports_psychologist", "Sports Psychologist", "Psicologia deportiva", "emotional climate, burnout, player trust", "basketball_ops", [
    ["locker_room_climate", "Clima"],
    ["emotional_alert", "Alerta"],
  ]),
  agent("video_coordinator", "Video Coordinator", "Video Coordinator", "film, scouting clips, opponent tendencies", "basketball_ops"),
  agent("international_scout", "International Scout", "International Scout", "international scouting, cultural fit, overseas markets", "basketball_ops"),
  agent(
    "ceo_business_ops",
    "CEO Business Ops",
    "CEO / President of Business Operations",
    "business strategy, fan experience, city relationships",
    "business_ops",
  ),
  agent("cfo", "CFO", "Finanzas", "budget, salary cap, financial risk", "business_ops", [
    ["financial_trust", "Fin."],
    ["budget_alert", "Alerta"],
  ]),
  agent("marketing_director", "Marketing & Brand", "Director de Marketing & Brand", "brand, fanbase, campaigns, marketability", "business_ops"),
  agent("ticket_sales_director", "Ticket Sales", "Director de Ticket Sales", "attendance, pricing, season tickets", "business_ops"),
  agent(
    "partnerships_director",
    "Corporate Partnerships",
    "Director de Corporate Partnerships & Sponsors",
    "sponsors, corporate relationships, activations",
    "business_ops",
  ),
  agent("pr_director", "PR & Communications", "Director de PR & Communications", "media strategy, crisis response, public narrative", "business_ops"),
  agent("arena_operations_director", "Arena Operations", "Director de Arena Operations", "arena logistics, events, maintenance, fan experience", "business_ops"),
  agent("legal_counsel", "Legal Counsel", "Legal Counsel", "contracts, legal risk, regulatory questions", "business_ops"),
  agent("mayor", "Alcalde", "Alcalde", "politics, permits, city agenda", "city"),
  agent("police_chief", "Jefe de Policia", "Jefe de Policia", "stadium security, logistics, public safety", "city"),
  agent(
    "chamber_commerce_president",
    "Camara de Comercio",
    "Presidente de la Camara de Comercio",
    "local business, sponsors, district economy",
    "city",
  ),
  agent("urbanism_director", "Director de Urbanismo", "Director de Urbanismo", "permits, planning, zoning process", "city"),
  agent("press", "La Prensa", "Agente colectivo", "coverage, public sentiment, dominant narrative", "press"),
];

export const coreAgentOrder = canonicalAgents.filter((agent) =>
  ["owner", "head_coach", "cfo", "scouting_director", "sports_psychologist"].includes(agent.id),
);

export function agentDefinitionFor(
  agentId: string,
  agentStates: AgentClientStates,
  rosterStates: RosterClientStates,
): CoreAgentDefinition {
  const canonical = canonicalAgents.find((agent) => agent.id === agentId);
  if (canonical) {
    return canonical;
  }

  const player = rosterStates[agentId];
  if (player) {
    return {
      id: player.player_id,
      label: player.full_name ?? player.player_id,
      role: player.position ? `Jugador del roster · ${player.position}` : "Jugador del roster",
      domain: "rendimiento, rol, vestuario y experiencia personal dentro del roster",
      category: "roster",
      metrics: [
        { key: "satisfaction", label: "Sat." },
        { key: "loyalty", label: "Leal." },
      ],
    };
  }

  const liveState = agentStates[agentId];
  if (liveState) {
    return {
      id: liveState.agent_id,
      label: liveState.agent_id,
      role: "Agente",
      domain: liveState.summary || "estado general de la franquicia",
      category: "basketball_ops",
      metrics: [],
    };
  }

  return canonicalAgents[0];
}

function agent(
  id: string,
  label: string,
  role: string,
  domain: string,
  category: AgentDirectoryCategory,
  metrics: [string, string][] = [],
): CoreAgentDefinition {
  return {
    id,
    label,
    role,
    domain,
    category,
    metrics: metrics.map(([key, metricLabel]) => ({ key, label: metricLabel })),
  };
}
