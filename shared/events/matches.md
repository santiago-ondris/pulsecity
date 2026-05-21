# Match Events

Contratos de partidos para Milestone 2.

## Subjects

- `partido.programado`
- `partido.iniciando`
- `partido.terminado`

## Campos comunes

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:00Z",
  "schema_version": 1
}
```

## Modelo de equipo

```json
{
  "team_id": "pulsecity|rival-generated-id",
  "name": "PulseCity Astrals",
  "abbreviation": "PCA",
  "rating": 78,
  "offense_rating": 80,
  "defense_rating": 76,
  "pace": 99,
  "home_court_advantage": 3
}
```

## `partido.programado`

Publicado por `team-service` cuando un partido queda listo para ser simulado por fecha.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:00Z",
  "schema_version": 1,
  "match_id": "uuid",
  "simulated_date": "2026-10-22",
  "home_team": {},
  "away_team": {},
  "seed": 123456789
}
```

## `partido.iniciando`

Publicado por `match-service` al aceptar una simulacion.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:01Z",
  "schema_version": 1,
  "match_id": "uuid",
  "simulated_date": "2026-10-22"
}
```

## `partido.terminado`

Publicado por `match-service`. Consumido por `team-service`, `agent-service`, `city-service`, `analytics-service`, `narrative-service` y `gateway`.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:02Z",
  "schema_version": 1,
  "match_id": "uuid",
  "simulated_date": "2026-10-22",
  "home_team": {},
  "away_team": {},
  "home_score": 112,
  "away_score": 106,
  "winner_team_id": "pulsecity",
  "seed": 123456789,
  "box_score": [
    {
      "player_id": "uuid",
      "team_id": "pulsecity",
      "minutes": 32,
      "points": 24,
      "rebounds": 8,
      "assists": 6,
      "steals": 1,
      "blocks": 0,
      "turnovers": 3
    }
  ],
  "key_moments": [
    {
      "quarter": 4,
      "clock": "02:14",
      "kind": "clutch_shot",
      "description": "El base titular anota un triple para romper la igualdad.",
      "team_id": "pulsecity",
      "player_id": "uuid"
    }
  ]
}
```

Reglas:

- `match-service` es stateless y no consulta otros servicios.
- mismo input + misma `seed` debe producir el mismo `partido.terminado`.
- los consumidores deben tratar `event_id` y `match_id` como claves de idempotencia.
