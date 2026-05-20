import type { CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import { initialScenarios, type ScenarioId } from "../constants";

interface ScenarioPageProps {
  selectedScenario: ScenarioId;
  onBack: () => void;
  onContinue: () => void;
  onSelect: (value: ScenarioId) => void;
}

const scenarioMeta: Record<
  ScenarioId,
  {
    code: string;
    pulse: string;
    accent: string;
  }
> = {
  rebuild: {
    code: "RB",
    pulse: "Paciencia",
    accent: "#00C896",
  },
  contention: {
    code: "PO",
    pulse: "Urgencia",
    accent: "#FFAA00",
  },
  decline: {
    code: "HD",
    pulse: "Presion",
    accent: "#FF6B2B",
  },
  expansion: {
    code: "EX",
    pulse: "Horizonte",
    accent: "#7B8CDE",
  },
};

export function ScenarioPage(props: ScenarioPageProps) {
  const selectedScenario = initialScenarios.find((scenario) => scenario.id === props.selectedScenario) ?? initialScenarios[0];
  const selectedMeta = scenarioMeta[selectedScenario.id];

  return (
    <section
      className="scenario-builder"
      style={{ "--scenario-backdrop": `url("${skylineBackdrop}")`, "--scenario-accent": selectedMeta.accent } as CSSProperties}
    >
      <div className="scenario-builder__image" />
      <div className="scenario-builder__shade" />

      <header className="scenario-builder__topbar">
        <div>
          <p className="eyebrow">Paso 2 de 4</p>
          <strong>Escenario inicial</strong>
        </div>
        <span>{selectedMeta.code}</span>
      </header>

      <main className="scenario-builder__main">
        <section className="scenario-summary" aria-label="Escenario seleccionado">
          <p className="eyebrow">Punto de partida</p>
          <h1>{selectedScenario.label}</h1>
          <div className="scenario-summary__pulse">
            <span>{selectedMeta.pulse}</span>
            <strong>{selectedScenario.pressure}</strong>
          </div>
          <div className="scenario-summary__details">
            <p>{selectedScenario.roster}</p>
            <p>{selectedScenario.city}</p>
          </div>
        </section>

        <section className="scenario-picker" aria-label="Elegir escenario inicial">
          <div className="scenario-picker__heading">
            <p className="eyebrow">Elegir contexto</p>
            <h2>Desde donde empieza la historia</h2>
          </div>

          <div className="scenario-option-list">
            {initialScenarios.map((scenario) => {
              const meta = scenarioMeta[scenario.id];
              const active = scenario.id === props.selectedScenario;

              return (
                <button
                  key={scenario.id}
                  type="button"
                  className={active ? "scenario-option active" : "scenario-option"}
                  style={{ "--option-accent": meta.accent } as CSSProperties}
                  onClick={() => props.onSelect(scenario.id)}
                >
                  <span className="scenario-option__code">{meta.code}</span>
                  <span className="scenario-option__body">
                    <strong>{scenario.label}</strong>
                    <small>{scenario.roster}</small>
                  </span>
                  <span className="scenario-option__status">{meta.pulse}</span>
                </button>
              );
            })}
          </div>

          <div className="scenario-actions">
            <button type="button" className="secondary-action" onClick={props.onBack}>
              Volver
            </button>
            <button type="button" className="primary-action" onClick={props.onContinue}>
              Continuar al gobierno de ciudad
            </button>
          </div>
        </section>
      </main>
    </section>
  );
}
