# SESION21MAYO

## Milestone 2 — inicio operativo

Arranque practico de `Milestone 2`: pasar del mundo fundacional de M1 a la base tecnica del sistema que va a latir solo.

El foco de la sesion fue preparar los servicios, contratos y estado de simulacion necesarios antes de implementar el loop hibrido.

## Mini milestones trabajados

- `M2.1 — Scaffolding de servicios M2`
- `M2.2 — Contratos compartidos M2`
- `M2.3 — Estado de simulacion`

Los tres quedaron marcados como `YA REALIZADO` en `INICIOM2.MD`.

## Objetivo general de la sesion

Dejar listo el piso tecnico para que `agent-service` pueda convertirse en dueño real del tiempo simulado en el siguiente corte.

Esto implico:

- crear los servicios nuevos de M2
- formalizar los contratos de eventos que van a conectar el sistema
- persistir el estado basico de simulacion por partida
- mantener el scope chico, sin implementar todavia gameplay ni loop de 100ms

## M2.1 — Scaffolding de servicios

Se crearon los servicios nuevos definidos para M2:

- `services/agent-service` en Rust
- `services/match-service` en Rust
- `services/team-service` en Go
- `services/city-service` en Go
- `services/analytics-service` en Go

Cada servicio quedo con estructura minima real:

- entrypoint
- conexion a NATS por `NATS_URL`
- estructura interna base
- tests simples
- integracion en `Makefile`

En Go se mantuvo la estructura canonica:

```text
cmd/
internal/domain/
internal/handlers/
internal/nats/
```

En Rust se dejo:

```text
src/main.rs
src/lib.rs
```

## M2.2 — Contratos compartidos

Se documentaron contratos canonicos en `shared/events`.

Archivos agregados:

- `shared/events/time.md`
- `shared/events/matches.md`
- `shared/events/city.md`
- `shared/events/agents.md`
- `shared/events/narrative.md`
- `shared/events/websocket_deltas.md`

Tambien se agregaron tipos internos alineados en:

- `agent-service`
- `match-service`
- `team-service`
- `city-service`
- `analytics-service`

## Decision sobre contratos compartidos

Por ahora `shared/events` queda como fuente canonica humana.

No se creo todavia un paquete compartido importable Go/Rust ni tooling de generacion de tipos. La razon es mantener bajo el costo operativo mientras todavia estamos construyendo el primer flujo sistemico.

Cada servicio tiene sus propios tipos internos alineados con los documentos de `shared/events`.

## Regla comun de eventos M2

Los eventos NATS de M2 que representan hechos del sistema incluyen:

- `event_id`
- `game_id`
- `occurred_at`
- `schema_version`

Esto queda pensado desde el principio para:

- idempotencia
- debugging
- trazabilidad
- separar hechos de deltas visuales

## M2.3 — Estado de simulacion

Se implemento el primer estado real de simulacion en `agent-service`.

Campos persistidos:

- `game_id`
- `current_simulated_date`
- `speed`
- `paused`
- `session_active`
- `last_tick_processed_at`

Se agrego la migracion:

- `db/migrations/009_create_agent_simulation_state.sql`

Y los modulos:

- `services/agent-service/src/simulation.rs`
- `services/agent-service/src/persistence.rs`

## Estado tecnico alcanzado

`agent-service` ahora puede:

- conectar a NATS
- conectar a PostgreSQL
- asegurar el schema propio de simulacion
- cargar estado existente por `game_id`
- inicializar estado si no existe
- guardar estado
- evaluar si el tiempo puede avanzar con `can_advance`

La regla actual es:

```text
puede avanzar = session_active && !paused
```

## Decision sobre el loop

El loop de 100ms no se implemento en esta sesion.

Queda deliberadamente para `M2.4`, porque este corte buscaba cerrar solamente persistencia y reglas base de estado. Separar esto evita mezclar:

- estructura de estado
- persistencia
- consumo de eventos de control
- acumulacion temporal
- publicacion de `tiempo.dia_avanzado`

## Dependencia nueva

Se agrego `tokio-postgres` a `agent-service` para persistencia desde Rust.

Fue necesario permitir acceso a red para:

- actualizar `Cargo.lock`
- descargar crates nuevas
- correr `make test-rust` con la dependencia nueva

## Verificacion

Se corrieron correctamente:

```bash
make test-rust
make test-go
make build
```

## Documentacion relacionada

Ademas de esta nota de sesion, quedaron notas granulares en:

- `docs/02_Progress/2026-05-21_M2_1_Scaffolding_Servicios.md`
- `docs/02_Progress/2026-05-21_M2_2_Contratos_Compartidos.md`
- `docs/02_Progress/2026-05-21_M2_3_Estado_Simulacion.md`

## Estado real al cierre

M2 ya tiene base estructural para avanzar:

- servicios nuevos creados
- contratos principales definidos
- `agent-service` con estado persistente de simulacion
- `Makefile` cubriendo build/test/run de los servicios M2
- `INICIOM2.MD` actualizado hasta `M2.3`

## Pendiente siguiente

`M2.4 — Loop hibrido`

Proximo objetivo recomendado:

- hacer que `agent-service` escuche eventos de sesion, pausa y velocidad
- actualizar `agent_simulation_state`
- correr un tick cada 100ms
- publicar `tiempo.dia_avanzado` cuando corresponda
- mantener la regla: sin sesion activa, el tiempo no avanza
