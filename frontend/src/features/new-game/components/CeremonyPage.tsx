import { useEffect, useMemo, useState, type CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import type {
  AgentClientStates,
  ChatClientMessages,
  CityClientState,
  FinanceClientState,
  MapClientState,
  NarrativeEvent,
  RealtimeEvent,
  RosterClientStates,
  SeasonClientState,
  SeasonMatchSummary,
  TimeClientState,
} from "../../../types";
import { agentDefinitionFor } from "./ceremony/agents";
import { CommandCenterOverview } from "./ceremony/CommandCenterOverview";
import { CeremonySidePanel } from "./ceremony/CeremonySidePanel";
import { CeremonyTopbar } from "./ceremony/CeremonyTopbar";
import { SeasonKickoffPanel } from "./ceremony/SeasonKickoffPanel";
import { WorldGenerationPanel } from "./ceremony/WorldGenerationPanel";
import type { CeremonySharedProps, CeremonyTab } from "./ceremony/types";
import "./ceremony/commandCenter.css";

interface CeremonyPageProps {
  abbreviation: string;
  cityName: string;
  currentStage: {
    label: string;
    title: string;
    description: string;
  };
  events: RealtimeEvent[];
  financeState: FinanceClientState;
  franchiseName: string;
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
  onStartSeason: () => Promise<boolean>;
  onOpenMedicalCenter: () => void;
  onOpenTradeCenter: () => void;
  onSendAgentChatMessage: (agentId: string, message: string, conversationId?: string) => Promise<string>;
}

export function CeremonyPage(props: CeremonyPageProps) {
  const [activeTab, setActiveTab] = useState<CeremonyTab>("overview");
  const [selectedAgentId, setSelectedAgentId] = useState("owner");
  const [draftMessage, setDraftMessage] = useState("");
  const [kickoffDismissed, setKickoffDismissed] = useState(false);
  const [startingSeason, setStartingSeason] = useState(false);
  const [kickoffError, setKickoffError] = useState<string | null>(null);
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
    financeState: props.financeState,
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
  const gamesPlayed = props.seasonState.wins + props.seasonState.losses;
  const seasonHasAdvanced = props.timeState.days_processed > 0 || gamesPlayed > 0;
  const worldReady = props.mapState.stage === "complete";
  const kickoffEligible = worldReady
    && Boolean(props.ownerIntroResponseLabel)
    && !seasonHasAdvanced;
  const showKickoff = kickoffEligible && !kickoffDismissed;

  useEffect(() => {
    setActiveTab("overview");
    setKickoffDismissed(false);
    setKickoffError(null);
    setStartingSeason(false);
  }, [props.gameId]);

  function submitMessage() {
    const message = draftMessage.trim();
    if (!message) {
      return;
    }

    setDraftMessage("");
    void props.onSendAgentChatMessage(selectedAgentId, message, activeConversationId);
  }

  async function startSeason(): Promise<void> {
    if (startingSeason) {
      return;
    }

    setStartingSeason(true);
    setKickoffError(null);
    setKickoffDismissed(true);
    const started = await props.onStartSeason();
    if (!started) {
      setKickoffDismissed(false);
      setKickoffError("No pudimos iniciar la simulación. La temporada sigue pausada.");
    }
    setStartingSeason(false);
  }

  function openStaff(): void {
    setKickoffDismissed(true);
    setActiveTab("staff");
  }

  return (
    <section
      className="ceremony-builder ceremony-builder--command"
      style={{ "--ceremony-backdrop": `url("${skylineBackdrop}")` } as CSSProperties}
    >
      <div className="ceremony-builder__image" />
      <div className="ceremony-builder__shade" />

      <CeremonyTopbar
        abbreviation={props.abbreviation}
        activeAlerts={props.narrativeInbox.length}
        cityName={props.cityName}
        data={sharedData}
        franchiseName={props.franchiseName}
        mapProgress={props.mapState.progress}
        mode={!worldReady ? "generation" : showKickoff ? "kickoff" : "active"}
        onSetPaused={props.onSetPaused}
        onSetSpeed={props.onSetSpeed}
      />

      {!worldReady ? (
        <WorldGenerationPanel
          cityName={props.cityName}
          cityState={props.cityState}
          currentStage={props.currentStage}
          mapState={props.mapState}
        />
      ) : showKickoff ? (
        <SeasonKickoffPanel
          cityName={props.cityName}
          cityState={props.cityState}
          error={kickoffError}
          financeState={props.financeState}
          franchiseName={props.franchiseName}
          mapState={props.mapState}
          narrativeInbox={props.narrativeInbox}
          ownerIntroResponseLabel={props.ownerIntroResponseLabel}
          rosterCount={Object.keys(props.rosterStates).length}
          seasonState={props.seasonState}
          starting={startingSeason}
          onOpenMedicalCenter={props.onOpenMedicalCenter}
          onOpenStaff={openStaff}
          onOpenTradeCenter={props.onOpenTradeCenter}
          onStartSeason={() => void startSeason()}
        />
      ) : (
        <main className="ceremony-builder__main ceremony-builder__main--command">
          <CommandCenterOverview
            cityName={props.cityName}
            data={sharedData}
            franchiseName={props.franchiseName}
            rosterCount={Object.keys(props.rosterStates).length}
            showStartPrompt={kickoffEligible}
            startingSeason={startingSeason}
            onOpenMedicalCenter={props.onOpenMedicalCenter}
            onOpenStaff={openStaff}
            onOpenTradeCenter={props.onOpenTradeCenter}
            onStartSeason={() => void startSeason()}
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
      )}
    </section>
  );
}
