package entity

import (
	"time"

	"github.com/google/uuid"
)

// FoodCategory represents a category tag on a food item.
type FoodCategory struct {
	id        uuid.UUID
	name      string
	createdAt time.Time
}

// NewFoodCategory creates a new FoodCategory with validation.
func NewFoodCategory(name string) (FoodCategory, error) {
	if name == "" {
		return FoodCategory{}, ErrCategoryNameRequired
	}
	return FoodCategory{
		id:        uuid.New(),
		name:      name,
		createdAt: time.Now(),
	}, nil
}

// ReconstructFoodCategory rebuilds a FoodCategory from persistent storage (no validation).
func ReconstructFoodCategory(id uuid.UUID, name string, createdAt time.Time) FoodCategory {
	return FoodCategory{id: id, name: name, createdAt: createdAt}
}

// ID returns the category's unique identifier.
func (c FoodCategory) ID() uuid.UUID { return c.id }

// Name returns the category name.
func (c FoodCategory) Name() string { return c.name }

// CreatedAt returns the creation timestamp.
func (c FoodCategory) CreatedAt() time.Time { return c.createdAt }
