# NutriTrack PWA — 7-Phase Implementation Plan

**Author:** Joe (Technical Lead)  
**Date:** April 19, 2026  
**Version:** 1.0  
**Tech Stack:** Go + PostgreSQL + Nuxt 4 + Docker/Hetzner

---

## Executive Summary

This document breaks the NutriTrack PRD into seven sequential implementation phases. Each phase is designed to produce a deployable, testable increment of the platform. The phases are ordered by dependency (infrastructure first, then core domain, then progressively richer features) and by risk (hardest architectural decisions early, polish and optimization last).

Each phase includes: scope definition, implementation guidance, specific deliverables, and a validation checklist that serves as the phase gate. A phase is considered complete only when every item on its validation checklist passes. No phase should begin until the previous phase is validated.

### Phase Overview

| Phase | Name | Focus | Est. Duration |
|-------|------|-------|---------------|
| 1 | Foundation & Infrastructure | Project scaffolding, CI/CD, database, auth | 2–3 weeks |
| 2 | Core Data Domain | Food DB, Medication DB, Super Admin panel | 2–3 weeks |
| 3 | Diet Plan Engine | Plan CRUD, nutritional computation, client views | 3–4 weeks |
| 4 | Client Tracking Suite | Food logs, water, sleep, exercise, body, meds | 2–3 weeks |
| 5 | Communication Layer | Messaging, food requests, lab uploads, notifications | 2–3 weeks |
| 6 | Offline & PWA | Service worker, IndexedDB, sync manager, web push | 2–3 weeks |
| 7 | Hardening & Launch | Security audit, perf tuning, monitoring, deployment | 1–2 weeks |

### Dependency Graph

Phase 1 → Phase 2 → Phase 3 → Phase 4 (parallel-safe with Phase 5 after Phase 3) → Phase 5 → Phase 6 → Phase 7. Phases 4 and 5 can run concurrently if the team is split, since they share no data model dependencies beyond what Phase 3 establishes.

---

## Phase 1: Foundation & Infrastructure

**Estimated Duration:** 2–3 weeks

**Goal:** Establish the complete development and deployment infrastructure so that every subsequent phase has a stable, automated pipeline to build on. By the end of this phase, you should have a running Go API server and Nuxt 4 frontend shell, behind Traefik, with PostgreSQL live, CI/CD green, and three user roles authenticating successfully.

### 1.1 Scope

#### Backend (Go)

- Initialize Go module with project structure: `cmd/`, `internal/handler/`, `internal/service/`, `internal/repository/`, `internal/model/`, `internal/middleware/`, `pkg/`
- Choose and configure HTTP framework (Fiber or Echo) with structured JSON logging
- Set up database connection pool with pgx/pgxpool to PostgreSQL
- Implement database migration system using Atlas or golang-migrate
- Create initial migration: `users` table with role enum (`super_admin`, `nutritionist`, `client`), all columns per PRD Section 9
- Implement JWT middleware: access token (15 min) + refresh token (30 days), bcrypt password hashing (cost 12)
- Implement OTP flow for client auth: generate 6-digit code, store with 2-min TTL, max 3 attempts, rate limit 3 requests/phone/10 min
- Build SMS gateway adapter interface with a mock implementation (log to stdout) and a real implementation (Kavenegar)
- Seed super admin account via CLI command or migration
- Health check endpoint: `GET /health` returning DB connectivity status
- CORS configuration restricted to frontend domain
- Input validation middleware using struct tags

#### Frontend (Nuxt 4)

- Initialize Nuxt 4 project with TypeScript, Tailwind CSS (RTL plugin), Pinia
- Configure project as Persian-only RTL: set `dir='rtl'` globally, Persian font stack (Vazirmatn or IRANSans)
- Set up mobile-first viewport meta, max-width container (no desktop breakpoints)
- Build auth pages: Super Admin login (email/password), Nutritionist login (email/password), Client login (mobile + OTP)
- Implement auth store in Pinia: token storage, refresh logic, role-based route guards
- Create layout shells: `AdminLayout`, `NutritionistLayout`, `ClientLayout` with bottom navigation stubs
- Set up `ofetch`/`useFetch` composable with auth header injection and 401 redirect

#### Infrastructure

- Docker Compose: Go API, PostgreSQL 16, Traefik reverse proxy with HTTPS (Let's Encrypt)
- Dockerfile for Go (multi-stage build) and Nuxt (SSR or static generation)
- GitLab CI/CD pipeline: lint, test, build, deploy stages
- Structured logging to stdout in JSON format (ready for Loki collection)
- Environment configuration via `.env` files with validation on startup

### 1.2 Implementation Guidance

Start with the Go project structure and database migration before touching the frontend. The `users` table is the foundation everything else references. Get the auth flow fully working with integration tests before building the frontend auth pages.

For the OTP system, build behind an SMS adapter interface from day one. The mock adapter logs OTPs to stdout during development. Swapping to Kavenegar later is a single environment variable change. Store OTP state in PostgreSQL (not Redis) to keep the stack simple — a simple `otp_codes` table with `phone`, `code`, `expires_at`, `attempts` columns works fine at this scale.

For JWT, use a well-tested Go library (`golang-jwt/jwt`). Store refresh tokens in the database with a `token_family` column to support rotation and revocation. Access tokens are stateless.

On the frontend, Jalali calendar conversion should be configured from the start using `jalaali-js`. All date displays throughout the entire app will use Shamsi dates, so bake this into a `useShamsiDate` composable early.

### 1.3 Validation Checklist

Phase 1 is complete when ALL of the following pass:

| # | Validation Criterion | How to Verify |
|---|---------------------|---------------|
| 1 | Docker Compose starts all services (Go API, PostgreSQL, Traefik) with a single command | Run `docker compose up -d`, all containers healthy within 60s |
| 2 | Database migrations run without errors, users table exists with correct schema | Run `migrate up`, then inspect table via psql |
| 3 | Super Admin can log in via email/password and receives JWT tokens | `POST /api/auth/login` with seeded credentials, verify 200 + tokens |
| 4 | Super Admin can create a nutritionist account | `POST /api/admin/nutritionists`, verify 201 + account in DB |
| 5 | Nutritionist can log in via email/password | `POST /api/auth/login` with new nutritionist credentials |
| 6 | Nutritionist can register a client (name, mobile, DOB, height, gender) | `POST /api/nutritionist/clients`, verify 201 + client in DB |
| 7 | Client receives OTP via mock SMS adapter (logged to stdout) | `POST /api/auth/otp/request`, check server logs for 6-digit code |
| 8 | Client can verify OTP and receive JWT tokens | `POST /api/auth/otp/verify` with correct code, verify tokens |
| 9 | OTP rate limiting works: 4th request within 10 min returns 429 | Send 4 rapid OTP requests, verify 429 on 4th |
| 10 | Frontend renders in Persian RTL, all three login flows work end-to-end in mobile viewport | Manual test in Chrome DevTools mobile emulator |
| 11 | JWT refresh flow works: expired access token triggers refresh, new tokens issued | Wait 15 min or mock expiry, verify seamless refresh |
| 12 | Role-based route guards: client cannot access /admin, admin cannot access /client | Attempt cross-role navigation, verify redirect to correct login |
| 13 | CI/CD pipeline runs: lint passes, tests pass, Docker image builds | Push to main branch, verify pipeline green in GitLab |
| 14 | Health check endpoint returns 200 with DB status | `GET /health`, verify response |
| 15 | HTTPS termination works via Traefik on staging domain | `curl -I https://staging.domain`, verify TLS |

---

## Phase 2: Core Data Domain

**Estimated Duration:** 2–3 weeks

**Goal:** Build the shared food and medication databases that the diet plan engine (Phase 3) depends on, and complete the Super Admin panel. After this phase, nutritionists and admins can populate the platform with food items and medications, and the Super Admin has full operational control.

### 2.1 Scope

#### Database Migrations

- `foods` table with all columns per PRD Section 9 (name, calories, protein, carbs, fat, fiber, sugar, sodium, measurement_unit, measurement_amount, description, is_active, created_by)
- `food_categories` junction table (food_id, category enum)
- `medications` table (name, generic_name, form enum, dosage_unit, description, is_active, created_by)
- Persian full-text search configuration using `pg_trgm` extension for fuzzy name matching
- Indexes on: `foods.name` (trigram), `foods.is_active`, `medications.name`, `medications.is_active`

#### Backend API Endpoints

- Foods CRUD: `POST/GET/PUT/DELETE /api/foods` with pagination (20/page), category filter, active/inactive filter, full-text search
- Food category management: attach/detach categories (many-to-many)
- Medications CRUD: `POST/GET/PUT/DELETE /api/medications` with pagination, search, form filter
- Super Admin endpoints: `GET /api/admin/stats` (total nutritionists, clients, foods, active plans), `GET /api/admin/nutritionists` (list), `PATCH /api/admin/nutritionists/:id/status` (activate/deactivate)
- Authorization middleware: ensure Super Admin and Nutritionist roles can CRUD food/medication; Super Admin can edit items created by others
- Audit trail: `created_by` and `updated_at` tracked on all food/medication records

#### Frontend

- Super Admin dashboard: platform statistics cards, nutritionist list with search/filter, create/edit nutritionist modal
- Food database browser: searchable list with category chips, infinite scroll or pagination, add/edit food form with all nutritional fields and multi-category selector
- Medication database browser: searchable list, add/edit medication form with form-type selector
- Shared components: `SearchInput` (Persian-aware), `PaginatedList`, `CategoryChip`, `NutritionLabel` (displays macro breakdown)

### 2.2 Implementation Guidance

Persian full-text search is a critical infrastructure decision. PostgreSQL does not ship with a Persian text search dictionary. The pragmatic approach is to use `pg_trgm` (trigram matching) with `ILIKE` or `similarity()` for fuzzy search. Create a GIN index on `foods.name` using `gin_trgm_ops`. This handles partial matches and Persian character variations without needing a custom dictionary.

For the food category many-to-many relationship, use a simple junction table rather than a PostgreSQL array column. This keeps queries straightforward (`JOIN + WHERE category = ?`) and supports future indexing if the food database grows large.

Build the `NutritionLabel` component as a reusable card early — it will appear in the food list, diet plan builder, and client tracking views. It should accept calories, protein, carbs, fat as props and render a compact macro summary in Persian.

The measurement unit enum should be defined as a shared constant in both Go (for validation) and Nuxt (for form selects). Consider a shared types file or generate TypeScript types from Go structs.

### 2.3 Validation Checklist

| # | Validation Criterion | How to Verify |
|---|---------------------|---------------|
| 1 | Super Admin can create, edit, activate, and deactivate nutritionist accounts | End-to-end test through admin panel |
| 2 | Nutritionist can add a food item with all required fields and multiple categories | `POST /api/foods` with valid payload, verify in DB and UI |
| 3 | Food search returns results for partial Persian text input within 200ms | Search for partial food name, verify results + timing |
| 4 | Food list pagination works: 20 items per page, next/prev navigation | Add 25+ food items, verify pagination controls |
| 5 | Food category filter correctly filters (e.g., show only 'breakfast' items) | Apply filter, verify only matching items shown |
| 6 | Nutritionist can add, edit, and soft-delete medications | Full CRUD cycle through UI |
| 7 | Super Admin can edit/delete food items created by any nutritionist | Admin edits a nutritionist-created food item |
| 8 | Super Admin dashboard shows correct live statistics | Cross-reference counts with direct DB queries |
| 9 | Authorization enforced: client role cannot access food/medication CRUD endpoints | Send request with client JWT, verify 403 |
| 10 | Soft delete works: deactivated foods do not appear in active list but exist in DB | Deactivate food, verify absence in list, presence in DB |
| 11 | All forms validate required fields and display Persian error messages | Submit empty forms, verify validation messages |
| 12 | Measurement units display correctly in Persian throughout | Visual inspection of food items in list and detail views |

---

## Phase 3: Diet Plan Engine

**Estimated Duration:** 3–4 weeks

**Goal:** This is the most complex phase and the core domain of the application. Build the complete diet plan creation, management, and viewing system. After this phase, a nutritionist can create a multi-day diet plan with meals, options, food items, exercise recommendations, and medication prescriptions — and a client can view their active plan on mobile.

### 3.1 Scope

#### Database Migrations

- `diet_plans` table (client_id, nutritionist_id, start_date, end_date, daily_water_target_ml, notes, status enum)
- `plan_days` table (diet_plan_id, day_number)
- `meals` table (plan_day_id, title, scheduled_time, display_order)
- `meal_options` table (meal_id, option_number)
- `meal_option_items` table (meal_option_id, food_id FK, quantity, measurement_unit, notes)
- `prescribed_medications` table (diet_plan_id, medication_id FK, dosage, frequency, times JSONB, instructions, start_date, end_date)
- `exercise_recommendations` table (plan_day_id, exercise_name, duration_minutes, description, calories_burn_estimate)
- Unique constraint: only one active plan per client (enforced at DB + application level)

#### Backend API

- Diet Plan CRUD: `POST /api/nutritionist/clients/:id/plans` (auto-archives previous active plan), `GET /api/nutritionist/clients/:id/plans` (list with status filter), `GET /api/plans/:id` (full plan tree), `PUT /api/plans/:id`, `DELETE /api/plans/:id`
- Plan Day management: add/remove/reorder days within a plan
- Meal management: CRUD meals within a day, reorder meals
- Meal Option management: CRUD options within a meal, add/remove food items to options
- Nutritional computation service: calculate totals for each option (sum items), each meal (min/max across options), each day (sum meals), and the full plan
- Prescribed Medication management: CRUD within a plan, with medication DB lookup
- Exercise Recommendation management: CRUD within a plan day
- Client plan view endpoint: `GET /api/client/plan` (returns active plan with computed nutrition) — optimized single query with JOINs, not N+1
- Plan archival logic: creating a new plan auto-archives the current active plan with `status='archived'`

#### Frontend — Nutritionist Side

- Diet Plan Builder: multi-step form or tabbed interface for creating a plan
- Day management UI: add days, switch between days via tabs or swipe
- Meal builder: add meals to a day, set title/time/order, add multiple options per meal
- Food item picker: search food database, select item, set quantity and unit — inline add to option
- Real-time nutrition display: as food items are added, show computed calories/protein/carbs/fat per option, per meal, per day (PRD Section 5.3 display format)
- Exercise recommendation form per plan day
- Medication prescription form: pick from medication DB, set dosage/frequency/times
- Plan history view: list of archived plans with date ranges

#### Frontend — Client Side

- Active plan view: daily view showing today's meals, each meal's options with food items and nutrition, scheduled times
- Day navigation: swipe or tap to view other plan days
- Medication schedule view: list of prescribed medications with timing
- Exercise recommendation view for the current day
- Water target display (if set by nutritionist)

### 3.2 Implementation Guidance

The diet plan is a deeply nested entity: Plan → Days → Meals → Options → Items. The biggest risk here is N+1 queries. On the backend, build a dedicated repository method that loads the entire plan tree in 2–3 queries using JOINs or batch loading, not one query per entity. Consider loading `plan_days + meals` in one query, then `meal_options + items` in another, and assembling in-memory.

For nutritional computation, build a pure function in the service layer: given a slice of `meal_option_items` (each with food reference), compute totals. This function is called server-side for the API response and should be fast enough to call on every plan read without caching. At ~50 food items per plan, this is trivial compute.

The plan builder UI is the most complex frontend component. Consider breaking it into composables: `usePlanDays`, `useMealBuilder`, `useFoodPicker`, `useNutritionCalc`. The food picker should be a modal/bottom sheet that searches the food database inline. When a food item is selected, the quantity/unit form appears inline.

For the repeating day pattern (e.g., 7-day cycle), implement this as a UI convenience: the nutritionist builds 7 days, and the frontend maps calendar dates to `day_number` using modulo. The backend stores only the template days, not 30+ individual days.

The "only one active plan" constraint should be enforced both at the database level (a partial unique index on `client_id WHERE status='active'`) and in the application layer (archive before insert). Belt and suspenders.

### 3.3 Validation Checklist

| # | Validation Criterion | How to Verify |
|---|---------------------|---------------|
| 1 | Nutritionist can create a complete diet plan with 7 days, each having 4–5 meals, each meal having 2–3 options with food items | Create a full plan through the UI, verify DB structure |
| 2 | Nutritional totals compute correctly: option total = sum of items, meal shows min/max across options, day = sum of meals | Add known food items (e.g., 100g rice = 130 cal), verify computed totals match manual calculation |
| 3 | Creating a new plan for a client auto-archives the previous active plan | Create two plans for same client, verify first has `status='archived'` |
| 4 | Only one active plan per client constraint holds at DB level | Attempt to manually INSERT two active plans for same client via SQL, verify constraint violation |
| 5 | Full plan tree loads via API in a single request under 500ms | `GET /api/client/plan`, measure response time with curl |
| 6 | Client mobile view displays today's meals with all options and nutritional info | Log in as client on mobile viewport, verify plan renders correctly in RTL |
| 7 | Exercise recommendations display on the correct plan day | Add exercises to specific days, verify client sees them on the right day |
| 8 | Prescribed medications show with dosage, frequency, and timing | Prescribe a medication, verify it appears in client's medication schedule view |
| 9 | Plan history shows archived plans with correct date ranges | Archive 2–3 plans, verify history list shows all with dates |
| 10 | Food picker search works within plan builder: search, select, set quantity | Search for a food while building a plan, add it to an option |
| 11 | Day navigation works: client can view all days of the plan | Swipe/tap through days, verify content changes correctly |
| 12 | Authorization: nutritionist can only manage plans for their own clients | Attempt to create plan for another nutritionist's client, verify 403 |
| 13 | Water target (if set) displays on client's daily view | Set water target in plan, verify client sees it |
| 14 | Plan with repeating 7-day pattern maps correctly to calendar dates | Create 7-day plan spanning 21 days, verify day_number cycling |

---

## Phase 4: Client Tracking Suite

**Estimated Duration:** 2–3 weeks

**Goal:** Build all client-side daily tracking features: food logging, water intake, sleep, exercise, medication intake, body measurements, and weight. After this phase, clients can record their daily activities and nutritionists can view comprehensive tracking history for each client.

### 4.1 Scope

#### Database Migrations

- `food_logs` table (client_id, date, meal_id FK, selected_option_id FK nullable, notes, local_id UUID)
- `water_logs` table (client_id, date, amount_ml, time nullable, local_id)
- `sleep_logs` table (client_id, date unique, sleep_time, wake_time, quality enum, notes, local_id)
- `exercise_logs` table (client_id, date, exercise_name, duration_minutes, calories_burned nullable, notes, local_id)
- `medication_logs` table (client_id, date, prescribed_medication_id FK nullable, medication_name, dosage, taken_at, notes, local_id)
- `body_measurements` table (client_id, date, weight_kg, waist_cm, hip_cm, abdomen_cm, thigh_cm, chest_cm, wrist_cm, recorded_by FK, local_id)
- All tracking tables include `local_id` UUID column for offline sync deduplication (Phase 6)

#### Backend API

- Food log: `POST /api/client/food-logs`, `GET /api/client/food-logs?date=`, `GET /api/nutritionist/clients/:id/food-logs?from=&to=`
- Water log: POST, GET (daily total + entries), target comparison endpoint
- Sleep log: POST/PUT (upsert per date), GET history
- Exercise log: POST, GET by date range
- Medication log: POST (mark dose taken), GET history with prescribed vs. self-reported distinction
- Body measurements: POST (client or nutritionist), GET history, weight trend endpoint
- Deduplication: on POST, if `local_id` already exists, return existing record (idempotent upsert)
- Nutritionist read endpoints: all tracking data for their own clients with date range filtering

#### Frontend — Client Side

- Daily Dashboard: shows today's plan summary, quick-log buttons for water/food/exercise/sleep/medication, progress indicators
- Food logging: tap a meal → select which option was eaten or mark as skipped
- Water tracker: tap-to-add with configurable glass size (200ml/250ml), progress bar toward daily target, timestamped entries
- Sleep log form: sleep time + wake time pickers, optional quality rating, auto-computed duration
- Exercise log form: exercise name, duration, optional calorie burn, free-text notes
- Medication checklist: list of prescribed medications with scheduled times, tap to mark taken (auto-fills timestamp), plus form for self-reported medications
- Body measurements form: weight (prominent), expandable section for waist/hip/abdomen/thigh/chest/wrist
- Weight chart: simple line chart showing weight over time using Chart.js

#### Frontend — Nutritionist Side

- Client profile: history tabs for each tracking category (food logs, water, sleep, exercise, meds, body measurements)
- Weight and measurement charts per client
- Ability for nutritionist to record body measurements for a client (`recorded_by` = nutritionist)

### 4.2 Implementation Guidance

Every tracking table follows the same pattern: `client_id`, `date`, data fields, `local_id`, `created_at`. Build a generic repository pattern for these: `CreateOrUpsertByLocalID`, `ListByClientAndDateRange`, `ListByDate`. The `local_id` upsert is critical infrastructure for Phase 6 offline sync — build it now even though offline support comes later.

The daily dashboard is the client's home screen. Design it as a single scrollable view with the day's plan at top, then quick-action cards for each tracking type below. Each card shows current status (e.g., "3 of 8 glasses", "Breakfast: logged", "Weight: not logged today"). Use Pinia stores per tracking domain: `useWaterStore`, `useFoodLogStore`, etc.

For the medication checklist, pre-populate from `prescribed_medications` in the active diet plan. Each prescribed medication with `times` array generates individual checklist items. The client taps to mark each dose taken, which creates a `medication_log` entry with `prescribed_medication_id` filled.

Weight and measurement charts should use Chart.js (already in the tech stack). Keep charts simple: line chart with Shamsi date labels on x-axis, values on y-axis. The measurement chart can be a multi-line chart (one line per body part) or separate small charts.

### 4.3 Validation Checklist

| # | Validation Criterion | How to Verify |
|---|---------------------|---------------|
| 1 | Client can log food for each meal: select an option or mark as skipped | Log food for all meals in a day, verify records in DB |
| 2 | Client can add water entries, daily total displays correctly vs. target | Add 5 water entries (200ml each), verify total shows 1000ml against target |
| 3 | Client can log sleep with time pickers, duration auto-computes | Log sleep 23:00 → wake 07:00, verify 8 hours displayed |
| 4 | Client can log exercise entries, multiple per day | Add 3 exercises in one day, all appear in list |
| 5 | Medication checklist populates from prescribed medications | Prescribe 3 meds in diet plan, verify 3 checklist items appear for client |
| 6 | Client can mark medication doses as taken with timestamp | Tap to mark taken, verify `medication_log` created with correct time |
| 7 | Client can log body measurements (weight + optional others) | Log weight + waist + hip, verify all stored correctly |
| 8 | Weight chart displays historical data with Shamsi dates | Log weight over 7 days, verify chart renders correctly |
| 9 | Nutritionist can view all tracking data for their client | Log data as client, switch to nutritionist view, verify all data visible |
| 10 | Nutritionist can record body measurements for a client (recorded_by shows nutritionist) | Record measurement as nutritionist, verify `recorded_by` field |
| 11 | `local_id` deduplication works: POSTing same local_id twice returns same record | Send identical POST twice, verify only one record in DB |
| 12 | Date range filtering works on all tracking endpoints | Query with from/to params, verify only matching records returned |
| 13 | Daily dashboard shows correct status summary for all tracking types | Log some items, leave others empty, verify dashboard reflects state |

---

## Phase 5: Communication & Collaboration

**Estimated Duration:** 2–3 weeks

**Goal:** Build the messaging system between clients and nutritionists, the food request workflow, lab result uploads, and the client management dashboard. After this phase, the platform supports full bidirectional communication and the nutritionist has a complete client management experience.

### 5.1 Scope

#### Database Migrations

- `messages` table (sender_id, receiver_id, content, attachment_type, attachment_path, sent_at, read_at)
- `food_requests` table (food_name, description, status enum, rejection_reason, requested_by FK, reviewed_by FK, created_at, updated_at)
- `lab_results` table (client_id, title, type enum, test_date, file_path, link, notes)

#### Messaging System

- Backend: `POST /api/messages` (send), `GET /api/messages/:conversationPartnerId` (list with pagination), `PATCH /api/messages/:id/read` (mark read), `GET /api/messages/unread-count`
- File upload endpoint for attachments: images (JPG/PNG, max 5MB), files (PDF, max 10MB) with content-type validation and content sniffing protection
- Polling endpoint: `GET /api/messages/new?since=` (returns new messages since timestamp, used by 10-second polling)
- Authorization: client can only message their assigned nutritionist; nutritionist can only message their own clients
- Frontend chat UI: conversation view with message bubbles, text input, attachment button, unread badge on navigation tab, auto-scroll to bottom, 10-second polling when chat is open

#### Food Request System

- Client: `POST /api/client/food-requests` (submit), `GET /api/client/food-requests` (my requests with status)
- Nutritionist: `GET /api/nutritionist/food-requests` (pending requests from own clients), `PATCH /api/nutritionist/food-requests/:id` (approve with food item creation, or reject with reason)
- Approval flow: approving a food request navigates nutritionist to the food creation form pre-filled with the request name, nutritionist adds nutritional data and saves
- Frontend: client sees request status history; nutritionist sees pending requests list with approve/reject actions

#### Lab Results

- Client: `POST /api/client/lab-results` (upload with file or link), `GET /api/client/lab-results` (my results)
- Nutritionist: `GET /api/nutritionist/clients/:id/lab-results` (view), file download endpoint
- File storage: save to `/data/uploads/lab-results/{client_id}/{uuid}.{ext}`, path stored in DB
- Validation: at least one of file or link required, max 10MB, accepted formats PDF/JPG/PNG

#### Client Management Dashboard (Nutritionist)

- Client list: name, mobile, status, current plan status, last activity, search/filter/sort
- Client profile: personal info, current plan summary, history tabs (weight chart, food logs, exercise, water, sleep, meds, lab results, archived plans), quick actions (new plan, message, deactivate)

### 5.2 Implementation Guidance

The messaging system uses polling, not WebSocket (per PRD Decision #7). Implement polling with a `useMessagePolling` composable that fires every 10 seconds when the chat view is active and stops when the user navigates away. The polling endpoint should be lightweight: query messages `WHERE sent_at > :since AND (sender_id = :partnerId OR receiver_id = :partnerId)`. Index on `(sender_id, sent_at)` and `(receiver_id, sent_at)`.

For file uploads (messages + lab results), use Go's multipart form handling. Validate: check Content-Type header, sniff the first 512 bytes to verify actual file type (prevent disguised executables), enforce size limits, and store files outside the web root. Serve files through an authenticated download endpoint, not direct filesystem access.

The food request approval flow is a two-step process: nutritionist approves the request (`status='approved'`), then creates the actual food item. In the UI, the "Approve" button should open the food creation form with the `food_name` pre-filled. This avoids creating incomplete food items and keeps the nutritionist in control of nutritional data.

The client management dashboard is the nutritionist's primary workspace. The client list should load quickly (`last_activity` can be a denormalized column on the users table, updated via a trigger or application-level update on any tracking log insert). The client profile is a tab-based view reusing the chart components from Phase 4.

### 5.3 Validation Checklist

| # | Validation Criterion | How to Verify |
|---|---------------------|---------------|
| 1 | Client can send text messages to their nutritionist, messages appear in conversation | Send message, verify it appears on nutritionist side |
| 2 | Nutritionist can send text messages to any of their own clients | Send message from nutritionist, verify client sees it |
| 3 | Image attachments work: upload, display thumbnail in chat, tap to view full size | Send image (JPG, <5MB), verify display in chat |
| 4 | PDF attachment works: upload and download | Send PDF (<10MB), verify download works |
| 5 | File validation blocks oversized and wrong-type files | Attempt 15MB upload, attempt .exe upload, verify 400 errors |
| 6 | Polling delivers new messages within 10 seconds | Send message, measure time until it appears on other side |
| 7 | Unread badge shows correct count, clears when conversation opened | Send 3 messages, verify badge shows 3, open chat, verify badge clears |
| 8 | Client cannot message a nutritionist other than their assigned one | Attempt via API, verify 403 |
| 9 | Client can submit food request, sees status updates | Submit request, verify pending status in client UI |
| 10 | Nutritionist can approve food request and create food item from it | Approve request, verify food created in shared DB |
| 11 | Nutritionist can reject food request with reason, client sees rejection | Reject with reason, verify client sees rejection reason |
| 12 | Client can upload lab results (PDF or link), nutritionist can view/download | Upload PDF lab result, verify nutritionist can access it |
| 13 | Client management list shows correct data and supports search/filter/sort | Add 10+ clients, verify search by name, filter by status, sort works |
| 14 | Client profile shows all tracking history tabs with correct data | View client profile, verify all tabs populated from Phase 4 data |
| 15 | Authorization: nutritionist cannot access another nutritionist's clients | Attempt via API with wrong nutritionist JWT, verify 403 |

---

## Phase 6: Offline Support & PWA

**Estimated Duration:** 2–3 weeks

**Goal:** Transform the application into a fully functional PWA with offline capabilities for clients. After this phase, clients can view their diet plan and log all tracking data while offline, with automatic synchronization when connectivity returns. Push notifications are also implemented.

### 6.1 Scope

#### Service Worker & Caching

- Configure `@vite-pwa/nuxt` with a custom service worker strategy
- Cache-first strategy for static assets (JS, CSS, fonts, images)
- Network-first strategy for API responses (diet plan, messages) with stale fallback
- PWA manifest with Persian app name, icons, theme color, `display: standalone`
- App shell caching: all route components cached for instant offline navigation

#### IndexedDB Storage (Dexie.js)

- Define Dexie schema for offline stores: `activePlan`, `foodLogs`, `waterLogs`, `sleepLogs`, `exerciseLogs`, `medicationLogs`, `bodyMeasurements`, `messages`, `syncQueue`
- On first load (or plan update), cache the entire active diet plan structure to IndexedDB
- Cache last 50 messages per conversation
- Cache plan-related food items for offline reference

#### Sync Manager

- Offline queue: when network is unavailable, all tracking log POST requests are saved to `syncQueue` table in IndexedDB with `local_id`, `entity_type`, `payload`, `created_at`, `status` (pending/syncing/failed)
- On network reconnect (`navigator.onLine` event + online event listener), process syncQueue in FIFO order
- Deduplication: backend uses `local_id` to prevent duplicate entries (already built in Phase 4)
- Conflict resolution: last-write-wins with server timestamp
- Retry logic: exponential backoff (1s, 2s, 4s), max 3 retries, then mark as "failed" for manual retry
- Background Sync API registration where supported, with polling fallback on app open
- Sync status UI: indicator showing "Syncing...", "All synced", or "X items pending" with manual retry button for failed items

#### Push Notifications (Web Push)

- Database: `push_subscriptions` table (user_id, endpoint, p256dh, auth), `notification_preferences` table
- Backend: generate VAPID keys, implement Web Push sending via `webpush-go` library
- Subscribe endpoint: `POST /api/push/subscribe` (save subscription), `DELETE /api/push/subscribe` (unsubscribe)
- Notification triggers: new message, new diet plan, food request result, meal reminders, medication reminders, water reminders
- Reminder scheduling: based on diet plan meal times and medication times, use a background Go goroutine or cron job to send reminders
- Client notification preferences: toggle each reminder type in settings
- Frontend: permission prompt on first login, notification handler in service worker to display and handle click (navigate to relevant page)

### 6.2 Implementation Guidance

The offline strategy is client-only (per PRD Decision #2). Nutritionist and admin roles always require connectivity. Wrap the existing `ofetch`/`useFetch` composable with an offline-aware layer: when a POST fails due to network error, serialize the request to IndexedDB syncQueue instead of showing an error. Show a toast message: "Saved offline, will sync when connected."

For the sync manager, build a `useSyncManager` composable that watches `navigator.onLine` and processes the queue. Process items one at a time (not in parallel) to maintain order and avoid overwhelming the API on reconnect. Each successful sync deletes the item from IndexedDB and emits an event so the UI can update.

Diet plan caching is the most important offline feature. When the client loads their plan (online), serialize the entire plan tree to IndexedDB. On subsequent visits, load from IndexedDB first (instant), then fetch from API in the background and update if changed. Use an ETag or `last_modified` header on the plan endpoint to skip the update if nothing changed.

For push notifications, the reminder scheduler needs to run server-side. A simple approach: every minute, a goroutine queries for plans with meal/medication times within the next minute and sends push notifications. Use a `processed_reminders` table or in-memory set to avoid duplicate sends. Alternatively, compute all reminders for the day at midnight and schedule them with `time.AfterFunc`.

### 6.3 Validation Checklist

| # | Validation Criterion | How to Verify |
|---|---------------------|---------------|
| 1 | PWA installs on Android Chrome and iOS Safari as a standalone app | Add to home screen on both platforms, verify standalone launch |
| 2 | Client can view their diet plan with no network connection | Enable airplane mode after loading plan, navigate plan days/meals |
| 3 | Client can log food intake offline, entry appears in syncQueue | Disable network, log food, check IndexedDB for queued entry |
| 4 | Client can log water, sleep, exercise, medications, and body measurements offline | Disable network, log each type, verify all queued |
| 5 | On reconnect, all queued entries sync to server within 30 seconds | Re-enable network, verify all entries appear in server DB within 30s |
| 6 | Duplicate prevention: same local_id synced twice does not create duplicate records | Manually trigger sync twice for same item, verify single DB record |
| 7 | Failed sync retries with exponential backoff, marks as failed after 3 attempts | Simulate server error (500), verify retry behavior and failed status |
| 8 | Sync status indicator shows correct state (syncing/synced/pending) | Observe indicator during offline → online transition |
| 9 | Cached messages are viewable offline | Load chat, go offline, verify last 50 messages still visible |
| 10 | Queued outgoing messages sync when back online | Send message offline, go online, verify message delivered |
| 11 | Push notification received for new message | Send message to client, verify push notification on client device |
| 12 | Meal time reminder notification fires at scheduled time | Set meal at specific time, verify notification fires |
| 13 | Medication reminder fires at prescribed time | Prescribe medication at specific time, verify notification |
| 14 | Client can toggle notification preferences (enable/disable each type) | Disable meal reminders, verify no meal notifications fire |
| 15 | Static assets load from cache (app works on slow 3G with <3s initial load after first visit) | Throttle to 3G in DevTools, reload app, measure load time |
| 16 | Service worker updates correctly: new version detected, prompt to reload | Deploy new version, verify update prompt on next visit |

---

## Phase 7: Hardening & Launch Preparation

**Estimated Duration:** 1–2 weeks

**Goal:** Perform security hardening, performance optimization, monitoring setup, and final polish. After this phase, the application is production-ready and can be deployed to the live Hetzner environment with confidence.

### 7.1 Scope

#### Security Hardening

- SQL injection audit: verify all queries use parameterized statements (no string concatenation)
- XSS prevention: sanitize all user-generated content rendered in the frontend (food names, messages, notes)
- CORS: restrict to exact production domain only (no wildcards)
- Rate limiting: apply to all public endpoints (login, OTP, API), use sliding window
- File upload hardening: re-verify content sniffing, directory traversal protection, file size enforcement
- JWT security: verify refresh token rotation works, revocation on password change, token family tracking
- Row-level authorization audit: verify nutritionist can ONLY access their own clients' data across ALL endpoints
- HTTPS-only: verify all HTTP redirects to HTTPS, HSTS header set
- Dependency vulnerability scan: `go vet`, `govulncheck`, `npm audit`

#### Performance Optimization

- API response time audit: identify and optimize any endpoint exceeding 200ms
- Database query analysis: `EXPLAIN ANALYZE` on all heavy queries, add missing indexes
- Diet plan load optimization: verify full plan loads under 500ms
- Frontend bundle analysis: check for oversized chunks, lazy-load routes
- Image optimization: compress uploaded images on save (if not already done)
- Connection pooling: verify pgxpool settings are appropriate for expected load (500 concurrent users)

#### Monitoring & Observability

- Grafana dashboards: API response times, error rates, active users, database connection pool utilization
- Loki log aggregation: verify structured JSON logs flow from Docker to Loki
- Health check: verify `/health` endpoint is monitored by an uptime checker
- Alerting: set up alerts for high error rates (>1% 5xx), slow responses (>1s p95), disk space warnings
- PostgreSQL backup: verify daily automated backups via `pg_dump` or WAL archiving, test restore procedure
- File backup: verify weekly backup of `/data/uploads/` directory

#### Final Polish

- End-to-end test: full user journey for each role (admin creates nutritionist → nutritionist creates client → client uses app)
- Persian text review: all UI strings, error messages, and notifications reviewed for natural Persian language
- Mobile UX audit: test on actual Android and iOS devices (not just emulators), check touch targets, scroll behavior, keyboard interactions
- Error handling: verify all API errors return meaningful Persian messages, no stack traces leak to client
- Loading states: verify all async operations show loading indicators
- Empty states: verify all list views show meaningful empty state messages (not blank screens)
- Production environment setup: final Docker Compose for production, environment variables, domain configuration

### 7.2 Implementation Guidance

This phase is about checklists, not features. Resist the urge to add new functionality. The security audit should be systematic: go through every handler in the Go codebase and verify it checks authorization, validates input, and uses parameterized queries. Use a spreadsheet to track each endpoint's security status.

For row-level authorization, the most critical test is: can Nutritionist A access Nutritionist B's client data? Write a specific integration test for every endpoint that accepts a `client_id` parameter, using a JWT from a different nutritionist. This is the most likely security vulnerability in the system.

For performance, the diet plan load endpoint is the most complex query. Run `EXPLAIN ANALYZE` with a realistic plan (7 days, 5 meals each, 3 options each, 4 items each = 420 `meal_option_items`). If it exceeds 500ms, consider materializing the plan as a JSON column or adding a Redis cache with short TTL.

The backup restore test is non-negotiable. Actually restore a `pg_dump` backup to a separate database and verify data integrity. Many teams discover their backups are broken only when they need them.

### 7.3 Validation Checklist

| # | Validation Criterion | How to Verify |
|---|---------------------|---------------|
| 1 | Zero SQL injection vectors: all queries use parameterized statements | Code audit of every repository method |
| 2 | Row-level authorization: integration test for every endpoint verifying cross-nutritionist access denied | Automated test suite with cross-tenant assertions |
| 3 | `govulncheck` and `npm audit` report no critical vulnerabilities | Run both tools, resolve any critical findings |
| 4 | All API endpoints respond under 200ms at p95 (load test with 100 concurrent users) | Run k6 or vegeta load test, verify p95 latencies |
| 5 | Diet plan full load under 500ms with realistic data (420+ items) | Seed realistic plan, measure endpoint response |
| 6 | PWA initial load under 3 seconds on simulated 3G (after first visit) | Chrome DevTools throttling, measure DOMContentLoaded |
| 7 | Grafana dashboard shows live metrics for API response times and error rates | Generate traffic, verify metrics appear in dashboard |
| 8 | Loki receives structured JSON logs from all containers | Verify logs visible in Grafana Loki explorer |
| 9 | PostgreSQL backup restore: full restore to separate DB, verify data integrity | Run backup, restore, query key tables, compare row counts |
| 10 | File upload backup: `/data/uploads/` backup verified restorable | Restore from backup, verify files accessible |
| 11 | Alert fires when error rate exceeds threshold | Generate 500 errors, verify alert notification received |
| 12 | Full end-to-end journey passes for all three roles on real mobile devices | Manual test on Android Chrome and iOS Safari |
| 13 | All Persian text is natural and correctly displayed RTL | Native Persian speaker review of all UI screens |
| 14 | No API stack traces leak to client: all errors return clean JSON with Persian messages | Trigger various error conditions, inspect responses |
| 15 | Production Docker Compose starts cleanly with production environment variables | Deploy to staging with production config, verify all services healthy |

---

## Appendix: Summary & Risk Assessment

### Total Estimated Timeline

- **Minimum:** 14 weeks (solo developer, sequential phases)
- **Maximum:** 21 weeks (with buffer for complexity in Phases 3 and 6)
- **With parallelism** (Phases 4 & 5): can save 2–3 weeks if team splits

### Risk Heatmap

| Phase | Highest Risk Area | Mitigation |
|-------|-------------------|------------|
| 1 | OTP rate limiting / JWT refresh edge cases | Write integration tests for every auth edge case before moving on |
| 2 | Persian full-text search performance | Benchmark pg_trgm early with 1000+ food items; have Meilisearch as fallback |
| 3 | Diet plan builder UI complexity | Build the data model and API first; iterate on UI separately |
| 4 | Tracking data volume at scale | Partition tracking tables by month if client count exceeds expectations |
| 5 | File upload security | Content sniffing + virus scan before storage; serve files through authenticated proxy |
| 6 | Offline sync conflicts and data loss | Comprehensive integration tests simulating offline → online transitions with concurrent edits |
| 7 | Performance under realistic load | Load test early in this phase with k6; don't defer until the last day |

### Phase Gate Process

Before starting the next phase, the current phase must pass a formal review:

- All validation checklist items marked as PASS with evidence (screenshot, test output, or log)
- Code reviewed and merged to main branch
- CI/CD pipeline green on main
- Any deferred items documented as tech debt with planned resolution phase
- Brief retrospective: what went well, what to improve for the next phase
