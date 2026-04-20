# Phase 4: Client Tracking Suite - Context

**Gathered:** 2026-04-20
**Status:** Ready for planning
**Mode:** --auto (sub-agent, no interactive session)

<domain>
## Phase Boundary

Phase 4 delivers the full client tracking layer on top of the completed Diet Plan Engine: food intake, water, sleep, exercise, medication intake, body measurements, weight trends, and lab results. Clients can record daily activity from a mobile-first dashboard, and nutritionists can review comprehensive client tracking history with date filtering and charts.

Offline queueing and background sync are explicitly deferred to Phase 6, but this phase must lay the server-side foundation for that work by adding `local_id` deduplication to all tracking writes.

</domain>

<decisions>
## Implementation Decisions

### Tracking Data Model

- **D-01:** Create seven Phase 4 tables: `food_logs`, `water_logs`, `sleep_logs`, `exercise_logs`, `medication_logs`, `body_measurements`, and `lab_results`.
- **D-02:** Every writeable Phase 4 table includes `id UUID`, `client_id UUID`, `local_id UUID`, `created_at TIMESTAMPTZ`, and `updated_at TIMESTAMPTZ` where updates are possible. `local_id` is `UNIQUE NOT NULL` for idempotent deduplication in Phase 6.
- **D-03:** `food_logs` stores `date DATE`, `meal_id UUID`, `selected_option_id UUID NULL`, `is_skipped BOOLEAN NOT NULL DEFAULT false`, and `notes TEXT NULL`. Exactly one log per `(client_id, date, meal_id)`; later submissions update the same logical record.
- **D-04:** `water_logs` stores `date DATE`, `amount_ml INTEGER`, `logged_time TIME NULL`. Multiple entries per day are allowed.
- **D-05:** `sleep_logs` stores `date DATE`, `sleep_time TIME`, `wake_time TIME`, `quality sleep_quality`, and `notes TEXT NULL`. Unique constraint on `(client_id, date)`; writes behave as upsert.
- **D-06:** `exercise_logs` stores `date DATE`, `exercise_name VARCHAR(200)`, `duration_minutes INTEGER`, `calories_burned INTEGER NULL`, and `notes TEXT NULL`. Multiple entries per day are allowed.
- **D-07:** `medication_logs` stores `date DATE`, `prescribed_medication_id UUID NULL`, `medication_name VARCHAR(200)`, `dosage VARCHAR(100) NULL`, `taken_at TIME`, `notes TEXT NULL`, and `is_self_reported BOOLEAN NOT NULL DEFAULT false`. Prescribed taps fill `prescribed_medication_id`; manual entries set `is_self_reported=true`.
- **D-08:** `body_measurements` stores `date DATE`, `weight_kg NUMERIC(5,2) NULL`, `waist_cm NUMERIC(5,2) NULL`, `hip_cm NUMERIC(5,2) NULL`, `abdomen_cm NUMERIC(5,2) NULL`, `thigh_cm NUMERIC(5,2) NULL`, `chest_cm NUMERIC(5,2) NULL`, `wrist_cm NUMERIC(5,2) NULL`, and `recorded_by UUID`. Unique constraint on `(client_id, date)`; later writes update the existing row.
- **D-09:** `lab_results` stores `title VARCHAR(200)`, `lab_type lab_result_type`, `test_date DATE`, `file_path TEXT NULL`, `external_link TEXT NULL`, `original_filename VARCHAR(255) NULL`, `mime_type VARCHAR(100) NULL`, `file_size_bytes BIGINT NULL`, and `uploaded_by UUID`. Validation enforces at least one of `file_path` or `external_link`.

### Idempotency, Authorization, and Query Pattern

- **D-10:** All client POST endpoints accept `local_id`; repositories implement `CreateOrUpsertByLocalID` or `UpsertByUniqueKey` methods so duplicate offline submissions return the existing record instead of creating a second one.
- **D-11:** Nutritionist read endpoints always enforce row-level ownership in SQL by joining the client to the authenticated nutritionist ID. Never trust `client_id` from the URL by itself.
- **D-12:** All list endpoints support `from`/`to` filtering (or `date` for single-day reads) and order newest-first for history views, except charts which return oldest-first to simplify rendering.
- **D-13:** Use sqlc-generated queries for standard CRUD/list operations; reserve raw pgx only for richer aggregate reads like the daily dashboard summary or measurement chart series.

### Client Dashboard and Logging UX

- **D-14:** `/client/tracking` becomes the client's daily dashboard. It shows today's plan-linked summary first, then quick-action cards for food, water, sleep, exercise, medication, body measurements, and lab results.
- **D-15:** Food logging is meal-centric. The client opens a meal, sees the plan meal title and all available options from the active diet plan, then selects the eaten option or marks the meal as skipped.
- **D-16:** Water tracking uses quick-add buttons for common sizes (200ml, 250ml, 500ml) plus a custom amount field, and shows progress toward the active plan's `daily_water_target_ml`.
- **D-17:** Sleep form uses bedtime and wake-time pickers with automatic duration calculation in the UI; quality choices are `good`, `fair`, `poor` but always displayed with Persian labels.
- **D-18:** Medication checklist is pre-populated from active plan prescriptions. Each prescribed medication time becomes a tappable checklist item that creates a `medication_log` with the linked `prescribed_medication_id`.
- **D-19:** Body measurements emphasize weight first, with additional fields hidden behind an expandable section to keep the mobile form compact.
- **D-20:** Lab results allow either file upload (PDF/JPG/PNG, max 10MB) or external link. File upload is part of Phase 4 because roadmap success criteria require client upload plus nutritionist download.

### Nutritionist Review UX

- **D-21:** Nutritionist tracking access lives under the existing client area, not as a separate global module. Implement a client profile route with tabs for overview, tracking categories, measurements, and lab results.
- **D-22:** Nutritionist can create or update body measurements for a client using the same API as clients; `recorded_by` stores the nutritionist user ID.
- **D-23:** Weight history is visualized with a simple Chart.js line chart using Shamsi date labels. Other measurements remain tabular in Phase 4 unless a multi-line chart is trivial once weight is working.
- **D-24:** Nutritionist lab result view supports list + download. File downloads use `Content-Disposition: attachment` and never inline rendering.

### Storage, Validation, and Error Handling

- **D-25:** Lab result files are stored on the local filesystem under a dedicated `lab-results` directory with UUID filenames and preserved original filename metadata in the database.
- **D-26:** File validation is mandatory: accept only PDF/JPG/PNG, max 10MB, and validate with MIME sniffing / magic bytes before persisting.
- **D-27:** All API errors remain Persian. Validation mirrors existing Gin binding + service sentinel error patterns from Phases 1–3.
- **D-28:** Tailwind logical RTL utilities, Persian numerals, Jalali dates, and mobile card layouts remain mandatory in all new UI.

### Agent's Discretion

- Exact quick-action card visual design on the dashboard
- Whether lab results appear on the main dashboard or in a secondary section below core daily logs
- Whether measurement history beyond weight uses mini-cards or a compact list layout
- Exact folder path configuration mechanism for lab-result file storage, as long as it remains local filesystem based

</decisions>

<specifics>
## Specific Ideas

- Reuse Phase 3 plan data directly in the client dashboard so meal logging and medication checklist are always anchored to the active diet plan rather than free-form entries.
- Keep tracking forms shallow and focused; mobile-first means most interactions should fit within one screen without long scrolling forms.
- `local_id` support is non-negotiable even though offline UX itself is Phase 6.

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Product Scope
- `.planning/ROADMAP.md` §Phase 4 — Goal, success criteria, and dependency position
- `.planning/REQUIREMENTS.md` §Client Tracking, §Lab Results — TRACK-01 through TRACK-13 and LAB-01 through LAB-05
- `docs/phases.md` §Phase 4 — Scope breakdown, implementation guidance, validation checklist

### Existing Architecture and Phase Contracts
- `.planning/research/STACK.md` — Confirmed stack choices including Chart.js and local filesystem storage
- `.planning/phases/03-diet-plan-engine/03-CONTEXT.md` — Active plan model, prescribed medications, water target, and client plan routes that Phase 4 builds on
- `.planning/phases/03-diet-plan-engine/03-RESEARCH.md` — Existing diet plan loading and frontend route patterns
- `.planning/phases/02-core-data-domain/02-CONTEXT.md` — Established CRUD, sqlc, repository/service/handler, and list UI patterns

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `frontend/app/layouts/client.vue` — already exposes `/client/tracking` in the client bottom navigation
- `frontend/app/stores/clientPlan.ts` — store structure to mirror for tracking stores and daily dashboard data loading
- `frontend/app/pages/client/plan.vue` — tabbed mobile client page pattern and active-plan rendering
- `frontend/app/components/plan/MedicationCard.vue` — reusable medication presentation for checklist rows
- `frontend/app/components/plan/ExerciseCard.vue` — reusable plan exercise display card
- `frontend/app/composables/useApi.ts` — standard authenticated fetch wrapper with refresh retry
- `frontend/app/composables/useShamsiDate.ts` — canonical Jalali display helper for history and charts
- `backend/internal/repository/diet_plan_repo.go` / `service` / `handler` — exact backend layering pattern to copy
- `backend/cmd/api/main.go` — route groups already separate client vs nutritionist access

### Established Patterns
- Repository-level authorization checks instead of trusting handler-level filtering
- sqlc query generation for standard CRUD/list operations
- Persian-only validation and error responses
- Mobile card UIs over desktop tables
- Shared DTOs in `backend/internal/model/dto`

### Integration Points
- Food logs attach to `meals` and `meal_options` from Phase 3
- Medication checklist derives from `plan_medications.times` from Phase 3
- Water progress reads `daily_water_target_ml` from the active diet plan
- Nutritionist tracking history should connect to the existing nutritionist client routes rather than creating a new standalone navigation surface

</code_context>

<deferred>
## Deferred Ideas

- Offline queueing, sync retries, and IndexedDB caching — Phase 6
- Push reminders for meals, medications, and water — Phase 6
- Rich multi-line body measurement charts beyond weight if they add substantial scope
- AI-derived adherence scoring or weekly tracking summaries — future backlog / v2

</deferred>

---

*Phase: 04-client-tracking-suite*
*Context gathered: 2026-04-20*
