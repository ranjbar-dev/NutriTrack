# Phase 1: Foundation & Infrastructure - Context

**Gathered:** 2025-07-18
**Status:** Ready for planning

<domain>
## Phase Boundary

All three user roles (Super Admin, Nutritionist, Client) can authenticate and access their respective dashboards through a deployed, CI/CD-backed Persian RTL application. This phase establishes the monorepo structure, authentication system, Persian/RTL foundation, Docker deployment, and CI/CD pipeline.

</domain>

<decisions>
## Implementation Decisions

### Authentication Architecture
- **D-01:** JWT tokens stored in httpOnly secure cookies (not localStorage) — prevents XSS token theft. Access token (15 min) in cookie, refresh token (30 days) in separate httpOnly cookie.
- **D-02:** Token refresh uses a mutex/queue pattern on the frontend — when access token expires, the first request triggers refresh while subsequent requests wait, preventing concurrent refresh storms and mass logout (AUTH-09).
- **D-03:** Super Admin account seeded via a CLI command (`go run cmd/seed/main.go`) executed during initial deployment — no migration-based seeding. Email/password provided via environment variables.
- **D-04:** OTP delivery via Kavenegar REST API. OTP service abstracts the SMS gateway behind an interface (`SMSSender`) for testability. In development, OTPs log to stdout instead of sending SMS.
- **D-05:** Client registration is nutritionist-initiated only (AUTH-12). Nutritionist creates client with mobile number; client can then request OTP to that number.

### Project Structure
- **D-06:** Monorepo layout — `backend/` (Go) and `frontend/` (Nuxt 4) in a single repository. Shared Docker Compose orchestrates both.
- **D-07:** Go backend follows handler → service → repository layered architecture. Directory structure: `backend/cmd/api/` (entrypoint), `backend/cmd/seed/` (admin seeder), `backend/internal/handler/`, `backend/internal/service/`, `backend/internal/repository/`, `backend/internal/model/`, `backend/internal/middleware/`, `backend/internal/config/`, `backend/db/migrations/`, `backend/db/queries/` (sqlc SQL files).
- **D-08:** Nuxt 4 uses the `app/` directory structure: `frontend/app/pages/`, `frontend/app/components/`, `frontend/app/composables/`, `frontend/app/stores/`, `frontend/app/layouts/`, `frontend/app/middleware/`, `frontend/app/plugins/`, `frontend/app/assets/css/`.
- **D-09:** sqlc for all database queries — write raw SQL in `backend/db/queries/`, generate type-safe Go code. No ORM.

### Persian RTL Foundation
- **D-10:** `dir="rtl"` and `lang="fa"` set on the root `<html>` element. Tailwind CSS v4 logical properties used everywhere (`ms-`, `me-`, `ps-`, `pe-`, `text-start`, `text-end`) — no physical LTR properties (`ml-`, `mr-`, `pl-`, `pr-`).
- **D-11:** Vazirmatn font loaded via npm package (`vazirmatn`), imported as a variable font in the global CSS. Applied as the default `font-family` via Tailwind theme extension.
- **D-12:** `useShamsiDate()` composable wrapping `jalaali-js` for all date displays. All dates stored as Gregorian in PostgreSQL, converted to Shamsi only at the display layer.
- **D-13:** `toPersianDigits(value)` utility function converts Latin digits (0-9) to Persian (۰-۹). Used globally via a Vue directive (`v-fa-digits`) or applied in display composables.
- **D-14:** Mobile-only viewport — `<meta name="viewport" content="width=device-width, initial-scale=1">`. No desktop breakpoints. Max width constrained to mobile proportions.

### Role-based Layouts & Navigation
- **D-15:** Three separate Nuxt layouts: `admin.vue`, `nutritionist.vue`, `client.vue`, plus `auth.vue` for login/OTP pages. Each layout has its own bottom navigation bar.
- **D-16:** Mobile bottom navigation bar pattern for all roles. Super Admin: 3-4 tabs (nutritionists, foods, stats). Nutritionist: 4-5 tabs (clients, foods, medications, messages, profile). Client: 4-5 tabs (plan, tracking, messages, profile).
- **D-17:** Nuxt route middleware for role-based access control. `auth` middleware checks JWT cookie validity. `role` middleware redirects unauthorized roles. Unauthenticated users redirect to login.
- **D-18:** After login, each role redirects to their default dashboard page: Super Admin → `/admin`, Nutritionist → `/nutritionist/clients`, Client → `/client/plan`.

### Deployment & Infrastructure
- **D-19:** Docker Compose with 4 services: `api` (Go, multi-stage build ~20MB image), `frontend` (Nuxt, Node-based), `postgres` (PostgreSQL 16 with named volume), `traefik` (v3, auto Let's Encrypt via ACME HTTP challenge).
- **D-20:** Traefik configured via Docker labels — no separate config file. HTTPS redirect enforced. API routes under `api.domain.com` or `domain.com/api/`, frontend under `domain.com`.
- **D-21:** GitLab CI/CD pipeline with 4 stages: `lint` (golangci-lint + ESLint), `test` (Go tests + Vitest), `build` (Docker images), `deploy` (SSH to Hetzner, docker compose pull + up).
- **D-22:** Health check endpoint at `GET /api/health` returns `{"status": "ok", "timestamp": "..."}` with 200. Used by Docker health checks and monitoring.
- **D-23:** Structured JSON logging via zerolog to stdout. Log fields: `timestamp`, `level`, `method`, `path`, `status`, `duration_ms`, `request_id`. Collected by Loki via Docker log driver.

### Security Foundation
- **D-24:** CORS restricted to the frontend domain only. Credentials allowed (for httpOnly cookies).
- **D-25:** Input validation via `go-playground/validator` struct tags on all request DTOs. Custom validators for Iranian mobile format (`^09[0-9]{9}$`), Persian text fields.
- **D-26:** All SQL queries via sqlc — parameterized by design, preventing SQL injection (SEC-03).
- **D-27:** Rate limiting middleware on OTP endpoints: max 3 requests per phone per 10 minutes, max 3 verification attempts per code (AUTH-06, AUTH-07, SEC-07). In-memory rate limiter sufficient for v1 scale.

### Agent's Discretion
- Loading skeleton/spinner design for pages
- Exact Tailwind color palette and design tokens
- Error page styling (404, 500, network error)
- Exact CI/CD runner tags and Docker registry choice
- Go module path naming
- Exact health check response fields beyond status

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Product Requirements
- `docs/PRD.md` — Full product requirements, data model, feature specs, and decision log
- `.planning/REQUIREMENTS.md` — 130 v1 requirements with traceability to phases

### Architecture & Stack
- `.planning/research/ARCHITECTURE.md` — System architecture, layered Go backend, Nuxt 4 frontend patterns
- `.planning/research/STACK.md` — Technology stack with verified library versions and installation commands
- `.planning/research/PITFALLS.md` — Known pitfalls and risks for the chosen stack

### Phase Scope
- `.planning/ROADMAP.md` §Phase 1 — Phase 1 requirements list, success criteria, and dependencies

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- None — greenfield project. All code will be created from scratch.

### Established Patterns
- None — this phase establishes the foundational patterns that all subsequent phases will follow.

### Integration Points
- None — this is the first phase. It creates the integration points (auth middleware, layout system, API structure) that subsequent phases will consume.

</code_context>

<specifics>
## Specific Ideas

- User explicitly requires **Gin** framework (not Fiber or Echo) for all API routes
- User explicitly requires **RTL** layout for all UI design
- UI/UX guidance should follow the skill from `https://github.com/nextlevelbuilder/ui-ux-pro-max-skill`
- Persian numerals must display throughout the entire app (not just dates)
- Vazirmatn is the designated Persian font
- Kavenegar is the SMS gateway for Iranian OTP delivery

</specifics>

<deferred>
## Deferred Ideas

- Persian pg_trgm search validation spike — Phase 2 concern (noted in STATE.md blockers)
- Plan builder UI state management complexity — Phase 3 concern
- iOS PWA storage eviction testing — Phase 6 concern
- Chart.js weight/measurement visualizations — Phase 4
- WebPush notification infrastructure — Phase 6

</deferred>

---

*Phase: 01-foundation-infrastructure*
*Context gathered: 2025-07-18*
