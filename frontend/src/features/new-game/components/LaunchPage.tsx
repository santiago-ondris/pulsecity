import type { CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import { managementModeLabel, scenarioById } from "../helpers";
import type { NewGameDraft } from "../types";

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
  const fullName = `${props.draft.cityName} ${props.draft.franchiseName}`;

  const launchStyle = {
    "--launch-backdrop": `url("${skylineBackdrop}")`,
    "--franchise-primary": props.draft.primaryColor,
    "--franchise-secondary": props.draft.secondaryColor,
    "--franchise-accent": props.draft.accentColor,
  } as CSSProperties;

  return (
    <section className="launch-builder" style={launchStyle}>
      <div className="launch-builder__image" />
      <div className="launch-builder__shade" />

      <header className="launch-builder__topbar">
        <div>
          <p className="eyebrow">Paso 4 de 4</p>
          <strong>Revision final</strong>
        </div>
        <span>{safeAbbreviation(props.draft.abbreviation)}</span>
      </header>

      <main className="launch-builder__main">
        <section className="launch-identity" aria-label="Franquicia a fundar">
          <div className="launch-identity__crest">
            <span>{safeAbbreviation(props.draft.abbreviation)}</span>
          </div>
          <div className="launch-identity__name">
            <p>{props.draft.cityName}</p>
            <h1>{props.draft.franchiseName}</h1>
          </div>
          <div className="launch-identity__palette">
            <span style={{ background: props.draft.primaryColor }} />
            <span style={{ background: props.draft.secondaryColor }} />
            <span style={{ background: props.draft.accentColor }} />
          </div>
        </section>

        <section className="launch-confirmation" aria-label="Confirmar fundacion">
          <div className="launch-confirmation__heading">
            <p className="eyebrow">Listo para fundar</p>
            <h2>{fullName}</h2>
          </div>

          <div className="launch-summary-grid">
            <SummaryItem label="Escenario" value={scenario.label} />
            <SummaryItem label="Presion inicial" value={scenario.pressure} />
            <SummaryItem label="Ciudad" value={scenario.city} />
            <SummaryItem label="Gobierno" value={managementModeLabel(props.draft.cityManagementMode)} />
            <SummaryItem
              label="Owner"
              value={props.ownerIntroResponseLabel ?? "Se define despues de la fundacion"}
            />
            <SummaryItem label="Sistema" value={props.status} />
          </div>

          <div className="launch-actions">
            <button type="button" className="secondary-action" onClick={props.onBack}>
              Volver
            </button>
            <button
              type="button"
              className="primary-action launch-primary-action"
              onClick={props.onLaunch}
              disabled={props.creatingGame}
            >
              {props.creatingGame ? "Fundando..." : "Fundar ciudad y generar mapa"}
            </button>
          </div>
        </section>
      </main>
    </section>
  );
}

function SummaryItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="launch-summary-item">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function safeAbbreviation(value: string) {
  const normalized = value.trim().toUpperCase();
  return normalized.length > 0 ? normalized : "PCY";
}
