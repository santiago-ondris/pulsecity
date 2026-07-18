# M3.P2a — Catalogo y semillas de agentes

## Objetivo de la sesion

Iniciar la particion de monolitos de `M3.P2` con un corte de refactor puro: separar la data-as-code de `agents.rs` sin modificar el comportamiento del `agent-service`.

## Cambios realizados

- Se creo `services/agent-service/src/agents/catalog.rs`.
- Se movieron al nuevo modulo:
  - los templates canonicos de los 30 agentes individuales
  - las semillas iniciales de relaciones inter-agente
  - el catalogo de equipos y la generacion determinista de los 30 GMs rivales
- Se conservaron en `agents.rs` las funciones publicas existentes para evitar cambios en consumidores.
- Los helpers compartidos por ambos modulos quedaron con visibilidad limitada al modulo padre.

## Decision tomada

El catalogo permanece dentro del dominio `agents`, como submodulo privado. No se creo un crate, una capa de repositorio ni una abstraccion adicional porque este corte solo necesita separar datos canonicos de logica operativa.

`M3.P2` continua pendiente: todavia deben particionarse las reacciones por dominio en `agent-service`, los handlers del gateway, el store de `team-service` y `useNewGameFlow.ts`.

## Resultado

- `agents.rs`: aproximadamente 2778 -> 1752 lineas.
- `agents/catalog.rs`: 1051 lineas de catalogo y semillas aisladas.
- Sin cambios en contratos, persistencia o eventos NATS.
- Los 26 tests del `agent-service` continuan pasando.

## Verificacion

```bash
cargo fmt --manifest-path services/agent-service/Cargo.toml -- --check
cargo test --manifest-path services/agent-service/Cargo.toml
cargo build --manifest-path services/agent-service/Cargo.toml
```

## Pendiente siguiente

Elegir un segundo corte pequeno de `M3.P2`. La opcion de menor riesgo dentro de `agent-service` es separar las reacciones a partidos de las decisiones medicas, trades y salary cap.
