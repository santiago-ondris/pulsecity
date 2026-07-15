# M3.8 — Modelo de partidos v2

Fecha: 2026-07-14

## Objetivo

Cerrar el modelo de partidos v2: que el resultado deje de depender solo del rating agregado y empiece a usar rotacion, tactica, carga reciente y estado emocional.

## Cambios

- `partido.programado` incluye `home_tactics`, `away_tactics` y `expected_minutes`.
- `team-service` deriva minutos esperados por slot de rotacion para roster propio y rivales abstractos.
- `team-service` deriva tactica del Head Coach segun rival y presion de record:
  - rivales de alto pace/ofensiva fuerzan `defensive_grind`
  - rivales vulnerables defensivamente habilitan `pace_and_space`
  - mala racha usa rotacion `top_heavy`
  - ventaja amplia de record usa rotacion `deep`
- `team-service` guarda una proyeccion propia de `roster.patch` para usar estado emocional de jugadores sin leer tablas de `agent-service`.
- `team-service` calcula fatiga desde minutos recientes guardados en `team_player_box_scores`.
- `match-service` pondera fuerza por minutos esperados, rating, fatiga, estado emocional y stamina.
- El box score reparte puntos y minutos segun rol/minutos esperados.
- `match-service` mantiene compatibilidad con payloads M2 sin tactica ni minutos explicitos.

## Decision

La emocion del jugador sigue siendo propiedad de `agent-service`.

`team-service` no lee tablas ajenas: consume `roster.patch` por NATS y mantiene una proyeccion local minima para preparar partidos. Esto respeta ownership estricto y deja el simulador stateless.

## Pruebas

- `cargo test --manifest-path services/match-service/Cargo.toml`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`

## Resultado

`M3.8 — Modelo de partidos v2` queda cerrado funcionalmente.

El siguiente corte natural es `M3.9`: usar esta carga reciente como base para riesgo de lesion, alertas del Medico y decisiones del GM sobre descanso/alta.
