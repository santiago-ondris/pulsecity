# Time Events

Contratos de tiempo para Milestone 2. El dueño del loop es `agent-service`.

## Subjects

- `tiempo.sesion_iniciada`
- `tiempo.sesion_terminada`
- `tiempo.velocidad_cambiada`
- `tiempo.pausa_activada`
- `tiempo.dia_avanzado`

## Campos comunes

Todo evento de tiempo incluye:

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:00Z",
  "schema_version": 1
}
```

## `tiempo.sesion_iniciada`

Publicado por `gateway` cuando se abre una sesion WebSocket activa para una partida.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:00Z",
  "schema_version": 1,
  "session_id": "uuid",
  "client_id": "browser-session-id"
}
```

## `tiempo.sesion_terminada`

Publicado por `gateway` cuando se cierra la ultima conexion activa de esa sesion.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:05:00Z",
  "schema_version": 1,
  "session_id": "uuid",
  "reason": "client_closed|timeout|server_shutdown"
}
```

## `tiempo.velocidad_cambiada`

Publicado por `gateway`. Consumido por `agent-service`.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:01:00Z",
  "schema_version": 1,
  "speed": 5
}
```

Velocidades validas en M2: `1`, `5`, `20`.

## `tiempo.pausa_activada`

Publicado por `gateway`. Consumido por `agent-service`.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:02:00Z",
  "schema_version": 1,
  "paused": true
}
```

## `tiempo.dia_avanzado`

Publicado por `agent-service` cuando procesa uno o mas dias simulados.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:03:00Z",
  "schema_version": 1,
  "simulated_date": "2026-10-22",
  "speed": 5,
  "days_processed": 1
}
```

Notas:

- `simulated_date` representa el dia ya procesado.
- `days_processed` normalmente es `1`; queda explicito para soportar acumulaciones futuras.
- el frontend no recibe este evento crudo: el `gateway` lo traduce a `time.patch`.
