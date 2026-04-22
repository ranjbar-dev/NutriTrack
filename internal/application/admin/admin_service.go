package admin

import (
	"context"

	domainadmin "github.com/ranjbar-dev/nutritrack/internal/domain/admin"
)

// AdminService provides platform-level statistics for the super-admin dashboard.
type AdminService struct {
	repo domainadmin.AdminRepository
}

// NewAdminService constructs an AdminService backed by the given AdminRepository.
func NewAdminService(repo domainadmin.AdminRepository) *AdminService {
	return &AdminService{repo: repo}
}

// GetStats returns aggregated platform counts (nutritionists, clients, foods, diet plans).
func (s *AdminService) GetStats(ctx context.Context) (domainadmin.AdminStats, error) {
	return s.repo.GetStats(ctx)
}
