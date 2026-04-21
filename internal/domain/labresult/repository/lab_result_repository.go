package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/labresult/entity"
)

// LabResultRepository defines persistence operations for lab results.
type LabResultRepository interface {
	Create(ctx context.Context, result *entity.LabResult) (*entity.LabResult, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.LabResult, error)
	ListByClientID(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.LabResult, int64, error)
}
