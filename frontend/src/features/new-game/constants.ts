export const stageSequence = ["terrain", "zoning", "stadium", "complete"] as const;

export const stageMeta: Record<
  string,
  {
    label: string;
    title: string;
    description: string;
  }
> = {
  idle: {
    label: "En espera",
    title: "El mundo todavia no fue fundado.",
    description: "La ceremonia empieza solo cuando la franquicia queda confirmada.",
  },
  terrain: {
    label: "Terreno",
    title: "Primero aparece la geografia.",
    description: "Costa, relieve y base territorial emergen antes de cualquier lectura urbana.",
  },
  zoning: {
    label: "Zonificacion",
    title: "Los distritos toman forma.",
    description: "La ciudad empieza a especializar barrios y declarar vocaciones de suelo.",
  },
  stadium: {
    label: "Estadio",
    title: "La franquicia encuentra su centro.",
    description: "El estadio fija el punto emocional alrededor del cual se organiza el mapa.",
  },
  complete: {
    label: "Completo",
    title: "El mundo quedo listo.",
    description: "La base de la partida ya existe y el cliente solo tiene que seguir sus deltas.",
  },
};

export const initialScenarios = [
  {
    id: "rebuild",
    label: "Reconstruccion",
    roster: "Joven, rating bajo, alto potencial.",
    pressure: "Dueño paciente, presión mediática baja.",
    city: "Ciudad pequeña, economía modesta y fanbase leal.",
  },
  {
    id: "contention",
    label: "Ventana de contencion",
    roster: "2-3 estrellas, veteranos y contratos grandes.",
    pressure: "Dueño exigente, playoffs este año.",
    city: "Ciudad desarrollada, estadio lleno y fanbase caliente.",
  },
  {
    id: "decline",
    label: "Historica en declive",
    roster: "Viejas glorias mezcladas con jovenes sin rumbo.",
    pressure: "Dueño nostalgico, comparación constante con el pasado.",
    city: "Ciudad grande, medios encima y frustración alta.",
  },
  {
    id: "expansion",
    label: "Expansion pura",
    roster: "Draft de expansion, sin compromisos heredados.",
    pressure: "Dueño visionario, horizonte largo y paciencia total.",
    city: "Ciudad virgen, todo por construir desde cero.",
  },
] as const;

export type ScenarioId = (typeof initialScenarios)[number]["id"];

export const colorPresets = [
  { label: "Signal Green", value: "#00C896" },
  { label: "Steel Blue", value: "#7B8CDE" },
  { label: "Arena Gold", value: "#FFAA00" },
  { label: "Burnt Orange", value: "#FF6B2B" },
  { label: "Infra Red", value: "#E05555" },
  { label: "Night Glass", value: "#1A1A1E" },
] as const;

export const cityManagementModes = [
  {
    id: "owner_influence",
    label: "Dueño con influencia",
    description:
      "Sos el GM de la franquicia. Tu poder sobre la ciudad es indirecto: lobby, financiamiento de proyectos y presión política.",
    impact:
      "El alcalde tiene agenda propia y la ciudad reacciona como un organismo independiente.",
  },
  {
    id: "dual_figure",
    label: "Figura dual",
    description:
      "Controlás tanto la franquicia como la ciudad directamente. Dos sombreros, control total.",
    impact:
      "La ciudad queda bajo tu manejo directo y el basket convive con ese control urbano.",
  },
] as const;

export type CityManagementModeId = (typeof cityManagementModes)[number]["id"];

export const flowPages = [
  "session",
  "home",
  "identity",
  "scenario",
  "management",
  "launch",
  "ceremony",
  "trade-center",
  "medical-center",
] as const;

export type FlowPage = (typeof flowPages)[number];

export const pagePaths: Record<FlowPage, string> = {
  session: "/session",
  home: "/",
  identity: "/new-game/identity",
  scenario: "/new-game/scenario",
  management: "/new-game/management",
  launch: "/new-game/review",
  ceremony: "/new-game/ceremony",
  "trade-center": "/franchise/trades",
  "medical-center": "/franchise/medical",
};

export const initialDraft = {
  cityName: "Nueva Aurora",
  franchiseName: "Lighthouses",
  abbreviation: "NAR",
  primaryColor: "#00C896",
  secondaryColor: "#7B8CDE",
  accentColor: "#FFAA00",
  selectedScenario: "expansion" as ScenarioId,
  cityManagementMode: "owner_influence" as CityManagementModeId,
};
