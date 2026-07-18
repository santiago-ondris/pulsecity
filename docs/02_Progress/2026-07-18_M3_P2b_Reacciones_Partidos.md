# M3.P2b — Reacciones a partidos

## Objetivo de la sesion

Continuar la particion de `agents.rs` con un corte de refactor puro dedicado a las consecuencias de `partido.terminado`.

## Cambios realizados

- Se creo `services/agent-service/src/agents/match_reactions.rs`.
- Se movieron al nuevo modulo:
  - reacciones de Owner, Head Coach, CFO, Scouting Director y Sports Psychologist
  - reacciones emocionales de los jugadores segun box score y resultado
  - movimientos de relaciones inter-agente causados por el partido
  - calculo de contexto local/visitante, margen, partido cerrado y blowout
  - construccion de eventos y deltas derivados del partido
- `agents.rs` reexporta las tres funciones publicas existentes, por lo que sus consumidores no cambiaron.
- Los helpers usados tambien por salary cap, medicina o trades permanecieron en el padre.

## Decision tomada

Las tres consecuencias de `partido.terminado` viven juntas porque comparten el mismo contexto y source event. Separarlas por tipo de entidad habria duplicado el calculo y fragmentado un solo dominio jugable.

No se movieron tests ni se modificaron asserts. Los tests existentes siguen validando la API publica desde `agents.rs`.

## Resultado

- `agents.rs`: 1752 -> 1206 lineas.
- `agents/match_reactions.rs`: 563 lineas.
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

`M3.P2` continua pendiente. El siguiente corte pequeno puede aislar medicina y relaciones causadas por decisiones del GM, antes de abordar trades y salary cap.
