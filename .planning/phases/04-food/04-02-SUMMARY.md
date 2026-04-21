---
phase: "04"
plan: "02"
subsystem: food
tags: [food-categories, many-to-many, category-filter, search, soft-delete]
dependency_graph:
  requires: [04-01]
  provides: [food-category-crud, category-filtered-search]
  affects: [food-search, food-handler, food-service, wire, router]
tech_stack:
  added: []
  patterns: [many-to-many-join-query, optional-filter-param, repository-extension]
key_files:
  created:
    - internal/application/food/food_category_service.go
    - internal/interfaces/http/handler/food_category_handler.go
  modified:
    - db/queries/foods.sql
    - internal/infrastructure/persistence/sqlc/foods.sql.go
    - internal/domain/food/entity/food.go
    - internal/domain/food/repository/food_repository.go
    - internal/infrastructure/persistence/food/mapper.go
    - internal/infrastructure/persistence/food/pg_food_repository.go
    - internal/application/food/food_service.go
    - internal/interfaces/http/handler/food_handler.go
    - bootstrap/wire.go
    - internal/interfaces/http/router/router.go
decisions:
  - FoodCategory.CreatedAt added to domain entity to support Create response
  - categoryID passed as *uuid.UUID so nil means no filter (backward-compatible)
  - adminGroup reused from existing router scope for superadmin category routes
metrics:
  duration: "~10 minutes"
  completed: "2026-04-21"
  tasks_completed: 9
  files_changed: 12
---

# Phase 04 Plan 02: Food Categories CRUD + Category Filter on Food Search Summary

**One-liner:** Food categories many-to-many CRUD with superadmin-gated write routes and optional category_id filter on paginated food search.

## What Was Built

### SQL Queries (Task 1)
Added two new queries to `db/queries/foods.sql`:
- `SearchFoodsByCategory :many` — JOIN foods with food_category_mappings, filter by `category_id`, apply pg_trgm similarity search, paginated
- `CountSearchFoodsByCategory :one` — same JOIN/filter without LIMIT/OFFSET for total count

`sqlc generate` produced corresponding Go functions in `foods.sql.go` with `SearchFoodsByCategoryParams` (CategoryID, Query, Off, Lim) and `CountSearchFoodsByCategoryParams` (CategoryID, Query).

### FoodRepository Interface (Task 2)
Added two methods to `internal/domain/food/repository/food_repository.go`:
- `SearchByCategory(ctx, categoryID, query, limit, offset) ([]*entity.Food, error)`
- `CountByCategory(ctx, categoryID, query) (int64, error)`

### PgFoodRepository (Task 3)
Implemented `SearchByCategory` and `CountByCategory` using the new sqlc functions, following the same pattern as existing `Search`/`CountSearch` methods.

### FoodService.SearchFoods (Task 4)
Signature updated to accept `categoryID *uuid.UUID`. When non-nil, routes to category-filtered repo calls; otherwise uses existing full-search path. NormalizePersian applied to query in both branches.

### FoodCategoryService (Task 5)
New service at `internal/application/food/food_category_service.go`:
- `Create(ctx, name)` — normalizes name with NormalizePersian, checks for conflict, delegates to categoryRepo
- `ListAll(ctx)` — returns all categories
- `Delete(ctx, id)` — verifies existence before deletion, returns ErrNotFound if missing

### FoodHandler.Search (Task 6)
Updated to parse optional `category_id` query parameter. Valid UUID → passes `&categoryID`; absent/empty → passes `nil`. Invalid UUID format → returns ErrValidation.

### FoodCategoryHandler (Task 7)
New handler at `internal/interfaces/http/handler/food_category_handler.go`:
- `POST /admin/food-categories` — requires name; calls Create; returns 201 with id, name, created_at
- `GET /food-categories` — lists all categories; accessible to any authenticated user
- `DELETE /admin/food-categories/:id` — parses UUID param; calls Delete; returns success message

### wire.go (Task 8)
Added `FoodCategoryService *appFood.FoodCategoryService` to Container struct and wired `catSvc := appFood.NewFoodCategoryService(pgCategoryRepo)`.

### router.go (Task 9)
Registered routes:
- `GET /api/v1/food-categories` — protected (any auth role) → `catHandler.ListAll`
- `POST /api/v1/admin/food-categories` — superadmin only → `catHandler.Create`
- `DELETE /api/v1/admin/food-categories/:id` — superadmin only → `catHandler.Delete`

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 2 - Missing Critical Functionality] Added CreatedAt to FoodCategory domain entity**
- **Found during:** Task 7 (FoodCategoryHandler Create response uses `cat.CreatedAt`)
- **Issue:** `entity.FoodCategory` only had `ID` and `Name`; no `CreatedAt`. The sqlc `FoodCategory` model has `CreatedAt time.Time`, and `categoryToDomain` mapper was discarding it.
- **Fix:** Added `CreatedAt time.Time` to `entity.FoodCategory`; updated `categoryToDomain` in mapper.go to populate it.
- **Files modified:** `internal/domain/food/entity/food.go`, `internal/infrastructure/persistence/food/mapper.go`
- **Commit:** 9391548

## Known Stubs

None — all endpoints return live data from the database.

## Self-Check: PASSED
