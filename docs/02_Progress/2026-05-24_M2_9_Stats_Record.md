# M2.9 — Stats y record

Fecha: 2026-05-24

## Objetivo

Hacer que `team-service` consuma `partido.terminado`, persista el resultado, actualice record deportivo y publique un delta `season.patch` para el frontend.

## Cambios hechos

### Team-service

- `team-service` ahora escucha `partido.terminado`.
- Se agrego `ApplyMatchFinished`.
- Al recibir un resultado:
  - marca el partido como `final`
  - guarda `home_score`
  - guarda `away_score`
  - guarda `winner_team_id`
  - guarda `played_at`
  - persiste box scores por jugador
  - recalcula el record desde partidos finalizados
  - actualiza `team_records`
  - publica `season.patch`

### Persistencia

- `team_schedule` ahora guarda resultado final.
- Se agrego `team_player_box_scores`.
- Se agrego `team_records`.
- La migracion `010_create_team_season_tables.sql` quedo ampliada con estas tablas/columnas.

### Gateway y frontend

- `gateway` escucha `season.patch` y lo reenvia por WebSocket.
- Frontend agrega `SeasonClientState`.
- `CeremonyPage` muestra el record real.

## Decision tomada

El record se recalcula desde partidos finalizados despues de cada resultado.

Esto es mas robusto que sumar victorias/derrotas de forma incremental, porque si llega un retry de `partido.terminado`, el sistema no duplica estado. El costo es minimo en M2 porque son solo 82 partidos de una franquicia.

## Archivos tocados

- `services/team-service/cmd/main.go`
- `services/team-service/internal/domain/events.go`
- `services/team-service/internal/domain/events_test.go`
- `services/team-service/internal/persistence/store.go`
- `services/gateway/cmd/main.go`
- `services/gateway/internal/domain/map_generation.go`
- `frontend/src/types.ts`
- `frontend/src/features/new-game/hooks/useNewGameFlow.ts`
- `frontend/src/features/new-game/NewGameFlow.tsx`
- `frontend/src/features/new-game/components/CeremonyPage.tsx`
- `frontend/src/features/new-game/helpers.ts`
- `db/migrations/010_create_team_season_tables.sql`
- `docs/Sesiones/MILESTONE2/INICIOM2.MD`

## Verificacion

Se corrieron correctamente:

```bash
GOCACHE=/tmp/pulsecity-team-gocache go test ./...
npm run build --prefix frontend
make test-go
make test-rust
make build
```

## Pendiente siguiente

`M2.10 — Ciudad late v1`

Proximo objetivo recomendado:

- `city-service` consume `partido.terminado`
- ajusta fan sentiment, ticket sales, economia local y valor del suelo
- publica `ciudad.economia_cambio` y/o `ciudad.suelo_actualizado`
