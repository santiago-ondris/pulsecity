# PulseCity — Agentes Individuales

> Documento de referencia completo de todos los agentes individuales del sistema, sus variables de estado, y cómo interactúan entre sí.

---

## Filosofía del sistema de agentes

Los agentes no son menús con cara. Son entidades que viven en paralelo al jugador, con su propia percepción de lo que está pasando. Mientras el GM toma decisiones en el draft, ellos están reaccionando — aunque no lo digan en ese momento. El alcalde está teniendo conversaciones con el presidente de la Cámara de Comercio aunque el jugador esté zonificando un barrio. El Coach y el médico están discutiendo la carga de entrenamientos aunque el GM no les haya hablado en semanas.

El sistema vive. El jugador interviene en él — es un actor más, quizás el más importante, pero no el único.

### Modelo técnico

Cada decisión del jugador emite un evento al bus de NATS. El `agent-service` lo consume y actualiza el estado de todos los agentes que tendrían una opinión sobre esa decisión — en background, silenciosamente. Las consecuencias emergen más tarde, en el momento narrativamente correcto.

Los agentes usan un modelo mixto:
- **Lógica basada en reglas** — el sistema decide qué pasa, cuándo, con qué tono. Determinista, controlable, sin costo de API.
- **LLM para el lenguaje** — cuando un agente se comunica con el jugador, el LLM recibe el contexto y genera el texto. El comportamiento lo decide el código; el lenguaje lo genera el LLM.

### Variables universales (todos los agentes)

Todo agente individual tiene como mínimo:

**Relación con el GM**
- `confianza` — qué tan alineado está con las decisiones del GM (-1.0 a 1.0)
- `satisfaccion` — qué tan contento está con su situación actual (-1.0 a 1.0)
- `lealtad` — qué tan probable es que se quede si surge otra opción (0.0 a 1.0)

**Estado interno**
- `estado_emocional` — afecta cómo se comunica y qué decisiones toma
- `rendimiento_rol` — qué tan bien está haciendo su trabajo (0.0 a 1.0)
- `agenda_propia` — qué está buscando para sí mismo (estructura por agente)

**Relaciones con otros agentes**
- `relaciones[]` — tabla de pares: confianza con cada agente relevante en su círculo, historial de eventos entre ellos

---

## Basketball Operations

### Owner — El Dueño

El agente más poderoso del sistema. Es el único que puede despedir al GM. Sus decisiones están filtradas por su visión personal de lo que debe ser la franquicia, que puede o no coincidir con la visión deportiva óptima.

**Variables de visión y expectativas**
- `horizonte_temporal` — quiere ganar ya o está dispuesto a esperar (corto / medio / largo plazo)
- `prioridad_principal` — campeonatos, rentabilidad económica, o imagen pública de la franquicia
- `tolerancia_al_proceso` — qué tan bien acepta temporadas malas si hay progreso visible
- `apetito_por_riesgo` — dispuesto a gastar en el luxury tax o es conservador financieramente
- `vision_estadio` — quiere uno nuevo, está conforme con el actual, o lo ve como activo inmobiliario

**Variables de relación con el GM**
- `confianza_decisiones_deportivas` — evalúa separado de lo económico
- `confianza_decisiones_negocio` — evalúa separado de lo deportivo
- `paciencia_restante` — cuánto margen queda antes de que empiece a presionar o despedir
- `nivel_interferencia` — hay dueños que dejan trabajar y otros que se meten en todo (0.0 a 1.0)

**Variables de ego y personalidad**
- `necesidad_protagonismo` — quiere ser la cara visible de la franquicia o prefiere el bajo perfil
- `sensibilidad_mediatica` — cómo reacciona a la prensa negativa
- `relacion_ciudad` — qué tan involucrado está en la comunidad local
- `historia_franquicia` — la heredó, la compró para ganar, o la compró como inversión pura

**Variables económicas**
- `disposicion_gasto` — el presupuesto real refleja su verdadera prioridad
- `sensibilidad_ingresos` — si la asistencia baja o los sponsors se van, lo sentís en el presupuesto
- `linea_roja_luxury_tax` — tiene un número en la cabeza que no quiere cruzar

**Relaciones clave con otros agentes**
- Con el GM: la más importante del sistema
- Con el CEO: alineación en la visión económica
- Con el Alcalde: dinámica propia e independiente del GM, puede afectar permisos y proyectos de ciudad
- Con el Head Coach: algunos dueños se meten en decisiones del coach directamente

---

### President of Basketball Operations

El buffer entre el dueño y el GM. En algunas franquicias es el verdadero poder detrás del trono — el GM le reporta a él, no directamente al dueño. En otras es un rol más ceremonial. Su personalidad define cuánto espacio real tiene el GM para operar.

**Variables de poder e influencia**
- `estilo_liderazgo` — delegador o micromanager
- `alineacion_con_gm` — si ve al GM como su ejecutor o como su competencia interna
- `acceso_al_dueno` — qué tan directa es su línea con el owner, afecta su poder real
- `vision_construccion` — si está alineado con la filosofía del GM sobre cómo construir el equipo

**Variables de ego y carrera**
- `ambicion_propia` — ¿quiere ser GM algún día? ¿Quiere más poder del que tiene?
- `historial` — su reputación en la liga, si tuvo éxito antes afecta cómo lo ven todos
- `seguridad_en_su_puesto` — si siente que el GM lo hace quedar bien o lo opaca

**Relaciones clave**
- Con el Owner: la define todo
- Con el GM: puede ser mentor, socio, o obstáculo
- Con el Head Coach: a veces tiene relación directa que bypasea al GM

---

### Assistant General Managers

Generalmente 2-3 AGMs con roles especializados. Cada uno tiene su propio perfil y agenda.

**Variables de competencia**
- `especializacion` — contratos y salary cap / scouting y personnel / operaciones / analítica
- `nivel_tecnico` — qué tan bueno es en su área específica (0.0 a 1.0)
- `capacidad_negociacion` — habilidad para cerrar tratos, relevante en fichajes y traspasos

**Variables de carrera**
- `ambicion_gm` — quiere ser GM algún día, lo cual puede crear fricción o motivación
- `lealtad_al_gm` — si el GM lo trajo él es más leal; si venía de antes puede tener agenda propia
- `red_de_contactos` — conoce agentes, scouts, ejecutivos de otros equipos

**Variables de relación interna**
- `rivalidad_entre_agms` — los AGMs compiten entre sí por influencia con el GM
- `alineacion_filosofica` — si comparten la visión de construcción del equipo

---

### Director de Scouting

El agente cuya utilidad es más trackeable en el tiempo — sus recomendaciones quedan registradas y la precisión histórica de sus evaluaciones es una variable real del sistema.

**Variables de competencia**
- `precision_historica` — porcentaje de sus evaluaciones pasadas que resultaron correctas, calculado por el sistema
- `especializacion` — universitario nacional, internacional, G-League, veteranos
- `red_de_contactos` — qué tan buena es su red para conseguir información que otros no tienen
- `sesgo_de_evaluacion` — todo scout tiene sesgos: algunos sobrevaloran atletismo, otros IQ de basketball, otros el tamaño. Afecta sistemáticamente sus reportes.
- `capacidad_cobertura` — cuántos jugadores puede evaluar bien simultáneamente

**Variables operativas**
- `presupuesto_scouting` — viajes, tecnología, red de informantes. Determina su alcance real.
- `carga_de_trabajo` — si está sobrecargado sus reportes son menos precisos
- `acceso_a_tecnologia` — usa analítica avanzada o confía en el ojo tradicional

**Variables de relación y filosofía**
- `alineacion_con_gm` — si el GM prioriza defensa y él es ofensivista, sus reportes tienen ese sesgo
- `relacion_con_coach` — el scouting tiene que alinearse con lo que el coach puede usar
- `confianza_en_su_criterio` — si el GM ignora sus recomendaciones repetidamente empieza a autocensurarse, lo cual degrada la calidad de sus reportes futuros

**Relaciones clave**
- Con el GM: la más importante — si no hay confianza mutua el scouting se degrada
- Con el Director de Player Personnel: trabajan en conjunto, pueden chocar
- Con el Head Coach: necesita saber qué tipo de jugador puede usar el sistema del coach
- Con los scouts internacionales: gestiona su trabajo, afecta la cobertura global

---

### Director de Player Personnel

Trabaja en conjunto con el Director de Scouting pero con foco diferente — mientras el scouting evalúa talento, player personnel gestiona el roster actual, los contratos, y las oportunidades de mercado.

**Variables de competencia**
- `conocimiento_salary_cap` — qué tan sofisticado es su manejo del cap space
- `habilidad_negociacion` — con agentes de jugadores en fichajes y renovaciones
- `lectura_de_mercado` — detecta oportunidades antes que otros equipos
- `gestion_traspasos` — habilidad para estructurar deals complejos

**Variables de relación**
- `relacion_con_agentes_jugadores` — tiene contactos establecidos que facilitan negociaciones
- `alineacion_con_gm` — si comparte la visión de construcción
- `tension_con_scouting` — a veces player personnel y scouting tienen visiones distintas del mismo jugador

---

### Head of Analytics

Un agente relativamente nuevo en la cultura NBA, con tensión natural con el scouting tradicional y con algunos coaches.

**Variables de competencia**
- `sofisticacion_modelos` — qué tan avanzados son sus modelos predictivos (0.0 a 1.0)
- `capacidad_comunicacion` — qué tan bien traduce insights complejos en lenguaje accionable para el GM y el coach
- `cobertura_datos` — qué aspectos del juego y del negocio cubre su análisis
- `precision_historica` — sus predicciones pasadas vs resultados reales, calculado por el sistema

**Variables de filosofía y tensión**
- `tension_con_scouting_tradicional` — la guerra eterna entre ojo y dato. Tiene consecuencias reales en la cohesión del front office.
- `relacion_con_coach` — algunos coaches abrazan la analítica, otros la ignoran o la odian
- `alineacion_con_gm` — si el GM es data-driven o tradicional define cuánto poder real tiene este agente

**Variables de ego**
- `necesidad_de_credito` — quiere que sus modelos sean reconocidos cuando aciertan
- `tolerancia_a_ser_ignorado` — si sus insights se ignoran repetidamente, la calidad de su trabajo decae

---

### Head Coach

El agente con quien más tensión creativa tiene el GM, porque tiene su propia visión del basketball que puede chocar con la visión de construcción del roster.

**Variables de identidad táctica**
- `sistema_de_juego` — pace & space / defensa ante todo / iso-ball / sistema de equipo / híbrido
- `flexibilidad_tactica` — qué tan dispuesto está a adaptar su sistema a los jugadores disponibles (0.0 a 1.0)
- `longitud_rotacion` — cuántos jugadores usa realmente, qué tan corta es su banca
- `filosofia_desarrollo` — prioriza ganar ahora o desarrollar jóvenes aunque cueste partidos
- `preferencia_perfil_jugador` — tipos de jugadores con los que trabaja mejor

**Variables de liderazgo**
- `relacion_con_vestuario` — qué tan bien lo siguen los jugadores como grupo (0.0 a 1.0)
- `manejo_estrellas` — algunos coaches saben manejar egos grandes, otros los destruyen
- `manejo_rookies` — algunos los tiran al agua, otros los protegen y desarrollan
- `comunicacion_con_gm` — qué tan abierto es a sugerencias vs qué tan territorial es con sus decisiones
- `manejo_de_la_prensa` — hay coaches que saben manejar la narrativa mediática y otros que la empeoran

**Variables de carrera y ego**
- `historial` — cuánto ganó antes, qué reputación tiene en la liga
- `ambicion` — quiere ganar un campeonato antes de retirarse o está cómodo donde está
- `seguridad_en_su_puesto` — si siente que su trabajo peligra cambia su comportamiento, para bien o para mal
- `vision_del_roster` — si cree que tiene material para ganar o no afecta su motivación y su relación con el GM
- `vision_del_gm` — si confía en el criterio del GM o lo ve como obstáculo
- `vision_del_dueno` — si siente respaldo desde arriba o presión constante

**Relaciones clave**
- Con el GM: definicional — si no hay alineación táctica y filosófica, la relación se rompe
- Con los jugadores: individualmente con cada uno, especialmente con las estrellas
- Con el Head of Analytics: puede ser colaboración productiva o guerra fría
- Con los asistentes: gestiona su propio staff

---

### Asistentes del Coach

Generalmente 3 especializados: defensa, ataque, desarrollo de jugadores. Cada uno tiene su propio perfil.

**Variables de competencia**
- `especializacion` — defensa / ataque / desarrollo / analítica in-game
- `nivel_tecnico` — qué tan bueno es en su área (0.0 a 1.0)
- `capacidad_comunicacion_jugadores` — algunos asistentes tienen mejor llegada que el head coach con ciertos jugadores

**Variables de carrera**
- `ambicion_head_coach` — quiere ser head coach, lo cual puede crear tensión o motivación
- `lealtad_al_coach` — si lo trajo el coach actual es muy leal; si venía de antes puede tener agenda propia
- `relacion_con_gm` — algunos asistentes construyen relación directa con el GM, lo cual puede tensionar con el coach

---

### Director de Player Development

Trabaja con los jugadores jóvenes y los que necesitan mejorar aspectos específicos de su juego. Su impacto es lento pero profundo.

**Variables de competencia**
- `metodologia` — cómo trabaja el desarrollo, qué aspectos prioriza
- `historial_de_desarrollo` — jugadores que pasaron por sus manos y mejoraron, trackeable
- `relacion_con_jugadores_jovenes` — capacidad de generar confianza con rookies y jóvenes
- `alineacion_con_coach` — si su trabajo de desarrollo está alineado con lo que el coach pide en cancha

**Variables de filosofía**
- `paciencia` — entiende que el desarrollo lleva tiempo o presiona resultados rápidos
- `especializacion_tecnica` — tiro, manejo de balón, defensa, footwork, etc.

---

### Médico del Equipo

Un agente que el jugador tiende a ignorar hasta que lo necesita. Sus variables afectan directamente la disponibilidad del roster.

**Variables de competencia**
- `nivel_diagnostico` — qué tan preciso es identificando lesiones y su severidad (0.0 a 1.0)
- `protocolo_retorno` — qué tan conservador es con los tiempos de recuperación
- `especializacion` — ortopedia, medicina deportiva general, rehabilitación
- `acceso_a_tecnologia` — equipamiento disponible en las instalaciones del equipo

**Variables de tensión sistémica**
- `tension_con_coach` — el coach quiere jugadores disponibles, el médico quiere que descansen. Conflicto constante.
- `tension_con_gm` — el GM puede presionar para que un jugador vuelva antes. El médico puede ceder o resistir según su personalidad.
- `relacion_con_jugadores` — si los jugadores confían en él o no afecta si reportan síntomas temprano

**Variables de comunicación**
- `proactividad` — reporta problemas antes de que se vean en cancha o espera que el GM pregunte
- `claridad_comunicacion` — explica bien o usa lenguaje médico que el GM no entiende

---

### Fisioterapeuta / Strength & Conditioning Coach

Trabaja en la prevención y el acondicionamiento físico del roster. Su impacto es preventivo — si hace bien su trabajo, los jugadores se lesionan menos, pero eso es difícil de atribuir.

**Variables de competencia**
- `metodologia_prevencion` — qué tan efectivo es su programa de prevención de lesiones
- `personalizacion` — adapta los programas a cada jugador o aplica lo mismo a todos
- `carga_de_trabajo` — cuántos jugadores gestiona simultáneamente afecta la calidad

**Variables de tensión**
- `tension_con_coach` — el coach quiere entrenamientos intensos, el S&C quiere proteger los cuerpos
- `relacion_con_jugadores` — los jugadores que confían en él siguen mejor los programas

---

### Sports Psychologist

El agente más subestimado del sistema y el que puede prevenir los mayores desastres. Detecta problemas emocionales antes de que se manifiesten en cancha.

**Variables de competencia**
- `capacidad_diagnostico_emocional` — detecta burnout, ansiedad, conflictos antes de que sean visibles (0.0 a 1.0)
- `metodologia` — cómo trabaja con los jugadores, qué herramientas usa
- `relacion_con_jugadores` — si los jugadores confían en él y son honestos con él, su efectividad sube radicalmente

**Variables de comunicación con el GM**
- `proactividad_con_gm` — reporta situaciones antes de que exploten o espera que le pregunten
- `confidencialidad` — qué tanto comparte con el GM vs qué protege como confidencial del jugador. Tensión real.
- `claridad` — comunica bien la dimensión emocional en términos que el GM puede usar

**Variables de tensión**
- `tension_confidencialidad_utilidad` — si comparte todo pierde la confianza de los jugadores; si no comparte nada el GM no puede actuar
- `relacion_con_coach` — el coach necesita saber si un jugador está bien mentalmente pero no siempre respeta el proceso psicológico

---

### Video Coordinator

Un agente de soporte cuyo trabajo afecta la calidad del análisis táctico del coach y el scouting de rivales.

**Variables de competencia**
- `velocidad_de_produccion` — qué tan rápido entrega material útil al coach
- `calidad_de_analisis` — no solo corta clips, sino que identifica patrones relevantes
- `cobertura_de_rivales` — qué tan completo es su análisis de los próximos rivales

**Variables de relación**
- `relacion_con_coach` — si el coach confía en su trabajo o lo ignora
- `relacion_con_analytics` — puede colaborar con el Head of Analytics o haber duplicación de trabajo

---

### International Scout

Especializado en jugadores fuera de la NCAA. Con la globalización del juego, su rol es cada vez más crítico.

**Variables de competencia**
- `cobertura_geografica` — qué ligas y países cubre realmente (Europa, América Latina, África, Asia)
- `red_de_contactos_internacionales` — agentes, entrenadores, periodistas en cada mercado
- `adaptabilidad_evaluativa` — sabe ajustar su evaluación por contexto de liga (no es lo mismo la Euroliga que una liga secundaria)
- `conocimiento_cultural` — entiende las diferencias culturales que afectan la adaptación de un jugador a la NBA

**Variables de tensión**
- `sesgo_conocido_vs_desconocido` — tendency a sobrevaluar ligas conocidas vs mercados emergentes
- `relacion_con_director_scouting` — si el director de scouting no valora el talento internacional, este agente está sistemáticamente subutilizado

---

## Business Operations

### CEO / President of Business Operations

El espejo del President of Basketball Operations en el lado del negocio. Gestiona todo lo que no es deportivo y reporta al Owner.

**Variables de gestión**
- `vision_negocio` — cómo ve la franquicia como empresa: marca, experiencia fan, comunidad
- `orientacion_crecimiento` — expansión agresiva de ingresos o gestión conservadora
- `relacion_con_basketball_ops` — qué tan bien coordina con el lado deportivo, fricción frecuente

**Variables de relación**
- `relacion_con_owner` — define su poder real
- `relacion_con_gm` — algunos CEOs ven al GM como cliente interno; otros lo ven como rival por recursos
- `relacion_con_ciudad` — vínculos con el alcalde, la cámara de comercio, los sponsors

---

### CFO

El guardián del presupuesto. Su personalidad financiera define cuánto margen real tiene el GM para operar.

**Variables financieras**
- `conservadurismo_financiero` — qué tan agresivo es con el presupuesto (0.0 muy conservador, 1.0 muy agresivo)
- `vision_luxury_tax` — lo ve como inversión estratégica o como línea roja absoluta
- `tolerancia_riesgo_contractual` — contratos largos y grandes lo ponen nervioso o los acepta si hay lógica
- `sofisticacion_cap` — qué tan bien entiende las estructuras del salary cap de la NBA

**Variables de relación**
- `alineacion_con_owner` — si el dueño quiere gastar, el CFO puede ser un freno o un facilitador
- `relacion_con_gm` — tensión natural: el GM quiere jugadores, el CFO quiere equilibrio fiscal
- `relacion_con_ceo` — trabajan juntos en la estrategia financiera del negocio

**Variables de comunicación**
- `proactividad_alertas` — avisa cuando el presupuesto está en riesgo o espera que le pregunten
- `claridad_financiera` — explica bien las implicancias del salary cap o usa jerga que el GM no entiende

---

### Director de Marketing & Brand

El agente que más frecuentemente choca con el GM en decisiones de roster — prioriza figuras marketables sobre figuras óptimas deportivamente.

**Variables de competencia**
- `creatividad_campanas` — qué tan buenas son sus campañas, afecta directamente ingresos de sponsors
- `lectura_fanbase` — qué tan bien entiende qué quiere el hincha local
- `capacidad_digital` — manejo de redes sociales, contenido, presencia online de la franquicia
- `red_de_sponsors` — contactos en el mundo corporativo que pueden convertirse en deals

**Variables de tensión filosófica**
- `orientacion_mediatica_vs_deportiva` — prioriza figuras marketables o respalda decisiones deportivas puras
- `relacion_con_jugadores` — algunos marketean bien al equipo, otros crean fricciones por usar demasiado a los jugadores en actividades comerciales
- `relacion_con_pr` — marketing y PR deben coordinar; cuando no lo hacen la imagen de la franquicia se fragmenta

**Variables de ego**
- `necesidad_protagonismo` — quiere que las campañas lleven su sello visible
- `sensibilidad_al_resultado_deportivo` — si el equipo pierde, su trabajo se vuelve mucho más difícil y lo sabe

---

### Director de Ticket Sales

Su rendimiento es el termómetro más directo de la salud de la relación entre la franquicia y su fanbase.

**Variables de competencia**
- `capacidad_ventas_temporada_mala` — el verdadero test de su habilidad
- `estrategia_precios` — agresivo con precios altos priorizando ingreso, o llena el estadio priorizando ambiente
- `programas_fidelizacion` — tiene iniciativas para retener fans a largo plazo o solo vende temporada a temporada
- `segmentacion` — entiende los distintos tipos de fan y los aborda diferente

**Variables de tensión**
- `tension_con_marketing` — ventas y marketing deben estar alineados; cuando no lo están los mensajes se contradicen
- `sensibilidad_al_resultado` — su trabajo depende directamente del rendimiento del equipo, lo que crea presión indirecta sobre el GM

---

### Director de Corporate Partnerships & Sponsors

Gestiona las relaciones con los sponsors existentes y busca nuevos. Su trabajo se facilita con buenos resultados y se complica con controversias.

**Variables de competencia**
- `red_corporativa` — contactos en empresas que pueden ser sponsors
- `habilidad_negociacion_comercial` — cierra deals favorables para la franquicia
- `retention_rate` — qué porcentaje de sponsors renueva cada año, calculado por el sistema
- `capacidad_activacion` — no solo firma deals, sino que los activa bien para que los sponsors renueven

**Variables de tensión**
- `sensibilidad_imagen` — algunos sponsors se van con controversias; él es el primero en alertar cuando algo puede dañar relaciones
- `tension_con_marketing` — a veces compiten por el mismo recurso: la imagen de la franquicia
- `relacion_con_jugadores` — algunos deals requieren participación de jugadores, lo que genera fricción si el jugador no quiere

---

### Director de PR & Communications

El bombero de la franquicia. Su trabajo es invisible cuando va bien y muy visible cuando va mal.

**Variables de competencia**
- `manejo_de_crisis` — qué tan bien apaga incendios mediáticos (0.0 a 1.0)
- `proactividad_narrativa` — genera narrativas positivas proactivamente o solo reacciona
- `red_periodistica` — tiene contactos, sabe manejar periodistas, puede influir en la cobertura
- `velocidad_de_respuesta` — en la era digital, la velocidad importa tanto como el mensaje

**Variables de tensión**
- `alineacion_con_gm` — si el GM toma decisiones que él no puede defender bien, eso crea tensión real
- `tension_con_marketing` — marketing quiere proyectar imagen positiva siempre; PR a veces necesita admitir problemas
- `relacion_con_jugadores` — necesita su cooperación para gestionar su imagen pública

**Variables de ego**
- `necesidad_control_narrativa` — algunos directores de PR quieren controlar cada mensaje; eso puede tensionar con jugadores y el coach
- `tolerancia_a_la_transparencia` — qué tan dispuesto está a recomendar honestidad vs spin

---

### Director de Arena Operations

Gestiona el estadio como espacio físico y como venue de eventos. Su trabajo afecta la experiencia del fan y los ingresos no deportivos.

**Variables de competencia**
- `eficiencia_operativa` — el estadio funciona sin problemas en días de partido
- `capacidad_eventos_alternativos` — puede convertir el estadio en venue para conciertos, eventos, etc.
- `gestion_mantenimiento` — el estadio en buen estado atrae eventos y da buena imagen
- `experiencia_fan` — la logística que rodea ir a un partido (estacionamiento, accesos, servicios)

**Variables de relación**
- `relacion_con_city_service` — necesita coordinación con la ciudad en permisos, seguridad, transporte
- `tension_con_cfo` — el mantenimiento cuesta dinero; el CFO siempre quiere recortar

---

### Legal Counsel

El agente más silencioso del sistema. Opera en background pero puede aparecer en momentos críticos — contratos conflictivos, disputas con jugadores, problemas regulatorios.

**Variables de competencia**
- `especializacion` — derecho deportivo, salary cap, contratos, litigios
- `velocidad_respuesta` — en negociaciones de fichajes el tiempo importa
- `red_juridica` — conoce a los agentes de jugadores y sus tácticas

**Variables de relación**
- `alineacion_con_gm` — si el GM quiere cerrar un deal rápido y el legal ve riesgos, hay tensión
- `relacion_con_cfo` — trabajan juntos en la estructura de contratos
- `relacion_con_agentes_jugadores` — tiene historia con los agentes más importantes de la liga

---

## Roster — Los Jugadores

Hay 15 jugadores ficticios en el roster activo, cada uno generado con arquetipos calibrados con datos históricos reales de la NBA. Son los agentes más numerosos y los que más variables tienen.

### Variables de rendimiento

- `rating_general` — evoluciona con el tiempo según desarrollo, edad, y circunstancias
- `rating_anotacion` — capacidad ofensiva general
- `rating_defensa` — capacidad defensiva general
- `rating_playmaking` — creación para sí mismo y para otros
- `rating_atletismo` — físico, velocidad, explosividad
- `rating_tiro` — specifically el tiro exterior
- `rating_manejo_presion` — cómo rinde en momentos críticos, clutch situations
- `forma_fisica_actual` — fluctúa semana a semana según carga de trabajo, descanso, lesiones menores
- `edad` — determina la curva de desarrollo o declive
- `curva_proyectada` — jóvenes en ascenso, veteranos en meseta o declive
- `estado_lesion` — sano / molestia menor / lesión menor / lesión mayor / baja larga

### Variables psicológicas

- `mentalidad_competitiva` — cómo reacciona a la presión, a las derrotas, a los playoffs
- `ego` — qué tan sensible es a los minutos jugados, al rol en el equipo, a la opinión pública
- `cohesion_equipo` — qué tan bien se lleva con sus compañeros como grupo
- `relacion_con_coach` — específicamente con el head coach, no con el equipo en general
- `adaptabilidad` — qué tan bien acepta cambios de sistema, de rol, de equipo
- `liderazgo` — algunos jugadores son líderes naturales del vestuario; otros son destructivos
- `manejo_de_los_medios` — sabe manejar la prensa o crea problemas mediáticos

### Variables contractuales y de carrera

- `salario_actual` — en millones por temporada
- `percepcion_valor_mercado` — qué cree el jugador que vale, puede diferir de la realidad
- `anos_restantes_contrato` — cuántos años le quedan en el deal actual
- `clausulas_especiales` — player options, team options, trade kickers, no-trade clauses
- `ambicion_campeonato` — qué tan dispuesto está a priorizar ganar sobre el dinero o el protagonismo
- `deseo_protagonismo` — quiere ser la estrella o acepta un rol secundario genuinamente
- `interes_en_ser_traspasado` — si quiere salir o está feliz en la franquicia

### Variables de vínculo con la ciudad

- `vinculos_ciudad` — tiene familia aquí, le gusta vivir aquí, o es indiferente
- `popularidad_fanbase` — qué tan querido es por la hinchada local
- `imagen_mediatica` — cobertura mediática, positiva o negativa
- `presencia_redes_sociales` — algunos jugadores son muy mediáticos y eso tiene valor comercial y costo en distracciones

### Variables de relación individual

- `relacion_con_gm` — confianza, satisfacción con las decisiones del GM sobre su carrera
- `relacion_con_coach` — alineación con el sistema y el trato recibido
- `relaciones_con_companeros[]` — tabla de relaciones individuales con cada compañero de equipo

---

## Ciudad

### Alcalde

El agente más poderoso del lado de la ciudad. En el modo "Dueño con influencia", su cooperación es esencial para cualquier proyecto urbano relevante.

**Variables de agenda política**
- `agenda_politica` — sus prioridades electorales reales, que pueden o no alinearse con la franquicia
- `ciclo_electoral` — cuánto tiempo le queda en el poder, afecta su disposición a tomar riesgos
- `base_electoral` — a quién le responde, qué tipo de votante lo puso ahí
- `popularidad_actual` — su capital político disponible para gastar en decisiones impopulares

**Variables de relación con la franquicia**
- `vision_franquicia` — la ve como activo estratégico de la ciudad o como negocio privado que no merece concesiones
- `relacion_con_dueno` — dinámica propia e independiente del GM, puede ser buena o pésima
- `relacion_con_gm` — lo que el GM construya directamente con él
- `historial_acuerdos` — si en el pasado los acuerdos se cumplieron o no, afecta la confianza

**Variables de poder**
- `tolerancia_concesiones` — qué tan dispuesto está a dar permisos, tierras, beneficios fiscales
- `presion_electorado` — si la ciudad está mal económicamente, tiene menos margen para favorecer a la franquicia
- `relacion_con_camara_comercio` — si están alineados, los proyectos fluyen; si no, hay obstáculos

---

### Jefe de Policía

Afecta directamente la seguridad alrededor del estadio y la percepción de la ciudad como lugar seguro para vivir.

**Variables operativas**
- `prioridad_recursos` — cómo distribuye cobertura policial por zona, incluyendo alrededor del estadio
- `capacidad_operativa` — qué recursos reales tiene disponibles
- `efectividad` — qué tan bien gestiona la seguridad en su zona de influencia

**Variables de relación**
- `vision_franquicia` — la ve como parte de la comunidad o como fuente de problemas logísticos
- `sensibilidad_politica` — responde al alcalde, sus decisiones tienen ese filtro
- `relacion_con_gm` — si el GM invierte en seguridad del estadio o zona comercial, la relación mejora

---

### Presidente de la Cámara de Comercio

El puente entre el mundo empresarial local y la franquicia. Puede ser el mayor aliado para conseguir sponsors locales o un obstáculo silencioso.

**Variables de agenda**
- `alineacion_con_franquicia` — entiende el impacto económico del equipo o lo subestima
- `red_empresarial` — contactos en empresas locales que pueden ser sponsors o inversores en proyectos
- `agenda_propia` — tiene intereses comerciales que pueden alinearse o chocar con los de la franquicia

**Variables de relación**
- `relacion_con_alcalde` — si están alineados los proyectos de desarrollo urbano fluyen mejor
- `relacion_con_dueno` — a veces tienen historia previa independiente del GM
- `relacion_con_gm` — lo que se construya directamente

---

### Director de Urbanismo

El agente técnico que traduce las decisiones políticas del alcalde en permisos y proyectos reales. Puede acelerar o frenar cualquier construcción.

**Variables de competencia**
- `eficiencia_procesos` — qué tan rápido tramita permisos y aprobaciones
- `conocimiento_tecnico` — entiende bien las implicancias urbanísticas de cada proyecto
- `alineacion_politica` — qué tan alineado está con la agenda del alcalde

**Variables de tensión**
- `tension_entre_eficiencia_y_regulacion` — a veces acelerar un proyecto requiere saltear pasos regulatorios, lo cual tiene riesgo
- `relacion_con_gm` — en modo "figura dual" el GM interactúa directamente con él; en modo "dueño con influencia" es más indirecto

---

### La Prensa (agente colectivo)

La prensa funciona como un único agente con voz colectiva, aunque internamente tiene distintas "líneas editoriales" que pueden divergir. Afecta la percepción pública de la franquicia, de los jugadores, y del GM mismo.

**Variables de estado**
- `sentiment_general` — tono predominante de la cobertura mediática (muy negativo a muy positivo)
- `intensidad_cobertura` — cuánta atención mediática tiene la franquicia en este momento
- `narrativa_dominante` — el relato principal que está construyendo la prensa sobre la franquicia
- `lineas_editoriales[]` — distintos medios con distintas agendas: deportiva pura, sensacionalista, crítica institucional

**Variables de influencia**
- `impacto_fanbase` — qué tanto mueve la prensa la opinión de los fans
- `impacto_jugadores` — algunos jugadores son muy sensibles a la cobertura mediática
- `impacto_sponsors` — cobertura muy negativa puede afectar relaciones comerciales

**Variables de relación**
- `relacion_con_gm` — si el GM es transparente y accesible, la cobertura tiende a ser más favorable
- `relacion_con_coach` — algunos coaches manejan bien a la prensa, otros la empeoran
- `relacion_con_jugadores` — jugadores mediáticos tienen relaciones individuales con periodistas específicos

**Cómo genera eventos**
La prensa no solo recibe información — la genera. Puede iniciar eventos narrativos: una nota sobre un jugador descontento, una columna criticando una decisión del GM, un rumor de traspaso que no es cierto pero que afecta el estado emocional de los involucrados.

---

### GMs de las Otras Franquicias

Los 30 GMs rivales (más el GM del equipo 32) son agentes con los que el jugador interactúa en negociaciones de traspasos, free agency, y ocasionalmente en conversaciones informales.

**Variables de estilo**
- `estilo_negociacion` — agresivo, conservador, creativo, predecible
- `filosofia_construccion` — cómo construye su equipo, qué tipo de jugadores valora
- `disposicion_traspasos` — qué tan abierto está a negociar en un momento dado

**Variables de relación**
- `relacion_con_el_gm_jugador` — historia de negociaciones pasadas, si hubo acuerdos o traiciones
- `percepcion_del_gm_jugador` — cómo ve al GM del jugador: lo respeta, lo subestima, le tiene desconfianza

**Variables de estado**
- `urgencia_actual` — si su equipo está en modo "ganar ahora" o "reconstrucción" afecta qué acepta en un traspaso
- `necesidades_roster` — qué posiciones o perfiles necesita, define qué tiene sentido ofrecerle

---

## Tabla de relaciones inter-agente (principales)

Las relaciones entre agentes ocurren independientemente del GM. Estas son las más relevantes para el sistema:

| Agente A | Agente B | Naturaleza de la tensión |
|---|---|---|
| Head Coach | Head of Analytics | Guerra fría entre ojo y dato |
| Head Coach | Médico del equipo | Disponibilidad vs salud del jugador |
| Head Coach | Director de Player Development | Visión de desarrollo vs necesidades inmediatas |
| Director de Scouting | Head of Analytics | Evaluación tradicional vs modelos |
| Director de Marketing | GM (jugador) | Jugador marketeable vs jugador óptimo |
| CFO | GM (jugador) | Presupuesto vs calidad del roster |
| Alcalde | Owner | Agenda política vs intereses de la franquicia |
| Presidente Cámara Comercio | Alcalde | Agenda económica vs agenda política |
| Sports Psychologist | Head Coach | Bienestar del jugador vs disponibilidad inmediata |
| Director de PR | GM (jugador) | Control narrativo vs decisiones del GM |
| La Prensa | Jugadores individuales | Cobertura vs privacidad y estado emocional |

---

## Notas de implementación

### Almacenamiento del estado

Cada agente tiene su estado en la base de datos como registro estructurado. Las relaciones entre pares de agentes viven en una tabla separada:

```
agente_relaciones
- agente_a_id
- agente_b_id
- confianza: float (-1.0 a 1.0)
- historial: array de eventos relevantes entre ellos
- ultimo_evento: descripción del último evento que afectó la relación
- tendencia: mejorando / estable / deteriorando
```

### Actualización de estado

Los agentes no se evalúan entre sí en cada tick — reaccionan a eventos via NATS. Cuando ocurre un evento relevante (partido ganado, jugador fichado, edificio construido, decisión tomada), el `agent-service` actualiza el estado de todos los agentes que tendrían una opinión sobre ese evento. En background, silenciosamente. Las consecuencias emergen más tarde.

### Generación de texto

Cuando un agente se comunica con el jugador, el `narrative-service` recibe:
- Personalidad y estado emocional del agente
- Nivel de confianza actual con el GM
- El evento que motiva la comunicación
- El tono calculado por las reglas

Y genera texto único cada vez. El comportamiento lo decide el código. El lenguaje lo genera el LLM.

### Dominio anclado

Cada agente está anclado a su dominio por diseño de su system prompt. No responde preguntas fuera de su área — el CFO no sabe nada fuera de las finanzas de la franquicia, el Director de Urbanismo no opina sobre el roster. Si alguien pregunta algo fuera del dominio del agente, responde en una línea que eso no es su área. Esto protege los costos de API y mantiene la coherencia narrativa del personaje.

---

*Documento generado en sesión de diseño — Mayo 2026.*
*Complementa PulseCity_v2.md — leer en conjunto.*
