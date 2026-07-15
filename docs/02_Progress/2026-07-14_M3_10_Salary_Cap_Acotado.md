# M3.10 — Salary Cap Acotado

Fecha: 2026-07-14

## Objetivo

Hacer que el dinero empiece a importar dentro de M3: cap base, luxury tax line, salarios por jugador, reaccion financiera de agentes y visibilidad minima en frontend.

## Cambios realizados

- `team-service` ahora calcula un `SalaryCapSnapshot` con:
  - `cap_base`
  - `luxury_tax_line`
  - `committed_salary`
  - `cap_space`
  - `luxury_tax_space`
  - `status`: `under_cap`, `over_cap`, `luxury_tax`
  - `projected_tax_payment`
- Se agrego persistencia `team_salary_cap` como fuente de verdad financiera del equipo.
- Los salarios iniciales del roster dejaron de ser montos cosmeticos bajos y pasan a derivarse de rating, tier de rol y posicion.
- `team-service` publica:
  - `salary_cap.calculado` para servicios internos
  - `finance.patch` para delta WebSocket
- `gateway` escucha `finance.patch` y lo reenvia al browser.
- El frontend mantiene `financeState` y muestra cap comprometido, cap space, luxury tax space y estado en la pestaña `Temporada`.
- `agent-service` escucha `salary_cap.calculado`:
  - el CFO sube `budget_alert` y baja `financial_trust` si hay luxury tax
  - el Owner pierde `patience_remaining` y `business_trust` si hay luxury tax
- Se documentaron contratos en:
  - `shared/events/finance.md`
  - `shared/events/websocket_deltas.md`

## Decision de alcance

M3.10 deja lista la fuente de verdad financiera y la reaccion sistemica. El flujo real de trades todavia no existe; por eso la validacion concreta al aceptar/rechazar trades se consume en `M3.12`, cuando se implemente `trade.aceptada` y el cierre operativo de una negociacion.

## Verificacion

- `GOCACHE=/tmp/pulsecity-team-gocache go -C services/team-service test ./...`
- `GOCACHE=/tmp/pulsecity-gateway-gocache go -C services/gateway test ./...`
- `cargo test --manifest-path services/agent-service/Cargo.toml`
- `npm run build --prefix frontend`

