# M3.P2g-c — Propuestas de trade

## Objetivo de la sesion

Continuar la particion de `team-service/internal/persistence/store.go` aislando el flujo de propuestas de trade.

## Cambios realizados

- Se creo `services/team-service/internal/persistence/trade_proposals.go`.
- Se movieron desde `store.go`:
  - `ApplyTradeProposal`
  - `loadTradeStatus`
  - `insertProposedTrade`
  - `insertRejectedTrade`
  - `insertRejectedTradeWithPlayer`
  - `parsePositiveInt`
- Los helpers de roster y salary cap se siguen consumiendo desde sus archivos dedicados dentro del mismo paquete.

## Decision tomada

`tradeRejectedFromDecision` permanece compartido en `store.go` porque la aceptacion tambien construye rechazos. Se movera junto con el dominio trade cuando el flujo de aceptacion sea extraido.

No se incluyeron `loadTrade` ni `tradeRosterPatch`: pertenecen al corte de aceptacion y moverlos ahora mezclaria dos unidades de trabajo.

## Resultado

- `store.go`: 1417 -> 1251 lineas.
- `trade_proposals.go`: 176 lineas.
- Sin cambios en schema, SQL, transacciones, API, ownership o comportamiento.
- Suite, build y vet de `team-service` pasan.

## Verificacion

```bash
gofmt -w services/team-service/internal/persistence/store.go \
  services/team-service/internal/persistence/trade_proposals.go
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...
git diff --check
```

## Pendiente siguiente

La particion de `team-service` continua pendiente. El siguiente corte recomendado es `M3.P2g-d`: separar aceptacion de trades, `loadTrade`, el tipo persistido y la construccion del `roster.patch`.
