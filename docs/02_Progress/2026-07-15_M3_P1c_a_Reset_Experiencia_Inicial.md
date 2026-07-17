# M3.P1c-a — Reset de experiencia inicial

## Origen

El primer intento de playtest de `M3.P1c` se detuvo en la pantalla inicial. La captura a 1920 px mostro que el problema no era todavia el game feel de trades o medicina:

- no habia una accion dominante para comenzar la temporada
- los controles `Pausa / x1 / x5 / x20` no explicaban el loop
- socket y progreso del mapa ocupaban la jerarquia principal
- el directorio completo de agentes aparecia abierto por defecto
- la grilla del mapa se extendia por debajo del panel lateral
- existia scroll horizontal incluso en desktop ancho
- textos de ceremonia/debug seguian visibles en una partida lista

Decision: pausar el playtest funcional y corregir la base antes de evaluar mecanicas.

## Entrada guiada

Una partida nueva entra pausada en un kickoff con:

- identidad de franquicia
- mandato inicial del Owner
- progreso de temporada
- cap disponible o estado de espera del CFO
- roster observado y alertas
- CTA dominante `Comenzar temporada`
- revisiones opcionales de Trade Center, Centro Medico y Staff

El CTA llama al endpoint de control temporal existente con velocidad x1. No se agrega un nuevo estado de temporada. Si el request falla, `updateTimeControl` restaura el estado previo y el kickoff vuelve con error visible.

## Estados de entrada

- mapa en generacion: ceremonia contenida con progreso y snapshot
- mapa completo + Owner respondido + cero dias/partidos: kickoff guiado
- partida avanzada: Command Center operativo directo

## Command Center operativo

- `Resumen` reemplaza a `Agentes` como tab inicial
- tabs: `Resumen`, `Inbox`, `Staff`, `Sistema`
- Staff contiene el buscador, directorio y chat existentes
- Sistema contiene pipeline y eventos tecnicos
- el area principal muestra pulso de temporada, metricas, operaciones, resultados y snapshot de ciudad
- Trade Center y Centro Medico siguen viviendo en paginas propias

## Correccion de layout

- el Command Center ya no usa margenes negativos para ocupar `100vw`
- ancho maximo comun de 1600 px
- columnas con `minmax(0, ...)` y children con `min-width: 0`
- snapshot con proporcion fija y recorte interno
- celdas sin `aspect-ratio` capaz de imponer ancho al grid
- sidebar contenido y sticky solo en desktop
- colapso a una columna en tablet/mobile
- targets de interaccion de 44 px en controles principales
- estilos de reduced motion

## Estructura

Componentes nuevos:

- `SeasonKickoffPanel`
- `WorldGenerationPanel`
- `CommandCenterOverview`
- `CommandSummaryPanel`
- `CitySnapshotCard`

Componentes reemplazados y eliminados:

- `CeremonyMapPanel`
- `SeasonPanel`

## Verificacion

- `npm run build --prefix frontend`
- `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`
- `make build`
- `git diff --check`
- auditoria estructural de CSS sin `100vw`, min-width fijos grandes ni margenes negativos en el nuevo layout

Todos finalizaron correctamente.

El entorno no dispone de navegador headless. La validacion visual final se realiza al recargar la app en el navegador del GM; cualquier ajuste que surja pertenece a este mismo corte.

## Pendiente siguiente

`M3.P1c-b`: retomar el playtest desde el nuevo kickoff, completar tres propuestas de trade y responder decisiones medicas para evaluar game feel real.
