package dietplan

import (
	"context"
	"time"

	"github.com/google/uuid"
	dietRepo "github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/repository"
	"github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/shared"
	userRepo "github.com/ranjbar-dev/nutritrack/internal/domain/user/repository"
)

// CreatePlanRequest contains fields required to create a new diet plan.
type CreatePlanRequest struct {
	ClientID           uuid.UUID
	NutritionistID     uuid.UUID
	Title              string
	StartDate          time.Time
	EndDate            time.Time
	Notes              string
	DailyWaterTargetML int
}

// AddDayRequest contains fields required to add a day to a diet plan.
type AddDayRequest struct {
	PlanID     uuid.UUID
	DayNumber  int
	CallerID   uuid.UUID
	CallerRole string
}

// AddMealRequest contains fields required to add a meal to a diet plan day.
type AddMealRequest struct {
	DayID         uuid.UUID
	Title         string
	ScheduledTime string
	DisplayOrder  int
	CallerID      uuid.UUID
	CallerRole    string
}

// AddOptionRequest contains fields required to add an option to a meal.
type AddOptionRequest struct {
	MealID       uuid.UUID
	OptionNumber int
	CallerID     uuid.UUID
	CallerRole   string
}

// DietPlanService provides business logic for diet plans.
type DietPlanService struct {
	planRepo dietRepo.DietPlanRepository
	userRepo userRepo.UserRepository
}

// NewDietPlanService creates a new DietPlanService.
func NewDietPlanService(planRepo dietRepo.DietPlanRepository, userRepo userRepo.UserRepository) *DietPlanService {
	return &DietPlanService{planRepo: planRepo, userRepo: userRepo}
}

// CreatePlan creates a new diet plan for a client.
// It automatically archives any existing active plan for the client.
func (s *DietPlanService) CreatePlan(ctx context.Context, req CreatePlanRequest) (*entity.DietPlan, error) {
	// Verify client exists.
	client, err := s.userRepo.FindByID(ctx, req.ClientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, shared.ErrUserNotFound
	}

	plan := &entity.DietPlan{
		ClientID:           req.ClientID,
		NutritionistID:     req.NutritionistID,
		Title:              req.Title,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		Notes:              req.Notes,
		DailyWaterTargetML: req.DailyWaterTargetML,
		Status:             entity.PlanStatusActive,
	}

	if err := s.planRepo.CreateWithArchive(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// GetFullPlan retrieves a diet plan with its full day/meal/option tree.
func (s *DietPlanService) GetFullPlan(ctx context.Context, planID, callerID uuid.UUID, callerRole string) (*entity.DietPlan, error) {
	plan, err := s.planRepo.FindByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, shared.ErrPlanNotFound
	}

	// Authorization: caller must be the nutritionist, the client, or a superadmin.
	if callerRole != "superadmin" && callerID != plan.NutritionistID && callerID != plan.ClientID {
		return nil, shared.ErrForbidden
	}

	// Load full tree.
	days, err := s.planRepo.ListDays(ctx, plan.ID)
	if err != nil {
		return nil, err
	}

	for _, day := range days {
		meals, err := s.planRepo.ListMeals(ctx, day.ID)
		if err != nil {
			return nil, err
		}
		for _, meal := range meals {
			options, err := s.planRepo.ListOptions(ctx, meal.ID)
			if err != nil {
				return nil, err
			}
			for _, option := range options {
				items, err := s.planRepo.ListItemsWithFood(ctx, option.ID)
				if err != nil {
					return nil, err
				}
				option.Items = items
				option.Totals = computeOptionTotals(items)
			}
			meal.Options = options
			meal.TotalRange = computeMealRange(meal.Options)
		}
		day.Meals = meals
		day.TotalRange = computeDayRange(day.Meals)

		exercises, err := s.planRepo.ListExercises(ctx, day.ID)
		if err != nil {
			return nil, err
		}
		day.Exercises = exercises

		prescriptions, err := s.planRepo.ListPrescriptionsWithMedication(ctx, day.ID)
		if err != nil {
			return nil, err
		}
		day.Prescriptions = prescriptions
	}
	plan.Days = days

	return plan, nil
}

// AddDay adds a new day to a diet plan.
func (s *DietPlanService) AddDay(ctx context.Context, req AddDayRequest) (*entity.DietPlanDay, error) {
	plan, err := s.planRepo.FindByID(ctx, req.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, shared.ErrPlanNotFound
	}

	// Only nutritionist owner or superadmin may add days.
	if req.CallerRole != "superadmin" && req.CallerID != plan.NutritionistID {
		return nil, shared.ErrForbidden
	}

	day := &entity.DietPlanDay{
		PlanID:    req.PlanID,
		DayNumber: req.DayNumber,
	}

	if err := s.planRepo.AddDay(ctx, day); err != nil {
		return nil, err
	}

	return day, nil
}

// AddMeal adds a new meal to a diet plan day.
func (s *DietPlanService) AddMeal(ctx context.Context, req AddMealRequest) (*entity.DietMeal, error) {
	day, err := s.planRepo.FindDayByID(ctx, req.DayID)
	if err != nil {
		return nil, err
	}
	if day == nil {
		return nil, shared.ErrPlanNotFound
	}

	// Fetch plan for ownership check.
	plan, err := s.planRepo.FindByID(ctx, day.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, shared.ErrPlanNotFound
	}

	if req.CallerRole != "superadmin" && req.CallerID != plan.NutritionistID {
		return nil, shared.ErrForbidden
	}

	meal := &entity.DietMeal{
		DayID:         req.DayID,
		Title:         req.Title,
		ScheduledTime: req.ScheduledTime,
		DisplayOrder:  req.DisplayOrder,
	}

	if err := s.planRepo.AddMeal(ctx, meal); err != nil {
		return nil, err
	}

	return meal, nil
}

// AddOption adds a new option to a diet meal.
func (s *DietPlanService) AddOption(ctx context.Context, req AddOptionRequest) (*entity.MealOption, error) {
	meal, err := s.planRepo.FindMealByID(ctx, req.MealID)
	if err != nil {
		return nil, err
	}
	if meal == nil {
		return nil, shared.ErrPlanNotFound
	}

	// Fetch day and plan for ownership check.
	day, err := s.planRepo.FindDayByID(ctx, meal.DayID)
	if err != nil {
		return nil, err
	}
	if day == nil {
		return nil, shared.ErrPlanNotFound
	}

	plan, err := s.planRepo.FindByID(ctx, day.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, shared.ErrPlanNotFound
	}

	if req.CallerRole != "superadmin" && req.CallerID != plan.NutritionistID {
		return nil, shared.ErrForbidden
	}

	option := &entity.MealOption{
		MealID:       req.MealID,
		OptionNumber: req.OptionNumber,
	}

	if err := s.planRepo.AddOption(ctx, option); err != nil {
		return nil, err
	}

	return option, nil
}

// ListClientPlans returns a paginated list of diet plans for a client.
func (s *DietPlanService) ListClientPlans(ctx context.Context, clientID, callerID uuid.UUID, callerRole string, limit, offset int32) ([]*entity.DietPlan, int64, error) {
	// Authorization: caller must be the client, their nutritionist, or superadmin.
	if callerRole != "superadmin" && callerID != clientID {
		// Allow nutritionist access (nutritionist can see their client's plans).
		if callerRole != "nutritionist" {
			return nil, 0, shared.ErrForbidden
		}
	}

	plans, err := s.planRepo.ListByClientID(ctx, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.planRepo.CountByClientID(ctx, clientID)
	if err != nil {
		return nil, 0, err
	}

	return plans, total, nil
}

// DeletePlan deletes a diet plan by ID.
func (s *DietPlanService) DeletePlan(ctx context.Context, planID, callerID uuid.UUID, callerRole string) error {
	plan, err := s.planRepo.FindByID(ctx, planID)
	if err != nil {
		return err
	}
	if plan == nil {
		return shared.ErrPlanNotFound
	}

	if callerRole != "superadmin" && callerID != plan.NutritionistID {
		return shared.ErrForbidden
	}

	return s.planRepo.Delete(ctx, planID)
}

// AddItemRequest contains fields required to add a food item to a meal option.
type AddItemRequest struct {
	OptionID   uuid.UUID
	FoodID     uuid.UUID
	Quantity   float64
	Unit       string
	Notes      string
	CallerID   uuid.UUID
	CallerRole string
}

// AddItem adds a food item to a meal option.
func (s *DietPlanService) AddItem(ctx context.Context, req AddItemRequest) (*entity.MealOptionItem, error) {
	// Walk up: option → meal → day → plan for ownership check
	option, err := s.planRepo.FindOptionByID(ctx, req.OptionID)
	if err != nil {
		return nil, err
	}
	if option == nil {
		return nil, shared.ErrPlanNotFound
	}

	meal, err := s.planRepo.FindMealByID(ctx, option.MealID)
	if err != nil {
		return nil, err
	}
	if meal == nil {
		return nil, shared.ErrPlanNotFound
	}

	day, err := s.planRepo.FindDayByID(ctx, meal.DayID)
	if err != nil {
		return nil, err
	}
	if day == nil {
		return nil, shared.ErrPlanNotFound
	}

	plan, err := s.planRepo.FindByID(ctx, day.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, shared.ErrPlanNotFound
	}

	if req.CallerRole != "superadmin" && req.CallerID != plan.NutritionistID {
		return nil, shared.ErrForbidden
	}

	item := &entity.MealOptionItem{
		OptionID: req.OptionID,
		FoodID:   req.FoodID,
		Quantity: req.Quantity,
		Unit:     req.Unit,
		Notes:    req.Notes,
	}
	if err := s.planRepo.AddItem(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

// RemoveItem removes a food item from a meal option.
func (s *DietPlanService) RemoveItem(ctx context.Context, itemID, callerID uuid.UUID, callerRole string) error {	item, err := s.planRepo.FindItemByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return shared.ErrPlanNotFound
	}

	option, err := s.planRepo.FindOptionByID(ctx, item.OptionID)
	if err != nil {
		return err
	}
	if option == nil {
		return shared.ErrPlanNotFound
	}

	meal, err := s.planRepo.FindMealByID(ctx, option.MealID)
	if err != nil {
		return err
	}
	if meal == nil {
		return shared.ErrPlanNotFound
	}

	day, err := s.planRepo.FindDayByID(ctx, meal.DayID)
	if err != nil {
		return err
	}
	if day == nil {
		return shared.ErrPlanNotFound
	}

	plan, err := s.planRepo.FindByID(ctx, day.PlanID)
	if err != nil {
		return err
	}
	if plan == nil {
		return shared.ErrPlanNotFound
	}

	if callerRole != "superadmin" && callerID != plan.NutritionistID {
		return shared.ErrForbidden
	}

	return s.planRepo.DeleteItem(ctx, itemID)
}

// AddExerciseRequest contains fields required to add an exercise recommendation to a diet plan day.
type AddExerciseRequest struct {
	DayID                uuid.UUID
	ExerciseName         string
	DurationMinutes      int
	Description          string
	CaloriesBurnEstimate int
	CallerID             uuid.UUID
	CallerRole           string
}

// AddPrescriptionRequest contains fields required to add a medication prescription to a diet plan day.
type AddPrescriptionRequest struct {
	DayID        uuid.UUID
	MedicationID uuid.UUID
	Dosage       string
	Frequency    string
	Times        []string
	Instructions string
	StartDate    *time.Time
	EndDate      *time.Time
	CallerID     uuid.UUID
	CallerRole   string
}

// AddExercise adds an exercise recommendation to a diet plan day.
func (s *DietPlanService) AddExercise(ctx context.Context, req AddExerciseRequest) (*entity.ExerciseRecommendation, error) {
	day, err := s.planRepo.FindDayByID(ctx, req.DayID)
	if err != nil {
		return nil, err
	}
	if day == nil {
		return nil, shared.ErrPlanNotFound
	}

	plan, err := s.planRepo.FindByID(ctx, day.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, shared.ErrPlanNotFound
	}

	if req.CallerRole != "superadmin" && req.CallerID != plan.NutritionistID {
		return nil, shared.ErrForbidden
	}

	ex := &entity.ExerciseRecommendation{
		DayID:                req.DayID,
		ExerciseName:         req.ExerciseName,
		DurationMinutes:      req.DurationMinutes,
		Description:          req.Description,
		CaloriesBurnEstimate: req.CaloriesBurnEstimate,
	}

	if err := s.planRepo.AddExercise(ctx, ex); err != nil {
		return nil, err
	}
	return ex, nil
}

// RemoveExercise removes an exercise recommendation.
func (s *DietPlanService) RemoveExercise(ctx context.Context, exerciseID, callerID uuid.UUID, callerRole string) error {
	ex, err := s.planRepo.FindExerciseByID(ctx, exerciseID)
	if err != nil {
		return err
	}
	if ex == nil {
		return shared.ErrPlanNotFound
	}

	day, err := s.planRepo.FindDayByID(ctx, ex.DayID)
	if err != nil {
		return err
	}
	if day == nil {
		return shared.ErrPlanNotFound
	}

	plan, err := s.planRepo.FindByID(ctx, day.PlanID)
	if err != nil {
		return err
	}
	if plan == nil {
		return shared.ErrPlanNotFound
	}

	if callerRole != "superadmin" && callerID != plan.NutritionistID {
		return shared.ErrForbidden
	}

	return s.planRepo.DeleteExercise(ctx, exerciseID)
}

// AddPrescription adds a medication prescription to a diet plan day.
func (s *DietPlanService) AddPrescription(ctx context.Context, req AddPrescriptionRequest) (*entity.PrescribedMedication, error) {
	day, err := s.planRepo.FindDayByID(ctx, req.DayID)
	if err != nil {
		return nil, err
	}
	if day == nil {
		return nil, shared.ErrPlanNotFound
	}

	plan, err := s.planRepo.FindByID(ctx, day.PlanID)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, shared.ErrPlanNotFound
	}

	if req.CallerRole != "superadmin" && req.CallerID != plan.NutritionistID {
		return nil, shared.ErrForbidden
	}

	rx := &entity.PrescribedMedication{
		DayID:        req.DayID,
		MedicationID: req.MedicationID,
		Dosage:       req.Dosage,
		Frequency:    req.Frequency,
		Times:        req.Times,
		Instructions: req.Instructions,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
	}

	if err := s.planRepo.AddPrescription(ctx, rx); err != nil {
		return nil, err
	}
	return rx, nil
}

// RemovePrescription removes a medication prescription.
func (s *DietPlanService) RemovePrescription(ctx context.Context, prescriptionID, callerID uuid.UUID, callerRole string) error {
	rx, err := s.planRepo.FindPrescriptionByID(ctx, prescriptionID)
	if err != nil {
		return err
	}
	if rx == nil {
		return shared.ErrPlanNotFound
	}

	day, err := s.planRepo.FindDayByID(ctx, rx.DayID)
	if err != nil {
		return err
	}
	if day == nil {
		return shared.ErrPlanNotFound
	}

	plan, err := s.planRepo.FindByID(ctx, day.PlanID)
	if err != nil {
		return err
	}
	if plan == nil {
		return shared.ErrPlanNotFound
	}

	if callerRole != "superadmin" && callerID != plan.NutritionistID {
		return shared.ErrForbidden
	}

	return s.planRepo.DeletePrescription(ctx, prescriptionID)
}
