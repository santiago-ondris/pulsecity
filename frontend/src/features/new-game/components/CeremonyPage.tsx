import { useMemo, useState, type CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import type {
  AgentClientStates,
  ChatClientMessages,
  CityClientState,
  MapClientState,
  NarrativeEvent,
  RealtimeEvent,
  RosterClientStates,
  SeasonClientState,
  SeasonMatchSummary,
  TimeClientState,
} from "../../../types";
import { agentDefinitionFor } from "./ceremony/agents";
import { CeremonyMapPanel } from "./ceremony/CeremonyMapPanel";
import { CeremonySidePanel } from "./ceremony/CeremonySidePanel";
import { CeremonyTopbar } from "./ceremony/CeremonyTopbar";
import type { CeremonySharedProps, CeremonyTab } from "./ceremony/types";

interface CeremonyPageProps {
  currentStage: {
    label: string;
    title: string;
    description: string;
  };
  events: RealtimeEvent[];
  gameId: string;
  agentStates: AgentClientStates;
  chatMessages: ChatClientMessages;
  cityState: CityClientState;
  mapState: MapClientState;
  narrativeInbox: NarrativeEvent[];
  ownerIntroResponseLabel: string | null;
  recentResults: SeasonMatchSummary[];
  rosterStates: RosterClientStates;
  seasonState: SeasonClientState;
  socketStatus: string;
  status: string;
  timeState: TimeClientState;
  onSetPaused: (paused: boolean) => void;
  onSetSpeed: (speed: 1 | 5 | 20) => void;
  onSendAgentChatMessage: (agentId: string, message: string, conversationId?: string) => Promise<string>;
}

export function CeremonyPage(props: CeremonyPageProps) {
  const [activeTab, setActiveTab] = useState<CeremonyTab>("agents");
  const [selectedAgentId, setSelectedAgentId] = useState("owner");
  const [draftMessage, setDraftMessage] = useState("");
  const activeConversationId = `chat-local-${selectedAgentId}`;
  const activeMessages = props.chatMessages[activeConversationId] ?? [];
  const selectedAgent = useMemo(
    () => agentDefinitionFor(selectedAgentId, props.agentStates, props.rosterStates),
    [props.agentStates, props.rosterStates, selectedAgentId],
  );
  const sharedData: CeremonySharedProps = {
    agentStates: props.agentStates,
    cityState: props.cityState,
    currentStage: props.currentStage,
    events: props.events,
    gameId: props.gameId,
    mapState: props.mapState,
    narrativeInbox: props.narrativeInbox,
    ownerIntroResponseLabel: props.ownerIntroResponseLabel,
    recentResults: props.recentResults,
    rosterStates: props.rosterStates,
    seasonState: props.seasonState,
    socketStatus: props.socketStatus,
    status: props.status,
    timeState: props.timeState,
  };

  function submitMessage() {
    const message = draftMessage.trim();
    if (!message) {
      return;
    }

    setDraftMessage("");
    void props.onSendAgentChatMessage(selectedAgentId, message, activeConversationId);
  }

  return (
    <section
      className="ceremony-builder ceremony-builder--command"
      style={{ "--ceremony-backdrop": `url("${skylineBackdrop}")` } as CSSProperties}
    >
      <div className="ceremony-builder__image" />
      <div className="ceremony-builder__shade" />

      <CeremonyTopbar
        data={sharedData}
        onSetPaused={props.onSetPaused}
        onSetSpeed={props.onSetSpeed}
      />

      <main className="ceremony-builder__main ceremony-builder__main--command">
        <CeremonyMapPanel
          cityState={props.cityState}
          currentStage={props.currentStage}
          mapState={props.mapState}
        />
        <CeremonySidePanel
          activeTab={activeTab}
          chat={{
            activeConversationId,
            activeMessages,
            draftMessage,
            selectedAgent,
            selectedAgentId,
            setDraftMessage,
            setSelectedAgentId,
            submitMessage,
          }}
          data={sharedData}
          setActiveTab={setActiveTab}
        />
      </main>
    </section>
  );
}
