# DDD Fix Report: internal/interfaces/http
Layer: interfaces
Fixed: 2025-07-14
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | avatar_handler.go reads raw string context keys instead of typed constants | HIGH (SECURITY) | handler/avatar_handler.go | SAFE | FIXED |
| 2 | middleware/rate_limit.go accepts `*redis.Client` directly — no interface abstraction | HIGH | middleware/rate_limit.go, infrastructure/redis/rate_limiter.go | SAFE | FIXED |
| 3 | router/router.go accepts raw `*pgxpool.Pool` and `*redis.Client` infrastructure types | HIGH | router/router.go, bootstrap/wire.go, cmd/server/main.go | SAFE | FIXED |
| 4 | Nine handlers import domain entity packages for response mapping | HIGH | handler/*.go, application/*/response.go (new files) | SAFE | FIXED |

## Changes Applied

### Fix 1: avatar_handler.go — wrong gin context keys
**File:** `handler/avatar_handler.go`  
**Change:** Added `middleware` import; replaced raw string keys `"userID"` and `"userRole"` with `middleware.AuthUserIDKey` and `middleware.AuthUserRoleKey`.  
**Before:**
```go
callerIDVal, _ := c.Get("userID")
callerRoleVal, _ := c.Get("userRole")
```
**After:**
```go
callerIDVal, _ := c.Get(middleware.AuthUserIDKey)
callerRoleVal, _ := c.Get(middleware.AuthUserRoleKey)
```
**Build:** PASS

---

### Fix 2: RateLimiter interface abstraction
**Files:**
- `middleware/rate_limit.go` — removed `*redis.Client` dependency, defined local `RateLimiter` interface
- `internal/infrastructure/redis/rate_limiter.go` — new file: `RedisRateLimiter` struct implementing the interface using Incr+Expire pattern

**Middleware interface defined:**
```go
type RateLimiter interface {
    Allow(ctx context.Context, key string, max int64, window time.Duration) (bool, error)
}
func RateLimitByIP(limiter RateLimiter, max int64) gin.HandlerFunc
```

**Infrastructure implementation:**
```go
type RedisRateLimiter struct { client *redis.Client }
func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter
func (r *RedisRateLimiter) Allow(ctx, key, max, window) (bool, error)
```
**Build:** PASS

---

### Fix 3: Router accepts *bootstrap.Container instead of raw infra types
**Files:**
- `router/router.go` — signature changed to `func New(container *bootstrap.Container)`; removed `*pgxpool.Pool`, `*redis.Client`, `*configs.Config` parameters; uses `container.Cfg` and `container.RateLimiter`
- `bootstrap/wire.go` — added `Cfg *configs.Config` and `RateLimiter *redisInfra.RedisRateLimiter` to `Container` struct; wires `rateLimiter` in `NewContainer`
- `cmd/server/main.go` — updated call to `router.New(bootstrap.NewContainer(db, rdb, cfg))`

**Before:**
```go
func New(db *pgxpool.Pool, rdb *redis.Client, cfg *configs.Config)
```
**After:**
```go
func New(container *bootstrap.Container)
```
**Build:** PASS

---

### Fix 4: Move response mappers from handlers to application layer

Nine handlers were importing domain entity packages to build JSON response maps. All mapping logic was extracted into `response.go` files in the corresponding application package, and handler entity imports were removed.

#### New application-layer response files created:

| File | Exported Functions |
|------|--------------------|
| `internal/application/user/response.go` | `MapClientResponse`, `MapNutritionistResponse` |
| `internal/application/food/response.go` | `MapFoodResponse` |
| `internal/application/medication/response.go` | `MapMedicationResponse` |
| `internal/application/foodrequest/response.go` | `MapFoodRequestResponse` |
| `internal/application/labresult/response.go` | `MapLabResultResponse` |
| `internal/application/message/response.go` | `MapMessageResponse` |
| `internal/application/tracking/response.go` | `MapFoodLog`, `MapWaterLog`, `MapSleepLog`, `MapExerciseLog`, `MapMedicationLog`, `MapBodyMeasurement` |
| `internal/application/dietplan/response.go` | `MapDietPlan`, `MapDietPlanFull` (with unexported helpers for nested structs) |

#### Handler files updated (entity imports removed, local mapper functions deleted):

| Handler | Entity import removed | Local function removed |
|---------|----------------------|----------------------|
| `handler/client_handler.go` | `domain/user/entity` | `toClientResponse` |
| `handler/nutritionist_handler.go` | `domain/user/entity` | `toNutritionistResponse` |
| `handler/food_handler.go` | `domain/food/entity` | `toFoodResponse` |
| `handler/medication_handler.go` | `domain/medication/entity` | `toMedicationResponse` |
| `handler/food_request_handler.go` | `domain/foodrequest/entity` | `foodRequestToMap` |
| `handler/lab_result_handler.go` | `domain/labresult/entity` | `labResultToMap` |
| `handler/message_handler.go` | `domain/message/entity` | `messageToMap` |
| `handler/tracking_handler.go` | `domain/tracking/entity` | `foodLogToMap`, `waterLogToMap`, `sleepLogToMap`, `exerciseLogToMap`, `medicationLogToMap`, `bodyMeasurementToMap` |
| `handler/diet_plan_handler.go` | `domain/dietplan/entity` | `planToMap`, `planFullToMap`, `summaryToMap`, `rangeToMap`, `itemsToSlice`, `exercisesToSlice`, `prescriptionsToSlice` |

**Build:** PASS

---

## Deferred Items

None — all HIGH findings from DDD-AUDIT.md have been resolved.

## Final Build Status
PASS — `go build ./...` after all fixes  
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None — all CRITICAL and HIGH findings resolved.
