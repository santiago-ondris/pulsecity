import type { SeasonClientState, SeasonMatchSummary } from "../../../../types";

export function formatSimulatedDate(value: string) {
  const date = new Date(`${value}T00:00:00Z`);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat("es", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    timeZone: "UTC",
  }).format(date);
}

export function formatPointDifferential(season: SeasonClientState) {
  const games = season.wins + season.losses;
  if (games === 0) {
    return "DIF 0.0";
  }

  const differential = (season.points_for - season.points_against) / games;
  const sign = differential > 0 ? "+" : "";
  return `DIF ${sign}${differential.toFixed(1)}`;
}

export function formatMatchScore(result: SeasonMatchSummary) {
  const ownHome = result.home_team_id === "pulsecity";
  const ownScore = ownHome ? result.home_score : result.away_score;
  const opponentScore = ownHome ? result.away_score : result.home_score;
  const venue = ownHome ? "vs" : "@";
  const opponent = ownHome ? result.away_team_id : result.home_team_id;
  return `${ownScore}-${opponentScore} ${venue} ${opponent}`;
}

export function formatAgentMetric(value: number | undefined) {
  if (value === undefined) {
    return "--";
  }

  return value.toFixed(2);
}
