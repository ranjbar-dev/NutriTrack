package entity

import "errors"

var (
	ErrFoodNotFound         = errors.New("food not found")
	ErrFoodNameRequired     = errors.New("food name is required")
	ErrFoodUnitRequired     = errors.New("food unit is required")
	ErrFoodInvalidNutrition = errors.New("nutrition values must be non-negative")
	ErrCategoryNameRequired = errors.New("category name is required")
	ErrCategoryNotFound     = errors.New("food category not found")
)
