---
phase: "04"
plan: "01"
subsystem: food
tags: [food, domain, ddd, pg_trgm, search, crud]
dependency_graph:
  requires: [01-foundation, 02-auth, 03-clients]
  provides: [food-aggregate, food-repository, food-service, food-handler, food-search]
  affects: [bootstrap/wire.go, router.go, sqlc/models.go]
tech_stack:
  added: [pg_trgm similarity search, food aggregate]
  patterns: [DDD aggregate, soft-delete, row-level ownership, pgtype.Numeric conversion]
key_files:
  created:
    - migrations/000003_foods.up.sql
    - migrations/000003_foods.down.sql
    - db/queries/foods.sql
    - db/queries/food_categories.sql
    - internal/infrastructure/persistence/sqlc/foods.sql.go
    - internal/infrastructure/persistence/sqlc/food_categories.sql.go
    - internal/domain/food/entity/food.go
    - internal/domain/food/repository/food_repository.go
    - internal/domain/food/repository/food_category_repository.go
    - internal/infrastructure/persistence/food/mapper.go
    - internal/infrastructure/persistence/food/pg_food_repository.go
    - internal/infrastructure/persistence/food/pg_food_category_repository.go
    - internal/application/food/food_service.go
    - internal/interfaces/http/handler/food_handler.go
  modified:
    - internal/infrastructure/persistence/sqlc/models.go
    - bootstrap/wire.go
    - internal/interfaces/http/router/router.go
decisions:
  - "Soft delete (is_active=false) for nutritionist-owned foods; hard DELETE for superadmin only"
  - "name_normalized column stores NormalizePersian output; gin_trgm_ops index on it for similarity search"
  - "Row-level ownership check in service layer (not handler) to enforce DDD boundary"
  - "Categories passed via food.Categories slice into Create/Update — replaced atomically in Update via RemoveFoodCategories + re-insert"
metrics:
  duration: "~15 min"
  completed: "2026-04-21"
  tasks_completed: 9
  files_created: 14
  files_modified: 3
---

# Phase 04 Plan 01: Food Domain Aggregate + pg_trgm Search + CRUD Handlers Summary

**One-liner:** Food DDD aggregate with pg_trgm similarity search on Persian normalized names, row-level ownership RBAC, and soft/hard delete by role.

## What Was Built

### Migration 000003
- `food_categories` — UUID PK, unique name
- `foods` — full macronutrient schema, `name_normalized` text column, `is_active` soft-delete flag, `created_by` FK to users
- `food_category_mappings` — join table for many-to-many
- `CREATE INDEX ... USING gin(name_normalized gin_trgm_ops)` — powers similarity search

### sqlc Queries
- `SearchFoods` — hybrid pg_trgm similarity + ILIKE fallback, ordered by similarity score then recency
- `CountSearchFoods` — same predicate for pagination total
- Full CRUD + `DeactivateFood` (soft delete) + category mapping operations

### Domain Layer
- `entity.Food` aggregate with `[]FoodCategory`, `*uuid.UUID CreatedBy`, `IsActive bool`
- `FoodRepository` + `FoodCategoryRepository` interfaces — no pgx/gin imports

### Infrastructure Layer
- `mapper.go` — `float64ToNumeric` / `numericToFloat64` helpers for `pgtype.Numeric`; `uuidToPgtypeUUID` for nullable FK; `foodToDomain` / `categoryToDomain`
- `PgFoodRepository` — Create (insert + category mappings), FindByID (food + categories), Update (atomic category replace), Delete, Deactivate, Search, CountSearch
- `PgFoodCategoryRepository` — CRUD for food categories

### Application Layer
- `FoodService.CreateFood` — normalizes name → NormalizePersian, resolves category IDs, creates food with createdBy
- `FoodService.GetFood` — hidden if `!IsActive`
- `FoodService.SearchFoods` — normalizes query before hitting pg_trgm
- `FoodService.UpdateFood` — nutritionist ownership check, superadmin bypass
- `FoodService.DeleteFood` — superadmin=hard delete, nutritionist=soft deactivate+ownership check, others=ErrForbidden

### HTTP Handler
- `FoodHandler` with Create/Search/GetOne/Update/Delete
- Uses `middleware.AuthUserIDKey` / `AuthUserRoleKey` constants (no magic strings)
- `toFoodResponse` includes all macro fields + categories array + nullable `created_by`
- Pagination via `dto.ParsePagination(c)` + `?q=` query param

### Wiring
- `bootstrap/wire.go` — `FoodService` added to Container, wired with `pgFoodRepo` + `pgCategoryRepo`
- `router.go` — 5 food routes registered in `protected` group (no additional role middleware — RBAC in service)

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None — all food fields are wired end-to-end from DB → domain → handler response.

## Self-Check: PASSED

**Files verified:**
- FOUND: migrations/000003_foods.up.sql
- FOUND: migrations/000003_foods.down.sql
- FOUND: db/queries/foods.sql
- FOUND: db/queries/food_categories.sql
- FOUND: internal/infrastructure/persistence/sqlc/foods.sql.go
- FOUND: internal/infrastructure/persistence/sqlc/food_categories.sql.go
- FOUND: internal/domain/food/entity/food.go
- FOUND: internal/domain/food/repository/food_repository.go
- FOUND: internal/domain/food/repository/food_category_repository.go
- FOUND: internal/infrastructure/persistence/food/mapper.go
- FOUND: internal/infrastructure/persistence/food/pg_food_repository.go
- FOUND: internal/infrastructure/persistence/food/pg_food_category_repository.go
- FOUND: internal/application/food/food_service.go
- FOUND: internal/interfaces/http/handler/food_handler.go

**Commit verified:** 0a71f49 — feat(04-01): food domain aggregate, pg_trgm search, CRUD handlers

**Build:** `go build ./...` exits 0 — no compile errors.
