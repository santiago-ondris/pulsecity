# M3.12b — Aceptacion y Mutacion de Roster en Trades

Fecha: 2026-07-15

## Objetivo

Cerrar el flujo operativo de trades: una propuesta abierta o contraoferta puede aceptarse, `team-service` muta roster/cap de forma transaccional y los agentes reaccionan emocionalmente al cambio.

## Cambios realizados

- `gateway` agrega `POST /api/v1/games/{gameID}/trades/acceptances`.
- El endpoint publica `decision.gm_registrada` con `kind = trade_acceptance`.
- `team-service` consume la decision y valida:
  - propuesta existente
  - propuesta todavia abierta
  - jugador saliente activo
  - salary cap proyectado dentro del limite acotado
- `team-service` marca al jugador saliente como `traded`.
- `team-service` materializa un jugador entrante deterministico desde:
  - `proposal_id`
  - posicion requerida
  - salario entrante
  - rating del jugador saliente
- `team-service` actualiza `team_trades` con estado `accepted`.
- `team-service` publica:
  - `trade.aceptada`
  - `roster.patch`
  - `salary_cap.calculado`
  - `finance.patch`
- `agent-service` consume `trade.aceptada`.
- `agent-service` actualiza la capa emocional del jugador saliente.
- `agent-service` crea la capa emocional del jugador entrante.
- `agent-service` publica `roster.patch` emocional.
- `narrative-service` consume `trade.aceptada` y genera narrativa templateada `post_trade`.
- `gateway` transforma `trade.aceptada` en `trade.patch` con `status = accepted`.
- Se agrega `agent_processed_trades` para idempotencia de reacciones emocionales.
- Se agrega migracion `023_create_agent_processed_trades.sql`.

## Decision de arquitectura

El cierre del trade queda en `team-service` porque es el dueño de contratos, roster y salary cap. `agent-service` no decide ni escribe el roster: solo reacciona a `trade.aceptada` para actualizar estados emocionales.

## Limite

No hay UI especifica de negociacion, assets reales de draft ni impacto economico profundo de ciudad. `accepted_additional_asset` queda como string acotado para representar la concesion aceptada en M3.12.

## Verificacion

- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`
- `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`
- `GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...`
- `cargo test --manifest-path services/agent-service/Cargo.toml`
