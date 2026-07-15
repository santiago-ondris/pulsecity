# M3.12a — Propuesta y Evaluacion Rival de Trades

Fecha: 2026-07-15

## Objetivo

Abrir el primer tramo operativo de trades sin mutar roster todavia: el GM puede iniciar una propuesta, `team-service` valida roster/cap y `agent-service` evalua la respuesta del GM rival usando perfiles persistidos en `M3.11`.

## Cambios realizados

- `gateway` agrega `POST /api/v1/games/{gameID}/trades/proposals`.
- El endpoint publica `decision.gm_registrada` con `kind = trade_proposal`.
- `team-service` consume la decision y valida:
  - jugador ofrecido existe en el roster de PulseCity
  - jugador ofrecido esta activo
  - salario entrante no empuja al equipo por encima de luxury tax line
- Se agrega tabla `team_trades`.
- Se agrega migracion `021_create_team_trades.sql`.
- Si la validacion local falla, `team-service` publica `trade.rechazada`.
- Si la validacion local pasa, `team-service` publica `trade.propuesta_enviada`.
- `agent-service` consume `trade.propuesta_enviada`.
- `agent-service` evalua la propuesta con `rival_gms` usando:
  - necesidades de roster
  - estilo de negociacion
  - urgencia actual
  - confianza inicial con el GM jugador
- Se agrega tabla `rival_gm_trade_evaluations` para idempotencia.
- Se agrega migracion `022_create_rival_gm_trade_evaluations.sql`.
- `agent-service` publica `trade.rechazada` o `trade.contraoferta`.
- `gateway` escucha `trade.*` y emite `trade.patch` al frontend.

## Decision de arquitectura

`team-service` no lee `rival_gms` porque esa tabla pertenece a `agent-service`. La negociacion queda coreografiada por NATS:

```text
gateway -> decision.gm_registrada
team-service -> trade.propuesta_enviada | trade.rechazada
agent-service -> trade.contraoferta | trade.rechazada
gateway -> trade.patch
```

Esto mantiene ownership estricto: `team-service` valida roster/cap, `agent-service` evalua personalidad del rival.

## Limite

Este corte no implementa `trade.aceptada`, aceptacion de contraoferta, materializacion del jugador recibido ni mutacion de roster. Eso queda como `M3.12b`.

## Verificacion

- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`
- `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`
- `cargo test --manifest-path services/agent-service/Cargo.toml`
