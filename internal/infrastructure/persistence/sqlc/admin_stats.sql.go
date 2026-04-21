package db

import (
	"context"
)

// AdminStats holds aggregated platform-level counts for the super-admin dashboard.
type AdminStats struct {
	TotalNutritionists    int64 `db:"total_nutritionists"`
	ActiveNutritionists   int64 `db:"active_nutritionists"`
	InactiveNutritionists int64 `db:"inactive_nutritionists"`
	TotalClients          int64 `db:"total_clients"`
	TotalFoods            int64 `db:"total_foods"`
	ActiveDietPlans       int64 `db:"active_diet_plans"`
}

const getAdminStats = `
SELECT
    COUNT(*) FILTER (WHERE role = 'nutritionist') AS total_nutritionists,
    COUNT(*) FILTER (WHERE role = 'nutritionist' AND is_active = true) AS active_nutritionists,
    COUNT(*) FILTER (WHERE role = 'nutritionist' AND is_active = false) AS inactive_nutritionists,
    COUNT(*) FILTER (WHERE role = 'client') AS total_clients,
    (SELECT COUNT(*) FROM foods) AS total_foods,
    (SELECT COUNT(*) FROM diet_plans WHERE is_active = true) AS active_diet_plans
FROM users
`

// GetAdminStats fetches aggregated platform statistics in a single query.
func (q *Queries) GetAdminStats(ctx context.Context) (AdminStats, error) {
	row := q.db.QueryRow(ctx, getAdminStats)
	var s AdminStats
	err := row.Scan(
		&s.TotalNutritionists,
		&s.ActiveNutritionists,
		&s.InactiveNutritionists,
		&s.TotalClients,
		&s.TotalFoods,
		&s.ActiveDietPlans,
	)
	return s, err
}
