# M3.P1d-b — Estado dinamico coherente en chat

Fecha: 2026-07-17

## Objetivo

Garantizar que el contexto del chat use el mismo estado emocional dinamico que el frontend muestra para los cinco agentes core.

## Cambios realizados

- Al persistir un cambio en `agent_core_states`, `agent-service` actualiza `agent_individual_states.emotional_state` dentro de la misma transaccion.
- Al iniciar, `EnsureSchema` reconcilia de forma idempotente estados legacy ya existentes con la capa individual.
- `narrative-service` mantiene su lectura actual desde `agent_individual_states`; no incorpora conocimiento sobre tablas internas adicionales.
- El test del prompt verifica explicitamente que el estado emocional recibido forma parte del contexto enviado al responder.

## Decision tecnica

El arreglo se realiza en `agent-service`, dueño del estado emocional. Resolverlo mediante un join dentro de `narrative-service` hubiera filtrado el desacople interno hacia otro servicio y dejado dos fuentes visibles para el lector.

Este corte no profundiza la conducta de los agentes ni agrega datos operativos al prompt. Solo garantiza coherencia entre el mood calculado por reglas, el estado visible y el contexto lingüistico.

## Compatibilidad con partidas existentes

La reconciliacion de startup actualiza las partidas ya creadas. Basta reiniciar `agent-service`; no hace falta fundar una partida nueva ni ejecutar una migracion manual.

## Pruebas

```bash
cargo test --manifest-path services/agent-service/Cargo.toml
GOCACHE=/tmp/pulsecity-narrative-gocache go -C services/narrative-service test ./...
```

Tambien se ejecuto la reconciliacion contra el Postgres local dentro de una transaccion descartable. El `UPDATE` encontro cinco agentes core desincronizados y la consulta posterior devolvio cero diferencias. La transaccion termino con `ROLLBACK`, por lo que la validacion no modifico la partida.

## Resultado

Un agente core mostrado como `concerned`, `pressured`, `frustrated` u otro mood dinamico entrega ese mismo valor al prompt del chat después de persistir el cambio o reconciliar el save al iniciar.

## Pendiente siguiente

Con las correcciones del playtest cerradas, el siguiente mini milestone es `M3.P2 — Particion de monolitos`.
