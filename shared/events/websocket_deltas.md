# WebSocket Deltas

Contratos de salida del `gateway` al frontend para Milestone 2.

Regla: WebSocket envia deltas, no estado completo. La rehidratacion se hace por snapshot REST o snapshot inicial especifico del flujo.

## Tipos esperados

- `time.patch`
- `season.patch`
- `match.result`
- `city.patch`
- `agent.patch`
- `roster.patch`
- `relations.patch`
- `narrative.event`
- `chat.message`

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
    "state": {
      "gm_trust": 0.16,
      "roster_satisfaction": 0.08,
      "results_pressure": 0.22
    },
    "summary": "El coach gana confianza tras una victoria cerrada.",
    "simulated_date": "2026-10-22",
    "source_event_id": "uuid",
    "source_subject": "partido.terminado"
  }
}
```

## `roster.patch`

```json
{
  "type": "roster.patch",
  "game_id": "uuid",
  "patch": {
    "simulated_date": "2026-10-22",
    "source_event_id": "match-finished-match-1",
    "source_subject": "partido.terminado",
    "players": [
      {
        "player_id": "game-player-01",
        "emotional_state": "confident",
        "satisfaction": 0.12,
        "loyalty": 0.64,
        "ego": 0.58,
        "competitive_drive": 0.72,
        "city_connection": 0.38,
        "summary": "Mateo Cross procesa la victoria con rol alto y 26 puntos."
      }
    ]
  }
}
```

## `relations.patch`

```json
{
  "type": "relations.patch",
  "game_id": "uuid",
  "patch": {
    "simulated_date": "2026-10-22",
    "source_event_id": "match-finished-match-1",
    "source_subject": "partido.terminado",
    "relationships": [
      {
        "agent_a_id": "head_coach",
        "agent_b_id": "head_analytics",
        "trust": -0.22,
        "trend": "deteriorating",
        "last_event": "La derrota amplia reabre la tension entre datos y decisiones de cancha.",
        "short_history": []
      }
    ]
  }
}
```

## `narrative.event`

Este contrato reutiliza el payload visible de `narrativa.evento_generado` con `type = "narrative.event"`.

## `chat.message`

Delta emitido por `gateway` cuando `narrative-service` publica una respuesta de chat.

```json
{
  "type": "chat.message",
  "subject": "agente.respuesta_generada",
  "game_id": "uuid",
  "conversation_id": "chat-uuid",
  "message_id": "agent-response-uuid",
  "agent_id": "head_coach",
  "sender": "agent",
  "body": "[stub M3.5] Soy ...",
  "metadata": {
    "generation": "stub"
  },
  "created_at": "2026-05-29T00:00:05Z"
}
```

Notas:

- WebSocket sigue siendo delta-only. La rehidratacion de historial se agrega por REST cuando M3.5 se cierre por completo.
- El mensaje del GM se persiste en `agent_chat_history`, pero el frontend puede mostrarlo optimistamente al enviar.
