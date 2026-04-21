# Requirements: NutriTrack Backend API

**Defined:** 2026-04-21
**Core Value:** A nutritionist must be able to create a diet plan and assign it to a client — everything else serves this workflow.

---

## v1 Requirements

### Authentication & Authorization

- [ ] **AUTH-01**: Super admin can log in with email and password, receiving JWT access + refresh tokens
- [ ] **AUTH-02**: Nutritionist can log in with email and password, receiving JWT access + refresh tokens
- [ ] **AUTH-03**: Client can request an OTP sent to their registered Iranian mobile number
- [ ] **AUTH-04**: Client can verify OTP and receive JWT access + refresh tokens (OTP valid 2 min, max 3 attempts)
- [ ] **AUTH-05**: Any authenticated user can refresh their access token using a valid refresh token
- [ ] **AUTH-06**: Any authenticated user can log out, invalidating their refresh token in Redis
- [ ] **AUTH-07**: OTP requests are rate-limited (max 3 per phone per 10 minutes) enforced via Redis
- [ ] **AUTH-08**: All protected routes validate JWT and enforce role-based access (super_admin, nutritionist, client)
- [ ] **AUTH-09**: Nutritionist can only access data belonging to their own clients (row-level authorization)

### User & Client Management

- [ ] **USER-01**: Super admin can create a nutritionist account (name, email, password)
- [ ] **USER-02**: Super admin can activate or deactivate a nutritionist account
- [ ] **USER-03**: Super admin can list all nutritionists with pagination and filters (status, search)
- [ ] **USER-04**: Super admin can view a nutritionist's client list (read-only)
- [ ] **USER-05**: Nutritionist can register a new client (full name, mobile, date of birth, height, gender, optional notes)
- [ ] **USER-06**: Nutritionist can activate or deactivate their own clients
- [ ] **USER-07**: Nutritionist can list their clients with search (name/mobile) and filters (active/inactive)
- [ ] **USER-08**: Nutritionist can view a client's full profile (personal info, current plan, history tabs)
- [ ] **USER-09**: Nutritionist can update client profile fields (height, date of birth editable by nutritionist only)

### Food Database

- [ ] **FOOD-01**: Super admin and nutritionist can create a food item with all required fields (name, categories, calories, protein, carbs, fat, measurement unit/amount)
- [ ] **FOOD-02**: Super admin and nutritionist can update food items
- [ ] **FOOD-03**: Super admin can delete/deactivate any food item; nutritionist can deactivate items they created
- [ ] **FOOD-04**: All authenticated users can list and search food items (Persian full-text search via pg_trgm, filter by category, paginated 20 per page)
- [ ] **FOOD-05**: A food item supports multiple categories (many-to-many junction table)
- [ ] **FOOD-06**: Food items include soft delete (`is_active` flag)

### Medication Database

- [ ] **MED-01**: Super admin and nutritionist can create a medication (name, optional generic name, form, dosage unit, description)
- [ ] **MED-02**: Super admin and nutritionist can update medications
- [ ] **MED-03**: Super admin can delete/deactivate any medication; nutritionist can deactivate their own
- [ ] **MED-04**: All authenticated users can list and search medications with pagination

### Diet Plan Management

- [ ] **PLAN-01**: Nutritionist can create a diet plan for a client (period, notes, optional daily water target)
- [ ] **PLAN-02**: Creating a new diet plan automatically archives the client's current active plan
- [ ] **PLAN-03**: Nutritionist can add plan days (numbered), each containing meals with ordered display
- [ ] **PLAN-04**: Each meal can have multiple options (alternatives); client picks one per meal
- [ ] **PLAN-05**: Each meal option contains food items with quantity, unit, and optional notes
- [ ] **PLAN-06**: API returns computed nutritional totals per option, per meal (min/max range), and per day
- [ ] **PLAN-07**: Nutritionist can add exercise recommendations per plan day (name, duration, optional calorie estimate)
- [ ] **PLAN-08**: Nutritionist can prescribe medications on a diet plan with dosage, frequency, times array, and optional date range
- [ ] **PLAN-09**: Client can retrieve their active diet plan (full structure with all days, meals, options, and food details)
- [ ] **PLAN-10**: Client and nutritionist can view archived (historical) diet plans
- [ ] **PLAN-11**: Nutritionist can update or delete an in-progress diet plan

### Daily Tracking — Client

- [ ] **TRACK-01**: Client can log a food entry for a specific date/meal (select option eaten, or mark skipped); idempotent via `local_id`
- [ ] **TRACK-02**: Client can log water intake entries (amount_ml, optional time); idempotent via `local_id`
- [ ] **TRACK-03**: Client can log a sleep record per date (sleep/wake timestamps, optional quality and notes); idempotent via `local_id`
- [ ] **TRACK-04**: Client can log exercise entries per date (name, duration, optional calories, notes); idempotent via `local_id`
- [ ] **TRACK-05**: Client can log medication intake (prescribed or self-reported, with timestamp); idempotent via `local_id`
- [ ] **TRACK-06**: Client and nutritionist can record body measurements per date (weight, waist, hip, abdomen, thigh, chest, wrist); idempotent via `local_id`
- [ ] **TRACK-07**: Nutritionist can view all tracking history for their clients (food logs, water, sleep, exercise, medications, measurements)
- [ ] **TRACK-08**: Bulk sync endpoint accepts an array of offline-queued tracking entries and processes them idempotently

### Lab Results

- [ ] **LAB-01**: Client can upload a lab result (title, type, test date, optional PDF/image file up to 10 MB or external link)
- [ ] **LAB-02**: Files are stored on the local filesystem; file path stored in database
- [ ] **LAB-03**: Nutritionist can list and download all lab results for their clients
- [ ] **LAB-04**: File type validation enforced (PDF, JPG, PNG only)

### Messaging System

- [ ] **MSG-01**: Client can send a text message to their assigned nutritionist
- [ ] **MSG-02**: Nutritionist can send a text message to any of their clients
- [ ] **MSG-03**: Messages support image attachments (JPG/PNG, max 5 MB) and file attachments (PDF, max 10 MB)
- [ ] **MSG-04**: Client can retrieve their conversation with their nutritionist (paginated, chronological)
- [ ] **MSG-05**: Nutritionist can retrieve conversation with any of their clients
- [ ] **MSG-06**: Message read status is tracked (`read_at` timestamp updated when recipient fetches)
- [ ] **MSG-07**: Unread message count returned per conversation

### Food Addition Requests

- [ ] **REQ-01**: Client can submit a food addition request (food name, optional description)
- [ ] **REQ-02**: Nutritionist can list pending food requests from their clients
- [ ] **REQ-03**: Nutritionist can approve a food request (creates food item in shared database)
- [ ] **REQ-04**: Nutritionist can reject a food request with optional reason
- [ ] **REQ-05**: Client receives status update on their food request (via notification)

### Push Notifications

- [ ] **NOTIF-01**: Client and nutritionist can register a Web Push subscription (VAPID endpoint, p256dh, auth)
- [ ] **NOTIF-02**: Backend sends push notification on new message received
- [ ] **NOTIF-03**: Backend sends push notification when new diet plan is assigned to client
- [ ] **NOTIF-04**: Backend sends push notification when food request is approved/rejected
- [ ] **NOTIF-05**: Scheduled push notifications for meal time reminders (based on diet plan scheduled times)
- [ ] **NOTIF-06**: Scheduled push notifications for medication reminders (based on prescribed times)
- [ ] **NOTIF-07**: User can update notification preferences (enable/disable each reminder type)

### Platform Statistics (Super Admin)

- [ ] **STAT-01**: Super admin can retrieve platform-wide statistics (total nutritionists, clients, food items, active diet plans)

### Infrastructure & Cross-Cutting

- [ ] **INFRA-01**: Docker Compose configuration for all services (app, postgres, redis) with `TZ=Asia/Tehran` environment variable
- [ ] **INFRA-02**: All API error responses return a `message` field in Persian
- [ ] **INFRA-03**: Database migrations via golang-migrate with versioned SQL files
- [ ] **INFRA-04**: Health check endpoint (`GET /health`) for monitoring
- [ ] **INFRA-05**: Structured JSON logging (request logs, error logs) to stdout
- [ ] **INFRA-06**: CORS restricted to configured frontend domain
- [ ] **INFRA-07**: All timestamps stored in UTC; API accepts/returns UTC; services run with Asia/Tehran TZ

---

## v2 Requirements

### Enhanced Tracking

- **TRACK-V2-01**: Weight and measurements chart data endpoint (time-series for frontend charting)
- **TRACK-V2-02**: Nutrition summary for a given date (total calories/macros consumed vs plan)

### Advanced Notifications

- **NOTIF-V2-01**: Water intake reminders (configurable interval)
- **NOTIF-V2-02**: Weekly progress summary notification

### Admin Features

- **ADMIN-V2-01**: Audit log of food/medication database modifications
- **ADMIN-V2-02**: Platform-level message broadcasting to all clients

---

## Out of Scope

| Feature | Reason |
|---------|--------|
| Frontend / PWA / Nuxt.js | Backend only — user's explicit instruction |
| Real-time WebSocket chat | PRD decision #7 — polling is simpler, offline-compatible |
| Payment processing | PRD non-goal — no financial transactions |
| AI diet recommendations | PRD non-goal |
| OAuth / social login | YAGNI — email+password and OTP sufficient |
| Multi-language / i18n | Persian only — PRD decision #11 |
| External health device APIs | PRD non-goal |
| Calorie detection from photos | PRD non-goal |
| Desktop UI concerns | Mobile-first frontend concern, not backend |
| Nutritionist self-registration | PRD decision #5 — super admin controlled |

---

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| INFRA-01 | Phase 1 — Foundation | Pending |
| INFRA-02 | Phase 1 — Foundation | Pending |
| INFRA-03 | Phase 1 — Foundation | Pending |
| INFRA-04 | Phase 1 — Foundation | Pending |
| INFRA-05 | Phase 1 — Foundation | Pending |
| INFRA-06 | Phase 1 — Foundation | Pending |
| INFRA-07 | Phase 1 — Foundation | Pending |
| AUTH-01 | Phase 2 — Authentication & Authorization | Pending |
| AUTH-02 | Phase 2 — Authentication & Authorization | Pending |
| AUTH-03 | Phase 2 — Authentication & Authorization | Pending |
| AUTH-04 | Phase 2 — Authentication & Authorization | Pending |
| AUTH-05 | Phase 2 — Authentication & Authorization | Pending |
| AUTH-06 | Phase 2 — Authentication & Authorization | Pending |
| AUTH-07 | Phase 2 — Authentication & Authorization | Pending |
| AUTH-08 | Phase 2 — Authentication & Authorization | Pending |
| AUTH-09 | Phase 2 — Authentication & Authorization | Pending |
| USER-01 | Phase 3 — User & Client Management | Pending |
| USER-02 | Phase 3 — User & Client Management | Pending |
| USER-03 | Phase 3 — User & Client Management | Pending |
| USER-04 | Phase 3 — User & Client Management | Pending |
| USER-05 | Phase 3 — User & Client Management | Pending |
| USER-06 | Phase 3 — User & Client Management | Pending |
| USER-07 | Phase 3 — User & Client Management | Pending |
| USER-08 | Phase 3 — User & Client Management | Pending |
| USER-09 | Phase 3 — User & Client Management | Pending |
| FOOD-01 | Phase 4 — Food & Medication Databases | Pending |
| FOOD-02 | Phase 4 — Food & Medication Databases | Pending |
| FOOD-03 | Phase 4 — Food & Medication Databases | Pending |
| FOOD-04 | Phase 4 — Food & Medication Databases | Pending |
| FOOD-05 | Phase 4 — Food & Medication Databases | Pending |
| FOOD-06 | Phase 4 — Food & Medication Databases | Pending |
| MED-01 | Phase 4 — Food & Medication Databases | Pending |
| MED-02 | Phase 4 — Food & Medication Databases | Pending |
| MED-03 | Phase 4 — Food & Medication Databases | Pending |
| MED-04 | Phase 4 — Food & Medication Databases | Pending |
| PLAN-01 | Phase 5 — Diet Plan Management | Pending |
| PLAN-02 | Phase 5 — Diet Plan Management | Pending |
| PLAN-03 | Phase 5 — Diet Plan Management | Pending |
| PLAN-04 | Phase 5 — Diet Plan Management | Pending |
| PLAN-05 | Phase 5 — Diet Plan Management | Pending |
| PLAN-06 | Phase 5 — Diet Plan Management | Pending |
| PLAN-07 | Phase 5 — Diet Plan Management | Pending |
| PLAN-08 | Phase 5 — Diet Plan Management | Pending |
| PLAN-09 | Phase 5 — Diet Plan Management | Pending |
| PLAN-10 | Phase 5 — Diet Plan Management | Pending |
| PLAN-11 | Phase 5 — Diet Plan Management | Pending |
| TRACK-01 | Phase 6 — Daily Tracking & Lab Results | Pending |
| TRACK-02 | Phase 6 — Daily Tracking & Lab Results | Pending |
| TRACK-03 | Phase 6 — Daily Tracking & Lab Results | Pending |
| TRACK-04 | Phase 6 — Daily Tracking & Lab Results | Pending |
| TRACK-05 | Phase 6 — Daily Tracking & Lab Results | Pending |
| TRACK-06 | Phase 6 — Daily Tracking & Lab Results | Pending |
| TRACK-07 | Phase 6 — Daily Tracking & Lab Results | Pending |
| TRACK-08 | Phase 6 — Daily Tracking & Lab Results | Pending |
| LAB-01 | Phase 6 — Daily Tracking & Lab Results | Pending |
| LAB-02 | Phase 6 — Daily Tracking & Lab Results | Pending |
| LAB-03 | Phase 6 — Daily Tracking & Lab Results | Pending |
| LAB-04 | Phase 6 — Daily Tracking & Lab Results | Pending |
| MSG-01 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| MSG-02 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| MSG-03 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| MSG-04 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| MSG-05 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| MSG-06 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| MSG-07 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| REQ-01 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| REQ-02 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| REQ-03 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| REQ-04 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| REQ-05 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| NOTIF-01 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| NOTIF-02 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| NOTIF-03 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| NOTIF-04 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| NOTIF-05 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| NOTIF-06 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| NOTIF-07 | Phase 7 — Messaging, Food Requests & Notifications | Pending |
| STAT-01 | Phase 8 — Admin, Stats & Hardening | Pending |

**Coverage:**
- v1 requirements: 78 total
- Mapped to phases: 78
- Unmapped: 0 ✓

---
*Requirements defined: 2026-04-21*
*Last updated: 2026-04-21 after initial definition*
