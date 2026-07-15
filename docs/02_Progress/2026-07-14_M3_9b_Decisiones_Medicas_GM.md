# M3.9b — Decisiones medicas del GM

Fecha: 2026-07-14

## Objetivo

Cerrar el flujo medico de M3: una lesion no solo baja a un jugador, tambien habilita una decision del GM y consecuencias sobre salud, disponibilidad y relaciones.

## Cambios

- `gateway` agrega endpoint `POST /api/v1/games/{gameID}/medical-decisions`.
- El endpoint valida ownership de partida y publica `decision.gm_registrada` con `kind = medical_decision`.
- Opciones validas:
  - `rest`
  - `reduce_minutes`
  - `ignore_doctor`
  - `force_return`
- `team-service` consume `decision.gm_registrada`.
- `team-service` registra la decision medica en `team_injuries`.
- Si la decision es `force_return`:
  - marca `forced_return_at`
  - activa al jugador inmediatamente
  - publica `jugador.recuperado`
- Si un jugador con alta forzada juega antes de `expected_recovery_date`, `team-service` genera una nueva lesion con `reason = forced_return_reaggravation`.
- `agent-service` mueve relaciones canonicas ante `medical_decision`:
  - GM ↔ Medico
  - Head Coach ↔ Medico
- `agent-service` publica `agente.relacion_cambio`, y el gateway ya lo convierte en `relations.patch`.

## Contratos

- `shared/events/decisions.md`
- `shared/events/players.md`

## Decision

No se agrego UI nueva en este corte. El endpoint queda listo para conectarse desde el inbox/panel roster en M3.22/M3.23.

## Pruebas

- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`
- `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`
- `cargo test --manifest-path services/agent-service/Cargo.toml`

## Resultado

`M3.9 — Carga, lesiones y flujo medico` queda cerrado funcionalmente.
