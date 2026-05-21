# City Events

Contratos de ciudad para Milestone 2. El dueño de estas metricas es `city-service`.

## Subjects

- `ciudad.economia_cambio`
- `ciudad.suelo_actualizado`

## `ciudad.economia_cambio`

Publicado por `city-service` como reaccion a hechos del juego, principalmente `partido.terminado`.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:03Z",
  "schema_version": 1,
  "simulated_date": "2026-10-22",
  "source_event_id": "uuid",
  "source_subject": "partido.terminado",
  "fan_sentiment_delta": 0.04,
  "ticket_sales_delta": 0.03,
  "local_economy_delta": 0.01,
  "reason": "home_win|home_loss|winning_streak|losing_streak"
}
```

## `ciudad.suelo_actualizado`

Publicado por `city-service` cuando cambia el valor de suelo de una zona.

```json
{
  "event_id": "uuid",
  "game_id": "uuid",
  "occurred_at": "2026-05-21T00:00:03Z",
  "schema_version": 1,
  "simulated_date": "2026-10-22",
  "zone_id": "stadium_district",
  "land_value_delta": 0.02,
  "new_land_value": 1.12,
  "source_event_id": "uuid",
  "reason": "winning_streak"
}
```

Notas:

- M2 usa metricas globales y zona del estadio.
- `city-service` reacciona a hechos publicados; no recibe ordenes directas de `match-service`.
- el `gateway` traduce estos eventos a `city.patch`.
