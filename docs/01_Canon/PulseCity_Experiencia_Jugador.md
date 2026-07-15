# PulseCity — Experiencia del Jugador

> Documento de referencia de la capa de juego. Define qué hace el jugador concretamente, cómo está organizada la interfaz, y cómo transcurre una sesión. Complementa PulseCity_v2.md, PulseCity_Agentes.md, PulseCity_CicloJugable.md y PulseCity_Arquitectura_Tecnica.md — leer en conjunto.

---

## Filosofía de la experiencia

Los otros documentos describen el mundo — sus agentes, su arquitectura, su ciclo. Este documento describe al jugador dentro de ese mundo.

La pregunta central que responde es simple: abrís PulseCity, ¿qué hacés?

La respuesta no es "gestionás un equipo de basketball". Es algo más específico: **tomás decisiones, hablás con personas, y observás las consecuencias**. Esos tres verbos definen el loop completo. Todo lo demás — el roster, la ciudad, las finanzas, el salary cap — son el contexto en el que eso ocurre.

Dos tensiones definen la experiencia:

**Proactivo vs reactivo.** A veces el jugador va a buscar — quiere hablar con el Director de Scouting antes de tomar una decisión, quiere revisar el cap space antes de hacer una oferta. A veces el mundo lo interrumpe — el médico alerta, el dueño llama, un evento obliga a decidir. Ambos modos coexisten en todo momento y la interfaz los soporta sin fricciones.

**Inmersión vs eficiencia.** El jugador a veces quiere sentir el mundo — navegar el estadio, ver quién está en cada zona, charlar con un agente sin urgencia. A veces quiere ejecutar rápido — abrir el roster, hacer un movimiento, avanzar el tiempo. La interfaz no lo fuerza a elegir uno de los dos modos permanentemente. Los dos caminos llevan al mismo lugar.

---

## Las acciones del jugador — catálogo completo

Todo lo que el jugador puede hacer en PulseCity cabe en seis categorías. Dentro de cada categoría, las acciones son atómicas — cada una tiene un inicio, un flujo, y un resultado claro.

### 1. Gestión del roster

Las decisiones sobre los jugadores del equipo. Son las más frecuentes durante la temporada y las de mayor impacto narrativo en los agentes.

- **Iniciar negociación de traspaso** — el jugador propone un deal a un GM rival. Puede incluir jugadores, picks, y assets. El GM rival responde según su estilo, su urgencia actual, y su relación previa con el jugador.
- **Responder oferta de traspaso entrante** — un GM rival propone algo. El jugador acepta, rechaza, o contraoferta.
- **Hacer oferta a jugador libre** — durante la free agency, el jugador hace una oferta formal a un jugador disponible en el mercado. El jugador libre evalúa según su percepción de valor, su ambición de campeonato, y su vínculo con la ciudad.
- **Renovar contrato de jugador propio** — antes de que expire el contrato de un jugador, el jugador inicia o responde negociaciones de extensión. La satisfacción acumulada del jugador durante la temporada determina cuánto cuesta retenerlo.
- **Cortar jugador** — liberar a un jugador del roster con las implicancias de cap correspondientes.
- **Reclamar jugador en waivers** — cuando otro equipo corta a alguien, el jugador puede reclamarlo antes de que llegue al mercado libre.
- **Asignar jugador al G-League / traerlo de vuelta** — mover un jugador entre el roster activo y el desarrollo, con efectos en su satisfacción y en el trabajo del Director de Player Development.
- **Elegir pick en el draft** — la decisión más visible del offseason. Puede seguir la recomendación del Director de Scouting o ignorarla, con consecuencias en la relación y en la calidad futura de sus reportes.

### 2. Gestión del staff

Las decisiones sobre las personas que rodean al equipo. Menos frecuentes que las de roster, pero con consecuencias más profundas y lentas.

- **Contratar miembro del staff** — cuando hay una vacante o el jugador decide reemplazar a alguien, busca y contrata entre candidatos disponibles.
- **Despedir miembro del staff** — terminar la relación con un miembro del staff. Tiene efectos en la ansiedad general del equipo — los demás lo notan.
- **Renovar contrato de staff** — en el cierre de temporada, decidir quién sigue y en qué condiciones.

### 3. Recursos internos de la franquicia

Decisiones sobre el dinero y las operaciones internas del equipo que no son directamente roster o staff.

- **Asignar presupuesto a scouting** — cuánto dinero tiene el Director de Scouting para viajes, tecnología, y su red de informantes. Afecta directamente la cobertura y precisión de sus reportes.
- **Aprobar mejora del estadio** — capacidad, tecnología, instalaciones de práctica, experiencia fan. Cada mejora tiene costo y retorno en ingresos y satisfacción de jugadores.
- **Aprobar deal de sponsor** — el Director de Corporate Partnerships trae oportunidades. El jugador evalúa y decide.
- **Aprobar evento en el estadio** — el Director de Arena Operations propone eventos en fechas sin partido. Ingresos extra con trade-offs logísticos.

### 4. Gestión de la ciudad

Las decisiones urbanas que definen el contexto en el que vive la franquicia. Varían según el modo de ciudad elegido al inicio.

- **Construir servicio urbano** — transporte, salud, seguridad, educación, espacios verdes, zona comercial alrededor del estadio.
- **Construir edificio especial** — las construcciones únicas que definen la identidad de la ciudad. One-time builds con trade-offs permanentes.
- **Ajustar impuestos por zona** — disponible en modo Figura Dual. Afecta la economía agregada de cada barrio.
- **Negociar proyecto con el Alcalde** — disponible en modo Dueño con Influencia. El Alcalde tiene su propia agenda y puede decir que no. La relación acumulada con él determina qué tan receptivo es.

### 5. Interacción con agentes

El corazón narrativo de la experiencia. Lo que diferencia a PulseCity de un juego de gestión tradicional: cuando querés información o querés gestionar una relación, hablás con alguien.

- **Consultar agente proactivamente** — el jugador va al directorio o busca por nombre, abre el perfil del agente que necesita, y le pregunta lo que quiera dentro de su dominio. Sin que haya un evento que lo dispare. Esta es la acción más libre del sistema — podés hablar con el Sports Psychologist para entender el estado del vestuario, con el CFO para que te explique las implicancias de un contrato, con el Director de Scouting para que te dé su opinión sobre un prospecto específico antes del draft. La frecuencia de estas consultas afecta la relación — agentes consultados regularmente son más proactivos y abiertos.
- **Responder a un agente que te contactó** — el agente inició la conversación via un evento o notificación. El jugador lee el mensaje y responde. En algunos casos hay opciones predefinidas con consecuencias distintas; en otros es conversación abierta dentro del dominio del agente.
- **Tomar decisión en un evento obligatorio** — algunos eventos pausan el juego independientemente del modo de notificaciones elegido. El jugador no puede ignorarlos: un jugador que pide traspaso formal, el dueño que convoca reunión de emergencia, un agente que da ultimátum con deadline real.

### 6. Control del tiempo

La capa meta que envuelve todo lo demás.

- **Cambiar velocidad de simulación** — pausa / x1 / x5 / x20. El jugador controla el ritmo de la experiencia en cualquier momento.
- **Pausar / reanudar** — pausa manual cuando necesita tiempo para pensar o tomar decisiones sin presión.

---

## La interfaz — dos vistas principales

PulseCity tiene dos vistas principales entre las que el jugador alterna según lo que esté gestionando. No son tabs — son espacios distintos con identidades propias.

### Vista Ciudad

El mapa ocupa el 100% de la pantalla. La ciudad vive ahí — las zonas, el estadio, los barrios, los ciudadanos moviéndose en días de partido. Es la vista de Cities Skylines: el mundo físico, observable y modificable.

**Siempre visible:**
- Controles de tiempo (pausa / x1 / x5 / x20) — abajo al centro
- Botón de alternancia a Vista Franquicia
- Indicador de fecha simulada y récord del equipo — esquina superior

**Paneles que se abren sobre el mapa:**
- **Construcción** — servicios urbanos y edificios especiales disponibles, con sus costos y efectos descriptos
- **Economía** — valor del suelo por zona, ingresos municipales, presupuesto disponible
- **Ciudad** — stats agregadas de población, densidad, satisfacción por barrio

**Atajos contextuales a agentes:** cuando el jugador está mirando una zona específica, aparecen accesos directos a los agentes relevantes para esa zona. Mirando el estadio → acceso al Director de Arena Operations. Mirando una zona comercial en caída → acceso al Alcalde. Son sugerencias contextuales, no obligaciones — el jugador puede ignorarlos y navegar al agente desde la Vista Franquicia cuando quiera.

---

### Vista Franquicia

El estadio como espacio físico navegable. Es donde vive todo lo que no es gestión urbana — el roster, el staff, los agentes, las finanzas, la narrativa.

**Siempre visible:**
- **Header** — métricas críticas en una línea: récord del equipo, cap space disponible, próximo partido, alertas activas con indicador de urgencia
- **Buscador de agentes** — campo de texto siempre accesible. Escribís el nombre o el rol y llegás directo al agente sin navegar el espacio. Es el seguro contra la desorientación con 50 agentes en el sistema.
- **Toolbar lateral** — íconos de acceso rápido a los paneles más usados: roster, finanzas/salary cap, calendario de partidos, bandeja de eventos. Siempre visible independientemente de en qué zona del estadio estés.
- Controles de tiempo
- Botón de alternancia a Vista Ciudad

**El estadio como espacio navegable:**

El estadio está dividido en zonas. Cada zona tiene sus agentes. Cuando el jugador entra a una zona, aparecen los 3-4 agentes que pertenecen a ese espacio — con animación simple que comunica su estado emocional actual sin necesidad de abrir el perfil.

| Zona | Agentes presentes |
|---|---|
| La cancha | Head Coach, Asistentes del coach, Video Coordinator |
| Vestuario | Jugadores del roster (rotación activa) |
| Sala médica | Médico del equipo, Fisioterapeuta / S&C Coach, Sports Psychologist |
| Oficinas de Basketball Ops | Director de Scouting, Director de Player Personnel, Head of Analytics, Director de Player Development, International Scout |
| Oficinas del GM | President of Basketball Operations, Assistant General Managers |
| Palcos / Suite ejecutiva | Owner, CEO, CFO, Legal Counsel |
| Zona de prensa | Director de PR & Communications, La Prensa |
| Oficinas de Business Ops | Director de Marketing, Director de Ticket Sales, Director de Corporate Partnerships, Director de Arena Operations |
| Sala de negociaciones | GMs rivales aparecen acá cuando hay una negociación activa |

Desde cualquier agente en cualquier zona: abrir perfil completo (estado, métricas, historial de relación), iniciar consulta directa, o ver alertas activas que ese agente generó.

**Los paneles de acceso rápido** — accesibles desde la toolbar lateral sin importar en qué zona estés:

- **Roster** — los 15 jugadores del roster activo con su estado, contrato, minutos asignados, y estado emocional resumido. Desde acá se inician traspasos, cortes, asignaciones al G-League.
- **Finanzas / Salary Cap** — el cap space actual, los contratos de cada jugador, las excepciones disponibles, la proyección de los próximos años. El CFO puede explicar cualquier línea en lenguaje humano si el jugador lo consulta.
- **Calendario** — el calendario de partidos de la temporada. Próximos rivales, resultados recientes, fases del ciclo (cuándo es el trade deadline, cuándo abre la free agency).
- **Bandeja de eventos** — todos los eventos pendientes ordenados por urgencia. Los eventos obligatorios aparecen marcados y no desaparecen hasta que el jugador los atiende.

---

## El flujo de una sesión

PulseCity está diseñado para sesiones largas — el jugador que se sienta a jugar una hora, o el que quiere simular una temporada entera de un tirón. El sistema lo soporta sin fricciones.

### Entrada

El jugador abre PulseCity y el mundo está exactamente donde lo dejó. El tiempo no avanzó mientras no había sesión activa.

Lo primero que ve depende de en qué vista estaba cuando cerró. Si estaba en la Vista Ciudad, el mapa. Si estaba en la Vista Franquicia, el estadio.

El header le dice de entrada si hay algo urgente — alertas del médico, eventos obligatorios pendientes, un agente que necesita respuesta con deadline. Si hay algo urgente, el jugador decide si atenderlo primero o explorar un momento antes. Si no hay nada urgente, el mundo está tranquilo y el jugador arranca en modo libre.

### El loop central

El tiempo corre. El jugador alterna entre dos modos según el momento:

**Modo reactivo** — el mundo trae algo. Un evento en la bandeja, una alerta en el header, una notificación de que un partido terminó. El jugador pausa, atiende, decide, y retoma.

**Modo proactivo** — el jugador quiere hacer algo. Va a buscar información, inicia una negociación, consulta a un agente, construye algo en la ciudad, revisa el roster. El tiempo puede seguir corriendo o puede pausarlo mientras gestiona.

Los dos modos se mezclan naturalmente. El jugador puede estar en medio de una negociación de traspaso (proactivo) cuando llega una alerta del Sports Psychologist (reactivo). Pausa, atiende la alerta, retoma la negociación.

### Los momentos de pausa natural

El juego genera sus propios ritmos. No todo es urgente todo el tiempo — hay períodos de la temporada más tranquilos (mitad de la temporada regular sin rachas ni lesiones) y períodos de alta intensidad (trade deadline, playoffs, cierre de temporada). El jugador siente esa variación sin que el sistema se lo explique.

Los partidos son el marcador de ritmo más constante — cada partido que termina es un momento natural de pausa. Ver el box score, leer el evento narrativo post-partido, evaluar el estado del vestuario. Después, reanudar.

Las fases del ciclo generan pausas automáticas con UI expandida — Draft, Free Agency, Trade Deadline, Cierre de Temporada. El juego lleva al jugador a ese estado sin que tenga que acordarse de pausar. La UI se adapta para mostrar exactamente lo que se necesita en ese momento.

### La sesión larga — cómo se siente

Una sesión típica de una hora a velocidad x1 cubre aproximadamente cuatro semanas simuladas — ocho o diez partidos, los eventos narrativos que generan, las reacciones de los agentes, los movimientos de mercado si los hay.

A velocidad x5, esa misma hora cubre una temporada completa. El jugador que quiere simular rápido y parar solo en lo importante puede hacer eso — subir la velocidad, pausar cuando el header indica algo urgente, atender y retomar.

No hay una manera correcta de jugar una sesión. El sistema soporta al jugador que quiere controlar cada detalle a x1 y al que quiere ver resultados rápido a x20. La velocidad es una decisión personal que puede cambiar en cualquier momento.

### Salida

El jugador pausa o simplemente cierra. El mundo queda exactamente donde está. Cuando vuelva, todo sigue desde ese punto — los agentes con su estado acumulado, la ciudad con su economía, la temporada en el día exacto donde quedó.

---

## La consulta proactiva — por qué es central

Vale la pena dedicarle una sección propia porque es la acción que más define la identidad de PulseCity como experiencia.

En un juego de gestión tradicional, cuando necesitás información sobre algo, navegás un menú. Abrís el panel de finanzas, buscás la línea que necesitás, cerrás el panel. La información está ahí pero el mundo se siente como una base de datos.

En PulseCity, cuando necesitás información sobre algo, hablás con alguien. Querés entender el estado del vestuario antes de los playoffs — hablás con el Sports Psychologist. Querés saber si un prospecto del draft vale la pena — hablás con el Director de Scouting. Querés entender las implicancias de un contrato largo — hablás con el CFO. La información llega en lenguaje humano, con el contexto y la personalidad del agente que la da.

Esa diferencia hace que el mundo se sienta habitado, no administrado.

Y tiene una consecuencia de diseño importante: **la consulta proactiva tiene que ser siempre accesible**. No escondida en un panel, no requiriendo varios pasos de navegación. El buscador de agentes en la Vista Franquicia existe exactamente para eso — en cualquier momento, desde cualquier zona, el jugador puede escribir "marketing" y estar hablando con el Director de Marketing en dos segundos.

La accesibilidad no es un detalle de UX — es parte de la filosofía del juego. Si consultar a un agente requiere fricción, el jugador lo hace menos. Y si lo hace menos, el mundo se siente menos vivo.

---

## Tabla de acciones por vista

| Acción | Vista Ciudad | Vista Franquicia |
|---|---|---|
| Construir servicio urbano | ✓ (panel Construcción) | — |
| Construir edificio especial | ✓ (panel Construcción) | — |
| Ajustar impuestos | ✓ (panel Economía) | — |
| Ver stats de ciudad | ✓ (panel Ciudad) | — |
| Gestionar roster | — | ✓ (panel Roster) |
| Gestionar finanzas / cap | — | ✓ (panel Finanzas) |
| Ver calendario | — | ✓ (panel Calendario) |
| Atender bandeja de eventos | — | ✓ (panel Bandeja) |
| Consultar agente proactivamente | Atajo contextual | ✓ (buscador + zonas del estadio) |
| Responder evento de agente | — | ✓ (panel Bandeja) |
| Controlar velocidad del tiempo | ✓ | ✓ |
| Negociar con el Alcalde | Atajo contextual | ✓ (zona de negociaciones) |

---

## Notas de diseño

### El buscador como seguro

Con ~50 agentes en el sistema, la navegación espacial por zonas del estadio es la experiencia rica pero no puede ser el único camino. El buscador de agentes existe para garantizar que el jugador nunca esté perdido — siempre puede llegar a quien necesita en dos segundos sin saber en qué zona del estadio vive ese agente.

### Los atajos contextuales en la Vista Ciudad

Los atajos a agentes en la Vista Ciudad no son acceso completo — abren una versión reducida del chat con ese agente, anclada al contexto de lo que el jugador estaba mirando. "Estoy mirando esta zona comercial que está bajando, hablame de eso." El agente responde desde ese contexto. Para una conversación más amplia, el jugador va a la Vista Franquicia.

### La toolbar lateral no reemplaza las zonas

Los paneles de acceso rápido (roster, finanzas, calendario, bandeja) son para ejecución. Las zonas del estadio son para inmersión y para llegar a los agentes. No son redundantes — son dos modos de uso distintos que coexisten. El jugador que quiere hacer un traspaso rápido usa la toolbar. El jugador que quiere entender qué está pasando con su equipo antes de tomar esa decisión navega el estadio y habla con quien necesita.

### Eventos obligatorios — por qué interrumpen siempre

Los eventos obligatorios (jugador que pide traspaso, dueño que convoca reunión de emergencia, ultimátum con deadline) interrumpen el tiempo independientemente del modo de notificaciones elegido porque tienen consecuencias que el jugador no puede ignorar. No son interrupciones arbitrarias — son el mundo diciéndole que algo importante necesita su atención ahora. La frecuencia es baja por diseño; si todo fuera obligatorio, el modo de notificaciones no tendría sentido.

---

*Documento generado en sesión de diseño — Mayo 2026.*
*Complementa PulseCity_v2.md, PulseCity_Agentes.md, PulseCity_CicloJugable.md y PulseCity_Arquitectura_Tecnica.md — leer en conjunto.*
