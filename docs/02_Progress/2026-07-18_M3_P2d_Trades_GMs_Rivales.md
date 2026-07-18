# M3.P2d — Trades y GMs rivales

## Objetivo de la sesion

Continuar la particion de `agents.rs` aislando la negociacion de GMs rivales y las consecuencias emocionales de un trade aceptado.

## Cambios realizados

- Se creo `services/agent-service/src/agents/trade_reactions.rs`.
- Se movieron al nuevo modulo:
  - `RivalGMTradeEvaluation`
  - evaluacion de encaje con necesidades, diferencia salarial, urgencia, confianza y estilo de negociacion
  - construccion de `trade.rechazada` y `trade.contraoferta`
  - reaccion emocional del jugador saliente y del entrante
  - construccion del `roster.patch` derivado de `trade.aceptada`
- `agents.rs` reexporta el enum y las funciones publicas, manteniendo sin cambios a persistence y runtime.
- Los imports usados solo por tests se hicieron explicitos dentro del modulo de tests.

## Decision tomada

La evaluacion del GM rival y la reaccion emocional del roster comparten el dominio jugable trade, pero no escriben estado contractual. `team-service` conserva ownership de roster, salarios y contratos; `agent-service` conserva ownership de perfiles rivales y estado emocional.

El catalogo de los 30 GMs rivales permanece en `catalog.rs` porque es data-as-code, no logica de negociacion.

## Resultado

- `agents.rs`: 1112 -> 933 lineas.
- `agents/trade_reactions.rs`: 193 lineas.
- Sin cambios en eventos NATS, payloads, persistencia o comportamiento.
- Los 26 tests del `agent-service` pasan sin warnings.

## Verificacion

```bash
cargo fmt --manifest-path services/agent-service/Cargo.toml -- --check
cargo test --manifest-path services/agent-service/Cargo.toml
cargo build --manifest-path services/agent-service/Cargo.toml
git diff --check
```

## Pendiente siguiente

`M3.P2` continua pendiente. En `agent-service` resta separar salary cap y decidir si los defaults/tipos compartidos justifican otro corte antes de pasar a gateway, team-service y frontend.
