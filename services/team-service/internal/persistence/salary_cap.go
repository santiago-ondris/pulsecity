package persistence

import (
	"context"
	"fmt"

	"github.com/pulsecity/services/team-service/internal/domain"
)

func saveSalaryCap(ctx context.Context, q queryer, snapshot domain.SalaryCapSnapshot) error {
	if _, err := q.Exec(ctx, `
INSERT INTO team_salary_cap (
	game_id, simulated_date, cap_base, luxury_tax_line, committed_salary, cap_space,
	luxury_tax_space, roster_count, status, near_luxury_tax, projected_tax_payment,
	source_event_id, source_subject, updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW())
ON CONFLICT (game_id) DO UPDATE SET
	simulated_date = EXCLUDED.simulated_date,
	cap_base = EXCLUDED.cap_base,
	luxury_tax_line = EXCLUDED.luxury_tax_line,
	committed_salary = EXCLUDED.committed_salary,
	cap_space = EXCLUDED.cap_space,
	luxury_tax_space = EXCLUDED.luxury_tax_space,
	roster_count = EXCLUDED.roster_count,
	status = EXCLUDED.status,
	near_luxury_tax = EXCLUDED.near_luxury_tax,
	projected_tax_payment = EXCLUDED.projected_tax_payment,
	source_event_id = EXCLUDED.source_event_id,
	source_subject = EXCLUDED.source_subject,
	updated_at = NOW();
`, snapshot.GameID, snapshot.SimulatedDate, snapshot.CapBase, snapshot.LuxuryTaxLine, snapshot.CommittedSalary, snapshot.CapSpace, snapshot.LuxuryTaxSpace, int16(snapshot.RosterCount), snapshot.Status, snapshot.NearLuxuryTax, snapshot.ProjectedTaxPayment, snapshot.SourceEventID, snapshot.SourceSubject); err != nil {
		return fmt.Errorf("save salary cap %s: %w", snapshot.GameID, err)
	}

	return nil
}
