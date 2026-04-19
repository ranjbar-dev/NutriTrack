package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// UserRepository defines operations on the users table.
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*sqlc.User, error)
	GetByMobile(ctx context.Context, mobile string) (*sqlc.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*sqlc.User, error)
	Create(ctx context.Context, params sqlc.CreateUserParams) (*sqlc.User, error)
	GetClientsByNutritionist(ctx context.Context, nutritionistID uuid.UUID) ([]sqlc.User, error)
	UpdateActive(ctx context.Context, id uuid.UUID, active bool) error
}

// userRepository implements UserRepository using sqlc-generated queries.
type userRepository struct {
	q *sqlc.Queries
}

// NewUserRepository creates a new UserRepository backed by the given sqlc.DBTX.
func NewUserRepository(db sqlc.DBTX) UserRepository {
	return &userRepository{q: sqlc.New(db)}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*sqlc.User, error) {
	user, err := r.q.GetUserByEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByMobile(ctx context.Context, mobile string) (*sqlc.User, error) {
	user, err := r.q.GetUserByMobile(ctx, pgtype.Text{String: mobile, Valid: true})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*sqlc.User, error) {
	user, err := r.q.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, params sqlc.CreateUserParams) (*sqlc.User, error) {
	user, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetClientsByNutritionist(ctx context.Context, nutritionistID uuid.UUID) ([]sqlc.User, error) {
	return r.q.GetClientsByNutritionistID(ctx, pgtype.UUID{Bytes: nutritionistID, Valid: true})
}

func (r *userRepository) UpdateActive(ctx context.Context, id uuid.UUID, active bool) error {
	return r.q.UpdateUserActive(ctx, sqlc.UpdateUserActiveParams{
		ID:       pgtype.UUID{Bytes: id, Valid: true},
		IsActive: active,
	})
}
