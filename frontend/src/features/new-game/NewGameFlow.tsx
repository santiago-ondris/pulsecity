import "./newGame.css";

import { LandingPage } from "./components/LandingPage";
import { IdentityPage } from "./components/IdentityPage";
import { ScenarioPage } from "./components/ScenarioPage";
import { ManagementPage } from "./components/ManagementPage";
import { LaunchPage } from "./components/LaunchPage";
import { CeremonyPage } from "./components/CeremonyPage";
import { OwnerIntroModal } from "./components/OwnerIntroModal";
import { useNewGameFlow } from "./hooks/useNewGameFlow";

export function NewGameFlow() {
  const flow = useNewGameFlow();

  return (
    <main className="new-game-shell">
      {flow.currentPage === "home" ? <LandingPage onStart={flow.startNewGame} /> : null}

      {flow.currentPage === "identity" ? (
        <IdentityPage
          accentColor={flow.draft.accentColor}
          abbreviation={flow.draft.abbreviation}
          cityName={flow.draft.cityName}
          franchiseName={flow.draft.franchiseName}
          primaryColor={flow.draft.primaryColor}
          secondaryColor={flow.draft.secondaryColor}
          onAbbreviationChange={(value) => flow.updateDraft("abbreviation", value)}
          onAccentColorChange={(value) => flow.updateDraft("accentColor", value)}
          onCityNameChange={(value) => flow.updateDraft("cityName", value)}
          onContinue={flow.completeIdentityStep}
          onFranchiseNameChange={(value) => flow.updateDraft("franchiseName", value)}
          onPrimaryColorChange={(value) => flow.updateDraft("primaryColor", value)}
          onSecondaryColorChange={(value) => flow.updateDraft("secondaryColor", value)}
        />
      ) : null}

      {flow.currentPage === "scenario" ? (
        <ScenarioPage
          selectedScenario={flow.draft.selectedScenario}
          onBack={flow.goBack}
          onContinue={flow.completeScenarioStep}
          onSelect={(value) => flow.updateDraft("selectedScenario", value)}
        />
      ) : null}

      {flow.currentPage === "management" ? (
        <ManagementPage
          cityManagementMode={flow.draft.cityManagementMode}
          onBack={flow.goBack}
          onContinue={flow.completeManagementStep}
          onSelect={(value) => flow.updateDraft("cityManagementMode", value)}
        />
      ) : null}

      {flow.currentPage === "launch" ? (
        <LaunchPage
          creatingGame={flow.creatingGame}
          draft={flow.draft}
          ownerIntroResponseLabel={flow.ownerIntroResponse?.label ?? null}
          status={flow.status}
          onBack={flow.goBack}
          onLaunch={flow.createGame}
        />
      ) : null}

      {flow.currentPage === "ceremony" ? (
        <CeremonyPage
          currentStage={flow.currentStage}
          events={flow.events}
          gameId={flow.gameId}
          mapState={flow.mapState}
          ownerIntroResponseLabel={flow.ownerIntroResponse?.label ?? null}
          socketStatus={flow.socketStatus}
          status={flow.status}
        />
      ) : null}

      {flow.activeNarrativeEvent ? (
        <OwnerIntroModal
          event={flow.activeNarrativeEvent}
          submitting={flow.submittingNarrativeChoice}
          onSelect={(choice) => void flow.submitOwnerIntroChoice(choice)}
        />
      ) : null}
    </main>
  );
}
