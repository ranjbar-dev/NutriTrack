package user

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranjbar-dev/nutritrack/internal/domain/user/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// toDomain converts a sqlc-generated User to a domain entity.
// This adapter lives at the repository boundary — domain never sees sqlc types.
func toDomain(u db.User) *entity.User {
	result := &entity.User{
		ID:        u.ID,
		Role:      u.Role,
		Mobile:    u.Mobile,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	if u.Email != nil {
		result.Email = *u.Email
	}
	if u.PasswordHash != nil {
		result.PasswordHash = *u.PasswordHash
	}
	if u.Gender != nil {
		result.Gender = *u.Gender
	}
	if u.BirthDate.Valid {
		t := u.BirthDate.Time
		result.BirthDate = &t
	}
	if u.Height.Valid {
		if f, err := u.Height.Float64Value(); err == nil && f.Valid {
			v := f.Float64
			result.Height = &v
		}
	}
	if u.Weight.Valid {
		if f, err := u.Weight.Float64Value(); err == nil && f.Valid {
			v := f.Float64
			result.Weight = &v
		}
	}
	if u.AvatarUrl != nil {
		result.AvatarURL = *u.AvatarUrl
	}
	if u.NutritionistID.Valid {
		id := uuid.UUID(u.NutritionistID.Bytes)
		result.NutritionistID = &id
	}

	return result
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
