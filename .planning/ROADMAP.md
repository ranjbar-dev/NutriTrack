# Roadmap: NutriTrack

## Overview

NutriTrack delivers a Persian-only, mobile-first PWA for nutritionist–client management in 7 phases following a strict dependency chain: foundation (auth, infrastructure, RTL) → shared data domains (food, medication, admin) → the diet plan engine (highest-complexity core) → parallel feature layers (tracking + communication) → offline/PWA wrapper → hardening for production launch. Each phase produces a deployable, verifiable increment. The diet plan engine (Phase 3) is the technical complexity nexus; offline support (Phase 6) is the highest-value differentiator.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [ ] **Phase 1: Foundation & Infrastructure** - Auth for all 3 roles, RTL/Persian base, Docker deployment, CI/CD pipeline
- [ ] **Phase 2: Core Data Domain** - Shared food & medication databases, Super Admin panel with platform management
- [ ] **Phase 3: Diet Plan Engine** - Nested plan builder (Plan→Days→Meals→Options→Items), nutritional computation, client plan view
- [ ] **Phase 4: Client Tracking Suite** - Six tracking dimensions (food, water, sleep, exercise, medication, body) + lab results
- [ ] **Phase 5: Communication & Collaboration** - Messaging, food requests, lab results, client management dashboard
- [ ] **Phase 6: Offline & PWA** - Service worker, IndexedDB sync queue, push notifications, PWA install
- [ ] **Phase 7: Hardening & Launch** - Security audit, performance validation, monitoring, backups, production deployment

## Phase Details

### Phase 1: Foundation & Infrastructure
**Goal**: All three user roles can authenticate and access their respective dashboards through a deployed, CI/CD-backed Persian RTL application
**Depends on**: Nothing (first phase)
**Requirements**: AUTH-01, AUTH-02, AUTH-03, AUTH-04, AUTH-05, AUTH-06, AUTH-07, AUTH-08, AUTH-09, AUTH-10, AUTH-11, AUTH-12, CLNT-01, UI-01, UI-02, UI-03, UI-04, UI-05, INFRA-01, INFRA-02, INFRA-03, INFRA-04, INFRA-06, SEC-01, SEC-02, SEC-03, SEC-06, SEC-07
**Success Criteria** (what must be TRUE):
  1. Super Admin can log in with email/password and access the admin layout; nutritionist accounts can only be created by Super Admin
  2. Nutritionist can log in with email/password, register a new client (name, mobile, DOB, height, gender), and access the nutritionist layout
  3. Client can request OTP via SMS, verify it, and access the client layout — with rate limiting enforced (max 3 requests per phone per 10 minutes, max 3 attempts per code)
  4. All pages render in Persian RTL layout with Vazirmatn font, Shamsi/Jalali dates, Persian numerals, and mobile-only viewport — JWT refresh handles concurrent requests without mass logout
  5. Application deploys via Docker Compose (Go API + PostgreSQL + Traefik HTTPS), health check returns OK, GitLab CI/CD pipeline runs green, structured JSON logs flow to stdout
**Plans**: 6 plans

Plans:
- [x] 01-01-PLAN.md — Monorepo scaffold, Go backend bootstrap, DB migrations, sqlc setup
- [x] 01-02-PLAN.md — Nuxt 4 frontend, Tailwind v4 RTL, Vazirmatn, Persian utilities, layouts
- [x] 01-03-PLAN.md — JWT package, SMS sender, auth & infrastructure middleware suite
- [ ] 01-04-PLAN.md — Auth services, handlers, router wiring, Super Admin seeder
- [ ] 01-05-PLAN.md — Frontend auth store, login/OTP pages, route guards
- [ ] 01-06-PLAN.md — Docker, Traefik, docker-compose, GitLab CI/CD

**UI hint**: yes

### Phase 2: Core Data Domain
**Goal**: Nutritionists and Super Admin can populate and manage the shared food and medication databases, and Super Admin has full platform operational control
**Depends on**: Phase 1
**Requirements**: FOOD-01, FOOD-02, FOOD-03, FOOD-04, FOOD-05, FOOD-06, FOOD-07, FOOD-08, FOOD-09, FOOD-10, MED-01, MED-02, MED-03, MED-04, MED-05, ADMIN-01, ADMIN-02, ADMIN-03, ADMIN-04, ADMIN-05, ADMIN-06, ADMIN-07, ADMIN-08
**Success Criteria** (what must be TRUE):
  1. Nutritionist can add, edit, search (Persian fuzzy via pg_trgm with character normalization ی/ي, ک/ك), and soft-delete food items with full nutritional data (calories, protein, carbs, fat, fiber, sugar, sodium), multiple categories, and 12 measurement units
  2. Nutritionist can add, edit, search, and soft-delete medications with form types (tablet, capsule, syrup, etc.)
  3. Super Admin can create/activate/deactivate nutritionist accounts, view their client lists (read-only), view platform statistics (total nutritionists, clients, food items, active plans), and edit/delete food and medication items created by any user with full audit trail
  4. Food list supports pagination (20 items/page), category filtering, active/inactive filtering, and audit tracking of who created each item
**Plans**: TBD
**UI hint**: yes

### Phase 3: Diet Plan Engine
**Goal**: Nutritionists can create complete multi-day diet plans with the full nested structure and clients can view their active plan with real-time computed nutritional data
**Depends on**: Phase 2
**Requirements**: DIET-01, DIET-02, DIET-03, DIET-04, DIET-05, DIET-06, DIET-07, DIET-08, DIET-09, DIET-10, DIET-11, DIET-12, INFRA-10
**Success Criteria** (what must be TRUE):
  1. Nutritionist can build a complete diet plan with the full nested structure: days (with repeating patterns) → meals (title, time, order) → options (client picks one) → food items (from shared DB, with quantity/unit) — plus exercise recommendations per day and medication prescriptions per plan
  2. Nutritional totals (calories, protein, carbs, fat, fiber) compute correctly and display in real-time at option, meal, and day levels as food items are added or modified
  3. Creating a new active plan auto-archives the previous one; the one-active-plan-per-client constraint holds at both application and database level; archived plans remain viewable for history
  4. Client can view their active diet plan on mobile with day navigation, meal options with food items and nutrition info, exercise recommendations, medication schedule, and water target
  5. Full diet plan aggregate loads in ≤500ms via batch queries (≤5 queries, no N+1) — validated with realistic data (7 days × 5 meals × 3 options × 4 items)
**Plans**: TBD
**UI hint**: yes

### Phase 4: Client Tracking Suite
**Goal**: Clients can log all daily health activities and nutritionists can view comprehensive tracking history for each client
**Depends on**: Phase 3
**Requirements**: TRACK-01, TRACK-02, TRACK-03, TRACK-04, TRACK-05, TRACK-06, TRACK-07, TRACK-08, TRACK-09, TRACK-10, TRACK-11, TRACK-12, TRACK-13, LAB-01, LAB-02, LAB-03, LAB-04, LAB-05
**Success Criteria** (what must be TRUE):
  1. Client can log food intake per meal (select option eaten or skip), water intake with timestamps and visual progress toward daily target, sleep with time pickers and auto-computed duration, exercise entries with optional calorie burn, and medication doses (tap-to-mark prescribed medications at scheduled times + self-reported)
  2. Client or nutritionist can record body measurements (weight, waist, hip, abdomen, thigh, chest, wrist) with one record per date per field; weight history displays as a chart with Shamsi dates
  3. Client can upload lab results (PDF/JPG/PNG up to 10MB, or link) with title and type; nutritionist can view and download all client lab results
  4. Client daily dashboard shows today's summary across all tracking types with quick-log actions; nutritionist can view all client tracking data with date range filtering
  5. All tracking entries support local_id (UUID) for offline sync deduplication — POSTing the same local_id twice returns the existing record (idempotent upsert infrastructure for Phase 6)
**Plans**: TBD
**UI hint**: yes

### Phase 5: Communication & Collaboration
**Goal**: Clients and nutritionists can communicate via messaging, clients can request food additions, and nutritionists have a complete client management workspace
**Depends on**: Phase 3 (parallel-safe with Phase 4 — no shared data dependencies beyond Phase 3)
**Requirements**: MSG-01, MSG-02, MSG-03, MSG-04, MSG-05, MSG-06, MSG-07, FREQ-01, FREQ-02, FREQ-03, FREQ-04, CLNT-02, CLNT-03, CLNT-04, CLNT-05, CLNT-06, CLNT-07, SEC-04, SEC-05, SEC-08
**Success Criteria** (what must be TRUE):
  1. Client and assigned nutritionist can exchange text messages with optional attachments (images JPG/PNG ≤5MB, files PDF ≤10MB); messages appear via 10-second polling with unread badge count and read receipts; messages are chronological, no editing or deletion
  2. Client can submit food addition requests (name + optional description); nutritionist can approve (creates food item in shared DB) or reject (with optional reason); client receives notification of the result
  3. Nutritionist client list shows name, mobile, status, plan status, last activity — searchable by name/mobile, filterable by active/inactive, sortable; client profile shows personal info, current plan summary, and history tabs (all tracking data from Phase 4, archived plans)
  4. Quick actions from client profile work: create new diet plan, send message, activate/deactivate client; height and date of birth editable only by nutritionist
  5. File uploads validated via magic byte verification, size limits enforced, UUID filenames used, Content-Disposition: attachment on downloads, per-client storage limits enforced
**Plans**: TBD
**UI hint**: yes

### Phase 6: Offline & PWA
**Goal**: Clients can view diet plans and log all tracking data while offline with automatic sync on reconnect, and receive push notification reminders for meals, medications, and messages
**Depends on**: Phase 4, Phase 5 (offline wraps all API endpoints from both phases)
**Requirements**: OFFL-01, OFFL-02, OFFL-03, OFFL-04, OFFL-05, OFFL-06, OFFL-07, OFFL-08, OFFL-09, OFFL-10, OFFL-11, OFFL-12, NOTIF-01, NOTIF-02, NOTIF-03, NOTIF-04, NOTIF-05, NOTIF-06, NOTIF-07, NOTIF-08, UI-06, UI-07
**Success Criteria** (what must be TRUE):
  1. PWA installs as standalone app on Android Chrome and iOS Safari; service worker caches static assets with autoUpdate strategy; PWA manifest displays Persian app name with correct icons and theme
  2. Client can view their active diet plan and last 50 cached messages per conversation with no network connection; diet plan cached on first load, refreshed on app open if online
  3. All client tracking (food, water, sleep, exercise, medication, body measurements) and outgoing messages work offline — entries queue in IndexedDB (Dexie.js) and sync on reconnect within 30 seconds, with server-side deduplication via local_id, exponential backoff retry (max 3), and Background Sync API where supported (polling fallback)
  4. Push notifications fire for new messages, new diet plans, food request results, meal time reminders, medication reminders, and water intake reminders — client can enable/disable each reminder type in notification preferences
  5. iOS PWA storage eviction handled gracefully: re-fetch on data loss, show pending sync count; sync status indicator shows syncing/synced/pending with manual retry for failed items
**Plans**: TBD
**UI hint**: yes

### Phase 7: Hardening & Launch
**Goal**: Application is production-ready with verified security, validated performance under load, live monitoring, automated backups, and polished end-to-end user experience
**Depends on**: Phase 6
**Requirements**: INFRA-05, INFRA-07, INFRA-08, INFRA-09, INFRA-11, UI-08
**Success Criteria** (what must be TRUE):
  1. Security audit passes: row-level authorization verified across all 30+ endpoints (no cross-nutritionist data access), all queries parameterized, dependency vulnerability scans clean (govulncheck, npm audit)
  2. All API endpoints respond under 200ms at p95 under load (500 concurrent users simulated); PWA initial load under 3 seconds on simulated 3G after first visit
  3. Grafana dashboards show live API response times, error rates, active users, and DB pool utilization; Loki receives structured JSON logs from all containers; alerts fire on threshold breaches (>1% 5xx, >1s p95)
  4. Daily automated PostgreSQL backups and weekly file storage backups run on schedule and are verified restorable (actual restore test passes)
  5. Full end-to-end user journey passes for all three roles on real Android and iOS mobile devices with natural Persian text, proper error handling, loading states, and empty states
**Plans**: TBD

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6 → 7
Note: Phases 4 and 5 are parallel-safe (no shared data dependencies beyond Phase 3) but execute sequentially for solo developer workflow.

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation & Infrastructure | 0/6 | Planned | - |
| 2. Core Data Domain | 0/? | Not started | - |
| 3. Diet Plan Engine | 0/? | Not started | - |
| 4. Client Tracking Suite | 0/? | Not started | - |
| 5. Communication & Collaboration | 0/? | Not started | - |
| 6. Offline & PWA | 0/? | Not started | - |
| 7. Hardening & Launch | 0/? | Not started | - |

---
*Roadmap created: 2025-07-18*
*Last updated: 2025-07-19*
