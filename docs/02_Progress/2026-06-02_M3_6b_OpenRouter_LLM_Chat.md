# 2026-06-02 — M3.6b OpenRouter LLM Chat

## Objetivo

Conectar un provider LLM real al chat directo GM ↔ agente sin romper el modo local con stub.

## Cambios

- Se eligio `OpenRouter` como provider de exploracion para M3.6.
- Se agrego provider `openrouter` detras de `ChatResponder` en `narrative-service`.
- El modelo default queda en `google/gemini-2.5-flash`, configurable por `LLM_MODEL`.
- El servicio llama `POST /api/v1/chat/completions` con request OpenAI-compatible.
- Se agregaron env vars:
  - `OPENROUTER_API_KEY`
  - `OPENROUTER_BASE_URL`
  - `OPENROUTER_APP_URL`
  - `OPENROUTER_APP_TITLE`
  - `LLM_MODEL`
  - `LLM_MAX_COMPLETION_TOKENS`
  - `LLM_REQUEST_TIMEOUT_SECONDS`
- `LLM_PROVIDER=stub` sigue siendo el default para desarrollo local sin key.
- Si falta key, falla OpenRouter o la respuesta viene vacia, el flujo existente usa fallback templateado.
- La metadata de respuesta registra provider, modelo y usage de tokens cuando OpenRouter lo informa.
- El `Makefile` carga `.env.local` automaticamente si existe y exporta sus variables a `make run-*` y `make dev-app`.

## Verificacion

- `GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...`
- `make -n run-narrative-service`
- `make -n dev-app`
- Smoke real desde la app local: chat con CFO respondio con texto natural via provider real, sin `[stub M3.6]` y sin fallback.

## Resultado

`M3.6 — LLM real conectado` queda cerrado funcionalmente. El siguiente mini milestone recomendado es `M3.7 — Buscador de agentes y vista de directorio`.
