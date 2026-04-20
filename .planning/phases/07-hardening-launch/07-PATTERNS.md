# Phase 7: Hardening & Launch — Pattern Map

**Mapped:** 2026-04-20
**Files analyzed:** 22 new/modified files across 5 plans
**Analogs found:** 18 / 22 (4 have no codebase analog — covered below)

---

## File Classification

| New / Modified File | Role | Data Flow | Closest Analog | Match Quality |
|---------------------|------|-----------|----------------|---------------|
| `backend/internal/middleware/security_headers.go` | middleware | request-response | itself | exact (1-line addition) |
| `backend/internal/middleware/auth.go` | middleware | request-response | itself | exact (cookie attribute audit) |
| `backend/internal/handler/authorization_test.go` (new) | test | request-response | `backend/internal/service/push_service_test.go` | role-match |
| `backend/internal/repository/diet_plan_repo_test.go` | test | CRUD / integration | `backend/internal/repository/tracking_repo_test.go` | exact |
| `backend/internal/service/diet_plan_service.go` | service | CRUD | `backend/internal/handler/diet_plan_handler.go` | role-match (ownership check pattern) |
| `backend/internal/service/communication_service.go` | service | event-driven | `backend/internal/handler/diet_plan_handler.go` | partial-match |
| `backend/scripts/load_test.js` (new) | utility | request-response | none | no analog |
| `docker-compose.yml` | config | — | itself + `docker-compose.dev.yml` | exact |
| `monitoring/loki-config.yml` (new) | config | — | `docker-compose.dev.yml` service blocks | partial |
| `monitoring/promtail-config.yml` (new) | config | — | none in codebase | no analog |
| `monitoring/grafana/provisioning/datasources/loki.yml` (new) | config | — | none in codebase | no analog |
| `monitoring/grafana/dashboards/nutritrack.json` (new) | config | — | none in codebase | no analog |
| `scripts/backup-db.sh` (new) | utility | file-I/O | none in codebase | no analog |
| `scripts/backup-uploads.sh` (new) | utility | file-I/O | `scripts/backup-db.sh` (peer) | peer |
| `docs/restore-procedure.md` (new) | doc | — | `.planning/ROADMAP.md` (structure) | partial |
| `backend/db/migrations/000011_add_tracking_indexes.up.sql` (new, if needed) | migration | CRUD | `backend/db/migrations/000009_create_communication.up.sql` | exact |
| `frontend/app/pages/client/messages.vue` | page/component | event-driven | `frontend/app/pages/client/plan.vue` | role-match |
| `frontend/app/pages/client/tracking/*.vue` (7 files) | page/component | CRUD | `frontend/app/pages/client/plan.vue` | role-match |
| `frontend/app/pages/client/food-requests.vue` | page/component | CRUD | `frontend/app/pages/client/plan.vue` | role-match |
| `frontend/app/pages/nutritionist/clients.vue` | page/component | CRUD | `frontend/app/pages/client/plan.vue` | role-match |
| `frontend/app/pages/nutritionist/clients/[clientId]/index.vue` | page/component | CRUD | `frontend/app/pages/client/plan.vue` | role-match |
| `frontend/app/layouts/client.vue` | layout | request-response | `frontend/app/pages/client/messages.vue` | partial |

---

## Pattern Assignments

### `backend/internal/middleware/security_headers.go` (middleware, request-response)

**Analog:** itself (`backend/internal/middleware/security_headers.go`)
**Why best:** This is a 1-line addition to an existing file; no new structure needed.

**Existing file — full content** (lines 1–16):
```go
package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adds standard security headers to all responses.
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Next()
    }
}
```

**Required change — add HSTS after the CSP line:**
```go
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Next()
```

**Analog for the surrounding middleware pattern:** `backend/internal/middleware/role_guard.go`
```go
// RoleGuard checks that the authenticated user's role (set by Auth middleware)
// matches one of the allowed roles. Returns 403 if unauthorized.
func RoleGuard(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        for _, allowed := range allowedRoles {
            if role == allowed {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "دسترسی غیرمجاز"})
    }
}
```

---

### `backend/internal/middleware/auth.go` (middleware, request-response — audit)

**Analog:** itself (`backend/internal/middleware/auth.go`)
**Why best:** Audit target, not a new file. Verify `setAuthCookies` in `auth_handler.go`.

**Cookie-setting pattern** (`backend/internal/handler/auth_handler.go` lines 131–143):
```go
// setAuthCookies sets httpOnly secure cookies for access and refresh tokens (D-01).
// access_token: path=/api, maxAge=900 (15min), secure=true, httpOnly=true
// refresh_token: path=/api/auth/refresh, maxAge=2592000 (30d), secure=true, httpOnly=true
func setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
    c.SetCookie("access_token", accessToken, 900, "/api", "", true, true)
    c.SetCookie("refresh_token", refreshToken, 2592000, "/api/auth/refresh", "", true, true)
}
```

**Gin's `SetCookie` signature:** `(name, value string, maxAge int, path, domain string, secure, httpOnly bool)`
- `secure=true` (5th positional `true`) ✅ — already correct
- `httpOnly=true` (6th positional `true`) ✅ — already correct
- `SameSite` is **NOT** set via `SetCookie` — requires `c.SetSameSite(http.SameSiteStrictMode)` **before** calling `SetCookie`

**Fix pattern:**
```go
func setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
    c.SetSameSite(http.SameSiteStrictMode) // add this line
    c.SetCookie("access_token", accessToken, 900, "/api", "", true, true)
    c.SetCookie("refresh_token", refreshToken, 2592000, "/api/auth/refresh", "", true, true)
}
```

---

### `backend/internal/handler/authorization_test.go` (test, integration — new)

**Analog:** `backend/internal/service/push_service_test.go`
**Why best:** The only concrete, fully-implemented test file in the backend. Shows the project's mock pattern, `testify/mock`, `testify/assert`, package naming convention (`package service_test`), and how to wire a service with a mock repo.

**Imports pattern** (`push_service_test.go` lines 1–16):
```go
package service_test

import (
    "context"
    "testing"
    "time"

    "github.com/rs/zerolog"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/ranjbar-dev/nutritrack/backend/internal/config"
    "github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
    "github.com/ranjbar-dev/nutritrack/backend/internal/repository"
    "github.com/ranjbar-dev/nutritrack/backend/internal/service"
)
```

**Mock repository pattern** (`push_service_test.go` lines 19–62):
```go
type mockPushRepo struct{ mock.Mock }

func (m *mockPushRepo) MethodName(ctx context.Context, ...) (ReturnType, error) {
    args := m.Called(ctx, ...)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(ReturnType), args.Error(1)
}
```

**Test function pattern** (`push_service_test.go` lines 73–97):
```go
func TestSendToClient_SkipsWhenPreferenceDisabled(t *testing.T) {
    repo := &mockPushRepo{}
    repo.On("MethodName", mock.Anything, "param-value").Return(
        &repository.SomeType{Field: value},
        nil,
    )
    svc := makeService(repo)
    err := svc.MethodUnderTest(context.Background(), "param-value", ...)
    assert.NoError(t, err)
    repo.AssertNotCalled(t, "UnexpectedMethod")
}
```

**Integration test build tag pattern** (`diet_plan_repo_test.go` line 1):
```go
//go:build integration
```
> Use `//go:build integration` at the top of all cross-tenant handler tests that require a real DB. Run with `go test ./... -tags integration`.

**New file skeleton for `authorization_test.go`:**
```go
//go:build integration

package handler_test

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/ranjbar-dev/nutritrack/backend/internal/handler"
    "github.com/ranjbar-dev/nutritrack/backend/internal/service"
)

// TestCrossTenantAccessDenied — seeds two nutritionists, verifies nutritionist-B
// cannot access nutritionist-A's client's resources. Expects 403 or 404.
func TestCrossTenantAccessDenied(t *testing.T) {
    // Seed two nutritionists and their clients
    // Use nutritionist-B's token to call endpoints for nutritionist-A's clients
    // Assert http.StatusForbidden or http.StatusNotFound
}

// TestClientCannotAccessOtherClientPlan — seeds two clients, verifies client-A
// cannot call GET /api/diet-plans/{planB_id}.
func TestClientCannotAccessOtherClientPlan(t *testing.T) {
    // ...
}
```

---

### `backend/internal/repository/diet_plan_repo_test.go` (test, integration — implement stub)

**Analog:** `backend/internal/repository/tracking_repo_test.go`
**Why best:** Same package (`repository_test`), same build tag (`//go:build integration`), same stub shape — the tracking test is the direct peer in the same directory.

**Current stub** (lines 1–22, entire file):
```go
//go:build integration

package repository_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestPlanAggregateLoadTime — Dim 2: ≤500ms SLA for 7×5×3×4 plan
func TestPlanAggregateLoadTime(t *testing.T) {
    t.Skip("stub — Plan 03-02 implements")
    _ = context.Background()
    _ = time.Millisecond
    _ = require.New(t)
    _ = assert.New(t)
}
```

**Implementation target pattern** — remove `t.Skip`, seed realistic data (7 days × 5 meals × 3 options × 4 items = 420 items), call repo method, assert duration < 500ms:
```go
func TestPlanAggregateLoadTime(t *testing.T) {
    ctx := context.Background()
    require := require.New(t)
    assert := assert.New(t)

    // Seed plan with 7 days × 5 meals × 3 options × 4 items
    planID := seedRealisticPlan(t, ctx)

    start := time.Now()
    _, err := repo.GetPlanAggregate(ctx, planID)
    elapsed := time.Since(start)

    require.NoError(err)
    assert.Less(elapsed, 500*time.Millisecond, "diet plan aggregate must load in < 500ms")
}
```

---

### `backend/internal/service/diet_plan_service.go` (service, CRUD — ownership audit)

**Analog:** `backend/internal/handler/diet_plan_handler.go`
**Why best:** The handler already demonstrates the ownership threading pattern: `userID` is extracted from JWT context and passed to service methods alongside `planID`. Auditing that the service actually enforces this is the goal.

**Ownership pattern in handler** (`diet_plan_handler.go` lines 56–83):
```go
func (h *DietPlanHandler) GetPlanAggregate(c *gin.Context) {
    planID, err := uuid.Parse(c.Param("id"))
    // ...
    userID, err := uuid.Parse(c.GetString("user_id"))
    // ...
    role := c.GetString("role")

    var resp *dto.DietPlanResponse
    if role == "client" {
        resp, err = h.planService.GetPlanAggregateForClient(c.Request.Context(), planID, userID)
    } else {
        resp, err = h.planService.GetPlanAggregate(c.Request.Context(), planID, userID)
    }
    if err != nil {
        h.handlePlanError(c, err)
        return
    }
    c.JSON(http.StatusOK, resp)
}
```

**What to verify in the service:** Every `GetPlanAggregate`, `UpdatePlanHeader`, `DeletePlan`, and sub-resource mutation must check that the calling `userID` owns the resource. The service must return a sentinel error (`service.ErrForbidden` or `service.ErrNotFound`) when ownership fails.

**Error response pattern** (`diet_plan_handler.go` lines 43–52):
```go
resp, err := h.planService.CreatePlan(c.Request.Context(), nutritionistID, req)
if err != nil {
    switch {
    case errors.Is(err, service.ErrPlanInvalidDateRange):
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
    default:
        h.handlePlanError(c, err)
    }
    return
}
```

---

### `backend/db/migrations/000011_add_tracking_indexes.up.sql` (migration, CRUD — new if needed)

**Analog:** `backend/db/migrations/000009_create_communication.up.sql`
**Why best:** Most recent migration; shows naming convention, `CREATE INDEX` pattern, and partial index syntax used in this project.

**Index naming convention** (lines 18–21):
```sql
CREATE INDEX idx_messages_sender_sent   ON messages (sender_id, sent_at);
CREATE INDEX idx_messages_receiver_sent ON messages (receiver_id, sent_at);
CREATE INDEX idx_messages_conversation  ON messages (LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id), sent_at);
CREATE INDEX idx_messages_unread        ON messages (receiver_id, read_at) WHERE read_at IS NULL;
```

**Pattern:** `idx_{table}_{column_description}` — composite columns joined by `_`; partial indexes use `WHERE`.

**Target index for tracking performance:**
```sql
-- Migration 000011: Add tracking performance indexes
-- Composite indexes for (client_id, date) range queries on all tracking tables

CREATE INDEX idx_food_logs_client_date    ON food_logs    (client_id, date DESC);
CREATE INDEX idx_water_logs_client_date   ON water_logs   (client_id, date DESC);
CREATE INDEX idx_sleep_logs_client_date   ON sleep_logs   (client_id, date DESC);
CREATE INDEX idx_exercise_logs_client_date ON exercise_logs (client_id, date DESC);
CREATE INDEX idx_medication_logs_client_date ON medication_logs (client_id, date DESC);
```

---

### `docker-compose.yml` (config — monitoring + uploads volume additions)

**Analog:** itself + `docker-compose.dev.yml`
**Why best:** The file is being modified; `docker-compose.dev.yml` shows healthcheck pattern and service structure used in this project.

**Existing service structure** (production `docker-compose.yml` lines 22–36):
```yaml
api:
  build:
    context: ./backend
    dockerfile: Dockerfile
  env_file: .env
  labels:
    - "traefik.enable=true"
    - "traefik.http.routers.api.rule=Host(`${DOMAIN}`) && PathPrefix(`/api`)"
    - "traefik.http.routers.api.entrypoints=websecure"
    - "traefik.http.routers.api.tls.certresolver=letsencrypt"
    - "traefik.http.services.api.loadbalancer.server.port=8080"
  depends_on:
    postgres:
      condition: service_healthy
  restart: unless-stopped
```

**Uploads volume addition — copy this pattern for the `api` service volumes block:**
```yaml
api:
  # ... existing keys ...
  volumes:
    - uploads:/app/uploads
  environment:
    UPLOADS_DIR: /app/uploads
```

**New monitoring services to append — copy `restart: unless-stopped` from every existing service:**
```yaml
  loki:
    image: grafana/loki:2.9.8
    command: -config.file=/etc/loki/local-config.yaml
    volumes:
      - ./monitoring/loki-config.yml:/etc/loki/local-config.yaml:ro
      - loki_data:/loki
    restart: unless-stopped

  promtail:
    image: grafana/promtail:2.9.8
    command: -config.file=/etc/promtail/config.yml
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./monitoring/promtail-config.yml:/etc/promtail/config.yml:ro
    depends_on:
      - loki
    restart: unless-stopped

  grafana:
    image: grafana/grafana:10.4.2
    environment:
      GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_ADMIN_PASSWORD}
    volumes:
      - grafana_data:/var/lib/grafana
      - ./monitoring/grafana/provisioning:/etc/grafana/provisioning:ro
    ports:
      - "127.0.0.1:3001:3000"   # bind to localhost ONLY — not exposed via Traefik
    depends_on:
      - loki
    restart: unless-stopped
```

**Volumes block addition:**
```yaml
volumes:
  pgdata:
  letsencrypt:
  uploads:       # add
  loki_data:     # add
  grafana_data:  # add
```

---

### `monitoring/loki-config.yml` (config — new)

**Analog:** RESEARCH.md §Backup Pattern (same YAML config style)
**No codebase analog.** Copy from RESEARCH.md §Monitoring Architecture Pattern. Use pinned version `grafana/loki:2.9.8` matching promtail.

**Canonical config:**
```yaml
auth_enabled: false

server:
  http_listen_port: 3100

ingester:
  lifecycler:
    address: 127.0.0.1
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
    final_sleep: 0s
  chunk_idle_period: 5m
  chunk_retain_period: 30s

schema_config:
  configs:
    - from: 2024-01-01
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: /loki/index
    cache_location: /loki/cache
    shared_store: filesystem
  filesystem:
    directory: /loki/chunks

limits_config:
  retention_period: 720h    # 30 days

chunk_store_config:
  max_look_back_period: 720h

table_manager:
  retention_deletes_enabled: true
  retention_period: 720h
```

---

### `monitoring/promtail-config.yml` (config — new)

**No codebase analog.** Must be written from scratch using standard Promtail Docker socket scrape pattern.

**Canonical config (pin version to match loki:2.9.8):**
```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: docker
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
    relabel_configs:
      - source_labels: [__meta_docker_container_name]
        target_label: container
      - source_labels: [__meta_docker_container_label_com_docker_compose_service]
        target_label: service
    pipeline_stages:
      - json:
          expressions:
            level: level
            request_id: request_id
            path: path
            status: status
            duration_ms: duration_ms
      - labels:
          level:
          path:
```

> **Key:** The `duration_ms` field from `backend/internal/middleware/logger.go` (line 26: `.Dur("duration_ms", duration)`) is what feeds the LogQL p95 calculation. No code changes to logger.go are needed.

---

### `monitoring/grafana/provisioning/datasources/loki.yml` (config — new)

**No codebase analog.**

```yaml
apiVersion: 1

datasources:
  - name: Loki
    type: loki
    access: proxy
    url: http://loki:3100
    isDefault: true
    version: 1
    editable: false
```

---

### `scripts/backup-db.sh` (utility, file-I/O — new)

**No codebase analog.** Copy shell script pattern from RESEARCH.md §Backup Pattern.

**Canonical script:**
```bash
#!/bin/bash
# backup-db.sh — Daily PostgreSQL backup via docker exec
# Run from cron: 0 2 * * * root /path/to/scripts/backup-db.sh

set -euo pipefail

BACKUP_DIR="/backups/db"
DATE=$(date +%Y-%m-%d)
CONTAINER="nutritrack-postgres-1"
PGUSER="${POSTGRES_USER:-nutritrack}"
PGDB="${POSTGRES_DB:-nutritrack}"
RETENTION_DAYS=30

mkdir -p "$BACKUP_DIR"

docker exec "$CONTAINER" pg_dump -U "$PGUSER" "$PGDB" \
  | gzip > "${BACKUP_DIR}/${DATE}.sql.gz"

echo "Backup completed: ${BACKUP_DIR}/${DATE}.sql.gz"

# Prune backups older than retention window
find "$BACKUP_DIR" -name "*.sql.gz" -mtime "+${RETENTION_DAYS}" -delete
```

---

### `frontend/app/pages/client/messages.vue` (page, event-driven — loading/empty/error hardening)

**Analog:** `frontend/app/pages/client/plan.vue`
**Why best:** `plan.vue` is the canonical model explicitly called out in the UI-SPEC: it has the full skeleton-pulse loading pattern, emoji empty-state pattern, and `pageLoading` ref lifecycle pattern. `messages.vue` currently uses bare text loading which must be replaced.

**Canonical loading skeleton pattern** (`client/plan.vue` lines 73–75):
```vue
<div v-if="pageLoading" class="space-y-3 p-4">
  <div v-for="i in 4" :key="i" class="h-24 animate-pulse rounded-2xl bg-white shadow-sm" />
</div>
```
> Replace `messages.vue` line 88: `<div v-if="messageStore.loading" class="text-center text-gray-400 text-sm">در حال بارگذاری...</div>` with skeleton pattern (3 skeleton rows for chat messages).

**Canonical pageLoading lifecycle pattern** (`client/plan.vue` lines 17, 32–44):
```vue
<script setup lang="ts">
const pageLoading = ref(true)

onMounted(async () => {
  pageLoading.value = true
  await Promise.all([
    store.fetchActivePlan(),
    store.fetchMyPlans(),
  ])
  pageLoading.value = false
})
</script>
```

**Canonical empty-state pattern** (`client/plan.vue` lines 181–187):
```vue
<div v-else class="px-4 py-16 text-center">
  <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 text-2xl">
    🍽️
  </div>
  <h2 class="mt-4 font-bold text-gray-800">برنامه‌ای فعال ندارید</h2>
  <p class="mt-2 text-sm text-gray-500">برای دریافت برنامه جدید با کارشناس تغذیه خود در تماس باشید.</p>
</div>
```
> Apply this pattern with appropriate emoji + copy for each screen. Refer to UI-SPEC §Empty States table for the exact Persian copy per screen.

**Fetch-error state pattern (must add to messages.vue and all ❓ screens):**
```vue
<script setup lang="ts">
const fetchError = ref<string | null>(null)

onMounted(async () => {
  pageLoading.value = true
  fetchError.value = null
  try {
    await store.fetchData()
  } catch {
    fetchError.value = 'خطا در بارگذاری اطلاعات. دوباره تلاش کنید.'
  } finally {
    pageLoading.value = false
  }
})
</script>

<template>
  <!-- after skeleton v-if block -->
  <div v-else-if="fetchError" class="px-4 py-8 text-center">
    <p class="text-sm text-red-600">{{ fetchError }}</p>
    <button
      class="mt-3 rounded-lg bg-emerald-500 px-4 py-2 text-sm text-white"
      @click="retryFetch"
    >
      تلاش مجدد
    </button>
  </div>
</template>
```

**Send button touch-target fix** (`messages.vue` line 133 — `px-3 py-1` → `px-4 py-2.5`):
```vue
<!-- BEFORE -->
<button class="bg-emerald-500 text-white rounded-lg px-3 py-1 text-sm disabled:opacity-50" ...>
  ارسال
</button>

<!-- AFTER (44px touch target) -->
<button class="bg-emerald-500 text-white rounded-lg px-4 py-2.5 text-sm min-h-[44px] disabled:opacity-50" ...>
  ارسال
</button>
```

**Message error text pattern** (`messages.vue` line 111 — already correct, keep):
```vue
<div v-if="errorMsg" class="text-red-500 text-sm">{{ errorMsg }}</div>
```

---

### `frontend/app/pages/client/tracking/*.vue` (7 pages, CRUD — full state audit)

**Analog:** `frontend/app/pages/client/plan.vue`
**Why best:** Same role (client page), same data flow (fetch on mount, display list), same layout conventions (`pb-24`, `bg-gray-50`, sticky header).

**Copy these patterns from `plan.vue` to every tracking sub-page:**

1. **Skeleton loading** (lines 73–75) — replace any bare `در حال بارگذاری...` with skeleton cards
2. **Fetch-error state** — add `fetchError` ref and retry button (see pattern above)
3. **Empty state** — add `v-else` with emoji icon + Persian heading + body using UI-SPEC copy
4. **`pageLoading` ref lifecycle** (lines 17, 32–44) — standard onMounted async pattern

**Tracking-specific empty state copy (from UI-SPEC §Empty States):**
```
Client: tracking (no entries today)
  heading: هنوز چیزی ثبت نشده
  body:    امروز اولین ثبت خود را انجام دهید.
```

---

### `frontend/app/pages/nutritionist/clients.vue` (page, CRUD — skeleton upgrade)

**Analog:** `frontend/app/pages/client/plan.vue`
**Why best:** Plan.vue canonical skeleton is what must replace the text-only loading.

**Current pattern** (`nutritionist/clients.vue` line 76 — must replace):
```vue
<div v-if="clientStore.loading" class="text-center text-gray-400 text-sm py-8">در حال بارگذاری...</div>
```

**Replacement — skeleton rows:**
```vue
<div v-if="clientStore.loading" class="space-y-3 p-4">
  <div v-for="i in 5" :key="i" class="h-16 animate-pulse rounded-2xl bg-white shadow-sm" />
</div>
```

**Empty state upgrade** (`nutritionist/clients.vue` line 77–79 — currently bare text):
```vue
<!-- BEFORE -->
<div v-else-if="clientStore.clients.length === 0" class="bg-white rounded-xl p-6 shadow-sm text-center text-gray-500 text-sm">
  هیچ مراجعی یافت نشد.
</div>

<!-- AFTER — canonical empty state pattern from plan.vue -->
<div v-else-if="clientStore.clients.length === 0" class="px-4 py-16 text-center">
  <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 text-2xl">
    👤
  </div>
  <h2 class="mt-4 font-bold text-gray-800">مراجعی یافت نشد</h2>
  <p class="mt-2 text-sm text-gray-500">جستجو را تغییر دهید یا مراجع جدیدی ثبت کنید.</p>
</div>
```

---

### `frontend/app/pages/client/food-requests.vue` (page, CRUD — full state audit)

**Analog:** `frontend/app/pages/client/plan.vue`

**Empty state copy (from UI-SPEC):**
```
heading: درخواستی ثبت نشده
body:    درخواست‌های شما برای افزودن غذا اینجا نمایش داده می‌شود.
```

Apply the same three patterns: skeleton loading, fetch-error + retry, empty state with emoji.

---

## Shared Patterns

### Loading Skeleton (apply to ALL pages with async fetch)

**Source:** `frontend/app/pages/client/plan.vue` lines 73–75
**Apply to:** `messages.vue`, all 7 `tracking/*.vue`, `food-requests.vue`, `nutritionist/clients.vue`, `nutritionist/clients/[clientId]/tracking/*`

```vue
<!-- Full-page skeleton: use 3–4 cards matching the content card height -->
<div v-if="pageLoading" class="space-y-3 p-4">
  <div v-for="i in 4" :key="i" class="h-24 animate-pulse rounded-2xl bg-white shadow-sm" />
</div>

<!-- List skeleton (for nutritionist pages): shorter rows -->
<div v-if="pageLoading" class="space-y-3 p-4">
  <div v-for="i in 5" :key="i" class="h-16 animate-pulse rounded-2xl bg-white shadow-sm" />
</div>
```

**Rule from UI-SPEC:** Skeleton pulse is for full-page loads. `LoadingSpinner.vue` is for button-level inline loading only.

---

### Empty State (apply to ALL list/content pages)

**Source:** `frontend/app/pages/client/plan.vue` lines 181–187
**Apply to:** Every `v-else` on a list that could be empty.

```vue
<div v-else class="px-4 py-16 text-center">
  <div class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-gray-100 text-2xl">
    {EMOJI}
  </div>
  <h2 class="mt-4 font-bold text-gray-800">{HEADING_FROM_UI_SPEC}</h2>
  <p class="mt-2 text-sm text-gray-500">{BODY_FROM_UI_SPEC}</p>
</div>
```

**Typography enforcement (UI-SPEC §Typography):**
- Empty-state heading: `font-bold text-gray-800` at `text-base` (16 px) ✅
- Empty-state body: `text-sm text-gray-500` ✅

---

### Fetch-Error State (apply to ALL pages with API calls)

**Source:** synthesized from `messages.vue` error pattern + UI-SPEC §Error States
**Apply to:** Every page with `onMounted` async fetch that doesn't already have a fetch-error handler.

```vue
<script setup lang="ts">
const fetchError = ref<string | null>(null)

onMounted(async () => {
  pageLoading.value = true
  fetchError.value = null
  try {
    await store.fetchData()
  } catch {
    fetchError.value = 'خطا در بارگذاری اطلاعات. دوباره تلاش کنید.'
  } finally {
    pageLoading.value = false
  }
})

async function retryFetch() {
  fetchError.value = null
  pageLoading.value = true
  try {
    await store.fetchData()
  } catch {
    fetchError.value = 'خطا در بارگذاری اطلاعات. دوباره تلاش کنید.'
  } finally {
    pageLoading.value = false
  }
}
</script>

<template>
  <div v-if="pageLoading" class="space-y-3 p-4">
    <div v-for="i in 4" :key="i" class="h-24 animate-pulse rounded-2xl bg-white shadow-sm" />
  </div>
  <div v-else-if="fetchError" class="px-4 py-8 text-center">
    <p class="text-sm text-red-600">{{ fetchError }}</p>
    <button
      class="mt-3 rounded-lg bg-emerald-500 px-4 py-2 text-sm text-white min-h-[44px]"
      @click="retryFetch"
    >
      تلاش مجدد
    </button>
  </div>
  <div v-else>
    <!-- content or empty state -->
  </div>
</template>
```

---

### Backend Handler Error Response (apply to ALL handler methods)

**Source:** `backend/internal/handler/diet_plan_handler.go` lines 29–52
**Source:** `backend/internal/handler/auth_handler.go` lines 26–39

```go
// JSON bind failure — always 400
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
    return
}

// UUID parse failure — always 401
userID, err := uuid.Parse(c.GetString("user_id"))
if err != nil {
    c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "توکن نامعتبر است"})
    return
}

// Service error — switch on sentinel errors, delegate unknown errors to handleXxxError()
if err != nil {
    switch {
    case errors.Is(err, service.ErrSomeSentinel):
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
    default:
        h.handlePlanError(c, err)
    }
    return
}
```

---

### Backend Integration Test Build Tag

**Source:** `backend/internal/repository/diet_plan_repo_test.go` line 1
**Apply to:** `authorization_test.go` and any other new test requiring a live DB.

```go
//go:build integration

package handler_test  // or repository_test / service_test
```

Run with: `go test ./... -tags integration -v`

---

### Go Module Import Path

**Source:** `backend/internal/service/push_service_test.go` lines 13–16
**Apply to:** All new Go files.

```go
import (
    "github.com/ranjbar-dev/nutritrack/backend/internal/config"
    "github.com/ranjbar-dev/nutritrack/backend/internal/model/dto"
    "github.com/ranjbar-dev/nutritrack/backend/internal/repository"
    "github.com/ranjbar-dev/nutritrack/backend/internal/service"
)
```

Module root: `github.com/ranjbar-dev/nutritrack/backend`

---

### Docker Compose Service Conventions

**Source:** `docker-compose.yml` lines 22–36 and `docker-compose.dev.yml` lines 1–39

- All services use `restart: unless-stopped` (production) — copy this on every new service
- Health checks: `pg_isready` pattern for postgres; monitoring services don't need healthchecks
- No `version:` key — already correct (Docker Compose v2 format)
- Named volumes: declare in top-level `volumes:` block with no options (filesystem default)
- Internal-only services (Grafana): bind port to `127.0.0.1:{host}:{container}` to prevent public exposure

---

### Persian Error Copy Rules (apply to ALL handler error messages)

**Source:** UI-SPEC §Error States, §Copywriting Contract
**Apply to:** All `dto.ErrorResponse{Error: "..."}` strings in handlers.

| Scenario | Required copy |
|----------|--------------|
| JSON bind error | `اطلاعات ورودی نامعتبر است` |
| Invalid UUID / token | `توکن نامعتبر است` |
| Not authenticated | `احراز هویت الزامی است` |
| Forbidden (wrong owner) | `دسترسی غیرمجاز` |
| Not found | `مورد یافت نشد` |
| Generic server error | `خطای داخلی سرور` |

**Rule:** No English text, raw error objects, or Go error strings must reach API JSON responses.

---

## No Analog Found

| File | Role | Data Flow | Reason |
|------|------|-----------|--------|
| `backend/scripts/load_test.js` | utility | request-response | No load test scripts exist; use k6 standard syntax from RESEARCH.md |
| `monitoring/promtail-config.yml` | config | — | No log shipping config exists; use Promtail Docker socket pattern from RESEARCH.md |
| `monitoring/grafana/provisioning/datasources/loki.yml` | config | — | No Grafana provisioning exists; use standard Grafana YAML provisioning format |
| `monitoring/grafana/dashboards/nutritrack.json` | config | — | No dashboards exist; generate via Grafana UI then export as JSON |

---

## Metadata

**Analog search scope:**
- `backend/internal/middleware/` — all 8 middleware files
- `backend/internal/handler/` — all 10 handler files (sampled diet_plan, auth)
- `backend/internal/service/` — all test files (push_service_test.go fully read)
- `backend/internal/repository/` — all test files
- `backend/db/migrations/` — last 3 migrations
- `backend/cmd/api/main.go` — full file
- `frontend/app/pages/client/plan.vue` — full file (canonical model)
- `frontend/app/pages/client/messages.vue` — full file
- `frontend/app/pages/nutritionist/clients.vue` — lines 1–80
- `docker-compose.yml` — full file
- `docker-compose.dev.yml` — full file

**Files scanned:** ~22 source files read
**Pattern extraction date:** 2026-04-20
