# M3.P2f-c — Narrativa y chat del gateway

## Objetivo de la sesion

Continuar la particion de `gateway/internal/handlers/http.go` aislando los endpoints del dominio narrativo.

## Cambios realizados

- Se creo `services/gateway/internal/handlers/narrative.go`.
- Se movieron desde `http.go`:
  - `answerOwnerIntro`
  - `startAgentChat`
  - `findNarrativeChoice`
- `RegisterRoutes` conserva exactamente las mismas rutas y metodos.
- Los handlers continuan compartiendo autorizacion, store y bus mediante `Dependencies` dentro del mismo paquete.

## Decision tomada

La respuesta inicial rule-based y el chat libre permanecen en el mismo archivo porque ambos son entradas HTTP del dominio narrativo. Esto no mezcla sus mecanismos internos: el Owner intro publica respuesta y memoria GM, mientras el chat solo publica `agente.consulta_iniciada` para que `narrative-service` genere la respuesta.

No se movieron medicina ni trades, aunque tambien registren decisiones del GM, porque sus propositos jugables y contratos son distintos.

## Resultado

- `http.go`: 1266 -> 1050 lineas.
- `narrative.go`: 227 lineas.
- Sin cambios en rutas HTTP, eventos NATS, payloads, ownership o comportamiento.
- Se conservaron limites y validaciones del chat.
- Suite, build y vet del gateway pasan.

## Verificacion

```bash
gofmt -w services/gateway/internal/handlers/http.go \
  services/gateway/internal/handlers/narrative.go
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...
git diff --check
```

## Pendiente siguiente

La particion del gateway continua pendiente. El siguiente corte recomendado es `M3.P2f-d`: separar decisiones medicas (`answerMedicalDecision` y `medicalDecisionLabel`) antes de trades y control de tiempo.
