package entity

import (
	"time"

	"github.com/google/uuid"
)

// FoodLog represents a client's food consumption record.
type FoodLog struct {
	id         uuid.UUID
	clientID   uuid.UUID
	localID    string
	loggedAt   time.Time
	loggedDate time.Time
	foodID     *uuid.UUID
	foodName   string
	quantity   float64
	unit       string
	calories   float64
	protein    float64
	carbs      float64
	fat        float64
	notes      string
	createdAt  time.Time
}

// NewFoodLog creates a new FoodLog (ID assigned by DB).
func NewFoodLog(clientID uuid.UUID, localID string, loggedAt, loggedDate time.Time, foodID *uuid.UUID, foodName string, quantity float64, unit string, calories, protein, carbs, fat float64, notes string) *FoodLog {
	return &FoodLog{clientID: clientID, localID: localID, loggedAt: loggedAt, loggedDate: loggedDate, foodID: foodID, foodName: foodName, quantity: quantity, unit: unit, calories: calories, protein: protein, carbs: carbs, fat: fat, notes: notes}
}

// ReconstituteFoodLog rebuilds a FoodLog from stored data.
func ReconstituteFoodLog(id, clientID uuid.UUID, localID string, loggedAt, loggedDate time.Time, foodID *uuid.UUID, foodName string, quantity float64, unit string, calories, protein, carbs, fat float64, notes string, createdAt time.Time) *FoodLog {
	return &FoodLog{id: id, clientID: clientID, localID: localID, loggedAt: loggedAt, loggedDate: loggedDate, foodID: foodID, foodName: foodName, quantity: quantity, unit: unit, calories: calories, protein: protein, carbs: carbs, fat: fat, notes: notes, createdAt: createdAt}
}

func (f *FoodLog) ID() uuid.UUID         { return f.id }
func (f *FoodLog) ClientID() uuid.UUID   { return f.clientID }
func (f *FoodLog) LocalID() string       { return f.localID }
func (f *FoodLog) LoggedAt() time.Time   { return f.loggedAt }
func (f *FoodLog) LoggedDate() time.Time { return f.loggedDate }
func (f *FoodLog) FoodID() *uuid.UUID    { return f.foodID }
func (f *FoodLog) FoodName() string      { return f.foodName }
func (f *FoodLog) Quantity() float64     { return f.quantity }
func (f *FoodLog) Unit() string          { return f.unit }
func (f *FoodLog) Calories() float64     { return f.calories }
func (f *FoodLog) Protein() float64      { return f.protein }
func (f *FoodLog) Carbs() float64        { return f.carbs }
func (f *FoodLog) Fat() float64          { return f.fat }
func (f *FoodLog) Notes() string         { return f.notes }
func (f *FoodLog) CreatedAt() time.Time  { return f.createdAt }

// WaterLog represents a client's water intake record.
type WaterLog struct {
	id         uuid.UUID
	clientID   uuid.UUID
	localID    string
	loggedAt   time.Time
	loggedDate time.Time
	amountMl   int
	notes      string
	createdAt  time.Time
}

// NewWaterLog creates a new WaterLog (ID assigned by DB).
func NewWaterLog(clientID uuid.UUID, localID string, loggedAt, loggedDate time.Time, amountMl int, notes string) *WaterLog {
	return &WaterLog{clientID: clientID, localID: localID, loggedAt: loggedAt, loggedDate: loggedDate, amountMl: amountMl, notes: notes}
}

// ReconstituteWaterLog rebuilds a WaterLog from stored data.
func ReconstituteWaterLog(id, clientID uuid.UUID, localID string, loggedAt, loggedDate time.Time, amountMl int, notes string, createdAt time.Time) *WaterLog {
	return &WaterLog{id: id, clientID: clientID, localID: localID, loggedAt: loggedAt, loggedDate: loggedDate, amountMl: amountMl, notes: notes, createdAt: createdAt}
}

func (w *WaterLog) ID() uuid.UUID         { return w.id }
func (w *WaterLog) ClientID() uuid.UUID   { return w.clientID }
func (w *WaterLog) LocalID() string       { return w.localID }
func (w *WaterLog) LoggedAt() time.Time   { return w.loggedAt }
func (w *WaterLog) LoggedDate() time.Time { return w.loggedDate }
func (w *WaterLog) AmountMl() int         { return w.amountMl }
func (w *WaterLog) Notes() string         { return w.notes }
func (w *WaterLog) CreatedAt() time.Time  { return w.createdAt }

// SleepLog represents a client's sleep record.
type SleepLog struct {
	id              uuid.UUID
	clientID        uuid.UUID
	localID         string
	loggedDate      time.Time
	sleepStart      time.Time
	sleepEnd        time.Time
	durationMinutes int
	quality         int
	notes           string
	createdAt       time.Time
}

// NewSleepLog creates a new SleepLog (ID assigned by DB).
func NewSleepLog(clientID uuid.UUID, localID string, loggedDate, sleepStart, sleepEnd time.Time, durationMinutes, quality int, notes string) *SleepLog {
	return &SleepLog{clientID: clientID, localID: localID, loggedDate: loggedDate, sleepStart: sleepStart, sleepEnd: sleepEnd, durationMinutes: durationMinutes, quality: quality, notes: notes}
}

// ReconstituteSleepLog rebuilds a SleepLog from stored data.
func ReconstituteSleepLog(id, clientID uuid.UUID, localID string, loggedDate, sleepStart, sleepEnd time.Time, durationMinutes, quality int, notes string, createdAt time.Time) *SleepLog {
	return &SleepLog{id: id, clientID: clientID, localID: localID, loggedDate: loggedDate, sleepStart: sleepStart, sleepEnd: sleepEnd, durationMinutes: durationMinutes, quality: quality, notes: notes, createdAt: createdAt}
}

func (s *SleepLog) ID() uuid.UUID         { return s.id }
func (s *SleepLog) ClientID() uuid.UUID   { return s.clientID }
func (s *SleepLog) LocalID() string       { return s.localID }
func (s *SleepLog) LoggedDate() time.Time { return s.loggedDate }
func (s *SleepLog) SleepStart() time.Time { return s.sleepStart }
func (s *SleepLog) SleepEnd() time.Time   { return s.sleepEnd }
func (s *SleepLog) DurationMinutes() int  { return s.durationMinutes }
func (s *SleepLog) Quality() int          { return s.quality }
func (s *SleepLog) Notes() string         { return s.notes }
func (s *SleepLog) CreatedAt() time.Time  { return s.createdAt }

// ExerciseLog represents a client's exercise session.
type ExerciseLog struct {
	id              uuid.UUID
	clientID        uuid.UUID
	localID         string
	loggedAt        time.Time
	loggedDate      time.Time
	exerciseName    string
	durationMinutes int
	caloriesBurned  int
	notes           string
	createdAt       time.Time
}

// NewExerciseLog creates a new ExerciseLog (ID assigned by DB).
func NewExerciseLog(clientID uuid.UUID, localID string, loggedAt, loggedDate time.Time, exerciseName string, durationMinutes, caloriesBurned int, notes string) *ExerciseLog {
	return &ExerciseLog{clientID: clientID, localID: localID, loggedAt: loggedAt, loggedDate: loggedDate, exerciseName: exerciseName, durationMinutes: durationMinutes, caloriesBurned: caloriesBurned, notes: notes}
}

// ReconstituteExerciseLog rebuilds an ExerciseLog from stored data.
func ReconstituteExerciseLog(id, clientID uuid.UUID, localID string, loggedAt, loggedDate time.Time, exerciseName string, durationMinutes, caloriesBurned int, notes string, createdAt time.Time) *ExerciseLog {
	return &ExerciseLog{id: id, clientID: clientID, localID: localID, loggedAt: loggedAt, loggedDate: loggedDate, exerciseName: exerciseName, durationMinutes: durationMinutes, caloriesBurned: caloriesBurned, notes: notes, createdAt: createdAt}
}

func (e *ExerciseLog) ID() uuid.UUID         { return e.id }
func (e *ExerciseLog) ClientID() uuid.UUID   { return e.clientID }
func (e *ExerciseLog) LocalID() string       { return e.localID }
func (e *ExerciseLog) LoggedAt() time.Time   { return e.loggedAt }
func (e *ExerciseLog) LoggedDate() time.Time { return e.loggedDate }
func (e *ExerciseLog) ExerciseName() string  { return e.exerciseName }
func (e *ExerciseLog) DurationMinutes() int  { return e.durationMinutes }
func (e *ExerciseLog) CaloriesBurned() int   { return e.caloriesBurned }
func (e *ExerciseLog) Notes() string         { return e.notes }
func (e *ExerciseLog) CreatedAt() time.Time  { return e.createdAt }

// MedicationLog represents a client's medication intake record.
type MedicationLog struct {
	id             uuid.UUID
	clientID       uuid.UUID
	localID        string
	loggedAt       time.Time
	loggedDate     time.Time
	medicationID   *uuid.UUID
	medicationName string
	dosage         string
	notes          string
	createdAt      time.Time
}

// NewMedicationLog creates a new MedicationLog (ID assigned by DB).
func NewMedicationLog(clientID uuid.UUID, localID string, loggedAt, loggedDate time.Time, medicationID *uuid.UUID, medicationName, dosage, notes string) *MedicationLog {
	return &MedicationLog{clientID: clientID, localID: localID, loggedAt: loggedAt, loggedDate: loggedDate, medicationID: medicationID, medicationName: medicationName, dosage: dosage, notes: notes}
}

// ReconstituteMedicationLog rebuilds a MedicationLog from stored data.
func ReconstituteMedicationLog(id, clientID uuid.UUID, localID string, loggedAt, loggedDate time.Time, medicationID *uuid.UUID, medicationName, dosage, notes string, createdAt time.Time) *MedicationLog {
	return &MedicationLog{id: id, clientID: clientID, localID: localID, loggedAt: loggedAt, loggedDate: loggedDate, medicationID: medicationID, medicationName: medicationName, dosage: dosage, notes: notes, createdAt: createdAt}
}

func (m *MedicationLog) ID() uuid.UUID            { return m.id }
func (m *MedicationLog) ClientID() uuid.UUID      { return m.clientID }
func (m *MedicationLog) LocalID() string          { return m.localID }
func (m *MedicationLog) LoggedAt() time.Time      { return m.loggedAt }
func (m *MedicationLog) LoggedDate() time.Time    { return m.loggedDate }
func (m *MedicationLog) MedicationID() *uuid.UUID { return m.medicationID }
func (m *MedicationLog) MedicationName() string   { return m.medicationName }
func (m *MedicationLog) Dosage() string           { return m.dosage }
func (m *MedicationLog) Notes() string            { return m.notes }
func (m *MedicationLog) CreatedAt() time.Time     { return m.createdAt }

// BodyMeasurement represents a client's body measurement snapshot.
type BodyMeasurement struct {
	id           uuid.UUID
	clientID     uuid.UUID
	localID      string
	measuredAt   time.Time
	measuredDate time.Time
	weightKg     *float64
	heightCm     *float64
	waistCm      *float64
	hipCm        *float64
	chestCm      *float64
	armCm        *float64
	notes        string
	createdAt    time.Time
}

// NewBodyMeasurement creates a new BodyMeasurement (ID assigned by DB).
func NewBodyMeasurement(clientID uuid.UUID, localID string, measuredAt, measuredDate time.Time, weightKg, heightCm, waistCm, hipCm, chestCm, armCm *float64, notes string) *BodyMeasurement {
	return &BodyMeasurement{clientID: clientID, localID: localID, measuredAt: measuredAt, measuredDate: measuredDate, weightKg: weightKg, heightCm: heightCm, waistCm: waistCm, hipCm: hipCm, chestCm: chestCm, armCm: armCm, notes: notes}
}

// ReconstituteBodyMeasurement rebuilds a BodyMeasurement from stored data.
func ReconstituteBodyMeasurement(id, clientID uuid.UUID, localID string, measuredAt, measuredDate time.Time, weightKg, heightCm, waistCm, hipCm, chestCm, armCm *float64, notes string, createdAt time.Time) *BodyMeasurement {
	return &BodyMeasurement{id: id, clientID: clientID, localID: localID, measuredAt: measuredAt, measuredDate: measuredDate, weightKg: weightKg, heightCm: heightCm, waistCm: waistCm, hipCm: hipCm, chestCm: chestCm, armCm: armCm, notes: notes, createdAt: createdAt}
}

func (b *BodyMeasurement) ID() uuid.UUID           { return b.id }
func (b *BodyMeasurement) ClientID() uuid.UUID     { return b.clientID }
func (b *BodyMeasurement) LocalID() string         { return b.localID }
func (b *BodyMeasurement) MeasuredAt() time.Time   { return b.measuredAt }
func (b *BodyMeasurement) MeasuredDate() time.Time { return b.measuredDate }
func (b *BodyMeasurement) WeightKg() *float64      { return b.weightKg }
func (b *BodyMeasurement) HeightCm() *float64      { return b.heightCm }
func (b *BodyMeasurement) WaistCm() *float64       { return b.waistCm }
func (b *BodyMeasurement) HipCm() *float64         { return b.hipCm }
func (b *BodyMeasurement) ChestCm() *float64       { return b.chestCm }
func (b *BodyMeasurement) ArmCm() *float64         { return b.armCm }
func (b *BodyMeasurement) Notes() string           { return b.notes }
func (b *BodyMeasurement) CreatedAt() time.Time    { return b.createdAt }
