# WebSocket Deltas

Contratos de salida del `gateway` al frontend para Milestone 2.

Regla: WebSocket envia deltas, no estado completo. La rehidratacion se hace por snapshot REST o snapshot inicial especifico del flujo.

## Tipos esperados

- `time.patch`
- `season.patch`
- `match.result`
- `city.patch`
- `agent.patch`
- `narrative.event`

## `time.patch`

```json
{
  "type": "time.patch",
  "game_id": "uuid",
  "patch": {
    "simulated_date": "2026-10-22",
    "speed": 5,
    "paused": false,
    "session_active": true
  }
}
```

## `season.patch`

```json
{
  "type": "season.patch",
  "game_id": "uuid",
  "patch": {
    "record": { "wins": 7, "losses": 5 },
    "next_match_id": "uuid"
  }
}
```

## `match.result`

```json
{
  "type": "match.result",
  "game_id": "uuid",
  "match_id": "uuid",
  "result": {
    "simulated_date": "2026-10-22",
    "home_score": 112,
    "away_score": 106,
    "winner_team_id": "pulsecity",
    "key_moments": []
  }
}
```

## `city.patch`

```json
{
  "type": "city.patch",
  "game_id": "uuid",
  "patch": {
    "fan_sentiment": 0.58,
    "ticket_sales": 0.63,
    "stadium_district_land_value": 1.12
  }
}
```

## `agent.patch`

```json
{
  "type": "agent.patch",
  "game_id": "uuid",
  "agent_id": "head_coach",
  "patch": {
    "mood": "calm",
    "trust": 0.16,
    "satisfaction": 0.08,
    "pressure": 0.22,
    "summary": "El coach gana confianza tras una victoria cerrada."
  }
}
```

## `narrative.event`

Este contrato reutiliza el payload visible de `narrativa.evento_generado` con `type = "narrative.event"`.
