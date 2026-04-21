package dietplan

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/dietplan/repository"
)

const (
	activePlanCacheTTL    = 2 * time.Minute
	activePlanCachePrefix = "dietplan:active:"
)

// CachedDietPlanRepository wraps a DietPlanRepository and caches the active plan per client in Redis.
// Only FindActiveByClientID is cached; all other methods delegate directly to the inner repository.
// Write operations that may change the active plan (CreateWithArchive, Update, Delete) invalidate
// the relevant cache entry.
type CachedDietPlanRepository struct {
	inner repository.DietPlanRepository
	rdb   *redis.Client
}

// NewCachedDietPlanRepository returns a CachedDietPlanRepository backed by the given inner repo and Redis client.
func NewCachedDietPlanRepository(inner repository.DietPlanRepository, rdb *redis.Client) *CachedDietPlanRepository {
	return &CachedDietPlanRepository{inner: inner, rdb: rdb}
}

func (r *CachedDietPlanRepository) activePlanKey(clientID uuid.UUID) string {
	return fmt.Sprintf("%s%s", activePlanCachePrefix, clientID.String())
}

func (r *CachedDietPlanRepository) invalidateActiveCache(ctx context.Context, clientID uuid.UUID) {
	r.rdb.Del(ctx, r.activePlanKey(clientID))
}

// --- Cached read ---

// FindActiveByClientID retrieves the active plan for a client, using a 2-minute Redis cache.
func (r *CachedDietPlanRepository) FindActiveByClientID(ctx context.Context, clientID uuid.UUID) (*entity.DietPlan, error) {
	key := r.activePlanKey(clientID)

	if cached, err := r.rdb.Get(ctx, key).Bytes(); err == nil {
		// Distinguish "nil plan" sentinel from a real plan.
		if string(cached) == "nil" {
			return nil, nil
		}
		var plan entity.DietPlan
		if json.Unmarshal(cached, &plan) == nil {
			return &plan, nil
		}
	}

	plan, err := r.inner.FindActiveByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if plan == nil {
		// Cache the "no active plan" result to avoid repeated DB hits.
		r.rdb.Set(ctx, key, "nil", activePlanCacheTTL)
	} else {
		if data, mErr := json.Marshal(plan); mErr == nil {
			r.rdb.Set(ctx, key, data, activePlanCacheTTL)
		}
	}

	return plan, nil
}

// --- Write operations (with cache invalidation) ---

// CreateWithArchive atomically archives the existing active plan and creates a new one,
// then invalidates the active-plan cache for the client.
func (r *CachedDietPlanRepository) CreateWithArchive(ctx context.Context, plan *entity.DietPlan) error {
	if err := r.inner.CreateWithArchive(ctx, plan); err != nil {
		return err
	}
	r.invalidateActiveCache(ctx, plan.ClientID)
	return nil
}

// Update persists updated diet plan fields and invalidates the active-plan cache.
func (r *CachedDietPlanRepository) Update(ctx context.Context, plan *entity.DietPlan) error {
	if err := r.inner.Update(ctx, plan); err != nil {
		return err
	}
	r.invalidateActiveCache(ctx, plan.ClientID)
	return nil
}

// Delete removes a diet plan. It fetches the plan first to obtain the clientID for cache invalidation.
func (r *CachedDietPlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Fetch before delete to get clientID for cache invalidation.
	plan, _ := r.inner.FindByID(ctx, id)
	if err := r.inner.Delete(ctx, id); err != nil {
		return err
	}
	if plan != nil {
		r.invalidateActiveCache(ctx, plan.ClientID)
	}
	return nil
}

// --- Pass-through delegations ---

func (r *CachedDietPlanRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.DietPlan, error) {
	return r.inner.FindByID(ctx, id)
}

func (r *CachedDietPlanRepository) ListByClientID(ctx context.Context, clientID uuid.UUID, limit, offset int32) ([]*entity.DietPlan, error) {
	return r.inner.ListByClientID(ctx, clientID, limit, offset)
}

func (r *CachedDietPlanRepository) CountByClientID(ctx context.Context, clientID uuid.UUID) (int64, error) {
	return r.inner.CountByClientID(ctx, clientID)
}

func (r *CachedDietPlanRepository) AddDay(ctx context.Context, day *entity.DietPlanDay) error {
	return r.inner.AddDay(ctx, day)
}

func (r *CachedDietPlanRepository) FindDayByID(ctx context.Context, id uuid.UUID) (*entity.DietPlanDay, error) {
	return r.inner.FindDayByID(ctx, id)
}

func (r *CachedDietPlanRepository) ListDays(ctx context.Context, planID uuid.UUID) ([]*entity.DietPlanDay, error) {
	return r.inner.ListDays(ctx, planID)
}

func (r *CachedDietPlanRepository) DeleteDay(ctx context.Context, id uuid.UUID) error {
	return r.inner.DeleteDay(ctx, id)
}

func (r *CachedDietPlanRepository) AddMeal(ctx context.Context, meal *entity.DietMeal) error {
	return r.inner.AddMeal(ctx, meal)
}

func (r *CachedDietPlanRepository) FindMealByID(ctx context.Context, id uuid.UUID) (*entity.DietMeal, error) {
	return r.inner.FindMealByID(ctx, id)
}

func (r *CachedDietPlanRepository) ListMeals(ctx context.Context, dayID uuid.UUID) ([]*entity.DietMeal, error) {
	return r.inner.ListMeals(ctx, dayID)
}

func (r *CachedDietPlanRepository) DeleteMeal(ctx context.Context, id uuid.UUID) error {
	return r.inner.DeleteMeal(ctx, id)
}

func (r *CachedDietPlanRepository) AddOption(ctx context.Context, option *entity.MealOption) error {
	return r.inner.AddOption(ctx, option)
}

func (r *CachedDietPlanRepository) FindOptionByID(ctx context.Context, id uuid.UUID) (*entity.MealOption, error) {
	return r.inner.FindOptionByID(ctx, id)
}

func (r *CachedDietPlanRepository) ListOptions(ctx context.Context, mealID uuid.UUID) ([]*entity.MealOption, error) {
	return r.inner.ListOptions(ctx, mealID)
}

func (r *CachedDietPlanRepository) DeleteOption(ctx context.Context, id uuid.UUID) error {
	return r.inner.DeleteOption(ctx, id)
}

func (r *CachedDietPlanRepository) AddItem(ctx context.Context, item *entity.MealOptionItem) error {
	return r.inner.AddItem(ctx, item)
}

func (r *CachedDietPlanRepository) FindItemByID(ctx context.Context, id uuid.UUID) (*entity.MealOptionItem, error) {
	return r.inner.FindItemByID(ctx, id)
}

func (r *CachedDietPlanRepository) ListItems(ctx context.Context, optionID uuid.UUID) ([]*entity.MealOptionItem, error) {
	return r.inner.ListItems(ctx, optionID)
}

func (r *CachedDietPlanRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	return r.inner.DeleteItem(ctx, id)
}

func (r *CachedDietPlanRepository) DeleteItemsByOption(ctx context.Context, optionID uuid.UUID) error {
	return r.inner.DeleteItemsByOption(ctx, optionID)
}

func (r *CachedDietPlanRepository) ListItemsWithFood(ctx context.Context, optionID uuid.UUID) ([]*entity.MealOptionItem, error) {
	return r.inner.ListItemsWithFood(ctx, optionID)
}

func (r *CachedDietPlanRepository) AddExercise(ctx context.Context, ex *entity.ExerciseRecommendation) error {
	return r.inner.AddExercise(ctx, ex)
}

func (r *CachedDietPlanRepository) FindExerciseByID(ctx context.Context, id uuid.UUID) (*entity.ExerciseRecommendation, error) {
	return r.inner.FindExerciseByID(ctx, id)
}

func (r *CachedDietPlanRepository) ListExercises(ctx context.Context, dayID uuid.UUID) ([]*entity.ExerciseRecommendation, error) {
	return r.inner.ListExercises(ctx, dayID)
}

func (r *CachedDietPlanRepository) DeleteExercise(ctx context.Context, id uuid.UUID) error {
	return r.inner.DeleteExercise(ctx, id)
}

func (r *CachedDietPlanRepository) AddPrescription(ctx context.Context, rx *entity.PrescribedMedication) error {
	return r.inner.AddPrescription(ctx, rx)
}

func (r *CachedDietPlanRepository) FindPrescriptionByID(ctx context.Context, id uuid.UUID) (*entity.PrescribedMedication, error) {
	return r.inner.FindPrescriptionByID(ctx, id)
}

func (r *CachedDietPlanRepository) ListPrescriptionsWithMedication(ctx context.Context, dayID uuid.UUID) ([]*entity.PrescribedMedication, error) {
	return r.inner.ListPrescriptionsWithMedication(ctx, dayID)
}

func (r *CachedDietPlanRepository) DeletePrescription(ctx context.Context, id uuid.UUID) error {
	return r.inner.DeletePrescription(ctx, id)
}
