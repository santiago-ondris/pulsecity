export function LandingPage({ onStart }: { onStart: () => void }) {
  return (
    <section className="screen landing-screen">
      <div className="landing-hero">
        <p className="eyebrow">PulseCity</p>
        <h1>La ciudad nace después de una sola decisión: fundar la franquicia.</h1>
        <p className="landing-copy">
          Este ya no es un formulario largo disfrazado de homepage. Desde acá solo hay un camino:
          empezar una nueva partida y atravesar el onboarding en orden.
        </p>

        <div className="landing-actions">
          <button type="button" className="primary-action landing-action" onClick={onStart}>
            Nueva partida
          </button>
        </div>
      </div>

      <div className="landing-grid">
        <article className="landing-note">
          <span>01</span>
          <strong>Identidad</strong>
          <p>Nombre, sigla y colores antes de cualquier otra decisión.</p>
        </article>
        <article className="landing-note">
          <span>02</span>
          <strong>Escenario</strong>
          <p>El arranque competitivo define el tono inicial del mundo.</p>
        </article>
        <article className="landing-note">
          <span>03</span>
          <strong>Gobierno</strong>
          <p>Después se fija la relación real entre franquicia y ciudad.</p>
        </article>
        <article className="landing-note">
          <span>04</span>
          <strong>Fundación</strong>
          <p>Recién ahí empieza la ceremonia del mapa y la llamada del Owner.</p>
        </article>
      </div>
    </section>
  );
}
