import { colorPresets } from "../constants";
import { ColorField, FranchisePreview } from "./common";

interface IdentityPageProps {
  accentColor: string;
  abbreviation: string;
  cityName: string;
  franchiseName: string;
  primaryColor: string;
  secondaryColor: string;
  onAbbreviationChange: (value: string) => void;
  onAccentColorChange: (value: string) => void;
  onCityNameChange: (value: string) => void;
  onContinue: () => void;
  onFranchiseNameChange: (value: string) => void;
  onPrimaryColorChange: (value: string) => void;
  onSecondaryColorChange: (value: string) => void;
}

export function IdentityPage(props: IdentityPageProps) {
  return (
    <section className="screen step-screen">
      <div className="step-copy">
        <p className="eyebrow">Paso 1 de 4</p>
        <h1>Primero definí la identidad visible de la franquicia.</h1>
        <p>
          En esta pantalla solo se resuelve eso. Nombre, sigla y paleta. Nada más compite por tu
          atención.
        </p>
      </div>

      <div className="step-layout">
        <article className="step-card">
          <div className="field-grid">
            <label className="field">
              <span>Ciudad</span>
              <input value={props.cityName} onChange={(event) => props.onCityNameChange(event.target.value)} />
            </label>

            <label className="field">
              <span>Franquicia</span>
              <input
                value={props.franchiseName}
                onChange={(event) => props.onFranchiseNameChange(event.target.value)}
              />
            </label>

            <label className="field field-short">
              <span>Abreviatura</span>
              <input
                value={props.abbreviation}
                maxLength={3}
                onChange={(event) => props.onAbbreviationChange(event.target.value.toUpperCase())}
              />
            </label>
          </div>

          <div className="panel-section">
            <div className="section-title-row">
              <div>
                <p className="eyebrow">Paleta</p>
                <h3>Colores fundacionales</h3>
              </div>
              <p className="microcopy">
                Esta identidad ya se persiste de verdad en backend cuando la partida queda creada.
              </p>
            </div>

            <div className="color-stack">
              <ColorField
                label="Primario"
                value={props.primaryColor}
                onChange={props.onPrimaryColorChange}
                presets={colorPresets}
              />
              <ColorField
                label="Secundario"
                value={props.secondaryColor}
                onChange={props.onSecondaryColorChange}
                presets={colorPresets}
              />
              <ColorField
                label="Acento"
                value={props.accentColor}
                onChange={props.onAccentColorChange}
                presets={colorPresets}
              />
            </div>
          </div>

          <div className="page-actions">
            <button type="button" className="primary-action" onClick={props.onContinue}>
              Continuar al escenario inicial
            </button>
          </div>
        </article>

        <FranchisePreview
          accentColor={props.accentColor}
          abbreviation={props.abbreviation}
          cityName={props.cityName}
          contextLabel="Preview viva"
          franchiseName={props.franchiseName}
          primaryColor={props.primaryColor}
          secondaryColor={props.secondaryColor}
        />
      </div>
    </section>
  );
}
