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

export type MapEvent = MapSnapshotEnvelope | MapPatchEnvelope;

export interface GameSetup {
  game_id: string;
  city_name: string;
  franchise_name: string;
  abbreviation: string;
  primary_color: string;
  secondary_color: string;
  accent_color: string;
  initial_scenario: string;
  status: string;
  created_at?: string;
  updated_at?: string;
}
