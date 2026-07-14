# Player Events

Contratos de disponibilidad de jugadores para Milestone 3.

## Subjects

- `jugador.lesionado`
- `jugador.recuperado`

## `jugador.lesionado`

Publicado por `team-service` cuando un jugador propio sufre una lesion sistemica por carga acumulada.

```json
{
  "event_id": "player-injured-game-1-match-004-game-1-player-01",
  "game_id": "game-1",
  "occurred_at": "2026-10-28T00:00:00Z",
  "schema_version": 1,
  "injury_id": "injury-game-1-match-004-game-1-player-01",
  "player_id": "game-1-player-01",
  "severity": "minor",
  "estimated_days_out": 5,
  "injured_on": "2026-10-28",
  "expected_recovery_date": "2026-11-02",
  "reason": "workload_accumulation",
  "source_match_id": "game-1-match-004",
  "workload_score": 112
}
```

Campos:

- `severity`: `minor`, `moderate` o `major`.
- `estimated_days_out`: dias simulados estimados fuera.
- `expected_recovery_date`: fecha simulada en la que `team-service` puede publicar recuperacion automatica.
- `reason`: motivo sistemico; por ahora `workload_accumulation`.
- `reason` tambien puede ser `forced_return_reaggravation` cuando el GM fuerza un alta antes de la fecha recomendada y el jugador vuelve a lesionarse.
- `workload_score`: minutos recientes usados para evaluar riesgo.

## `jugador.recuperado`

Publicado por `team-service` cuando una lesion activa llega a su fecha de recuperacion.

```json
{
  "event_id": "player-recovered-injury-game-1-match-004-game-1-player-01",
  "game_id": "game-1",
  "occurred_at": "2026-11-02T00:00:00Z",
  "schema_version": 1,
  "injury_id": "injury-game-1-match-004-game-1-player-01",
  "player_id": "game-1-player-01",
  "recovered_on": "2026-11-02"
}
```
