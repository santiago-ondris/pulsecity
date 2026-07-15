# 2026-06-02 — M3.7a Ceremony Command Center

## Objetivo

Corregir la base UX de `new-game/ceremony` antes de avanzar con el buscador completo de agentes.

## Problema detectado

La pantalla mezclaba mapa, tiempo, temporada, resultados, inbox, agentes, chat, pipeline y eventos recientes en una sola columna derecha. El mapa quedaba empujado verticalmente y el chat no era evidente para el jugador.

## Cambios

- `CeremonyPage` quedo como orquestador de estado local y layout.
- Se creo `frontend/src/features/new-game/components/ceremony/`.
- Se separaron componentes para topbar, mapa, panel lateral, agentes, chat, inbox, temporada y sistema.
- El mapa vuelve a ser el foco del primer viewport.
- El panel derecho usa tabs:
  - `Agentes`
  - `Inbox`
  - `Temporada`
  - `Sistema`
- El chat queda integrado a la tab `Agentes`.
- `Pipeline` y `Eventos recientes` quedan en `Sistema`, con jerarquia secundaria.
- El badge del chat ahora dice `LLM real`, no `Stub M3.5`.
- Se agregaron reglas responsive basicas para desktop/tablet/mobile.

## Continuidad con M3.7

Despues de este reset, se completo `M3.7 — Buscador de agentes y vista de directorio` en la misma sesion:

- la tab `Agentes` recibio buscador por nombre, rol, dominio o `agent_id`
- se agregaron filtros por categoria: Todos, Basketball, Business, Ciudad y Roster
- el frontend incorporo el catalogo canonico de agentes individuales sembrado en `agent-service`
- seleccionar cualquier agente actualiza el foco y abre chat contra ese `agent_id`
- el roster se suma dinamicamente cuando llegan `roster.patch`

## Verificacion

- `npm run build --prefix frontend`

## Resultado

`M3.7a` fue el paso de UX necesario para que `M3.7` no agregara mas complejidad sobre una pantalla saturada. `M3.7` quedo cerrado funcionalmente despues de este reset.
