# M3.P2g-g — Resultado de partidos y record

## Objetivo

Partir de `team-service/internal/persistence/store.go` el flujo transaccional que consume `partido.terminado`, persiste el box score y actualiza el record de temporada.

## Cambios

- Se creo `services/team-service/internal/persistence/match_results.go`.
- Se movio `ApplyMatchFinished`.
- Se movieron `recalculateSeasonRecord` y `loadSeasonRecord`.
- `createInjuriesForMatch` permanece en `injuries.go` y sigue participando de la misma transaccion.
- La interfaz `queryer` queda en `store.go` porque es compartida por todos los archivos de persistencia.

## Decision de frontera

El resultado deportivo, sus box scores y el record forman un unico flujo: el record solo se recalcula despues de persistir el resultado, y un retry debe devolver el record existente sin repetir lesiones ni mutaciones.

La generacion de lesiones no se duplico en este archivo. Permanece en su dominio y se invoca usando la misma transaccion, conservando atomicidad sin mezclar responsabilidades.

## Resultado

- `store.go`: 517 → 346 lineas.
- `match_results.go`: 180 lineas.
- `store.go` queda concentrado en conexion, schema e inicializacion de temporada.
- Se preservan idempotencia por `status <> 'final'`, upsert por `(match_id, player_id)` y commit unico.
- Se conserva el calculo de victorias, derrotas, puntos a favor/en contra y ultimo partido segun localia.
- No cambiaron schema, SQL, eventos, ownership ni asserts de tests.

## Verificacion

- `gofmt`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`
- `git diff --check`

## Siguiente corte recomendado

Cerrar `M3.P2g` con una revision integral de la nueva particion de `team-service`: mapa de responsabilidades, dependencias internas, conteo de lineas y suite completa. No hace falta seguir partiendo `store.go`; sus 346 lineas restantes corresponden a schema e inicializacion cohesivos.
