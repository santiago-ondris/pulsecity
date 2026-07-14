# Finance Events

Contratos financieros acotados para Milestone 3.

## Subjects

- `salary_cap.calculado`
- `finance.patch`

## `salary_cap.calculado`

Publicado por `team-service` cuando recalcula la situacion de salary cap de la franquicia.

```json
{
  "event_id": "salary-cap-game-1-initial-season-game-1",
  "game_id": "game-1",
  "occurred_at": "2026-10-22T00:00:00Z",
  "schema_version": 1,
  "simulated_date": "2026-10-22",
  "cap_base": 141000000,
  "luxury_tax_line": 171000000,
  "committed_salary": 78500000,
  "cap_space": 62500000,
  "luxury_tax_space": 92500000,
  "roster_count": 15,
  "status": "under_cap",
  "near_luxury_tax": false,
  "projected_tax_payment": 0
}
```

Estados:

- `under_cap`
- `over_cap`
- `luxury_tax`

## `finance.patch`

Delta reenviado por `gateway` al frontend.

```json
{
  "type": "finance.patch",
  "subject": "finance.patch",
  "game_id": "game-1",
  "patch": {
    "simulated_date": "2026-10-22",
    "source_event_id": "salary-cap-game-1-initial-season-game-1",
    "source_subject": "salary_cap.calculado",
    "cap_base": 141000000,
    "luxury_tax_line": 171000000,
    "committed_salary": 78500000,
    "cap_space": 62500000,
    "luxury_tax_space": 92500000,
    "roster_count": 15,
    "status": "under_cap",
    "near_luxury_tax": false,
    "projected_tax_payment": 0
  }
}
```
