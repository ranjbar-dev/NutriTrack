package tracking

import (
	"github.com/ranjbar-dev/nutritrack/internal/domain/tracking/entity"
)

// MapFoodLog converts a FoodLog entity to a JSON-serialisable map.
func MapFoodLog(l *entity.FoodLog) map[string]any {
	return map[string]any{
		"id":          l.ID(),
		"client_id":   l.ClientID(),
		"local_id":    l.LocalID(),
		"logged_at":   l.LoggedAt(),
		"logged_date": l.LoggedDate().Format("2006-01-02"),
		"food_id":     l.FoodID(),
		"food_name":   l.FoodName(),
		"quantity":    l.Quantity(),
		"unit":        l.Unit(),
		"calories":    l.Calories(),
		"protein":     l.Protein(),
		"carbs":       l.Carbs(),
		"fat":         l.Fat(),
		"notes":       l.Notes(),
		"created_at":  l.CreatedAt(),
	}
}

// MapWaterLog converts a WaterLog entity to a JSON-serialisable map.
func MapWaterLog(l *entity.WaterLog) map[string]any {
	return map[string]any{
		"id":          l.ID(),
		"client_id":   l.ClientID(),
		"local_id":    l.LocalID(),
		"logged_at":   l.LoggedAt(),
		"logged_date": l.LoggedDate().Format("2006-01-02"),
		"amount_ml":   l.AmountMl(),
		"notes":       l.Notes(),
		"created_at":  l.CreatedAt(),
	}
}

// MapSleepLog converts a SleepLog entity to a JSON-serialisable map.
func MapSleepLog(l *entity.SleepLog) map[string]any {
	return map[string]any{
		"id":               l.ID(),
		"client_id":        l.ClientID(),
		"local_id":         l.LocalID(),
		"logged_date":      l.LoggedDate().Format("2006-01-02"),
		"sleep_start":      l.SleepStart(),
		"sleep_end":        l.SleepEnd(),
		"duration_minutes": l.DurationMinutes(),
		"quality":          l.Quality(),
		"notes":            l.Notes(),
		"created_at":       l.CreatedAt(),
	}
}

// MapExerciseLog converts an ExerciseLog entity to a JSON-serialisable map.
func MapExerciseLog(l *entity.ExerciseLog) map[string]any {
	return map[string]any{
		"id":               l.ID(),
		"client_id":        l.ClientID(),
		"local_id":         l.LocalID(),
		"logged_at":        l.LoggedAt(),
		"logged_date":      l.LoggedDate().Format("2006-01-02"),
		"exercise_name":    l.ExerciseName(),
		"duration_minutes": l.DurationMinutes(),
		"calories_burned":  l.CaloriesBurned(),
		"notes":            l.Notes(),
		"created_at":       l.CreatedAt(),
	}
}

// MapMedicationLog converts a MedicationLog entity to a JSON-serialisable map.
func MapMedicationLog(l *entity.MedicationLog) map[string]any {
	return map[string]any{
		"id":              l.ID(),
		"client_id":       l.ClientID(),
		"local_id":        l.LocalID(),
		"logged_at":       l.LoggedAt(),
		"logged_date":     l.LoggedDate().Format("2006-01-02"),
		"medication_id":   l.MedicationID(),
		"medication_name": l.MedicationName(),
		"dosage":          l.Dosage(),
		"notes":           l.Notes(),
		"created_at":      l.CreatedAt(),
	}
}

// MapBodyMeasurement converts a BodyMeasurement entity to a JSON-serialisable map.
func MapBodyMeasurement(m *entity.BodyMeasurement) map[string]any {
	return map[string]any{
		"id":            m.ID(),
		"client_id":     m.ClientID(),
		"local_id":      m.LocalID(),
		"measured_at":   m.MeasuredAt(),
		"measured_date": m.MeasuredDate().Format("2006-01-02"),
		"weight_kg":     m.WeightKg(),
		"height_cm":     m.HeightCm(),
		"waist_cm":      m.WaistCm(),
		"hip_cm":        m.HipCm(),
		"chest_cm":      m.ChestCm(),
		"arm_cm":        m.ArmCm(),
		"notes":         m.Notes(),
		"created_at":    m.CreatedAt(),
	}
}
