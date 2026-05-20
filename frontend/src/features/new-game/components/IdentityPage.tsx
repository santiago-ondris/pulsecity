import type { CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import { colorPresets } from "../constants";

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
  const themeStyle = {
    "--identity-backdrop": `url("${skylineBackdrop}")`,
    "--franchise-primary": props.primaryColor,
    "--franchise-secondary": props.secondaryColor,
    "--franchise-accent": props.accentColor,
  } as CSSProperties;

  return (
    <section className="identity-builder" style={themeStyle}>
      <div className="identity-builder__image" />
      <div className="identity-builder__shade" />

      <header className="identity-builder__topbar">
        <div>
          <p className="eyebrow">Paso 1 de 4</p>
          <strong>Identidad de franquicia</strong>
        </div>
        <span>{safeAbbreviation(props.abbreviation)}</span>
      </header>

      <main className="identity-builder__main">
        <section className="identity-live-card" aria-label="Vista previa de la identidad">
          <div className="identity-live-card__crest">
            <span>{safeAbbreviation(props.abbreviation)}</span>
          </div>
          <div className="identity-live-card__name">
            <p>{props.cityName || "Ciudad"}</p>
            <h1>{props.franchiseName || "Franquicia"}</h1>
          </div>
          <div className="identity-live-card__palette">
            <span style={{ background: props.primaryColor }} />
            <span style={{ background: props.secondaryColor }} />
            <span style={{ background: props.accentColor }} />
          </div>
        </section>

        <section className="identity-editor" aria-label="Editar identidad de franquicia">
          <div className="identity-editor__heading">
            <p className="eyebrow">Crear marca</p>
            <h2>Nombre, sigla y colores</h2>
          </div>

          <div className="identity-editor__fields">
            <label className="identity-input">
              <span>Ciudad</span>
              <input
                value={props.cityName}
                onChange={(event) => props.onCityNameChange(event.target.value)}
                placeholder="Nueva Aurora"
              />
            </label>

            <label className="identity-input">
              <span>Franquicia</span>
              <input
                value={props.franchiseName}
                onChange={(event) => props.onFranchiseNameChange(event.target.value)}
                placeholder="Lighthouses"
              />
            </label>

            <label className="identity-input identity-input--short">
              <span>Sigla</span>
              <input
                value={props.abbreviation}
                maxLength={3}
                onChange={(event) => props.onAbbreviationChange(event.target.value.toUpperCase())}
                placeholder="NAR"
              />
            </label>
          </div>

          <div className="identity-editor__palette" aria-label="Elegir paleta de colores">
            <PaletteRow label="Primario" value={props.primaryColor} onChange={props.onPrimaryColorChange} />
            <PaletteRow label="Secundario" value={props.secondaryColor} onChange={props.onSecondaryColorChange} />
            <PaletteRow label="Acento" value={props.accentColor} onChange={props.onAccentColorChange} />
          </div>

          <button type="button" className="primary-action identity-builder__continue" onClick={props.onContinue}>
            Continuar al escenario inicial
          </button>
        </section>
      </main>
    </section>
  );
}

function PaletteRow({
  label,
  value,
  onChange,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
}) {
  return (
    <div className="identity-palette-row">
      <div className="identity-palette-row__label">
        <span>{label}</span>
        <strong>{value}</strong>
      </div>

      <div className="identity-palette-row__options">
        {colorPresets.map((preset) => (
          <button
            key={`${label}-${preset.value}`}
            type="button"
            className={
              value.toLowerCase() === preset.value.toLowerCase()
                ? "identity-color-dot active"
                : "identity-color-dot"
            }
            onClick={() => onChange(preset.value)}
            aria-label={`${label} ${preset.label}`}
            title={preset.label}
          >
            <span style={{ background: preset.value }} />
          </button>
        ))}

        <label className="identity-color-picker" title={`Elegir color ${label.toLowerCase()}`}>
          <input type="color" value={value} onChange={(event) => onChange(event.target.value)} />
        </label>
      </div>
    </div>
  );
}

function safeAbbreviation(value: string) {
  const normalized = value.trim().toUpperCase();
  return normalized.length > 0 ? normalized : "PCY";
}
