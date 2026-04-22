package foodrequest

import (
	frEntity "github.com/ranjbar-dev/nutritrack/internal/domain/foodrequest/entity"
)

// MapFoodRequestResponse converts a FoodRequest entity to a JSON-serialisable map.
// Defined here so that handlers in the interfaces layer do not need to import the entity package.
func MapFoodRequestResponse(r *frEntity.FoodRequest) map[string]any {
	return map[string]any{
		"id":               r.GetID(),
		"client_id":        r.GetClientID(),
		"nutritionist_id":  r.GetNutritionistID(),
		"food_name":        r.GetFoodName(),
		"status":           string(r.GetStatus()),
		"rejection_reason": r.GetRejectionReason(),
		"created_food_id":  r.GetCreatedFoodID(),
		"created_at":       r.GetCreatedAt(),
		"updated_at":       r.GetUpdatedAt(),
	}
}
