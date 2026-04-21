package entity

import (
	"time"

	"github.com/google/uuid"
)

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
	ID        uuid.UUID
	PlanID    uuid.UUID
	DayNumber int
	Meals     []*DietMeal
	CreatedAt time.Time
}

type DietMeal struct {
	ID            uuid.UUID
	DayID         uuid.UUID
	Title         string
	ScheduledTime string // "HH:MM"
	DisplayOrder  int
	Options       []*MealOption
	CreatedAt     time.Time
}

type MealOption struct {
	ID           uuid.UUID
	MealID       uuid.UUID
	OptionNumber int
	Items        []*MealOptionItem
	CreatedAt    time.Time
}

type MealOptionItem struct {
	ID        uuid.UUID
	OptionID  uuid.UUID
	FoodID    uuid.UUID
	Quantity  float64
	Unit      string
	Notes     string
	CreatedAt time.Time
}
