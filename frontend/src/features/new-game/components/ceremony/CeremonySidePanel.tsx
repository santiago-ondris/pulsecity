import { AgentDirectoryPanel } from "./AgentDirectoryPanel";
import { CommandSummaryPanel } from "./CommandSummaryPanel";
import { InboxPanel } from "./InboxPanel";
import { SystemPanel } from "./SystemPanel";
import type { AgentChatState, CeremonySharedProps, CeremonyTab } from "./types";

interface CeremonySidePanelProps {
  activeTab: CeremonyTab;
  chat: AgentChatState;
  data: CeremonySharedProps;
  setActiveTab: (tab: CeremonyTab) => void;
}

const tabs: { id: CeremonyTab; label: string }[] = [
  { id: "overview", label: "Resumen" },
  { id: "inbox", label: "Inbox" },
  { id: "staff", label: "Staff" },
  { id: "system", label: "Sistema" },
];

export function CeremonySidePanel({ activeTab, chat, data, setActiveTab }: CeremonySidePanelProps) {
  return (
    <aside className="ceremony-command-panel" aria-label="Panel de control">
      <div className="ceremony-command-panel__tabs" role="tablist" aria-label="Secciones de la partida">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            type="button"
            className={activeTab === tab.id ? "active" : ""}
            role="tab"
            aria-selected={activeTab === tab.id}
            onClick={() => setActiveTab(tab.id)}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div className="ceremony-command-panel__body">
        {activeTab === "overview" ? (
          <CommandSummaryPanel data={data} onOpenStaff={() => setActiveTab("staff")} />
        ) : null}
        {activeTab === "inbox" ? <InboxPanel events={data.narrativeInbox} /> : null}
        {activeTab === "staff" ? <AgentDirectoryPanel chat={chat} data={data} /> : null}
        {activeTab === "system" ? <SystemPanel events={data.events} mapState={data.mapState} /> : null}
      </div>
    </aside>
  );
}
