# DDD Audit: internal/infrastructure/persistence
Layer: infrastructure
Audited: 2026-04-22
Files reviewed: 23 files across 11 sub-packages (sqlc/ excluded — auto-generated)

## Summary
- CRITICAL: 0
- HIGH: 1
- MEDIUM: 3
- LOW: 2
- PASS: 6 sub-packages (dietplan, food, medication, tracking, user, messaging)

---

## Findings

### [HIGH] Raw DB errors leak through domain boundary in 5 sub-packages

**Files:** `foodrequest/pg_food_request_repository.go:30,43,57,66`, `labresult/pg_lab_result_repository.go:41,56`, `message/pg_message_repository.go:48,62,67`, `notification/pg_notification_preference_repository.go:34,44`, `push/pg_push_subscription_repository.go:29,34,42`
**Issue:** Five sub-packages (`foodrequest`, `labresult`, `message`, `notification`, `push`) return raw `pgx` driver errors directly to callers, leaking infrastructure internals (connection errors, constraint violation messages, pgx sentinel types) through the repository boundary. The established pattern in `dietplan`, `food`, `medication`, `tracking`, and `user` correctly wraps all failures with `shared.ErrInternal`.
**DDD Rule:** `internal/infrastructure/` MUST NOT expose infrastructure details to outer layers. Errors crossing the repository boundary must be expressed in domain terms.
**Fix:** Wrap all non-not-found errors in `shared.ErrInternal`. Apply the same pattern to every `return err` in the 5 sub-packages.

---

### [MEDIUM] Constructor factory functions return concrete types instead of domain repository interfaces

**Files:** All `New*Repository(...)` constructors across all sub-packages.
**Issue:** All constructors return concrete pointer types (e.g., `*PgDietPlanRepository`) instead of the domain repository interface (e.g., `repository.DietPlanRepository`). Violates dependency inversion principle.
**DDD Rule:** Infrastructure factory functions should return the domain repository interface.
**Fix:** Change all `New*Repository` return types to the corresponding domain interface. Example: `func NewPgDietPlanRepository(pool *pgxpool.Pool) repository.DietPlanRepository`.

---

### [MEDIUM] Infrastructure directly mutates exported fields of domain aggregates after DB insert

**Files:** `dietplan/pg_diet_plan_repository.go:64–66, 280`, `food/pg_food_repository.go:52–55`, `foodrequest/pg_food_request_repository.go:34–37`, `tracking/pg_tracking_repository.go`
**Issue:** Infrastructure sets aggregate fields directly post-insert (e.g., `plan.ID = created.ID; plan.Status = ...; plan.CreatedAt = ...`). Bypasses domain encapsulation; aggregate can be in intermediate inconsistent state.
**DDD Rule:** Aggregates should have unexported fields and setter/factory methods. Infrastructure should call domain methods, not set fields directly.
**Fix (long-term):** Move domain aggregates toward unexported fields with controlled constructors. Short-term: batch all post-insert field assignments together; consider returning a new aggregate instance from `Create` operations.

---

### [MEDIUM] Missing `NewFromAggregate()` / reverse-mapper functions

**Files:** All `mapper.go` files across all sub-packages.
**Issue:** All `mapper.go` files only provide `toDomain()` (DB row → domain entity) conversion. Reverse mapping (domain entity → DB insert params) is performed inline inside each repository method by accessing aggregate exported fields directly, scattering mapping logic.
**DDD Rule:** Infrastructure SHOULD have centralised `NewFromAggregate()` / `ToAggregate()` conversion pattern.
**Fix:** Add reverse-mapping function to each `mapper.go` (e.g., `createParamsFromFood(f *entity.Food) db.CreateFoodParams`). Call it from repository methods instead of building params inline.

---

### [LOW] `isNotFound()` uses direct `==` comparison instead of `errors.Is()`

**Files:** `foodrequest/mapper.go:13`, `labresult/mapper.go:20`, `message/mapper.go:14`
**Issue:** `return err == pgx.ErrNoRows` — Go 1.13+ requires `errors.Is()` for wrapped error comparisons.
**Fix:** Replace `err == pgx.ErrNoRows` with `errors.Is(err, pgx.ErrNoRows)`.

---

### [LOW] No per-package internal DTO struct; sqlc types used directly as DB mapping layer

**Files:** All `mapper.go` files.
**Issue:** All packages use shared `db.*` types from `sqlc/` directly. A schema change regenerates all shared types simultaneously.
**Mitigation already in place:** `sqlc/` is isolated within `internal/infrastructure/persistence/` and does not leak into domain or application layers. Risk is low in practice.
**Fix (optional):** For larger aggregates that may diverge from DB schema, introduce package-local DTOs.

---

## Compliant Patterns Found

- **6 sub-packages fully compliant**: dietplan, food, medication, tracking, user, messaging. ✓
- **`isNotFound()` helper correctly used in most packages** for not-found disambiguation. ✓
- **No domain/application imports in infrastructure** — unidirectional dependencies maintained. ✓
- **`sqlc/` properly contained** within `persistence/` — no direct DB types leak to outer layers. ✓

## Fix Priority Order
1. **[HIGH]** Wrap raw DB errors in `shared.ErrInternal` across 5 sub-packages
2. **[MEDIUM]** Change constructor return types to domain repository interfaces
3. **[MEDIUM]** Extract aggregate-to-params mapper functions into `mapper.go` files
4. **[LOW]** Replace `err == pgx.ErrNoRows` with `errors.Is()` in 3 mapper files
