import "./newGame.css";

import { SessionPage } from "./components/SessionPage";
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
      {flow.currentPage === "session" ? (
        <SessionPage
          activeAuthKind={flow.activeAuthKind}
          authenticatingUser={flow.authenticatingUser}
          creatingGuestSession={flow.creatingGuestSession}
          guestToken={flow.guestToken}
          restoringSession={flow.restoringSession}
          status={flow.status}
          userSession={flow.userSession}
          onClearAllAccess={flow.clearAllAccess}
          onCreateGuestSession={() => void flow.createGuestSession()}
          onForgotPassword={flow.forgotPassword}
          onLogin={(email, password) => void flow.login(email, password)}
          onLogoutUser={() => void flow.logoutUser()}
          onRegister={(email, displayName, password) =>
            void flow.register(email, displayName, password)}
          onSwitchToGuestSession={() => void flow.switchToGuestSession()}
        />
      ) : null}

      {flow.currentPage === "home" ? (
        <LandingPage
          activeAuthKind={flow.activeAuthKind}
          games={flow.games}
          restoringSession={flow.restoringSession}
          selectedGame={flow.selectedGame}
          selectedGameId={flow.selectedGameId}
          onContinueSelectedGame={flow.continueSelectedGame}
          onLogoutUser={() => void flow.logoutUser()}
          onSelectGame={flow.setSelectedGameId}
          onStart={flow.startNewGame}
          userSession={flow.userSession}
        />
      ) : null}

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
          agentStates={flow.agentStates}
          cityState={flow.cityState}
          mapState={flow.mapState}
          narrativeInbox={flow.narrativeInbox}
          ownerIntroResponseLabel={flow.ownerIntroResponse?.label ?? null}
          recentResults={flow.recentResults}
          seasonState={flow.seasonState}
          socketStatus={flow.socketStatus}
          status={flow.status}
          timeState={flow.timeState}
          onSetPaused={(paused) => void flow.updateTimeControl({ paused })}
          onSetSpeed={(speed) => void flow.updateTimeControl({ speed, paused: false })}
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
