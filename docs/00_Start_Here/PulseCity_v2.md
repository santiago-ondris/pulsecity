# PulseCity — Documento de Diseño Completo

> Una ciudad ficticia que late al ritmo de su franquicia de NBA.

---

## ¿Qué es PulseCity?

PulseCity es una simulación interactiva browser-based de una ciudad ficticia generada proceduralmente, cuya economía, crecimiento urbano y vida social están directamente ligados al rendimiento de su franquicia de basketball ficticia en la NBA.

No es un juego en el sentido tradicional, ni una simple visualización. Es un sistema vivo: la ciudad respira, los ciudadanos se mueven, la economía fluctúa, y todo eso reacciona a las decisiones que tomás como General Manager de tu franquicia y a los resultados en la cancha.

Ganás un campeonato → el barrio alrededor del estadio se desarrolla, sube el valor del suelo, llegan nuevos ciudadanos. Temporadas malas → emigración, caída económica, presión del dueño. Construís transporte al estadio → más asistencia, más ingresos. Recortás salarios → fans descontentos, comercio en baja.

Los dos ejes inspiracionales y fundacionales del proyecto son **Cities Skylines 2** (city builder, gestión urbana, UI, emergencia de sistemas) y **NBA 2K** (profundidad de franquicia, negociaciones, gestión de roster). PulseCity es la fusión de ambos mundos.

---

## ¿Por qué existe este proyecto?

### Motivación técnica

PulseCity nace como un proyecto personal sin límites artificiales — sin "eso es overkill", sin "eso es demasiado complejo". El objetivo es familiarizarse con sistemas distribuidos reales, lenguajes nuevos (Rust, Go), orquestación con Kubernetes, generación procedural, y arquitecturas event-driven, todo dentro de un dominio narrativo coherente y visualmente impresionante.

### Motivación personal
Los intereses personales que convergen en este proyecto:

- **NBA y basketball** — pasión por el deporte, sus stats, su dinámica
- **Cities Skylines 2** — fascinación por ver sistemas urbanos interactuar, por el diseño y la emergencia
- **Simulación y sistemas vivos** — la satisfacción de ver un sistema evolucionar solo
- **Visualización interactiva** — algo que se pueda mostrar, explorar, y que cuente una historia visualmente

Siempre es mejor construir algo que genuinamente te interesa.

---

## El mundo de PulseCity

### La franquicia
PulseCity no usa equipos reales de la NBA. El jugador crea su propia franquicia ficticia — con nombre, ciudad, colores e identidad propios — como si fuera una **expansión 33 de la NBA**. Existe un segundo equipo de expansión ficticio (equipo 32), manejado por IA, que compite en la conferencia opuesta.

Los 30 equipos reales de la NBA existen como contexto y rivales simulados, calibrados con datos reales históricos de rendimiento.

Los jugadores del roster son **ficticios calibrados con arquetipos reales** — no hay licencia de jugadores reales (NBPA), pero cada jugador ficticio está generado con stats, físico y tendencias de declive basadas en datos históricos reales de la NBA. Esto preserva la agencia del GM: el scouting tiene sentido genuino y el meta del juego no es obvio.

### La ciudad
La ciudad es **ficticia y generada proceduralmente** usando algoritmos de **Perlin Noise** (terreno, elevación, geografía) y **Voronoi** (división en zonas y barrios). Cada partida genera una ciudad única. En el futuro, este módulo puede ser reemplazado por un loader de ciudades reales via OpenStreetMap sin cambiar el resto del sistema.

### Los agentes
La población de la ciudad usa un modelo **híbrido**:

- **Población agregada** — zonas con stats de habitantes, densidad, economía. No son entidades individuales, son stats que reaccionan a servicios y eventos. Sin pathfinding individual.
- **Agentes visibles** — emergen cuando hacen algo relevante (ciudadanos yendo al estadio en día de partido como grupo animado, inversores apareciendo en zonas de boom económico). Son sprites animados siguiendo paths precalculados, no entidades físicas independientes.
- **Agentes individuales con memoria** — los personajes únicos de la franquicia y la ciudad. Son aproximadamente 50 entidades: todo el staff de Basketball Operations, Business Operations, el roster de jugadores, y figuras clave de la ciudad. Cada uno tiene estado propio, personalidad, historia, relaciones con otros agentes, y es interactuable individualmente. Son el corazón narrativo del sistema.

### Los agentes individuales — estructura completa

**Basketball Operations (~15 agentes)**
- Owner (dueño de la franquicia) — tu jefe máximo, puede despedirte
- President of Basketball Operations — entre el dueño y el GM
- Assistant General Managers
- Director de Scouting
- Director de Player Personnel
- Head of Analytics
- Head Coach
- Asistentes del coach (defensa, ataque, desarrollo)
- Director de Player Development
- Médico del equipo
- Fisioterapeuta / Strength & Conditioning Coach
- Sports Psychologist
- Video Coordinator
- International Scout

**Business Operations (~8 agentes)**
- CEO / President of Business Operations
- CFO
- Director de Marketing & Brand
- Director de Ticket Sales
- Director de Corporate Partnerships & Sponsors
- Director de PR & Communications
- Director de Arena Operations
- Legal Counsel

**Roster (~15 agentes)**
- 15 jugadores ficticios con nombre, perfil físico, personalidad, contrato, y estado de relación con el GM

**Ciudad y liga (~8-10 agentes)**
- Alcalde
- Jefe de Policía
- Presidente de la Cámara de Comercio
- Director de Urbanismo
- La Prensa (agente colectivo con voz única pero líneas editoriales internas divergentes)
- GMs de las otras 30 franquicias + el GM del equipo 32 (agentes con quienes se negocia traspasos y free agency, con memoria de cada negociación)

Cada agente tiene en base de datos: atributos de personalidad, nivel de competencia, historial de decisiones, y una tabla de relaciones con otros agentes (nivel de confianza de -1 a 1, historial de eventos relevantes entre ellos).

### El sistema de agentes — cómo funcionan

Los agentes usan un modelo **mixto**:

- **Lógica basada en reglas (código)** — el sistema decide qué pasa, cuándo, sobre quién, y con qué tono. Si el nivel de confianza entre el Director de Scouting y el GM está en -0.8 y el GM pide un favor, la probabilidad de cooperación es baja. Esto es determinista, controlable, performante, y no tiene costo de API.
- **LLM para el lenguaje** — cuando un agente necesita comunicarse con el jugador, el sistema le pasa al LLM el contexto del agente (personalidad, relación actual con el GM, el evento en cuestión, el tono calculado por las reglas) y el LLM genera el texto del mensaje. Cada interacción se siente única y humana sin que el LLM tome decisiones de gameplay.

Este modelo tiene un costo de API negligible — aproximadamente 15 llamadas por semana simulada, menos de $0.05 por temporada completa.

### El tiempo

La unidad mínima de tiempo simulado es **un día**. Toda la lógica del sistema — partidos, eventos de agentes, efectos en la ciudad — se ancla a días simulados.

La simulación corre en **tiempo continuo acelerado** con velocidades ajustables (pausa / x1 / x5 / x20), igual que Cities Skylines. Una temporada completa (~365 días simulados) dura aproximadamente 10 minutos en tiempo real a velocidad x1.

| Velocidad | Tiempo real por temporada |
|---|---|
| x1 | ~10 minutos |
| x5 | ~2 minutos |
| x20 | ~30 segundos |

El loop de simulación es **híbrido**: corre continuamente cuando hay una sesión activa (conexión WebSocket abierta), y entra en sleep cuando no hay nadie conectado. El tiempo no avanza si el jugador no está — cuando vuelve, el mundo reanuda exactamente donde quedó.

### Los partidos
Los partidos se **simulan** — el motor calcula resultados usando modelos calibrados con stats históricas reales de la NBA. No se importan resultados reales, lo que preserva la agencia del GM: tus decisiones impactan el rendimiento del equipo.

---

## Flujo completo del usuario

### 1. Pantalla de entrada — Sesión

El jugador llega a PulseCity y ve tres opciones:

- **Iniciar sesión** — ya tiene cuenta
- **Registrarse** — cuenta nueva
- **Jugar como invitado** — sin registro, sin fricción. La partida vive en el backend asociada a un guest token almacenado en el browser. En cualquier momento desde adentro del juego puede crear una cuenta y la partida migra automáticamente.

Las partidas guardadas viven en el **backend** desde el día uno, asociadas a cuenta o guest token. Esto permite soporte multi-usuario nativo sin migración futura.

### 2. Menú principal

- **Nueva Partida**
- **Cargar Partida** — lista las saves asociadas a la cuenta o token del jugador

### 3. Creación de partida

La creación de partida tiene tres decisiones fundacionales:

#### 3a. Identidad visual de la franquicia

El momento donde el jugador siente que la franquicia es suya. Todos los cambios se reflejan en tiempo real en una preview mientras se diseña.

**Nombre**
- Nombre de la ciudad (también nombre del mapa generado)
- Nombre de la franquicia
- Abreviatura de 3 letras (estilo NBA: LAL, GSW, etc.)

**Colores**
- Color primario, secundario y de acento
- Se aplican en tiempo real a la UI mientras se eligen — el panel ya se tiñe con los colores propios

**Logo**
- Constructor de logo con capas: forma de fondo (escudo, círculo, hexágono, diamante, etc.), ícono central de una biblioteca curada (animales, símbolos, formas abstractas), tipografía del nombre de la franquicia
- Todo en tiempo real con los colores elegidos
- Opción de "subir tu propio logo" que reemplaza el constructor

**Uniformes**
- Tres variantes: Local, Visitante, y Alternativa (estilo City Edition de la NBA)
- Zonas del uniforme coloreables independientemente: cuerpo, mangas, pantalón, franja lateral
- Patrones aplicables: sólido, franjas horizontales, franjas verticales, degradé
- Número y apellido de ejemplo en tiempo real

#### 3b. Estado inicial de la franquicia

Define el punto de partida narrativo y operativo de toda la partida. Cada opción determina el estado del roster, el presupuesto, el perfil de la ciudad inicial, la personalidad del dueño, y la presión mediática.

**Reconstrucción**
- Roster joven, rating bajo, alto potencial
- Presupuesto ajustado
- Ciudad pequeña, economía modesta, fanbase chica pero leal
- Dueño paciente — da 3-4 temporadas antes de presionar
- Presión mediática baja

**Ventana de contención**
- 2-3 jugadores de élite, veteranos, contratos grandes
- Presupuesto alto pero comprometido en su mayoría
- Ciudad desarrollada, estadio lleno, fanbase caliente
- Dueño exigente — quiere playoffs este año, puede despedir rápido
- Presión mediática alta

**Franquicia histórica en declive**
- Mezcla rara de viejas glorias y juveniles sin dirección
- Presupuesto medio con deudas contractuales heredadas
- Ciudad grande, fanbase enorme pero frustrada, medios encima
- Dueño nostálgico — compara todo con el pasado glorioso, difícil de manejar
- Presión mediática muy alta

**Expansión pura**
- Draft de expansión — lo peor de cada equipo
- Presupuesto fresco, sin compromisos
- Ciudad chica y virgen, todo por construir
- Dueño visionario — horizonte largo, libertad total
- Presión mediática nula, nadie te conoce todavía

#### 3c. Modo de gestión de la ciudad

Define cómo el jugador se relaciona con la ciudad. Ambas opciones se presentan con su descripción exacta para que el jugador elija conscientemente.

**Dueño con influencia** — Sos el GM de la franquicia. Tu poder sobre la ciudad es indirecto: lobby, financiamiento de proyectos, presión política. El alcalde tiene su propia agenda y puede decirte que no. La ciudad se siente como un organismo independiente que reacciona a vos.

**Figura dual** — Controlás tanto la franquicia como la ciudad directamente. Dos sombreros, control total. La ciudad es tuya para construir como querés, y el basket convive con eso.

En ambos modos los servicios, edificios y efectos cruzados existen igual. Lo que cambia es si tenés acceso directo a los controles de la ciudad o si tenés que negociarlos con el alcalde.

### 4. Ceremonia de generación del mundo

Al confirmar la creación de partida, el backend comienza a generar el mapa. El proceso no es instantáneo — hay trabajo real de cómputo. En lugar de un spinner genérico, el mapa **se construye visualmente en tiempo real frente al jugador**:

1. Aparece el terreno con elevaciones (Perlin Noise)
2. Emergen las zonas Voronoi coloreándose
3. Se ubica el estadio
4. Aparecen los primeros barrios
5. La ciudad queda completa

El backend emite eventos por WebSocket conforme avanza ("terreno_listo", "zonas_calculadas", "estadio_ubicado") y el frontend los dibuja en secuencia. El jugador está viendo nacer su ciudad — el primer vínculo emocional con el mapa antes de tomar una sola decisión.

### 5. Primer evento narrativo — La llamada del dueño

Inmediatamente después de que la ciudad aparece completa, el dueño de la franquicia llama al GM (el jugador). Es el primer evento LLM del juego y es **obligatorio** — no se puede saltear. Dura poco, pero establece el tono de toda la relación.

El contenido de la llamada varía completamente según el estado inicial elegido:

- En **Reconstrucción**: el dueño es tranquilo, te da contexto de la situación, te marca el horizonte largo
- En **Ventana de contención**: el dueño es directo y exigente, te habla de expectativas de playoffs
- En **Franquicia histórica en declive**: el dueño está frustrado, te habla del pasado glorioso, te presiona con la historia
- En **Expansión pura**: el dueño es visionario y entusiasta, te da libertad total, habla del futuro

### 6. Pantalla principal — El juego

#### UI y layout

Inspirada directamente en Cities Skylines 2: **el mapa ocupa el 100% de la pantalla**. No hay paneles fijos laterales. Sobre el mapa, flotando, viven elementos mínimos siempre visibles:

- **Arriba izquierda** — widget de información general: fecha simulada, récord del equipo, temporada actual, presupuesto disponible. Se abre en panel completo al clickear, igual que CS2.
- **Abajo centro** — controles de tiempo: pausa / x1 / x5 / x20
- **Accesos a paneles** — botones discretos para abrir el Panel de GM y el Panel de Ciudad

Los paneles se abren sobre el mapa cuando el jugador los necesita. No roban espacio permanente.

#### Onboarding

- **Tutorial opcional** — el jugador decide si lo activa. Guía por los paneles principales paso a paso.
- **Wiki interna** — siempre accesible desde cualquier pantalla. Explica cada mecánica en profundidad. El jugador consulta cuando quiere, no antes.

La llamada del dueño funciona como tutorial disfrazado de narrativa — el jugador aprende el contexto del juego sin sentir que está en un tutorial.

---

## Loop de juego

El jugador opera en dos capas simultáneas:

### Como GM de la franquicia

**Gestión del roster**
- Fichajes, traspasos, extensiones de contrato, salary cap
- Negociaciones con agentes de jugadores
- Draft anual

**Gestión del staff**
- Contratación y despido de todo el staff de Basketball Ops y Business Ops
- Cada miembro del staff es un agente individual con personalidad y agenda propias

**Infraestructura del estadio**
- Capacidad de asientos
- Tecnología (pantallas, sonido, experiencia fan)
- Estacionamiento
- Instalaciones de práctica anexas

**Finanzas del equipo**
- Presupuesto operativo
- Patrocinadores y partnerships
- Salary cap management

**Relaciones con agentes individuales**
- Cada miembro del staff y cada jugador es interactuable individualmente
- Las relaciones tienen memoria y consecuencias a largo plazo

### Como autoridad de la ciudad

**Servicios urbanos** (construcción y gestión de presupuesto)
- Transporte: líneas de metro / tren ligero, paradas de autobús, autopistas y accesos
- Salud: hospitales, clínicas de barrio
- Seguridad: comisarías, cobertura policial por zona
- Educación: escuelas, universidades
- Espacios verdes y recreación: parques, canchas públicas de basket
- Zona comercial alrededor del estadio: bares, restaurantes, tiendas de merchandise

Los servicios básicos de infraestructura (electricidad, agua, basura) existen como stats de fondo que afectan la economía agregada, pero no se construyen manualmente. No agregan nada interesante al loop de PulseCity.

**Edificios especiales — one-time builds**

Construcciones únicas que definen la identidad de la ciudad y generan efectos permanentes con trade-offs reales:

| Edificio | Plus | Trade-off |
|---|---|---|
| Casino | Ingresos altos, turismo, valor del suelo | Mayor criminalidad, ahuyenta sponsors corporativos serios |
| Estadio de conciertos / arena multiuso | Ingresos extra, turismo, imagen | Compite con fechas del estadio principal, logística compleja |
| Universidad deportiva | Pipeline de talento local, cultura basket, atrae jugadores que quieren estudiar | Cara, tarda temporadas en dar frutos |
| Centro de convenciones | Turismo de negocios, sponsors corporativos, imagen profesional | Requiere buena infraestructura de transporte previa |
| Atracción turística icónica | Turismo, orgullo ciudadano | Muy cara, retorno lento |
| Distrito cultural | Atrae jugadores que quieren vivir en ciudad vibrante, mejor retención de estrellas | Gentrificación — sube valor del suelo, desplaza fans locales de bajos ingresos |
| Puerto deportivo / waterfront | Turismo premium, valor del suelo en la costa | Solo disponible si el mapa generado tiene agua |
| Hospital de élite | Jugadores top más dispuestos a firmar, atrae ciudadanos de alto poder adquisitivo | Requiere infraestructura de salud básica previa |
| Zona industrial / parque tecnológico | Empleos, economía robusta, más impuestos | Contamina zonas residenciales cercanas, baja calidad de vida |
| Academia de basket juvenil | En 5-6 temporadas genera prospectos locales con bonus "hijo de la ciudad" | Lenta, cara, devastadora si el jugador se va a otro equipo |
| Aeropuerto internacional | Conectividad total, turismo masivo, sponsors internacionales | Requiere mucho espacio, aleja zonas residenciales por ruido, muy cara |
| Centro de rehabilitación deportiva | Jugadores lesionados se recuperan más rápido, atrae veteranos | Cara, efecto invisible hasta que lo necesitás |
| Barrio histórico protegido | Turismo cultural, identidad única, orgullo ciudadano | Congela el desarrollo urbano en esa zona para siempre |
| Puerto comercial | Empleos masivos, economía industrial fuerte | Tráfico pesado, puede congestionar accesos al estadio |
| Complejo olímpico | Cada ciertos años podés postular a eventos internacionales con ingresos enormes | Construcción brutal, si no ganás la postulación el dinero se pierde igual |
| Zoológico / acuario | Turismo familiar, imagen amigable, mejora calidad de vida | Mantenimiento constante, primero en cerrar en crisis económica |
| Mercado central / distrito gastronómico | Vida urbana vibrante, atrae ciudadanos de distintos niveles económicos | Congestión peatonal y vehicular en la zona |
| Centro espacial / hub tecnológico | Imagen innovadora, atrae sponsors tech, ciudadanos de alto nivel educativo | Requiere universidad previa, impacto muy lento |
| Estadio alternativo / cancha de desarrollo | G-League o juveniles tienen donde jugar, ingresos extra | Compite por atención mediática local con la franquicia principal |
| Monumento a la franquicia | Orgullo ciudadano enorme, fanbase más leal, turismo nostálgico | Si se construye sin haber ganado nada, la prensa te destruye |
| Distrito de embajadas | Eventos internacionales, sponsors multinacionales, imagen global | Requiere seguridad elevada, costo operativo alto |
| Playa artificial / lago urbano | Calidad de vida disparada, valor del suelo en toda la zona | Solo tiene sentido si el mapa no tiene ya costa natural |

**Gestión económica urbana**
- Impuestos por zona
- Presupuesto municipal
- Valor del suelo (variable, afectado por servicios y por el rendimiento de la franquicia)

### Los efectos cruzados (el corazón del sistema)

| Decisión GM | Efecto en la ciudad |
|---|---|
| Fichás una estrella | Boom de turismo, sube valor de suelo cerca del estadio |
| Recortás salarios | Fans descontentos, cae comercio en el barrio |
| Mejorás infraestructura del estadio | Atrae eventos, genera empleo, nuevo barrio crece |
| Temporada mala | Emigración, bajan impuestos recaudados |
| Campeonato | Expansión urbana, nuevo presupuesto disponible |

| Decisión Ciudad | Efecto en el equipo |
|---|---|
| Construís transporte al estadio | Más asistencia, más ingresos |
| Zona residencial cerca del estadio | Base de fans local más grande |
| Crisis de servicios | Jugadores exigen irse, imagen del equipo cae |
| Expansión de la ciudad | Más impuestos, más presupuesto para el GM |
| Hospital de élite | Jugadores top más dispuestos a firmar |
| Zona insegura cerca del estadio | Menos asistencia nocturna, menos ingresos |
| Academia de basket juvenil (madura) | Prospectos locales en el draft con fanbase que los adora |
| Mala infraestructura médica general | Jugadores top no quieren vivir en la ciudad |
| Ciudad educada y sofisticada | Mejores sponsors, fanbase más comprometida |
| Canchas públicas de basket | Cultura de basket orgánica, fanbase más leal y duradera |

---

## El sistema de eventos narrativos

### Filosofía de diseño

Los eventos no son notificaciones genéricas. Son **causales** — cada evento tiene origen, contexto, y consecuencias que dependen de cómo el jugador responde. El sistema tiene memoria: una decisión tomada en la temporada 2 puede tener consecuencias en la temporada 5.

### Modo de notificaciones (elegible por el jugador)

**Interrupciones** — los eventos pausan el juego y aparecen en pantalla. El jugador los atiende en el momento.

**Bandeja de entrada** — los eventos se acumulan y el jugador los atiende cuando quiere.

Esta elección es de preferencia personal — ambos modos ven los mismos eventos. Sin embargo, algunos eventos son **obligatoriamente interrupciones** independientemente del modo elegido:

- Un jugador pide un traspaso
- El dueño convoca una reunión de emergencia
- Un agente da un ultimátum con deadline real (ej: 24 horas para cerrar contrato)
- Eventos de crisis que requieren decisión inmediata

### Tipos de eventos

**Eventos de Basketball Operations**
- Reportes de scouting de jugadores y prospectos
- Alertas del médico sobre estado físico de jugadores
- El Sports Psychologist detecta problemas emocionales en un jugador antes de que se noten en la cancha
- El Coach pide reunión para discutir el sistema de juego
- Un jugador pide extensión de contrato
- Un agente de jugador rival ofrece un traspaso
- Conflicto interno entre dos jugadores del roster
- Un jugador quiere irse — solicitud de traspaso formal
- El Director de Scouting trae información sobre un prospecto del draft
- El Analytics team detecta una tendencia preocupante en el rendimiento

**Eventos de Business Operations**
- El CFO alerta sobre exceso de gasto
- El Director de Marketing presiona para fichar una estrella mediática (puede no ser la mejor decisión deportiva)
- Un sponsor importante amenaza con irse si los resultados no mejoran
- Propuesta de nuevo sponsor que tiene condiciones específicas
- El Director de PR alerta sobre una crisis de imagen
- Oportunidades de partnership con empresas locales

**Eventos de la ciudad**
- El alcalde bloquea un proyecto de construcción
- El jefe de policía reporta aumento de criminalidad cerca del estadio
- El presidente de la Cámara de Comercio propone una alianza
- Crisis de servicios urbanos que impacta en el equipo
- Propuesta de expansión urbana que requiere aprobación

**Eventos aleatorios emergentes**
- Te cruzás a un empleado del estadio que te dice algo sobre un jugador
- Un periodista te pregunta algo comprometedor antes de un partido importante
- Un jugador aparece tarde a un entrenamiento sin explicación
- Rumores de que otro equipo está interesado en tu mejor jugador
- Un ex-jugador de la franquicia hace declaraciones públicas sobre la organización

### Consecuencias con memoria

El sistema recuerda. Ejemplos concretos:

- Si traicionás a un agente en una negociación, ese agente bloquea a sus jugadores de ir a tu franquicia por las próximas temporadas
- Si ignorás repetidamente las alertas del Sports Psychologist, un jugador colapsa emocionalmente en el peor momento
- Si el alcalde y el dueño acumulan conflictos, el alcalde puede bloquear sistemáticamente proyectos de expansión del estadio
- Si construís el Casino antes que el Centro de Convenciones, algunos sponsors corporativos serios nunca se interesan en tu franquicia

---

## La franquicia en la liga

### La liga simulada

PulseCity usa los **30 equipos reales de la NBA** como contexto y rivales, calibrados con datos históricos reales de rendimiento de cada franquicia. El jugador es el equipo 33 (expansión), compitiendo en la conferencia que eligió durante la creación de partida.

Existe un equipo 32 (segunda expansión), manejado por IA, que compite en la conferencia opuesta. Es el rival narrativo principal — otro GM expandiéndose al mismo tiempo, con sus propias decisiones y ciudad.

### Año de inicio

La simulación comienza en un año real (punto de partida: temporada 2025-26 como referencia), con los historiales reales de los 30 equipos hasta ese punto. Desde ahí la simulación diverge — los resultados futuros son propios del mundo de PulseCity, no del mundo real. Los datos de equipos y partidos históricos son públicos y gratuitos (balldontlie.io, Basketball Reference).

### Los jugadores

Ficticios, calibrados con arquetipos reales. Cada jugador tiene:
- Nombre y apellido propios
- Perfil físico (posición, altura, peso, atlétismo)
- Rating general y por habilidades específicas
- Personalidad (afecta su comportamiento como agente)
- Contrato con salario y años restantes
- Historia dentro de la franquicia
- Relación actual con el GM

---

## Arquitectura técnica

### Visión general

PulseCity es un sistema distribuido donde cada subsistema vive en su propio microservicio, orquestado con Kubernetes. El browser es una ventana a ese sistema — recibe estado en tiempo real via WebSocket y permite interacciones que disparan eventos en el backend.

```
Browser (React + WebGL)
        ↓ WebSocket + REST
   API Gateway (Go)
        ↓
   NATS Event Bus
        ↓
┌──────────────────────────────────────────────┐
│  map-service      city-service               │
│  (Rust)           (Go)                       │
│                                              │
│  agent-service    team-service               │
│  (Rust)           (Go)                       │
│                                              │
│  match-service    analytics-service          │
│  (Rust)           (Go + TimescaleDB)         │
│                                              │
│  narrative-service                           │
│  (Go + LLM API)                              │
└──────────────────────────────────────────────┘
```

### Microservicios

**`map-service` (Rust)**
Generación procedural del mapa usando Perlin Noise y Voronoi. Estado del mapa, zonas, topografía. Es el servicio más CPU-intensivo en la inicialización; Rust es la elección natural. Emite eventos por WebSocket durante la generación para la ceremonia visual.

**`city-service` (Go)**
Economía urbana, servicios públicos, zonificación, presupuesto municipal, edificios especiales, valor del suelo. Lógica de negocio compleja pero no CPU-intensiva — Go es ideal.

**`agent-service` (Rust)**
Motor de agentes híbridos. Población agregada como stats por zona. Agentes visibles como entidades con paths precalculados. Los ~50 agentes individuales con memoria, estado relacional y lógica basada en reglas. Performance crítica para la simulación — Rust.

**`team-service` (Go)**
Franquicia, roster, jugadores ficticios, finanzas del equipo, salary cap, staff. Coordinación y lógica de negocio — Go.

**`match-service` (Rust)**
Simulación de partidos. Corre modelos de rendimiento calibrados con datos reales. Puede simular cientos de partidos por temporada en segundos — Rust.

**`analytics-service` (Go + TimescaleDB)**
Ingesta y consulta de series temporales: valor del suelo por tick, asistencia por partido, ingresos por semana, rendimiento de jugadores por temporada. TimescaleDB es PostgreSQL con extensiones para series temporales.

**`narrative-service` (Go + LLM API)**
Servicio nuevo respecto al diseño original. Gestiona la generación de texto para los eventos narrativos. Recibe del `agent-service` la decisión de qué evento ocurre, con qué contexto y tono, y genera el texto visible al jugador via LLM. También gestiona la bandeja de entrada y el historial de eventos.

**`gateway` (Go)**
Punto de entrada único desde el browser. Maneja WebSocket, REST, autenticación (login, registro, guest tokens), routing a servicios internos.

**NATS (Event Bus)**
Message broker liviano y rápido. Todos los servicios publican y consumen eventos — "partido_terminado", "jugador_fichado", "zona_expandida", "evento_narrativo_generado". Es el pegamento del sistema. Elegido sobre Kafka por simplicidad operacional sin sacrificar el paradigma event-driven.

### Frontend

**React + WebGL (Three.js / custom renderer)**
Vista 2D top-down minimalista del mapa. Zonas coloreadas, agentes animados, estadio destacado, flujo de ciudadanos en días de partido. Diseño limpio y legible. El browser no simula nada — es un display inteligente que recibe deltas de estado por WebSocket y los dibuja.

Incluye: panel de creación de identidad visual (logo builder, uniforme builder), controles de tiempo estilo CS2, paneles de GM y ciudad que se abren sobre el mapa, sistema de notificaciones y bandeja de entrada de eventos narrativos, wiki interna, y panel de personajes — el directorio de todos los agentes individuales con su estado visible y canal de chat directo para consultas dentro de su dominio.

Aspiración futura: vista isométrica como segunda capa visual sin cambiar la lógica.

### Performance — consideraciones clave

El mayor riesgo de performance no es la cantidad de agentes ni la complejidad de la simulación — es **qué se manda por WebSocket y con qué frecuencia**.

Reglas de diseño:
- El backend manda **deltas**, no estado completo. Solo lo que cambió en cada tick.
- El mapa visual se actualiza con throttling inteligente — no más de una vez por segundo en velocidad x1, solo zonas que cambiaron.
- Los agentes visibles son sprites con paths precalculados, no entidades con física propia.
- Los ~50 agentes individuales no se evalúan entre sí en cada tick — reaccionan a eventos (event-driven, igual que el resto del sistema).

### Infraestructura

**Kubernetes (k3s local en desarrollo)**
Cada microservicio en su propio deployment. HPA para escalar `match-service` durante simulaciones intensivas. El objetivo es operar un cluster real, entender service discovery, config maps, secrets, y networking entre pods.

**Railway (deploy público)**
Deploy del frontend y servicios backend para demo pública. Cualquier persona puede abrir un link y ver PulseCity corriendo.

**GitHub Actions (CI/CD)**
Pipeline de build, test y deploy automático. Un push a main despliega.

**TimescaleDB**
PostgreSQL con extensión de series temporales. Una sola base de datos familiar con superpoderes para el volumen de datos que genera la simulación. Almacena también el estado de agentes y sus relaciones.

---

## Stack completo

| Componente | Tecnología | Razón |
|---|---|---|
| Motor de simulación | Rust | Performance, seguridad de memoria |
| Generación procedural | Rust | CPU-intensivo |
| Agentes | Rust | Performance en simulación continua |
| Servicios de negocio | Go | I/O, coordinación, simplicidad |
| API Gateway | Go | Concurrencia, HTTP, WebSocket |
| Servicio narrativo | Go + LLM API | Lógica en Go, lenguaje en LLM |
| Event Bus | NATS | Liviano, rápido, sin dependencias externas. PulseCity no necesita replay de eventos ni retención de mensajes — el estado vive en TimescaleDB. La complejidad operacional de Kafka no está justificada para este volumen. |
| Base de datos | TimescaleDB (PostgreSQL) | Series temporales + familiaridad |
| Frontend | React + WebGL | Browser-based, visualización en tiempo real |
| Orquestación | Kubernetes (k3s) | El punto del proyecto |
| CI/CD | GitHub Actions | Integración con GitHub, simple |
| Deploy | Railway | Ya en uso, soporta Docker containers |
| Datos NBA históricos | balldontlie.io + Kaggle datasets | Gratuitos, suficientes para calibrar modelos |
| LLM para narrativa | API externa (Claude / GPT-4o mini) | Texto emergente para eventos de agentes |

---

## Milestones de desarrollo

PulseCity no tiene un MVP en el sentido tradicional — no hay usuarios externos ni deadlines. El proyecto se construye en **milestones incrementales**: cada uno es jugable, mostrable en redes, y construye sobre el anterior sin tirar nada. El progreso es continuo y visible desde el primer día.

### Milestone 1 — El mundo nace
El jugador crea su franquicia, ve su ciudad generarse en tiempo real, y recibe la llamada del dueño. La primera experiencia emocional del sistema.

**Qué funciona:** autenticación (cuenta + guest token), creación de franquicia con identidad visual completa, ceremonia de generación del mapa via WebSocket, primer evento narrativo LLM obligatorio.

### Milestone 2 — El sistema late
El loop central funcionando. Partidos simulándose, tiempo avanzando, ciudad reaccionando a resultados. Los 5 agentes core activos con personalidad y estado real.

**Qué funciona:** loop de tiempo híbrido con todas las velocidades, 82 partidos simulados con box scores y momentos clave narrativos, `city-service` básico reaccionando a resultados, Owner / Head Coach / CFO / Director de Scouting / Sports Psychologist como agentes activos, eventos narrativos básicos post-partido.

### Milestone 3 — Los agentes viven
El corazón narrativo del sistema en su totalidad. Todos los agentes, todas las relaciones, causalidad real.

**Qué funciona:** los agentes individuales completos + roster de 15 jugadores, sistema de relaciones entre agentes (Coach vs Médico, Analytics vs Scouting, etc.), causalidad real con memoria de decisiones pasadas, canal de consulta directa con cada agente, ciclo jugable completo (Draft → Playoffs → Cierre).

### Milestone 4 — La ciudad respira
El loop ciudad-franquicia en toda su profundidad. Edificios especiales, alcalde con agenda propia, efectos cruzados completos.

**Qué funciona:** edificios especiales one-time con todos sus trade-offs, modo "Dueño con influencia" con el alcalde como agente real, todos los agentes de ciudad activos, analytics avanzados con visualizaciones históricas.

### Milestone 5 — Polish y experiencia final
La experiencia pulida como producto terminado. Convierte a PulseCity en una app madura, compartible públicamente, antes de abrir las features grandes del juego completo (multijugador, vista isométrica, ciudades reales).

**Qué funciona:** onboarding profundo + tutorial opcional contextual, wiki interna, microinteracciones y feedback visual fino, sonido y atmósfera, animaciones narrativas pulidas, accesibilidad (keyboard nav, contrast, screen readers), settings de usuario, performance percibido (skeletons, loading states, transiciones), estados vacíos/error pulidos, auditoría completa del design system aplicada a toda la app, test de usabilidad con personas externas.

**Filosofía:** cada milestone previo cierra con UX mínima viable para poder testear sus features. M5 toma la responsabilidad de la experiencia como un todo, lo que requiere ver la app completa. No incluye nuevas mecánicas de gameplay.

### Juego completo
Equipo rival (equipo 32) manejado por IA con su propia ciudad, liga completa con 32 equipos, modo multijugador, vista isométrica, ciudades reales via OpenStreetMap, sistema de legado.

> Para el detalle técnico de qué servicios están activos en cada milestone, ver **PulseCity_Arquitectura_Tecnica.md**.

---

## Escalabilidad del proyecto

PulseCity tiene **features infinitas**. La arquitectura event-driven y de microservicios permite agregar funcionalidad nueva sin reescribir lo existente. Cada feature nueva es un servicio nuevo o una extensión de uno existente, suscrito a los mismos eventos que ya corren.

Algunas direcciones posibles más allá del Milestone 4:
- Agente rival con IA real de toma de decisiones (equipo 32)
- Liga completa con 32 equipos simulados con sus propias ciudades
- Modo multijugador (dos GMs, dos ciudades, misma liga)
- Eventos históricos globales (recesión económica, desastre natural, pandemia)
- Editor de ciudad manual encima de la generación procedural
- Vista isométrica como segunda capa visual
- Ciudades reales via OpenStreetMap
- API pública para que otros devs construyan encima
- Mobile viewer (solo observación, sin gestión)
- Sistema de legado — varios GMs a lo largo de décadas en la misma franquicia

---

## Build in public

PulseCity se desarrolla abiertamente. El proceso de aprendizaje — implementar Perlin Noise en Rust por primera vez, ver la ciudad reaccionar a un resultado de partido por primera vez, configurar un service mesh real, integrar un LLM en un sistema event-driven — es parte del valor del proyecto.

Repositorio open source desde el día uno.

---

*Documento actualizado en sesión de diseño — Mayo 2026.*
*v1: Documento inicial de arquitectura y concepto.*
*v2: Flujo completo del usuario, sistema de agentes individuales, estructura real de franquicia NBA, edificios especiales, sistema de eventos narrativos, modelo mixto reglas + LLM.*
*v3: Modelo de tiempo simulado definido (1 día, loop híbrido, velocidades calibradas), state ownership entre servicios, modelo de partidos con inputs y outputs concretos, estructura de milestones reemplazando MVP. Ver PulseCity_Arquitectura_Tecnica.md para el detalle técnico completo.*
