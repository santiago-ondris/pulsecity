import { cityManagementModes, type CityManagementModeId } from "../constants";

interface ManagementPageProps {
  cityManagementMode: CityManagementModeId;
  onBack: () => void;
  onContinue: () => void;
  onSelect: (value: CityManagementModeId) => void;
}

export function ManagementPage(props: ManagementPageProps) {
  return (
    <section className="screen step-screen">
      <div className="step-copy">
        <p className="eyebrow">Paso 3 de 4</p>
        <h1>Después decidí cómo se cruza tu poder con la ciudad.</h1>
        <p>
          Acá se define si la ciudad es un sistema que influís desde afuera o un tablero que
          controlás directamente.
        </p>
      </div>

      <article className="step-card wide-card">
        <div className="management-grid dual-grid">
          {cityManagementModes.map((mode) => {
            const active = mode.id === props.cityManagementMode;

            return (
              <button
                key={mode.id}
                type="button"
                className={active ? "management-card active" : "management-card"}
                onClick={() => props.onSelect(mode.id)}
              >
                <span className="scenario-label">{mode.label}</span>
                <strong>{mode.description}</strong>
                <p>{mode.impact}</p>
              </button>
            );
          })}
        </div>

        <div className="page-actions split-actions">
          <button type="button" className="secondary-action" onClick={props.onBack}>
            Volver
          </button>
          <button type="button" className="primary-action" onClick={props.onContinue}>
            Continuar a la revisión final
          </button>
        </div>
      </article>
    </section>
  );
}
