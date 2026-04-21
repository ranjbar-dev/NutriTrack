package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
)

// UserRepository defines persistence operations for the User aggregate.
// Implementations live in internal/infrastructure/persistence/user/.
type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByMobile(ctx context.Context, mobile string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindClientsByNutritionist(ctx context.Context, nutritionistID uuid.UUID, limit, offset int32) ([]*entity.User, error)
	CountClientsByNutritionist(ctx context.Context, nutritionistID uuid.UUID) (int64, error)
	FindAllNutritionists(ctx context.Context, limit, offset int32) ([]*entity.User, error)
	CountAllNutritionists(ctx context.Context) (int64, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByMobile(ctx context.Context, mobile string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
