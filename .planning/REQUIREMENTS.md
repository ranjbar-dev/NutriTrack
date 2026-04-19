# Requirements: NutriTrack

**Defined:** 2025-07-18
**Core Value:** Digitalize the nutritionist–client workflow in Iran, replacing WhatsApp + Excel + paper with a structured, offline-capable PWA

## v1 Requirements

Requirements for initial release. Each maps to roadmap phases.

### Authentication & Authorization

- [x] **AUTH-01**: Super Admin can log in with email and password
- [x] **AUTH-02**: Super Admin account is seeded via backend migration or CLI command (no self-registration)
- [x] **AUTH-03**: Nutritionist can log in with email and password
- [x] **AUTH-04**: Nutritionist accounts are created exclusively by Super Admin (no self-registration)
- [x] **AUTH-05**: Client receives OTP via SMS (Kavenegar) to registered mobile number for login
- [x] **AUTH-06**: OTP is 6 digits, valid for 2 minutes, max 3 attempts per code
- [x] **AUTH-07**: OTP rate limit: max 3 requests per phone per 10 minutes
- [x] **AUTH-08**: JWT access token (15 min) + refresh token (30 days) issued on successful auth
- [x] **AUTH-09**: JWT refresh handles concurrent requests without mass logout (refresh queue pattern)
- [x] **AUTH-10**: Passwords hashed with bcrypt (cost factor 12)
- [x] **AUTH-11**: Row-level authorization: nutritionist can only access own clients' data (repository-level enforcement)
- [x] **AUTH-12**: Client cannot self-register — only nutritionist can register clients

### Food Database

- [ ] **FOOD-01**: Shared platform-wide food database accessible to Super Admin and all Nutritionists
- [ ] **FOOD-02**: Food item CRUD with fields: name, categories, calories, protein, carbs, fat, fiber, sugar, sodium, measurement_unit, measurement_amount, description
- [ ] **FOOD-03**: Food categories (many-to-many): breakfast, lunch, dinner, snack, fruit, beverage, supplement, other
- [ ] **FOOD-04**: 12 measurement units supported (gram, kg, tablespoon, teaspoon, cup, piece, slice, palm, matchbox, bowl, ml, liter)
- [ ] **FOOD-05**: Persian full-text search on food name using pg_trgm with correct UTF-8 locale
- [ ] **FOOD-06**: Persian character normalization (ی/ي, ک/ك) at storage and query boundaries
- [ ] **FOOD-07**: Filter by category, active/inactive status
- [ ] **FOOD-08**: Pagination with 20 items per page
- [ ] **FOOD-09**: Soft delete via is_active flag
- [ ] **FOOD-10**: Audit trail: created_by tracks who added each item

### Medication Database

- [ ] **MED-01**: Shared platform-wide medication database accessible to Super Admin and all Nutritionists
- [ ] **MED-02**: Medication CRUD with fields: name, generic_name, form, dosage_unit, description
- [ ] **MED-03**: Medication forms: tablet, capsule, syrup, injection, drop, powder, other
- [ ] **MED-04**: Search on medication name (Persian-aware)
- [ ] **MED-05**: Soft delete via is_active flag

### Diet Plan Management

- [ ] **DIET-01**: Nutritionist can create a diet plan for own client with: start_date, end_date, notes, daily_water_target_ml, status
- [ ] **DIET-02**: Only one active diet plan per client at any time; new plan auto-archives previous
- [ ] **DIET-03**: Diet plan has nested structure: Plan → Days → Meals → Options → Items
- [ ] **DIET-04**: Plan days with day_number supporting unique-per-date and repeating patterns
- [ ] **DIET-05**: Meals with title, scheduled_time, display_order
- [ ] **DIET-06**: Meal options (client picks one per meal) with option_number
- [ ] **DIET-07**: Meal option items link to food database with quantity, measurement_unit, notes
- [ ] **DIET-08**: Real-time computed nutritional totals (calories, protein, carbs, fat, fiber) per option, meal, and day
- [ ] **DIET-09**: Exercise recommendations per plan day: exercise_name, duration_minutes, description, calories_burn_estimate
- [ ] **DIET-10**: Prescribed medications per diet plan: medication, dosage, frequency, times (JSONB), instructions, date range
- [ ] **DIET-11**: Archived plans remain viewable for history (both nutritionist and client)
- [ ] **DIET-12**: Diet plan aggregate loaded in ≤5 queries using batch loading (pgx SendBatch) — no N+1

### Client Tracking

- [ ] **TRACK-01**: Client can log food intake per meal: select option eaten or mark as skipped (optional)
- [ ] **TRACK-02**: Client can log water intake: amount_ml, optional time, multiple entries per day
- [ ] **TRACK-03**: Water intake shows daily total vs target (if set) with visual progress indicator
- [ ] **TRACK-04**: Client can log sleep: sleep_time, wake_time, quality (good/fair/poor), one entry per date
- [ ] **TRACK-05**: Sleep duration computed from sleep_time and wake_time
- [ ] **TRACK-06**: Client can log exercise: exercise_name, duration_minutes, optional calories_burned, multiple per day
- [ ] **TRACK-07**: Client can log medication intake: prescribed or self-reported, dosage, taken_at time
- [ ] **TRACK-08**: Client sees prescribed medications with scheduled times, taps to mark as taken
- [ ] **TRACK-09**: Client or nutritionist can record body measurements: weight, waist, hip, abdomen, thigh, chest, wrist
- [ ] **TRACK-10**: Body measurement history viewable as list and chart (weight over time, measurements over time)
- [ ] **TRACK-11**: One body measurement record per date per field (last write wins or update existing)
- [ ] **TRACK-12**: All tracking entries have local_id (UUID) for offline sync deduplication
- [ ] **TRACK-13**: Nutritionist can view all client tracking data: food logs, weight, measurements, exercise, sleep, water, medication

### Lab Results

- [ ] **LAB-01**: Client can upload lab results with: title, type (blood_test, urine_test, thyroid, hormone, allergy, other), date, file or link
- [ ] **LAB-02**: At least one of file or link must be provided
- [ ] **LAB-03**: Accepted formats: PDF, JPG, PNG; max 10 MB
- [ ] **LAB-04**: Nutritionist can view and download all client lab results
- [ ] **LAB-05**: Files stored on Hetzner server filesystem

### Messaging

- [ ] **MSG-01**: Chat-style messaging between client and assigned nutritionist
- [ ] **MSG-02**: Client can only message their assigned nutritionist; nutritionist messages own clients
- [ ] **MSG-03**: Text messages with optional attachments (images: JPG/PNG max 5 MB, files: PDF max 10 MB)
- [ ] **MSG-04**: Polling-based delivery (every 10 seconds when chat is open)
- [ ] **MSG-05**: Unread message count shown as badge
- [ ] **MSG-06**: Messages ordered chronologically, no editing or deletion
- [ ] **MSG-07**: Read receipts (read_at timestamp)

### Food Requests

- [ ] **FREQ-01**: Client can submit food addition request with name and optional description
- [ ] **FREQ-02**: Request goes to client's assigned nutritionist
- [ ] **FREQ-03**: Nutritionist can approve (creates food item in shared DB) or reject (with optional reason)
- [ ] **FREQ-04**: Client receives notification of approval/rejection

### Client Management

- [x] **CLNT-01**: Nutritionist can register client: full_name, mobile (unique Iranian format), date_of_birth, height_cm, gender, notes
- [ ] **CLNT-02**: Client list view: name, mobile, status, current plan status, last activity, searchable by name/mobile
- [ ] **CLNT-03**: Client list filterable by active/inactive, sortable by name or last activity
- [ ] **CLNT-04**: Client profile view with personal info, current plan summary, and history tabs (weight, food, exercise, water, sleep, medication, lab results, archived plans)
- [ ] **CLNT-05**: Quick actions from client profile: new diet plan, send message, deactivate client
- [ ] **CLNT-06**: Nutritionist can activate/deactivate clients
- [ ] **CLNT-07**: Height and date of birth editable only by nutritionist

### Super Admin Panel

- [ ] **ADMIN-01**: List all nutritionists: name, email, status, client count, created date
- [ ] **ADMIN-02**: Create new nutritionist: name, email, password
- [ ] **ADMIN-03**: Activate/deactivate nutritionist accounts
- [ ] **ADMIN-04**: View nutritionist's client list (read-only)
- [ ] **ADMIN-05**: Full CRUD on food items including ability to edit/delete items created by others
- [ ] **ADMIN-06**: Full CRUD on medications
- [ ] **ADMIN-07**: View audit log of who created/modified food and medication items
- [ ] **ADMIN-08**: Platform statistics: total nutritionists, clients (active/inactive), food items, active diet plans

### Offline & Sync

- [ ] **OFFL-01**: Service Worker caches static assets and API responses for client role only
- [ ] **OFFL-02**: Active diet plan fully viewable offline
- [ ] **OFFL-03**: All tracking (food, water, exercise, sleep, medication, body measurements) works offline with queue
- [ ] **OFFL-04**: Queued messages sendable offline
- [ ] **OFFL-05**: IndexedDB via Dexie.js stores: active plan, pending logs, cached messages, queued outgoing messages
- [ ] **OFFL-06**: Sync manager pushes queued entries on reconnect with last-write-wins conflict resolution
- [ ] **OFFL-07**: Each queued entry has local_id (UUID) + synced_at; server deduplicates via ON CONFLICT (local_id) DO NOTHING
- [ ] **OFFL-08**: Failed sync retries with exponential backoff (max 3 retries, then flag for manual retry)
- [ ] **OFFL-09**: Background Sync API where supported, fallback to polling on app open
- [ ] **OFFL-10**: Diet plan cached on first load, refreshed on app open if online
- [ ] **OFFL-11**: Cached last 50 messages per conversation, fetch new on open
- [ ] **OFFL-12**: Handle iOS PWA storage eviction gracefully (re-fetch on data loss, show pending count)

### Push Notifications

- [ ] **NOTIF-01**: Web Push notifications via VAPID keys using webpush-go
- [ ] **NOTIF-02**: Client subscribes to push on first login (permission prompt)
- [ ] **NOTIF-03**: Notification triggers: new message, new diet plan assigned, food request result
- [ ] **NOTIF-04**: Meal time reminders based on diet plan scheduled times
- [ ] **NOTIF-05**: Medication reminders based on prescribed medication times
- [ ] **NOTIF-06**: Water intake reminders
- [ ] **NOTIF-07**: Client can enable/disable each reminder type in notification preferences
- [ ] **NOTIF-08**: Notification payload includes: title, body, action URL, icon

### UI/UX & PWA

- [x] **UI-01**: Persian-only RTL layout using Tailwind CSS v4 logical properties (ms-, me-, ps-, pe-, text-start, text-end)
- [x] **UI-02**: Mobile-only viewport design (no desktop optimization)
- [x] **UI-03**: Vazirmatn font for all Persian text
- [x] **UI-04**: Shamsi/Jalali calendar for all date displays using jalaali-js
- [x] **UI-05**: Persian numeral display throughout the app
- [ ] **UI-06**: PWA manifest with install prompt
- [ ] **UI-07**: Service worker with registerType: autoUpdate for stale cache prevention
- [ ] **UI-08**: Initial load < 3 seconds on 3G

### Infrastructure & Performance

- [ ] **INFRA-01**: Docker + Docker Compose deployment on Hetzner
- [ ] **INFRA-02**: Traefik reverse proxy with HTTPS (Let's Encrypt)
- [ ] **INFRA-03**: GitLab CI/CD pipeline for automated testing and deployment
- [x] **INFRA-04**: Structured JSON logging to stdout, collected by Loki
- [ ] **INFRA-05**: Grafana dashboards for monitoring
- [x] **INFRA-06**: Health check endpoint
- [ ] **INFRA-07**: Daily automated PostgreSQL backups
- [ ] **INFRA-08**: Weekly file storage backups
- [ ] **INFRA-09**: API response time < 200ms for standard CRUD
- [ ] **INFRA-10**: Diet plan load time < 500ms
- [ ] **INFRA-11**: Support 50 nutritionists, 10,000 clients, ~500 concurrent users

### Security

- [ ] **SEC-01**: All traffic over HTTPS (TLS 1.2+)
- [x] **SEC-02**: Input validation and sanitization on all endpoints
- [x] **SEC-03**: SQL injection prevention via parameterized queries (sqlc)
- [ ] **SEC-04**: File upload validation: type checking, size limits, magic byte verification, UUID filenames
- [ ] **SEC-05**: Content-Disposition: attachment on file downloads to prevent content sniffing
- [x] **SEC-06**: CORS restricted to app domain only
- [x] **SEC-07**: OTP brute force protection (rate limiting + attempt limiting)
- [ ] **SEC-08**: Per-client file storage limits

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Enhanced Tracking

- **TRACK-V2-01**: Chart.js visualizations for all tracking dimensions with trend lines
- **TRACK-V2-02**: Weekly/monthly tracking summaries exportable as PDF
- **TRACK-V2-03**: Client progress score based on adherence to plan

### Enhanced Admin

- **ADMIN-V2-01**: Advanced platform analytics (most used foods, average plan duration, etc.)
- **ADMIN-V2-02**: Bulk food import from CSV
- **ADMIN-V2-03**: System health dashboard in admin panel

### Enhanced Communication

- **MSG-V2-01**: Voice message support
- **MSG-V2-02**: Message search within conversations
- **MSG-V2-03**: Adaptive polling (increase interval when no activity)

## Out of Scope

Explicitly excluded. Documented to prevent scope creep.

| Feature | Reason |
|---------|--------|
| Desktop-optimized UI | Mobile-only design per PRD; YAGNI |
| Multi-language / i18n | Persian-only; target audience is Iranian |
| Real-time video/voice consultation | High complexity, not core value |
| Payment processing / billing | No financial transactions in v1 |
| Wearable / health device integration | High complexity, low priority |
| AI-powered diet recommendations | Explicitly excluded in PRD non-goals |
| Calorie auto-detection from photos | Explicitly excluded in PRD non-goals |
| OAuth / social login | Email/password + OTP sufficient for target users |
| Native mobile app | PWA covers all requirements |
| Real-time WebSocket chat | Polling sufficient per PRD Decision #7 |
| Redis / caching layer | PostgreSQL sufficient at target scale |
| Microservices architecture | Monolith appropriate for scale |
| Recipe management | Not part of nutritionist workflow |
| Grocery list generation | Adds complexity without core value |
| Social features / gamification | Not aligned with professional tool positioning |
| Detailed micronutrients beyond fiber/sugar/sodium | Diminishing returns; can add later |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| AUTH-01 | Phase 1 | ✅ Complete |
| AUTH-02 | Phase 1 | ✅ Complete |
| AUTH-03 | Phase 1 | ✅ Complete |
| AUTH-04 | Phase 1 | ✅ Complete |
| AUTH-05 | Phase 1 | ✅ Complete |
| AUTH-06 | Phase 1 | ✅ Complete |
| AUTH-07 | Phase 1 | ✅ Complete |
| AUTH-08 | Phase 1 | ✅ Complete |
| AUTH-09 | Phase 1 | ✅ Complete |
| AUTH-10 | Phase 1 | ✅ Complete |
| AUTH-11 | Phase 1 | ✅ Complete |
| AUTH-12 | Phase 1 | ✅ Complete |
| FOOD-01 | Phase 2 | Pending |
| FOOD-02 | Phase 2 | Pending |
| FOOD-03 | Phase 2 | Pending |
| FOOD-04 | Phase 2 | Pending |
| FOOD-05 | Phase 2 | Pending |
| FOOD-06 | Phase 2 | Pending |
| FOOD-07 | Phase 2 | Pending |
| FOOD-08 | Phase 2 | Pending |
| FOOD-09 | Phase 2 | Pending |
| FOOD-10 | Phase 2 | Pending |
| MED-01 | Phase 2 | Pending |
| MED-02 | Phase 2 | Pending |
| MED-03 | Phase 2 | Pending |
| MED-04 | Phase 2 | Pending |
| MED-05 | Phase 2 | Pending |
| DIET-01 | Phase 3 | Pending |
| DIET-02 | Phase 3 | Pending |
| DIET-03 | Phase 3 | Pending |
| DIET-04 | Phase 3 | Pending |
| DIET-05 | Phase 3 | Pending |
| DIET-06 | Phase 3 | Pending |
| DIET-07 | Phase 3 | Pending |
| DIET-08 | Phase 3 | Pending |
| DIET-09 | Phase 3 | Pending |
| DIET-10 | Phase 3 | Pending |
| DIET-11 | Phase 3 | Pending |
| DIET-12 | Phase 3 | Pending |
| TRACK-01 | Phase 4 | Pending |
| TRACK-02 | Phase 4 | Pending |
| TRACK-03 | Phase 4 | Pending |
| TRACK-04 | Phase 4 | Pending |
| TRACK-05 | Phase 4 | Pending |
| TRACK-06 | Phase 4 | Pending |
| TRACK-07 | Phase 4 | Pending |
| TRACK-08 | Phase 4 | Pending |
| TRACK-09 | Phase 4 | Pending |
| TRACK-10 | Phase 4 | Pending |
| TRACK-11 | Phase 4 | Pending |
| TRACK-12 | Phase 4 | Pending |
| TRACK-13 | Phase 4 | Pending |
| LAB-01 | Phase 4 | Pending |
| LAB-02 | Phase 4 | Pending |
| LAB-03 | Phase 4 | Pending |
| LAB-04 | Phase 4 | Pending |
| LAB-05 | Phase 4 | Pending |
| MSG-01 | Phase 5 | Pending |
| MSG-02 | Phase 5 | Pending |
| MSG-03 | Phase 5 | Pending |
| MSG-04 | Phase 5 | Pending |
| MSG-05 | Phase 5 | Pending |
| MSG-06 | Phase 5 | Pending |
| MSG-07 | Phase 5 | Pending |
| FREQ-01 | Phase 5 | Pending |
| FREQ-02 | Phase 5 | Pending |
| FREQ-03 | Phase 5 | Pending |
| FREQ-04 | Phase 5 | Pending |
| CLNT-01 | Phase 1 | ✅ Complete |
| CLNT-02 | Phase 5 | Pending |
| CLNT-03 | Phase 5 | Pending |
| CLNT-04 | Phase 5 | Pending |
| CLNT-05 | Phase 5 | Pending |
| CLNT-06 | Phase 5 | Pending |
| CLNT-07 | Phase 5 | Pending |
| ADMIN-01 | Phase 2 | Pending |
| ADMIN-02 | Phase 2 | Pending |
| ADMIN-03 | Phase 2 | Pending |
| ADMIN-04 | Phase 2 | Pending |
| ADMIN-05 | Phase 2 | Pending |
| ADMIN-06 | Phase 2 | Pending |
| ADMIN-07 | Phase 2 | Pending |
| ADMIN-08 | Phase 2 | Pending |
| OFFL-01 | Phase 6 | Pending |
| OFFL-02 | Phase 6 | Pending |
| OFFL-03 | Phase 6 | Pending |
| OFFL-04 | Phase 6 | Pending |
| OFFL-05 | Phase 6 | Pending |
| OFFL-06 | Phase 6 | Pending |
| OFFL-07 | Phase 6 | Pending |
| OFFL-08 | Phase 6 | Pending |
| OFFL-09 | Phase 6 | Pending |
| OFFL-10 | Phase 6 | Pending |
| OFFL-11 | Phase 6 | Pending |
| OFFL-12 | Phase 6 | Pending |
| NOTIF-01 | Phase 6 | Pending |
| NOTIF-02 | Phase 6 | Pending |
| NOTIF-03 | Phase 6 | Pending |
| NOTIF-04 | Phase 6 | Pending |
| NOTIF-05 | Phase 6 | Pending |
| NOTIF-06 | Phase 6 | Pending |
| NOTIF-07 | Phase 6 | Pending |
| NOTIF-08 | Phase 6 | Pending |
| UI-01 | Phase 1 | ✅ Complete |
| UI-02 | Phase 1 | ✅ Complete |
| UI-03 | Phase 1 | ✅ Complete |
| UI-04 | Phase 1 | ✅ Complete |
| UI-05 | Phase 1 | ✅ Complete |
| UI-06 | Phase 6 | Pending |
| UI-07 | Phase 6 | Pending |
| UI-08 | Phase 7 | Pending |
| INFRA-01 | Phase 1 | Pending |
| INFRA-02 | Phase 1 | Pending |
| INFRA-03 | Phase 1 | Pending |
| INFRA-04 | Phase 1 | ✅ Complete |
| INFRA-05 | Phase 7 | Pending |
| INFRA-06 | Phase 1 | ✅ Complete |
| INFRA-07 | Phase 7 | Pending |
| INFRA-08 | Phase 7 | Pending |
| INFRA-09 | Phase 7 | Pending |
| INFRA-10 | Phase 3 | Pending |
| INFRA-11 | Phase 7 | Pending |
| SEC-01 | Phase 1 | Pending |
| SEC-02 | Phase 1 | ✅ Complete |
| SEC-03 | Phase 1 | ✅ Complete |
| SEC-04 | Phase 5 | Pending |
| SEC-05 | Phase 5 | Pending |
| SEC-06 | Phase 1 | ✅ Complete |
| SEC-07 | Phase 1 | ✅ Complete |
| SEC-08 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 130 total (corrected from initial estimate of 103)
- Mapped to phases: 130
- Unmapped: 0 ✓

**Per-phase distribution:**
| Phase | Count | Categories |
|-------|-------|------------|
| Phase 1 | 28 | AUTH (12), CLNT (1), UI (5), INFRA (5), SEC (5) |
| Phase 2 | 23 | FOOD (10), MED (5), ADMIN (8) |
| Phase 3 | 13 | DIET (12), INFRA (1) |
| Phase 4 | 18 | TRACK (13), LAB (5) |
| Phase 5 | 20 | MSG (7), FREQ (4), CLNT (6), SEC (3) |
| Phase 6 | 22 | OFFL (12), NOTIF (8), UI (2) |
| Phase 7 | 6 | INFRA (5), UI (1) |

---
*Requirements defined: 2025-07-18*
*Last updated: 2025-07-18 after roadmap creation — full traceability mapped*
