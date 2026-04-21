# Roadmap: NutriTrack — Go DDD Backend API

## Overview

NutriTrack is built in eight sequential phases, each delivering a coherent, testable capability. The
order is strictly dependency-driven: the DDD scaffold and Iran-specific infrastructure baseline must
be correct before any feature work begins, authentication gates every protected endpoint, shared food
and medication databases are prerequisites for diet plans, and diet plans must be stable before
tracking and messaging workflows are layered on top. Hardening closes the milestone after all
functionality is proven in production conditions.

---

## Phases

**Phase Numbering:**
- Integer phases (1–8): Planned milestone work executed in order
- Decimal phases (e.g. 2.1): Urgent insertions created via `/gsd-insert-phase`

- [ ] **Phase 1: Foundation** — DDD scaffold, Docker infra, pg_trgm, Tehran timezone baseline
- [ ] **Phase 2: Authentication & Authorization** — JWT + OTP + Redis token store + RBAC middleware
- [ ] **Phase 3: User & Client Management** — User profiles, client–nutritionist relationships
- [ ] **Phase 4: Food & Medication Databases** — Shared nutrition/medication reference data + Persian search
- [ ] **Phase 5: Diet Plan Management** — 5-level plan hierarchy, computed totals, auto-archive
- [ ] **Phase 6: Daily Tracking & Lab Results** — 6 tracking types, offline sync idempotency, lab uploads
- [ ] **Phase 7: Messaging, Food Requests & Notifications** — Chat, food request workflow, VAPID push
- [ ] **Phase 8: Admin, Stats & Hardening** — Super admin panel, analytics, security audit, production polish

---

## Phase Details

### Phase 1: Foundation
**Goal**: Establish the load-bearing DDD scaffold, Docker infrastructure, and all Iran-specific
baseline concerns so that every subsequent phase builds on a correct, contamination-free foundation.
**Depends on**: Nothing (first phase)
**Requirements**: INFRA-01, INFRA-02, INFRA-03, INFRA-04, INFRA-05, INFRA-06, INFRA-07
**Complexity**: High
**Success Criteria** (what must be TRUE):
  1. `docker compose up` starts the app, PostgreSQL 16, and Redis 7 with health-check-gated ordering; the app responds `200 OK` on `GET /health`
  2. Migration 001 runs cleanly on a fresh database and `pg_trgm` extension is confirmed active (`SELECT * FROM pg_extension WHERE extname = 'pg_trgm'` returns a row)
  3. The app binary runs correctly inside Alpine (`time.LoadLocation("Asia/Tehran")` returns no error; `TodayTehran()` returns today's date in Iran Standard/Daylight Time)
  4. All API error responses (including 404, 422, 500) return JSON with a `message` field containing a Persian string drawn from the `AppError` catalog
  5. Structured JSON logs (request, error, startup) are emitted to stdout; each request log includes a unique `request_id` field
**Plans**: TBD

Plans:
- [ ] 01-01: Go module init, DDD folder skeleton, Viper config, zerolog setup
- [ ] 01-02: Dockerfile (multi-stage, Alpine + tzdata), docker-compose.yml (PG 16, Redis 7, TZ=Asia/Tehran), health-check endpoint
- [ ] 01-03: golang-migrate setup, Migration 001 (pg_trgm + uuid-ossp), sqlc v2 config
- [ ] 01-04: Domain shared utilities (timeutil.go, apperror.go), Gin router groups, request-ID middleware, centralized error middleware

### Phase 2: Authentication & Authorization
**Goal**: Deliver fully working JWT + OTP authentication for all three roles with Redis-backed token
invalidation, RBAC middleware, and rate-limited OTP to gate every protected endpoint in subsequent
phases.
**Depends on**: Phase 1
**Requirements**: AUTH-01, AUTH-02, AUTH-03, AUTH-04, AUTH-05, AUTH-06, AUTH-07, AUTH-08, AUTH-09
**Complexity**: High
**Success Criteria** (what must be TRUE):
  1. Super admin and nutritionist can POST `/auth/login` with email + password and receive a JWT access token (15 min TTL) and a refresh token; invalid credentials return `401` with a Persian error message
  2. Client can POST `/auth/otp/send` with an Iranian mobile number and receive an OTP via SMS (Kavenegar adapter); a fourth request within 10 minutes returns `429` with a Persian rate-limit message
  3. Client can POST `/auth/otp/verify` with the correct OTP and receive JWT tokens; a wrong OTP after 3 attempts invalidates the OTP and returns `401`
  4. Any authenticated user can POST `/auth/refresh` with a valid refresh token to receive a new access token; a used or revoked refresh token returns `401`
  5. Any authenticated user can POST `/auth/logout`; the refresh token is removed from Redis and subsequent refresh attempts with it return `401`
  6. A request to any protected route without a valid JWT returns `401`; a request with a valid JWT but wrong role returns `403`; nutritionist endpoints reject requests targeting clients not belonging to that nutritionist
**Plans**: TBD

Plans:
- [ ] 02-01: User domain aggregate, Iranian mobile validator, bcrypt password hashing, JWT service (access + refresh generation/validation)
- [ ] 02-02: Redis OTP store (atomic INCR rate limiting, 2-min TTL, 3-attempt cap), Kavenegar SMS adapter
- [ ] 02-03: Auth application service (login, OTP send/verify, refresh, logout), Redis token blacklist
- [ ] 02-04: Auth HTTP handlers + public router group, RBAC middleware (RequireRole), row-level auth guard

### Phase 3: User & Client Management
**Goal**: Enable super admin to manage nutritionist accounts and nutritionists to register, view, and
update their client profiles, with BMI computation and optional profile picture upload.
**Depends on**: Phase 2
**Requirements**: USER-01, USER-02, USER-03, USER-04, USER-05, USER-06, USER-07, USER-08, USER-09
**Complexity**: Medium
**Success Criteria** (what must be TRUE):
  1. Super admin can `POST /admin/nutritionists` to create a nutritionist account and receive the created record; the new nutritionist can immediately log in
  2. Super admin can `PATCH /admin/nutritionists/:id/status` to activate or deactivate a nutritionist; a deactivated nutritionist's login attempt returns `401`
  3. Super admin can `GET /admin/nutritionists` with pagination, status filter, and name search returning correct subsets
  4. Nutritionist can `POST /clients` to register a client (full name, mobile, DOB, height, gender); the client appears in `GET /clients` for that nutritionist but not for another nutritionist
  5. Nutritionist can `GET /clients/:id/profile` and receive the client's full profile including computed BMI; `PATCH /clients/:id` updates editable fields and returns the updated profile
  6. A nutritionist cannot access, update, or list clients belonging to another nutritionist (returns `403`)
**Plans**: TBD

Plans:
- [ ] 03-01: Nutritionist management — super admin CRUD + status toggle + list/search (sqlc queries, application service, HTTP handlers)
- [ ] 03-02: Client management — registration, list/search/filter, full profile view, BMI calculation, profile update
- [ ] 03-03: Profile picture upload (magic byte MIME validation, local filesystem storage, path stored in DB)

### Phase 4: Food & Medication Databases
**Goal**: Deliver shared, searchable food and medication reference databases with Persian full-text
search, Arabic/Persian character normalisation, and the food addition request workflow so that diet
plans and prescriptions have concrete items to reference.
**Depends on**: Phase 3
**Requirements**: FOOD-01, FOOD-02, FOOD-03, FOOD-04, FOOD-05, FOOD-06, MED-01, MED-02, MED-03, MED-04
**Complexity**: Medium
**Success Criteria** (what must be TRUE):
  1. Nutritionist can `POST /foods` with name (in Persian), categories, calories, macros, and unit; the item appears immediately in search results
  2. `GET /foods?q=مرغ` returns relevant items using pg_trgm similarity; searching with Arabic Kaf (ك) returns the same results as Persian Kaf (ک) due to normalisation at insert and query time
  3. `GET /foods` paginates at 20 items per page and supports `category` filter; soft-deleted items (`is_active=false`) do not appear in results
  4. Super admin can delete any food item; nutritionist can only deactivate items they created; attempting to deactivate another nutritionist's item returns `403`
  5. Nutritionist can `POST /medications`, `PATCH /medications/:id`, and `GET /medications?q=` with the same CRUD + search behaviour as foods
  6. A food item can belong to multiple categories; assigning two categories at creation returns both in the item's response
**Plans**: TBD

Plans:
- [ ] 04-01: Food domain aggregate, Arabic/Persian char normalisation utility, pg_trgm sqlc queries (similarity search + ILIKE fallback), food CRUD HTTP handlers
- [ ] 04-02: Food categories (many-to-many), soft delete, pagination, category filter
- [ ] 04-03: Medication domain aggregate, CRUD + search, soft delete, row-level write isolation

### Phase 5: Diet Plan Management
**Goal**: Deliver the core value of the platform — nutritionists can create 5-level diet plans for
clients, plans auto-archive atomically, nutritional totals bubble up through the hierarchy, and
clients can retrieve their active plan with all structure and food details.
**Depends on**: Phase 4
**Requirements**: PLAN-01, PLAN-02, PLAN-03, PLAN-04, PLAN-05, PLAN-06, PLAN-07, PLAN-08, PLAN-09, PLAN-10, PLAN-11
**Complexity**: Very High
**Success Criteria** (what must be TRUE):
  1. Nutritionist can `POST /clients/:id/plans` to create a diet plan; if the client already has an active plan, it is atomically archived (status = `archived`) within the same transaction before the new plan is activated — no window where two plans are active
  2. Nutritionist can add days, meals (with ordered display), meal options (alternatives), and food items with quantity and unit; `GET /plans/:id` returns the full nested structure with all levels populated
  3. `GET /plans/:id` returns computed nutritional totals: calories + macros per meal option, min/max range per meal across options, and aggregated totals per day
  4. Client can `GET /plans/active` and receive their current plan with all days, meals, options, and food item details; a client with no active plan receives a clear `404` with a Persian message
  5. Nutritionist can `GET /clients/:id/plans` to list all plans (active + archived) with summary data; client can retrieve their own archived plans
  6. Nutritionist can add exercise recommendations and prescribe medications (with dosage, frequency, and times array) on a plan day; these appear in the plan's day response
**Plans**: TBD

Plans:
- [ ] 05-01: DietPlan domain aggregate (Plan + Days + Meals + Options), CreateWithArchive domain service (atomic transaction), sqlc queries for plan tree assembly
- [ ] 05-02: MealFood aggregate (option + items), computed nutritional totals (bubble-up calculation), plan template cloning
- [ ] 05-03: Exercise recommendations + prescribed medications per day
- [ ] 05-04: Client plan retrieval (active plan, archived plans list), nutritionist plan management (update, delete in-progress plan)

### Phase 6: Daily Tracking & Lab Results
**Goal**: Clients can log all six daily tracking types with offline-sync idempotency guaranteed at the
database level, nutritionists can view full tracking history, bulk sync processes replay-safe batches,
and clients can upload lab results with validated file types.
**Depends on**: Phase 5
**Requirements**: TRACK-01, TRACK-02, TRACK-03, TRACK-04, TRACK-05, TRACK-06, TRACK-07, TRACK-08, LAB-01, LAB-02, LAB-03, LAB-04
**Complexity**: High
**Success Criteria** (what must be TRUE):
  1. Client can submit a food log, water log, sleep record, exercise entry, medication intake, and body measurement — all with a `local_id`; submitting the exact same `local_id` a second time returns `200 OK` with the original record (not `409`), with no duplicate row in the database
  2. `POST /sync/batch` accepts an array of mixed tracking entries; already-synced entries (matching `local_id`) are silently skipped; the response reports total received, inserted, and skipped counts
  3. All "today's" date comparisons use `Asia/Tehran` timezone, not UTC; a log submitted at 23:00 UTC (02:30 Tehran next day) is recorded under tomorrow's Tehran date
  4. Nutritionist can `GET /clients/:id/tracking?type=food&date=2025-06-01` and receive the client's tracking entries for that type and date; querying a client not assigned to the nutritionist returns `403`
  5. Client can `POST /lab-results` with a PDF or image file (max 10 MB); the file is stored on the local filesystem with a UUID-based path; files with invalid magic bytes (not PDF/JPEG/PNG) are rejected with `422` and a Persian error message
  6. Nutritionist can `GET /clients/:id/lab-results` to list all lab results and `GET /lab-results/:id/download` to retrieve the file with `Content-Disposition: attachment`
**Plans**: TBD

Plans:
- [ ] 06-01: Tracking domain aggregates (FoodLog, WaterLog, SleepLog, ExerciseLog, MedicationLog, BodyMeasurement), UNIQUE(client_id, local_id) migrations, ON CONFLICT DO NOTHING sqlc queries
- [ ] 06-02: Tracking HTTP handlers (individual + bulk sync endpoint), Tehran date resolution, nutritionist read access
- [ ] 06-03: Lab result domain aggregate, file upload handler (magic byte validation, UUID paths, 10 MB limit), filesystem storage adapter, download endpoint

### Phase 7: Messaging, Food Requests & Notifications
**Goal**: Clients and nutritionists can exchange messages with file attachments via polling, clients
can submit food addition requests that nutritionists approve into the shared database, and VAPID push
notifications fire on new messages, plan assignments, request outcomes, and scheduled meal/medication
reminders.
**Depends on**: Phase 6
**Requirements**: MSG-01, MSG-02, MSG-03, MSG-04, MSG-05, MSG-06, MSG-07, REQ-01, REQ-02, REQ-03, REQ-04, REQ-05, NOTIF-01, NOTIF-02, NOTIF-03, NOTIF-04, NOTIF-05, NOTIF-06, NOTIF-07
**Complexity**: Very High
**UI hint**: no
**Success Criteria** (what must be TRUE):
  1. Client can `POST /messages` to send a text message to their nutritionist; nutritionist can `POST /clients/:id/messages` to reply; `GET /messages` (client) and `GET /clients/:id/messages` (nutritionist) return the conversation paginated and in chronological order
  2. Messages support image (JPG/PNG ≤ 5 MB) and PDF (≤ 10 MB) attachments validated by magic bytes; invalid files are rejected with `422`; `read_at` is set when the recipient fetches the conversation; unread count is returned per conversation
  3. Client can `POST /food-requests` with a food name; nutritionist can `GET /food-requests` to see pending requests from their clients; approving a request via `POST /food-requests/:id/approve` creates the food item in the shared database; rejecting it via `POST /food-requests/:id/reject` stores the reason
  4. A Web Push subscription registered via `POST /push/subscribe` (VAPID endpoint, p256dh, auth) receives a push notification within 5 seconds when a new message arrives for the subscribed user
  5. A scheduled meal reminder push notification fires at the meal's scheduled time in `Asia/Tehran` timezone (not UTC); DST transitions do not cause reminders to fire an hour early or late
  6. User can `PATCH /notifications/preferences` to disable meal reminders, medication reminders, or message notifications; disabled notification types are not dispatched
**Plans**: TBD

Plans:
- [ ] 07-01: Message domain aggregate, conversation/pagination sqlc queries, read_at tracking, unread count, file attachment upload (magic byte validation)
- [ ] 07-02: Food request domain aggregate (pending → approved/rejected lifecycle), approve handler creates food item in transaction
- [ ] 07-03: VAPID push infrastructure (webpush-go, VAPID key storage, device token registration), event-driven push on message + plan assignment + food request status change
- [ ] 07-04: Notification scheduler (goroutine-based, Asia/Tehran-aware), meal reminders + medication reminders, notification preferences CRUD

### Phase 8: Admin, Stats & Hardening
**Goal**: Super admin has full platform visibility and management APIs, production security concerns
are fully addressed (non-root Docker, rate limiting coverage, input validation, sqlc CI check), and
the system is confirmed stable at target scale (~500 concurrent users).
**Depends on**: Phase 7
**Requirements**: STAT-01
**Complexity**: Medium
**Success Criteria** (what must be TRUE):
  1. Super admin can `GET /admin/stats` and receive platform-wide counts: total nutritionists (active/inactive), total clients, total food items, and total active diet plans — values match the actual database counts
  2. Super admin can view, activate/deactivate, and delete any nutritionist or food/medication item through admin-namespaced endpoints protected by `RequireRole(super_admin)`; any request from a nutritionist or client to an admin endpoint returns `403`
  3. The Docker container runs as a non-root user; `docker inspect` confirms `User` is not `root` or `0`
  4. `make sqlc-check` in CI exits non-zero if the generated sqlc code is out of sync with SQL query files (prevents stale generated code reaching production)
  5. A load test at 500 concurrent users against the five highest-traffic endpoints (food search, active plan retrieval, tracking log submission, message fetch, health check) shows p99 latency ≤ 500ms with no 5xx errors
**Plans**: TBD

Plans:
- [ ] 08-01: Super admin stats endpoint, admin nutritionist management APIs, admin food/medication management APIs
- [ ] 08-02: Non-root Docker user, sqlc CI freshness check, `.env.example` documentation, Makefile hardening targets
- [ ] 08-03: Security audit (RBAC coverage, row-level auth boundaries, JWT blacklist, file upload validation, rate limiting gaps), N+1 query audit, Redis cache for hot paths
- [ ] 08-04: Load test at target scale, fix any p99 regressions, final production readiness review

---

## Progress

**Execution Order:** 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation | 0/4 | Not started | - |
| 2. Authentication & Authorization | 0/4 | Not started | - |
| 3. User & Client Management | 0/3 | Not started | - |
| 4. Food & Medication Databases | 0/3 | Not started | - |
| 5. Diet Plan Management | 0/4 | Not started | - |
| 6. Daily Tracking & Lab Results | 0/3 | Not started | - |
| 7. Messaging, Food Requests & Notifications | 0/4 | Not started | - |
| 8. Admin, Stats & Hardening | 0/4 | Not started | - |
