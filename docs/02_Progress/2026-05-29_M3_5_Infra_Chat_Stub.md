# M3.5 — Infra de chat con stub

Fecha: 2026-05-29

## Objetivo

Cerrar el primer camino end-to-end del chat directo GM ↔ agente sin conectar todavia un LLM real.

## Cambios hechos

- `gateway` agrega `POST /api/v1/games/{gameID}/agent-chat`.
- El endpoint valida ownership de la partida y publica `agente.consulta_iniciada`.
- `narrative-service` consume la consulta y queda como servicio responsable del chat.
- `narrative-service` lee contexto real desde tablas de `agent-service`: estado individual/jugador, relacion con GM y ultimas decisiones del GM.
- Se agrega persistencia `agent_chat_history`.
- Se genera respuesta stub deterministica por agente.
- `narrative-service` publica `chat.message` con subject `agente.respuesta_generada`.
- `gateway` reenvia `chat.message` por WebSocket.
- Frontend agrega un panel minimo de chat en `CeremonyPage` para los agentes core.
- El mensaje del GM se muestra de forma optimista; la respuesta del agente llega por delta WebSocket.
- Se documentan los contratos en `shared/events/agents.md` y `shared/events/websocket_deltas.md`.

## Decision tomada

El chat vive en `narrative-service`. `agent-service` sigue siendo dueño del estado emocional, relaciones y memoria; `narrative-service` solo lee ese contexto para generar lenguaje y persistir historial conversacional.

## Limite del corte

No se agrego REST de rehidratacion del historial. La tabla ya existe y los mensajes quedan persistidos, pero la UI actual mantiene la conversacion viva por estado local + deltas WebSocket.

## Pruebas corridas

```bash
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...
npm run build --prefix frontend
make build
```

## Siguiente recomendado

M3.6 puede conectar el proveedor LLM real sobre el mismo contrato. Antes de eso conviene decidir proveedor/costo y agregar caps de tokens/turnos.
