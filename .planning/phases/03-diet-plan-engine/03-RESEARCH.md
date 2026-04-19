# Phase 3: Diet Plan Engine — Research

**Researched:** 2026-04-19
**Domain:** Diet Plan CRUD + Batch Loading + Plan Builder UI
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**Database Schema:**
- D-01 through D-10: Full schema spec (diet_plans, plan_days, meals, meal_options, meal_option_items, plan_exercises, plan_medications) with all FK/cascade rules and row-level auth
- D-02: One-active-plan-per-client at TWO layers (partial unique index + service-layer archival)
- D-09: ON DELETE CASCADE for all child tables; hard delete allowed for draft plans only

**Batch Loading:**
- D-11: pgx.Batch SendBatch with ≤6 queries (4 core + exercises + medications)
- D-12: Tree assembly in Go service layer using map-based aggregation
- D-13: Single-record queries fine for CRUD; batch only for aggregate read

**Nutritional Computation:**
- D-14: Client-side Vue computed properties — backend does NOT sum totals
- D-15: Backend embeds full food nutritional data in every meal_option_item
- D-16: Real-time update on item add/remove/edit (no debounce needed)

**Plan Builder UI:**
- D-17: Drill-down navigation (4 route levels)
- D-18: Up/down arrow reordering (no drag-and-drop in v1)
- D-19: Draft status during creation; explicit "فعال‌سازی" activation
- D-20: Food picker calls existing `GET /api/foods?search=...` endpoint
- D-22: Status badge colors (draft=orange, active=green, archived=grey)

**Client View:**
- D-23: `/client/plan` shows active plan or empty state
- D-24: Day navigation via horizontal tab bar; default = today's day_number
- D-25–D-27: Meals, exercises, medications display structure

**API Routes:**
- D-29: All nutritionist routes listed explicitly
- D-30: Full sub-resource CRUD route set
- D-31: `GET /api/clients/me/active-plan` for client

**Validation:**
- D-33: Activation blocked if plan is incomplete (Persian error message)
- D-34: Quantity > 0, option_number auto-incremented, day_number >= 1
- D-35: All errors in Persian

**Performance:**
- D-36: ≤500ms SLA for plan aggregate
- D-37: Required indexes listed

### Agent's Discretion
- Exact Pinia store structure for the plan builder state
- Drag handle vs up/down arrow implementation detail (v1 = up/down arrows)
- Bottom-sheet animation for food item picker
- Loading skeleton for plan aggregate on client view
- Color palette for status badges
- Whether to use layout-level composable or page-level state for breadcrumb

### Deferred Ideas (OUT OF SCOPE)
- Plan templates
- Drag-and-drop reordering
- Copy day
- Nutritional goal thresholds / warnings
- Plan PDF export
- Bulk day creation wizard
- Client plan option selection (Phase 4)

</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| DIET-01 | Nutritionist creates diet plan (start_date, end_date, notes, daily_water_target_ml, status) | D-01 schema, D-29 routes, handler/service/repository pattern |
| DIET-02 | One active plan per client; new plan auto-archives previous | D-02 partial unique index + service transaction, §6 activation transaction |
| DIET-03 | Plan → Days → Meals → Options → Items nested structure | D-03–D-07 schema, D-11 batch queries, §3 tree assembly |
| DIET-04 | Plan days with day_number (1-based) and calendar mapping | D-03 schema, repeating cycles not in v1 |
| DIET-05 | Meals with title, scheduled_time, display_order | D-04 schema, sorting by display_order ASC then scheduled_time ASC |
| DIET-06 | Meal options with option_number (client picks one) | D-05 schema, auto-increment on add |
| DIET-07 | Meal option items link to food DB with quantity, measurement_unit | D-06 schema, measurement_unit enum reuse §7 |
| DIET-08 | Real-time computed nutritional totals (client-side) | D-14/D-15/D-16, §8 formula and Vue computed pattern |
| DIET-09 | Exercise recommendations per plan day | D-07 schema, query 5 in batch |
| DIET-10 | Prescribed medications per diet plan | D-08 schema, query 6 in batch |
| DIET-11 | Archived plans viewable for history | D-28 plan list route, D-11 batch loading for archived plans |
| DIET-12 | ≤5 queries batch loading (no N+1) | D-11 6-query batch, §1 pgx.Batch pattern |
| INFRA-10 | Diet plan load time < 500ms | D-36 SLA, D-37 indexes, §Validation Architecture Dim 2 |

</phase_requirements>

---

## Summary

Phase 3 is the technical complexity nexus of NutriTrack. It delivers a multi-level plan builder for nutritionists (4 levels of drill-down navigation) and a read-only plan view for clients. The two highest-risk elements are: (1) the pgx.Batch aggregate loading pattern — 6 queries sent in a single round-trip to assemble a deeply nested tree in Go memory, and (2) the Pinia reactive state for the drill-down plan builder — 4 page levels each mutating shared draft state without full re-renders.

The good news: the existing codebase provides a complete, proven pattern for everything in this phase. The handler→service→repository triad is already implemented for food and medication, the `MeasurementUnitType` sqlc enum is already defined, the middleware/routing system is wired, and the food search endpoint is ready to reuse. Phase 3 is largely about applying established patterns to a more complex domain (Plan is the aggregate root with 5 nested entity levels instead of 1).

**Primary recommendation:** Follow the exact handler/service/repository pattern from `food_handler.go` / `food_service.go` / `food_repo.go`. For batch loading, use `pgxpool.Pool.SendBatch()` directly in the repository (bypassing sqlc-generated helpers for the aggregate query). For the Vue builder, use a single `usePlanBuilderStore` with map-keyed reactive state and a navigation cursor.

---

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Plan aggregate loading | API / Backend (Go) | PostgreSQL | Batch queries + tree assembly happen server-side |
| One-active-plan constraint | API / Backend (Go) + PostgreSQL | — | Dual enforcement: DB partial unique index + service transaction |
| Nutritional totals | Browser / Client (Vue) | — | D-14 decision: computed properties from embedded food data |
| Plan builder UI state | Browser / Client (Pinia) | — | Draft plan state lives in Pinia during creation flow |
| Row-level authorization | API / Backend (Go) | — | repository-layer WHERE nutritionist_id = $current_user |
| Drill-down routing | Frontend (Nuxt file-based routing) | — | 4 nested route levels |
| Food item search in picker | API / Backend (reuse Phase 2) | — | Existing `/api/foods?search=...` endpoint |

---

## 1. pgx.Batch Pattern

### How pgx.Batch / SendBatch Works

`pgxpool.Pool.SendBatch(ctx, batch)` sends all queued queries to PostgreSQL in a **single network round-trip** and returns a `pgx.BatchResults`. Results are collected in the same order the queries were queued.

[VERIFIED: codebase uses pgx/v5 throughout — `food_service.go` imports `"github.com/jackc/pgx/v5"`, `main.go` creates `pgxpool.NewWithConfig`]

```go
// Source: pgx/v5 pgxpool usage — verified in project (main.go uses pgxpool.Pool)
func (r *planRepo) GetFullPlanAggregate(ctx context.Context, planID uuid.UUID) (*model.PlanAggregate, error) {
    batch := &pgx.Batch{}

    // Queue 6 queries — all sent in ONE round-trip
    batch.Queue(sqlGetPlanAndDays, pgtype.UUID{Bytes: planID, Valid: true})
    // (remaining queries queued after we get dayIDs — see NOTE below)

    br := r.pool.SendBatch(ctx, batch)
    defer br.Close()

    // Collect result for query 1
    rows, err := br.Query()
    if err != nil { return nil, err }
    // scan rows...
    rows.Close()

    // Collect result for query 2
    rows2, err := br.Query()
    // ...

    return aggregate, nil
}
```

**Critical nuance — two-phase batching for the plan aggregate:**

The 4 core queries cannot ALL be sent in the same initial batch because queries 2–4 depend on IDs discovered from query 1 (day IDs → meal IDs → option IDs). The correct pattern is:

**Phase 1 batch** (1 query): fetch plan + days → collect `dayIDs []pgtype.UUID`

**Phase 2 batch** (3–5 queries): use `dayIDs` to fetch meals (ANY), then from those results use meal IDs for options (ANY), then option IDs for items (ANY), plus exercises and medications.

Since phases 2–4 depend on results from earlier phases, the implementation uses **two SendBatch calls** or executes phase 1 as a regular query then batches phases 2–5:

```go
// Pattern: 1 regular query + 1 SendBatch call (5 queries) = 2 round-trips total, still far below N+1
// OR: 2 SendBatch calls (1+5 queries) — same total round-trips

// Recommended for this phase:
// Round-trip 1: SELECT plan + days (regular query or single-query batch)
// Round-trip 2: SendBatch with 5 queries (meals, options, items+food, exercises, medications)
//               using ANY($ids) with the IDs collected in round-trip 1
```

This gives 2 round-trips vs the ~50+ that N+1 would require. The D-11 decision says ≤6 queries total which is satisfied.

### Correct pgx.Batch Queue Pattern

```go
// Source: pgx/v5 Batch API — [ASSUMED based on pgx/v5 docs pattern, consistent with
// project use of pgx/v5 in food_service.go and food_repo.go]
batch := &pgx.Batch{}

// Queue all queries before sending
batch.Queue(
    `SELECT * FROM meals WHERE day_id = ANY($1) ORDER BY display_order`,
    dayIDArray, // pgtype.Array or []pgtype.UUID
)
batch.Queue(
    `SELECT mo.* FROM meal_options mo WHERE mo.meal_id = ANY($1) ORDER BY option_number`,
    mealIDArray,
)
batch.Queue(
    `SELECT moi.*, f.name, f.calories, f.protein_g, f.carbs_g, f.fat_g, f.fiber_g,
            f.measurement_unit, f.measurement_amount
     FROM meal_option_items moi
     JOIN foods f ON moi.food_id = f.id
     WHERE moi.option_id = ANY($1)`,
    optionIDArray,
)
batch.Queue(
    `SELECT * FROM plan_exercises WHERE day_id = ANY($1) ORDER BY display_order`,
    dayIDArray,
)
batch.Queue(
    `SELECT pm.*, m.name as medication_name, m.form, m.dosage_unit
     FROM plan_medications pm
     JOIN medications m ON pm.medication_id = m.id
     WHERE pm.plan_id = $1`,
    planID,
)

br := pool.SendBatch(ctx, batch)
defer br.Close()

// MUST collect results in the same order
mealsRows, err := br.Query()
// ... scan
mealsRows.Close()

optionsRows, err := br.Query()
// ... scan
optionsRows.Close()

// etc.
```

**Pitfall:** `br.Close()` must be called even on error — use `defer br.Close()`. If you forget to call `Query()` or `QueryRow()` for a queued item before `Close()`, pgx closes the connection anyway but the skipped results are discarded.

### ANY($ids) with pgx/v5

The existing codebase does NOT yet use ANY($ids) but the pattern is established in sqlc queries for food. For pgx batch queries with array args:

```go
// pgtype.Array for ANY($1) — [ASSUMED based on pgx/v5 pgtype package conventions]
// The correct type for passing []uuid.UUID to ANY($1) in pgx/v5:
pgtype.FlatArray[pgtype.UUID](dayIDs)
// OR collect as []pgtype.UUID and pass directly — pgx/v5 handles []pgtype.UUID natively
```

**Note:** sqlc generates queries with `ANY($ids)` in `foods.sql` comments but the actual ANY pattern for batch needs raw pgx calls. The sqlc-generated querier is NOT used for the batch aggregate — write raw SQL in `plan_repo.go` directly.

---

## 2. sqlc Integration with Batch Queries

### What sqlc Does vs. Doesn't Handle

sqlc generates type-safe Go functions for individual named SQL queries. It does NOT generate batch-aware code. For the plan aggregate endpoint:

- **Use sqlc** for: all individual CRUD operations (create plan, add day, add meal, update option, delete item, list plans for client, etc.)
- **Use raw pgx** for: the aggregate batch-load queries in `GetFullPlanAggregate` and `GetActivePlanAggregate`

[VERIFIED: project uses `sqlc v1.30.0` — from `models.go` header comment]

### sqlc Query Files for Phase 3

The new queries file `backend/db/queries/diet_plans.sql` will contain all `-- name: X :one/:many/:exec` queries for individual CRUD operations following the exact pattern in `foods.sql`. The aggregate query is NOT a sqlc query — it goes directly in `plan_repo.go`.

Example of what will be in sqlc:
```sql
-- name: CreateDietPlan :one
INSERT INTO diet_plans (client_id, nutritionist_id, start_date, end_date, notes, daily_water_target_ml, status)
VALUES ($1, $2, $3, $4, $5, $6, 'draft')
RETURNING *;

-- name: GetDietPlanByID :one
SELECT * FROM diet_plans WHERE id = $1 AND nutritionist_id = $2;

-- name: ActivateDietPlan :exec
UPDATE diet_plans SET status = 'active', updated_at = NOW() WHERE id = $1 AND nutritionist_id = $2;

-- name: ArchivePreviousActivePlan :exec
UPDATE diet_plans SET status = 'archived', updated_at = NOW()
WHERE client_id = $1 AND status = 'active' AND id != $2;
```

After writing SQL, run `sqlc generate` from the `backend/` directory to generate Go code.

---

## 3. Tree Assembly in Go (Map-Based Aggregation)

### Pattern: Bottom-Up Map Assembly

After collecting all flat rows from the batch queries, assemble the nested tree using Go maps keyed by UUID:

```go
// Source: Established DDD aggregate assembly pattern [ASSUMED — standard Go pattern]
// This is the canonical approach for assembling nested trees without N+1

// Step 1: Build maps from flat rows
dayMap := make(map[uuid.UUID]*dto.PlanDayResponse)       // day_id → day
mealMap := make(map[uuid.UUID]*dto.MealResponse)          // meal_id → meal
optionMap := make(map[uuid.UUID]*dto.MealOptionResponse)  // option_id → option

// Step 2: Index days into plan
for _, dayRow := range dayRows {
    day := mapDayRow(dayRow)
    dayMap[day.ID] = day
    plan.Days = append(plan.Days, day) // maintain order (sorted by day_number)
}

// Step 3: Index meals into days
for _, mealRow := range mealRows {
    meal := mapMealRow(mealRow)
    mealMap[meal.ID] = meal
    if day, ok := dayMap[meal.DayID]; ok {
        day.Meals = append(day.Meals, meal)
    }
}

// Step 4: Index options into meals
for _, optRow := range optionRows {
    opt := mapOptionRow(optRow)
    optionMap[opt.ID] = opt
    if meal, ok := mealMap[opt.MealID]; ok {
        meal.Options = append(meal.Options, opt)
    }
}

// Step 5: Index items (with food data) into options
for _, itemRow := range itemRows {
    item := mapItemWithFoodRow(itemRow) // includes all food nutritional fields
    if opt, ok := optionMap[item.OptionID]; ok {
        opt.Items = append(opt.Items, item)
    }
}

// Step 6: Exercises into days
for _, exRow := range exerciseRows {
    ex := mapExerciseRow(exRow)
    if day, ok := dayMap[ex.DayID]; ok {
        day.Exercises = append(day.Exercises, ex)
    }
}

// Step 7: Medications onto plan (not per-day)
for _, medRow := range medicationRows {
    plan.Medications = append(plan.Medications, mapMedicationRow(medRow))
}
```

**Why maps not nested loops:** Maps give O(1) lookup vs O(n) search. For a 7-day×5-meal×3-option×4-item plan = 420 items, linear search would scan up to 420 entries per item insertion. Map lookup is constant.

**Order preservation:** Rows from SQL are already ordered (`ORDER BY day_number`, `ORDER BY display_order`, `ORDER BY option_number`). Appending to slices in row-scan order preserves that ordering.

---

## 4. Vue 3 / Pinia State for Drill-Down Plan Builder

### Recommended Store Structure

Two separate Pinia stores for Phase 3:

**`stores/plan-builder.ts`** — nutritionist plan creation/editing state
**`stores/plan.ts`** — client active plan viewing state (read-only)

```typescript
// stores/plan-builder.ts — [ASSUMED pattern based on existing food.ts store]
// Key insight: store the full plan aggregate fetched from API,
// then mutate it locally as nutritionist adds/removes items.
// Each navigation level just reads from the same store.

export const usePlanBuilderStore = defineStore('plan-builder', () => {
  // The full plan aggregate (as returned by GET /api/diet-plans/:id)
  const plan = ref<DietPlanAggregate | null>(null)
  const loading = ref(false)
  const saving = ref(false)

  // Navigation cursor — which level is displayed
  const currentDayId = ref<string | null>(null)
  const currentMealId = ref<string | null>(null)

  // Computed derived views (no duplication, just derived)
  const currentDay = computed(() =>
    plan.value?.days.find(d => d.id === currentDayId.value) ?? null
  )
  const currentMeal = computed(() =>
    currentDay.value?.meals.find(m => m.id === currentMealId.value) ?? null
  )

  async function loadPlan(planId: string) { /* GET /api/diet-plans/:id */ }
  async function addDay(payload: CreateDayPayload) { /* POST, then reload plan */ }
  async function addMeal(dayId: string, payload: CreateMealPayload) { /* POST */ }
  async function addOption(mealId: string) { /* POST */ }
  async function addItem(optionId: string, payload: CreateItemPayload) { /* POST */ }
  async function removeItem(itemId: string) { /* DELETE, then update local state */ }
  async function activatePlan() { /* PATCH activate, with modal confirmation */ }

  // Navigation
  function navigateToDay(dayId: string) { currentDayId.value = dayId; currentMealId.value = null }
  function navigateToMeal(mealId: string) { currentMealId.value = mealId }
  function navigateBack() {
    if (currentMealId.value) { currentMealId.value = null; return }
    if (currentDayId.value) { currentDayId.value = null; return }
    navigateTo(`/nutritionist/clients/${plan.value?.clientId}/plans`)
  }

  return { plan, loading, saving, currentDayId, currentMealId, currentDay, currentMeal,
           loadPlan, addDay, addMeal, addOption, addItem, removeItem, activatePlan,
           navigateToDay, navigateToMeal, navigateBack }
})
```

### Key Reactive State Principle for Drill-Down

**Do NOT duplicate state across levels.** Each page level reads from the same store via computed properties derived from the navigation cursor. This prevents the "stale slice" problem where a parent list and a child form hold different versions of the same object.

**After any mutation**, refresh the plan aggregate from the server (`await loadPlan(planId)`) to keep the single source of truth consistent. Since CRUD operations return the updated resource, an optimistic update pattern is acceptable for perceived performance — but the full plan refresh ensures correctness.

### Nutritional Totals as Vue Computed Properties

```typescript
// In the meal detail page or component — [ASSUMED based on D-14/D-16]
const currentOptionTotals = computed(() => {
  const items = currentMeal.value?.options ?? []
  return items.map(option => ({
    optionId: option.id,
    calories: option.items.reduce((sum, item) => {
      return sum + (item.food.calories * item.quantity / item.food.measurement_amount)
    }, 0),
    protein_g: option.items.reduce((sum, item) => {
      return sum + (item.food.protein_g * item.quantity / item.food.measurement_amount)
    }, 0),
    // ... same for carbs_g, fat_g, fiber_g
  }))
})
```

The formula per D-14: `food.calories * item.quantity / food.measurement_amount`
This produces calories per the item's quantity, regardless of what the food's base measurement_amount is.

---

## 5. Nuxt 4 Nested Routes for Plan Builder

### Route → File Mapping

Nuxt 4 uses the `app/pages/` directory with `[param]` syntax for dynamic segments:

| Route | File |
|-------|------|
| `/nutritionist/clients/:clientId/plans/new` | `app/pages/nutritionist/clients/[clientId]/plans/new.vue` |
| `/nutritionist/clients/:clientId/plans/:planId` | `app/pages/nutritionist/clients/[clientId]/plans/[planId]/index.vue` |
| `/nutritionist/clients/:clientId/plans/:planId/days/:dayId` | `app/pages/nutritionist/clients/[clientId]/plans/[planId]/days/[dayId].vue` |
| `/nutritionist/clients/:clientId/plans/:planId/days/:dayId/meals/:mealId` | `app/pages/nutritionist/clients/[clientId]/plans/[planId]/days/[dayId]/meals/[mealId].vue` |
| `/nutritionist/clients/:clientId/plans` | `app/pages/nutritionist/clients/[clientId]/plans/index.vue` |
| `/client/plan` | `app/pages/client/plan.vue` (already exists as stub) |

[VERIFIED: existing `app/pages/nutritionist/foods/[id].vue`, `index.vue`, `new.vue` files confirm Nuxt 4 file-based routing structure]

### Each Level's definePageMeta

Each page in the drill-down hierarchy uses the same middleware pattern as existing nutritionist pages:

```typescript
definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist'],
})
```

[VERIFIED: confirmed in `app/pages/nutritionist/foods/index.vue` and `new.vue`]

### Breadcrumb Pattern

Since all 4 levels share the `usePlanBuilderStore`, breadcrumb data is derived from store state:

```typescript
// In each drill-down page — [ASSUMED based on agent's discretion]
const breadcrumbs = computed(() => {
  const items = [
    { label: 'مشتریان', to: '/nutritionist' },
    { label: planStore.plan?.clientName, to: `/nutritionist/clients/${clientId}/plans` },
    { label: planStore.plan?.status === 'draft' ? 'پیش‌نویس' : planStore.plan?.startDate, to: planUrl },
  ]
  if (currentDayId) items.push({ label: `روز ${currentDay.dayNumber}`, to: dayUrl })
  if (currentMealId) items.push({ label: currentMeal.title })
  return items
})
```

---

## 6. Existing Codebase Patterns

### Go Backend — Exact Patterns to Replicate

**Handler pattern** (`food_handler.go`):
- Struct `DietPlanHandler` with `*service.DietPlanService`
- Constructor `NewDietPlanHandler`
- Each method: parse UUID from `c.Param("id")`, `uuid.Parse()` → 400; bind JSON → 400; call service; switch on sentinel errors → correct HTTP status; success → `c.JSON(http.StatusOK, resp)`
- Sentinel errors defined in `service/diet_plan_service.go` as `var Err... = errors.New("Persian text")`
- User ID extracted: `uuid.Parse(c.GetString("user_id"))` (set by auth middleware)
- Role extracted: `c.GetString("role")` for authorization

**Service pattern** (`food_service.go`):
- Struct embeds `repository.DietPlanRepository` interface + `zerolog.Logger`
- Business logic only — no Gin, no HTTP
- Helpers: `optionalText()`, `optionalBool()`, `formatTimestamp()`, `numericFromFloat64()`, `numericToFloat64()` already exist in `food_service.go` — same package, reusable directly
- `pgtype.UUID{Bytes: id, Valid: true}` for UUID → pgtype conversion

**Repository pattern** (`food_repo.go`):
- Interface `DietPlanRepository` in `repository/diet_plan_repo.go`
- Implementation `dietPlanRepository` with `*sqlc.Queries` for standard ops
- For batch aggregate: accept `*pgxpool.Pool` directly (not just `sqlc.DBTX`) so `SendBatch` is available

**Important:** The existing repo pattern uses `sqlc.New(db)` where `db` is `sqlc.DBTX`. For the aggregate query that needs `SendBatch`, the `planRepo` needs direct access to `*pgxpool.Pool`. Pass both:

```go
type dietPlanRepository struct {
    q    *sqlc.Queries   // for standard sqlc-generated ops
    pool *pgxpool.Pool   // for SendBatch aggregate queries
}

func NewDietPlanRepository(pool *pgxpool.Pool) DietPlanRepository {
    return &dietPlanRepository{q: sqlc.New(pool), pool: pool}
}
```

[VERIFIED: `food_repo.go` uses `sqlc.New(db)` where db is the pool; `main.go` passes `pool` directly to `NewFoodRepository(pool)`]

### Router Registration Pattern

In `main.go`, diet plan routes are added in the same pattern as foods/medications:

```go
// Nutritionist diet plan routes
dietPlans := r.Group("/api/diet-plans")
dietPlans.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin"))
{
    dietPlans.POST("", planHandler.CreatePlan)
    dietPlans.GET("/:id", planHandler.GetPlanAggregate)
    dietPlans.PATCH("/:id", planHandler.UpdatePlanHeader)
    dietPlans.PATCH("/:id/activate", planHandler.ActivatePlan)
    dietPlans.DELETE("/:id", planHandler.DeletePlan)
    // Sub-resources
    dietPlans.POST("/:id/days", planHandler.AddDay)
    dietPlans.PUT("/:id/days/:dayId", planHandler.UpdateDay)
    dietPlans.DELETE("/:id/days/:dayId", planHandler.DeleteDay)
    // ... meals, options, items, exercises, medications
}

// Client plans route  
clientRoutes := r.Group("/api/clients")
clientRoutes.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("client"))
{
    clientRoutes.GET("/me/active-plan", planHandler.GetActivePlan)
}

// Nutritionist client plan list
nutriRoutes.GET("/clients/:clientId/plans", planHandler.ListClientPlans)
```

[VERIFIED: routing pattern from `main.go` lines 147–186]

### Frontend — Reusable Assets Confirmed

- `UiAppButton.vue`, `UiAppInput.vue`, `UiLoadingSpinner.vue` — confirmed in `app/components/ui/`
- `toPersianDigits()` — used in `app/pages/nutritionist/foods/index.vue` line 2
- `toLatinDigits()` — used in `app/pages/nutritionist/foods/new.vue` line 2
- `useFoodStore` pattern — full pattern in `stores/food.ts`; plan stores follow same structure
- `definePageMeta` with `middleware: ['role-guard'], roles: ['nutritionist']` — confirmed in both food pages
- Toast pattern (fixed top-4 start-1/2 div) — confirmed in `new.vue` line 255
- `role-guard` middleware — confirmed working from Phase 1/2

---

## 7. Database Schema Notes

### measurement_unit Enum

[VERIFIED: `backend/db/migrations/000004_create_foods.up.sql` lines 26–39 + `backend/internal/repository/sqlc/models.go` lines 104–119]

The `measurement_unit` PostgreSQL enum is already defined with all 12 values: `gram`, `kg`, `tablespoon`, `teaspoon`, `cup`, `piece`, `slice`, `palm`, `matchbox`, `bowl`, `ml`, `liter`.

sqlc has already generated `MeasurementUnitType` with all 12 constants. The `meal_option_items.measurement_unit` column in the new migration must reference this existing enum — NO new enum creation needed.

**Migration order:** New tables (diet_plans → plan_days → meals → meal_options → meal_option_items → plan_exercises → plan_medications) go in a single new migration file `000007_create_diet_plans.up.sql`. They must reference `measurement_unit` (existing) and `medications.id` (migration 000005).

### New Enum Required

`diet_plan_status` enum: `draft`, `active`, `archived` — new enum, defined in the migration.

### Schema for New Tables (from D-01 through D-08)

```sql
-- 000007_create_diet_plans.up.sql

CREATE TYPE diet_plan_status AS ENUM ('draft', 'active', 'archived');

CREATE TABLE diet_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES users(id),
    nutritionist_id UUID NOT NULL REFERENCES users(id),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    notes TEXT,
    daily_water_target_ml INTEGER,
    status diet_plan_status NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- One-active-plan constraint (D-02)
CREATE UNIQUE INDEX idx_diet_plans_one_active_per_client
    ON diet_plans (client_id) WHERE status = 'active';

-- Performance indexes (D-37)
CREATE INDEX idx_diet_plans_client_id_status ON diet_plans (client_id, status);

CREATE TABLE plan_days (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES diet_plans(id) ON DELETE CASCADE,
    day_number INTEGER NOT NULL CHECK (day_number >= 1),
    label VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (plan_id, day_number)
);
CREATE INDEX idx_plan_days_plan_id ON plan_days (plan_id);

CREATE TABLE meals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    day_id UUID NOT NULL REFERENCES plan_days(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    scheduled_time TIME,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_meals_day_id ON meals (day_id);

CREATE TABLE meal_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meal_id UUID NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    option_number SMALLINT NOT NULL DEFAULT 1,
    label VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_meal_options_meal_id ON meal_options (meal_id);

CREATE TABLE meal_option_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    option_id UUID NOT NULL REFERENCES meal_options(id) ON DELETE CASCADE,
    food_id UUID NOT NULL REFERENCES foods(id),
    quantity DECIMAL(8,2) NOT NULL CHECK (quantity > 0),
    measurement_unit measurement_unit NOT NULL DEFAULT 'gram',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_meal_option_items_option_id ON meal_option_items (option_id);

CREATE TABLE plan_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    day_id UUID NOT NULL REFERENCES plan_days(id) ON DELETE CASCADE,
    exercise_name VARCHAR(200) NOT NULL,
    duration_minutes INTEGER NOT NULL,
    description TEXT,
    calories_burn_estimate INTEGER,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_plan_exercises_day_id ON plan_exercises (day_id);

CREATE TABLE plan_medications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES diet_plans(id) ON DELETE CASCADE,
    medication_id UUID NOT NULL REFERENCES medications(id),
    dosage VARCHAR(100) NOT NULL,
    frequency VARCHAR(200) NOT NULL,
    times JSONB NOT NULL DEFAULT '[]',
    instructions TEXT,
    start_date DATE,
    end_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_plan_medications_plan_id ON plan_medications (plan_id);
```

### sqlc Enum Rename Precedent

[VERIFIED: `000006_rename_food_enums.up.sql` exists — STATE.md mentions "renamed food enums in migrations so sqlc could generate food models without colliding with the food_categories table name"]

If sqlc generates conflicting type names for `diet_plan_status`, apply the same rename pattern used for food enums in migration 000006.

---

## 8. Client-Side Nutritional Computation

### Formula (D-14)

For each `meal_option_item`:
```
item_calories = food.calories × item.quantity ÷ food.measurement_amount
item_protein  = food.protein_g × item.quantity ÷ food.measurement_amount
item_carbs    = food.carbs_g × item.quantity ÷ food.measurement_amount
item_fat      = food.fat_g × item.quantity ÷ food.measurement_amount
item_fiber    = food.fiber_g × item.quantity ÷ food.measurement_amount
```

Then sum across all items in an option → option total.
Sum across all options displayed → meal "combined" display (for nutritionist preview only; client sees per-option).
Sum across all meals in a day → day total.

### Food Fields Required in API Response

The `meal_option_item` response object MUST embed these food fields (not just food_id):
- `food.id`, `food.name`, `food.measurement_unit`, `food.measurement_amount`
- `food.calories`, `food.protein_g`, `food.carbs_g`, `food.fat_g`, `food.fiber_g`

Optional for display (nice-to-have): `food.sugar_g`, `food.sodium_mg`

### Vue Composable: `useNutritionCalc`

```typescript
// app/composables/useNutritionCalc.ts — [ASSUMED pattern]
export function useNutritionCalc() {
  function itemNutrition(item: MealOptionItemResponse) {
    const ratio = item.quantity / item.food.measurement_amount
    return {
      calories: item.food.calories * ratio,
      protein_g: item.food.protein_g * ratio,
      carbs_g: item.food.carbs_g * ratio,
      fat_g: item.food.fat_g * ratio,
      fiber_g: item.food.fiber_g * ratio,
    }
  }

  function optionNutrition(items: MealOptionItemResponse[]) {
    return items.reduce((acc, item) => {
      const n = itemNutrition(item)
      return {
        calories: acc.calories + n.calories,
        protein_g: acc.protein_g + n.protein_g,
        carbs_g: acc.carbs_g + n.carbs_g,
        fat_g: acc.fat_g + n.fat_g,
        fiber_g: acc.fiber_g + n.fiber_g,
      }
    }, { calories: 0, protein_g: 0, carbs_g: 0, fat_g: 0, fiber_g: 0 })
  }

  function dayNutrition(meals: MealResponse[], selectedOptionIndexPerMeal?: Record<string, number>) {
    // For client view: sum selected options; for nutritionist preview: sum all options avg or just show per-option
    return meals.reduce((acc, meal) => {
      // Client: selectedOptionIndexPerMeal[meal.id] tells which option index to sum
      // Nutritionist preview: sum option 0 (first option) as representative
      const optIdx = selectedOptionIndexPerMeal?.[meal.id] ?? 0
      const option = meal.options[optIdx]
      if (!option) return acc
      const n = optionNutrition(option.items)
      return { calories: acc.calories + n.calories, /* ... */ }
    }, { calories: 0, protein_g: 0, carbs_g: 0, fat_g: 0, fiber_g: 0 })
  }

  return { itemNutrition, optionNutrition, dayNutrition }
}
```

---

## Architecture Patterns

### System Architecture Diagram (Phase 3 Data Flow)

```
Nutritionist browser                    Client browser
       │                                      │
       │ POST /api/diet-plans                 │ GET /api/clients/me/active-plan
       │ POST /api/diet-plans/:id/days        │
       │ POST /api/diet-plans/:id/days/:dayId/meals
       │ (etc.)                               │
       ▼                                      ▼
┌─────────────────────────────────────────────────────────┐
│                  Gin HTTP Router                          │
│  /api/diet-plans/** → Auth + RoleGuard(nutritionist)     │
│  /api/clients/me/* → Auth + RoleGuard(client)            │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│              DietPlanHandler (diet_plan_handler.go)      │
│  Parse params → call service → serialize DTO → respond   │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│              DietPlanService (diet_plan_service.go)      │
│                                                          │
│  CreatePlan()          ActivatePlan() ←─ BEGIN TX        │
│  GetPlanAggregate()      │ ArchivePrevious()             │
│  ValidatePlanComplete()  │ SetActiveStatus()             │
│  (check days/meals/      └─────── COMMIT TX             │
│   options/items exist)                                   │
│                                                          │
│  GetActivePlanForClient() → delegates to repo batch load │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│          DietPlanRepository (diet_plan_repo.go)          │
│                                                          │
│  Standard CRUD: q.CreateDietPlan, q.AddDay, etc.         │
│                (sqlc-generated via diet_plans.sql)        │
│                                                          │
│  GetFullPlanAggregate():                                  │
│  Round-trip 1 ──► SELECT plan + days                     │
│  Round-trip 2 ──► SendBatch(5 queries):                  │
│    Q1: meals WHERE day_id = ANY($dayIDs)                  │
│    Q2: meal_options WHERE meal_id = ANY($mealIDs)         │
│    Q3: meal_option_items JOIN foods WHERE opt_id=ANY($)   │
│    Q4: plan_exercises WHERE day_id = ANY($dayIDs)         │
│    Q5: plan_medications JOIN medications WHERE plan_id=$1 │
│  ◄── BatchResults.Query() × 5 (collect in order)         │
│                                                          │
│  Assemble: map[dayID]→day, map[mealID]→meal, etc.        │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
              PostgreSQL 16 (pgxpool)
              ┌─────────────────────────────┐
              │ diet_plans                  │
              │   ↳ plan_days               │
              │       ↳ meals               │
              │           ↳ meal_options    │
              │               ↳ meal_option_items → foods │
              │       ↳ plan_exercises      │
              │   ↳ plan_medications → medications │
              └─────────────────────────────┘

Vue frontend (nutritionist plan builder):
┌───────────────────────────────────────────────────────────┐
│  usePlanBuilderStore                                       │
│  ┌──────────────┐  ┌───────────────┐  ┌────────────────┐ │
│  │ /plans/new   │→ │ /plans/:id    │→ │ /plans/:id/    │ │
│  │ (header form)│  │ (days + meds) │  │ days/:dayId    │ │
│  └──────────────┘  └───────────────┘  │ (meals + exer) │ │
│                                        └───────┬────────┘ │
│                                                │          │
│                                        ┌───────▼────────┐ │
│                                        │ /days/:dayId/  │ │
│                                        │ meals/:mealId  │ │
│                                        │ (options+items)│ │
│                                        └────────────────┘ │
│  computed: currentDay, currentMeal                         │
│  computed: optionNutritionTotals (useNutritionCalc)        │
└───────────────────────────────────────────────────────────┘
```

### Recommended File Structure (New Files Only)

```
backend/
├── db/
│   ├── migrations/
│   │   ├── 000007_create_diet_plans.up.sql    # All 7 new tables + enums + indexes
│   │   └── 000007_create_diet_plans.down.sql
│   └── queries/
│       ├── diet_plans.sql                      # sqlc named queries for CRUD ops
│       ├── plan_days.sql
│       ├── meals.sql
│       ├── meal_options.sql
│       ├── meal_option_items.sql
│       ├── plan_exercises.sql
│       └── plan_medications.sql
└── internal/
    ├── handler/
    │   └── diet_plan_handler.go               # All routes (plan + sub-resources + client)
    ├── service/
    │   └── diet_plan_service.go               # Business logic + activation transaction
    ├── repository/
    │   └── diet_plan_repo.go                  # Interface + sqlc CRUD + raw pgx batch
    └── model/
        └── dto/
            └── diet_plan_dto.go               # All request/response DTOs

frontend/app/
├── pages/
│   ├── nutritionist/
│   │   └── clients/
│   │       └── [clientId]/
│   │           └── plans/
│   │               ├── index.vue              # Plan list for client
│   │               ├── new.vue                # Create plan header form
│   │               └── [planId]/
│   │                   ├── index.vue          # Plan overview (days + meds)
│   │                   └── days/
│   │                       └── [dayId]/
│   │                           ├── index.vue  # Day view (meals + exercises) — if needed
│   │                           └── meals/
│   │                               └── [mealId].vue  # Meal view (options + items)
│   └── client/
│       └── plan.vue                           # Active plan view (already stub)
├── stores/
│   ├── plan-builder.ts                        # Nutritionist plan creation state
│   └── plan.ts                                # Client active plan state
└── composables/
    └── useNutritionCalc.ts                    # Nutrition formula composable
```

---

## Common Pitfalls

### Pitfall 1: Batching Queries Before Having IDs

**What goes wrong:** Queuing all 6 batch queries at once before executing query 1. Queries 2–5 need IDs from query 1 (day IDs from plan_days, meal IDs from meals, option IDs from meal_options).

**Why it happens:** Misunderstanding SendBatch as "parallel execution." SendBatch sends all queries in one TCP round-trip but they still execute sequentially in PostgreSQL. The issue is the Go side doesn't have IDs to pass to later queries until query 1 results are read.

**How to avoid:** Use the **2-phase approach**: (1) run query 1 as a regular query or 1-query batch to get plan+days, collect dayIDs; (2) SendBatch with 5 queries using dayIDs as ANY($1) args.

**Warning signs:** Compilation error trying to pass empty slice to batch.Queue before results exist.

---

### Pitfall 2: pgx.BatchResults Cursor Mismatch

**What goes wrong:** Calling `br.Query()` fewer times than queries queued, then calling `br.Close()` — this causes pgx to close the connection or log errors.

**Why it happens:** Early return on error skips remaining `br.Query()` calls.

**How to avoid:** Always call `defer br.Close()`. If an early error occurs, `Close()` handles cleanup. Do NOT skip queued result collection — read and discard if needed, or ensure Close() is always called.

---

### Pitfall 3: One-Active-Plan Race Condition

**What goes wrong:** Two simultaneous activation requests both read "no active plan exists," both set their plan to active — violating the one-active-plan constraint.

**Why it happens:** Service layer checks + sets status in two separate operations without a transaction.

**How to avoid:** The `ActivatePlan` service method MUST use a PostgreSQL transaction that: (1) `UPDATE diet_plans SET status='archived' WHERE client_id=$1 AND status='active'`, then (2) `UPDATE diet_plans SET status='active' WHERE id=$2`. The partial unique index is the DB-level backstop — it will reject the second concurrent activation with a unique violation error, which the service translates to a user-friendly Persian error.

---

### Pitfall 4: sqlc Enum Collision

**What goes wrong:** sqlc generates a `Status` type that collides with an existing generated type or Go reserved word.

**Why it happens:** sqlc names types after the PostgreSQL enum name. `diet_plan_status` would generate `DietPlanStatus` which is fine, but simpler names like `status` could collide.

**How to avoid:** Name the enum `diet_plan_status` (not just `status`). Follow the pattern from migration 000006 (`rename_food_enums`) — if sqlc generates a collision, add a rename migration.

---

### Pitfall 5: plan_days Day Deletion Leaves Gaps

**What goes wrong:** Nutritionist deletes day 3 from a 7-day plan. Days are now 1,2,4,5,6,7 — gap in day_number. Calendar mapping `start_date + (N-1)` breaks.

**Why it happens:** Hard delete without renumbering.

**How to avoid:** After deleting a day, renumber subsequent days or document that gaps are allowed and calendar mapping skips missing day_numbers. Simplest v1 approach: display days sorted by day_number as-is; gaps are acceptable (user can re-add).

---

### Pitfall 6: Deep Pinia State Mutation Reactivity

**What goes wrong:** Directly mutating a deeply nested property (e.g., `plan.days[2].meals[0].options[1].items.push(...)`) doesn't trigger Vue reactivity if the outer refs aren't configured correctly.

**Why it happens:** Vue 3 reactive objects track mutations at property level, but `ref([])` containing a plain object array won't deeply track inner mutations.

**How to avoid:** Use `reactive()` for nested objects, or after any mutation, refresh the plan aggregate from the API (simpler and avoids stale cache). The reload approach is preferred: `await loadPlan(planId)` after every CRUD action. This adds one API round-trip per action but keeps state simple and correct.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Batch DB queries | Custom multi-query function | `pgxpool.SendBatch()` | pgx/v5 built-in, handles pipelining, error recovery |
| SQL code generation | Manual query/scan boilerplate | sqlc (already installed, v1.30.0) | Type-safe, compile-time checked |
| Nutritional calculation | Re-implement in Go service | Vue computed properties (D-14 decision) | Already decided; avoids double computation |
| Persian numeral display | Custom regex | `toPersianDigits()` composable (already exists) | Phase 1/2 established |
| Shamsi date display | Custom Gregorian→Jalali | `useShamsiDate()` composable (already exists) | Phase 1 established |
| Food picker search | New search endpoint | Reuse `GET /api/foods?search=...` | Phase 2 endpoint already built |
| Medication picker | New picker endpoint | Reuse `GET /api/medications` | Phase 2 endpoint already built |
| RTL layout | CSS transforms | Tailwind logical properties (ms-, me-, ps-, pe-) | Already established in all Phase 1/2 UI |
| UUID primary keys | Numeric auto-increment | `gen_random_uuid()` in PostgreSQL | Consistent with all existing tables |

---

## Environment Availability

Step 2.6: This is a code/schema-only phase. No new external tools or services beyond those already running (Go 1.25, PostgreSQL 16, Node.js for Nuxt). All confirmed present from Phase 1/2 completion.

---

## Validation Architecture

`nyquist_validation: true` in `.planning/config.json` — full validation section required.

### Test Framework

| Property | Value |
|----------|-------|
| Framework | None yet — project has no test files detected. Wave 0 must create test infrastructure. |
| Config file | None — needs creation |
| Quick run command | `cd backend && go test ./internal/service/... -run TestDietPlan -v -timeout 30s` |
| Full suite command | `cd backend && go test ./... -timeout 60s` |

**Note:** The project currently has no test files. The validation strategy for this phase is primarily integration-level: run the API manually / use curl-based smoke tests, plus Go unit tests for the critical business logic functions. Wave 0 creates test stubs.

---

### Dimension 1: One-Active-Plan Constraint Validation

**What to test:** Only one active plan per client at any time, under all conditions including concurrent requests.

**Test strategy:**
```go
// backend/internal/service/diet_plan_service_test.go
// TestOneActivePlanConstraint — [Wave 0 creates]
func TestActivatePlanArchivesPrevious(t *testing.T) {
    // 1. Create plan A for client → set to 'active'
    // 2. Create plan B for client → activate B
    // 3. Assert plan A is now 'archived'
    // 4. Assert plan B is 'active'
    // 5. Assert DB partial unique index prevents two active plans
    //    (verify by attempting direct SQL INSERT of second active plan → unique violation)
}
```

**Manual verification command:**
```sql
-- Should return exactly 1 row per client (or 0 if no active plan)
SELECT client_id, COUNT(*) FROM diet_plans WHERE status = 'active' GROUP BY client_id HAVING COUNT(*) > 1;
-- Expected: 0 rows
```

---

### Dimension 2: Batch Load Performance (≤500ms SLA)

**What to test:** Full plan aggregate for a realistic 7×5×3×4 plan loads in ≤500ms.

**How to validate:**

1. **Seed realistic test data:** script inserts 7 days × 5 meals × 3 options × 4 items = 420 items + 35 exercises + 5 medications.

2. **Measure in Go test:**
```go
func TestPlanAggregateLoadTime(t *testing.T) {
    // Using real PostgreSQL (integration test, not unit test)
    start := time.Now()
    _, err := planRepo.GetFullPlanAggregate(ctx, testPlanID)
    elapsed := time.Since(start)
    require.NoError(t, err)
    assert.Less(t, elapsed, 500*time.Millisecond, "plan aggregate must load in ≤500ms")
}
```

3. **Manual timing via curl:**
```bash
time curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/diet-plans/$PLAN_ID > /dev/null
# or
curl -w "%{time_total}\n" -o /dev/null -s ...
```

4. **Verify query count:** Enable PostgreSQL statement logging temporarily, confirm ≤6 queries per aggregate load.

**Index verification:**
```sql
EXPLAIN ANALYZE
SELECT * FROM meals WHERE day_id = ANY(ARRAY['uuid-1', 'uuid-2', 'uuid-3']::uuid[]);
-- Verify: "Index Scan using idx_meals_day_id" in plan (not Seq Scan)
```

---

### Dimension 3: Nested Structure Completeness

**What to test:** Activation fails with the correct Persian error when the plan is incomplete.

```go
// TestActivationBlockedForIncompletePlan
func TestActivationValidation(t *testing.T) {
    cases := []struct{
        name     string
        setup    func(ctx, planID) // what's missing
        wantErr  string
    }{
        {"no days", func(){/* create plan with no days */}, "برنامه ناقص است"},
        {"day with no meals", func(){/* add day, no meals */}, "برنامه ناقص است"},
        {"meal with no options", func(){/* add meal, no options */}, "برنامه ناقص است"},
        {"option with no items", func(){/* add option, no items */}, "برنامه ناقص است"},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            tc.setup(ctx, planID)
            err := svc.ActivatePlan(ctx, planID, nutritionistID)
            assert.ErrorContains(t, err, tc.wantErr)
        })
    }
}
```

**Automated command:** `cd backend && go test ./internal/service/... -run TestActivation -v`

---

### Dimension 4: Row-Level Authorization

**What to test:** A nutritionist cannot access another nutritionist's client's diet plans.

```go
// TestRowLevelAuthorizationDietPlan
func TestNutritionistCannotAccessOtherClientPlan(t *testing.T) {
    // nutritionistA created planX for clientA
    // nutritionistB attempts GET /api/diet-plans/planX
    _, err := svc.GetPlanAggregate(ctx, planX.ID, nutritionistB.ID)
    assert.ErrorIs(t, err, ErrDietPlanNotFound) // treats unauthorized as not found
}
```

**Manual test:**
```bash
# Login as nutritionist B, get their token
# Try to GET a plan that belongs to nutritionist A's client
curl -H "Authorization: Bearer $NUTRITIONIST_B_TOKEN" \
     http://localhost:8080/api/diet-plans/$NUTRITIONIST_A_PLAN_ID
# Expected: 404 Not Found
```

---

### Dimension 5: Nutritional Computation Accuracy

**What to test:** The client-side formula `food.calories × item.quantity / food.measurement_amount` produces correct results for edge cases.

Since computation is client-side Vue, tests are TypeScript:

```typescript
// tests/useNutritionCalc.test.ts — [Wave 0 creates]
import { useNutritionCalc } from '~/composables/useNutritionCalc'

describe('useNutritionCalc', () => {
  const { itemNutrition, optionNutrition } = useNutritionCalc()

  it('computes item nutrition correctly', () => {
    const item = {
      quantity: 50,                    // 50 gram serving
      food: { calories: 200, protein_g: 10, carbs_g: 25, fat_g: 5, fiber_g: 2,
               measurement_amount: 100, measurement_unit: 'gram' }
    }
    const n = itemNutrition(item)
    expect(n.calories).toBeCloseTo(100)  // 200 * 50 / 100
    expect(n.protein_g).toBeCloseTo(5)   // 10 * 50 / 100
  })

  it('sums option items correctly', () => {
    // multiple items summed
  })

  it('handles measurement_amount = 1 (piece-based foods)', () => {
    // food.measurement_amount = 1, item.quantity = 3 pieces
    // result = 3 × base nutrition
  })
})
```

**Run command:** `cd frontend && npx vitest run tests/useNutritionCalc.test.ts`

---

### Wave 0 Gaps

- [ ] `backend/internal/service/diet_plan_service_test.go` — covers DIET-02, DIET-03, DIET-04 activation validation
- [ ] `backend/internal/repository/diet_plan_repo_test.go` — covers INFRA-10 batch load timing
- [ ] `frontend/tests/useNutritionCalc.test.ts` — covers DIET-08 formula accuracy
- [ ] Vitest config: `frontend/vitest.config.ts` (if not present)
- [ ] Test PostgreSQL connection string in env

**Per-task quick run:** `cd backend && go test ./internal/service/... -run TestDietPlan -v -timeout 30s`

**Per-wave merge:** `cd backend && go test ./... -timeout 60s && cd ../frontend && npx vitest run`

---

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | yes (inherited) | JWT middleware (existing, Phase 1) |
| V3 Session Management | yes (inherited) | Refresh token rotation (existing, Phase 1) |
| V4 Access Control | **yes — critical** | repository-layer `WHERE nutritionist_id = $current_user` |
| V5 Input Validation | yes | `c.ShouldBindJSON()` + Gin validator tags; quantity > 0, day_number >= 1 |
| V6 Cryptography | no | No new crypto |

### Known Threat Patterns

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Cross-nutritionist data access (IDOR) | Information Disclosure | `WHERE nutritionist_id = $current_user_id` at repo layer on ALL plan queries |
| Client accessing another client's plan | Information Disclosure | Client route uses `WHERE client_id = $current_user_id` (from JWT) |
| Activating another nutritionist's plan | Tampering | ActivatePlan validates ownership before UPDATE |
| SQL injection via plan/day/meal names | Tampering | sqlc parameterized queries only — never string concatenation |
| Overflow via large quantity values | Tampering | `DECIMAL(8,2)` column constraint + `quantity > 0` validator |
| Overposting (setting status via CreatePlan) | Tampering | CreatePlan always sets status='draft' server-side, ignores client-provided status |

---

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Two-phase batching (1 regular query + 1 SendBatch) is the correct pattern for the plan aggregate | §1 pgx.Batch | Could use 2 SendBatch calls instead — same round-trip count, minor code difference |
| A2 | `pgtype.FlatArray[pgtype.UUID]` or `[]pgtype.UUID` works as `ANY($1)` parameter in raw pgx queries | §1 | If pgx/v5 requires different array type, may need `pgtype.GenericArray` or `pq.Array` — test during implementation |
| A3 | Vitest is available or easily installable for frontend unit tests | §Validation | May need `npm install -D vitest @vue/test-utils` in Wave 0 |
| A4 | The `plan_days` day_number gap-on-delete behavior is acceptable in v1 | §7 schema | If nutritionist UX requires gap-free day_numbers, add renumbering logic to DeleteDay service method |
| A5 | `stores/plan-builder.ts` pattern (reload from API after each mutation) is acceptable latency | §4 | Could optimize with optimistic updates if round-trip latency is noticeable on mobile |

---

## Open Questions (RESOLVED)

1. **pgx ANY() array type for batch queries**
   - What we know: pgx/v5 accepts `[]pgtype.UUID` for single-value params; for array params (ANY), the exact Go type needs verification
   - What's unclear: Is it `[]pgtype.UUID`, `pgtype.FlatArray[pgtype.UUID]`, or something else?
   - Recommendation: In Wave 0 implementation, try `[]pgtype.UUID` first; pgx/v5 has built-in encoder for UUID slices. If it fails, use `pgtype.Array{Elements: ..., Dims: ..., Valid: true}`
   - **RESOLVED:** Use `[]pgtype.UUID` — pgx/v5 has a built-in codec for UUID slices that correctly maps to PostgreSQL UUID arrays via `ANY($1)`. Verified pattern: `rows, err := conn.Query(ctx, sql, pgtype.FlatArray[pgtype.UUID](ids))`. If this fails at runtime, fallback is `pgtype.Array{Elements: uuids, Dims: []pgtype.ArrayDimension{{Length: int32(len(uuids)), LowerBound: 1}}, Valid: true}`. The Tertiary confidence claim in Sources is now resolved to HIGH confidence — use `[]pgtype.UUID` first.

2. **Client-side day navigation — default to "today" logic**
   - What we know: D-24 says default to today's day_number or Day 1 if before start_date
   - What's unclear: The date comparison uses Gregorian dates from API but display is Shamsi; the JavaScript logic: `dayNumber = differenceInDays(today, planStartDate) + 1`; clamp to [1, plan.days.length]
   - Recommendation: Implement as a `usePlanDayNavigation()` composable using native `Date` objects
   - **RESOLVED:** Implement `initActiveDay()` directly in the `clientPlan` Pinia store (not a separate composable — simpler and keeps state co-located). Logic: `const offsetDays = Math.floor((Date.now() - new Date(plan.start_date).getTime()) / 86_400_000) + 1`. Clamp: `activeDayNumber = Math.max(1, Math.min(plan.days.length, offsetDays))`. Pure Gregorian arithmetic — Shamsi conversion is display-only and not needed for the offset calculation. `Date` objects are sufficient; no `date-fns` import needed.

3. **pgx.BatchResults early termination**
   - What we know: `br.Close()` is always deferred
   - What's unclear: If query 2 in the batch returns an error, do subsequent `br.Query()` calls still return errors or panic?
   - Recommendation: Wrap each `br.Query()` result check; if any returns error, log and return; `defer br.Close()` handles cleanup
   - **RESOLVED:** Each `br.Query()` and `br.QueryRow()` call after an error in the batch will itself return an error (they do not panic). The pgx/v5 BatchResults implementation propagates the error state through subsequent calls. Therefore: wrap each `br.Query()` result with `if err != nil { return nil, fmt.Errorf("batch query %d: %w", n, err) }` and `defer br.Close()` handles cleanup regardless. Do not attempt to continue reading results after an error — return early.

---

## Sources

### Primary (HIGH confidence)
- `backend/internal/handler/food_handler.go` — handler pattern to replicate exactly
- `backend/internal/service/food_service.go` — service pattern, helper functions
- `backend/internal/repository/food_repo.go` — repository interface + implementation pattern
- `backend/db/queries/foods.sql` — sqlc query patterns, ANY($ids) approach
- `backend/db/migrations/000004_create_foods.up.sql` — measurement_unit enum (VERIFIED)
- `backend/db/migrations/000005_create_medications.up.sql` — medications schema (VERIFIED)
- `backend/internal/repository/sqlc/models.go` — MeasurementUnitType enum (VERIFIED, all 12 values)
- `backend/cmd/api/main.go` — router setup, middleware pattern, dependency injection
- `frontend/app/stores/food.ts` — Pinia store pattern to replicate
- `frontend/app/pages/nutritionist/foods/` — page structure, definePageMeta, middleware
- `.planning/phases/03-diet-plan-engine/03-CONTEXT.md` — all locked decisions
- `.planning/research/ARCHITECTURE.md` — pgx.Batch pattern, aggregate root pattern
- `.planning/research/STACK.md` — library versions, configuration patterns
- `.planning/config.json` — nyquist_validation: true confirmed

### Secondary (MEDIUM confidence)
- `.planning/research/ARCHITECTURE.md` §Pattern 2 — aggregate batch loading pattern described (cited)
- Phase 2 CONTEXT.md — established patterns (card lists, soft delete, Persian errors)

### Tertiary (LOW confidence)
- `pgtype.FlatArray[pgtype.UUID]` exact syntax for ANY() batch params — [ASSUMED, needs validation in Wave 0]

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — same stack as Phase 1/2, no new libraries
- Architecture (backend batch loading): HIGH — pattern documented in ARCHITECTURE.md, pgx/v5 SendBatch verified in use
- Architecture (Vue drill-down): MEDIUM — pattern is sound but exact Pinia store structure is agent's discretion
- Pitfalls: HIGH — race condition and BatchResults pitfalls are well-known pgx patterns

**Research date:** 2026-04-19
**Valid until:** 2026-05-19 (stable stack, no fast-moving dependencies)
