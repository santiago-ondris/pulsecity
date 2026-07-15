# PulseCity — Ciclo Jugable

> Documento de referencia completo del ciclo de una temporada en PulseCity. Cubre cada fase, sus decisiones disponibles, y los efectos en los agentes y la ciudad.

---

## Filosofía del ciclo

El ciclo jugable de PulseCity no es una secuencia de pantallas — es un sistema vivo donde el tiempo avanza y el mundo reacciona constantemente, con o sin la intervención del jugador.

Dos principios rigen todo el ciclo:

**Las decisiones son causales.** Cada acción tiene consecuencias reales que pueden manifestarse inmediatamente o semanas simuladas después. El sistema tiene memoria. Una decisión tomada en el Draft puede tener consecuencias en los Playoffs de esa misma temporada o en la free agency de la siguiente.

**Los agentes viven en paralelo.** Mientras el jugador gestiona una fase, los agentes están reaccionando, interactuando entre sí, y acumulando estado. El Coach y el médico tienen su tensión propia aunque el GM no los esté mirando. El Alcalde llama al Owner directamente. La Prensa construye narrativas sin esperar instrucciones.

---

## Control del tiempo

**Tiempo continuo con pausa manual** — el modo base. El jugador controla la velocidad (pausa / x1 / x5 / x20) en cualquier momento. La ciudad y la temporada regular corren así.

**Pausa automática + UI expandida** — se activa en momentos de decisión crítica que el juego detecta: Draft, Free Agency, Trade Deadline, Cierre de Temporada. El juego lleva al jugador a ese estado sin que tenga que acordarse de pausar. La UI se adapta para mostrar exactamente lo que se necesita en ese momento.

**Eventos obligatorios** — algunos eventos interrumpen el tiempo independientemente del modo de notificaciones elegido: un jugador pide traspaso formal, el Owner convoca reunión de emergencia, un agente da ultimátum con deadline real.

---

## El canal de consulta directa

En cualquier momento del ciclo, desde cualquier fase, el jugador tiene acceso al panel de personajes — todos los agentes listados con su estado actual visible. Desde cualquier perfil puede abrir un chat directo y consultarle lo que necesite dentro de su dominio.

El CFO explica opciones de salary cap en lenguaje humano. El Director de Scouting da su opinión sobre un prospecto específico. El Sports Psychologist da su lectura del vestuario. El Director de Urbanismo informa el estado de permisos pendientes.

La frecuencia de interacción afecta la relación — agentes consultados regularmente son más proactivos y abiertos. Agentes ignorados hasta que hay una crisis responden con más frialdad.

Cada agente está anclado a su dominio por diseño. No responde preguntas fuera de su área — una línea, se acabó.

---

## Las reglas complejas — filosofía de diseño

Las reglas del salary cap, los tipos de contratos, las excepciones, las reglas de traspasos — todo existe en el sistema con su complejidad real. Pero el jugador no necesita aprenderlas.

Los agentes las manejan y las explican en lenguaje humano cuando son relevantes. El CFO te dice "podemos usar la mid-level exception para fichar a este jugador sin romper el cap" — no necesitás saber qué es la mid-level exception. Tu Director de Player Personnel estructura los deals de traspaso — vos decidís si querés al jugador, él te dice qué hace falta para lograrlo.

Este principio se extiende a toda complejidad del sistema: si algo es complejo, hay un agente que lo entiende y te lo traduce.

---

## El ciclo completo

```
Draft
  ↓
Free Agency
  ↓
Traspasos y armado de roster
  ↓
Pretemporada
  ↓
Temporada Regular
  ↓
Playoffs
  ↓
Cierre de Temporada
  ↓
Draft (nueva temporada)
```

---

## OFFSEASON — Fase 1: El Draft

**Modo:** Pausa automática + UI expandida el día del draft. Tiempo continuo antes y después.

El Draft tiene dos rondas, 30 picks por ronda. La posición del jugador en el draft depende del récord de la temporada anterior — peor récord, pick más alto en la lotería.

### Antes del Draft

**Decisiones disponibles:**

Mandar scouts a ver prospectos específicos antes del draft.
- Agentes afectados: Director de Scouting (carga de trabajo sube, precisión de reportes mejora si tiene tiempo suficiente), International Scout si los prospectos son internacionales.

Negociar traspasos de picks con otros equipos.
- Agentes afectados: GMs rivales (reaccionan según su estilo de negociación y relación previa), Director de Player Personnel (satisfacción sube si se le da autonomía para negociar), CFO (evalúa el impacto financiero de assets incluidos).

Decidir subir o bajar en el draft.
- Agentes afectados: Owner (si quiere ganar ya, bajar en el draft lo frustra), Head Coach (tiene preferencias sobre el perfil de jugador que necesita su sistema).

### El día del Draft — UI expandida

El board completo de prospectos visible. Cada pick que se acerca el sistema muestra los jugadores disponibles con las evaluaciones del equipo de scouting.

**Decisiones disponibles:**

Elegir el jugador recomendado por el Director de Scouting.
- Agentes afectados: Director de Scouting (confianza sube, próximas interacciones más fluidas y proactivas), Head Coach (evalúa si el perfil encaja en su sistema).

Ignorar la recomendación del Director de Scouting.
- Agentes afectados: Director de Scouting (confianza baja, empieza a autocensurarse en reportes futuros — degradación real de la calidad de la información que recibís), Head Coach (puede estar de acuerdo o no con la elección alternativa).

Elegir el jugador más marketeable.
- Agentes afectados: Director de Marketing (satisfacción sube), Head Coach (puede frustrarse si no encaja en su sistema), Owner (reacción según su prioridad principal).

Elegir según el sistema del Coach.
- Agentes afectados: Head Coach (confianza en el GM sube), Director de Player Development (recibe el tipo de jugador que puede trabajar mejor).

Elegir un prospecto internacional.
- Agentes afectados: International Scout (satisfacción sube, siente que su trabajo importa), Director de Scouting (evalúa la decisión según su alineación filosófica con el talento internacional).

Traspasar el pick en el momento.
- Agentes afectados: GM rival involucrado (relación afectada según los términos), Owner (reacciona según si el deal tiene sentido para su visión).

### Después del Draft

**Decisiones disponibles:**

Mandar un rookie al roster activo.
- Agentes afectados: Head Coach (tiene que integrarlo, puede frustrarse si no estaba en sus planes), Director de Player Development (pierde un jugador para trabajar), el jugador rookie (satisfecho de estar en el roster activo).

Mandar un rookie a desarrollarse.
- Agentes afectados: Director de Player Development (tiene trabajo nuevo, satisfacción sube), Head Coach (contento de no tener presión de usar al rookie), el jugador rookie (puede sentirse subestimado según su ego).

---

## OFFSEASON — Fase 2: Free Agency

**Modo:** Pausa automática + UI expandida cuando abre el mercado. Tiempo continuo antes.

La Free Agency es caótica — varios jugadores disponibles, varios equipos compitiendo, decisiones en paralelo con deadlines reales. A diferencia del Draft que tiene orden secuencial, acá todo pasa al mismo tiempo.

### Antes de que abra el mercado

**Decisiones disponibles:**

Definir qué jugadores propios querés renovar y en qué términos.
- Agentes afectados: el jugador involucrado (¿está contento con la oferta? Su historial de satisfacción acumulada durante la temporada determina su respuesta), CFO (evalúa el impacto en el cap), Director de Player Personnel (gestiona la negociación con el agente del jugador).

Identificar targets en el mercado libre.
- Agentes afectados: Director de Player Personnel (lidera la identificación), Head of Analytics (presenta modelos de valor real vs percibido de cada jugador), Head Coach (tiene su lista de perfiles que necesita su sistema — puede coincidir o no con Analytics).

### Cuando abre el mercado — UI expandida

**Decisiones disponibles:**

Hacer ofertas a jugadores libres.
- La decisión del jugador libre depende de: salario ofrecido vs su percepción de su valor de mercado, ambición de campeonato vs dinero, vínculos con la ciudad, relación previa con el GM si existió, percepción de la franquicia como lugar para ganar.
- Agentes afectados: Legal Counsel y Director de Player Personnel (gestionan la negociación), CFO (alerta si se superan umbrales), Owner (puede presionar para fichar a alguien específico por razones propias).

Pagar de más por un jugador para cerrarlo antes que un rival.
- Agentes afectados: CFO (reacción según su conservadurismo financiero), Owner (reacción según su disposición al gasto), GM rival que también lo quería (relación se enfría).

Dejar pasar un jugador que el Director de Scouting recomendaba.
- Agentes afectados: Director de Scouting (lo registra, confianza baja gradualmente si se repite el patrón).

Firmar un jugador muy mediático que el Marketing quería.
- Agentes afectados: Director de Marketing (satisfacción sube), Head Coach (evalúa si encaja en su sistema), el jugador fichado (llega con sus propias variables de ego y expectativas).

### Efectos en la ciudad

Fichar una estrella → boom inmediato en ticket sales, valor del suelo cerca del estadio sube, la Prensa genera narrativa positiva, Director de Ticket Sales reporta aumento en ventas de abonos.

Free agency decepcionante → la Prensa empieza a construir narrativa crítica, Director de Ticket Sales alerta sobre ventas flojas, Owner baja su satisfacción si esperaba movimientos importantes.

---

## OFFSEASON — Fase 3: Traspasos y armado del roster final

**Modo:** Tiempo continuo con pausa manual. Los traspasos ocurren en negociaciones que el jugador puede iniciar en cualquier momento.

Esta fase convive parcialmente con la Free Agency y se extiende hasta el Trade Deadline durante la temporada regular. En el offseason es cuando hay más flexibilidad para mover el roster completo.

### Traspasos

**Decisiones disponibles:**

Iniciar conversaciones con GMs rivales.
- Agentes afectados: GM rival (reacciona según su estilo de negociación, su urgencia actual, y su relación previa con el jugador), Director de Player Personnel (gestiona la negociación, satisfacción sube si se le da autonomía).

Ofrecer un jugador que el Coach no quiere perder.
- Agentes afectados: Head Coach (confianza en el GM baja si el traspaso se cierra sin consultarle), el jugador traspasado (relación con el GM termina), jugadores del roster que tenían relación con él.

Cerrar un traspaso que perjudica a un rival.
- Agentes afectados: GM rival (lo recuerda, próximas negociaciones más frías o directamente se niega a negociar).

Traspasar un jugador querido por la fanbase.
- Agentes afectados: Prensa (genera narrativa negativa), Director de Ticket Sales (alerta sobre posible impacto en ventas), Director de PR (entra en modo gestión de la reacción pública), la ciudad (el valor del suelo en la zona del estadio puede bajar levemente).

### Armado del roster final

**Decisiones disponibles:**

Definir los 15 del roster activo y los contratos especiales.
- Agentes afectados: Head Coach (da su opinión sobre la rotación, confianza baja si se ignora sistemáticamente), Director de Player Development (pide mantener algún juvenil con potencial), cada jugador del roster (reacciona según su rol asignado y sus expectativas).

### Efectos en la ciudad

Roster final con estrellas → expectativas de la fanbase suben, Prensa construye narrativa de esperanza, Director de Ticket Sales reporta ventas anticipadas positivas.

Roster decepcionante → Director de Ticket Sales alerta, Owner empieza a observar más de cerca, Prensa construye narrativa de escepticismo.

---

## PRETEMPORADA

**Modo:** Tiempo continuo con pausa manual. No hay UI expandida automática — es más observación que acción.

La pretemporada es la fase de transición. El roster está armado, los partidos empiezan pero no cuentan. El Coach trabaja su sistema, los jugadores se conocen, y empezás a ver si lo que armaste tiene sentido en la cancha.

### Decisiones disponibles

**Roster y rotación:**

Definir minutos y roles para cada jugador.
- Agentes afectados: cada jugador reacciona según su ego y sus expectativas. Un jugador que esperaba más minutos y no los tiene empieza a acumular insatisfacción silenciosa — no lo dice, pero el Sports Psychologist puede detectarlo.

Cortar jugadores del training camp.
- Agentes afectados: Director de Player Development (puede estar en desacuerdo si cortás a alguien con potencial que él quería trabajar), el jugador cortado, jugadores del roster que tenían relación con él.

Decidir si un rookie arranca en el roster o va a desarrollarse.
- Agentes afectados: el jugador rookie (lo siente, afecta su satisfacción inicial con la franquicia), Head Coach (tiene su preferencia sobre cuándo usar rookies).

**Sistema de juego:**

Dar libertad total al Coach para implementar su sistema.
- Agentes afectados: Head Coach (satisfacción sube, se siente respaldado), Head of Analytics (puede frustrarse si sus recomendaciones no llegan al Coach).

Meterse en decisiones tácticas.
- Agentes afectados: Head Coach (confianza en el GM baja si siente que lo micromanageás), asistentes del coach (observan la dinámica y ajustan su relación con el GM).

Proponer ajustes del Analytics team al sistema.
- Agentes afectados: Head Coach (reacciona según su relación con la analítica — puede abrirse o cerrarse), Head of Analytics (satisfacción sube si el GM lo usa como canal), Director de Scouting (observa si la analítica gana influencia sobre el scouting tradicional).

### Lo que el sistema genera solo

- Primeros reportes del Sports Psychologist sobre el clima del vestuario
- El médico reporta el estado físico de cada jugador después de los primeros entrenamientos intensos
- El S&C Coach alerta si algún jugador está siendo sobreexigido
- Primeras fricciones o buenas noticias del vestuario que emergen como eventos narrativos

### Efectos en la ciudad

La pretemporada genera expectativa — la Prensa empieza a cubrir al equipo, el Director de Ticket Sales ve si las ventas de abonos están bien. Si hay una figura nueva y emocionante en el roster, la ciudad lo siente antes del primer partido oficial.

---

## TEMPORADA REGULAR

**Modo:** Tiempo continuo con pausa manual. El calendario de partidos siempre visible. Trade Deadline con pausa automática + UI expandida.

La fase más larga — aproximadamente 82 partidos entre octubre y abril. Cada semana simulada tiene entre 2 y 4 partidos.

### Decisiones de gestión del roster — ongoing

Gestionar minutos y carga de jugadores.
- Agentes afectados: S&C Coach y médico alertan si algún jugador necesita descanso. Ignorarlos aumenta riesgo de lesión. Hacerles caso puede costar partidos en ventanas importantes.

Responder cuando un jugador pide más rol.
- Opciones: darle más minutos, explicarle el plan, ignorarlo.
- Agentes afectados: el jugador (satisfacción según la respuesta), Head Coach (puede estar de acuerdo o no con el cambio de rotación), Sports Psychologist (observa el estado emocional resultante).

Gestionar una lesión.
- El médico da diagnóstico y tiempo estimado. El GM puede presionar para que el jugador vuelva antes.
- Agentes afectados: médico (puede ceder o resistir según su personalidad y relación con el GM), el jugador lesionado (si vuelve antes y se reagrava, su confianza en la organización cae dramáticamente), Head Coach (quiere al jugador disponible).

Responder cuando un jugador está en mal momento emocional.
- El Sports Psychologist alerta. Opciones: intervenir directamente, dejar que el psicólogo lo maneje, ignorarlo.
- Agentes afectados: el jugador (evolución de su estado emocional según la respuesta), Sports Psychologist (satisfacción según el nivel de autonomía que se le da), Head Coach (necesita saber si el jugador está disponible mentalmente).

### Decisiones de mercado — ongoing

**Trade Deadline — pausa automática + UI expandida:**

El momento de mayor actividad de traspasos en la temporada. Los GMs rivales están activos simultáneamente — hay oportunidades y presiones en paralelo.

- Decisiones equivalentes a las de la fase de traspasos del offseason pero con mayor urgencia temporal
- Agentes afectados: los mismos que en traspasos, con mayor intensidad emocional por la presión del deadline

**Waiver wire:**

Jugadores cortados por otros equipos disponibles. El Director de Player Personnel alerta si hay alguien interesante.
- Agentes afectados: Director de Player Personnel (satisfacción si se actúa sobre sus recomendaciones), Head Coach (evalúa si el jugador encaja en su sistema).

**Contratos de corto plazo:**

Podés firmar jugadores por períodos cortos para cubrir lesiones o probar opciones.
- Agentes afectados: CFO (evalúa el impacto presupuestario), Head Coach (tiene su opinión sobre el jugador específico).

### Decisiones de negocio — ongoing

Extensiones de contrato con jugadores propios.
- El agente del jugador contacta al GM cuando el contrato está próximo a expirar. Si se espera demasiado el precio sube o el jugador decide explorar el mercado.
- Agentes afectados: el jugador (su satisfacción acumulada determina qué tan receptivo es a renovar y a qué precio), Director de Player Personnel (gestiona la negociación), CFO (evalúa el impacto en el cap a largo plazo), Legal Counsel (estructura el contrato).

Aceptar o rechazar propuestas de sponsors.
- El Director de Corporate Partnerships trae oportunidades con condiciones específicas.
- Agentes afectados: Director de Corporate Partnerships (satisfacción según si se actúa sobre sus deals), Director de Marketing (evalúa alineación con la imagen de la franquicia), jugadores involucrados si el deal requiere su participación.

Aprobar eventos alternativos en el estadio.
- El Director de Arena Operations propone eventos en fechas sin partido.
- Agentes afectados: Director de Arena Operations (satisfacción según la autonomía que se le da), CFO (evalúa ingresos proyectados), Director de PR (evalúa si los eventos están alineados con la imagen de la franquicia).

### Decisiones de ciudad — ongoing

- Construís servicios o edificios especiales en cualquier momento de la temporada
- El Alcalde aparece con propuestas o bloqueos según el modo de gestión elegido
- La economía de la ciudad fluctúa semana a semana según resultados del equipo

### Lo que el sistema genera solo

Los agentes interactúan entre sí independientemente del GM:
- El Coach y el médico tienen su tensión propia sobre disponibilidad de jugadores
- El Director de Marketing presiona al Director de PR por el manejo de algún jugador mediático
- El Alcalde puede llamar al Owner directamente sobre algún proyecto de ciudad
- Los jugadores acumulan estado emocional. Un jugador insatisfecho con sus minutos empieza a hablar con compañeros — eso eventualmente emerge como evento narrativo

La ciudad reacciona en tiempo real:
- Partidos ganados en casa → más ciudadanos cerca del estadio el día siguiente
- Racha mala → el valor del suelo en la zona comercial del estadio empieza a bajar lentamente
- Racha ganadora → ticket sales suben, la Prensa construye narrativa positiva que atrae sponsors

### Antes y después de cada partido

**Antes:**
- Reporte del Video Coordinator sobre el próximo rival
- Opinión del Coach sobre el plan de juego
- Estado físico actualizado del roster

**Después:**
- Resultado y stats
- Reacción de la Prensa
- Estado emocional del vestuario según el resultado
- Si perdiste feo, el Owner puede aparecer

---

## PLAYOFFS

**Modo:** Tiempo continuo con pausa manual entre partidos. La presión sube para todos los agentes simultáneamente.

Bracket de 16 equipos, series al mejor de 7, cuatro rondas hasta el campeonato. El bracket completo siempre visible. Todo lo que construiste durante la temporada regular se pone a prueba acá — las relaciones con los agentes, el estado emocional del roster, las decisiones de carga que tomaste o no tomaste.

### Lo que cambia respecto a la temporada regular

- El Owner aparece más frecuentemente — cada derrota en una serie lo activa
- La Prensa intensifica su cobertura — cada partido genera narrativa nueva
- El estado emocional de los jugadores fluctúa más rápido — una actuación mala en playoffs pesa más que diez en temporada regular
- Las relaciones entre jugadores se tensan o se fortalecen bajo presión — emergen como eventos narrativos más frecuentes

### Decisiones disponibles

**Gestión de la serie:**

Responder a ajustes tácticos propuestos por el Coach entre partidos.
- Agentes afectados: Head Coach (confianza según el nivel de apoyo recibido), Head of Analytics (puede proponer ajustes alternativos basados en datos del rival), asistentes del coach.

Gestionar la rotación acortada.
- En playoffs el Coach naturalmente acorta la rotación. Jugadores fuera de la rotación acumulan insatisfacción más rápido que en temporada regular.
- Agentes afectados: jugadores excluidos de la rotación (insatisfacción acumulada rápida), Sports Psychologist (monitorea el estado del vestuario activamente).

Gestionar carga en series largas.
- Si una serie llega a game 6 o 7 los jugadores están agotados. Médico y S&C Coach alertan con más urgencia.
- Agentes afectados: médico (sus recomendaciones tienen más peso — si un jugador se lesiona por ignorarlas, la relación con el GM cae drásticamente), S&C Coach, el jugador en cuestión, Head Coach.

**Gestión emocional:**

Responder cuando un jugador explota bajo presión.
- El Sports Psychologist alerta antes de que sea público. Opciones: intervenir directamente con el jugador, dejar que el Coach lo maneje, pedirle al psicólogo que trabaje con él.
- Cada camino tiene consecuencias distintas en el jugador, el Coach, y el psicólogo.

Responder a un conflicto en el vestuario.
- Series largas generan tensiones entre jugadores. El Coach puede manejarlo o no según su habilidad de liderazgo.
- Agentes afectados: los jugadores involucrados, Head Coach (su manejo del conflicto afecta su credibilidad con el vestuario completo), Sports Psychologist.

Gestionar la narrativa de un jugador estrella en racha negativa.
- La Prensa lo destruye. Su estado emocional cae. Opciones: defenderlo públicamente, guardar silencio, pedir al Director de PR que maneje la narrativa.
- Agentes afectados: el jugador (su estado emocional según el apoyo recibido), Director de PR (satisfacción si se le da autonomía para gestionar), Prensa (reacciona a la postura del GM), fanbase.

**Gestión mediática:**

Dar conferencias de prensa.
- El sistema presenta opciones de respuesta. Cada una construye una narrativa distinta con la Prensa.
- Agentes afectados: Prensa (narrativa resultante), Owner (evalúa si el GM se comunica bien públicamente), Director de PR (satisfacción según si el GM sigue sus recomendaciones o improvisa).

### Lo que el sistema genera solo

- Los GMs rivales reaccionan a ser eliminados — esa relación se afecta
- La ciudad reacciona partido a partido con mayor intensidad que en temporada regular
- El bracket avanza solo mientras el tiempo corre

### Los tres desenlaces

**Eliminación temprana:**
El Owner convoca una conversación seria. Según su paciencia restante puede ser advertencia, exigencia de cambios, o despido. La ciudad entra en modo decepción — emigración leve, ticket sales en riesgo, Prensa construye narrativa de fracaso. Los jugadores con contratos expirando evalúan si quieren quedarse.

**Final sin campeonato:**
El Owner evalúa según sus expectativas iniciales — si eras equipo de reconstrucción llegar a la final es un éxito; si eras favorito es una decepción. La ciudad celebra haber llegado lejos, el impacto económico es positivo aunque menor que un campeonato. Los jugadores sienten el sabor de la cercanía — puede aumentar su ambición de campeonato y su disposición a quedarse.

**Campeonato:**
El evento de mayor impacto de todo el juego. Un terremoto en todos los sistemas simultáneamente:
- La ciudad explota económicamente — valor del suelo sube en toda el área del estadio, nuevos ciudadanos llegan, zona comercial se expande
- El Owner da presupuesto expandido y confianza renovada
- Los jugadores suben su lealtad — quieren defender el título
- Los sponsors hacen cola — el Director de Corporate Partnerships tiene el mejor momento de su carrera
- La Prensa construye narrativa de dinastía que tiene peso en las próximas temporadas
- Los agentes de jugadores free agents empiezan a llamar — la franquicia es un destino atractivo
- Los GMs rivales te ven diferente en las negociaciones

---

## CIERRE DE TEMPORADA

**Modo:** Pausa automática + UI expandida. Momento de reflexión y decisiones que definen el punto de partida de la próxima temporada.

### Evaluación formal con el Owner

La conversación más importante de esta fase. El Owner convoca al GM — no es opcional. El tono depende de tres factores: el resultado deportivo, el estado económico de la franquicia, y cuánta paciencia le queda según su perfil inicial.

**Posibles salidas:**

Confianza renovada → más presupuesto, más libertad, horizonte claro para la próxima temporada.

Advertencia formal → sigue confiando pero con condiciones explícitas. El jugador sabe exactamente qué tiene que lograr la próxima temporada.

Presión con ultimátum → una temporada más para mostrar resultados concretos. El Owner pone métricas específicas.

Despido → game over para este GM. El juego ofrece empezar una nueva partida o cargar un save anterior.

### Evaluación del staff

Revisión del rendimiento de cada miembro del staff con sus métricas calculadas por el sistema durante la temporada:
- Precisión histórica del Director de Scouting
- Retention rate de sponsors del Director de Corporate Partnerships
- Track record del Director de Player Development
- Efectividad del Sports Psychologist (lesiones evitadas, situaciones manejadas)

**Decisiones disponibles:**

Renovar o no renovar al Coach.
- Si llegó lejos en playoffs su posición es fuerte. Si fue eliminado temprano puede ser vulnerable.
- Agentes afectados: el Coach (su reacción según cómo se manejó la conversación), asistentes del coach (evalúan su propia seguridad), jugadores (algunos tienen relación fuerte con el coach y reaccionan a su salida), Owner (tiene su propia opinión sobre el Coach).

Hacer cambios en el front office o cuerpo técnico.
- Agentes afectados: cada despido afecta la relación con el resto del staff — algunos lo ven como señal de que nadie está seguro, lo cual sube la ansiedad general del equipo.

### Jugadores con contratos expirando

Antes de que abra la free agency hay una ventana para retener a los propios jugadores en condiciones más favorables.

La decisión de cada jugador de quedarse o explorar el mercado depende de su satisfacción acumulada durante toda la temporada — no solo el resultado final. Un jugador que acumuló insatisfacción durante meses no se queda aunque ganen el campeonato si siente que no fue valorado.

Agentes afectados: el jugador involucrado, su agente, Director de Player Personnel (gestiona las negociaciones), CFO (evalúa el impacto en el cap de la próxima temporada).

### La ciudad procesa la temporada

Los efectos económicos de los playoffs se estabilizan — el boom post-campeonato o la caída post-eliminación encuentran su nuevo nivel base.

El valor del suelo, la asistencia proyectada, y los ingresos de sponsors se recalculan para la próxima temporada. La Prensa cierra su narrativa de la temporada y empieza a construir expectativas para la siguiente — esa narrativa llega a la próxima temporada con peso real.

---

## Efectos cruzados por fase — resumen

| Decisión | Fase | Efectos ciudad | Efectos agentes |
|---|---|---|---|
| Draft exitoso (prospecto recomendado) | Offseason | Expectativa fanbase sube | Scouting +confianza, Coach evalúa |
| Free agency con estrella | Offseason | Boom ticket sales, suelo sube | Marketing +satisfacción, CFO alerta |
| Free agency decepcionante | Offseason | Ticket sales flojas | Prensa narrativa negativa, Owner observa |
| Traspaso jugador querido | Offseason / TR | Suelo baja leve, Prensa negativa | Coach -confianza, PR en gestión |
| Racha ganadora | Temporada Regular | Suelo sube, ciudadanos activos | Owner +satisfacción, Prensa positiva |
| Racha perdedora | Temporada Regular | Suelo baja, ticket sales caen | Owner -paciencia, Prensa narrativa crítica |
| Ignorar médico → lesión | Temporada Regular | — | Médico -confianza, jugador -lealtad |
| Campeonato | Playoffs | Boom económico total | Todos los agentes +satisfacción, nuevos sponsors |
| Eliminación temprana | Playoffs | Decepción, emigración leve | Owner evalúa despido, jugadores reconsideran |
| Despido de staff | Cierre | — | Ansiedad general sube en todo el equipo |
| Renovación del Coach | Cierre | Estabilidad percibida | Coach +lealtad, jugadores +seguridad |

---

*Documento generado en sesión de diseño — Mayo 2026.*
*Complementa PulseCity_v2.md y PulseCity_Agentes.md — leer en conjunto.*
