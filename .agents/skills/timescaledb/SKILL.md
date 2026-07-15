---
name: timescaledb
description: >
  TimescaleDB (PostgreSQL con extensión de series temporales) para PulseCity.
  Úsala cuando trabajes en analytics-service o cualquier código que lea/escriba
  en la base de datos. Cubre: hypertables, continuous aggregates, retención,
  compresión, y queries de series temporales para stats del juego.
  También cubre uso con sqlx y pgx en Go, y la extensión diesel-async en Rust.
---

# TimescaleDB para PulseCity

PulseCity usa TimescaleDB como única base de datos — PostgreSQL familiar
con superpoderes para el volumen de series temporales que genera la simulación.

## Setup inicial

```sql
-- Habilitar la extensión (una vez por base de datos)
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Verificar versión
SELECT extversion FROM pg_extension WHERE extname = 'timescaledb';
```

## Tablas principales de PulseCity

### Eventos de simulación (serie temporal principal)

```sql
CREATE TABLE sim_events (
    time        TIMESTAMPTZ NOT NULL,  -- tiempo real del insert
    sim_day     DATE        NOT NULL,  -- día simulado
    event_type  TEXT        NOT NULL,  -- "jugador.firmado", "partido.terminado", etc.
    partida_id  UUID        NOT NULL,
    payload     JSONB       NOT NULL
);

-- Convertir en hypertable particionada por tiempo
SELECT create_hypertable('sim_events', by_range('time'));

-- Índice para queries por partida y tipo
CREATE INDEX ON sim_events (partida_id, event_type, time DESC);
```

### Stats de jugadores por partido

```sql
CREATE TABLE player_game_stats (
    time        TIMESTAMPTZ NOT NULL,
    sim_day     DATE        NOT NULL,
    partida_id  UUID        NOT NULL,
    player_id   UUID        NOT NULL,
    partido_id  UUID        NOT NULL,
    puntos      INT,
    rebotes     INT,
    asistencias INT,
    minutos     FLOAT,
    plus_minus  INT,
    rating_rendimiento FLOAT  -- calculado por match-service
);

SELECT create_hypertable('player_game_stats', by_range('time'));
CREATE INDEX ON player_game_stats (partida_id, player_id, time DESC);
```

### Estado de agentes en el tiempo

```sql
CREATE TABLE agent_state_history (
    time        TIMESTAMPTZ NOT NULL,
    sim_day     DATE        NOT NULL,
    partida_id  UUID        NOT NULL,
    agent_id    UUID        NOT NULL,
    confianza   FLOAT,
    satisfaccion FLOAT,
    lealtad     FLOAT
);

SELECT create_hypertable('agent_state_history', by_range('time'));
CREATE INDEX ON agent_state_history (partida_id, agent_id, time DESC);
```

### Economía de la ciudad

```sql
CREATE TABLE city_economy (
    time            TIMESTAMPTZ NOT NULL,
    sim_day         DATE        NOT NULL,
    partida_id      UUID        NOT NULL,
    zona_id         TEXT,          -- NULL = indicador global
    indicador       TEXT        NOT NULL,  -- "valor_suelo", "ticket_sales", "empleo"
    valor           FLOAT       NOT NULL,
    valor_anterior  FLOAT
);

SELECT create_hypertable('city_economy', by_range('time'));
```

## Continuous Aggregates — pre-calcular stats frecuentes

```sql
-- Promedio móvil de rendimiento del equipo (últimos 10 partidos)
CREATE MATERIALIZED VIEW team_rolling_avg
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 day', time) AS bucket,
    partida_id,
    AVG(puntos) OVER (
        PARTITION BY partida_id
        ORDER BY time_bucket('1 day', time)
        ROWS BETWEEN 9 PRECEDING AND CURRENT ROW
    ) AS avg_puntos_10g,
    AVG(rating_rendimiento) OVER (
        PARTITION BY partida_id
        ORDER BY time_bucket('1 day', time)
        ROWS BETWEEN 9 PRECEDING AND CURRENT ROW
    ) AS avg_rating_10g
FROM player_game_stats
GROUP BY time_bucket('1 day', time), partida_id;

-- Refresh automático
SELECT add_continuous_aggregate_policy('team_rolling_avg',
    start_offset => INTERVAL '30 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour'
);
```

## Retención y compresión

```sql
-- Comprimir datos más viejos de 30 días simulados
ALTER TABLE player_game_stats SET (
    timescaledb.compress,
    timescaledb.compress_orderby = 'time DESC',
    timescaledb.compress_segmentby = 'partida_id, player_id'
);

SELECT add_compression_policy('player_game_stats', INTERVAL '30 days');

-- Retención: borrar datos de partidas completadas después de 1 año real
SELECT add_retention_policy('sim_events', INTERVAL '365 days');
```

## Queries típicas para analytics-service

```sql
-- Stats de un jugador en los últimos N partidos
SELECT
    sim_day,
    partido_id,
    puntos,
    rebotes,
    asistencias,
    minutos,
    plus_minus
FROM player_game_stats
WHERE partida_id = $1
  AND player_id = $2
  AND time > NOW() - INTERVAL '90 days'
ORDER BY time DESC
LIMIT 20;

-- Evolución del estado emocional de un agente
SELECT
    time_bucket('7 days', time) AS semana,
    AVG(confianza) AS confianza_avg,
    AVG(satisfaccion) AS satisfaccion_avg,
    AVG(lealtad) AS lealtad_avg
FROM agent_state_history
WHERE partida_id = $1
  AND agent_id = $2
GROUP BY semana
ORDER BY semana DESC;

-- Top jugadores por rendimiento en la temporada actual
SELECT
    player_id,
    AVG(puntos) AS ppg,
    AVG(rebotes) AS rpg,
    AVG(asistencias) AS apg,
    COUNT(*) AS partidos_jugados
FROM player_game_stats
WHERE partida_id = $1
  AND sim_day BETWEEN $2 AND $3
GROUP BY player_id
ORDER BY ppg DESC;

-- Valor del suelo por zona en el tiempo
SELECT
    sim_day,
    zona_id,
    valor
FROM city_economy
WHERE partida_id = $1
  AND indicador = 'valor_suelo'
  AND zona_id = ANY($2)  -- array de zonas
ORDER BY sim_day;
```

## Go — sqlx con TimescaleDB

```go
import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func NewDB(dsn string) (*sqlx.DB, error) {
    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("connect timescaledb: %w", err)
    }
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    return db, nil
}

// Insertar stats de partido (batch para performance)
func InsertPlayerGameStats(ctx context.Context, db *sqlx.DB, stats []PlayerGameStat) error {
    query := `
        INSERT INTO player_game_stats
            (time, sim_day, partida_id, player_id, partido_id, puntos, rebotes, asistencias, minutos, plus_minus, rating_rendimiento)
        VALUES
            (:time, :sim_day, :partida_id, :player_id, :partido_id, :puntos, :rebotes, :asistencias, :minutos, :plus_minus, :rating_rendimiento)
    `
    _, err := db.NamedExecContext(ctx, query, stats)
    return err
}

// Query con time_bucket
type TeamAvgRow struct {
    Bucket   time.Time `db:"bucket"`
    AvgPuntos float64  `db:"avg_puntos"`
}

func GetTeamAvgByWeek(ctx context.Context, db *sqlx.DB, partidaID string, desde time.Time) ([]TeamAvgRow, error) {
    var rows []TeamAvgRow
    err := db.SelectContext(ctx, &rows, `
        SELECT
            time_bucket('7 days', time) AS bucket,
            AVG(puntos) AS avg_puntos
        FROM player_game_stats
        WHERE partida_id = $1 AND time >= $2
        GROUP BY bucket
        ORDER BY bucket DESC
    `, partidaID, desde)
    return rows, err
}
```

## Reglas para PulseCity específicamente

- Siempre incluir `partida_id` en todas las queries — cada partida es aislada.
- Usar `sim_day` (DATE) para queries lógicas del juego; `time` (TIMESTAMPTZ) para retention y compresión.
- El `analytics-service` es el ÚNICO que escribe en TimescaleDB. Los demás servicios
  publican eventos en NATS y analytics-service los persiste.
- Para lecturas del frontend, el gateway puede hacer REST al analytics-service directamente
  (excepción al patrón NATS-first, igual que las lecturas de team-service y agent-service).
- Nunca usar `SELECT *` en queries de producción. Siempre columnas explícitas.
- Para inserts de alta frecuencia (cada día simulado en x20 = 12.5 inserts/segundo),
  usar batch inserts — no insert individual por evento.
