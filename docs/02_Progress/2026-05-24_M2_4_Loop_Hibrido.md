# M2.4 — Loop hibrido

Fecha: 2026-05-24

## Objetivo

Implementar el loop de 100ms en `agent-service`, manteniendo el ownership del tiempo dentro del servicio de agentes y publicando `tiempo.dia_avanzado` cuando corresponda.

## Cambios hechos

- Se agrego `services/agent-service/src/runtime.rs`.
- `main.rs` ahora carga estado persistido y entrega el control al runtime.
- El runtime se suscribe a:
  - `tiempo.sesion_iniciada`
  - `tiempo.sesion_terminada`
  - `tiempo.velocidad_cambiada`
  - `tiempo.pausa_activada`
- Se agrego `SimulationAccumulator` para convertir tiempo real en dias simulados.
- Se agrego avance de fecha simulada con soporte para cambios de mes y años bisiestos.
- Al procesar dias:
  - se actualiza `current_simulated_date`
  - se actualiza `last_tick_processed_at`
  - se persiste `agent_simulation_state`
  - se publica `tiempo.dia_avanzado`

## Decision tomada

El runtime corre en una sola tarea con `tokio::select!`.

Esto deja el estado mutable de simulacion en un unico lugar y evita introducir locks compartidos antes de necesitarlos. Para M2.4 alcanza porque `agent-service` solo necesita reaccionar a eventos de control y al tick interno.

## Reglas implementadas

- Sin sesion activa, el tiempo no avanza.
- Con pausa activa, el tiempo no avanza.
- Al cerrar sesion o pausar, se limpia la acumulacion parcial.
- Velocidades validas: x1, x5, x20.
- Velocidades invalidas se ignoran y quedan registradas en logs.
- El evento publicado sigue el contrato `tiempo.dia_avanzado`.

## Archivos tocados

- `services/agent-service/Cargo.toml`
- `services/agent-service/src/lib.rs`
- `services/agent-service/src/main.rs`
- `services/agent-service/src/runtime.rs`
- `services/agent-service/src/simulation.rs`
- `docs/Sesiones/MILESTONE2/INICIOM2.MD`

## Verificacion

Se corrieron correctamente:

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
make test-rust
make build
```

## Pendiente siguiente

`M2.5 — HUD de tiempo`

Proximo objetivo recomendado:

- exponer controles de pausa/velocidad desde frontend hacia gateway
- traducir cambios del backend a `time.patch`
- mostrar fecha simulada y velocidad real en el HUD
