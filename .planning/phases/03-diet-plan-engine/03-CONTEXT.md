# Phase 3: Diet Plan Engine - Context

**Gathered:** 2026-04-19
**Status:** Ready for planning
**Mode:** --auto (sub-agent, no interactive session)

<domain>
## Phase Boundary

Nutritionists can create complete multi-day diet plans with the full nested structure (Plan → Days → Meals → Options → Items) plus exercise recommendations and medication prescriptions. Clients can view their active diet plan on mobile with real-time computed nutritional totals per option, meal, and day. The phase also delivers the one-active-plan-per-client constraint, automatic archival of prior plans, and performance-safe batch loading (≤5 queries, no N+1).

Phase 4 owns all client *tracking* (marking food eaten, water, sleep). Phase 3 is plan creation (nutritionist) + plan viewing (client, read-only).

</domain>

<decisions>
## Implementation Decisions

### Database Schema

- **D-01:** `diet_plans` table: `id`, `client_id` (FK users), `nutritionist_id` (FK users), `start_date` DATE, `end_date` DATE, `notes` TEXT, `daily_water_target_ml` INTEGER, `status` ENUM('draft','active','archived'), `created_at`, `updated_at`.
- **D-02:** One-active-plan-per-client enforced at TWO layers: (1) PostgreSQL partial unique index `CREATE UNIQUE INDEX ON diet_plans (client_id) WHERE status = 'active'`; (2) service layer sets previous active plan to 'archived' before activating the new plan.
- **D-03:** `plan_days` table: `id`, `plan_id` (FK diet_plans), `day_number` INTEGER (1-based, e.g. 1..7 for a 7-day plan), `label` VARCHAR(100) NULLABLE (e.g. "روز اول"), `created_at`. Unique constraint: `(plan_id, day_number)`. No repeating pattern complexity in v1 — days are simply numbered cycles. Calendar mapping: day N maps to `start_date + (N-1) days`.
- **D-04:** `meals` table: `id`, `day_id` (FK plan_days), `title` VARCHAR(200), `scheduled_time` TIME NULLABLE, `display_order` INTEGER. Client sees meals sorted by `display_order` ASC then `scheduled_time` ASC.
- **D-05:** `meal_options` table: `id`, `meal_id` (FK meals), `option_number` SMALLINT (1-based, client picks ONE option per meal), `label` VARCHAR(100) NULLABLE (e.g. "گزینه اول"), `created_at`. Client sees all options; selects one per meal during tracking (Phase 4).
- **D-06:** `meal_option_items` table: `id`, `option_id` (FK meal_options), `food_id` (FK foods), `quantity` DECIMAL(8,2), `measurement_unit` measurement_unit (same enum as foods), `notes` TEXT NULLABLE. Links to the shared food DB from Phase 2.
- **D-07:** `plan_exercises` table: `id`, `day_id` (FK plan_days), `exercise_name` VARCHAR(200), `duration_minutes` INTEGER, `description` TEXT NULLABLE, `calories_burn_estimate` INTEGER NULLABLE, `display_order` INTEGER.
- **D-08:** `plan_medications` table: `id`, `plan_id` (FK diet_plans), `medication_id` (FK medications), `dosage` VARCHAR(100), `frequency` VARCHAR(200), `times` JSONB (array of time strings, e.g. `["08:00","20:00"]`), `instructions` TEXT NULLABLE, `start_date` DATE NULLABLE, `end_date` DATE NULLABLE.
- **D-09:** All FKs use `ON DELETE CASCADE` for child tables (days cascade-delete meals, meals cascade-delete options, options cascade-delete items, days cascade-delete exercises). Plan medications cascade-delete if plan deleted (hard delete is allowed for plans in 'draft' status only; active/archived plans are never hard deleted).
- **D-10:** Row-level authorization: nutritionist can only access diet plans for own clients (enforced at repository layer: `WHERE nutritionist_id = $current_user_id`). Client can only access own active plan. Super Admin has no diet plan management role.

### Batch Loading (DIET-12)

- **D-11:** Full plan aggregate uses `pgx.Batch` (SendBatch) with exactly **4 queries** sent in a single batch round-trip:
  1. `SELECT * FROM diet_plans JOIN plan_days WHERE plan_id = $1 ORDER BY day_number`
  2. `SELECT * FROM meals WHERE day_id = ANY($day_ids) ORDER BY display_order`
  3. `SELECT * FROM meal_options WHERE meal_id = ANY($meal_ids) ORDER BY option_number`
  4. `SELECT moi.*, f.* FROM meal_option_items moi JOIN foods f ON moi.food_id = f.id WHERE moi.option_id = ANY($option_ids)`
  Exercises and plan medications fetched in two additional simple queries (total ≤6, within ≤500ms SLA). Adjust target to ≤6 queries to include exercises and medications naturally.
- **D-12:** Tree assembly happens in the Go service layer using a map-based aggregation pattern (not N+1 loops). All IDs collected in-memory and used for the next `ANY($ids)` query in batch.
- **D-13:** For nutritionist CRUD operations (adding/editing individual nodes), single-record queries are fine — batch loading is only for the aggregate read endpoint used by the client view and nutritionist preview.

### Nutritional Computation (DIET-08)

- **D-14:** Nutritional totals are computed **client-side** in the Vue frontend using food data embedded in the plan aggregate response. The backend stores quantities and food IDs; the frontend multiplies `(food.calories * item.quantity / food.measurement_amount)` to compute per-item calories, then sums to option/meal/day totals.
- **D-15:** The plan aggregate API response includes full food nutritional data for every `meal_option_item` (not just food_id). This enables offline computation and avoids any N+1 backend calls.
- **D-16:** Real-time update in the plan builder: when nutritionist adds/removes/edits an item's quantity, the totals at option, meal, and day level update immediately in the UI (Vue computed properties derived from reactive item list, no debounce needed for computation).

### Plan Builder UI (Nutritionist Side)

- **D-17:** Mobile-first **drill-down navigation** pattern for plan creation. The UI navigates into progressively deeper levels with a back button and breadcrumb:
  - `/nutritionist/clients/:clientId/plans/new` — Plan header form (dates, water target, notes)
  - `/nutritionist/clients/:clientId/plans/:planId` — Plan overview: days list + medication prescriptions card
  - `/nutritionist/clients/:clientId/plans/:planId/days/:dayId` — Day view: meals list + exercise recommendations
  - `/nutritionist/clients/:clientId/plans/:planId/days/:dayId/meals/:mealId` — Meal view: options with items
- **D-18:** Each level shows a scrollable list of children with an "+ اضافه کردن" (Add) FAB or button at the bottom. Items are reorderable via drag handles (meals have display_order; options by option_number). For v1, reorder via up/down arrow buttons (no drag-and-drop complexity).
- **D-19:** The plan stays in 'draft' status during creation. Nutritionist explicitly activates it ("فعال‌سازی برنامه") with a confirmation modal. Activation archives the previous active plan.
- **D-20:** Food item search within the option item picker: calls `GET /api/foods?search=...` (existing Phase 2 endpoint). Shows a bottom-sheet search modal. Nutritionist selects food, sets quantity+unit, taps "افزودن".
- **D-21:** Nutritional totals display on the meal/option view: small inline badges showing "کالری: ۲۵۰ | پروتئین: ۱۵ گرم" computed live.
- **D-22:** Plan status indicator shown in plan header: "پیش‌نویس" (badge, orange), "فعال" (badge, green), "آرشیو" (badge, grey). Archived plans are read-only.

### Client Plan View (DIET-03, DIET-11)

- **D-23:** Client plan view route: `/client/plan` shows active plan. If no active plan: empty state "برنامه‌ای فعال ندارید" + contact nutritionist message.
- **D-24:** Day navigation: scrollable horizontal tab bar at top showing "روز ۱", "روز ۲", etc. (Persian digits via `toPersianDigits()`). Tapping a tab shows that day's content. Default: today's day_number or Day 1 if before start_date.
- **D-25:** Each day view shows: (1) Meals in display_order with their options stacked vertically; options not yet interactive (Phase 4 adds tracking). (2) Exercise recommendations card at the bottom of the day. Water target shown in plan header bar (not per-day).
- **D-26:** Each meal option shows all its food items with quantities and computed nutritional totals. Options are numbered (گزینه ۱، گزینه ۲). Client can collapse/expand individual options via accordion.
- **D-27:** Medication prescriptions shown as a dedicated "داروها" card on the plan overview page (not per-day). Lists: medication name, dosage, frequency, scheduled times.
- **D-28:** Archived plan history: nutritionist can navigate to `/nutritionist/clients/:clientId/plans` to see all plans (active + archived). Client can see archived plans under a "تاریخچه برنامه‌ها" tab.

### API Routes

- **D-29:** Nutritionist API routes (all require `nutritionist` or `super_admin` role + own-client check):
  - `POST /api/diet-plans` — Create plan in draft status
  - `GET /api/diet-plans/:id` — Get plan aggregate (full nested structure for builder)
  - `PATCH /api/diet-plans/:id` — Update plan header fields
  - `PATCH /api/diet-plans/:id/activate` — Activate plan (archives previous active)
  - `DELETE /api/diet-plans/:id` — Hard delete (draft only)
  - `GET /api/clients/:clientId/plans` — List all plans for a client
  - Sub-resource routes for days, meals, options, items, exercises, medications (nested CRUD)
- **D-30:** Sub-resource CRUD routes:
  - `POST /api/diet-plans/:id/days`, `PUT /api/diet-plans/:id/days/:dayId`, `DELETE /api/diet-plans/:id/days/:dayId`
  - `POST /api/diet-plans/:id/days/:dayId/meals`, `PUT .../:mealId`, `DELETE .../:mealId`
  - `POST .../:mealId/options`, `PUT .../:optId`, `DELETE .../:optId`
  - `POST .../:optId/items`, `PUT .../:itemId`, `DELETE .../:itemId`
  - `POST /api/diet-plans/:id/days/:dayId/exercises`, `PUT .../:exId`, `DELETE .../:exId`
  - `POST /api/diet-plans/:id/medications`, `PUT .../:medId`, `DELETE .../:medId`
- **D-31:** Client API route: `GET /api/clients/me/active-plan` — Returns the full nested plan aggregate (batch-loaded). Used for both online view and offline cache (Phase 6).
- **D-32:** Plan list response (nutritionist): `{data: [{id, status, start_date, end_date, day_count, created_at}], total, page}`. Pagination 20/page.

### Validation & Error Handling

- **D-33:** Plan activation blocked if plan has no days, or any day has no meals, or any meal has no options, or any option has no items. Return descriptive Persian error: "برنامه ناقص است — حداقل یک روز با یک وعده و یک گزینه الزامی است".
- **D-34:** day_number must be >= 1. Meals require title (non-empty). option_number assigned auto-incrementally if not provided. Item quantity must be > 0.
- **D-35:** All error messages in Persian following project convention.

### Performance

- **D-36:** Plan aggregate endpoint SLA: ≤500ms for a 7-day × 5-meal × 3-option × 4-item plan (DIET-12 success criteria). Use pgx.Batch for parallel DB round-trips.
- **D-37:** Indexes needed: `idx_plan_days_plan_id`, `idx_meals_day_id`, `idx_meal_options_meal_id`, `idx_meal_option_items_option_id`, `idx_diet_plans_client_id_status`.

### Agent's Discretion

- Exact Pinia store structure for the plan builder state
- Drag handle vs up/down arrow ordering implementation detail
- Exact bottom-sheet animation for food item picker
- Loading skeleton for plan aggregate on client view
- Color palette for status badges (draft/active/archived)
- Whether to use a Nuxt layout-level composable or page-level state for drill-down breadcrumb

</decisions>

<specifics>
## Specific Ideas

- **Highest technical risk** (flagged in STATE.md): batch loading queries and plan builder UI state management. Researcher should investigate pgx.Batch usage patterns and Vue reactive state patterns for deeply nested forms before planning.
- **Drill-down pattern** (not a single-page form): The deeply nested Plan→Day→Meal→Option→Item hierarchy is too complex for a single form on mobile. Navigation between levels with back/breadcrumb is the right UX pattern.
- **No drag-and-drop in v1**: Reordering via up/down buttons only — keeps the MVP scope clean and avoids touch-drag complexity on mobile.
- **Client plan view is read-only in Phase 3**: Tracking (marking food eaten, logging water, etc.) is Phase 4 work. Phase 3 client view = browsing the plan structure only.
- **Food picker uses existing Phase 2 search endpoint**: `GET /api/foods?search=...` already exists and supports Persian fuzzy search — reuse directly in the plan builder's food picker bottom sheet.

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requirements
- `.planning/REQUIREMENTS.md` §Diet Plan Management (DIET-01 through DIET-12) — Full diet plan requirements, nested structure spec, nutritional computation, batch loading, exercise, medication prescriptions
- `.planning/ROADMAP.md` §Phase 3 — Phase goals, success criteria, and dependencies

### Architecture & Patterns
- `.planning/research/ARCHITECTURE.md` — Layered Go backend architecture, pgx.Batch patterns
- `.planning/research/STACK.md` — Library versions (pgx/v5, sqlc, Nuxt 4, Pinia, Dexie.js)
- `.planning/phases/01-foundation-infrastructure/01-CONTEXT.md` — Established patterns: handler→service→repository, sqlc queries, Persian utilities, mobile-first UI, Nuxt route middleware, role layouts
- `.planning/phases/02-core-data-domain/02-CONTEXT.md` — Food/medication schema and APIs reused by diet plan items; established card-based list pattern, Persian search, soft delete patterns

### Prior Phase Summaries (for integration)
- `.planning/phases/02-core-data-domain/02-02-SUMMARY.md` — Food CRUD API endpoints (reused by plan builder food picker)
- `.planning/phases/02-core-data-domain/02-04-SUMMARY.md` — Medication CRUD API endpoints (reused by plan medication prescriptions)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- **`GET /api/foods?search=...`**: Phase 2 food search endpoint (food_handler.go) — directly reusable in plan builder food picker. Supports Persian fuzzy search, pagination, category filter.
- **`GET /api/medications`**: Phase 2 medication list endpoint — directly reusable in plan medication prescription picker.
- **Handler→Service→Repository pattern**: `food_handler.go`, `food_service.go`, `food_repo.go` — exact pattern to follow for diet plan handlers.
- **sqlc query pattern**: `backend/db/queries/foods.sql` — shows how to write parameterized queries with `ANY($ids)` for batch fetching.
- **`normalize_persian()` DB function**: Available in all migrations via Phase 2. Not needed for plan data but available if needed for exercise/plan names.
- **`measurement_unit` enum**: Already defined in DB — reuse directly for `meal_option_items.measurement_unit`.
- **`AppButton.vue`, `AppInput.vue`, `LoadingSpinner.vue`**: Phase 1/2 UI components — reuse in plan builder forms.
- **`useShamsiDate()` composable**: For converting plan start_date/end_date to Shamsi for display.
- **`toPersianDigits()`**: For day numbers, option numbers in client view.
- **Pinia store pattern** (`stores/auth.ts`): Follow same pattern for diet plan stores (separate stores for builder state and client view state).
- **Role middleware** (`middleware: ['auth', 'role.nutritionist']`): All nutritionist plan pages use this middleware pattern.

### Established Patterns
- **All errors in Persian**: Every error response uses Persian text. Continue throughout Phase 3.
- **Mobile-first cards, no tables**: All list views use card-based layout, not desktop tables.
- **Soft delete**: NOT applicable to diet plan items — nutrition plans and their children can be hard deleted (draft plans) or are immutable once archived (active/archived plans).
- **RTL Tailwind logical properties**: `ms-`, `me-`, `ps-`, `pe-`, `text-start`, `text-end` everywhere.
- **`created_by` audit**: Add `created_by` to `diet_plans` table (maps to `nutritionist_id` — same field serves both purposes).

### Integration Points
- **Nutritionist bottom nav**: Phase 2 added Foods and Medications tabs. Phase 3 adds "برنامه‌ها" (Plans) access from client profile page — NOT a new bottom nav tab (plan builder is accessed via client profile, not standalone).
- **Client bottom nav**: `layouts/client.vue` has a "Plan" tab already routing to `/client/plan`. Phase 3 implements that page.
- **`/client/plan` page**: Exists as a stub from Phase 1 — needs implementation in Phase 3.
- **Client handler**: `backend/internal/handler/client_handler.go` exists — add `GetActivePlan` endpoint there or in a new `diet_plan_handler.go`.

</code_context>

<deferred>
## Deferred Ideas

- **Plan templates**: Nutritionist saves a plan as a template to reuse for future clients — out of Phase 3 scope, defer to backlog.
- **Drag-and-drop reordering**: Touch-based drag handles for meal/option reordering — defer to Phase 7 polish or backlog.
- **Copy day**: Copy one day's structure to another day in the plan — useful feature, defer to backlog.
- **Nutritional goals/thresholds**: Warning when a day's totals exceed/fall below target ranges — Phase 3 shows totals only, no threshold logic. Defer to backlog.
- **Plan PDF export**: Export the diet plan as a PDF for the client — defer to backlog.
- **Bulk day creation**: Wizard to auto-generate N days with a template — defer to backlog.
- **Client plan option selection persistence**: Client marks which option they'll eat per meal — this is Phase 4 (food intake tracking), NOT Phase 3.

</deferred>

---

*Phase: 03-diet-plan-engine*
*Context gathered: 2026-04-19*
