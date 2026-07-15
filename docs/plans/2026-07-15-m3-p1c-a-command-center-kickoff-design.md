# M3.P1c-a — Reset de experiencia inicial del Command Center

## Hallazgo del playtest

La primera captura del playtest de `M3.P1c` mostro que la pantalla inicial no comunica el loop jugable. No existe una accion dominante, los controles de tiempo no explican que reanudan la simulacion, el contenido tecnico ocupa jerarquia principal, el directorio completo de agentes aparece sin contexto y el mapa desborda su columna hasta superponerse con el panel lateral.

El playtest funcional de trades y medicina se pausa hasta corregir esta base. Continuar sobre la pantalla actual contaminaria el feedback de game feel con problemas de navegacion y layout.

## Entrada guiada

Una partida nueva entra pausada en un estado de kickoff dentro del Command Center.

La pantalla presenta:

- identidad y contexto minimo de la franquicia
- mandato inicial del Owner
- estado financiero disponible
- progreso de temporada y alertas pendientes
- CTA dominante `Comenzar temporada`
- accesos secundarios a Trade Center, Centro Medico y Staff

`Comenzar temporada` usa el control temporal existente con `paused: false` y `speed: 1`. No agrega un estado de gameplay ni un endpoint nuevo. Si la solicitud falla, el kickoff vuelve a mostrarse y comunica el error.

El kickoff aparece cuando la partida tiene cero dias procesados y cero partidos jugados. Una partida avanzada entra directamente al Command Center operativo.

## Command Center operativo

La vista inicial pasa de `Agentes` a `Resumen`.

La navegacion secundaria queda formada por:

- `Resumen`: pulso de temporada, decisiones pendientes, operaciones y snapshot de ciudad
- `Inbox`: narrativa y alertas
- `Staff`: buscador, directorio y chat de agentes
- `Sistema`: pipeline, socket y eventos tecnicos

Trade Center y Centro Medico conservan sus paginas propias.

## Layout

- el Command Center deja de usar expansion full-viewport basada en margenes negativos
- todos los tracks de grid y flex children usan limites que permiten encogerse
- el mapa vive dentro de un frame con proporcion fija y `overflow: hidden`
- el panel lateral nunca se superpone con el mapa
- desktop usa dos columnas; tablet y mobile, una columna
- no debe existir scroll horizontal a 375, 768, 1024, 1440 ni 1920 px

La grilla CSS del mapa se trata como snapshot contextual interino. La Vista Ciudad profunda con Three.js permanece en M4.

## Criterio de done

- la accion inicial se entiende en menos de tres segundos
- `Comenzar temporada` activa x1 mediante el flujo real
- un fallo restaura el kickoff con feedback visible
- `Resumen` es la seccion predeterminada y Staff se abre bajo demanda
- el contenido tecnico queda confinado a Sistema
- mapa y paneles permanecen contenidos sin overflow horizontal
- frontend y build general compilan limpios
- el hallazgo y la correccion quedan registrados en la documentacion de M3
