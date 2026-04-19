package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"

	"github.com/ranjbar-dev/nutritrack/backend/internal/model"
	"github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
	"github.com/ranjbar-dev/nutritrack/backend/internal/repository/sqlc"
)

var (
	ErrFoodDuplicate          = errors.New("غذا با این نام قبلاً ثبت شده است")
	ErrFoodNotFound           = errors.New("غذا یافت نشد")
	ErrFoodUnauthorizedEdit   = errors.New("شما مجوز ویرایش این غذا را ندارید")
	ErrFoodUnauthorizedDelete = errors.New("شما مجوز حذف این غذا را ندارید")
)

// FoodService handles food management business logic.
type FoodService struct {
	foodRepo repository.FoodRepository
	logger   zerolog.Logger
}

// NewFoodService creates a new FoodService with the given dependencies.
func NewFoodService(foodRepo repository.FoodRepository, logger zerolog.Logger) *FoodService {
	return &FoodService{
		foodRepo: foodRepo,
		logger:   logger,
	}
}

func (s *FoodService) CreateFood(ctx context.Context, userID uuid.UUID, req dto.CreateFoodRequest) (*dto.FoodResponse, error) {
	isDuplicate, err := s.foodRepo.CheckDuplicateName(ctx, req.Name, nil)
	if err != nil {
		s.logger.Error().Err(err).Str("name", req.Name).Msg("failed to check duplicate food name")
		return nil, fmt.Errorf("check duplicate food name: %w", err)
	}
	if isDuplicate {
		return nil, ErrFoodDuplicate
	}

	params, err := createFoodParamsFromRequest(userID, req)
	if err != nil {
		return nil, err
	}

	food, err := s.foodRepo.Create(ctx, params)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID.String()).Msg("failed to create food")
		return nil, fmt.Errorf("create food: %w", err)
	}

	foodID := uuid.UUID(food.ID.Bytes)
	for _, category := range req.Categories {
		if err := s.foodRepo.AddCategory(ctx, foodID, category); err != nil {
			s.logger.Error().Err(err).Str("food_id", foodID.String()).Str("category", category).Msg("failed to add food category")
			return nil, fmt.Errorf("add food category: %w", err)
		}
	}

	return s.GetFood(ctx, foodID)
}

func (s *FoodService) GetFood(ctx context.Context, foodID uuid.UUID) (*dto.FoodResponse, error) {
	food, err := s.foodRepo.GetByID(ctx, foodID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFoodNotFound
		}
		s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to get food")
		return nil, fmt.Errorf("get food: %w", err)
	}

	categories, err := s.foodRepo.GetCategories(ctx, foodID)
	if err != nil {
		s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to get food categories")
		return nil, fmt.Errorf("get food categories: %w", err)
	}

	return foodRowToResponse(food, categories)
}

func (s *FoodService) ListFoods(ctx context.Context, query dto.FoodListQueryParams) (*dto.FoodListResponse, error) {
	page := query.Page
	if page < 1 {
		page = 1
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	listParams := sqlc.ListFoodsParams{
		IsActive:  optionalBool(query.IsActive),
		Search:    optionalText(query.Search),
		Category:  optionalText(query.Category),
		OffsetVal: int64((page - 1) * limit),
		LimitVal:  int64(limit),
	}

	countParams := sqlc.CountFoodsParams{
		IsActive: optionalBool(query.IsActive),
		Search:   optionalText(query.Search),
		Category: optionalText(query.Category),
	}

	foods, err := s.foodRepo.List(ctx, listParams)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list foods")
		return nil, fmt.Errorf("list foods: %w", err)
	}

	total, err := s.foodRepo.Count(ctx, countParams)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to count foods")
		return nil, fmt.Errorf("count foods: %w", err)
	}

	data := make([]dto.FoodResponse, 0, len(foods))
	for _, food := range foods {
		foodID := uuid.UUID(food.ID.Bytes)
		categories, err := s.foodRepo.GetCategories(ctx, foodID)
		if err != nil {
			s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to get categories while listing foods")
			return nil, fmt.Errorf("get food categories: %w", err)
		}

		resp, err := listRowToResponse(food, categories)
		if err != nil {
			return nil, err
		}
		data = append(data, *resp)
	}

	return &dto.FoodListResponse{
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
		HasMore: int64(page*limit) < total,
	}, nil
}

func (s *FoodService) UpdateFood(ctx context.Context, foodID, userID uuid.UUID, role string, req dto.UpdateFoodRequest) (*dto.FoodResponse, error) {
	current, err := s.foodRepo.GetByID(ctx, foodID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFoodNotFound
		}
		s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to load food for update")
		return nil, fmt.Errorf("load food for update: %w", err)
	}

	if role == string(model.RoleNutritionist) && uuid.UUID(current.CreatedBy.Bytes) != userID {
		return nil, ErrFoodUnauthorizedEdit
	}

	isDuplicate, err := s.foodRepo.CheckDuplicateName(ctx, req.Name, &foodID)
	if err != nil {
		s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to check duplicate food name for update")
		return nil, fmt.Errorf("check duplicate food name: %w", err)
	}
	if isDuplicate {
		return nil, ErrFoodDuplicate
	}

	params, err := updateFoodParamsFromRequest(foodID, req)
	if err != nil {
		return nil, err
	}

	if err := s.foodRepo.DeleteCategories(ctx, foodID); err != nil {
		s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to delete existing food categories")
		return nil, fmt.Errorf("delete food categories: %w", err)
	}

	for _, category := range req.Categories {
		if err := s.foodRepo.AddCategory(ctx, foodID, category); err != nil {
			s.logger.Error().Err(err).Str("food_id", foodID.String()).Str("category", category).Msg("failed to add food category during update")
			return nil, fmt.Errorf("add food category: %w", err)
		}
	}

	if _, err := s.foodRepo.Update(ctx, params); err != nil {
		s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to update food")
		return nil, fmt.Errorf("update food: %w", err)
	}

	return s.GetFood(ctx, foodID)
}

func (s *FoodService) DeleteFood(ctx context.Context, foodID, userID uuid.UUID, role string) error {
	current, err := s.foodRepo.GetByID(ctx, foodID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrFoodNotFound
		}
		s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to load food for delete")
		return fmt.Errorf("load food for delete: %w", err)
	}

	switch role {
	case string(model.RoleSuperAdmin):
		if err := s.foodRepo.SoftDelete(ctx, foodID); err != nil {
			s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to soft delete food")
			return fmt.Errorf("soft delete food: %w", err)
		}
	case string(model.RoleNutritionist):
		if uuid.UUID(current.CreatedBy.Bytes) != userID {
			return ErrFoodUnauthorizedDelete
		}
		if err := s.foodRepo.SoftDeleteByOwner(ctx, foodID, userID); err != nil {
			s.logger.Error().Err(err).Str("food_id", foodID.String()).Str("user_id", userID.String()).Msg("failed to soft delete own food")
			return fmt.Errorf("soft delete food by owner: %w", err)
		}
	default:
		return ErrFoodUnauthorizedDelete
	}

	return nil
}

func createFoodParamsFromRequest(userID uuid.UUID, req dto.CreateFoodRequest) (sqlc.CreateFoodParams, error) {
	calories, err := numericFromFloat64(req.Calories)
	if err != nil {
		return sqlc.CreateFoodParams{}, fmt.Errorf("convert calories: %w", err)
	}
	protein, err := numericFromFloat64(req.ProteinG)
	if err != nil {
		return sqlc.CreateFoodParams{}, fmt.Errorf("convert protein: %w", err)
	}
	carbs, err := numericFromFloat64(req.CarbsG)
	if err != nil {
		return sqlc.CreateFoodParams{}, fmt.Errorf("convert carbs: %w", err)
	}
	fat, err := numericFromFloat64(req.FatG)
	if err != nil {
		return sqlc.CreateFoodParams{}, fmt.Errorf("convert fat: %w", err)
	}
	fiber, err := numericFromFloat64(req.FiberG)
	if err != nil {
		return sqlc.CreateFoodParams{}, fmt.Errorf("convert fiber: %w", err)
	}
	sugar, err := numericFromFloat64(req.SugarG)
	if err != nil {
		return sqlc.CreateFoodParams{}, fmt.Errorf("convert sugar: %w", err)
	}
	sodium, err := numericFromFloat64(req.SodiumMg)
	if err != nil {
		return sqlc.CreateFoodParams{}, fmt.Errorf("convert sodium: %w", err)
	}
	measurementAmount, err := numericFromFloat64(req.MeasurementAmount)
	if err != nil {
		return sqlc.CreateFoodParams{}, fmt.Errorf("convert measurement amount: %w", err)
	}

	params := sqlc.CreateFoodParams{
		Name:              req.Name,
		Description:       optionalDescription(req.Description),
		Calories:          calories,
		ProteinG:          protein,
		CarbsG:            carbs,
		FatG:              fat,
		FiberG:            fiber,
		SugarG:            sugar,
		SodiumMg:          sodium,
		MeasurementUnit:   sqlc.MeasurementUnitType(req.MeasurementUnit),
		MeasurementAmount: measurementAmount,
		CreatedBy:         pgtype.UUID{Bytes: userID, Valid: true},
	}

	return params, nil
}

func updateFoodParamsFromRequest(foodID uuid.UUID, req dto.UpdateFoodRequest) (sqlc.UpdateFoodParams, error) {
	calories, err := numericFromFloat64(req.Calories)
	if err != nil {
		return sqlc.UpdateFoodParams{}, fmt.Errorf("convert calories: %w", err)
	}
	protein, err := numericFromFloat64(req.ProteinG)
	if err != nil {
		return sqlc.UpdateFoodParams{}, fmt.Errorf("convert protein: %w", err)
	}
	carbs, err := numericFromFloat64(req.CarbsG)
	if err != nil {
		return sqlc.UpdateFoodParams{}, fmt.Errorf("convert carbs: %w", err)
	}
	fat, err := numericFromFloat64(req.FatG)
	if err != nil {
		return sqlc.UpdateFoodParams{}, fmt.Errorf("convert fat: %w", err)
	}
	fiber, err := numericFromFloat64(req.FiberG)
	if err != nil {
		return sqlc.UpdateFoodParams{}, fmt.Errorf("convert fiber: %w", err)
	}
	sugar, err := numericFromFloat64(req.SugarG)
	if err != nil {
		return sqlc.UpdateFoodParams{}, fmt.Errorf("convert sugar: %w", err)
	}
	sodium, err := numericFromFloat64(req.SodiumMg)
	if err != nil {
		return sqlc.UpdateFoodParams{}, fmt.Errorf("convert sodium: %w", err)
	}
	measurementAmount, err := numericFromFloat64(req.MeasurementAmount)
	if err != nil {
		return sqlc.UpdateFoodParams{}, fmt.Errorf("convert measurement amount: %w", err)
	}

	return sqlc.UpdateFoodParams{
		ID:                pgtype.UUID{Bytes: foodID, Valid: true},
		Name:              req.Name,
		Description:       optionalDescription(req.Description),
		Calories:          calories,
		ProteinG:          protein,
		CarbsG:            carbs,
		FatG:              fat,
		FiberG:            fiber,
		SugarG:            sugar,
		SodiumMg:          sodium,
		MeasurementUnit:   sqlc.MeasurementUnitType(req.MeasurementUnit),
		MeasurementAmount: measurementAmount,
	}, nil
}

func foodRowToResponse(food *sqlc.GetFoodByIDRow, categories []string) (*dto.FoodResponse, error) {
	calories, err := numericToFloat64(food.Calories)
	if err != nil {
		return nil, err
	}
	protein, err := numericToFloat64(food.ProteinG)
	if err != nil {
		return nil, err
	}
	carbs, err := numericToFloat64(food.CarbsG)
	if err != nil {
		return nil, err
	}
	fat, err := numericToFloat64(food.FatG)
	if err != nil {
		return nil, err
	}
	fiber, err := numericToFloat64(food.FiberG)
	if err != nil {
		return nil, err
	}
	sugar, err := numericToFloat64(food.SugarG)
	if err != nil {
		return nil, err
	}
	sodium, err := numericToFloat64(food.SodiumMg)
	if err != nil {
		return nil, err
	}
	measurementAmount, err := numericToFloat64(food.MeasurementAmount)
	if err != nil {
		return nil, err
	}

	resp := &dto.FoodResponse{
		ID:                uuid.UUID(food.ID.Bytes).String(),
		Name:              food.Name,
		Categories:        categories,
		Calories:          calories,
		ProteinG:          protein,
		CarbsG:            carbs,
		FatG:              fat,
		FiberG:            fiber,
		SugarG:            sugar,
		SodiumMg:          sodium,
		MeasurementUnit:   string(food.MeasurementUnit),
		MeasurementAmount: measurementAmount,
		IsActive:          food.IsActive,
		CreatedBy:         uuid.UUID(food.CreatedBy.Bytes).String(),
		CreatorName:       food.CreatorName,
		CreatedAt:         formatTimestamp(food.CreatedAt),
		UpdatedAt:         formatTimestamp(food.UpdatedAt),
	}
	if food.Description.Valid {
		resp.Description = &food.Description.String
	}

	return resp, nil
}

func listRowToResponse(food sqlc.ListFoodsRow, categories []string) (*dto.FoodResponse, error) {
	calories, err := numericToFloat64(food.Calories)
	if err != nil {
		return nil, err
	}
	protein, err := numericToFloat64(food.ProteinG)
	if err != nil {
		return nil, err
	}
	carbs, err := numericToFloat64(food.CarbsG)
	if err != nil {
		return nil, err
	}
	fat, err := numericToFloat64(food.FatG)
	if err != nil {
		return nil, err
	}
	fiber, err := numericToFloat64(food.FiberG)
	if err != nil {
		return nil, err
	}
	sugar, err := numericToFloat64(food.SugarG)
	if err != nil {
		return nil, err
	}
	sodium, err := numericToFloat64(food.SodiumMg)
	if err != nil {
		return nil, err
	}
	measurementAmount, err := numericToFloat64(food.MeasurementAmount)
	if err != nil {
		return nil, err
	}

	resp := &dto.FoodResponse{
		ID:                uuid.UUID(food.ID.Bytes).String(),
		Name:              food.Name,
		Categories:        categories,
		Calories:          calories,
		ProteinG:          protein,
		CarbsG:            carbs,
		FatG:              fat,
		FiberG:            fiber,
		SugarG:            sugar,
		SodiumMg:          sodium,
		MeasurementUnit:   string(food.MeasurementUnit),
		MeasurementAmount: measurementAmount,
		IsActive:          food.IsActive,
		CreatedBy:         uuid.UUID(food.CreatedBy.Bytes).String(),
		CreatorName:       food.CreatorName,
		CreatedAt:         formatTimestamp(food.CreatedAt),
		UpdatedAt:         formatTimestamp(food.UpdatedAt),
	}
	if food.Description.Valid {
		resp.Description = &food.Description.String
	}

	return resp, nil
}

func numericFromFloat64(value float64) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	if err := numeric.Scan(strconv.FormatFloat(value, 'f', -1, 64)); err != nil {
		return pgtype.Numeric{}, err
	}
	return numeric, nil
}

func numericToFloat64(value pgtype.Numeric) (float64, error) {
	if !value.Valid {
		return 0, nil
	}
	floatValue, err := value.Float64Value()
	if err != nil {
		return 0, err
	}
	return floatValue.Float64, nil
}

func optionalDescription(description *string) pgtype.Text {
	if description == nil || *description == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *description, Valid: true}
}

func optionalText(value *string) pgtype.Text {
	if value == nil || *value == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func optionalBool(value *bool) pgtype.Bool {
	if value == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *value, Valid: true}
}

func formatTimestamp(value pgtype.Timestamptz) string {
	if !value.Valid {
		return ""
	}
	return value.Time.Format(time.RFC3339)
}
