import type {
  AgentClientStates,
  ChatMessageEvent,
  CityClientState,
  FinanceClientState,
  MapClientState,
  NarrativeEvent,
  RealtimeEvent,
  RosterClientStates,
  SeasonClientState,
  SeasonMatchSummary,
  TimeClientState,
} from "../../../../types";

export type CeremonyTab = "agents" | "inbox" | "season" | "system";
export type AgentDirectoryCategory = "basketball_ops" | "business_ops" | "city" | "press" | "roster";

export interface CoreAgentDefinition {
  id: string;
  label: string;
  role: string;
  domain: string;
  category: AgentDirectoryCategory;
  metrics: {
    key: string;
    label: string;
  }[];
}

export interface CeremonySharedProps {
  agentStates: AgentClientStates;
  cityState: CityClientState;
  currentStage: {
    label: string;
    title: string;
    description: string;
  };
  events: RealtimeEvent[];
  financeState: FinanceClientState;
  gameId: string;
  mapState: MapClientState;
  narrativeInbox: NarrativeEvent[];
  ownerIntroResponseLabel: string | null;
  recentResults: SeasonMatchSummary[];
  rosterStates: RosterClientStates;
  seasonState: SeasonClientState;
  socketStatus: string;
  status: string;
  timeState: TimeClientState;
}

export interface AgentChatState {
  activeConversationId: string;
  activeMessages: ChatMessageEvent[];
  draftMessage: string;
  selectedAgent: CoreAgentDefinition;
  selectedAgentId: string;
  setDraftMessage: (value: string) => void;
  setSelectedAgentId: (value: string) => void;
  submitMessage: () => void;
}
