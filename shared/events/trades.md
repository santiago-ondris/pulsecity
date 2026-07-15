# Trade Events

Contratos de negociacion de trades para Milestone 3.

## Subjects

- `trade.propuesta_enviada`
- `trade.rechazada`
- `trade.contraoferta`
- `trade.aceptada`

## `trade.propuesta_enviada`

Publicado por `team-service` despues de validar que la propuesta del GM refiere a un jugador propio activo y no rompe el salary cap acotado.

```json
{
  "event_id": "trade-proposed-trade-uuid",
  "game_id": "game-1",
  "occurred_at": "2026-11-01T00:00:04Z",
  "schema_version": 1,
  "proposal_id": "trade-uuid",
  "simulated_date": "2026-11-01",
  "rival_team_id": "bos",
  "offered_player_id": "game-1-player-06",
  "offered_player_name": "Adrian Vale",
  "offered_salary": 9000000,
  "requested_position": "PG",
  "incoming_salary": 12000000,
  "cap_space_after": -15000000
}
```

## `trade.rechazada`

Publicado por `team-service` si la propuesta falla validacion local, o por `agent-service` si el GM rival rechaza segun perfil y necesidades.

```json
{
  "event_id": "trade-rejected-trade-uuid",
  "game_id": "game-1",
  "occurred_at": "2026-11-01T00:00:05Z",
  "schema_version": 1,
  "proposal_id": "trade-uuid",
  "simulated_date": "2026-11-01",
  "rival_team_id": "bos",
  "reason": "rival_needs_mismatch",
  "detail": "Elliot Walsh no ve encaje claro con sus necesidades actuales."
}
```

## `trade.contraoferta`

Publicado por `agent-service` cuando el GM rival no acepta el paquete inicial pero mantiene la negociacion abierta.

```json
{
  "event_id": "trade-countered-trade-uuid",
  "game_id": "game-1",
  "occurred_at": "2026-11-01T00:00:05Z",
  "schema_version": 1,
  "proposal_id": "trade-uuid",
  "simulated_date": "2026-11-01",
  "rival_team_id": "bos",
  "requested_position": "PG",
  "additional_asset_required": "second_round_pick",
  "detail": "Elliot Walsh no acepta el paquete inicial, pero deja abierta una contraoferta."
}
```

## `trade.aceptada`

Publicado por `team-service` despues de persistir el cierre del trade, marcar al jugador saliente como `traded`, materializar al jugador entrante y recalcular salary cap.

```json
{
  "event_id": "trade-accepted-trade-uuid",
  "game_id": "game-1",
  "occurred_at": "2026-11-01T00:01:05Z",
  "schema_version": 1,
  "proposal_id": "trade-uuid",
  "simulated_date": "2026-11-01",
  "rival_team_id": "bos",
  "outgoing_player_id": "game-1-player-06",
  "outgoing_player_name": "Adrian Vale",
  "incoming_player_id": "trade-uuid-incoming",
  "incoming_player_name": "Jalen Warren",
  "incoming_position": "PG",
  "incoming_rating": 76,
  "incoming_salary": 12000000,
  "accepted_additional_asset": "second_round_pick"
}
```

## Ownership

- `team-service` es dueño de `team_trades`, validacion de roster propio y salary cap.
- `agent-service` es dueño de `rival_gms` y de la evaluacion del GM rival.
- `gateway` solo publica la decision inicial del GM y transforma `trade.*` en `trade.patch` para WebSocket.

## Mutacion de roster

M3 no modela rosters completos de rivales. Al aceptar un trade:

- el jugador saliente queda `roster_status = traded`
- el jugador entrante se materializa deterministamente desde `proposal_id`, posicion requerida, salario entrante y rating del jugador saliente
- `team-service` publica `roster.patch`, `salary_cap.calculado` y `finance.patch`
- `agent-service` escucha `trade.aceptada`, crea/actualiza la capa emocional de los jugadores afectados y publica otro `roster.patch` emocional

## Limite M3.12

No hay UI especifica de negociacion avanzada ni assets reales de draft. `accepted_additional_asset` queda como string acotado para que el flujo sea jugable y para profundizarlo en cortes posteriores.
