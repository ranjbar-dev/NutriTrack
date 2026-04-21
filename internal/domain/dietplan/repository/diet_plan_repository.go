package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/entity"
)

type DietPlanRepository interface {
	// CreateWithArchive atomically archives any existing active plan for the client,
	// then inserts the new plan in a single DB transaction.
	CreateWithArchive(ctx context.Context, plan *entity.DietPlan) error

	FindByID(ctx context.Context, id uuid.UUID) (*entity.DietPlan, error)
	FindActiveByClientID(ctx context.Context, clientID uuid.UUID) (*entity.DietPlan, error)
	ListByClientID(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.DietPlan, error)
	CountByClientID(ctx context.Context, clientID uuid.UUID) (int64, error)
	Update(ctx context.Context, plan *entity.DietPlan) error
	Delete(ctx context.Context, id uuid.UUID) error

	AddDay(ctx context.Context, day *entity.DietPlanDay) error
	FindDayByID(ctx context.Context, id uuid.UUID) (*entity.DietPlanDay, error)
	ListDays(ctx context.Context, planID uuid.UUID) ([]*entity.DietPlanDay, error)
	DeleteDay(ctx context.Context, id uuid.UUID) error

	AddMeal(ctx context.Context, meal *entity.DietMeal) error
	FindMealByID(ctx context.Context, id uuid.UUID) (*entity.DietMeal, error)
	ListMeals(ctx context.Context, dayID uuid.UUID) ([]*entity.DietMeal, error)
	DeleteMeal(ctx context.Context, id uuid.UUID) error

	AddOption(ctx context.Context, option *entity.MealOption) error
	FindOptionByID(ctx context.Context, id uuid.UUID) (*entity.MealOption, error)
	ListOptions(ctx context.Context, mealID uuid.UUID) ([]*entity.MealOption, error)
	DeleteOption(ctx context.Context, id uuid.UUID) error

	AddItem(ctx context.Context, item *entity.MealOptionItem) error
	FindItemByID(ctx context.Context, id uuid.UUID) (*entity.MealOptionItem, error)
	ListItems(ctx context.Context, optionID uuid.UUID) ([]*entity.MealOptionItem, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	DeleteItemsByOption(ctx context.Context, optionID uuid.UUID) error
}
