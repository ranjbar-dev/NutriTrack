---
phase: "01"
depth: standard
status: issues_found
files_reviewed: 15
files_reviewed_list:
  - go.mod
  - configs/config.go
  - configs/migrate.go
  - bootstrap/logger.go
  - bootstrap/database.go
  - bootstrap/redis.go
  - bootstrap/migrator.go
  - cmd/server/main.go
  - internal/domain/shared/apperror.go
  - internal/domain/shared/timeutil.go
  - internal/domain/shared/normalize.go
  - internal/interfaces/http/router/router.go
  - internal/interfaces/http/middleware/request_id.go
  - internal/interfaces/http/middleware/logger.go
  - internal/interfaces/http/middleware/recovery.go
  - internal/interfaces/http/middleware/not_found.go
  - Dockerfile
findings:
  critical: 2
  warning: 5
  info: 6
  total: 13
reviewed_at: "2026-04-22"
---

# Code Review: Phase 01 — Foundation

## Summary

The foundation layer is well-structured with clean separation between config, bootstrap, domain errors, and HTTP middleware. Two critical issues require immediate attention: a Go toolchain version mismatch that may silently break Docker builds, and an unsafe type assertion inside the panic-recovery middleware that would cause a secondary panic if the recovery handler is ever triggered before the RequestID middleware runs. Five warnings cover log injection via unvalidated headers, an uploads path mismatch in Docker, an insecure SSL default, and missing startup timeouts.

---

## Findings

### CR-001 — Go toolchain version mismatch: Dockerfile 1.24 vs go.mod 1.25.0

**File:** `Dockerfile` (line 1), `go.mod` (line 3)  
**Severity:** CRITICAL  
**Line:** ~1 / ~3  
**Issue:** `go.mod` declares `go 1.25.0` but the Docker builder stage pulls `golang:1.24-alpine`. With `GOTOOLCHAIN=auto` (the Go default since 1.21), the 1.24 toolchain will attempt to download Go 1.25 at build time. In a CI environment with no egress or with `GOTOOLCHAIN=off`, the build will fail outright with *"go.mod requires go >= 1.25"*. The mismatch is silent locally when Go 1.25 is installed on the host but breaks reproducible Docker builds.  
**Recommendation:** Pin the builder image to match go.mod exactly:
```dockerfile
FROM golang:1.25-alpine AS builder
```
Alternatively, add `GOTOOLCHAIN=local` as a build arg after confirming 1.24 is intentional and go.mod is downgraded.

---

### CR-002 — Unsafe type assertion in Recovery middleware — panic inside panic handler

**File:** `internal/interfaces/http/middleware/recovery.go` (line ~10–11)  
**Severity:** CRITICAL  
**Line:** ~10  
**Issue:** `requestID.(string)` is an unchecked type assertion. `c.Get(RequestIDKey)` returns `(any, bool)` — the bool is discarded. If the `RequestID()` middleware is ever absent from the chain, not yet executed, or if the value stored is not a `string`, this assertion panics inside `gin.CustomRecovery`. A panic inside the recovery handler is not re-caught by Gin, causing the goroutine to crash and the HTTP response to be dropped entirely — the exact scenario that recovery is meant to prevent.
```go
// Current — panics if key is missing or not a string
requestID, _ := c.Get(RequestIDKey)
log.Error().Str("request_id", requestID.(string))...
```
**Recommendation:** Use the two-value form with a fallback:
```go
requestID, _ := c.Get(RequestIDKey)
rid, _ := requestID.(string)   // zero value "" on failure — safe
if rid == "" {
    rid = "unknown"
}
log.Error().Str("request_id", rid)...
```

---

### WR-001 — Unsafe type assertion in Logger middleware

**File:** `internal/interfaces/http/middleware/logger.go` (line ~20)  
**Severity:** WARNING  
**Line:** ~20  
**Issue:** Same unchecked type assertion pattern as CR-002 — `requestID.(string)` will panic if the key is not set or holds a non-string value. While the middleware ordering in `router.go` currently guarantees `RequestID()` runs first, any future refactor that registers `Logger()` on a sub-group without `RequestID()` will produce an unrecovered panic.  
**Recommendation:**
```go
requestID, _ := c.Get(RequestIDKey)
rid, _ := requestID.(string)
if rid == "" {
    rid = "unknown"
}
log.Info().Str("request_id", rid)...
```

---

### WR-002 — Client-supplied X-Request-ID accepted without validation

**File:** `internal/interfaces/http/middleware/request_id.go` (line ~12–15)  
**Severity:** WARNING  
**Line:** ~12  
**Issue:** Any non-empty `X-Request-ID` header value from the client is accepted verbatim and stored as the request ID. While zerolog JSON-encodes values (neutralising newline injection into structured logs), the raw string is also echoed back in the response header and could carry arbitrary content. If another log consumer reads the raw string (e.g., a log aggregator that ingests the header value), log-forging or SIEM alert suppression is possible. Best practice for distributed tracing is to validate that client-supplied request IDs conform to UUID format and generate a new one if they do not.  
**Recommendation:**
```go
requestID := c.GetHeader(RequestIDHeader)
if requestID == "" {
    requestID = uuid.New().String()
} else {
    // Validate it is a well-formed UUID to prevent injection
    if _, err := uuid.Parse(requestID); err != nil {
        requestID = uuid.New().String()
    }
}
```

---

### WR-003 — Uploads directory path mismatch between Dockerfile and router

**File:** `Dockerfile` (line ~40), `internal/interfaces/http/router/router.go` (line ~33)  
**Severity:** WARNING  
**Line:** Dockerfile ~40 / router.go ~33  
**Issue:** The Dockerfile creates `/uploads` at the filesystem root (`RUN mkdir -p /uploads`) with `WORKDIR /app`. The router serves static files from the **relative** path `"./uploads"`, which resolves to `/app/uploads` at runtime (WORKDIR is `/app` for the final stage). Files written to `/app/uploads` by the upload service are not accessible from `/uploads`, and vice versa. In the current container, the runtime upload directory (`/app/uploads`) does not exist and is not owned by `appuser`, so file writes will fail with a permission error the first time an upload is attempted.  
**Recommendation:** Create the directory at the correct absolute path and ensure the router path is consistent:
```dockerfile
# Dockerfile — fix path to match WORKDIR
RUN mkdir -p /app/uploads && chown appuser:appgroup /app/uploads
```
Or switch the router to an absolute path:
```go
r.Static("/uploads", "/uploads") // matches Dockerfile mkdir
```
Pick one convention and align both sides.

---

### WR-004 — Default SSL mode is "disable" for database connections

**File:** `configs/config.go` (line ~73)  
**Severity:** WARNING  
**Line:** ~73  
**Issue:** When `DB_SSLMODE` is not set, the default falls back to `"disable"`, meaning all database traffic is unencrypted. This is safe only when the app and database are co-located on a private network with no external exposure. If a production deployment connects to a managed database (e.g., RDS, Cloud SQL), unencrypted connections are a security risk and typically violate compliance requirements. The default should be at minimum `"require"` so that deployments fail loudly rather than running silently insecure.  
**Recommendation:**
```go
if cfg.Database.SSLMode == "" {
    cfg.Database.SSLMode = "require"   // fail loudly if cert is missing
}
```
For deployments that genuinely need `disable` (local Docker Compose), set `DB_SSLMODE=disable` explicitly in the environment.

---

### WR-005 — No startup timeout for database and Redis pings

**File:** `bootstrap/database.go` (line ~12–14), `bootstrap/redis.go` (line ~14)  
**Severity:** WARNING  
**Line:** database.go ~12 / redis.go ~14  
**Issue:** Both `pool.Ping(context.Background())` and `client.Ping(context.Background())` use `context.Background()`, which has no deadline. If the database or Redis is unreachable at startup, the application will block indefinitely — no timeout, no error, no pod restart signal. In Kubernetes or similar orchestration, this prevents the container from transitioning to a failed state and triggering a restart.  
**Recommendation:**
```go
// bootstrap/database.go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := pool.Ping(ctx); err != nil {
    return nil, fmt.Errorf("failed to ping postgres: %w", err)
}
```
Apply the same pattern to the Redis ping.

---

### IN-001 — panic() in package-level init() in timeutil.go

**File:** `internal/domain/shared/timeutil.go` (line ~8–10)  
**Severity:** INFO  
**Line:** ~8  
**Issue:** `panic(...)` inside an `init()` function cannot be caught by the application's recovery middleware and will crash the entire process before `main()` runs. The comment acknowledges this should never happen because `time/tzdata` is embedded in `main.go`, but the dependency between two separate packages (shared domain package relying on main's import side-effect) is fragile. Any binary that imports `shared` without also importing `time/tzdata` will crash on startup.  
**Recommendation:** Consider a lazy-init pattern using `sync.Once` that returns an error rather than panicking, or document the mandatory `import _ "time/tzdata"` requirement in the package-level doc comment. At minimum, change the panic message to explicitly name the missing import:
```go
panic(`failed to load Asia/Tehran: add import _ "time/tzdata" to your main package`)
```

---

### IN-002 — Hardcoded TimeZone=Asia/Tehran in database DSN

**File:** `configs/config.go` (line ~32)  
**Severity:** INFO  
**Line:** ~32  
**Issue:** `TimeZone=Asia/Tehran` is hardcoded in `DatabaseConfig.DSN()` rather than derived from the configurable `App.TimeZone` field. If the app is ever adapted for a different region or the timezone is changed in config, the DSN will silently remain on Tehran time.  
**Recommendation:**
```go
func (d DatabaseConfig) DSN() string {
    tz := "Asia/Tehran" // keep as default but allow override
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
        d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, tz)
}
```
Or thread `AppConfig.TimeZone` into `DSN()` so the two stay in sync.

---

### IN-003 — Hardcoded version string in health endpoint

**File:** `internal/interfaces/http/router/router.go` (line ~28)  
**Severity:** INFO  
**Line:** ~28  
**Issue:** `"version": "1.0.0"` is hardcoded in the health check handler. The version will drift from the actual deployed artifact as the project evolves, making the health endpoint unreliable for deployment verification.  
**Recommendation:** Inject the version at build time via `-ldflags`:
```go
// package main
var Version = "dev"
```
```makefile
-ldflags="-X main.Version=$(VERSION)"
```
Then pass it into the container via `router.New(container, version)` or a global config field.

---

### IN-004 — Raw query string logged verbatim — potential sensitive data exposure

**File:** `internal/interfaces/http/middleware/logger.go` (line ~17)  
**Severity:** INFO  
**Line:** ~17  
**Issue:** `c.Request.URL.RawQuery` is logged as-is. If any endpoint accepts sensitive parameters via query string (e.g., `?token=...`, `?api_key=...`), these values will appear in plain text in log output, violating least-privilege logging principles.  
**Recommendation:** Either scrub known sensitive query param names before logging, or adopt a policy of never putting sensitive data in query strings (enforce via code review). For the logger itself, consider logging only the path without the query, or stripping known sensitive keys:
```go
// Minimal change: don't log query at all, or redact known keys
query := redactSensitiveQuery(c.Request.URL.RawQuery)
```

---

### IN-005 — Global time.Local mutation in main()

**File:** `cmd/server/main.go` (line ~30)  
**Severity:** INFO  
**Line:** ~30  
**Issue:** `time.Local = loc` mutates the global `time.Local` variable. While done once at startup before any goroutines are spawned, `time.Local` is a package-level global shared across the entire process. Any library or package that uses `time.Now()` or `time.Local` directly will silently be affected. The explicit helper `shared.NowTehran()` exists precisely to avoid relying on `time.Local`, making this mutation redundant while adding global state risk.  
**Recommendation:** Remove `time.Local = loc` from `main.go` and rely exclusively on `shared.NowTehran()` and `shared.ToTehran()` throughout the codebase. This makes timezone intent explicit at the call site.

---

### IN-006 — Logger only distinguishes "development" — staging/test environments get production log format

**File:** `bootstrap/logger.go` (line ~14)  
**Severity:** INFO  
**Line:** ~14  
**Issue:** `if env == "development"` is the only branch that enables pretty-printed console output and debug level. Any other env value (e.g., `"staging"`, `"test"`, `"local"`) silently gets JSON production-level logging. This makes it harder to debug staging issues, and test environments with DEBUG expectations will receive INFO-level logs.  
**Recommendation:** Expand the condition or make log level independently configurable:
```go
func InitLogger(env string, level ...string) {
    // Allow LOG_LEVEL override from config
    switch env {
    case "development", "local", "test":
        log.Logger = log.Output(zerolog.ConsoleWriter{...})
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    default:
        log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    }
}
```

---

## Files Reviewed

| # | File |
|---|------|
| 1 | `go.mod` |
| 2 | `configs/config.go` |
| 3 | `configs/migrate.go` |
| 4 | `bootstrap/logger.go` |
| 5 | `bootstrap/database.go` |
| 6 | `bootstrap/redis.go` |
| 7 | `bootstrap/migrator.go` |
| 8 | `cmd/server/main.go` |
| 9 | `internal/domain/shared/apperror.go` |
| 10 | `internal/domain/shared/timeutil.go` |
| 11 | `internal/domain/shared/normalize.go` |
| 12 | `internal/interfaces/http/router/router.go` |
| 13 | `internal/interfaces/http/middleware/request_id.go` |
| 14 | `internal/interfaces/http/middleware/logger.go` |
| 15 | `internal/interfaces/http/middleware/recovery.go` |
| 16 | `internal/interfaces/http/middleware/not_found.go` |
| 17 | `Dockerfile` |

---

_Reviewed: 2026-04-22_  
_Reviewer: gsd-code-reviewer_  
_Depth: standard_
