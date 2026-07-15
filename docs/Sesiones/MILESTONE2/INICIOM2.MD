# Milestone 2 — Inicio

## El sistema late

`Milestone 2` arranca con una decision importante: no se trata de sumar features sueltas, sino de activar la primera cadena sistemica real de PulseCity.

La cadena base del milestone es:

```text
tiempo -> calendario -> partido -> resultado -> reacciones -> narrativa -> deltas al frontend
```

El objetivo no es todavia dar al jugador gestion deportiva profunda. El objetivo es que el mundo empiece a avanzar solo, que los partidos ocurran en fechas simuladas, que la ciudad reaccione, que los agentes acumulen estado, y que el jugador pueda observar el pulso del sistema partido a partido.

## Criterio de done de M2

Una partida creada en `Milestone 1` puede avanzar una temporada regular completa de 82 partidos del equipo propio, con:

- loop de tiempo hibrido
- pausa, x1, x5 y x20
- calendario de temporada regular
- partidos simulados con box score
- record actualizado del equipo
- ciudad reaccionando a resultados y rachas
- 5 agentes core con estado real
- eventos narrativos post-partido
- frontend mostrando HUD, calendario/resultados, inbox y pulso visual de la ciudad

## Decisiones cerradas

### Estrategia de servicios

M2 se implementa con servicios separados desde el primer mini milestone.

No vamos a meter logica nueva de M2 temporalmente dentro de `gateway` para despues extraerla. Eso respeta el canon de ownership desde el inicio y evita deuda tecnica temprana.

Servicios activos de M2:

- `agent-service`
- `team-service`
- `match-service`
- `city-service`
- `analytics-service`
- `narrative-service`
- `gateway`
- `frontend`

`gateway` queda como punto de entrada/salida:

- recibe acciones del frontend
- publica eventos NATS
- mantiene WebSocket
- reenvia deltas al frontend
- puede hacer lecturas REST internas para ensamblar respuestas

`gateway` no calcula dias, no simula partidos, no decide reacciones y no escribe estado de otros servicios.

### Control activo del jugador

M2 queda en modo `observacion gestionable`.

El jugador puede:

- controlar el tiempo
- pausar/reanudar
- cambiar velocidad
- ver calendario
- ver resultados
- leer eventos narrativos
- ver estado resumido de agentes
- ver la ciudad reaccionar

El jugador no tiene todavia:

- trades
- free agency
- salary cap profundo
- rotacion editable
- gestion manual de minutos
- chat libre con agentes
- directorio completo de personajes
- decisiones complejas de ciudad

Esto queda diferido a milestones posteriores para no mezclar el primer loop sistemico con gestion deportiva profunda.

### Alcance de liga

M2 simula 82 partidos del equipo propio contra rivales abstractos generados.

No se modelan todavia:

- standings completos de 32 equipos
- calendario global de toda la NBA
- ciudades rivales
- agentes de otros GMs
- trades de otros equipos
- IA completa de la liga

Cada rival abstracto necesita solo lo suficiente para simular:

- nombre
- abreviatura
- rating general
- perfil tactico simple
- rating ofensivo
- rating defensivo
- ritmo
- fuerza de localia si aplica

Esto permite cumplir el espiritu de una temporada regular completa sin abrir todavia la complejidad de toda la liga.

### Ownership del tiempo

El loop de tiempo vive en `agent-service`.

Razon:

- el canon ya define que `agent-service` contiene el simulation loop
- el loop esta ligado a agentes, sesion activa y avance diario
- evita crear un `time-service` antes de que el sistema lo justifique
- evita poner logica de simulacion en `gateway`

`gateway` publica:

- `tiempo.sesion_iniciada`
- `tiempo.sesion_terminada`
- `tiempo.velocidad_cambiada`
- `tiempo.pausa_activada`

`agent-service` consume esos eventos, mantiene el estado de simulacion y publica:

- `tiempo.dia_avanzado`

### Persistencia

M2 usa una unica base Postgres/Timescale local, con ownership logico por tablas.

Esto significa:

- todos los servicios pueden conectarse a la misma DB fisica en desarrollo
- cada tabla tiene un servicio dueño
- ningun servicio escribe tablas de otro servicio
- compartir DB fisica no rompe el ownership mientras respetemos esta regla

Ejemplo:

- `agent-service` escribe estado de tiempo y agentes
- `team-service` escribe calendario, roster, partidos y stats acumuladas
- `city-service` escribe metricas urbanas
- `narrative-service` escribe eventos narrativos
- `analytics-service` escribe series historicas append-only

### Narrativa

M2 mantiene narrativa templateada/rule-based.

No se integra proveedor LLM real todavia.

La decision no es tecnica, es de foco:

- M1 ya dejo asentado que el LLM real se difiere por costo
- M2 necesita probar primero que el contexto sistemico sea real
- un texto templateado con estado real aporta mas que un LLM conectado a datos pobres

El `narrative-service` sigue siendo el dueño de los eventos narrativos visibles.

Para eventos post-partido:

- escucha `partido.terminado`
- espera 250-500ms
- asume que `agent-service` ya proceso consecuencias
- genera un evento narrativo templateado usando resultado, momentos clave y estado actualizado
- publica `narrativa.evento_generado`

### Ciudad

La ciudad en M2 reacciona con metricas globales y zona del estadio.

Incluye:

- ticket sales
- fan sentiment
- economia local simple
- valor del suelo cerca del estadio
- reaccion a victorias, derrotas y rachas

No incluye todavia:

- construccion de edificios
- zoning editable
- presupuesto municipal profundo
- permisos del alcalde
- efectos multi-zona complejos
- city-builder completo

Eso queda para `Milestone 4`.

### Simulador de partidos

`match-service` usa un simulador deterministico con seed.

Mismo input debe producir mismo output.

Razon:

- facilita tests
- facilita debugging
- permite reproducir una temporada
- evita resultados imposibles de explicar

No se implementa todavia un modelo posesion a posesion detallado. El simulador M2 debe ser suficientemente rico para generar box score, resultado y momentos clave, pero no debe intentar ser el simulador definitivo del juego.

### Frontend

La UI minima de M2 es:

- HUD de fecha simulada
- controles de tiempo
- record del equipo
- proximo partido
- calendario/resultados
- inbox narrativo
- resumen visible de agentes core
- pulso visual de ciudad alrededor del estadio

No se implementa todavia:

- Panel GM completo
- gestion de roster profunda
- finanzas completas
- directorio completo de agentes
- chat con agentes
- charts avanzados

### Agentes core

M2 activa 5 agentes:

- Owner
- Head Coach
- CFO
- Director de Scouting
- Sports Psychologist

Cada uno debe tener:

- identidad
- personalidad base
- variables reales de estado
- reaccion a partidos y rachas
- estado visible resumido para el jugador

No se implementa todavia el sistema completo de relaciones ni causalidad profunda de M3.

### Analytics

`analytics-service` en M2 ingesta datos, pero no expone charts visibles todavia.

Debe guardar:

- resultados de partidos
- box scores relevantes
- metricas de ciudad por dia o evento
- cambios de valor de suelo
- evolucion de ticket sales y fan sentiment

La UI analitica avanzada queda para mas adelante.

## Interfaces y eventos principales

### Tiempo

```text
tiempo.sesion_iniciada
tiempo.sesion_terminada
tiempo.velocidad_cambiada
tiempo.pausa_activada
tiempo.dia_avanzado
```

`tiempo.dia_avanzado` debe incluir:

- `game_id`
- `fecha_simulada`
- `velocidad`
- `dias_procesados`

### Partidos

```text
partido.programado
partido.iniciando
partido.terminado
```

`partido.terminado` debe incluir:

- `game_id`
- `partido_id`
- `fecha_simulada`
- equipo local
- equipo visitante
- puntos local
- puntos visitante
- box score
- momentos clave
- seed usada

### Ciudad

```text
ciudad.economia_cambio
ciudad.suelo_actualizado
```

La ciudad no recibe ordenes del `match-service`; reacciona a hechos publicados.

### Agentes

```text
agente.estado_cambio
agente.evento_critico
```

En M2, `agente.estado_cambio` alcanza para mostrar que los agentes estan vivos sin abrir todavia toda la matriz relacional.

### Narrativa

```text
narrativa.evento_generado
```

El evento narrativo post-partido debe poder mostrar:

- quien habla
- por que habla
- urgencia
- titulo
- cuerpo
- metadata del partido
- opciones simples si aplica

### WebSocket

Deltas nuevos esperados:

```text
time.patch
season.patch
match.result
city.patch
agent.patch
narrative.event
```

Regla no negociable:

El WebSocket sigue enviando deltas, no estado completo.

## Mini milestones propuestos

### M2.0 — Documento rector de M2 — YA REALIZADO

Objetivo:

Dejar asentado en `docs/Sesiones/MILESTONE2/INICIOM2.MD` el alcance, decisiones, mini milestones, interfaces, edge cases y criterios de done.

Done:

- este documento existe
- las decisiones principales estan escritas
- el milestone queda listo para ejecutarse por partes

### M2.1 — Scaffolding de servicios M2 — YA REALIZADO

Objetivo:

Crear estructura minima real para los servicios vacios.

Servicios:

- `agent-service` en Rust
- `match-service` en Rust
- `team-service` en Go
- `city-service` en Go
- `analytics-service` en Go

Done:

- cada servicio compila
- cada servicio tiene entrypoint
- cada servicio puede conectarse a NATS
- Makefile incluye build/test/run basicos
- no hay logica de gameplay todavia

Resultado 2026-05-21:

- creados `agent-service` y `match-service` en Rust con `Cargo.toml`, `src/main.rs` y `src/lib.rs`
- creados `team-service`, `city-service` y `analytics-service` en Go con `cmd`, `internal/domain`, `internal/handlers` e `internal/nats`
- cada servicio nuevo tiene entrypoint minimo que lee `NATS_URL` y conecta a NATS
- `Makefile` actualizado para build/test/run de los servicios M2
- pruebas corridas: `make test-rust`, `make test-go`, `make build`

### M2.2 — Contratos compartidos M2 — YA REALIZADO

Objetivo:

Definir payloads canonicos para eventos de tiempo, partidos, ciudad, agentes y narrativa.

Done:

- eventos documentados
- tipos internos alineados por servicio
- nombres respetan `entidad.accion`
- no se rompen eventos de M1

Resultado 2026-05-21:

- documentados contratos M2 en `shared/events`
- definidos payloads canonicos para tiempo, partidos, ciudad, agentes, narrativa y deltas WebSocket
- agregados tipos internos alineados en `agent-service`, `match-service`, `team-service`, `city-service` y `analytics-service`
- se mantiene M1 sin cambios de comportamiento: `mapa.*` y owner intro siguen vigentes
- pruebas corridas: `make test-rust`, `make test-go`, `make build`

### M2.3 — Estado de simulacion — YA REALIZADO

Objetivo:

Persistir estado basico de tiempo por partida.

Incluye:

- fecha simulada actual
- velocidad actual
- pausa
- sesion activa
- ultimo tick procesado

Done:

- `agent-service` puede cargar/guardar estado de simulacion
- una partida reanuda donde quedo
- sin sesion activa no avanza

Resultado 2026-05-21:

- creada tabla `agent_simulation_state` via migracion `009_create_agent_simulation_state.sql`
- `agent-service` tiene modulo `simulation` con estado base, velocidades validas y regla `can_advance`
- `agent-service` tiene modulo `persistence` para asegurar schema, cargar, guardar e inicializar estado por `game_id`
- el binario carga/inicializa estado usando `DATABASE_URL` y `GAME_ID`
- el loop de 100ms queda explicitamente diferido a M2.4
- pruebas corridas: `make test-rust`, `make test-go`, `make build`

### M2.4 — Loop hibrido — YA REALIZADO

Objetivo:

Implementar el loop de 100ms en `agent-service`.

Done:

- con sesion activa avanza el tiempo
- sin sesion activa duerme
- pausa detiene avance
- x1/x5/x20 modifican acumulacion
- emite `tiempo.dia_avanzado` cuando corresponde

Resultado 2026-05-24:

- agregado runtime de simulacion en `agent-service` con tick cada 100ms
- `agent-service` escucha `tiempo.sesion_iniciada`, `tiempo.sesion_terminada`, `tiempo.velocidad_cambiada` y `tiempo.pausa_activada`
- el estado persistido se actualiza ante cambios de sesion, pausa y velocidad
- la acumulacion temporal respeta x1, x5 y x20
- sin sesion activa o con pausa, el acumulador no avanza y se limpia la acumulacion parcial
- al completar dias simulados, se avanza `current_simulated_date`, se guarda `last_tick_processed_at` y se publica `tiempo.dia_avanzado`
- agregados tests unitarios para acumulacion por velocidad, pausa/sesion inactiva y avance de fechas
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`, `make test-rust`, `make build`

### M2.5 — HUD de tiempo — YA REALIZADO

Objetivo:

Mostrar en frontend fecha simulada y controles de tiempo.

Done:

- el usuario puede cambiar velocidad
- el usuario puede pausar
- el frontend recibe `time.patch`
- el HUD refleja el estado real del backend

Resultado 2026-05-24:

- `gateway` publica `tiempo.sesion_iniciada` al abrir WebSocket validado para una partida
- `gateway` publica `tiempo.sesion_terminada` al cerrar la ultima conexion WebSocket activa de la partida
- agregado endpoint `POST /api/v1/games/{gameID}/time-control`
- el endpoint valida ownership de partida y publica `tiempo.velocidad_cambiada` y/o `tiempo.pausa_activada`
- `gateway` escucha `tiempo.dia_avanzado` y lo traduce a delta WebSocket `time.patch`
- el frontend mantiene `TimeClientState` y aplica `time.patch`
- `CeremonyPage` muestra HUD de fecha simulada, pausa, x1, x5 y x20
- los controles del HUD envian cambios al gateway y el gateway reenvia el delta aceptado a clientes conectados
- agregado test para que multiples clientes de la misma partida no apaguen el loop hasta cerrar la ultima conexion
- pruebas corridas: `go test ./...` en `services/gateway`, `npm run build --prefix frontend`, `make test-go`, `make build`

### M2.6 — Temporada inicial — YA REALIZADO

Objetivo:

`team-service` genera la estructura minima de temporada regular.

Incluye:

- roster ficticio minimo del equipo propio
- rivales abstractos
- calendario de 82 partidos
- partidos con fecha simulada

Done:

- al crear/fundar partida existe calendario
- se puede consultar proximo partido
- no hay standings globales todavia

Resultado 2026-05-24:

- `gateway` amplia `mapa.generacion_iniciada` con `franchise_name` y `abbreviation`
- `team-service` escucha `mapa.generacion_iniciada` como hecho de partida nueva
- agregado generador deterministico de temporada inicial en `team-service`
- se genera franquicia propia, roster ficticio de 15 jugadores, 30 rivales abstractos y calendario de 82 partidos
- el calendario alterna local/visitante y empieza el `2026-10-22`
- se agregan seeds deterministicas por partido para reproducibilidad futura del simulador
- agregada persistencia propia de `team-service` para franquicia, roster, rivales y calendario
- agregada migracion `010_create_team_season_tables.sql`
- se agrego helper de dominio para consultar el proximo partido programado desde una lista de calendario
- pruebas corridas: `GOCACHE=/tmp/pulsecity-team-gocache go test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go test ./...`, `make test-go`, `make build`

### M2.7 — Match-service v1 — YA REALIZADO

Objetivo:

Simular un partido deterministico con input completo.

Done:

- recibe payload completo
- no consulta otros servicios
- genera resultado
- genera box score
- genera 3-5 momentos clave
- mismo input produce mismo output

Resultado 2026-05-24:

- agregado simulador deterministico en `match-service`
- definido `MatchSimulationInput` con equipos, jugadores, fecha, partido y seed
- el simulador falla si falta `game_id`, `match_id`, fecha simulada o rotacion minima de 5 jugadores por equipo
- genera marcador sin empates, ganador, box score y 4 momentos clave
- el box score distribuye puntos, minutos, rebotes, asistencias, robos, tapones y perdidas
- mismo input y seed producen exactamente el mismo resultado
- seed distinta produce resultados distintos
- `match-service` consume `partido.programado`, publica `partido.iniciando` y luego `partido.terminado`
- `partido.programado` queda ampliado con `players` para cumplir el contrato de input completo
- pruebas corridas: `cargo test --manifest-path services/match-service/Cargo.toml`, `make test-rust`, `make test-go`, `make build`

### M2.8 — Orquestacion de partido — YA REALIZADO

Objetivo:

Conectar fecha simulada con calendario y simulacion.

Flujo:

```text
agent-service publica tiempo.dia_avanzado
team-service detecta partido del dia
team-service publica partido.programado o request equivalente
match-service publica partido.iniciando
match-service simula
match-service publica partido.terminado
```

Done:

- un partido ocurre automaticamente en su fecha
- no se duplica si el evento se reintenta
- el gateway puede reenviar resultado al frontend

Resultado 2026-05-24:

- `team-service` escucha `tiempo.dia_avanzado`
- al recibir un dia, busca un partido `scheduled` para `game_id` + `simulated_date`
- si no hay partido ese dia, no publica nada
- si hay partido, lo marca como `scheduled_dispatched` antes de publicar para evitar duplicados ante retries
- arma `partido.programado` con equipos, seed, roster propio y jugadores abstractos del rival
- los jugadores abstractos del rival se derivan determinísticamente del rival, `match_id` y `game_id`
- `match-service` recibe `partido.programado`, publica `partido.iniciando`, simula y publica `partido.terminado`
- pruebas corridas: `GOCACHE=/tmp/pulsecity-team-gocache go test ./...`, `make test-go`, `make test-rust`, `make build`

### M2.9 — Stats y record — YA REALIZADO

Objetivo:

`team-service` consume `partido.terminado` y actualiza estado deportivo.

Incluye:

- victorias
- derrotas
- puntos a favor/en contra
- box scores acumulados basicos
- ultimo resultado
- proximo partido

Done:

- el HUD puede mostrar record real
- el calendario puede mostrar resultados jugados

Resultado 2026-05-24:

- `team-service` consume `partido.terminado`
- marca el partido como `final` en `team_schedule`
- persiste `home_score`, `away_score`, `winner_team_id` y `played_at`
- persiste box scores en `team_player_box_scores`
- recalcula record desde partidos finalizados para mantener idempotencia
- mantiene resumen en `team_records`
- publica `season.patch` con victorias, derrotas, puntos a favor/en contra y ultimo resultado
- `gateway` escucha `season.patch` y lo reenvia por WebSocket
- frontend mantiene `SeasonClientState` y muestra el record en `CeremonyPage`
- si llega dos veces el mismo `partido.terminado`, no duplica record ni vuelve a publicar patch
- pruebas corridas: `GOCACHE=/tmp/pulsecity-team-gocache go test ./...`, `npm run build --prefix frontend`, `make test-go`, `make test-rust`, `make build`

### M2.10 — Ciudad late v1 — YA REALIZADO

Objetivo:

`city-service` reacciona a resultados.

Incluye:

- victoria local mejora ticket sales y fan sentiment
- derrota local reduce momentum
- rachas positivas suben valor del suelo cerca del estadio
- rachas negativas enfrían economia local

Done:

- se publican `ciudad.economia_cambio` y/o `ciudad.suelo_actualizado`
- el frontend puede mostrar pulso visual sobre zona estadio
- las metricas persisten

Resultado 2026-05-24:

- `city-service` consume `partido.terminado`
- calcula reacciones urbanas sobre fan sentiment, ticket sales, economia local y valor del suelo del distrito del estadio
- victorias y derrotas impactan distinto si el partido fue local o visitante
- rachas positivas desde 3 victorias elevan el valor del suelo cerca del estadio
- rachas negativas desde 3 derrotas enfrían el distrito del estadio
- las metricas persisten en `city_metrics`
- la idempotencia se controla con `city_processed_matches` por `(game_id, match_id)`
- se agrego migracion `011_create_city_metrics_tables.sql`
- `city-service` publica `ciudad.economia_cambio`, `ciudad.suelo_actualizado` y `city.patch`
- `gateway` escucha `city.patch` y lo reenvia por WebSocket
- frontend mantiene `CityClientState`, muestra animo/suelo/demanda y activa pulso visual sobre el estadio cuando llega un cambio urbano
- pruebas corridas: `GOCACHE=/tmp/pulsecity-city-gocache go test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go test ./...`, `npm run build --prefix frontend`, `make test-go`, `make test-rust`, `make build`

### M2.11 — Agentes core v1 — YA REALIZADO

Objetivo:

Crear estado real para los 5 agentes core.

Variables iniciales sugeridas:

Owner:

- confianza deportiva
- confianza negocio
- paciencia restante
- satisfaccion general

Head Coach:

- confianza en GM
- satisfaccion con roster
- presion por resultados
- relacion con vestuario

CFO:

- confianza financiera
- alerta de presupuesto
- conservadurismo financiero

Director de Scouting:

- confianza en criterio
- motivacion
- precision percibida

Sports Psychologist:

- clima de vestuario
- nivel de alerta emocional
- confianza de jugadores

Done:

- reaccionan a `partido.terminado`
- publican `agente.estado_cambio`
- el frontend puede mostrar resumen de estado

Resultado 2026-05-24:

- agregado modulo `agents` en `agent-service`
- definidos estados iniciales para Owner, Head Coach, CFO, Director de Scouting y Sports Psychologist
- `agent-service` escucha `partido.terminado`
- cada resultado actualiza variables especificas por agente segun victoria/derrota, localia y margen
- el estado se persiste en `agent_core_states`
- la idempotencia se controla con `agent_processed_matches` por `(game_id, match_id)`
- se publica un `agente.estado_cambio` por agente core actualizado
- `gateway` traduce `agente.estado_cambio` a `agent.patch`
- frontend mantiene `AgentClientStates` y muestra resumen compacto en `CeremonyPage`
- agregada migracion `012_create_agent_core_state_tables.sql`
- contratos actualizados en `shared/events/agents.md` y `shared/events/websocket_deltas.md`
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`, `GOCACHE=/tmp/pulsecity-gateway-gocache go test ./...`, `npm run build --prefix frontend`

### M2.12 — Narrativa post-partido — YA REALIZADO

Objetivo:

Generar eventos narrativos basicos despues de partidos.

Done:

- `narrative-service` escucha `partido.terminado`
- espera 250-500ms
- genera texto segun resultado, racha y momentos clave
- publica `narrativa.evento_generado`
- gateway reenvia al inbox

Resultado 2026-05-24:

- `narrative-service` escucha `partido.terminado`
- el handler espera 250-500ms antes de generar narrativa
- la espera ocurre en goroutine para no bloquear el dispatcher de NATS
- se genera `narrativa.evento_generado` templateado/rule-based
- el texto contempla victoria/derrota, localia, margen, maximo anotador propio, momentos clave y racha si esta disponible
- se persiste en `narrative_events`
- la idempotencia se controla con indice unico por `(game_id, source_match_id, kind)`
- se agrego migracion `013_create_narrative_events_tables.sql`
- contrato actualizado en `shared/events/narrative.md`
- pruebas corridas: `GOCACHE=/tmp/pulsecity-narrative-gocache go test ./...`

### M2.13 — Frontend jugable de M2 — YA REALIZADO

Objetivo:

Convertir el loop backend en experiencia visible.

Incluye:

- HUD de fecha/velocidad/record
- calendario o lista de partidos
- resultados recientes
- inbox narrativo
- resumen de agentes core
- pulso visual cerca del estadio

Done:

- el jugador puede mirar una temporada avanzar
- cada partido genera feedback visible
- la ciudad y los agentes se sienten conectados al resultado

Resultado 2026-05-24:

- `CeremonyPage` reorganizada como pantalla observable de M2
- agregado panel de temporada viva con record, partidos jugados sobre 82 y diferencial promedio
- agregado panel de resultados recientes alimentado por `season.patch`
- agregado inbox narrativo alimentado por `narrative.event`
- corregido flujo frontend para que eventos narrativos post-partido no se descarten despues del owner intro
- los eventos tecnicos quedan separados del inbox narrativo
- el panel de agentes core sigue visible junto al estado urbano
- al crear nueva partida se resetean tiempo, temporada, resultados, inbox, ciudad y agentes
- pruebas corridas: `npm run build --prefix frontend`

### M2.14 — Analytics basico — YA REALIZADO

Objetivo:

`analytics-service` ingesta series sin UI.

Incluye:

- resultados
- box scores
- ticket sales
- fan sentiment
- valor del suelo
- economia cambio

Done:

- datos quedan persistidos
- no se duplican ante reintentos
- queda base lista para charts futuros

Resultado 2026-05-24:

- `analytics-service` ahora conecta a Postgres/TimescaleDB usando `pgx`
- agregado paquete `internal/persistence`
- agregado schema propio para series analiticas
- consume `partido.terminado`
- persiste resultados y box scores por jugador
- consume `ciudad.economia_cambio`
- persiste fan sentiment, ticket sales y economia local
- consume `ciudad.suelo_actualizado`
- persiste valor del suelo por zona
- consume `agente.estado_cambio`
- persiste historial de estado de agentes para charts futuros
- handlers corren en goroutine con timeout corto para no bloquear NATS
- idempotencia por claves primarias de partido/evento/metrica
- agregada migracion `014_create_analytics_timeseries_tables.sql`
- pruebas corridas: `GOCACHE=/tmp/pulsecity-analytics-gocache go test ./...`

### M2.15 — Temporada completa — YA REALIZADO

Objetivo:

Cerrar M2 corriendo una temporada regular completa.

Done:

- 82 partidos simulados
- record final consistente
- ciudad tiene historial de cambios
- agentes acumularon estado
- narrativa post-partido existe
- pausa/reanudacion funciona
- reconexion frontend no rompe estado
- WebSocket sigue siendo delta-only

Resultado 2026-05-24:

- se corrio smoke end-to-end con servicios locales reales
- partida validada: `0ba7bc3d-f6ca-49db-9597-5196485b0e65`
- `team_schedule`: 82 partidos finalizados de 82
- `team_records`: 37-45, total 82
- `city_processed_matches`: 82
- `agent_processed_matches`: 82
- `narrative_events` post-partido: 82
- `analytics_match_results`: 82
- `analytics_player_box_scores`: 1640
- `analytics_city_metric_points`: 246
- `analytics_land_value_points`: 82
- `analytics_agent_state_points`: 1394
- bug encontrado: a x20 `tiempo.dia_avanzado` podia traer `days_processed > 1` y `team-service` solo procesaba la fecha final
- fix aplicado: `team-service` expande todas las fechas cubiertas por `days_processed` antes de buscar partidos
- agregado test `CoveredDayAdvancedEvents`
- pruebas corridas: `GOCACHE=/tmp/pulsecity-team-gocache go test ./...`, `make test-go`, `make test-rust`, `npm run build --prefix frontend`, `make build`

### M2.16 — DX local para probar la app — YA REALIZADO

Objetivo:

Hacer que levantar PulseCity en local sea simple para probar la app como usuario real, sin ritual manual de multiples terminales ni pasos fragiles.

Problema actual:

- `make up` levanta solo infraestructura: Postgres/TimescaleDB y NATS
- `make dev` hoy solo hace `docker compose up -d` y muestra logs de Docker
- los servicios de app deben levantarse manualmente con `make run-*`
- `agent-service` todavia requiere `GAME_ID` fijo para avanzar una partida concreta
- para probar una temporada completa desde frontend hay que crear/abrir partida, copiar el `game_id`, levantar `agent-service` con ese id y recien ahi controlar tiempo desde la UI

Done esperado:

- existe un comando claro para entorno local jugable, por ejemplo `make dev-app` o `make dev-full`
- queda documentado que servicios levanta y en que puertos
- el flujo de prueba desde browser queda reducido a pasos concretos
- se resuelve o mejora el acople manual de `agent-service` con `GAME_ID`
- no cambia la arquitectura de ownership ni la regla WebSocket delta-only

Notas:

- este mini milestone no bloquea el cierre funcional de M2
- conviene resolverlo antes de empezar M3 para que cada prueba jugable sea mas rapida y menos propensa a errores

Resultado 2026-05-25:

- agregado `make dev-app` para levantar infraestructura y servicios de app desde un solo comando
- documentado el flujo local jugable en `README.md`
- `agent-service` ahora puede correr sin `GAME_ID`
- en modo normal, `agent-service` actua como supervisor dinamico y carga/inicializa la simulacion cuando recibe `tiempo.sesion_iniciada`
- se mantiene compatibilidad con `GAME_ID=<id> make run-agent-service` para debug de una partida concreta
- el cambio no altera ownership: el tiempo sigue en `agent-service`, gateway solo publica eventos de sesion/control, WebSocket sigue enviando deltas
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`, `make -n dev-app`

## Edge cases a contemplar

### Tiempo y sesion

- Si no hay WebSocket activo, el tiempo no avanza.
- Al volver el jugador, no se compensan dias offline.
- Si se cierra la sesion en medio de acumulacion parcial de dia, no debe publicarse dia duplicado.
- Si cambia la velocidad durante una acumulacion, la nueva velocidad aplica desde el siguiente tick.
- Si la pausa se activa, no se procesan nuevos dias aunque haya sesion activa.

### Partidos

- Si un partido ya fue simulado, no se puede simular de nuevo por retry.
- Si `partido.terminado` llega dos veces, `team-service`, `city-service`, `agent-service`, `analytics-service` y `narrative-service` deben ser idempotentes.
- Si el calendario no tiene partido en un dia, no pasa nada deportivo.
- Si hay varios partidos el mismo dia en M2, se puede decidir no permitirlo para el equipo propio.
- Si falta input para `match-service`, el partido debe fallar de forma visible en logs, no inventar datos silenciosamente.

### Multi partida y multi cliente

- Todo evento debe incluir `game_id`.
- El frontend filtra por `game_id`.
- El gateway no debe mezclar deltas de partidas distintas.
- Dos clientes mirando la misma partida deben ver los mismos deltas.
- Dos partidas activas no deben compartir loop, calendario ni estado.

### Servicios caidos

- Si un consumidor esta caido, no debe escribir estado corrupto al volver.
- Si `narrative-service` cae, los partidos igual se simulan.
- Si `analytics-service` cae, el gameplay no se frena.
- Si `match-service` cae, el partido del dia no puede completarse y debe quedar pendiente o fallido de forma clara.

### WebSocket

- El backend no manda estado completo por WebSocket.
- Ante reconexion, el frontend debe rehidratar snapshot por REST o flujo equivalente y despues continuar con deltas.
- Los deltas deben ser chicos y especificos.

## Test plan

### Unit tests

- calculo de acumulacion de dias por velocidad
- pausa y reanudacion
- sesion activa/inactiva
- generacion de calendario de 82 partidos
- simulacion deterministica de partidos
- actualizacion de record
- reaccion de ciudad por resultado
- reaccion de agentes por resultado
- generacion templateada de narrativa post-partido

### Integration tests

Flujo minimo esperado:

```text
tiempo.dia_avanzado
team-service detecta partido
match-service simula
partido.terminado
team-service actualiza record
city-service actualiza ciudad
agent-service actualiza agentes
narrative-service genera evento
gateway reenvia delta
```

### Idempotencia

Probar que reintentar estos eventos no duplica estado:

- `tiempo.dia_avanzado`
- `partido.terminado`
- `ciudad.economia_cambio`
- `agente.estado_cambio`
- `narrativa.evento_generado`

### Build/checks

Comandos esperados al cierre de mini milestones:

```bash
go test ./...
cargo test
npm run build
```

Adaptados por servicio segun Makefile.

## Supuestos explicitos

- M2 no implementa playoffs.
- M2 no implementa Draft.
- M2 no implementa Free Agency.
- M2 no implementa trades.
- M2 no implementa salary cap real.
- M2 no implementa lesiones profundas.
- M2 no implementa LLM real.
- M2 no implementa ciudades rivales.
- M2 no implementa agentes de otros equipos.
- M2 no implementa graficos analiticos visibles.
- M2 no implementa city-builder completo.

## Filosofia de ejecucion

Cada sesion de trabajo dentro de M2 debe elegir un mini milestone chico y cerrarlo.

No avanzar en tareas gigantes tipo "hacer M2".

Formato recomendado por sesion:

```text
Objetivo de la sesion:
Mini milestone:
Archivos tocados:
Decision tomada:
Pruebas corridas:
Pendiente siguiente:
```

La prioridad de M2 es orden, no velocidad.

El mundo ya nace. Ahora tiene que empezar a latir.
