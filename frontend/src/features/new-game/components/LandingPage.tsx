import { useState, type CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import { managementModeLabel, scenarioById } from "../helpers";
import type { GameSummary, UserSession } from "../../../types";

type LandingPanel = "overview" | "library";

interface LandingPageProps {
  activeAuthKind: "none" | "guest" | "user";
  games: GameSummary[];
  restoringSession: boolean;
  selectedGame: GameSummary | null;
  selectedGameId: string;
  userSession: UserSession | null;
  onContinueSelectedGame: () => void;
  onLogoutUser: () => void;
  onSelectGame: (gameId: string) => void;
  onStart: () => void;
}

export function LandingPage(props: LandingPageProps) {
  const [panel, setPanel] = useState<LandingPanel>("overview");

  const selectedScenario = props.selectedGame
    ? scenarioById(props.selectedGame.initial_scenario)
    : null;

  const identityLabel =
    props.activeAuthKind === "user"
      ? props.userSession?.user.display_name ?? "Cuenta activa"
      : props.activeAuthKind === "guest"
        ? "Invitado activo"
        : "Sin sesion";

  const sessionSummary =
    props.activeAuthKind === "user"
      ? props.userSession?.user.email ?? "Biblioteca vinculada a cuenta"
      : props.activeAuthKind === "guest"
        ? "Entrada temporal lista para fundar una ciudad"
        : "Todavia no elegiste una identidad para entrar";

  function handleLoad() {
    setPanel("library");
  }

  return (
    <section
      className="screen landing-screen landing-cinematic"
      style={{ "--landing-backdrop": `url("${skylineBackdrop}")` } as CSSProperties}
    >
      <div className="landing-cinematic__backdrop" />
      <div className="landing-cinematic__veil" />

      <header className="landing-topbar">
        <div>
          <p className="eyebrow">PulseCity</p>
          <strong className="landing-brand">Franchise and City Control</strong>
        </div>

        <div className="topbar-status">
          <span>{identityLabel}</span>
          <small>{props.restoringSession ? "Restaurando..." : sessionSummary}</small>
        </div>
      </header>

      <div className="landing-cinematic__content landing-cinematic__content--menu">
        <section className="landing-menu-panel">
          <div className="hero-copy-block landing-menu-copy">
            <p className="eyebrow">Basketball civic simulation</p>
            <h1>Una ciudad entera espera tu primera orden.</h1>
            <p className="landing-copy">
              Elegi como entrar. Lo demas aparece despues.
            </p>
          </div>

          <div className="landing-primary-actions landing-primary-actions--stacked">
              <button
                type="button"
                className="landing-cta landing-cta--primary"
                onClick={props.onStart}
                disabled={props.restoringSession}
              >
                <span>Nueva partida</span>
                <small>Empezar una franquicia desde cero</small>
            </button>

            <button
              type="button"
              className="landing-cta"
              onClick={handleLoad}
              disabled={props.restoringSession}
            >
              <span>Cargar partida</span>
              <small>
                {props.games.length > 0
                  ? `${props.games.length} mundo(s) guardado(s) listo(s) para retomar`
                  : "Abrir biblioteca de ciudades guardadas"}
              </small>
            </button>

          </div>
          <div className="landing-menu-meta">
            <div className="landing-menu-meta__row">
              <span>Cuenta</span>
              <strong>{identityLabel}</strong>
            </div>
            <div className="landing-menu-meta__row">
              <span>Biblioteca</span>
              <strong>{props.games.length} partida(s)</strong>
            </div>
            <div className="landing-menu-meta__row">
              <span>Estado</span>
              <strong>{props.restoringSession ? "Restaurando" : "Lista"}</strong>
            </div>
          </div>
          {props.activeAuthKind === "user" ? (
            <div className="landing-inline-actions">
              <button
                type="button"
                className="secondary-action"
                onClick={props.onLogoutUser}
                disabled={props.restoringSession}
              >
                Cerrar sesion
              </button>
            </div>
          ) : null}
        </section>
      </div>

      {panel === "library" ? (
        <div className="landing-overlay">
          <aside className="landing-side-panel__section landing-modal-panel">
            <>
              <div className="panel-heading">
                <div>
                  <p className="eyebrow">Biblioteca</p>
                  <h2>Cargar partida</h2>
                </div>
                <div className="panel-heading-actions">
                  <span className="hero-chip">
                  {props.games.length > 0 ? `${props.games.length} disponibles` : "Placeholder"}
                  </span>
                  <button
                    type="button"
                    className="panel-close"
                    onClick={() => setPanel("overview")}
                  >
                    Cerrar
                  </button>
                </div>
              </div>

              {props.games.length === 0 ? (
                <div className="landing-empty-card">
                  <strong>Todavia no hay mundos guardados.</strong>
                  <p>
                    Esta vista despues puede crecer como pagina propia con autosaves, previews y
                    filtros. Por ahora ya queda separada de la landing principal.
                  </p>
                </div>
              ) : (
                <div className="landing-library-list">
                  {props.games.map((game) => {
                    const scenario = scenarioById(game.initial_scenario);
                    const active = game.game_id === props.selectedGameId;

                    return (
                      <button
                        key={game.game_id}
                        type="button"
                        className={active ? "library-card active" : "library-card"}
                        onClick={() => props.onSelectGame(game.game_id)}
                      >
                        <div>
                          <p className="eyebrow">{game.city_name}</p>
                          <h3>{game.franchise_name}</h3>
                          <p>
                            {scenario.label} · {managementModeLabel(game.city_management_mode)}
                          </p>
                        </div>
                        <small>{formatUpdatedAt(game.updated_at)}</small>
                      </button>
                    );
                  })}
                </div>
              )}

              {props.selectedGame && selectedScenario ? (
                <div className="landing-selected-save">
                  <div>
                    <p className="eyebrow">{props.selectedGame.city_name}</p>
                    <h3>{props.selectedGame.franchise_name}</h3>
                    <p>
                      {selectedScenario.label} ·{" "}
                      {managementModeLabel(props.selectedGame.city_management_mode)}
                    </p>
                  </div>
                  <div className="landing-metadata-list">
                    <div>
                      <span>Estado</span>
                      <strong>{props.selectedGame.status}</strong>
                    </div>
                    <div>
                      <span>Dueño</span>
                      <strong>{props.selectedGame.owner_kind === "user" ? "Cuenta" : "Invitado"}</strong>
                    </div>
                    <div>
                      <span>Actualizada</span>
                      <strong>{formatUpdatedAt(props.selectedGame.updated_at)}</strong>
                    </div>
                  </div>
                  <button
                    type="button"
                    className="primary-action"
                    onClick={props.onContinueSelectedGame}
                    disabled={props.restoringSession}
                  >
                    Continuar partida
                  </button>
                </div>
              ) : null}
            </>
          </aside>
        </div>
      ) : null}
    </section>
  );
}

function formatUpdatedAt(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "Sin fecha";
  }

  return date.toLocaleString("es-AR", {
    dateStyle: "short",
    timeStyle: "short",
  });
}
