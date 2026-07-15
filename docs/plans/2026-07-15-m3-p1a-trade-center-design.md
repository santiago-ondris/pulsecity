# M3.P1a — Trade Center

## Objetivo

Hacer operativo el flujo de trades desde el frontend sin sumar otra mecánica al panel lateral del Command Center.

## Decisión de navegación

`Trade Center` vive en una página propia. El Command Center conserva su rol de resumen y ofrece una entrada visible al espacio de negociación. La página tiene retorno explícito al Command Center y mantiene el contexto de la partida activa.

Esta decisión inaugura una regla permanente: las mecánicas con flujo propio no se acumulan como pestañas, modales o bloques dentro de una única pantalla. Centro Médico, roster, finanzas, calendario y playoffs podrán seguir el mismo patrón cuando su alcance lo justifique.

## Arquitectura frontend

- Ruta: `/franchise/trades`.
- Página presentacional: `TradeCenterPage`.
- Estado de dominio y acciones REST: `useTradeOperations`.
- Contrato en cliente: `trade.patch`, indexado por `proposal_id` y aplicado como delta.
- Acciones: proponer trade y aceptar contraoferta.
- El hook general conserva la conexión WebSocket, pero delega el estado y las mutaciones de trades al hook de dominio.
- No se agrega dependencia nueva.

## Experiencia

La página separa dos tareas:

1. Mesa de propuesta: jugador propio, rival, posición solicitada e incoming salary.
2. Negociaciones: historial vivo con estado `proposed`, `countered`, `rejected` o `accepted`.

Una contraoferta muestra el asset adicional requerido y permite aceptarlo. El servidor sigue siendo la autoridad sobre roster, salary cap y evaluación rival. La UI muestra estados de envío, errores cercanos a la acción y el resultado recibido por WebSocket.

## Criterio de done

- se puede proponer un trade sin `curl`
- los estados llegan por `trade.patch` y son legibles
- una contraoferta se puede aceptar desde la página
- existe navegación clara Command Center ↔ Trade Center
- el frontend compila sin errores
- el playtest de tres propuestas queda para `M3.P1c`, después de incorporar Centro Médico en `M3.P1b`
