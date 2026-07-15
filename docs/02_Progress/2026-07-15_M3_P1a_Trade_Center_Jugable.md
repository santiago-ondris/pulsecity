# M3.P1a — Trade Center jugable

## Objetivo de la sesion

Hacer operable el flujo de trades desde el frontend sin seguir acumulando mecanicas dentro de la pagina principal.

## Decision de UX permanente

El Command Center pasa a entenderse como resumen y punto de navegacion. Las mecanicas con flujo propio viven en paginas o espacios de trabajo propios; no se agregan como otra pestaña, modal o bloque por conveniencia tecnica.

La regla quedo registrada en `AGENTS.md`. Trade Center es la primera aplicacion formal. Centro Medico seguira el mismo principio en `M3.P1b`.

## Implementacion

- nueva ruta `/franchise/trades`
- nueva pagina `TradeCenterPage`
- acceso Command Center → Trade Center y retorno explicito
- propuesta de trade con jugador propio, rival, posicion e incoming salary
- catalogo frontend de los 30 equipos y GMs rivales canonicos
- estado local de negociaciones por `proposal_id`
- aplicacion incremental de `trade.patch`
- estados visibles: enviada, contraoferta, rechazada y aceptada
- aceptacion de contraoferta con su asset adicional
- loading y error feedback junto a cada accion
- hook de dominio `useTradeOperations`
- CSS de dominio separado y responsive

## Ownership y flujo

El frontend no valida gameplay ni muta roster/cap por su cuenta. Publica decisiones mediante el gateway y espera los deltas WebSocket. `team-service` sigue validando roster y cap; `agent-service` sigue evaluando el perfil del GM rival.

## Verificacion

- `npm run build --prefix frontend`
- `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`
- `make build`

Todos finalizaron correctamente.

## Limite conocido

El selector usa jugadores recibidos en vivo por `roster.patch`. La rehidratacion REST contractual completa queda para `M3.22`; esta UI no inventa una segunda fuente de verdad.

## Pendiente siguiente

`M3.P1b`: Centro Medico como pagina propia, disponibilidad del roster y las cuatro decisiones medicas. Despues, `M3.P1c` ejecuta el playtest conjunto y documenta game feel para `M3.13`.
