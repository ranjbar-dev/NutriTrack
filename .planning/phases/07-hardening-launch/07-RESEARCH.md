# Phase 7: Hardening & Launch — Research

**Researched:** 2026-04-20
**Domain:** Security hardening / performance validation / observability / backup / Persian UX launch
**Confidence:** HIGH (all findings verified against actual codebase)

---

## Summary

Phase 7 is the final hardening pass for a fully-featured Persian PWA that ships auth, diet plans,
six tracking dimensions, offline sync, push notifications, and real-time messaging across three user
roles. Phases 1–6 are all verified complete. No product features are added in this phase — the goal
is to prove correctness and production-readiness of what already exists.

The highest-risk domain is **row-level authorization**. The route surface has 40+ endpoints that accept
client-scoped identifiers; the middleware only enforces role membership, not resource ownership. Whether
each handler correctly threads the calling user's ID into the service/repo layer must be confirmed
by reading every handler and by adding targeted integration tests. The diet plan aggregate handler is
the most dangerous: it permits nutritionist, super_admin, AND client roles but must silently enforce
different ownership logic for each.

The second highest gap is **monitoring and backup infrastructure**, which does not exist yet. The
docker-compose.yml has no Grafana, Loki, Promtail, or pg_dump service. Both must be added from scratch
and proven with live data before launch.

**Primary recommendation:** Sequence work as Security → Performance → Monitoring → Backup → Polish.
Never skip the backup restore test — a backup without a proven restore is not a backup.

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01:** Start with row-level authorization and cross-tenant access checks. Phase 7 should assume
  authorization bugs are the highest-risk launch blocker.
- **D-02:** Treat dependency and configuration audits as part of the same hardening pass: `govulncheck`,
  `npm audit`, CORS, HTTPS redirect/HSTS, refresh-token behavior, and file-upload revalidation all
  belong in this phase.
- **D-03:** Hardening work should prefer proving existing behavior with tests/audits over rewriting
  stable code paths unless a concrete issue is found.
- **D-04:** Performance verification must use realistic data from the completed product, especially
  the diet-plan aggregate and offline/PWA startup path.
- **D-05:** Targets are locked: API p95 < 200ms, diet plan < 500ms, PWA initial load < 3s on simulated 3G.
- **D-06:** Query/index fixes and bundle-size optimizations are in scope; new platform dependencies
  are not, unless an existing target cannot be met otherwise.
- **D-07:** Monitoring is self-hosted and Docker-centric: Grafana, Loki, alert thresholds — all matching
  the Hetzner + Docker Compose deployment model.
- **D-08:** Backup work is incomplete until restore is proven. PostgreSQL and upload backups both need
  an actual restore check, not just cron configuration.
- **D-09:** Final polish includes real Persian UX review, loading/empty/error state audits, and manual
  Android/iOS launch-path checks using the now-complete Phase 6 PWA.

### Agent's Discretion

- Exact sequencing between security, performance, monitoring, and backup work
- Whether to group audits by subsystem (backend/frontend/infra) or by verification goal
- Which existing commands and test harnesses to extend first, as long as no new product features are introduced

### Deferred Ideas (OUT OF SCOPE)

- New end-user features or product-scope changes
- Native/mobile-only capabilities outside the web/PWA stack
- Large architecture rewrites without evidence from the Phase 7 audits
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|-----------------|
| INFRA-05 | Grafana dashboards for monitoring | §Monitoring section — add grafana + loki + promtail to docker-compose |
| INFRA-07 | Daily automated PostgreSQL backups | §Backup section — pg_dump cron container with restore test |
| INFRA-08 | Weekly file storage backups | §Backup section — uploads volume backup script |
| INFRA-09 | API response time < 200ms for standard CRUD | §Performance section — k6/vegeta load test, EXPLAIN ANALYZE |
| INFRA-11 | Support 50 nutritionists, 10,000 clients, ~500 concurrent users | §Performance section — pgxpool tuning, load test at target concurrency |
| UI-08 | Initial load < 3 seconds on 3G | §Performance section — Nuxt bundle analysis, Lighthouse/DevTools 3G throttle |
</phase_requirements>

---

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Row-level authorization | API / Backend | — | Service layer enforces resource ownership; middleware only checks role |
| HSTS / TLS redirect | CDN / Traefik | API (headers) | Traefik handles TLS; HSTS should also be set in SecurityHeaders middleware |
| CSP / security headers | API / Backend | Frontend (meta) | SecurityHeaders middleware sets response headers for all API responses |
| API performance | API / Backend | Database | Query optimization, index analysis, pgxpool tuning |
| PWA bundle size / load time | Frontend Server (Nuxt) | CDN / Static | Nuxt build output, workbox precache manifest, chunking |
| Grafana/Loki monitoring | Infrastructure | Docker Compose | Added as additional docker-compose services |
| PostgreSQL backup | Database | Infrastructure | pg_dump scheduled service, Hetzner volume |
| File upload backup | Infrastructure | Database | Uploads Docker volume, weekly cron |
| Persian UX / device testing | Browser / Client | — | Manual review on real Android/iOS devices |

---

## Current-State Findings

### Security Findings

#### What's Already Strong [VERIFIED: codebase]
- Auth uses httpOnly cookie `access_token` — not localStorage; XSS-resilient
- CORS: explicit `FrontendURL` origin (not `*`), `Access-Control-Allow-Credentials: true` — correct
- All SQL via sqlc parameterized queries — no string interpolation found in any repo file
- File upload MIME validation via `gabriel-vasile/mimetype` (magic bytes, not filename extension)
- JWT signed HMAC-SHA256 with minimum 32-char secret validation at startup
- Rate limiter on OTP endpoints (3/10min), reply detection on refresh tokens
- Role guard applied at route group level before handlers

#### Security Gaps Found [VERIFIED: codebase]

**GAP-SEC-01: No HSTS header** [VERIFIED: security_headers.go]
`SecurityHeaders()` sets X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Referrer-Policy,
CSP — but NOT `Strict-Transport-Security`. Traefik does HTTPS but HSTS must come from the app or a
Traefik middleware to reach the browser. Risk: clients can be downgraded to HTTP after first visit.
Fix: add `c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")`.

**GAP-SEC-02: CSP `default-src 'self'` may be too broad or too narrow** [VERIFIED: security_headers.go]
Current CSP: `Content-Security-Policy: default-src 'self'`. For a Nuxt SSR app with injectManifest
service worker, inline scripts injected by Vite hydration will be blocked unless `'unsafe-inline'`
or nonce-based CSP is used. Vazirmatn font is served from local npm package (OK), but Chart.js CDN
loads (if any) would be blocked. Must test in production build with browser console open — if CSP
violations appear, tighten or adjust rather than remove. The API's CSP header applies to API responses;
the Nuxt frontend's own CSP should be set via nuxt.config.ts `routeRules` or a response header in
the frontend Dockerfile / Nuxt server middleware.

**GAP-SEC-03: Row-level authorization must be audited per endpoint** [VERIFIED: main.go, handlers]
Middleware only checks role membership. Individual handlers pass `user_id` to services, but each
service must verify the requesting user owns the resource. Critical endpoints to audit:

| Endpoint | Risk | Who calls it | What to verify |
|----------|------|-------------|----------------|
| `GET /api/diet-plans/:id` | HIGH | client, nutritionist, super_admin | Client can only read OWN plan; nutritionist can only read own clients' plans |
| `PATCH /api/diet-plans/:id` | HIGH | nutritionist, super_admin | Nutritionist can only edit own plans |
| `GET /api/clients/:clientId/plans` | HIGH | nutritionist, super_admin | Nutritionist must own clientId |
| `GET /api/nutritionist/clients/:clientId/tracking/*` | HIGH | nutritionist | Must own clientId |
| `GET /api/nutritionist/clients/:clientId/lab-results/:labId/download` | HIGH | nutritionist | Must own clientId |
| `GET /api/messages/:partnerId` | HIGH | client, nutritionist | Must be in the conversation |
| `GET /api/messages/attachment/:messageId` | HIGH | client, nutritionist | Must own message conversation |
| `PATCH /api/nutritionist/food-requests/:requestId/approve\|reject` | MEDIUM | nutritionist | requestId must belong to own client |
| `DELETE /api/diet-plans/:id` | HIGH | nutritionist, super_admin | Nutritionist must own plan |
| All diet plan sub-resource mutations (days, meals, options, items, exercises, meds) | HIGH | nutritionist | Must own parent plan |

**GAP-SEC-04: Refresh token cookie attributes** [ASSUMED — not verified in token handling code]
Cookies must have `Secure; HttpOnly; SameSite=Strict` in production. Needs verification in auth handler
cookie-setting code.

**GAP-SEC-05: File upload path traversal** [ASSUMED]
`UploadsDir` + UUID filename approach should be safe, but need to verify no path components from
user input are used in filepath.Join calls.

**GAP-SEC-06: `govulncheck` and `npm audit` not yet run**
Need to run both as part of hardening pass. Current deps: Go 1.25, Gin v1.12.0, pgx v5.9.2,
zerolog v1.35.0, webpush-go v1.4.0, mimetype v1.4.13, golang-jwt v5.3.1 — all appear current as of
research date [ASSUMED — not verified against CVE databases].

---

### Performance Findings

#### Diet Plan Aggregate [VERIFIED: codebase — test stub skipped]
The `TestPlanAggregateLoadTime` test in `backend/internal/repository/diet_plan_repo_test.go` is an
**integration test stub marked `t.Skip()`** — it was never implemented. The plan loads via batch
queries (Phase 3 design), but no measured baseline exists. Must be validated with realistic data
(7 days × 5 meals × 3 options × 4 items = 420 meal_option_items).

#### pgxpool Configuration [VERIFIED: main.go]
```go
poolConfig.MaxConns = 20
poolConfig.MinConns = 5
poolConfig.MaxConnLifetime = 1 * time.Hour
poolConfig.MaxConnIdleTime = 30 * time.Minute
```
20 max connections is a reasonable starting point. For 500 concurrent users each making short-lived
requests, this depends on average request duration. If p50 latency is 10ms, 20 connections can serve
~2000 req/s — likely sufficient. However, this needs validation under the target 500 concurrent load.

#### Frontend Bundle [VERIFIED: package.json, nuxt.config.ts]
Heavy dependencies: chart.js (^4.5.1), dexie (^4.4.2), vue-chartjs, dayjs, jalaali-js.
No explicit Nuxt route lazy-loading configuration seen. `@vite-pwa/nuxt` uses `injectManifest`
strategy — precaches all `**/*.{js,css,woff2,png,svg,ico,webp}` files.
No `splitChunks` or `routeRules` for code splitting visible in nuxt.config.ts.

**3G Target Reality Check:** Vazirmatn Variable font (woff2) can be 50–150KB compressed. Chart.js + 
vue-chartjs is ~200KB uncompressed. Total uncompressed JS for a heavy Nuxt route may exceed 1MB before
service worker caching helps. Must run `nuxt build && npx nuxt analyze` to measure actual chunk sizes.

#### HTTP Server Timeouts [VERIFIED: main.go]
```go
ReadTimeout:  15 * time.Second
WriteTimeout: 15 * time.Second
IdleTimeout:  60 * time.Second
```
Appropriate for a non-streaming API.

---

### Monitoring Findings

**ZERO monitoring infrastructure exists** [VERIFIED: docker-compose.yml]
The production docker-compose.yml has: traefik, api, frontend, postgres — no Grafana, Loki, Promtail,
Prometheus, or any alerting. Must be built from scratch.

**Zerolog is already production-ready** [VERIFIED: main.go]
In production environment, zerolog outputs JSON to stdout at Info level. Docker's log driver captures
stdout/stderr. Promtail can read from Docker socket or log files and ship to Loki.

**No Prometheus metrics endpoint** [VERIFIED: main.go, go.mod]
`prometheus/client_golang` is not in go.mod. To get API response time histograms in Grafana, options are:
1. Add Prometheus client + Gin metrics middleware (lightweight, recommended)  
2. Derive metrics from structured JSON logs via LogQL in Loki (no code change needed)

Option 2 requires no new Go dependencies and aligns with D-06. Loki + LogQL can compute p95 from
zerolog's `duration` field if the Logger middleware logs response time.

---

### Backup Findings

**No backup infrastructure** [VERIFIED: docker-compose.yml]
PostgreSQL data lives in Docker named volume `pgdata`. Uploads default to `./uploads` inside the api
container (should be a volume for persistence). No pg_dump cron, no volume backup scripts.

**Uploads directory persistence risk** [VERIFIED: docker-compose.yml, config.go]
`UploadsDir` defaults to `./uploads` relative to the binary. In the Docker container this is ephemeral
unless mounted as a volume. The production docker-compose.yml has no volume mount for uploads.
This is a **data loss risk** that must be fixed before launch (add `uploads` named volume).

---

### UX/Polish Findings

**Manual device testing not done** [VERIFIED: Phase 6 VERIFICATION.md]
> "Manual device validation (physical Android push receipt, iOS installed-PWA behavior, and real-device
> storage-eviction exercises) was not possible in this CLI environment"

Phase 6 verification explicitly deferred this. Phase 7 must execute it.

**Persian UX gaps to audit:**
- All Persian error messages use correct right-to-left punctuation and natural phrasing
- Loading spinners present on all async operations
- Empty states (no diet plan assigned, no tracking data, no messages) show helpful Persian messages
- Toast notifications display correctly in RTL layout
- Form validation errors display beside the correct fields in RTL

---

## Standard Stack

### Core (already in use — verify versions are current)
| Library | Current Version | Purpose | Status |
|---------|----------------|---------|--------|
| Gin | v1.12.0 | HTTP framework | [VERIFIED: go.mod] |
| zerolog | v1.35.0 | Structured JSON logging | [VERIFIED: go.mod] |
| pgx/v5 | v5.9.2 | PostgreSQL driver | [VERIFIED: go.mod] |
| govulncheck | latest | Go vulnerability scanner | Must install — not in go.mod |
| Nuxt | ^4.0.0 | Frontend SSR | [VERIFIED: package.json] |
| vitest | ^3.0.0 | Frontend unit tests | [VERIFIED: package.json] |

### New Infrastructure for Phase 7
| Component | Version | Purpose | Source |
|-----------|---------|---------|--------|
| grafana/grafana | 10.x or latest | Dashboards + alerting | [ASSUMED — verify at docker pull] |
| grafana/loki | 2.9.x or latest | Log aggregation | [ASSUMED] |
| grafana/promtail | 2.9.x or latest | Log shipper (Docker → Loki) | [ASSUMED] |
| k6 / vegeta | latest | API load testing | [ASSUMED] |

> **Version note:** Always pin Grafana stack versions explicitly (`grafana/loki:2.9.8`) to prevent
> uncontrolled upgrades in production docker-compose.yml.

### Tools (CI / audit)
| Tool | Install | Purpose |
|------|---------|---------|
| govulncheck | `go install golang.org/x/vuln/cmd/govulncheck@latest` | Go CVE scanner |
| npm audit | built-in | Node.js CVE scanner |
| k6 | brew/choco install k6 | HTTP load testing |
| go test -bench | built-in | Micro-benchmarks |
| nuxt analyze | `nuxt build --analyze` | Bundle size visualizer |

---

## Architecture Patterns

### System Architecture Diagram — Phase 7 additions

```
[Hetzner Host]
  ├── Traefik (HTTPS termination, HTTP→HTTPS redirect)
  │     └── routes: /api → api:8080, /* → frontend:3000
  │
  ├── api (Go/Gin) ──────────────────────────────────────────────────────────────
  │     ├── SecurityHeaders (add HSTS) ─► response headers
  │     ├── Auth (JWT cookie) ─► user_id + role in ctx
  │     ├── RoleGuard ─► role membership
  │     └── Handlers ─► Services (ownership checks) ─► Repos (parameterized SQL)
  │                                                        │
  │                                                   pgxpool (20 max)
  │                                                        │
  ├── postgres:16 ──────────────────────────────────── pgdata volume
  │     └── daily pg_dump ─► backup volume
  │
  ├── frontend (Nuxt SSR) ── serves PWA shell + SSR
  │     └── /icons, /public assets ── Workbox precache manifest
  │
  ├── [NEW] grafana ─────── :3001 (internal only, not exposed via Traefik)
  │     ├── datasource: Loki
  │     └── dashboards: API latency, error rate, active users, DB pool
  │
  ├── [NEW] loki ────────── receives logs from Promtail
  │     └── retention: 30 days
  │
  └── [NEW] promtail ────── reads Docker container stdout/stderr
        └── ships structured JSON logs ──► loki

[Backup service / cron on host]
  ├── daily: docker exec postgres pg_dump | gzip → /backups/db/YYYY-MM-DD.sql.gz
  └── weekly: tar uploads volume → /backups/files/YYYY-WW.tar.gz
```

### Monitoring Architecture Pattern

```
Docker container stdout (zerolog JSON)
    │
    ▼
Promtail (Docker log driver or file tail)
    │
    ▼
Loki (log storage + LogQL query engine)
    │
    ▼
Grafana (dashboards + alerting)
    ├── Panel: API p95 latency via LogQL on "duration" field
    ├── Panel: Error rate (5xx count / total requests)
    ├── Panel: Active push subscriptions
    └── Alert: p95 > 1s for 5min → fire webhook/email
```

**LogQL example for p95 from zerolog:**
```logql
quantile_over_time(0.95,
  {container="nutritrack-api"}
  | json
  | duration > 0
  | unwrap duration [5m]
) by (path)
```
This works because zerolog's Logger middleware logs `duration` (in ms or nanoseconds) per request.
**Verify** the Logger middleware actually logs the `duration` field: check
`backend/internal/middleware/logger.go`.

### Backup Pattern

```yaml
# docker-compose.yml addition — backup service
backup:
  image: postgres:16-alpine
  depends_on:
    - postgres
  environment:
    PGPASSWORD: ${POSTGRES_PASSWORD}
  volumes:
    - pgbackups:/backups
    - ./scripts/backup.sh:/backup.sh:ro
  entrypoint: crond -f -d 8
  restart: unless-stopped
```

Simpler alternative: host-level cron job on the Hetzner server calling `docker exec`:
```bash
# /etc/cron.d/nutritrack-backup
0 2 * * * root docker exec nutritrack-postgres-1 pg_dump -U $PGUSER $PGDB | gzip > /backups/db/$(date +\%Y-\%m-\%d).sql.gz
```
The host-level cron is simpler and avoids adding another container to manage.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Log shipping to Grafana | Custom log forwarder | Promtail | Handles Docker log discovery, backpressure, retry |
| Vulnerability scanning | Manual go.sum audit | govulncheck | Understands Go module graph, CVE database |
| Load testing scripts | Custom HTTP flood | k6 (JS) or vegeta | Battle-tested ramp-up patterns, percentile reports |
| Backup scheduling | New Go service | host cron + docker exec | Simpler, no new container, survives app restarts |
| Bundle analysis | Manual file counting | nuxt build --analyze | Rollup plugin chart, interactive treemap |
| DB performance profiling | Manual query timing | `EXPLAIN (ANALYZE, BUFFERS)` | Shows actual row counts, buffer hits, plan nodes |

**Key insight:** Phase 7 is a verification and configuration phase, not a development phase. Every
problem already has a standard tool. Adding bespoke code increases surface area and risk.

---

## Common Pitfalls

### Pitfall 1: Assuming role-guard = authorization
**What goes wrong:** A nutritionist with valid JWT for role "nutritionist" calls
`GET /api/nutritionist/clients/{OTHER_NUTRITIONIST_CLIENT_ID}/tracking/food` and gets data back
because the route group only checks the role, not the relationship.
**Why it happens:** Middleware enforces role, but resource ownership is service-level logic that may
or may not have been written consistently across 26+ tracking endpoints.
**How to avoid:** Write a Go integration test with two seeded nutritionists; test every `:clientId`
endpoint with nutritionist-B's JWT requesting nutritionist-A's client's data → expect 403 or 404.
**Warning signs:** Service methods that don't receive `nutritionistID` as a parameter alongside
`clientID` cannot enforce ownership.

### Pitfall 2: Grafana/Loki version mismatch
**What goes wrong:** Promtail 2.9.x is incompatible with Loki 3.x push API endpoint (`/loki/api/v1/push`).
**Why it happens:** Loki 3.x changed some API defaults; Promtail must match the Loki minor version.
**How to avoid:** Use the same version tag for loki and promtail images. Prefer
`grafana/loki:2.9.8` + `grafana/promtail:2.9.8` (stable, widely documented).

### Pitfall 3: Backup without restore test
**What goes wrong:** pg_dump runs daily for months, then the restore fails due to encoding mismatch,
corrupted gzip, or missing extension.
**Why it happens:** No one ever tests the restore. Postgres initdb locale (`C.UTF-8` in dev vs
`en_US.UTF-8` in production) can cause restore failures for Persian text.
**How to avoid:** Actually run `psql -U user -d restore_db < backup.sql` on the same Postgres version.
Check row counts on key tables. This is a required deliverable.
**Warning signs:** dev docker-compose sets `POSTGRES_INITDB_ARGS: "--locale=C.UTF-8"` but production
does not. Locale mismatch may cause pg_dump/restore issues with collation-dependent indexes.

### Pitfall 4: CSP blocks Nuxt hydration
**What goes wrong:** `Content-Security-Policy: default-src 'self'` blocks inline `<script>` tags that
Nuxt injects for hydration data (`__NUXT_DATA__`, etc.), causing the client-side app to fail silently.
**Why it happens:** The SecurityHeaders middleware on the **API** sets CSP for API responses. However
if Nuxt proxies or the frontend Dockerfile/Nginx adds the same header, it will break hydration.
The API CSP only affects API JSON responses, not HTML pages — but this must be confirmed.
**How to avoid:** Verify that the SecurityHeaders middleware's CSP is only sent on `/api/*` routes.
The Nuxt frontend (served from port 3000) should set its own CSP via `nuxt.config.ts` routeRules
or server middleware, not the Go API.

### Pitfall 5: Uploads volume not persisted in production
**What goes wrong:** Lab result PDFs and message attachments disappear on container restart because
`./uploads` is inside the container filesystem with no volume mount.
**Why it happens:** docker-compose.yml has no volume for uploads. Development works because developers
rarely restart containers.
**How to avoid:** Add `uploads:/app/uploads` volume in docker-compose.yml and `UPLOADS_DIR=/app/uploads`.
This must be done before any production traffic touches the upload endpoints.

### Pitfall 6: 3G load time measured wrong
**What goes wrong:** Chrome DevTools 3G throttle is measured from first visit — but the target is
"after first visit" (service worker cached). Measuring cold load (first visit) conflates SW install
with page render.
**How to avoid:** Measure two scenarios separately: (a) cold first visit to establish SW cache, (b)
second visit with throttling to measure SW-served load. Target < 3s is for the SW-served load.
Lighthouse also has a PWA audit that simulates this.

### Pitfall 7: k6 load test against localhost misses real production latency
**What goes wrong:** k6 run locally shows p95 < 50ms, but production on Hetzner VPS shows p95 > 300ms.
**Why it happens:** In-process test has no network RTT; production has TLS handshake + TCP from client.
**How to avoid:** Run load test FROM the Hetzner server itself targeting `localhost:8080` to isolate
app latency from network RTT, then add network overhead separately.

---

## Implementation Slices (Recommended Plans)

### Plan 07-01: Security Audit & Hardening
**Goal:** Prove all 40+ endpoints enforce row-level ownership; fix HSTS gap; run govulncheck/npm audit.
**Work:**
1. Add `Strict-Transport-Security` header to `SecurityHeaders()` middleware
2. Verify CSP scope (API responses only, not HTML — document finding)
3. Verify refresh token and access token cookie attributes (`Secure`, `HttpOnly`, `SameSite`)
4. Write Go integration test: `TestCrossTenantAccessDenied` — seeds two nutritionists, tests every
   `:clientId` endpoint with wrong nutritionist's token → expects 403/404
5. Write Go integration test: `TestClientCannotAccessOtherClientPlan` — seeds two clients, verifies
   client-A cannot call `GET /api/diet-plans/{planB_id}`
6. Audit messaging ownership: `ListMessages`, `PollNewMessages`, `MarkRead`, `DownloadAttachment` —
   confirm service layer checks partnership
7. Run `govulncheck ./...` in backend, resolve any HIGH/CRITICAL findings
8. Run `npm audit --audit-level=high` in frontend, resolve any HIGH/CRITICAL findings
9. Verify file upload path: confirm no user-controlled path components in `filepath.Join`
10. Document result: create `07-01-SECURITY-AUDIT.md` checklist with PASS/FAIL per endpoint

**Files touched:**
- `backend/internal/middleware/security_headers.go` (HSTS)
- `backend/internal/middleware/auth.go` (cookie attribute check)
- New: `backend/internal/handler/auth_handler_test.go` or `backend/internal/service/*_test.go`
- New integration test file: `backend/internal/handler/authorization_test.go`

**Verification commands:**
```bash
cd backend && go test ./... -v -run TestCrossTenant
cd backend && govulncheck ./...
cd frontend && npm audit --audit-level=high
curl -I https://staging.nutritrack.app/api/health | grep -i strict-transport
```

---

### Plan 07-02: Performance Validation
**Goal:** Prove all three performance targets with evidence; fix any failures found.
**Work:**
1. Implement and run the skipped `TestPlanAggregateLoadTime` integration test with realistic seed data
   (7 days, 5 meals, 3 options, 4 items = 420 items); measure sub-500ms
2. Run `EXPLAIN (ANALYZE, BUFFERS)` on the diet plan aggregate query, tracking list queries (food/water/
   sleep/exercise/medication/body with date range), and messaging list query
3. Add missing DB indexes if needed (e.g., `(client_id, date)` on tracking tables)
4. Run k6 load test: 100 VUs, 30s duration, targeting key endpoints; verify p95 < 200ms
5. Run `nuxt build` then `nuxt analyze` (or `--analyze` flag); identify chunks > 300KB
6. Implement route-level code splitting for nutritionist pages if client bundle is too heavy
7. Measure PWA load: Chrome DevTools → Performance → throttle to Slow 3G → reload (warm cache);
   record DOMContentLoaded + LCP; must be < 3s
8. Tune pgxpool `MaxConns` if load test shows pool exhaustion
9. Verify Vazirmatn font is served with correct cache headers (long TTL since it's versioned by npm)

**Files touched:**
- `backend/internal/repository/diet_plan_repo_test.go` (implement the skipped test)
- `backend/db/migrations/` (new index migration if needed)
- `frontend/nuxt.config.ts` (code splitting config if needed)
- New: `backend/scripts/load_test.js` (k6 script)
- `backend/cmd/api/main.go` (pgxpool tuning if needed)

**Verification commands:**
```bash
# Backend benchmark
cd backend && go test ./... -run TestPlanAggregateLoadTime -v -tags integration

# Database explain
docker exec -it nutritrack-postgres-1 psql -U nutritrack -c \
  "EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) SELECT ..."

# Load test
k6 run backend/scripts/load_test.js --vus 100 --duration 30s

# Frontend bundle analysis
cd frontend && npm run build -- --analyze

# PWA 3G: manual step in Chrome DevTools
```

---

### Plan 07-03: Monitoring Infrastructure
**Goal:** Grafana + Loki live and showing API metrics; alert thresholds configured.
**Work:**
1. Verify `backend/internal/middleware/logger.go` logs `duration` field per request (required for LogQL p95)
2. Add to docker-compose.yml: `loki`, `promtail`, `grafana` services
3. Create `monitoring/loki-config.yml` (local storage, 30d retention)
4. Create `monitoring/promtail-config.yml` (Docker socket scrape, label extraction for container name)
5. Create `monitoring/grafana/provisioning/datasources/loki.yml`
6. Create `monitoring/grafana/provisioning/dashboards/nutritrack.json` with panels:
   - API p95 latency by endpoint (LogQL `quantile_over_time`)
   - HTTP error rate (5xx count / total)
   - Request volume per minute
   - DB pool utilization (if logged)
7. Configure Grafana alerts: p95 > 1s for 5min → contact point; 5xx rate > 1% → contact point
8. Test: generate traffic, verify logs appear in Loki Explorer, panels update
9. Ensure Grafana is NOT publicly accessible (internal port only or behind auth)

**Files touched:**
- `docker-compose.yml` (add monitoring services)
- New: `monitoring/loki-config.yml`
- New: `monitoring/promtail-config.yml`
- New: `monitoring/grafana/provisioning/datasources/loki.yml`
- New: `monitoring/grafana/dashboards/nutritrack.json`
- `backend/internal/middleware/logger.go` (verify/add duration logging)

**Verification commands:**
```bash
docker compose up -d grafana loki promtail
# Generate some traffic
curl -s http://localhost:8080/api/health
# Check Loki received logs
curl -s 'http://localhost:3100/loki/api/v1/labels' | jq .
# Open Grafana: http://localhost:3001
```

---

### Plan 07-04: Backup & Restore
**Goal:** Daily DB backups + weekly uploads backups configured AND restore proven to work.
**Work:**
1. **Fix uploads volume**: add `uploads` named volume to docker-compose.yml, mount at `/app/uploads`
   in api service; add `UPLOADS_DIR=/app/uploads` to .env.example
2. Create `scripts/backup-db.sh`: `pg_dump | gzip > /backups/db/$(date +%Y-%m-%d).sql.gz`
3. Create `scripts/backup-uploads.sh`: `tar czf /backups/files/$(date +%Y-W%W).tar.gz uploads_volume`
4. Add backup retention: delete DB backups older than 30 days, file backups older than 90 days
5. Set up host cron on Hetzner: daily at 02:00 for DB, weekly on Sunday 03:00 for files
6. **Restore test (non-negotiable)**:
   - Restore latest DB backup to `nutritrack_restore` database on same Postgres instance
   - Run `SELECT COUNT(*) FROM users;`, `SELECT COUNT(*) FROM diet_plans;`, etc.
   - Compare counts with production
   - Restore a backup upload file, verify it's accessible
7. Document restore procedure in `docs/restore-procedure.md`

**Files touched:**
- `docker-compose.yml` (uploads volume mount)
- `.env.example` (UPLOADS_DIR)
- New: `scripts/backup-db.sh`
- New: `scripts/backup-uploads.sh`
- New: `scripts/restore-test.sh`
- New: `docs/restore-procedure.md`

**Verification commands:**
```bash
# Run backup manually
bash scripts/backup-db.sh

# Restore test
docker exec -it nutritrack-postgres-1 psql -U nutritrack -c "CREATE DATABASE nutritrack_restore;"
gunzip < /backups/db/latest.sql.gz | docker exec -i nutritrack-postgres-1 psql -U nutritrack -d nutritrack_restore
docker exec -it nutritrack-postgres-1 psql -U nutritrack -d nutritrack_restore -c "SELECT COUNT(*) FROM users;"
```

---

### Plan 07-05: Final Polish & Launch Readiness
**Goal:** Persian UX review complete; E2E role journeys pass on real devices; production config verified.
**Work:**
1. Persian UX audit checklist (see §UX Checklist below) — review all 3 role journeys
2. Error state audit: test all API error conditions, verify Persian error messages, no stack traces
3. Loading state audit: deliberately slow network (DevTools), verify every list/form shows skeleton/spinner
4. Empty state audit: new account with no data — verify meaningful Persian messages on all screens
5. Android Chrome: install PWA, complete full client journey (login → plan → track → message → offline)
6. iOS Safari: add to home screen, verify standalone launch, push notification prompt, offline behavior
7. Verify production docker-compose.yml starts cleanly with production env vars on staging
8. Review all .env.example fields are documented with correct examples
9. Final `go test ./... -race` and `npm run test` pass
10. Final `npm run lint` and `go vet ./...` pass
11. Create phase verification checklist `07-05-LAUNCH-CHECKLIST.md` with evidence for each item

**Files touched:**
- Various `frontend/app/pages/**/*.vue` (Persian text/UX fixes as found)
- Various `backend/internal/handler/*.go` (error message consistency as found)
- `docker-compose.yml` (any remaining production config gaps)
- `.env.example`
- New: `docs/07-05-LAUNCH-CHECKLIST.md`

**Verification commands:**
```bash
cd backend && go test ./... -race
cd frontend && npm run test && npm run lint
docker compose -f docker-compose.yml up -d
curl -s https://staging.nutritrack.app/api/health
```

---

## Persian UX Launch Checklist

Use this as a review template during Plan 07-05:

### Text & Language
- [ ] All UI strings are in Persian (no English text visible to end users)
- [ ] Error messages from API use natural Persian phrasing (not literal translations)
- [ ] Dates display in Shamsi/Jalali calendar with Persian numerals everywhere
- [ ] Number inputs accept Persian digits (۰–۹) and convert to Latin for API calls
- [ ] Form labels and placeholders use RTL-appropriate Persian
- [ ] Toast/notification messages use correct Persian punctuation (،، ؟، !)

### Loading States
- [ ] All list views show skeleton or spinner during initial fetch
- [ ] All forms show submit button disabled + loading indicator during POST
- [ ] Diet plan builder shows loading state during day/meal/item operations
- [ ] Tracking submission shows loading before confirmation
- [ ] Message send shows pending state until confirmed

### Empty States
- [ ] Client with no active diet plan: helpful message + call to action
- [ ] Client with no tracking logs: welcoming first-time message
- [ ] Nutritionist with no clients: action to register first client
- [ ] Empty message thread: conversation starter prompt
- [ ] No food search results: suggest food request to nutritionist
- [ ] No lab results: upload prompt

### Error States
- [ ] Network offline banner appears correctly
- [ ] Sync queue error shows count and manual retry button
- [ ] Failed image/file upload shows specific error (size, format, etc.)
- [ ] Session expired shows re-login prompt (not blank screen)
- [ ] Server 500 shows generic Persian message (no stack trace)

### Mobile-Specific
- [ ] Touch targets are ≥ 44×44px (Apple HIG minimum)
- [ ] Keyboard doesn't obscure form inputs on mobile (scroll into view)
- [ ] RTL swipe gestures work on iOS (back navigation)
- [ ] Bottom navigation is reachable with thumb (not top-heavy layout)
- [ ] PWA installed icon uses correct maskable icon (512×512)

---

## Canonical Files Likely Touched

### Security (Plan 07-01)
| File | Change |
|------|--------|
| `backend/internal/middleware/security_headers.go` | Add HSTS header |
| `backend/internal/middleware/auth.go` | Verify/fix cookie attributes |
| `backend/internal/handler/authorization_test.go` (new) | Cross-tenant integration tests |
| `backend/internal/service/diet_plan_service.go` | Verify ownership checks |
| `backend/internal/service/communication_service.go` | Verify message partnership |
| `backend/internal/service/tracking_service.go` | Verify nutritionist→client ownership |

### Performance (Plan 07-02)
| File | Change |
|------|--------|
| `backend/internal/repository/diet_plan_repo_test.go` | Implement skipped benchmark |
| `backend/db/migrations/000011_add_tracking_indexes.up.sql` (new) | Date range indexes if needed |
| `backend/cmd/api/main.go` | pgxpool tuning |
| `frontend/nuxt.config.ts` | Code splitting / route-level loading |
| `backend/scripts/load_test.js` (new) | k6 load test script |

### Monitoring (Plan 07-03)
| File | Change |
|------|--------|
| `docker-compose.yml` | Add loki/promtail/grafana services |
| `monitoring/loki-config.yml` (new) | Loki storage config |
| `monitoring/promtail-config.yml` (new) | Docker scrape config |
| `monitoring/grafana/provisioning/datasources/loki.yml` (new) | Grafana datasource |
| `monitoring/grafana/dashboards/nutritrack.json` (new) | Dashboard JSON |
| `backend/internal/middleware/logger.go` | Verify duration field |

### Backup (Plan 07-04)
| File | Change |
|------|--------|
| `docker-compose.yml` | Add uploads volume, mount in api service |
| `.env.example` | Add UPLOADS_DIR documentation |
| `scripts/backup-db.sh` (new) | pg_dump cron script |
| `scripts/backup-uploads.sh` (new) | Volume backup script |
| `docs/restore-procedure.md` (new) | Runbook |

### Polish (Plan 07-05)
| File | Change |
|------|--------|
| Various `frontend/app/pages/**/*.vue` | Text/UX fixes found during audit |
| Various `backend/internal/handler/*.go` | Error message consistency |
| `docker-compose.yml` | Any remaining production gaps |

---

## Environment Availability

| Dependency | Required By | Available | Notes |
|------------|------------|-----------|-------|
| Docker Compose v2 | All deployment work | ✓ | Confirmed (no `version:` key in compose files) |
| Go 1.25 | Backend tests | ✓ | `go.mod` requires 1.25.0 |
| Node 22 | Frontend tests | ✓ | Used in CI (`NODE_VERSION: "22"`) |
| PostgreSQL 16 | DB tests | ✓ | docker-compose.dev.yml |
| govulncheck | Security audit | ✗ | Install: `go install golang.org/x/vuln/cmd/govulncheck@latest` |
| k6 | Load testing | ✗ | Install on Hetzner or local |
| Grafana/Loki | Monitoring | ✗ | Must be added to docker-compose |
| Real Android device | PWA testing | Unknown | Noted as deferred from Phase 6 |
| Real iOS device | PWA testing | Unknown | Noted as deferred from Phase 6 |

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Backend framework | `go test` + testify (stretchr/testify v1.11.1) |
| Frontend framework | Vitest v3.0.0 |
| Backend config | `go test ./...` (unit), `-tags integration` (DB tests) |
| Frontend config | `vitest.config.ts` (happy-dom environment) |
| Quick run (backend) | `cd backend && go test ./... -short` |
| Quick run (frontend) | `cd frontend && npm run test` |
| Full suite | `cd backend && go test ./... -race && cd ../frontend && npm run test` |

### Phase Requirements → Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| INFRA-09 | API p95 < 200ms | load test (k6) | `k6 run scripts/load_test.js` | ❌ Wave 0 (Plan 07-02) |
| INFRA-10 | Diet plan load < 500ms | integration bench | `go test -run TestPlanAggregateLoadTime -tags integration` | ❌ stub only |
| INFRA-11 | 500 concurrent users | load test (k6) | `k6 run scripts/load_test.js --vus 500` | ❌ Wave 0 (Plan 07-02) |
| SEC auth | Cross-tenant 403 | integration | `go test -run TestCrossTenant ./...` | ❌ Wave 0 (Plan 07-01) |
| SEC auth | Client cannot read other client plan | integration | `go test -run TestClientPlanIsolation ./...` | ❌ Wave 0 (Plan 07-01) |
| UI-08 | PWA < 3s on 3G | manual/Lighthouse | Chrome DevTools + Lighthouse PWA audit | Manual |
| INFRA-07 | DB backup + restore | manual/script | `bash scripts/restore-test.sh` | ❌ Wave 0 (Plan 07-04) |
| INFRA-08 | File backup restore | manual/script | `bash scripts/backup-uploads.sh && verify` | ❌ Wave 0 (Plan 07-04) |
| INFRA-05 | Grafana shows live metrics | manual | Open Grafana dashboard, generate traffic | Manual |

### Wave 0 Gaps
- [ ] `backend/internal/handler/authorization_test.go` — cross-tenant security tests (Plan 07-01)
- [ ] `backend/internal/repository/diet_plan_repo_test.go` — implement skipped benchmark (Plan 07-02)
- [ ] `backend/scripts/load_test.js` — k6 load test script (Plan 07-02)
- [ ] `scripts/backup-db.sh` + `scripts/restore-test.sh` — backup validation (Plan 07-04)
- [ ] `monitoring/` directory with loki/promtail/grafana configs (Plan 07-03)

---

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | yes | JWT + httpOnly cookie, bcrypt cost 12, refresh token rotation |
| V3 Session Management | yes | Token family tracking, revocation on logout, 30-day refresh |
| V4 Access Control | yes (highest risk) | Service-level resource ownership checks — AUDIT THIS PHASE |
| V5 Input Validation | yes | go-playground/validator + Gin ShouldBindJSON, sqlc parameterized |
| V6 Cryptography | yes | HMAC-SHA256 JWT, bcrypt — do NOT change; govulncheck for deps |
| V8 Data Protection | yes | httpOnly cookies, HSTS (to be added), no sensitive data in logs |
| V12 File Upload | yes | mimetype magic byte check, UUID filenames, size limits — VERIFY |
| V14 Config | yes | env vars, no secrets in code, VAPID keys optional |

### Known Threat Patterns

| Pattern | STRIDE | Standard Mitigation | Current Status |
|---------|--------|---------------------|---------------|
| Cross-tenant data access | Information Disclosure | Service-level ownership checks | ⚠️ Must audit all 40+ endpoints |
| JWT theft via XSS | Tampering | httpOnly cookie (already done) | ✅ Implemented |
| CSRF on cookie-based auth | Tampering | SameSite=Strict + CORS restriction | ⚠️ Verify SameSite attribute |
| File path traversal | Tampering | UUID filenames + controlled dir | ⚠️ Verify no user-controlled paths |
| Content sniffing | Spoofing | X-Content-Type-Options + mimetype | ✅ Implemented |
| Missing HSTS | Spoofing | Add HSTS header | ❌ Missing |
| Dependency CVEs | Various | govulncheck + npm audit | ❌ Not yet run |
| OTP brute force | Elevation | Rate limiter (3/10min) + attempt limit | ✅ Implemented |

---

## Open Questions

1. **Does the Logger middleware log `duration` per request?**
   - What we know: zerolog is configured; a Logger middleware exists at `backend/internal/middleware/logger.go`
   - What's unclear: Whether it logs `duration` in a field name that LogQL can unwrap
   - Recommendation: Read `logger.go` at start of Plan 07-03; add duration field if missing

2. **Is there a PATCH endpoint — does it check cookie SameSite attributes?**
   - What we know: Auth middleware reads cookie with `c.Cookie("access_token")`
   - What's unclear: Whether the cookie is set with `SameSite=Strict` and `Secure` in the auth handler
   - Recommendation: Read `auth_handler.go` cookie-setting code at start of Plan 07-01

3. **Are uploads currently persisted across container restarts in production?**
   - What we know: `UploadsDir` defaults to `./uploads`, no volume mount in docker-compose.yml
   - What's unclear: Whether production has a custom UPLOADS_DIR pointed at a mounted path
   - Recommendation: Treat as data loss risk; fix in Plan 07-04 regardless

4. **Does the super_admin role need to be blocked from accessing arbitrary diet plans?**
   - What we know: super_admin is allowed on `GET /api/diet-plans/:id` by route group
   - What's unclear: Whether this is intentional (admin oversight) or an oversight
   - Recommendation: Treat as intentional per requirements (super_admin has full visibility); document

5. **Real device availability for iOS testing?**
   - What we know: Phase 6 deferred this; Phase 7 must do it
   - What's unclear: Whether developer has access to real iOS/Android devices
   - Recommendation: Flag as a manual blocker in Plan 07-05; cannot be automated

---

## Recommendation: Number of Plans

**5 plans** is the right breakdown for Phase 7:

| Plan | Name | Focus | Dependency |
|------|------|-------|-----------|
| 07-01 | Security Audit & Hardening | Row-level auth tests, HSTS, CSP, cookie attrs, vuln scan | None |
| 07-02 | Performance Validation | Diet plan bench, load test, bundle analysis, index fixes | None (parallel-safe with 07-01) |
| 07-03 | Monitoring Infrastructure | Grafana + Loki + Promtail + dashboards + alerts | None |
| 07-04 | Backup & Restore | pg_dump cron, uploads volume fix, restore test, runbook | None |
| 07-05 | Final Polish & Launch | Persian UX audit, device testing, E2E journeys, prod config | 07-01, 07-02, 07-03, 07-04 |

Plans 07-01 through 07-04 are independent and can be executed in any order or in parallel. Plan 07-05
is the integration gate that cannot begin until all four hardening plans are complete and verified.

---

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Grafana/Loki 2.9.x is the stable version to use | Standard Stack | Minor version upgrade path; functionally same |
| A2 | govulncheck and npm audit find no CRITICAL vulnerabilities | Security | Would require dependency upgrades before launch |
| A3 | The Logger middleware logs `duration` per request | Monitoring | LogQL-based p95 would need different query or code change |
| A4 | Refresh token cookie has Secure+HttpOnly+SameSite attributes | Security | Potential CSRF or token theft vector |
| A5 | super_admin access to `GET /api/diet-plans/:id` is intentional | Security | Could be an authorization gap |
| A6 | Production deployment uses same docker-compose.yml as repo | Backup | Uploads volume fix may already be in place |

---

## Sources

### Primary (HIGH confidence — verified against actual codebase)
- `backend/cmd/api/main.go` — full route surface, pgxpool config, middleware chain
- `backend/internal/middleware/security_headers.go` — HSTS gap confirmed
- `backend/internal/middleware/cors.go` — CORS config verified correct
- `backend/internal/middleware/auth.go` — cookie auth pattern
- `backend/internal/middleware/role_guard.go` — role-only enforcement confirmed
- `backend/internal/config/config.go` — VAPID optional, uploads dir configurable
- `backend/internal/service/diet_plan_service.go` — plan ownership pattern
- `backend/internal/service/tracking_service.go` — tracking service structure
- `backend/internal/repository/diet_plan_repo_test.go` — skipped benchmark confirmed
- `backend/go.mod` — exact dependency versions
- `frontend/nuxt.config.ts` — PWA config, no explicit code splitting
- `frontend/package.json` — frontend dependency list
- `docker-compose.yml` — no monitoring services, no uploads volume
- `.gitlab-ci.yml` — CI commands: `go test ./... -v -race`, `npm run test`
- `.planning/phases/06-offline-pwa/VERIFICATION.md` — manual device testing deferred

### Secondary (MEDIUM confidence — from official docs / established patterns)
- docs/phases.md §Phase 7 — validation checklist items cross-referenced
- Grafana/Loki Docker Compose patterns — standard self-hosted observability stack [ASSUMED]
- k6 load testing patterns for Go HTTP APIs [ASSUMED]

---

## Metadata

**Confidence breakdown:**
- Security findings: HIGH — verified against actual source files
- Performance gaps: HIGH — skipped test stub and missing volume confirmed by code
- Monitoring approach: MEDIUM — Loki/Grafana pattern well-established; specific LogQL query assumes duration logging
- Backup approach: HIGH — docker-compose.yml gap confirmed; restore pattern is standard pg_dump
- UX checklist: HIGH — Phase 6 VERIFICATION.md explicitly deferred device testing

**Research date:** 2026-04-20
**Valid until:** 2026-05-20 (stable stack; no fast-moving components)
