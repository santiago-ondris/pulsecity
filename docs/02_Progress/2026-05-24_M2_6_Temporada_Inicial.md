# M2.6 — Temporada inicial

Fecha: 2026-05-24

## Objetivo

Hacer que `team-service` cree la estructura minima de temporada regular para una partida: franquicia propia, roster ficticio, rivales abstractos y calendario de 82 partidos.

## Cambios hechos

### Gateway

- `mapa.generacion_iniciada` ahora incluye:
  - `franchise_name`
  - `abbreviation`
- `map-service` sigue pudiendo consumir el evento porque solo necesita `game_id` y `city_name`.

### Team-service

- Se agrego `GameStartedEvent` para consumir `mapa.generacion_iniciada`.
- Se agrego generador deterministico de temporada inicial.
- Se genera:
  - equipo propio
  - 15 jugadores ficticios
  - 30 rivales abstractos
  - calendario de 82 partidos
- Cada partido tiene:
  - `match_id`
  - fecha simulada
  - local/visitante
  - rival
  - seed deterministica
  - estado `scheduled`
- Se agrego persistencia propia para:
  - `team_franchises`
  - `team_roster_players`
  - `team_opponents`
  - `team_schedule`

## Decision tomada

Para M2.6, `team-service` reacciona a `mapa.generacion_iniciada`.

No es el nombre perfecto para una partida creada, pero hoy es el primer evento canonico que se publica cuando nace una partida. Evita agregar un evento nuevo antes de necesitarlo y mantiene el milestone chico. Si mas adelante aparece un flujo de creacion mas rico, podemos introducir `partida.creada` o equivalente y dejar `mapa.generacion_iniciada` solo para `map-service`.

## Archivos tocados

- `services/gateway/internal/domain/map_generation.go`
- `services/gateway/internal/handlers/http.go`
- `services/team-service/cmd/main.go`
- `services/team-service/go.mod`
- `services/team-service/go.sum`
- `services/team-service/internal/domain/season.go`
- `services/team-service/internal/domain/season_test.go`
- `services/team-service/internal/persistence/store.go`
- `db/migrations/010_create_team_season_tables.sql`
- `shared/events/map_generation.md`
- `docs/Sesiones/MILESTONE2/INICIOM2.MD`

## Verificacion

Se corrieron correctamente:

```bash
GOCACHE=/tmp/pulsecity-team-gocache go test ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go test ./...
make test-go
make build
```

## Pendiente siguiente

`M2.7 — Match-service v1`

Proximo objetivo recomendado:

- implementar simulacion deterministica en `match-service`
- producir resultado, box score y 3-5 momentos clave con el input completo
- mantener `match-service` stateless y testeable
