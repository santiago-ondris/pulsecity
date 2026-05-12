import type { CSSProperties } from "react";

import { safeAbbreviation } from "../helpers";

export function StatusBadge({
  label,
  tone,
}: {
  label: string;
  tone: "primary" | "info" | "urgent";
}) {
  return (
    <span
      className={
        tone === "primary"
          ? "badge badge-primary"
          : tone === "urgent"
            ? "badge badge-urgent"
            : "badge badge-info"
      }
    >
      <span className="badge-dot" />
      {label}
    </span>
  );
}

export function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="metric">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

export function Swatch({ label, value }: { label: string; value: string }) {
  return (
    <div className="swatch">
      <span className="swatch-chip" style={{ background: value }} />
      <div>
        <span>{label}</span>
        <strong>{value}</strong>
      </div>
    </div>
  );
}

export function ColorField({
  label,
  value,
  onChange,
  presets,
}: {
  label: string;
  value: string;
  onChange: (value: string) => void;
  presets: readonly { label: string; value: string }[];
}) {
  return (
    <div className="color-field">
      <label className="field field-inline">
        <span>{label}</span>
        <div className="color-input-row">
          <input
            className="color-picker"
            type="color"
            value={value}
            onChange={(event) => onChange(event.target.value)}
          />
          <input value={value} onChange={(event) => onChange(event.target.value)} />
        </div>
      </label>

      <div className="preset-row">
        {presets.map((preset) => (
          <button
            key={`${label}-${preset.value}`}
            type="button"
            className={value.toLowerCase() === preset.value.toLowerCase() ? "preset active" : "preset"}
            onClick={() => onChange(preset.value)}
            aria-label={`${label} ${preset.label}`}
          >
            <span style={{ background: preset.value }} />
          </button>
        ))}
      </div>
    </div>
  );
}

export function FranchisePreview({
  cityName,
  franchiseName,
  abbreviation,
  primaryColor,
  secondaryColor,
  accentColor,
  contextLabel,
}: {
  cityName: string;
  franchiseName: string;
  abbreviation: string;
  primaryColor: string;
  secondaryColor: string;
  accentColor: string;
  contextLabel: string;
}) {
  const previewStyle = {
    "--franchise-primary": primaryColor,
    "--franchise-secondary": secondaryColor,
    "--franchise-accent": accentColor,
  } as CSSProperties;

  return (
    <article className="preview-card" style={previewStyle}>
      <p className="eyebrow">{contextLabel}</p>
      <div className="franchise-crest">
        <span className="crest-mark">{safeAbbreviation(abbreviation)}</span>
      </div>
      <div className="franchise-headline">
        <p>{cityName}</p>
        <h3>{franchiseName}</h3>
      </div>
      <div className="swatch-row">
        <Swatch label="Primario" value={primaryColor} />
        <Swatch label="Secundario" value={secondaryColor} />
        <Swatch label="Acento" value={accentColor} />
      </div>
    </article>
  );
}
