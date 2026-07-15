# M2.8 — Orquestacion de partido

Fecha: 2026-05-24

## Objetivo

Conectar el avance de tiempo con el calendario deportivo: cuando `agent-service` publica `tiempo.dia_avanzado`, `team-service` detecta si hay partido ese dia y publica `partido.programado` para que `match-service` lo simule.

## Cambios hechos

### Team-service

- `team-service` ahora escucha `tiempo.dia_avanzado`.
- Se agrego `DispatchScheduledMatchForDate`.
- El flujo del handler es:
  - decodificar `tiempo.dia_avanzado`
  - buscar partido `scheduled` para `game_id` + `simulated_date`
  - marcarlo como `scheduled_dispatched`
  - cargar franquicia propia, rival y roster
  - armar `partido.programado`
  - publicar en NATS
- Si no hay partido para ese dia, no pasa nada.
- Si el evento se reintenta, el partido ya marcado no se vuelve a publicar.

### Payload completo

- `partido.programado` incluye jugadores.
- El roster propio sale de `team_roster_players`.
- El roster rival se genera de forma abstracta y deterministica a partir de:
  - `game_id`
  - `match_id`
  - `opponent_team_id`
  - rating del rival
- Esto mantiene `match-service` stateless sin abrir todavia el modelado profundo de la liga.

## Decision tomada

El estado `scheduled_dispatched` se guarda antes de publicar.

Esto prioriza no duplicar simulaciones ante retries del mismo `tiempo.dia_avanzado`. Si el publish falla despues de marcar, queda visible en logs y mas adelante se puede agregar una recuperacion de partidos `scheduled_dispatched` sin resultado. Para M2.8, la prioridad es evitar duplicados y cerrar el flujo feliz.

## Archivos tocados

- `services/team-service/cmd/main.go`
- `services/team-service/internal/domain/season.go`
- `services/team-service/internal/domain/season_test.go`
- `services/team-service/internal/persistence/store.go`
- `docs/Sesiones/MILESTONE2/INICIOM2.MD`

## Verificacion

Se corrieron correctamente:

```bash
GOCACHE=/tmp/pulsecity-team-gocache go test ./...
make test-go
make test-rust
make build
```

## Pendiente siguiente

`M2.9 — Stats y record`

Proximo objetivo recomendado:

- `team-service` consume `partido.terminado`
- actualiza victorias, derrotas, puntos a favor/en contra y estado del partido
- deja base para `season.patch` y el HUD deportivo
