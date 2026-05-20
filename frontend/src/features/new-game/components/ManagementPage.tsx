import type { CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import { cityManagementModes, type CityManagementModeId } from "../constants";

interface ManagementPageProps {
  cityManagementMode: CityManagementModeId;
  onBack: () => void;
  onContinue: () => void;
  onSelect: (value: CityManagementModeId) => void;
}

const managementMeta: Record<
  CityManagementModeId,
  {
    code: string;
    stance: string;
    accent: string;
  }
> = {
  owner_influence: {
    code: "IF",
    stance: "Influencia indirecta",
    accent: "#00C896",
  },
  dual_figure: {
    code: "DF",
    stance: "Control directo",
    accent: "#FFAA00",
  },
};

export function ManagementPage(props: ManagementPageProps) {
  const selectedMode = cityManagementModes.find((mode) => mode.id === props.cityManagementMode) ?? cityManagementModes[0];
  const selectedMeta = managementMeta[selectedMode.id];

  return (
    <section
      className="management-builder"
      style={{ "--management-backdrop": `url("${skylineBackdrop}")`, "--management-accent": selectedMeta.accent } as CSSProperties}
    >
      <div className="management-builder__image" />
      <div className="management-builder__shade" />

      <header className="management-builder__topbar">
        <div>
          <p className="eyebrow">Paso 3 de 4</p>
          <strong>Gobierno de ciudad</strong>
        </div>
        <span>{selectedMeta.code}</span>
      </header>

      <main className="management-builder__main">
        <section className="management-summary" aria-label="Modo de gestion seleccionado">
          <p className="eyebrow">Modelo de poder</p>
          <h1>{selectedMode.label}</h1>
          <div className="management-summary__stance">
            <span>{selectedMeta.stance}</span>
            <strong>{selectedMode.description}</strong>
          </div>
          <p>{selectedMode.impact}</p>
        </section>

        <section className="management-picker" aria-label="Elegir modo de gestion">
          <div className="management-picker__heading">
            <p className="eyebrow">Elegir rol</p>
            <h2>Como se cruza tu poder con la ciudad</h2>
          </div>

          <div className="management-option-list">
            {cityManagementModes.map((mode) => {
              const meta = managementMeta[mode.id];
              const active = mode.id === props.cityManagementMode;

              return (
                <button
                  key={mode.id}
                  type="button"
                  className={active ? "management-option active" : "management-option"}
                  style={{ "--option-accent": meta.accent } as CSSProperties}
                  onClick={() => props.onSelect(mode.id)}
                >
                  <span className="management-option__code">{meta.code}</span>
                  <span className="management-option__body">
                    <strong>{mode.label}</strong>
                    <small>{meta.stance}</small>
                    <p>{mode.impact}</p>
                  </span>
                </button>
              );
            })}
          </div>

          <div className="management-actions">
            <button type="button" className="secondary-action" onClick={props.onBack}>
              Volver
            </button>
            <button type="button" className="primary-action" onClick={props.onContinue}>
              Continuar a la revision final
            </button>
          </div>
        </section>
      </main>
    </section>
  );
}
