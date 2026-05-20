import { useState, type CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import { managementModeLabel, scenarioById } from "../helpers";
import type { GameSummary, UserSession } from "../../../types";

type LandingPanel = "overview" | "access" | "library";
type AuthMode = "login" | "register";

interface LandingPageProps {
  activeAuthKind: "none" | "guest" | "user";
  authenticatingUser: boolean;
  creatingGuestSession: boolean;
  games: GameSummary[];
  guestToken: string;
  restoringSession: boolean;
  selectedGame: GameSummary | null;
  selectedGameId: string;
  status: string;
  userSession: UserSession | null;
  onClearAllAccess: () => void;
  onContinueSelectedGame: () => void;
  onCreateGuestSession: () => void;
  onLogin: (email: string, password: string) => void;
  onLogoutUser: () => void;
  onRegister: (email: string, displayName: string, password: string) => void;
  onSelectGame: (gameId: string) => void;
  onStart: () => void;
  onSwitchToGuestSession: () => void;
}

export function LandingPage(props: LandingPageProps) {
  const [panel, setPanel] = useState<LandingPanel>("overview");
  const [authMode, setAuthMode] = useState<AuthMode>("login");
  const [registerEmail, setRegisterEmail] = useState("");
  const [registerDisplayName, setRegisterDisplayName] = useState("");
  const [registerPassword, setRegisterPassword] = useState("");
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");

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

  function handleStart() {
    if (props.activeAuthKind === "none") {
      setPanel("access");
      return;
    }

    props.onStart();
  }

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
              onClick={handleStart}
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

            <button
              type="button"
              className="landing-cta"
              onClick={props.onCreateGuestSession}
              disabled={props.creatingGuestSession || props.restoringSession}
            >
              <span>
                {props.activeAuthKind === "guest" ? "Seguir como invitado" : "Jugar como invitado"}
              </span>
              <small>
                {props.creatingGuestSession
                  ? "Creando acceso temporal..."
                  : "Entrar sin cuenta y resolver identidad mas adelante"}
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
        </section>
      </div>

      {(panel === "access" || panel === "library") ? (
        <div className="landing-overlay">
          <aside className="landing-side-panel__section landing-modal-panel">
            {panel === "access" ? (
              <>
              <div className="panel-heading">
                <div>
                  <p className="eyebrow">Cuenta</p>
                  <h2>Acceso del jugador</h2>
                </div>
                <button
                  type="button"
                  className="panel-close"
                  onClick={() => setPanel("overview")}
                >
                  Cerrar
                </button>
              </div>

              <div className="panel-toggle-row">
                <button
                  type="button"
                  className={authMode === "login" ? "panel-chip active" : "panel-chip"}
                  onClick={() => setAuthMode("login")}
                >
                  Login
                </button>
                <button
                  type="button"
                  className={authMode === "register" ? "panel-chip active" : "panel-chip"}
                  onClick={() => setAuthMode("register")}
                >
                  Registro
                </button>
              </div>

              <div className="landing-summary-card">
                <strong>{identityLabel}</strong>
                <p>{sessionSummary}</p>
              </div>

              {authMode === "login" ? (
                <div className="panel-toggle-row">
                  <div className="landing-form">
                    <label className="field">
                      <span>Email</span>
                      <input
                        type="email"
                        value={loginEmail}
                        onChange={(event) => setLoginEmail(event.target.value)}
                        placeholder="gm@pulsecity.test"
                      />
                    </label>
                    <label className="field">
                      <span>Password</span>
                      <input
                        type="password"
                        value={loginPassword}
                        onChange={(event) => setLoginPassword(event.target.value)}
                        placeholder="Tu password"
                      />
                    </label>
                    <button
                      type="button"
                      className="primary-action"
                      disabled={props.authenticatingUser || props.restoringSession}
                      onClick={() => props.onLogin(loginEmail, loginPassword)}
                    >
                      {props.authenticatingUser ? "Ingresando..." : "Iniciar sesion"}
                    </button>
                  </div>
                </div>
              ) : (
                <div className="landing-form">
                  <label className="field">
                    <span>Nombre visible</span>
                    <input
                      value={registerDisplayName}
                      onChange={(event) => setRegisterDisplayName(event.target.value)}
                      placeholder="Jordan Vale"
                    />
                  </label>
                  <label className="field">
                    <span>Email</span>
                    <input
                      type="email"
                      value={registerEmail}
                      onChange={(event) => setRegisterEmail(event.target.value)}
                      placeholder="gm@pulsecity.test"
                    />
                  </label>
                  <label className="field">
                    <span>Password</span>
                    <input
                      type="password"
                      value={registerPassword}
                      onChange={(event) => setRegisterPassword(event.target.value)}
                      placeholder="Minimo 8 caracteres"
                    />
                  </label>
                  <button
                    type="button"
                    className="primary-action"
                    disabled={props.authenticatingUser || props.restoringSession}
                    onClick={() =>
                      props.onRegister(registerEmail, registerDisplayName, registerPassword)}
                  >
                    {props.authenticatingUser ? "Creando cuenta..." : "Crear cuenta"}
                  </button>
                </div>
              )}

              <div className="landing-utility-actions">
                <button
                  type="button"
                  className="secondary-action"
                  onClick={props.onCreateGuestSession}
                  disabled={props.creatingGuestSession || props.restoringSession}
                >
                  {props.creatingGuestSession ? "Creando invitado..." : "Crear acceso invitado"}
                </button>

                {props.activeAuthKind === "user" ? (
                  <button
                    type="button"
                    className="secondary-action"
                    onClick={props.onLogoutUser}
                    disabled={props.restoringSession}
                  >
                    Cerrar sesion
                  </button>
                ) : null}

                {props.activeAuthKind === "user" && props.guestToken ? (
                  <button
                    type="button"
                    className="secondary-action"
                    onClick={props.onSwitchToGuestSession}
                    disabled={props.restoringSession}
                  >
                    Cambiar a invitado
                  </button>
                ) : null}

                <button
                  type="button"
                  className="secondary-action"
                  onClick={props.onClearAllAccess}
                  disabled={props.restoringSession}
                >
                  Limpiar sesiones locales
                </button>
              </div>
              </>
            ) : null}

            {panel === "library" ? (
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
            ) : null}
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
