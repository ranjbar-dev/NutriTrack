# Feature Landscape

**Domain:** Persian Nutrition Management Backend API (Nutritionist ↔ Client)
**Project:** NutriTrack — Go Backend API
**Researched:** 2026-04-21
**Source:** PRD v1.0 + PROJECT.md

---

## Table Stakes

Features where absence means the platform cannot fulfill its core purpose.

| Feature | Why Expected | Complexity | Offline Sync | RBAC Scope |
|---------|--------------|------------|--------------|------------|
| JWT + OTP Authentication | No login = no platform | Medium | ❌ No | All roles |
| Diet Plan CRUD (nested) | Core value proposition | **High** | ❌ No | Nutritionist creates, Client reads |
| Client Registration (by nutritionist) | Clients cannot self-register | Low | ❌ No | Nutritionist only |
| Food Database CRUD + Persian search | Diet plans reference food items | Medium | Partial (cached) | Admin + Nutritionist write; Client reads |
| Medication Database CRUD | Prescriptions reference medications | Low | Partial (cached) | Admin + Nutritionist write; Client reads |
| Daily Food Log | Core client-side tracking | Low | ✅ `local_id` | Client writes; Nutritionist reads |
| Body Measurement Tracking | Progress monitoring | Low | ✅ `local_id` | Client + Nutritionist write; Nutritionist reads |
| Water Intake Tracking | Common nutritionist prescription | Low | ✅ `local_id` | Client writes; Nutritionist reads |
| Sleep Tracking | Holistic health data | Low | ✅ `local_id` | Client writes; Nutritionist reads |
| Exercise Tracking | Plan compliance verification | Low | ✅ `local_id` | Client writes; Nutritionist reads |
| Medication Intake Logging | Prescription adherence | Low | ✅ `local_id` | Client writes; Nutritionist reads |
| Lab Results Upload | Nutritionist diagnostic input | Medium | ❌ No | Client uploads; Nutritionist views/downloads |
| Role-Based Access Control | 3 roles with strict isolation | Medium | N/A | Enforced everywhere |
| Persian Error Messages | All API `message` fields in Farsi | Low | N/A | All endpoints |
| Asia/Tehran Timezone Handling | Iranian users; DST-aware | Medium | N/A | All timestamps |

---

## Differentiators

Features that make this platform special for the Iranian nutritionist-client market.

| Feature | Value Proposition | Complexity | Offline Sync | RBAC Scope |
|---------|-------------------|------------|--------------|------------|
| OTP via Iranian SMS (Kavenegar/Melipayamak) | Native Iranian auth UX — no email required for clients | Medium | ❌ No | Client auth only |
| Persian full-text search with `pg_trgm` | Search food/medication names in Farsi script | **High** | N/A | All roles |
| Jalali (Shamsi) date support at API level | Dates meaningful to Iranian users | Medium | N/A | Optional query param |
| Offline-first idempotent tracking API | Field use without stable mobile data | **High** | ✅ All tracking tables | Client only |
| `local_id` deduplication for sync | Prevents duplicate entries on reconnect | Medium | ✅ Required | Client tracking endpoints |
| Computed nutritional totals per option/meal/day | Real-time diet plan building feedback for nutritionists | Medium | N/A | Nutritionist plan builder |
| Multi-option meal structure (client picks ONE) | Flexible real-world diets | Medium | N/A | Nutritionist creates; Client reads |
| Nutritional min/max range across meal options | Shows total daily range based on choices | Medium | N/A | Nutritionist plan builder |
| Diet plan day templates (repeating patterns) | Efficient 7-day cycle creation | Medium | N/A | Nutritionist only |
| Food addition request workflow | Client-driven database enrichment | Low | ❌ No | Client requests; Nutritionist approves |
| Web Push notifications (VAPID) | PWA notifications on Android/iOS | Medium | ❌ No | All roles receive; Backend sends |
| Notification preferences per user | Opt-out granular control | Low | N/A | Client self-service |
| Polling-based chat with file attachments | Works offline, simpler than WebSocket | Medium | ✅ Queue messages | Client ↔ Nutritionist (own clients only) |
| Client view of own nutritionist (isolated) | Strict row-level data isolation | Medium | N/A | Row-level authorization |
| Super Admin nutritionist management | Platform operator controls professional access | Low | ❌ No | Super admin only |
| Platform statistics dashboard | Basic operational visibility | Low | ❌ No | Super admin only |
| Scheduled meal + medication reminders | Time-based push from plan data | **High** | ❌ No | Backend scheduler; Client receives |
| Bcrypt password hashing (cost 12) | Strong password security | Low | N/A | Admin + Nutritionist accounts |
| SMS OTP rate limiting (3 per 10 min) | Protects against SMS cost abuse | Low | N/A | Client auth only |
| Refresh token rotation (Redis-backed) | Secure long-lived sessions | Medium | N/A | All roles |

---

## Anti-Features

Things explicitly **out of scope** per PRD. Do not build, do not architect for.

| Anti-Feature | Why Explicitly Avoided | What to Do Instead |
|--------------|------------------------|-------------------|
| WebSocket real-time chat | Polling is sufficient; WebSocket adds complexity and offline issues | 10-second polling when chat open |
| Payment/billing API | Out of scope per PRD §2 | No payment models, no billing tables |
| AI diet recommendations | Out of scope per PRD §2 | Manual nutritionist-created plans only |
| External health device integrations | Out of scope per PRD §2 | Manual client tracking entry only |
| i18n / multi-language | Persian-only platform; YAGNI | Persian strings hardcoded; no translation layer |
| Calorie detection from food photos | Out of scope per PRD §2 | Manual food item entry only |
| Self-registration for any role | All accounts created by higher authority | Nutritionist → super admin; Client → nutritionist |
| Message editing or deletion | Explicit PRD decision | Messages are immutable append-only |
| Jalali-to-Gregorian conversion on backend | Frontend handles Jalali display | Backend stores Gregorian dates; frontend converts |
| Client-side diet plan creation | Nutritionist-only workflow | Clients are plan consumers, not creators |
| Multiple concurrent active plans per client | Explicit PRD decision §6 | Auto-archive previous plan on new plan creation |
| Desktop-optimized endpoints | Mobile-only platform | No special desktop response shaping needed |

---

## Feature Groups with API Endpoint Patterns

### 1. Authentication & Session Management
**Complexity: Medium**  
**Role implications:** Nutritionist/Admin use email+password; Clients use mobile+OTP only

```
POST   /api/v1/auth/login                    → email+password (admin/nutritionist)
POST   /api/v1/auth/otp/send                 → send OTP to mobile (client)
POST   /api/v1/auth/otp/verify               → verify OTP, issue JWT pair (client)
POST   /api/v1/auth/refresh                  → exchange refresh token for new access token
POST   /api/v1/auth/logout                   → invalidate refresh token in Redis
```

**Implementation notes:**
- OTP: 6-digit, 2-minute TTL, max 3 verification attempts, max 3 requests per phone per 10 minutes → stored in Redis
- JWT access token: 15-minute expiry; refresh token: 30-day expiry stored in Redis
- Role embedded in JWT claims (`role: super_admin | nutritionist | client`)
- OTP rate limiting: Redis counter with TTL, atomic INCR
- Kavenegar/Melipayamak adapter pattern — swap SMS provider without code changes
- Iranian mobile format validation: `09[0-9]{9}` (10-digit starting with 09)

---

### 2. Food Database
**Complexity: Medium**  
**Persian search:** `pg_trgm` trigram index on `foods.name`; GIN index recommended

```
GET    /api/v1/foods                         → list with search, category filter, pagination (all authenticated)
GET    /api/v1/foods/:id                     → single food item (all authenticated)
POST   /api/v1/foods                         → create (admin, nutritionist)
PUT    /api/v1/foods/:id                     → update (admin: any; nutritionist: own items only)
DELETE /api/v1/foods/:id                     → soft delete (admin: any; nutritionist: own items only)
```

**Query parameters:**
- `?q=نان` — Persian full-text search via `pg_trgm` `ILIKE '%نان%'` or `similarity()`
- `?category=breakfast` — filter by food category
- `?active=true|false` — filter by active status (admin/nutritionist only)
- `?page=1&per_page=20` — pagination (default 20 per page)

**Nutritional totals computation:**  
When plan items reference food, computed calories/protein/carbs/fat = `(quantity / measurement_amount) * nutrient_per_unit`. This computation happens at the service layer, not in SQL.

---

### 3. Medication Database
**Complexity: Low**  
Same access pattern as food database.

```
GET    /api/v1/medications                   → list with search, pagination (all authenticated)
GET    /api/v1/medications/:id               → single medication (all authenticated)
POST   /api/v1/medications                   → create (admin, nutritionist)
PUT    /api/v1/medications/:id               → update (admin: any; nutritionist: own items only)
DELETE /api/v1/medications/:id               → soft delete (admin: any; nutritionist: own items only)
```

**Query parameters:** `?q=`, `?form=tablet`, `?active=`, `?page=`, `?per_page=`

---

### 4. Client Management (Nutritionist-facing)
**Complexity: Low**  
Row-level isolation: nutritionist can only see their own clients.

```
GET    /api/v1/clients                       → list own clients (nutritionist)
GET    /api/v1/clients/:id                   → client profile (nutritionist: own clients only)
POST   /api/v1/clients                       → register new client (nutritionist)
PUT    /api/v1/clients/:id                   → update client profile (nutritionist: own clients only)
PATCH  /api/v1/clients/:id/status            → activate/deactivate (nutritionist: own clients only)
```

**Query parameters:** `?q=` (name/mobile), `?status=active|inactive`, `?sort=name|last_activity`

---

### 5. Diet Plan Management
**Complexity: HIGH — most complex feature in the system**  
Deeply nested structure; computed totals traverse 4 levels of aggregation.

```
# Plan-level
GET    /api/v1/clients/:id/diet-plans        → list plans (active first; nutritionist)
GET    /api/v1/diet-plans/:id                → full plan with all nested data (nutritionist + client:own)
POST   /api/v1/clients/:id/diet-plans        → create plan (auto-archives previous) (nutritionist)
PUT    /api/v1/diet-plans/:id                → update plan metadata (nutritionist: own clients)
DELETE /api/v1/diet-plans/:id                → archive plan (nutritionist: own clients)

# Plan days
POST   /api/v1/diet-plans/:id/days           → add day (nutritionist)
PUT    /api/v1/plan-days/:id                 → update day (nutritionist)
DELETE /api/v1/plan-days/:id                 → delete day (nutritionist)

# Meals
POST   /api/v1/plan-days/:id/meals           → add meal (nutritionist)
PUT    /api/v1/meals/:id                     → update meal (nutritionist)
DELETE /api/v1/meals/:id                     → delete meal (nutritionist)

# Meal options
POST   /api/v1/meals/:id/options             → add option to meal (nutritionist)
PUT    /api/v1/meal-options/:id              → update option (nutritionist)
DELETE /api/v1/meal-options/:id              → delete option (nutritionist)

# Meal option items
POST   /api/v1/meal-options/:id/items        → add food item to option (nutritionist)
PUT    /api/v1/meal-option-items/:id         → update item quantity/unit (nutritionist)
DELETE /api/v1/meal-option-items/:id         → delete item (nutritionist)

# Prescribed medications (plan-level, not day-level)
POST   /api/v1/diet-plans/:id/medications    → prescribe medication (nutritionist)
PUT    /api/v1/prescribed-medications/:id    → update prescription (nutritionist)
DELETE /api/v1/prescribed-medications/:id    → remove prescription (nutritionist)

# Exercise recommendations (day-level)
POST   /api/v1/plan-days/:id/exercises       → add exercise recommendation (nutritionist)
PUT    /api/v1/exercise-recommendations/:id  → update recommendation (nutritionist)
DELETE /api/v1/exercise-recommendations/:id  → delete recommendation (nutritionist)
```

**Computed nutritional totals — returned inline with plan response:**
- Per `meal_option_item`: `calories = (quantity / food.measurement_amount) * food.calories`
- Per `meal_option`: sum of all items
- Per `meal`: `{min: min(options), max: max(options)}` range across all options
- Per `plan_day`: sum of meal min/max ranges
- All computations done in application layer (service/domain), not SQL

**Auto-archive rule:** `POST /clients/:id/diet-plans` must atomically set all existing `active` plans to `archived` before inserting new plan.

**Client access:** `GET /diet-plans/:id` — client can only fetch their own active plan (or history by plan ID if previously fetched).

---

### 6. Daily Tracking (Client-facing)
**Complexity: Low per endpoint; Medium for sync orchestration**  
**ALL tracking endpoints require `local_id` in request body for idempotency.**

#### Idempotency Pattern (applies to all tracking endpoints)
```
INSERT INTO <table> (..., local_id)
VALUES (...)
ON CONFLICT (client_id, local_id) DO NOTHING
RETURNING *;
```
If `local_id` already exists → return the existing record (200 OK), not an error. This allows clients to safely replay the same record on reconnect.

#### Food Log
```
GET    /api/v1/tracking/food-logs            → list food logs by date range (client: own; nutritionist: own clients)
POST   /api/v1/tracking/food-logs            → log meal selection (client only) [requires local_id]
PUT    /api/v1/tracking/food-logs/:id        → update log entry (client: own only)
DELETE /api/v1/tracking/food-logs/:id        → delete entry (client: own only)
```
Request body includes: `date`, `meal_id`, `selected_option_id` (nullable = skipped), `notes`, `local_id`

#### Water Intake
```
GET    /api/v1/tracking/water                → list water logs by date (client: own; nutritionist: own clients)
POST   /api/v1/tracking/water                → log water intake entry (client only) [requires local_id]
DELETE /api/v1/tracking/water/:id            → delete entry (client: own only)
```
Response includes daily total vs. `daily_water_target_ml` from active plan.

#### Sleep
```
GET    /api/v1/tracking/sleep                → list sleep logs by date range (client: own; nutritionist: own clients)
POST   /api/v1/tracking/sleep                → create/update sleep entry for date (client only) [requires local_id]
PUT    /api/v1/tracking/sleep/:id            → update sleep entry (client: own only)
```
One entry per date per client. `POST` with same `date` should upsert.

#### Exercise
```
GET    /api/v1/tracking/exercise             → list exercise logs by date range (client: own; nutritionist: own clients)
POST   /api/v1/tracking/exercise             → log exercise session (client only) [requires local_id]
PUT    /api/v1/tracking/exercise/:id         → update session (client: own only)
DELETE /api/v1/tracking/exercise/:id         → delete session (client: own only)
```

#### Medication Intake
```
GET    /api/v1/tracking/medications          → list medication logs by date range (client: own; nutritionist: own clients)
POST   /api/v1/tracking/medications          → log medication taken (client only) [requires local_id]
DELETE /api/v1/tracking/medications/:id      → delete entry (client: own only)
```
Supports both prescribed (`prescribed_medication_id` present) and self-reported (only `medication_name` + `dosage`).

#### Body Measurements
```
GET    /api/v1/tracking/measurements         → list measurements by date range (client: own; nutritionist: own clients)
POST   /api/v1/tracking/measurements         → record measurement (client + nutritionist) [requires local_id]
PUT    /api/v1/tracking/measurements/:id     → update measurement (recorded_by only)
```
`recorded_by` auto-set from authenticated user. Both roles can create but data belongs to client.

---

### 7. Lab Results
**Complexity: Medium (file upload validation + storage)**

```
GET    /api/v1/clients/:id/lab-results       → list lab results (nutritionist: own clients; client: own)
GET    /api/v1/lab-results/:id               → single result with file download URL (nutritionist + client:own)
POST   /api/v1/clients/:id/lab-results       → upload lab result (client only; multipart/form-data)
DELETE /api/v1/lab-results/:id               → delete (client: own only)
GET    /api/v1/lab-results/:id/download      → serve file (authenticated, authorized)
```

**File validation:**
- Accepted MIME types: `application/pdf`, `image/jpeg`, `image/png`
- Max size: 10 MB
- Content sniffing (do not trust `Content-Type` header alone — inspect magic bytes)
- Storage path: `/data/uploads/lab-results/{client_id}/{uuid}.{ext}`
- At least one of `file` or `link` must be provided

---

### 8. Messaging System
**Complexity: Medium**  
Polling-based. No WebSocket. Clients can only message their own nutritionist.

```
GET    /api/v1/messages                      → fetch conversation messages (paginated, newest-first cursor)
POST   /api/v1/messages                      → send message (text and/or attachment) [requires local_id for offline queue]
PATCH  /api/v1/messages/read                 → mark messages as read (bulk IDs)
GET    /api/v1/messages/unread-count         → badge count endpoint (polling target)
```

**Polling strategy:**  
Frontend polls `GET /api/v1/messages?since=<last_message_id>` every 10 seconds when chat is open. Backend returns only new messages since cursor. `unread-count` is polled globally for badge.

**File attachments (multipart):**  
- Images: JPG, PNG, max 5 MB → `/data/uploads/messages/{sender_id}/{uuid}.{ext}`
- Files: PDF, max 10 MB → same path pattern

**Push notification trigger:**  
On `POST /api/v1/messages`, backend immediately enqueues Web Push notification to the receiver (if subscribed).

**Authorization:** Client may only send to/receive from their `nutritionist_id`. Nutritionist may only send to/receive from own clients. Verified at handler level.

---

### 9. Food Addition Requests
**Complexity: Low**

```
GET    /api/v1/food-requests                 → list requests (nutritionist: pending requests for own clients)
GET    /api/v1/food-requests/:id             → single request
POST   /api/v1/food-requests                 → submit request (client only)
PATCH  /api/v1/food-requests/:id/approve     → approve + create food item (nutritionist only)
PATCH  /api/v1/food-requests/:id/reject      → reject with optional reason (nutritionist only)
```

**Approval flow:**  
`PATCH /approve` is a **domain transaction**: atomically creates the `foods` record AND sets `food_requests.status = approved` AND sets `reviewed_by`. Triggers push notification to requesting client.

---

### 10. Web Push Notifications
**Complexity: Medium (VAPID key management + scheduled reminders = High)**

```
POST   /api/v1/push/subscribe                → register push subscription (endpoint, p256dh, auth)
DELETE /api/v1/push/subscribe                → unregister subscription (logout or permission revoked)
GET    /api/v1/push/preferences              → get notification preferences
PUT    /api/v1/push/preferences              → update notification preferences (enable/disable per type)
```

**Notification triggers (backend-initiated):**
| Event | Target | Trigger Point |
|-------|--------|---------------|
| New message received | Client / Nutritionist | `POST /messages` handler |
| New diet plan assigned | Client | `POST /diet-plans` handler |
| Food request approved/rejected | Client | `PATCH /food-requests/:id/approve|reject` |
| Meal time reminder | Client | Scheduler goroutine (cron) |
| Medication reminder | Client | Scheduler goroutine (cron) |
| Water intake reminder | Client | Scheduler goroutine (cron) |

**Scheduler notes:**  
Reminders based on `meals.scheduled_time` and `prescribed_medications.times[]` from the client's active diet plan. Cron job wakes every minute, loads due reminders, sends push. Must respect `notification_preferences` per user. Asia/Tehran timezone critical — scheduler must convert UTC times to Tehran local before comparing.

---

### 11. Super Admin Panel
**Complexity: Low**

```
# Nutritionist management
GET    /api/v1/admin/nutritionists           → list all nutritionists
POST   /api/v1/admin/nutritionists           → create nutritionist account
PATCH  /api/v1/admin/nutritionists/:id/status → activate/deactivate
GET    /api/v1/admin/nutritionists/:id/clients → view nutritionist's client list (read-only)

# Platform statistics
GET    /api/v1/admin/stats                   → platform-wide counts
```

---

### 12. File Serving
**Complexity: Low (infrastructure concern)**

```
GET    /api/v1/files/:path                   → serve uploaded file (authenticated + authorized)
```

Authorization rule: User can only access files they own (lab results for own client_id; messages they sent/received). Nutritionists can access their clients' lab results and shared message files.

---

### 13. Health Check
**Complexity: Trivial**

```
GET    /health                               → returns 200 OK with version + DB ping result (no auth)
```

---

## Feature Dependencies

```
Food Database ──────────────────────────────────────────┐
Medication Database ─────────────────────────────────── │
Client Registration (by nutritionist) ──────────────── │
                                                        ▼
                                              Diet Plan CRUD
                                                        │
                          ┌─────────────────────────────┤
                          │                             │
                          ▼                             ▼
                 Daily Tracking APIs           Notification Scheduler
                 (food/water/sleep/            (meal/medication reminders
                  exercise/meds/               derived from plan times)
                  measurements)
                          │
                          ▼
                 Offline Sync (local_id)
                          │
                          ▼
                 Messaging (can reference plan context)

Authentication ──────────────────────────────────────────► EVERYTHING

Food Addition Requests ──► Food Database (on approval)
                       ──► Push Notification (on status change)

Lab Results Upload ──────► File Storage + Download Endpoint

Web Push Subscribe ──────► Notification Scheduler
                       ──► Message Send handler
```

---

## MVP Recommendation

**Build in this order to validate core value quickly:**

**Phase 1 — Auth + User Foundation**
1. JWT + OTP authentication (all three roles)
2. Client registration by nutritionist
3. RBAC middleware

**Phase 2 — Diet Plan Core**
4. Food database with Persian search
5. Diet plan CRUD (full nested structure)
6. Nutritional totals computation

**Phase 3 — Daily Tracking**
7. All 6 tracking endpoints (food log, water, sleep, exercise, meds, body measurements) with `local_id`
8. Lab results upload

**Phase 4 — Communication**
9. Messaging system (polling-based)
10. Food addition request workflow
11. Medication database

**Phase 5 — Engagement**
12. Web Push notification subscription + triggered events
13. Scheduled meal/medication reminders
14. Notification preferences

**Phase 6 — Admin**
15. Super admin panel (nutritionist management + stats)

**Defer to post-launch:**
- Scheduled reminder cron (complex timezone handling; validate demand first)
- Platform statistics beyond simple counts

---

## Persian-Specific API Requirements

| Concern | Requirement | Where It Applies |
|---------|-------------|-----------------|
| **Error messages** | All `message` fields in JSON error responses MUST be in Farsi | Every endpoint |
| **Iranian mobile format** | Validate `09[0-9]{9}` pattern for client mobile numbers | Client registration, OTP send |
| **Jalali dates** | Backend stores and accepts Gregorian (ISO 8601) dates; frontend handles Jalali display — Assumption #5 | No backend conversion needed |
| **Asia/Tehran timezone** | All cron/scheduler logic in Tehran local time (UTC+3:30, DST-aware); `TZ=Asia/Tehran` in all containers | Reminder scheduler, OTP TTL, sleep time parsing |
| **Persian full-text search** | `pg_trgm` extension + GIN index on `foods.name` and `medications.name`; use `similarity()` or `ILIKE '%query%'` — pure trigram works better than `to_tsvector` for Persian since Postgres has no Persian stemmer | Food/medication search |
| **RTL field length** | Persian strings can be longer in bytes (UTF-8 multibyte); ensure `varchar` lengths account for 3 bytes/char in Persian | `foods.name`, `medications.name`, `meals.title` |
| **SMS gateway abstraction** | Kavenegar/Melipayamak behind adapter interface; configurable via env vars | OTP sending |
| **No Gregorian assumption in responses** | Include raw ISO dates in responses, not formatted strings — let frontend format per Jalali | All `date` and `timestamp` fields |

---

## Offline Sync: Endpoints Requiring `local_id` Deduplication

| Table | Endpoint | `local_id` Required | Conflict Strategy |
|-------|----------|--------------------|--------------------|
| `food_logs` | `POST /tracking/food-logs` | ✅ Yes | `ON CONFLICT (client_id, local_id) DO NOTHING` |
| `water_logs` | `POST /tracking/water` | ✅ Yes | `ON CONFLICT (client_id, local_id) DO NOTHING` |
| `sleep_logs` | `POST /tracking/sleep` | ✅ Yes | `ON CONFLICT (client_id, local_id) DO NOTHING` |
| `exercise_logs` | `POST /tracking/exercise` | ✅ Yes | `ON CONFLICT (client_id, local_id) DO NOTHING` |
| `medication_logs` | `POST /tracking/medications` | ✅ Yes | `ON CONFLICT (client_id, local_id) DO NOTHING` |
| `body_measurements` | `POST /tracking/measurements` | ✅ Yes | `ON CONFLICT (client_id, local_id) DO NOTHING` |
| `messages` | `POST /messages` | ✅ Yes | `ON CONFLICT (sender_id, local_id) DO NOTHING` |
| All other tables | — | ❌ No | Standard constraint errors |

**Response behavior on duplicate `local_id`:** Return `200 OK` with the existing record (not `409 Conflict`). This is critical — the client cannot distinguish "first sync" from "retry after network error."

---

## Role Permission Matrix

| Endpoint Group | super_admin | nutritionist | client |
|----------------|:-----------:|:------------:|:------:|
| Auth (all) | ✅ | ✅ | ✅ (OTP only) |
| Food DB read | ✅ | ✅ | ✅ |
| Food DB write | ✅ any | ✅ own | ❌ (request only) |
| Medication DB read | ✅ | ✅ | ✅ |
| Medication DB write | ✅ any | ✅ own | ❌ |
| Client management | ❌ | ✅ own clients | ❌ |
| Diet plan read | ❌ | ✅ own clients | ✅ own plans |
| Diet plan write | ❌ | ✅ own clients | ❌ |
| Daily tracking write | ❌ | ❌ (measurements only) | ✅ own data |
| Daily tracking read | ❌ | ✅ own clients | ✅ own data |
| Lab results upload | ❌ | ❌ | ✅ own |
| Lab results download | ❌ | ✅ own clients | ✅ own |
| Messaging | ❌ | ✅ own clients only | ✅ own nutritionist only |
| Food requests submit | ❌ | ❌ | ✅ |
| Food requests review | ❌ | ✅ own clients | ❌ |
| Push subscribe | ❌ | ✅ | ✅ |
| Admin: nutritionist mgmt | ✅ | ❌ | ❌ |
| Admin: stats | ✅ | ❌ | ❌ |
| Health check | ✅ (no auth) | ✅ (no auth) | ✅ (no auth) |

---

## Complexity Notes by Feature

| Feature | Complexity | Reason |
|---------|------------|--------|
| OTP Auth | Medium | Redis TTL, rate limiting, SMS adapter, Iranian number validation |
| JWT + Refresh | Medium | Redis-backed invalidation, token rotation, 3-role claims |
| Food DB + pg_trgm | Medium | Persian trigram config, GIN index, search ranking |
| Diet Plan nested CRUD | **High** | 5-level hierarchy, transactional writes, auto-archive, computed totals at 4 levels |
| Nutritional totals computation | Medium | Unit conversion, quantity scaling across 4 aggregate levels |
| Daily tracking × 6 | Low each | Simple CRUD + `ON CONFLICT`; complexity is in the pattern repetition |
| Offline sync `local_id` | Medium | Database unique constraint design, conflict response semantics |
| File upload (lab + messages) | Medium | MIME validation, magic byte sniffing, path sanitization, download auth |
| Polling chat | Medium | Cursor-based pagination, unread count, attachment file handling |
| Food request workflow | Low | Simple state machine: pending → approved/rejected + side effect |
| Web Push VAPID | Medium | Key generation, subscription storage, payload construction |
| Reminder scheduler | **High** | Timezone math (Asia/Tehran DST), diet plan time parsing, goroutine safety, user preference checks |
| Super admin panel | Low | Simple CRUD with elevated role checks |
| RBAC middleware | Medium | 3-role JWT, row-level checks (nutritionist→client ownership) at domain service layer |
| Persian error messages | Low | Message catalog in Farsi; no logic complexity |

---

## Sources

- `docs/PRD.md` v1.0 (April 19, 2026) — PRIMARY source, all feature specs derived from here
- `.planning/PROJECT.md` — Confirmed stack decisions, constraints, and out-of-scope items
- PostgreSQL documentation on `pg_trgm` — trigram search for non-Latin scripts (HIGH confidence: well-established for Persian)
- PRD Decision Log §11 — All anti-features traceable to explicit product decisions
