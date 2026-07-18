# M3.P2f-f — Control de tiempo y shell HTTP

## Objetivo de la sesion

Cerrar la particion del gateway separando control de tiempo y data-as-code visual del shell HTTP.

## Cambios realizados

- Se creo `services/gateway/internal/handlers/time_control.go`.
- Se movieron desde `http.go`:
  - `updateTimeControl`
  - `validTimeSpeed`
- Se creo `services/gateway/internal/handlers/debug_page.go`.
- Se movio `debugHTML` sin modificar su contenido.
- `http.go` conserva:
  - `Dependencies`
  - `RegisterRoutes`
  - `debugPage`
  - `healthz`
  - `writeJSON`

## Decision tomada

El HTML debug se trata como data-as-code y vive separado del wiring HTTP. No se reviso ni modifico su diseño en este corte porque el objetivo es un refactor de estructura con comportamiento identico.

El control de tiempo permanece en el gateway como entrada del frontend, pero no persiste ni calcula simulacion: valida y publica eventos, y emite solo el delta correspondiente.

## Resultado

- `http.go`: 770 -> 57 lineas.
- `time_control.go`: 108 lineas.
- `debug_page.go`: 613 lineas de data-as-code visual.
- La particion del gateway dentro de `M3.P2` queda cerrada.
- Sin cambios en rutas HTTP, eventos NATS, payloads, ownership, contenido visual o comportamiento.
- Suite, build y vet del gateway pasan.

## Verificacion

```bash
gofmt -w services/gateway/internal/handlers/http.go \
  services/gateway/internal/handlers/time_control.go \
  services/gateway/internal/handlers/debug_page.go
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...
git diff --check
```

Tambien se comparo el bloque `debugHTML` nuevo contra el original y no hubo diferencias.

## Pendiente siguiente

`M3.P2` continua pendiente por `team-service/internal/persistence/store.go` y `frontend/src/features/new-game/hooks/useNewGameFlow.ts`. El siguiente corte recomendado es iniciar `team-service` por un agregado pequeno, comenzando por salary cap.
