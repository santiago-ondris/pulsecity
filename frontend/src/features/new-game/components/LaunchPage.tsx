import { managementModeLabel, scenarioById } from "../helpers";
import type { NewGameDraft } from "../types";
import { FranchisePreview, StatusBadge } from "./common";

interface LaunchPageProps {
  creatingGame: boolean;
  draft: NewGameDraft;
  ownerIntroResponseLabel: string | null;
  status: string;
  onBack: () => void;
  onLaunch: () => void;
}

export function LaunchPage(props: LaunchPageProps) {
  const scenario = scenarioById(props.draft.selectedScenario);

  return (
    <section className="screen step-screen">
      <div className="step-copy">
        <p className="eyebrow">Paso 4 de 4</p>
        <h1>Última página antes de fundar el mundo.</h1>
        <p>
          Ya no estás editando campos. Solo revisás la lectura completa de la partida y confirmás
          el inicio.
        </p>
      </div>

      <div className="step-layout">
        <article className="step-card">
          <div className="launch-dossier">
            <div className="dossier-row">
              <span>Franquicia</span>
              <strong>
                {props.draft.cityName} {props.draft.franchiseName}
              </strong>
            </div>
            <div className="dossier-row">
              <span>Sigla</span>
              <strong>{props.draft.abbreviation}</strong>
            </div>
            <div className="dossier-row">
              <span>Escenario</span>
              <strong>{scenario.label}</strong>
            </div>
            <div className="dossier-row">
              <span>Presión inicial</span>
              <strong>{scenario.pressure}</strong>
            </div>
            <div className="dossier-row">
              <span>Gestión de ciudad</span>
              <strong>{managementModeLabel(props.draft.cityManagementMode)}</strong>
            </div>
            <div className="dossier-row">
              <span>Dirección del Owner</span>
              <strong>{props.ownerIntroResponseLabel ?? "Se define después de la fundación"}</strong>
            </div>
          </div>

          <div className="launch-status">
            <StatusBadge label={props.status} tone="info" />
          </div>

          <div className="page-actions split-actions">
            <button type="button" className="secondary-action" onClick={props.onBack}>
              Volver
            </button>
            <button
              type="button"
              className="primary-action"
              onClick={props.onLaunch}
              disabled={props.creatingGame}
            >
              {props.creatingGame ? "Fundando..." : "Confirmar franquicia y generar mapa"}
            </button>
          </div>
        </article>

        <FranchisePreview
          accentColor={props.draft.accentColor}
          abbreviation={props.draft.abbreviation}
          cityName={props.draft.cityName}
          contextLabel="Dossier final"
          franchiseName={props.draft.franchiseName}
          primaryColor={props.draft.primaryColor}
          secondaryColor={props.draft.secondaryColor}
        />
      </div>
    </section>
  );
}
