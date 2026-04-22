package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// PgUserRepository implements domain/user/repository.UserRepository using PostgreSQL + sqlc.
type PgUserRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPgUserRepository(pool *pgxpool.Pool) *PgUserRepository {
	return &PgUserRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *PgUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	u, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return toDomain(u), nil
}

func (r *PgUserRepository) FindByMobile(ctx context.Context, mobile string) (*entity.User, error) {
	u, err := r.queries.GetUserByMobile(ctx, mobile)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return toDomain(u), nil
}

func (r *PgUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	u, err := r.queries.GetUserByEmail(ctx, &email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, shared.ErrInternal
	}
	return toDomain(u), nil
}

func (r *PgUserRepository) FindClientsByNutritionist(ctx context.Context, nutritionistID uuid.UUID, limit, offset int32) ([]*entity.User, error) {
	users, err := r.queries.ListClientsByNutritionist(ctx, db.ListClientsByNutritionistParams{
		NutritionistID: pgtype.UUID{Bytes: [16]byte(nutritionistID), Valid: true},
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	return toDomainList(users), nil
}

func (r *PgUserRepository) CountClientsByNutritionist(ctx context.Context, nutritionistID uuid.UUID) (int64, error) {
	count, err := r.queries.CountClientsByNutritionist(ctx, pgtype.UUID{Bytes: [16]byte(nutritionistID), Valid: true})
	if err != nil {
		return 0, shared.ErrInternal
	}
	return count, nil
}

func (r *PgUserRepository) FindAllNutritionists(ctx context.Context, limit, offset int32) ([]*entity.User, error) {
	users, err := r.queries.ListNutritionists(ctx, db.ListNutritionistsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, shared.ErrInternal
	}
	return toDomainList(users), nil
}

func (r *PgUserRepository) CountAllNutritionists(ctx context.Context) (int64, error) {
	count, err := r.queries.CountNutritionists(ctx)
	if err != nil {
		return 0, shared.ErrInternal
	}
	return count, nil
}

func (r *PgUserRepository) Create(ctx context.Context, user *entity.User) error {
	var birthDate pgtype.Date
	if user.GetBirthDate() != nil {
		birthDate = pgtype.Date{Time: *user.GetBirthDate(), Valid: true}
	}

	created, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Role:           string(user.GetRole()),
		Mobile:         user.GetMobile(),
		Email:          strPtrOrNil(user.GetEmail()),
		PasswordHash:   strPtrOrNil(user.GetPasswordHash()),
		FirstName:      user.GetFirstName(),
		LastName:       user.GetLastName(),
		Gender:         strPtrOrNil(user.GetGender()),
		BirthDate:      birthDate,
		Height:         float64ToNumeric(user.GetHeight()),
		Weight:         float64ToNumeric(user.GetWeight()),
		AvatarUrl:      strPtrOrNil(user.GetAvatarURL()),
		IsActive:       user.GetIsActive(),
		NutritionistID: uuidToPgtypeUUID(user.GetNutritionistID()),
	})
	if err != nil {
		return shared.ErrInternal
	}

	// Populate entity with DB-generated values.
	user.SetID(created.ID)
	user.SetCreatedAt(created.CreatedAt)
	user.SetUpdatedAt(created.UpdatedAt)
	return nil
}

func (r *PgUserRepository) Update(ctx context.Context, user *entity.User) error {
	var birthDate pgtype.Date
	if user.GetBirthDate() != nil {
		birthDate = pgtype.Date{Time: *user.GetBirthDate(), Valid: true}
	}

	updated, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:           user.GetID(),
		FirstName:    user.GetFirstName(),
		LastName:     user.GetLastName(),
		Gender:       strPtrOrNil(user.GetGender()),
		BirthDate:    birthDate,
		Height:       float64ToNumeric(user.GetHeight()),
		Weight:       float64ToNumeric(user.GetWeight()),
		AvatarUrl:    strPtrOrNil(user.GetAvatarURL()),
		IsActive:     user.GetIsActive(),
		PasswordHash: strPtrOrNil(user.GetPasswordHash()),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shared.ErrUserNotFound
		}
		return shared.ErrInternal
	}

	user.SetUpdatedAt(updated.UpdatedAt)
	return nil
}

func (r *PgUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeleteUser(ctx, id); err != nil {
		return shared.ErrInternal
	}
	return nil
}

func (r *PgUserRepository) ExistsByMobile(ctx context.Context, mobile string) (bool, error) {
	exists, err := r.queries.ExistsByMobile(ctx, mobile)
	if err != nil {
		return false, shared.ErrInternal
	}
	return exists, nil
}

func (r *PgUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.queries.ExistsByEmail(ctx, &email)
	if err != nil {
		return false, shared.ErrInternal
	}
	return exists, nil
}
