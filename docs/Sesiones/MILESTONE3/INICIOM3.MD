# Milestone 3 — Inicio

## Los agentes viven

`Milestone 3` arranca con un cambio de corazon: el sistema deja de tener "5 agentes que reaccionan a partidos" y pasa a ser un mundo poblado de personajes con memoria, relaciones, lenguaje propio y consecuencias diferidas.

La cadena base del milestone es:

```text
agentes vivos -> relaciones evolucionan -> decisiones del GM tienen memoria
   -> trades / playoffs / cierre con presion real
   -> narrativa rule-based ampliada + chat libre con LLM
   -> una temporada completa con desenlace y veredicto del Owner
```

El objetivo no es todavia abrir Draft ni Free Agency. El objetivo es que los 30 agentes del canon (mas el roster como agentes) existan con estado real, que sus relaciones se muevan ante eventos, que el jugador pueda hablar con cualquiera de ellos via chat real, y que la temporada se cierre con Trade Deadline, Playoffs y evaluacion formal del Owner.

## Criterio de done de M3

Una partida iniciada en `Milestone 2` puede vivir una temporada completa con:

- los 30 agentes individuales del canon activos y con estado persistido
- los 15 jugadores del roster como agentes con estado emocional/relacional
- una tabla viva de relaciones entre agentes que evoluciona con eventos
- una memoria registrada de decisiones del GM, referenciable por agentes
- chat libre con cada agente, anclado a su dominio, via LLM real
- modelo de partidos v2 donde ratings + minutos del Coach + estado emocional + cansancio mueven el resultado
- lesiones operativas con flujo Medico ↔ GM ↔ Coach
- salary cap acotado con cap base, luxury tax line y salarios reales por jugador
- trades operativos durante la temporada regular con GMs rivales reaccionando con perfil propio
- Trade Deadline funcional con pausa automatica y UI expandida
- Playoffs con bracket de 16 equipos, series al mejor de 7 y 4 rondas
- cierre formal con evaluacion del Owner: confianza renovada, advertencia, ultimatum o despido
- frontend con Vista Franquicia minima jugable: estadio por zonas, toolbar, buscador de agentes, paneles roster/finanzas/calendario/bandeja y bracket de playoffs
- WebSocket sigue enviando deltas, no estado completo

## Decisiones cerradas

### Estrategia general

M3 se construye sobre los servicios existentes de M2. No se crean servicios nuevos en el primer corte. La regla sigue siendo: ownership estricto, comunicacion via NATS, deltas por WebSocket.

`agent-service` crece para alojar a todos los agentes individuales (no solo los 5 core), la capa emocional de los jugadores del roster, la tabla de relaciones inter-agente y el log de decisiones del GM. `narrative-service` crece para soportar dos canales: eventos sistemicos templateados (igual que M2) y chat libre con LLM.

La pregunta de si conviene separar un `chat-service` propio queda diferida. Se evalua despues de cerrar M3.6 con LLM real. Si la complejidad se vuelve incomoda dentro de `narrative-service`, se separa antes de M3.7. Si no, queda dentro hasta M4.

### LLM hibrido

M3 conecta LLM real, pero solo en el canal de chat directo agente ↔ jugador.

- chat libre con agentes: LLM real
- narrativa post-partido, post-trade, post-lesion, post-evento critico: rule-based / templateada
- eventos automaticos del sistema (resumenes, alertas, reportes): rule-based / templateada

Razon:

- el chat libre es la experiencia central del milestone, "el mundo se siente habitado, no administrado". Sin LLM se siente plastico.
- meter LLM en cada evento automatico explota costo sin retorno proporcional. La narrativa templateada de M2 ya funciona bien con estado real.
- mantener LLM acotado al chat permite controlar costo, prompt engineering y manejo de fallos en un solo canal.

Detalle operativo:

- el chat solo se activa cuando el jugador inicia conversacion
- nunca hay LLM corriendo en background
- cap de tokens por turno
- cap de turnos por sesion de chat
- fallback templateado si el LLM falla, en personaje, sin romper la sesion
- cache de respuestas frecuentes se difiere a M4

El proveedor concreto (Claude API vs GPT-4o mini) se decide al iniciar `M3.6`. M3.5 implementa la infraestructura con stub determinista; M3.6 conecta el proveedor real.

### Alcance de la liga

M3 mantiene los rivales abstractos de M2 a nivel de simulacion de partidos. No se generan los 30 rosters rivales completos.

Lo que si entra:

- los 30 GMs rivales pasan a ser agentes con perfil real (estilo de negociacion, urgencia actual, necesidades de roster, memoria de relaciones con el GM jugador)
- los GMs rivales reaccionan en trades segun su perfil
- la red de relaciones con GMs rivales acumula historia

Lo que no entra:

- roster ficticio generado para cada uno de los 30 equipos rivales
- salary cap real para los rivales
- IA de roster de los rivales
- trades nominales (es decir, "te doy a este jugador especifico de mi roster por este especifico del tuyo")

Cuando se acepta un trade, el jugador recibido se genera con ratings y perfil calibrados a la necesidad declarada del rival. No tiene historia previa en otro roster real.

Razon: el canon de M3 se llama "los agentes viven", no "la liga vive". Generar 30 rosters completos + IA de roster por equipo es trabajo del milestone final. Para que los trades y la negociacion se sientan reales alcanza con un GM rival con personalidad y necesidades.

### Ciclo jugable cubierto por M3

M3 cubre, en orden de aparicion durante una temporada:

- Pretemporada acotada (sin UI propia, solo como contexto)
- Temporada Regular con gestion ongoing del roster, carga, lesiones
- Trades durante la temporada regular
- Trade Deadline con pausa automatica + UI expandida
- Playoffs (bracket completo, 16 equipos, mejor de 7, 4 rondas)
- Cierre de Temporada con evaluacion del Owner

M3 no cubre:

- Draft (queda para `M3.5`)
- Free Agency (queda para `M3.5`)
- Pool global de jugadores libres
- Renovaciones de contrato profundas (M3 cubre renovaciones in-line basicas de jugadores propios que estan por expirar contrato; las negociaciones complejas de tipos de contrato quedan para `M3.5`)

Razon: Trade Deadline, Playoffs y Cierre son la prueba de fuego de los agentes vivos. Ahi explotan emociones, relaciones, presion del Owner. Si no metemos eso, los agentes existen pero no se usan. Draft y Free Agency son fases offseason grandes con UI propia y mezcladas en M3 lo vuelven inabarcable.

### Subsistemas profundos incluidos

M3 incluye, ademas de los agentes:

- modelo de partidos v2 (ratings, minutos del Coach, estado emocional y cansancio mueven el resultado)
- lesiones operativas con flujo Medico ↔ GM ↔ Coach
- salary cap acotado: cap base + luxury tax line + salarios por jugador

M3 no incluye:

- player development (curvas de progresion/declive por edad, mejora de ratings via Director de Player Development) → queda para `M4`
- exceptions del salary cap (mid-level, bi-annual, etc.), player options, team options, no-trade clauses → quedan para `M3.5`
- edificios especiales nuevos ni alcalde como agente con agenda → quedan para `M4`

Razon:

- sin modelo v2 los ratings son cosmeticos y las decisiones del Coach sobre minutos no mueven el resultado → no hay causalidad real, contradice el espiritu del milestone
- sin lesiones el Medico es un agente vacio
- sin cap acotado los trades son juguete
- player development se ve "en temporadas, no en semanas" → necesitaria simular varias temporadas para verlo funcionar, incompatible con el alcance de M3

### Ownership

`team-service` sigue siendo dueño de:

- contrato (salario, anios restantes, tipo basico)
- posicion en roster (activo / inactivo / lesionado)
- salary cap del equipo
- box scores acumulados
- calendario, resultados, record
- bracket de playoffs (porque el bracket se deriva del record)

`agent-service` es dueño de:

- estado real de los 30 agentes individuales (Basketball Ops + Business Ops + Ciudad + Prensa)
- estado emocional/relacional de los 15 jugadores del roster (capa agente, no contractual)
- 30 perfiles de GMs rivales
- tabla `agent_relationships`
- log `gm_decisions_log`
- estado de simulacion (tiempo, sesion, pausa, velocidad) — sin cambios respecto a M2

`city-service` sigue como en M2.

`narrative-service` es dueño de:

- eventos narrativos sistemicos templateados (post-partido, post-trade, post-lesion, post-evento critico, post-playoffs, post-cierre)
- historial de chat libre por `(game_id, agent_id)`
- llamadas al LLM y manejo de fallos

`analytics-service` sigue ingestando, ampliando series a relaciones, lesiones, trades, playoffs.

`gateway` no calcula ni decide nada de gameplay. Sigue siendo punto de entrada y salida, mantiene WebSocket, publica eventos, reenvia deltas.

### Persistencia

Misma base Postgres/Timescale local, ownership logico por tabla. Cada servicio agrega sus migraciones propias.

Tablas nuevas previstas en M3:

- `agent_individual_states` (los 30 + roster como agentes)
- `agent_relationships`
- `gm_decisions_log`
- `agent_chat_history`
- `rival_gms`
- `team_injuries`
- `team_salary_cap`
- `team_trades`
- `team_playoffs_bracket`
- `team_playoffs_series`
- `season_close_evaluations`

Reglas que no cambian:

- ningun servicio escribe tablas de otro
- todo evento incluye `event_id`, `game_id`, `occurred_at`, `schema_version`
- idempotencia por claves naturales de evento

## Interfaces y eventos principales

### Agentes y relaciones

```text
agente.estado_cambio        (ampliado: cubre los 30 agentes + 15 jugadores)
agente.relacion_cambio
agente.evento_critico
agente.consulta_iniciada    (jugador inicia chat con un agente)
agente.respuesta_generada   (chat: respuesta lista, generada por LLM o stub)
```

### Decisiones y memoria

```text
decision.gm_registrada
decision.referenciada_por_agente
```

`decision.gm_registrada` debe incluir:

- `game_id`
- `decision_id`
- `kind` (trade_propuesta, trade_aceptada, ignorar_medico, rotacion_cambiada, etc.)
- `payload` minimo necesario para reconstruir contexto
- `simulated_date`
- `agents_affected` (lista de `agent_id`)

### Roster, contratos, finanzas

```text
jugador.solicita_traspaso
jugador.lesionado
jugador.recuperado
jugador.contrato_extendido
salary_cap.calculado
finanzas.actualizada
```

### Trades

```text
trade.propuesta_enviada
trade.contraoferta
trade.rechazada
trade.aceptada
trade.deadline_activado
trade.deadline_cerrado
```

### Playoffs

```text
playoffs.bracket_generado
playoffs.serie_iniciada
playoffs.partido_resultado
playoffs.serie_terminada
playoffs.ronda_avanzada
playoffs.campeon_definido
```

### Cierre de temporada

```text
temporada.cerrada
owner.evaluacion_iniciada
owner.veredicto_emitido      (renovado | advertencia | ultimatum | despido)
staff.renovacion_decidida
```

### Narrativa rule-based ampliada

```text
narrativa.evento_generado    (mismo evento de M2, ahora cubre trades, lesiones, conflictos, playoffs, cierre, etc.)
```

### Deltas WebSocket nuevos

```text
roster.patch
relations.patch
inbox.patch                  (ampliado para chat + eventos)
finance.patch
trade.patch
playoffs.patch
season_close.patch
chat.message                 (mensaje nuevo en una conversacion abierta)
```

Regla no negociable: WebSocket sigue enviando deltas, no estado completo. Reconexion rehidrata por REST y continua con deltas.

## Mini milestones propuestos

Orden recomendado A → B → C → D → E → F → G. Bloques pueden intercalarse si conviene mantener motivacion, pero las dependencias duras quedan marcadas.

### Bloque A — Fundacion de agentes vivos

#### M3.0 — Documento rector M3 — YA REALIZADO

Objetivo:

Dejar asentado en `docs/Sesiones/MILESTONE3/INICIOM3.MD` el alcance, decisiones, mini milestones, interfaces, edge cases y criterios de done de M3.

Done:

- este documento existe
- las decisiones principales estan escritas
- el milestone queda listo para ejecutarse por partes

#### M3.1 — Catalogo completo de 30 agentes — YA REALIZADO

Objetivo:

Sembrar estado inicial real para los agentes del canon que todavia no existen.

Incluye:

- Basketball Ops: President of Basketball Operations, AGMs, Director de Player Personnel, Head of Analytics, Asistentes del Coach, Director de Player Development, Medico, Fisioterapeuta/S&C Coach, Video Coordinator, International Scout
- Business Ops: CEO, Director de Marketing & Brand, Director de Ticket Sales, Director de Corporate Partnerships, Director de PR & Communications, Director de Arena Operations, Legal Counsel
- Ciudad: Alcalde, Jefe de Policia, Presidente de la Camara de Comercio, Director de Urbanismo
- Prensa (agente colectivo)

Done:

- cada agente tiene estado inicial persistido segun variables canon
- existe migracion `agent_individual_states`
- `agent-service` carga e inicializa esos agentes al fundar partida
- todavia no reaccionan a eventos: es solo semilla
- los 5 core de M2 se mantienen sin cambios y se integran al nuevo schema

Resultado 2026-05-25:

- agregado catalogo canonico de 30 agentes individuales en `agent-service`
- los 5 core de M2 (`owner`, `head_coach`, `cfo`, `scouting_director`, `sports_psychologist`) quedan tambien sembrados en `agent_individual_states`
- cada agente tiene categoria, rol, dominio, estado emocional, confianza, satisfaccion, lealtad, rendimiento de rol, estado numerico inicial y agenda inicial
- agregada migracion `015_create_agent_individual_states.sql`
- `agent-service` asegura el schema `agent_individual_states`
- `agent-service` escucha `mapa.generacion_iniciada` y siembra agentes de forma idempotente al fundar partida
- el supervisor dinamico tambien asegura agentes al cargar/inicializar una simulacion por `game_id`
- no se agregaron reacciones nuevas: los 30 agentes quedan sembrados, y las reacciones siguen limitadas a los 5 core de M2
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`

#### M3.2 — Roster como agentes — YA REALIZADO

Objetivo:

Cada uno de los 15 jugadores propios deja de ser solo un registro contractual y se convierte en agente con capa emocional.

Done:

- `agent-service` agrega capa por jugador con estado emocional, satisfaccion, lealtad, ego, mentalidad competitiva, vinculos con la ciudad
- `player_id` compartido con `team-service`
- ningun servicio escribe en tablas del otro
- estados iniciales sembrados al crear partida
- los jugadores reaccionan a `partido.terminado` aplicando deltas sobre su capa emocional segun su actuacion y el resultado
- delta `roster.patch` lleva estado emocional resumido al frontend

Resultado 2026-05-25:

- agregada tabla propia `agent_player_states` para la capa emocional de jugadores
- agregada migracion `016_create_agent_player_states.sql`
- `agent-service` lee `team_roster_players` solo para sembrar `player_id`, nombre, posicion y rating inicial
- `agent-service` escribe solo su tabla emocional, respetando ownership de `team-service`
- al fundar/cargar partida se siembran los 15 jugadores propios como agentes emocionales
- `partido.terminado` ahora incluye box score en el contrato interno de `agent-service`
- los jugadores reaccionan a box score propio ajustando `emotional_state`, `satisfaction`, `loyalty`, `ego`, `competitive_drive` y `city_connection`
- `agent-service` publica `roster.patch` con resumen emocional por jugador actualizado
- `gateway` escucha `roster.patch` y lo reenvia por WebSocket
- frontend tipa y aplica `roster.patch` en `RosterClientStates`
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `npm run build --prefix frontend`, `make build`

#### M3.3 — Tabla de relaciones inter-agente — YA REALIZADO

Objetivo:

Modelar la red de relaciones del canon.

Incluye:

- tabla `agent_relationships` con `(agente_a_id, agente_b_id, confianza, ultimo_evento, tendencia, historial_corto)`
- semillas iniciales con tensiones predefinidas del canon (Coach ↔ Medico, Scouting ↔ Analytics, Director de Marketing ↔ GM, etc.)
- eventos `agente.relacion_cambio` cuando un delta supera umbral
- delta `relations.patch` al frontend

Done:

- al crear partida quedan sembradas las relaciones canonicas
- eventos sistemicos relevantes mueven relaciones (ej: ignorar al Medico empeora Medico ↔ GM)
- la idempotencia se controla por hash del evento que disparo el cambio

Resultado 2026-05-25:

- agregada tabla `agent_relationships`
- agregada tabla `agent_relationship_event_hashes` para idempotencia por `(game_id, relationship_key, source_event_id)`
- agregada migracion `017_create_agent_relationships.sql`
- sembradas relaciones canonicas iniciales: Coach ↔ Analytics, Coach ↔ Medico, Coach ↔ Player Development, Scouting ↔ Analytics, Marketing ↔ GM, CFO ↔ GM, Alcalde ↔ Owner, Camara de Comercio ↔ Alcalde, Sports Psychologist ↔ Coach, PR ↔ GM y Prensa ↔ roster colectivo
- `agent-service` asegura relaciones al fundar/cargar partida
- `partido.terminado` mueve un subconjunto de relaciones segun resultado, localia y margen, sin abrir todavia mecanicas de decisiones GM
- `agent-service` publica `agente.relacion_cambio`
- `gateway` traduce `agente.relacion_cambio` a `relations.patch`
- frontend tipa y aplica `relations.patch` en `RelationshipClientStates`
- contratos documentados en `shared/events/agents.md` y `shared/events/websocket_deltas.md`
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `npm run build --prefix frontend`, `make build`

#### M3.4 — Memoria de decisiones del GM — YA REALIZADO

Objetivo:

El sistema empieza a recordar lo que el GM decide para que los agentes lo referencien despues.

Done:

- existe `gm_decisions_log` append-only
- cada decision relevante (trade, ignorar alerta, cambio de rotacion, renovar staff, etc.) publica `decision.gm_registrada`
- `agent-service` consume y persiste
- los agentes pueden consultar las ultimas N decisiones del GM para armar contexto
- el log no se expone directamente al frontend, solo via reads agregadas del gateway

Resultado 2026-05-25:

- agregada tabla append-only `gm_decisions_log`
- agregada migracion `018_create_gm_decisions_log.sql`
- documentado contrato `decision.gm_registrada` en `shared/events/decisions.md`
- `gateway` publica `decision.gm_registrada` cuando el GM responde el evento inicial del Owner
- `agent-service` escucha `decision.gm_registrada` y persiste la decision de forma idempotente
- idempotencia por `(game_id, decision_id)` y `event_id`
- `agent-service` expone lectura interna `latest_gm_decisions(game_id, limit)` para que M3.5 pueda armar contexto de chat
- el log no se expone al frontend
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `make build`

### Bloque B — Chat directo con LLM

#### M3.5 — Infra de chat con stub — YA REALIZADO

Objetivo:

Dejar todo el camino del chat funcionando end-to-end, sin LLM real.

Incluye:

- endpoint del `gateway` para iniciar/continuar conversacion con un agente
- ruteo a `narrative-service` (o `agent-service`, decision interna a tomar en este corte)
- ensamblado de contexto del agente: estado actual + relacion con GM + ultimas decisiones relevantes
- respuesta stub determinista por agente
- persistencia de historial por `(game_id, agent_id, conversation_id)`
- frontend minimo: panel de chat que se abre desde el perfil del agente
- delta `chat.message` por WebSocket

Done:

- el chat funciona end-to-end con respuesta stub
- existe migracion `agent_chat_history`
- el contexto enviado al "modelo" ya es el real, lo unico que falta es el modelo

Resultado 2026-05-29:

- agregado endpoint `POST /api/v1/games/{gameID}/agent-chat` en `gateway`
- `gateway` valida ownership de partida, normaliza `agent_id/message/conversation_id` y publica `agente.consulta_iniciada`
- `narrative-service` consume `agente.consulta_iniciada`
- `narrative-service` ensambla contexto real leyendo `agent_individual_states` o `agent_player_states`, relacion con GM y ultimas decisiones de `gm_decisions_log`
- agregado stub deterministico de respuesta por agente, con metadata `generation=stub`
- agregada tabla/migracion `agent_chat_history`
- `narrative-service` persiste mensaje del GM y respuesta del agente por `(game_id, agent_id, conversation_id)`
- `narrative-service` publica delta `chat.message` con subject `agente.respuesta_generada`
- `gateway` escucha `chat.message` y lo reenvia por WebSocket
- frontend agrega panel minimo de chat en `CeremonyPage` desde agentes core
- frontend muestra mensaje del GM de forma optimista y aplica respuestas `chat.message`
- contratos documentados en `shared/events/agents.md` y `shared/events/websocket_deltas.md`
- decision tomada: el chat vive en `narrative-service`; `agent-service` sigue siendo dueño del estado, relaciones y memoria, leidos como contexto
- limite del corte: no se agrego REST de rehidratacion de historial; el historial queda persistido para exponerlo cuando M3.5 se refine o al iniciar M3.6
- pruebas corridas: `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...`, `npm run build --prefix frontend`, `make build`

#### M3.6 — LLM real conectado — YA REALIZADO

Objetivo:

Reemplazar el stub por LLM real, anclado a dominio.

Incluye:

- decision final del proveedor (Claude API vs GPT-4o mini vs Azure vs Gemini vs Opciones de modelos mas baratas) tomada en este corte y documentada
- cliente del proveedor en el servicio responsable
- prompt template por agente con anclaje de dominio (el CFO no opina del roster, el Director de Urbanismo no opina de partidos, etc.)
- cap de tokens por turno
- cap de turnos por sesion de chat
- manejo de fallos: fallback templateado en personaje
- secrets gestionados via env / k8s secrets en su momento

Done:

- chat libre real con cualquier agente
- si el LLM falla, el chat no se rompe
- el costo por sesion esta acotado y se puede estimar
- los anclajes de dominio funcionan: preguntar al CFO sobre el roster devuelve "no es mi area" en una linea

Avance parcial 2026-05-29 — M3.6a:

- `narrative-service` ya no llama al stub de forma directa; ahora usa una interfaz `ChatResponder`
- agregado `ChatRuntimeConfig` configurable por env:
  - `LLM_PROVIDER` (default `stub`)
  - `LLM_MAX_PROMPT_CHARS` (default `6000`)
  - `LLM_MAX_RESPONSE_CHARS` (default `900`)
  - `LLM_MAX_TURNS_PER_CONVERSATION` (default `12`)
- agregado prompt canonico de chat con contexto real del agente, relacion con GM y ultimas decisiones
- el prompt incluye guardrail de dominio: si el GM pregunta fuera del area, el agente debe redirigir
- agregado limite de turnos por conversacion; al superarlo se persiste/publica respuesta `turn_limit`
- si el provider configurado no esta implementado o falla, se usa fallback templateado sin romper la conversacion
- el provider real queda pendiente para completar M3.6; `stub` sigue siendo el unico provider implementado
- pruebas corridas: `GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...`

Avance parcial 2026-06-02 — M3.6b:

- decision de proveedor para exploracion: `OpenRouter`
- modelo default inicial: `google/gemini-2.5-flash`
- `stub` sigue siendo el default local si no se configura `LLM_PROVIDER`
- agregado provider `openrouter` detras de `ChatResponder`
- agregado cliente HTTP directo a `POST /api/v1/chat/completions`, sin SDK nuevo
- secrets y configuracion via env:
  - `OPENROUTER_API_KEY`
  - `OPENROUTER_BASE_URL` (default `https://openrouter.ai/api/v1`)
  - `OPENROUTER_APP_URL` (opcional)
  - `OPENROUTER_APP_TITLE` (default `PulseCity`)
  - `LLM_MODEL` (default `google/gemini-2.5-flash`)
  - `LLM_MAX_COMPLETION_TOKENS` (default `300`)
  - `LLM_REQUEST_TIMEOUT_SECONDS` (default `8`)
- se mantiene cap de caracteres de prompt/respuesta y cap de turnos ya existentes
- si falta API key, OpenRouter falla o responde vacio, el flujo superior conserva fallback templateado
- metadata de respuesta registra `generation=llm`, `provider=openrouter`, `model` y tokens reportados por el provider
- pruebas unitarias cubren request, headers, metadata y error de provider sin abrir sockets locales
- pruebas corridas: `GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...`

Resultado 2026-06-02:

- smoke real ejecutado desde la app local con `.env.local` y `OPENROUTER_API_KEY`
- el chat con el CFO respondio con texto natural del agente, sin prefijo `[stub M3.6]` y sin fallback templateado
- `Makefile` carga `.env.local` automaticamente si existe, por lo que `make dev-app` propaga `LLM_PROVIDER=openrouter` y la key al `narrative-service`
- `M3.6` queda cerrado funcionalmente; el provider puede cambiarse por env sin tocar codigo

#### M3.7a — Reset UX de Ceremony / Command Center — YA REALIZADO

Objetivo:

Reordenar la pantalla jugable inicial antes de agregar el buscador completo de agentes.

Resultado 2026-06-02:

- `CeremonyPage` dejo de ser un archivo monolitico con toda la UI embebida
- se agrego carpeta `frontend/src/features/new-game/components/ceremony/`
- componentes extraidos:
  - `CeremonyTopbar`
  - `CeremonyMapPanel`
  - `CeremonySidePanel`
  - `AgentDirectoryPanel`
  - `InboxPanel`
  - `SeasonPanel`
  - `SystemPanel`
- el mapa pasa a ser el protagonista del primer viewport en desktop
- el panel derecho deja de apilar tiempo, temporada, inbox, agentes, chat, pipeline y eventos uno debajo del otro
- se agregan tabs internas: `Agentes`, `Inbox`, `Temporada`, `Sistema`
- el chat vive dentro de la tab `Agentes`, asociado al agente seleccionado
- se reemplazo el badge viejo `Stub M3.5` por `LLM real`
- `Pipeline` y `Eventos recientes` bajan a la tab `Sistema`, con jerarquia de debug
- responsive basico: en pantallas chicas el mapa queda primero y el panel baja debajo
- prueba corrida: `npm run build --prefix frontend`

Limite:

- el buscador real y el directorio completo se cerraron en `M3.7`

#### M3.7 — Buscador de agentes y vista de directorio — YA REALIZADO

Objetivo:

Hacer que llegar a cualquier agente sea trivial.

Done:

- buscador siempre accesible desde la Vista Franquicia
- directorio de los 30 agentes + roster con estado resumido
- desde cualquier perfil se puede abrir chat
- atajos contextuales por zona (placeholder mientras no exista vista por zonas de M3.21)

Resultado 2026-06-02:

- la tab `Agentes` incorpora buscador por nombre, rol, dominio o `agent_id`
- se agregaron filtros por categoria:
  - Todos
  - Basketball
  - Business
  - Ciudad
  - Roster
- el frontend contiene el catalogo canonico de agentes individuales sembrado en `agent-service`
- el directorio lista los agentes de Basketball Ops, Business Ops, Ciudad y Prensa
- el roster aparece como agentes cuando llegan `roster.patch` del backend
- seleccionar cualquier agente actualiza el foco y abre chat contra ese `agent_id`
- el chat conserva el flujo LLM real de `M3.6`
- los agentes no core usan contexto real desde `narrative-service`; si no tienen estado visual en frontend, igual pueden chatear por `agent_id`
- limite conocido: los jugadores del roster hoy se muestran con `player_id`, porque el delta emocional no trae nombre/posicion contractual al frontend
- prueba corrida: `npm run build --prefix frontend`

### Bloque C — Subsistemas profundos

#### M3.8 — Modelo de partidos v2 — YA REALIZADO

Objetivo:

Hacer que las decisiones reales muevan el resultado.

Incluye:

- ratings individuales reales por jugador
- minutos asignados por el Coach (derivados de su sistema y filosofia)
- estado emocional y cansancio acumulado modulan rendimiento
- contexto tactico simple del Coach (sistema, flexibilidad)
- mantiene determinismo por seed: mismo input + mismo seed = mismo output
- tests de regresion contra el simulador M2 (no es retrocompatible, pero se documenta el cambio)

Done:

- decisiones de roster impactan visiblemente en el resultado a lo largo de la temporada
- mismos seeds producen mismos resultados
- box score y momentos clave siguen siendo el output

Avance parcial 2026-06-06 — M3.8a — Base de rotacion y tactica — YA REALIZADO:

- `partido.programado` ahora puede incluir `home_tactics` y `away_tactics` con `system`, `rotation_preference` y `flexibility`
- `MatchPlayer` ahora puede incluir `expected_minutes`
- `team-service` publica minutos esperados deterministas para la rotacion propia y rivales abstractos
- `team-service` publica contexto tactico inicial `balanced` para ambos equipos
- `match-service` mantiene compatibilidad con payloads M2 sin tactica ni minutos explicitos
- `match-service` deriva una rotacion deterministica si faltan `expected_minutes`
- la fuerza del equipo ahora pondera rating, fatiga, estado emocional y stamina por minutos esperados
- el box score reparte puntos y minutos segun rol/minutos esperados, no solo por ranking de rating
- los sistemas tacticos ajustan pace y pequenos bonuses ofensivos/defensivos
- pruebas nuevas cubren compatibilidad M2, determinismo, cambio visible al modificar rotacion y efecto de minutos sobre box score
- pruebas corridas: `cargo test --manifest-path services/match-service/Cargo.toml`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`

Resultado 2026-07-14:

- `team-service` agrega `MatchPreparation` para preparar el payload del partido con estado vivo antes de publicar `partido.programado`
- `team-service` deriva tactica del Head Coach segun rival y presion de temporada:
  - rivales de alto pace/ofensiva fuerzan `defensive_grind`
  - defensas rivales vulnerables habilitan `pace_and_space`
  - mala racha usa rotacion `top_heavy`
  - ventaja amplia de record usa rotacion `deep`
- `team-service` calcula fatiga desde minutos recientes persistidos en `team_player_box_scores`
- `team-service` consume `roster.patch` publicado por `agent-service` y mantiene una proyeccion local minima en `team_player_match_states`
- `emotional_state` de jugadores propios entra al payload de partido desde esa proyeccion, sin leer tablas de `agent-service`
- `match-service` sigue stateless: solo consume el payload completo y mantiene determinismo por seed
- contratos documentados en `shared/events/matches.md`
- nota de progreso creada en `docs/02_Progress/2026-07-14_M3_8_Modelo_Partidos_V2.md`
- pruebas corridas: `cargo test --manifest-path services/match-service/Cargo.toml`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`

Limite:

- la filosofia del Head Coach queda modelada como politica deterministica inicial en `team-service`; cuando `agent-service` exponga filosofia editable del Coach, esta politica se reemplaza por esa fuente
- lesiones y riesgo medico quedan para `M3.9`

#### M3.9 — Carga, lesiones y flujo medico — YA REALIZADO

Objetivo:

Hacer que el Medico y el S&C Coach sean agentes utiles.

Incluye:

- carga acumulada por jugador segun minutos jugados
- riesgo probabilistico de lesion creciente con carga
- evento `jugador.lesionado` con severidad y dias estimados
- evento `jugador.recuperado` cuando termina la baja
- flujo: alerta del Medico → opciones del GM → decision → consecuencias
- riesgo de reagravacion si el GM presiona alta antes de tiempo
- delta `roster.patch` lleva estado de disponibilidad

Done:

- las lesiones ocurren de forma sistemica, no manualmente
- ignorar al Medico cuesta caro
- la relacion Medico ↔ GM se mueve segun decisiones

Avance parcial 2026-07-14 — M3.9a — Lesiones sistemicas — YA REALIZADO:

- `team-service` agrega tabla `team_injuries`
- `team-service` evalua riesgo de lesion al aplicar `partido.terminado`
- el riesgo usa carga reciente, minutos del partido y estado emocional negativo
- la evaluacion es deterministica por `game_id`, `match_id` y `player_id`
- si hay lesion, el jugador pasa a `roster_status = injured`
- `team-service` publica `jugador.lesionado` con severidad, dias estimados, fecha esperada de recuperacion y carga usada
- en cada `tiempo.dia_avanzado`, `team-service` recupera lesiones vencidas antes de despachar partidos del dia
- si hay recuperacion, el jugador vuelve a `roster_status = active`
- `team-service` publica `jugador.recuperado`
- `gateway` escucha `jugador.lesionado` y `jugador.recuperado`
- `gateway` transforma ambos eventos en `roster.patch` delta-only con disponibilidad del jugador
- frontend tipa campos opcionales de disponibilidad: `availability`, `injury_id`, `severity`, `expected_recovery_date`, `estimated_days_out`, `availability_changed_on`
- contratos documentados en `shared/events/players.md` y `shared/events/websocket_deltas.md`
- nota de progreso creada en `docs/02_Progress/2026-07-14_M3_9a_Lesiones_Sistemicas.md`
- pruebas corridas: `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `npm run build --prefix frontend`

Avance parcial 2026-07-14 — M3.9b — Decisiones medicas del GM — YA REALIZADO:

- `gateway` agrega endpoint `POST /api/v1/games/{gameID}/medical-decisions`
- el endpoint publica `decision.gm_registrada` con `kind = medical_decision`
- opciones validas:
  - `rest`
  - `reduce_minutes`
  - `ignore_doctor`
  - `force_return`
- `team-service` consume `decision.gm_registrada`
- `team-service` registra la decision medica en `team_injuries`
- si la decision es `force_return`, el jugador vuelve a `roster_status = active` y se publica `jugador.recuperado`
- si un jugador con alta forzada juega antes de `expected_recovery_date`, `team-service` genera una nueva lesion con `reason = forced_return_reaggravation`
- `agent-service` agrega relacion canonica GM ↔ Medico
- `agent-service` mueve relaciones GM ↔ Medico y Head Coach ↔ Medico ante decisiones medicas
- `agent-service` publica `agente.relacion_cambio` y `gateway` lo convierte en `relations.patch`
- contratos actualizados en `shared/events/decisions.md` y `shared/events/players.md`
- nota de progreso creada en `docs/02_Progress/2026-07-14_M3_9b_Decisiones_Medicas_GM.md`
- pruebas corridas: `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `cargo test --manifest-path services/agent-service/Cargo.toml`

Limite:

- no se agrego UI nueva para responder decisiones medicas; el endpoint queda listo para conectarse desde inbox/panel roster en `M3.22` / `M3.23`

#### M3.10 — Salary cap acotado — YA REALIZADO

Objetivo:

Hacer que el dinero importe lo suficiente para que los trades tengan sentido.

Incluye:

- cap base + luxury tax line por temporada
- salarios reales por jugador derivados de su rating, edad y arquetipo
- recalculo del cap cuando hay altas/bajas
- alerta del CFO cuando se acerca a la linea roja
- evento `salary_cap.calculado` y `finanzas.actualizada`
- panel financiero basico en el frontend

Done:

- el cap se calcula y muestra en frontend
- queda disponible la base de validacion de cap que `M3.12` usara antes de cerrar trades
- exceder la linea roja del luxury tax mueve la paciencia del Owner

Resultado 2026-07-14:

- `team-service` calcula salary cap con cap base, luxury tax line, salarios por jugador y estado `under_cap` / `over_cap` / `luxury_tax`
- los salarios iniciales del roster ahora se derivan de rating, rol implicito por tier y posicion, con montos NBA-like suficientes para que el cap importe
- `team-service` persiste `team_salary_cap` y publica `salary_cap.calculado`
- `team-service` tambien publica `finance.patch` como delta listo para WebSocket
- `gateway` escucha `finance.patch` y lo reenvia al frontend sin calcular gameplay
- el frontend guarda `financeState` y muestra salarios comprometidos, cap space, luxury tax space y estado en la pestaña `Temporada`
- `agent-service` escucha `salary_cap.calculado`; si la franquicia entra en luxury tax, el CFO sube alerta presupuestaria y el Owner pierde paciencia
- contratos documentados en `shared/events/finance.md` y `shared/events/websocket_deltas.md`
- nota de progreso creada en `docs/02_Progress/2026-07-14_M3_10_Salary_Cap_Acotado.md`
- pruebas corridas: `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `cargo test --manifest-path services/agent-service/Cargo.toml`, `npm run build --prefix frontend`

Limite:

- no se implemento flujo operativo de trades en este corte; la validacion concreta al aceptar/rechazar trades queda para `M3.12`, usando este estado financiero como fuente de verdad

### Bloque D — Ciclo dentro de la temporada

#### M3.11 — GMs rivales como agentes — YA REALIZADO

Objetivo:

Materializar los 30 GMs rivales con perfil propio.

Done:

- existe `rival_gms` con estilo de negociacion, urgencia actual, filosofia de construccion, necesidades de roster, historia con el GM jugador
- generacion deterministica al crear partida
- los GMs rivales no negocian solos en background: se "activan" cuando el GM jugador inicia un trade o cuando un rival quiere algo del jugador

Resultado 2026-07-15:

- `agent-service` agrega `RivalGMProfile` y catalogo deterministico de 30 GMs rivales, uno por franquicia NBA real
- cada GM rival tiene `rival_team_id`, `gm_agent_id`, nombre, equipo, estilo de negociacion, urgencia actual, filosofia de construccion, necesidades de roster, confianza inicial e historial inicial con el GM de PulseCity
- se agrega tabla `rival_gms` con `PRIMARY KEY (game_id, rival_team_id)` y `UNIQUE (game_id, gm_agent_id)`
- se agrega migracion `020_create_rival_gms.sql`
- `agent-service` siembra los GMs rivales de forma idempotente al fundar/cargar partida junto con agentes individuales, jugadores-agente y relaciones
- la generacion varia por `game_id`, pero es estable para la misma partida
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`

Limite:

- no se implemento evaluacion de propuestas ni mutacion de roster; los GMs quedan materializados como base para `M3.12`

#### M3.12 — Trades operativos — YA REALIZADO

Objetivo:

Permitir trades reales durante la temporada regular.

Flujo:

```text
GM jugador propone trade
team-service valida contra cap
rival GM evalua segun perfil + necesidades + relacion previa
trade.contraoferta | trade.rechazada | trade.aceptada
si aceptada -> jugador recibido se materializa con ratings calibrados
team-service actualiza roster
agent-service actualiza estado emocional de los jugadores afectados
narrative-service genera evento templateado de trade
city-service evalua impacto en fanbase si era jugador querido
```

Done:

- un trade puede proponerse, contraofertarse y cerrarse
- la validacion de cap bloquea trades imposibles
- el Player Personnel (agente) media la negociacion narrativamente
- el sistema es idempotente ante retries de `trade.aceptada`

Avance parcial 2026-07-15 — M3.12a — Propuesta y evaluacion rival:

- `gateway` agrega endpoint `POST /api/v1/games/{gameID}/trades/proposals`
- el endpoint publica `decision.gm_registrada` con `kind = trade_proposal`
- `team-service` consume la decision, valida que el jugador ofrecido pertenezca al roster activo y calcula salary cap proyectado
- se agrega tabla `team_trades` y migracion `021_create_team_trades.sql`
- si la propuesta falla validacion local, `team-service` persiste y publica `trade.rechazada`
- si pasa validacion local, `team-service` persiste y publica `trade.propuesta_enviada`
- `agent-service` consume `trade.propuesta_enviada`, lee su tabla propia `rival_gms` y evalua la respuesta segun necesidades, urgencia, estilo de negociacion y confianza inicial
- se agrega tabla `rival_gm_trade_evaluations` y migracion `022_create_rival_gm_trade_evaluations.sql` para idempotencia de evaluacion
- `agent-service` publica `trade.rechazada` o `trade.contraoferta`
- `gateway` escucha `trade.*` y transforma los eventos en `trade.patch` delta-only para WebSocket
- contratos documentados en `shared/events/trades.md`, `shared/events/decisions.md` y `shared/events/websocket_deltas.md`
- pruebas corridas: `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `cargo test --manifest-path services/agent-service/Cargo.toml`

Limite M3.12a:

- todavia no se implemento `trade.aceptada`, aceptacion de contraoferta ni mutacion de roster; eso queda para `M3.12b`

Resultado 2026-07-15 — M3.12b — Aceptacion y mutacion de roster:

- `gateway` agrega endpoint `POST /api/v1/games/{gameID}/trades/acceptances`
- el endpoint publica `decision.gm_registrada` con `kind = trade_acceptance`
- `team-service` consume la aceptacion, valida que la propuesta exista, siga abierta y que el jugador ofrecido siga activo
- `team-service` revalida salary cap al momento de aceptar
- si la aceptacion ya no es posible, `team-service` publica `trade.rechazada` con motivo explicito
- si la aceptacion es valida, `team-service` marca al jugador saliente como `roster_status = traded`
- `team-service` materializa un jugador entrante deterministico calibrado por `proposal_id`, posicion requerida, salario entrante y rating del jugador saliente
- `team-service` actualiza `team_trades` con `status = accepted`, jugador entrante y asset adicional aceptado
- `team-service` publica `trade.aceptada`, `roster.patch`, `salary_cap.calculado` y `finance.patch`
- `agent-service` consume `trade.aceptada`, actualiza la capa emocional del jugador saliente, crea la capa emocional del jugador entrante y publica `roster.patch` emocional
- `narrative-service` consume `trade.aceptada`, espera delay deterministico y genera narrativa templateada `post_trade`
- se agrega migracion `023_create_agent_processed_trades.sql` para idempotencia de reacciones emocionales ante trades aceptados
- `gateway` transforma `trade.aceptada` en `trade.patch` con `status = accepted`
- contratos actualizados en `shared/events/trades.md`, `shared/events/decisions.md`, `shared/events/narrative.md` y `shared/events/websocket_deltas.md`
- pruebas corridas: `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...`, `cargo test --manifest-path services/agent-service/Cargo.toml`

Limite:

- no se agrego UI especifica de negociacion, assets reales de draft ni impacto economico profundo de ciudad; `accepted_additional_asset` queda como string acotado

### M3 — PAUSA ANALISIS ANTES DE SEGUIR — 15 JULIO 2026

El 15 de julio de 2026, con `M3.12b` recien cerrado, se hizo una revision completa del proyecto: canon, notas de progreso, los 8 servicios, frontend, tests, migraciones y contratos de eventos. Todos los builds compilan y todos los tests pasan (26 tests de `agent-service`, 9 de `match-service`, 53 funciones de test en Go, build limpio del frontend).

Lo que la revision confirmo que esta bien y no se toca:

- la disciplina de proceso (mini milestones, notas de progreso, este documento) es el activo mas valioso del proyecto
- los principios de arquitectura se respetan en el codigo, no solo en el papel: ownership estricto verificado (ej: `team_player_match_states` como proyeccion local en vez de leer tablas ajenas), idempotencia real via tablas dedicadas, `match-service` stateless, deltas por WebSocket
- el scoping de M3 (rivales abstractos, LLM solo en chat, draft/FA diferidos) es correcto y no se reabre
- la arquitectura general no se reabre: es un proyecto de aprendizaje de sistemas distribuidos y la complejidad es el punto

Lo que la revision detecto para corregir ahora, antes de seguir con el Bloque D:

1. **Gap de validacion de gameplay.** ~17.000 lineas de backend contra ~4.300 de frontend. Decisiones medicas (M3.9b): endpoint listo, sin UI. Trades (M3.12): flujo completo de propuesta/contraoferta/aceptacion, sin UI de negociacion. Salary cap: solo numeros en la tab Temporada. El modelo de partidos v2, las lesiones, el cap y los trades funcionan segun los tests, pero ningun humano los jugo nunca. Un test unitario dice que el trade es idempotente; no dice si la contraoferta del GM rival se siente viva o robotica. Con el plan original, la Vista Franquicia (M3.21-23) llegaria con ~10 sistemas sin tocar por un humano, y descubrir ahi que algo se siente mal implicaria reabrir sistemas cerrados hace meses.
2. **Archivos escalando a monolitos.** `agents.rs` (2778 lineas), `http.go` del gateway (1641), `store.go` de team-service (1545) y el god-hook `useNewGameFlow.ts` (1313) contradicen la regla de modulos por dominio. Partirlos hoy cuesta una sesion; en M3.20 costaria cinco.
3. **Los tests cubren lo que menos se rompe.** El dominio esta bien testeado, pero `persistence/`, `nats/` y `handlers/` de casi todos los servicios Go tienen cero tests, y ningun test cruza dos servicios. En un sistema distribuido por eventos, los bugs caros viven entre servicios: un typo en un subject NATS o un campo renombrado en un payload rompe la conversacion sin que ningun test unitario se entere. Los cuatro flujos e2e del test plan hoy son solo texto.
4. **Drift entre AGENTS.md y la realidad.** AGENTS.md promete Three.js (no hay: el mapa es grilla CSS), k3s (no existe `infra/k8s`), GitHub Actions (no hay workflows) y Railway (no hay deploy). Un agente que lea AGENTS.md hoy asume cosas falsas.

Nota: la revision tambien detecto que `docs/` y `.agents/` no estaban trackeados en git. Eso se resuelve directamente con un commit fuera de este documento y no se lista como mini milestone.

Decision:

El Bloque D se pausa antes de `M3.13`. Los cuatro cortes siguientes (`M3.P1` → `M3.P4`) se ejecutan en orden antes de retomar Trade Deadline. Razon central: `M3.13` es UI de presion sobre el flujo de trades, y ese flujo nunca fue jugado; construir el deadline encima de un sistema no validado invierte el orden natural.

Regla nueva vigente desde esta pausa:

Todo mini milestone futuro que agregue una mecanica nueva incluye en su done "operable desde el frontend, aunque sea minimo". En M2 esta regla existia de facto (`M2.13` fue "frontend jugable"); en M3 se perdio y el costo fue acumular sistemas sin validar. La regla queda tambien escrita en AGENTS.md como principio permanente.

#### M3.P1 — UI minima jugable de trades y decisiones medicas — PENDIENTE

Objetivo:

Cerrar el gap de validacion de gameplay: que el GM pueda operar trades y decisiones medicas desde el frontend, aunque la UI sea nivel debug, antes de apilar mas sistemas encima.

Incluye:

- panel de trades (calidad debug aceptable): proponer trade eligiendo jugador propio, posicion requerida y asset adicional
- visualizacion de los estados del trade via `trade.patch`: `propuesta_enviada`, `contraoferta`, `rechazada`, `aceptada`
- aceptar una contraoferta desde la UI (`POST /api/v1/games/{gameID}/trades/acceptances`)
- responder decisiones medicas desde inbox o panel roster: `rest`, `reduce_minutes`, `ignore_doctor`, `force_return`
- el roster muestra disponibilidad: `availability`, severidad, fecha estimada de retorno

Done:

- una temporada jugada a mano con al menos 3 trades propuestos y todas las decisiones medicas respondidas desde el frontend, sin tocar curl ni scripts
- fricciones y sensaciones documentadas en nota de `02_Progress/`: ¿se entiende el estado del trade? ¿la respuesta del rival se siente coherente con su perfil? ¿la decision medica tiene peso?
- lo aprendido queda registrado como input de diseño para `M3.13`

Avance 2026-07-15 — M3.P1a — Trade Center jugable — YA REALIZADO:

- se decide que las mecanicas principales no se acumulan en una sola pagina; el Command Center resume y navega, mientras cada dominio con flujo propio vive en una pagina o espacio de trabajo propio
- la regla queda registrada de forma permanente en `AGENTS.md`
- se agrega ruta `/franchise/trades` y pagina independiente `TradeCenterPage`
- el Command Center incorpora un acceso visible al Trade Center y la pagina ofrece retorno explicito
- se agrega formulario de propuesta con jugador propio, franquicia rival, posicion solicitada e incoming salary
- se materializa el catalogo canonico de 30 franquicias y GMs rivales en el frontend
- el cliente tipa y aplica `trade.patch` de forma delta-only, indexando negociaciones por `proposal_id`
- la sala de respuestas muestra estados `proposed`, `countered`, `rejected` y `accepted`
- las contraofertas muestran el asset adicional requerido y pueden aceptarse desde la UI
- estado, requests y errores del dominio viven en `useTradeOperations`; no se agrega la mecanica al god-hook mas alla del dispatch WebSocket y la composicion temporal previa a `M3.P2`
- estilos propios en `components/trades/tradeCenter.css`, responsivos y alineados al design system
- diseño validado registrado en `docs/plans/2026-07-15-m3-p1a-trade-center-design.md`
- pruebas corridas: `npm run build --prefix frontend`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `make build`

Limites M3.P1a:

- el selector de jugador usa estado vivo recibido por `roster.patch`; la rehidratacion REST completa del roster sigue reservada para `M3.22`
- el playtest manual de al menos tres propuestas y la documentacion de game feel se cierran en `M3.P1c`, despues de agregar Centro Medico en `M3.P1b`
- `M3.P1` permanece pendiente hasta cerrar medicina jugable y el playtest conjunto

#### M3.P2 — Particion de monolitos — PENDIENTE

Objetivo:

Partir los archivos que escalaron a monolitos antes de que crezcan mas, respetando la regla de modulos por dominio que el propio AGENTS.md declara.

Estado actual:

- `services/agent-service/src/agents.rs`: 2778 lineas (≈750 son el catalogo de templates de los 30 agentes: data disfrazada de codigo)
- `services/gateway/internal/handlers/http.go`: 1641 lineas
- `services/team-service/internal/persistence/store.go`: 1545 lineas
- `frontend/src/features/new-game/hooks/useNewGameFlow.ts`: 1313 lineas — god-hook que concentra todo el estado del cliente, el WebSocket y el REST

Incluye:

- `agents.rs`: extraer catalogo/templates/seeds a modulo propio (`catalog.rs` o similar); separar reacciones por dominio (partidos, decisiones medicas, trades, salary cap)
- `http.go`: separar handlers por recurso (auth, games, chat, trades, medical, ws)
- `store.go`: partir por agregado (roster, box scores, injuries, trades, salary cap, season)
- `useNewGameFlow.ts`: separar el dispatch del WebSocket, el cliente REST y el estado por dominio (map, time, season, finance, roster, relations, chat). Este corte es prerequisito duro de `M3.21` — es el paso cero del Bloque F
- refactor puro: cero cambios de comportamiento

Done:

- ningun modulo mezcla mas de un dominio
- data-as-code (catalogos, seeds) vive separada de la logica
- todos los tests siguen pasando sin modificar asserts
- `make build`, `make test` y `npm run build --prefix frontend` limpios

#### M3.P3 — Smoke de integracion end-to-end — PENDIENTE

Objetivo:

Cubrir el hueco mas caro del test plan: hoy ningun test cruza dos servicios, y ahi es donde viven los bugs de un sistema por eventos (subjects, payloads, timing).

Incluye:

- script de smoke contra el entorno `docker-compose` (`make smoke` o similar) que:
  - funda una partida nueva
  - simula N dias
  - verifica que llegaron deltas WebSocket de cada tipo activo (mapa, tiempo, roster, relations, finance)
  - propone un trade y verifica el ciclo completo hasta `trade.patch`
  - dispara una decision medica y verifica `roster.patch` + `relations.patch`
  - abre un chat con un agente usando provider `stub`
- sin framework nuevo: un binario Go o script shell que hable con el gateway alcanza

Done:

- `make smoke` corre el flujo completo y falla con mensaje claro si un subject o payload se rompio
- se agrega a los comandos esperados al cierre de cada mini milestone (seccion Build/checks de este documento)
- los flujos e2e de playoffs y cierre se incorporan al smoke cuando existan (`M3.16`+)

#### M3.P4 — Sincronizar AGENTS.md con la realidad — PENDIENTE

Objetivo:

Eliminar el drift entre lo que AGENTS.md promete y lo que existe, para que cualquier agente (IA o humano) que lo lea obtenga una imagen fiel del proyecto.

Decisiones tomadas en la revision del 15 de julio:

- **React + Three.js sigue en el plan** (confirmado por el GM). La grilla CSS actual del mapa es interina; la vista WebGL se aborda como parte de la Vista Ciudad expandida en `M4`. Documentarlo asi en AGENTS.md.
- k3s, GitHub Actions y Railway siguen en el plan pero no existen todavia: marcarlos como pendientes con milestone tentativo, no como presente
- la regla anti-monolito queda escrita en AGENTS.md como principio permanente (agregada en la sesion del 15 de julio)

Incluye:

- actualizar tabla de stack y catalogo de servicios de AGENTS.md distinguiendo estado real vs planeado
- registrar que `gatewayBaseUrl`/`socketBaseUrl` hardcodeados en `useNewGameFlow.ts` pasan a env (`VITE_*`) a mas tardar al primer deploy — puede diferirse, pero queda anotado

Done:

- AGENTS.md no afirma nada que el repo contradiga
- un clone fresco + lectura de AGENTS.md da una imagen fiel del proyecto

### Bloque D — Ciclo dentro de la temporada (continuacion)

#### M3.13 — Trade Deadline — PENDIENTE

Objetivo:

Implementar la fecha del Trade Deadline como evento dramatico real.

Done:

- el sistema detecta automaticamente el dia del deadline en el calendario
- pausa automatica + UI expandida
- el mercado queda abierto en condiciones de mayor presion temporal
- evento `trade.deadline_activado` y `trade.deadline_cerrado`
- despues del cierre, no se aceptan mas trades hasta proxima offseason

#### M3.14 — Eventos criticos de agentes — PENDIENTE

Objetivo:

Implementar los eventos que interrumpen el tiempo si o si.

Casos:

- jugador pide traspaso formal (`jugador.solicita_traspaso`)
- Owner convoca reunion de emergencia
- ultimatum de un agente con deadline real (24-48hs simuladas)
- conflicto detectado por el Sports Psychologist que pasa a critico

Done:

- estos eventos pausan el tiempo independientemente del modo de notificaciones
- evento `agente.evento_critico` con `urgencia` y `deadline` opcional
- frontend bloquea avance hasta que se responde

#### M3.15 — Narrativa rule-based ampliada — PENDIENTE

Objetivo:

`narrative-service` cubre todos los eventos sistemicos nuevos de M3.

Incluye:

- post-trade (origen, recepcion, impacto en fanbase)
- post-lesion (alerta inicial, declaracion de baja, alta)
- post-conflicto detectado
- post-evento critico
- post-resultado de serie de playoffs
- post-cierre con veredicto del Owner

Done:

- cada evento genera narrativa templateada con estado real del agente involucrado
- el LLM no se activa para estos eventos
- los chats libres siguen siendo el unico canal con LLM real

### Bloque E — Playoffs y cierre

#### M3.16 — Bracket de playoffs — PENDIENTE

Objetivo:

Implementar el bracket completo.

Done:

- al terminar la temporada regular se genera bracket de 16 equipos segun record
- series al mejor de 7
- 4 rondas hasta el campeonato
- evento `playoffs.bracket_generado` y `playoffs.serie_iniciada` por cada serie
- `match-service` simula partidos de playoffs igual que regular (mismo simulador v2)

#### M3.17 — Intensidad emocional en playoffs — PENDIENTE

Objetivo:

Que los playoffs se sientan distintos de la temporada regular.

Incluye:

- rotacion acortada automaticamente por el Coach
- Sports Psychologist mas activo (mas alertas)
- multiplicador de impacto emocional en momentos clave
- Owner aparece mas frecuentemente
- la prensa intensifica cobertura (cambio de tono de eventos narrativos templateados)
- jugadores fuera de rotacion acumulan insatisfaccion mas rapido

Done:

- las relaciones se mueven mas fuerte en playoffs que en regular
- el frontend muestra el bracket avanzando partido a partido

#### M3.18 — Tres desenlaces de playoffs — PENDIENTE

Objetivo:

Implementar los tres caminos del canon: eliminacion temprana, final sin campeonato, campeonato.

Done:

- eliminacion temprana: Owner activa, ciudad entra en modo decepcion, jugadores con contrato expirando reevaluan
- final sin campeonato: Owner evalua segun expectativas iniciales del escenario, ciudad celebra parcialmente, jugadores suben ambicion
- campeonato: terremoto en todos los sistemas — ciudad explota economicamente, Owner expande presupuesto, jugadores suben lealtad, sponsors entran en cola, prensa construye narrativa de dinastia, GMs rivales recalculan reputacion
- evento `playoffs.campeon_definido` con `team_id` y contexto completo

#### M3.19 — Cierre de temporada con Owner — PENDIENTE

Objetivo:

Implementar la evaluacion formal con el Owner como evento obligatorio de cierre.

Done:

- evento `temporada.cerrada` dispara `owner.evaluacion_iniciada`
- el sistema calcula veredicto segun resultado deportivo + estado economico + paciencia restante + perfil inicial del Owner
- evento `owner.veredicto_emitido` con uno de: renovado, advertencia, ultimatum (con condiciones explicitas), despido
- si es despido, la partida entra en estado terminal (game over) — opcion de empezar nueva partida o cargar save
- evento narrativo templateado de cierre

#### M3.20 — Retencion simple de staff y jugadores — PENDIENTE

Objetivo:

Permitir renovaciones in-line sin abrir Free Agency real.

Done:

- jugadores propios con contrato expirando aparecen en una lista de "renovar antes de offseason"
- decision binaria simple: ofrecer extension o no
- terminos basicos derivados de su satisfaccion acumulada y su rating
- jugadores no renovados quedan fuera del roster al iniciar nueva temporada
- mismo flujo para staff
- sin pool de FA: lo que entra y sale entre temporadas se modela solo via retencion + trades

### Bloque F — Frontend completo M3

> Este bloque cubre UX minima viable para M3, no polish. El polish integral de toda la app (onboarding profundo, wiki interna, sonido, animaciones narrativas pulidas, accesibilidad, settings de usuario, test de usabilidad) queda diferido a `M5 — Polish y Experiencia Final`.

#### M3.21 — Vista Franquicia funcional — PENDIENTE

Objetivo:

Construir la pantalla principal del jugador segun el canon de Experiencia.

Done:

- estadio por zonas con agentes presentes en cada zona
- toolbar lateral con accesos a paneles
- header con metricas criticas: record, cap space, proximo partido, alertas activas
- alternancia con la Vista Ciudad existente (CeremonyPage / nueva)
- buscador de agentes siempre accesible

#### M3.22 — Paneles de acceso rapido — PENDIENTE

Objetivo:

Implementar los paneles que la toolbar abre.

Incluye:

- panel Roster: los 15 jugadores con estado, contrato, minutos, estado emocional resumido
- panel Finanzas/Salary Cap: cap actual, contratos, proyeccion
- panel Calendario: proximos rivales, resultados recientes, fases del ciclo
- panel Bandeja de eventos: eventos pendientes por urgencia
- panel Playoffs: bracket cuando playoffs estan activos

Done:

- cada panel rehidrata por REST y se actualiza por deltas
- los paneles son accesibles desde cualquier zona del estadio

#### M3.23 — Inbox narrativo unificado — PENDIENTE

Objetivo:

Que el jugador tenga un solo lugar donde vivir su narrativa.

Done:

- eventos rule-based + chats activos + eventos criticos pendientes conviven en el inbox
- diferenciacion visual clara entre evento sistemico y chat
- eventos criticos quedan marcados y bloquean avance del tiempo

### Bloque G — Cierre integrador

#### M3.24 — Temporada de prueba end-to-end — PENDIENTE

Objetivo:

Smoke real con todo M3 activo.

Done:

- una partida nueva atraviesa una temporada completa
- los 30 agentes vivos generaron eventos durante la temporada
- el jugador inicio al menos N chats con agentes distintos (objetivo concreto en este corte)
- ocurrieron al menos M lesiones, K trades, los playoffs llegaron a un desenlace, el Owner emitio veredicto
- las metricas finales se documentan en el resultado de este corte
- bug-hunt sobre idempotencia, reconexion, multi-cliente, fallos del LLM

Con esto `Milestone 3` queda cerrado funcionalmente y queda preparado el camino para `M3.5` (Draft + Free Agency) antes de `M4`.

## Edge cases a contemplar

### Chat con LLM

- si el LLM no responde, fallback templateado en personaje
- si el LLM responde fuera de dominio, se descarta y se reintenta con prompt mas estricto, sin pedir al usuario que reintente
- cap de tokens por turno respetado siempre
- cap de turnos por sesion respetado siempre
- dos clientes mirando la misma partida pueden abrir chats distintos en paralelo sin chocar
- no se borra historial automaticamente: el jugador puede consultar conversaciones anteriores

### Trades

- una propuesta no aceptada no muta roster
- un trade aceptado se persiste antes de publicar evento para evitar inconsistencia ante crash
- un retry de `trade.aceptada` no duplica jugadores ni dispara doble narrativa
- el cap se valida antes de aceptar; si falla, el trade se rechaza con motivo explicito
- los jugadores afectados acumulan estado emocional una sola vez por trade

### Lesiones

- una lesion no se duplica si llega dos veces el mismo evento
- la recuperacion se publica una sola vez
- si el GM presiona alta y hay reagravacion, la nueva lesion es un evento nuevo, no una mutacion de la anterior

### Playoffs

- el bracket se genera una sola vez al cerrar la regular
- si una serie ya termino, recibir de nuevo `playoffs.serie_terminada` no avanza la ronda dos veces
- partidos de playoffs son idempotentes igual que los de regular
- si la partida se cierra a mitad de una serie, al volver el bracket se rehidrata por REST y los deltas siguientes son delta-only

### Cierre

- el veredicto del Owner se calcula una sola vez por temporada
- si el veredicto es despido, ningun evento posterior a `owner.veredicto_emitido` debe mutar estado
- el jugador puede iniciar nueva partida desde game over

### Multi partida y multi cliente

- todo evento sigue incluyendo `game_id`
- el chat se filtra por `game_id` y por `(game_id, agent_id)` en sus deltas
- dos partidas activas no comparten LLM history ni relations ni log de decisiones

### Servicios caidos

- si `narrative-service` cae, el chat queda inaccesible pero los partidos siguen
- si `agent-service` cae, el sistema queda sin reacciones de agentes pero los partidos se simulan
- al recuperar, los eventos no consumidos pueden reprocesarse gracias a la idempotencia
- ningun servicio escribe estado corrupto al volver

### WebSocket

- regla no negociable de M2 se mantiene: solo deltas, nunca estado completo
- ante reconexion, snapshot inicial via REST + continuacion con deltas
- chat reconectado rehidrata historial via REST por `(game_id, agent_id)`

## Test plan

### Unit tests

- inicializacion de estado para cada uno de los 30 agentes nuevos
- inicializacion de capa emocional para cada arquetipo de jugador
- evolucion de relaciones ante eventos canonicos
- modelo de partidos v2: ratings altos vs ratings bajos producen distribucion esperada de resultados
- impacto de minutos del Coach en rendimiento
- impacto de estado emocional en rendimiento
- generacion de lesiones por carga
- riesgo de reagravacion por alta temprana
- calculo de salary cap, luxury tax line
- validacion de trade contra cap
- generacion de bracket de playoffs segun record
- veredicto del Owner segun matriz (resultado, economia, paciencia, perfil inicial)
- generacion templateada de narrativa para cada nuevo tipo de evento

### Integration tests

Flujos minimos esperados:

```text
trade end-to-end:
GM jugador propone -> rival evalua -> contraoferta -> aceptacion
-> cap actualizado -> roster actualizado -> agentes reaccionan
-> narrativa generada -> deltas al frontend

lesion end-to-end:
jugador acumula carga -> riesgo sube -> lesion publicada
-> medico alerta -> GM decide -> consecuencias persistidas

playoffs end-to-end:
regular cerrada -> bracket generado -> serie iniciada
-> partidos simulados -> serie terminada -> ronda avanzada
-> campeon definido o eliminacion -> efectos sistemicos

chat end-to-end:
jugador abre chat con agente -> contexto se ensambla
-> LLM responde -> mensaje persistido -> delta WebSocket
-> agente registra interaccion en relacion con GM
```

### Idempotencia

Probar que reintentar estos eventos no duplica estado:

- `decision.gm_registrada`
- `jugador.lesionado`
- `jugador.recuperado`
- `trade.aceptada`
- `trade.deadline_cerrado`
- `playoffs.serie_terminada`
- `playoffs.campeon_definido`
- `owner.veredicto_emitido`
- `agente.relacion_cambio`

### Build/checks

Comandos esperados al cierre de mini milestones:

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
cargo test --manifest-path services/match-service/Cargo.toml
GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...
GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...
GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...
make test-rust
make test-go
make build
npm run build --prefix frontend
```

## Supuestos explicitos

- M3 no implementa Draft.
- M3 no implementa Free Agency.
- M3 no implementa pool global de jugadores libres.
- M3 no implementa pretemporada con UI propia.
- M3 no implementa 30 rosters rivales completos.
- M3 no implementa trades nominales (das X recibis Y de un roster rival real).
- M3 no implementa player development (curvas de progresion/declive por edad).
- M3 no implementa exceptions del salary cap, ni player/team options, ni no-trade clauses.
- M3 no implementa edificios especiales nuevos ni alcalde con agenda independiente.
- M3 no implementa cache de respuestas LLM ni optimizacion avanzada de costo.
- M3 no implementa modo multijugador.
- M3 no implementa Vista Ciudad expandida nueva (la de M1/M2 sigue tal cual hasta M4).
- M3 no implementa polish profundo (onboarding completo, wiki interna, sonido, animaciones narrativas pulidas, accesibilidad, settings de usuario, test de usabilidad con personas externas). Eso queda para `M5 — Polish y Experiencia Final`.

## Filosofia de ejecucion

Cada sesion de trabajo dentro de M3 elige un mini milestone chico y lo cierra. No avanzar en tareas gigantes tipo "hacer Bloque C de un saque".

Formato recomendado por sesion (igual que M2):

```text
Objetivo de la sesion:
Mini milestone:
Archivos tocados:
Decision tomada:
Pruebas corridas:
Pendiente siguiente:
```

Al cerrar cada mini milestone, dejarlo en este documento como `YA REALIZADO` con su bloque `Resultado <fecha>:`, y registrar nota granular en `docs/02_Progress/` con prefijo `M3_<n>_<nombre>`.

La prioridad de M3 es orden y profundidad, no velocidad. Los agentes recien empiezan a vivir.

El sistema ya late. Ahora aprende a hablar.
