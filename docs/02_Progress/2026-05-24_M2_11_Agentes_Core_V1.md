# 2026-05-24 — M2.11 Agentes core v1

## Objetivo

Crear estado real para los 5 agentes core de `Milestone 2` y conectarlos al loop sistemico ya existente.

## Implementado

- `agent-service` ahora tiene modulo de dominio `agents`.
- Se definieron estados iniciales para:
  - Owner
  - Head Coach
  - CFO
  - Director de Scouting
  - Sports Psychologist
- `agent-service` escucha `partido.terminado`.
- Cada partido actualiza variables emocionales y operativas por agente.
- Cada agente publica `agente.estado_cambio`.
- Se agrego idempotencia por `(game_id, match_id)` con `agent_processed_matches`.
- Se persiste estado en `agent_core_states`.
- `gateway` traduce `agente.estado_cambio` a delta WebSocket `agent.patch`.
- El frontend mantiene `AgentClientStates` y muestra resumen compacto de agentes core en `CeremonyPage`.

## Variables base

- Owner: confianza deportiva, confianza de negocio, paciencia restante, satisfaccion.
- Head Coach: confianza en GM, satisfaccion con roster, presion por resultados, relacion con vestuario.
- CFO: confianza financiera, alerta de presupuesto, conservadurismo financiero.
- Director de Scouting: confianza en criterio, motivacion, precision percibida.
- Sports Psychologist: clima de vestuario, alerta emocional, confianza de jugadores.

## Decisiones

- Los valores se mantienen en rango `-1.0` a `1.0`.
- El frontend recibe `state` como mapa anidado dentro de `agent.patch`, porque cada agente tiene variables distintas.
- No se implementa todavia matriz de relaciones ni causalidad profunda; eso queda para M3.
- No se implementa narrativa post-partido en este mini milestone; queda como siguiente bloque natural (`M2.12`).

## Verificacion

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
GOCACHE=/tmp/pulsecity-gateway-gocache go test ./...
npm run build --prefix frontend
```
