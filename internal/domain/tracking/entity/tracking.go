package entity

import (
	"time"

	"github.com/google/uuid"
)

// FoodLog represents a client's food consumption record.
type FoodLog struct {
	ID         uuid.UUID
	ClientID   uuid.UUID
	LocalID    string
	LoggedAt   time.Time
	LoggedDate time.Time // date portion (Asia/Tehran)
	FoodID     *uuid.UUID
	FoodName   string
	Quantity   float64
	Unit       string
	Calories   float64
	Protein    float64
	Carbs      float64
	Fat        float64
	Notes      string
	CreatedAt  time.Time
}

// WaterLog represents a client's water intake record.
type WaterLog struct {
	ID         uuid.UUID
	ClientID   uuid.UUID
	LocalID    string
	LoggedAt   time.Time
	LoggedDate time.Time
	AmountMl   int
	Notes      string
	CreatedAt  time.Time
}

// SleepLog represents a client's sleep record.
type SleepLog struct {
	ID              uuid.UUID
	ClientID        uuid.UUID
	LocalID         string
	LoggedDate      time.Time
	SleepStart      time.Time
	SleepEnd        time.Time
	DurationMinutes int
	Quality         int // 1-5
	Notes           string
	CreatedAt       time.Time
}

// ExerciseLog represents a client's exercise session.
type ExerciseLog struct {
	ID              uuid.UUID
	ClientID        uuid.UUID
	LocalID         string
	LoggedAt        time.Time
	LoggedDate      time.Time
	ExerciseName    string
	DurationMinutes int
	CaloriesBurned  int
	Notes           string
	CreatedAt       time.Time
}

// MedicationLog represents a client's medication intake record.
type MedicationLog struct {
	ID             uuid.UUID
	ClientID       uuid.UUID
	LocalID        string
	LoggedAt       time.Time
	LoggedDate     time.Time
	MedicationID   *uuid.UUID
	MedicationName string
	Dosage         string
	Notes          string
	CreatedAt      time.Time
}

// BodyMeasurement represents a client's body measurement snapshot.
type BodyMeasurement struct {
	ID           uuid.UUID
	ClientID     uuid.UUID
	LocalID      string
	MeasuredAt   time.Time
	MeasuredDate time.Time
	WeightKg     *float64
	HeightCm     *float64
	WaistCm      *float64
	HipCm        *float64
	ChestCm      *float64
	ArmCm        *float64
	Notes        string
	CreatedAt    time.Time
}
