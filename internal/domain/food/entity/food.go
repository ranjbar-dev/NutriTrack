package entity

import (
	"time"

	"github.com/google/uuid"
)

// Food is the food domain aggregate.
// All fields are unexported; use getters for read access and domain methods for mutations.
type Food struct {
	id             uuid.UUID
	name           string
	nameNormalized string
	unit           string
	calories       float64
	protein        float64
	carbohydrate   float64
	fat            float64
	fiber          float64
	sugar          float64
	sodium         float64
	amount         float64
	createdBy      *uuid.UUID
	isActive       bool
	categories     []FoodCategory
	createdAt      time.Time
	updatedAt      time.Time
}

// NewFood creates a new Food aggregate with input validation.
func NewFood(
	name, nameNormalized, unit string,
	calories, protein, carbohydrate, fat, fiber, sugar, sodium, amount float64,
	createdBy *uuid.UUID,
	categories []FoodCategory,
) (*Food, error) {
	if name == "" {
		return nil, ErrFoodNameRequired
	}
	if unit == "" {
		return nil, ErrFoodUnitRequired
	}
	return &Food{
		id:             uuid.New(),
		name:           name,
		nameNormalized: nameNormalized,
		unit:           unit,
		calories:       calories,
		protein:        protein,
		carbohydrate:   carbohydrate,
		fat:            fat,
		fiber:          fiber,
		sugar:          sugar,
		sodium:         sodium,
		amount:         amount,
		createdBy:      createdBy,
		isActive:       true,
		categories:     categories,
	}, nil
}

// ReconstructFood rebuilds a Food aggregate from persistent storage without validation.
// For infrastructure use only.
func ReconstructFood(
	id uuid.UUID,
	name, nameNormalized, unit string,
	calories, protein, carbohydrate, fat, fiber, sugar, sodium, amount float64,
	createdBy *uuid.UUID,
	isActive bool,
	categories []FoodCategory,
	createdAt, updatedAt time.Time,
) *Food {
	return &Food{
		id:             id,
		name:           name,
		nameNormalized: nameNormalized,
		unit:           unit,
		calories:       calories,
		protein:        protein,
		carbohydrate:   carbohydrate,
		fat:            fat,
		fiber:          fiber,
		sugar:          sugar,
		sodium:         sodium,
		amount:         amount,
		createdBy:      createdBy,
		isActive:       isActive,
		categories:     categories,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// --- Getters ---

func (f *Food) ID() uuid.UUID              { return f.id }
func (f *Food) Name() string               { return f.name }
func (f *Food) NameNormalized() string     { return f.nameNormalized }
func (f *Food) Unit() string               { return f.unit }
func (f *Food) Calories() float64          { return f.calories }
func (f *Food) Protein() float64           { return f.protein }
func (f *Food) Carbohydrate() float64      { return f.carbohydrate }
func (f *Food) Fat() float64               { return f.fat }
func (f *Food) Fiber() float64             { return f.fiber }
func (f *Food) Sugar() float64             { return f.sugar }
func (f *Food) Sodium() float64            { return f.sodium }
func (f *Food) Amount() float64            { return f.amount }
func (f *Food) CreatedBy() *uuid.UUID      { return f.createdBy }
func (f *Food) IsActive() bool             { return f.isActive }
func (f *Food) Categories() []FoodCategory { return f.categories }
func (f *Food) CreatedAt() time.Time       { return f.createdAt }
func (f *Food) UpdatedAt() time.Time       { return f.updatedAt }

// --- Domain methods ---

// Update applies validated field changes to the food aggregate.
func (f *Food) Update(
	name, nameNormalized, unit string,
	calories, protein, carbohydrate, fat, fiber, sugar, sodium, amount float64,
	categories []FoodCategory,
) error {
	if name == "" {
		return ErrFoodNameRequired
	}
	if unit == "" {
		return ErrFoodUnitRequired
	}
	f.name = name
	f.nameNormalized = nameNormalized
	f.unit = unit
	f.calories = calories
	f.protein = protein
	f.carbohydrate = carbohydrate
	f.fat = fat
	f.fiber = fiber
	f.sugar = sugar
	f.sodium = sodium
	f.amount = amount
	f.categories = categories
	return nil
}

// --- Infrastructure-only setters ---

// SetPersistedState populates DB-generated fields after an INSERT (infrastructure use only).
func (f *Food) SetPersistedState(id uuid.UUID, isActive bool, createdAt, updatedAt time.Time) {
	f.id = id
	f.isActive = isActive
	f.createdAt = createdAt
	f.updatedAt = updatedAt
}

// SetUpdatedAt records the DB-assigned updated_at after an UPDATE (infrastructure use only).
func (f *Food) SetUpdatedAt(updatedAt time.Time) {
	f.updatedAt = updatedAt
}

// SetCategories sets the loaded categories on a reconstructed aggregate (infrastructure use only).
func (f *Food) SetCategories(cats []FoodCategory) {
	f.categories = cats
}
