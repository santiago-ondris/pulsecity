# GM Decision Events

Contratos de memoria de decisiones del GM para Milestone 3.

## Subjects

- `decision.gm_registrada`

## `decision.gm_registrada`

Publicado cuando el jugador toma una decision relevante que los agentes podran referenciar mas adelante.

```json
{
  "event_id": "decision-owner-intro-game-1",
  "game_id": "game-1",
  "occurred_at": "2026-05-25T00:00:04Z",
  "schema_version": 1,
  "decision_id": "owner-intro-game-1",
  "kind": "owner_intro_response",
  "payload": {
    "choice_id": "win_now",
    "choice_label": "Competir desde el dia uno"
  },
  "simulated_date": "2026-10-01",
  "agents_affected": ["owner", "president_basketball_ops", "ceo_business_ops"],
  "source_event_id": "owner-intro-game-1",
  "source_subject": "narrativa.respuesta_gm"
}
```

Notas:

- `agent-service` persiste el evento en `gm_decisions_log`.
- `payload` debe contener solo el contexto minimo necesario para reconstruir la decision.
- `gm_decisions_log` es append-only e idempotente por `(game_id, decision_id)`.
- El log no se expone directamente al frontend; se usara por lecturas agregadas para contexto de agentes.
