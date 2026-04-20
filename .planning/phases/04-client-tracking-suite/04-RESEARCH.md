# Phase 4: Client Tracking Suite — Research

**Researched:** 2026-04-20
**Domain:** Go tracking repositories · Nuxt 4 Pinia tracking stores · PostgreSQL upsert deduplication · Chart.js RTL charts · local filesystem file uploads
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

| # | Decision |
|---|----------|
| D-01 | Seven Phase 4 tables: `food_logs`, `water_logs`, `sleep_logs`, `exercise_logs`, `medication_logs`, `body_measurements`, `lab_results` |
| D-02 | Every writeable table: `id UUID`, `client_id UUID`, `local_id UUID UNIQUE NOT NULL`, `created_at TIMESTAMPTZ`, `updated_at TIMESTAMPTZ` where applicable |
| D-03 | `food_logs`: `date DATE`, `meal_id UUID`, `selected_option_id UUID NULL`, `is_skipped BOOLEAN DEFAULT false`, `notes TEXT NULL`; unique on `(client_id, date, meal_id)` |
| D-04 | `water_logs`: `date DATE`, `amount_ml INTEGER`, `logged_time TIME NULL`; multiple per day allowed |
| D-05 | `sleep_logs`: `date DATE`, `sleep_time TIME`, `wake_time TIME`, `quality sleep_quality`, `notes TEXT NULL`; unique on `(client_id, date)` |
| D-06 | `exercise_logs`: `date DATE`, `exercise_name VARCHAR(200)`, `duration_minutes INTEGER`, `calories_burned INTEGER NULL`, `notes TEXT NULL`; multiple per day |
| D-07 | `medication_logs`: `date DATE`, `prescribed_medication_id UUID NULL`, `medication_name VARCHAR(200)`, `dosage VARCHAR(100) NULL`, `taken_at TIME`, `notes TEXT NULL`, `is_self_reported BOOLEAN DEFAULT false` |
| D-08 | `body_measurements`: `date DATE`, weight + six circumference fields `NUMERIC(5,2) NULL`, `recorded_by UUID`; unique on `(client_id, date)` |
| D-09 | `lab_results`: `title VARCHAR(200)`, `lab_type lab_result_type`, `test_date DATE`, `file_path TEXT NULL`, `external_link TEXT NULL`, `original_filename VARCHAR(255) NULL`, `mime_type VARCHAR(100) NULL`, `file_size_bytes BIGINT NULL`, `uploaded_by UUID`; at least one of `file_path`/`external_link` required |
| D-10 | All client POST endpoints accept `local_id`; repositories implement `CreateOrUpsertByLocalID` or `UpsertByUniqueKey` returning the existing record on duplicate |
| D-11 | Nutritionist read endpoints: row-level ownership enforced in SQL by JOINing client to `users.nutritionist_id = $nutritionist_id`; never trust URL `client_id` alone |
| D-12 | List endpoints: `from`/`to` filtering, newest-first for history; chart endpoints oldest-first |
| D-13 | sqlc for standard CRUD/list; raw pgx for aggregate dashboard and chart series reads |
| D-14 | `/client/tracking` = client daily dashboard |
| D-15 | Food logging is meal-centric from active plan |
| D-16 | Water: quick-add buttons (200ml/250ml/500ml) + custom amount; progress toward `daily_water_target_ml` |
| D-17 | Sleep form: bedtime/wake pickers; duration auto-computed in UI; quality Persian-labelled |
| D-18 | Medication checklist pre-populated from active plan `plan_medications`; tap creates `medication_log` with `prescribed_medication_id` |
| D-19 | Body measurements: weight prominent, extra fields collapsible |
| D-20 | Lab results: file upload (PDF/JPG/PNG, max 10 MB) or external link; both in Phase 4 |
| D-21 | Nutritionist tracking under existing client area, not a separate global module |
| D-22 | Nutritionist can create/update body measurements for client; `recorded_by` = nutritionist UUID |
| D-23 | Weight history: Chart.js line chart + Shamsi date labels; other measurements tabular in Phase 4 |
| D-24 | Lab file downloads: `Content-Disposition: attachment`; no inline rendering |
| D-25 | Lab files: local filesystem, dedicated `lab-results/` directory, UUID filenames |
| D-26 | File validation: accept PDF/JPG/PNG, max 10 MB, MIME sniffing / magic bytes before persist |
| D-27 | All API errors Persian; Gin binding + service sentinel error pattern |
| D-28 | Tailwind logical RTL utilities, Persian numerals, Jalali dates, mobile card layouts mandatory |

### Agent's Discretion
- Exact quick-action card visual design on the daily dashboard
- Whether lab results appear on the main dashboard or a secondary section below core daily logs
- Whether measurement history beyond weight uses mini-cards or a compact list layout
- Exact folder path configuration mechanism for lab-result file storage (local filesystem only)

### Deferred Ideas (OUT OF SCOPE)
- Offline queueing, sync retries, IndexedDB caching → Phase 6
- Push reminders for meals, medications, and water → Phase 6
- Rich multi-line body measurement charts (scope risk)
- AI-derived adherence scoring or weekly summaries → v2 backlog
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| TRACK-01 | Food intake logging per meal (select option or skip) | food_logs schema + D-15 meal-centric UX |
| TRACK-02 | Water logging with amounts + timestamps | water_logs schema + D-16 quick-add UX |
| TRACK-03 | Water daily total vs target with visual progress | Aggregate query from water_logs; target from active plan `daily_water_target_ml` |
| TRACK-04 | Sleep logging: sleep_time/wake_time, quality, one per date | sleep_logs upsert on `(client_id, date)` |
| TRACK-05 | Sleep duration auto-computed from sleep/wake time | Frontend computation only; no backend column needed |
| TRACK-06 | Exercise logging: name, duration, optional calorie burn | exercise_logs schema; multiple per day |
| TRACK-07 | Medication intake logging (prescribed + self-reported) | medication_logs with `prescribed_medication_id` linkage |
| TRACK-08 | Client sees prescribed medication checklist from plan | Query `plan_medications` of active plan; each `times[]` entry = one checklist item |
| TRACK-09 | Body measurements: weight + 6 circumference fields | body_measurements upsert on `(client_id, date)` |
| TRACK-10 | Measurement history as list and weight chart | Chart.js line + Shamsi dates; tabular for others |
| TRACK-11 | One body measurement record per date (last write wins) | `ON CONFLICT (client_id, date) DO UPDATE` |
| TRACK-12 | All tracking entries have `local_id` UUID for offline dedup | `local_id UUID UNIQUE NOT NULL` + `ON CONFLICT (local_id) DO NOTHING/UPDATE` |
| TRACK-13 | Nutritionist views all client tracking with date range filter | `from`/`to` params; JOIN enforces nutritionist ownership |
| LAB-01 | Client uploads lab results (title, type, date, file or link) | `lab_results` table; multipart form handler |
| LAB-02 | At least one of file or link required | Service-layer validation sentinel error |
| LAB-03 | Accepted: PDF/JPG/PNG; max 10 MB | `gabriel-vasile/mimetype` (already in go.mod) + size guard |
| LAB-04 | Nutritionist views and downloads client lab results | `GET /api/nutritionist/clients/:id/lab-results` + download endpoint with `Content-Disposition: attachment` |
| LAB-05 | Files stored on Hetzner server filesystem | Local directory under `UPLOADS_DIR`; UUID filename in DB |
</phase_requirements>

---

## Summary

Phase 4 adds seven PostgreSQL tables, roughly 18 API endpoints, and six new client UI sections on top of the Diet Plan Engine delivered in Phase 3. The codebase patterns are stable and well-established: migrations in `backend/db/migrations/`, sqlc queries in `backend/db/queries/`, repository interfaces in `backend/internal/repository/`, services in `backend/internal/service/`, handlers wired into `backend/cmd/api/main.go` route groups. On the frontend, Pinia stores in `frontend/app/stores/`, composables in `frontend/app/composables/`, and pages under `frontend/app/pages/client/` and `frontend/app/pages/nutritionist/clients/[clientId]/`.

The most technically nuanced pieces are: (1) the `local_id` idempotent upsert pattern — which must be implemented consistently across all seven tables so Phase 6 offline sync requires zero schema changes; (2) row-level authorization for nutritionist read endpoints — enforced entirely at the SQL/repository level by JOINing through the `users.nutritionist_id` column; and (3) the daily dashboard data flow — a single aggregate read that combines the active plan's meals/medications/water target with today's logs to produce a status summary. Chart.js (`chart.js` + `vue-chartjs`) is **not yet installed** in `frontend/package.json` and must be added. The `gabriel-vasile/mimetype` package is already present in `go.mod`, making MIME-sniffing file validation a single import away.

The recommended plan decomposition is: (A) one migration plan creating all seven tables and enums; (B) one backend plan per tracking domain (food+water, sleep+exercise, medication, body+lab) handling sqlc queries + repo + service + handler; (C) one backend plan wiring routes and the daily dashboard aggregate; (D) client UI plans per domain; (E) nutritionist tracking view plan. Targeting 8–10 plans total at the project's standard granularity.

**Primary recommendation:** Implement all `local_id` upserts as `ON CONFLICT (local_id) DO NOTHING RETURNING *` with a follow-up `SELECT` on no-rows, never as INSERT-then-catch — this avoids advisory locks and is safe under concurrent submissions.

---

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Tracking data persistence | API / Backend | Database | All writes validated and authorized before touching DB |
| `local_id` deduplication | Database | API / Backend | `ON CONFLICT (local_id)` is a DB constraint; service layer reads RETURNING result |
| Nutritionist row-level auth | Database / SQL | API / Backend | SQL JOIN enforces ownership at query layer, not in handler |
| Daily dashboard aggregation | API / Backend | Frontend | Server computes today's summary; frontend renders |
| Medication checklist population | Frontend (reads active plan) | — | Plan already cached in `useClientPlanStore`; no extra API call needed |
| Sleep duration computation | Frontend | — | `wake_time - sleep_time` is pure arithmetic; no backend column |
| Water progress indicator | Frontend | — | Frontend: sum of today's `water_logs` vs. `active_plan.daily_water_target_ml` |
| Lab file storage | API / Backend | Filesystem | Go multipart handler → filesystem; path in DB |
| File MIME validation | API / Backend | — | Magic-byte sniff before persist; Content-Type header is untrusted |
| Weight chart rendering | Frontend (Browser) | — | Chart.js in browser; data fetched from chart-series endpoint |

---

## Standard Stack

### Core (already installed — no new backend dependencies needed)

| Library | Version (verified) | Purpose | Reference |
|---------|--------------------|---------|-----------|
| `github.com/gin-gonic/gin` | v1.12.0 | HTTP router + binding | `go.mod` [VERIFIED: go.mod] |
| `github.com/jackc/pgx/v5` | v5.9.2 | PostgreSQL driver + pgxpool | `go.mod` [VERIFIED: go.mod] |
| `github.com/jackc/pgerrcode` | 0.0.0-20220416 | Postgres error code constants | `go.mod` [VERIFIED: go.mod] |
| `github.com/gabriel-vasile/mimetype` | v1.4.13 | Magic-byte MIME sniffing | `go.mod` [VERIFIED: go.mod] |
| `github.com/google/uuid` | v1.6.0 | UUID generation for `local_id` | `go.mod` [VERIFIED: go.mod] |
| `github.com/rs/zerolog` | v1.35.0 | Structured logging | `go.mod` [VERIFIED: go.mod] |

### New Frontend Dependency (must install)

| Library | Purpose | Install |
|---------|---------|---------|
| `chart.js` | Line chart engine for weight history | `npm install chart.js vue-chartjs` |
| `vue-chartjs` | Vue 3 component wrappers for Chart.js | included above |

> **[VERIFIED: frontend/package.json]** — `chart.js` and `vue-chartjs` are absent from the current `package.json`. All other listed frontend libraries (`pinia`, `jalaali-js`, `dayjs`, `nuxt`, `tailwindcss`) are already installed.

### New Backend Configuration (env var)

Add to `backend/internal/config/config.go` and `.env.example`:

```
UPLOADS_DIR=/data/uploads   # base directory for lab-result file storage
```

---

## Architecture Patterns

### System Architecture Diagram

```
Client Mobile Browser
    │
    │  POST /api/client/food-logs        (local_id in body)
    │  POST /api/client/water-logs
    │  POST /api/client/sleep-logs
    │  POST /api/client/exercise-logs
    │  POST /api/client/medication-logs
    │  POST /api/client/body-measurements
    │  POST /api/client/lab-results      (multipart/form-data)
    │  GET  /api/client/tracking/daily   (today's dashboard summary)
    ▼
Gin Router  ──► middleware.Auth  ──► middleware.RoleGuard("client")
    │
    ▼
TrackingHandler (one handler struct per domain)
    │  bind JSON/multipart  →  call Service
    ▼
TrackingService
    │  sentinel errors (Persian)
    │  validate: at least file OR link (lab_results)
    │  validate: MIME type + size (lab_results)
    │  emit filesystem write (lab_results)
    ▼
TrackingRepository (interface)
    │  CreateOrUpsertByLocalID  ──►  ON CONFLICT (local_id) DO NOTHING  ──► pgx RETURNING *
    │  UpsertByUniqueKey        ──►  ON CONFLICT (client_id, date[, meal_id]) DO UPDATE
    │  ListByClientAndDateRange ──►  SELECT ... WHERE client_id = $1 AND date BETWEEN $2 AND $3
    │  DashboardSummary         ──►  raw pgx: JOINs today's logs to active plan
    ▼
PostgreSQL  (7 new tables in migration 000008)
    │
    ├── food_logs          UNIQUE (local_id), UNIQUE (client_id, date, meal_id)
    ├── water_logs         UNIQUE (local_id)
    ├── sleep_logs         UNIQUE (local_id), UNIQUE (client_id, date)
    ├── exercise_logs      UNIQUE (local_id)
    ├── medication_logs    UNIQUE (local_id)
    ├── body_measurements  UNIQUE (local_id), UNIQUE (client_id, date)
    └── lab_results        (no date unique; multiple per client per day allowed)

Nutritionist Browser
    │
    │  GET /api/nutritionist/clients/:id/food-logs?from=&to=
    │  GET /api/nutritionist/clients/:id/water-logs?from=&to=
    │  GET /api/nutritionist/clients/:id/sleep-logs?from=&to=
    │  GET /api/nutritionist/clients/:id/exercise-logs?from=&to=
    │  GET /api/nutritionist/clients/:id/medication-logs?from=&to=
    │  GET /api/nutritionist/clients/:id/body-measurements?from=&to=
    │  GET /api/nutritionist/clients/:id/lab-results
    │  GET /api/nutritionist/clients/:id/lab-results/:labId/download
    ▼
Gin Router ──► middleware.Auth ──► middleware.RoleGuard("nutritionist")
    ▼
TrackingHandler.NutritionistGetXxx
    ▼
TrackingRepository.ListForNutritionist
    │  SQL: JOIN users ON users.id = $client_id AND users.nutritionist_id = $nutritionist_id
    │  Only returns rows when JOIN succeeds — prevents cross-nutritionist data access
    ▼
PostgreSQL
```

### Recommended Project Structure

```
backend/
├── db/
│   ├── migrations/
│   │   ├── 000008_create_tracking.up.sql     ← all 7 tracking tables + 2 enums
│   │   └── 000008_create_tracking.down.sql
│   └── queries/
│       ├── food_logs.sql
│       ├── water_logs.sql
│       ├── sleep_logs.sql
│       ├── exercise_logs.sql
│       ├── medication_logs.sql
│       ├── body_measurements.sql
│       └── lab_results.sql
├── internal/
│   ├── model/dto/
│   │   └── tracking_dto.go                   ← all 7 request/response DTOs
│   ├── repository/
│   │   ├── sqlc/                              ← regenerated after 000008 + queries
│   │   └── tracking_repo.go                  ← TrackingRepository interface + impl
│   ├── service/
│   │   └── tracking_service.go               ← TrackingService (all 7 domains)
│   └── handler/
│       └── tracking_handler.go               ← TrackingHandler (all routes)
│
frontend/app/
├── stores/
│   ├── tracking.ts                           ← useTrackingStore (daily dashboard state)
│   ├── waterLog.ts                           ← useWaterLogStore
│   ├── foodLog.ts                            ← useFoodLogStore
│   ├── sleepLog.ts                           ← useSleepLogStore
│   ├── exerciseLog.ts                        ← useExerciseLogStore
│   ├── medicationLog.ts                      ← useMedicationLogStore
│   ├── bodyMeasurement.ts                    ← useBodyMeasurementStore
│   └── labResult.ts                          ← useLabResultStore
├── types/
│   └── tracking.types.ts                     ← TypeScript interfaces for all 7 domains
├── components/
│   └── tracking/
│       ├── DailyDashboard.vue
│       ├── WaterProgressBar.vue
│       ├── MedicationChecklistItem.vue       ← wraps plan/MedicationCard.vue
│       ├── WeightChart.vue                   ← Chart.js line chart component
│       ├── FoodLogMealCard.vue
│       ├── BodyMeasurementForm.vue
│       └── LabResultUploadForm.vue
├── pages/
│   ├── client/
│   │   └── tracking/
│   │       ├── index.vue                     ← daily dashboard
│   │       ├── water.vue
│   │       ├── food.vue
│   │       ├── sleep.vue
│   │       ├── exercise.vue
│   │       ├── medication.vue
│   │       ├── body.vue
│   │       └── lab-results.vue
│   └── nutritionist/
│       └── clients/
│           └── [clientId]/
│               └── tracking/
│                   ├── index.vue             ← tab container for nutritionist review
│                   ├── food.vue
│                   ├── water.vue
│                   ├── sleep.vue
│                   ├── exercise.vue
│                   ├── medication.vue
│                   ├── body.vue
│                   └── lab-results.vue
```

---

## Key Technical Patterns

### Pattern 1: `local_id` Idempotent Upsert (TRACK-12, D-10)

**What:** Every tracking POST accepts a client-generated UUID `local_id`. If the same `local_id` is submitted twice (e.g., offline replay in Phase 6), the second call returns the existing record instead of creating a duplicate.

**Implementation:** Two-step pattern — attempt INSERT with `ON CONFLICT (local_id) DO NOTHING RETURNING *`. If zero rows returned, the `local_id` was already used; run a follow-up `SELECT ... WHERE local_id = $1` and return that record.

**Why not `ON CONFLICT DO UPDATE`?** For `local_id` upserts we want idempotency, not overwrite — the original data is authoritative. Overwriting would corrupt the record if the payload changes (shouldn't happen, but defensive). For `unique key` upserts (sleep, food log, body measurements) we DO use `ON CONFLICT (...) DO UPDATE SET ...` because "last write wins" is the intended semantics.

```sql
-- db/queries/water_logs.sql
-- name: CreateWaterLog :one
INSERT INTO water_logs (id, client_id, local_id, date, amount_ml, logged_time)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
ON CONFLICT (local_id) DO NOTHING
RETURNING *;

-- name: GetWaterLogByLocalID :one
SELECT * FROM water_logs WHERE local_id = $1 AND client_id = $2;
```

```go
// backend/internal/repository/tracking_repo.go
func (r *trackingRepository) CreateWaterLog(ctx context.Context, p CreateWaterLogParams) (*sqlc.WaterLog, error) {
    row, err := r.q.CreateWaterLog(ctx, sqlc.CreateWaterLogParams{...})
    if err != nil {
        return nil, err
    }
    // ON CONFLICT DO NOTHING returns zero rows when local_id already exists
    if row == nil {
        existing, err := r.q.GetWaterLogByLocalID(ctx, sqlc.GetWaterLogByLocalIDParams{
            LocalID:  p.LocalID,
            ClientID: p.ClientID,
        })
        if err != nil {
            return nil, err
        }
        return &existing, nil
    }
    return row, nil
}
```

> **Note on sqlc + `ON CONFLICT DO NOTHING`:** When `DO NOTHING` fires, pgx returns zero rows and sqlc's `:one` annotation would normally error with `pgx.ErrNoRows`. Declare the query as `:one` and handle `pgx.ErrNoRows` as the "conflict detected" branch, then do the follow-up SELECT. [VERIFIED: backend/internal/repository/diet_plan_repo.go pattern for pgx.ErrNoRows handling]

### Pattern 2: Unique-Key Upsert (sleep_logs, food_logs, body_measurements)

These tables have a business-logic unique key beyond `local_id`:
- `sleep_logs`: `(client_id, date)` — one sleep record per day
- `food_logs`: `(client_id, date, meal_id)` — one log per meal per day
- `body_measurements`: `(client_id, date)` — one measurement record per day

```sql
-- db/queries/sleep_logs.sql
-- name: UpsertSleepLog :one
INSERT INTO sleep_logs (id, client_id, local_id, date, sleep_time, wake_time, quality, notes)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (client_id, date)
DO UPDATE SET
    sleep_time = EXCLUDED.sleep_time,
    wake_time  = EXCLUDED.wake_time,
    quality    = EXCLUDED.quality,
    notes      = EXCLUDED.notes,
    updated_at = NOW()
RETURNING *;
```

The `local_id` UNIQUE constraint on these tables still fires for true duplicate offline submissions. The business-key upsert handles the "user edited and re-submitted same day" case.

### Pattern 3: Nutritionist Row-Level Authorization (D-11)

**Never trust `client_id` from the URL alone.** Enforce ownership in SQL by joining to `users`:

```sql
-- db/queries/water_logs.sql
-- name: ListWaterLogsForNutritionist :many
SELECT wl.*
FROM water_logs wl
JOIN users u ON u.id = wl.client_id
    AND u.nutritionist_id = @nutritionist_id   -- ownership check in JOIN condition
WHERE wl.client_id = @client_id
  AND wl.date BETWEEN @from_date AND @to_date
ORDER BY wl.date DESC, wl.logged_time DESC;
```

This is identical to how `diet_plan_repo.go` enforces nutritionist ownership via `GetDietPlanByIDParams{NutritionistID: ...}` in all nutritionist read queries. [VERIFIED: backend/internal/repository/diet_plan_repo.go `GetPlanByID` requires both `planID` and `nutritionistID`]

### Pattern 4: Daily Dashboard Aggregate (D-14, D-13)

The dashboard at `/api/client/tracking/daily?date=YYYY-MM-DD` returns a single summary object. This is a raw pgx query (per D-13), not sqlc, because it joins five tables in one round trip:

```sql
-- Used in trackingRepository.GetDailyDashboard (raw pgx, not sqlc)
SELECT
    -- Water
    COALESCE(SUM(wl.amount_ml), 0)   AS water_total_ml,
    -- Sleep (latest entry)
    sl.sleep_time, sl.wake_time, sl.quality,
    -- Exercise count
    COUNT(DISTINCT el.id)             AS exercise_count,
    -- Medication taken count
    COUNT(DISTINCT ml.id)             AS medication_taken_count,
    -- Body measurement logged today?
    (bm.id IS NOT NULL)               AS body_logged_today
FROM diet_plans dp
LEFT JOIN water_logs wl       ON wl.client_id = $1 AND wl.date = $2
LEFT JOIN sleep_logs sl        ON sl.client_id = $1 AND sl.date = $2
LEFT JOIN exercise_logs el     ON el.client_id = $1 AND el.date = $2
LEFT JOIN medication_logs ml   ON ml.client_id = $1 AND ml.date = $2
LEFT JOIN body_measurements bm ON bm.client_id = $1 AND bm.date = $2
WHERE dp.client_id = $1 AND dp.status = 'active'
GROUP BY sl.sleep_time, sl.wake_time, sl.quality, bm.id, dp.daily_water_target_ml;
```

Water target comes from the active plan's `daily_water_target_ml`. Food log status (how many meals logged) is fetched separately via the sqlc `ListFoodLogsByDate` query to keep the aggregate simple.

### Pattern 5: Medication Checklist from Plan Prescriptions (D-18, TRACK-08)

**No extra API call needed.** The `useClientPlanStore` already fetches the full plan aggregate including `medications[]`. Each `PlanMedicationResponse` has a `times: string[]` array (e.g., `["08:00", "13:00", "21:00"]`).

**Checklist construction (frontend):**
```typescript
// frontend/app/stores/medicationLog.ts
const checklistItems = computed(() => {
  const plan = clientPlanStore.activePlan
  if (!plan) return []
  return plan.medications.flatMap(med =>
    (med.times ?? []).map(time => ({
      prescribedMedicationId: med.id,
      medicationName: med.medication_name,
      dosage: med.dosage,
      time,
      takenAt: todayLogs.value.find(
        l => l.prescribed_medication_id === med.id && l.taken_at === time
      )?.taken_at ?? null,
    }))
  )
})
```

When the client taps a checklist item, the store calls:
```typescript
POST /api/client/medication-logs
{
  local_id: crypto.randomUUID(),
  date: today,
  prescribed_medication_id: item.prescribedMedicationId,
  medication_name: item.medicationName,
  dosage: item.dosage,
  taken_at: item.time,   // or current time if tapped outside the scheduled window
  is_self_reported: false
}
```

> **[VERIFIED: backend/internal/repository/sqlc/models.go]** `PlanMedication.Times` is stored as `[]byte` (JSONB) in the DB model and deserialized to `[]string` in `backend/internal/model/dto/diet_plan_dto.go PlanMedicationResponse.Times`.

### Pattern 6: Chart.js Weight History with Shamsi Labels (D-23, TRACK-10)

Chart.js is installed via `npm install chart.js vue-chartjs`. The `WeightChart.vue` component uses `vue-chartjs` Line chart with Shamsi-converted date labels:

```typescript
// frontend/app/components/tracking/WeightChart.vue
<script setup lang="ts">
import { Line } from 'vue-chartjs'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Tooltip } from 'chart.js'
import { useShamsiDate } from '~/composables/useShamsiDate'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip)

const { formatShamsi } = useShamsiDate()
const props = defineProps<{ measurements: BodyMeasurementPoint[] }>()

const chartData = computed(() => ({
  labels: props.measurements.map(m => formatShamsi(m.date, 'short')),  // Jalali labels
  datasets: [{
    label: 'وزن (کیلوگرم)',
    data: props.measurements.map(m => m.weight_kg),
    borderColor: '#22c55e',
    tension: 0.3,
  }]
}))
</script>
```

The backend chart-series endpoint returns measurements **oldest-first** (per D-12) so Chart.js renders a left-to-right time progression without frontend reversal.

```sql
-- db/queries/body_measurements.sql
-- name: GetWeightHistory :many
SELECT date, weight_kg FROM body_measurements
WHERE client_id = $1 AND weight_kg IS NOT NULL
ORDER BY date ASC;  -- oldest-first for chart rendering
```

> **RTL caveat:** Chart.js does not natively flip axes for `dir="rtl"`. For a simple line chart of weight-over-time this is acceptable (time always flows left→right universally). If the product owner requires RTL axis mirroring later, use `chart.js` plugin `chartjs-plugin-datalabels` or a custom axis plugin — defer to Phase 7 polish. [ASSUMED — based on Chart.js RTL behavior knowledge; acceptable for Phase 4 scope]

### Pattern 7: Lab Result File Upload (LAB-01–05, D-25–26)

Go multipart handling with `gabriel-vasile/mimetype` (already in `go.mod`):

```go
// backend/internal/handler/tracking_handler.go (UploadLabResult)
func (h *TrackingHandler) UploadLabResult(c *gin.Context) {
    clientID := uuid.MustParse(c.GetString("user_id"))

    // 1. Parse multipart (max 10MB)
    if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "فایل بیش از حد مجاز است"})
        return
    }

    var req dto.CreateLabResultRequest
    if err := c.ShouldBind(&req); err != nil { // binds title, lab_type, test_date, external_link
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
        return
    }

    // 2. Handle optional file upload
    var filePath, origFilename, mimeType string
    var fileSize int64
    if fh, err := c.FormFile("file"); err == nil {
        f, _ := fh.Open()
        defer f.Close()

        // 3. Detect MIME from magic bytes (not Content-Type header)
        kind, _ := mimetype.DetectReader(f)
        if !allowedMIME(kind.String()) {
            c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "فرمت فایل مجاز نیست"})
            return
        }
        f.Seek(0, 0) // reset after MIME detection

        // 4. Save with UUID filename
        ext := extensionFromMIME(kind.String()) // ".pdf", ".jpg", ".png"
        uuidName := uuid.New().String() + ext
        dest := filepath.Join(h.uploadsDir, "lab-results", clientID.String(), uuidName)
        os.MkdirAll(filepath.Dir(dest), 0755)
        // save file...
        filePath = dest
        origFilename = fh.Filename
        mimeType = kind.String()
        fileSize = fh.Size
    }

    resp, err := h.trackingService.CreateLabResult(c.Request.Context(), clientID, dto.CreateLabResultRequest{
        FilePath: filePath, OriginalFilename: origFilename, MimeType: mimeType, FileSize: fileSize,
        ExternalLink: req.ExternalLink, Title: req.Title, LabType: req.LabType, TestDate: req.TestDate,
    })
    // ...
}
```

The `uploadsDir` is passed to the handler from `config.UploadsDir`. Download endpoint adds:

```go
c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, record.OriginalFilename))
c.File(record.FilePath)
```

### Pattern 8: Pinia Store Structure (mirrors clientPlan.ts)

Each tracking domain gets its own store following the `useClientPlanStore` composition-API pattern:

```typescript
// frontend/app/stores/waterLog.ts  — representative example
export const useWaterLogStore = defineStore('waterLog', () => {
  const logs = ref<WaterLogEntry[]>([])
  const totalMl = computed(() => logs.value.reduce((s, l) => s + l.amount_ml, 0))
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchToday() {
    loading.value = true
    try {
      const { apiFetch } = useApi()
      const today = new Date().toISOString().slice(0, 10)
      logs.value = await apiFetch<WaterLogEntry[]>(`/client/water-logs?date=${today}`)
    } catch (e: unknown) {
      error.value = (e as any).data?.error ?? 'خطا در بارگذاری'
    } finally {
      loading.value = false
    }
  }

  async function addLog(amountMl: number, loggedTime?: string) {
    const { apiFetch } = useApi()
    const entry = await apiFetch<WaterLogEntry>('/client/water-logs', {
      method: 'POST',
      body: JSON.stringify({
        local_id: crypto.randomUUID(),   // TRACK-12: client-generated UUID
        date: new Date().toISOString().slice(0, 10),
        amount_ml: amountMl,
        logged_time: loggedTime,
      }),
    })
    logs.value.push(entry)
  }

  function $reset() { logs.value = []; error.value = null }
  return { logs, totalMl, loading, error, fetchToday, addLog, $reset }
})
```

> **[VERIFIED: frontend/app/stores/clientPlan.ts]** — exact pattern reproduced: `defineStore('name', () => {...})`, `apiFetch`, error handling with `(e as { data?: { error?: string } })`, and `$reset()`.

---

## Migration Design: 000008_create_tracking.up.sql

Next migration number: **000008** (current highest: `000007_create_diet_plans.up.sql`). [VERIFIED: backend/db/migrations/]

```sql
-- Migration 000008: Client Tracking Suite tables
-- Creates: sleep_quality enum, lab_result_type enum,
--          food_logs, water_logs, sleep_logs, exercise_logs,
--          medication_logs, body_measurements, lab_results

CREATE TYPE sleep_quality     AS ENUM ('good', 'fair', 'poor');
CREATE TYPE lab_result_type   AS ENUM ('blood_test', 'urine_test', 'thyroid', 'hormone', 'allergy', 'other');

-- food_logs: one per (client, date, meal) — upsert on meal_id, idempotent on local_id
CREATE TABLE food_logs (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id          UUID NOT NULL REFERENCES users(id),
    local_id           UUID NOT NULL,
    date               DATE NOT NULL,
    meal_id            UUID NOT NULL REFERENCES meals(id),
    selected_option_id UUID REFERENCES meal_options(id),
    is_skipped         BOOLEAN NOT NULL DEFAULT false,
    notes              TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_food_logs_local_id       UNIQUE (local_id),
    CONSTRAINT uq_food_logs_client_date_meal UNIQUE (client_id, date, meal_id)
);
CREATE INDEX idx_food_logs_client_date ON food_logs (client_id, date);

-- water_logs: multiple per day, idempotent on local_id
CREATE TABLE water_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id   UUID NOT NULL REFERENCES users(id),
    local_id    UUID NOT NULL UNIQUE,
    date        DATE NOT NULL,
    amount_ml   INTEGER NOT NULL CHECK (amount_ml > 0),
    logged_time TIME,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_water_logs_client_date ON water_logs (client_id, date);

-- sleep_logs: one per day — upsert on (client_id, date)
CREATE TABLE sleep_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id   UUID NOT NULL REFERENCES users(id),
    local_id    UUID NOT NULL UNIQUE,
    date        DATE NOT NULL,
    sleep_time  TIME NOT NULL,
    wake_time   TIME NOT NULL,
    quality     sleep_quality NOT NULL,
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_sleep_logs_client_date UNIQUE (client_id, date)
);
CREATE INDEX idx_sleep_logs_client_date ON sleep_logs (client_id, date);

-- exercise_logs: multiple per day, idempotent on local_id
CREATE TABLE exercise_logs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id        UUID NOT NULL REFERENCES users(id),
    local_id         UUID NOT NULL UNIQUE,
    date             DATE NOT NULL,
    exercise_name    VARCHAR(200) NOT NULL,
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
    calories_burned  INTEGER CHECK (calories_burned >= 0),
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_exercise_logs_client_date ON exercise_logs (client_id, date);

-- medication_logs: multiple per day; links to plan_medications or self-reported
CREATE TABLE medication_logs (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id                UUID NOT NULL REFERENCES users(id),
    local_id                 UUID NOT NULL UNIQUE,
    date                     DATE NOT NULL,
    prescribed_medication_id UUID REFERENCES plan_medications(id),
    medication_name          VARCHAR(200) NOT NULL,
    dosage                   VARCHAR(100),
    taken_at                 TIME NOT NULL,
    notes                    TEXT,
    is_self_reported         BOOLEAN NOT NULL DEFAULT false,
    created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_medication_logs_client_date ON medication_logs (client_id, date);

-- body_measurements: one per day — upsert on (client_id, date)
CREATE TABLE body_measurements (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id    UUID NOT NULL REFERENCES users(id),
    local_id     UUID NOT NULL UNIQUE,
    date         DATE NOT NULL,
    weight_kg    NUMERIC(5,2) CHECK (weight_kg > 0),
    waist_cm     NUMERIC(5,2) CHECK (waist_cm > 0),
    hip_cm       NUMERIC(5,2) CHECK (hip_cm > 0),
    abdomen_cm   NUMERIC(5,2) CHECK (abdomen_cm > 0),
    thigh_cm     NUMERIC(5,2) CHECK (thigh_cm > 0),
    chest_cm     NUMERIC(5,2) CHECK (chest_cm > 0),
    wrist_cm     NUMERIC(5,2) CHECK (wrist_cm > 0),
    recorded_by  UUID NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_body_measurements_client_date UNIQUE (client_id, date)
);
CREATE INDEX idx_body_measurements_client_date ON body_measurements (client_id, date);

-- lab_results: multiple per client, no date uniqueness (D-09)
CREATE TABLE lab_results (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id         UUID NOT NULL REFERENCES users(id),
    uploaded_by       UUID NOT NULL REFERENCES users(id),
    title             VARCHAR(200) NOT NULL,
    lab_type          lab_result_type NOT NULL,
    test_date         DATE NOT NULL,
    file_path         TEXT,
    external_link     TEXT,
    original_filename VARCHAR(255),
    mime_type         VARCHAR(100),
    file_size_bytes   BIGINT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_lab_results_client_id ON lab_results (client_id);
```

---

## Route Group Design

Extend `backend/cmd/api/main.go` following existing group patterns. [VERIFIED: backend/cmd/api/main.go]

```go
// Add to existing `client` group (already has middleware.Auth + RoleGuard("client"))
client.GET("/tracking/daily",        trackingHandler.GetDailyDashboard)
client.POST("/food-logs",            trackingHandler.LogFood)
client.GET("/food-logs",             trackingHandler.ListFoodLogs)       // ?date=
client.POST("/water-logs",           trackingHandler.LogWater)
client.GET("/water-logs",            trackingHandler.ListWaterLogs)      // ?date=
client.POST("/sleep-logs",           trackingHandler.UpsertSleep)
client.GET("/sleep-logs",            trackingHandler.GetSleepLog)        // ?date=
client.POST("/exercise-logs",        trackingHandler.LogExercise)
client.GET("/exercise-logs",         trackingHandler.ListExerciseLogs)   // ?date=
client.POST("/medication-logs",      trackingHandler.LogMedication)
client.GET("/medication-logs",       trackingHandler.ListMedicationLogs) // ?date=
client.POST("/body-measurements",    trackingHandler.UpsertBodyMeasurement)
client.GET("/body-measurements",     trackingHandler.GetBodyMeasurement) // ?date=
client.GET("/body-measurements/history", trackingHandler.GetMeasurementHistory) // ?from=&to=
client.POST("/lab-results",          trackingHandler.UploadLabResult)
client.GET("/lab-results",           trackingHandler.ListLabResults)

// Add to existing `nutri` group (already has middleware.Auth + RoleGuard("nutritionist"))
// Nutritionist tracking reads for their clients:
nutri.GET("/clients/:clientId/tracking/food-logs",         trackingHandler.NutriListFoodLogs)
nutri.GET("/clients/:clientId/tracking/water-logs",        trackingHandler.NutriListWaterLogs)
nutri.GET("/clients/:clientId/tracking/sleep-logs",        trackingHandler.NutriListSleepLogs)
nutri.GET("/clients/:clientId/tracking/exercise-logs",     trackingHandler.NutriListExerciseLogs)
nutri.GET("/clients/:clientId/tracking/medication-logs",   trackingHandler.NutriListMedicationLogs)
nutri.GET("/clients/:clientId/tracking/body-measurements", trackingHandler.NutriListBodyMeasurements)
nutri.GET("/clients/:clientId/tracking/weight-history",    trackingHandler.NutriGetWeightHistory)
nutri.POST("/clients/:clientId/body-measurements",         trackingHandler.NutriUpsertBodyMeasurement)
nutri.GET("/clients/:clientId/lab-results",                trackingHandler.NutriListLabResults)
nutri.GET("/clients/:clientId/lab-results/:labId/download",trackingHandler.NutriDownloadLabResult)
```

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| MIME type detection from file bytes | Custom magic-byte reader | `gabriel-vasile/mimetype` (already in go.mod) | Handles PDF/JPG/PNG and 1,200+ other types; Content-Type header is spoofable |
| UUID generation for `local_id` | `fmt.Sprintf("%d", rand...)` | `crypto.randomUUID()` (browser) / `github.com/google/uuid` (Go) | Already used throughout codebase; cryptographically random |
| Jalali date labels for Chart.js | Custom date formatter | `useShamsiDate()` composable in `frontend/app/composables/useShamsiDate.ts` | Already exists and tested in Phase 3 plan views |
| File download with safe headers | Setting headers manually | Gin `c.Header(...)` + `c.File(...)` with `Content-Disposition: attachment` | One-liner; avoids content-sniffing attacks |
| Turkish/Farsi month names in chart axis | Custom lookup table | `useShamsiDate.formatShamsi()` already returns Persian month names | Already exists; reuse |
| Date arithmetic (sleep duration) | Backend computed column | Frontend: `(wakeH*60+wakeM) - (sleepH*60+sleepM)` in component | Zero DB cost; display-only value |

**Key insight:** Phase 3 established every pattern this phase needs — the only genuine additions are the 7 new tables, the file upload handler, and Chart.js installation.

---

## Common Pitfalls

### Pitfall 1: Trusting `local_id` for Unique-Key Upserts

**What goes wrong:** Developer uses `ON CONFLICT (local_id) DO UPDATE SET sleep_time=...` for `sleep_logs` instead of `ON CONFLICT (client_id, date) DO UPDATE`. The conflict on `local_id` never fires when the user edits today's sleep record (because the new submission has a new `local_id`), creating a second record for the same date.

**Why it happens:** Confusion between idempotency semantics (same offline submission → same record) and business-key semantics (same day → same record, update it).

**How to avoid:** Two separate `UNIQUE` constraints on every table that has a business-unique key:
1. `UNIQUE (local_id)` — for offline deduplication (DO NOTHING)
2. `UNIQUE (client_id, date[, meal_id])` — for business logic (DO UPDATE)

Each upsert query targets the relevant constraint. [VERIFIED: migration design above]

### Pitfall 2: `ON CONFLICT DO NOTHING` Returning No Rows in sqlc

**What goes wrong:** sqlc `:one` query panics or returns `pgx.ErrNoRows` when `DO NOTHING` fires, crashing the handler instead of returning the existing record.

**How to avoid:** Use the two-step pattern: `:one` → catch `pgx.ErrNoRows` → run the follow-up `:one` SELECT by `local_id`. This is idiomatic and already done in diet_plan_repo.go for similar edge cases. [VERIFIED: backend/internal/repository/diet_plan_repo.go error handling]

### Pitfall 3: File Path Stored Without Config-Relative Root

**What goes wrong:** Lab result file stored as absolute path like `/data/uploads/lab-results/...` in the database. When `UPLOADS_DIR` changes or the container volume mounts differently, all existing paths break.

**How to avoid:** Store only the **relative path** from `UPLOADS_DIR` root, e.g., `lab-results/{client_id}/{uuid}.pdf`. The download handler joins `cfg.UploadsDir + "/" + record.FilePath`. This mirrors how Phase 2's food photos (if any) would be stored. [ASSUMED — best practice for configurable file storage]

### Pitfall 4: MIME Type Validated from Content-Type Header Only

**What goes wrong:** Handler checks `file.Header.Get("Content-Type") == "application/pdf"` — trivially bypassed by renaming an `.exe` to `.pdf`.

**How to avoid:** Use `mimetype.DetectReader(f)` after opening the multipart file. The `gabriel-vasile/mimetype` library reads the first 512 bytes and identifies the actual format. [VERIFIED: go.mod — `github.com/gabriel-vasile/mimetype v1.4.13` already imported]

### Pitfall 5: Chart.js Leaks Globally if Not Tree-Shaken

**What goes wrong:** Import entire Chart.js via `import ChartJS from 'chart.js/auto'` — includes all chart types (~200KB gzipped). Bloats the client bundle.

**How to avoid:** Register only needed controllers:
```typescript
import { Chart, CategoryScale, LinearScale, PointElement, LineElement, Tooltip } from 'chart.js'
Chart.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip)
```
This keeps the chart bundle under ~30KB for a simple line chart. [ASSUMED — standard Chart.js tree-shaking guidance]

### Pitfall 6: Medication Checklist Re-Fetches Plan on Every Log

**What goes wrong:** Each `POST /client/medication-logs` triggers a re-fetch of the full plan aggregate to refresh the checklist state, causing 500ms+ round trips per tap.

**How to avoid:** The `useMedicationLogStore` maintains a local `todayLogs` ref. On successful POST, push the response directly to `todayLogs` without re-fetching the plan. The checklist computed property (Pattern 5 above) derives its checked state from `todayLogs`, so it updates immediately. [VERIFIED: frontend/app/stores/clientPlan.ts — same optimistic-update pattern used for plan state]

### Pitfall 7: Missing `updated_at` on Tables That Are Upserted

**What goes wrong:** `sleep_logs` and `body_measurements` use `ON CONFLICT DO UPDATE` but forget to include `updated_at = NOW()` in the SET clause. The record retains the original `created_at` as the modification timestamp, making audit trails inaccurate.

**How to avoid:** Include `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()` in the table definition and always set it explicitly in the upsert's DO UPDATE clause. [VERIFIED: migration design above already includes this]

---

## Validation Architecture

**Framework:** Vitest 3.x (frontend), Go testing + `//go:build integration` tags (backend)

### Test Framework

| Property | Value |
|----------|-------|
| Frontend framework | Vitest 3.x + happy-dom |
| Frontend config | `frontend/vitest.config.ts` [VERIFIED] |
| Frontend run command | `cd frontend && npm test` |
| Backend integration tag | `//go:build integration` (see `diet_plan_service_test.go`) |
| Backend unit run | `cd backend && go test ./internal/...` |
| Backend integration run | `cd backend && go test -tags=integration ./internal/...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | File |
|--------|----------|-----------|------|
| TRACK-12 | `local_id` idempotent: POST same UUID twice → same record returned | Integration | `backend/internal/repository/tracking_repo_test.go` (Wave 0 stub) |
| TRACK-11 | Body measurement last-write-wins: second POST for same date updates fields | Integration | `backend/internal/repository/tracking_repo_test.go` |
| TRACK-05 | Sleep duration = wake minus sleep | Unit (frontend) | `frontend/tests/useSleepDuration.test.ts` (Wave 0 stub) |
| TRACK-13 | Nutritionist cannot access another nutritionist's client tracking | Integration | `backend/internal/repository/tracking_repo_test.go` |
| LAB-03 | MIME sniffing rejects non-PDF/JPG/PNG even with correct Content-Type | Unit (backend) | `backend/internal/service/tracking_service_test.go` |
| LAB-02 | Neither file nor link provided → 400 error | Unit (backend) | `backend/internal/service/tracking_service_test.go` |
| TRACK-03 | Water total = sum of today's logs; compares to plan target | Unit (frontend) | `frontend/tests/useWaterLog.test.ts` (Wave 0 stub) |

### Wave 0 Gaps (must create before implementation)

- [ ] `backend/internal/repository/tracking_repo_test.go` — stubs for TRACK-12, TRACK-11, TRACK-13 (integration build tag)
- [ ] `backend/internal/service/tracking_service_test.go` — stubs for LAB-02, LAB-03
- [ ] `frontend/tests/useWaterLog.test.ts` — water total computed test stub
- [ ] `frontend/tests/useSleepDuration.test.ts` — duration arithmetic test stub

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| PostgreSQL | All tracking tables | ✓ | 16 (Docker Compose) | — |
| `gabriel-vasile/mimetype` | Lab file validation | ✓ | v1.4.13 in go.mod | — |
| `chart.js` + `vue-chartjs` | Weight chart | ✗ | Not installed | None — must install |
| `UPLOADS_DIR` env var | Lab file storage | ✗ | Not in config.go | Add to Config struct + .env.example |
| `crypto.randomUUID()` | Client-side local_id | ✓ | Browser native (Chrome 92+, FF 95+, Safari 15.4+) | `uuid` npm package fallback |

**Missing with no fallback:**
- `chart.js` + `vue-chartjs`: Install required before implementing WeightChart.vue
- `UPLOADS_DIR` config field: Must add to `backend/internal/config/config.go` before lab upload handler

---

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V4 Access Control | YES — critical | Repository-level JOIN with `nutritionist_id`; never handler-level `if` check |
| V5 Input Validation | YES | Gin binding tags on all DTOs; service sentinel errors |
| V12 File Upload | YES | `mimetype.DetectReader()` magic bytes; size limit 10MB; UUID filenames; `Content-Disposition: attachment` |
| V2 Authentication | Inherited | JWT middleware already applied to all route groups |
| V6 Cryptography | NO — no new crypto | UUID filenames are not secret; no new keys |

### Known Threat Patterns for This Stack

| Pattern | STRIDE | Mitigation |
|---------|--------|------------|
| Client submits `client_id` of another client in POST body | Elevation of Privilege | Use `c.GetString("user_id")` from JWT — never trust body `client_id` for writes |
| Nutritionist accesses another nutritionist's client via URL param | Elevation of Privilege | JOIN enforcement: `users.nutritionist_id = $authenticated_nutritionist_id` |
| Malicious file disguised as PDF (content-type spoofed) | Tampering | `gabriel-vasile/mimetype` magic-byte detection |
| Path traversal in file download (`../../etc/passwd`) | Information Disclosure | Download handler uses stored UUID path from DB; never user-supplied path |
| Oversized file exhausting disk | Denial of Service | `ParseMultipartForm(10 << 20)` before reading |
| Serving file inline (XSS via SVG/HTML file) | XSS | `Content-Disposition: attachment` always; MIME type enforcement blocks non-PDF/JPG/PNG |

---

## Practical Plan Decomposition

Recommended 8-plan split at standard granularity:

| Plan | Title | Scope |
|------|-------|-------|
| **04-01** | Migration + sqlc skeleton | `000008_create_tracking.up/down.sql`, all 7 sqlc query files, `sqlc generate`, `tracking_dto.go` |
| **04-02** | TrackingRepository + Wave 0 tests | `tracking_repo.go` (interface + impl for all 7 domains), integration test stubs |
| **04-03** | TrackingService + route wiring | `tracking_service.go` (all 7 domains + lab file handling + dashboard aggregate), handler, `main.go` additions |
| **04-04** | Client daily dashboard UI | `frontend/app/pages/client/tracking/index.vue`, `useTrackingStore`, `DailyDashboard.vue`, `tracking.types.ts` |
| **04-05** | Client food + water + sleep UI | Pages + stores for food log (meal picker), water log (quick-add + progress bar), sleep log (time pickers) |
| **04-06** | Client exercise + medication UI | Pages + stores for exercise log, medication checklist (`MedicationChecklistItem.vue` reusing `MedicationCard.vue`) |
| **04-07** | Client body measurements + weight chart | Body measurement form, `WeightChart.vue` (chart.js install + component), history page |
| **04-08** | Client lab results + nutritionist tracking views | Lab result upload form, lab result list; nutritionist tab-based tracking history pages for all 7 domains |

> The nutritionist view (plan 04-08) is intentionally last because it reuses all the list components built in 04-05/06/07 — it's primarily a data-fetch + layout composition task.

---

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Chart.js does not natively flip x-axis for RTL in Phase 4; left-to-right time flow is acceptable | Pattern 6 | Product owner may require RTL axis — add chartjs plugin; low effort |
| A2 | Relative file path stored in `lab_results.file_path` (not absolute) | Pattern 7 / Pitfall 3 | Absolute paths still work but break if UPLOADS_DIR changes |
| A3 | `crypto.randomUUID()` is available on target mobile browsers (Chrome 92+, Safari 15.4+) | Pattern 8 | If older iOS is targeted, add `uuid` npm package as polyfill |

---

## Sources

### Primary (HIGH confidence — verified in this codebase)
- `backend/db/migrations/000007_create_diet_plans.up.sql` — migration numbering convention and column patterns
- `backend/internal/repository/diet_plan_repo.go` — raw pgx aggregate query pattern, `pgx.ErrNoRows` handling, batch query idiom
- `backend/internal/repository/sqlc/models.go` — existing enum types (`DietPlanStatus`, `MedicationForm`, etc.) and struct patterns
- `backend/internal/service/diet_plan_service.go` — sentinel error pattern, service struct, Persian error strings
- `backend/cmd/api/main.go` — route group structure, middleware chaining, dependency injection
- `backend/internal/config/config.go` — env var config loading pattern to extend for `UPLOADS_DIR`
- `backend/go.mod` — confirmed `gabriel-vasile/mimetype v1.4.13` already present
- `frontend/app/stores/clientPlan.ts` — Pinia store composition API pattern
- `frontend/app/composables/useApi.ts` — fetch wrapper + 401 refresh pattern
- `frontend/app/composables/useShamsiDate.ts` — Jalali date conversion for chart labels
- `frontend/app/components/plan/MedicationCard.vue` — reusable medication display component
- `frontend/package.json` — confirmed `chart.js` and `vue-chartjs` absent; must install
- `frontend/vitest.config.ts` — test framework configuration
- `.planning/research/STACK.md` — Chart.js and vue-chartjs confirmed as standard stack choice
- `.planning/phases/04-client-tracking-suite/04-CONTEXT.md` — all locked decisions D-01 through D-28

### Secondary (MEDIUM confidence — from docs/phases.md and REQUIREMENTS.md)
- `docs/phases.md §4.1–4.3` — scope table, validation checklist, implementation guidance
- `.planning/REQUIREMENTS.md §Client Tracking, §Lab Results` — TRACK-01–13, LAB-01–05

---

## Metadata

**Confidence breakdown:**
- Schema design: HIGH — all decisions locked in CONTEXT.md; migration pattern verified from 000007
- `local_id` upsert: HIGH — pattern verified from pgx docs + existing error handling in diet_plan_repo.go
- Authorization: HIGH — verified from existing `GetPlanByID` nutritionist-ownership JOIN pattern
- Chart.js integration: MEDIUM — library confirmed in STACK.md; one assumption about RTL axis behavior
- File upload: HIGH — `gabriel-vasile/mimetype` in go.mod; pattern is standard Go multipart

**Research date:** 2026-04-20
**Valid until:** 2026-05-20 (stable Go/Nuxt versions; Chart.js API stable)
