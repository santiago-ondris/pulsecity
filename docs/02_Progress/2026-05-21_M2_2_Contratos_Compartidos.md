# M2.2 — Contratos compartidos M2

## Objetivo de la sesion

Definir contratos canonicos para los eventos principales de Milestone 2 antes de implementar logica de simulacion.

## Mini milestone

M2.2 — Contratos compartidos M2.

## Archivos tocados

- `shared/events/time.md`
- `shared/events/matches.md`
- `shared/events/city.md`
- `shared/events/agents.md`
- `shared/events/narrative.md`
- `shared/events/websocket_deltas.md`
- `services/agent-service/src/events.rs`
- `services/match-service/src/events.rs`
- `services/team-service/internal/domain/events.go`
- `services/city-service/internal/domain/events.go`
- `services/analytics-service/internal/domain/events.go`
- `docs/Sesiones/MILESTONE2/INICIOM2.MD`

## Decision tomada

Los contratos quedan documentados en `shared/events` como fuente canonica humana. Por ahora no se crea un paquete compartido importable multi-lenguaje: cada servicio mantiene sus tipos internos alineados con esos documentos. Esto evita agregar tooling de generacion Go/Rust antes de tener el primer flujo sistemico funcionando.

Todos los eventos NATS de M2 que representan hechos incluyen `event_id`, `game_id`, `occurred_at` y `schema_version`. `event_id` y los IDs de dominio como `match_id` quedan disponibles para idempotencia.

## Pruebas corridas

```bash
make test-rust
make test-go
make build
```

## Pendiente siguiente

M2.3 — persistir estado basico de simulacion en `agent-service`.
