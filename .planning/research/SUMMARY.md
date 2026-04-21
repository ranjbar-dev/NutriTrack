# Project Research Summary

**Project:** NutriTrack — Go Backend REST API  
**Domain:** Persian Nutrition Management Platform (Nutritionist ↔ Client)  
**Researched:** 2026-04-21  
**Confidence:** HIGH

---

## Executive Summary

NutriTrack is a pure Go backend REST API serving a Persian-language nutrition management platform targeting Iranian nutritionists and their clients (~50 nutritionists, ~10,000 clients, ~500 concurrent users). The correct approach is a **layered Domain-Driven Design (DDD)** architecture with four strictly separated layers: Interface (Gin handlers) → Application (use cases) → Domain (entities, value objects, repo interfaces) → Infrastructure (sqlc/PostgreSQL, Redis, SMS, WebPush, filesystem). The golden rule is inward-only dependencies: `domain` imports nothing from this project; `application` imports only `domain`; `infrastructure` implements `domain` interfaces; `interface` wires everything together via manual DI. This architecture must be established correctly in Phase 1 because violations permeate every subsequent phase.

The stack is fully locked by project requirements: **Go 1.24 + Gin v1.12.0** for HTTP, **sqlc v1.31.0** for type-safe query generation (no ORM), **pgx/v5** for the PostgreSQL driver, **go-redis/v9** for Redis, **zerolog** for structured logging, **golang-migrate** for versioned migrations, and **webpush-go** for VAPID push notifications. Three Iran-specific constraints are non-negotiable and must be implemented in Phase 1: (1) `import _ "time/tzdata"` in main.go plus `apk add tzdata` in Dockerfile to make `Asia/Tehran` work on Alpine; (2) `CREATE EXTENSION IF NOT EXISTS pg_trgm` must be in Migration 001 because Persian full-text search depends on trigram similarity (PostgreSQL has no Persian FTS dictionary); (3) all API error `message` fields must be Persian strings, centralised in an `AppError` catalog rather than hardcoded per handler.

The principal risks are (a) DDD layer contamination — sqlc-generated flat structs must never leak into the domain layer, only mapped through repository adapters; (b) offline sync race conditions — `UNIQUE(client_id, local_id)` must be a DB-level constraint, not application-level check; (c) Diet Plan complexity — a 5-level aggregate (Plan → Days → Meals → Options → Items) that must be split into two aggregates with transaction support to remain manageable; and (d) timezone correctness — Asia/Tehran is UTC+3:30 with DST, and using a fixed offset or missing tzdata silently breaks daily tracking and push notification scheduling.

---

## Key Findings

### Recommended Stack

The stack is fully specified in `STACK.md` and all versions are verified against `proxy.golang.org`. No contested decisions remain — project constraints lock the major choices. The notable research decision was **zerolog over zap** (both are zero-allocation loggers; zerolog wins on API simplicity and chaining ergonomics for a new project). SMS integration uses a raw `net/http` adapter pattern rather than an SDK dependency because Kavenegar's Go SDK is unmaintained.

**Core technologies:**
- `github.com/gin-gonic/gin` v1.12.0 — HTTP router/middleware — project requirement; best-in-class Go REST framework
- `github.com/sqlc-dev/sqlc` v1.31.0 — type-safe SQL code gen CLI — project requirement; zero reflection, compile-time safety
- `github.com/jackc/pgx/v5` v5.9.2 — PostgreSQL driver — required by sqlc v2 config; faster than lib/pq
- `github.com/redis/go-redis/v9` v9.18.0 — Redis client — OTP, rate limiting, JWT token store, caching
- `github.com/rs/zerolog` v1.35.1 — structured JSON logging — zero-allocation, simpler API than zap
- `github.com/golang-jwt/jwt/v5` v5.3.1 — JWT access/refresh tokens — RFC-compliant successor to dgrijalva/jwt-go
- `github.com/golang-migrate/migrate/v4` v4.19.1 — versioned SQL migrations — works alongside sqlc
- `github.com/SherClockHolmes/webpush-go` v1.4.0 — VAPID push notifications — handles key management automatically
- `github.com/spf13/viper` v1.21.0 — config management — 12-factor env vars + .env file in one pass
- PostgreSQL 17-alpine (Docker) — primary data store with pg_trgm for Persian search
- Redis 7-alpine (Docker) — OTP TTL, rate limiting, refresh token store
- **Go 1.24+** — required for range-over-integers, improved type inference, slices/maps packages

**Critical version/config requirement:**  
`sqlc.yaml` must set `emit_pointers_for_null_types: true` — otherwise sqlc generates `pgtype.Text` for nullable columns, which serialises as `{"String":"","Valid":false}` in JSON responses and causes runtime panics.

### Expected Features

See `FEATURES.md` for full endpoint inventory. Summary below:

**Must have (table stakes):**
- JWT + OTP authentication — three roles: super_admin (email+password), nutritionist (email+password), client (mobile+OTP only)
- Diet Plan CRUD with nested structure (5 levels: Plan → Days → Meals → Options → Items) + computed nutritional totals
- Food + Medication shared databases with Persian full-text search (`pg_trgm` trigram similarity)
- Six daily tracking types: food logs, water, sleep, exercise, medications, body measurements — all with `local_id` offline sync
- Lab result upload (PDF/image, local filesystem, UUID-based paths)
- Role-based access control with strict row-level isolation (nutritionist sees only own clients)
- Persian error messages on all API responses
- Asia/Tehran timezone handling throughout

**Should have (differentiators for Iranian market):**
- OTP via Iranian SMS (Kavenegar/Melipayamak adapter) — native Iranian auth UX
- `pg_trgm` Persian full-text search with Arabic/Persian character normalisation at insert + query time
- `local_id` deduplication for offline-first sync (DB-level `UNIQUE` constraint)
- Multi-option meal structure — nutritionist creates alternatives; client picks one per meal
- Computed min/max nutritional range across meal options per day
- Auto-archive previous plan on new plan creation (atomic transaction)
- Polling-based chat with file attachments (10-second polling; no WebSocket)
- Web Push notifications (VAPID) for meal/medication reminders
- Food addition request workflow (client → nutritionist approval)
- Super admin panel for nutritionist management and platform statistics
- Refresh token rotation (Redis-backed, rotate on every use)

**Defer to v2+:**
- WebSocket real-time chat (polling is PRD decision and sufficient for this scale)
- Jalali/Shamsi date conversion on backend (frontend handles display)
- AI diet recommendations
- External health device integrations
- Payment/billing

### Architecture Approach

The architecture follows a **4-layer DDD pattern** with seven bounded contexts (aggregates): User, Food, Medication, DietPlan (deepest — 5 levels), TrackingRecord (6 thin per-type aggregates), Message, and FoodRequest. Dependency injection is manual (no framework) — `internal/interface/bootstrap/wire.go` assembles all services at startup. The folder structure is fully documented in `ARCHITECTURE.md` with explicit package-level boundaries enforced by Go's visibility rules. sqlc-generated code lives in `internal/infrastructure/sqlc/` and is treated as read-only; repository adapters in `internal/infrastructure/repository/` map sqlc flat structs to domain entities. Application-layer ports (`internal/application/ports/`) define interfaces for SMS, push notifications, file storage, and JWT management that infrastructure adapters implement.

**Major components:**
1. `internal/domain/` — 7 aggregate roots with behavior methods, value objects, repository interfaces; zero external imports
2. `internal/application/` — use-case services (commands + queries) orchestrating domain via repository interfaces and ports; imports domain only
3. `internal/infrastructure/` — concrete adapters: sqlc-backed repositories, Redis OTP/token store, Kavenegar SMS adapter, webpush adapter, local filesystem storage
4. `internal/interface/http/` — Gin handlers, middleware chain (Recovery → RequestID → Logging → RateLimit → Auth → RoleCheck), DTOs, centralized error middleware
5. `pkg/` — project-agnostic utilities: Persian text normalisation, Iranian mobile validation, zerolog setup, AppError catalog
6. `migrations/` — golang-migrate versioned SQL files; run as a separate Docker service, never from app startup
7. `docker-compose.yml` — App + PostgreSQL 17 + Redis 7 + Traefik; health-check-gated startup order

### Critical Pitfalls

1. **sqlc structs leaking into domain layer** — Define domain entities as separate structs with behavior; map from `db.*` to `domain.*` exclusively at the repository boundary. Never let `pgtype.*` appear in application or domain code. *Applies: Phase 1 — must get right before any feature work.*

2. **`pg_trgm` extension missing from Migration 001** — Without `CREATE EXTENSION IF NOT EXISTS pg_trgm` in the first migration, Persian food/medication search silently fails in every clean environment (CI, production). Not retrofittable without downtime. *Applies: Phase 1 schema setup.*

3. **Alpine Docker image without tzdata** — `time.LoadLocation("Asia/Tehran")` returns error on Alpine minimal images. Fix: add `import _ "time/tzdata"` to `main.go` (embeds tzdata in binary, +450KB) AND `apk add tzdata` + `ENV TZ=Asia/Tehran` in Dockerfile. Using a fixed offset (`time.FixedZone`) breaks DST — Iran observes IRDT (UTC+4:30) summer/IRST (UTC+3:30) winter. *Applies: Phase 1 — silent data corruption for daily tracking.*

4. **`local_id` offline sync dedup at app level only** — Application-level SELECT-then-INSERT has a race condition under concurrent sync retries. Must enforce with `UNIQUE(client_id, local_id)` DB constraint + `ON CONFLICT DO NOTHING RETURNING *` in all tracking INSERT queries. *Applies: Phase 5 — catastrophic duplicate tracking data.*

5. **OTP attempt counter race condition** — GET+SET pattern on OTP attempt counter allows concurrent requests to bypass rate limits. Must use Redis atomic `INCR` with `EXPIRE` in a pipeline or Lua script. *Applies: Phase 2 auth.*

6. **Centralized Persian error handling** — All API `message` fields must be Persian strings. Define a typed `AppError{Code, MessageFA, HTTPStatus}` catalog in `pkg/apperror/` and use a Gin error-handling middleware. Never hardcode Persian strings in individual handlers. *Applies: Phase 1 — retrofit is 50-file refactor.*

7. **Diet Plan aggregate over-engineering** — Loading the full 5-level tree for every item-level operation (add food to option) causes 6-table JOINs per request. Split into two aggregates: `DietPlan` (root + days + meals + options) and `MealOptionItems` (option + items). Use multiple focused sqlc queries assembled in Go, not one mega-JOIN. *Applies: Phase 4.*

---

## Implications for Roadmap

Based on combined research, the build order is strictly dependency-driven: infrastructure foundation must precede every feature, auth must precede all protected endpoints, shared databases must precede diet plans (which reference them), and tracking must follow diet plans (which tracking logs reference).

### Phase 1: Foundation
**Rationale:** Every subsequent phase depends on this scaffold. DDD layer boundaries, sqlc config, pg_trgm extension, tzdata, Docker Compose with health checks, centralized error handling, and the AppError catalog must all be established here. Mistakes here require 50-file refactors later.  
**Delivers:** Compilable DDD project scaffold, Docker Compose with PG+Redis+health checks, Migration 001 (`pg_trgm` + extensions + `uuid-ossp`), zerolog structured logging, Viper config with fail-fast validation, Persian `AppError` catalog, centralized error middleware, Gin router groups (public/protected), request ID middleware, Makefile targets (`migrate-up`, `sqlc-gen`, `lint`, `test`).  
**Addresses:** All infrastructure setup, timezone handling, Docker startup sequencing.  
**Avoids:** DDD layer contamination, missing pg_trgm, Alpine tzdata failure, migration-on-app-start anti-pattern, secrets in docker-compose.yml.

### Phase 2: Authentication & Session Management
**Rationale:** Every protected endpoint needs auth. JWT middleware, OTP flow, and Redis token store must work before any business logic can be implemented.  
**Delivers:** OTP send/verify (Kavenegar adapter), email+password login (nutritionist/admin), JWT access tokens (15 min), refresh token rotation (Redis-backed, 30 days), logout with JWT blacklist, RBAC middleware (RequireRole), Iranian mobile validation, OTP rate limiting (3 per 10 min, atomic INCR).  
**Uses:** go-redis/v9, golang-jwt/v5, golang.org/x/crypto (bcrypt cost 12), webpush-go (VAPID key generation for later).  
**Avoids:** OTP race condition (atomic INCR), JWT secret validation on startup, no refresh token rotation, auth middleware on wrong router group.

### Phase 3: Shared Databases (Food + Medication)
**Rationale:** Diet plans, tracking logs, and food requests all reference food and medication entities. These shared databases must exist before complex features can be built.  
**Delivers:** Food CRUD with pg_trgm Persian search + Arabic/Persian char normalisation, pagination, category filtering; Medication CRUD with search; soft-delete pattern; nutritionist-owned items (row-level write isolation).  
**Implements:** `domain/food`, `domain/medication`, `application/food`, `application/medication`, sqlc queries with `similarity()` + `ILIKE` fallback.  
**Avoids:** tsvector for Persian (wrong — use pg_trgm), Arabic/Persian Unicode normalization missing at insert+query time, PostgreSQL ENUM for extensible categories (use TEXT+CHECK).

### Phase 4: Diet Plan Management
**Rationale:** Core value proposition of the platform. The most complex phase — 5-level aggregate, transaction logic for auto-archive, computed nutritional totals, and 20+ REST endpoints.  
**Delivers:** Full diet plan CRUD (Plan + Days + Meals + Options + Items), plan auto-archive on new plan creation (transaction), computed nutritional totals per option/meal/day, exercise recommendations per day, prescribed medications per plan, client access to own active plan.  
**Uses:** TxManager pattern (domain/ports), two-aggregate split (DietPlan + MealOptionItems), multiple focused sqlc queries assembled in Go.  
**Avoids:** Single large aggregate loading full tree for item operations, missing transaction support in repositories, mega-JOIN sqlc queries.

### Phase 5: Daily Tracking (Client-facing)
**Rationale:** Six tracking types with offline sync idempotency. All reference diet plan data (meal IDs, option IDs). Relatively simple per endpoint but requires careful DB constraint setup for `local_id`.  
**Delivers:** FoodLog, WaterLog, SleepLog, ExerciseLog, MedicationLog, BodyMeasurement endpoints; `UNIQUE(client_id, local_id)` constraint on all tracking tables; `ON CONFLICT DO NOTHING RETURNING *` idempotency; Tehran date handling; nutritionist read access to client tracking data.  
**Avoids:** App-level-only `local_id` dedup (race condition), UTC date used for "today" (use Tehran timezone), TIMESTAMP without TZ in schema.

### Phase 6: Messaging + Lab Results + Food Requests
**Rationale:** Supporting features that require auth, user data, and file storage infrastructure from prior phases.  
**Delivers:** Polling-based chat between client and nutritionist (10-second poll pattern, immutable messages), file attachment upload (magic byte MIME validation, UUID-based server paths, `http.MaxBytesReader`), lab result upload/download (PDF/image, `Content-Disposition: attachment`), food addition request workflow (submit → approve/reject).  
**Avoids:** MIME type spoofing (magic bytes, not Content-Type header), path traversal (UUID paths only), file size DoS (`http.MaxBytesReader`), serving files without Content-Disposition.

### Phase 7: Notifications + Admin Panel
**Rationale:** Push notifications require VAPID keys, push subscription storage, and a scheduler that references diet plan data. Admin panel is straightforward CRUD with super_admin role guard.  
**Delivers:** Web Push subscription registration/management (VAPID), scheduled meal/medication reminder push notifications, notification preferences per user, super admin APIs for nutritionist CRUD and platform statistics.  
**Uses:** webpush-go v1.4.0, background scheduler goroutine (cron-style), `domain/notification` entities.  
**Avoids:** TZ=Asia/Tehran missing from scheduler (reminders fire at wrong time), DST miscalculation using FixedZone.

### Phase 8: Hardening
**Rationale:** Security audit, performance validation, and production readiness before launch.  
**Delivers:** Security audit (RBAC coverage, row-level auth verification, JWT blacklist correctness), CI enforcement of `sqlc generate` freshness (`git diff --exit-code`), rate limiting coverage, performance testing at target scale (500 concurrent), Docker non-root user, volume permissions, structured logging completeness, `.env.example` documentation.  
**Avoids:** All "MINOR" pitfalls deferred from earlier phases, secrets committed to git, pgtype leaking into response JSON.

---

### Phase Ordering Rationale

- **Foundation before everything:** DDD layer structure, pg_trgm, tzdata, and Docker Compose startup sequencing are load-bearing for all subsequent phases. Retrofitting these is extremely expensive.
- **Auth before any protected endpoints:** JWT middleware and RBAC are prerequisites for every business-logic endpoint.
- **Food + Medication before Diet Plans:** Diet plan meals reference food items; prescriptions reference medications. Referential integrity requires these tables to exist first.
- **Diet Plans before Tracking:** Tracking food logs reference meal IDs and option IDs from diet plans. Some tracking types can be built without plans (water, sleep, exercise) but the full feature requires plan data.
- **Messaging/Labs/Requests after core features:** These are supporting workflows that depend on User and DietPlan being stable.
- **Notifications after Diet Plans:** Scheduled reminders are generated from diet plan meal times and prescribed medication schedules.
- **Hardening last:** Polish and security audit after all functionality is proven.

---

### Research Flags

**Phases needing deeper research before planning:**
- **Phase 7 (Notifications scheduler):** Background goroutine scheduling strategy in Go (cron vs. ticker vs. dedicated scheduler library) and VAPID key rotation policy are not fully specified. Needs `/gsd-research-phase` to validate push scheduler architecture before implementation.
- **Phase 6 (Messaging polling):** Exact polling strategy (long-poll vs. short-poll, unread count endpoint caching) may benefit from a quick design spike.

**Phases with well-documented patterns (skip research-phase):**
- **Phase 1 (Foundation):** Go DDD scaffold, Docker Compose health checks, sqlc config, zerolog — all established patterns with HIGH confidence documentation.
- **Phase 2 (Auth):** JWT + OTP with Redis is a standard Go pattern. Every decision is documented in PITFALLS.md with code examples.
- **Phase 3 (Shared DBs):** pg_trgm search, CRUD with soft delete — standard patterns.
- **Phase 5 (Tracking):** `ON CONFLICT DO NOTHING` idempotency is a well-understood PostgreSQL pattern.
- **Phase 8 (Hardening):** Security checklist items are standard; no research needed, just execution.

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | All versions verified against proxy.golang.org; no contested decisions |
| Features | HIGH | Derived directly from PRD v1.0 + PROJECT.md; all requirements explicit |
| Architecture | HIGH | DDD Go patterns are well-established; folder structure fully specified in ARCHITECTURE.md |
| Pitfalls | HIGH | All pitfalls grounded in stack-specific known issues with code-level prevention strategies |

**Overall confidence:** HIGH

### Gaps to Address

- **Push notification scheduler implementation:** No concrete implementation decided (cron library vs. custom goroutine scheduler). Validate during Phase 7 planning.
- **Conflict resolution rules for offline sync:** Partially defined in PITFALLS.md (last-write-wins vs. additive per entity type). Must be documented as explicit requirements before Phase 5 implementation starts to avoid mid-sprint redesign.
- **SMS provider failover:** Kavenegar + Melipayamak adapter pattern is specified but failover logic (retry on provider failure, switch to backup) is not defined. Decide before Phase 2 implementation.
- **Jalali date support scope:** PROJECT.md says "frontend handles Jalali display" but FEATURES.md lists Jalali as a differentiator. Confirm: backend returns Gregorian only, no Jalali conversion endpoints needed.

---

## Sources

### Primary (HIGH confidence)
- `STACK.md` — All library versions verified against proxy.golang.org; zerolog vs. zap decision documented
- `FEATURES.md` — Full endpoint inventory derived from PRD v1.0 + PROJECT.md; all 80+ endpoints listed
- `ARCHITECTURE.md` — Complete DDD folder structure, 7 aggregate boundaries, 4 data flow diagrams
- `PITFALLS.md` — 35+ pitfalls across 10 domains with Go code examples; phase-specific warning matrix
- `PROJECT.md` — Authoritative requirements, constraints, and key decisions for NutriTrack

### Secondary (MEDIUM confidence)
- Go DDD community patterns — Matt Boyle's DDD in Go; go-ddd reference implementations
- sqlc GitHub issues — nullable pgtype handling, pgx/v5 migration notes
- PostgreSQL documentation — pg_trgm operator reference, TIMESTAMPTZ behavior, ENUM vs CHECK tradeoffs
- Redis documentation — INCR/EXPIRE atomicity, pipeline/TxPipeline semantics

### Tertiary (LOW confidence — needs validation during implementation)
- Push notification scheduler strategy — multiple approaches possible; needs spike
- SMS provider failover behavior — Kavenegar-specific failure modes not documented

---
*Research completed: 2026-04-21*  
*Ready for roadmap: yes*
