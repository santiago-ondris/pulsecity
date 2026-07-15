# M2.7 — Match-service v1

Fecha: 2026-05-24

## Objetivo

Implementar un simulador de partidos deterministico en `match-service`, manteniendo el servicio stateless y listo para ser conectado por `team-service` en `M2.8`.

## Cambios hechos

- Se agrego `services/match-service/src/simulator.rs`.
- Se agrego `services/match-service/src/runtime.rs`.
- `match-service` ahora escucha `partido.programado`.
- Al recibir un partido valido:
  - publica `partido.iniciando`
  - simula el partido
  - publica `partido.terminado`
- El simulador genera:
  - marcador local/visitante
  - ganador
  - box score por jugador de rotacion
  - 4 momentos clave
- El contrato `partido.programado` ahora incluye `players`, porque `match-service` no consulta otros servicios.

## Decision tomada

El modelo es intencionalmente de equipo/rotacion, no posesion a posesion.

Para M2 necesitamos reproducibilidad, box score plausible y momentos narrativos suficientes para activar ciudad/agentes/narrativa. El simulador profundo queda para una evolucion posterior, cuando el loop completo ya este funcionando.

## Reglas implementadas

- Mismo input + misma seed = mismo output.
- Seed distinta cambia el resultado.
- No se permiten empates.
- El box score suma exactamente el puntaje de cada equipo.
- Si un equipo tiene menos de 5 jugadores, la simulacion falla.
- `match-service` no persiste estado y no consulta otros servicios.

## Archivos tocados

- `services/match-service/Cargo.toml`
- `services/match-service/src/events.rs`
- `services/match-service/src/lib.rs`
- `services/match-service/src/main.rs`
- `services/match-service/src/runtime.rs`
- `services/match-service/src/simulator.rs`
- `services/team-service/internal/domain/events.go`
- `shared/events/matches.md`
- `docs/Sesiones/MILESTONE2/INICIOM2.MD`

## Verificacion

Se corrieron correctamente:

```bash
cargo test --manifest-path services/match-service/Cargo.toml
make test-rust
make test-go
make build
```

## Pendiente siguiente

`M2.8 — Orquestacion de partido`

Proximo objetivo recomendado:

- `team-service` escucha `tiempo.dia_avanzado`
- detecta si hay partido programado para la fecha
- arma el payload completo con roster propio y rival abstracto
- publica `partido.programado`
- evita duplicar simulaciones si el evento se reintenta
