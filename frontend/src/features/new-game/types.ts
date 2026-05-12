import type { CityManagementModeId, FlowPage, ScenarioId } from "./constants";

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
