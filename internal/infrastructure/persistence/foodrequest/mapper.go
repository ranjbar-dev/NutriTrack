package foodrequest

import (
	"github.com/jackc/pgx/v5"
	"github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// toDomain converts a db.FoodRequest row to a domain *entity.FoodRequest.
func toDomain(row db.FoodRequest) *entity.FoodRequest {
	return &entity.FoodRequest{
		ID:              row.ID,
		ClientID:        row.ClientID,
		NutritionistID:  row.NutritionistID,
		FoodName:        row.FoodName,
		Status:          entity.FoodRequestStatus(row.Status),
		RejectionReason: row.RejectionReason,
		CreatedFoodID:   row.CreatedFoodID,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func isNotFound(err error) bool {
	return err == pgx.ErrNoRows
}
