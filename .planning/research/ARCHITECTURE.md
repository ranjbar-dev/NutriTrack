# Architecture Patterns

**Domain:** Nutrition Management PWA (Nutritionist–Client Platform)
**Researched:** 2025-07-18
**Confidence:** HIGH — Stack is well-understood (Go/Gin, PostgreSQL, Nuxt 4), PRD provides exhaustive data model and feature specs.

---

## System Overview

```
┌──────────────────────────────────────────────────────────────┐
│                   Client Devices (Mobile)                     │
│  ┌────────────────────────────────────────────────────────┐  │
│  │              Nuxt 4 PWA (app/ directory)               │  │
│  │  ┌──────────┐ ┌──────────┐ ┌────────────────────────┐ │  │
│  │  │  Pinia   │ │Composable│ │  Service Worker        │ │  │
│  │  │  Stores  │ │  Layer   │ │  (@vite-pwa/nuxt)      │ │  │
│  │  └────┬─────┘ └────┬─────┘ └─────────┬──────────────┘ │  │
│  │       │             │                 │                │  │
│  │  ┌────▼─────────────▼─────┐  ┌───────▼──────────────┐ │  │
│  │  │   ofetch / useFetch    │  │    IndexedDB          │ │  │
│  │  │   (API Client Layer)   │  │    (Dexie.js)         │ │  │
│  │  │   + Offline Queue      │  │    + Sync Queue       │ │  │
│  │  └────────────┬───────────┘  └───────────────────────┘ │  │
│  └───────────────┼───────────────────────────────────────┘  │
└──────────────────┼──────────────────────────────────────────┘
                   │ HTTPS (REST JSON)
                   ▼
┌──────────────────────────────────────────────────────────────┐
│              Traefik Reverse Proxy (Docker)                   │
│              TLS termination (Let's Encrypt)                  │
│              Rate limiting, CORS, HTTPS redirect              │
└──────────────────┬───────────────────────────────────────────┘
                   │
                   ▼
┌──────────────────────────────────────────────────────────────┐
│                    Go API Server (Gin)                        │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                   Middleware Chain                       │ │
│  │  CORS → Logger → Recovery → RateLimit → JWT Auth        │ │
│  │  → Role Guard → Request Validation                      │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                              │
│  ┌───────────────────┐  ┌───────────────────────────────┐   │
│  │   Route Groups    │  │    Background Workers         │   │
│  │                   │  │                               │   │
│  │  /api/auth/*      │  │  • Push Notification Sender   │   │
│  │  /api/admin/*     │  │  • Reminder Scheduler         │   │
│  │  /api/nutritionist│  │    (meal/med/water times)     │   │
│  │  /api/client/*    │  │  • Stale OTP Cleanup          │   │
│  │  /api/foods/*     │  │                               │   │
│  │  /api/medications │  └───────────────────────────────┘   │
│  │  /api/messages/*  │                                      │
│  │  /api/push/*      │                                      │
│  └────────┬──────────┘                                      │
│           │                                                  │
│  ┌────────▼──────────┐                                      │
│  │   Handler Layer   │  ← HTTP concern: parse, validate,    │
│  │   (internal/      │     serialize, status codes           │
│  │    handler/)      │                                      │
│  └────────┬──────────┘                                      │
│           │                                                  │
│  ┌────────▼──────────┐                                      │
│  │   Service Layer   │  ← Business logic: authorization,    │
│  │   (internal/      │     computation, orchestration,       │
│  │    service/)      │     domain rules                     │
│  └────────┬──────────┘                                      │
│           │                                                  │
│  ┌────────▼──────────┐                                      │
│  │  Repository Layer │  ← Data access: SQL queries,         │
│  │  (internal/       │     batch loading, pgx/pgxpool       │
│  │   repository/)    │                                      │
│  └────────┬──────────┘                                      │
│           │                                                  │
└───────────┼──────────────────────────────────────────────────┘
            │
     ┌──────┴───────┐
     │              │
     ▼              ▼
┌──────────┐  ┌───────────┐
│PostgreSQL│  │   File    │
│  16      │  │  Storage  │
│          │  │ /data/    │
│ pg_trgm  │  │ uploads/  │
│ UUID ext │  │           │
└──────────┘  └───────────┘
```

---

## Component Responsibilities

### 1. Traefik Reverse Proxy

| Aspect | Detail |
|--------|--------|
| **Responsibility** | TLS termination, HTTPS redirect, routing to Go API and Nuxt SSR, basic rate limiting |
| **Communicates with** | Go API Server, Nuxt SSR server (if SSR mode), Let's Encrypt ACME |
| **Boundary** | External world ↔ Internal Docker network. No application logic. |
| **Key config** | Auto-discovery via Docker labels; wildcard routing `Host(...)` rules |

### 2. Go API Server (Gin)

The backend is the largest component. It follows a **3-layer architecture** (Handler → Service → Repository) which is the idiomatic Go pattern for domain-driven web services.

#### 2a. Handler Layer (`internal/handler/`)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | HTTP concerns only: parse request bodies (Gin `ShouldBindJSON`), extract path/query params, call service methods, serialize responses, set HTTP status codes |
| **Communicates with** | Service layer (calls), Gin context (reads/writes) |
| **Boundary** | Knows HTTP, knows JSON. Does NOT contain business logic or SQL. |
| **Key pattern** | One file per domain: `auth_handler.go`, `food_handler.go`, `plan_handler.go`, `tracking_handler.go`, `message_handler.go` |

#### 2b. Service Layer (`internal/service/`)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Business rules, authorization enforcement (nutritionist owns client), domain computations (nutritional totals), workflow orchestration (archive old plan → create new plan) |
| **Communicates with** | Repository layer (calls for data), other services (e.g., PlanService calls FoodRepository for nutrition data), SMS adapter, Push notification adapter |
| **Boundary** | Knows domain models. Does NOT know HTTP or SQL. Receives and returns domain structs. |
| **Key pattern** | Interfaces for external dependencies (SMS, Push) enabling mock testing |

#### 2c. Repository Layer (`internal/repository/`)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | All SQL queries via pgx/pgxpool. Parameterized queries only. Translates between database rows and domain models. |
| **Communicates with** | PostgreSQL (via pgxpool connection) |
| **Boundary** | Knows SQL and pgx. Returns domain model structs. Does NOT know HTTP or business rules. |
| **Key pattern** | One file per aggregate root: `user_repo.go`, `food_repo.go`, `plan_repo.go` (includes days/meals/options/items), `tracking_repo.go`, `message_repo.go` |

#### 2d. Model Layer (`internal/model/`)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Domain struct definitions, enums, request/response DTOs |
| **Communicates with** | Referenced by all layers |
| **Key pattern** | Separate domain models from API DTOs. A `DietPlan` domain model has nested `PlanDay` slices; a `CreatePlanRequest` DTO is a flat input structure. |

#### 2e. Middleware Layer (`internal/middleware/`)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Cross-cutting concerns: JWT extraction/validation, role-based access control, request logging, rate limiting, CORS |
| **Communicates with** | Gin router (injected via `Use()`), JWT/token service |
| **Key pattern** | Gin route groups with stacked middleware: public group (no auth), admin group (JWT + admin role), nutritionist group (JWT + nutritionist role), client group (JWT + client role) |

#### 2f. Background Workers

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Scheduled tasks that don't run in request context: reminder notifications (meal times, medication times, water), OTP cleanup, push notification delivery |
| **Communicates with** | Repository layer (reads schedules), Push adapter (sends notifications) |
| **Key pattern** | Go goroutines launched at server startup with `context.Context` for graceful shutdown. Minute-tick loop for reminder checks. |

### 3. PostgreSQL 16

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Primary persistent data store. 20+ tables. Full-text search via pg_trgm. |
| **Communicates with** | Go API Server only (via pgxpool TCP connection within Docker network) |
| **Boundary** | No direct external access. All access through repository layer. |
| **Key extensions** | `pg_trgm` (Persian fuzzy text search), `uuid-ossp` or `gen_random_uuid()` (UUID PKs) |
| **Key constraints** | Partial unique index on `diet_plans(client_id) WHERE status='active'` (one active plan per client). Foreign keys cascade on plan deletion. `local_id` unique constraints for offline dedup. |

### 4. Nuxt 4 PWA Frontend

The frontend runs as a mobile-first SPA/SSR hybrid with PWA capabilities.

#### 4a. Pages & Layouts (`app/pages/`, `app/layouts/`)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Route-based views, role-specific layouts (admin, nutritionist, client) with bottom navigation |
| **Key pattern** | File-based routing. Three layout files: `admin.vue`, `nutritionist.vue`, `client.vue`. Route middleware redirects based on JWT role. |

#### 4b. Composables (`app/composables/`)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Reusable stateful logic: `useAuth()`, `useApi()` (ofetch wrapper with auth headers + offline queue), `useSyncManager()`, `useShamsiDate()`, `usePlanDays()`, `useMealBuilder()`, `useNutritionCalc()`, `useMessagePolling()` |
| **Communicates with** | Pinia stores, Dexie.js, ofetch/useFetch |
| **Key pattern** | Composables are the bridge between components and data. They manage side effects (API calls, IndexedDB reads) and expose reactive state. |

#### 4c. Pinia Stores (`app/stores/` or `app/composables/`)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Client-side state management per domain: `authStore`, `planStore`, `foodStore`, `waterStore`, `sleepStore`, `exerciseStore`, `medicationStore`, `measurementStore`, `messageStore` |
| **Communicates with** | API composables, Dexie.js (for offline data) |
| **Key pattern** | Each store follows: `state` (reactive data), `actions` (fetch from API or IndexedDB, submit), `getters` (computed summaries like daily totals). |

#### 4d. Service Worker + Offline Layer

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Cache static assets (cache-first), cache API responses (network-first with stale fallback), intercept failed POST requests and queue to IndexedDB |
| **Communicates with** | Dexie.js (IndexedDB), browser Cache API, Push API |
| **Key pattern** | `@vite-pwa/nuxt` configures Workbox-based service worker. Custom sync logic in `useSyncManager()` composable. Background Sync API registration with polling fallback. |

#### 4e. Dexie.js (IndexedDB)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Client-side persistent storage for offline: active diet plan cache, pending tracking logs (syncQueue), cached messages, plan-related food items |
| **Communicates with** | Composables/stores read/write; sync manager drains queue to API |
| **Schema** | Tables: `activePlan`, `foodLogs`, `waterLogs`, `sleepLogs`, `exerciseLogs`, `medicationLogs`, `bodyMeasurements`, `messages`, `syncQueue` |

### 5. File Storage (`/data/uploads/`)

| Aspect | Detail |
|--------|--------|
| **Responsibility** | Persistent storage for user-uploaded files (lab results, message attachments) |
| **Communicates with** | Go API server reads/writes; served through authenticated download endpoints (never directly exposed) |
| **Structure** | `/data/uploads/lab-results/{client_id}/{uuid}.{ext}`, `/data/uploads/messages/{conversation_id}/{uuid}.{ext}` |

---

## Recommended Project Structure

### Go Backend

```
nutritrack-api/
├── cmd/
│   └── server/
│       └── main.go                    # Entry point: init config, DB, router, start server
├── internal/
│   ├── config/
│   │   └── config.go                  # Env-based config with validation
│   ├── handler/                       # HTTP handlers (one per domain)
│   │   ├── auth_handler.go            # Login, OTP request/verify, refresh
│   │   ├── admin_handler.go           # Nutritionist CRUD, stats
│   │   ├── food_handler.go            # Food CRUD, search, categories
│   │   ├── medication_handler.go      # Medication CRUD
│   │   ├── plan_handler.go            # Diet plan CRUD, days/meals/options
│   │   ├── tracking_handler.go        # All 6 tracking logs (food, water, sleep, exercise, med, body)
│   │   ├── message_handler.go         # Send, list, poll, read receipts
│   │   ├── food_request_handler.go    # Client requests, nutritionist approval
│   │   ├── lab_result_handler.go      # Upload, list, download
│   │   ├── push_handler.go            # Subscribe, unsubscribe, preferences
│   │   └── upload_handler.go          # File upload utility (shared by messages, labs)
│   ├── service/                       # Business logic (one per domain)
│   │   ├── auth_service.go            # JWT issue/refresh, OTP generate/verify, password hash
│   │   ├── user_service.go            # User CRUD, role management
│   │   ├── food_service.go            # Food CRUD with authorization
│   │   ├── medication_service.go
│   │   ├── plan_service.go            # Plan lifecycle, archive, nutrition compute
│   │   ├── tracking_service.go        # All tracking with local_id dedup
│   │   ├── message_service.go         # Authorization (own clients only), unread counts
│   │   ├── food_request_service.go    # Request workflow, approval → food creation
│   │   ├── notification_service.go    # Push send, reminder scheduling logic
│   │   └── file_service.go            # File validation, storage, path generation
│   ├── repository/                    # Data access (one per aggregate)
│   │   ├── user_repo.go
│   │   ├── food_repo.go               # Includes pg_trgm search queries
│   │   ├── medication_repo.go
│   │   ├── plan_repo.go               # Complex: batch loads plan tree in 2-3 queries
│   │   ├── tracking_repo.go           # Generic pattern for all 6 log tables
│   │   ├── message_repo.go
│   │   ├── food_request_repo.go
│   │   ├── lab_result_repo.go
│   │   ├── push_repo.go               # Push subscriptions, notification prefs
│   │   └── otp_repo.go                # OTP codes with TTL
│   ├── model/                         # Domain models + DTOs
│   │   ├── user.go
│   │   ├── food.go
│   │   ├── medication.go
│   │   ├── plan.go                    # DietPlan, PlanDay, Meal, MealOption, MealOptionItem
│   │   ├── tracking.go               # All log types
│   │   ├── message.go
│   │   ├── notification.go
│   │   └── dto/                       # Request/response DTOs (separate from domain)
│   │       ├── auth_dto.go
│   │       ├── food_dto.go
│   │       ├── plan_dto.go
│   │       └── ...
│   ├── middleware/
│   │   ├── auth.go                    # JWT extraction + validation
│   │   ├── role_guard.go              # Role-based route protection
│   │   ├── rate_limit.go              # Sliding window rate limiter
│   │   ├── logger.go                  # Structured JSON logging
│   │   └── cors.go
│   └── worker/                        # Background goroutines
│       ├── reminder_worker.go         # Meal/medication/water reminder scheduler
│       └── cleanup_worker.go          # Expired OTP cleanup
├── pkg/                               # Shared utilities (no domain knowledge)
│   ├── sms/                           # SMS adapter interface + implementations
│   │   ├── sms.go                     # interface Sender { Send(phone, message) error }
│   │   ├── mock.go                    # Logs to stdout
│   │   └── kavenegar.go               # Production SMS gateway
│   ├── push/                          # Web Push adapter
│   │   └── webpush.go                 # VAPID key management, push sending
│   ├── jwt/                           # JWT utilities
│   │   └── jwt.go
│   └── validator/                     # Custom validation helpers
│       └── iranian_mobile.go          # Iranian phone format validation
├── migrations/                        # SQL migration files (golang-migrate format)
│   ├── 001_users.up.sql
│   ├── 001_users.down.sql
│   ├── 002_foods_medications.up.sql
│   ├── 003_diet_plans.up.sql
│   ├── 004_tracking_tables.up.sql
│   ├── 005_messages_requests_labs.up.sql
│   └── 006_push_notifications.up.sql
├── Dockerfile                         # Multi-stage build
├── docker-compose.yml
├── .env.example
└── go.mod
```

### Nuxt 4 Frontend

```
nutritrack-web/
├── app/
│   ├── assets/
│   │   ├── css/
│   │   │   └── main.css               # Tailwind + RTL + Persian fonts
│   │   └── fonts/
│   │       └── Vazirmatn/             # Persian font files
│   ├── components/
│   │   ├── ui/                        # Generic UI components
│   │   │   ├── AppButton.vue
│   │   │   ├── AppInput.vue
│   │   │   ├── AppModal.vue
│   │   │   ├── AppBottomSheet.vue
│   │   │   ├── PaginatedList.vue
│   │   │   ├── SearchInput.vue        # Persian-aware search with debounce
│   │   │   ├── LoadingSpinner.vue
│   │   │   └── EmptyState.vue
│   │   ├── nutrition/                 # Domain-specific shared components
│   │   │   ├── NutritionLabel.vue     # Macro display card (cal/protein/carb/fat)
│   │   │   ├── FoodPicker.vue         # Search + select food item modal
│   │   │   ├── CategoryChip.vue
│   │   │   └── MealCard.vue
│   │   ├── tracking/
│   │   │   ├── WaterProgress.vue      # Water glass fill animation
│   │   │   ├── MedicationChecklist.vue
│   │   │   ├── WeightChart.vue        # Chart.js line chart
│   │   │   └── DailySummaryCard.vue
│   │   ├── chat/
│   │   │   ├── ChatBubble.vue
│   │   │   ├── ChatInput.vue
│   │   │   └── AttachmentPreview.vue
│   │   └── sync/
│   │       └── SyncStatusIndicator.vue # "Syncing..." / "All synced" / "X pending"
│   ├── composables/
│   │   ├── useApi.ts                  # ofetch wrapper: auth headers, 401 refresh, offline queue
│   │   ├── useAuth.ts                 # Login, logout, token management
│   │   ├── useSyncManager.ts          # Offline queue processing, retry logic
│   │   ├── useOfflineDetect.ts        # navigator.onLine + event listeners
│   │   ├── useShamsiDate.ts           # Jalali date formatting (jalaali-js)
│   │   ├── usePlanDays.ts             # Day navigation for diet plan view
│   │   ├── useMealBuilder.ts          # Plan builder orchestration (nutritionist)
│   │   ├── useNutritionCalc.ts        # Client-side nutrition sum computation
│   │   ├── useMessagePolling.ts       # 10s interval polling, start/stop lifecycle
│   │   └── usePersianNumbers.ts       # Latin → Persian numeral conversion
│   ├── layouts/
│   │   ├── admin.vue                  # Super admin layout
│   │   ├── nutritionist.vue           # Nutritionist layout with bottom nav
│   │   ├── client.vue                 # Client layout with bottom nav
│   │   └── auth.vue                   # Login/registration pages (no nav)
│   ├── middleware/
│   │   ├── auth.global.ts             # Redirect unauthenticated to login
│   │   └── role-guard.ts              # Check JWT role matches route prefix
│   ├── pages/
│   │   ├── auth/
│   │   │   ├── login.vue              # Role selector → email/password or OTP
│   │   │   └── otp.vue                # OTP input for clients
│   │   ├── admin/
│   │   │   ├── index.vue              # Dashboard with stats
│   │   │   ├── nutritionists.vue      # Nutritionist management
│   │   │   ├── foods.vue              # Food database
│   │   │   └── medications.vue        # Medication database
│   │   ├── nutritionist/
│   │   │   ├── index.vue              # Client list
│   │   │   ├── clients/
│   │   │   │   └── [id].vue           # Client profile with tabs
│   │   │   ├── plans/
│   │   │   │   ├── create.vue         # Diet plan builder
│   │   │   │   └── [id].vue           # Plan detail/edit
│   │   │   ├── foods.vue              # Food database
│   │   │   ├── medications.vue
│   │   │   ├── messages.vue           # Conversation list
│   │   │   ├── messages/
│   │   │   │   └── [clientId].vue     # Chat with specific client
│   │   │   └── food-requests.vue      # Pending requests
│   │   └── client/
│   │       ├── index.vue              # Daily dashboard (home screen)
│   │       ├── plan.vue               # Active diet plan view
│   │       ├── tracking/
│   │       │   ├── food.vue           # Food log
│   │       │   ├── water.vue          # Water tracker
│   │       │   ├── sleep.vue
│   │       │   ├── exercise.vue
│   │       │   ├── medication.vue     # Medication checklist
│   │       │   └── measurements.vue   # Weight + body measurements
│   │       ├── messages.vue           # Chat with nutritionist
│   │       ├── lab-results.vue        # Upload and list
│   │       └── settings.vue           # Notification preferences
│   ├── plugins/
│   │   ├── dexie.client.ts            # Initialize Dexie DB (client-only)
│   │   └── pwa.client.ts              # PWA registration and update prompt
│   ├── stores/                        # Pinia stores
│   │   ├── auth.ts
│   │   ├── plan.ts
│   │   ├── food.ts
│   │   ├── tracking.ts               # Composite store for all tracking types
│   │   └── message.ts
│   ├── utils/
│   │   └── constants.ts               # Measurement units, categories, enums (mirrors Go)
│   ├── app.vue                        # Root component
│   └── app.config.ts
├── shared/
│   ├── types/
│   │   ├── api.ts                     # API response types
│   │   ├── food.ts                    # Food, Category, MeasurementUnit
│   │   ├── plan.ts                    # DietPlan, PlanDay, Meal, MealOption, MealOptionItem
│   │   ├── tracking.ts               # All log types
│   │   └── user.ts
│   └── utils/
│       └── nutrition.ts               # Shared nutrition calculation (isomorphic)
├── server/                            # Nuxt server (minimal — Go is the real API)
│   └── api/                           # Only if SSR proxy needed; otherwise empty
├── public/
│   ├── icons/                         # PWA icons (multiple sizes)
│   └── manifest.webmanifest
├── nuxt.config.ts                     # Tailwind, PWA, RTL, API proxy config
├── tailwind.config.ts                 # RTL plugin, Persian font, mobile-first
├── Dockerfile
└── package.json
```

---

## Architectural Patterns to Follow

### Pattern 1: Layered Service Architecture (Go Backend)

**What:** Three clean layers — Handler → Service → Repository — with dependency injection via constructor functions. Each layer depends only on the layer below it via interfaces.

**Why:** This is the dominant Go web service pattern. It keeps business logic testable without HTTP or database dependencies. Gin's middleware + route group system maps cleanly to this.

**Implementation:**

```go
// internal/repository/food_repo.go
type FoodRepository interface {
    Create(ctx context.Context, food *model.Food) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.Food, error)
    Search(ctx context.Context, query string, filters FoodFilters, page int) ([]model.Food, int, error)
}

type foodRepo struct {
    pool *pgxpool.Pool
}

func NewFoodRepository(pool *pgxpool.Pool) FoodRepository {
    return &foodRepo{pool: pool}
}

// internal/service/food_service.go
type FoodService struct {
    repo   repository.FoodRepository
    // No HTTP knowledge, no SQL knowledge
}

func NewFoodService(repo repository.FoodRepository) *FoodService {
    return &FoodService{repo: repo}
}

// internal/handler/food_handler.go
type FoodHandler struct {
    svc *service.FoodService
}

func NewFoodHandler(svc *service.FoodService) *FoodHandler {
    return &FoodHandler{svc: svc}
}

func (h *FoodHandler) RegisterRoutes(rg *gin.RouterGroup) {
    foods := rg.Group("/foods")
    foods.GET("", h.List)
    foods.POST("", h.Create)
    foods.GET("/:id", h.GetByID)
    foods.PUT("/:id", h.Update)
    foods.DELETE("/:id", h.Delete)
}
```

**Confidence:** HIGH — Standard Go web architecture verified via Gin docs, pgx docs, and ecosystem conventions.

### Pattern 2: Aggregate Root Loading for Nested Entities (Diet Plan)

**What:** The diet plan entity (Plan → Days → Meals → Options → Items) is loaded as a complete aggregate in 2–3 SQL queries using batch loading, NOT one query per entity level and NOT a single massive JOIN.

**Why:** A single JOIN across 5 tables produces a cartesian explosion (420+ rows for a realistic plan). N+1 queries (one per meal, option, etc.) would mean 50+ queries. The sweet spot: load flat entity lists in 2–3 queries, then assemble the tree in Go memory.

**Implementation:**

```go
// internal/repository/plan_repo.go
func (r *planRepo) GetFullPlan(ctx context.Context, planID uuid.UUID) (*model.DietPlanFull, error) {
    // Query 1: Plan + Days + Exercise Recommendations
    // SELECT p.*, pd.*, er.* FROM diet_plans p
    //   JOIN plan_days pd ON pd.diet_plan_id = p.id
    //   LEFT JOIN exercise_recommendations er ON er.plan_day_id = pd.id
    //   WHERE p.id = $1 ORDER BY pd.day_number

    // Query 2: Meals + Options + Items (with food data)
    // SELECT m.*, mo.*, moi.*, f.name, f.calories, f.protein, f.carbohydrates, f.fat
    //   FROM meals m
    //   JOIN meal_options mo ON mo.meal_id = m.id
    //   JOIN meal_option_items moi ON moi.meal_option_id = mo.id
    //   JOIN foods f ON f.id = moi.food_id
    //   WHERE m.plan_day_id = ANY($1)  -- pass all day IDs from query 1
    //   ORDER BY m.display_order, mo.option_number

    // Query 3: Prescribed medications (with medication names)
    // SELECT pm.*, med.name FROM prescribed_medications pm
    //   JOIN medications med ON med.id = pm.medication_id
    //   WHERE pm.diet_plan_id = $1

    // Assemble in memory: map day_id → meals → options → items
}
```

**Alternative considered:** Using pgx `SendBatch` to execute all 3 queries in a single round-trip. This reduces network overhead further. **Recommendation:** Use `SendBatch` from day one — pgx supports it natively and it keeps the code clean.

**Confidence:** HIGH — pgx batch queries verified via Context7 docs. Aggregate loading is a well-established DDD pattern.

### Pattern 3: Route Group Middleware Stacking (Gin)

**What:** Use Gin's route group system to apply auth middleware at the group level, creating four distinct API zones with different access rules.

**Why:** Verified via Gin docs — route groups with `Use()` apply middleware to all routes in the group. This eliminates per-handler auth checks and ensures no endpoint is accidentally left unprotected.

**Implementation:**

```go
func SetupRouter(
    authMiddleware gin.HandlerFunc,
    roleGuard func(roles ...string) gin.HandlerFunc,
    handlers *Handlers,
) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery(), middleware.Logger(), middleware.CORS())

    // Public routes — no auth
    public := r.Group("/api")
    {
        public.GET("/health", handlers.Health)
        public.POST("/auth/login", handlers.Auth.Login)
        public.POST("/auth/otp/request", handlers.Auth.RequestOTP)
        public.POST("/auth/otp/verify", handlers.Auth.VerifyOTP)
        public.POST("/auth/refresh", handlers.Auth.Refresh)
    }

    // Admin routes — JWT + super_admin role
    admin := r.Group("/api/admin")
    admin.Use(authMiddleware, roleGuard("super_admin"))
    {
        handlers.Admin.RegisterRoutes(admin)
    }

    // Nutritionist routes — JWT + nutritionist role
    nutritionist := r.Group("/api/nutritionist")
    nutritionist.Use(authMiddleware, roleGuard("nutritionist"))
    {
        handlers.Plan.RegisterRoutes(nutritionist)
        handlers.ClientMgmt.RegisterRoutes(nutritionist)
        handlers.FoodRequest.RegisterNutritionistRoutes(nutritionist)
    }

    // Client routes — JWT + client role
    client := r.Group("/api/client")
    client.Use(authMiddleware, roleGuard("client"))
    {
        handlers.Tracking.RegisterRoutes(client)
        handlers.ClientPlan.RegisterRoutes(client)
        handlers.FoodRequest.RegisterClientRoutes(client)
    }

    // Shared authenticated routes — any role with JWT
    shared := r.Group("/api")
    shared.Use(authMiddleware)
    {
        handlers.Food.RegisterRoutes(shared)     // /api/foods
        handlers.Medication.RegisterRoutes(shared)
        handlers.Message.RegisterRoutes(shared)
        handlers.Push.RegisterRoutes(shared)
    }

    return r
}
```

**Confidence:** HIGH — Verified directly against Gin documentation via Context7.

### Pattern 4: Offline-First Sync Queue (Frontend)

**What:** An offline-aware API layer that intercepts failed POST/PUT requests, serializes them to Dexie.js IndexedDB, and replays them in FIFO order when connectivity returns. Uses `local_id` (UUID) for server-side deduplication.

**Why:** Client users (mobile, potentially poor connectivity in Iran) need to log meals, water, exercise etc. without waiting for network. The sync queue pattern is the standard approach for offline-first apps.

**Implementation:**

```typescript
// app/composables/useApi.ts
export function useApi() {
  const { isOnline } = useOfflineDetect()

  async function post<T>(url: string, body: any): Promise<T> {
    // Attach local_id for dedup if not present
    if (!body.local_id) {
      body.local_id = crypto.randomUUID()
    }

    if (!isOnline.value) {
      await db.syncQueue.add({
        id: crypto.randomUUID(),
        method: 'POST',
        url,
        body,
        createdAt: new Date(),
        status: 'pending',
        retries: 0,
      })
      return body as T // Return optimistic response
    }

    try {
      return await $fetch<T>(url, { method: 'POST', body })
    } catch (err) {
      if (isNetworkError(err)) {
        await db.syncQueue.add({ /* same as above */ })
        return body as T
      }
      throw err
    }
  }

  return { get, post, put, delete: del }
}

// app/composables/useSyncManager.ts
export function useSyncManager() {
  const processingQueue = ref(false)
  const pendingCount = ref(0)

  async function processQueue() {
    const items = await db.syncQueue
      .where('status').equals('pending')
      .sortBy('createdAt')

    for (const item of items) {
      try {
        await db.syncQueue.update(item.id, { status: 'syncing' })
        await $fetch(item.url, { method: item.method, body: item.body })
        await db.syncQueue.delete(item.id)
      } catch (err) {
        const retries = item.retries + 1
        if (retries >= 3) {
          await db.syncQueue.update(item.id, { status: 'failed', retries })
        } else {
          await db.syncQueue.update(item.id, { status: 'pending', retries })
          await delay(Math.pow(2, retries) * 1000) // Exponential backoff
        }
      }
    }
  }

  // Watch online status
  watch(isOnline, (online) => {
    if (online) processQueue()
  })
}
```

**Confidence:** HIGH — Dexie.js schema/query patterns verified via Context7. Standard offline-first pattern.

### Pattern 5: Domain-Scoped Pinia Stores + Composable Pairs

**What:** Each feature domain gets one Pinia store (state container) and one composable (side-effect orchestrator). The store holds reactive state; the composable handles API calls, IndexedDB sync, and complex logic.

**Why:** Separating "what data we have" (store) from "how we get/update data" (composable) keeps components clean and enables the offline layer to transparently swap data sources.

**Example:**

```typescript
// app/stores/tracking.ts — State only
export const useTrackingStore = defineStore('tracking', () => {
  const waterLogs = ref<WaterLog[]>([])
  const dailyWaterTotal = computed(() =>
    waterLogs.value.reduce((sum, log) => sum + log.amount_ml, 0)
  )
  const waterTarget = ref<number | null>(null)

  return { waterLogs, dailyWaterTotal, waterTarget }
})

// app/composables/useWaterTracking.ts — Side effects
export function useWaterTracking() {
  const store = useTrackingStore()
  const api = useApi()

  async function addWater(amountMl: number) {
    const entry = { amount_ml: amountMl, date: today(), local_id: crypto.randomUUID() }
    // Optimistic: add to store immediately
    store.waterLogs.push(entry)
    // Persist (online or queued)
    await api.post('/api/client/water-logs', entry)
  }

  async function fetchToday() {
    // Try API, fallback to IndexedDB
    // ...
  }

  return { addWater, fetchToday }
}
```

**Confidence:** HIGH — Nuxt 4 + Pinia is the official recommended state management pattern per Nuxt docs.

---

## Data Flow

### Flow 1: Client Views Diet Plan (Online → Cached)

```
Client opens app
  → auth.global.ts middleware checks JWT
  → client/plan.vue mounts
  → usePlanDays() composable → planStore.fetchActivePlan()
    → useApi().get('/api/client/plan')
      → Go Handler: extracts client_id from JWT
      → Go Service: verifies active plan exists
      → Go Repository: batch loads full plan tree (3 queries via pgx SendBatch)
      → Response: 200 + full plan JSON with computed nutrition
    → Store receives plan, updates reactive state
    → Dexie.js: serialize full plan to `activePlan` table (cache for offline)
  → Vue renders plan days with swipeable navigation
```

### Flow 2: Client Logs Water Intake (Offline-Capable)

```
Client taps "Add Water" button
  → useWaterTracking().addWater(250)
    → Generates local_id UUID
    → Optimistically adds to Pinia store (instant UI update)
    → useApi().post('/api/client/water-logs', { amount_ml: 250, local_id: "..." })
      IF ONLINE:
        → Go Handler: validates, extracts client_id from JWT
        → Go Service: calls repo
        → Go Repository: INSERT with local_id, ON CONFLICT (local_id) DO NOTHING
        → Response: 201 Created
      IF OFFLINE:
        → Catches network error
        → Writes to Dexie.js syncQueue: { method: 'POST', url: '...', body: {...}, status: 'pending' }
        → Shows toast: "ذخیره شد — بعد از اتصال همگام‌سازی می‌شود"
        LATER, when online:
        → useSyncManager detects navigator.onLine → true
        → Processes syncQueue FIFO
        → POST to /api/client/water-logs
        → Backend deduplicates via local_id
        → Deletes from syncQueue on success
```

### Flow 3: Nutritionist Creates Diet Plan

```
Nutritionist opens plan builder
  → nutritionist/plans/create.vue mounts
  → useMealBuilder() composable manages local plan state
  → For each day:
    → Add meals with title, time, order
    → For each meal:
      → Add options
      → For each option:
        → FoodPicker.vue → searches /api/foods?q=... (pg_trgm fuzzy search)
        → Select food → set quantity + unit
        → useNutritionCalc() recomputes totals (client-side, instant)
  → "Save Plan" button
    → useApi().post('/api/nutritionist/clients/:id/plans', fullPlanPayload)
      → Go Handler: validates deeply nested JSON payload
      → Go Service:
        1. Verify nutritionist owns this client
        2. Archive current active plan (UPDATE status='archived' WHERE client_id AND status='active')
        3. Insert new plan + days + meals + options + items in a single DB transaction
        4. Insert prescribed medications
        5. Insert exercise recommendations
        6. Trigger push notification: "برنامه غذایی جدید" to client
      → Go Repository: single pgx transaction with multiple INSERTs
      → Response: 201 + plan ID
```

### Flow 4: Messaging (Polling-Based)

```
Client opens chat
  → client/messages.vue mounts
  → Loads existing messages: GET /api/messages/:nutritionistId?limit=50
  → useMessagePolling() starts 10-second interval
    → Every 10s: GET /api/messages/new?since={lastMessageTimestamp}
      → Go Repository: SELECT WHERE sent_at > $since AND (sender_id OR receiver_id)
      → Returns only new messages (lightweight query)
    → New messages appended to Pinia store + Dexie.js cache
  → Client sends message
    → POST /api/messages { receiver_id, content }
      → Go Service: verifies sender→receiver relationship (client→own nutritionist)
      → Go Repository: INSERT message
      → Go NotificationService: send Web Push to receiver
    → Optimistic: message appears in chat immediately
  → On component unmount: useMessagePolling().stop() clears interval
```

### Flow 5: Push Notification (Meal Reminder)

```
Server-side (background worker):
  → ReminderWorker goroutine ticks every 60 seconds
  → Queries active diet plans for meals/medications with times in next 60 seconds:
    SELECT u.id, ps.endpoint, ps.p256dh, ps.auth, m.title, m.scheduled_time
    FROM meals m
    JOIN plan_days pd ON pd.id = m.plan_day_id
    JOIN diet_plans dp ON dp.id = pd.diet_plan_id
    JOIN users u ON u.id = dp.client_id
    JOIN push_subscriptions ps ON ps.user_id = u.id
    JOIN notification_preferences np ON np.user_id = u.id
    WHERE dp.status = 'active'
      AND np.meal_reminders = true
      AND m.scheduled_time BETWEEN NOW()::time AND (NOW() + INTERVAL '1 minute')::time
      AND (dp.day_number logic matching today's date)
  → For each match:
    → webpush-go sends push to endpoint
    → Service worker receives push event
    → Displays notification: "وقت صبحانه 🍽️"
    → On click: navigates to /client/plan (today's meals)
```

---

## Anti-Patterns to Avoid

### Anti-Pattern 1: N+1 Queries on Diet Plan Load

**What:** Loading the plan, then each day, then each meal per day, then each option per meal, then each item per option — resulting in 50–100+ queries.
**Why bad:** At 420 items, this means ~100 queries per plan load. Response time exceeds 500ms target by 5–10x.
**Instead:** Batch load in 2–3 queries using JOINs or `WHERE id = ANY($1)` with pgx SendBatch, then assemble the tree in Go memory.

### Anti-Pattern 2: Putting Business Logic in Handlers

**What:** Handlers directly querying the database, checking authorization inline, computing nutrition, etc.
**Why bad:** Untestable without spinning up HTTP server. Duplicated logic when multiple endpoints need the same checks (e.g., "nutritionist owns client").
**Instead:** Handlers call service methods. Services contain all business logic. Repositories handle data access. Test services with mock repositories.

### Anti-Pattern 3: Single Massive SQL JOIN for Plan Tree

**What:** `SELECT * FROM diet_plans JOIN plan_days JOIN meals JOIN meal_options JOIN meal_option_items JOIN foods` — a single query joining 6 tables.
**Why bad:** Cartesian product explosion. A plan with 7 days × 5 meals × 3 options × 4 items = 420 leaf rows, but the JOIN produces rows with redundant plan/day/meal columns. Result set becomes huge, deserialization is complex, and the query planner may struggle with this.
**Instead:** 2–3 focused queries (plan+days, meals+options+items, medications) assembled in application code.

### Anti-Pattern 4: Storing Sync Queue in localStorage

**What:** Using `localStorage` for the offline sync queue instead of IndexedDB.
**Why bad:** localStorage has a 5–10MB limit, is synchronous (blocks main thread), and can be cleared by the browser without warning. A queue of messages with attachments or a cached diet plan can easily exceed this.
**Instead:** Dexie.js with IndexedDB provides structured storage, async access, and much larger capacity (typically 50–100MB+).

### Anti-Pattern 5: WebSocket for Messaging

**What:** Implementing WebSocket for real-time messaging instead of polling.
**Why bad for this project:** Adds significant complexity (connection management, reconnection, Docker/Traefik WebSocket proxy config, offline handling) for a feature where 10-second latency is acceptable. Only ~500 concurrent users. WebSocket connection pooling on a single server is wasted infrastructure.
**Instead:** 10-second polling with `GET /api/messages/new?since=` is simple, stateless, and works perfectly with the existing REST architecture and offline queue.

### Anti-Pattern 6: Direct File Serving from Upload Directory

**What:** Mapping `/data/uploads/` directly as a static file route.
**Why bad:** No authentication. Any user can access any other user's lab results or message attachments by guessing file paths.
**Instead:** Serve files through authenticated Go handler endpoints that verify the requesting user has permission to access the file.

---

## Integration Points

### External Services

| Service | Integration Type | Adapter Pattern | Fallback |
|---------|-----------------|----------------|----------|
| SMS Gateway (Kavenegar) | HTTP API call from Go | Interface `sms.Sender` with mock + real impl | Mock adapter logs to stdout in dev |
| Let's Encrypt (ACME) | Traefik auto-config | Docker labels on Traefik container | Self-signed cert for local dev |
| Web Push (VAPID) | `webpush-go` library | Direct library call in notification service | Log notification to stdout in dev |

### Internal Service Communication

| From | To | Method | Notes |
|------|-----|--------|-------|
| Nuxt PWA | Go API | REST JSON over HTTPS | All communication via `/api/*` routes |
| Go API | PostgreSQL | pgxpool TCP | Connection pool (25 max conns for 500 concurrent users) |
| Go API | File Storage | OS filesystem | `/data/uploads/` Docker volume mount |
| Service Worker | Go API | REST (from SW fetch handler) | Same API, but from SW context for background sync |
| Service Worker | Browser | Push API + Notification API | Receives push events, shows notifications |
| Reminder Worker | Push Subscriptions | Web Push protocol | Outbound HTTPS to push endpoints (browser vendor servers) |

### Key Database Indexes (Build-Order Critical)

These indexes directly affect query performance for core features and should be created with their respective migration:

| Table | Index | Purpose | Phase |
|-------|-------|---------|-------|
| `users` | `(mobile) WHERE mobile IS NOT NULL` | OTP lookup | 1 |
| `users` | `(nutritionist_id) WHERE role='client'` | Client list for nutritionist | 1 |
| `foods.name` | GIN `gin_trgm_ops` | Persian fuzzy text search | 2 |
| `foods` | `(is_active)` | Active food filter | 2 |
| `diet_plans` | `(client_id) WHERE status='active'` UNIQUE | One active plan per client | 3 |
| `plan_days` | `(diet_plan_id, day_number)` | Plan tree loading | 3 |
| `meals` | `(plan_day_id, display_order)` | Meal ordering in plan | 3 |
| `meal_options` | `(meal_id, option_number)` | Option ordering | 3 |
| `meal_option_items` | `(meal_option_id)` | Batch loading items per option | 3 |
| All tracking tables | `(client_id, date)` | Date-range queries | 4 |
| All tracking tables | `(local_id) UNIQUE` | Offline dedup | 4 |
| `messages` | `(sender_id, sent_at)` | Polling query | 5 |
| `messages` | `(receiver_id, sent_at)` | Polling query | 5 |
| `messages` | `(receiver_id) WHERE read_at IS NULL` | Unread count | 5 |

---

## Scalability Considerations

| Concern | At 50 users (launch) | At 500 concurrent users (target) | At 10K clients |
|---------|---------------------|----------------------------------|----------------|
| **DB connections** | pgxpool: 10 max | pgxpool: 25 max | pgxpool: 50 max, consider PgBouncer |
| **Plan loading** | 2–3 queries, <100ms | Same, <200ms | Same, cache hot plans with ETag |
| **Message polling** | Negligible | ~50 polls/sec at peak | Add `last_activity` index, consider SSE |
| **Push notifications** | Direct send in goroutine | Buffered channel, 10 workers | Dedicated notification queue (Redis) |
| **File storage** | Local disk adequate | Local disk adequate (~50GB/yr) | S3-compatible (MinIO) if >200GB |
| **Search** | pg_trgm fine for 1K foods | pg_trgm fine for 10K foods | pg_trgm fine for 100K; beyond → dedicated search |
| **Offline sync** | Trivial | FIFO queue handles burst | Consider batch sync API endpoint |

**Bottom line:** The target scale (50 nutritionists, 10K clients, 500 concurrent) is well within PostgreSQL + single Go server capability. No microservices, no Redis, no message broker needed. Keep it simple.

---

## Suggested Build Order (Dependency Graph)

The architecture drives a clear build order based on data dependencies:

```
Phase 1: Foundation
  ├── Go project structure (cmd/, internal/{handler,service,repository,model,middleware}/)
  ├── PostgreSQL + migrations system + users table
  ├── Gin router skeleton with middleware chain
  ├── JWT auth + OTP flow
  ├── Nuxt 4 project with layouts, auth pages, RTL setup
  └── Docker Compose (Go + PG + Traefik)
       │
       ▼
Phase 2: Core Data Domain
  ├── foods table + pg_trgm search + food_categories junction
  ├── medications table
  ├── Food/Medication CRUD (all 3 layers)
  ├── Super Admin panel (frontend)
  └── NutritionLabel, FoodPicker, SearchInput components
       │
       ▼
Phase 3: Diet Plan Engine  ← MOST COMPLEX, highest risk
  ├── 5 new tables: diet_plans, plan_days, meals, meal_options, meal_option_items
  ├── 2 more tables: prescribed_medications, exercise_recommendations
  ├── Aggregate root loading pattern (batch queries)
  ├── Nutrition computation service
  ├── Plan builder UI (nutritionist) — most complex frontend component
  └── Plan viewer UI (client)
       │
       ├─────────────────────┐
       ▼                     ▼
Phase 4: Tracking       Phase 5: Communication
  ├── 6 tracking tables    ├── messages table
  ├── local_id dedup       ├── food_requests table
  ├── Tracking CRUD        ├── lab_results table
  ├── Daily dashboard      ├── Chat UI + 10s polling
  ├── Weight charts        ├── File upload handler
  └── (Chart.js)           ├── Food request workflow
                           └── Client mgmt dashboard
       │                     │
       └─────────┬───────────┘
                 ▼
Phase 6: Offline & PWA
  ├── Service Worker (@vite-pwa/nuxt)
  ├── Dexie.js schema + IndexedDB stores
  ├── Sync queue (useApi offline wrapper)
  ├── Sync manager (useSyncManager)
  ├── Push notifications (webpush-go + SW handler)
  ├── Reminder worker (background goroutine)
  └── Notification preferences
       │
       ▼
Phase 7: Hardening & Launch
  ├── Security audit (row-level auth, SQL injection, XSS)
  ├── Performance (EXPLAIN ANALYZE, load testing)
  ├── Monitoring (Grafana + Loki)
  └── Backup verification
```

**Key dependency insight:** Phases 4 and 5 can run in parallel after Phase 3 because they share no data model dependencies. Phase 6 (offline) must come after all data entry features (Phases 4+5) exist, because the offline layer wraps their API calls.

**Highest risk phase:** Phase 3 (Diet Plan Engine). It has the most complex data model (5-level nesting), the highest-risk query patterns (N+1 trap), and the most complex frontend component (plan builder). Allocate extra time and do a design spike on the aggregate loading pattern before starting implementation.

---

## Sources

- **Gin web framework documentation** — Context7 `/websites/gin-gonic_en`: Route grouping, middleware, model binding/validation patterns. **HIGH confidence.**
- **pgx v5 PostgreSQL driver** — Context7 `/websites/pkg_go_dev_github_com_jackc_pgx_v5`: SendBatch for multi-query round-trips, pgxpool connection management, transaction patterns. **HIGH confidence.**
- **Nuxt 4 documentation** — Context7 `/websites/nuxt_4_x`: Directory structure (`app/` as srcDir), composables, middleware, Pinia state management, server directory. **HIGH confidence.**
- **Dexie.js documentation** — Context7 `/websites/dexie`: Schema versioning, table definitions, query patterns for IndexedDB. **HIGH confidence.**
- **NutriTrack PRD** — Sections 8 (Tech Stack), 9 (Data Model), 6 (Offline Strategy), 7 (Notifications): Definitive source for all data entities, relationships, and architectural decisions. **HIGH confidence.**
- **NutriTrack phases.md** — Implementation guidance on batch loading, pg_trgm search, sync manager pattern, polling vs WebSocket decision. **HIGH confidence.**
