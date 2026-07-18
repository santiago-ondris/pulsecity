# M3.P2f-d — Decisiones medicas del gateway

## Objetivo de la sesion

Continuar la particion de `gateway/internal/handlers/http.go` aislando la entrada HTTP de decisiones medicas.

## Cambios realizados

- Se creo `services/gateway/internal/handlers/medical.go`.
- Se movieron desde `http.go`:
  - `answerMedicalDecision`
  - `medicalDecisionLabel`
- `RegisterRoutes` conserva exactamente la misma ruta y metodo.
- El handler continua compartiendo autorizacion, store y bus mediante `Dependencies` dentro del mismo paquete.

## Decision tomada

El archivo contiene solo la traduccion de una accion del GM a `decision.gm_registrada`. No calcula riesgo, recuperacion ni cambios de confianza: esas consecuencias siguen en `team-service` y `agent-service` segun ownership.

La tabla de etiquetas viaja con el handler porque forma parte de la validacion/representacion del request, no de la mecanica de lesiones.

## Resultado

- `http.go`: 1050 -> 944 lineas.
- `medical.go`: 116 lineas.
- Sin cambios en rutas HTTP, eventos NATS, payloads, ownership o comportamiento.
- Suite, build y vet del gateway pasan.

## Verificacion

```bash
gofmt -w services/gateway/internal/handlers/http.go \
  services/gateway/internal/handlers/medical.go
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...
git diff --check
```

## Pendiente siguiente

La particion del gateway continua pendiente. El siguiente corte recomendado es `M3.P2f-e`: separar propuestas y aceptaciones de trades antes de aislar control de tiempo y el shell HTTP restante.
