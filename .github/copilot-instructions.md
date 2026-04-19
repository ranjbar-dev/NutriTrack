<!-- GSD:project-start source:PROJECT.md -->
## Project

**NutriTrack**

NutriTrack is a Persian-only, mobile-first Progressive Web Application (PWA) for managing the relationship between nutritionists (کارشناس تغذیه) and their clients. Nutritionists create personalized diet plans with meals, options, food items, exercise recommendations, and medication prescriptions. Clients view their plans, track daily intake (food, water, sleep, exercise, weight, body measurements, medications), and communicate with their nutritionist — all with full offline support for viewing and data entry.

**Core Value:** A structured, digital workflow for nutritionists to manage clients and diet plans, and for clients to track daily health activities with full offline capability — in Persian, on mobile.

### Constraints

- **Tech stack**: Go (Gin framework) + PostgreSQL + Nuxt 4 (Vue 3) — chosen by product owner
- **Design**: Persian-only RTL, mobile-first viewport only, no desktop breakpoints
- **Hosting**: Hetzner dedicated/VPS with Docker + Docker Compose + Traefik
- **CI/CD**: GitLab CI/CD pipeline
- **Auth**: JWT (15-min access + 30-day refresh) with OTP for clients (6-digit, 2-min TTL)
- **Offline**: Client-side only (nutritionist/admin always requires connectivity)
- **Chat**: Polling every 10 seconds (not WebSocket)
- **Files**: Local filesystem storage (not S3/MinIO), max 5MB images, 10MB PDFs
- **Performance**: API < 200ms, diet plan load < 500ms, PWA initial load < 3s on 3G
- **Security**: bcrypt cost 12, parameterized queries, row-level authorization, CORS restricted
<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->
## Technology Stack

## Recommended Stack
### Core Technologies
| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| **Go** | 1.25+ | Backend language | High performance, excellent concurrency for handling ~500 concurrent users, single-binary deployment simplifies Docker/Hetzner setup, strong stdlib for HTTP/JSON. Gin requires Go 1.25+. | HIGH — verified via Context7 (Gin README states Go 1.25 prerequisite) |
| **Gin** | v1 (`github.com/gin-gonic/gin`) | HTTP framework | Zero-allocation router, built-in JSON validation via `go-playground/validator`, middleware chaining matches the auth layer needs (JWT + role guards), largest Go web framework community. PRD mentions Fiber/Echo but Gin has the most battle-tested middleware ecosystem and the project context specifies Gin. | HIGH — verified via Context7 |
| **PostgreSQL** | 16 | Primary database | Supports JSONB (for medication `times` array), full-text search (for Persian food search via `pg_trgm`), UUID primary keys natively, excellent Docker support. At 10K clients + 50 nutritionists, single PostgreSQL handles this easily. | HIGH — PRD requirement |
| **Nuxt** | 4.x (`nuxt@^4.0.0`) | Frontend framework | Vue 3-based SSR/SPA framework with file-based routing, auto-imports, and new `app/` directory structure. Nuxt 4 is stable (v4.1.3+ available). Provides Nitro server engine for SSR, which improves PWA first-load performance. | HIGH — verified via Context7 (v4.0.0 and v4.1.3 tagged) |
| **Tailwind CSS** | 4.x | CSS framework | CSS-first configuration (no `tailwind.config.js`), built-in logical properties (`ms-`, `me-`, `ps-`, `pe-`, `text-start`, `text-end`) that auto-flip for RTL via `dir="rtl"` — eliminates the need for a separate RTL plugin. `@import "tailwindcss"` replaces old `@tailwind` directives. | HIGH — verified via Context7 |
| **Docker / Docker Compose** | Latest | Containerization | Multi-service orchestration (Go API + PostgreSQL + Traefik). Single `docker compose up` for dev and prod. Multi-stage Go builds produce ~20MB images. | HIGH — PRD requirement |
| **Traefik** | v3.x | Reverse proxy / HTTPS | Auto Let's Encrypt via ACME HTTP challenge, Docker provider auto-discovers services via labels, built-in HTTP-to-HTTPS redirect. Much simpler than nginx for Docker-native setups. | HIGH — verified via Context7 |
### Backend Libraries
| Library | Import Path | Purpose | Why | Confidence |
|---------|------------|---------|-----|------------|
| **pgx/v5** | `github.com/jackc/pgx/v5` | PostgreSQL driver | Pure Go, high-performance, supports connection pooling via `pgxpool`, LISTEN/NOTIFY, COPY protocol. Direct pgx (not database/sql wrapper) for best performance. Pool config supports health checks, min/max conns, idle timeouts. | HIGH — verified via Context7 |
| **sqlc** | CLI tool (`sqlc-dev/sqlc`) | SQL code generation | Write SQL, generate type-safe Go code. Integrates with pgx/v5 via `sql_package: "pgx/v5"` config. Catches SQL errors at build time, not runtime. Superior to GORM for this project because the data model is well-defined upfront and raw SQL gives full control over PostgreSQL-specific features (full-text search, jsonb). | HIGH — verified via Context7 |
| **golang-migrate** | `github.com/golang-migrate/migrate/v4` | Database migrations | Version-controlled SQL migration files, supports both CLI and Go library usage, PostgreSQL driver built-in. Simpler and more transparent than Atlas for a project this size. | HIGH — verified via Context7 (v4.18.3) |
| **golang-jwt** | `github.com/golang-jwt/jwt/v5` | JWT auth | Industry-standard Go JWT library. Supports custom claims structs, HMAC-SHA256 signing, token parsing with validation. Handles access (15 min) + refresh (30 day) token pattern. | HIGH — verified via Context7 |
| **go-playground/validator** | `github.com/go-playground/validator/v10` | Input validation | Struct tag-based validation, built into Gin's binding system. Supports custom validators for Iranian mobile format, Persian text, enum values. | HIGH — verified via Context7 (v10.27.0) |
| **zerolog** | `github.com/rs/zerolog` | Structured logging | Zero-allocation JSON logger, outputs to stdout for Loki collection. Fastest structured logger in Go. Context-based logger propagation fits the middleware pattern. | HIGH — verified via Context7 |
| **webpush-go** | `github.com/SherClockHolmes/webpush-go` | Web Push notifications | VAPID-based Web Push for PWA notifications. The standard Go library for this. Handles subscription management, payload encryption, and push delivery. | MEDIUM — PRD-specified, verified repository exists |
| **bcrypt** | `golang.org/x/crypto/bcrypt` | Password hashing | Part of Go's extended stdlib. Cost factor 12 as specified in PRD. Battle-tested, no external dependency concerns. | HIGH — Go stdlib |
| **uuid** | `github.com/google/uuid` | UUID generation | v4 UUID generation for all primary keys. Google-maintained, zero-dependency. | HIGH |
### Frontend Libraries
| Library | Package | Purpose | Why | Confidence |
|---------|---------|---------|-----|------------|
| **Pinia** | `pinia` | State management | Official Vue 3 state manager. Stores for auth tokens, active diet plan, tracking data, offline queue. SSR-compatible with Nuxt (auto-hydration). Plugin system for persistence. | HIGH — verified via Context7 |
| **@vite-pwa/nuxt** | `@vite-pwa/nuxt` | PWA module | Zero-config PWA for Nuxt. Handles service worker registration, manifest generation, offline readiness detection, and update prompts via `$pwa` composable. Wraps Workbox under the hood. | HIGH — verified via Context7 |
| **Dexie.js** | `dexie` | IndexedDB wrapper | Simplifies IndexedDB for offline storage of diet plans, pending log entries, cached messages, and sync queue. Supports versioned schemas and complex queries. Do NOT use Dexie Cloud addon (paid service) — use Dexie.js standalone for local-only IndexedDB. | HIGH — verified via Context7 |
| **ofetch** | `ofetch` | HTTP client | Built into Nuxt via `useFetch`/`$fetch`. Automatic JSON parsing, interceptors for auth header injection, 401 handling. Wrap with offline queue logic for sync. | HIGH — Nuxt built-in |
| **Vazirmatn** | `vazirmatn` (npm) | Persian font | Open-source Persian/Arabic web font. Variable font support for optimal loading. The standard Persian web font — free, well-maintained, covers all Persian glyphs. | HIGH |
| **jalaali-js** | `jalaali-js` | Jalali calendar | Converts Gregorian ↔ Jalali (Shamsi) dates. Lightweight (~3KB), no dependencies. All date displays in the app use Shamsi calendar. Wrap in a `useShamsiDate()` composable. | HIGH |
| **Chart.js** | `chart.js` + `vue-chartjs` | Data visualization | Weight trends, body measurements over time, water intake progress. Lightweight, mobile-friendly, RTL-compatible. `vue-chartjs` provides Vue 3 component wrappers. | HIGH — verified via Context7 |
| **@tailwindcss/postcss** | `@tailwindcss/postcss` | Tailwind integration | PostCSS plugin for Tailwind v4 in Nuxt. Direct integration without the `@nuxtjs/tailwindcss` module (which may lag behind Tailwind v4 support). Add to `postcss.config.mjs`. | HIGH — verified via Context7 |
### Development & Infrastructure Tools
| Tool | Version | Purpose | Why | Confidence |
|------|---------|---------|-----|------------|
| **GitLab CI/CD** | N/A | CI/CD pipeline | Lint → Test → Build → Deploy stages. Docker-in-Docker for building images. SSH deploy to Hetzner. PRD requirement. | HIGH — PRD specified |
| **Grafana** | Latest | Monitoring dashboards | Visualize API response times, error rates, user activity. PRD specifies existing Grafana stack. | HIGH — PRD specified |
| **Loki** | Latest | Log aggregation | Collects JSON logs from stdout. Pairs with Grafana for log search/analysis. Go's zerolog JSON output → Docker log driver → Loki. | HIGH — PRD specified |
| **golangci-lint** | Latest | Go linting | Meta-linter combining 50+ linters. Run in CI. Catches bugs, enforces style consistency. | HIGH |
| **sqlc CLI** | Latest | SQL codegen | `sqlc generate` in CI to verify SQL ↔ schema consistency. Install via `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`. | HIGH |
| **ESLint** | 9.x | JS/TS linting | With `@nuxt/eslint-config` for Nuxt-specific rules. Flat config format in ESLint 9. | HIGH |
## Installation
### Backend (Go module init)
# Initialize Go module
# Core dependencies
# Install CLI tools
### Frontend (Nuxt 4 project)
# Initialize Nuxt 4 project
# Core dependencies
# Nuxt modules
# Tailwind CSS v4 (direct PostCSS integration)
# Dev dependencies
### Infrastructure (Docker Compose)
# docker-compose.yml — key services
## Alternatives Considered
| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| **Go HTTP Framework** | Gin | Echo | Gin has larger ecosystem, more middleware libraries, built-in validator integration. Echo is equally capable but Gin's community size means faster problem-solving. |
| **Go HTTP Framework** | Gin | Fiber | Fiber uses fasthttp (not net/http), which breaks compatibility with most Go middleware. Gin uses net/http, ensuring compatibility with the entire Go ecosystem. |
| **SQL Layer** | sqlc | GORM | GORM adds runtime overhead and magic. sqlc generates type-safe code from SQL at build time. For a well-defined schema like NutriTrack, raw SQL with sqlc gives full control over PostgreSQL features (FTS, JSONB, CTEs) without ORM abstraction leaks. |
| **SQL Layer** | sqlc | squirrel | squirrel is a query builder, not a code generator. Still requires manual struct mapping. sqlc eliminates this entirely. |
| **Migrations** | golang-migrate | Atlas | Atlas is more powerful (declarative schema, drift detection) but adds complexity. golang-migrate's sequential SQL files are simpler, transparent, and sufficient for this project's scale. |
| **Frontend Framework** | Nuxt 4 | Next.js (React) | Vue 3 / Nuxt is the PRD choice. Nuxt's auto-imports, file-based routing, and built-in SSR are ideal for this project. Switching to React would be a complete rewrite for no gain. |
| **CSS** | Tailwind v4 | UnoCSS | UnoCSS is faster at build time but Tailwind v4's logical property utilities (`ms-`, `pe-`, etc.) are purpose-built for RTL, and its ecosystem (components, examples) is vastly larger. |
| **CSS RTL** | Tailwind logical properties | tailwindcss-rtl plugin | Tailwind v4 has native logical property support. A separate RTL plugin is unnecessary and would add complexity. Use `ms-` instead of `ml-`, `ps-` instead of `pl-`, `text-start` instead of `text-left`. |
| **Offline Storage** | Dexie.js | idb | Dexie provides a higher-level API with versioned schemas, compound indexes, and better query ergonomics. idb is thinner but requires more boilerplate for the complex offline data model (diet plans, queued logs, cached messages). |
| **Charting** | Chart.js | ECharts | ECharts is more powerful but significantly heavier (~1MB). Chart.js (~60KB) is sufficient for simple line/bar charts (weight trends, water intake). Mobile-first requires small bundles. |
| **Logging (Go)** | zerolog | zap | Both are zero-alloc. zerolog has a simpler API and produces cleaner JSON output for Loki. zap requires more configuration for equivalent output. |
| **Reverse Proxy** | Traefik v3 | Nginx | Traefik auto-discovers Docker services via labels, auto-renews Let's Encrypt certs, and requires zero config file management. Nginx requires manual config and certbot cron jobs. |
| **Push Notifications** | webpush-go | Firebase FCM | FCM requires Google account dependency and adds external service coupling. VAPID-based Web Push is a web standard, works with any browser, and keeps everything self-hosted on Hetzner. |
## What NOT to Use
| Technology | Why Avoid |
|------------|-----------|
| **GORM** | Adds ORM abstraction over PostgreSQL that leaks in complex queries (diet plan joins, FTS, aggregations). sqlc generates cleaner code from raw SQL. GORM's auto-migration is dangerous in production. |
| **Redis** | Unnecessary at this scale. OTP storage, session data, and caching can all live in PostgreSQL. Adding Redis doubles infrastructure complexity for zero benefit at 500 concurrent users. The PRD explicitly recommends against it. |
| **WebSockets** | PRD explicitly chose polling (10s interval) for chat. WebSocket adds connection management complexity, doesn't work offline, and the chat latency tolerance is high enough for polling. |
| **Dexie Cloud** | Paid cloud sync service. Use Dexie.js standalone for local IndexedDB only. Build custom sync logic with the offline queue pattern described in the PRD. |
| **@nuxtjs/tailwindcss module** | May lag behind Tailwind v4 support. Direct `@tailwindcss/postcss` integration is simpler, more reliable, and gives full control over Tailwind v4's CSS-first configuration. |
| **tailwindcss-rtl plugin** | Completely unnecessary with Tailwind v4. Use logical properties natively: `ms-` (margin-inline-start), `me-` (margin-inline-end), `ps-`, `pe-`, `text-start`, `text-end`, `float-start`, `border-s-`. These auto-flip based on `dir="rtl"`. |
| **i18n libraries (vue-i18n, @nuxtjs/i18n)** | Persian-only app. No i18n needed. Hardcode Persian strings directly in components. Adding i18n infrastructure for a single language is pure overhead. |
| **Moment.js / Day.js** | Not needed. `jalaali-js` handles Jalali conversion. Native `Date` + `Intl.DateTimeFormat('fa-IR')` handles Persian number formatting. These libraries add unnecessary bundle size. |
| **MinIO / S3** | Overkill for local file storage on Hetzner. Simple filesystem paths stored in PostgreSQL is the right approach at this scale. Can migrate to object storage later if needed. |
| **GraphQL** | REST JSON API is simpler and sufficient. The data access patterns (CRUD + simple relations) don't benefit from GraphQL's query flexibility. GraphQL would add unnecessary schema complexity. |
## Version Compatibility Matrix
| Component | Min Version | Recommended | Notes |
|-----------|-------------|-------------|-------|
| Go | 1.25 | 1.25+ | Required by Gin (verified via Context7) |
| PostgreSQL | 15 | 16 | v16 has better JSON, parallelism. Use `-alpine` Docker image. |
| Node.js | 18.x | 22.x LTS | Required for Nuxt 4 build tooling |
| Nuxt | 4.0.0 | 4.1.x | Stable release; uses `app/` directory structure by default |
| Tailwind CSS | 4.0 | 4.x | CSS-first config; `@import "tailwindcss"` |
| Docker Compose | v2 | v2.x | Compose v2 is the Docker CLI plugin (not standalone `docker-compose`) |
| Traefik | 3.0 | 3.4+ | v3 has improved Docker provider and Let's Encrypt support |
| pgx | v5 | v5.x | Must match sqlc `sql_package: "pgx/v5"` config |
| golang-migrate | v4 | v4.18+ | PostgreSQL driver support, `file://` source |
| Dexie.js | 3.x | 4.x | v4 has improved TypeScript support and performance |
## Key Configuration Patterns
### Tailwind v4 RTL Setup (No Plugin Needed)
### Nuxt 4 PostCSS Config
### sqlc Configuration
# sqlc.yaml
### pgxpool Configuration
## Sources
- **Gin Framework**: Context7 `/gin-gonic/gin` and `/websites/gin-gonic_en` — Go 1.25 requirement, middleware patterns, JWT example
- **pgx/v5**: Context7 `/jackc/pgx` and `/websites/pkg_go_dev_github_com_jackc_pgx_v5` — pgxpool configuration, connection management
- **sqlc**: Context7 `/websites/sqlc_dev_en` — pgx/v5 integration, configuration format
- **golang-migrate**: Context7 `/golang-migrate/migrate` (v4.18.3) — PostgreSQL setup, Go library usage
- **golang-jwt**: Context7 `/golang-jwt/jwt` — token creation/parsing, custom claims
- **Nuxt 4**: Context7 `/nuxt/nuxt` (v4.0.0, v4.1.3) and `/websites/nuxt_4_x` — directory structure, upgrade path, compatibilityVersion
- **Tailwind CSS v4**: Context7 `/websites/tailwindcss` — RTL logical properties, CSS-first config, PostCSS setup
- **Pinia**: Context7 `/vuejs/pinia` — SSR hydration, store plugins, Nuxt integration
- **@vite-pwa/nuxt**: Context7 `/websites/vite-pwa-org_netlify_app` — Nuxt module setup, `$pwa` composable, service worker lifecycle
- **Dexie.js**: Context7 `/websites/dexie` — schema versioning, offline-first patterns (standalone, NOT Dexie Cloud)
- **Traefik**: Context7 `/websites/doc_traefik_io_traefik` — Docker Compose setup, Let's Encrypt ACME, HTTPS configuration
- **Chart.js**: Context7 `/chartjs/chart.js` and `/websites/chartjs`
- **go-playground/validator**: Context7 `/go-playground/validator` (v10.27.0)
- **zerolog**: Context7 `/rs/zerolog`
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->
## Project Skills

No project skills found. Add skills to any of: `.github/skills/`, `.agents/skills/`, `.cursor/skills/`, `.github/skills/`, or `.codex/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
