# Sesion 25 Mayo 2026 — Arranque M3 y fundacion de agentes vivos

## Objetivo de la sesion

Cerrar el ultimo pendiente ergonomico de M2 y arrancar `Milestone 3 — Los agentes viven` con mini milestones chicos, dejando base real de agentes, jugadores, relaciones y memoria del GM.

## Estado al iniciar

M2 ya estaba cerrado a nivel funcional:

- temporada regular completa de 82 partidos simulada
- ciudad reaccionando a resultados y rachas
- 5 agentes core reaccionando a partidos
- narrativa post-partido persistida
- analytics basico funcionando
- frontend mostrando HUD, calendario/resultados, ciudad, agentes e inbox

Quedaba un pendiente recomendado antes de M3:

```text
M2.16 — DX local para probar la app
```

## M2.16 — DX local para probar la app

Se agrego `make dev-app` para levantar el entorno jugable completo desde un solo comando.

Servicios levantados:

```bash
make dev-app
```

Puertos:

- frontend: `http://localhost:5173`
- gateway: `http://localhost:8080`
- Postgres/TimescaleDB: `localhost:5433`
- NATS: `localhost:4222`

Decision importante:

`agent-service` ya no necesita `GAME_ID` para el flujo normal. Ahora corre como supervisor dinamico: escucha `tiempo.sesion_iniciada`, carga o inicializa el estado de simulacion para ese `game_id`, y procesa pausa/velocidad desde ahi.

Se mantiene el modo anterior para debug puntual:

```bash
GAME_ID=<id-de-la-partida> make run-agent-service
```

Validacion manual:

Se probo desde browser con una partida nueva. El log confirmo:

- mapa generado
- temporada inicial creada con roster 15, rivales 30 y calendario 82
- `agent-service` cargo la simulacion automaticamente al abrir WebSocket
- controles de pausa/x1/x5/x20 llegaron correctamente
- al cerrar el browser se publico `tiempo.sesion_terminada`

## M3.1 — Catalogo completo de 30 agentes

Se agrego el catalogo canonico de 30 agentes individuales en `agent-service`.

Cada agente tiene:

- `agent_id`
- nombre visible
- categoria
- rol
- dominio
- estado emocional inicial
- confianza
- satisfaccion
- lealtad
- rendimiento de rol
- estado numerico inicial
- agenda inicial

Los 5 core de M2 quedan tambien integrados al nuevo schema:

- `owner`
- `head_coach`
- `cfo`
- `scouting_director`
- `sports_psychologist`

Persistencia:

- tabla `agent_individual_states`
- migracion `015_create_agent_individual_states.sql`
- seed idempotente al fundar/cargar partida

Decision tecnica:

M3.1 fue solo semilla. No se agregaron reacciones nuevas todavia.

## M3.2 — Roster como agentes

Los 15 jugadores propios dejaron de existir solo como registros contractuales de `team-service` y ahora tienen capa emocional en `agent-service`.

Tabla nueva:

```text
agent_player_states
```

Variables emocionales iniciales:

- `emotional_state`
- `satisfaction`
- `loyalty`
- `ego`
- `competitive_drive`
- `city_connection`

Ownership:

- `team-service` sigue siendo dueño de contrato, rating, posicion contractual y roster.
- `agent-service` lee `team_roster_players` solo para sembrar identidad compartida.
- `agent-service` escribe unicamente estado emocional en `agent_player_states`.

Flujo nuevo:

```text
partido.terminado
  -> agent-service lee box score
  -> actualiza capa emocional de jugadores propios
  -> publica roster.patch
  -> gateway reenvia delta
  -> frontend aplica RosterClientStates
```

Migracion:

```text
016_create_agent_player_states.sql
```

## M3.3 — Relaciones inter-agente

Se creo la primera red persistida de relaciones entre agentes.

Tablas nuevas:

```text
agent_relationships
agent_relationship_event_hashes
```

Relaciones sembradas:

- Head Coach ↔ Head of Analytics
- Head Coach ↔ Medico del Equipo
- Head Coach ↔ Director de Player Development
- Director de Scouting ↔ Head of Analytics
- Director de Marketing ↔ GM
- CFO ↔ GM
- Alcalde ↔ Owner
- Presidente Camara de Comercio ↔ Alcalde
- Sports Psychologist ↔ Head Coach
- Director de PR ↔ GM
- Prensa ↔ roster colectivo

Flujo nuevo:

```text
partido.terminado
  -> agent-service mueve relaciones relevantes
  -> publica agente.relacion_cambio
  -> gateway traduce a relations.patch
  -> frontend aplica RelationshipClientStates
```

Decision tecnica:

No se abrieron todavia decisiones GM complejas. Por eso solo `partido.terminado` mueve relaciones donde el resultado tiene causalidad directa.

Migracion:

```text
017_create_agent_relationships.sql
```

## M3.4 — Memoria de decisiones del GM

Se agrego la primera memoria append-only de decisiones del GM.

Tabla nueva:

```text
gm_decisions_log
```

Evento nuevo:

```text
decision.gm_registrada
```

Primer emisor real:

El `gateway` publica `decision.gm_registrada` cuando el GM responde el evento inicial del Owner.

Flujo:

```text
GM responde owner intro
  -> gateway persiste owner_intro_response
  -> gateway publica narrativa.respuesta_gm
  -> gateway publica decision.gm_registrada
  -> agent-service consume y persiste en gm_decisions_log
```

Decision tecnica:

El log no se expone directamente al frontend. Es memoria sistemica interna para que los agentes puedan construir contexto mas adelante.

Se agrego lectura interna:

```text
latest_gm_decisions(game_id, limit)
```

Esto queda listo para `M3.5 — Infra de chat con stub`.

Migracion:

```text
018_create_gm_decisions_log.sql
```

## Contratos nuevos documentados

Se agregaron o actualizaron contratos en:

- `shared/events/agents.md`
- `shared/events/websocket_deltas.md`
- `shared/events/decisions.md`

Eventos/deltas nuevos:

```text
roster.patch
agente.relacion_cambio
relations.patch
decision.gm_registrada
```

## Pruebas corridas

Durante la sesion se corrieron, segun mini milestone:

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
npm run build --prefix frontend
make build
```

Todos los checks pasaron.

## Archivos de progreso creados

- `docs/02_Progress/2026-05-25_M2_16_DX_Local.md`
- `docs/02_Progress/2026-05-25_M3_1_Catalogo_Completo_30_Agentes.md`
- `docs/02_Progress/2026-05-25_M3_2_Roster_Como_Agentes.md`
- `docs/02_Progress/2026-05-25_M3_3_Relaciones_Inter_Agente.md`
- `docs/02_Progress/2026-05-25_M3_4_Memoria_Decisiones_GM.md`

## Estado al cerrar

M3 queda iniciado con el Bloque A practicamente armado:

- agentes individuales sembrados
- jugadores como agentes emocionales
- relaciones canonicas persistidas
- memoria inicial de decisiones del GM

Todavia no hay:

- chat con agentes
- LLM real
- trades
- lesiones
- salary cap acotado
- playoffs
- vista franquicia completa

## Pendiente siguiente

Siguiente mini milestone recomendado:

```text
M3.5 — Infra de chat con stub
```

Objetivo chico para la proxima sesion:

- endpoint del gateway para iniciar/continuar conversacion con un agente
- persistencia de historial de chat
- respuesta stub deterministica
- contexto armado con estado del agente + relaciones + ultimas decisiones del GM
- delta `chat.message` por WebSocket

## Cierre

La sesion cerro la transicion real entre M2 y M3.

M2 dejo al sistema latiendo. M3 ya empezo a darle memoria, relaciones y vida interna a los personajes.
