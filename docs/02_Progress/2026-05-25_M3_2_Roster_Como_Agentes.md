# 2026-05-25 — M3.2 Roster como agentes

## Objetivo

Convertir los 15 jugadores propios del roster en agentes con capa emocional propia, sin mover ownership contractual fuera de `team-service`.

## Cambios realizados

- Se agrego `agent_player_states` como tabla propia de `agent-service`.
- Se agrego la migracion `016_create_agent_player_states.sql`.
- `agent-service` lee `team_roster_players` solo para sembrar identidad compartida (`player_id`, nombre, posicion, rating).
- Los estados iniciales incluyen `emotional_state`, `satisfaction`, `loyalty`, `ego`, `competitive_drive` y `city_connection`.
- Al cargar/inicializar una partida, `agent-service` asegura los 15 jugadores como agentes emocionales.
- Al recibir `partido.terminado`, `agent-service` usa el box score para actualizar la capa emocional de cada jugador propio.
- `agent-service` publica `roster.patch`.
- `gateway` reenvia `roster.patch` por WebSocket.
- El frontend tipa y aplica el delta en `RosterClientStates`.

## Decision tecnica

No se escriben tablas de `team-service` desde `agent-service`.

La tabla `team_roster_players` se usa como fuente de lectura para sincronizar identidad. Todo estado emocional vive en `agent_player_states`, propiedad de `agent-service`.

## Pruebas

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
npm run build --prefix frontend
make build
```

## Pendiente siguiente

Siguiente mini milestone recomendado: `M3.3 — Tabla de relaciones inter-agente`.
