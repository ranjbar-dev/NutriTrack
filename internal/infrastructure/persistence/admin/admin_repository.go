package admin

import (
	"context"

	domainadmin "github.com/ranjbar-dev/nutritrack/internal/domain/admin"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgAdminRepository implements domain/admin.AdminRepository using SQLC.
type PgAdminRepository struct {
	queries *db.Queries
}

// NewPgAdminRepository creates a new PgAdminRepository.
func NewPgAdminRepository(q *db.Queries) *PgAdminRepository {
	return &PgAdminRepository{queries: q}
}

// GetStats fetches aggregated platform statistics and maps to the domain DTO.
func (r *PgAdminRepository) GetStats(ctx context.Context) (domainadmin.AdminStats, error) {
	row, err := r.queries.GetAdminStats(ctx)
	if err != nil {
		return domainadmin.AdminStats{}, err
	}
	return domainadmin.AdminStats{
		TotalNutritionists:    row.TotalNutritionists,
		ActiveNutritionists:   row.ActiveNutritionists,
		InactiveNutritionists: row.InactiveNutritionists,
		TotalClients:          row.TotalClients,
		TotalFoods:            row.TotalFoods,
		ActiveDietPlans:       row.ActiveDietPlans,
	}, nil
}
