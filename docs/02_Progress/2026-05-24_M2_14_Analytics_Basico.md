# 2026-05-24 — M2.14 Analytics basico

## Objetivo

Activar `analytics-service` como consumidor de series sin UI visible todavia.

## Implementado

- `analytics-service` conecta a Postgres/TimescaleDB usando `pgx`.
- Se agrego persistencia propia en `internal/persistence`.
- El servicio asegura schema al arrancar.
- Consume eventos NATS:
  - `partido.terminado`
  - `ciudad.economia_cambio`
  - `ciudad.suelo_actualizado`
  - `agente.estado_cambio`
- Persiste:
  - resultados de partidos
  - box scores por jugador
  - ticket sales
  - fan sentiment
  - economia local
  - valor de suelo por zona
  - estado historico de agentes
- Cada handler corre en goroutine con timeout corto para no bloquear el dispatcher de NATS.
- La idempotencia se controla con claves primarias por evento/partido/metric.

## Migracion

- `014_create_analytics_timeseries_tables.sql`

## Decision

Las tablas se crean como PostgreSQL normal para no exigir que la extension TimescaleDB exista en todos los entornos locales. Quedan modeladas como series temporales y listas para convertirse a hypertables cuando activemos formalmente TimescaleDB en infra.

No se agrego UI ni endpoints de consulta. Eso queda para charts futuros.

## Verificacion

```bash
GOCACHE=/tmp/pulsecity-analytics-gocache go test ./...
```
