# Phase 1: Foundation & Infrastructure - Research

**Researched:** 2025-07-18
**Domain:** Go/Gin API + Nuxt 4 PWA Foundation — Auth, RTL/Persian, Docker, CI/CD
**Confidence:** HIGH

## Summary

Phase 1 establishes the entire project foundation: monorepo structure, Go backend with 3-layer architecture (Handler → Service → Repository), Nuxt 4 frontend with RTL Persian layout, complete authentication for 3 roles (Super Admin email/password, Nutritionist email/password, Client OTP via Kavenegar), Docker Compose deployment with Traefik HTTPS, and GitLab CI/CD pipeline. This phase produces 28 requirements across AUTH (12), CLNT (1), UI (5), INFRA (5), and SEC (5) categories.

The tech stack is well-established and all libraries are verified current. The primary risk is the **Go version requirement**: Gin v1.12.0 requires Go 1.25+, but the local environment has Go 1.24.2. Go 1.25.9 is available for download and upgrading is required. The secondary risk is the JWT refresh token race condition (Pitfall 7 from PITFALLS.md) — the frontend refresh queue pattern must be implemented correctly from day one to prevent mass logouts.

**Primary recommendation:** Scaffold monorepo with `backend/` (Go 1.25+, Gin v1.12.0) and `frontend/` (Nuxt 4.4.x), implement auth flows with httpOnly cookie JWT first, then build out RTL/Persian foundation, then Docker/CI. Upgrade Go to 1.25+ before any backend code.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** JWT tokens stored in httpOnly secure cookies (not localStorage) — access token (15 min), refresh token (30 days) in separate httpOnly cookies
- **D-02:** Token refresh uses mutex/queue pattern on frontend — first request triggers refresh, others wait
- **D-03:** Super Admin seeded via CLI command (`go run cmd/seed/main.go`) — env vars for email/password
- **D-04:** OTP via Kavenegar REST API with `SMSSender` interface abstraction; dev mode logs to stdout
- **D-05:** Client registration is nutritionist-initiated only (AUTH-12)
- **D-06:** Monorepo: `backend/` (Go) + `frontend/` (Nuxt 4) in single repository
- **D-07:** Handler → Service → Repository layers. Structure: `backend/cmd/api/`, `backend/cmd/seed/`, `backend/internal/{handler,service,repository,model,middleware,config}/`, `backend/db/{migrations,queries}/`
- **D-08:** Nuxt 4 `app/` directory: `frontend/app/{pages,components,composables,stores,layouts,middleware,plugins,assets/css}/`
- **D-09:** sqlc for all database queries — SQL in `backend/db/queries/`, generated Go code. No ORM.
- **D-10:** `dir="rtl"` + `lang="fa"` on `<html>`. Tailwind CSS v4 logical properties only (`ms-`, `me-`, `ps-`, `pe-`, `text-start`, `text-end`). No physical LTR properties.
- **D-11:** Vazirmatn font via npm `vazirmatn` package, variable font in global CSS, default font-family via Tailwind theme
- **D-12:** `useShamsiDate()` composable wrapping `jalaali-js`. Dates stored Gregorian in PostgreSQL, converted to Shamsi at display layer only.
- **D-13:** `toPersianDigits(value)` utility. Vue directive `v-fa-digits` or composable usage.
- **D-14:** Mobile-only viewport. No desktop breakpoints. Max width constrained.
- **D-15:** Three Nuxt layouts: `admin.vue`, `nutritionist.vue`, `client.vue`, plus `auth.vue`
- **D-16:** Bottom navigation bar for all roles
- **D-17:** Nuxt route middleware for role-based access control. `auth` checks JWT, `role` redirects unauthorized.
- **D-18:** Post-login redirects: Admin → `/admin`, Nutritionist → `/nutritionist/clients`, Client → `/client/plan`
- **D-19:** Docker Compose: `api` (Go multi-stage ~20MB), `frontend` (Nuxt Node), `postgres` (PG 16 named volume), `traefik` (v3 ACME HTTP challenge)
- **D-20:** Traefik via Docker labels — no separate config file. HTTPS redirect enforced.
- **D-21:** GitLab CI/CD: lint → test → build → deploy stages
- **D-22:** Health check at `GET /api/health` returns `{"status":"ok","timestamp":"..."}`
- **D-23:** zerolog structured JSON logging to stdout. Fields: timestamp, level, method, path, status, duration_ms, request_id.
- **D-24:** CORS restricted to frontend domain only. Credentials allowed for httpOnly cookies.
- **D-25:** Input validation via `go-playground/validator` struct tags. Custom validators for Iranian mobile (`^09[0-9]{9}$`), Persian text.
- **D-26:** All SQL via sqlc — parameterized by design (SEC-03).
- **D-27:** In-memory rate limiter on OTP endpoints: max 3 requests/phone/10min, max 3 verification attempts per code.

### Agent's Discretion
- Loading skeleton/spinner design for pages
- Exact Tailwind color palette and design tokens
- Error page styling (404, 500, network error)
- Exact CI/CD runner tags and Docker registry choice
- Go module path naming
- Exact health check response fields beyond status

### Deferred Ideas (OUT OF SCOPE)
- Persian pg_trgm search validation spike — Phase 2
- Plan builder UI state management complexity — Phase 3
- iOS PWA storage eviction testing — Phase 6
- Chart.js weight/measurement visualizations — Phase 4
- WebPush notification infrastructure — Phase 6
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| AUTH-01 | Super Admin login with email/password | Gin route handler + bcrypt verify + JWT cookie issue |
| AUTH-02 | Super Admin seeded via CLI command | `cmd/seed/main.go` using env vars; migration creates schema |
| AUTH-03 | Nutritionist login with email/password | Same handler as Admin with role differentiation |
| AUTH-04 | Nutritionist accounts created by Super Admin only | Admin handler + service with role guard middleware |
| AUTH-05 | Client OTP via SMS (Kavenegar) | SMSSender interface + Kavenegar REST adapter + dev stdout mock |
| AUTH-06 | OTP: 6 digits, 2-min validity, max 3 attempts | `otp_codes` table with hashed OTP, expiry, attempt counter |
| AUTH-07 | OTP rate limit: 3 req/phone/10 min | In-memory sliding window rate limiter middleware |
| AUTH-08 | JWT access (15min) + refresh (30 days) | golang-jwt/v5 with custom claims, dual httpOnly cookie pattern |
| AUTH-09 | JWT refresh handles concurrency | Frontend Pinia auth store with isRefreshing mutex + retry queue |
| AUTH-10 | Passwords hashed with bcrypt cost 12 | `golang.org/x/crypto/bcrypt` with cost factor 12 |
| AUTH-11 | Row-level auth: nutritionist only owns clients | Repository WHERE clauses join on nutritionist_id; return 404 not 403 |
| AUTH-12 | Client cannot self-register | No public client registration endpoint; nutritionist-only flow |
| CLNT-01 | Nutritionist registers client | Handler + service for creating user with role=client + nutritionist_id |
| UI-01 | Persian RTL layout with Tailwind v4 logical props | `@import "tailwindcss"` + `@theme` Vazirmatn + `dir="rtl"` |
| UI-02 | Mobile-only viewport | `<meta viewport>` + max-width constraint in CSS |
| UI-03 | Vazirmatn font | npm `vazirmatn` imported as variable font in global CSS |
| UI-04 | Shamsi/Jalali dates via jalaali-js | `useShamsiDate()` composable wrapping jalaali-js |
| UI-05 | Persian numeral display | `toPersianDigits()` utility replacing 0-9 → ۰-۹ |
| INFRA-01 | Docker + Docker Compose on Hetzner | Multi-service compose file with multi-stage Go build |
| INFRA-02 | Traefik with HTTPS Let's Encrypt | Traefik v3 with Docker labels, ACME HTTP challenge |
| INFRA-03 | GitLab CI/CD pipeline | 4-stage: lint → test → build → deploy |
| INFRA-04 | Structured JSON logging to stdout | zerolog middleware with request_id, duration_ms fields |
| INFRA-06 | Health check endpoint | `GET /api/health` returning status + timestamp |
| SEC-01 | All traffic HTTPS (TLS 1.2+) | Traefik enforces HTTPS redirect; `Secure` flag on cookies |
| SEC-02 | Input validation on all endpoints | go-playground/validator struct tags on all DTOs |
| SEC-03 | SQL injection prevention | sqlc generates parameterized queries by design |
| SEC-06 | CORS restricted to app domain | Gin CORS middleware with explicit origin, credentials=true |
| SEC-07 | OTP brute force protection | Rate limiter + attempt counter + same error message for all failures |
</phase_requirements>

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| JWT authentication (issue/verify/refresh) | API / Backend | — | Token signing, validation, cookie management are server concerns |
| OTP generation & verification | API / Backend | — | Server generates, stores hash, verifies; SMS sent server-side |
| Role-based route protection (backend) | API / Backend | — | Gin middleware guards route groups by JWT role claim |
| Role-based route protection (frontend) | Frontend Server (SSR) | Browser / Client | Nuxt route middleware runs on SSR and client navigation |
| Persian RTL layout & typography | Browser / Client | — | CSS dir="rtl", font loading, layout — purely client-side |
| Shamsi date display | Browser / Client | — | Conversion from Gregorian happens at display layer only |
| Persian numeral conversion | Browser / Client | — | UI concern; server stores/returns standard numerals |
| Input validation (server) | API / Backend | — | go-playground/validator on all request DTOs |
| Input validation (client) | Browser / Client | — | Form validation UX before submission |
| Client registration | API / Backend | Browser / Client | Backend creates user record; frontend provides form UI |
| Docker orchestration | CDN / Static | — | Traefik reverse proxy + Docker Compose service mesh |
| CI/CD pipeline | Infrastructure | — | GitLab CI/CD — lint, test, build, deploy stages |
| Structured logging | API / Backend | — | zerolog outputs JSON to stdout; collected by Loki |
| Rate limiting | API / Backend | — | In-memory rate limiter middleware on OTP endpoints |
| CORS | API / Backend | — | Gin CORS middleware restricts origins |
| Database schema & migrations | Database / Storage | — | PostgreSQL 16 with golang-migrate SQL files |

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go | 1.25+ (1.25.9 available) | Backend language | Required by Gin v1.12.0 [VERIFIED: go.dev/dl] |
| Gin | v1.12.0 | HTTP framework | Zero-alloc router, built-in validator integration, largest Go web ecosystem [VERIFIED: Go module proxy] |
| PostgreSQL | 16-alpine | Primary database | JSONB, pg_trgm, UUID, Docker support [CITED: PRD §8] |
| pgx/v5 | v5.9.2 | PostgreSQL driver | Pure Go, pgxpool, SendBatch for batch queries [VERIFIED: Go module proxy] |
| sqlc | v1.30.0 (CLI) | SQL → Go codegen | Type-safe queries from raw SQL, pgx/v5 integration [VERIFIED: sqlc version] |
| golang-migrate | v4.19.1 | DB migrations | Sequential SQL migration files, Go library + CLI [VERIFIED: Go module proxy] |
| golang-jwt | v5.3.1 | JWT tokens | Custom claims, HMAC-SHA256, industry standard [VERIFIED: Go module proxy] |
| go-playground/validator | v10.30.x | Input validation | Struct tags, built into Gin's ShouldBindJSON [VERIFIED: Go module proxy] |
| zerolog | v1.35.0 | Structured logging | Zero-alloc JSON, context-based propagation [VERIFIED: Go module proxy] |
| bcrypt | golang.org/x/crypto | Password hashing | Go extended stdlib, cost factor 12 [VERIFIED: Go module proxy] |
| google/uuid | v1.6.0 | UUID generation | v4 UUIDs for all PKs [VERIFIED: Go module proxy] |
| Nuxt | 4.4.2 | Frontend framework | Vue 3 SSR/SPA, file-based routing, app/ directory [VERIFIED: npm registry] |
| Tailwind CSS | 4.2.2 | CSS framework | Native logical properties for RTL, CSS-first config [VERIFIED: npm registry] |
| @tailwindcss/postcss | 4.2.2 | PostCSS integration | Direct Tailwind v4 integration for Nuxt [VERIFIED: npm registry] |
| Pinia | 3.0.4 | State management | Official Vue 3 state manager, SSR hydration [VERIFIED: npm registry] |
| Vazirmatn | 33.0.3 | Persian font | Standard Persian variable web font [VERIFIED: npm registry] |
| jalaali-js | 1.2.8 | Jalali calendar | Gregorian ↔ Shamsi conversion, ~3KB [VERIFIED: npm registry] |
| Traefik | v3.4 | Reverse proxy | Auto Let's Encrypt, Docker labels, HTTPS redirect [CITED: STACK.md] |
| Docker Compose | v5.0.2 | Orchestration | Multi-service deployment [VERIFIED: docker compose version] |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| @nuxt/eslint | 1.15.2 | JS/TS linting | CI lint stage, Nuxt-specific rules [VERIFIED: npm registry] |
| Vitest | 4.1.4 | Frontend tests | Unit tests for composables, stores [VERIFIED: npm registry] |
| golangci-lint | v1.64.8 | Go linting | CI lint stage, 50+ linters [VERIFIED: golangci-lint --version] |
| testify | latest | Go test assertions | Assertion helpers for Go tests [ASSUMED] |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Gin v1.12.0 | Gin v1.11.0 | Compatible with Go 1.24.2 but misses improvements; upgrade Go instead |
| sqlc | GORM | GORM adds runtime overhead, magic; sqlc gives full SQL control |
| golang-migrate | Atlas | Atlas is more powerful but adds complexity for this project scale |
| @tailwindcss/postcss | @nuxtjs/tailwindcss | Module may lag behind Tailwind v4; direct PostCSS is simpler |

**Installation:**

Backend (Go):
```bash
# Upgrade Go to 1.25+ first
go install golang.org/dl/go1.25.9@latest  # or download from go.dev/dl

cd backend
go mod init github.com/ranjbar-dev/nutritrack/backend
go get github.com/gin-gonic/gin@v1.12.0
go get github.com/jackc/pgx/v5@v5.9.2
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/golang-jwt/jwt/v5@v5.3.1
go get github.com/golang-migrate/migrate/v4@v4.19.1
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
go get github.com/go-playground/validator/v10
go get github.com/rs/zerolog@v1.35.0
go get github.com/google/uuid@v1.6.0
go get golang.org/x/crypto
```

Frontend (Nuxt 4):
```bash
cd frontend
npx nuxi@latest init . --template nuxt4
npm install pinia vazirmatn jalaali-js
npm install tailwindcss @tailwindcss/postcss
npm install -D @nuxt/eslint vitest @vue/test-utils
```

## Architecture Patterns

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                  Mobile Browser (PWA Shell)                   │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              Nuxt 4 SPA (app/ directory)               │  │
│  │                                                        │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐ │  │
│  │  │  Pinia   │  │Composable│  │     Layouts          │ │  │
│  │  │authStore │  │useAuth() │  │admin/nutri/client/   │ │  │
│  │  │          │  │useApi()  │  │auth                  │ │  │
│  │  └────┬─────┘  └────┬─────┘  └──────────────────────┘ │  │
│  │       │              │                                  │  │
│  │  ┌────▼──────────────▼──────┐                          │  │
│  │  │  ofetch / useFetch       │  Refresh mutex/queue     │  │
│  │  │  + 401 interceptor       │  pattern on 401          │  │
│  │  └──────────┬───────────────┘                          │  │
│  └─────────────┼──────────────────────────────────────────┘  │
└────────────────┼─────────────────────────────────────────────┘
                 │ HTTPS (REST JSON)
                 ▼
┌────────────────────────────────────────────────────────────────┐
│              Traefik v3 (Docker Container)                      │
│    TLS termination (Let's Encrypt ACME HTTP challenge)         │
│    HTTP → HTTPS redirect                                       │
│    Route: domain.com → frontend, domain.com/api/ → api         │
└────────────┬──────────────────────┬────────────────────────────┘
             │                      │
    ┌────────▼──────┐     ┌────────▼──────┐
    │  Nuxt SSR     │     │  Go API (Gin) │
    │  (frontend)   │     │  Port :8080    │
    │  Port :3000   │     │               │
    └───────────────┘     │  ┌──────────────────────────────┐
                          │  │    Middleware Chain           │
                          │  │  Recovery → Logger → CORS    │
                          │  │  → RateLimit → JWT Auth      │
                          │  │  → Role Guard                │
                          │  └──────────────────────────────┘
                          │               │
                          │  ┌────────────▼─────────────┐
                          │  │     Route Groups          │
                          │  │  /api/auth/*    (public)  │
                          │  │  /api/health    (public)  │
                          │  │  /api/admin/*   (admin)   │
                          │  │  /api/nutritionist/* (nut)│
                          │  │  /api/client/*  (client)  │
                          │  └────────────┬─────────────┘
                          │               │
                          │  ┌────────────▼─────────────┐
                          │  │   Handler Layer           │
                          │  │   Parse request, call svc │
                          │  └────────────┬─────────────┘
                          │               │
                          │  ┌────────────▼─────────────┐
                          │  │   Service Layer           │
                          │  │   Business logic, auth    │
                          │  │   checks, SMS adapter     │
                          │  └────────────┬─────────────┘
                          │               │
                          │  ┌────────────▼─────────────┐
                          │  │   Repository Layer        │
                          │  │   sqlc-generated queries  │
                          │  │   + custom compositions   │
                          │  └────────────┬─────────────┘
                          └───────────────┼───────────────┘
                                          │
                                   ┌──────▼──────┐
                                   │ PostgreSQL  │
                                   │ 16-alpine   │
                                   │             │
                                   │ users table │
                                   │ otp_codes   │
                                   │ refresh_    │
                                   │  tokens     │
                                   └─────────────┘
```

### Recommended Project Structure

```
NutriTrack/                          # Git root (monorepo)
├── backend/
│   ├── cmd/
│   │   ├── api/
│   │   │   └── main.go              # Entry point: config, DB, router, start
│   │   └── seed/
│   │       └── main.go              # Super Admin seeder CLI
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go            # Env-based config with validation
│   │   ├── handler/
│   │   │   ├── auth_handler.go      # Login, OTP, refresh endpoints
│   │   │   ├── admin_handler.go     # Nutritionist CRUD by admin
│   │   │   ├── client_handler.go    # Client registration by nutritionist
│   │   │   └── health_handler.go    # Health check endpoint
│   │   ├── service/
│   │   │   ├── auth_service.go      # JWT issue/refresh, OTP verify, bcrypt
│   │   │   └── user_service.go      # User CRUD, role management
│   │   ├── repository/
│   │   │   ├── sqlc/                # sqlc-generated code (DO NOT EDIT)
│   │   │   │   ├── db.go
│   │   │   │   ├── models.go
│   │   │   │   └── queries.sql.go
│   │   │   ├── user_repo.go         # Wraps sqlc + custom compositions
│   │   │   └── otp_repo.go          # OTP CRUD with expiry
│   │   ├── model/
│   │   │   ├── user.go              # Domain models + enums
│   │   │   └── dto/
│   │   │       └── auth_dto.go      # Request/response DTOs
│   │   ├── middleware/
│   │   │   ├── auth.go              # JWT extraction from httpOnly cookie
│   │   │   ├── role_guard.go        # Role-based access control
│   │   │   ├── rate_limit.go        # Sliding window in-memory limiter
│   │   │   ├── logger.go            # zerolog request logging middleware
│   │   │   ├── cors.go              # CORS configuration
│   │   │   └── request_id.go        # UUID request ID injection
│   │   └── config/
│   │       └── config.go            # Env loading + validation
│   ├── db/
│   │   ├── migrations/
│   │   │   ├── 000001_create_users.up.sql
│   │   │   ├── 000001_create_users.down.sql
│   │   │   ├── 000002_create_otp_codes.up.sql
│   │   │   ├── 000002_create_otp_codes.down.sql
│   │   │   ├── 000003_create_refresh_tokens.up.sql
│   │   │   └── 000003_create_refresh_tokens.down.sql
│   │   └── queries/
│   │       ├── users.sql             # sqlc: user CRUD queries
│   │       └── otp.sql               # sqlc: OTP CRUD queries
│   ├── pkg/
│   │   ├── sms/
│   │   │   ├── sms.go               # interface SMSSender
│   │   │   ├── mock.go              # Logs to stdout
│   │   │   └── kavenegar.go         # HTTP REST adapter
│   │   └── jwt/
│   │       └── jwt.go               # Token create/parse helpers
│   ├── sqlc.yaml                    # sqlc configuration
│   ├── Dockerfile                   # Multi-stage build
│   ├── .env.example
│   └── go.mod
├── frontend/
│   ├── app/
│   │   ├── assets/
│   │   │   └── css/
│   │   │       └── main.css         # Tailwind + Vazirmatn import
│   │   ├── components/
│   │   │   └── ui/
│   │   │       ├── AppButton.vue
│   │   │       ├── AppInput.vue
│   │   │       ├── BottomNav.vue
│   │   │       ├── LoadingSpinner.vue
│   │   │       └── EmptyState.vue
│   │   ├── composables/
│   │   │   ├── useAuth.ts           # Login, logout, token management
│   │   │   ├── useApi.ts            # ofetch wrapper + 401 refresh queue
│   │   │   ├── useShamsiDate.ts     # Jalali date formatting
│   │   │   └── usePersianDigits.ts  # Latin → Persian numeral conversion
│   │   ├── layouts/
│   │   │   ├── admin.vue            # Admin layout + bottom nav
│   │   │   ├── nutritionist.vue     # Nutritionist layout + bottom nav
│   │   │   ├── client.vue           # Client layout + bottom nav
│   │   │   └── auth.vue             # Login/OTP pages, no nav
│   │   ├── middleware/
│   │   │   ├── auth.global.ts       # Redirect unauthenticated to login
│   │   │   └── role-guard.ts        # Check JWT role matches route prefix
│   │   ├── pages/
│   │   │   ├── auth/
│   │   │   │   ├── login.vue        # Email/password for admin+nutritionist
│   │   │   │   └── otp.vue          # OTP flow for clients
│   │   │   ├── admin/
│   │   │   │   └── index.vue        # Admin dashboard (placeholder)
│   │   │   ├── nutritionist/
│   │   │   │   ├── clients.vue      # Client list + register form
│   │   │   │   └── index.vue        # Redirect to clients
│   │   │   └── client/
│   │   │       ├── plan.vue          # Plan view (placeholder)
│   │   │       └── index.vue         # Redirect to plan
│   │   ├── plugins/
│   │   │   └── api.ts               # Configure ofetch defaults
│   │   ├── stores/
│   │   │   └── auth.ts              # Auth Pinia store
│   │   ├── utils/
│   │   │   ├── persian-digits.ts    # toPersianDigits() function
│   │   │   └── constants.ts         # Role enums, route maps
│   │   ├── app.vue                  # Root: <html dir="rtl" lang="fa">
│   │   └── app.config.ts
│   ├── shared/
│   │   └── types/
│   │       ├── api.ts               # API response types
│   │       └── user.ts              # User, Role types
│   ├── nuxt.config.ts               # Tailwind, API proxy, CSS config
│   ├── postcss.config.mjs           # @tailwindcss/postcss
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml               # api + frontend + postgres + traefik
├── docker-compose.dev.yml           # Dev overrides (no TLS, exposed ports)
├── .gitlab-ci.yml                   # CI/CD pipeline definition
└── .gitignore
```

### Pattern 1: Layered Service Architecture (Go Backend)

**What:** Handler → Service → Repository, each layer depends only on the layer below via interfaces.
**When to use:** Every endpoint in the Go API.

```go
// Source: Context7 /gin-gonic/gin, ARCHITECTURE.md Pattern 1

// internal/repository/user_repo.go — Interface + implementation
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    GetByMobile(ctx context.Context, mobile string) (*model.User, error)
    GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type userRepo struct {
    q *sqlc.Queries  // sqlc-generated
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
    return &userRepo{q: sqlc.New(pool)}
}

// internal/service/auth_service.go — Business logic
type AuthService struct {
    userRepo  repository.UserRepository
    otpRepo   repository.OTPRepository
    smsSender sms.Sender
    jwtSecret []byte
}

func NewAuthService(ur repository.UserRepository, or repository.OTPRepository, s sms.Sender, secret []byte) *AuthService {
    return &AuthService{userRepo: ur, otpRepo: or, smsSender: s, jwtSecret: secret}
}

// internal/handler/auth_handler.go — HTTP concerns only
type AuthHandler struct {
    svc *service.AuthService
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "داده‌های ورودی نامعتبر است"})
        return
    }
    tokens, user, err := h.svc.LoginWithPassword(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        c.JSON(401, gin.H{"error": "ایمیل یا رمز عبور اشتباه است"})
        return
    }
    setAuthCookies(c, tokens)
    c.JSON(200, gin.H{"user": user})
}
```

### Pattern 2: JWT httpOnly Cookie Auth

**What:** Access and refresh tokens stored in separate httpOnly secure cookies, not localStorage.
**When to use:** All auth token management.

```go
// Source: Context7 /gin-gonic/gin cookie docs

func setAuthCookies(c *gin.Context, tokens *model.TokenPair) {
    // Access token cookie — 15 minutes
    c.SetCookie(
        "access_token",       // name
        tokens.AccessToken,   // value
        900,                  // maxAge (15 min)
        "/api",               // path — only sent to API
        "",                   // domain (auto from request)
        true,                 // secure (HTTPS only)
        true,                 // httpOnly
    )
    // Refresh token cookie — 30 days
    c.SetCookie(
        "refresh_token",
        tokens.RefreshToken,
        2592000,              // 30 days
        "/api/auth/refresh",  // path — only sent to refresh endpoint
        "",
        true,
        true,
    )
}

func clearAuthCookies(c *gin.Context) {
    c.SetCookie("access_token", "", -1, "/api", "", true, true)
    c.SetCookie("refresh_token", "", -1, "/api/auth/refresh", "", true, true)
}

// Middleware: extract JWT from cookie
func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr, err := c.Cookie("access_token")
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }
        claims, err := jwt.ParseToken(tokenStr, jwtSecret)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

### Pattern 3: Frontend JWT Refresh Queue (AUTH-09)

**What:** When multiple concurrent requests detect 401, only one triggers refresh; others wait on a promise.
**When to use:** The `useApi()` composable's response interceptor.

```typescript
// Source: PITFALLS.md Pitfall 7, CONTEXT.md D-02

// app/composables/useApi.ts
let isRefreshing = false
let refreshPromise: Promise<void> | null = null

export function useApi() {
  const config = useRuntimeConfig()

  async function refreshTokens(): Promise<void> {
    await $fetch('/api/auth/refresh', {
      method: 'POST',
      baseURL: config.public.apiBase,
      credentials: 'include',
    })
  }

  async function request<T>(url: string, options: any = {}): Promise<T> {
    try {
      return await $fetch<T>(url, {
        baseURL: config.public.apiBase,
        credentials: 'include', // send httpOnly cookies
        ...options,
      })
    } catch (err: any) {
      if (err?.statusCode === 401 && !options._retry) {
        // Mutex: only first request triggers refresh
        if (!isRefreshing) {
          isRefreshing = true
          refreshPromise = refreshTokens()
            .finally(() => {
              isRefreshing = false
              refreshPromise = null
            })
        }
        // All requests wait for the same refresh
        await refreshPromise
        // Retry original request
        return request<T>(url, { ...options, _retry: true })
      }
      throw err
    }
  }

  return {
    get: <T>(url: string, opts?: any) => request<T>(url, { method: 'GET', ...opts }),
    post: <T>(url: string, body: any, opts?: any) => request<T>(url, { method: 'POST', body, ...opts }),
    put: <T>(url: string, body: any, opts?: any) => request<T>(url, { method: 'PUT', body, ...opts }),
    del: <T>(url: string, opts?: any) => request<T>(url, { method: 'DELETE', ...opts }),
  }
}
```

### Pattern 4: Gin Route Groups with Middleware Stacking

**What:** Route groups apply auth middleware at the group level; four API zones.
**When to use:** Router setup in `cmd/api/main.go`.

```go
// Source: Context7 /gin-gonic/gin route groups

func SetupRouter(deps *Dependencies) *gin.Engine {
    r := gin.New()
    r.Use(middleware.Recovery())
    r.Use(middleware.RequestID())
    r.Use(middleware.Logger(deps.Logger))
    r.Use(middleware.CORS(deps.Config.FrontendURL))

    // Public routes — no auth
    pub := r.Group("/api")
    {
        pub.GET("/health", deps.Handlers.Health)
        pub.POST("/auth/login", deps.Handlers.Auth.Login)
        pub.POST("/auth/otp/request", middleware.RateLimit(deps.RateLimiter), deps.Handlers.Auth.RequestOTP)
        pub.POST("/auth/otp/verify", middleware.RateLimit(deps.RateLimiter), deps.Handlers.Auth.VerifyOTP)
        pub.POST("/auth/refresh", deps.Handlers.Auth.Refresh)
    }

    // Admin routes
    admin := r.Group("/api/admin")
    admin.Use(middleware.Auth(deps.JWTSecret), middleware.RoleGuard("super_admin"))
    {
        admin.POST("/nutritionists", deps.Handlers.Admin.CreateNutritionist)
        // ... more admin routes in Phase 2+
    }

    // Nutritionist routes
    nutri := r.Group("/api/nutritionist")
    nutri.Use(middleware.Auth(deps.JWTSecret), middleware.RoleGuard("nutritionist"))
    {
        nutri.POST("/clients", deps.Handlers.Client.Register)
        // ... more nutritionist routes in Phase 2+
    }

    // Client routes
    client := r.Group("/api/client")
    client.Use(middleware.Auth(deps.JWTSecret), middleware.RoleGuard("client"))
    {
        // placeholder for Phase 2+ routes
    }

    return r
}
```

### Pattern 5: sqlc Query Definition + Code Generation

**What:** Write SQL in `backend/db/queries/`, generate Go code via `sqlc generate`.
**When to use:** All database access.

```sql
-- Source: Context7 /websites/sqlc_dev_en

-- backend/db/queries/users.sql

-- name: GetUserByEmail :one
SELECT id, role, full_name, email, password_hash, mobile,
       date_of_birth, height_cm, gender, nutritionist_id,
       is_active, created_at, updated_at
FROM users
WHERE email = $1 AND is_active = true;

-- name: GetUserByMobile :one
SELECT id, role, full_name, email, password_hash, mobile,
       date_of_birth, height_cm, gender, nutritionist_id,
       is_active, created_at, updated_at
FROM users
WHERE mobile = $1 AND is_active = true;

-- name: CreateUser :one
INSERT INTO users (id, role, full_name, email, password_hash, mobile,
                   date_of_birth, height_cm, gender, nutritionist_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetClientsByNutritionistID :many
SELECT id, full_name, mobile, is_active, created_at
FROM users
WHERE nutritionist_id = $1 AND role = 'client'
ORDER BY created_at DESC;
```

```yaml
# backend/sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/queries/"
    schema: "db/migrations/"
    gen:
      go:
        package: "sqlc"
        out: "internal/repository/sqlc"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_empty_slices: true
```

### Pattern 6: Nuxt 4 Layouts with Bottom Navigation

**What:** Role-specific layouts with bottom nav bar, layout selected via `definePageMeta`.
**When to use:** Every page file.

```vue
<!-- Source: Context7 /websites/nuxt_4_x layout docs -->

<!-- app/layouts/nutritionist.vue -->
<template>
  <div class="min-h-screen bg-gray-50 pb-16" dir="rtl">
    <slot />
    <BottomNav :items="navItems" />
  </div>
</template>

<script setup lang="ts">
const navItems = [
  { label: 'کلاینت‌ها', icon: 'users', to: '/nutritionist/clients' },
  { label: 'غذاها', icon: 'food', to: '/nutritionist/foods' },
  { label: 'پیام‌ها', icon: 'chat', to: '/nutritionist/messages' },
  { label: 'پروفایل', icon: 'user', to: '/nutritionist/profile' },
]
</script>

<!-- app/pages/nutritionist/clients.vue -->
<script setup lang="ts">
definePageMeta({
  layout: 'nutritionist',
  middleware: ['role-guard'],
})
</script>
```

### Anti-Patterns to Avoid

- **Physical CSS properties in RTL:** NEVER use `ml-`, `mr-`, `pl-`, `pr-`, `text-left`, `text-right`. ALWAYS use `ms-`, `me-`, `ps-`, `pe-`, `text-start`, `text-end`. [CITED: D-10]
- **Business logic in handlers:** Handlers ONLY parse/validate/serialize. No SQL, no business rules, no authorization checks inline. [CITED: ARCHITECTURE.md Anti-Pattern 2]
- **JWT secret in code:** Load from environment variable; minimum 256-bit random key. [CITED: PITFALLS.md Security Mistakes]
- **Different error messages for auth failures:** Return same error for "user not found" and "wrong password/OTP" to prevent enumeration. [CITED: PITFALLS.md]
- **Storing OTP in plaintext:** Hash OTP with SHA-256 before storing in database. [CITED: PITFALLS.md]
- **Using Nuxt 3 directory structure:** Must use `app/` directory for all frontend code (pages, components, composables). [CITED: PITFALLS.md Pitfall 14]

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| JWT creation/validation | Custom token signing | `golang-jwt/jwt/v5` | Edge cases in claim validation, algorithm confusion attacks |
| Password hashing | Custom hash function | `golang.org/x/crypto/bcrypt` cost 12 | bcrypt handles salt, cost factor, timing-safe comparison |
| Input validation | Manual field checks | `go-playground/validator/v10` struct tags | Integrated with Gin's ShouldBindJSON; handles nested, custom validators |
| SQL queries | String concatenation | sqlc generated parameterized queries | Prevents SQL injection by design; compile-time type safety |
| UUID generation | Custom ID schemes | `google/uuid` v4 | Cryptographically random, collision-safe, standard format |
| CSS RTL layout | Manual direction flipping | Tailwind v4 logical properties | `ms-`, `me-` auto-flip based on `dir="rtl"`; no plugin needed |
| Jalali dates | Manual calendar math | `jalaali-js` composable wrapper | Leap year, month length, boundary rules are complex |
| Persian numerals | Regex per-component | Global utility + Vue directive | One function `toPersianDigits()` handles all cases |
| DB migrations | Manual schema changes | `golang-migrate/migrate/v4` | Version-controlled, reversible, CLI + Go library |
| Reverse proxy + TLS | Nginx + certbot cron | Traefik v3 Docker labels | Auto Let's Encrypt, auto-discovery, zero config files |

**Key insight:** Phase 1 has zero novel algorithms. Every capability maps to a well-tested library. The value is in correct integration and consistent patterns, not custom solutions.

## Common Pitfalls

### Pitfall 1: Go Version Incompatibility

**What goes wrong:** Gin v1.12.0 requires Go 1.25 (`go 1.25.0` in go.mod). The local environment has Go 1.24.2. Running `go get github.com/gin-gonic/gin@v1.12.0` will fail or produce a broken build.
**Why it happens:** Go 1.25 was released recently; developers often have stale Go installations.
**How to avoid:** Upgrade Go to 1.25+ before any backend initialization. Go 1.25.9 is current stable. [VERIFIED: go.dev/dl and Go module proxy]
**Warning signs:** `go build` fails with "requires go >= 1.25"

### Pitfall 2: JWT Refresh Token Race Condition (AUTH-09)

**What goes wrong:** Multiple concurrent API requests discover expired access token simultaneously, all try to refresh, only first succeeds if rotation is used → mass logout.
**Why it happens:** Nuxt makes parallel API calls on page load (e.g., load user + load data).
**How to avoid:** Implement mutex/queue pattern in `useApi()` composable — first 401 triggers refresh, all others wait on the promise. [CITED: PITFALLS.md Pitfall 7, CONTEXT.md D-02]
**Warning signs:** Users report "random logouts" after 15+ minutes of inactivity

### Pitfall 3: OTP Rate Limit Bypass via Phone Format

**What goes wrong:** Iranian numbers formatted as `09123456789`, `9123456789`, `+989123456789`, `00989123456789` all bypass rate limiter keyed on raw input.
**Why it happens:** Rate limiter uses raw input as key; attacker sends different formats.
**How to avoid:** Normalize phone to canonical 10-digit format (`9XXXXXXXXX`) before rate limiting, OTP generation, and user lookup. [CITED: PITFALLS.md Pitfall 16]
**Warning signs:** Rate limiter tests pass but manual testing with different formats bypasses limits

### Pitfall 4: Cookie SameSite and CORS Misconfiguration

**What goes wrong:** httpOnly cookies not sent cross-origin if SameSite and CORS aren't properly configured. Frontend at `domain.com` can't reach API at `api.domain.com` or `domain.com/api/` without correct setup.
**Why it happens:** Browser security policies require `credentials: 'include'` on fetch AND `Access-Control-Allow-Credentials: true` + explicit origin (not `*`) in CORS. SameSite must be `Lax` or `None` for cross-origin.
**How to avoid:** If API is on same domain (recommended via Traefik path routing), use `SameSite=Strict`. Configure Gin CORS middleware with explicit frontend URL and `AllowCredentials: true`. [ASSUMED]
**Warning signs:** Login succeeds but subsequent requests fail with 401; cookies not sent in browser DevTools

### Pitfall 5: Nuxt 4 Directory Structure Mismatch

**What goes wrong:** Using Nuxt 3 directory structure (components at root) instead of `app/` directory causes auto-imports to fail, pages not found, middleware not running.
**Why it happens:** Most tutorials/community examples still use Nuxt 3 structure.
**How to avoid:** Use `app/` directory from day one per Nuxt 4 convention. Don't set `srcDir: '.'`. [CITED: PITFALLS.md Pitfall 14]
**Warning signs:** `nuxi dev` shows "component not found" warnings; composables return undefined

### Pitfall 6: Timestamp Without Timezone in PostgreSQL

**What goes wrong:** Using `timestamp` instead of `timestamptz` causes ambiguous times that break when server timezone changes or is different from client timezone (Tehran UTC+3:30).
**Why it happens:** `timestamp` is the shorter type name; developers default to it.
**How to avoid:** Always use `timestamptz` for all timestamp columns. Set `time.LoadLocation("Asia/Tehran")` in Go backend for "today" calculations. [CITED: PITFALLS.md Technical Debt Patterns]
**Warning signs:** Dates off by one day for events near midnight

## Code Examples

### Database Schema — Users Table (Phase 1 Migration)

```sql
-- Source: PRD §9 Data Model + ARCHITECTURE.md indexes

-- backend/db/migrations/000001_create_users.up.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM ('super_admin', 'nutritionist', 'client');
CREATE TYPE gender_type AS ENUM ('male', 'female');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role user_role NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    mobile VARCHAR(15) UNIQUE,
    date_of_birth DATE,
    height_cm REAL,
    gender gender_type,
    nutritionist_id UUID REFERENCES users(id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for OTP lookup by mobile
CREATE INDEX idx_users_mobile ON users (mobile) WHERE mobile IS NOT NULL;
-- Index for nutritionist's client list
CREATE INDEX idx_users_nutritionist ON users (nutritionist_id) WHERE role = 'client';
-- Index for admin login
CREATE INDEX idx_users_email ON users (email) WHERE email IS NOT NULL;
```

### OTP Storage Table

```sql
-- backend/db/migrations/000002_create_otp_codes.up.sql
CREATE TABLE otp_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mobile VARCHAR(15) NOT NULL,
    code_hash VARCHAR(255) NOT NULL,  -- SHA-256 hashed OTP
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,
    expires_at TIMESTAMPTZ NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_otp_mobile_active ON otp_codes (mobile, expires_at)
    WHERE verified = false;
```

### Refresh Token Storage

```sql
-- backend/db/migrations/000003_create_refresh_tokens.up.sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    family_id UUID NOT NULL,  -- for rotation theft detection
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens (user_id) WHERE revoked = false;
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens (token_hash) WHERE revoked = false;
```

### SMS Sender Interface (Kavenegar Abstraction)

```go
// Source: CONTEXT.md D-04

// backend/pkg/sms/sms.go
package sms

type Sender interface {
    SendOTP(phone string, code string) error
}

// backend/pkg/sms/mock.go
type MockSender struct {
    logger zerolog.Logger
}

func (m *MockSender) SendOTP(phone, code string) error {
    m.logger.Info().Str("phone", phone).Str("code", code).Msg("OTP sent (mock)")
    return nil
}

// backend/pkg/sms/kavenegar.go
type KavenegarSender struct {
    apiKey   string
    template string
    client   *http.Client
}

func (k *KavenegarSender) SendOTP(phone, code string) error {
    url := fmt.Sprintf(
        "https://api.kavenegar.com/v1/%s/verify/lookup.json?receptor=%s&token=%s&template=%s",
        k.apiKey, phone, code, k.template,
    )
    resp, err := k.client.Get(url)
    if err != nil {
        return fmt.Errorf("kavenegar request failed: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != 200 {
        return fmt.Errorf("kavenegar returned %d", resp.StatusCode)
    }
    return nil
}
```

### Tailwind v4 + Vazirmatn RTL Setup

```css
/* Source: Context7 /websites/tailwindcss, STACK.md Key Configuration Patterns */

/* frontend/app/assets/css/main.css */
@import "tailwindcss";

@theme {
  --font-family-sans: 'Vazirmatn', system-ui, -apple-system, sans-serif;
}
```

```javascript
// frontend/postcss.config.mjs
const config = {
  plugins: {
    "@tailwindcss/postcss": {},
  },
};
export default config;
```

```typescript
// frontend/nuxt.config.ts
export default defineNuxtConfig({
  compatibilityDate: '2025-07-18',
  future: { compatibilityVersion: 4 },
  css: ['~/assets/css/main.css', 'vazirmatn/vazirmatn-font-face.css'],
  app: {
    head: {
      htmlAttrs: { dir: 'rtl', lang: 'fa' },
      meta: [
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      ],
    },
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
    },
  },
  modules: ['@pinia/nuxt', '@nuxt/eslint'],
})
```

### useShamsiDate Composable

```typescript
// Source: CONTEXT.md D-12, jalaali-js docs

// frontend/app/composables/useShamsiDate.ts
import { toJalaali, isValidJalaaliDate } from 'jalaali-js'

const persianMonths = [
  'فروردین', 'اردیبهشت', 'خرداد', 'تیر', 'مرداد', 'شهریور',
  'مهر', 'آبان', 'آذر', 'دی', 'بهمن', 'اسفند'
]

export function useShamsiDate() {
  function toShamsi(date: Date | string): { jy: number; jm: number; jd: number } {
    const d = new Date(date)
    return toJalaali(d.getFullYear(), d.getMonth() + 1, d.getDate())
  }

  function formatShamsi(date: Date | string, format: 'short' | 'long' = 'short'): string {
    const { jy, jm, jd } = toShamsi(date)
    const pd = toPersianDigits
    if (format === 'long') {
      return `${pd(jd)} ${persianMonths[jm - 1]} ${pd(jy)}`
    }
    return `${pd(jy)}/${pd(jm.toString().padStart(2, '0'))}/${pd(jd.toString().padStart(2, '0'))}`
  }

  function todayTehran(): Date {
    // Ensure "today" is calculated in Tehran timezone (UTC+3:30)
    return new Date(new Date().toLocaleString('en-US', { timeZone: 'Asia/Tehran' }))
  }

  return { toShamsi, formatShamsi, todayTehran }
}
```

### toPersianDigits Utility

```typescript
// Source: CONTEXT.md D-13

// frontend/app/utils/persian-digits.ts
const persianDigits = ['۰', '۱', '۲', '۳', '۴', '۵', '۶', '۷', '۸', '۹']

export function toPersianDigits(value: string | number): string {
  return String(value).replace(/[0-9]/g, (d) => persianDigits[parseInt(d)])
}
```

### Docker Compose Configuration

```yaml
# Source: CONTEXT.md D-19, D-20; Context7 Traefik docs

# docker-compose.yml
services:
  traefik:
    image: traefik:v3.4
    command:
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entryPoints.web.address=:80"
      - "--entryPoints.websecure.address=:443"
      - "--entryPoints.web.http.redirections.entryPoint.to=websecure"
      - "--entryPoints.web.http.redirections.entryPoint.scheme=https"
      - "--certificatesresolvers.le.acme.email=${ACME_EMAIL}"
      - "--certificatesresolvers.le.acme.storage=/letsencrypt/acme.json"
      - "--certificatesresolvers.le.acme.httpchallenge.entrypoint=web"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - letsencrypt:/letsencrypt
    restart: unless-stopped

  api:
    build:
      context: ./backend
      dockerfile: Dockerfile
    environment:
      - DATABASE_URL=postgres://nutritrack:${DB_PASSWORD}@postgres:5432/nutritrack?sslmode=disable
      - JWT_SECRET=${JWT_SECRET}
      - SMS_API_KEY=${SMS_API_KEY}
      - SMS_TEMPLATE=${SMS_TEMPLATE}
      - FRONTEND_URL=${FRONTEND_URL}
      - ENVIRONMENT=${ENVIRONMENT:-production}
    depends_on:
      postgres:
        condition: service_healthy
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(`${DOMAIN}`) && PathPrefix(`/api`)"
      - "traefik.http.routers.api.entrypoints=websecure"
      - "traefik.http.routers.api.tls.certresolver=le"
      - "traefik.http.services.api.loadbalancer.server.port=8080"
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 5s
      retries: 3
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    environment:
      - NUXT_PUBLIC_API_BASE=https://${DOMAIN}
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.frontend.rule=Host(`${DOMAIN}`)"
      - "traefik.http.routers.frontend.entrypoints=websecure"
      - "traefik.http.routers.frontend.tls.certresolver=le"
      - "traefik.http.services.frontend.loadbalancer.server.port=3000"
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_DB=nutritrack
      - POSTGRES_USER=nutritrack
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_INITDB_ARGS=--locale=en_US.UTF-8 --encoding=UTF8
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U nutritrack"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

volumes:
  pgdata:
  letsencrypt:
```

### GitLab CI/CD Pipeline

```yaml
# Source: CONTEXT.md D-21

# .gitlab-ci.yml
stages:
  - lint
  - test
  - build
  - deploy

variables:
  GO_VERSION: "1.25"
  NODE_VERSION: "22"

# --- LINT ---
lint:go:
  stage: lint
  image: golangci/golangci-lint:v1.64
  script:
    - cd backend
    - golangci-lint run ./...

lint:frontend:
  stage: lint
  image: node:${NODE_VERSION}-alpine
  script:
    - cd frontend
    - npm ci
    - npx eslint .

# --- TEST ---
test:go:
  stage: test
  image: golang:${GO_VERSION}
  services:
    - postgres:16-alpine
  variables:
    POSTGRES_DB: nutritrack_test
    POSTGRES_USER: test
    POSTGRES_PASSWORD: test
    DATABASE_URL: "postgres://test:test@postgres:5432/nutritrack_test?sslmode=disable"
  script:
    - cd backend
    - go test ./... -v -count=1 -race

test:frontend:
  stage: test
  image: node:${NODE_VERSION}-alpine
  script:
    - cd frontend
    - npm ci
    - npx vitest run

# --- BUILD ---
build:api:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t ${CI_REGISTRY_IMAGE}/api:${CI_COMMIT_SHORT_SHA} ./backend
    - docker push ${CI_REGISTRY_IMAGE}/api:${CI_COMMIT_SHORT_SHA}
  only:
    - main

build:frontend:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t ${CI_REGISTRY_IMAGE}/frontend:${CI_COMMIT_SHORT_SHA} ./frontend
    - docker push ${CI_REGISTRY_IMAGE}/frontend:${CI_COMMIT_SHORT_SHA}
  only:
    - main

# --- DEPLOY ---
deploy:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache openssh-client
    - eval $(ssh-agent -s)
    - echo "$SSH_PRIVATE_KEY" | ssh-add -
  script:
    - ssh -o StrictHostKeyChecking=no ${DEPLOY_USER}@${DEPLOY_HOST} "
        cd /opt/nutritrack &&
        docker compose pull &&
        docker compose up -d --remove-orphans"
  only:
    - main
  environment:
    name: production
```

### In-Memory Rate Limiter

```go
// Source: CONTEXT.md D-27

// backend/internal/middleware/rate_limit.go
package middleware

import (
    "net/http"
    "sync"
    "time"
    "github.com/gin-gonic/gin"
)

type entry struct {
    count    int
    windowStart time.Time
}

type RateLimiter struct {
    mu      sync.Mutex
    entries map[string]*entry
    max     int
    window  time.Duration
}

func NewRateLimiter(max int, window time.Duration) *RateLimiter {
    rl := &RateLimiter{
        entries: make(map[string]*entry),
        max:     max,
        window:  window,
    }
    // Cleanup goroutine
    go func() {
        for range time.Tick(window) {
            rl.mu.Lock()
            now := time.Now()
            for k, e := range rl.entries {
                if now.Sub(e.windowStart) > window {
                    delete(rl.entries, k)
                }
            }
            rl.mu.Unlock()
        }
    }()
    return rl
}

func (rl *RateLimiter) Allow(key string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    e, exists := rl.entries[key]
    if !exists || now.Sub(e.windowStart) > rl.window {
        rl.entries[key] = &entry{count: 1, windowStart: now}
        return true
    }
    if e.count >= rl.max {
        return false
    }
    e.count++
    return true
}

func RateLimit(limiter *RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Key: normalized phone number from request body
        key := c.GetString("rate_limit_key")
        if key == "" {
            key = c.ClientIP()
        }
        if !limiter.Allow(key) {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "تعداد درخواست‌ها بیش از حد مجاز است. لطفاً چند دقیقه صبر کنید.",
            })
            return
        }
        c.Next()
    }
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Tailwind CSS v3 config file | Tailwind v4 CSS-first `@import "tailwindcss"` + `@theme` | Jan 2025 | No `tailwind.config.js` needed; logical properties built-in |
| Nuxt 3 root directory | Nuxt 4 `app/` directory structure | Dec 2024 | `~` alias resolves to `app/`; auto-imports from `app/` |
| `@nuxtjs/tailwindcss` module | `@tailwindcss/postcss` direct integration | 2025 | Module may lag Tailwind v4; PostCSS is simpler |
| Gin v1.10.x (Go 1.20+) | Gin v1.12.0 (Go 1.25+) | Jun 2025 | Requires Go 1.25; adds QUIC support, improved validation |
| `tailwindcss-rtl` plugin | Native logical properties | Tailwind v4 | No plugin needed; `ms-`, `me-`, `ps-`, `pe-` built-in |
| `docker-compose` (v1 standalone) | `docker compose` (v2 CLI plugin) | 2023 | Use `docker compose` not `docker-compose` |

**Deprecated/outdated:**
- `tailwindcss-rtl` plugin — unnecessary with Tailwind v4 native logical properties
- `@nuxtjs/tailwindcss` — may not support Tailwind v4 properly; use `@tailwindcss/postcss`
- Nuxt 3 directory structure (`pages/` at root) — Nuxt 4 uses `app/pages/`
- `docker-compose` v1 standalone binary — Docker Compose v2 is the CLI plugin

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Cookie SameSite=Strict works when API and frontend share domain via Traefik path routing | Pitfall 4 | If using subdomains (api.domain.com), need SameSite=Lax or None; CORS config changes |
| A2 | testify is the standard Go test assertion library | Standard Stack | Low risk — can use stdlib testing alone |
| A3 | Kavenegar REST API uses `/verify/lookup.json` URL pattern for template-based OTP | Code Examples | Wrong URL pattern means OTP delivery fails; verify with Kavenegar docs during implementation |
| A4 | Vazirmatn npm package exports `vazirmatn-font-face.css` for import | Code Examples | Might be different path; check actual npm package contents |

## Open Questions

1. **Go Module Path**
   - What we know: STACK.md suggests `github.com/ranjbar-dev/nutritrack`; D-06 says monorepo with `backend/` subdir
   - What's unclear: Should the Go module be `github.com/ranjbar-dev/nutritrack/backend` or `github.com/ranjbar-dev/nutritrack`?
   - Recommendation: Use `github.com/ranjbar-dev/nutritrack/backend` since it's a subdirectory. Agent's discretion per CONTEXT.md.

2. **Refresh Token Rotation Strategy**
   - What we know: D-01 specifies refresh token in httpOnly cookie; PITFALLS.md recommends grace period for rotation
   - What's unclear: Full rotation (invalidate on use) vs. stateless (long-lived, no rotation)?
   - Recommendation: Use stateful rotation with `refresh_tokens` table + 30-second grace period for concurrent requests. Safer but adds DB query on every refresh.

3. **Kavenegar API Details**
   - What we know: REST API for Iranian SMS OTP delivery
   - What's unclear: Exact API URL pattern, template format, error codes
   - Recommendation: Use interface abstraction (D-04); implement mock first, real adapter later. Verify Kavenegar docs during implementation.

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | Backend | ✓ (wrong version) | 1.24.2 | Upgrade to 1.25.9 (REQUIRED for Gin v1.12.0) |
| Node.js | Frontend | ✓ | v24.14.1 | — |
| npm | Frontend | ✓ | 11.11.0 | — |
| Docker | Deployment | ✓ | 29.2.0 | — |
| Docker Compose | Deployment | ✓ | v5.0.2 | — |
| Git | Version control | ✓ | 2.53.0 | — |
| sqlc | SQL codegen | ✓ | v1.30.0 | — |
| golangci-lint | CI linting | ✓ | v1.64.8 | — |
| PostgreSQL (local) | Dev testing | ✗ | — | Use Docker: `postgres:16-alpine` |
| golang-migrate CLI | Migrations | ✗ | — | Use Go library integration (run migrations in `cmd/api/main.go` startup) |

**Missing dependencies with no fallback:**
- Go 1.25+ — **MUST upgrade** from 1.24.2 to 1.25.9 for Gin v1.12.0 compatibility

**Missing dependencies with fallback:**
- PostgreSQL not installed locally — run in Docker container (standard approach)
- golang-migrate CLI not installed — run migrations programmatically in Go startup code

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework (Go) | Go `testing` + testify v1.11+ |
| Framework (Frontend) | Vitest 4.1.4 + @vue/test-utils |
| Config file (Go) | None needed (Go convention) |
| Config file (Frontend) | `vitest.config.ts` (auto-detected by Nuxt) |
| Quick run command (Go) | `cd backend && go test ./internal/... -v -count=1` |
| Quick run command (Frontend) | `cd frontend && npx vitest run` |
| Full suite command | `cd backend && go test ./... -v -race && cd ../frontend && npx vitest run` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| AUTH-01 | Admin login with email/password | integration | `go test ./internal/handler/ -run TestAdminLogin -v` | ❌ Wave 0 |
| AUTH-02 | Admin seeded via CLI | integration | `go test ./cmd/seed/ -run TestSeedAdmin -v` | ❌ Wave 0 |
| AUTH-03 | Nutritionist login | integration | `go test ./internal/handler/ -run TestNutritionistLogin -v` | ❌ Wave 0 |
| AUTH-04 | Admin creates nutritionist | integration | `go test ./internal/handler/ -run TestCreateNutritionist -v` | ❌ Wave 0 |
| AUTH-05 | Client OTP flow | integration | `go test ./internal/handler/ -run TestClientOTP -v` | ❌ Wave 0 |
| AUTH-06 | OTP 6 digits, 2min, 3 attempts | unit | `go test ./internal/service/ -run TestOTPValidation -v` | ❌ Wave 0 |
| AUTH-07 | Rate limit 3 req/10min | unit | `go test ./internal/middleware/ -run TestRateLimit -v` | ❌ Wave 0 |
| AUTH-08 | JWT access + refresh tokens | unit | `go test ./pkg/jwt/ -run TestTokenPair -v` | ❌ Wave 0 |
| AUTH-09 | Concurrent refresh no mass logout | unit (frontend) | `npx vitest run --filter useApi` | ❌ Wave 0 |
| AUTH-10 | bcrypt cost 12 | unit | `go test ./internal/service/ -run TestPasswordHash -v` | ❌ Wave 0 |
| AUTH-11 | Row-level auth | integration | `go test ./internal/handler/ -run TestCrossNutritionistAccess -v` | ❌ Wave 0 |
| AUTH-12 | No client self-registration | integration | `go test ./internal/handler/ -run TestNoClientSelfReg -v` | ❌ Wave 0 |
| CLNT-01 | Nutritionist registers client | integration | `go test ./internal/handler/ -run TestRegisterClient -v` | ❌ Wave 0 |
| UI-01–05 | RTL, font, dates, numerals | manual-only | Visual inspection on mobile viewport | N/A |
| INFRA-01 | Docker Compose works | smoke | `docker compose up -d && curl http://localhost:8080/api/health` | ❌ Wave 0 |
| INFRA-06 | Health check endpoint | integration | `go test ./internal/handler/ -run TestHealthCheck -v` | ❌ Wave 0 |
| SEC-02 | Input validation | unit | `go test ./internal/handler/ -run TestInputValidation -v` | ❌ Wave 0 |
| SEC-07 | OTP brute force protection | integration | `go test ./internal/handler/ -run TestOTPBruteForce -v` | ❌ Wave 0 |

### Sampling Rate
- **Per task commit:** `cd backend && go test ./internal/... -v -count=1 -short`
- **Per wave merge:** Full Go test suite + Vitest
- **Phase gate:** All tests green, Docker compose health check passing

### Wave 0 Gaps
- [ ] `backend/internal/handler/auth_handler_test.go` — covers AUTH-01 to AUTH-12
- [ ] `backend/internal/middleware/rate_limit_test.go` — covers AUTH-07, SEC-07
- [ ] `backend/internal/service/auth_service_test.go` — covers AUTH-06, AUTH-08, AUTH-10
- [ ] `backend/pkg/jwt/jwt_test.go` — covers AUTH-08
- [ ] `frontend/app/composables/useApi.test.ts` — covers AUTH-09 (refresh queue)
- [ ] `frontend/app/utils/persian-digits.test.ts` — covers UI-05
- [ ] `frontend/app/composables/useShamsiDate.test.ts` — covers UI-04
- [ ] Framework install: `npm install -D vitest @vue/test-utils` + testify for Go

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | yes | bcrypt cost 12, OTP SHA-256 hashed, 2-min expiry, 3 attempts |
| V3 Session Management | yes | httpOnly Secure cookies, SameSite, 15-min access + 30-day refresh |
| V4 Access Control | yes | Gin route group middleware, role-based guards, repository-level ownership checks |
| V5 Input Validation | yes | go-playground/validator v10 struct tags on all DTOs; custom Iranian mobile validator |
| V6 Cryptography | no (Phase 1) | N/A — no encryption beyond bcrypt/JWT signing |

### Known Threat Patterns for Go/Gin + PostgreSQL + Nuxt

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| SQL injection | Tampering | sqlc parameterized queries (by design) |
| XSS via token theft | Information Disclosure | httpOnly cookies (not localStorage); CSP headers |
| OTP brute force | Elevation of Privilege | Rate limiter (3/10min) + attempt limiter (3 per code) + same error message |
| Phone number enumeration | Information Disclosure | Same error for "not registered" and "wrong OTP" |
| JWT secret compromise | Spoofing | Load from env var; minimum 256-bit random key |
| CORS misconfiguration | Tampering | Explicit origin, credentials=true, no wildcards |
| CSRF on cookie-auth endpoints | Tampering | SameSite=Strict on cookies; validate Origin header on mutations |
| Refresh token theft | Elevation of Privilege | Stateful rotation with family_id theft detection |

## Sources

### Primary (HIGH confidence)
- [VERIFIED: Go module proxy] — Gin v1.12.0 (Go 1.25 requirement), pgx v5.9.2, golang-jwt v5.3.1, golang-migrate v4.19.1, validator v10.30.x, zerolog v1.35.0, uuid v1.6.0
- [VERIFIED: npm registry] — nuxt 4.4.2, tailwindcss 4.2.2, @tailwindcss/postcss 4.2.2, pinia 3.0.4, vazirmatn 33.0.3, jalaali-js 1.2.8, dexie 4.4.2, vitest 4.1.4
- [VERIFIED: go.dev/dl] — Go 1.25.9 and Go 1.26.2 available for download
- [VERIFIED: local environment] — Go 1.24.2, Node v24.14.1, Docker 29.2.0, Docker Compose v5.0.2, sqlc v1.30.0, golangci-lint v1.64.8
- [Context7 /gin-gonic/gin] — Cookie management, middleware patterns, route groups
- [Context7 /websites/nuxt_4_x] — Middleware, layouts, definePageMeta, directory structure
- [Context7 /websites/sqlc_dev_en] — pgx/v5 integration, configuration format
- [Context7 /golang-migrate/migrate] — Go library integration, PostgreSQL driver
- [Context7 /websites/doc_traefik_io_traefik] — Docker labels, ACME Let's Encrypt, HTTPS config

### Secondary (MEDIUM confidence)
- `.planning/research/ARCHITECTURE.md` — System architecture, layered patterns, project structure
- `.planning/research/STACK.md` — Technology stack with rationale
- `.planning/research/PITFALLS.md` — 18 pitfalls with prevention strategies

### Tertiary (LOW confidence)
- Kavenegar API URL pattern — [ASSUMED] based on common API documentation references; verify during implementation

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all versions verified against registries
- Architecture: HIGH — patterns from Context7 + established project research docs
- Pitfalls: HIGH — documented in PITFALLS.md with verification sources
- Go version issue: HIGH — Gin v1.12.0 go.mod confirmed requires Go 1.25.0

**Research date:** 2025-07-18
**Valid until:** 2025-08-18 (stable stack, 30 days)
