# 2026-05-24 — M2.15 Temporada completa

## Objetivo

Cerrar `Milestone 2` validando que una temporada regular completa pueda correr con la cadena sistemica activa:

```text
tiempo -> calendario -> partido -> resultado -> record -> ciudad -> agentes -> narrativa -> analytics
```

## Validacion smoke end-to-end

Se levanto el entorno local con:

- NATS
- TimescaleDB/PostgreSQL
- gateway
- map-service
- team-service
- match-service
- city-service
- agent-service
- narrative-service
- analytics-service

Se creo una partida guest real:

- `game_id`: `0ba7bc3d-f6ca-49db-9597-5196485b0e65`
- calendario generado: 82 partidos
- velocidad usada: x20

Resultado en DB:

```text
team_schedule finalizados:      82 / 82
team_records:                   37-45
city_processed_matches:         82
agent_processed_matches:        82
narrative_events post_match:    82
analytics_match_results:        82
analytics_player_box_scores:    1640
analytics_city_metric_points:   246
analytics_land_value_points:    82
analytics_agent_state_points:   1394
```

## Bug encontrado y corregido

Durante el primer smoke, a velocidad x20 el `agent-service` podia publicar `tiempo.dia_avanzado` con `days_processed > 1`.

`team-service` solo buscaba partidos en `simulated_date`, la fecha final del evento, y se salteaba partidos en dias intermedios.

Ejemplo:

```text
days_processed = 3
simulated_date = 2026-10-24
```

Antes solo se procesaba `2026-10-24`. Ahora `team-service` expande el rango:

```text
2026-10-22
2026-10-23
2026-10-24
```

Con eso no se pierden partidos al correr en x20.

## Cambios de codigo

- `team-service` agrega `CoveredDayAdvancedEvents`.
- `team-service` procesa cada fecha cubierta por `tiempo.dia_avanzado`.
- Se agregaron tests para expansion de fechas y rechazo de fecha invalida.

## Verificacion

```bash
GOCACHE=/tmp/pulsecity-team-gocache go test ./...
make test-go
make test-rust
npm run build --prefix frontend
make build
```

## Estado

`Milestone 2` queda cerrado funcionalmente: una temporada de 82 partidos corre completa, con record consistente, ciudad reaccionando, agentes acumulando estado, narrativa post-partido y analytics persistido.
