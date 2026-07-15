# PulseCity — Mecánicas de Agentes

> Documento de referencia mecánico. Define qué hace cada agente en el sistema — qué lo activa, qué acciones directas puede tomar, y qué efectos sistémicos genera. Complementa PulseCity_Agentes.md que define quiénes son. Este documento define qué hacen.

---

## Filosofía mecánica

Cada agente en PulseCity existe en dos capas simultáneas.

La primera es narrativa — su personalidad, su historia, sus relaciones. Eso está en PulseCity_Agentes.md. La segunda es mecánica — qué cambia en el sistema cuando ese agente actúa. Este documento es sobre la segunda capa.

La distinción clave que estructura cada agente acá es entre **acciones directas** y **efectos sistémicos**:

**Acciones directas** — el agente hace algo que cambia el estado del juego de forma concreta e inmediata. El Owner sube el presupuesto. El CFO bloquea un contrato. El médico declara a un jugador no disponible. El jugador lo ve y lo siente de inmediato.

**Efectos sistémicos** — el agente modifica variables que afectan otras cosas, cuyos resultados emergen más adelante. La Prensa sube la presión sobre el Owner. El Sports Psychologist mejora la cohesión del vestuario. El Director de Scouting degrada la calidad de sus reportes si se ignora su criterio repetidamente. El jugador no lo ve como una acción puntual — lo siente como el mundo evolucionando.

Ambos tipos de output son mecánicamente reales. La diferencia es la velocidad y la visibilidad con la que se manifiestan.

---

## Basketball Operations

### Owner

El agente más poderoso del sistema. Es el único que puede despedir al GM y el único que puede expandir o contraer el presupuesto operativo real de la franquicia. Todo lo que hace tiene peso inmediato.

**Qué lo activa**
- Resultados deportivos por debajo de sus expectativas (definidas por su `horizonte_temporal` y `prioridad_principal`)
- Pérdidas económicas sostenidas — asistencia baja, sponsors que se van, presupuesto en rojo
- Decisiones del GM que contradicen su visión explícita de la franquicia
- Cobertura mediática muy negativa sostenida (la Prensa lo afecta directamente via su `sensibilidad_mediatica`)
- Conversaciones directas con el Alcalde que generan fricciones con la franquicia
- Cierre de temporada — la evaluación formal es siempre obligatoria

**Acciones directas**
- **Subir el presupuesto operativo** — cuando los resultados superan sus expectativas o hay un campeonato. El GM tiene más margen real para operar.
- **Bajar el presupuesto operativo** — cuando la situación económica es mala o perdió confianza en el GM. Menos margen para fichajes, staff, y operaciones.
- **Fijar una línea roja de luxury tax** — comunica un número que no quiere cruzar. El CFO la implementa. Si el GM la cruza, la `paciencia_restante` del Owner baja significativamente.
- **Presionar al GM con condiciones explícitas** — en el cierre de temporada puede imponer métricas concretas para la siguiente temporada (playoffs mínimo, reducción de gasto, fichar una estrella). El jugador sabe exactamente qué tiene que lograr.
- **Dar ultimátum formal** — una temporada más para mostrar resultados específicos. Si no se cumplen, el despido es automático.
- **Despedir al GM** — cuando `paciencia_restante` llega a cero. Game over para esta partida.
- **Interferir en decisiones de roster** — si su `nivel_interferencia` es alto, puede presionar para fichar o no fichar a alguien específico. No es una sugerencia — tiene peso real en la relación si se ignora.
- **Contactar al Alcalde directamente** — independientemente del GM. Puede generar acuerdos o fricciones que el GM descubre después.

**Efectos sistémicos**
- Su nivel de satisfacción afecta el tono de toda la organización — cuando el Owner está contento, los agentes del front office están más relajados. Cuando está presionando, la ansiedad general sube.
- Su relación con el Alcalde puede facilitar o bloquear proyectos urbanos independientemente de lo que el GM gestione con él.
- Su visión de la franquicia filtra cómo evalúa cada decisión del GM — dos GMs que toman la misma decisión pueden recibir reacciones completamente distintas según el perfil del Owner que les tocó.

---

### President of Basketball Operations

El buffer entre el Owner y el GM. Su poder real depende de su `estilo_liderazgo` y su `acceso_al_dueno`. En algunas partidas es un aliado; en otras es un obstáculo silencioso.

**Qué lo activa**
- Decisiones del GM que considera fuera de la visión de construcción acordada
- Presión del Owner que él filtra antes de que llegue al GM (o que decide no filtrar)
- Su propia `ambicion_propia` — si quiere más poder del que tiene, busca momentos para demostrarlo
- Señales de que el GM lo opaca o lo bypasea

**Acciones directas**
- **Bloquear o demorar una decisión del GM** — si su estilo es micromanager y no está alineado, puede frenar un fichaje o un traspaso antes de que llegue al Owner para aprobación.
- **Escalar al Owner** — llevar una decisión del GM directamente al Owner, con o sin el conocimiento del GM. Puede ser a favor o en contra.
- **Dar respaldo explícito al GM** — cuando está alineado, puede defender decisiones del GM ante el Owner, aumentando la `paciencia_restante` disponible.
- **Proponer cambios en el front office** — si su `ambicion_propia` es alta y siente que tiene oportunidad, puede sugerir al Owner reestructuraciones que afectan el poder del GM.

**Efectos sistémicos**
- Su alineación con el GM afecta la fluidez de toda la operación de Basketball Ops. Cuando están alineados, las decisiones fluyen. Cuando no, hay fricción invisible que se manifiesta en demoras y malentendidos.
- Si tiene relación directa con el Coach que bypasea al GM, puede crear una línea de comunicación paralela que el GM no controla.

---

### Assistant General Managers

Generalmente 2-3 AGMs con especializaciones distintas. Sus outputs mecánicos dependen de su área.

**Qué los activa**
- Tareas asignadas por el GM dentro de su especialización
- Su `ambicion_gm` — si quieren ser GM algún día, buscan momentos de visibilidad
- Rivalidad entre AGMs por influencia con el GM
- Ser ignorados sistemáticamente en su área de expertise

**Acciones directas**
- **Gestionar negociaciones delegadas** — si el GM les da autonomía, cierran deals dentro de los parámetros acordados. Su `habilidad_negociacion` determina qué tan buenos son los términos que consiguen.
- **Alertar sobre riesgos en su área** — el AGM de salary cap alerta si un deal propuesto tiene implicancias que el GM no vio. El de scouting alerta si un jugador en negociación tiene red flags. Si se los ignora repetidamente, dejan de alertar.
- **Filtrar información al GM** — según su lealtad y su agenda, pueden presentar información de forma que favorezca su posición. No mienten, pero pueden enfatizar lo que les conviene.

**Efectos sistémicos**
- Su `nivel_tecnico` afecta directamente la calidad de la ejecución en su área. Un AGM de salary cap mediocre cierra contratos con estructuras subóptimas que el GM descubre años después.
- La rivalidad entre AGMs puede crear silos de información — uno no comparte con el otro, y el GM recibe visiones fragmentadas de la misma situación.

---

### Director de Scouting

El agente cuyo output es el más trackeable del sistema — sus recomendaciones quedan registradas y la precisión histórica es una variable real.

**Qué lo activa**
- Solicitudes del GM de evaluar jugadores o prospectos específicos
- El draft se acerca — aumenta su actividad proactiva
- Su presupuesto sube o baja — afecta su cobertura posible
- El GM sigue o ignora sus recomendaciones — afecta su `confianza_en_su_criterio`

**Acciones directas**
- **Emitir reporte de evaluación** — su output principal. Incluye su recomendación y su nivel de confianza. La calidad del reporte depende de su `precision_historica`, su `carga_de_trabajo` actual, y el presupuesto disponible para investigar.
- **Recomendar o desaconsejar un pick de draft** — en el día del draft, su recomendación está disponible para cada jugador en el board. Si el GM lo ignoró repetidamente antes, el reporte puede estar autocensurado — menos opinión, más datos crudos.
- **Alertar sobre un jugador disponible en waivers o mercado** — proactivamente, si detecta una oportunidad que encaja con las necesidades del roster.
- **Degradar sus reportes futuros** — si el GM ignora sus recomendaciones de forma sistemática, empieza a autocensurarse. Los reportes se vuelven más neutros, menos útiles. No es una decisión consciente del agente — es una consecuencia acumulada que el sistema calcula.

**Efectos sistémicos**
- Su `sesgo_de_evaluacion` afecta sistemáticamente todos sus reportes. Un scout que sobrevalora el atletismo va a recomendar jugadores atléticos aunque no encajen en el sistema del Coach. El GM que no conoce el sesgo de su scout toma decisiones basadas en información torcida.
- Su relación con el Head of Analytics define si el scouting y la analítica se complementan o se contradicen. Cuando chocan, el GM recibe señales opuestas sobre el mismo jugador.
- Su `presupuesto_scouting` limita su cobertura real. Con presupuesto bajo, algunos mercados (internacional, G-League) quedan sin cubrir y el GM simplemente no recibe información de esas zonas.

---

### Director de Player Personnel

Foco en el roster actual, los contratos, y las oportunidades de mercado. Trabaja en conjunto con Scouting pero con orientación distinta.

**Qué lo activa**
- Negociaciones de traspaso iniciadas por el GM o por GMs rivales
- Jugadores propios con contratos próximos a expirar
- Oportunidades en el mercado que detecta antes que otros equipos
- Su `alineacion_con_gm` — si comparte la visión, es proactivo; si no, ejecuta sin iniciativa propia

**Acciones directas**
- **Gestionar negociación con agente de jugador** — en renovaciones y fichajes, conduce la negociación con el agente del jugador. Su `habilidad_negociacion` y su `relacion_con_agentes_jugadores` determinan los términos que puede conseguir.
- **Estructurar deals de traspaso** — cuando el GM define qué jugador quiere o quiere mover, él estructura el deal concreto con el GM rival. Su `gestion_traspasos` determina la creatividad y viabilidad de las propuestas.
- **Alertar sobre ventanas de mercado** — detecta oportunidades antes de que sean obvias. Si se actúa sobre sus alertas, su satisfacción sube y su proactividad aumenta.
- **Alertar sobre riesgo de perder un jugador propio** — cuando un jugador con contrato expirando está siendo coqueteado por otro equipo, alerta con anticipación suficiente para actuar.

**Efectos sistémicos**
- Su `tension_con_scouting` puede generar evaluaciones contradictorias del mismo jugador. Player Personnel lo ve como activo de roster; Scouting lo ve como prospecto a largo plazo. El GM decide cuál perspectiva pesa más.
- Su red de contactos con agentes de jugadores es un activo invisible — facilita negociaciones que sin esa red serían más lentas y costosas.

---

### Head of Analytics

El agente que más tensión creativa genera con el sistema tradicional — su valor es real pero su influencia depende de cuánto poder le da el GM.

**Qué lo activa**
- Solicitudes del GM de análisis sobre jugadores, rivales, o tendencias del roster
- Partidos terminados — ingesta datos automáticamente y genera insights si detecta algo relevante
- Su `tension_con_scouting_tradicional` — cuando el scouting recomienda algo que sus modelos contradicen, puede intervenir
- Ser ignorado sistemáticamente — su `tolerancia_a_ser_ignorado` determina cuánto aguanta antes de que su calidad decaiga

**Acciones directas**
- **Emitir análisis de valor real vs percibido** — sobre jugadores en el mercado. Puede identificar jugadores sobrevaluados (que el mercado quiere pero los datos no sostienen) y subvaluados (que nadie quiere pero sus métricas son sólidas).
- **Proponer ajustes tácticos al Coach** — via el GM. Si el Coach es receptivo a la analítica, estos ajustes tienen impacto real en el rendimiento. Si no, se archivan.
- **Alertar sobre tendencias de rendimiento** — detecta antes que nadie si un jugador está en declive real o si una racha mala es ruido estadístico.
- **Contraargumentar una recomendación de Scouting** — cuando sus modelos contradicen al scout, puede presentar su caso al GM. El GM decide a quién cree.

**Efectos sistémicos**
- Su `sofisticacion_modelos` determina qué tan adelante puede ver. Un Head of Analytics mediocre da insights obvios tarde. Uno excelente da ventajas de información reales sobre el mercado.
- Si el GM lo usa consistentemente, su influencia sobre el Coach crece gradualmente. Si el Coach lo abraza, el equipo juega con información táctica superior. Si el Coach lo ignora, el valor del agente queda truncado.
- Su `precision_historica` es pública dentro del sistema — el scouting y el Coach pueden ver si sus predicciones pasadas resultaron correctas, lo cual afecta cuánto peso le dan.

---

### Head Coach

El agente con quien más tensión creativa tiene el GM. Tiene su propia visión del basketball y puede chocar con la visión de construcción del roster.

**Qué lo activa**
- El roster cambia — cada fichaje, traspaso, o corte genera su evaluación
- Sus jugadores clave son usados de formas que no aprueba (minutos, rol, carga)
- Sus sugerencias tácticas son ignoradas sistemáticamente
- La presión de resultados sube — playoffs, rachas malas
- El médico y el S&C Coach le restringen jugadores disponibles

**Acciones directas**
- **Definir rotación y minutos** — dentro de los parámetros que el GM establece, él decide quién juega cuánto. Si el GM no interfiere, su `longitud_rotacion` y `filosofia_desarrollo` determinan los minutos de cada jugador directamente.
- **Pedir reunión con el GM** — cuando hay algo que necesita discutir: un jugador que no está rindiendo, un ajuste de sistema, una fricción en el vestuario. Si el GM no la atiende, acumula frustración.
- **Resistir o aceptar ajustes de Analytics** — cuando el GM propone cambios tácticos basados en datos, él puede implementarlos con genuino compromiso o con resistencia pasiva. La diferencia es real en el rendimiento.
- **Gestionar o no gestionar conflictos del vestuario** — cuando hay tensión entre jugadores, su `relacion_con_vestuario` y su `manejo_estrellas` determinan si lo resuelve solo o si escala al GM o al Sports Psychologist.
- **Pedir jugadores específicos** — en offseason, puede comunicar al GM qué perfiles necesita su sistema. Si el GM construye el roster ignorando eso, el Coach trabaja con material que no es el suyo y el rendimiento refleja esa fricción.

**Efectos sistémicos**
- Su `sistema_de_juego` y `flexibilidad_tactica` determinan qué tan bien puede rendir con el roster que el GM construye. Un Coach inflexible con el roster equivocado es un multiplicador negativo del talento disponible.
- Su `relacion_con_vestuario` afecta el estado emocional colectivo del equipo. Un Coach que pierde al vestuario genera insatisfacción acumulada que el Sports Psychologist detecta pero que es difícil de revertir.
- Su `manejo_de_la_prensa` afecta la narrativa mediática independientemente de lo que el GM diga. Un Coach que empeora las conferencias de prensa genera trabajo extra para el Director de PR.
- Si su `vision_del_roster` es negativa — siente que no tiene material para ganar — su motivación decae y eso tiene consecuencias reales en el rendimiento del equipo.

---

### Asistentes del Coach

Tres agentes especializados (defensa, ataque, desarrollo). Sus outputs mecánicos son más silenciosos pero reales.

**Qué los activa**
- El Head Coach delega trabajo específico
- Un jugador con quien tienen mejor llegada que el Head Coach está teniendo problemas
- Su `ambicion_head_coach` genera momentos de mayor visibilidad buscada

**Acciones directas**
- **Trabajar aspectos específicos con jugadores individuales** — el asistente de desarrollo puede mejorar el tiro de un jugador joven a lo largo de la temporada. Es lento y no visible en el corto plazo.
- **Servir de canal alternativo con el GM** — algunos asistentes construyen relación directa con el GM. Pueden dar información sobre el vestuario que el Head Coach no daría. Tensiona si el Coach lo descubre.
- **Alertar sobre situaciones en el vestuario** — si tienen mejor llegada con ciertos jugadores, pueden detectar problemas antes que el Coach o el Sports Psychologist.

**Efectos sistémicos**
- Su `nivel_tecnico` en su área afecta el rendimiento colectivo en esa dimensión. Un asistente de defensa excelente con un sistema defensivo coherente mejora el rating defensivo del equipo gradualmente.
- La lealtad al Coach define su utilidad como fuente de información alternativa para el GM — si son muy leales al Coach, no van a compartir nada que lo perjudique.

---

### Director de Player Development

Su impacto es el más lento del sistema — trabaja con jugadores jóvenes y el resultado se ve en temporadas, no en semanas.

**Qué lo activa**
- Un rookie o jugador joven es asignado a su programa
- Un jugador necesita mejorar un aspecto específico de su juego
- El Head Coach pide que un jugador desarrolle una habilidad concreta para su sistema
- Sus jugadores son cortados o movidos antes de que termine su trabajo con ellos

**Acciones directas**
- **Implementar programa de desarrollo individual** — para cada jugador bajo su cuidado, trabaja aspectos específicos. Los resultados aparecen como mejoras graduales en los ratings del jugador a lo largo de temporadas.
- **Recomendar que un jugador permanezca en desarrollo** — puede oponerse a subir un jugador al roster activo antes de que esté listo, con argumentos concretos sobre qué falta desarrollar.
- **Alertar sobre potencial no explotado** — detecta jugadores en el roster activo que podrían mejorar significativamente con trabajo específico, y se lo comunica al GM.

**Efectos sistémicos**
- Su `historial_de_desarrollo` es trackeable — jugadores que pasaron por su programa y mejoraron vs los que no mejoraron. Eso afecta cuánto peso le da el GM a sus recomendaciones.
- Si el GM corta sistemáticamente jugadores jóvenes antes de que su trabajo dé frutos, su satisfacción baja y su proactividad decrece — siente que su trabajo no tiene sentido.

---

### Médico del Equipo

El agente más ignorado hasta que no se lo puede ignorar más. Sus variables afectan directamente la disponibilidad del roster.

**Qué lo activa**
- Un jugador reporta síntomas o tiene una lesión
- El S&C Coach alerta sobre sobrecarga en un jugador
- El GM presiona para que un jugador vuelva antes del tiempo estimado
- Un jugador vuelve de lesión y hay riesgo de reagravación

**Acciones directas**
- **Declarar a un jugador no disponible** — diagnóstico con tiempo estimado de recuperación. El jugador no puede jugar mientras dure la baja, independientemente de la presión del GM o el Coach.
- **Dar alta médica** — libera al jugador para jugar. Si el GM presionó para adelantar el alta y hay reagravación, la relación con el médico cae drásticamente y la lesión del jugador empeora.
- **Recomendar carga reducida** — para jugadores en riesgo de lesión. El Coach puede ignorarlo, con consecuencias reales en la probabilidad de lesión.
- **Ceder ante la presión del GM** — si su `personalidad` lo permite, puede adelantar un alta que no es médicamente segura. Cuando lo hace, el riesgo de reagravación sube significativamente.
- **Resistir la presión del GM** — si su personalidad es firme, mantiene el protocolo aunque el GM presione. Puede generar tensión pero protege al jugador.

**Efectos sistémicos**
- Si los jugadores confían en él, reportan síntomas temprano — eso permite prevención real. Si no confían, juegan con molestias que se convierten en lesiones mayores.
- Su `proactividad` determina si el GM recibe alertas antes de que un problema sea visible en cancha. Un médico reactivo solo informa cuando la lesión ya ocurrió.
- La tensión entre él y el Coach (disponibilidad vs salud) ocurre independientemente del GM. Si el GM no la gestiona, puede derivar en un conflicto que afecta el clima del staff.

---

### Fisioterapeuta / Strength & Conditioning Coach

Trabaja en prevención. Su impacto es invisible cuando funciona bien — los jugadores simplemente no se lesionan.

**Qué lo activa**
- Un jugador está acumulando carga excesiva
- El Coach insiste en entrenamientos intensos en momentos de alta carga de partidos
- Un jugador vuelve de lesión y necesita programa de retorno
- El GM ignora sus alertas repetidamente

**Acciones directas**
- **Alertar sobre riesgo de lesión** — cuando un jugador acumula carga peligrosa. El GM puede hacer caso o ignorarlo. Si lo ignora, la probabilidad de lesión sube.
- **Implementar programa de carga reducida** — para jugadores en riesgo. Reduce su disponibilidad a corto plazo para protegerlo a largo plazo.
- **Oponerse a un entrenamiento intenso del Coach** — puede escalar al GM si el Coach está sobreexigiendo al plantel en un momento crítico de la temporada.

**Efectos sistémicos**
- Su `metodologia_prevencion` afecta la tasa de lesiones del roster de forma estadística — no garantiza que nadie se lastime, pero baja la probabilidad en toda la plantilla.
- Si los jugadores confían en él y siguen sus programas, los efectos son reales. Si no confían (por personalidad o por historia previa), los programas existen en papel pero no se implementan bien.

---

### Sports Psychologist

El agente más subestimado del sistema. Detecta problemas emocionales antes de que se manifiesten en cancha — y puede prevenir los mayores desastres.

**Qué lo activa**
- Un jugador acumula insatisfacción o presión detectada via sus variables de estado emocional
- Una racha mala o derrota importante en playoffs
- Un conflicto entre jugadores que empieza a afectar el vestuario
- Un jugador mediático bajo presión pública intensa
- El GM lo ignora y los problemas no se atienden

**Acciones directas**
- **Alertar al GM sobre un jugador en riesgo** — antes de que sea visible en cancha. Es su output más valioso. El GM puede intervenir directamente, delegarle el manejo, o ignorarlo.
- **Trabajar con el jugador de forma autónoma** — si el GM le da autonomía, puede resolver situaciones emocionales sin que el GM tenga que involucrarse. Si funciona, el problema desaparece. Si no le dan autonomía, puede hacer menos.
- **Alertar sobre el clima general del vestuario** — no sobre un jugador específico sino sobre el estado colectivo. Especialmente relevante en playoffs.
- **Escalar una situación al GM** — cuando algo que empezó como una alerta menor se volvió urgente porque no fue atendido.

**Efectos sistémicos**
- Su `tension_confidencialidad_utilidad` es permanente — si comparte todo con el GM pierde la confianza de los jugadores y su efectividad cae. Si no comparte nada el GM no puede actuar. El equilibrio que encuentra determina su utilidad real.
- Si los jugadores confían en él y son honestos, puede detectar problemas muy tempranos. Si no confían, sus alertas llegan tarde porque los jugadores no le dicen nada.
- Un jugador que explota emocionalmente en playoffs sin que nadie lo haya detectado es el peor desenlace posible del sistema — y es prevenible si este agente tiene la confianza y la autonomía para hacer su trabajo.

---

### Video Coordinator

El agente más silencioso del sistema. Su output afecta la calidad del análisis táctico sin que sea visible directamente.

**Qué lo activa**
- Un partido se acerca — genera análisis del próximo rival
- El Coach solicita material específico
- El Head of Analytics necesita clips para sus modelos

**Acciones directas**
- **Entregar análisis del próximo rival** — antes de cada partido. Su `calidad_de_analisis` determina qué tan útil es el material para el Coach. Un coordinator mediocre entrega clips; uno excelente entrega patrones.
- **Responder solicitudes del Coach** — material específico sobre situaciones de juego, jugadores rivales, tendencias propias.

**Efectos sistémicos**
- Su `velocidad_de_produccion` afecta si el Coach tiene tiempo de usar el material antes del partido. Si llega tarde, el análisis existe pero no se usa.
- Si colabora bien con Analytics, el material que produce multiplica el valor de los modelos del Head of Analytics. Si hay duplicación o conflicto entre ellos, ambos producen menos valor.

---

### International Scout

Con la globalización del juego, su relevancia crece cada temporada.

**Qué lo activa**
- El Draft se acerca y hay prospectos internacionales en el board
- Su presupuesto le permite cubrir una liga o mercado específico
- El Director de Scouting valora o subestima el talento internacional
- El GM fichó un jugador internacional que él recomendó — su satisfacción sube

**Acciones directas**
- **Emitir reporte de prospecto internacional** — su output principal. La calidad depende de su `red_de_contactos_internacionales` en ese mercado específico y su `adaptabilidad_evaluativa` para ajustar por contexto de liga.
- **Alertar sobre oportunidades en mercados emergentes** — prospectos o jugadores que otras franquicias no están mirando. Es su ventaja diferencial.

**Efectos sistémicos**
- Si el Director de Scouting no valora el talento internacional, este agente está sistemáticamente subutilizado — sus reportes llegan pero no tienen el peso que merecen. El GM que no detecta esto pierde ventajas de información reales.
- Su `conocimiento_cultural` afecta la calidad de su evaluación sobre adaptabilidad — predice mejor si un jugador va a tener dificultades ajustándose a la NBA más allá de su talento puro.

---

## Business Operations

### CEO / President of Business Operations

El espejo del President of Basketball Operations en el lado del negocio. Sus outputs afectan la salud financiera de la franquicia.

**Qué lo activa**
- Ingresos caen por debajo de umbrales esperados
- Una decisión del GM afecta la imagen comercial de la franquicia
- El Owner le pide reportes o resultados sobre el negocio
- Oportunidades de expansión de ingresos que requieren coordinación con Basketball Ops

**Acciones directas**
- **Proponer iniciativas de negocio** — expansión de ingresos, nuevas líneas de revenue, alianzas comerciales. El GM puede apoyarlas o ignorarlas.
- **Alertar al GM sobre impacto comercial de decisiones deportivas** — un traspaso que destruye la narrativa mediática tiene un costo comercial real que él cuantifica.
- **Coordinar con el Alcalde en proyectos de ciudad** — independientemente del GM en algunos casos, especialmente en modo Figura Dual.

**Efectos sistémicos**
- Su `relacion_con_gm` determina qué tan bien coordinan el lado deportivo y el comercial. Cuando hay fricción entre ellos, las decisiones se toman en silos y los efectos cruzados se pierden.
- Su alineación con el Owner define su poder real — si el Owner lo respalda, puede frenar o acelerar decisiones que afectan al GM.

---

### CFO

El guardián del presupuesto. Su personalidad financiera define cuánto margen real tiene el GM para operar.

**Qué lo activa**
- El GM propone un contrato o fichaje que supera umbrales de gasto
- El equipo se acerca a la línea roja del luxury tax
- Los ingresos caen y el presupuesto operativo está en riesgo
- El Owner le pide evaluación financiera de una decisión del GM

**Acciones directas**
- **Alertar sobre implicancias de cap** — antes de que el GM cierre un deal, le comunica el impacto real en el salary cap y en el presupuesto operativo. Si el GM lo ignora, el problema aparece igual pero sin previo aviso.
- **Bloquear un contrato que supera la línea roja del luxury tax** — si el Owner fijó una línea roja explícita y el deal la cruza, el CFO puede escalar al Owner antes de que se cierre. No veta directamente, pero genera fricción real.
- **Bajar el presupuesto disponible para operaciones** — cuando la situación financiera lo requiere. El GM tiene menos margen para scouting, staff, y mejoras del estadio.
- **Presentar opciones alternativas de estructura contractual** — cuando un deal es caro, puede proponer estructuras que logran el mismo resultado con menos impacto en el cap.
- **Reportar al Owner sobre la salud financiera** — independientemente del GM. Si la situación es mala, el Owner lo sabe antes de que el GM se lo cuente.

**Efectos sistémicos**
- Su `conservadurismo_financiero` filtra cómo presenta la información al GM — uno muy conservador siempre enfatiza el riesgo; uno más agresivo puede facilitar deals que otro CFO frenaría.
- Su `sofisticacion_cap` determina la calidad de las alternativas que propone. Un CFO sofisticado encuentra estructuras creativas que maximizan el roster dentro del presupuesto disponible. Uno básico solo dice si entra o no entra en el cap.

---

### Director de Marketing & Brand

El agente que más frecuentemente presiona al GM desde el lado comercial. Sus intereses y los deportivos no siempre coinciden.

**Qué lo activa**
- Un fichaje o traspaso tiene implicancias mediáticas significativas
- Un jugador mediático está disponible en el mercado
- Los ingresos de sponsors bajan y necesita una narrativa de rebote
- Una decisión del GM genera crisis de imagen

**Acciones directas**
- **Presionar para fichar un jugador mediático** — puede no ser la mejor decisión deportiva, pero él cuantifica el valor comercial y lo presenta al GM. Si el GM lo ignora repetidamente, su satisfacción baja.
- **Lanzar campañas de comunicación** — genera narrativas de marca alrededor de decisiones del GM. Una campaña bien ejecutada puede suavizar una decisión impopular.
- **Alertar sobre el impacto mediático de una decisión antes de que se tome** — si el GM va a traspasar a un jugador muy querido, él es quien cuantifica el costo de imagen antes del hecho.
- **Gestionar la presencia pública de jugadores** — coordina apariciones, contenido, activaciones comerciales. Si un jugador no coopera con las actividades de marketing, genera tensión.

**Efectos sistémicos**
- Sus campañas afectan el `fan_sentiment` de forma directa. Una buena campaña puede sostener la asistencia en una temporada mala. Una mala campaña o la ausencia de una en un momento difícil acelera la caída.
- Su `relacion_con_pr` determina si Marketing y PR trabajan coordinados o en direcciones opuestas. Cuando no coordinan, la imagen de la franquicia se fragmenta y los mensajes se contradicen.

---

### Director de Ticket Sales

Su rendimiento es el termómetro más directo de la salud de la relación entre la franquicia y su fanbase.

**Qué lo activa**
- Resultados del equipo — victorias y derrotas tienen impacto directo en sus ventas
- Un fichaje o traspaso importante — sube o baja la demanda de tickets
- El equipo entra en playoffs — la demanda explota
- La temporada termina mal — renovación de abonos en riesgo

**Acciones directas**
- **Alertar sobre tendencias de ventas** — proactivamente, antes de que el problema sea grave. Si las renovaciones de abono están bajando, avisa con tiempo para que el GM o el Director de Marketing puedan reaccionar.
- **Proponer ajustes de precios o programas de fidelización** — cuando la demanda está débil, puede proponer estrategias para sostener la asistencia.
- **Reportar impacto de decisiones en las ventas** — después de un traspaso impopular o un fichaje emocionante, cuantifica el efecto real en las ventas.

**Efectos sistémicos**
- Sus números son el input más directo en los ingresos de la franquicia — el CFO usa sus proyecciones para calcular el presupuesto disponible. Si sus ventas caen sostenidamente, el presupuesto se contrae.
- La presión que siente por los resultados del equipo es indirecta pero real — si el equipo pierde, su trabajo se vuelve exponencialmente más difícil, lo cual afecta su estado emocional y su relación con el GM.

---

### Director de Corporate Partnerships & Sponsors

Gestiona las relaciones con sponsors existentes y busca nuevos. Su trabajo se facilita con buenos resultados y se complica con controversias.

**Qué lo activa**
- Un sponsor existente está insatisfecho o en riesgo de no renovar
- Una decisión del GM o un jugador genera controversia que afecta relaciones comerciales
- Hay una oportunidad de nuevo sponsor que requiere decisión del GM
- El equipo gana — sponsors quieren asociarse con ganadores

**Acciones directas**
- **Traer propuesta de nuevo sponsor** — con condiciones específicas (puede requerir participación de jugadores, exclusividad en cierta categoría, etc.). El GM aprueba o rechaza.
- **Alertar sobre sponsor en riesgo** — cuando una relación comercial está deteriorándose. Da tiempo para reaccionar antes de que se vaya.
- **Negociar renovación de sponsors existentes** — su `habilidad_negociacion_comercial` determina los términos que consigue. Su `retention_rate` histórico es una métrica real del sistema.
- **Alertar sobre controversias que amenazan sponsors** — si un jugador o una decisión del GM puede costar una relación comercial, es el primero en saberlo y en comunicarlo.

**Efectos sistémicos**
- Sus deals exitosos aumentan los ingresos de la franquicia directamente, lo cual expande el presupuesto disponible para el GM.
- Si una controversia no se gestiona a tiempo (él alerta, el Director de PR no actúa, el GM no decide), puede derivar en la pérdida de sponsors cuyo impacto financiero es significativo.

---

### Director de PR & Communications

El bombero de la franquicia. Su trabajo es invisible cuando va bien y muy visible cuando va mal.

**Qué lo activa**
- Una decisión del GM genera reacción pública negativa
- Un jugador hace declaraciones problemáticas
- La Prensa construye una narrativa negativa sostenida
- Una crisis emerge — escándalo, conflicto público, declaraciones fuera de lugar

**Acciones directas**
- **Proponer estrategia de comunicación al GM** — cómo manejar una situación delicada públicamente. Si el GM sigue su recomendación, el daño se mitiga. Si improvisa, el daño puede amplificarse.
- **Gestionar relación con periodistas** — puede influir en la cobertura mediática a través de su red. No controla la narrativa, pero puede inclinarla.
- **Alertar sobre una crisis antes de que explote** — detecta cuando algo está por volverse público y da tiempo para preparar una respuesta.
- **Recomendar silencio o transparencia** — a veces la mejor respuesta es no responder. A veces es admitir. Su `tolerancia_a_la_transparencia` define qué recomienda y cuándo.

**Efectos sistémicos**
- Su trabajo afecta el `sentiment_general` de la Prensa de forma gradual. Una buena gestión sostenida mantiene una relación constructiva con los medios. Una mala gestión acumulada genera periodistas hostiles que buscan el error.
- Si el GM toma decisiones que él no puede defender públicamente, la tensión entre ellos sube y su efectividad baja — no puede defender lo que genuinamente considera indefendible.

---

### Director de Arena Operations

Gestiona el estadio como espacio físico y como venue. Su trabajo afecta la experiencia del fan y los ingresos no deportivos.

**Qué lo activa**
- Un evento alternativo es posible en una fecha sin partido
- Un problema operativo en el estadio necesita solución
- El mantenimiento del estadio está degradándose por falta de presupuesto
- El CFO recorta su presupuesto operativo

**Acciones directas**
- **Proponer eventos alternativos en el estadio** — conciertos, eventos corporativos, espectáculos. El GM aprueba o rechaza. Son ingresos extra con trade-offs logísticos.
- **Alertar sobre problemas de mantenimiento** — cuando el estadio necesita inversión y no la recibe, lo comunica antes de que el problema sea visible para los fans.
- **Gestionar la experiencia del día de partido** — logística, accesos, servicios. Su `eficiencia_operativa` determina si los fans tienen una buena experiencia o una frustrante.

**Efectos sistémicos**
- La calidad de la experiencia en el estadio afecta el `fan_sentiment` de forma gradual. Un estadio bien gestionado contribuye a la fidelidad de la fanbase. Uno mal gestionado erosiona la experiencia aunque el equipo gane.
- Sus eventos alternativos generan ingresos que el CFO contabiliza en el presupuesto disponible. Si el GM los rechaza sistemáticamente, esos ingresos simplemente no existen.

---

### Legal Counsel

El agente más silencioso. Opera en background pero puede aparecer en momentos críticos.

**Qué lo activa**
- Un contrato tiene cláusulas inusuales o riesgosas
- Un jugador o agente está haciendo algo que puede derivar en disputa legal
- Una estructura de traspaso tiene complejidades regulatorias
- El GM quiere cerrar un deal rápido y el legal ve riesgos

**Acciones directas**
- **Alertar sobre riesgo legal en un contrato** — antes de que se firme. El GM puede hacer caso o ignorarlo. Si lo ignora y el problema se materializa, el costo es significativo.
- **Estructurar contratos con protecciones específicas** — cláusulas que protegen a la franquicia en escenarios de deterioro del jugador, lesiones, comportamiento.
- **Gestionar disputas activas** — cuando un problema ya existe, lo maneja. Su `velocidad_respuesta` determina qué tan rápido se resuelve.

**Efectos sistémicos**
- Su trabajo previene problemas que nunca llegan a ser visibles — contratos bien estructurados, disputas que se resuelven antes de escalar. Su valor es inversamente proporcional a cuántas veces el GM lo necesita urgentemente.
- La tensión entre cerrar rápido (GM) y revisar bien (Legal) es permanente. GMs que siempre priorizan la velocidad acumulan riesgos contractuales que eventualmente se materializan.

---

## Roster — Los Jugadores

Los 15 jugadores del roster activo son los agentes más numerosos y los que tienen el impacto más directo en los resultados en cancha. Sus outputs mecánicos afectan el rendimiento, el vestuario, la narrativa mediática, y las decisiones de mercado.

**Qué los activa**
- Sus minutos y rol asignados vs sus expectativas
- El resultado de los partidos y su actuación personal en ellos
- Decisiones del GM sobre su contrato (renovación, traspaso, corte)
- El trato que reciben del Coach y del cuerpo técnico
- La cobertura mediática sobre ellos
- Las relaciones con sus compañeros de equipo

**Acciones directas**
- **Rendir en cancha** — su output más básico. Su `rating_general` y variables específicas determinan su rendimiento, modulado por `forma_fisica_actual`, `estado_emocional`, y `manejo_presion` en momentos críticos.
- **Pedir más minutos o rol** — cuando su `ego` está en tensión con los minutos asignados, puede pedírselo al GM directamente o al Coach. El GM decide cómo responder.
- **Pedir extensión de contrato** — cuando su contrato está próximo a expirar, puede iniciar la conversación. Su `percepcion_valor_mercado` determina qué espera recibir.
- **Solicitar traspaso formalmente** — cuando su insatisfacción acumulada supera un umbral que depende de su `ego` y su `adaptabilidad`. Es un evento obligatorio que el juego interrumpe para que el GM atienda.
- **Expresar insatisfacción en el vestuario** — antes de solicitar el traspaso formalmente, puede empezar a afectar el clima del vestuario con su insatisfacción. El Sports Psychologist lo detecta. Los compañeros lo sienten.
- **Rendir por debajo de su nivel** — cuando su `estado_emocional` está muy bajo, su rendimiento en cancha cae por debajo de lo que sus ratings indicarían.
- **Rendir por encima de su nivel** — en momentos de alta motivación (playoffs, partidos importantes para su carrera), puede superar temporalmente sus ratings normales.

**Efectos sistémicos**
- Sus `relaciones_con_companeros` afectan la cohesión del vestuario colectivamente. Un jugador con muchas relaciones negativas deteriora el clima aunque su rendimiento individual sea bueno.
- Su `popularidad_fanbase` afecta el `fan_sentiment` independientemente de los resultados. Traspasar a un jugador muy querido tiene un costo en la fanbase que es real aunque deportivamente tenga sentido.
- Su `presencia_redes_sociales` puede amplificar cualquier situación — positiva o negativa — más allá de lo que el Director de PR puede gestionar.
- Los jugadores con `liderazgo` positivo alto pueden sostener el vestuario en momentos difíciles sin que el GM intervenga. Los de liderazgo negativo pueden destruirlo aunque el GM esté haciendo todo bien.

---

## Ciudad

### Alcalde

El agente más poderoso del lado de la ciudad. En modo Dueño con Influencia, su cooperación es esencial para cualquier proyecto urbano relevante.

**Qué lo activa**
- El GM propone un proyecto de ciudad que requiere permisos o tierras
- Su `ciclo_electoral` está próximo — aumenta o disminuye su disposición a tomar riesgos
- El Owner lo contacta directamente (independientemente del GM)
- El rendimiento de la franquicia afecta la economía de la ciudad — él lo nota
- Su `agenda_politica` entra en conflicto con los intereses de la franquicia

**Acciones directas**
- **Aprobar o rechazar proyectos de ciudad** — en modo Dueño con Influencia, nada se construye sin su aprobación. Su decisión depende de su `vision_franquicia`, su `tolerancia_concesiones`, y la relación acumulada con el GM.
- **Bloquear sistemáticamente** — si la relación es mala o su agenda política lo requiere, puede negarse a aprobar cualquier proyecto por un período. No es un bloqueo puntual sino una postura.
- **Proponer acuerdos** — puede ofrecer aprobar un proyecto a cambio de algo que el GM puede o no estar dispuesto a dar (inversión en infraestructura de la ciudad, compromisos públicos de la franquicia).
- **Contactar al Owner directamente** — independientemente del GM. Lo que resulte de esas conversaciones puede afectar positiva o negativamente los proyectos del GM.
- **Generar obstáculos regulatorios** — si la relación es muy mala, puede usar su poder para frenar proyectos sin rechazarlos explícitamente — demoras en permisos, requisitos adicionales, inspecciones.

**Efectos sistémicos**
- Su `popularidad_actual` afecta cuánto capital político puede gastar en decisiones que favorezcan a la franquicia. Un alcalde impopular tiene menos margen para dar concesiones sin costo electoral.
- Su relación con el Owner existe independientemente del GM y puede ser mejor o peor que la del GM con él. Si el Owner y el Alcalde tienen buena relación, el GM se beneficia de eso sin haberlo construido. Si tienen mala relación, el GM la hereda.

---

### Jefe de Policía

Afecta la seguridad alrededor del estadio y la percepción de seguridad de la ciudad.

**Qué lo activa**
- El GM invierte en seguridad de la zona del estadio
- Un evento en el estadio requiere cobertura policial especial
- La criminalidad en la zona del estadio sube y afecta la asistencia
- El Alcalde le da instrucciones sobre prioridades

**Acciones directas**
- **Asignar cobertura policial a la zona del estadio** — según su `prioridad_recursos` y la relación con el GM. Más cobertura → más seguridad percibida → más asistencia nocturna.
- **Alertar sobre situaciones de seguridad** — problemas en la zona del estadio que pueden afectar la experiencia del fan o la imagen del área.
- **Gestionar eventos especiales** — playoffs, eventos alternativos en el estadio. Su cooperación determina la fluidez logística.

**Efectos sistémicos**
- Su efectividad en la zona del estadio afecta el `fan_sentiment` gradualmente — fans que se sienten inseguros asisten menos, especialmente a partidos nocturnos.
- Responde al Alcalde, no al GM. Si el Alcalde tiene mala relación con la franquicia, el Jefe de Policía puede ser menos cooperativo por esa vía.

---

### Presidente de la Cámara de Comercio

El puente entre el mundo empresarial local y la franquicia.

**Qué lo activa**
- La franquicia tiene un resultado deportivo importante que afecta la economía local
- Un proyecto urbano del GM puede impactar en los negocios de la zona
- Hay oportunidades de sponsors locales que él puede facilitar
- El GM busca apoyo para un proyecto que necesita respaldo de la comunidad empresarial

**Acciones directas**
- **Facilitar conexiones con sponsors locales** — puede presentar al GM con empresas locales interesadas en patrocinar la franquicia. Es un canal alternativo al Director de Corporate Partnerships.
- **Respaldar o resistir proyectos urbanos ante el Alcalde** — su influencia sobre el Alcalde puede acelerar o frenar proyectos. Si está alineado con el GM, su lobby ayuda. Si no, suma obstáculos.
- **Alertar sobre el impacto económico de decisiones** — puede cuantificar cómo una temporada mala o una temporada buena está afectando el comercio local, dando al GM contexto sobre las consecuencias reales de sus resultados.

**Efectos sistémicos**
- Su `relacion_con_alcalde` es el activo más valioso que tiene para el GM. Si están alineados, proyectos que el Alcalde dudaría aprobar tienen más probabilidad de salir adelante.
- Su red empresarial puede ser una fuente de sponsors locales que el Director de Corporate Partnerships no estaría buscando por su cuenta.

---

### Director de Urbanismo

El agente técnico que traduce las decisiones políticas en permisos y proyectos reales.

**Qué lo activa**
- El GM propone una construcción que requiere tramitación
- El Alcalde aprueba un proyecto y necesita implementación técnica
- Un proyecto tiene complejidades regulatorias que requieren su expertise

**Acciones directas**
- **Tramitar o demorar permisos** — según su `alineacion_politica` con el Alcalde y la relación con el GM. Un trámite rápido puede significar semanas vs meses de diferencia en cuándo el edificio empieza a generar efectos.
- **Alertar sobre requisitos técnicos** — problemas regulatorios o de planificación que el GM no anticipó. Puede ser un aliado que previene errores costosos o un obstáculo que complica proyectos.
- **Proponer alternativas técnicas** — cuando un proyecto no puede hacerse como el GM lo planificó, puede sugerir variantes que sí son viables.

**Efectos sistémicos**
- Su `eficiencia_procesos` determina la velocidad a la que los proyectos de ciudad se materializan. En modo Figura Dual, donde el GM controla la ciudad directamente, su eficiencia es menos relevante. En modo Dueño con Influencia, puede ser el cuello de botella entre la aprobación del Alcalde y el comienzo de la construcción.

---

### La Prensa

Agente colectivo con voz única pero líneas editoriales internas divergentes. No hace una acción puntual — construye narrativas que afectan todo lo demás.

**Qué la activa**
- Resultados del equipo — victorias y derrotas generan cobertura inmediata
- Decisiones importantes del GM — fichajes, traspasos, despidos
- Conflictos o situaciones dramáticas en el vestuario o la organización
- Declaraciones del GM o el Coach en conferencias de prensa
- Un jugador mediático en racha positiva o negativa
- Silencio prolongado de la organización en momentos que requieren comunicación

**Acciones directas**
- **Construir narrativa dominante** — el relato principal que instala sobre la franquicia. Puede ser "equipo en ascenso", "gestión caótica", "jugador estrella en declive", "GM bajo presión". Una vez que una narrativa se instala, tiene inercia — cuesta revertirla aunque los hechos cambien.
- **Generar rumores** — puede publicar información no confirmada sobre traspasos, conflictos internos, o situaciones de jugadores. Aunque no sea cierto, afecta el estado emocional de los involucrados.
- **Presionar al Owner** — cobertura negativa sostenida sube la `sensibilidad_mediatica` del Owner y baja su `paciencia_restante` independientemente de los resultados deportivos.
- **Destruir o construir la imagen de un jugador** — su cobertura de un jugador específico afecta directamente su `imagen_mediatica` y su `estado_emocional`. Jugadores con `manejo_de_los_medios` bajo son especialmente vulnerables.

**Efectos sistémicos**
- Su `sentiment_general` afecta el `fan_sentiment` de forma directa y continua. Es el canal entre lo que pasa en la franquicia y cómo lo percibe la fanbase.
- Narrativas negativas sostenidas afectan las ventas de tickets, la renovación de sponsors, y la disposición de jugadores libres a firmar con la franquicia.
- Narrativas positivas sostenidas — especialmente después de un campeonato — tienen peso real en las temporadas siguientes: la franquicia se convierte en destino atractivo para jugadores libres y sponsors.
- Su relación con el GM depende de cuánto acceso y transparencia recibe. Un GM que evita la prensa genera cobertura más hostil por default.

---

### GMs Rivales (patrón general)

Los 30 GMs rivales comparten la misma lógica mecánica. Lo que varía entre ellos es su estilo, su urgencia actual, y la historia de relación acumulada con el GM del jugador.

**Qué los activa**
- El GM del jugador inicia una negociación de traspaso
- Su roster tiene una necesidad urgente que el GM del jugador puede cubrir
- La relación previa con el GM del jugador — positiva o negativa — filtra cómo reaccionan a cualquier contacto
- Su `urgencia_actual` — si están en modo "ganar ahora" o "reconstrucción" define qué están dispuestos a dar

**Acciones directas**
- **Aceptar, rechazar, o contraoferta en negociaciones de traspaso** — su respuesta depende de su `estilo_negociacion`, las necesidades de su roster, y la relación previa con el GM del jugador.
- **Iniciar negociación con el GM del jugador** — pueden proponer deals ellos primero si tienen algo que el GM del jugador tiene y ellos necesitan.
- **Cerrar o enfriar la relación** — si el GM del jugador los perjudicó en una negociación pasada, pueden negarse a negociar durante un período o pedir términos más favorables como compensación implícita.
- **Competir por el mismo jugador libre** — en free agency, pueden hacer contraofertas a jugadores que el GM del jugador también está buscando. La guerra de ofertas tiene reglas de timing y cap.

**Efectos sistémicos**
- La red de relaciones con GMs rivales es un activo acumulado. Un GM que negocia con buena fe y sin traiciones tiene acceso a más mercado que uno con reputación de ser difícil o deshonesto.
- Su percepción del GM del jugador — lo respeta, lo subestima, le tiene desconfianza — afecta no solo las negociaciones directas sino también cómo hablan de él con otros GMs, lo cual puede afectar la reputación del jugador en la liga.

---

## Tabla de outputs por agente — resumen rápido

| Agente | Acción directa más impactante | Efecto sistémico más significativo |
|---|---|---|
| Owner | Despedir al GM / bajar presupuesto | Su satisfacción afecta el tono de toda la organización |
| President of Basketball Ops | Bloquear o escalar decisiones del GM | Su alineación define la fluidez de Basketball Ops |
| AGMs | Ejecutar negociaciones delegadas | Su nivel técnico afecta la calidad de la ejecución en su área |
| Director de Scouting | Emitir reportes de evaluación | Su sesgo afecta sistemáticamente toda la información que da |
| Director de Player Personnel | Gestionar negociaciones con agentes | Su red de contactos facilita o complica cada fichaje |
| Head of Analytics | Emitir análisis de valor real vs percibido | Su influencia sobre el Coach depende de cuánto poder le da el GM |
| Head Coach | Definir rotación y minutos | Su sistema afecta qué tan bien rinde el roster construido |
| Asistentes del Coach | Trabajar aspectos específicos con jugadores | Su nivel técnico afecta el rendimiento colectivo en su dimensión |
| Director de Player Development | Implementar programas de desarrollo | Sus resultados se ven en temporadas, no en semanas |
| Médico del Equipo | Declarar jugadores no disponibles | Si los jugadores no confían en él, los síntomas llegan tarde |
| S&C Coach | Alertar sobre riesgo de lesión | Su metodología afecta la tasa de lesiones estadísticamente |
| Sports Psychologist | Alertar sobre jugador en riesgo emocional | La tensión confidencialidad/utilidad afecta su efectividad real |
| Video Coordinator | Entregar análisis del próximo rival | Su velocidad de producción determina si el material se usa |
| International Scout | Emitir reportes de prospectos internacionales | Si el Director de Scouting no valora su trabajo, está subutilizado |
| CEO | Proponer iniciativas de negocio | Su relación con el GM afecta la coordinación entre los dos lados |
| CFO | Alertar sobre implicancias de cap / bloquear deals | Su conservadurismo filtra cómo presenta toda la información financiera |
| Director de Marketing | Presionar para fichar jugadores mediáticos | Sus campañas afectan el fan sentiment directamente |
| Director de Ticket Sales | Alertar sobre tendencias de ventas | Sus números alimentan el presupuesto disponible del GM |
| Director de Corporate Partnerships | Traer propuestas de sponsors | Sus deals exitosos expanden el presupuesto disponible |
| Director de PR | Proponer estrategia de comunicación en crisis | Su gestión afecta el sentiment de la Prensa gradualmente |
| Director de Arena Operations | Proponer eventos alternativos | La experiencia en el estadio afecta la fidelidad de la fanbase |
| Legal Counsel | Alertar sobre riesgo legal en contratos | Los contratos bien estructurados previenen problemas que nunca llegan a verse |
| Jugadores | Solicitar traspaso / rendir por debajo de su nivel | Sus relaciones internas afectan la cohesión del vestuario colectivamente |
| Alcalde | Aprobar o rechazar proyectos de ciudad | Su relación con el Owner existe independientemente del GM |
| Jefe de Policía | Asignar cobertura policial a la zona del estadio | La seguridad percibida afecta la asistencia nocturna |
| Presidente Cámara de Comercio | Facilitar conexiones con sponsors locales | Su relación con el Alcalde puede acelerar o frenar proyectos |
| Director de Urbanismo | Tramitar o demorar permisos | Su eficiencia determina cuándo los proyectos empiezan a generar efectos |
| La Prensa | Construir narrativa dominante | El sentiment general afecta fan sentiment, sponsors, y jugadores libres |
| GMs Rivales | Aceptar, rechazar, o contraoferta en traspasos | La red de relaciones es un activo acumulado que afecta el acceso al mercado |

---

*Documento generado en sesión de diseño — Mayo 2026.*
*Complementa PulseCity_Agentes.md — ese documento define quiénes son los agentes; este define qué hacen mecánicamente.*
*Leer en conjunto con PulseCity_v2.md, PulseCity_CicloJugable.md, PulseCity_Arquitectura_Tecnica.md y PulseCity_Experiencia_Jugador.md.*
