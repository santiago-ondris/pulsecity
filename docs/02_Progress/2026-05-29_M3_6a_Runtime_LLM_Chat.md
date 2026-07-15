# M3.6a — Runtime LLM de chat

Fecha: 2026-05-29

## Objetivo

Preparar `narrative-service` para reemplazar el stub por un LLM real sin tocar de nuevo el flujo NATS/WebSocket.

## Cambios hechos

- Se agrego `ChatResponder` como interfaz de generacion de respuesta.
- Se agrego `ChatRuntimeConfig` configurable por env:
  - `LLM_PROVIDER`
  - `LLM_MAX_PROMPT_CHARS`
  - `LLM_MAX_RESPONSE_CHARS`
  - `LLM_MAX_TURNS_PER_CONVERSATION`
- `processAgentConsultation` ahora arma prompt canonico con contexto real antes de generar.
- El prompt incluye estado del agente, dominio, relacion con GM y decisiones recientes.
- Se agrego guardrail de dominio: el agente debe redirigir cuando la pregunta no pertenece a su area.
- Se agrego limite de turnos por conversacion, contando mensajes del GM persistidos.
- Si el provider configurado no esta implementado o falla, la respuesta cae a fallback templateado.
- `stub` sigue siendo el unico provider implementado; el siguiente corte agrega cliente real.

## Decision tomada

No se conecta proveedor real todavia. La razon es evitar mezclar tres decisiones en una sola sesion: proveedor/costo, secrets y cliente HTTP. Este corte deja la arquitectura lista para que M3.6 real sea un cambio acotado.

## Pruebas corridas

```bash
GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...
```

## Siguiente recomendado

Elegir proveedor y modelo con costo bajo para chat. Despues implementar el cliente real detras de `ChatResponder`, usando los caps ya existentes y manteniendo fallback.
