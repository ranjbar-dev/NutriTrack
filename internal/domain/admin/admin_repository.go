package admin

import "context"

// AdminStats holds aggregated platform statistics.
// This is the domain DTO — no db/json struct tags.
type AdminStats struct {
	TotalNutritionists    int64
	ActiveNutritionists   int64
	InactiveNutritionists int64
	TotalClients          int64
	TotalFoods            int64
	ActiveDietPlans       int64
}

// AdminRepository defines the query port for admin statistics.
// Implemented in internal/infrastructure/persistence/admin.
type AdminRepository interface {
	GetStats(ctx context.Context) (AdminStats, error)
}
