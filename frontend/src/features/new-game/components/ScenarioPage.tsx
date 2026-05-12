import { initialScenarios, type ScenarioId } from "../constants";

interface ScenarioPageProps {
  selectedScenario: ScenarioId;
  onBack: () => void;
  onContinue: () => void;
  onSelect: (value: ScenarioId) => void;
}

export function ScenarioPage(props: ScenarioPageProps) {
  return (
    <section className="screen step-screen">
      <div className="step-copy">
        <p className="eyebrow">Paso 2 de 4</p>
        <h1>Ahora elegí desde dónde empieza la historia.</h1>
        <p>
          Esta página define la presión inicial de la franquicia y el tono con el que la ciudad va
          a leer tus primeros movimientos.
        </p>
      </div>

      <article className="step-card wide-card">
        <div className="scenario-grid">
          {initialScenarios.map((scenario) => {
            const active = scenario.id === props.selectedScenario;

            return (
              <button
                key={scenario.id}
                type="button"
                className={active ? "scenario-card active" : "scenario-card"}
                onClick={() => props.onSelect(scenario.id)}
              >
                <span className="scenario-label">{scenario.label}</span>
                <strong>{scenario.roster}</strong>
                <p>{scenario.pressure}</p>
                <p>{scenario.city}</p>
              </button>
            );
          })}
        </div>

        <div className="page-actions split-actions">
          <button type="button" className="secondary-action" onClick={props.onBack}>
            Volver
          </button>
          <button type="button" className="primary-action" onClick={props.onContinue}>
            Continuar al gobierno de ciudad
          </button>
        </div>
      </article>
    </section>
  );
}
