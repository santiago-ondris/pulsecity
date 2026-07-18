# M3.P2g-f — Despacho de calendario

## Objetivo

Partir de `team-service/internal/persistence/store.go` el flujo cohesivo que reclama un partido programado y prepara el evento completo para el simulador, sin incluir el procesamiento posterior del resultado.

## Cambios

- Se creo `services/team-service/internal/persistence/schedule_dispatch.go`.
- Se movio `DispatchScheduledMatchForDate`.
- Se movieron `storedScheduleMatch`, `loadFranchise`, `loadOpponent` y `loadPlayerMatchStates`.
- `loadRoster` y `loadSeasonRecord` siguen definidos en sus dominios y se reutilizan desde el nuevo archivo.
- `ApplyMatchFinished` permanece sin cambios en `store.go`.

## Decision de frontera

El corte termina cuando se construye `MatchScheduledEvent`. Procesar `partido.terminado`, persistir box scores y recalcular el record forman otro flujo transaccional y quedan para el siguiente mini milestone.

`loadPlayerMatchStates` vive con la preparacion porque ensambla carga reciente y estado emocional para el payload stateless de `match-service`. El modulo de lesiones puede seguir reutilizandolo dentro del mismo paquete.

## Resultado

- `store.go`: 725 → 517 lineas.
- `schedule_dispatch.go`: 218 lineas.
- Se mantiene el claim atomico de un unico partido por fecha y partida.
- Se preservan el orden determinista, el caso sin partido y el commit unico posterior al armado del evento.
- El payload conserva equipos local/visitante, rival, roster, estados de jugadores y record de temporada.
- No cambiaron schema, SQL, eventos, ownership ni asserts de tests.

## Verificacion

- `gofmt`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`
- `git diff --check`

## Siguiente corte recomendado

`M3.P2g-g — Resultado de partidos y record`: separar `ApplyMatchFinished`, la persistencia de box scores y el recalculo/lectura del record de temporada. Con ese corte, `store.go` quedaria concentrado en schema e inicializacion de temporada.
