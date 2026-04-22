package entity

import (
	"time"

	"github.com/google/uuid"
)

// NutritionalSummary holds computed macro totals.
type NutritionalSummary struct {
	Calories float64
	Protein  float64
	Carbs    float64
	Fat      float64
	Fiber    float64
}

// NutritionalRange holds min/max nutritional totals across options.
type NutritionalRange struct {
	Min NutritionalSummary
	Max NutritionalSummary
}

// FoodSnapshot carries the food details needed for nutritional computation.
type FoodSnapshot struct {
	ID           uuid.UUID
	Name         string
	Unit         string
	Calories     float64
	Protein      float64
	Carbohydrate float64
	Fat          float64
	Fiber        float64
}

// MedicationSnapshot carries medication reference data for a prescription.
type MedicationSnapshot struct {
	ID   uuid.UUID
	Name string
	Unit string
}

type PlanStatus string

const (
	PlanStatusActive   PlanStatus = "active"
	PlanStatusArchived PlanStatus = "archived"
	PlanStatusDraft    PlanStatus = "draft"
)

// DietPlan is the root aggregate for a nutritionist's diet prescription.
type DietPlan struct {
	id                 uuid.UUID
	clientID           uuid.UUID
	nutritionistID     uuid.UUID
	title              string
	startDate          time.Time
	endDate            time.Time
	notes              string
	dailyWaterTargetML int
	status             PlanStatus
	days               []*DietPlanDay
	createdAt          time.Time
	updatedAt          time.Time
}

// NewDietPlan creates a new DietPlan aggregate.
func NewDietPlan(clientID, nutritionistID uuid.UUID, title string, startDate, endDate time.Time, notes string, dailyWaterTargetML int) *DietPlan {
	return &DietPlan{
		clientID:           clientID,
		nutritionistID:     nutritionistID,
		title:              title,
		startDate:          startDate,
		endDate:            endDate,
		notes:              notes,
		dailyWaterTargetML: dailyWaterTargetML,
		status:             PlanStatusActive,
		days:               []*DietPlanDay{},
	}
}

// ReconstituteDietPlan rebuilds a DietPlan from stored data.
func ReconstituteDietPlan(id, clientID, nutritionistID uuid.UUID, title string, startDate, endDate time.Time, notes string, dailyWaterTargetML int, status PlanStatus, days []*DietPlanDay, createdAt, updatedAt time.Time) *DietPlan {
	return &DietPlan{id: id, clientID: clientID, nutritionistID: nutritionistID, title: title, startDate: startDate, endDate: endDate, notes: notes, dailyWaterTargetML: dailyWaterTargetML, status: status, days: days, createdAt: createdAt, updatedAt: updatedAt}
}

func (p *DietPlan) ID() uuid.UUID             { return p.id }
func (p *DietPlan) ClientID() uuid.UUID       { return p.clientID }
func (p *DietPlan) NutritionistID() uuid.UUID { return p.nutritionistID }
func (p *DietPlan) Title() string             { return p.title }
func (p *DietPlan) StartDate() time.Time      { return p.startDate }
func (p *DietPlan) EndDate() time.Time        { return p.endDate }
func (p *DietPlan) Notes() string             { return p.notes }
func (p *DietPlan) DailyWaterTargetML() int   { return p.dailyWaterTargetML }
func (p *DietPlan) Status() PlanStatus        { return p.status }
func (p *DietPlan) Days() []*DietPlanDay      { return p.days }
func (p *DietPlan) CreatedAt() time.Time      { return p.createdAt }
func (p *DietPlan) UpdatedAt() time.Time      { return p.updatedAt }

func (p *DietPlan) SetID(id uuid.UUID)           { p.id = id }
func (p *DietPlan) SetStatus(s PlanStatus)       { p.status = s }
func (p *DietPlan) SetCreatedAt(t time.Time)     { p.createdAt = t }
func (p *DietPlan) SetUpdatedAt(t time.Time)     { p.updatedAt = t }
func (p *DietPlan) SetTitle(title string)        { p.title = title }
func (p *DietPlan) SetNotes(notes string)        { p.notes = notes }
func (p *DietPlan) SetDailyWaterTargetML(ml int) { p.dailyWaterTargetML = ml }
func (p *DietPlan) SetDays(days []*DietPlanDay)  { p.days = days }

func (p *DietPlan) IsActive() bool { return p.status == PlanStatusActive }

// DietPlanDay represents a single day within a diet plan.
type DietPlanDay struct {
	id            uuid.UUID
	planID        uuid.UUID
	dayNumber     int
	meals         []*DietMeal
	totalRange    *NutritionalRange
	exercises     []*ExerciseRecommendation
	prescriptions []*PrescribedMedication
	createdAt     time.Time
}

// NewDietPlanDay creates a new DietPlanDay.
func NewDietPlanDay(planID uuid.UUID, dayNumber int) *DietPlanDay {
	return &DietPlanDay{planID: planID, dayNumber: dayNumber, meals: []*DietMeal{}, exercises: []*ExerciseRecommendation{}, prescriptions: []*PrescribedMedication{}}
}

// ReconstituteDietPlanDay rebuilds a DietPlanDay from stored data.
func ReconstituteDietPlanDay(id, planID uuid.UUID, dayNumber int, meals []*DietMeal, totalRange *NutritionalRange, exercises []*ExerciseRecommendation, prescriptions []*PrescribedMedication, createdAt time.Time) *DietPlanDay {
	return &DietPlanDay{id: id, planID: planID, dayNumber: dayNumber, meals: meals, totalRange: totalRange, exercises: exercises, prescriptions: prescriptions, createdAt: createdAt}
}

func (d *DietPlanDay) ID() uuid.UUID                          { return d.id }
func (d *DietPlanDay) PlanID() uuid.UUID                      { return d.planID }
func (d *DietPlanDay) DayNumber() int                         { return d.dayNumber }
func (d *DietPlanDay) Meals() []*DietMeal                     { return d.meals }
func (d *DietPlanDay) TotalRange() *NutritionalRange          { return d.totalRange }
func (d *DietPlanDay) Exercises() []*ExerciseRecommendation   { return d.exercises }
func (d *DietPlanDay) Prescriptions() []*PrescribedMedication { return d.prescriptions }
func (d *DietPlanDay) CreatedAt() time.Time                   { return d.createdAt }

func (d *DietPlanDay) SetID(id uuid.UUID)                          { d.id = id }
func (d *DietPlanDay) SetCreatedAt(t time.Time)                    { d.createdAt = t }
func (d *DietPlanDay) SetMeals(meals []*DietMeal)                  { d.meals = meals }
func (d *DietPlanDay) SetTotalRange(r *NutritionalRange)           { d.totalRange = r }
func (d *DietPlanDay) SetExercises(ex []*ExerciseRecommendation)   { d.exercises = ex }
func (d *DietPlanDay) SetPrescriptions(rx []*PrescribedMedication) { d.prescriptions = rx }

// DietMeal represents a single meal within a diet plan day.
type DietMeal struct {
	id            uuid.UUID
	dayID         uuid.UUID
	title         string
	scheduledTime string
	displayOrder  int
	options       []*MealOption
	totalRange    *NutritionalRange
	createdAt     time.Time
}

// NewDietMeal creates a new DietMeal.
func NewDietMeal(dayID uuid.UUID, title, scheduledTime string, displayOrder int) *DietMeal {
	return &DietMeal{dayID: dayID, title: title, scheduledTime: scheduledTime, displayOrder: displayOrder, options: []*MealOption{}}
}

// ReconstituteDietMeal rebuilds a DietMeal from stored data.
func ReconstituteDietMeal(id, dayID uuid.UUID, title, scheduledTime string, displayOrder int, options []*MealOption, totalRange *NutritionalRange, createdAt time.Time) *DietMeal {
	return &DietMeal{id: id, dayID: dayID, title: title, scheduledTime: scheduledTime, displayOrder: displayOrder, options: options, totalRange: totalRange, createdAt: createdAt}
}

func (m *DietMeal) ID() uuid.UUID                 { return m.id }
func (m *DietMeal) DayID() uuid.UUID              { return m.dayID }
func (m *DietMeal) Title() string                 { return m.title }
func (m *DietMeal) ScheduledTime() string         { return m.scheduledTime }
func (m *DietMeal) DisplayOrder() int             { return m.displayOrder }
func (m *DietMeal) Options() []*MealOption        { return m.options }
func (m *DietMeal) TotalRange() *NutritionalRange { return m.totalRange }
func (m *DietMeal) CreatedAt() time.Time          { return m.createdAt }

func (m *DietMeal) SetID(id uuid.UUID)                { m.id = id }
func (m *DietMeal) SetCreatedAt(t time.Time)          { m.createdAt = t }
func (m *DietMeal) SetOptions(opts []*MealOption)     { m.options = opts }
func (m *DietMeal) SetTotalRange(r *NutritionalRange) { m.totalRange = r }

// MealOption represents one alternative option within a meal.
type MealOption struct {
	id           uuid.UUID
	mealID       uuid.UUID
	optionNumber int
	items        []*MealOptionItem
	totals       *NutritionalSummary
	createdAt    time.Time
}

// NewMealOption creates a new MealOption.
func NewMealOption(mealID uuid.UUID, optionNumber int) *MealOption {
	return &MealOption{mealID: mealID, optionNumber: optionNumber, items: []*MealOptionItem{}}
}

// ReconstituteMealOption rebuilds a MealOption from stored data.
func ReconstituteMealOption(id, mealID uuid.UUID, optionNumber int, items []*MealOptionItem, totals *NutritionalSummary, createdAt time.Time) *MealOption {
	return &MealOption{id: id, mealID: mealID, optionNumber: optionNumber, items: items, totals: totals, createdAt: createdAt}
}

func (o *MealOption) ID() uuid.UUID               { return o.id }
func (o *MealOption) MealID() uuid.UUID           { return o.mealID }
func (o *MealOption) OptionNumber() int           { return o.optionNumber }
func (o *MealOption) Items() []*MealOptionItem    { return o.items }
func (o *MealOption) Totals() *NutritionalSummary { return o.totals }
func (o *MealOption) CreatedAt() time.Time        { return o.createdAt }

func (o *MealOption) SetID(id uuid.UUID)               { o.id = id }
func (o *MealOption) SetCreatedAt(t time.Time)         { o.createdAt = t }
func (o *MealOption) SetItems(items []*MealOptionItem) { o.items = items }
func (o *MealOption) SetTotals(s *NutritionalSummary)  { o.totals = s }

// MealOptionItem is a single food item within a meal option.
type MealOptionItem struct {
	id        uuid.UUID
	optionID  uuid.UUID
	foodID    uuid.UUID
	quantity  float64
	unit      string
	notes     string
	food      *FoodSnapshot
	computed  *NutritionalSummary
	createdAt time.Time
}

// NewMealOptionItem creates a new MealOptionItem.
func NewMealOptionItem(optionID, foodID uuid.UUID, quantity float64, unit, notes string) *MealOptionItem {
	return &MealOptionItem{optionID: optionID, foodID: foodID, quantity: quantity, unit: unit, notes: notes}
}

// ReconstituteMealOptionItem rebuilds a MealOptionItem from stored data.
func ReconstituteMealOptionItem(id, optionID, foodID uuid.UUID, quantity float64, unit, notes string, food *FoodSnapshot, computed *NutritionalSummary, createdAt time.Time) *MealOptionItem {
	return &MealOptionItem{id: id, optionID: optionID, foodID: foodID, quantity: quantity, unit: unit, notes: notes, food: food, computed: computed, createdAt: createdAt}
}

func (i *MealOptionItem) ID() uuid.UUID                 { return i.id }
func (i *MealOptionItem) OptionID() uuid.UUID           { return i.optionID }
func (i *MealOptionItem) FoodID() uuid.UUID             { return i.foodID }
func (i *MealOptionItem) Quantity() float64             { return i.quantity }
func (i *MealOptionItem) Unit() string                  { return i.unit }
func (i *MealOptionItem) Notes() string                 { return i.notes }
func (i *MealOptionItem) Food() *FoodSnapshot           { return i.food }
func (i *MealOptionItem) Computed() *NutritionalSummary { return i.computed }
func (i *MealOptionItem) CreatedAt() time.Time          { return i.createdAt }

func (i *MealOptionItem) SetID(id uuid.UUID)                { i.id = id }
func (i *MealOptionItem) SetCreatedAt(t time.Time)          { i.createdAt = t }
func (i *MealOptionItem) SetFood(f *FoodSnapshot)           { i.food = f }
func (i *MealOptionItem) SetComputed(s *NutritionalSummary) { i.computed = s }

// ExerciseRecommendation is a day-level exercise suggestion.
type ExerciseRecommendation struct {
	id                   uuid.UUID
	dayID                uuid.UUID
	exerciseName         string
	durationMinutes      int
	description          string
	caloriesBurnEstimate int
	createdAt            time.Time
}

// NewExerciseRecommendation creates a new ExerciseRecommendation.
func NewExerciseRecommendation(dayID uuid.UUID, exerciseName string, durationMinutes int, description string, caloriesBurnEstimate int) *ExerciseRecommendation {
	return &ExerciseRecommendation{dayID: dayID, exerciseName: exerciseName, durationMinutes: durationMinutes, description: description, caloriesBurnEstimate: caloriesBurnEstimate}
}

// ReconstituteExerciseRecommendation rebuilds an ExerciseRecommendation from stored data.
func ReconstituteExerciseRecommendation(id, dayID uuid.UUID, exerciseName string, durationMinutes int, description string, caloriesBurnEstimate int, createdAt time.Time) *ExerciseRecommendation {
	return &ExerciseRecommendation{id: id, dayID: dayID, exerciseName: exerciseName, durationMinutes: durationMinutes, description: description, caloriesBurnEstimate: caloriesBurnEstimate, createdAt: createdAt}
}

func (e *ExerciseRecommendation) ID() uuid.UUID             { return e.id }
func (e *ExerciseRecommendation) DayID() uuid.UUID          { return e.dayID }
func (e *ExerciseRecommendation) ExerciseName() string      { return e.exerciseName }
func (e *ExerciseRecommendation) DurationMinutes() int      { return e.durationMinutes }
func (e *ExerciseRecommendation) Description() string       { return e.description }
func (e *ExerciseRecommendation) CaloriesBurnEstimate() int { return e.caloriesBurnEstimate }
func (e *ExerciseRecommendation) CreatedAt() time.Time      { return e.createdAt }

func (e *ExerciseRecommendation) SetID(id uuid.UUID)       { e.id = id }
func (e *ExerciseRecommendation) SetCreatedAt(t time.Time) { e.createdAt = t }

// PrescribedMedication is a day-level medication prescription.
type PrescribedMedication struct {
	id           uuid.UUID
	dayID        uuid.UUID
	medicationID uuid.UUID
	medication   *MedicationSnapshot
	dosage       string
	frequency    string
	times        []string
	instructions string
	startDate    *time.Time
	endDate      *time.Time
	createdAt    time.Time
}

// NewPrescribedMedication creates a new PrescribedMedication.
func NewPrescribedMedication(dayID, medicationID uuid.UUID, dosage, frequency string, times []string, instructions string, startDate, endDate *time.Time) *PrescribedMedication {
	return &PrescribedMedication{dayID: dayID, medicationID: medicationID, dosage: dosage, frequency: frequency, times: times, instructions: instructions, startDate: startDate, endDate: endDate}
}

// ReconstitutePrescribedMedication rebuilds a PrescribedMedication from stored data.
func ReconstitutePrescribedMedication(id, dayID, medicationID uuid.UUID, medication *MedicationSnapshot, dosage, frequency string, times []string, instructions string, startDate, endDate *time.Time, createdAt time.Time) *PrescribedMedication {
	return &PrescribedMedication{id: id, dayID: dayID, medicationID: medicationID, medication: medication, dosage: dosage, frequency: frequency, times: times, instructions: instructions, startDate: startDate, endDate: endDate, createdAt: createdAt}
}

func (r *PrescribedMedication) ID() uuid.UUID                   { return r.id }
func (r *PrescribedMedication) DayID() uuid.UUID                { return r.dayID }
func (r *PrescribedMedication) MedicationID() uuid.UUID         { return r.medicationID }
func (r *PrescribedMedication) Medication() *MedicationSnapshot { return r.medication }
func (r *PrescribedMedication) Dosage() string                  { return r.dosage }
func (r *PrescribedMedication) Frequency() string               { return r.frequency }
func (r *PrescribedMedication) Times() []string                 { return r.times }
func (r *PrescribedMedication) Instructions() string            { return r.instructions }
func (r *PrescribedMedication) StartDate() *time.Time           { return r.startDate }
func (r *PrescribedMedication) EndDate() *time.Time             { return r.endDate }
func (r *PrescribedMedication) CreatedAt() time.Time            { return r.createdAt }

func (r *PrescribedMedication) SetID(id uuid.UUID)                  { r.id = id }
func (r *PrescribedMedication) SetCreatedAt(t time.Time)            { r.createdAt = t }
func (r *PrescribedMedication) SetMedication(m *MedicationSnapshot) { r.medication = m }

// --- Domain computation methods ---

// ComputeTotals sums item computed nutritional values for this option.
// This is domain logic and must not be duplicated in the application layer.
func (o *MealOption) ComputeTotals() *NutritionalSummary {
	s := &NutritionalSummary{}
	for _, it := range o.items {
		if it.Computed() != nil {
			s.Calories += it.Computed().Calories
			s.Protein += it.Computed().Protein
			s.Carbs += it.Computed().Carbs
			s.Fat += it.Computed().Fat
			s.Fiber += it.Computed().Fiber
		}
	}
	return s
}

// ComputeTotalRange derives min/max NutritionalSummary across this meal's options.
// This is domain logic and must not be duplicated in the application layer.
func (m *DietMeal) ComputeTotalRange() *NutritionalRange {
	if len(m.options) == 0 {
		return nil
	}
	first := true
	r := &NutritionalRange{}
	for _, opt := range m.options {
		if opt.Totals() == nil {
			continue
		}
		t := opt.Totals()
		if first {
			r.Min = *t
			r.Max = *t
			first = false
			continue
		}
		if t.Calories < r.Min.Calories {
			r.Min.Calories = t.Calories
		}
		if t.Calories > r.Max.Calories {
			r.Max.Calories = t.Calories
		}
		if t.Protein < r.Min.Protein {
			r.Min.Protein = t.Protein
		}
		if t.Protein > r.Max.Protein {
			r.Max.Protein = t.Protein
		}
		if t.Carbs < r.Min.Carbs {
			r.Min.Carbs = t.Carbs
		}
		if t.Carbs > r.Max.Carbs {
			r.Max.Carbs = t.Carbs
		}
		if t.Fat < r.Min.Fat {
			r.Min.Fat = t.Fat
		}
		if t.Fat > r.Max.Fat {
			r.Max.Fat = t.Fat
		}
		if t.Fiber < r.Min.Fiber {
			r.Min.Fiber = t.Fiber
		}
		if t.Fiber > r.Max.Fiber {
			r.Max.Fiber = t.Fiber
		}
	}
	if first {
		return nil
	}
	return r
}

// ComputeTotalRange sums meal min/max ranges for this day.
// This is domain logic and must not be duplicated in the application layer.
func (d *DietPlanDay) ComputeTotalRange() *NutritionalRange {
	day := &NutritionalRange{}
	for _, meal := range d.meals {
		if meal.TotalRange() == nil {
			continue
		}
		day.Min.Calories += meal.TotalRange().Min.Calories
		day.Min.Protein += meal.TotalRange().Min.Protein
		day.Min.Carbs += meal.TotalRange().Min.Carbs
		day.Min.Fat += meal.TotalRange().Min.Fat
		day.Min.Fiber += meal.TotalRange().Min.Fiber
		day.Max.Calories += meal.TotalRange().Max.Calories
		day.Max.Protein += meal.TotalRange().Max.Protein
		day.Max.Carbs += meal.TotalRange().Max.Carbs
		day.Max.Fat += meal.TotalRange().Max.Fat
		day.Max.Fiber += meal.TotalRange().Max.Fiber
	}
	return day
}
