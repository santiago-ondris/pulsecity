# M3.P2c — Decisiones medicas y relaciones

## Objetivo de la sesion

Continuar la particion de `agents.rs` aislando las consecuencias relacionales de las decisiones medicas del GM.

## Cambios realizados

- Se creo `services/agent-service/src/agents/medical_reactions.rs`.
- Se movieron al nuevo modulo:
  - el consumo logico de decisiones GM con `kind = medical_decision`
  - la matriz de deltas para `rest`, `reduce_minutes`, `ignore_doctor` y `force_return`
  - la construccion de `agente.relacion_cambio` derivada de la decision
- `agents.rs` reexporta `apply_gm_decision_to_relationships`, manteniendo sin cambios a sus consumidores.
- El helper que aplica confianza, tendencia e historial permanece en el padre y se comparte con las reacciones a partidos.

## Decision tomada

Este modulo cubre solo la reaccion emocional/relacional dentro de `agent-service`. El calculo de lesiones, las recomendaciones del Medico y las mutaciones contractuales siguen en sus servicios y modulos dueños.

No se introdujo un tipo nuevo para las decisiones ni se modificaron strings de `choice_id`, porque este corte es exclusivamente estructural.

## Resultado

- `agents.rs`: 1206 -> 1112 lineas.
- `agents/medical_reactions.rs`: 106 lineas.
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

`M3.P2` continua pendiente. Dentro de `agent-service`, los siguientes dominios separables son trades y salary cap; despues corresponde abordar los monolitos de gateway, team-service y frontend.
