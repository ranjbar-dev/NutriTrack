package entity

import (
	"time"

	"github.com/google/uuid"
)

// Food is the food domain aggregate.
type Food struct {
	ID             uuid.UUID
	Name           string
	NameNormalized string
	Unit           string
	Calories       float64
	Protein        float64
	Carbohydrate   float64
	Fat            float64
	Fiber          float64
	CreatedBy      *uuid.UUID
	IsActive       bool
	Categories     []FoodCategory
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// FoodCategory represents a category tag on a food item.
type FoodCategory struct {
	ID   uuid.UUID
	Name string
}
