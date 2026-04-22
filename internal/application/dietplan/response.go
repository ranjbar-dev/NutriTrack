package dietplan

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/entity"
)

// MapDietPlan builds a flat plan response map (no days tree).
func MapDietPlan(p *entity.DietPlan) map[string]any {
	return map[string]any{
		"id":                    p.ID(),
		"client_id":             p.ClientID(),
		"nutritionist_id":       p.NutritionistID(),
		"title":                 p.Title(),
		"start_date":            p.StartDate().Format("2006-01-02"),
		"end_date":              p.EndDate().Format("2006-01-02"),
		"notes":                 p.Notes(),
		"daily_water_target_ml": p.DailyWaterTargetML(),
		"status":                p.Status(),
		"created_at":            p.CreatedAt(),
		"updated_at":            p.UpdatedAt(),
	}
}

// MapDietPlanFull builds a full plan response including the nested days tree.
func MapDietPlanFull(plan *entity.DietPlan) map[string]any {
	days := make([]any, len(plan.Days()))
	for i, day := range plan.Days() {
		meals := make([]any, len(day.Meals()))
		for j, meal := range day.Meals() {
			options := make([]any, len(meal.Options()))
			for k, opt := range meal.Options() {
				options[k] = map[string]any{
					"id":            opt.ID(),
					"option_number": opt.OptionNumber(),
					"totals":        mapNutritionalSummary(opt.Totals()),
					"items":         mapMealOptionItems(opt.Items()),
				}
			}
			meals[j] = map[string]any{
				"id":             meal.ID(),
				"title":          meal.Title(),
				"scheduled_time": meal.ScheduledTime(),
				"display_order":  meal.DisplayOrder(),
				"total_range":    mapNutritionalRange(meal.TotalRange()),
				"options":        options,
			}
		}
		days[i] = map[string]any{
			"id":            day.ID(),
			"day_number":    day.DayNumber(),
			"total_range":   mapNutritionalRange(day.TotalRange()),
			"meals":         meals,
			"exercises":     mapExerciseRecommendations(day.Exercises()),
			"prescriptions": mapPrescribedMedications(day.Prescriptions()),
		}
	}

	return map[string]any{
		"id":                    plan.ID(),
		"client_id":             plan.ClientID(),
		"nutritionist_id":       plan.NutritionistID(),
		"title":                 plan.Title(),
		"start_date":            plan.StartDate().Format("2006-01-02"),
		"end_date":              plan.EndDate().Format("2006-01-02"),
		"notes":                 plan.Notes(),
		"daily_water_target_ml": plan.DailyWaterTargetML(),
		"status":                plan.Status(),
		"days":                  days,
		"created_at":            plan.CreatedAt(),
		"updated_at":            plan.UpdatedAt(),
	}
}

func mapNutritionalSummary(s *entity.NutritionalSummary) map[string]any {
	if s == nil {
		return nil
	}
	return map[string]any{
		"calories": s.Calories,
		"protein":  s.Protein,
		"carbs":    s.Carbs,
		"fat":      s.Fat,
		"fiber":    s.Fiber,
	}
}

func mapNutritionalRange(r *entity.NutritionalRange) map[string]any {
	if r == nil {
		return nil
	}
	return map[string]any{
		"min": mapNutritionalSummary(&r.Min),
		"max": mapNutritionalSummary(&r.Max),
	}
}

func mapMealOptionItems(items []*entity.MealOptionItem) []any {
	result := make([]any, len(items))
	for i, item := range items {
		m := map[string]any{
			"id":         item.ID(),
			"option_id":  item.OptionID(),
			"food_id":    item.FoodID(),
			"quantity":   item.Quantity(),
			"unit":       item.Unit(),
			"notes":      item.Notes(),
			"computed":   mapNutritionalSummary(item.Computed()),
			"created_at": item.CreatedAt(),
		}
		if item.Food() != nil {
			m["food"] = map[string]any{
				"id":   item.Food().ID,
				"name": item.Food().Name,
				"unit": item.Food().Unit,
			}
		}
		result[i] = m
	}
	return result
}

func mapExerciseRecommendations(exercises []*entity.ExerciseRecommendation) []any {
	result := make([]any, len(exercises))
	for i, e := range exercises {
		result[i] = map[string]any{
			"id":                     e.ID(),
			"exercise_name":          e.ExerciseName(),
			"duration_minutes":       e.DurationMinutes(),
			"description":            e.Description(),
			"calories_burn_estimate": e.CaloriesBurnEstimate(),
			"created_at":             e.CreatedAt(),
		}
	}
	return result
}

func mapPrescribedMedications(prescriptions []*entity.PrescribedMedication) []any {
	result := make([]any, len(prescriptions))
	for i, rx := range prescriptions {
		m := map[string]any{
			"id":            rx.ID(),
			"medication_id": rx.MedicationID(),
			"dosage":        rx.Dosage(),
			"frequency":     rx.Frequency(),
			"times":         rx.Times(),
			"instructions":  rx.Instructions(),
			"start_date":    rx.StartDate(),
			"end_date":      rx.EndDate(),
			"created_at":    rx.CreatedAt(),
		}
		if rx.Medication() != nil {
			m["medication"] = map[string]any{
				"id":   rx.Medication().ID,
				"name": rx.Medication().Name,
				"unit": rx.Medication().Unit,
			}
		}
		result[i] = m
	}
	return result
}
