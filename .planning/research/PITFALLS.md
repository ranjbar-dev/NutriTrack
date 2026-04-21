# Domain Pitfalls — NutriTrack Go Backend

**Domain:** Go DDD Backend · Gin · sqlc · Redis · PostgreSQL · Persian locale  
**Researched:** 2026-04-21  
**Stack:** Go + Gin + sqlc + Redis + PostgreSQL + Docker  
**Confidence:** HIGH (stack-specific, grounded in known Go/sqlc/DDD issues)

---

## Table of Contents

1. [Go DDD Architecture](#1-go-ddd-architecture)
2. [sqlc Usage](#2-sqlc-usage)
3. [Gin Middleware Patterns](#3-gin-middleware-patterns)
4. [Redis OTP + Rate Limiting](#4-redis-otp--rate-limiting)
5. [Persian Text in PostgreSQL](#5-persian-text-in-postgresql)
6. [JWT Security](#6-jwt-security)
7. [Timezone Handling (Asia/Tehran)](#7-timezone-handling-asiatehran)
8. [Offline Sync Idempotency](#8-offline-sync-idempotency)
9. [File Upload Security](#9-file-upload-security)
10. [Docker + Docker Compose](#10-docker--docker-compose)
11. [Phase-Specific Warning Matrix](#11-phase-specific-warning-matrix)

---

## 1. Go DDD Architecture

### CRITICAL — Anemic Domain Models

**What goes wrong:**  
Entities become bare structs with only public fields and zero behavior. All business logic ends up scattered across application services, making the domain layer a pure data layer (the anti-pattern DDD exists to prevent). This is especially tempting with sqlc because generated query result structs map 1:1 to DB columns.

**Why it happens:**  
sqlc generates flat Go structs for query results. Developers copy these directly into the domain layer, ending up with a `User` entity that is literally the same struct as `db.User`. No methods, no invariants, no behavior.

**Consequences:**  
- Business rules duplicated across handlers and services
- Impossible to enforce invariants at compile time
- Domain layer adds no value; might as well use raw SQL

**Prevention:**  
Define domain entities as separate structs with behavior methods. Map from sqlc result to domain entity at the repository boundary:

```go
// WRONG — anemic, mirrors the DB struct
type Client struct {
    ID        uuid.UUID
    Name      string
    IsActive  bool
    CreatedAt time.Time
}

// RIGHT — domain entity with behavior
type Client struct {
    id        uuid.UUID  // unexported fields enforce encapsulation
    name      string
    isActive  bool
    createdAt time.Time
}

func (c *Client) Deactivate() error {
    if !c.isActive {
        return ErrClientAlreadyInactive
    }
    c.isActive = false
    return nil
}

func (c *Client) FullName() string { return c.name }
func (c *Client) IsActive() bool   { return c.isActive }

// Repository maps: db.Client → domain.Client (never expose db.Client to domain)
func toDomainClient(row db.Client) *domain.Client {
    return domain.NewClient(row.ID, row.Name, row.IsActive, row.CreatedAt)
}
```

**Warning signs:** Entity fields are all public. Service methods start with `entity.Field = value`. No methods on entity types.

**Applies to:** Phase 1 (foundation), every subsequent phase

---

### CRITICAL — Infrastructure Leaking into Domain Layer

**What goes wrong:**  
Domain entities import `database/sql`, `jackc/pgx`, Redis clients, or `gin` context. The domain knows about HTTP status codes, SQL null types (`pgtype.Text`), or Redis keys. This makes the domain untestable without a running database.

**Why it happens:**  
Shortcuts during early development. sqlc's `pgtype.NullString` used directly in domain structs. `*gin.Context` passed into application services for "convenience."

**Consequences:**  
- Cannot unit test domain logic without mocking infrastructure
- Database schema changes break domain entity definitions
- Circular import cycles when domain imports db package

**Prevention:**  
Enforce strict layer boundaries via Go package structure:

```
internal/
  domain/          ← NO external imports (stdlib only)
    client/
      entity.go    ← domain entities
      repository.go ← repository INTERFACES (not implementations)
      service.go   ← domain services
  application/     ← imports domain, NOT infrastructure
    client/
      service.go   ← orchestrates domain + calls repo interfaces
  infrastructure/
    postgres/      ← implements domain repository interfaces
    redis/
  handler/         ← imports application services only
    client.go
```

Go's package visibility enforces this if you structure correctly. Run `go vet` + `golang.org/x/tools/cmd/deadcode` in CI to catch violations.

**Warning signs:** `import "github.com/jackc/pgx"` appears in `internal/domain/`. Domain entities have `json:""` tags (they shouldn't — that's a DTO concern).

**Applies to:** Phase 1 (foundation) — fix this early or it permeates everything

---

### MODERATE — Over-Engineering Aggregates for the Diet Plan

**What goes wrong:**  
The Diet Plan hierarchy (Plan → Days → Meals → Options → Items) gets modeled as a single large aggregate with 5 levels of nesting, loaded entirely on every operation. Loading a full diet plan to add one food item becomes a 6-table JOIN.

**Why it happens:**  
Purist DDD thinking: "the plan is one aggregate." But DDD also says aggregates should be as small as possible.

**Prevention:**  
Split into two aggregates with eventual consistency:

```
Aggregate 1: DietPlan (root: DietPlan, includes: PlanDays, Meals, MealOptions)
Aggregate 2: MealOptionItems (root: MealOption, includes: MealOptionItems)
```

Operations that only touch items (add/remove food from option) load only Aggregate 2. Cross-aggregate references use IDs only — `MealOption.ID` not `MealOption` embedded in DietPlan.

**Warning signs:** A handler that "just adds a food item" runs 6 queries to load the whole plan tree before saving.

**Applies to:** Phase 2 (diet plan management)

---

### MODERATE — Missing Repository Transaction Support

**What goes wrong:**  
Repository interfaces have no transaction support. Creating a diet plan requires inserting into 5 tables atomically — without transaction support in repositories, you either leak transactions into application services or roll your own with raw `*sql.Tx`.

**Prevention:**  
Use the Unit of Work pattern with a transaction factory:

```go
// domain/ports.go
type TxManager interface {
    WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// infrastructure/postgres/tx.go
type pgTxManager struct { pool *pgxpool.Pool }

func (m *pgTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
    tx, err := m.pool.Begin(ctx)
    if err != nil { return err }
    if err := fn(pgx.WithTx(ctx, tx)); err != nil {
        tx.Rollback(ctx)
        return err
    }
    return tx.Commit(ctx)
}

// Repository reads tx from context (pgx supports this natively)
func (r *clientRepo) Create(ctx context.Context, c *domain.Client) error {
    q := r.queries.WithTx(pgx.TxFromContext(ctx)) // if no tx, uses pool
    // ...
}
```

**Applies to:** Phase 1 (design), Phase 2 (diet plan creation)

---

## 2. sqlc Usage

### CRITICAL — Nullable Column Handling with pgtype vs. Go Pointers

**What goes wrong:**  
sqlc generates `pgtype.Text` for nullable `VARCHAR` columns and `pgtype.Int4` for nullable `INTEGER`. Code that treats these as Go `*string` or `*int32` panics at runtime or silently loses data.

**Why it happens:**  
sqlc's default behavior with pgx/v5 driver uses `pgtype.*` for all nullable columns. Developers assume `*string` because that's the Go idiom.

**Consequences:**  
- `pgtype.Text{String: "", Valid: false}` serialized as `{}` in JSON, not `null`
- Accessing `.String` on an invalid pgtype panics if not checked
- API responses contain unexpected `{"String":"","Valid":false}` objects

**Prevention — configure sqlc.yaml correctly:**

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./queries"
    schema: "./migrations"
    gen:
      go:
        package: "db"
        out: "internal/infrastructure/postgres/db"
        emit_json_tags: false        # don't expose db tags in JSON
        emit_pointers_for_null_types: true   # use *string instead of pgtype.Text
        null_style: "omit_zero"
```

Or map manually at the repository boundary — never let pgtype leak past the repository.

**Warning signs:** `pgtype.Text` appears in HTTP handler code. JSON responses contain `"Valid":false` fields.

**Applies to:** Phase 1 (sqlc setup), all subsequent phases

---

### CRITICAL — Migration + sqlc Drift (The Out-of-Sync Problem)

**What goes wrong:**  
A migration adds/renames a column. Developer forgets to run `sqlc generate`. Code compiles (old generated files still present) but queries fail at runtime with "column does not exist."

**Prevention:**  
Enforce in CI — never allow `sqlc generate` output to be stale:

```yaml
# .github/workflows/ci.yml
- name: Check sqlc is up to date
  run: |
    sqlc generate
    git diff --exit-code internal/infrastructure/postgres/db/
```

Also add a `make generate` target that runs both `sqlc generate` and `go generate ./...` together, so they're always run as one step.

**Warning signs:** CI does not check sqlc output. Team members commit `*.go` and `*.sql` changes in separate PRs.

**Applies to:** All phases whenever schema changes

---

### MODERATE — Complex Join Queries Returning Partial Results

**What goes wrong:**  
sqlc `GetDietPlanWithDetails` query joins 5 tables. When there are no meal options yet, the LEFT JOIN returns NULLs and sqlc silently omits the row — or generates a struct where every nested field is nullable `pgtype.*`. The handler incorrectly concludes "plan has no meals."

**Prevention:**  
For complex hierarchical queries, prefer multiple simple queries over one massive JOIN:

```go
// application/diet_plan_service.go
// Load plan header
plan, err := r.planRepo.GetByID(ctx, planID)

// Load days separately
days, err := r.dayRepo.ListByPlanID(ctx, planID)

// Load meals for all days in one query (WHERE day_id = ANY($1))
meals, err := r.mealRepo.ListByDayIDs(ctx, dayIDs)

// Assemble in Go — predictable, no NULL surprises
```

Use sqlc for simple, focused queries. Assemble nested structures in Go code. This also makes caching individual levels easier.

**Applies to:** Phase 2 (diet plan), Phase 3 (tracking)

---

### MODERATE — PostgreSQL Enum Types and sqlc Code Generation

**What goes wrong:**  
Adding a new value to a PostgreSQL ENUM type (e.g., adding `powder` to `medication_form`) requires both a migration AND `sqlc generate`. But more critically, the Go enum constants become stale — existing code comparing `db.MedicationFormTablet` compiles but doesn't cover the new value.

**Prevention:**  
For ENUMs with frequent additions (food categories, measurement units), use `TEXT` with a `CHECK` constraint instead:

```sql
-- AVOID for extensible enums:
CREATE TYPE food_category AS ENUM ('breakfast', 'lunch', 'dinner');

-- PREFER for extensible enums:
ALTER TABLE foods ADD COLUMN category TEXT NOT NULL
  CHECK (category IN ('breakfast', 'lunch', 'dinner', 'snack', 'fruit', 'beverage', 'supplement', 'other'));
```

Reserve PostgreSQL ENUMs only for truly fixed-value types (e.g., `user_role`: `super_admin`, `nutritionist`, `client`).

In Go, define typed constants for TEXT-based enums and validate at the application layer:

```go
type FoodCategory string
const (
    FoodCategoryBreakfast FoodCategory = "breakfast"
    FoodCategoryLunch     FoodCategory = "lunch"
    // ...
)

func (c FoodCategory) IsValid() bool {
    switch c {
    case FoodCategoryBreakfast, FoodCategoryLunch /* ... */ :
        return true
    }
    return false
}
```

**Applies to:** Phase 1 (schema design)

---

### MINOR — sqlc `RETURNING` Clauses Skipped

**What goes wrong:**  
INSERT queries don't use `RETURNING id, created_at`. Developer does `INSERT ... ; SELECT ...` as two queries, creating a race condition window and doubling DB round-trips.

**Prevention:**  
Always use `RETURNING` for INSERT/UPDATE:

```sql
-- queries/clients.sql
-- name: CreateClient :one
INSERT INTO clients (id, name, mobile, is_active, created_at)
VALUES ($1, $2, $3, true, NOW())
RETURNING *;
```

**Applies to:** Phase 1 onward

---

## 3. Gin Middleware Patterns

### CRITICAL — Authentication Middleware Registered on Wrong Router Group

**What goes wrong:**  
Auth middleware applied globally catches requests to `/health`, `/metrics`, and `/auth/login` — causing the health check to return 401 and breaking deployment automation. Or conversely, auth middleware added only to some routes and skipped accidentally on sensitive endpoints.

**Prevention:**  
Use explicit router groups with clear auth requirement:

```go
func SetupRoutes(r *gin.Engine, authMiddleware gin.HandlerFunc) {
    // Public — NO auth
    public := r.Group("/api/v1")
    {
        public.POST("/auth/otp/send", handlers.SendOTP)
        public.POST("/auth/otp/verify", handlers.VerifyOTP)
        public.POST("/auth/login", handlers.Login)
        public.GET("/health", handlers.Health)
    }

    // Protected — requires auth
    protected := r.Group("/api/v1")
    protected.Use(authMiddleware)
    {
        // Role-specific sub-groups
        client := protected.Group("/client")
        client.Use(RequireRole(domain.RoleClient))
        client.GET("/plan", handlers.GetActivePlan)

        nutritionist := protected.Group("/nutritionist")
        nutritionist.Use(RequireRole(domain.RoleNutritionist))
        nutritionist.GET("/clients", handlers.ListClients)
    }
}
```

**Warning signs:** Health check endpoint returns 401. Auth middleware registered with `r.Use()` at engine level.

**Applies to:** Phase 1 (auth foundation)

---

### CRITICAL — Error Handling Without Centralized Middleware

**What goes wrong:**  
Each handler calls `c.JSON(http.StatusBadRequest, ...)` with its own error format. Persian error messages hardcoded in 50 different handlers. When the error format changes, 50 files need updates.

**Prevention:**  
Use Gin's error system + a centralized error-handling middleware:

```go
// domain/errors.go
type AppError struct {
    Code       string // machine-readable: "INVALID_OTP"
    MessageFA  string // Persian: "کد تأیید نامعتبر است"
    HTTPStatus int
}
func (e *AppError) Error() string { return e.Code }

var (
    ErrInvalidOTP     = &AppError{"INVALID_OTP", "کد تأیید نامعتبر است", 400}
    ErrOTPExpired     = &AppError{"OTP_EXPIRED", "کد تأیید منقضی شده است", 400}
    ErrUnauthorized   = &AppError{"UNAUTHORIZED", "دسترسی غیرمجاز", 401}
    ErrNotFound       = &AppError{"NOT_FOUND", "مورد مورد نظر یافت نشد", 404}
    ErrRateLimited    = &AppError{"RATE_LIMITED", "تعداد درخواست‌ها بیش از حد مجاز است", 429}
)

// handler/middleware/error.go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if len(c.Errors) == 0 { return }

        err := c.Errors.Last().Err
        var appErr *domain.AppError
        if errors.As(err, &appErr) {
            c.JSON(appErr.HTTPStatus, gin.H{
                "success": false,
                "code":    appErr.Code,
                "message": appErr.MessageFA,
            })
            return
        }
        // Unexpected error — log and return generic Persian message
        log.Error().Err(err).Msg("unexpected error")
        c.JSON(500, gin.H{"success": false, "message": "خطای داخلی سرور"})
    }
}

// In handlers — use c.Error(), NOT c.JSON for errors
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
    if err := h.service.VerifyOTP(c.Request.Context(), req); err != nil {
        c.Error(err) // let middleware handle it
        return
    }
    c.JSON(200, response)
}
```

**Applies to:** Phase 1 — establish this pattern before writing any handlers

---

### MODERATE — Panic Recovery Hiding Real Errors

**What goes wrong:**  
Gin's default `gin.Recovery()` middleware catches panics and returns 500, but swallows the stack trace. In production, nil pointer dereferences in complex diet plan logic are silently recovered — users get 500, developers have no stack trace to debug.

**Prevention:**  
Replace default recovery with a structured logging recovery:

```go
func RecoveryWithLogger(logger zerolog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                logger.Error().
                    Interface("panic", r).
                    Bytes("stack", debug.Stack()).
                    Str("path", c.Request.URL.Path).
                    Msg("panic recovered")
                c.AbortWithStatusJSON(500, gin.H{
                    "success": false,
                    "message": "خطای داخلی سرور",
                })
            }
        }()
        c.Next()
    }
}
```

**Applies to:** Phase 1

---

### MODERATE — Middleware Execution Order with Role Checks

**What goes wrong:**  
Role-checking middleware runs before JWT validation middleware. JWT middleware sets user claims in `c.Set("user", claims)`. If auth middleware panics or short-circuits before setting claims, role middleware reads nil and panics.

**Prevention:**  
Always order middleware explicitly: `Recovery → RequestID → Logging → RateLimit → Auth → RoleCheck`

```go
protected.Use(
    middleware.RequestID(),   // 1. assign trace ID
    middleware.Logger(),      // 2. log request
    middleware.RateLimit(),   // 3. rate limit before auth (save resources)
    middleware.Authenticate(),// 4. validate JWT, set user claims
    // Role checks go on sub-groups, never on the same group as Authenticate
)
```

**Applies to:** Phase 1

---

## 4. Redis OTP + Rate Limiting

### CRITICAL — Race Condition in OTP Attempt Counting

**What goes wrong:**  
OTP attempt counter implemented as GET → increment → SET. Under concurrent requests (user double-taps "send OTP"), two goroutines read count=0, both increment to 1, both SET 1. Rate limit never triggers. An attacker can brute-force OTP with unlimited attempts.

**Prevention:**  
Use Redis atomic `INCR` with `EXPIRE` — never GET+SET:

```go
// infrastructure/redis/otp.go

const (
    otpKeyPrefix      = "otp:code:"      // otp:code:{phone}
    otpAttemptsPrefix = "otp:attempts:"  // otp:attempts:{phone}
    otpRatePrefix     = "otp:rate:"      // otp:rate:{phone}    (send rate limit)
)

// IncrementAttempts returns current attempt count atomically
func (r *RedisOTPRepo) IncrementAttempts(ctx context.Context, phone string) (int64, error) {
    key := otpAttemptsPrefix + phone
    pipe := r.client.TxPipeline()
    incr := pipe.Incr(ctx, key)
    pipe.Expire(ctx, key, 2*time.Minute) // match OTP TTL
    _, err := pipe.Exec(ctx)
    if err != nil { return 0, err }
    return incr.Val(), nil
}

// CheckSendRateLimit uses INCR + EXPIRE for send attempts too
func (r *RedisOTPRepo) CheckSendRateLimit(ctx context.Context, phone string) (bool, error) {
    key := otpRatePrefix + phone
    count, err := r.client.Incr(ctx, key).Result()
    if err != nil { return false, err }
    if count == 1 {
        r.client.Expire(ctx, key, 10*time.Minute)
    }
    return count <= 3, nil // max 3 OTP sends per 10 min
}
```

**Warning signs:** OTP code uses `GET` followed by `SET` for counters. No `TxPipeline` or `INCR` in OTP code.

**Applies to:** Phase 1 (auth)

---

### CRITICAL — Redis Key Namespace Collisions

**What goes wrong:**  
OTP key: `phone:09123456789`. Rate limit key: `09123456789`. Refresh token key: `token:abc123`. No consistent namespace. When the team grows or features are added, keys collide silently — an OTP delete accidentally deletes a rate limit counter.

**Prevention:**  
Enforce a key naming convention from day one:

```
otp:code:{phone}              TTL: 2min
otp:attempts:{phone}          TTL: 2min  
otp:rate:{phone}              TTL: 10min
session:refresh:{tokenID}     TTL: 30d
session:blacklist:{tokenID}   TTL: 15min (access token expiry)
ratelimit:api:{userID}:{endpoint}  TTL: 1min
cache:food:{foodID}           TTL: 1h
push:subscription:{userID}    TTL: permanent (no TTL)
```

Define these as typed constants in Go, never as inline strings:

```go
// infrastructure/redis/keys.go
func OTPCodeKey(phone string) string     { return "otp:code:" + phone }
func OTPAttemptsKey(phone string) string { return "otp:attempts:" + phone }
func RefreshTokenKey(id string) string   { return "session:refresh:" + id }
func BlacklistKey(jti string) string     { return "session:blacklist:" + jti }
```

**Warning signs:** Redis key strings inline in handler or service code. No dedicated keys package.

**Applies to:** Phase 1 onward

---

### MODERATE — TTL Not Set on Rate Limit Keys

**What goes wrong:**  
Rate limit key set with INCR but `EXPIRE` only called conditionally (e.g., only when count == 1). If the first INCR+EXPIRE call succeeds but subsequent calls within the window happen and `EXPIRE` isn't refreshed, it's fine. BUT if there's a Redis restart before `EXPIRE` is set, the key persists forever — permanently banning a phone number.

**Prevention:**  
Always use pipeline to set key + expiry atomically. For rate limiters, use `SET key value EX seconds NX` pattern or Lua scripts for atomic check-and-set:

```go
// Lua script for atomic OTP rate limit check
const otpRateLimitScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local ttl = tonumber(ARGV[2])
local current = redis.call('INCR', key)
if current == 1 then
    redis.call('EXPIRE', key, ttl)
end
if current > limit then
    return 0
end
return 1
`

var otpRateLimitScriptSHA = "" // cache after SCRIPT LOAD
```

**Applies to:** Phase 1 (auth)

---

### MINOR — Storing OTP as Plaintext in Redis

**What goes wrong:**  
OTP `"123456"` stored as plaintext string in Redis. If Redis is compromised (no requirepass, exposed port), all valid OTPs are readable.

**Prevention:**  
Store bcrypt hash of OTP in Redis. Cost 4-6 (lower than password bcrypt, OTPs are short-lived):

```go
func (r *RedisOTPRepo) StoreOTP(ctx context.Context, phone, otp string, ttl time.Duration) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(otp), 4)
    if err != nil { return err }
    return r.client.Set(ctx, OTPCodeKey(phone), hash, ttl).Err()
}

func (r *RedisOTPRepo) VerifyOTP(ctx context.Context, phone, otp string) (bool, error) {
    hash, err := r.client.Get(ctx, OTPCodeKey(phone)).Bytes()
    if err == redis.Nil { return false, domain.ErrOTPExpired }
    if err != nil { return false, err }
    err = bcrypt.CompareHashAndPassword(hash, []byte(otp))
    return err == nil, nil
}
```

**Applies to:** Phase 1

---

## 5. Persian Text in PostgreSQL

### CRITICAL — pg_trgm Extension Not Enabled in Migration

**What goes wrong:**  
The migration creates the `GIN` index on the `name` column using `pg_trgm` ops, but `CREATE EXTENSION IF NOT EXISTS pg_trgm` is missing from the migration. The migration fails in a fresh Docker container where the extension isn't pre-installed. Works on developer machine (which happens to have it), fails in CI and production.

**Prevention:**  
Always include extension creation as the first migration:

```sql
-- migrations/000001_init_extensions.up.sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;  -- for diacritic-insensitive search
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- migrations/000002_create_foods.up.sql
CREATE TABLE foods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    -- ...
);

-- Index MUST come after extension
CREATE INDEX idx_foods_name_trgm ON foods USING GIN (name gin_trgm_ops);
```

**Warning signs:** Extension creation in a comment or README instead of migration. Index creation in same migration as table creation but before extension migration.

**Applies to:** Phase 1 (schema setup)

---

### CRITICAL — Full-Text Search Failing on Persian (Arabic-Script) Text

**What goes wrong:**  
PostgreSQL's built-in full-text search (`tsvector`/`tsquery`) has NO Persian language dictionary. Using `to_tsvector('persian', name)` fails with "text search configuration 'persian' does not exist." Using `to_tsvector('english', name)` on Persian text produces nothing useful — Persian words aren't stemmed by English rules.

**Prevention:**  
For Persian text, use `pg_trgm` similarity search, NOT `tsvector`. This is the correct approach for NutriTrack:

```sql
-- WRONG — no Persian tsvector config exists
SELECT * FROM foods WHERE to_tsvector('persian', name) @@ to_tsquery('persian', 'مرغ');

-- RIGHT — trigram similarity search works with any script
SELECT * FROM foods
WHERE name % $1          -- % operator: similarity > threshold (default 0.3)
   OR name ILIKE '%' || $1 || '%'  -- ILIKE fallback for exact substring
ORDER BY similarity(name, $1) DESC, name
LIMIT 20 OFFSET $2;

-- Set similarity threshold (0.1 = more permissive, 0.6 = more strict)
-- For Persian with short food names, 0.2 is a good starting point
SET pg_trgm.similarity_threshold = 0.2;
```

```go
// In sqlc query:
-- name: SearchFoods :many
SELECT * FROM foods
WHERE is_active = true
  AND (
    $1::text = '' 
    OR name ILIKE '%' || $1 || '%'
    OR similarity(name, $1) > 0.2
  )
ORDER BY 
  CASE WHEN name ILIKE '%' || $1 || '%' THEN 0 ELSE 1 END,
  similarity(name, $1) DESC
LIMIT $2 OFFSET $3;
```

**Warning signs:** Any `to_tsvector` or `tsquery` used for food/medication name search. Search returns 0 results for Persian queries that clearly match.

**Applies to:** Phase 1 (food database), Phase 2 (medication database)

---

### MODERATE — Persian Text Normalization Not Applied

**What goes wrong:**  
Arabic Kaf (ك) vs Persian Keh (ک) and Arabic Yeh (ي) vs Persian Yeh (ی) are different Unicode code points but look identical. A food named with Arabic characters won't match a search with Persian characters. Iranian users type on both Persian and Arabic keyboards.

**Prevention:**  
Normalize at INSERT time. Create a PostgreSQL function or do it in Go before storing:

```go
// infrastructure/text/normalize.go
import "strings"

// NormalizePersian converts Arabic-script variants to their Persian equivalents
func NormalizePersian(s string) string {
    replacer := strings.NewReplacer(
        "ك", "ک",  // Arabic Kaf → Persian Keh
        "ي", "ی",  // Arabic Yeh → Persian Yeh  
        "ة", "ه",  // Ta Marbuta → He
        "ؤ", "و",  // Waw with Hamza → Waw
        "إ", "ا",  // Alef with Hamza below → Alef
        "أ", "ا",  // Alef with Hamza above → Alef
        "آ", "ا",  // Alef with Madda → Alef (optional, depends on requirements)
    )
    return replacer.Replace(s)
}
```

Apply normalization to BOTH stored data (on insert) AND search queries (before querying).

**Applies to:** Phase 1 (food database), all text search features

---

### MODERATE — Collation Issues with Persian Sorting

**What goes wrong:**  
`ORDER BY name` on Persian text uses C locale (byte-order). Persian letters sort by their Unicode code points, not Persian alphabetical order. Food list appears in random-looking order to Iranian users.

**Prevention:**  
For display ordering, either sort in Go using `golang.org/x/text/collate` with Persian locale, or accept the limitation for search results (where relevance sorting matters more than alphabetical):

```go
import (
    "golang.org/x/text/collate"
    "golang.org/x/text/language"
)

col := collate.New(language.Persian)
sort.SliceStable(foods, func(i, j int) bool {
    return col.CompareString(foods[i].Name, foods[j].Name) < 0
})
```

**Applies to:** Phase 1 (food database listing)

---

## 6. JWT Security

### CRITICAL — No Token Invalidation on Logout

**What goes wrong:**  
JWT access tokens are stateless — there's no way to invalidate them server-side once issued. A user logs out, but their 15-minute access token remains valid. If a nutritionist's device is stolen, the thief has 15 minutes of access with no way to revoke it.

**Prevention:**  
Maintain a Redis blacklist for invalidated access tokens using the `jti` (JWT ID) claim:

```go
// When issuing tokens, add jti
claims := jwt.MapClaims{
    "sub": userID,
    "jti": uuid.New().String(), // unique token ID
    "exp": time.Now().Add(15 * time.Minute).Unix(),
    "role": role,
}

// On logout — blacklist the jti
func (r *RedisTokenRepo) BlacklistToken(ctx context.Context, jti string, expiry time.Duration) error {
    return r.client.Set(ctx, BlacklistKey(jti), "1", expiry).Err()
}

// In auth middleware — check blacklist
func (m *AuthMiddleware) Authenticate(c *gin.Context) {
    token, claims, err := m.parseToken(c)
    // ...
    jti := claims["jti"].(string)
    blacklisted, err := m.tokenRepo.IsBlacklisted(c.Request.Context(), jti)
    if err != nil || blacklisted {
        c.AbortWithStatusJSON(401, gin.H{"message": "توکن نامعتبر است"})
        return
    }
    c.Set("userID", claims["sub"])
    c.Next()
}
```

Set blacklist TTL = token's remaining expiry time (parse `exp` claim), so Redis doesn't accumulate stale entries.

**Applies to:** Phase 1 (auth)

---

### CRITICAL — Refresh Token Not Rotated on Use

**What goes wrong:**  
Refresh token is stored in Redis, used once to get a new access token, but NOT invalidated. An attacker who steals a refresh token can use it indefinitely — even after the legitimate user has logged in multiple times.

**Prevention:**  
Implement refresh token rotation — each use invalidates the old token and issues a new one:

```go
func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
    // 1. Validate refresh token exists in Redis
    tokenData, err := s.tokenRepo.GetRefreshToken(ctx, refreshToken)
    if err != nil { return nil, domain.ErrInvalidRefreshToken }

    // 2. Delete OLD refresh token BEFORE issuing new one
    if err := s.tokenRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
        return nil, err
    }

    // 3. Issue new token pair
    newPair, err := s.issueTokenPair(ctx, tokenData.UserID, tokenData.Role)
    if err != nil { return nil, err }

    return newPair, nil
}
```

If a stolen refresh token is used after the legitimate user has already rotated it, the old token won't exist in Redis — automatic detection.

**Applies to:** Phase 1 (auth)

---

### MODERATE — JWT Secret in Environment Without Validation

**What goes wrong:**  
App starts with `JWT_SECRET=""` (missing env var). JWT tokens signed with empty string. Every crafted token is valid. This happens in local dev, slips into staging, then production.

**Prevention:**  
Validate all required env vars at startup, fail fast:

```go
// config/config.go
func Load() (*Config, error) {
    cfg := &Config{}
    if cfg.JWTSecret = os.Getenv("JWT_SECRET"); len(cfg.JWTSecret) < 32 {
        return nil, errors.New("JWT_SECRET must be at least 32 characters")
    }
    // ...
    return cfg, nil
}

// main.go
cfg, err := config.Load()
if err != nil {
    log.Fatal().Err(err).Msg("invalid configuration — refusing to start")
}
```

**Applies to:** Phase 1

---

## 7. Timezone Handling (Asia/Tehran)

### CRITICAL — Mixing time.Local, time.UTC, and Asia/Tehran

**What goes wrong:**  
`time.Now()` returns time in the server's local timezone. `time.Now().UTC()` returns UTC. `time.Now().In(tehranLoc)` returns Iran time. When these are mixed — storing `time.Now()` in one place and `time.Now().UTC()` in another — comparisons and queries produce wrong results. Persian date display shows wrong day if UTC is used (Iran is UTC+3:30, crossing midnight at 20:30 UTC in summer).

**Prevention:**  
One rule: **Always store UTC. Always display Tehran.** Load the timezone once at startup:

```go
// config/timezone.go
var TehranLocation *time.Location

func init() {
    var err error
    TehranLocation, err = time.LoadLocation("Asia/Tehran")
    if err != nil {
        // If TZ data not available (Alpine Linux), use fixed offset as fallback
        TehranLocation = time.FixedZone("IRST", 3*60*60+30*60)
        log.Warn().Msg("Asia/Tehran timezone not found in tzdata, using fixed IRST offset — DST will not be applied")
    }
}

// Store — always UTC
func nowUTC() time.Time { return time.Now().UTC() }

// Display / business logic — Tehran
func nowTehran() time.Time { return time.Now().In(TehranLocation) }

// "Today" in Tehran timezone — critical for daily tracking
func todayTehran() time.Time {
    t := time.Now().In(TehranLocation)
    return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, TehranLocation)
}
```

**Warning signs:** `time.Now()` used directly without `.UTC()` in database writes. `time.Local` in any code path. `"2006-01-02"` date formatting without specifying timezone location.

**Applies to:** All phases, especially Phase 3 (daily tracking)

---

### CRITICAL — Alpine Linux Docker Image Missing tzdata

**What goes wrong:**  
Go binary compiled on Ubuntu/Debian, deployed on Alpine Docker image. `time.LoadLocation("Asia/Tehran")` returns `unknown time zone Asia/Tehran` because Alpine minimal image has no timezone database. App starts, OTP 2-minute expiry works, but "today's" food log date is wrong — app thinks it's still yesterday.

**Prevention:**  
Two options — choose one:

**Option A: Add tzdata to Dockerfile**
```dockerfile
# Dockerfile
FROM alpine:3.19
RUN apk add --no-cache tzdata ca-certificates
ENV TZ=Asia/Tehran
COPY --from=builder /app/nutritrack /app/nutritrack
```

**Option B: Embed tzdata in the binary (zero-dependency)**
```go
import _ "time/tzdata" // Go 1.15+ — embeds timezone database in binary
```

Add this import to `main.go`. Binary size increases by ~450KB but works on any OS/Docker image.

**Recommendation:** Use Option B (embedded tzdata) — Docker images may change, binary always works.

**Warning signs:** `time.LoadLocation` error not checked. App deployed on Alpine without tzdata package.

**Applies to:** Phase 1 (infrastructure), Phase 3 (all date-dependent tracking)

---

### MODERATE — DST Handling for Iran (IRST/IRDT)

**What goes wrong:**  
Iran observes Daylight Saving Time (IRDT = UTC+4:30, active roughly March 21 – September 21). Using `time.FixedZone("IRST", 3*60*60+30*60)` misses DST. During summer months, all timestamps are off by 1 hour. Medication reminders scheduled at "08:00 Tehran" fire at 07:00 local time in summer.

**Prevention:**  
Always use `time.LoadLocation("Asia/Tehran")` with the embedded tzdata (as above) — it includes DST transitions. Never use `time.FixedZone` for Tehran as primary implementation (only as fallback).

**Applies to:** Phase 1 (auth/OTP timing), Phase 3 (tracking timestamps), Phase 4 (push notifications)

---

### MODERATE — PostgreSQL TIMESTAMPTZ vs TIMESTAMP

**What goes wrong:**  
Using `TIMESTAMP WITHOUT TIME ZONE` (plain `TIMESTAMP`) in PostgreSQL. Values stored as "2026-03-15 09:00:00" with no timezone info. When Go reads this back via pgx, it's interpreted as UTC. But if the INSERT came from a Tehran-timezone client without UTC conversion, you've silently stored a Tehran timestamp as UTC — now 3.5 hours off.

**Prevention:**  
Always use `TIMESTAMPTZ` (TIMESTAMP WITH TIME ZONE) in every table:

```sql
CREATE TABLE food_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    logged_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),  -- always TIMESTAMPTZ
    date DATE NOT NULL,  -- the Tehran date (stored explicitly, not derived)
    -- ...
);
```

Store `date` (the Tehran calendar date) as an explicit `DATE` column — don't derive it from `logged_at` because `logged_at::date` in PostgreSQL uses the DB server timezone, not Tehran.

**Applies to:** Phase 1 (schema design)

---

## 8. Offline Sync Idempotency

### CRITICAL — local_id Dedup Not Enforced at DB Level

**What goes wrong:**  
`local_id` deduplication implemented only in application code: "check if local_id exists, then insert." Under concurrent sync requests (user loses connection mid-sync, retries), two parallel requests both pass the check and insert duplicates. Food logs, water entries, exercise logs all duplicated.

**Prevention:**  
Enforce at the database level with a unique constraint:

```sql
CREATE TABLE food_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id),
    local_id TEXT,  -- client-generated UUID for offline sync
    date DATE NOT NULL,
    -- ...
    UNIQUE (client_id, local_id)  -- DB-enforced, handles race conditions
);
```

In Go, handle the constraint violation gracefully (return the existing record, not an error):

```go
func (r *FoodLogRepo) CreateWithLocalID(ctx context.Context, log *domain.FoodLog) (*domain.FoodLog, error) {
    existing, err := r.queries.GetFoodLogByLocalID(ctx, db.GetFoodLogByLocalIDParams{
        ClientID: log.ClientID,
        LocalID:  log.LocalID,
    })
    
    result, err := r.queries.CreateFoodLog(ctx, params)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
            // Already exists — return existing record (idempotent)
            return r.queries.GetFoodLogByLocalID(ctx, ...)
        }
        return nil, err
    }
    return toDomain(result), nil
}
```

**Warning signs:** No `UNIQUE` constraint on `(client_id, local_id)`. Sync code does SELECT-then-INSERT without DB constraint backup.

**Applies to:** Phase 3 (all tracking endpoints), Phase 5 (sync strategy)

---

### CRITICAL — Partial Sync Failures Leave Inconsistent State

**What goes wrong:**  
Client syncs 50 offline records in a single batch request. Record 27 fails validation (invalid food_id). Remaining 23 records are not processed. Client doesn't know which records synced — retries the whole batch. Records 1-26 are duplicated (if local_id dedup not working) or create errors.

**Prevention:**  
Design sync endpoint to process each record independently and return per-item results:

```go
// POST /api/v1/sync
type SyncRequest struct {
    FoodLogs    []FoodLogSync    `json:"food_logs"`
    WaterLogs   []WaterLogSync   `json:"water_logs"`
    ExerciseLogs []ExerciseLogSync `json:"exercise_logs"`
}

type SyncResponse struct {
    FoodLogs []SyncItemResult `json:"food_logs"`
    // ...
}

type SyncItemResult struct {
    LocalID string `json:"local_id"`
    Success bool   `json:"success"`
    ServerID string `json:"server_id,omitempty"`
    Error   string `json:"error,omitempty"` // Persian error if failed
}
```

Process each item in its own transaction. Return partial success — client only retries failed items (by local_id). This requires the DB-level unique constraint to make retries safe.

**Applies to:** Phase 5 (offline sync)

---

### MODERATE — Sync Conflict Resolution Not Defined

**What goes wrong:**  
Nutritionist updates a client's body measurement at 10:00. Client's offline data (recorded at 09:00) syncs at 11:00. The offline data is older but was entered by the client — should it overwrite the nutritionist's entry?

**Prevention:**  
Define and document conflict rules before implementation:

| Entity | Conflict Rule |
|--------|--------------|
| Body measurements | Last write wins by `recorded_at` timestamp |
| Food logs | Per-day-per-meal: client's choice wins; nutritionist cannot override |
| Sleep entries | Last write wins; only one entry per day |
| Water entries | Additive — both entries kept (each is a separate intake event) |
| Medication intake | Additive — each `taken_at` timestamp is a separate event |

Implement: sync endpoint receives `recorded_at` from client. Server uses `recorded_at` (not `synced_at`) for conflict resolution.

**Applies to:** Phase 5 (offline sync)

---

## 9. File Upload Security

### CRITICAL — MIME Type Spoofing via Content-Type Header

**What goes wrong:**  
File type validation checks only `Content-Type: application/pdf` header. Attacker renames `malware.exe` to `report.pdf` and sets `Content-Type: application/pdf`. The header passes validation, the .exe is stored on the filesystem.

**Prevention:**  
Read the first 512 bytes of the file and detect MIME type from magic bytes, NOT the header:

```go
// handler/upload.go
func validateFileType(file multipart.File, allowedTypes []string) (string, error) {
    // Read magic bytes
    buf := make([]byte, 512)
    n, err := file.Read(buf)
    if err != nil { return "", err }
    
    // Reset reader for subsequent storage
    file.Seek(0, io.SeekStart)
    
    // Detect real MIME type from content
    detectedType := http.DetectContentType(buf[:n])
    
    for _, allowed := range allowedTypes {
        if detectedType == allowed { return detectedType, nil }
    }
    return "", fmt.Errorf("نوع فایل مجاز نیست: %s", detectedType)
}

// For PDF specifically — http.DetectContentType returns "application/pdf"
// For images — returns "image/jpeg" or "image/png"
// For executables — returns "application/octet-stream" (rejected)

var labResultAllowedTypes = []string{
    "application/pdf",
    "image/jpeg",
    "image/png",
}
```

**Warning signs:** File type check only reads `r.Header.Get("Content-Type")`. No magic byte inspection.

**Applies to:** Phase 4 (lab results), Phase 3 (messaging attachments)

---

### CRITICAL — Path Traversal in File Storage

**What goes wrong:**  
Filename from client used directly to construct storage path: `/data/uploads/` + `filename`. Client sends `../../../../etc/cron.d/backdoor` as filename. File written to `/etc/cron.d/backdoor`.

**Prevention:**  
Never use client-provided filenames. Generate server-side UUID-based paths:

```go
// infrastructure/storage/local.go
func (s *LocalStorage) Store(ctx context.Context, file io.Reader, mimeType string) (string, error) {
    // Generate server-controlled path — never trust client filename
    id := uuid.New().String()
    ext := mimeToExt(mimeType) // ".pdf", ".jpg", ".png"
    
    // Partition by date to avoid huge directories
    now := time.Now().UTC()
    relPath := fmt.Sprintf("%d/%02d/%s%s", now.Year(), now.Month(), id, ext)
    fullPath := filepath.Join(s.basePath, relPath)
    
    // Ensure the path is within basePath (defense in depth)
    if !strings.HasPrefix(fullPath, filepath.Clean(s.basePath)+string(os.PathSeparator)) {
        return "", errors.New("invalid storage path")
    }
    
    if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
        return "", err
    }
    
    f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
    // ...
    return relPath, nil // store relative path in DB
}
```

**Applies to:** Phase 4 (lab results), Phase 3 (messaging)

---

### MODERATE — File Size Limit Not Enforced at Server Level

**What goes wrong:**  
`maxMemory` in `r.ParseMultipartForm(10 << 20)` controls how much is buffered in memory, but doesn't reject oversized files — it just spills to disk. A 2 GB file upload is happily stored before the Go handler has a chance to check `header.Size`.

**Prevention:**  
Use `http.MaxBytesReader` BEFORE parsing the multipart form:

```go
func (h *UploadHandler) UploadLabResult(c *gin.Context) {
    // Limit at the HTTP reader level — rejects large requests immediately
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10MB
    
    if err := c.Request.ParseMultipartForm(2 << 20); err != nil {
        if err.Error() == "http: request body too large" {
            c.Error(domain.ErrFileTooLarge) // "حجم فایل بیش از ۱۰ مگابایت است"
            return
        }
        c.Error(err)
        return
    }
}
```

**Applies to:** Phase 4 (lab results upload)

---

### MINOR — Uploaded Files Served Without Content-Disposition

**What goes wrong:**  
Lab result PDFs served via `/uploads/2026/03/uuid.pdf`. Browser opens the PDF inline instead of downloading. More critically, if an attacker uploads an HTML file (bypassing type check), the browser executes it as JavaScript (XSS).

**Prevention:**  
Serve all uploaded files with `Content-Disposition: attachment` and re-validate MIME type from disk:

```go
func (h *FileHandler) ServeFile(c *gin.Context) {
    // ... auth check, path resolution ...
    
    c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, safeFilename))
    c.Header("X-Content-Type-Options", "nosniff")
    c.Header("Content-Security-Policy", "default-src 'none'")
    c.File(fullPath)
}
```

**Applies to:** Phase 4 (file serving)

---

## 10. Docker + Docker Compose

### CRITICAL — Go App Starts Before PostgreSQL Is Ready

**What goes wrong:**  
`docker compose up` starts all services in parallel. Go app starts, tries to connect to PostgreSQL, gets "connection refused" or "database does not exist," crashes. Docker restarts it. By restart 3, PostgreSQL is ready, app starts fine. BUT: if `restart: unless-stopped` is not set, the app stays crashed.

This is worse with migrations — `golang-migrate` runs before PostgreSQL has initialized the `nutritrack` database.

**Prevention:**  
Use a proper health check + `depends_on` condition:

```yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: nutritrack
      POSTGRES_USER: nutritrack
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U nutritrack -d nutritrack"]
      interval: 5s
      timeout: 5s
      retries: 10
      start_period: 10s
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

  migrate:
    image: migrate/migrate:v4
    depends_on:
      postgres:
        condition: service_healthy
    command: ["-path=/migrations", "-database", "${DATABASE_URL}", "up"]
    volumes:
      - ./migrations:/migrations

  app:
    build: .
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      migrate:
        condition: service_completed_successfully  # wait for migrations
    environment:
      TZ: Asia/Tehran
    restart: unless-stopped
```

**Warning signs:** `depends_on` lists service names without `condition:`. No health checks on postgres/redis. App container crashes in logs on first start.

**Applies to:** Phase 1 (Docker setup)

---

### CRITICAL — Migration Run on Every App Restart

**What goes wrong:**  
Migration logic embedded in the Go app (`migrate.Up()` called in `main.go`). Every `docker compose restart app` re-runs migrations. While `golang-migrate` is idempotent for "already applied" migrations, two app instances starting simultaneously (during rolling update) both try to acquire the migration lock — one times out and fails.

**Prevention:**  
Run migrations as a separate one-shot service (as shown above). The app should NOT run migrations. Separation of concerns:

```go
// main.go — NO migration code here
func main() {
    cfg := config.MustLoad()
    db := database.MustConnect(cfg.DatabaseURL)
    // ... start server
}

// Separate: cmd/migrate/main.go (or use migrate/migrate Docker image)
```

**Applies to:** Phase 1

---

### MODERATE — Volume Permissions for Uploaded Files

**What goes wrong:**  
Docker container runs as root (default). Files uploaded to `/data/uploads/` are owned by root:root with mode 644. When Traefik or a backup script running as a different user tries to read them, permission denied. Worse: if the Go app is later changed to run as a non-root user (security hardening), it can't write to `/data/uploads/` anymore.

**Prevention:**  
Set consistent ownership in Dockerfile and docker-compose:

```dockerfile
# Dockerfile
FROM alpine:3.19
RUN apk add --no-cache tzdata ca-certificates
RUN addgroup -g 1001 nutritrack && adduser -u 1001 -G nutritrack -D nutritrack
RUN mkdir -p /data/uploads && chown -R nutritrack:nutritrack /data/uploads

USER nutritrack
COPY --from=builder --chown=nutritrack:nutritrack /app/nutritrack /app/nutritrack
```

```yaml
# docker-compose.yml
services:
  app:
    user: "1001:1001"
    volumes:
      - uploads_data:/data/uploads
```

**Applies to:** Phase 1 (Docker), Phase 4 (file uploads)

---

### MODERATE — TZ Environment Variable Not Set in All Services

**What goes wrong:**  
`TZ=Asia/Tehran` set in the Go app container but NOT in the PostgreSQL container. `NOW()` in PostgreSQL returns UTC (PostgreSQL default). `DEFAULT NOW()` timestamps stored in UTC — correct for `TIMESTAMPTZ`. BUT `CURRENT_DATE` and `CURRENT_TIME` return UTC date/time, not Tehran. If any query uses `CURRENT_DATE` for "today's tracking," it returns yesterday evening for Iranian users after 20:30 UTC.

**Prevention:**  
Set timezone in PostgreSQL container AND never use `CURRENT_DATE` in SQL queries — always pass the Tehran date from Go:

```yaml
services:
  postgres:
    environment:
      POSTGRES_DB: nutritrack
      TZ: Asia/Tehran
      PGTZ: Asia/Tehran  # also set PGTZ for pg-specific timezone
```

```go
// In Go — pass Tehran date explicitly to SQL
tehranDate := time.Now().In(config.TehranLocation).Format("2006-01-02")
// Pass as parameter to sqlc query, never rely on DB CURRENT_DATE
```

**Applies to:** Phase 1 (Docker), Phase 3 (daily tracking)

---

### MINOR — Secrets in docker-compose.yml

**What goes wrong:**  
`JWT_SECRET: "my-secret-key"` hardcoded in `docker-compose.yml` committed to git. Repository is public or gets leaked — all production tokens are now compromised.

**Prevention:**  
Use `.env` file (gitignored) with docker-compose's `env_file` or environment variable substitution:

```yaml
# docker-compose.yml
services:
  app:
    env_file:
      - .env.production  # NOT committed to git
```

```
# .env.example (committed — shows required vars without values)
JWT_SECRET=
POSTGRES_PASSWORD=
REDIS_PASSWORD=
KAVENEGAR_API_KEY=
VAPID_PRIVATE_KEY=
```

**Applies to:** Phase 1

---

## 11. Phase-Specific Warning Matrix

| Phase | Topic | Critical Pitfall | Mitigation |
|-------|-------|-----------------|------------|
| Phase 1 | Foundation & Auth | DDD layer separation violated by sqlc structs | Strict package structure; map at repo boundary |
| Phase 1 | Foundation & Auth | OTP race condition in attempt counting | Use Redis INCR, never GET+SET |
| Phase 1 | Foundation & Auth | JWT_SECRET validation on startup | Fail fast with config validation |
| Phase 1 | Foundation & Auth | No refresh token rotation | Delete old token before issuing new |
| Phase 1 | Foundation & Auth | Alpine + missing tzdata | `import _ "time/tzdata"` in main.go |
| Phase 1 | Docker setup | App starts before PG ready | `depends_on` with `condition: service_healthy` |
| Phase 1 | Docker setup | Migrations run on app start | Separate migrate service/container |
| Phase 1 | Schema | PostgreSQL ENUM vs TEXT | Use TEXT+CHECK for extensible enums |
| Phase 1 | Schema | `TIMESTAMP` vs `TIMESTAMPTZ` | Always `TIMESTAMPTZ` |
| Phase 1 | Schema | pg_trgm not in migration | First migration: `CREATE EXTENSION IF NOT EXISTS pg_trgm` |
| Phase 2 | Diet Plan | Over-engineering single large aggregate | Two aggregates: DietPlan + MealOptionItems |
| Phase 2 | Diet Plan | No transaction support in repos | TxManager pattern in domain/ports |
| Phase 2 | Food DB | Persian full-text search with tsvector | Use pg_trgm similarity, never tsvector for Persian |
| Phase 2 | Food DB | Arabic/Persian char normalization | Normalize at INSERT and search time |
| Phase 3 | Daily Tracking | "Today" date using UTC | Always compute date in TehranLocation |
| Phase 3 | Daily Tracking | DST transition breaks Tehran offset | Use LoadLocation not FixedZone |
| Phase 4 | File Uploads | MIME type spoofing | Magic byte detection, not Content-Type header |
| Phase 4 | File Uploads | Path traversal | UUID-based server-generated paths only |
| Phase 4 | File Uploads | Size limit not enforced | `http.MaxBytesReader` before ParseMultipartForm |
| Phase 5 | Offline Sync | local_id dup at DB level only | `UNIQUE (client_id, local_id)` constraint |
| Phase 5 | Offline Sync | Batch sync partial failure | Per-item response with per-item error |
| Phase 5 | Offline Sync | Conflict resolution undefined | Document rules before implementation |
| All | Gin middleware | Centralized error handling | AppError type + error middleware; no inline c.JSON for errors |
| All | Gin middleware | Persian errors hardcoded in handlers | Centralized AppError catalog in domain/errors.go |
| All | sqlc | Generated code diverges from schema | CI check: `sqlc generate && git diff --exit-code` |
| All | sqlc | Nullable pgtype leaks past repo | `emit_pointers_for_null_types: true` in sqlc.yaml |

---

## Sources & Confidence

| Area | Confidence | Basis |
|------|------------|-------|
| Go DDD patterns | HIGH | Well-documented anti-patterns in Go community; `go-ddd` repos; Matt Boyle's DDD in Go |
| sqlc nullable handling | HIGH | sqlc GitHub issues; pgx/v5 migration guides |
| Redis OTP race conditions | HIGH | Redis documentation on atomic operations; INCR/EXPIRE pattern |
| Persian pg_trgm | HIGH | PostgreSQL docs confirm no Persian FTS config; pg_trgm is the standard approach |
| tzdata on Alpine | HIGH | Known Go issue; `time/tzdata` package docs; multiple production post-mortems |
| JWT invalidation | HIGH | Standard security practice; well-documented |
| DST for Iran | HIGH | IANA timezone database; Iran Standard Time specification |
| local_id idempotency | HIGH | Standard distributed systems pattern; PostgreSQL constraint behavior |
| File magic bytes | HIGH | Go stdlib `http.DetectContentType` docs; OWASP file upload guidance |
| Docker health checks | HIGH | Docker Compose v2 specification; `condition: service_healthy` |
