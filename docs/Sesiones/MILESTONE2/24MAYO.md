# SESION24MAYO

## Milestone 2 — el sistema empieza a latir

Sesion larga de avance sobre `Milestone 2`, enfocada en activar la primera cadena sistemica real de PulseCity:

```text
tiempo -> calendario -> partido -> resultado -> record -> ciudad -> deltas al frontend
```

El foco no fue sumar gestion deportiva profunda, sino lograr que una partida pueda avanzar con tiempo simulado, disparar partidos automaticamente, persistir resultados y mostrar consecuencias visibles.

## Mini milestones trabajados

- `M2.4 — Loop hibrido`
- `M2.5 — HUD de tiempo`
- `M2.6 — Temporada inicial`
- `M2.7 — Match-service v1`
- `M2.8 — Orquestacion de partido`
- `M2.9 — Stats y record`
- `M2.10 — Ciudad late v1`

Todos quedaron marcados como `YA REALIZADO` en `INICIOM2.MD`.

## Objetivo general de la sesion

Convertir la base tecnica de M2 en un loop observable:

- el tiempo avanza solo cuando hay WebSocket activo
- el jugador puede pausar y cambiar velocidad
- `team-service` genera calendario, roster y rivales
- `match-service` simula partidos deterministas
- `team-service` persiste resultados, box scores y record
- `city-service` reacciona a resultados y rachas
- `gateway` reenvia deltas al frontend
- el frontend empieza a mostrar el pulso del sistema

## M2.4 — Loop hibrido

Se implemento el runtime de simulacion en `agent-service`.

El servicio ahora:

- escucha `tiempo.sesion_iniciada`
- escucha `tiempo.sesion_terminada`
- escucha `tiempo.velocidad_cambiada`
- escucha `tiempo.pausa_activada`
- corre un tick cada 100ms
- acumula dias segun velocidad x1, x5 o x20
- publica `tiempo.dia_avanzado` cuando corresponde
- no avanza si no hay sesion activa
- no avanza si esta en pausa

Decision importante: el loop duerme sin sesion activa. No hay compensacion offline.

## M2.5 — HUD de tiempo

Se conecto el control de tiempo al jugador.

Cambios principales:

- `gateway` publica `tiempo.sesion_iniciada` al abrir WebSocket validado
- `gateway` publica `tiempo.sesion_terminada` al cerrar la ultima conexion activa de una partida
- se agrego `POST /api/v1/games/{gameID}/time-control`
- el endpoint publica cambios de pausa y velocidad
- `gateway` escucha `tiempo.dia_avanzado` y reenvia `time.patch`
- frontend mantiene `TimeClientState`
- `CeremonyPage` muestra fecha simulada, pausa, x1, x5 y x20

Tambien se agrego proteccion para multiples clientes mirando la misma partida: el loop no se apaga hasta cerrar la ultima conexion.

## M2.6 — Temporada inicial

Se activo `team-service` como dueño del estado deportivo base.

Al crear/fundar una partida:

- `gateway` publica `mapa.generacion_iniciada` con `franchise_name` y `abbreviation`
- `team-service` escucha ese evento
- genera franquicia propia
- genera roster ficticio de 15 jugadores
- genera 30 rivales abstractos
- genera calendario de 82 partidos
- persiste franquicia, roster, rivales y calendario

La temporada arranca el `2026-10-22` y usa seeds deterministicas por partido.

## M2.7 — Match-service v1

Se implemento el primer simulador deterministico en `match-service`.

El servicio:

- consume `partido.programado`
- valida input completo
- exige rotacion minima de 5 jugadores por equipo
- publica `partido.iniciando`
- simula marcador sin empates
- genera box score
- genera momentos clave
- publica `partido.terminado`

Decision importante: `match-service` sigue siendo stateless. Todo lo que necesita llega en el payload.

## M2.8 — Orquestacion de partido

Se conecto fecha simulada con calendario y simulacion.

Flujo alcanzado:

```text
agent-service publica tiempo.dia_avanzado
team-service busca partido del dia
team-service marca scheduled_dispatched
team-service publica partido.programado
match-service publica partido.iniciando
match-service simula
match-service publica partido.terminado
```

La orquestacion quedo idempotente: si el mismo dia se procesa de nuevo, el partido ya despachado no se vuelve a publicar.

## M2.9 — Stats y record

`team-service` ahora consume `partido.terminado` y actualiza el estado deportivo persistido.

Se implemento:

- marca de partido como `final`
- persistencia de `home_score`, `away_score`, `winner_team_id` y `played_at`
- persistencia de box scores en `team_player_box_scores`
- recalculo de record desde partidos finalizados
- tabla `team_records`
- publicacion de `season.patch`

El frontend mantiene `SeasonClientState` y muestra el record real en la ceremonia.

Decision importante: el record se recalcula desde partidos finalizados para que un retry de `partido.terminado` no duplique victorias, derrotas ni puntos.

## M2.10 — Ciudad late v1

`city-service` empezo a reaccionar a resultados.

Se agregaron metricas urbanas por partida:

- fan sentiment
- ticket sales index
- local economy index
- stadium district land value
- win streak
- loss streak
- last match processed

Reglas iniciales:

- victoria local mejora animo, demanda de entradas y economia local
- victoria visitante tambien mejora, pero menos
- derrota local reduce momentum con mas fuerza
- derrota visitante tambien enfria, pero menos
- rachas positivas desde 3 victorias suben el valor del suelo cerca del estadio
- rachas negativas desde 3 derrotas enfrian el distrito

Persistencia:

- `city_metrics`
- `city_processed_matches`

La idempotencia se controla por `(game_id, match_id)`.

Eventos publicados:

- `ciudad.economia_cambio`
- `ciudad.suelo_actualizado`
- `city.patch`

El frontend mantiene `CityClientState`, muestra animo/suelo/demanda y activa un pulso visual sobre el estadio cuando llega un cambio urbano.

## Migraciones agregadas

- `010_create_team_season_tables.sql`
- `011_create_city_metrics_tables.sql`

## Contratos actualizados

- `shared/events/map_generation.md`
- `shared/events/matches.md`
- `shared/events/city.md`

## Estado tecnico alcanzado

Al cierre de esta sesion, el sistema ya puede:

- fundar una partida
- activar WebSocket
- arrancar/detener el loop segun sesion activa
- pausar y cambiar velocidad
- avanzar dias simulados
- encontrar partidos por fecha
- simular partidos automaticamente
- persistir resultados y box scores
- actualizar record
- publicar deltas deportivos
- actualizar metricas urbanas por resultado
- publicar deltas urbanos
- mostrar en frontend tiempo, record y pulso de ciudad

## Lo que todavia no hace el usuario

El usuario todavia no tiene gestion profunda.

No hay aun:

- trades
- free agency
- rotacion editable
- salary cap real
- decisiones de ciudad
- agentes core visibles con estado real
- narrativa post-partido conectada al resultado
- calendario completo visible como pantalla de gestion

Esto es consistente con el alcance actual: M2 esta activando el loop sistemico antes de abrir decisiones profundas.

## Verificacion

Se corrieron correctamente:

```bash
GOCACHE=/tmp/pulsecity-team-gocache go test ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go test ./...
GOCACHE=/tmp/pulsecity-city-gocache go test ./...
npm run build --prefix frontend
make test-go
make test-rust
make build
```

## Documentacion relacionada

Notas granulares agregadas en `docs/02_Progress`:

- `2026-05-24_M2_4_Loop_Hibrido.md`
- `2026-05-24_M2_5_HUD_Tiempo.md`
- `2026-05-24_M2_6_Temporada_Inicial.md`
- `2026-05-24_M2_7_Match_Service_V1.md`
- `2026-05-24_M2_8_Orquestacion_Partido.md`
- `2026-05-24_M2_9_Stats_Record.md`
- `2026-05-24_M2_10_Ciudad_Late_V1.md`

## Estado real al cierre

M2 ya paso de base estructural a loop funcional observable.

El sistema todavia no es una experiencia jugable completa, pero ya tiene una cadena viva:

```text
sesion activa -> tiempo -> calendario -> partido -> resultado -> record -> ciudad -> frontend
```

## Continuacion — M2.11 Agentes core v1

Despues del loop ciudad/temporada, se implemento el primer estado real de agentes core.

Cambios principales:

- `agent-service` agrego dominio `agents`
- se crearon estados iniciales para Owner, Head Coach, CFO, Director de Scouting y Sports Psychologist
- `agent-service` consume `partido.terminado`
- cada resultado actualiza variables por agente
- se persiste estado en `agent_core_states`
- se agrego idempotencia por `(game_id, match_id)` en `agent_processed_matches`
- se publica `agente.estado_cambio`
- `gateway` traduce esos eventos a `agent.patch`
- frontend muestra los agentes core en la ceremonia

Migracion agregada:

- `012_create_agent_core_state_tables.sql`

Documentacion relacionada:

- `docs/02_Progress/2026-05-24_M2_11_Agentes_Core_V1.md`

Verificacion:

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
GOCACHE=/tmp/pulsecity-gateway-gocache go test ./...
npm run build --prefix frontend
```

## Continuacion — M2.12 Narrativa post-partido

Se conecto el `narrative-service` al loop de partidos.

Cambios principales:

- `narrative-service` consume `partido.terminado`
- espera 250-500ms antes de generar texto
- genera narrativa templateada/rule-based segun resultado, localia, margen, momentos clave y racha disponible
- persiste eventos en `narrative_events`
- controla idempotencia por `(game_id, source_match_id, kind)`
- publica `narrativa.evento_generado`
- el gateway ya reenvia ese evento al frontend como `narrative.event`

Migracion agregada:

- `013_create_narrative_events_tables.sql`

Documentacion relacionada:

- `docs/02_Progress/2026-05-24_M2_12_Narrativa_Post_Partido.md`

Verificacion:

```bash
GOCACHE=/tmp/pulsecity-narrative-gocache go test ./...
```

## Continuacion — M2.13 Frontend jugable de M2

Se reorganizo la ceremonia para que M2 sea mas observable como loop jugable.

Cambios principales:

- panel de temporada viva con record, partidos jugados sobre 82 y diferencial promedio
- lista de resultados recientes alimentada por `season.patch`
- inbox narrativo alimentado por `narrative.event`
- correccion del frontend para no descartar narrativa post-partido luego de responder el owner intro
- separacion entre eventos tecnicos e inbox narrativo
- reseteo completo de estado visual al crear una nueva partida

Documentacion relacionada:

- `docs/02_Progress/2026-05-24_M2_13_Frontend_Jugable.md`

Verificacion:

```bash
npm run build --prefix frontend
```

## Continuacion — M2.14 Analytics basico

Se activo `analytics-service` como consumidor de series sin UI.

Cambios principales:

- conexion a Postgres/TimescaleDB con `pgx`
- schema propio de analytics
- ingesta de `partido.terminado`
- persistencia de resultados y box scores
- ingesta de `ciudad.economia_cambio`
- persistencia de fan sentiment, ticket sales y economia local
- ingesta de `ciudad.suelo_actualizado`
- persistencia de valor de suelo por zona
- ingesta de `agente.estado_cambio`
- persistencia historica de estado de agentes
- idempotencia por claves primarias de partido/evento/metrica

Migracion agregada:

- `014_create_analytics_timeseries_tables.sql`

Documentacion relacionada:

- `docs/02_Progress/2026-05-24_M2_14_Analytics_Basico.md`

Verificacion:

```bash
GOCACHE=/tmp/pulsecity-analytics-gocache go test ./...
```

## Cierre — M2.15 Temporada completa

Se corrio un smoke end-to-end real para cerrar M2.

Servicios levantados:

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

Partida validada:

- `game_id`: `0ba7bc3d-f6ca-49db-9597-5196485b0e65`
- velocidad: x20

Resultado en DB:

- `team_schedule`: 82 finalizados de 82
- `team_records`: 37-45
- `city_processed_matches`: 82
- `agent_processed_matches`: 82
- `narrative_events` post_match: 82
- `analytics_match_results`: 82
- `analytics_player_box_scores`: 1640
- `analytics_city_metric_points`: 246
- `analytics_land_value_points`: 82
- `analytics_agent_state_points`: 1394

Bug encontrado y corregido:

- a x20, `agent-service` podia publicar `tiempo.dia_avanzado` con `days_processed > 1`
- `team-service` solo procesaba la fecha final del evento
- eso salteaba partidos en fechas intermedias
- se agrego `CoveredDayAdvancedEvents`
- ahora `team-service` procesa todas las fechas cubiertas por el evento de tiempo

Documentacion relacionada:

- `docs/02_Progress/2026-05-24_M2_15_Temporada_Completa.md`

Verificacion:

```bash
GOCACHE=/tmp/pulsecity-team-gocache go test ./...
make test-go
make test-rust
npm run build --prefix frontend
make build
```

Con esto `Milestone 2` queda cerrado funcionalmente.

## Pendiente siguiente

`M2.11 — Agentes core v1`

Proximo objetivo recomendado:

- crear estado persistido para Owner, Head Coach, CFO, Director de Scouting y Sports Psychologist
- hacer que reaccionen a `partido.terminado`
- publicar `agente.estado_cambio`
- reenviar `agent.patch` por WebSocket
- mostrar resumen visible de agentes en frontend
