# 2026-05-25 — M3.3 Relaciones inter-agente

## Objetivo

Crear la red inicial de relaciones entre agentes del canon y permitir que eventos sistemicos empiecen a moverla de forma idempotente.

## Cambios realizados

- Se agrego `agent_relationships`.
- Se agrego `agent_relationship_event_hashes` para controlar idempotencia por relacion y evento fuente.
- Se agrego la migracion `017_create_agent_relationships.sql`.
- Se sembraron relaciones canonicas iniciales: Coach ↔ Analytics, Coach ↔ Medico, Coach ↔ Player Development, Scouting ↔ Analytics, Marketing ↔ GM, CFO ↔ GM, Alcalde ↔ Owner, Camara de Comercio ↔ Alcalde, Sports Psychologist ↔ Coach, PR ↔ GM y Prensa ↔ roster colectivo.
- `agent-service` asegura relaciones al fundar o cargar partida.
- `partido.terminado` mueve relaciones relevantes segun victoria/derrota, localia y margen.
- `agent-service` publica `agente.relacion_cambio`.
- `gateway` traduce ese evento a `relations.patch`.
- El frontend tipa y aplica el delta en `RelationshipClientStates`.

## Decision tecnica

M3.3 no abre todavia decisiones GM complejas. Por eso solo `partido.terminado` mueve relaciones donde el resultado tiene causalidad directa.

Las relaciones con `gm` y `roster_collective` quedan como nodos logicos, no como registros de `agent_individual_states`, porque representan al jugador y al grupo de jugadores.

## Pruebas

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
npm run build --prefix frontend
make build
```

## Pendiente siguiente

Siguiente mini milestone recomendado: `M3.4 — Memoria de decisiones del GM`.
