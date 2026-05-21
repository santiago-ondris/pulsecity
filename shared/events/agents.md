# Agent Events

Contratos de agentes para Milestone 2. El dueño del estado emocional es `agent-service`.

## Subjects

- `agente.estado_cambio`
- `agente.evento_critico`

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
