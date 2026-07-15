# 2026-05-25 — M3.4 Memoria de decisiones del GM

## Objetivo

Crear el primer registro persistente de decisiones del GM para que los agentes puedan referenciar decisiones pasadas en cortes posteriores.

## Cambios realizados

- Se agrego `gm_decisions_log` como tabla append-only de `agent-service`.
- Se agrego la migracion `018_create_gm_decisions_log.sql`.
- Se documento `decision.gm_registrada` en `shared/events/decisions.md`.
- `gateway` publica `decision.gm_registrada` cuando el GM responde el evento inicial del Owner.
- `agent-service` escucha `decision.gm_registrada` y persiste la decision.
- La persistencia es idempotente por `event_id` y `(game_id, decision_id)`.
- `agent-service` agrega lectura interna `latest_gm_decisions(game_id, limit)` para contexto futuro de agentes/chat.

## Decision tecnica

No se expone el log directamente al frontend.

El log es memoria sistemica interna. M3.5 podra usarlo para construir contexto de chat, pero el jugador lo vera indirectamente por lo que los agentes recuerden o mencionen.

## Pruebas

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
make build
```

## Pendiente siguiente

Siguiente mini milestone recomendado: `M3.5 — Infra de chat con stub`.
