# Agent Events

Contratos de agentes para Milestone 2. El dueño del estado emocional es `agent-service`.

## Subjects

- `agente.estado_cambio`
- `agente.relacion_cambio`
- `agente.evento_critico`
- `agente.consulta_iniciada`
- `agente.respuesta_generada`

## Agentes core M2

- `owner`
- `head_coach`
- `cfo`
- `scouting_director`
- `sports_psychologist`

## `agente.estado_cambio`

Publicado por `agent-service` cuando un agente core cambia estado por un hecho del juego.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:04Z",
  "schema_version": 1,
  "simulated_date": "2026-10-22",
  "agent_id": "head_coach",
  "source_event_id": "uuid",
  "source_subject": "partido.terminado",
  "mood": "calm|concerned|excited|frustrated|pressured",
  "state": {
    "trust": 0.16,
    "satisfaction": 0.08,
    "pressure": 0.22
  },
  "summary": "El coach gana confianza tras una victoria cerrada."
}
```

## `agente.evento_critico`

Publicado por `agent-service` cuando una variable cruza un umbral relevante.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:04Z",
  "schema_version": 1,
  "simulated_date": "2026-10-22",
  "agent_id": "owner",
  "severity": "warning|critical",
  "source_event_id": "uuid",
  "source_subject": "partido.terminado",
  "title": "La paciencia del owner cae",
  "summary": "Una racha negativa empieza a afectar la confianza deportiva."
}
```

Notas:

- `state` permite variables distintas por agente sin abrir todavia una matriz de relaciones.
- los valores numericos se mantienen en rango `-1.0` a `1.0`, salvo variables que el dominio documente distinto.
- el `gateway` traduce `agente.estado_cambio` a `agent.patch`.

## `agente.relacion_cambio`

Publicado por `agent-service` cuando una relacion canonica cambia por un evento sistemico.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-25T00:00:04Z",
  "schema_version": 1,
  "simulated_date": "2026-10-22",
  "agent_a_id": "head_coach",
  "agent_b_id": "head_analytics",
  "trust": -0.22,
  "trend": "deteriorating",
  "last_event": "La derrota amplia reabre la tension entre datos y decisiones de cancha.",
  "short_history": [
    "Guerra fria entre ojo y dato",
    "La derrota amplia reabre la tension entre datos y decisiones de cancha."
  ],
  "source_event_id": "match-finished-match-1",
  "source_subject": "partido.terminado"
}
```

Notas:

- la relacion se identifica por el par canonico `(agent_a_id, agent_b_id)`.
- la idempotencia se controla por `(game_id, relationship_key, source_event_id)`.
- el `gateway` traduce `agente.relacion_cambio` a `relations.patch`.

## `agente.consulta_iniciada`

Publicado por `gateway` cuando el GM envia un mensaje directo a un agente.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-29T00:00:04Z",
  "schema_version": 1,
  "conversation_id": "chat-uuid",
  "agent_id": "head_coach",
  "sender": "gm",
  "message": "Como ves el vestuario?"
}
```

Notas:

- `narrative-service` consume este evento, ensambla contexto real y persiste el historial.
- En M3.5 la respuesta es stub deterministica. En M3.6 este mismo contrato alimenta el LLM real.

## `agente.respuesta_generada`

Publicado por `narrative-service` cuando ya existe una respuesta del agente.

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
    "generation": "stub",
    "source_subject": "agente.consulta_iniciada"
  },
  "created_at": "2026-05-29T00:00:05Z"
}
```

## Variables iniciales M2.11

`owner`:

- `sporting_trust`
- `business_trust`
- `patience_remaining`
- `satisfaction`

`head_coach`:

- `gm_trust`
- `roster_satisfaction`
- `results_pressure`
- `locker_room_relationship`

`cfo`:

- `financial_trust`
- `budget_alert`
- `financial_conservatism`

`scouting_director`:

- `criteria_trust`
- `motivation`
- `perceived_precision`

`sports_psychologist`:

- `locker_room_climate`
- `emotional_alert`
- `player_trust`
