# 2026-05-25 — M3.1 Catalogo completo de 30 agentes

## Objetivo

Sembrar el estado inicial real de los 30 agentes individuales del canon, dejando lista la base de M3 para relaciones, memoria y chat.

## Cambios realizados

- Se agrego el catalogo canonico de 30 agentes en `agent-service`.
- Cada agente queda definido con categoria, rol, dominio, estado emocional, confianza, satisfaccion, lealtad, rendimiento de rol, variables numericas iniciales y agenda inicial.
- Se agrego la tabla `agent_individual_states`.
- Se agrego la migracion `015_create_agent_individual_states.sql`.
- `agent-service` escucha `mapa.generacion_iniciada` y siembra los agentes de forma idempotente al fundar partida.
- El supervisor dinamico tambien asegura los agentes cuando carga o inicializa una simulacion por `game_id`.
- Los 5 core de M2 se mantienen en su flujo actual y tambien quedan integrados al nuevo schema individual.

## Decision tecnica

M3.1 es solo semilla. No se agregaron reacciones nuevas ni eventos nuevos de gameplay.

El estado vivo amplio queda en `agent_individual_states`, mientras `agent_core_states` sigue soportando el comportamiento M2 hasta que los proximos cortes migren o amplien reacciones.

## Pruebas

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
```

## Pendiente siguiente

Siguiente mini milestone recomendado: `M3.2 — Roster como agentes`.
