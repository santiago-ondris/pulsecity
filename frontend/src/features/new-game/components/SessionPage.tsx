import { useState, type CSSProperties } from "react";

import skylineBackdrop from "../../assets/landing-city-night.svg";
import type { UserSession } from "../../../types";

type AuthMode = "login" | "register";

interface SessionPageProps {
  activeAuthKind: "none" | "guest" | "user";
  authenticatingUser: boolean;
  creatingGuestSession: boolean;
  guestToken: string;
  restoringSession: boolean;
  status: string;
  userSession: UserSession | null;
  onClearAllAccess: () => void;
  onCreateGuestSession: () => void;
  onForgotPassword: () => void;
  onLogin: (email: string, password: string) => void;
  onLogoutUser: () => void;
  onRegister: (email: string, displayName: string, password: string) => void;
  onSwitchToGuestSession: () => void;
}

export function SessionPage(props: SessionPageProps) {
  const [authMode, setAuthMode] = useState<AuthMode>("login");
  const [registerEmail, setRegisterEmail] = useState("");
  const [registerDisplayName, setRegisterDisplayName] = useState("");
  const [registerPassword, setRegisterPassword] = useState("");
  const [loginEmail, setLoginEmail] = useState("");
  const [loginPassword, setLoginPassword] = useState("");

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
        ? "Entrada temporal lista para crear o retomar una ciudad"
        : "Entrá con cuenta o invitado para empezar a jugar";

  return (
    <section
      className="screen landing-screen landing-cinematic session-screen"
      style={{ "--landing-backdrop": `url("${skylineBackdrop}")` } as CSSProperties}
    >
      <div className="landing-cinematic__backdrop" />
      <div className="landing-cinematic__veil" />

      <header className="landing-topbar">
        <div>
          <p className="eyebrow">PulseCity</p>
          <strong className="landing-brand">Session Control</strong>
        </div>

        <div className="topbar-status">
          <span>{identityLabel}</span>
          <small>{props.restoringSession ? "Restaurando..." : sessionSummary}</small>
        </div>
      </header>

      <div className="landing-cinematic__content session-cinematic__content">
        <section className="session-hero-panel">
          <div className="hero-copy-block">
            <p className="eyebrow">Identidad del jugador</p>
            <h1>Primero entra quien va a dirigir la ciudad.</h1>
            <p className="landing-copy">
              La cuenta resuelve biblioteca, continuidad y ownership. El acceso invitado
              mantiene baja la fricción cuando solo querés empezar.
            </p>
          </div>

          <div className="status-banner">
            <strong>{props.restoringSession ? "Restaurando sesión" : identityLabel}</strong>
            <span>{props.status}</span>
          </div>
        </section>

        <section className="session-access-panel">
          <div className="panel-heading">
            <div>
              <p className="eyebrow">Acceso</p>
              <h2>Entrar a PulseCity</h2>
            </div>
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
              <button
                type="button"
                className="text-action"
                disabled={props.restoringSession}
                onClick={props.onForgotPassword}
              >
                Olvide mi contraseña
              </button>
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

          <div className="session-divider" aria-hidden="true">
            <span />
            <small>o entrar sin cuenta</small>
            <span />
          </div>

          <button
            type="button"
            className="landing-cta session-guest-cta"
            onClick={props.onCreateGuestSession}
            disabled={props.creatingGuestSession || props.restoringSession}
          >
            <span>{props.activeAuthKind === "guest" ? "Seguir como invitado" : "Jugar como invitado"}</span>
            <small>
              {props.creatingGuestSession
                ? "Creando acceso temporal..."
                : "Entrar sin cuenta y migrar la partida mas adelante"}
            </small>
          </button>

          <div className="landing-utility-actions">
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
        </section>
      </div>
    </section>
  );
}
