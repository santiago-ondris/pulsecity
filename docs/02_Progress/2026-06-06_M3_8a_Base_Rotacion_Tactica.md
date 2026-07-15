# M3.8a — Base de rotacion y tactica

Fecha: 2026-06-06

## Objetivo

Abrir el modelo de partidos v2 sin reescribir todo el simulador: la rotacion y los minutos esperados empiezan a pesar en el resultado y en el box score.

## Cambios

- `partido.programado` acepta `home_tactics`, `away_tactics` y `expected_minutes` por jugador.
- `team-service` publica contexto tactico inicial y minutos esperados deterministas para roster propio y rivales abstractos.
- `match-service` mantiene compatibilidad con payloads M2 sin campos nuevos.
- `match-service` pondera fuerza de equipo por minutos esperados, rating, fatiga, estado emocional y stamina.
- El box score reparte produccion segun minutos esperados, no solo por top de rating.
- Los sistemas tacticos (`balanced`, `pace_and_space`, `defensive_grind`) ajustan pace y pequenos bonuses.

## Decision

El corte queda como `M3.8a`, no como cierre completo de `M3.8`.

Todavia falta conectar una filosofia real del Head Coach, estados emocionales/fatiga vivos desde los servicios dueños y decisiones reales de rotacion desde el GM/Coach.

## Pruebas

- `cargo test --manifest-path services/match-service/Cargo.toml`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`

## Pendiente siguiente

Continuar con `M3.8b`: conectar el input del modelo v2 con decisiones reales de Coach/GM y estado emocional/fatiga vivo.
