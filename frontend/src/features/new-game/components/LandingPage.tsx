import { useState } from "react";

import { managementModeLabel, scenarioById } from "../helpers";
import type { GameSummary, UserSession } from "../../../types";

interface LandingPageProps {
  activeAuthKind: "none" | "guest" | "user";
  authenticatingUser: boolean;
  creatingGuestSession: boolean;
  games: GameSummary[];
  guestToken: string;
  status: string;
  userSession: UserSession | null;
  onContinueGame: (gameId: string) => void;
  onCreateGuestSession: () => void;
  onLogin: (email: string, password: string) => void;
  onRegister: (email: string, displayName: string, password: string) => void;
  onStart: () => void;
}

export function LandingPage(props: LandingPageProps) {
  const [registerEmail, setRegisterEmail] = useState("");
  const [registerDisplayName, setRegisterDisplayName] = useState("");
  const [registerPassword, setRegisterPassword] = useState("");
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");

  return (
    <section className="screen landing-screen">
      <div className="landing-hero">
        <p className="eyebrow">PulseCity</p>
        <h1>La puerta de entrada ya distingue entre visitante ocasional y cuenta real.</h1>
        <p className="landing-copy">
          Este corte suma registro e inicio de sesión sin romper la base ya armada para invitados.
          La idea es simple: el jugador ya puede existir como usuario real, pero todavía sin
          mezclar migración de partidas guest hacia cuenta.
        </p>

        <div className="landing-actions">
          <button type="button" className="primary-action landing-action" onClick={props.onStart}>
            Nueva partida
          </button>
          {props.activeAuthKind === "none" ? (
            <button
              type="button"
              className="secondary-action landing-action"
              onClick={props.onCreateGuestSession}
              disabled={props.creatingGuestSession}
            >
              {props.creatingGuestSession ? "Creando invitado..." : "Jugar como invitado"}
            </button>
          ) : null}
        </div>
      </div>

      <div className="landing-grid auth-grid">
        <article className="landing-note">
          <span>01</span>
          <strong>Invitado</strong>
          <p>
            {props.guestToken
              ? `Token activo disponible: ${props.guestToken}`
              : "Todavía no hay sesión invitada persistida."}
          </p>
        </article>
        <article className="landing-note">
          <span>02</span>
          <strong>Cuenta</strong>
          <p>
            {props.userSession
              ? `${props.userSession.user.display_name} ya está autenticado.`
              : "Ahora podés registrarte o iniciar sesión desde esta misma entrada."}
          </p>
        </article>
        <article className="landing-note">
          <span>03</span>
          <strong>Ownership</strong>
          <p>Las partidas nuevas se crean bajo el actor activo, sin perder separación de dueños.</p>
        </article>
        <article className="landing-note">
          <span>04</span>
          <strong>Siguiente paso</strong>
          <p>Después de esto ya queda listo el terreno para migrar partidas guest a cuenta.</p>
        </article>
      </div>

      <div className="auth-access-grid">
        <section className="auth-panel">
          <div className="panel-header">
            <p className="eyebrow">Registro</p>
            <h2>Crear cuenta</h2>
          </div>
          <p className="copy">Cuenta real mínima para empezar a dejar partidas bajo usuario.</p>
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
            disabled={props.authenticatingUser}
            onClick={() => props.onRegister(registerEmail, registerDisplayName, registerPassword)}
          >
            {props.authenticatingUser ? "Creando cuenta..." : "Registrarme"}
          </button>
        </section>

        <section className="auth-panel">
          <div className="panel-header">
            <p className="eyebrow">Login</p>
            <h2>Iniciar sesión</h2>
          </div>
          <p className="copy">
            Si ya existe una cuenta, esta sesión pasa a trabajar bajo usuario autenticado.
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
            disabled={props.authenticatingUser}
            onClick={() => props.onLogin(loginEmail, loginPassword)}
          >
            {props.authenticatingUser ? "Ingresando..." : "Iniciar sesión"}
          </button>
        </section>
      </div>

      <section className="guest-session-panel">
        <div className="panel-header">
          <p className="eyebrow">
            {props.activeAuthKind === "user" ? "Sesión autenticada" : "Sesión activa"}
          </p>
          <h2>Partidas asociadas</h2>
        </div>
        <p className="copy">{props.status}</p>

        {props.userSession ? (
          <div className="account-summary">
            <strong>{props.userSession.user.display_name}</strong>
            <p>{props.userSession.user.email}</p>
          </div>
        ) : null}

        {props.games.length === 0 ? (
          <div className="guest-empty-state">
            <strong>Todavía no hay partidas para este actor.</strong>
            <p>La próxima fundación quedará registrada bajo la sesión actualmente activa.</p>
          </div>
        ) : (
          <ul className="guest-game-list">
            {props.games.map((game) => {
              const scenario = scenarioById(game.initial_scenario);

              return (
                <li key={game.game_id} className="guest-game-card">
                  <div>
                    <p className="eyebrow">{game.city_name}</p>
                    <h3>{game.franchise_name}</h3>
                    <p className="guest-game-meta">
                      {scenario.label} · {managementModeLabel(game.city_management_mode)}
                    </p>
                    <p className="guest-game-meta">
                      Dueño: {game.owner_kind === "user" ? "Cuenta" : "Invitado"} · Estado:{" "}
                      {game.status}
                    </p>
                  </div>
                  <div className="guest-game-actions">
                    <span>{formatUpdatedAt(game.updated_at)}</span>
                    <button
                      type="button"
                      className="secondary-action"
                      onClick={() => props.onContinueGame(game.game_id)}
                    >
                      Continuar
                    </button>
                  </div>
                </li>
              );
            })}
          </ul>
        )}
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
