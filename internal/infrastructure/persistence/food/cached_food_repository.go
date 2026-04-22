package food

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/entity"
	"github.com/ranjbar-dev/nutritrack/internal/domain/food/repository"
)

const (
	foodCacheTTL    = 5 * time.Minute
	foodCachePrefix = "food:cache:"
)

// CachedFoodRepository wraps a FoodRepository and caches search/count results in Redis.
// Write operations (Create, Update, Delete, Deactivate) invalidate all search cache entries.
type CachedFoodRepository struct {
	inner repository.FoodRepository
	rdb   *redis.Client
}

// NewCachedFoodRepository returns a CachedFoodRepository backed by the given inner repo and Redis client.
func NewCachedFoodRepository(inner repository.FoodRepository, rdb *redis.Client) *CachedFoodRepository {
	return &CachedFoodRepository{inner: inner, rdb: rdb}
}

// --- Cache key helpers ---

func (r *CachedFoodRepository) searchKey(query string, limit, offset int32) string {
	return fmt.Sprintf("%ssearch:noc:%s:%d:%d", foodCachePrefix, query, limit, offset)
}

func (r *CachedFoodRepository) countKey(query string) string {
	return fmt.Sprintf("%scount:noc:%s", foodCachePrefix, query)
}

func (r *CachedFoodRepository) searchByCategoryKey(categoryID uuid.UUID, query string, limit, offset int32) string {
	return fmt.Sprintf("%ssearch:cat:%s:%s:%d:%d", foodCachePrefix, categoryID, query, limit, offset)
}

func (r *CachedFoodRepository) countByCategoryKey(categoryID uuid.UUID, query string) string {
	return fmt.Sprintf("%scount:cat:%s:%s", foodCachePrefix, categoryID, query)
}

// invalidateSearchCache removes all food:cache:* keys from Redis using SCAN + DEL.
func (r *CachedFoodRepository) invalidateSearchCache(ctx context.Context) {
	var cursor uint64
	pattern := foodCachePrefix + "*"
	for {
		keys, nextCursor, err := r.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}
		if len(keys) > 0 {
			r.rdb.Del(ctx, keys...)
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}

// --- Read operations (with cache-aside) ---

// Search returns active foods matching query, from cache if available.
func (r *CachedFoodRepository) Search(ctx context.Context, query string, limit, offset int32) ([]*entity.Food, error) {
	key := r.searchKey(query, limit, offset)

	if cached, err := r.rdb.Get(ctx, key).Bytes(); err == nil {
		var dtos []foodCacheDTO
		if json.Unmarshal(cached, &dtos) == nil {
			result := make([]*entity.Food, len(dtos))
			for i, d := range dtos {
				result[i] = cacheToFood(d)
			}
			return result, nil
		}
	}

	foods, err := r.inner.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	dtos := make([]foodCacheDTO, len(foods))
	for i, f := range foods {
		dtos[i] = foodToCache(f)
	}
	if data, mErr := json.Marshal(dtos); mErr == nil {
		r.rdb.Set(ctx, key, data, foodCacheTTL)
	}
	return foods, nil
}

// CountSearch returns the total count of active foods matching query, from cache if available.
func (r *CachedFoodRepository) CountSearch(ctx context.Context, query string) (int64, error) {
	key := r.countKey(query)

	if cached, err := r.rdb.Get(ctx, key).Int64(); err == nil {
		return cached, nil
	}

	count, err := r.inner.CountSearch(ctx, query)
	if err != nil {
		return 0, err
	}

	r.rdb.Set(ctx, key, count, foodCacheTTL)
	return count, nil
}

// SearchByCategory returns active foods in a category matching query, from cache if available.
func (r *CachedFoodRepository) SearchByCategory(ctx context.Context, categoryID uuid.UUID, query string, limit, offset int32) ([]*entity.Food, error) {
	key := r.searchByCategoryKey(categoryID, query, limit, offset)

	if cached, err := r.rdb.Get(ctx, key).Bytes(); err == nil {
		var dtos []foodCacheDTO
		if json.Unmarshal(cached, &dtos) == nil {
			result := make([]*entity.Food, len(dtos))
			for i, d := range dtos {
				result[i] = cacheToFood(d)
			}
			return result, nil
		}
	}

	foods, err := r.inner.SearchByCategory(ctx, categoryID, query, limit, offset)
	if err != nil {
		return nil, err
	}

	dtos := make([]foodCacheDTO, len(foods))
	for i, f := range foods {
		dtos[i] = foodToCache(f)
	}
	if data, mErr := json.Marshal(dtos); mErr == nil {
		r.rdb.Set(ctx, key, data, foodCacheTTL)
	}
	return foods, nil
}

// CountByCategory returns the total count of active foods in a category matching query, from cache if available.
func (r *CachedFoodRepository) CountByCategory(ctx context.Context, categoryID uuid.UUID, query string) (int64, error) {
	key := r.countByCategoryKey(categoryID, query)

	if cached, err := r.rdb.Get(ctx, key).Int64(); err == nil {
		return cached, nil
	}

	count, err := r.inner.CountByCategory(ctx, categoryID, query)
	if err != nil {
		return 0, err
	}

	r.rdb.Set(ctx, key, count, foodCacheTTL)
	return count, nil
}

// FindByID is a point-lookup — not cached (low read pressure, needs freshness).
func (r *CachedFoodRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Food, error) {
	return r.inner.FindByID(ctx, id)
}

// --- Write operations (invalidate search cache on success) ---

// Create persists a new food and invalidates the search cache.
func (r *CachedFoodRepository) Create(ctx context.Context, food *entity.Food) error {
	if err := r.inner.Create(ctx, food); err != nil {
		return err
	}
	r.invalidateSearchCache(ctx)
	return nil
}

// Update persists updated food fields and invalidates the search cache.
func (r *CachedFoodRepository) Update(ctx context.Context, food *entity.Food) error {
	if err := r.inner.Update(ctx, food); err != nil {
		return err
	}
	r.invalidateSearchCache(ctx)
	return nil
}

// Delete hard-deletes a food and invalidates the search cache.
func (r *CachedFoodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.inner.Delete(ctx, id); err != nil {
		return err
	}
	r.invalidateSearchCache(ctx)
	return nil
}

// Deactivate soft-deletes a food and invalidates the search cache.
func (r *CachedFoodRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	if err := r.inner.Deactivate(ctx, id); err != nil {
		return err
	}
	r.invalidateSearchCache(ctx)
	return nil
}
