package tracking

import (
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranjbar-dev/nutritrack/internal/domain/tracking/entity"
	db "github.com/ranjbar-dev/nutritrack/internal/infrastructure/persistence/sqlc"
)

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

func float64ToNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(strconv.FormatFloat(f, 'f', 3, 64))
	return n
}

func optFloat64ToNumeric(f *float64) *pgtype.Numeric {
	if f == nil {
		return nil
	}
	n := float64ToNumeric(*f)
	return &n
}

func optNumericToFloat64(n *pgtype.Numeric) *float64 {
	if n == nil || !n.Valid {
		return nil
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return nil
	}
	v := f.Float64
	return &v
}

func foodLogToDomain(r db.FoodLog) *entity.FoodLog {
	return &entity.FoodLog{
		ID:         r.ID,
		ClientID:   r.ClientID,
		LocalID:    r.LocalID,
		LoggedAt:   r.LoggedAt,
		LoggedDate: r.LoggedDate,
		FoodID:     r.FoodID,
		FoodName:   r.FoodName,
		Quantity:   numericToFloat64(r.Quantity),
		Unit:       r.Unit,
		Calories:   numericToFloat64(r.Calories),
		Protein:    numericToFloat64(r.Protein),
		Carbs:      numericToFloat64(r.Carbs),
		Fat:        numericToFloat64(r.Fat),
		Notes:      r.Notes,
		CreatedAt:  r.CreatedAt,
	}
}

func waterLogToDomain(r db.WaterLog) *entity.WaterLog {
	return &entity.WaterLog{
		ID:         r.ID,
		ClientID:   r.ClientID,
		LocalID:    r.LocalID,
		LoggedAt:   r.LoggedAt,
		LoggedDate: r.LoggedDate,
		AmountMl:   int(r.AmountMl),
		Notes:      r.Notes,
		CreatedAt:  r.CreatedAt,
	}
}

func sleepLogToDomain(r db.SleepLog) *entity.SleepLog {
	return &entity.SleepLog{
		ID:              r.ID,
		ClientID:        r.ClientID,
		LocalID:         r.LocalID,
		LoggedDate:      r.LoggedDate,
		SleepStart:      r.SleepStart,
		SleepEnd:        r.SleepEnd,
		DurationMinutes: int(r.DurationMinutes),
		Quality:         int(r.Quality),
		Notes:           r.Notes,
		CreatedAt:       r.CreatedAt,
	}
}

func exerciseLogToDomain(r db.ExerciseLog) *entity.ExerciseLog {
	return &entity.ExerciseLog{
		ID:              r.ID,
		ClientID:        r.ClientID,
		LocalID:         r.LocalID,
		LoggedAt:        r.LoggedAt,
		LoggedDate:      r.LoggedDate,
		ExerciseName:    r.ExerciseName,
		DurationMinutes: int(r.DurationMinutes),
		CaloriesBurned:  int(r.CaloriesBurned),
		Notes:           r.Notes,
		CreatedAt:       r.CreatedAt,
	}
}

func medicationLogToDomain(r db.MedicationLog) *entity.MedicationLog {
	return &entity.MedicationLog{
		ID:             r.ID,
		ClientID:       r.ClientID,
		LocalID:        r.LocalID,
		LoggedAt:       r.LoggedAt,
		LoggedDate:     r.LoggedDate,
		MedicationID:   r.MedicationID,
		MedicationName: r.MedicationName,
		Dosage:         r.Dosage,
		Notes:          r.Notes,
		CreatedAt:      r.CreatedAt,
	}
}

func bodyMeasurementToDomain(r db.BodyMeasurement) *entity.BodyMeasurement {
	return &entity.BodyMeasurement{
		ID:           r.ID,
		ClientID:     r.ClientID,
		LocalID:      r.LocalID,
		MeasuredAt:   r.MeasuredAt,
		MeasuredDate: r.MeasuredDate,
		WeightKg:     optNumericToFloat64(r.WeightKg),
		HeightCm:     optNumericToFloat64(r.HeightCm),
		WaistCm:      optNumericToFloat64(r.WaistCm),
		HipCm:        optNumericToFloat64(r.HipCm),
		ChestCm:      optNumericToFloat64(r.ChestCm),
		ArmCm:        optNumericToFloat64(r.ArmCm),
		Notes:        r.Notes,
		CreatedAt:    r.CreatedAt,
	}
}
