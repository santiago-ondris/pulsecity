# M3.P2g-a — Salary cap del store

## Objetivo de la sesion

Iniciar la particion de `team-service/internal/persistence/store.go` con la operacion autonoma de salary cap.

## Analisis previo

Salary cap no tiene hoy un CRUD propio dentro del store. Su unica operacion de persistencia independiente es `saveSalaryCap`; el calculo vive en dominio y los usos aparecen dentro de fundacion y trades.

Por eso el corte se mantuvo pequeno y no arrastro logica de otros agregados solo para reducir mas lineas.

## Cambios realizados

- Se creo `services/team-service/internal/persistence/salary_cap.go`.
- Se movio `saveSalaryCap` desde `store.go`.
- La funcion conserva `queryer` para operar tanto con pool como dentro de una transaccion.
- Los callers existentes no cambiaron.

## Decision tomada

El DDL de `team_salary_cap` continua dentro de `EnsureSchema`. Separar la declaracion completa del schema se tratara como un corte propio para evitar mezclar organizacion del bootstrap con persistencia de agregados.

`domain.CalculateSalaryCap` tampoco se mueve: pertenece a logica de negocio, no a persistencia.

## Resultado

- `store.go`: 1545 -> 1516 lineas.
- `salary_cap.go`: 37 lineas.
- SQL y parametros identicos al bloque original.
- Sin cambios en schema, transacciones, API, ownership o comportamiento.
- Suite, build y vet de `team-service` pasan.

## Verificacion

```bash
gofmt -w services/team-service/internal/persistence/store.go \
  services/team-service/internal/persistence/salary_cap.go
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...
git diff --check
```

## Pendiente siguiente

La particion de `team-service` continua pendiente. El siguiente corte recomendado es `M3.P2g-b`: separar persistencia de roster (`ApplyRosterPatch`, `loadRoster` y `loadRosterPlayer`) antes de trades, lesiones y temporada.
