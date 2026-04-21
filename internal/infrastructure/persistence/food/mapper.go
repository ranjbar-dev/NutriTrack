package food

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// foodToDomain converts a sqlc Food row to a domain entity (without categories).
func foodToDomain(f db.Food) *entity.Food {
	result := &entity.Food{
		ID:             f.ID,
		Name:           f.Name,
		NameNormalized: f.NameNormalized,
		Unit:           f.Unit,
		Calories:       numericToFloat64(f.Calories),
		Protein:        numericToFloat64(f.Protein),
		Carbohydrate:   numericToFloat64(f.Carbohydrate),
		Fat:            numericToFloat64(f.Fat),
		Fiber:          numericToFloat64(f.Fiber),
		Sugar:          numericToFloat64(f.Sugar),
		Sodium:         numericToFloat64(f.Sodium),
		Amount:         numericToFloat64(f.Amount),
		IsActive:       f.IsActive,
		Categories:     []entity.FoodCategory{},
		CreatedAt:      f.CreatedAt,
		UpdatedAt:      f.UpdatedAt,
	}

	if f.CreatedBy.Valid {
		id := uuid.UUID(f.CreatedBy.Bytes)
		result.CreatedBy = &id
	}

	return result
}

// categoryToDomain converts a sqlc FoodCategory row to a domain entity.
func categoryToDomain(c db.FoodCategory) entity.FoodCategory {
	return entity.FoodCategory{
		ID:        c.ID,
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
	}
}

// categoriesToDomain converts a slice of sqlc FoodCategory to domain entities.
func categoriesToDomain(rows []db.FoodCategory) []entity.FoodCategory {
	result := make([]entity.FoodCategory, len(rows))
	for i, c := range rows {
		result[i] = categoryToDomain(c)
	}
	return result
}

// --- helpers: domain → sqlc param types ---

// float64ToNumeric converts a float64 value to pgtype.Numeric.
func float64ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(f, 'f', 2, 64))
	return n
}

// numericToFloat64 converts pgtype.Numeric to float64, returning 0 on error.
func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return 0
	}
	return f.Float64
}

// uuidToPgtypeUUID converts *uuid.UUID to pgtype.UUID.
func uuidToPgtypeUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: [16]byte(*id), Valid: true}
}
