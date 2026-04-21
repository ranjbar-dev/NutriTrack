package admin

import (
	"context"

	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// AdminService provides platform-level statistics for the super-admin dashboard.
type AdminService struct {
	queries *db.Queries
}

// NewAdminService constructs an AdminService backed by the given sqlc Queries.
func NewAdminService(q *db.Queries) *AdminService {
	return &AdminService{queries: q}
}

// GetStats returns aggregated platform counts (nutritionists, clients, foods, diet plans).
func (s *AdminService) GetStats(ctx context.Context) (db.AdminStats, error) {
	return s.queries.GetAdminStats(ctx)
}
