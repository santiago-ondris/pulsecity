# M2.3 — Estado de simulacion

## Objetivo de la sesion

Persistir el estado basico de tiempo por partida en `agent-service`, sin implementar todavia el loop hibrido.

## Mini milestone

M2.3 — Estado de simulacion.

## Archivos tocados

- `db/migrations/009_create_agent_simulation_state.sql`
- `services/agent-service/Cargo.toml`
- `services/agent-service/Cargo.lock`
- `services/agent-service/src/lib.rs`
- `services/agent-service/src/main.rs`
- `services/agent-service/src/simulation.rs`
- `services/agent-service/src/persistence.rs`
- `docs/Sesiones/MILESTONE2/INICIOM2.MD`

## Decision tomada

El estado de simulacion queda bajo ownership de `agent-service` en la tabla `agent_simulation_state`.

Campos persistidos:

- `game_id`
- `current_simulated_date`
- `speed`
- `paused`
- `session_active`
- `last_tick_processed_at`

El servicio puede cargar un estado existente o inicializar uno nuevo. La regla pura `can_advance` devuelve `true` solo cuando hay sesion activa y el juego no esta pausado. El procesamiento de ticks cada 100ms queda para M2.4.

## Dependencia nueva

Se agrego `tokio-postgres` para persistencia desde Rust. Fue necesario permitir acceso a red para actualizar `Cargo.lock` y descargar crates.

## Pruebas corridas

```bash
make test-rust
make test-go
make build
```

## Pendiente siguiente

M2.4 — implementar el loop hibrido de 100ms en `agent-service`.
