package food

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
)

// MapFoodResponse converts a domain Food to a JSON-serialisable map.
// Defined here so that handlers in the interfaces layer do not need to import the entity package.
func MapFoodResponse(f *entity.Food) map[string]any {
	cats := make([]map[string]any, len(f.Categories()))
	for i, c := range f.Categories() {
		cats[i] = map[string]any{
			"id":   c.ID(),
			"name": c.Name(),
		}
	}

	resp := map[string]any{
		"id":           f.ID(),
		"name":         f.Name(),
		"unit":         f.Unit(),
		"calories":     f.Calories(),
		"protein":      f.Protein(),
		"carbohydrate": f.Carbohydrate(),
		"fat":          f.Fat(),
		"fiber":        f.Fiber(),
		"sugar":        f.Sugar(),
		"sodium":       f.Sodium(),
		"amount":       f.Amount(),
		"is_active":    f.IsActive(),
		"categories":   cats,
		"created_at":   f.CreatedAt(),
		"updated_at":   f.UpdatedAt(),
	}

	if f.CreatedBy() != nil {
		resp["created_by"] = *f.CreatedBy()
	} else {
		resp["created_by"] = nil
	}

	return resp
}
