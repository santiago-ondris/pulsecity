# M3.P2e — Reacciones a salary cap

## Objetivo de la sesion

Cerrar la particion por dominios de `agent-service` aislando las reacciones de agentes al estado del salary cap.

## Cambios realizados

- Se creo `services/agent-service/src/agents/salary_cap_reactions.rs`.
- Se movieron al nuevo modulo:
  - cambios de `budget_alert` y `financial_trust` del CFO
  - cambios de `patience_remaining` y `business_trust` del Owner
  - seleccion de mood segun estado financiero
  - construccion de `agente.estado_cambio` con source `salary_cap.calculado`
- `agents.rs` reexporta `apply_salary_cap_to_core_agents`, manteniendo sin cambios a persistence y runtime.
- Los imports de eventos usados solo por tests quedaron explicitos dentro del modulo de tests.

## Decision tomada

`team-service` mantiene ownership del calculo de salary cap. El modulo nuevo solo contiene la reaccion emocional de datos propiedad de `agent-service`.

No se separaron los tests de `agents.rs`: despues de este corte, la logica de produccion termina en la linea 326 y las aproximadamente 535 lineas restantes son cobertura existente. Mover tests solo para bajar el contador no mejora la separacion de dominios.

## Resultado

- `agents.rs`: 933 -> 862 lineas fisicas; 326 corresponden a produccion.
- `agents/salary_cap_reactions.rs`: 74 lineas.
- La particion por dominio de `agent-service` dentro de `M3.P2` queda cerrada.
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

`M3.P2` continua pendiente por los monolitos restantes: handlers del gateway, store de `team-service` y `useNewGameFlow.ts` en frontend.
