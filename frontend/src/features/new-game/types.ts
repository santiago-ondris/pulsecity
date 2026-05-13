import type { CityManagementModeId, FlowPage, ScenarioId } from "./constants";
import type { GameSummary } from "../../types";

export interface NewGameDraft {
  cityName: string;
  franchiseName: string;
  abbreviation: string;
  primaryColor: string;
  secondaryColor: string;
  accentColor: string;
  selectedScenario: ScenarioId;
  cityManagementMode: CityManagementModeId;
}

export interface FlowProgress {
  currentPage: FlowPage;
  unlockedPage: FlowPage;
}

export interface GuestAccessState {
  guestToken: string;
  games: GameSummary[];
}
