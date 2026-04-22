# DDD Fix Report: internal/application/foodrequest
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Concrete sibling service injected instead of interface | HIGH | food_request_service.go | SAFE | FIXED |
| 2 | Struct literal construction of domain entity | MEDIUM | food_request_service.go | N/A | PRE-FIXED |
| 3 | Hard-coded role string literals | MEDIUM | food_request_service.go | N/A | PRE-FIXED |
| 4 | Raw exported field access on User entity | LOW | food_request_service.go | N/A | PRE-FIXED |
| 5 | ApproveRequest DTO duplication | LOW | food_request_service.go | DEFERRED | DEFERRED |

## Changes Applied

### Fix 1: Define FoodCreator interface; decouple from *appFood.FoodService

**File:** `internal/application/foodrequest/food_request_service.go`

**Before:**
```go
type FoodRequestService struct {
    ...
    foodSvc  *appFood.FoodService
}
func NewFoodRequestService(..., foodSvc *appFood.FoodService) *FoodRequestService
```

**After:**
```go
// FoodCreator is the port through which FoodRequestService creates foods.
type FoodCreator interface {
    CreateFood(ctx context.Context, req appFood.CreateFoodRequest) (*foodEntity.Food, error)
}

type FoodRequestService struct {
    ...
    foodSvc  FoodCreator
}
func NewFoodRequestService(..., foodSvc FoodCreator) *FoodRequestService
```

`*appFood.FoodService` satisfies `FoodCreator` implicitly (structural typing), so `bootstrap/wire.go` required no changes — `foodSvc` is passed as-is and Go resolves the interface at compile time.

**Build:** PASS

## Pre-Fixed Findings (verified)

- **NewFoodRequest factory**: `Submit` already calls `frEntity.NewFoodRequest(clientID, *client.GetNutritionistID(), foodName)`. ✓
- **IsNutritionist()**: `ListPending` already calls `user.IsNutritionist()` instead of raw string comparison. ✓
- **GetNutritionistID()**: `Submit` already calls `client.GetNutritionistID()` (getter method, not direct field). ✓

## Deferred Items

- **[LOW]** `ApproveRequest` DTO rationalisation — cosmetic, deferred.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None at HIGH severity.
