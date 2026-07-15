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

### `medical_decision`

Publicado por `gateway` cuando el GM responde una recomendacion medica.

```json
{
  "event_id": "decision-medical-game-1-injury-game-1-match-004-game-1-player-01",
  "game_id": "game-1",
  "occurred_at": "2026-10-28T00:00:04Z",
  "schema_version": 1,
  "decision_id": "medical-injury-game-1-match-004-game-1-player-01",
  "kind": "medical_decision",
  "payload": {
    "injury_id": "injury-game-1-match-004-game-1-player-01",
    "player_id": "game-1-player-01",
    "choice_id": "force_return",
    "choice_label": "Forzar alta anticipada"
  },
  "simulated_date": "2026-10-30",
  "agents_affected": ["team_doctor", "strength_conditioning_coach", "head_coach"],
  "source_event_id": "injury-game-1-match-004-game-1-player-01",
  "source_subject": "jugador.lesionado"
}
```

Opciones validas:

- `rest`
- `reduce_minutes`
- `ignore_doctor`
- `force_return`

### `trade_proposal`

Publicado por `gateway` cuando el GM inicia una propuesta de trade.

```json
{
  "event_id": "decision-trade-uuid",
  "game_id": "game-1",
  "occurred_at": "2026-11-01T00:00:04Z",
  "schema_version": 1,
  "decision_id": "trade-uuid",
  "kind": "trade_proposal",
  "payload": {
    "proposal_id": "trade-uuid",
    "rival_team_id": "bos",
    "offered_player_id": "game-1-player-06",
    "requested_position": "PG",
    "incoming_salary": "12000000"
  },
  "simulated_date": "2026-11-01",
  "agents_affected": ["director_player_personnel", "cfo", "rival_gm_bos"],
  "source_event_id": "trade-uuid",
  "source_subject": "trade.propuesta_gm"
}
```

Notas:

- `team-service` valida ownership de roster y salary cap antes de publicar `trade.propuesta_enviada`.
- `agent-service` evalua el perfil persistido del GM rival solo despues de recibir `trade.propuesta_enviada`.

### `trade_acceptance`

Publicado por `gateway` cuando el GM acepta cerrar una propuesta abierta o una contraoferta.

```json
{
  "event_id": "decision-accept-trade-uuid",
  "game_id": "game-1",
  "occurred_at": "2026-11-01T00:01:04Z",
  "schema_version": 1,
  "decision_id": "accept-trade-uuid",
  "kind": "trade_acceptance",
  "payload": {
    "proposal_id": "trade-uuid",
    "accepted_additional_asset": "second_round_pick"
  },
  "simulated_date": "2026-11-01",
  "agents_affected": ["director_player_personnel", "cfo"],
  "source_event_id": "trade-uuid",
  "source_subject": "trade.contraoferta"
}
```

Notas:

- `team-service` vuelve a validar que la propuesta siga abierta y que el jugador ofrecido siga activo.
- Si se acepta, `team-service` muta roster y cap antes de publicar `trade.aceptada`.
- Reintentar la misma aceptacion no duplica jugadores ni vuelve a publicar eventos desde `team-service`.
