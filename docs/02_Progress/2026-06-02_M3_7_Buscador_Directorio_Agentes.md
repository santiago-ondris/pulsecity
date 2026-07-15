# 2026-06-02 — M3.7 Buscador y directorio de agentes

## Objetivo

Hacer que llegar a un agente sea trivial desde la pantalla jugable inicial.

## Cambios

- Se amplio la tab `Agentes` del Command Center.
- Se agrego buscador por nombre, rol, dominio o `agent_id`.
- Se agregaron filtros:
  - Todos
  - Basketball
  - Business
  - Ciudad
  - Roster
- Se incorporo al frontend el catalogo canonico de agentes individuales sembrado en `agent-service`.
- El directorio cubre Basketball Ops, Business Ops, Ciudad y Prensa.
- El roster se suma dinamicamente cuando llegan estados emocionales por `roster.patch`.
- Seleccionar cualquier entrada del directorio cambia el agente activo y el chat apunta a ese `agent_id`.
- El chat sigue usando el provider real de `M3.6`.

## Limite conocido

Los jugadores del roster se muestran por `player_id`, porque el estado emocional que llega al frontend no incluye nombre ni posicion contractual. Resolver eso requiere que `gateway` o `team-service` expongan/combinen datos de roster para la Vista Franquicia.

## Verificacion

- `npm run build --prefix frontend`

## Resultado

`M3.7 — Buscador de agentes y vista de directorio` queda cerrado funcionalmente.
