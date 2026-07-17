# M3.P1c-b — Playtest funcional de trades y medicina

Fecha: 2026-07-17

## Objetivo

Validar manualmente desde el frontend que la experiencia inicial corregida, Trade Center y Centro Medico puedan recorrerse sin asistencia tecnica, y registrar claridad, friccion y game feel antes de avanzar hacia Trade Deadline.

## Recorrido realizado

- se inicio la temporada desde el kickoff nuevo
- se navego desde el Command Center hacia las mecanicas necesarias sin perder contexto
- se realizaron dos propuestas desde Trade Center
- llegaron contraofertas y se acepto una negociacion
- el roster se actualizo con los jugadores recibidos
- se respondio una decision medica desde Centro Medico
- se simulo la temporada regular y luego se dejo avanzar el calendario hasta el 30 de julio para observar la transicion de fase

El criterio original indicaba al menos tres propuestas. El muestreo manual concreto fue de dos; se conserva este dato sin reinterpretarlo. Aun asi, el recorrido cubrio propuesta, contraoferta, aceptacion y mutacion de roster, que eran los estados funcionales principales a validar.

## Resultado de experiencia

La navegacion fue intuitiva. El GM supo donde entrar cada vez que necesito operar una mecanica. El Command Center funciono como resumen y punto de acceso, mientras Trade Center y Centro Medico conservaron un proposito propio y entendible.

Las negociaciones se sintieron claras dentro de su profundidad actual: se propone, el rival contraoferta y la aceptacion actualiza el roster. No se evalua todavia una negociacion NBA profunda ni rosters rivales nominales, porque estan fuera del alcance actual.

La decision medica pudo responderse correctamente desde el frontend. El flujo fue comprensible y no requirio scripts, `curl` ni inspeccion tecnica.

## Hallazgos

### 1. Identidad del roster inicial

Los jugadores iniciales aparecen en distintas superficies mediante un hash o `player_id` en lugar de nombre y apellido. Los jugadores recibidos por trade si muestran nombre.

La causa es el contrato de datos actual: el `roster.patch` emocional inicial no incluye identidad contractual, mientras el parche emitido al aceptar un trade incluye `full_name`, posicion, rating y salario. El frontend mezcla ambos tipos de parche y solo puede mostrar el nombre cuando lo recibio.

Esto es un bug visible, no una decision de diseño. Se crea `M3.P1d-a` para corregirlo antes del refactor de `M3.P2`.

### 2. Utilidad y estado de los agentes

Muchos agentes todavia tienen poca utilidad mecanica. Parte de esto es esperable: el catalogo completo fue sembrado, pero eventos criticos, narrativa ampliada, intensidad de playoffs y actividad integral de los 30 agentes siguen pendientes en `M3.14`, `M3.15`, `M3.17` y `M3.24`.

Sin embargo, el chat tampoco refleja de forma confiable el estado visible del agente. Un agente puede aparecer como `frustrated`, `calm`, `pressured` o `concerned` y responder sin incorporar ese estado.

La investigacion encontro dos fuentes distintas: los moods dinamicos de los agentes core se persisten en `agent_core_states`, mientras `narrative-service` ensambla el chat leyendo `agent_individual_states`. Esa segunda capa puede conservar el estado inicial sembrado. Se crea `M3.P1d-b` para unificar el estado que ve el GM con el que recibe el prompt, sin agregar mecanicas nuevas.

El contexto operativo del chat tambien es acotado: identidad, dominio, estado, relacion con GM y decisiones recientes. La futura profundidad de cada agente debera incorporar informacion real de su dominio sin permitir que el LLM decida gameplay.

### 3. Fin silencioso de temporada regular

El calendario actual contiene 82 partidos sinteticos. Una vez agotados, el reloj puede continuar hasta julio sin Trade Deadline, Playoffs, All-Star, Draft, cierre ni una explicacion de fase.

Trade Deadline, Playoffs y Cierre pertenecen a `M3.13` y `M3.16`-`M3.19`. Draft y Free Agency estan diferidos a `M3.5`. NBA Cup y All-Star no tienen todavia un mini milestone asignado.

No se implementa postseason como correccion de este playtest. El hallazgo queda como criterio obligatorio de la futura transicion: al terminar la temporada regular, el juego no puede continuar silenciosamente fuera de fase; debe pausar o comunicar claramente el siguiente estado.

## Decision

No se pasa directamente a `M3.P2`. Primero se ejecutan dos microcorrecciones:

1. `M3.P1d-a` — identidad visible del roster inicial.
2. `M3.P1d-b` — estado dinamico coherente en el contexto del chat.

Luego se retoma `M3.P2` como refactor puro. El final silencioso de temporada queda registrado para la implementacion de la transicion a Playoffs, sin ampliar el alcance de estas correcciones.

## Archivos modificados en esta sesion

- `docs/Sesiones/MILESTONE3/INICIOM3.MD`
- `docs/02_Progress/2026-07-17_M3_P1c_b_Playtest_Funcional.md`

## Pruebas

No se modifico codigo. La verificacion de esta sesion fue el playtest manual informado por el GM y la inspeccion del flujo de datos relevante.
