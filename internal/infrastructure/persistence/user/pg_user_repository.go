package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
)

// PgUserRepository implements domain/user/repository.UserRepository using PostgreSQL + sqlc.
// Full implementation added in Phase 3.
type PgUserRepository struct {
	db *pgxpool.Pool
}

func NewPgUserRepository(db *pgxpool.Pool) *PgUserRepository {
	return &PgUserRepository{db: db}
}

func (r *PgUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return nil, nil // TODO Phase 3
}

func (r *PgUserRepository) FindByMobile(ctx context.Context, mobile string) (*entity.User, error) {
	return nil, nil // TODO Phase 3
}

func (r *PgUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil // TODO Phase 3
}

func (r *PgUserRepository) FindClientsByNutritionist(ctx context.Context, nutritionistID uuid.UUID, limit, offset int32) ([]*entity.User, error) {
	return nil, nil // TODO Phase 3
}

func (r *PgUserRepository) CountClientsByNutritionist(ctx context.Context, nutritionistID uuid.UUID) (int64, error) {
	return 0, nil // TODO Phase 3
}

func (r *PgUserRepository) FindAllNutritionists(ctx context.Context, limit, offset int32) ([]*entity.User, error) {
	return nil, nil // TODO Phase 3
}

func (r *PgUserRepository) CountAllNutritionists(ctx context.Context) (int64, error) {
	return 0, nil // TODO Phase 3
}

func (r *PgUserRepository) Create(ctx context.Context, user *entity.User) error {
	return nil // TODO Phase 3
}

func (r *PgUserRepository) Update(ctx context.Context, user *entity.User) error {
	return nil // TODO Phase 3
}

func (r *PgUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil // TODO Phase 3
}

func (r *PgUserRepository) ExistsByMobile(ctx context.Context, mobile string) (bool, error) {
	return false, nil // TODO Phase 3
}

func (r *PgUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil // TODO Phase 3
}
