# M3.P2g-h — Revision integral de persistencia

## Objetivo

Validar como conjunto la particion de `team-service/internal/persistence` antes de cerrar `M3.P2g`, sin introducir nuevos cambios salvo que la revision encontrara un problema real.

## Inventario final

| Archivo | Lineas | Responsabilidad |
|---|---:|---|
| `store.go` | 346 | Conexion, schema e inicializacion atomica de temporada |
| `injuries.go` | 318 | Decisiones medicas, recuperacion y lesiones |
| `trade_acceptance.go` | 228 | Aceptacion atomica y materializacion de trades |
| `schedule_dispatch.go` | 218 | Claim de calendario y payload para el simulador |
| `match_results.go` | 180 | Resultado, box scores y record de temporada |
| `trade_proposals.go` | 176 | Propuestas y rechazos de trades |
| `roster.go` | 107 | Patches emocionales y lecturas del roster |
| `salary_cap.go` | 37 | Persistencia del snapshot de salary cap |

El archivo original tenia 1545 lineas. La particion completa suma 1610 lineas fisicas: 65 lineas adicionales de overhead estructural por packages, imports y separacion de archivos.

## Dependencias internas revisadas

- `loadRoster` y `loadRosterPlayer` se comparten con trades y calendario.
- `saveSalaryCap` se comparte con fundacion y aceptacion de trades.
- `loadPlayerMatchStates` se comparte con despacho y lesiones.
- `loadSeasonRecord` se comparte con despacho y resultado idempotente.
- `createInjuriesForMatch` permanece en el dominio medico y participa en la transaccion de resultado.
- `queryer` es la unica interfaz de infraestructura compartida.

Estos cruces representan colaboracion entre dominios dentro del mismo paquete y no justifican crear subpackages ni nuevas interfaces.

## Findings

No se encontraron findings pendientes: no hay simbolos duplicados, helpers huerfanos, imports incorrectos, archivos con responsabilidades mezcladas ni crecimiento por encima del umbral definido en `AGENTS.md`.

No hubo cambios de codigo en este corte.

## Verificacion

- `gofmt -d services/team-service/internal/persistence`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service vet ./...`
- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service build ./...`
- `git diff --check`

## Pendiente siguiente

Continuar con el proximo objetivo de la pausa de analisis de M3. `M3.P2g` queda cerrado y `store.go` no necesita otra particion por ahora.
