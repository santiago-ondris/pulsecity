# M3.9a — Lesiones sistemicas por carga

Fecha: 2026-07-14

## Objetivo

Crear la base operativa de lesiones para que el Medico y el S&C Coach tengan un sistema real sobre el cual actuar.

## Cambios

- `team-service` agrega tabla `team_injuries`.
- `team-service` evalua riesgo de lesion al aplicar `partido.terminado`.
- El riesgo usa carga reciente, minutos del partido y estado emocional negativo como modificador.
- La evaluacion es deterministica por `game_id`, `match_id` y `player_id`.
- Si hay lesion:
  - se persiste en `team_injuries`
  - el jugador pasa a `roster_status = injured`
  - se publica `jugador.lesionado`
- En cada `tiempo.dia_avanzado`, `team-service` recupera lesiones vencidas antes de despachar partidos del dia.
- Si hay recuperacion:
  - la lesion recibe `recovered_on`
  - el jugador vuelve a `roster_status = active`
  - se publica `jugador.recuperado`
- `gateway` escucha `jugador.lesionado` y `jugador.recuperado`, y los transforma en `roster.patch` delta-only para el frontend.
- El frontend queda tipado para disponibilidad opcional de jugadores.

## Contratos

- `shared/events/players.md`
- `shared/events/websocket_deltas.md`

## Decision

Este corte no implementa todavia la decision del GM ante alerta medica.

Primero quedo cerrada la base sistemica: carga -> lesion -> baja -> recuperacion. El flujo narrativo Medico ↔ GM ↔ Coach se implementa encima en `M3.9b`, usando estos eventos como disparadores.

## Pruebas

- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`
- `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`
- `npm run build --prefix frontend`

## Resultado

`M3.9a` queda cerrado funcionalmente.

Siguiente paso: `M3.9b — alerta medica y decision del GM`.
