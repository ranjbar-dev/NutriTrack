package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/medication/entity"
)

// MedicationRepository defines persistence operations for the Medication aggregate.
type MedicationRepository interface {
	Create(ctx context.Context, med *entity.Medication) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Medication, error)
	Update(ctx context.Context, med *entity.Medication) error
	Delete(ctx context.Context, id uuid.UUID) error
	Deactivate(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, limit, offset int32) ([]*entity.Medication, error)
	CountSearch(ctx context.Context, query string) (int64, error)
}
