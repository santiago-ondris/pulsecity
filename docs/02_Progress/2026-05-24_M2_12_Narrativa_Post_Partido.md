# 2026-05-24 — M2.12 Narrativa post-partido

## Objetivo

Generar eventos narrativos basicos despues de partidos, sin LLM real todavia, conectados al resultado del sistema.

## Implementado

- `narrative-service` escucha `partido.terminado`.
- El handler espera entre 250ms y 500ms antes de generar narrativa.
- La espera se ejecuta en goroutine para no bloquear el dispatcher de NATS.
- Se genera `narrativa.evento_generado` con:
  - emisor
  - urgencia
  - titulo
  - cuerpo
  - metadata del partido
  - opcion simple `acknowledge`
- El texto usa:
  - victoria/derrota
  - localia
  - margen
  - maximo anotador propio
  - primer momento clave disponible
  - racha positiva/negativa si ya esta en `city_metrics`
- Se persiste en `narrative_events`.
- La idempotencia se controla con indice unico por `(game_id, source_match_id, kind)`.
- El gateway ya reenvia `narrativa.evento_generado` al frontend como `narrative.event`.

## Migracion

- `013_create_narrative_events_tables.sql`

## Decision

La narrativa sigue siendo templateada/rule-based. Esto mantiene el foco en probar que el contexto sistemico sea real antes de conectar un LLM.

`narrative-service` lee contexto de racha desde `city_metrics` en modo read-only despues del delay. No escribe estado de ciudad.

## Verificacion

```bash
GOCACHE=/tmp/pulsecity-narrative-gocache go test ./...
```
