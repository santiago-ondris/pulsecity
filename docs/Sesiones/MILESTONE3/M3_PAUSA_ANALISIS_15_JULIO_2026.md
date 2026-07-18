# M3 — PAUSA ANALISIS ANTES DE SEGUIR — 15 JULIO 2026

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

#### M3.P1 — UI minima jugable de trades y decisiones medicas — YA REALIZADO

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

Avance 2026-07-15 — M3.P1b — Centro Medico jugable — YA REALIZADO:

- se agrega ruta `/franchise/medical` y pagina independiente `MedicalCenterPage`
- el Command Center incorpora acceso visible al Centro Medico y conserva el retorno explicito
- la pagina deriva el parte medico desde el estado vivo recibido por `roster.patch`; no calcula lesiones ni disponibilidad en frontend
- se muestran jugadores disponibles, casos activos, severidad, dias estimados y fecha esperada de retorno
- cada lesion permite registrar desde la UI las cuatro decisiones de M3.9b: `rest`, `reduce_minutes`, `ignore_doctor` y `force_return`
- el estado de envio, errores y confirmaciones por `injury_id` vive en `useMedicalOperations`
- una respuesta HTTP aceptada confirma que la decision fue registrada, pero la UI no muta disponibilidad de forma optimista; espera el siguiente delta del backend
- la pantalla distingue visualmente protocolo, carga reducida, ignorar recomendacion y alta forzada segun su riesgo, usando los colores semanticos del design system
- se contemplan estados sin snapshot, roster sano, envio en curso, error y decision registrada
- diseño validado registrado en `docs/plans/2026-07-15-m3-p1b-centro-medico-design.md`
- nota de progreso creada en `docs/02_Progress/2026-07-15_M3_P1b_Centro_Medico_Jugable.md`
- pruebas corridas: `npm run build --prefix frontend`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `make build`

Limites M3.P1b:

- las confirmaciones de decisiones medicas viven en memoria de la sesion; su rehidratacion REST completa sigue reservada para `M3.22` / `M3.23`
- `M3.P1` permanece pendiente hasta completar `M3.P1c` con playtest manual conjunto y documentacion de game feel

#### M3.P1c — Playtest conjunto de trades y medicina — YA REALIZADO

Objetivo:

Jugar manualmente los dos flujos que abrio `M3.P1` y registrar si las decisiones se entienden y tienen peso antes de construir Trade Deadline encima.

Done:

- al menos tres propuestas de trade realizadas desde Trade Center
- todas las decisiones medicas que aparezcan durante la temporada respondidas desde Centro Medico
- fricciones, claridad de estados y sensacion de las respuestas rivales/medicas registradas en `docs/02_Progress/`
- aprendizajes concretos asentados como input de diseño para `M3.13`
- `M3.P1` pasa a `YA REALIZADO`

Avance 2026-07-15 — M3.P1c-a — Reset de experiencia inicial — YA REALIZADO:

- la primera captura del playtest revelo un bloqueo previo al game feel de trades/medicina: no habia accion inicial dominante, el mapa invadia el panel de agentes, existia scroll horizontal y el contenido tecnico ocupaba la jerarquia principal
- se pausa el playtest funcional hasta corregir la experiencia base; este hallazgo forma parte de `M3.P1c`, no es polish diferido a M5
- una partida nueva y lista entra en un kickoff guiado, pausado, con CTA dominante `Comenzar temporada`
- el CTA usa el control temporal existente (`paused = false`, `speed = 1`); no se agrega estado ni endpoint de gameplay
- si el request de tiempo falla, el estado optimista hace rollback y el kickoff vuelve a mostrarse con error
- durante la generacion del mapa se conserva una ceremonia propia y contenida; el kickoff aparece solo despues de `mapa.generacion_completa` y de responder al Owner
- una partida con dias o partidos procesados entra directamente al Command Center operativo
- la vista por defecto deja de ser el directorio completo de agentes y pasa a `Resumen`
- la navegacion secundaria queda en `Resumen`, `Inbox`, `Staff` y `Sistema`; socket, pipeline y eventos crudos quedan confinados a `Sistema`
- el resumen muestra pulso de temporada, record, cap, alertas, roster, operaciones y resultados recientes
- Trade Center y Centro Medico conservan sus paginas propias; Staff abre buscador, directorio y chat bajo demanda
- la grilla CSS pasa a ser un snapshot de ciudad con proporcion fija, tracks encogibles y `overflow: hidden`
- el layout nuevo elimina margenes negativos/full-viewport del Command Center y define breakpoints para desktop, tablet y mobile sin scroll horizontal estructural
- `CeremonyMapPanel` y `SeasonPanel` se eliminan al quedar reemplazados, evitando componentes huerfanos
- diseño validado registrado en `docs/plans/2026-07-15-m3-p1c-a-command-center-kickoff-design.md`
- nota de progreso creada en `docs/02_Progress/2026-07-15_M3_P1c_a_Reset_Experiencia_Inicial.md`
- pruebas corridas: `npm run build --prefix frontend`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `make build`

Limite de validacion M3.P1c-a:

- el entorno de desarrollo del agente no dispone de navegador headless para generar capturas automaticas por viewport; la proxima revision visual en navegador del GM valida el resultado real y alimenta ajustes del mismo corte si fueran necesarios

##### M3.P1c-b — Retomar playtest funcional — PLAYTEST REALIZADO, HALLAZGOS REGISTRADOS

Con la experiencia base corregida, retomar el recorrido original: iniciar temporada desde el kickoff, realizar al menos tres propuestas de trade y responder las decisiones medicas que aparezcan. Registrar claridad, friccion y game feel como input para `M3.13`.

Resultado de playtest 2026-07-17:

- el GM pudo iniciar la temporada y navegar intuitivamente: supo donde entrar para operar cada mecanica sin asistencia externa
- Trade Center resulto claro en su alcance actual; se realizaron dos propuestas, llegaron contraofertas y una negociacion pudo aceptarse y mutar el roster correctamente
- Centro Medico permitio responder una decision medica sin friccion y desde el frontend
- el rediseño de `M3.P1c-a` queda validado funcionalmente: Command Center resume y navega, mientras trades y medicina viven en paginas propias
- el muestreo quedo en dos propuestas frente al objetivo formal de tres; se registra de forma explicita y no se presenta como tres ejecuciones
- hallazgo de identidad: los jugadores iniciales aparecen mediante `player_id` en lugares donde deberian verse nombre y apellido; los jugadores recibidos por trade si muestran nombre porque el parche contractual del trade incluye `full_name`
- hallazgo de agentes: gran parte de los agentes todavia tiene poca profundidad mecanica, acorde al roadmap pendiente, pero el chat tampoco refleja de forma confiable moods visibles como `frustrated`, `pressured` o `concerned`
- se confirma un desacople tecnico: la UI muestra moods dinamicos provenientes de `agent_core_states`, mientras el contexto de chat se ensambla desde `agent_individual_states`, cuya capa puede conservar el estado inicial sembrado
- hallazgo de ciclo: al terminar los 82 partidos la fecha puede continuar hasta julio sin transicion, pausa ni explicacion; Trade Deadline, Playoffs y Cierre siguen pendientes, pero la futura transicion debe evitar que el calendario avance silenciosamente fuera de fase
- nota detallada: `docs/02_Progress/2026-07-17_M3_P1c_b_Playtest_Funcional.md`

Decision de orden posterior al playtest:

- antes de `M3.P2` se hacen dos microcorrecciones, sin ampliar mecanicas: identidad del roster inicial y coherencia entre estado dinamico visible y contexto del chat
- el final silencioso de la temporada se registra como requisito de transicion para `M3.16`; no se implementan playoffs dentro de este corte
- despues de esas correcciones se retoma `M3.P2` como refactor puro sobre comportamiento ya corregido
- con `M3.P1d-a` y `M3.P1d-b` cerrados, el GM da por suficiente la cobertura funcional del recorrido y `M3.P1` queda `YA REALIZADO`; la desviacion de dos propuestas frente a tres permanece documentada

##### M3.P1d-a — Identidad visible del roster inicial — YA REALIZADO

Objetivo:

Hacer que todo jugador del roster propio se muestre con nombre y apellido desde el primer snapshot disponible, conservando `player_id` solo como identidad tecnica.

Done:

- Command Center, Staff/directorio, Trade Center y Centro Medico muestran `full_name` para los jugadores iniciales y los recibidos por trade
- los deltas siguen siendo delta-only y la identidad contractual conserva ownership de `team-service`
- ningun componente usa un hash/ID como etiqueta principal cuando existe nombre contractual
- pruebas del contrato y build del frontend limpios

Resultado 2026-07-17:

- `PlayerEmotionalPatch` de `agent-service` incorpora `full_name` y `position`
- todos los parches emocionales de partido y trade retransmiten la identidad que `agent-service` ya proyecta desde `team_roster_players`
- `AgentDirectoryPanel` usa nombre y apellido como etiqueta principal y conserva `player_id` como clave tecnica
- Trade Center y Centro Medico ya tenian fallback a `full_name`, por lo que reciben nombres desde el primer `roster.patch` emocional disponible
- no se agrego snapshot completo, endpoint REST ni escritura cruzada: `team-service` conserva ownership contractual
- contrato actualizado en `shared/events/websocket_deltas.md`
- nota de progreso: `docs/02_Progress/2026-07-17_M3_P1d_a_Identidad_Roster.md`
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`, `npm run build --prefix frontend`

##### M3.P1d-b — Estado dinamico coherente en chat — YA REALIZADO

Objetivo:

Garantizar que el contexto enviado al chat use el mismo estado emocional dinamico que el frontend presenta al GM.

Done:

- un agente visible como `concerned`, `pressured`, `frustrated` u otro mood dinamico recibe ese estado en el prompt del chat
- se conserva un unico ownership en `agent-service`; `narrative-service` solo lee contexto
- se elimina el desacople entre la capa legacy de los cinco agentes core y `agent_individual_states` sin agregar profundidad mecanica nueva
- pruebas cubren persistencia, lectura de contexto y prompt resultante

Resultado 2026-07-17:

- cada guardado transaccional de `agent_core_states` sincroniza el `mood` en `agent_individual_states.emotional_state`
- `EnsureSchema` reconcilia de forma idempotente partidas existentes al iniciar `agent-service`, evitando esperar al proximo evento para corregir el chat
- `narrative-service` sigue leyendo una sola fuente de contexto y no conoce la tabla legacy
- el prompt conserva el estado emocional en el contexto; se agrego cobertura explicita para verificarlo
- no se agregaron reacciones, personalidad ni contexto operativo nuevo: el corte corrige coherencia, no profundidad
- validacion SQL sobre la partida existente dentro de una transaccion con `ROLLBACK`: se reconciliaron cinco agentes core y quedaron cero diferencias antes de revertir la prueba
- nota de progreso: `docs/02_Progress/2026-07-17_M3_P1d_b_Estado_Chat.md`
- pruebas corridas: `cargo test --manifest-path services/agent-service/Cargo.toml`, `GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...`

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

##### M3.P2a — Catalogo y semillas de agentes — YA REALIZADO

Objetivo:

Extraer de `agents.rs` la data-as-code de agentes individuales, relaciones canonicas y GMs rivales sin cambiar comportamiento.

Resultado 2026-07-18:

- agregado `services/agent-service/src/agents/catalog.rs` como modulo dedicado al catalogo canonico
- movidos los templates de los 30 agentes individuales, las semillas de relaciones y el catalogo/perfiles deterministas de los 30 GMs rivales
- `agents.rs` conserva las APIs publicas existentes y queda enfocado en tipos, defaults operativos y logica de reacciones
- las fronteras nuevas usan visibilidad interna minima (`pub(super)`) para no ampliar la API del crate
- `agents.rs` bajo de aproximadamente 2778 a 1752 lineas; las 1051 lineas de data-as-code quedaron aisladas en el modulo propio
- no cambiaron contratos, persistencia, eventos NATS ni asserts de tests
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2a_Catalogo_Agentes.md`
- pruebas corridas: `cargo fmt --manifest-path services/agent-service/Cargo.toml -- --check`, `cargo test --manifest-path services/agent-service/Cargo.toml`, `cargo build --manifest-path services/agent-service/Cargo.toml`

##### M3.P2b — Reacciones a partidos — YA REALIZADO

Objetivo:

Separar de `agents.rs` las reacciones emocionales y relacionales disparadas por `partido.terminado`, sin modificar comportamiento ni contratos.

Resultado 2026-07-18:

- agregado `services/agent-service/src/agents/match_reactions.rs` como modulo del dominio partido
- movidas las reacciones de los cinco agentes core, los jugadores del roster y las relaciones inter-agente
- movidos al mismo modulo el calculo de contexto del partido, performance individual, resumenes y construccion de `roster.patch`/`agente.relacion_cambio`
- `agents.rs` conserva mediante reexports las funciones publicas `apply_match_finished`, `apply_match_to_player_agents` y `apply_match_to_relationships`
- los helpers numericos y de relaciones compartidos con salary cap, trades y medicina permanecen en el modulo padre para evitar duplicacion
- `agents.rs` bajo de 1752 a 1206 lineas y `match_reactions.rs` quedo en 563 lineas
- no cambiaron eventos NATS, payloads, persistencia ni asserts de tests
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2b_Reacciones_Partidos.md`
- pruebas corridas: `cargo fmt --manifest-path services/agent-service/Cargo.toml -- --check`, `cargo test --manifest-path services/agent-service/Cargo.toml`, `cargo build --manifest-path services/agent-service/Cargo.toml`, `git diff --check`

##### M3.P2c — Decisiones medicas y relaciones — YA REALIZADO

Objetivo:

Separar de `agents.rs` las reacciones relacionales disparadas por decisiones medicas del GM, sin modificar comportamiento ni contratos.

Resultado 2026-07-18:

- agregado `services/agent-service/src/agents/medical_reactions.rs`
- movida la reaccion publica `apply_gm_decision_to_relationships`, la matriz de efectos por `choice_id` y la construccion de `agente.relacion_cambio`
- `agents.rs` conserva la API publica mediante reexport
- el mutador generico de relaciones permanece compartido en el modulo padre porque tambien lo usa `match_reactions.rs`
- se conservan el filtro por `medical_decision`, los deltas de confianza, historial corto, source event e idempotencia existente
- `agents.rs` bajo de 1206 a 1112 lineas y `medical_reactions.rs` quedo en 106 lineas
- no cambiaron eventos NATS, payloads, persistencia ni asserts de tests
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2c_Decisiones_Medicas_Relaciones.md`
- pruebas corridas: `cargo fmt --manifest-path services/agent-service/Cargo.toml -- --check`, `cargo test --manifest-path services/agent-service/Cargo.toml`, `cargo build --manifest-path services/agent-service/Cargo.toml`, `git diff --check`

##### M3.P2d — Trades y GMs rivales — YA REALIZADO

Objetivo:

Separar de `agents.rs` la evaluacion de propuestas por GMs rivales y las reacciones emocionales del roster ante trades aceptados, sin modificar comportamiento ni ownership.

Resultado 2026-07-18:

- agregado `services/agent-service/src/agents/trade_reactions.rs`
- movido `RivalGMTradeEvaluation` junto con la logica de rechazo/contraoferta segun necesidades, salarios, urgencia, confianza y estilo del GM rival
- movida la reaccion emocional de jugadores salientes y entrantes y la construccion de su `roster.patch`
- `agents.rs` conserva mediante reexports el enum y las funciones publicas existentes
- el catalogo determinista de GMs rivales permanece en `catalog.rs`; la mutacion contractual y del roster sigue siendo ownership de `team-service`
- `agents.rs` bajo de 1112 a 933 lineas y `trade_reactions.rs` quedo en 193 lineas
- no cambiaron eventos NATS, payloads, persistencia ni asserts de tests
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2d_Trades_GMs_Rivales.md`
- pruebas corridas: `cargo fmt --manifest-path services/agent-service/Cargo.toml -- --check`, `cargo test --manifest-path services/agent-service/Cargo.toml`, `cargo build --manifest-path services/agent-service/Cargo.toml`, `git diff --check`

##### M3.P2e — Reacciones a salary cap — YA REALIZADO

Objetivo:

Separar de `agents.rs` las reacciones de CFO y Owner a `salary_cap.calculado`, sin mover el calculo financiero fuera de su servicio dueño ni cambiar contratos.

Resultado 2026-07-18:

- agregado `services/agent-service/src/agents/salary_cap_reactions.rs`
- movidas las reglas de estado y mood de CFO/Owner y la construccion de `agente.estado_cambio`
- `agents.rs` conserva mediante reexport `apply_salary_cap_to_core_agents`
- `team-service` mantiene ownership del calculo de salary cap; `agent-service` solo reacciona al evento
- `agents.rs` bajo de 933 a 862 lineas y `salary_cap_reactions.rs` quedo en 74 lineas
- la logica de produccion de `agents.rs` termina en la linea 326; el resto del archivo son tests existentes, por lo que no se los movio solo para reducir lineas
- con catalogo, partidos, medicina, trades y salary cap separados, queda cerrada la particion por dominio de `agent-service` dentro de M3.P2
- no cambiaron eventos NATS, payloads, persistencia ni asserts de tests
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2e_Reacciones_Salary_Cap.md`
- pruebas corridas: `cargo fmt --manifest-path services/agent-service/Cargo.toml -- --check`, `cargo test --manifest-path services/agent-service/Cargo.toml`, `cargo build --manifest-path services/agent-service/Cargo.toml`, `git diff --check`

##### M3.P2f-a — Sesiones guest y WebSocket del gateway — YA REALIZADO

Objetivo:

Iniciar la particion de `gateway/internal/handlers/http.go` separando el ciclo de WebSocket y ubicando la creacion/lectura de sesiones guest dentro del recurso de autenticacion, sin cambiar rutas ni comportamiento.

Resultado 2026-07-18:

- agregado `services/gateway/internal/handlers/websocket.go`
- movidos `serveWebSocket`, `clientIDFromRequest` y `guestOwnsGame` al modulo del ciclo de conexion
- movidos `createGuestSession` y `guestTokenFromRequest` a `auth.go`, que ya contenia register, login, sesion actual y upgrade de guest
- se mantuvo el paquete `handlers`; no se crearon paquetes artificiales ni cambiaron los nombres usados por `RegisterRoutes`
- se conservaron rehidratacion inicial por REST/store, snapshot inicial del socket, activacion/desactivacion de simulacion y eventos `tiempo.sesion_iniciada`/`tiempo.sesion_terminada`
- `http.go` bajo de 1641 a 1492 lineas; `websocket.go` quedo en 135 lineas
- no cambiaron rutas HTTP, eventos NATS, payloads, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2f_a_Sesiones_WebSocket_Gateway.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...`, `git diff --check`

##### M3.P2f-b — Games y snapshots del gateway — YA REALIZADO

Objetivo:

Separar de `gateway/internal/handlers/http.go` la fundacion, listado, detalle y rehidratacion HTTP de partidas, junto con sus normalizadores exclusivos, sin cambiar rutas ni comportamiento.

Resultado 2026-07-18:

- agregado `services/gateway/internal/handlers/games.go`
- movidos `startGame`, `listGames`, `getGame` y `getSnapshot`
- movidos los normalizadores de nombre, abreviatura, colores, escenario inicial y modo de gestion usados al fundar partida
- `answerOwnerIntro` permanece fuera porque pertenece al dominio narrativo, aunque consulte el estado de la partida
- se mantuvo el paquete `handlers` y `RegisterRoutes` no cambio
- se conservaron autorizacion user/guest, ownership, persistencia, publicacion de `mapa.generacion_iniciada` y rehidratacion snapshot desde store/cache
- `http.go` bajo de 1492 a 1266 lineas y `games.go` quedo en 238 lineas
- no cambiaron rutas HTTP, eventos NATS, payloads, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2f_b_Games_Snapshots_Gateway.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...`, `git diff --check`

##### M3.P2f-c — Narrativa y chat del gateway — YA REALIZADO

Objetivo:

Separar de `gateway/internal/handlers/http.go` la respuesta inicial del Owner y el inicio de chat con agentes, junto con el helper de elecciones narrativas, sin cambiar rutas ni comportamiento.

Resultado 2026-07-18:

- agregado `services/gateway/internal/handlers/narrative.go`
- movidos `answerOwnerIntro`, `startAgentChat` y `findNarrativeChoice`
- se conservaron autorizacion y ownership de partida, persistencia de la respuesta inicial, publicacion de `narrativa.respuesta_gm` y `decision.gm_registrada`
- se conservaron validacion y normalizacion de chat, limite de 1200 caracteres, generacion de `conversation_id` y publicacion de `agente.consulta_iniciada`
- se mantuvo el paquete `handlers` y `RegisterRoutes` no cambio
- `http.go` bajo de 1266 a 1050 lineas y `narrative.go` quedo en 227 lineas
- no cambiaron rutas HTTP, eventos NATS, payloads, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2f_c_Narrativa_Chat_Gateway.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...`, `git diff --check`

##### M3.P2f-d — Decisiones medicas del gateway — YA REALIZADO

Objetivo:

Separar de `gateway/internal/handlers/http.go` el endpoint de decisiones medicas y su tabla de etiquetas, sin mover reglas de lesiones o relaciones al gateway ni cambiar contratos.

Resultado 2026-07-18:

- agregado `services/gateway/internal/handlers/medical.go`
- movidos `answerMedicalDecision` y `medicalDecisionLabel`
- se conservaron autorizacion y ownership de partida, normalizacion/validacion de request y publicacion de `decision.gm_registrada`
- se conservaron los cuatro `choice_id`, agentes afectados, `source_event_id` y source `jugador.lesionado`
- el gateway continua actuando solo como traductor/autorizador; consecuencias medicas y relacionales permanecen en los servicios dueños
- se mantuvo el paquete `handlers` y `RegisterRoutes` no cambio
- `http.go` bajo de 1050 a 944 lineas y `medical.go` quedo en 116 lineas
- no cambiaron rutas HTTP, eventos NATS, payloads, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2f_d_Decisiones_Medicas_Gateway.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...`, `git diff --check`

##### M3.P2f-e — Trades del gateway — YA REALIZADO

Objetivo:

Separar de `gateway/internal/handlers/http.go` las entradas HTTP de propuesta y aceptacion de trades, sin mover evaluacion rival ni mutacion contractual al gateway y sin cambiar contratos.

Resultado 2026-07-18:

- agregado `services/gateway/internal/handlers/trades.go`
- movidos `proposeTrade` y `acceptTrade`
- se conservaron autorizacion y ownership de partida, normalizacion/validacion de requests y publicacion de `decision.gm_registrada`
- se conservaron payload salarial, agentes afectados, IDs deterministas de decision y sources de propuesta/contraoferta
- el gateway continua actuando solo como entrada y traductor; evaluacion del GM rival permanece en `agent-service` y roster/cap en `team-service`
- se mantuvo el paquete `handlers` y `RegisterRoutes` no cambio
- `http.go` bajo de 944 a 770 lineas y `trades.go` quedo en 184 lineas
- no cambiaron rutas HTTP, eventos NATS, payloads, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2f_e_Trades_Gateway.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...`, `git diff --check`

##### M3.P2f-f — Control de tiempo y shell HTTP — YA REALIZADO

Objetivo:

Cerrar la particion de `gateway/internal/handlers/http.go` separando control de tiempo y el HTML debug embebido, dejando el archivo original como shell de dependencias, rutas, health y JSON.

Resultado 2026-07-18:

- agregado `services/gateway/internal/handlers/time_control.go`
- movidos `updateTimeControl` y `validTimeSpeed`
- agregado `services/gateway/internal/handlers/debug_page.go` con `debugHTML` byte-a-byte equivalente al bloque original
- se conservaron autorizacion/ownership, velocidades validas `1/5/20`, eventos `tiempo.velocidad_cambiada`/`tiempo.pausa_activada` y delta `time.patch`
- `http.go` queda limitado a `Dependencies`, `RegisterRoutes`, debug/health handlers y `writeJSON`
- `http.go` bajo de 770 a 57 lineas; `time_control.go` quedo en 108 y `debug_page.go` en 613 lineas de data-as-code visual
- con sesiones/WebSocket, games, narrativa/chat, medicina, trades, tiempo y debug separados, queda cerrada la particion del gateway dentro de M3.P2
- no cambiaron rutas HTTP, eventos NATS, payloads, ownership, contenido visual ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2f_f_Tiempo_Shell_HTTP.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway vet ./...`, `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway build ./...`, comparacion exacta de `debugHTML`, `git diff --check`

##### M3.P2g-a — Salary cap del store — YA REALIZADO

Objetivo:

Iniciar la particion de `team-service/internal/persistence/store.go` extrayendo la operacion autonoma de persistencia del snapshot de salary cap, sin cambiar SQL, transacciones ni API.

Resultado 2026-07-18:

- agregado `services/team-service/internal/persistence/salary_cap.go`
- movido `saveSalaryCap` con el upsert completo de `team_salary_cap`
- se conservaron los 13 parametros SQL, conflicto por `game_id`, actualizacion de source/idempotencia y wrapping de error
- `queryer` permanece compartido en `store.go`, permitiendo usar la misma funcion dentro de las transacciones de fundacion y aceptacion de trade
- el calculo `domain.CalculateSalaryCap` y sus callers permanecen en sus dominios actuales; este corte solo separa persistencia
- el DDL de `team_salary_cap` permanece temporalmente en `EnsureSchema`; la particion de schemas sera un corte distinto
- `store.go` bajo de 1545 a 1516 lineas y `salary_cap.go` quedo en 37 lineas
- no cambiaron schema, SQL, contratos, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2g_a_Salary_Cap_Store.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`, `git diff --check`

##### M3.P2g-b — Persistencia de roster — YA REALIZADO

Objetivo:

Separar de `team-service/internal/persistence/store.go` la aplicacion de patches emocionales y las lecturas contractuales del roster, sin cambiar SQL, transacciones ni API interna.

Resultado 2026-07-18:

- agregado `services/team-service/internal/persistence/roster.go`
- movidos `ApplyRosterPatch`, `loadRoster` y `loadRosterPlayer`
- se conservaron el upsert idempotente por fecha simulada, conversiones `int16`/`uint8`, orden por `sort_order` y wrapping de errores
- los callers de trades y calendario continúan usando las mismas funciones internas del paquete
- la siembra inicial del roster permanece dentro de `SaveInitialSeason` para no fragmentar su transaccion atomica
- `store.go` bajo de 1516 a 1417 lineas y `roster.go` quedo en 107 lineas
- no cambiaron schema, SQL, contratos, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2g_b_Persistencia_Roster.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`, `git diff --check`

##### M3.P2g-c — Propuestas de trade — YA REALIZADO

Objetivo:

Separar de `team-service/internal/persistence/store.go` la aplicacion y persistencia de propuestas de trade, moviendo solo helpers exclusivos y conservando compartidos con aceptacion.

Resultado 2026-07-18:

- agregado `services/team-service/internal/persistence/trade_proposals.go`
- movido `ApplyTradeProposal`
- movidos `loadTradeStatus`, `insertProposedTrade`, inserts de rechazo y `parsePositiveInt`
- se conservaron idempotencia por proposal existente, rechazos por jugador inexistente/no activo, validacion contra luxury tax y transaccion unica
- `tradeRejectedFromDecision` permanece temporalmente compartido en `store.go` porque tambien lo usa aceptacion
- carga del trade y construccion de `roster.patch` permanecen para el corte de aceptacion
- `store.go` bajo de 1417 a 1251 lineas y `trade_proposals.go` quedo en 176 lineas
- no cambiaron schema, SQL, contratos, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2g_c_Propuestas_Trade.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`, `git diff --check`

##### M3.P2g-d — Aceptacion de trades — YA REALIZADO

Objetivo:

Separar de `team-service/internal/persistence/store.go` la aceptacion atomica de trades, su lectura persistida y la construccion de eventos/deltas derivados.

Resultado 2026-07-18:

- agregado `services/team-service/internal/persistence/trade_acceptance.go`
- movido `ApplyTradeAcceptance`
- movidos `loadTrade`, `storedTrade`, `tradeRejectedFromDecision` y `tradeRosterPatch`
- `trade_proposals.go` reutiliza el constructor comun de rechazos dentro del mismo paquete
- se conservaron idempotencia por status, validacion de propuesta/jugador/cap, materializacion determinista del jugador entrante y transaccion unica
- se conservaron mutacion del roster, actualizacion de `team_trades`, persistencia de salary cap y emision de `trade.aceptada`, `roster.patch` y `salary_cap.calculado`
- `store.go` bajo de 1251 a 1033 lineas y `trade_acceptance.go` quedo en 228 lineas
- con propuestas y aceptacion separadas, queda cerrada la particion del dominio trade dentro del store
- no cambiaron schema, SQL, contratos, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2g_d_Aceptacion_Trades.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`, `git diff --check`

##### M3.P2g-e — Decisiones medicas y lesiones — YA REALIZADO

Objetivo:

Separar de `team-service/internal/persistence/store.go` las decisiones medicas, la recuperacion diaria y los helpers de creacion y reagravacion de lesiones.

Resultado 2026-07-18:

- agregado `services/team-service/internal/persistence/injuries.go`
- movidos `ApplyMedicalDecision`, `RecoverPlayersForDate`, `createInjuriesForMatch`, `forcedReturnReaggravation`, `insertInjury` y `parseDate`
- `loadPlayerMatchStates` permanece junto a la preparacion de partidos porque tambien alimenta el simulador
- `ApplyMatchFinished` conserva la llamada a `createInjuriesForMatch` dentro del mismo paquete
- se conservaron force return, recuperacion diaria, reagravacion determinista, idempotencia por `ON CONFLICT (injury_id) DO NOTHING` y transacciones unicas
- `store.go` bajo de 1033 a 725 lineas y `injuries.go` quedo en 318 lineas
- con decisiones, recuperacion y persistencia separadas, queda cerrada la particion del dominio medico dentro del store
- no cambiaron schema, SQL, contratos, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2g_e_Decisiones_Medicas_Lesiones.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`, `git diff --check`

##### M3.P2g-f — Despacho de calendario — YA REALIZADO

Objetivo:

Separar de `team-service/internal/persistence/store.go` el claim del partido programado y la construccion del payload completo que consume `match-service`.

Resultado 2026-07-18:

- agregado `services/team-service/internal/persistence/schedule_dispatch.go`
- movido `DispatchScheduledMatchForDate`
- movidos `storedScheduleMatch`, `loadFranchise`, `loadOpponent` y `loadPlayerMatchStates`
- `loadPlayerMatchStates` queda junto a la preparacion del partido y sigue disponible para el calculo de lesiones dentro del mismo paquete
- `loadRoster` y `loadSeasonRecord` permanecen en sus dominios y son reutilizados por el despacho
- se conservaron el claim atomico `scheduled` a `scheduled_dispatched`, el orden determinista, el caso sin partido y la transaccion unica
- se conservo el payload completo con equipos, roster, estado de jugadores y record; `ApplyMatchFinished` no fue modificado
- `store.go` bajo de 725 a 517 lineas y `schedule_dispatch.go` quedo en 218 lineas
- no cambiaron schema, SQL, contratos, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2g_f_Despacho_Calendario.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`, `git diff --check`

##### M3.P2g-g — Resultado de partidos y record — YA REALIZADO

Objetivo:

Separar de `team-service/internal/persistence/store.go` la aplicacion de `partido.terminado`, la persistencia de box scores y el recalculo del record de temporada.

Resultado 2026-07-18:

- agregado `services/team-service/internal/persistence/match_results.go`
- movido `ApplyMatchFinished`
- movidos `recalculateSeasonRecord` y `loadSeasonRecord`
- `createInjuriesForMatch` permanece en `injuries.go` y se reutiliza dentro de la misma transaccion de resultado
- `queryer` permanece en `store.go` como frontera minima compartida por los modulos de persistencia
- se conservaron idempotencia por estado final, upsert de box scores, calculo local/visitante, `last_match_id` y transaccion unica
- el retry de un partido final conserva la lectura del record sin volver a generar lesiones ni mutar box scores
- `store.go` bajo de 517 a 346 lineas y `match_results.go` quedo en 180 lineas
- `store.go` queda concentrado en conexion, schema e inicializacion de temporada
- no cambiaron schema, SQL, contratos, ownership ni asserts de tests
- revision Go completada sin findings pendientes
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2g_g_Resultados_Record.md`
- pruebas corridas: `gofmt`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`, `git diff --check`

##### M3.P2g-h — Revision integral de persistencia — YA REALIZADO

Objetivo:

Cerrar `M3.P2g` revisando de forma integral la nueva particion de `team-service/internal/persistence`, sus responsabilidades, dependencias internas, tamaños y verificaciones Go.

Resultado 2026-07-18:

- revisados `store.go`, `salary_cap.go`, `roster.go`, `trade_proposals.go`, `trade_acceptance.go`, `injuries.go`, `schedule_dispatch.go` y `match_results.go`
- el `store.go` original de 1545 lineas quedo reducido a 346 lineas cohesivas de conexion, schema e inicializacion
- la persistencia completa quedo en 1610 lineas fisicas distribuidas por dominio; el overhead de 65 lineas corresponde a packages, imports y separacion de archivos
- ningun archivo supera 346 lineas: `injuries.go` 318, `trade_acceptance.go` 228, `schedule_dispatch.go` 218, `match_results.go` 180, `trade_proposals.go` 176, `roster.go` 107 y `salary_cap.go` 37
- los cruces internos son intencionales: roster para trades/calendario, salary cap para fundacion/trades, estado de partido para calendario/lesiones y record para despacho/resultados
- `queryer` queda como unica frontera de infraestructura compartida; no se agregaron packages, interfaces ni abstracciones adicionales
- no se encontraron simbolos duplicados, helpers huerfanos, imports incorrectos ni findings Go pendientes
- `M3.P2g` queda cerrado; no se recomienda seguir partiendo `store.go` en este momento
- no hubo cambios de codigo en este corte de revision
- nota de progreso: `docs/02_Progress/2026-07-18_M3_P2g_h_Revision_Integral_Persistencia.md`
- pruebas corridas: `gofmt -d services/team-service/internal/persistence`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`, `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`, `git diff --check`

#### M3.P3 — Smoke de integracion end-to-end — PENDIENTE

##### Por que existe este corte

Los tests unitarios actuales validan reglas y contratos dentro de cada servicio, pero no prueban una cadena real que cruce procesos. En PulseCity, los defectos mas costosos suelen aparecer precisamente entre servicios:

- un subject NATS escrito de forma distinta entre publisher y subscriber
- un payload que compila en ambos servicios pero no coincide al serializarse
- una publicacion que ocurre antes de confirmar la persistencia
- un delta WebSocket que no llega, llega con tipo incorrecto o pierde `game_id`
- una dependencia que esta saludable en aislamiento pero no en `docker-compose`
- un retry que duplica una consecuencia narrativa, emocional o contractual

`M3.P3` crea un smoke reproducible que comprueba los caminos esenciales de M3 contra el sistema levantado, sin intentar reemplazar la suite unitaria ni convertirse en un framework e2e general.

##### Objetivo

Agregar un comando `make smoke` que levante o utilice el entorno local, ejecute una partida controlada a traves del gateway y falle con un diagnostico claro cuando se rompa un contrato entre servicios.

##### Alcance original conservado

El smoke debe:

1. fundar una partida nueva
2. observar la conexion WebSocket de esa partida
3. simular una cantidad acotada de dias
4. verificar deltas activos de mapa, tiempo, roster, relaciones y finanzas
5. proponer un trade y seguir el ciclo hasta `trade.patch`
6. disparar una decision medica y verificar sus deltas derivados
7. abrir un chat con un agente usando el provider `stub`

Cuando existan los flujos de `M3.16+`, el mismo runner incorporara playoffs y cierre de temporada sin cambiar de herramienta.

##### Fuera de alcance

- no es un test de carga ni de performance
- no valida game feel ni calidad visual
- no reemplaza tests unitarios de reglas de dominio
- no prueba proveedores LLM reales; el chat usa el stub determinista
- no introduce Playwright, Cypress ni otro framework de browser
- no inspecciona directamente tablas propiedad de cada servicio como mecanismo principal de validacion
- no agrega coordinacion HTTP entre servicios
- no cambia contratos NATS ni WebSocket para facilitar el test

##### Frontera del runner

El runner se comporta como un cliente real de PulseCity:

```text
runner smoke
  -> REST del gateway para comandos y snapshots
  -> WebSocket del gateway para observar deltas
  -> nunca llama servicios internos directamente
  -> nunca escribe tablas de servicios
```

La excepcion permitida es consultar endpoints de salud de infraestructura para producir mejores mensajes de error antes de iniciar el escenario.

##### Tecnologia recomendada

Implementar un binario Go pequeno dentro de una carpeta de tooling del repositorio y exponerlo mediante `make smoke`.

Razon:

- el gateway y sus contratos ya usan Go
- la biblioteca estandar cubre HTTP, JSON, contextos y timeouts
- permite errores estructurados y salida legible sin depender de herramientas instaladas globalmente
- puede crecer por escenarios sin convertir un script shell en un parser de JSON complejo

No se agrega una abstraccion reusable hasta que exista una segunda necesidad concreta.

##### Principios del smoke

###### Caja negra

Las verificaciones se hacen sobre respuestas REST y deltas WebSocket observables. Si el estado interno es correcto pero el jugador no puede verlo u operarlo, el flujo falla.

###### Determinismo

Cada corrida genera identificadores propios y usa el provider de chat `stub`. Los asserts deben validar contratos y consecuencias, no textos variables ni timestamps exactos.

###### Timeouts explicitos

Cada espera tiene un timeout corto con contexto. Ninguna espera puede quedar bloqueada indefinidamente.

###### Diagnostico antes que silencio

Una falla informa como minimo:

- escenario y paso
- evento o delta esperado
- `game_id`
- ultimo estado HTTP relevante
- subjects/tipos observados durante la espera
- timeout o payload que provoco la discrepancia

###### Limpieza segura

El runner crea una partida nueva por ejecucion. No borra partidas ni datos existentes. La aislacion se obtiene con `game_id`, no con limpieza destructiva de la base.

##### Escenario base

###### 1. Preflight

- comprobar que gateway, NATS, Postgres y servicios requeridos estan disponibles
- confirmar que el narrative provider efectivo es `stub`
- abrir el WebSocket antes de disparar acciones que deban observarse
- iniciar un buffer acotado de deltas con registro de tipos recibidos

###### 2. Fundacion

- crear una partida guest por la API publica
- capturar `game_id` y credenciales de sesion
- completar el flujo minimo de fundacion
- esperar los deltas de generacion de mapa exigidos por el contrato vigente
- rehidratar snapshot por REST y comprobar que pertenece al mismo `game_id`

###### 3. Tiempo y temporada

- iniciar o reanudar la simulacion desde el gateway
- observar al menos un `tiempo.dia_avanzado` traducido al delta correspondiente
- avanzar solo los dias necesarios para alcanzar un partido o estado util para los escenarios siguientes
- verificar que los deltas observados siempre incluyen el `game_id` correcto

###### 4. Trade

```text
accion REST del GM
  -> decision.gm_registrada
  -> trade.propuesta_enviada / respuesta rival
  -> aceptacion cuando corresponda
  -> roster y salary cap persistidos
  -> trade.patch + roster.patch + finance.patch observables
```

El escenario debe seleccionar datos validos desde el snapshot actual; no puede depender de IDs hardcodeados de una partida anterior.

###### 5. Decision medica

```text
lesion disponible en snapshot/evento
  -> accion medica del GM
  -> decision.gm_registrada
  -> mutacion owned por team-service
  -> roster.patch y relations.patch observables
```

Si el calendario acotado no produce una lesion de forma determinista, el primer corte del smoke debe definir un mecanismo de fixture autorizado a traves de una entrada publica de debug. No se permite insertar directamente en `team_injuries` desde el runner.

###### 6. Chat con stub

- elegir un agente real del snapshot
- iniciar una conversacion por el endpoint publico del gateway
- esperar `chat.message`
- validar `game_id`, `agent_id`, `conversation_id`, rol y contenido no vacio
- consultar el historial y comprobar que el mensaje quedo persistido

###### 7. Cierre de corrida

- pausar la simulacion
- cerrar WebSocket y recursos del runner
- imprimir resumen de pasos, duracion y tipos de delta observados
- devolver exit code distinto de cero si cualquier assert fallo

##### Contratos observables minimos

El catálogo exacto debe tomarse de `shared/events/` y del codigo vigente al implementar cada escenario. Como minimo, el runner registra y puede esperar:

- progreso/finalizacion de mapa
- avance y control de tiempo
- `roster.patch`
- `relations.patch`
- `finance.patch`
- `trade.patch`
- `chat.message`
- eventos de partido necesarios para atravesar la temporada

WebSocket conserva la regla delta-only. El snapshot inicial se obtiene por REST; el runner no espera estado completo por WebSocket.

##### Estructura sugerida

La ubicacion exacta se confirma al iniciar el primer corte, despues de revisar el tooling existente. Estructura minima esperada:

```text
tools/smoke/
├── main.go             # flags, ejecucion y exit code
├── client.go           # REST + WebSocket del gateway
├── observer.go         # buffer y esperas de deltas
└── scenarios.go        # pasos del escenario M3
```

Si el primer corte cabe con claridad en menos archivos, se prefiere la estructura mas simple. Ningun archivo debe mezclar transporte, asserts y escenarios al superar las señales de particion de `AGENTS.md`.

##### Configuracion prevista

El runner debe aceptar configuracion explicita con defaults locales:

- URL base HTTP del gateway
- URL WebSocket del gateway
- timeout global
- timeout por espera
- nivel de detalle de salida

No se agregan opciones para saltar arbitrariamente asserts esenciales. Los flags de diagnostico no cambian el significado de una corrida exitosa.

##### Mini milestones de M3.P3

###### M3.P3a — Preflight y cliente del gateway — PENDIENTE

Done:

- `make smoke` ejecuta un binario versionado
- valida configuracion y salud del entorno
- crea una partida guest por REST
- abre/cierra WebSocket con timeout y cancelacion correctos
- errores de infraestructura distinguen servicio caido de assert funcional

###### M3.P3b — Observer de deltas — PENDIENTE

Done:

- el runner decodifica el envelope WebSocket vigente
- filtra estrictamente por `game_id`
- puede esperar uno o varios tipos sin perder mensajes intermedios
- al vencer el timeout informa tipos observados y ultimo payload relevante
- existe cobertura unitaria para matching, buffering y timeout

###### M3.P3c — Fundacion, mapa y tiempo — PENDIENTE

Done:

- una partida nueva completa el flujo minimo de fundacion
- el snapshot REST y todos los deltas pertenecen a la misma partida
- se observa la ceremonia contractual del mapa
- se observa avance de tiempo sin recibir estado completo por WebSocket

###### M3.P3d — Trade end-to-end — PENDIENTE

Done:

- el runner descubre IDs validos desde el snapshot
- propone y completa un trade por entradas publicas
- observa `trade.patch`, `roster.patch` y `finance.patch`
- un fallo identifica el eslabon roto sin consultar tablas internas

###### M3.P3e — Medicina end-to-end — PENDIENTE

Done:

- existe una lesion alcanzable de forma determinista por una entrada publica autorizada
- la decision del GM atraviesa gateway, NATS, team-service y agent-service
- se observan `roster.patch` y `relations.patch`
- un retry no duplica consecuencias

###### M3.P3f — Chat stub end-to-end — PENDIENTE

Done:

- se abre chat con un agente existente
- llega `chat.message` con correlacion correcta
- el historial REST contiene el intercambio
- la corrida no requiere credenciales de un proveedor LLM real

###### M3.P3g — Integracion y comando obligatorio — PENDIENTE

Done:

- `make smoke` recorre todos los escenarios anteriores en orden
- el resumen final muestra duracion y deltas observados
- el comando se agrega a los checks esperados de los mini milestones posteriores
- la documentacion incluye troubleshooting de fallos frecuentes

###### M3.P3h — Playoffs y cierre — BLOQUEADO POR M3.16+

Done futuro:

- el smoke atraviesa generacion de bracket y al menos una serie
- valida idempotencia al avanzar rondas
- alcanza `owner.veredicto_emitido` y comprueba el estado terminal cuando corresponda

Este corte no bloquea el cierre inicial de `M3.P3`; se activa cuando las mecanicas de playoffs y cierre existan.

##### Criterio de done de M3.P3

`M3.P3` queda realizado cuando:

- `make smoke` ejecuta fundacion, mapa/tiempo, trade, medicina y chat stub contra el entorno local
- cada escenario usa exclusivamente contratos publicos del gateway
- los mensajes se correlacionan por `game_id` y por IDs de dominio cuando corresponda
- todas las esperas tienen timeout y cancelacion
- los fallos muestran el paso y contrato roto
- el runner no borra ni reescribe datos ajenos
- sus tests unitarios, `go vet`, build y `git diff --check` pasan
- el comando queda incorporado al checklist habitual del repositorio

##### Riesgos a resolver durante la implementacion

- confirmar si `docker-compose` puede exponer un healthcheck confiable por servicio
- confirmar el flujo publico exacto para fundar una partida sin UI
- decidir una via determinista y no destructiva para disponer de una lesion
- evitar que la velocidad de simulacion vuelva flaky la espera de partidos
- garantizar que el observer no pierda deltas mientras un escenario ejecuta requests REST
- conservar aislamiento cuando dos smoke runs se ejecuten en paralelo

Estas decisiones se resuelven en el mini milestone que las necesita y se registran en este documento; no se inventan abstracciones anticipadamente.

##### Comandos de verificacion previstos

```bash
make smoke
go test ./... # dentro del modulo del runner, si usa modulo propio
go vet ./...
go build ./...
git diff --check
```

Los comandos exactos se actualizan cuando exista la estructura definitiva del runner.

##### Registro de avance

###### 2026-07-18 — Plan operativo consolidado

- se consolido el plan operativo completo dentro de este documento de pausa
- se preservo el alcance original
- se definieron frontera de caja negra, escenarios, diagnostico y mini milestones
- no se implemento codigo del smoke en este corte

##### Proximo paso

Iniciar `M3.P3a — Preflight y cliente del gateway`: primero auditar endpoints publicos, autenticacion guest, URLs locales y tooling existente; despues elegir la estructura minima del runner.

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
