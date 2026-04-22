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
	return entity.ReconstituteDietPlan(
		p.ID,
		p.ClientID,
		p.NutritionistID,
		p.Title,
		p.StartDate,
		p.EndDate,
		p.Notes,
		int(p.DailyWaterTargetMl),
		entity.PlanStatus(p.Status),
		[]*entity.DietPlanDay{},
		p.CreatedAt,
		p.UpdatedAt,
	)
}

// dietPlanDayToDomain converts a sqlc DietPlanDay row to a domain entity.
func dietPlanDayToDomain(d db.DietPlanDay) *entity.DietPlanDay {
	return entity.ReconstituteDietPlanDay(
		d.ID,
		d.PlanID,
		int(d.DayNumber),
		[]*entity.DietMeal{},
		nil,
		[]*entity.ExerciseRecommendation{},
		[]*entity.PrescribedMedication{},
		d.CreatedAt,
	)
}

// dietMealToDomain converts a sqlc DietMeal row to a domain entity.
func dietMealToDomain(m db.DietMeal) *entity.DietMeal {
	return entity.ReconstituteDietMeal(
		m.ID,
		m.DayID,
		m.Title,
		pgtimeToString(m.ScheduledTime),
		int(m.DisplayOrder),
		[]*entity.MealOption{},
		nil,
		m.CreatedAt,
	)
}

// mealOptionToDomain converts a sqlc MealOption row to a domain entity.
func mealOptionToDomain(o db.MealOption) *entity.MealOption {
	return entity.ReconstituteMealOption(
		o.ID,
		o.MealID,
		int(o.OptionNumber),
		[]*entity.MealOptionItem{},
		nil,
		o.CreatedAt,
	)
}

// mealOptionItemToDomain converts a sqlc MealOptionItem row to a domain entity.
func mealOptionItemToDomain(i db.MealOptionItem) *entity.MealOptionItem {
	return entity.ReconstituteMealOptionItem(
		i.ID,
		i.OptionID,
		i.FoodID,
		numericToFloat64(i.Quantity),
		i.Unit,
		i.Notes,
		nil,
		nil,
		i.CreatedAt,
	)
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

// exerciseToDomain converts a sqlc ExerciseRecommendation row to a domain entity.
func exerciseToDomain(e db.ExerciseRecommendation) *entity.ExerciseRecommendation {
	return entity.ReconstituteExerciseRecommendation(
		e.ID,
		e.DayID,
		e.ExerciseName,
		int(e.DurationMinutes),
		e.Description,
		int(e.CaloriesBurnEstimate),
		e.CreatedAt,
	)
}

// prescriptionWithMedToDomain converts a join row to a domain PrescribedMedication with medication snapshot.
func prescriptionWithMedToDomain(r db.ListDayPrescribedMedicationsWithMedicationRow) *entity.PrescribedMedication {
	med := &entity.MedicationSnapshot{
		ID:   r.MedicationID,
		Name: r.MedicationName,
		Unit: r.MedicationUnit,
	}
	return entity.ReconstitutePrescribedMedication(
		r.ID,
		r.DayID,
		r.MedicationID,
		med,
		r.Dosage,
		r.Frequency,
		r.Times,
		r.Instructions,
		r.StartDate,
		r.EndDate,
		r.CreatedAt,
	)
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

	return entity.ReconstituteMealOptionItem(
		i.ID,
		i.OptionID,
		i.FoodID,
		qty,
		i.Unit,
		i.Notes,
		food,
		computed,
		i.CreatedAt,
	)
}
