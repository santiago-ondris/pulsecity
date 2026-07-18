# M3.P2f-e — Trades del gateway

## Objetivo de la sesion

Continuar la particion de `gateway/internal/handlers/http.go` aislando las entradas HTTP del dominio trade.

## Cambios realizados

- Se creo `services/gateway/internal/handlers/trades.go`.
- Se movieron desde `http.go`:
  - `proposeTrade`
  - `acceptTrade`
- `RegisterRoutes` conserva exactamente las mismas rutas y metodos.
- Los handlers continuan compartiendo autorizacion, store y bus mediante `Dependencies` dentro del mismo paquete.

## Decision tomada

El gateway solo valida la accion del GM y la traduce a `decision.gm_registrada`. No evalua la oferta, no genera al jugador entrante y no muta roster, contratos ni salary cap.

La evaluacion del GM rival sigue en `agent-service`; la aceptacion y mutacion contractual sigue en `team-service`, respetando ownership estricto.

## Resultado

- `http.go`: 944 -> 770 lineas.
- `trades.go`: 184 lineas.
- Sin cambios en rutas HTTP, eventos NATS, payloads, ownership o comportamiento.
- Suite, build y vet del gateway pasan.

## Verificacion

```bash
gofmt -w services/gateway/internal/handlers/http.go \
  services/gateway/internal/handlers/trades.go
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...
git diff --check
```

## Pendiente siguiente

La particion del gateway continua pendiente. El siguiente corte recomendado es `M3.P2f-f`: separar control de tiempo y luego dejar `http.go` como shell de rutas/health, moviendo el HTML debug embebido a un archivo propio.
