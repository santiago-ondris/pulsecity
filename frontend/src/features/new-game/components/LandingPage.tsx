import { useState } from "react";

import { managementModeLabel, scenarioById } from "../helpers";
import type { GameSummary, UserSession } from "../../../types";

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
  const [registerEmail, setRegisterEmail] = useState("");
  const [registerDisplayName, setRegisterDisplayName] = useState("");
  const [registerPassword, setRegisterPassword] = useState("");
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");

  const selectedScenario = props.selectedGame
    ? scenarioById(props.selectedGame.initial_scenario)
    : null;
  const sessionTone =
    props.activeAuthKind === "user" ? "Cuenta activa" : props.activeAuthKind === "guest" ? "Invitado activo" : "Sin sesión";

  return (
    <section className="screen landing-screen">
      <div className="landing-hero">
        <div className="hero-ribbon">
          <p className="eyebrow">PulseCity</p>
          <span className="hero-chip">{props.restoringSession ? "Rehidratando sesión" : sessionTone}</span>
        </div>
        <h1>Entrar, elegir identidad y decidir si el mundo sigue o vuelve a nacer.</h1>
        <p className="landing-copy">
          La entrada de Milestone 1 ya quedó completa como flujo: cuenta o invitado, biblioteca
          disponible, selección de partida y decisión explícita entre continuar una historia o
          fundar una ciudad nueva.
        </p>
      </div>

      <div className="status-banner">
        <strong>Estado actual</strong>
        <span>{props.restoringSession ? "Restaurando sesión guardada..." : props.status}</span>
      </div>

      <section className="entry-grid">
        <article className="entry-panel">
          <div className="panel-header">
            <div>
              <p className="eyebrow">Acceso</p>
              <h2>Identidad activa</h2>
            </div>
            <p className="microcopy">Una sola identidad manda la biblioteca y la creación de partidas.</p>
          </div>

          <div className="actor-stack">
            <div className={props.activeAuthKind === "user" ? "actor-card active" : "actor-card"}>
              <span>Cuenta</span>
              {props.userSession ? (
                <>
                  <strong>{props.userSession.user.display_name}</strong>
                  <p>{props.userSession.user.email}</p>
                </>
              ) : (
                <>
                  <strong>Sin sesión autenticada</strong>
                  <p>Registrate o iniciá sesión para dejar las partidas bajo una cuenta real.</p>
                </>
              )}
            </div>

            <div className={props.activeAuthKind === "guest" ? "actor-card active" : "actor-card"}>
              <span>Invitado</span>
              {props.guestToken ? (
                <>
                  <strong>Token disponible</strong>
                  <p>{props.guestToken}</p>
                </>
              ) : (
                <>
                  <strong>Sin invitado activo</strong>
                  <p>Podés entrar rápido sin cuenta y después migrar esa biblioteca a usuario.</p>
                </>
              )}
            </div>
          </div>

          <div className="landing-actions">
            <button
              type="button"
              className="primary-action landing-action"
              onClick={props.onStart}
              disabled={props.restoringSession}
            >
              Nueva partida
            </button>
            {props.activeAuthKind === "none" ? (
              <button
                type="button"
                className="secondary-action landing-action"
                onClick={props.onCreateGuestSession}
                disabled={props.creatingGuestSession || props.restoringSession}
              >
                {props.creatingGuestSession ? "Creando invitado..." : "Jugar como invitado"}
              </button>
            ) : null}
          </div>

          {props.activeAuthKind !== "none" ? (
            <div className="identity-actions">
              {props.activeAuthKind === "user" ? (
                <button
                  type="button"
                  className="secondary-action"
                  onClick={props.onLogoutUser}
                  disabled={props.restoringSession}
                >
                  Cerrar sesión de cuenta
                </button>
              ) : null}
              {props.activeAuthKind === "user" && props.guestToken ? (
                <button
                  type="button"
                  className="secondary-action"
                  onClick={props.onSwitchToGuestSession}
                  disabled={props.restoringSession}
                >
                  Cambiar a invitado guardado
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
          ) : null}
        </article>

        <article className="entry-panel">
          <div className="panel-header">
            <div>
              <p className="eyebrow">Cuenta</p>
              <h2>Registro y login</h2>
            </div>
            <p className="microcopy">La cuenta toma ownership del flujo actual y de las partidas migradas.</p>
          </div>

          <div className="auth-access-grid">
            <section className="auth-panel">
              <p className="copy">
                Crear cuenta para pasar del acceso temporal a una identidad persistente del jugador.
              </p>
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
                  placeholder="Mínimo 8 caracteres"
                />
              </label>
              <button
                type="button"
                className="primary-action"
                disabled={props.authenticatingUser || props.restoringSession}
                onClick={() =>
                  props.onRegister(registerEmail, registerDisplayName, registerPassword)}
              >
                {props.authenticatingUser ? "Creando cuenta..." : "Registrarme"}
              </button>
            </section>

            <section className="auth-panel">
              <p className="copy">
                Si la cuenta ya existe, esta entrada carga directamente su biblioteca autenticada.
              </p>
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
                  placeholder="Tu contraseña"
                />
              </label>
              <button
                type="button"
                className="secondary-action"
                disabled={props.authenticatingUser || props.restoringSession}
                onClick={() => props.onLogin(loginEmail, loginPassword)}
              >
                {props.authenticatingUser ? "Ingresando..." : "Iniciar sesión"}
              </button>
            </section>
          </div>
        </article>
      </section>

      <section className="library-grid">
        <article className="entry-panel">
          <div className="panel-header">
            <div>
              <p className="eyebrow">Biblioteca</p>
              <h2>Partidas disponibles</h2>
            </div>
            <p className="microcopy">
              {props.games.length === 0
                ? "Todavía no hay mundos guardados para esta identidad."
                : `${props.games.length} partida(s) disponible(s) para cargar.`}
            </p>
          </div>

          {props.games.length === 0 ? (
            <div className="guest-empty-state">
              <strong>Todavía no hay partidas para este actor.</strong>
              <p>Si querés cerrar Milestone 1 desde cero, la próxima fundación quedará asociada a esta identidad.</p>
            </div>
          ) : (
            <ul className="guest-game-list selectable">
              {props.games.map((game) => {
                const scenario = scenarioById(game.initial_scenario);
                const active = game.game_id === props.selectedGameId;

                return (
                  <li key={game.game_id}>
                    <button
                      type="button"
                      className={active ? "guest-game-card selectable active" : "guest-game-card selectable"}
                      onClick={() => props.onSelectGame(game.game_id)}
                    >
                      <div>
                        <p className="eyebrow">{game.city_name}</p>
                        <h3>{game.franchise_name}</h3>
                        <p className="guest-game-meta">
                          {scenario.label} · {managementModeLabel(game.city_management_mode)}
                        </p>
                        <p className="guest-game-meta">
                          Dueño: {game.owner_kind === "user" ? "Cuenta" : "Invitado"} · Estado: {game.status}
                        </p>
                      </div>
                      <span className="guest-game-date">{formatUpdatedAt(game.updated_at)}</span>
                    </button>
                  </li>
                );
              })}
            </ul>
          )}
        </article>

        <article className="entry-panel">
          <div className="panel-header">
            <div>
              <p className="eyebrow">Carga</p>
              <h2>Decisión actual</h2>
            </div>
            <p className="microcopy">Una selección clara antes de entrar evita mezclar contexto y progreso.</p>
          </div>

          {props.selectedGame && selectedScenario ? (
            <div className="loadout-panel">
              <div className="loadout-hero">
                <p className="eyebrow">{props.selectedGame.city_name}</p>
                <h3>{props.selectedGame.franchise_name}</h3>
                <p className="copy">
                  {selectedScenario.label} · {managementModeLabel(props.selectedGame.city_management_mode)}
                </p>
              </div>

              <div className="loadout-metrics">
                <div className="metric">
                  <span>Estado</span>
                  <strong>{props.selectedGame.status}</strong>
                </div>
                <div className="metric">
                  <span>Dueño</span>
                  <strong>{props.selectedGame.owner_kind === "user" ? "Cuenta" : "Invitado"}</strong>
                </div>
                <div className="metric">
                  <span>Actualizada</span>
                  <strong>{formatUpdatedAt(props.selectedGame.updated_at)}</strong>
                </div>
              </div>

              <div className="loadout-actions">
                <button
                  type="button"
                  className="primary-action"
                  onClick={props.onContinueSelectedGame}
                  disabled={props.restoringSession}
                >
                  Continuar partida
                </button>
                <button
                  type="button"
                  className="secondary-action"
                  onClick={props.onStart}
                  disabled={props.restoringSession}
                >
                  Fundar otra ciudad
                </button>
              </div>
            </div>
          ) : (
            <div className="guest-empty-state">
              <strong>No hay una partida seleccionada.</strong>
              <p>Elegí una fundación de la biblioteca para retomarla, o avanzá directo con una nueva partida.</p>
            </div>
          )}
        </article>
      </section>
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
