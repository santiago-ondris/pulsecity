# M3.P2g-b — Persistencia de roster

## Objetivo de la sesion

Continuar la particion de `team-service/internal/persistence/store.go` aislando operaciones propias del roster.

## Cambios realizados

- Se creo `services/team-service/internal/persistence/roster.go`.
- Se movieron desde `store.go`:
  - `ApplyRosterPatch`
  - `loadRoster`
  - `loadRosterPlayer`
- Los callers existentes de trades y calendario no cambiaron.
- Las funciones continúan usando `queryer`, por lo que operan dentro de las transacciones de sus callers.

## Decision tomada

La insercion del roster inicial permanece en `SaveInitialSeason`: forma parte de una transaccion mayor que crea franquicia, rivales, calendario, record y salary cap. Extraer solo ese fragmento no aportaria ownership adicional y haria menos legible la atomicidad.

El `roster.patch` emocional se persiste aqui porque `team-service` lo consume para su vista de estado de partido; no modifica el ownership emocional de `agent-service`.

## Resultado

- `store.go`: 1516 -> 1417 lineas.
- `roster.go`: 107 lineas.
- Sin cambios en schema, SQL, transacciones, API, ownership o comportamiento.
- Suite, build y vet de `team-service` pasan.

## Verificacion

```bash
gofmt -w services/team-service/internal/persistence/store.go \
  services/team-service/internal/persistence/roster.go
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...
git diff --check
```

## Pendiente siguiente

La particion de `team-service` continua pendiente. El siguiente corte recomendado es `M3.P2g-c`: separar propuestas de trade y sus helpers de persistencia antes de abordar aceptacion, lesiones y temporada.
