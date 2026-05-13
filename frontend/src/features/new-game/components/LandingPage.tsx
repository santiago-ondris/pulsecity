import { managementModeLabel, scenarioById } from "../helpers";
import type { GameSummary } from "../../../types";

interface LandingPageProps {
  creatingGuestSession: boolean;
  guestGames: GameSummary[];
  guestReady: boolean;
  guestToken: string;
  status: string;
  onContinueGame: (gameId: string) => void;
  onCreateGuestSession: () => void;
  onStart: () => void;
}

export function LandingPage(props: LandingPageProps) {
  return (
    <section className="screen landing-screen">
      <div className="landing-hero">
        <p className="eyebrow">PulseCity</p>
        <h1>Antes de fundar una ciudad, primero queda definido quién entra a jugar.</h1>
        <p className="landing-copy">
          El corte de hoy abre la puerta mínima de Milestone 1: sesión invitada real, token
          persistido y partidas asociadas correctamente desde el primer click.
        </p>

        <div className="landing-actions">
          {!props.guestReady ? (
            <button
              type="button"
              className="primary-action landing-action"
              onClick={props.onCreateGuestSession}
              disabled={props.creatingGuestSession}
            >
              {props.creatingGuestSession ? "Creando acceso invitado..." : "Jugar como invitado"}
            </button>
          ) : (
            <button type="button" className="primary-action landing-action" onClick={props.onStart}>
              Nueva partida
            </button>
          )}
        </div>
      </div>

      <div className="landing-grid">
        <article className="landing-note">
          <span>01</span>
          <strong>Acceso</strong>
          <p>{props.guestReady ? `Invitado activo: ${props.guestToken}` : "Todavía sin sesión."}</p>
        </article>
        <article className="landing-note">
          <span>02</span>
          <strong>Ownership</strong>
          <p>Las partidas nuevas quedan ligadas al token invitado y ya no flotan sueltas.</p>
        </article>
        <article className="landing-note">
          <span>03</span>
          <strong>Continuidad</strong>
          <p>La entrada puede listar y reanudar fundaciones ya asociadas a este invitado.</p>
        </article>
        <article className="landing-note">
          <span>04</span>
          <strong>Próximo paso</strong>
          <p>Sobre esta base ya se puede sumar login y registro sin rehacer el flujo.</p>
        </article>
      </div>

      {props.guestReady ? (
        <section className="guest-session-panel">
          <div className="panel-header">
            <p className="eyebrow">Sesión invitada</p>
            <h2>Partidas asociadas</h2>
          </div>
          <p className="copy">{props.status}</p>

          {props.guestGames.length === 0 ? (
            <div className="guest-empty-state">
              <strong>Este invitado todavía no tiene partidas.</strong>
              <p>La próxima fundación quedará registrada automáticamente bajo este token.</p>
            </div>
          ) : (
            <ul className="guest-game-list">
              {props.guestGames.map((game) => {
                const scenario = scenarioById(game.initial_scenario);

                return (
                  <li key={game.game_id} className="guest-game-card">
                    <div>
                      <p className="eyebrow">{game.city_name}</p>
                      <h3>{game.franchise_name}</h3>
                      <p className="guest-game-meta">
                        {scenario.label} · {managementModeLabel(game.city_management_mode)}
                      </p>
                      <p className="guest-game-meta">Estado: {game.status}</p>
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
