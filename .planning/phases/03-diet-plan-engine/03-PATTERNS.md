# Phase 3: Diet Plan Engine — Pattern Map

**Mapped:** 2026-04-19
**Files analyzed:** 22 new/modified files
**Analogs found:** 22 / 22

---

## File Classification

| New/Modified File | Role | Data Flow | Closest Analog | Match Quality |
|---|---|---|---|---|
| `backend/db/migrations/000007_create_diet_plans.up.sql` | migration | batch | `000005_create_medications.up.sql` | exact |
| `backend/db/migrations/000007_create_diet_plans.down.sql` | migration | batch | `000005_create_medications.down.sql` | exact |
| `backend/db/queries/diet_plans.sql` | query | CRUD | `backend/db/queries/foods.sql` | exact |
| `backend/db/queries/plan_days.sql` | query | CRUD | `backend/db/queries/medications.sql` | role-match |
| `backend/db/queries/meals.sql` | query | CRUD | `backend/db/queries/medications.sql` | role-match |
| `backend/db/queries/meal_options.sql` | query | CRUD | `backend/db/queries/medications.sql` | role-match |
| `backend/db/queries/meal_option_items.sql` | query | CRUD | `backend/db/queries/medications.sql` | role-match |
| `backend/db/queries/plan_exercises.sql` | query | CRUD | `backend/db/queries/medications.sql` | role-match |
| `backend/db/queries/plan_medications.sql` | query | CRUD | `backend/db/queries/medications.sql` | role-match |
| `backend/internal/model/dto/diet_plan_dto.go` | model | request-response | `backend/internal/model/dto/food_dto.go` | exact |
| `backend/internal/repository/diet_plan_repo.go` | repository | CRUD + batch | `backend/internal/repository/food_repo.go` | exact |
| `backend/internal/service/diet_plan_service.go` | service | CRUD + event-driven | `backend/internal/service/food_service.go` | exact |
| `backend/internal/handler/diet_plan_handler.go` | handler | request-response | `backend/internal/handler/food_handler.go` | exact |
| `backend/cmd/api/main.go` *(modified)* | config | request-response | existing `main.go` lines 100–186 | exact |
| `frontend/app/stores/plan-builder.ts` | store | event-driven | `frontend/app/stores/food.ts` | role-match |
| `frontend/app/stores/plan.ts` | store | request-response | `frontend/app/stores/food.ts` | role-match |
| `frontend/app/composables/useNutritionCalc.ts` | composable | transform | `frontend/app/composables/useApi.ts` | partial |
| `frontend/app/pages/nutritionist/clients/[clientId]/plans/index.vue` | page | CRUD | `frontend/app/pages/nutritionist/foods/index.vue` | exact |
| `frontend/app/pages/nutritionist/clients/[clientId]/plans/new.vue` | page | CRUD | `frontend/app/pages/nutritionist/foods/new.vue` | exact |
| `frontend/app/pages/nutritionist/clients/[clientId]/plans/[planId]/index.vue` | page | CRUD | `frontend/app/pages/nutritionist/foods/[id].vue` | role-match |
| `frontend/app/pages/nutritionist/clients/[clientId]/plans/[planId]/days/[dayId].vue` | page | CRUD | `frontend/app/pages/nutritionist/foods/[id].vue` | role-match |
| `frontend/app/pages/nutritionist/clients/[clientId]/plans/[planId]/days/[dayId]/meals/[mealId].vue` | page | CRUD | `frontend/app/pages/nutritionist/foods/[id].vue` | role-match |
| `frontend/app/pages/client/plan.vue` *(modified from stub)* | page | request-response | `frontend/app/pages/client/index.vue` | role-match |

---

## Pattern Assignments

---

### `backend/db/migrations/000007_create_diet_plans.up.sql` (migration)

**Analog:** `backend/db/migrations/000005_create_medications.up.sql`

**Enum definition pattern** (lines 1–10 of 000005):
```sql
CREATE TYPE medication_form AS ENUM (
    'tablet',
    'capsule',
    ...
);
```
➜ Follow same pattern for `diet_plan_status`:
```sql
CREATE TYPE diet_plan_status AS ENUM ('draft', 'active', 'archived');
```

**Table + index pattern** (lines 12–34 of 000005):
```sql
CREATE TABLE medications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ...
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_medications_is_active ON medications (is_active);
CREATE INDEX idx_medications_created_by ON medications (created_by);
```
➜ Apply identical structure for all 7 new tables. Key differences vs. foods migration:
- `measurement_unit` enum already exists — do **NOT** re-define it; reference it as a column type directly.
- Partial unique index for one-active-plan constraint (no analog exists yet — see RESEARCH.md §7):
```sql
CREATE UNIQUE INDEX idx_diet_plans_one_active_per_client
    ON diet_plans (client_id) WHERE status = 'active';
```
- All child FKs use `ON DELETE CASCADE` (foods/medications do not use cascade — this is new for Phase 3).

**Full schema:** See RESEARCH.md §7 for complete DDL (verified against D-01–D-10).

---

### `backend/db/queries/diet_plans.sql` (sqlc named queries)

**Analog:** `backend/db/queries/foods.sql`

**Insert pattern** (foods.sql lines 1–8):
```sql
-- name: CreateFood :one
INSERT INTO foods (name, name_normalized, ...) VALUES ($1, ...) RETURNING *;
```
➜ Diet plan insert (no name_normalized — no Persian fuzzy search on plan data):
```sql
-- name: CreateDietPlan :one
INSERT INTO diet_plans (client_id, nutritionist_id, start_date, end_date, notes,
    daily_water_target_ml, status)
VALUES ($1, $2, $3, $4, $5, $6, 'draft')
RETURNING *;
```

**Row-level filter pattern** (foods.sql lines 56–80): Use `WHERE nutritionist_id = $2` for ownership checks:
```sql
-- name: GetDietPlanByID :one
SELECT * FROM diet_plans WHERE id = $1 AND nutritionist_id = $2;
```

**List + count pair** (foods.sql lines 22–54):
```sql
-- name: ListClientPlans :many
SELECT id, status, start_date, end_date, created_at,
       (SELECT COUNT(*) FROM plan_days WHERE plan_id = diet_plans.id) AS day_count
FROM diet_plans
WHERE client_id = $1 AND nutritionist_id = $2
ORDER BY created_at DESC
LIMIT sqlc.arg('limit_val')::bigint OFFSET sqlc.arg('offset_val')::bigint;

-- name: CountClientPlans :one
SELECT COUNT(*) FROM diet_plans WHERE client_id = $1 AND nutritionist_id = $2;
```

**Exec (no return) pattern** — used for activate/archive (foods.sql uses `:one`, use `:exec` for status updates):
```sql
-- name: ActivateDietPlan :exec
UPDATE diet_plans SET status = 'active', updated_at = NOW() WHERE id = $1 AND nutritionist_id = $2;

-- name: ArchivePreviousActivePlan :exec
UPDATE diet_plans SET status = 'archived', updated_at = NOW()
WHERE client_id = $1 AND status = 'active' AND id != $2;
```

**Sub-resource queries** go in separate files (plan_days.sql, meals.sql, etc.) — see analogs below.

---

### `backend/db/queries/plan_days.sql` through `plan_medications.sql` (sub-resource sqlc queries)

**Analog:** `backend/db/queries/medications.sql`

**Standard CRUD pattern** (medications.sql — all 4 operations):
```sql
-- name: CreateMedication :one
INSERT INTO medications (...) VALUES (...) RETURNING *;

-- name: GetMedicationByID :one
SELECT ... FROM medications WHERE id = $1;

-- name: UpdateMedication :one
UPDATE medications SET ... WHERE id = $1 RETURNING *;

-- name: DeleteMedication :exec
DELETE FROM medications WHERE id = $1;
```
➜ Apply for each sub-resource (plan_days, meals, meal_options, meal_option_items, plan_exercises, plan_medications). Add parent-FK filter on every query:
- `WHERE plan_id = $2` for plan_days
- `WHERE day_id = $2` for meals, plan_exercises
- `WHERE meal_id = $2` for meal_options
- `WHERE option_id = $2` for meal_option_items
- `WHERE plan_id = $2` for plan_medications

**Note:** Aggregate batch queries (`GetFullPlanAggregate`) are **NOT** sqlc queries — they are written as raw SQL strings directly in `diet_plan_repo.go`. Do not add them to `.sql` files.

---

### `backend/internal/model/dto/diet_plan_dto.go` (model, request-response)

**Analog:** `backend/internal/model/dto/food_dto.go`

**Request struct pattern** (food_dto.go lines 1–16):
```go
type CreateFoodRequest struct {
    Name    string  `json:"name" binding:"required,max=200"`
    ...
    Calories float64 `json:"calories" binding:"gte=0,lte=9999.99"`
}
```
➜ Diet plan requests follow identical struct/binding convention:
```go
type CreateDietPlanRequest struct {
    ClientID           string  `json:"client_id" binding:"required,uuid"`
    StartDate          string  `json:"start_date" binding:"required"`   // YYYY-MM-DD
    EndDate            string  `json:"end_date" binding:"required"`
    Notes              *string `json:"notes" binding:"omitempty,max=2000"`
    DailyWaterTargetMl *int    `json:"daily_water_target_ml" binding:"omitempty,gte=0"`
}

type CreateDayRequest struct {
    DayNumber int     `json:"day_number" binding:"required,min=1"`
    Label     *string `json:"label" binding:"omitempty,max=100"`
}

type CreateMealRequest struct {
    Title         string  `json:"title" binding:"required,max=200"`
    ScheduledTime *string `json:"scheduled_time" binding:"omitempty"` // HH:MM
    DisplayOrder  int     `json:"display_order" binding:"gte=0"`
}

type CreateMealOptionItemRequest struct {
    FoodID          string  `json:"food_id" binding:"required,uuid"`
    Quantity        float64 `json:"quantity" binding:"required,gt=0"`
    MeasurementUnit string  `json:"measurement_unit" binding:"required,oneof=gram kg tablespoon teaspoon cup piece slice palm matchbox bowl ml liter"`
    Notes           *string `json:"notes" binding:"omitempty,max=500"`
}
```

**Response struct pattern** (food_dto.go lines 41–74):
```go
type FoodResponse struct {
    ID        string  `json:"id"`
    Name      string  `json:"name"`
    CreatedAt string  `json:"created_at"`
    UpdatedAt string  `json:"updated_at"`
}
```
➜ Nested aggregate response uses the same flat-field convention with embedded child slices:
```go
type DietPlanResponse struct {
    ID                 string           `json:"id"`
    ClientID           string           `json:"client_id"`
    NutritionistID     string           `json:"nutritionist_id"`
    StartDate          string           `json:"start_date"`
    EndDate            string           `json:"end_date"`
    Notes              *string          `json:"notes,omitempty"`
    DailyWaterTargetMl *int             `json:"daily_water_target_ml,omitempty"`
    Status             string           `json:"status"`
    Days               []PlanDayResponse `json:"days"`
    Medications        []PlanMedicationResponse `json:"medications"`
    CreatedAt          string           `json:"created_at"`
    UpdatedAt          string           `json:"updated_at"`
}

type MealOptionItemResponse struct {
    ID              string       `json:"id"`
    Quantity        float64      `json:"quantity"`
    MeasurementUnit string       `json:"measurement_unit"`
    Notes           *string      `json:"notes,omitempty"`
    Food            FoodEmbedded `json:"food"` // embedded food nutritional data
}

type FoodEmbedded struct {
    ID                string  `json:"id"`
    Name              string  `json:"name"`
    Calories          float64 `json:"calories"`
    ProteinG          float64 `json:"protein_g"`
    CarbsG            float64 `json:"carbs_g"`
    FatG              float64 `json:"fat_g"`
    FiberG            float64 `json:"fiber_g"`
    MeasurementUnit   string  `json:"measurement_unit"`
    MeasurementAmount float64 `json:"measurement_amount"`
}
```

**List + error response pattern** (food_dto.go lines 62–74):
```go
type FoodListResponse struct {
    Data    []FoodResponse `json:"data"`
    Total   int64          `json:"total"`
    Page    int            `json:"page"`
    Limit   int            `json:"limit"`
    HasMore bool           `json:"has_more"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
```
➜ `DietPlanListResponse` uses identical wrapper pattern with `DietPlanSummaryResponse` items (summary omits `Days` slice — just `day_count` int).

---

### `backend/internal/repository/diet_plan_repo.go` (repository, CRUD + batch)

**Analog:** `backend/internal/repository/food_repo.go`

**Interface + implementation pattern** (food_repo.go lines 12–35):
```go
// Interface
type FoodRepository interface {
    Create(ctx context.Context, params sqlc.CreateFoodParams) (*sqlc.Food, error)
    GetByID(ctx context.Context, id uuid.UUID) (*sqlc.GetFoodByIDRow, error)
    ...
}

// Implementation struct
type foodRepository struct {
    q *sqlc.Queries
}

// Constructor
func NewFoodRepository(db sqlc.DBTX) FoodRepository {
    return &foodRepository{q: sqlc.New(db)}
}
```
➜ Diet plan repo diverges at the constructor because it needs `*pgxpool.Pool` for `SendBatch`:
```go
type DietPlanRepository interface {
    // Standard CRUD (sqlc-backed)
    CreatePlan(ctx context.Context, params sqlc.CreateDietPlanParams) (*sqlc.DietPlan, error)
    GetPlanByID(ctx context.Context, planID, nutritionistID uuid.UUID) (*sqlc.DietPlan, error)
    ActivatePlan(ctx context.Context, planID, nutritionistID uuid.UUID) error
    ArchivePreviousActivePlan(ctx context.Context, clientID, exceptPlanID uuid.UUID) error
    // ... all sub-resource CRUD methods ...

    // Batch aggregate (raw pgx — NOT sqlc)
    GetFullPlanAggregate(ctx context.Context, planID uuid.UUID) (*dto.DietPlanResponse, error)
    GetActivePlanForClient(ctx context.Context, clientID uuid.UUID) (*dto.DietPlanResponse, error)
}

type dietPlanRepository struct {
    q    *sqlc.Queries  // for standard sqlc-generated operations
    pool *pgxpool.Pool  // for SendBatch aggregate queries
}

func NewDietPlanRepository(pool *pgxpool.Pool) DietPlanRepository {
    return &dietPlanRepository{q: sqlc.New(pool), pool: pool}
}
```

**Standard CRUD method pattern** (food_repo.go lines 37–44):
```go
func (r *foodRepository) Create(ctx context.Context, params sqlc.CreateFoodParams) (*sqlc.Food, error) {
    food, err := r.q.CreateFood(ctx, params)
    if err != nil {
        return nil, err
    }
    return &food, nil
}
```
➜ Copy exactly for all standard sub-resource CRUD operations.

**pgtype UUID conversion pattern** (food_repo.go lines 52–57):
```go
func (r *foodRepository) GetByID(ctx context.Context, id uuid.UUID) (*sqlc.GetFoodByIDRow, error) {
    food, err := r.q.GetFoodByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
    // ...
}
```
➜ Use `pgtype.UUID{Bytes: id, Valid: true}` for every UUID parameter — identical in diet plan repo.

**Batch aggregate pattern** (NEW — no direct analog; based on RESEARCH.md §1):
```go
func (r *dietPlanRepository) GetFullPlanAggregate(ctx context.Context, planID uuid.UUID) (*dto.DietPlanResponse, error) {
    // Round-trip 1: get plan + days (regular query)
    plan, dayRows, err := r.getPlanAndDays(ctx, planID)
    if err != nil { return nil, err }

    // Collect day IDs
    dayIDs := make([]pgtype.UUID, len(dayRows))
    for i, d := range dayRows { dayIDs[i] = d.ID }

    // Round-trip 2: SendBatch with 5 queries
    batch := &pgx.Batch{}
    batch.Queue(sqlGetMealsByDayIDs, dayIDs)
    batch.Queue(sqlGetOptionsByMealIDs, /* mealIDs collected after br.Query() call 1 */)
    // NOTE: meal IDs come from query 1 results, so either:
    //   (a) do a 3rd round-trip for options+items, OR
    //   (b) collect meal IDs in Round-trip 1 via a JOIN
    // Recommended: 2 SendBatch calls (see RESEARCH.md §1 for two-phase pattern)

    br := r.pool.SendBatch(ctx, batch)
    defer br.Close()

    // Collect in queue order — MUST match Queue() order
    mealsRows, err := br.Query()
    // ... scan meals, collect mealIDs
    mealsRows.Close()

    // ... repeat for options, items+food JOIN, exercises, medications

    // Assemble tree using maps (see RESEARCH.md §3)
    return assemblePlanAggregate(plan, dayRows, mealRows, optRows, itemRows, exRows, medRows)
}
```

**Tree assembly** (map-based, RESEARCH.md §3):
```go
// Build O(1) lookup maps then append children in row order
dayMap := make(map[uuid.UUID]*dto.PlanDayResponse)
mealMap := make(map[uuid.UUID]*dto.MealResponse)
optionMap := make(map[uuid.UUID]*dto.MealOptionResponse)

for _, row := range dayRows {
    d := mapDayRow(row)
    dayMap[d.ID] = d
    plan.Days = append(plan.Days, d)
}
for _, row := range mealRows {
    m := mapMealRow(row)
    mealMap[m.ID] = m
    if day, ok := dayMap[m.DayID]; ok {
        day.Meals = append(day.Meals, m)
    }
}
// ... same pattern for options, items, exercises, medications
```

---

### `backend/internal/service/diet_plan_service.go` (service, CRUD + event-driven)

**Analog:** `backend/internal/service/food_service.go`

**Sentinel error pattern** (food_service.go lines 22–28):
```go
var (
    ErrFoodDuplicate          = errors.New("غذا با این نام قبلاً ثبت شده است")
    ErrFoodNotFound           = errors.New("غذا یافت نشد")
    ErrFoodUnauthorizedEdit   = errors.New("شما مجوز ویرایش این غذا را ندارید")
)
```
➜ Diet plan sentinel errors (all Persian):
```go
var (
    ErrPlanNotFound            = errors.New("برنامه غذایی یافت نشد")
    ErrPlanUnauthorized        = errors.New("شما مجوز دسترسی به این برنامه را ندارید")
    ErrPlanNotDraft            = errors.New("فقط برنامه‌های پیش‌نویس قابل ویرایش هستند")
    ErrPlanIncomplete          = errors.New("برنامه ناقص است — حداقل یک روز با یک وعده و یک گزینه الزامی است")
    ErrPlanAlreadyActive       = errors.New("این برنامه قبلاً فعال شده است")
    ErrDayNotFound             = errors.New("روز یافت نشد")
    ErrMealNotFound            = errors.New("وعده یافت نشد")
    ErrOptionNotFound          = errors.New("گزینه یافت نشد")
    ErrItemNotFound            = errors.New("آیتم یافت نشد")
    ErrExerciseNotFound        = errors.New("تمرین یافت نشد")
    ErrMedicationPrescNotFound = errors.New("نسخه دارویی یافت نشد")
)
```

**Service struct + constructor pattern** (food_service.go lines 31–42):
```go
type FoodService struct {
    foodRepo repository.FoodRepository
    logger   zerolog.Logger
}

func NewFoodService(foodRepo repository.FoodRepository, logger zerolog.Logger) *FoodService {
    return &FoodService{foodRepo: foodRepo, logger: logger}
}
```

**Method pattern — fetch + ownership check + mutate + return** (food_service.go lines 166–218):
```go
func (s *FoodService) UpdateFood(ctx context.Context, foodID, userID uuid.UUID, role string, req dto.UpdateFoodRequest) (*dto.FoodResponse, error) {
    current, err := s.foodRepo.GetByID(ctx, foodID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, ErrFoodNotFound
        }
        s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to load food for update")
        return nil, fmt.Errorf("load food for update: %w", err)
    }
    if role == string(model.RoleNutritionist) && uuid.UUID(current.CreatedBy.Bytes) != userID {
        return nil, ErrFoodUnauthorizedEdit
    }
    // ... mutate + return
}
```
➜ Diet plan methods follow identical pattern with `nutritionist_id` check instead of `created_by`:
```go
func (s *DietPlanService) UpdatePlanHeader(ctx context.Context, planID, nutritionistID uuid.UUID, req dto.UpdateDietPlanRequest) (*dto.DietPlanResponse, error) {
    plan, err := s.planRepo.GetPlanByID(ctx, planID, nutritionistID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) { return nil, ErrPlanNotFound }
        s.logger.Error().Err(err).Str("plan_id", planID.String()).Msg("failed to load plan for update")
        return nil, fmt.Errorf("load plan for update: %w", err)
    }
    if plan.Status != sqlc.DietPlanStatusDraft { return nil, ErrPlanNotDraft }
    // ... mutate
}
```

**Activation transaction pattern** (NEW — no direct analog in food service; food uses no transactions):
```go
func (s *DietPlanService) ActivatePlan(ctx context.Context, planID, nutritionistID uuid.UUID) error {
    plan, err := s.planRepo.GetPlanByID(ctx, planID, nutritionistID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) { return ErrPlanNotFound }
        return fmt.Errorf("load plan for activation: %w", err)
    }
    if plan.Status == sqlc.DietPlanStatusActive { return ErrPlanAlreadyActive }

    // Validate plan completeness (D-33)
    if err := s.validatePlanComplete(ctx, planID); err != nil {
        return err // ErrPlanIncomplete with descriptive message
    }

    // Archive previous active plan (service-layer half of D-02)
    clientID := uuid.UUID(plan.ClientID.Bytes)
    if err := s.planRepo.ArchivePreviousActivePlan(ctx, clientID, planID); err != nil {
        s.logger.Error().Err(err).Msg("failed to archive previous active plan")
        return fmt.Errorf("archive previous plan: %w", err)
    }

    // Activate
    if err := s.planRepo.ActivatePlan(ctx, planID, nutritionistID); err != nil {
        s.logger.Error().Err(err).Msg("failed to activate plan")
        return fmt.Errorf("activate plan: %w", err)
    }
    return nil
}
```

**Logger pattern** (food_service.go lines 53–54, 67–68):
```go
s.logger.Error().Err(err).Str("food_id", foodID.String()).Msg("failed to create food")
s.logger.Info().Str("food_id", foodID.String()).Str("deleted_by", userID.String()).Msg("food deleted")
```
➜ Use identical structured logging with `.Str("plan_id", planID.String()).Str("user_id", nutritionistID.String())` fields.

**Helper functions reuse** (food_service.go lines 480–525 — same package):
```go
// These helpers are already in the service package — reuse directly, do NOT copy:
numericFromFloat64(value float64) (pgtype.Numeric, error)
numericToFloat64(value pgtype.Numeric) (float64, error)
optionalText(value *string) pgtype.Text
optionalBool(value *bool) pgtype.Bool
formatTimestamp(value pgtype.Timestamptz) string
```

---

### `backend/internal/handler/diet_plan_handler.go` (handler, request-response)

**Analog:** `backend/internal/handler/food_handler.go` (lines 1–161)

**Imports pattern** (food_handler.go lines 1–12):
```go
package handler

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
    "github.com/ranjbar-dev/nutritrack/backend/internal/service"
)
```

**Handler struct + constructor** (food_handler.go lines 15–22):
```go
type FoodHandler struct {
    foodService *service.FoodService
}
func NewFoodHandler(foodService *service.FoodService) *FoodHandler {
    return &FoodHandler{foodService: foodService}
}
```

**UUID param parse pattern** (food_handler.go lines 55–59):
```go
foodID, err := uuid.Parse(c.Param("id"))
if err != nil {
    c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "شناسه نامعتبر است"})
    return
}
```

**User ID extraction pattern** (food_handler.go lines 32–36):
```go
userID, err := uuid.Parse(c.GetString("user_id"))
if err != nil {
    c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
    return
}
```

**Role extraction** (food_handler.go line 113):
```go
c.GetString("role")  // passed to service for authorization decisions
```

**Error switch pattern** (food_handler.go lines 40–49):
```go
switch {
case errors.Is(err, service.ErrFoodNotFound):
    c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
case errors.Is(err, service.ErrFoodDuplicate):
    c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
default:
    c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در ثبت غذا"})
}
```
➜ Diet plan error switch (additional cases):
```go
switch {
case errors.Is(err, service.ErrPlanNotFound):
    c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
case errors.Is(err, service.ErrPlanUnauthorized):
    c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
case errors.Is(err, service.ErrPlanNotDraft):
    c.JSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
case errors.Is(err, service.ErrPlanIncomplete):
    c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})
default:
    c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطا در پردازش برنامه غذایی"})
}
```

**JSON bind pattern** (food_handler.go lines 27–30):
```go
var req dto.CreateFoodRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
    return
}
```

**Query param bind pattern** (food_handler.go lines 78–82):
```go
var query dto.FoodListQueryParams
if err := c.ShouldBindQuery(&query); err != nil {
    c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
    return
}
```

**Success responses**:
- Create: `c.JSON(http.StatusCreated, resp)` (201)
- Get/List/Update: `c.JSON(http.StatusOK, resp)` (200)
- Delete: `c.JSON(http.StatusOK, gin.H{"message": "برنامه با موفقیت حذف شد"})`

---

### `backend/cmd/api/main.go` *(modified — add diet plan route registration)*

**Analog:** `backend/cmd/api/main.go` lines 100–186

**Repository + service + handler init pattern** (main.go lines 100–117):
```go
foodRepo := repository.NewFoodRepository(pool)
// ...
foodService := service.NewFoodService(foodRepo, logger)
// ...
foodHandler := handler.NewFoodHandler(foodService)
```
➜ Add after `medHandler` init:
```go
planRepo := repository.NewDietPlanRepository(pool)       // NOTE: passes pool, not DBTX iface
planService := service.NewDietPlanService(planRepo, logger)
planHandler := handler.NewDietPlanHandler(planService)
```

**Route group pattern** (main.go lines 147–165):
```go
foods := r.Group("/api/foods")
foods.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin"))
{
    foods.GET("", foodHandler.List)
    foods.POST("", foodHandler.Create)
    foods.GET("/:id", foodHandler.Get)
    foods.PUT("/:id", foodHandler.Update)
    foods.DELETE("/:id", foodHandler.Delete)
}
```
➜ Diet plan routes (D-29, D-30, D-31):
```go
dietPlans := r.Group("/api/diet-plans")
dietPlans.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin"))
{
    dietPlans.POST("", planHandler.CreatePlan)
    dietPlans.GET("/:id", planHandler.GetPlanAggregate)
    dietPlans.PATCH("/:id", planHandler.UpdatePlanHeader)
    dietPlans.PATCH("/:id/activate", planHandler.ActivatePlan)
    dietPlans.DELETE("/:id", planHandler.DeletePlan)
    // Days
    dietPlans.POST("/:id/days", planHandler.AddDay)
    dietPlans.PUT("/:id/days/:dayId", planHandler.UpdateDay)
    dietPlans.DELETE("/:id/days/:dayId", planHandler.DeleteDay)
    // Meals
    dietPlans.POST("/:id/days/:dayId/meals", planHandler.AddMeal)
    dietPlans.PUT("/:id/days/:dayId/meals/:mealId", planHandler.UpdateMeal)
    dietPlans.DELETE("/:id/days/:dayId/meals/:mealId", planHandler.DeleteMeal)
    // Options
    dietPlans.POST("/:id/days/:dayId/meals/:mealId/options", planHandler.AddOption)
    dietPlans.DELETE("/:id/days/:dayId/meals/:mealId/options/:optId", planHandler.DeleteOption)
    // Items
    dietPlans.POST("/:id/days/:dayId/meals/:mealId/options/:optId/items", planHandler.AddItem)
    dietPlans.PUT("/:id/days/:dayId/meals/:mealId/options/:optId/items/:itemId", planHandler.UpdateItem)
    dietPlans.DELETE("/:id/days/:dayId/meals/:mealId/options/:optId/items/:itemId", planHandler.DeleteItem)
    // Exercises
    dietPlans.POST("/:id/days/:dayId/exercises", planHandler.AddExercise)
    dietPlans.PUT("/:id/days/:dayId/exercises/:exId", planHandler.UpdateExercise)
    dietPlans.DELETE("/:id/days/:dayId/exercises/:exId", planHandler.DeleteExercise)
    // Medications
    dietPlans.POST("/:id/medications", planHandler.AddMedication)
    dietPlans.PUT("/:id/medications/:medId", planHandler.UpdateMedication)
    dietPlans.DELETE("/:id/medications/:medId", planHandler.DeleteMedication)
}

// Client plans list (nutritionist viewing client's plans)
nutri.GET("/clients/:clientId/plans", planHandler.ListClientPlans)

// Client active plan (client self-service)
client.GET("/me/active-plan", planHandler.GetActivePlan)
```

---

### `frontend/app/stores/plan-builder.ts` (store, event-driven)

**Analog:** `frontend/app/stores/food.ts`

**Interface + defineStore(composition API) pattern** (food.ts lines 1–52):
```typescript
export interface FoodResponse { id: string; name: string; ... }
export interface CreateFoodPayload { name: string; ... }

export const useFoodStore = defineStore('food', () => {
  const foods = ref<FoodResponse[]>([])
  const loading = ref(false)
  // ...
  async function fetchFoods(reset = false) {
    loading.value = true
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<FoodListResponse>(`/foods?${buildQueryParams()}`)
      // ...
    } catch (error) {
      console.error('Failed to fetch foods', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  return { foods, loading, fetchFoods, ... }
})
```
➜ Plan builder store diverges by adding navigation cursor state and drill-down helpers:
```typescript
export const usePlanBuilderStore = defineStore('plan-builder', () => {
  // Single source of truth — full aggregate from GET /api/diet-plans/:id
  const plan = ref<DietPlanAggregate | null>(null)
  const loading = ref(false)
  const saving = ref(false)

  // Navigation cursor (which level is shown)
  const currentDayId = ref<string | null>(null)
  const currentMealId = ref<string | null>(null)

  // Derived views — no state duplication
  const currentDay = computed(() =>
    plan.value?.days.find(d => d.id === currentDayId.value) ?? null
  )
  const currentMeal = computed(() =>
    currentDay.value?.meals.find(m => m.id === currentMealId.value) ?? null
  )

  async function loadPlan(planId: string) {
    loading.value = true
    try {
      const { apiFetch } = useApi()
      plan.value = await apiFetch<DietPlanAggregate>(`/diet-plans/${planId}`)
    } catch (error) {
      console.error('Failed to load plan', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function addDay(planId: string, payload: CreateDayPayload) {
    saving.value = true
    try {
      const { apiFetch } = useApi()
      await apiFetch(`/diet-plans/${planId}/days`, { method: 'POST', body: JSON.stringify(payload) })
      await loadPlan(planId) // refresh aggregate
    } finally {
      saving.value = false
    }
  }

  // Navigation
  function navigateToDay(dayId: string) { currentDayId.value = dayId; currentMealId.value = null }
  function navigateToMeal(mealId: string) { currentMealId.value = mealId }
  function resetNavigation() { currentDayId.value = null; currentMealId.value = null }

  return { plan, loading, saving, currentDayId, currentMealId, currentDay, currentMeal,
           loadPlan, addDay, /* addMeal, addOption, addItem, ... */ navigateToDay, navigateToMeal, resetNavigation }
})
```

**apiFetch usage pattern** (food.ts lines 88–92):
```typescript
const { apiFetch } = useApi()
const data = await apiFetch<FoodListResponse>(`/foods?${buildQueryParams()}`)
```
➜ Use identical pattern in plan stores. Note: `useApi()` is a composable, call it inside async functions (not at store setup level).

---

### `frontend/app/stores/plan.ts` (store, request-response — client view)

**Analog:** `frontend/app/stores/food.ts`

Simpler than plan-builder — read-only active plan for client:
```typescript
export const usePlanStore = defineStore('plan', () => {
  const activePlan = ref<DietPlanAggregate | null>(null)
  const loading = ref(false)
  const currentDayIndex = ref(0) // which day tab is selected

  const currentDay = computed(() =>
    activePlan.value?.days[currentDayIndex.value] ?? null
  )

  async function fetchActivePlan() {
    loading.value = true
    try {
      const { apiFetch } = useApi()
      activePlan.value = await apiFetch<DietPlanAggregate>('/clients/me/active-plan')
    } catch (error: any) {
      if (error?.statusCode === 404) {
        activePlan.value = null // empty state — no active plan
      } else {
        throw error
      }
    } finally {
      loading.value = false
    }
  }

  function selectDay(index: number) { currentDayIndex.value = index }

  return { activePlan, loading, currentDayIndex, currentDay, fetchActivePlan, selectDay }
})
```

---

### `frontend/app/composables/useNutritionCalc.ts` (composable, transform)

**Analog:** `frontend/app/composables/useApi.ts` (structure only — this is a new pattern)

**Composable export pattern** (useApi.ts lines 23–100):
```typescript
export function useApi() {
  // ...
  return { apiFetch }
}
```
➜ Follow identical named export function pattern:
```typescript
export function useNutritionCalc() {
  function itemNutrition(item: MealOptionItemResponse) {
    const ratio = item.quantity / item.food.measurement_amount
    return {
      calories:   +(item.food.calories  * ratio).toFixed(1),
      protein_g:  +(item.food.protein_g * ratio).toFixed(1),
      carbs_g:    +(item.food.carbs_g   * ratio).toFixed(1),
      fat_g:      +(item.food.fat_g     * ratio).toFixed(1),
      fiber_g:    +(item.food.fiber_g   * ratio).toFixed(1),
    }
  }

  function optionNutrition(items: MealOptionItemResponse[]) {
    return items.reduce(
      (acc, item) => {
        const n = itemNutrition(item)
        return {
          calories:  acc.calories  + n.calories,
          protein_g: acc.protein_g + n.protein_g,
          carbs_g:   acc.carbs_g   + n.carbs_g,
          fat_g:     acc.fat_g     + n.fat_g,
          fiber_g:   acc.fiber_g   + n.fiber_g,
        }
      },
      { calories: 0, protein_g: 0, carbs_g: 0, fat_g: 0, fiber_g: 0 }
    )
  }

  return { itemNutrition, optionNutrition }
}
```

---

### `frontend/app/pages/nutritionist/clients/[clientId]/plans/index.vue` (page, CRUD — plan list)

**Analog:** `frontend/app/pages/nutritionist/foods/index.vue`

**definePageMeta pattern** (foods/index.vue lines 3–7):
```typescript
definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist'],
})
```
➜ **Identical** — copy verbatim for all 4 nutritionist plan pages.

**Store + onMounted fetch pattern** (foods/index.vue lines 10–55):
```typescript
const foodStore = useFoodStore()
onMounted(() => {
  foodStore.fetchFoods(true)
})
```
➜ Replace with `usePlanBuilderStore()` and `planStore.fetchClientPlans(clientId)`.

**List rendering + empty state pattern** (foods/index.vue lines 140–181):
```html
<div v-if="isInitialLoading" class="py-16"><UiLoadingSpinner size="lg" /></div>
<div v-else-if="emptyState === 'empty'" class="text-center py-12 space-y-3">
  <p class="text-gray-600">هیچ برنامه‌ای ثبت نشده</p>
  <UiAppButton ...>افزودن برنامه</UiAppButton>
</div>
<div v-else class="space-y-3">
  <!-- PlanPlanCard v-for items -->
</div>
```

**Header + action button pattern** (foods/index.vue lines 91–101):
```html
<header class="flex items-center justify-between gap-3">
  <div>
    <h1 class="text-xl font-bold text-gray-800">برنامه‌های غذایی</h1>
    <p class="text-xs text-gray-500 mt-1">{{ toPersianDigits(planStore.total) }} برنامه</p>
  </div>
  <UiAppButton class="w-auto" size="sm" @click="navigateTo(newPlanUrl)">
    برنامه جدید
  </UiAppButton>
</header>
```

**Toast pattern** (foods/index.vue lines 86–90):
```html
<div v-if="showToast"
  class="fixed top-4 start-1/2 -translate-x-1/2 bg-red-600 text-white px-4 py-2 rounded-lg shadow-lg z-50 text-sm">
  {{ toastMessage }}
</div>
```
➜ Copy verbatim — used in all form/list pages.

---

### `frontend/app/pages/nutritionist/clients/[clientId]/plans/new.vue` (page, CRUD — create plan header)

**Analog:** `frontend/app/pages/nutritionist/foods/new.vue`

**Reactive form + errors pattern** (foods/new.vue lines 46–74):
```typescript
const form = reactive({
  name: '', description: '', ...
})
const errors = reactive({
  name: '', description: '', ...
})
```
➜ Diet plan header form fields:
```typescript
const form = reactive({
  client_id: route.params.clientId as string, // pre-filled from route
  start_date: '',   // YYYY-MM-DD
  end_date: '',
  daily_water_target_ml: '',
  notes: '',
})
const errors = reactive({
  start_date: '', end_date: '', daily_water_target_ml: '', notes: ''
})
```

**validateForm + handleSubmit pattern** (foods/new.vue lines 131–220):
```typescript
function validateForm() {
  resetErrors()
  let isValid = true
  if (!form.name.trim()) {
    errors.name = 'نام غذا الزامی است'; isValid = false
  }
  return isValid
}

async function handleSubmit() {
  if (!validateForm()) { triggerToast('لطفاً خطاهای فرم را اصلاح کنید'); return }
  isSubmitting.value = true
  try {
    await foodStore.createFood(buildPayload())
    navigateTo('/nutritionist/foods')
  } catch (error) {
    triggerToast((error as Error).message || 'ثبت غذا با خطا مواجه شد')
  } finally {
    isSubmitting.value = false
  }
}
```
➜ After plan creation: navigate to `planStore.plan.id` plan overview page, not back to list.

**Submit button pattern** (foods/new.vue lines 420–423):
```html
<UiAppButton type="submit" :loading="isSubmitting" :disabled="isSubmitting">
  {{ isSubmitting ? 'در حال ذخیره...' : 'ذخیره غذا' }}
</UiAppButton>
```

**Section card pattern** (foods/new.vue lines 269–298):
```html
<section class="bg-white rounded-2xl p-4 shadow-sm space-y-4">
  <h2 class="text-sm font-semibold text-gray-700">اطلاعات پایه</h2>
  <UiAppInput v-model="form.name" label="نام غذا" :disabled="isSubmitting" :error="errors.name" />
</section>
```

---

### `frontend/app/pages/nutritionist/clients/[clientId]/plans/[planId]/index.vue` (page, CRUD — plan overview)

**Analog:** `frontend/app/pages/nutritionist/foods/[id].vue`

**Route param extraction** (foods/[id].vue lines 13, 84):
```typescript
const route = useRoute()
const foodId = computed(() => route.params.id as string)
```
➜ Multiple route params:
```typescript
const clientId = computed(() => route.params.clientId as string)
const planId = computed(() => route.params.planId as string)
```

**Load on mount + loading spinner** (foods/[id].vue lines 215–228, 275–278):
```typescript
onMounted(() => { loadFood() })
// template:
<div v-if="isLoading" class="py-16"><UiLoadingSpinner size="lg" /></div>
<form v-else ...>
```

**Back button pattern** (foods/[id].vue lines 264–271):
```html
<UiAppButton class="w-auto" size="sm" variant="secondary"
  @click="navigateTo(`/nutritionist/clients/${clientId}/plans`)">
  بازگشت
</UiAppButton>
```

**Plan overview specific additions** (no analog — new pattern):
- Days list: cards with day_number + label + meal count → tap to navigate to day
- Status badge: `<span :class="statusBadgeClass(plan.status)">{{ statusLabel(plan.status) }}</span>`
  - draft → `bg-orange-100 text-orange-700`
  - active → `bg-green-100 text-green-700`
  - archived → `bg-gray-100 text-gray-600`
- Activate button (shown only if draft): `<UiAppButton @click="handleActivate">فعال‌سازی برنامه</UiAppButton>`
- Medications card at bottom

---

### `frontend/app/pages/nutritionist/clients/[clientId]/plans/[planId]/days/[dayId].vue` (page, CRUD — day view)

**Analog:** `frontend/app/pages/nutritionist/foods/[id].vue`

Same `definePageMeta` + route param + load-on-mount + back-button patterns.

Day view specific (use `usePlanBuilderStore` — do NOT fetch separately):
```typescript
const planStore = usePlanBuilderStore()
const day = computed(() => planStore.plan?.days.find(d => d.id === dayId.value) ?? null)
// planStore.plan already loaded by parent plan page — no additional fetch needed
// If arriving directly (deep link), call planStore.loadPlan(planId) in onMounted guard
```

Meals list → tap → navigate to `[planId]/days/[dayId]/meals/[mealId]`.
Exercise list → inline add/edit form below meals.

---

### `frontend/app/pages/nutritionist/clients/[clientId]/plans/[planId]/days/[dayId]/meals/[mealId].vue` (page, CRUD — meal/options/items view)

**Analog:** `frontend/app/pages/nutritionist/foods/[id].vue`

Same patterns. Meal view specific:
```typescript
const meal = computed(() => planStore.currentMeal) // via usePlanBuilderStore cursor

// Nutritional totals (D-14, D-16)
const { optionNutrition } = useNutritionCalc()
const optionTotals = computed(() =>
  meal.value?.options.map(opt => ({
    optionId: opt.id,
    totals: optionNutrition(opt.items)
  })) ?? []
)
```

Food item picker bottom sheet (D-20):
- Opens when tapping "+ افزودن آیتم"
- `<input>` calls `apiFetch('/foods?search=...')` on input (debounced 300ms — copy debounce from foods/index.vue lines 36–47)
- Lists results as selectable rows
- On select: shows quantity + unit form → submit calls `planStore.addItem(optionId, payload)`

---

### `frontend/app/pages/client/plan.vue` *(modified from stub)*

**Analog:** `frontend/app/pages/client/index.vue` (stub) + `frontend/app/pages/nutritionist/foods/index.vue` (list pattern)

**definePageMeta** (client/plan.vue already has):
```typescript
definePageMeta({
  layout: 'client',
})
```

**Full implementation pattern** (no existing client data-page analog — closest is foods/index.vue structure):
```typescript
const planStore = usePlanStore()
const { toPersianDigits } = usePersianDigits()

onMounted(async () => {
  await planStore.fetchActivePlan()
  // Default to today's day_number
  if (planStore.activePlan) {
    const todayOffset = Math.floor(
      (Date.now() - new Date(planStore.activePlan.start_date).getTime()) / 86400000
    )
    const todayIndex = Math.min(
      Math.max(todayOffset, 0),
      planStore.activePlan.days.length - 1
    )
    planStore.selectDay(todayIndex)
  }
})
```

**Empty state pattern** (foods/index.vue lines 144–149):
```html
<div v-if="!planStore.activePlan && !planStore.loading" class="text-center py-16 space-y-3">
  <p class="text-gray-600">برنامه‌ای فعال ندارید</p>
  <p class="text-sm text-gray-400">با متخصص تغذیه خود تماس بگیرید</p>
</div>
```

**Day tab bar pattern** (new UI element — no analog; use horizontal scroll with `overflow-x-auto`):
```html
<div class="flex gap-2 overflow-x-auto pb-2 scrollbar-hide">
  <button v-for="(day, idx) in planStore.activePlan.days" :key="day.id"
    class="flex-shrink-0 px-4 py-2 rounded-full text-sm"
    :class="planStore.currentDayIndex === idx
      ? 'bg-emerald-600 text-white'
      : 'bg-gray-100 text-gray-600'"
    @click="planStore.selectDay(idx)">
    {{ toPersianDigits(`روز ${day.day_number}`) }}
  </button>
</div>
```

**Accordion option pattern** (new — use `<details>` / `v-show` toggle):
```html
<div v-for="option in day.meals[mealIdx].options" :key="option.id" class="bg-gray-50 rounded-xl p-3">
  <button class="w-full flex items-center justify-between" @click="toggleOption(option.id)">
    <span class="font-medium text-sm">{{ toPersianDigits(`گزینه ${option.option_number}`) }}</span>
    <span class="text-xs text-emerald-600">{{ formatNutritionBadge(option) }}</span>
  </button>
  <div v-show="expandedOptions.has(option.id)" class="mt-2 space-y-1">
    <div v-for="item in option.items" :key="item.id" class="text-sm text-gray-600">
      {{ item.food.name }} — {{ toPersianDigits(item.quantity) }} {{ item.measurement_unit }}
    </div>
  </div>
</div>
```

---

## Shared Patterns

### Authentication / Role Guard
**Source:** `frontend/app/middleware/role-guard.ts` (lines 1–26) + `backend/internal/middleware/` (auth + role guard)

**Apply to:** All 5 nutritionist plan pages, all backend diet plan routes

Frontend (copy to every nutritionist plan page):
```typescript
definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
  roles: ['nutritionist'],
})
```

Backend (set on route groups in main.go):
```go
dietPlans.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("nutritionist", "super_admin"))
client.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("client"))
```

### API Fetch
**Source:** `frontend/app/composables/useApi.ts` (lines 23–100)

**Apply to:** All store async methods (plan-builder.ts, plan.ts)

Always call `useApi()` inside the async function body — never at store setup scope:
```typescript
async function loadPlan(planId: string) {
  const { apiFetch } = useApi()  // call inside function
  const data = await apiFetch<DietPlanAggregate>(`/diet-plans/${planId}`)
}
```

### Persian Error Messages
**Source:** `backend/internal/service/food_service.go` (lines 22–28)

**Apply to:** All Go sentinel error variables in `diet_plan_service.go`

All error text MUST be Persian. No English in error messages. Pattern: `errors.New("Persian text here")`.

### Persian Digits Display
**Source:** `frontend/app/pages/nutritionist/foods/index.vue` (line 2)

**Apply to:** All count displays, day numbers, option numbers, quantity displays in plan pages

```typescript
import { toPersianDigits } from '~/utils/persian-digits'
// In template: {{ toPersianDigits(day.day_number) }}
```

### RTL Tailwind Logical Properties
**Source:** `frontend/app/pages/nutritionist/foods/index.vue` (line 87, toast)

**Apply to:** All new Vue templates in Phase 3

Use `start-` / `end-` / `ms-` / `me-` / `ps-` / `pe-` instead of `left-` / `right-` / `ml-` / `mr-` etc.

### pgtype UUID Conversion
**Source:** `backend/internal/repository/food_repo.go` (lines 53, 62, 84)

**Apply to:** All repository methods in `diet_plan_repo.go`

```go
pgtype.UUID{Bytes: someUUID, Valid: true}  // uuid.UUID → pgtype.UUID
uuid.UUID(pgRow.ID.Bytes)                  // pgtype.UUID → uuid.UUID
```

### zerolog Structured Logging
**Source:** `backend/internal/service/food_service.go` (lines 53, 67, 248)

**Apply to:** All service methods in `diet_plan_service.go`

```go
s.logger.Error().Err(err).Str("plan_id", planID.String()).Str("user_id", userID.String()).Msg("descriptive message")
s.logger.Info().Str("plan_id", planID.String()).Msg("plan activated")
```

### Error Wrapping
**Source:** `backend/internal/service/food_service.go` (lines 54–55, 68)

**Apply to:** All service methods — wrap internal errors with `fmt.Errorf("context: %w", err)`

```go
return nil, fmt.Errorf("create diet plan: %w", err)       // not: return nil, err
return nil, fmt.Errorf("load plan for activation: %w", err)
```

### Loading State + Spinner
**Source:** `frontend/app/pages/nutritionist/foods/[id].vue` (lines 12–14, 275–278)

**Apply to:** All plan pages that fetch data on mount

```typescript
const isLoading = ref(true)
// template:
<div v-if="isLoading" class="py-16"><UiLoadingSpinner size="lg" /></div>
<div v-else><!-- content --></div>
```

### Mobile-First Card Layout
**Source:** `frontend/app/pages/nutritionist/foods/index.vue` (lines 163–169)

**Apply to:** All list views (plan list, day list, meal list, option list)

```html
<div class="space-y-3">
  <SomeEntityCard v-for="item in items" :key="item.id" :item="item" @action="handler" />
</div>
```
No tables. All lists use vertical card stacks with `space-y-3`.

---

## No Analog Found

All files have codebase analogs. The following patterns are **new** with no prior analog — planner must use RESEARCH.md patterns:

| Pattern | Used In | Reason |
|---|---|---|
| `pgx.SendBatch` aggregate query | `diet_plan_repo.go` | No batch queries exist in codebase yet |
| Two-phase batch (plan+days → then batch 5) | `diet_plan_repo.go` | New to Phase 3 |
| Map-based tree assembly | `diet_plan_repo.go` | New to Phase 3 |
| Service-layer activation transaction (archive+activate) | `diet_plan_service.go` | No transactions used in food/med services |
| PostgreSQL partial unique index | migration 000007 | No partial indexes in prior migrations |
| Drill-down navigation cursor in Pinia store | `plan-builder.ts` | `food.ts` has no navigation cursor |
| Horizontal day tab bar | `client/plan.vue` | New UI component pattern |
| Accordion option expand/collapse | `client/plan.vue` | New UI component pattern |
| Food picker bottom sheet | `meals/[mealId].vue` | New UI pattern (reuses existing `/api/foods?search=` endpoint) |
| Up/down reorder buttons | meal/option lists | New interaction (no drag-and-drop) |
| `useNutritionCalc` composable | meal/option components | New computation composable |

---

## Metadata

**Analog search scope:** `backend/internal/`, `backend/db/`, `frontend/app/pages/`, `frontend/app/stores/`, `frontend/app/composables/`, `frontend/app/middleware/`
**Files scanned:** 24
**Pattern extraction date:** 2026-04-19
