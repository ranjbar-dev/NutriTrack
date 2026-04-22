# DDD Fix Report: internal/domain/food
Layer: domain
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
FAIL — `go build ./...` before fixes (pre-existing user domain failures in application/auth, application/user, application/foodrequest, application/labresult, application/message, infrastructure/persistence/user — none related to food domain)

Food-specific packages: PASS — `go build ./internal/domain/food/...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Food aggregate exposes all fields as exported | CRITICAL | entity/food.go, all callers | SAFE | FIXED |
| 2 | No `NewFood()` factory function | HIGH | entity/food.go | SAFE | FIXED |
| 3 | No domain errors defined in package | HIGH | entity/errors.go (new) | SAFE | FIXED |
| 4 | `FoodCategory` has exported fields and no factory | HIGH | entity/food_category.go (new) | SAFE | FIXED |
| 5 | `FoodCategory` defined in same file as `Food` | HIGH | entity/food.go → entity/food_category.go | SAFE | FIXED |
| 6 | `valueobject/` package empty — nutrition as raw primitives | MEDIUM | entity/food.go | DEFERRED | DEFERRED |

## Changes Applied

### Fix 1 + 2: Unexport Food fields, add NewFood factory and ReconstructFood
**File:** `internal/domain/food/entity/food.go`
**Change:** Replaced exported struct fields with unexported fields; added `NewFood()` factory with name/unit validation; added `ReconstructFood()` for infrastructure hydration; added getters for all fields; added `Update()` domain method; added infrastructure-only setters `SetPersistedState()`, `SetUpdatedAt()`, `SetCategories()`.
**Before:**
```go
type Food struct {
    ID   uuid.UUID
    Name string
    // ... all exported
}
```
**After:**
```go
type Food struct {
    id   uuid.UUID
    name string
    // ... all unexported
}
func NewFood(name, nameNormalized, unit string, ...) (*Food, error) { ... }
func ReconstructFood(id uuid.UUID, ...) *Food { ... }
func (f *Food) ID() uuid.UUID { return f.id }
// ... getters, Update(), SetPersistedState(), SetUpdatedAt(), SetCategories()
```
**Build:** PASS

### Fix 3: Domain sentinel errors
**File:** `internal/domain/food/entity/errors.go` (new)
**Change:** Created new file with sentinel error variables.
```go
var (
    ErrFoodNotFound         = errors.New("food not found")
    ErrFoodNameRequired     = errors.New("food name is required")
    ErrFoodUnitRequired     = errors.New("food unit is required")
    ErrFoodInvalidNutrition = errors.New("nutrition values must be non-negative")
    ErrCategoryNameRequired = errors.New("category name is required")
    ErrCategoryNotFound     = errors.New("food category not found")
)
```
**Build:** PASS

### Fix 4 + 5: FoodCategory unexported fields, factory, moved to own file
**File:** `internal/domain/food/entity/food_category.go` (new)
**Change:** Extracted `FoodCategory` from `food.go` into dedicated file. Unexported all fields. Added `NewFoodCategory()` factory with name validation. Added `ReconstructFoodCategory()` for infrastructure. Added `ID()`, `Name()`, `CreatedAt()` getters.
**Build:** PASS

### Caller updates: infrastructure/persistence/food/mapper.go
**Change:** Updated `foodToDomain` to call `entity.ReconstructFood(...)` instead of struct literal; updated `categoryToDomain` to call `entity.ReconstructFoodCategory(...)`.
**Build:** PASS

### Caller updates: infrastructure/persistence/food/pg_food_repository.go
**Change:** Updated `Create` to use getters (`food.Name()`, etc.) and `food.SetPersistedState()` after INSERT; updated `FindByID` to use `food.SetCategories()` instead of direct assignment; updated `Update` to use getters and `food.SetUpdatedAt()` after UPDATE; updated category loops to use `cat.ID()`.
**Build:** PASS

### Caller updates: infrastructure/persistence/food/cache_dto.go (new)
**Change:** Added `foodCacheDTO` and `catCacheDTO` structs with json tags for Redis serialisation; added `foodToCache()` and `cacheToFood()` conversion functions. This is required because `json.Marshal` cannot serialise unexported fields on the aggregate.
**Build:** PASS

### Caller updates: infrastructure/persistence/food/cached_food_repository.go
**Change:** Updated `Search` and `SearchByCategory` to marshal/unmarshal `[]foodCacheDTO` instead of `[]*entity.Food`, using `foodToCache()` / `cacheToFood()` helpers. Prevents silent cache corruption caused by JSON serialising zero-value structs.
**Build:** PASS

### Caller updates: application/food/food_service.go
**Change:** `CreateFood` now calls `entity.NewFood(...)` factory instead of struct literal. `GetFood` uses `food.IsActive()` getter. `UpdateFood` uses `food.IsActive()`, `food.CreatedBy()` getters and calls `food.Update(...)` domain method instead of direct field mutation. `DeleteFood` uses `food.CreatedBy()` getter.
**Build:** PASS

### Caller updates: interfaces/http/handler/food_handler.go
**Change:** `toFoodResponse` updated to use getters: `f.ID()`, `f.Name()`, `f.Unit()`, `f.Calories()`, etc. Category iteration uses `c.ID()` and `c.Name()`. `f.CreatedBy()` replaces `f.CreatedBy`.
**Build:** PASS

### Caller updates: interfaces/http/handler/food_category_handler.go
**Change:** Handler now uses `cat.ID()`, `cat.Name()`, `cat.CreatedAt()` getters in both `Create` (response) and `ListAll` (response).
**Build:** PASS

### Caller updates: application/foodrequest/food_request_service.go
**Change:** Fixed `&food.ID` which became invalid (ID is now a method) → replaced with inline helper `func() *uuid.UUID { id := food.ID(); return &id }()`.
**Build:** PASS (food-specific error resolved; remaining errors are pre-existing user domain failures)

## Deferred Items
- **[MEDIUM] valueobject/nutrition.go** — Grouping the eight nutritional float64 fields into an immutable `Nutrition` value object would require changing the `NewFood` / `ReconstructFood` signatures and all callers. Deferred as it is a MEDIUM finding and exceeds safe 3-file change boundary without further scoping.

## Final Build Status
PASS — `go build ./internal/domain/food/...` after all fixes
PASS — `go build ./internal/infrastructure/persistence/food/...` after all fixes
PASS — `go build ./internal/application/food/...` after all fixes
PASS — `go vet ./internal/domain/food/... ./internal/infrastructure/persistence/food/... ./internal/application/food/...` after all fixes

**Overall `go build ./...`:** FAIL — same pre-existing user domain errors as baseline (NutritionistID, BirthDate, Role, IsActive, PasswordHash fields unexported without getters in domain/user). Zero new failures introduced by food domain fixes.

## Remaining Violations
- **[MEDIUM]** Nutritional value object — deferred (see above)
- **[LOW]** `NameNormalized` derived from `Name` — kept as explicit field for now; computed internally by `NewFood`/`Update` via caller-provided normalized string. Future improvement: compute within the aggregate using `shared.NormalizePersian`.
