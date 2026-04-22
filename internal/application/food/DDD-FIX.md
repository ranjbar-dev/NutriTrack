# DDD Fix Report: internal/application/food
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes (food-specific packages were already building after domain refactor)

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Food aggregate constructed via raw struct literal — no factory called | CRITICAL | food_service.go | PRE-FIXED | PRE-FIXED: already calls entity.NewFood() |
| 2 | Application service directly mutates Food aggregate exported fields | CRITICAL | food_service.go | PRE-FIXED | PRE-FIXED: already calls food.Update() domain method |
| 3 | Authorization policy for FoodCategoryService.Create delegated to handler | HIGH | food_category_service.go | DEFERRED | DEFERRED: handler refactor required |
| 4 | Role identity expressed as untyped string literals | LOW | food_service.go | DEFERRED | DEFERRED: cross-cutting change |

## Changes Applied

### No changes required — CRITICAL violations were already fixed

The CRITICAL violations listed in the audit were resolved during the preceding domain entity refactor (`internal/domain/food/DDD-FIX.md`). At the time this fix pass runs, `food_service.go` already:

- Calls `entity.NewFood(name, normalized, unit, ...)` factory (not a struct literal)
- Calls `food.Update(name, normalized, unit, ...)` domain method (not direct field mutation)
- Accesses aggregate state via getters: `food.IsActive()`, `food.CreatedBy()`, `food.ID()`

No code changes were needed for this package.

## Deferred Items
- **[HIGH]** `FoodCategoryService.Create` — authorization enforcement ("superadmin only") is currently in the HTTP handler. Moving it into the service requires adding a `callerRole` parameter and propagating it from the handler. Deferred to avoid handler signature churn in the current pass.
- **[LOW]** Role string literals (`"nutritionist"`, `"superadmin"`) — define `shared.Role` type and constants to get compile-time safety. Deferred as a cross-cutting change that affects multiple packages.

## Final Build Status
PASS — `go build ./...` (no changes made)
PASS — `go vet ./internal/...` (no changes made)

## Remaining Violations
- HIGH: Auth policy for FoodCategoryService.Create in handler layer — deferred
- LOW: Untyped role string literals — deferred
