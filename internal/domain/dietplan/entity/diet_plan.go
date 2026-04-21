package entity

import (
	"time"

	"github.com/google/uuid"
)

// NutritionalSummary holds computed macro totals.
type NutritionalSummary struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}

// NutritionalRange holds min/max nutritional totals across options.
type NutritionalRange struct {
	Min NutritionalSummary `json:"min"`
	Max NutritionalSummary `json:"max"`
}

// FoodSnapshot carries the food details needed for nutritional computation.
type FoodSnapshot struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Unit         string    `json:"unit"`
	Calories     float64   `json:"calories_per_unit"`
	Protein      float64   `json:"protein_per_unit"`
	Carbohydrate float64   `json:"carbs_per_unit"`
	Fat          float64   `json:"fat_per_unit"`
	Fiber        float64   `json:"fiber_per_unit"`
}

type PlanStatus string

const (
	PlanStatusActive   PlanStatus = "active"
	PlanStatusArchived PlanStatus = "archived"
	PlanStatusDraft    PlanStatus = "draft"
)

type DietPlan struct {
	ID                 uuid.UUID
	ClientID           uuid.UUID
	NutritionistID     uuid.UUID
	Title              string
	StartDate          time.Time
	EndDate            time.Time
	Notes              string
	DailyWaterTargetML int
	Status             PlanStatus
	Days               []*DietPlanDay
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (p *DietPlan) IsActive() bool {
	return p.Status == PlanStatusActive
}

type DietPlanDay struct {
	ID            uuid.UUID
	PlanID        uuid.UUID
	DayNumber     int
	Meals         []*DietMeal
	TotalRange    *NutritionalRange
	Exercises     []*ExerciseRecommendation
	Prescriptions []*PrescribedMedication
	CreatedAt     time.Time
}

type DietMeal struct {
	ID            uuid.UUID
	DayID         uuid.UUID
	Title         string
	ScheduledTime string // "HH:MM"
	DisplayOrder  int
	Options       []*MealOption
	TotalRange    *NutritionalRange
	CreatedAt     time.Time
}

type MealOption struct {
	ID           uuid.UUID
	MealID       uuid.UUID
	OptionNumber int
	Items        []*MealOptionItem
	Totals       *NutritionalSummary
	CreatedAt    time.Time
}

type MealOptionItem struct {
	ID        uuid.UUID
	OptionID  uuid.UUID
	FoodID    uuid.UUID
	Quantity  float64
	Unit      string
	Notes     string
	Food      *FoodSnapshot
	Computed  *NutritionalSummary
	CreatedAt time.Time
}

// ExerciseRecommendation is a day-level exercise suggestion.
type ExerciseRecommendation struct {
	ID                   uuid.UUID
	DayID                uuid.UUID
	ExerciseName         string
	DurationMinutes      int
	Description          string
	CaloriesBurnEstimate int
	CreatedAt            time.Time
}

// MedicationSnapshot carries medication reference data for a prescription.
type MedicationSnapshot struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Unit string    `json:"unit"`
}

// PrescribedMedication is a day-level medication prescription.
type PrescribedMedication struct {
	ID           uuid.UUID
	DayID        uuid.UUID
	MedicationID uuid.UUID
	Medication   *MedicationSnapshot
	Dosage       string
	Frequency    string
	Times        []string   // "HH:MM" strings
	Instructions string
	StartDate    *time.Time
	EndDate      *time.Time
	CreatedAt    time.Time
}
