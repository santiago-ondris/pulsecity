import type { CSSProperties } from "react";

import type { MapCell, RealtimeEvent } from "../../types";
import { cityManagementModes, flowPages, initialScenarios, pagePaths, type CityManagementModeId, type FlowPage } from "./constants";

export function formatGameId(value: string) {
  if (!value) {
    return "sin partida";
  }

  if (value.length <= 18) {
    return value;
  }

  return `${value.slice(0, 8)}...${value.slice(-6)}`;
}

export function safeAbbreviation(value: string) {
  const trimmed = value.trim().toUpperCase();
  return trimmed.length > 0 ? trimmed : "NEW";
}

export function summarizeTerrain(cells: MapCell[][]) {
  const flat = cells.flat();
  if (flat.length === 0) {
    return { water: 0, forest: 0, plain: 0, hill: 0 };
  }

  const counts = flat.reduce(
    (acc, cell) => {
      acc[cell.terrain] += 1;
      return acc;
    },
    { water: 0, forest: 0, plain: 0, hill: 0 },
  );

  return {
    water: Math.round((counts.water / flat.length) * 100),
    forest: Math.round((counts.forest / flat.length) * 100),
    plain: Math.round((counts.plain / flat.length) * 100),
    hill: Math.round((counts.hill / flat.length) * 100),
  };
}

export function buildCellClassName({
  cell,
  showZones,
  showStadium,
}: {
  cell: MapCell;
  showZones: boolean;
  showStadium: boolean;
}) {
  return [
    "cell",
    `terrain-${cell.terrain}`,
    showZones && cell.zone ? `zone-${cell.zone}` : "",
    showStadium ? "stadium" : "",
  ]
    .filter(Boolean)
    .join(" ");
}

export function gridColumns(width: number): CSSProperties {
  return {
    gridTemplateColumns: `repeat(${width}, 1fr)`,
  };
}

export function managementModeLabel(value: CityManagementModeId | string) {
  return cityManagementModes.find((mode) => mode.id === value)?.label ?? "Dueño con influencia";
}

export function describeRealtimeEvent(event: RealtimeEvent) {
  if (event.type === "narrative.event") {
    return event.title;
  }

  if (event.type === "narrative.response") {
    return event.choice.label;
  }

  if (event.type === "map.snapshot") {
    return event.state.message;
  }

  if (event.type === "time.patch") {
    return event.patch.simulated_date
      ? `Fecha simulada ${event.patch.simulated_date}`
      : "Control de tiempo actualizado";
  }

  if (event.type === "season.patch") {
    return `${event.patch.wins ?? 0}-${event.patch.losses ?? 0}`;
  }

  if (event.type === "city.patch") {
    return event.patch.reason ?? "Pulso urbano actualizado";
  }

  return event.patch.message ?? "Patch recibido";
}

export function pageIndex(page: FlowPage) {
  return flowPages.indexOf(page);
}

export function clampPage(requested: FlowPage, unlocked: FlowPage) {
  return pageIndex(requested) <= pageIndex(unlocked) ? requested : unlocked;
}

export function pageFromPath(pathname: string): FlowPage {
  const found = (Object.entries(pagePaths) as Array<[FlowPage, string]>).find(
    ([, value]) => value === pathname,
  );

  return found?.[0] ?? "home";
}

export function scenarioById(id: string) {
  return initialScenarios.find((scenario) => scenario.id === id) ?? initialScenarios[0];
}
