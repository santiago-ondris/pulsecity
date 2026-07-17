# M3.P1d-a — Identidad visible del roster inicial

Fecha: 2026-07-17

## Objetivo

Mostrar nombre y apellido de los jugadores propios desde el primer `roster.patch` emocional disponible, conservando `player_id` como identidad tecnica y el ownership contractual en `team-service`.

## Cambios realizados

- `PlayerEmotionalPatch` incorpora `full_name` y `position`.
- Los parches generados tras partidos y trades incluyen esos campos desde `PlayerAgentState`.
- La proyeccion emocional sigue identificando al jugador por el `player_id` compartido.
- El directorio de agentes usa `full_name` como etiqueta visible y muestra la posicion en el rol.
- Trade Center y Centro Medico reciben la misma identidad sin agregar adaptaciones propias.
- Se actualizo el contrato de `roster.patch`.

## Decision tecnica

No se creo un snapshot WebSocket ni un endpoint REST nuevo. `agent-service` ya lee `team_roster_players` para sembrar su proyeccion emocional; el delta retransmite esa identidad de lectura junto al estado. `team-service` sigue siendo el unico dueño de los datos contractuales.

## Pruebas

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
npm run build --prefix frontend
```

## Resultado

Cuando existe nombre contractual, ninguna superficie principal necesita presentar el hash/ID como nombre del jugador. Los jugadores iniciales y los recibidos por trade siguen el mismo criterio visual.

## Pendiente siguiente

`M3.P1d-b` — sincronizar el mood dinamico visible con el contexto del chat.
