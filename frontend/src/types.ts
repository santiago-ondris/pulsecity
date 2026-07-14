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

export interface TimeClientState {
  simulated_date: string;
  speed: 1 | 5 | 20;
  paused: boolean;
  days_processed: number;
}

export interface TimeStatePatch {
  simulated_date?: string;
  speed?: 1 | 5 | 20;
  paused?: boolean;
  days_processed?: number;
}

export interface TimePatchEnvelope {
  type: "time.patch";
  subject: string;
  game_id: string;
  patch: TimeStatePatch;
}

export interface SeasonClientState {
  wins: number;
  losses: number;
  points_for: number;
  points_against: number;
  last_result?: SeasonMatchSummary;
}

export interface SeasonMatchSummary {
  match_id: string;
  simulated_date: string;
  home_team_id: string;
  away_team_id: string;
  home_score: number;
  away_score: number;
  winner_team_id: string;
}

export interface SeasonStatePatch {
  wins?: number;
  losses?: number;
  points_for?: number;
  points_against?: number;
  last_result?: SeasonMatchSummary;
}

export interface SeasonPatchEnvelope {
  type: "season.patch";
  subject: string;
  game_id: string;
  patch: SeasonStatePatch;
}

export interface CityClientState {
  fan_sentiment: number;
  ticket_sales_index: number;
  local_economy_index: number;
  stadium_district_land_value: number;
  win_streak: number;
  loss_streak: number;
  last_match_id?: string;
  reason?: string;
}

export interface CityStatePatch {
  fan_sentiment?: number;
  ticket_sales_index?: number;
  local_economy_index?: number;
  stadium_district_land_value?: number;
  win_streak?: number;
  loss_streak?: number;
  last_match_id?: string;
  reason?: string;
}

export interface CityPatchEnvelope {
  type: "city.patch";
  subject: string;
  game_id: string;
  patch: CityStatePatch;
}

export interface AgentClientState {
  agent_id: string;
  mood: string;
  state: Record<string, number>;
  summary: string;
  simulated_date?: string;
  source_event_id?: string;
  source_subject?: string;
}

export type AgentClientStates = Record<string, AgentClientState>;

export interface AgentStatePatch {
  mood?: string;
  state?: Record<string, number>;
  summary?: string;
  simulated_date?: string;
  source_event_id?: string;
  source_subject?: string;
}

export interface AgentPatchEnvelope {
  type: "agent.patch";
  subject: string;
  game_id: string;
  agent_id: string;
  patch: AgentStatePatch;
}

export interface PlayerEmotionalState {
  player_id: string;
  emotional_state: string;
  satisfaction: number;
  loyalty: number;
  ego: number;
  competitive_drive: number;
  city_connection: number;
  summary: string;
  availability?: "active" | "injured";
  injury_id?: string;
  severity?: "minor" | "moderate" | "major";
  expected_recovery_date?: string;
  estimated_days_out?: number;
  availability_changed_on?: string;
  simulated_date?: string;
  source_event_id?: string;
  source_subject?: string;
}

export type RosterClientStates = Record<string, PlayerEmotionalState>;

export interface RosterPatchEnvelope {
  type: "roster.patch";
  subject: string;
  game_id: string;
  patch: {
    simulated_date: string;
    source_event_id: string;
    source_subject: string;
    players: PlayerEmotionalState[];
  };
}

export interface RelationshipClientState {
  relationship_id: string;
  agent_a_id: string;
  agent_b_id: string;
  trust: number;
  trend: string;
  last_event: string;
  short_history: string[];
  simulated_date?: string;
  source_event_id?: string;
  source_subject?: string;
}

export type RelationshipClientStates = Record<string, RelationshipClientState>;

export interface RelationshipPatch {
  agent_a_id: string;
  agent_b_id: string;
  trust: number;
  trend: string;
  last_event: string;
  short_history: string[];
}

export interface RelationsPatchEnvelope {
  type: "relations.patch";
  subject: string;
  game_id: string;
  patch: {
    simulated_date: string;
    source_event_id: string;
    source_subject: string;
    relationships: RelationshipPatch[];
  };
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

export interface ChatMessageEvent {
  type: "chat.message";
  subject: string;
  game_id: string;
  conversation_id: string;
  message_id: string;
  agent_id: string;
  sender: "gm" | "agent";
  body: string;
  metadata?: Record<string, string>;
  created_at: string;
}

export type ChatClientMessages = Record<string, ChatMessageEvent[]>;

export type RealtimeEvent =
  | MapSnapshotEnvelope
  | MapPatchEnvelope
  | TimePatchEnvelope
  | SeasonPatchEnvelope
  | CityPatchEnvelope
  | AgentPatchEnvelope
  | RosterPatchEnvelope
  | RelationsPatchEnvelope
  | NarrativeEvent
  | NarrativeResponseEvent
  | ChatMessageEvent;

export interface GuestSession {
  guest_token: string;
  created_at?: string;
  last_seen_at?: string;
}

export interface User {
  user_id: string;
  email: string;
  display_name: string;
  created_at?: string;
}

export interface UserSession {
  session_token: string;
  user: User;
  created_at?: string;
  last_seen_at?: string;
}

export interface GuestUpgradeResult {
  user_session: UserSession;
  migrated_games: number;
  guest_token_used: string;
}

export interface GameSummary {
  game_id: string;
  city_name: string;
  franchise_name: string;
  owner_kind: "guest" | "user";
  initial_scenario: string;
  city_management_mode: string;
  status: string;
  updated_at: string;
}

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
