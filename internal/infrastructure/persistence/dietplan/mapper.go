package dietplan

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

// dietPlanToDomain converts a sqlc DietPlan row to a domain entity.
func dietPlanToDomain(p db.DietPlan) *entity.DietPlan {
	return &entity.DietPlan{
		ID:                 p.ID,
		ClientID:           p.ClientID,
		NutritionistID:     p.NutritionistID,
		Title:              p.Title,
		StartDate:          p.StartDate,
		EndDate:            p.EndDate,
		Notes:              p.Notes,
		DailyWaterTargetML: int(p.DailyWaterTargetMl),
		Status:             entity.PlanStatus(p.Status),
		Days:               []*entity.DietPlanDay{},
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

// dietPlanDayToDomain converts a sqlc DietPlanDay row to a domain entity.
func dietPlanDayToDomain(d db.DietPlanDay) *entity.DietPlanDay {
	return &entity.DietPlanDay{
		ID:        d.ID,
		PlanID:    d.PlanID,
		DayNumber: int(d.DayNumber),
		Meals:     []*entity.DietMeal{},
		CreatedAt: d.CreatedAt,
	}
}

// dietMealToDomain converts a sqlc DietMeal row to a domain entity.
func dietMealToDomain(m db.DietMeal) *entity.DietMeal {
	return &entity.DietMeal{
		ID:            m.ID,
		DayID:         m.DayID,
		Title:         m.Title,
		ScheduledTime: pgtimeToString(m.ScheduledTime),
		DisplayOrder:  int(m.DisplayOrder),
		Options:       []*entity.MealOption{},
		CreatedAt:     m.CreatedAt,
	}
}

// mealOptionToDomain converts a sqlc MealOption row to a domain entity.
func mealOptionToDomain(o db.MealOption) *entity.MealOption {
	return &entity.MealOption{
		ID:           o.ID,
		MealID:       o.MealID,
		OptionNumber: int(o.OptionNumber),
		Items:        []*entity.MealOptionItem{},
		CreatedAt:    o.CreatedAt,
	}
}

// mealOptionItemToDomain converts a sqlc MealOptionItem row to a domain entity.
func mealOptionItemToDomain(i db.MealOptionItem) *entity.MealOptionItem {
	return &entity.MealOptionItem{
		ID:        i.ID,
		OptionID:  i.OptionID,
		FoodID:    i.FoodID,
		Quantity:  numericToFloat64(i.Quantity),
		Unit:      i.Unit,
		Notes:     i.Notes,
		CreatedAt: i.CreatedAt,
	}
}

// pgtimeToString converts pgtype.Time (microseconds since midnight) to "HH:MM" string.
func pgtimeToString(t pgtype.Time) string {
	if !t.Valid {
		return "08:00"
	}
	totalSec := t.Microseconds / 1_000_000
	h := totalSec / 3600
	m := (totalSec % 3600) / 60
	return fmt.Sprintf("%02d:%02d", h, m)
}

// stringToPgtime converts a "HH:MM" string to pgtype.Time.
func stringToPgtime(s string) pgtype.Time {
	var h, m int
	fmt.Sscanf(s, "%d:%d", &h, &m)
	return pgtype.Time{Microseconds: int64(h*3600+m*60) * 1_000_000, Valid: true}
}

// numericToFloat64 converts pgtype.Numeric to float64, returning 0 on error.
func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return 0
	}
	return f.Float64
}

// float64ToNumeric converts a float64 value to pgtype.Numeric.
func float64ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(f, 'f', 2, 64))
	return n
}

// mealOptionItemWithFoodToDomain converts a join row to a domain entity with food snapshot and computed nutrition.
func mealOptionItemWithFoodToDomain(i db.ListMealOptionItemsWithFoodRow) *entity.MealOptionItem {
	qty := numericToFloat64(i.Quantity)
	cal := numericToFloat64(i.FoodCalories)
	prot := numericToFloat64(i.FoodProtein)
	carb := numericToFloat64(i.FoodCarbohydrate)
	fat := numericToFloat64(i.FoodFat)
	fib := numericToFloat64(i.FoodFiber)

	food := &entity.FoodSnapshot{
		ID:           i.FoodID,
		Name:         i.FoodName,
		Unit:         i.FoodUnit,
		Calories:     cal,
		Protein:      prot,
		Carbohydrate: carb,
		Fat:          fat,
		Fiber:        fib,
	}

	computed := &entity.NutritionalSummary{
		Calories: cal * qty,
		Protein:  prot * qty,
		Carbs:    carb * qty,
		Fat:      fat * qty,
		Fiber:    fib * qty,
	}

	return &entity.MealOptionItem{
		ID:        i.ID,
		OptionID:  i.OptionID,
		FoodID:    i.FoodID,
		Quantity:  qty,
		Unit:      i.Unit,
		Notes:     i.Notes,
		Food:      food,
		Computed:  computed,
		CreatedAt: i.CreatedAt,
	}
}
