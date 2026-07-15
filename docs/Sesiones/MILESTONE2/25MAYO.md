# Sesion 25 Mayo 2026 — Cierre M2 y DX local

## Objetivo de la sesion

Cerrar la sesion despues de M2 y dejar registrada la duda practica: como levantar PulseCity localmente para probar la app desde el punto de vista de un usuario real.

## Estado al iniciar

M2 quedo funcionalmente cerrado el 24 de mayo:

- 82 partidos simulados en una temporada completa
- ciudad reaccionando a resultados
- 5 agentes core con estado real
- narrativa post-partido generada
- analytics basico persistiendo series
- frontend mostrando HUD, resultados recientes, inbox narrativo, ciudad y agentes

## Servicios necesarios para probar la app hoy

Infraestructura:

```bash
make up
```

Servicios de app, actualmente en terminales separadas:

```bash
make run-gateway
make run-map-service
make run-team-service
make run-match-service
make run-city-service
make run-narrative-service
make run-analytics-service
make run-frontend
```

Para que avance el tiempo de una partida concreta:

```bash
GAME_ID=<id-de-la-partida> make run-agent-service
```

Detalle importante:

`agent-service` todavia esta atado a un `GAME_ID` por proceso. Para probar desde browser, el flujo actual es crear o abrir una partida, copiar el id visible en la pantalla de ceremonia y levantar `agent-service` con ese id.

## Que puede hacer un usuario real hoy

Un usuario puede:

- entrar como invitado
- registrarse o loguearse
- migrar partidas guest a una cuenta
- crear una nueva franquicia/partida
- elegir ciudad, nombre de franquicia, escenario inicial y modo de gestion
- ver la ceremonia de generacion del mapa
- recibir el primer evento narrativo del Owner
- responder la direccion inicial como GM
- entrar a la pantalla viva de M2
- controlar pausa y velocidad x1/x5/x20
- ver avanzar fecha simulada
- mirar record de temporada
- ver resultados recientes
- ver metricas urbanas basicas
- ver estado resumido de los 5 agentes core
- leer inbox narrativo post-partido

Todavia no puede:

- fichar jugadores
- hacer trades
- modificar roster
- construir edificios
- conversar libremente con agentes
- tomar decisiones profundas de GM
- ver dashboards historicos de analytics

## Decision tomada

Se agrego un milestone final propuesto a `INICIOM2.MD`:

```text
M2.16 — DX local para probar la app — PENDIENTE
```

La idea es mejorar la experiencia local antes de M3, agregando un comando claro para levantar el entorno jugable y reduciendo el acople manual de `agent-service` con `GAME_ID`.

## Pendiente siguiente

Primer mini milestone recomendado para la proxima sesion:

Implementar `M2.16 — DX local para probar la app`.

Objetivo chico:

- agregar un comando `make dev-app` o `make dev-full`
- documentar el flujo exacto de prueba local
- definir si `agent-service` sigue con `GAME_ID` manual por ahora o si pasa a manejar partidas activas de forma mas ergonomica

## Cierre

M2 queda cerrado a nivel funcional. Lo que falta es mejorar la ergonomia de desarrollo local para poder probarlo seguido sin pasos manuales innecesarios.
