export type Terrain = "water" | "plain" | "forest" | "hill";
export type Zone = "residential" | "commercial" | "industrial" | "park";

export interface MapCell {
  terrain: Terrain;
  zone?: Zone;
}

export interface MapData {
  width: number;
  height: number;
  cells: MapCell[][];
}

export interface GridPoint {
  x: number;
  y: number;
}

export interface MapClientState {
  game_id: string;
  stage: string;
  progress: number;
  message: string;
  map_data?: MapData;
  stadium?: GridPoint;
}

export interface MapSnapshotEnvelope {
  type: "map.snapshot";
  subject: string;
  state: MapClientState;
}

export interface MapStatePatch {
  stage?: string;
  progress?: number;
  message?: string;
  map_data?: MapData;
  stadium?: GridPoint;
}

export interface MapPatchEnvelope {
  type: "map.patch";
  subject: string;
  game_id: string;
  patch: MapStatePatch;
}

export interface NarrativeChoice {
  id: string;
  label: string;
}

export interface NarrativeEvent {
  event_id: string;
  game_id: string;
  type: "narrative.event";
  subject: string;
  emitter: string;
  kind: string;
  urgency: string;
  title: string;
  body: string;
  metadata?: Record<string, string>;
  choices?: NarrativeChoice[];
}

export interface NarrativeResponseEvent {
  type: "narrative.response";
  subject: string;
  game_id: string;
  event_id: string;
  choice: NarrativeChoice;
  emitter: string;
  metadata?: Record<string, string>;
  timestamp: string;
}

export type RealtimeEvent =
  | MapSnapshotEnvelope
  | MapPatchEnvelope
  | NarrativeEvent
  | NarrativeResponseEvent;

export interface GameSetup {
  game_id: string;
  city_name: string;
  franchise_name: string;
  abbreviation: string;
  primary_color: string;
  secondary_color: string;
  accent_color: string;
  initial_scenario: string;
  city_management_mode: string;
  owner_intro_event?: NarrativeEvent;
  owner_intro_response?: NarrativeChoice;
  status: string;
  created_at?: string;
  updated_at?: string;
}
