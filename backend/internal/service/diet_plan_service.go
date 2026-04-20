package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

// ─── Sentinel errors ─────────────────────────────────────────────────────────

var (
	ErrPlanNotFound            = errors.New("برنامه غذایی یافت نشد")
	ErrPlanUnauthorized        = errors.New("شما مجوز دسترسی به این برنامه را ندارید")
	ErrPlanNotDraft            = errors.New("فقط برنامه‌های پیش‌نویس قابل ویرایش هستند")
	ErrPlanIncomplete          = errors.New("برنامه ناقص است — حداقل یک روز با یک وعده و یک گزینه الزامی است")
	ErrPlanAlreadyActive       = errors.New("این برنامه قبلاً فعال شده است")
	ErrPlanInvalidDateRange    = errors.New("تاریخ پایان باید بعد از تاریخ شروع باشد")
	ErrDayNotFound             = errors.New("روز یافت نشد")
	ErrMealNotFound            = errors.New("وعده یافت نشد")
	ErrOptionNotFound          = errors.New("گزینه یافت نشد")
	ErrItemNotFound            = errors.New("آیتم یافت نشد")
	ErrExerciseNotFound        = errors.New("تمرین یافت نشد")
	ErrMedicationPrescNotFound = errors.New("نسخه دارویی یافت نشد")
)

// ─── Service ─────────────────────────────────────────────────────────────────

// DietPlanService handles diet plan business logic.
type DietPlanService struct {
	planRepo repository.DietPlanRepository
	logger   zerolog.Logger
}

// NewDietPlanService creates a new DietPlanService.
func NewDietPlanService(planRepo repository.DietPlanRepository, logger zerolog.Logger) *DietPlanService {
	return &DietPlanService{planRepo: planRepo, logger: logger}
}

// ─── Plan-level CRUD ──────────────────────────────────────────────────────────

// CreatePlan creates a new diet plan in draft status.
func (s *DietPlanService) CreatePlan(ctx context.Context, nutritionistID uuid.UUID, req dto.CreateDietPlanRequest) (*dto.DietPlanSummaryResponse, error) {
	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("invalid client_id: %w", err)
	}

	startDate, err := parseDateToDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}
	endDate, err := parseDateToDate(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}

	if startDate.Valid && endDate.Valid && endDate.Time.Before(startDate.Time) {
		return nil, ErrPlanInvalidDateRange
	}

	plan, err := s.planRepo.CreatePlan(ctx, sqlc.CreateDietPlanParams{
		ClientID:           pgtype.UUID{Bytes: clientID, Valid: true},
		NutritionistID:     pgtype.UUID{Bytes: nutritionistID, Valid: true},
		StartDate:          startDate,
		EndDate:            endDate,
		Notes:              optionalText(req.Notes),
		DailyWaterTargetMl: optionalInt4(req.DailyWaterTargetMl),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("nutritionist_id", nutritionistID.String()).Msg("failed to create diet plan")
		return nil, fmt.Errorf("create diet plan: %w", err)
	}

	resp := planToSummary(plan, 0)
	return &resp, nil
}

// GetPlanAggregate returns the full nested diet plan aggregate.
func (s *DietPlanService) GetPlanAggregate(ctx context.Context, planID, nutritionistID uuid.UUID) (*dto.DietPlanResponse, error) {
	// Verify ownership — GetPlanByID returns pgx.ErrNoRows if not found or wrong nutritionist.
	_, err := s.planRepo.GetPlanByID(ctx, planID, nutritionistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("get plan by id: %w", err)
	}

	result, err := s.planRepo.GetFullPlanAggregate(ctx, planID)
	if err != nil {
		s.logger.Error().Err(err).Str("plan_id", planID.String()).Msg("failed to get plan aggregate")
		return nil, fmt.Errorf("get plan aggregate: %w", err)
	}
	return result, nil
}

// GetPlanAggregateForClient returns a non-draft plan aggregate owned by the authenticated client.
func (s *DietPlanService) GetPlanAggregateForClient(ctx context.Context, planID, clientID uuid.UUID) (*dto.DietPlanResponse, error) {
	_, err := s.planRepo.GetPlanByIDForClient(ctx, planID, clientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("get client plan by id: %w", err)
	}

	result, err := s.planRepo.GetFullPlanAggregate(ctx, planID)
	if err != nil {
		s.logger.Error().Err(err).Str("plan_id", planID.String()).Msg("failed to get client plan aggregate")
		return nil, fmt.Errorf("get client plan aggregate: %w", err)
	}
	return result, nil
}

// ListClientPlans returns a paginated list of plans for a client.
func (s *DietPlanService) ListClientPlans(ctx context.Context, clientID, nutritionistID uuid.UUID, page, limit int) (*dto.DietPlanListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := int32((page - 1) * limit)

	plans, err := s.planRepo.ListClientPlans(ctx, sqlc.ListClientPlansParams{
		ClientID:       pgtype.UUID{Bytes: clientID, Valid: true},
		NutritionistID: pgtype.UUID{Bytes: nutritionistID, Valid: true},
		OffsetVal:      offset,
		LimitVal:       int32(limit),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("client_id", clientID.String()).Msg("failed to list client plans")
		return nil, fmt.Errorf("list client plans: %w", err)
	}

	total, err := s.planRepo.CountClientPlans(ctx, sqlc.CountClientPlansParams{
		ClientID:       pgtype.UUID{Bytes: clientID, Valid: true},
		NutritionistID: pgtype.UUID{Bytes: nutritionistID, Valid: true},
	})
	if err != nil {
		s.logger.Error().Err(err).Str("client_id", clientID.String()).Msg("failed to count client plans")
		return nil, fmt.Errorf("count client plans: %w", err)
	}

	data := make([]dto.DietPlanSummaryResponse, 0, len(plans))
	for _, p := range plans {
		planUUID := uuid.UUID(p.ID.Bytes)
		dayCount, err := s.planRepo.CountPlanDays(ctx, planUUID)
		if err != nil {
			s.logger.Error().Err(err).Str("plan_id", planUUID.String()).Msg("failed to count plan days")
			return nil, fmt.Errorf("count plan days: %w", err)
		}
		data = append(data, planToSummary(&p, dayCount))
	}

	return &dto.DietPlanListResponse{
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasMore: int64(page*limit) < total,
	}, nil
}

// UpdatePlanHeader updates the plan's header fields (dates, notes, water target).
func (s *DietPlanService) UpdatePlanHeader(ctx context.Context, planID, nutritionistID uuid.UUID, req dto.UpdateDietPlanRequest) (*dto.DietPlanSummaryResponse, error) {
	plan, err := s.planRepo.GetPlanByID(ctx, planID, nutritionistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("get plan by id: %w", err)
	}
	if plan.Status != sqlc.DietPlanStatusDraft {
		return nil, ErrPlanNotDraft
	}

	startDate, err := parseDateToDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}
	endDate, err := parseDateToDate(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}
	if startDate.Valid && endDate.Valid && endDate.Time.Before(startDate.Time) {
		return nil, ErrPlanInvalidDateRange
	}

	updated, err := s.planRepo.UpdatePlanHeader(ctx, sqlc.UpdateDietPlanHeaderParams{
		ID:                 pgtype.UUID{Bytes: planID, Valid: true},
		NutritionistID:     pgtype.UUID{Bytes: nutritionistID, Valid: true},
		StartDate:          startDate,
		EndDate:            endDate,
		Notes:              optionalText(req.Notes),
		DailyWaterTargetMl: optionalInt4(req.DailyWaterTargetMl),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("plan_id", planID.String()).Msg("failed to update plan header")
		return nil, fmt.Errorf("update plan header: %w", err)
	}

	dayCount, err := s.planRepo.CountPlanDays(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("count plan days: %w", err)
	}

	resp := planToSummary(updated, dayCount)
	return &resp, nil
}

// ActivatePlan validates completeness, archives the previous active plan, then activates.
func (s *DietPlanService) ActivatePlan(ctx context.Context, planID, nutritionistID uuid.UUID) error {
	plan, err := s.planRepo.GetPlanByID(ctx, planID, nutritionistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPlanNotFound
		}
		return fmt.Errorf("get plan by id: %w", err)
	}
	if plan.Status == sqlc.DietPlanStatusActive {
		return ErrPlanAlreadyActive
	}
	if plan.Status != sqlc.DietPlanStatusDraft {
		return ErrPlanNotDraft
	}

	if err := s.validatePlanComplete(ctx, planID); err != nil {
		return err
	}

	clientID := uuid.UUID(plan.ClientID.Bytes)

	// Archive any existing active plan for this client (D-02 constraint).
	if err := s.planRepo.ArchivePreviousActivePlan(ctx, clientID, planID); err != nil {
		s.logger.Error().Err(err).Str("client_id", clientID.String()).Msg("failed to archive previous active plan")
		return fmt.Errorf("archive previous active plan: %w", err)
	}

	if err := s.planRepo.ActivatePlan(ctx, planID, nutritionistID); err != nil {
		// Unique constraint violation = race condition backstop (T-03-03-C)
		if isUniqueViolation(err) {
			return ErrPlanAlreadyActive
		}
		s.logger.Error().Err(err).Str("plan_id", planID.String()).Msg("failed to activate plan")
		return fmt.Errorf("activate plan: %w", err)
	}

	s.logger.Info().
		Str("plan_id", planID.String()).
		Str("client_id", clientID.String()).
		Str("nutritionist_id", nutritionistID.String()).
		Msg("diet plan activated")

	return nil
}

// isUniqueViolation checks if a pgx error is a PostgreSQL unique constraint violation (23505).
func isUniqueViolation(err error) bool {
	return err != nil && (fmt.Sprintf("%v", err) == "ERROR: duplicate key value violates unique constraint" ||
		contains(err.Error(), "23505") || contains(err.Error(), "unique constraint"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsRune(s, substr))
}

func containsRune(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// validatePlanComplete validates that a plan has at least one day, each day has meals, each meal has options, each option has items.
func (s *DietPlanService) validatePlanComplete(ctx context.Context, planID uuid.UUID) error {
	plan, err := s.planRepo.GetFullPlanAggregate(ctx, planID)
	if err != nil {
		return fmt.Errorf("get plan aggregate for validation: %w", err)
	}

	if len(plan.Days) == 0 {
		return ErrPlanIncomplete
	}
	for _, day := range plan.Days {
		if len(day.Meals) == 0 {
			return ErrPlanIncomplete
		}
		for _, meal := range day.Meals {
			if len(meal.Options) == 0 {
				return ErrPlanIncomplete
			}
			for _, opt := range meal.Options {
				if len(opt.Items) == 0 {
					return ErrPlanIncomplete
				}
			}
		}
	}
	return nil
}

// DeletePlan hard-deletes a draft plan. Returns ErrPlanNotDraft for non-draft plans.
func (s *DietPlanService) DeletePlan(ctx context.Context, planID, nutritionistID uuid.UUID) error {
	plan, err := s.planRepo.GetPlanByID(ctx, planID, nutritionistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPlanNotFound
		}
		return fmt.Errorf("get plan by id: %w", err)
	}
	if plan.Status != sqlc.DietPlanStatusDraft {
		return ErrPlanNotDraft
	}
	if err := s.planRepo.DeletePlan(ctx, planID, nutritionistID); err != nil {
		s.logger.Error().Err(err).Str("plan_id", planID.String()).Msg("failed to delete plan")
		return fmt.Errorf("delete plan: %w", err)
	}
	return nil
}

// GetActivePlanForClient returns the active plan aggregate for the authenticated client.
func (s *DietPlanService) GetActivePlanForClient(ctx context.Context, clientID uuid.UUID) (*dto.DietPlanResponse, error) {
	result, err := s.planRepo.GetActivePlanForClient(ctx, clientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPlanNotFound
		}
		s.logger.Error().Err(err).Str("client_id", clientID.String()).Msg("failed to get active plan for client")
		return nil, fmt.Errorf("get active plan: %w", err)
	}
	return result, nil
}

// ListMyPlans returns all non-draft plans for the authenticated client.
func (s *DietPlanService) ListMyPlans(ctx context.Context, clientID uuid.UUID, page, limit int) (*dto.DietPlanListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := int32((page - 1) * limit)

	plans, err := s.planRepo.ListMyPlans(ctx, sqlc.ListMyPlansParams{
		ClientID:  pgtype.UUID{Bytes: clientID, Valid: true},
		OffsetVal: offset,
		LimitVal:  int32(limit),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("client_id", clientID.String()).Msg("failed to list client plan history")
		return nil, fmt.Errorf("list my plans: %w", err)
	}

	total, err := s.planRepo.CountMyPlans(ctx, clientID)
	if err != nil {
		s.logger.Error().Err(err).Str("client_id", clientID.String()).Msg("failed to count client plan history")
		return nil, fmt.Errorf("count my plans: %w", err)
	}

	data := make([]dto.DietPlanSummaryResponse, 0, len(plans))
	for _, p := range plans {
		planUUID := uuid.UUID(p.ID.Bytes)
		dayCount, err := s.planRepo.CountPlanDays(ctx, planUUID)
		if err != nil {
			s.logger.Error().Err(err).Str("plan_id", planUUID.String()).Msg("failed to count client plan days")
			return nil, fmt.Errorf("count plan days: %w", err)
		}
		data = append(data, planToSummary(&p, dayCount))
	}

	return &dto.DietPlanListResponse{
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasMore: int64(page*limit) < total,
	}, nil
}

// ─── Day CRUD ─────────────────────────────────────────────────────────────────

// AddDay adds a new day to a draft plan.
func (s *DietPlanService) AddDay(ctx context.Context, planID, nutritionistID uuid.UUID, req dto.CreateDayRequest) (*dto.PlanDayResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}

	day, err := s.planRepo.AddDay(ctx, sqlc.CreatePlanDayParams{
		PlanID:    pgtype.UUID{Bytes: planID, Valid: true},
		DayNumber: int32(req.DayNumber),
		Label:     optionalText(req.Label),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("plan_id", planID.String()).Msg("failed to add day")
		return nil, fmt.Errorf("add day: %w", err)
	}

	resp := dayToResponse(day)
	return &resp, nil
}

// UpdateDay updates a day's label.
func (s *DietPlanService) UpdateDay(ctx context.Context, planID, dayID, nutritionistID uuid.UUID, req dto.UpdateDayRequest) (*dto.PlanDayResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}

	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDayNotFound
		}
		return nil, fmt.Errorf("get day by id: %w", err)
	}

	day, err := s.planRepo.UpdateDay(ctx, sqlc.UpdatePlanDayParams{
		ID:    pgtype.UUID{Bytes: dayID, Valid: true},
		Label: optionalText(req.Label),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("day_id", dayID.String()).Msg("failed to update day")
		return nil, fmt.Errorf("update day: %w", err)
	}

	resp := dayToResponse(day)
	return &resp, nil
}

// DeleteDay removes a day from a draft plan.
func (s *DietPlanService) DeleteDay(ctx context.Context, planID, dayID, nutritionistID uuid.UUID) error {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return err
	}

	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDayNotFound
		}
		return fmt.Errorf("get day by id: %w", err)
	}

	return s.planRepo.DeleteDay(ctx, dayID, planID)
}

// ─── Meal CRUD ────────────────────────────────────────────────────────────────

// AddMeal adds a new meal to a plan day.
func (s *DietPlanService) AddMeal(ctx context.Context, planID, dayID, nutritionistID uuid.UUID, req dto.CreateMealRequest) (*dto.MealResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDayNotFound
		}
		return nil, fmt.Errorf("get day by id: %w", err)
	}

	scheduledTime, err := parseTimeToTime(req.ScheduledTime)
	if err != nil {
		return nil, fmt.Errorf("invalid scheduled_time: %w", err)
	}

	meal, err := s.planRepo.AddMeal(ctx, sqlc.CreateMealParams{
		DayID:         pgtype.UUID{Bytes: dayID, Valid: true},
		Title:         req.Title,
		ScheduledTime: scheduledTime,
		DisplayOrder:  int32(req.DisplayOrder),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("day_id", dayID.String()).Msg("failed to add meal")
		return nil, fmt.Errorf("add meal: %w", err)
	}

	resp := mealToResponse(meal)
	return &resp, nil
}

// UpdateMeal updates an existing meal.
func (s *DietPlanService) UpdateMeal(ctx context.Context, planID, dayID, mealID, nutritionistID uuid.UUID, req dto.UpdateMealRequest) (*dto.MealResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDayNotFound
		}
		return nil, fmt.Errorf("get day by id: %w", err)
	}
	if _, err := s.planRepo.GetMealByID(ctx, mealID, dayID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMealNotFound
		}
		return nil, fmt.Errorf("get meal by id: %w", err)
	}

	scheduledTime, err := parseTimeToTime(req.ScheduledTime)
	if err != nil {
		return nil, fmt.Errorf("invalid scheduled_time: %w", err)
	}

	meal, err := s.planRepo.UpdateMeal(ctx, sqlc.UpdateMealParams{
		ID:            pgtype.UUID{Bytes: mealID, Valid: true},
		Title:         req.Title,
		ScheduledTime: scheduledTime,
		DisplayOrder:  int32(req.DisplayOrder),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("meal_id", mealID.String()).Msg("failed to update meal")
		return nil, fmt.Errorf("update meal: %w", err)
	}

	resp := mealToResponse(meal)
	return &resp, nil
}

// DeleteMeal removes a meal from a plan day.
func (s *DietPlanService) DeleteMeal(ctx context.Context, planID, dayID, mealID, nutritionistID uuid.UUID) error {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return err
	}
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDayNotFound
		}
		return fmt.Errorf("get day by id: %w", err)
	}
	if _, err := s.planRepo.GetMealByID(ctx, mealID, dayID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMealNotFound
		}
		return fmt.Errorf("get meal by id: %w", err)
	}
	return s.planRepo.DeleteMeal(ctx, mealID, dayID)
}

// ReorderMeal changes a meal's display_order within its day.
func (s *DietPlanService) ReorderMeal(ctx context.Context, planID, dayID, mealID, nutritionistID uuid.UUID, newOrder int32) error {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return err
	}
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDayNotFound
		}
		return fmt.Errorf("get day by id: %w", err)
	}
	if _, err := s.planRepo.GetMealByID(ctx, mealID, dayID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMealNotFound
		}
		return fmt.Errorf("get meal by id: %w", err)
	}
	return s.planRepo.ReorderMeal(ctx, mealID, newOrder)
}

// ─── Option CRUD ──────────────────────────────────────────────────────────────

// AddOption appends a new option to a meal, auto-numbering via GetNextOptionNumber.
func (s *DietPlanService) AddOption(ctx context.Context, planID, dayID, mealID, nutritionistID uuid.UUID) (*dto.MealOptionResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDayNotFound
		}
		return nil, fmt.Errorf("get day by id: %w", err)
	}
	if _, err := s.planRepo.GetMealByID(ctx, mealID, dayID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMealNotFound
		}
		return nil, fmt.Errorf("get meal by id: %w", err)
	}

	nextNum, err := s.planRepo.GetNextOptionNumber(ctx, mealID)
	if err != nil {
		return nil, fmt.Errorf("get next option number: %w", err)
	}

	opt, err := s.planRepo.AddOption(ctx, sqlc.CreateMealOptionParams{
		MealID:       pgtype.UUID{Bytes: mealID, Valid: true},
		OptionNumber: int16(nextNum), // GetNextOptionNumber returns int32, cast to int16
		Label:        pgtype.Text{Valid: false},
	})
	if err != nil {
		s.logger.Error().Err(err).Str("meal_id", mealID.String()).Msg("failed to add option")
		return nil, fmt.Errorf("add option: %w", err)
	}

	resp := optionToResponse(opt)
	return &resp, nil
}

// DeleteOption removes an option from a meal.
func (s *DietPlanService) DeleteOption(ctx context.Context, planID, dayID, mealID, optID, nutritionistID uuid.UUID) error {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return err
	}
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDayNotFound
		}
		return fmt.Errorf("get day by id: %w", err)
	}
	if _, err := s.planRepo.GetMealByID(ctx, mealID, dayID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMealNotFound
		}
		return fmt.Errorf("get meal by id: %w", err)
	}
	if _, err := s.planRepo.GetOptionByID(ctx, optID, mealID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOptionNotFound
		}
		return fmt.Errorf("get option by id: %w", err)
	}
	return s.planRepo.DeleteOption(ctx, optID, mealID)
}

// ─── Item CRUD ────────────────────────────────────────────────────────────────

// AddItem adds a food item to a meal option.
func (s *DietPlanService) AddItem(ctx context.Context, planID, dayID, mealID, optID, nutritionistID uuid.UUID, req dto.CreateMealOptionItemRequest) (*dto.MealOptionItemResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}
	if err := s.verifyOptionChain(ctx, planID, dayID, mealID, optID); err != nil {
		return nil, err
	}

	foodID, err := uuid.Parse(req.FoodID)
	if err != nil {
		return nil, fmt.Errorf("invalid food_id: %w", err)
	}

	qty, err := numericFromFloat64(req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("convert quantity: %w", err)
	}

	item, err := s.planRepo.AddItem(ctx, sqlc.CreateMealOptionItemParams{
		OptionID:        pgtype.UUID{Bytes: optID, Valid: true},
		FoodID:          pgtype.UUID{Bytes: foodID, Valid: true},
		Quantity:        qty,
		MeasurementUnit: sqlc.MeasurementUnitType(req.MeasurementUnit),
		Notes:           optionalText(req.Notes),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("option_id", optID.String()).Msg("failed to add item")
		return nil, fmt.Errorf("add item: %w", err)
	}

	return itemToBasicResponse(item)
}

// UpdateItem updates an existing food item's quantity, unit and notes.
func (s *DietPlanService) UpdateItem(ctx context.Context, planID, dayID, mealID, optID, itemID, nutritionistID uuid.UUID, req dto.UpdateMealOptionItemRequest) (*dto.MealOptionItemResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}
	if err := s.verifyOptionChain(ctx, planID, dayID, mealID, optID); err != nil {
		return nil, err
	}
	if _, err := s.planRepo.GetItemByID(ctx, itemID, optID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("get item by id: %w", err)
	}

	qty, err := numericFromFloat64(req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("convert quantity: %w", err)
	}

	item, err := s.planRepo.UpdateItem(ctx, sqlc.UpdateMealOptionItemParams{
		ID:              pgtype.UUID{Bytes: itemID, Valid: true},
		Quantity:        qty,
		MeasurementUnit: sqlc.MeasurementUnitType(req.MeasurementUnit),
		Notes:           optionalText(req.Notes),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("item_id", itemID.String()).Msg("failed to update item")
		return nil, fmt.Errorf("update item: %w", err)
	}

	return itemToBasicResponse(item)
}

// DeleteItem removes a food item from a meal option.
func (s *DietPlanService) DeleteItem(ctx context.Context, planID, dayID, mealID, optID, itemID, nutritionistID uuid.UUID) error {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return err
	}
	if err := s.verifyOptionChain(ctx, planID, dayID, mealID, optID); err != nil {
		return err
	}
	if _, err := s.planRepo.GetItemByID(ctx, itemID, optID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrItemNotFound
		}
		return fmt.Errorf("get item by id: %w", err)
	}
	return s.planRepo.DeleteItem(ctx, itemID, optID)
}

// ─── Exercise CRUD ────────────────────────────────────────────────────────────

// AddExercise adds an exercise recommendation to a plan day.
func (s *DietPlanService) AddExercise(ctx context.Context, planID, dayID, nutritionistID uuid.UUID, req dto.CreateExerciseRequest) (*dto.PlanExerciseResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDayNotFound
		}
		return nil, fmt.Errorf("get day by id: %w", err)
	}

	ex, err := s.planRepo.AddExercise(ctx, sqlc.CreatePlanExerciseParams{
		DayID:                pgtype.UUID{Bytes: dayID, Valid: true},
		ExerciseName:         req.ExerciseName,
		DurationMinutes:      int32(req.DurationMinutes),
		Description:          optionalText(req.Description),
		CaloriesBurnEstimate: optionalInt4(req.CaloriesBurnEstimate),
		DisplayOrder:         int32(req.DisplayOrder),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("day_id", dayID.String()).Msg("failed to add exercise")
		return nil, fmt.Errorf("add exercise: %w", err)
	}

	resp := exerciseToResponse(ex)
	return &resp, nil
}

// UpdateExercise updates an existing exercise recommendation.
func (s *DietPlanService) UpdateExercise(ctx context.Context, planID, dayID, exID, nutritionistID uuid.UUID, req dto.UpdateExerciseRequest) (*dto.PlanExerciseResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDayNotFound
		}
		return nil, fmt.Errorf("get day by id: %w", err)
	}
	if _, err := s.planRepo.GetExerciseByID(ctx, exID, dayID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrExerciseNotFound
		}
		return nil, fmt.Errorf("get exercise by id: %w", err)
	}

	ex, err := s.planRepo.UpdateExercise(ctx, sqlc.UpdatePlanExerciseParams{
		ID:                   pgtype.UUID{Bytes: exID, Valid: true},
		ExerciseName:         req.ExerciseName,
		DurationMinutes:      int32(req.DurationMinutes),
		Description:          optionalText(req.Description),
		CaloriesBurnEstimate: optionalInt4(req.CaloriesBurnEstimate),
	})
	if err != nil {
		s.logger.Error().Err(err).Str("exercise_id", exID.String()).Msg("failed to update exercise")
		return nil, fmt.Errorf("update exercise: %w", err)
	}

	resp := exerciseToResponse(ex)
	return &resp, nil
}

// DeleteExercise removes an exercise from a plan day.
func (s *DietPlanService) DeleteExercise(ctx context.Context, planID, dayID, exID, nutritionistID uuid.UUID) error {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return err
	}
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDayNotFound
		}
		return fmt.Errorf("get day by id: %w", err)
	}
	if _, err := s.planRepo.GetExerciseByID(ctx, exID, dayID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrExerciseNotFound
		}
		return fmt.Errorf("get exercise by id: %w", err)
	}
	return s.planRepo.DeleteExercise(ctx, exID, dayID)
}

// ─── Medication CRUD ──────────────────────────────────────────────────────────

// AddMedication adds a medication prescription to a plan.
func (s *DietPlanService) AddMedication(ctx context.Context, planID, nutritionistID uuid.UUID, req dto.CreateMedicationPrescriptionRequest) (*dto.PlanMedicationResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}

	medID, err := uuid.Parse(req.MedicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid medication_id: %w", err)
	}

	timesJSON, err := json.Marshal(req.Times)
	if err != nil {
		return nil, fmt.Errorf("marshal times: %w", err)
	}

	startDate, err := parseDateToDateOptional(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}
	endDate, err := parseDateToDateOptional(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}

	med, err := s.planRepo.AddMedication(ctx, sqlc.CreatePlanMedicationParams{
		PlanID:       pgtype.UUID{Bytes: planID, Valid: true},
		MedicationID: pgtype.UUID{Bytes: medID, Valid: true},
		Dosage:       req.Dosage,
		Frequency:    req.Frequency,
		Times:        timesJSON,
		Instructions: optionalText(req.Instructions),
		StartDate:    startDate,
		EndDate:      endDate,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("plan_id", planID.String()).Msg("failed to add medication")
		return nil, fmt.Errorf("add medication: %w", err)
	}

	return medicationToResponse(med)
}

// UpdateMedication updates an existing medication prescription.
func (s *DietPlanService) UpdateMedication(ctx context.Context, planID, medID, nutritionistID uuid.UUID, req dto.UpdateMedicationPrescriptionRequest) (*dto.PlanMedicationResponse, error) {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return nil, err
	}
	if _, err := s.planRepo.GetMedicationByID(ctx, medID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMedicationPrescNotFound
		}
		return nil, fmt.Errorf("get medication by id: %w", err)
	}

	timesJSON, err := json.Marshal(req.Times)
	if err != nil {
		return nil, fmt.Errorf("marshal times: %w", err)
	}

	startDate, err := parseDateToDateOptional(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}
	endDate, err := parseDateToDateOptional(req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}

	med, err := s.planRepo.UpdateMedication(ctx, sqlc.UpdatePlanMedicationParams{
		ID:           pgtype.UUID{Bytes: medID, Valid: true},
		Dosage:       req.Dosage,
		Frequency:    req.Frequency,
		Times:        timesJSON,
		Instructions: optionalText(req.Instructions),
		StartDate:    startDate,
		EndDate:      endDate,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("med_id", medID.String()).Msg("failed to update medication")
		return nil, fmt.Errorf("update medication: %w", err)
	}

	return medicationToResponse(med)
}

// DeleteMedication removes a medication prescription from a plan.
func (s *DietPlanService) DeleteMedication(ctx context.Context, planID, medID, nutritionistID uuid.UUID) error {
	if err := s.requireDraft(ctx, planID, nutritionistID); err != nil {
		return err
	}
	if _, err := s.planRepo.GetMedicationByID(ctx, medID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMedicationPrescNotFound
		}
		return fmt.Errorf("get medication by id: %w", err)
	}
	return s.planRepo.DeleteMedication(ctx, medID, planID)
}

// ─── Private helpers ──────────────────────────────────────────────────────────

// requireDraft verifies that the plan exists, belongs to nutritionistID, and is in draft status.
func (s *DietPlanService) requireDraft(ctx context.Context, planID, nutritionistID uuid.UUID) error {
	plan, err := s.planRepo.GetPlanByID(ctx, planID, nutritionistID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPlanNotFound
		}
		return fmt.Errorf("get plan by id: %w", err)
	}
	if plan.Status != sqlc.DietPlanStatusDraft {
		return ErrPlanNotDraft
	}
	return nil
}

// verifyOptionChain verifies: day in plan, meal in day, option in meal.
func (s *DietPlanService) verifyOptionChain(ctx context.Context, planID, dayID, mealID, optID uuid.UUID) error {
	if _, err := s.planRepo.GetDayByID(ctx, dayID, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDayNotFound
		}
		return fmt.Errorf("get day by id: %w", err)
	}
	if _, err := s.planRepo.GetMealByID(ctx, mealID, dayID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMealNotFound
		}
		return fmt.Errorf("get meal by id: %w", err)
	}
	if _, err := s.planRepo.GetOptionByID(ctx, optID, mealID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOptionNotFound
		}
		return fmt.Errorf("get option by id: %w", err)
	}
	return nil
}

// ─── Conversion helpers ───────────────────────────────────────────────────────

// parseDateToDate converts "YYYY-MM-DD" string to pgtype.Date (required).
func parseDateToDate(s string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("parse date %q: %w", s, err)
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}

// parseDateToDateOptional converts an optional *string to pgtype.Date (nil → invalid).
func parseDateToDateOptional(s *string) (pgtype.Date, error) {
	if s == nil || *s == "" {
		return pgtype.Date{Valid: false}, nil
	}
	return parseDateToDate(*s)
}

// parseTimeToTime converts an optional *string "HH:MM" to pgtype.Time.
func parseTimeToTime(s *string) (pgtype.Time, error) {
	if s == nil || *s == "" {
		return pgtype.Time{Valid: false}, nil
	}
	t, err := time.Parse("15:04", *s)
	if err != nil {
		return pgtype.Time{}, fmt.Errorf("parse time %q: %w", *s, err)
	}
	micros := int64(t.Hour())*3_600_000_000 + int64(t.Minute())*60_000_000
	return pgtype.Time{Microseconds: micros, Valid: true}, nil
}

// formatDateVal formats pgtype.Date as "YYYY-MM-DD" string.
func formatDateVal(v pgtype.Date) string {
	if !v.Valid {
		return ""
	}
	return v.Time.Format("2006-01-02")
}

// formatTimeVal formats pgtype.Time as "HH:MM" string.
func formatTimeVal(v pgtype.Time) string {
	if !v.Valid {
		return ""
	}
	hours := v.Microseconds / 3_600_000_000
	minutes := (v.Microseconds % 3_600_000_000) / 60_000_000
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

// optionalInt4 converts *int to pgtype.Int4.
func optionalInt4(v *int) pgtype.Int4 {
	if v == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*v), Valid: true}
}

// planToSummary converts a sqlc.DietPlan + dayCount to DietPlanSummaryResponse.
func planToSummary(plan *sqlc.DietPlan, dayCount int64) dto.DietPlanSummaryResponse {
	return dto.DietPlanSummaryResponse{
		ID:        uuid.UUID(plan.ID.Bytes).String(),
		Status:    string(plan.Status),
		StartDate: formatDateVal(plan.StartDate),
		EndDate:   formatDateVal(plan.EndDate),
		DayCount:  dayCount,
		CreatedAt: formatTimestamp(plan.CreatedAt),
	}
}

// dayToResponse converts a sqlc.PlanDay (or pointer) to PlanDayResponse (no nested children).
func dayToResponse(day *sqlc.PlanDay) dto.PlanDayResponse {
	resp := dto.PlanDayResponse{
		ID:        uuid.UUID(day.ID.Bytes).String(),
		DayNumber: int(day.DayNumber),
		Meals:     []dto.MealResponse{},
		Exercises: []dto.PlanExerciseResponse{},
	}
	if day.Label.Valid {
		resp.Label = &day.Label.String
	}
	return resp
}

// mealToResponse converts a sqlc.Meal to MealResponse (no nested options).
func mealToResponse(meal *sqlc.Meal) dto.MealResponse {
	resp := dto.MealResponse{
		ID:           uuid.UUID(meal.ID.Bytes).String(),
		Title:        meal.Title,
		DisplayOrder: int(meal.DisplayOrder),
		Options:      []dto.MealOptionResponse{},
	}
	if meal.ScheduledTime.Valid {
		t := formatTimeVal(meal.ScheduledTime)
		resp.ScheduledTime = &t
	}
	return resp
}

// optionToResponse converts a sqlc.MealOption to MealOptionResponse (no nested items).
func optionToResponse(opt *sqlc.MealOption) dto.MealOptionResponse {
	resp := dto.MealOptionResponse{
		ID:           uuid.UUID(opt.ID.Bytes).String(),
		OptionNumber: int(opt.OptionNumber),
		Items:        []dto.MealOptionItemResponse{},
	}
	if opt.Label.Valid {
		resp.Label = &opt.Label.String
	}
	return resp
}

// exerciseToResponse converts a sqlc.PlanExercise to PlanExerciseResponse.
func exerciseToResponse(ex *sqlc.PlanExercise) dto.PlanExerciseResponse {
	resp := dto.PlanExerciseResponse{
		ID:              uuid.UUID(ex.ID.Bytes).String(),
		ExerciseName:    ex.ExerciseName,
		DurationMinutes: int(ex.DurationMinutes),
		DisplayOrder:    int(ex.DisplayOrder),
	}
	if ex.Description.Valid {
		resp.Description = &ex.Description.String
	}
	if ex.CaloriesBurnEstimate.Valid {
		c := int(ex.CaloriesBurnEstimate.Int32)
		resp.CaloriesBurnEstimate = &c
	}
	return resp
}

// medicationToResponse converts a sqlc.PlanMedication to PlanMedicationResponse.
// MedicationName and MedicationForm are not available in single-row mutations; they remain empty.
func medicationToResponse(med *sqlc.PlanMedication) (*dto.PlanMedicationResponse, error) {
	var times []string
	if len(med.Times) > 0 {
		if err := json.Unmarshal(med.Times, &times); err != nil {
			return nil, fmt.Errorf("unmarshal medication times: %w", err)
		}
	}
	if times == nil {
		times = []string{}
	}

	resp := &dto.PlanMedicationResponse{
		ID:           uuid.UUID(med.ID.Bytes).String(),
		MedicationID: uuid.UUID(med.MedicationID.Bytes).String(),
		Dosage:       med.Dosage,
		Frequency:    med.Frequency,
		Times:        times,
	}
	if med.Instructions.Valid {
		resp.Instructions = &med.Instructions.String
	}
	if med.StartDate.Valid {
		s := med.StartDate.Time.Format("2006-01-02")
		resp.StartDate = &s
	}
	if med.EndDate.Valid {
		e := med.EndDate.Time.Format("2006-01-02")
		resp.EndDate = &e
	}
	return resp, nil
}

// itemToBasicResponse converts a sqlc.MealOptionItem to MealOptionItemResponse.
// Food.ID is populated from FoodID; other food fields are zero (not available without a JOIN).
func itemToBasicResponse(item *sqlc.MealOptionItem) (*dto.MealOptionItemResponse, error) {
	qty, err := numericToFloat64(item.Quantity)
	if err != nil {
		return nil, fmt.Errorf("convert quantity: %w", err)
	}
	resp := &dto.MealOptionItemResponse{
		ID:              uuid.UUID(item.ID.Bytes).String(),
		Quantity:        qty,
		MeasurementUnit: string(item.MeasurementUnit),
		Food: dto.FoodEmbedded{
			ID: uuid.UUID(item.FoodID.Bytes).String(),
		},
	}
	if item.Notes.Valid {
		resp.Notes = &item.Notes.String
	}
	return resp, nil
}
