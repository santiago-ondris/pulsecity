# M3.P2g-d — Aceptacion de trades

## Objetivo de la sesion

Continuar la particion de `team-service/internal/persistence/store.go` aislando el flujo atomico de aceptacion de trades.

## Cambios realizados

- Se creo `services/team-service/internal/persistence/trade_acceptance.go`.
- Se movieron desde `store.go`:
  - `ApplyTradeAcceptance`
  - `loadTrade`
  - `storedTrade`
  - `tradeRejectedFromDecision`
  - `tradeRosterPatch`
- `trade_proposals.go` continua usando `tradeRejectedFromDecision` dentro del mismo paquete.
- Los helpers de roster y salary cap se consumen desde sus archivos dedicados.

## Decision tomada

La aceptacion permanece en un solo metodo transaccional: marcar al jugador saliente, insertar al entrante, cerrar la propuesta y recalcular/persistir salary cap deben confirmarse o revertirse juntos.

Los eventos y deltas derivados se construyen en el mismo modulo porque representan el resultado de esa transaccion, no nuevas mutaciones independientes.

## Resultado

- `store.go`: 1251 -> 1033 lineas.
- `trade_acceptance.go`: 228 lineas.
- La particion del dominio trade dentro del store queda cerrada.
- Sin cambios en schema, SQL, transacciones, API, ownership o comportamiento.
- Suite, build y vet de `team-service` pasan.

## Verificacion

```bash
gofmt -w services/team-service/internal/persistence/store.go \
  services/team-service/internal/persistence/trade_acceptance.go \
  services/team-service/internal/persistence/trade_proposals.go
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...
git diff --check
```

## Pendiente siguiente

La particion de `team-service` continua pendiente. El siguiente corte recomendado es `M3.P2g-e`: separar decisiones medicas y persistencia de lesiones antes de calendario, partidos y temporada.
