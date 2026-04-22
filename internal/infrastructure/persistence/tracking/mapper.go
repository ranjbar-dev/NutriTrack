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
	return entity.ReconstituteFoodLog(
		r.ID,
		r.ClientID,
		r.LocalID,
		r.LoggedAt,
		r.LoggedDate,
		r.FoodID,
		r.FoodName,
		numericToFloat64(r.Quantity),
		r.Unit,
		numericToFloat64(r.Calories),
		numericToFloat64(r.Protein),
		numericToFloat64(r.Carbs),
		numericToFloat64(r.Fat),
		r.Notes,
		r.CreatedAt,
	)
}

func waterLogToDomain(r db.WaterLog) *entity.WaterLog {
	return entity.ReconstituteWaterLog(
		r.ID,
		r.ClientID,
		r.LocalID,
		r.LoggedAt,
		r.LoggedDate,
		int(r.AmountMl),
		r.Notes,
		r.CreatedAt,
	)
}

func sleepLogToDomain(r db.SleepLog) *entity.SleepLog {
	return entity.ReconstituteSleepLog(
		r.ID,
		r.ClientID,
		r.LocalID,
		r.LoggedDate,
		r.SleepStart,
		r.SleepEnd,
		int(r.DurationMinutes),
		int(r.Quality),
		r.Notes,
		r.CreatedAt,
	)
}

func exerciseLogToDomain(r db.ExerciseLog) *entity.ExerciseLog {
	return entity.ReconstituteExerciseLog(
		r.ID,
		r.ClientID,
		r.LocalID,
		r.LoggedAt,
		r.LoggedDate,
		r.ExerciseName,
		int(r.DurationMinutes),
		int(r.CaloriesBurned),
		r.Notes,
		r.CreatedAt,
	)
}

func medicationLogToDomain(r db.MedicationLog) *entity.MedicationLog {
	return entity.ReconstituteMedicationLog(
		r.ID,
		r.ClientID,
		r.LocalID,
		r.LoggedAt,
		r.LoggedDate,
		r.MedicationID,
		r.MedicationName,
		r.Dosage,
		r.Notes,
		r.CreatedAt,
	)
}

func bodyMeasurementToDomain(r db.BodyMeasurement) *entity.BodyMeasurement {
	return entity.ReconstituteBodyMeasurement(
		r.ID,
		r.ClientID,
		r.LocalID,
		r.MeasuredAt,
		r.MeasuredDate,
		optNumericToFloat64(r.WeightKg),
		optNumericToFloat64(r.HeightCm),
		optNumericToFloat64(r.WaistCm),
		optNumericToFloat64(r.HipCm),
		optNumericToFloat64(r.ChestCm),
		optNumericToFloat64(r.ArmCm),
		r.Notes,
		r.CreatedAt,
	)
}
