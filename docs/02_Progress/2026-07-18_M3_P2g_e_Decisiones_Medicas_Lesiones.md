# M3.P2g-e — Decisiones medicas y lesiones

## Objetivo

Partir el siguiente dominio cohesivo de `team-service/internal/persistence/store.go` sin modificar comportamiento: decisiones medicas, recuperacion diaria y persistencia de lesiones.

## Cambios

- Se creo `services/team-service/internal/persistence/injuries.go`.
- Se movieron `ApplyMedicalDecision` y `RecoverPlayersForDate`.
- Se movieron `createInjuriesForMatch`, `forcedReturnReaggravation`, `insertInjury` y `parseDate`.
- `ApplyMatchFinished` sigue coordinando el resultado del partido y llama al helper medico dentro del mismo paquete.

## Decision de frontera

`loadPlayerMatchStates` permanece en `store.go`: aunque consulta estado medico, tambien prepara el payload completo que consume el simulador. Moverlo al modulo de lesiones mezclaria persistencia medica con preparacion de partidos.

## Resultado

- `store.go`: 1033 → 725 lineas.
- `injuries.go`: 318 lineas.
- El dominio medico queda separado sin agregar abstracciones ni modificar la API publica del paquete.
- Se preservan las transacciones, el force return, la recuperacion por fecha, la reagravacion determinista y la idempotencia de `injury_id`.
- No cambiaron schema, SQL, eventos, ownership ni asserts de tests.

## Verificacion

- `gofmt`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`
- `git diff --check`

## Siguiente corte recomendado

`M3.P2g-f — Despacho de calendario`: separar la seleccion del partido programado y la construccion del payload que se envia al simulador, manteniendo el procesamiento de resultados para un corte posterior.
