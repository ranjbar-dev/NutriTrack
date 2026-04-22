package foodrequest

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// toDomain converts a db.FoodRequest row to a domain *entity.FoodRequest.
func toDomain(row db.FoodRequest) *entity.FoodRequest {
	return entity.FromPersistence(
		row.ID,
		row.ClientID,
		row.NutritionistID,
		row.FoodName,
		entity.FoodRequestStatus(row.Status),
		row.RejectionReason,
		row.CreatedFoodID,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
