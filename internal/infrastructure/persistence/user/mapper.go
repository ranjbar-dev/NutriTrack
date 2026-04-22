package user

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// toDomain converts a sqlc-generated User to a domain entity.
// This adapter lives at the repository boundary — domain never sees sqlc types.
func toDomain(u db.User) *entity.User {
	var email, passwordHash, gender, avatarURL string
	var birthDate *time.Time
	var height, weight *float64
	var nutritionistID *uuid.UUID

	if u.Email != nil {
		email = *u.Email
	}
	if u.PasswordHash != nil {
		passwordHash = *u.PasswordHash
	}
	if u.Gender != nil {
		gender = *u.Gender
	}
	if u.BirthDate.Valid {
		t := u.BirthDate.Time
		birthDate = &t
	}
	if u.Height.Valid {
		if f, err := u.Height.Float64Value(); err == nil && f.Valid {
			v := f.Float64
			height = &v
		}
	}
	if u.Weight.Valid {
		if f, err := u.Weight.Float64Value(); err == nil && f.Valid {
			v := f.Float64
			weight = &v
		}
	}
	if u.AvatarUrl != nil {
		avatarURL = *u.AvatarUrl
	}
	if u.NutritionistID.Valid {
		id := uuid.UUID(u.NutritionistID.Bytes)
		nutritionistID = &id
	}

	return entity.Reconstitute(
		u.ID,
		entity.Role(u.Role),
		u.Mobile,
		email,
		passwordHash,
		u.FirstName,
		u.LastName,
		gender,
		birthDate,
		height,
		weight,
		avatarURL,
		u.IsActive,
		nutritionistID,
		u.CreatedAt,
		u.UpdatedAt,
	)
}

func toDomainList(users []db.User) []*entity.User {
	result := make([]*entity.User, len(users))
	for i, u := range users {
		result[i] = toDomain(u)
	}
	return result
}

// --- helpers: domain → sqlc param types ---

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func float64ToNumeric(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{}
	}
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(*f, 'f', 2, 64))
	return n
}

func uuidToPgtypeUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: [16]byte(*id), Valid: true}
}
