# M3.P1b — Centro Médico jugable

## Objetivo

Hacer operativo el flujo de decisiones médicas desde el frontend sin convertir el Command Center en una acumulación de mecánicas ni adelantar la rehidratación completa de roster prevista para `M3.22`.

## Decisión de navegación

`Centro Médico` vive en una página propia con ruta `/franchise/medical`. El Command Center conserva su rol de resumen y ofrece una entrada visible al espacio médico. La página mantiene el contexto de la partida activa y un retorno explícito al Command Center, siguiendo la regla inaugurada por Trade Center en `M3.P1a`.

## Arquitectura frontend

- Página presentacional: `MedicalCenterPage`.
- Estado de envío y acción REST: `useMedicalOperations`.
- Tipos de la acción médica: módulo propio del dominio médico.
- Fuente de estado jugable: `roster.patch`, ya aplicado en `rosterStates`.
- Acción: `POST /api/v1/games/{gameID}/medical-decisions`.
- El hook general conserva la conexión WebSocket, pero delega la mutación médica al hook de dominio.
- No se agrega dependencia ni endpoint nuevo.

El frontend no diagnostica, calcula riesgos ni cambia disponibilidad. Solo presenta el estado publicado por backend y registra la decisión del GM. `team-service` sigue siendo autoridad sobre lesiones y roster; `agent-service`, sobre sus consecuencias relacionales.

## Experiencia

La página se organiza como una sala médica operativa, con dos lecturas inmediatas:

1. Resumen del plantel: disponibles, lesionados y casos graves.
2. Casos activos: jugador, severidad, días estimados, fecha de retorno y cuatro decisiones posibles.

Las acciones se expresan con su consecuencia inmediata conocida: seguir protocolo, reducir carga al volver, ignorar la recomendación o forzar el alta. Las opciones de riesgo usan los colores semánticos del design system y no se presentan como equivalentes visualmente.

Si todavía no llegó estado de roster, la estructura permanece estable y explica que espera deltas. Si no hay lesiones, se muestra un estado tranquilo y el roster disponible sigue visible como contexto. Los errores aparecen junto al caso afectado y los botones quedan bloqueados solo mientras se envía la decisión de esa lesión.

## Criterio de done

- existe navegación clara Command Center ↔ Centro Médico
- disponibilidad, severidad y retorno estimado son legibles
- las cuatro decisiones pueden enviarse sin `curl`
- el frontend no inventa estado ni modifica disponibilidad de forma optimista
- el dominio médico no agrega más estado operativo al cuerpo de `useNewGameFlow`
- compilan frontend y gateway, y pasan los tests del gateway
- `M3.P1c` conserva el playtest conjunto de medicina y trades
