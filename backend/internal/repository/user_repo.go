package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// SearchClientsParams holds filter/sort parameters for client search.
type SearchClientsParams struct {
	NutritionistID uuid.UUID
	Query          string // ILIKE search against full_name and mobile
	Active         *bool  // nil = no filter
	SortBy         string // "name" or "created_at" (default)
	Limit          int32
	Offset         int32
}

// UserRepository defines operations on the users table.
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*sqlc.User, error)
	GetByMobile(ctx context.Context, mobile string) (*sqlc.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*sqlc.User, error)
	Create(ctx context.Context, params sqlc.CreateUserParams) (*sqlc.User, error)
	GetClientsByNutritionist(ctx context.Context, nutritionistID uuid.UUID) ([]sqlc.User, error)
	UpdateActive(ctx context.Context, id uuid.UUID, active bool) error
	SearchClients(ctx context.Context, params SearchClientsParams) ([]sqlc.User, int, error)
	GetClientByIDForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID) (*sqlc.User, error)
	UpdateClientProfile(ctx context.Context, id uuid.UUID, dateOfBirth pgtype.Date, heightCm pgtype.Float4) (*sqlc.User, error)
}

// userRepository implements UserRepository using sqlc-generated queries.
type userRepository struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

// NewUserRepository creates a new UserRepository backed by the given pool.
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool, q: sqlc.New(pool)}
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

func (r *userRepository) SearchClients(ctx context.Context, params SearchClientsParams) ([]sqlc.User, int, error) {
	// Build ORDER BY clause safely (not user-controlled string injection — limited to two known values)
	orderBy := "created_at DESC"
	if params.SortBy == "name" {
		orderBy = "full_name ASC"
	}

	// Build WHERE clause for active filter
	activeClause := ""
	args := []any{
		pgtype.UUID{Bytes: params.NutritionistID, Valid: true},
		params.Query,
	}
	argIdx := 3
	if params.Active != nil {
		activeClause = fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *params.Active)
		argIdx++
	}

	countSQL := fmt.Sprintf(`
		SELECT COUNT(*) FROM users
		WHERE nutritionist_id = $1
		  AND role = 'client'
		  AND ($2 = '' OR full_name ILIKE '%%' || $2 || '%%' OR mobile ILIKE '%%' || $2 || '%%')
		%s`, activeClause)

	var total int
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	selectSQL := fmt.Sprintf(`
		SELECT id, role, full_name, email, password_hash, mobile, date_of_birth, height_cm,
		       gender, nutritionist_id, is_active, notes, created_at, updated_at
		FROM users
		WHERE nutritionist_id = $1
		  AND role = 'client'
		  AND ($2 = '' OR full_name ILIKE '%%' || $2 || '%%' OR mobile ILIKE '%%' || $2 || '%%')
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`, activeClause, orderBy, argIdx, argIdx+1)

	args = append(args, params.Limit, params.Offset)

	rows, err := r.pool.Query(ctx, selectSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []sqlc.User
	for rows.Next() {
		var u sqlc.User
		if err := rows.Scan(
			&u.ID, &u.Role, &u.FullName, &u.Email, &u.PasswordHash,
			&u.Mobile, &u.DateOfBirth, &u.HeightCm, &u.Gender, &u.NutritionistID,
			&u.IsActive, &u.Notes, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) GetClientByIDForNutritionist(ctx context.Context, clientID, nutritionistID uuid.UUID) (*sqlc.User, error) {
	user, err := r.q.GetClientByIDForNutritionist(ctx, sqlc.GetClientByIDForNutritionistParams{
		ID:             pgtype.UUID{Bytes: clientID, Valid: true},
		NutritionistID: pgtype.UUID{Bytes: nutritionistID, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateClientProfile(ctx context.Context, id uuid.UUID, dateOfBirth pgtype.Date, heightCm pgtype.Float4) (*sqlc.User, error) {
	user, err := r.q.UpdateClientProfile(ctx, sqlc.UpdateClientProfileParams{
		ID:          pgtype.UUID{Bytes: id, Valid: true},
		DateOfBirth: dateOfBirth,
		HeightCm:    heightCm,
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}
