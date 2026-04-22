package food

import (
	"time"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
)

// foodCacheDTO is the Redis serialisation DTO for the Food aggregate.
// Exported fields carry json tags so encoding/json can round-trip the data.
type foodCacheDTO struct {
	ID             uuid.UUID     `json:"id"`
	Name           string        `json:"name"`
	NameNormalized string        `json:"name_normalized"`
	Unit           string        `json:"unit"`
	Calories       float64       `json:"calories"`
	Protein        float64       `json:"protein"`
	Carbohydrate   float64       `json:"carbohydrate"`
	Fat            float64       `json:"fat"`
	Fiber          float64       `json:"fiber"`
	Sugar          float64       `json:"sugar"`
	Sodium         float64       `json:"sodium"`
	Amount         float64       `json:"amount"`
	CreatedBy      *uuid.UUID    `json:"created_by"`
	IsActive       bool          `json:"is_active"`
	Categories     []catCacheDTO `json:"categories"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type catCacheDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func foodToCache(f *entity.Food) foodCacheDTO {
	cats := make([]catCacheDTO, len(f.Categories()))
	for i, c := range f.Categories() {
		cats[i] = catCacheDTO{ID: c.ID(), Name: c.Name(), CreatedAt: c.CreatedAt()}
	}
	return foodCacheDTO{
		ID:             f.ID(),
		Name:           f.Name(),
		NameNormalized: f.NameNormalized(),
		Unit:           f.Unit(),
		Calories:       f.Calories(),
		Protein:        f.Protein(),
		Carbohydrate:   f.Carbohydrate(),
		Fat:            f.Fat(),
		Fiber:          f.Fiber(),
		Sugar:          f.Sugar(),
		Sodium:         f.Sodium(),
		Amount:         f.Amount(),
		CreatedBy:      f.CreatedBy(),
		IsActive:       f.IsActive(),
		Categories:     cats,
		CreatedAt:      f.CreatedAt(),
		UpdatedAt:      f.UpdatedAt(),
	}
}

func cacheToFood(dto foodCacheDTO) *entity.Food {
	cats := make([]entity.FoodCategory, len(dto.Categories))
	for i, c := range dto.Categories {
		cats[i] = entity.ReconstructFoodCategory(c.ID, c.Name, c.CreatedAt)
	}
	return entity.ReconstructFood(
		dto.ID, dto.Name, dto.NameNormalized, dto.Unit,
		dto.Calories, dto.Protein, dto.Carbohydrate, dto.Fat,
		dto.Fiber, dto.Sugar, dto.Sodium, dto.Amount,
		dto.CreatedBy, dto.IsActive, cats,
		dto.CreatedAt, dto.UpdatedAt,
	)
}
