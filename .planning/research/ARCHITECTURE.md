# Architecture Patterns — NutriTrack Go Backend

**Domain:** Nutrition management platform (nutritionist ↔ client)
**Researched:** 2026-04-21
**Confidence:** HIGH (DDD Go patterns are well-established; specifics cross-checked against project constraints)

---

## Recommended Architecture: Layered DDD

Four concentric layers with a strict **inward-only dependency rule**: outer layers import inner layers, never the reverse.

```
┌─────────────────────────────────────────────────────────┐
│  Interface Layer      (HTTP handlers, middleware, DTOs)  │
│  internal/interface/                                     │
├─────────────────────────────────────────────────────────┤
│  Application Layer    (use cases, commands, queries)     │
│  internal/application/                                   │
├─────────────────────────────────────────────────────────┤
│  Domain Layer         (entities, value objects, repos)   │
│  internal/domain/                                        │
├─────────────────────────────────────────────────────────┤
│  Infrastructure Layer (sqlc, Redis, SMS, push, files)    │
│  internal/infrastructure/                                │
└─────────────────────────────────────────────────────────┘
```

**Golden rule:** `domain` imports nothing from this project. `application` imports only `domain`. `infrastructure` implements `domain` interfaces. `interface` wires them together via DI.

---

## Complete Folder Structure

```
nutritrack/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point: wire all layers, start server
│
├── internal/
│   │
│   ├── domain/                        # ── DOMAIN LAYER ──────────────────────────
│   │   │                              # Zero external dependencies
│   │   ├── user/
│   │   │   ├── entity.go              # User, Nutritionist, Client (embedded)
│   │   │   ├── value_objects.go       # Mobile, Email, Role, Gender, HeightCM
│   │   │   ├── repository.go          # UserRepository interface
│   │   │   ├── service.go             # UserDomainService (cross-entity rules)
│   │   │   └── errors.go              # ErrUserNotFound, ErrMobileAlreadyTaken…
│   │   │
│   │   ├── food/
│   │   │   ├── entity.go              # Food aggregate root
│   │   │   ├── value_objects.go       # NutritionalValues, MeasurementUnit, Category
│   │   │   ├── repository.go          # FoodRepository interface
│   │   │   └── errors.go
│   │   │
│   │   ├── medication/
│   │   │   ├── entity.go              # Medication aggregate root
│   │   │   ├── value_objects.go       # MedicationForm, DosageUnit
│   │   │   ├── repository.go          # MedicationRepository interface
│   │   │   └── errors.go
│   │   │
│   │   ├── dietplan/
│   │   │   ├── entity.go              # DietPlan aggregate root
│   │   │   ├── plan_day.go            # PlanDay entity (child of DietPlan)
│   │   │   ├── meal.go                # Meal entity
│   │   │   ├── meal_option.go         # MealOption entity
│   │   │   ├── meal_option_item.go    # MealOptionItem entity
│   │   │   ├── exercise_rec.go        # ExerciseRecommendation entity
│   │   │   ├── prescribed_med.go      # PrescribedMedication entity
│   │   │   ├── value_objects.go       # NutritionalTotals, PlanStatus, DailyRange
│   │   │   ├── service.go             # DietPlanDomainService (archive rules, totals)
│   │   │   ├── repository.go          # DietPlanRepository interface
│   │   │   └── errors.go
│   │   │
│   │   ├── tracking/
│   │   │   ├── food_log.go            # FoodLog entity
│   │   │   ├── water_log.go           # WaterLog entity
│   │   │   ├── sleep_log.go           # SleepLog entity
│   │   │   ├── exercise_log.go        # ExerciseLog entity
│   │   │   ├── medication_log.go      # MedicationLog entity
│   │   │   ├── measurement.go         # BodyMeasurement entity
│   │   │   ├── value_objects.go       # SleepQuality, MeasurementFields
│   │   │   ├── repository.go          # TrackingRepository interface (one per type)
│   │   │   └── errors.go
│   │   │
│   │   ├── labresult/
│   │   │   ├── entity.go              # LabResult aggregate root
│   │   │   ├── value_objects.go       # LabResultType, FileRef
│   │   │   ├── repository.go          # LabResultRepository interface
│   │   │   └── errors.go
│   │   │
│   │   ├── message/
│   │   │   ├── entity.go              # Message aggregate root
│   │   │   ├── value_objects.go       # AttachmentType
│   │   │   ├── repository.go          # MessageRepository interface
│   │   │   └── errors.go
│   │   │
│   │   ├── foodrequest/
│   │   │   ├── entity.go              # FoodRequest aggregate root
│   │   │   ├── value_objects.go       # RequestStatus
│   │   │   ├── repository.go          # FoodRequestRepository interface
│   │   │   └── errors.go
│   │   │
│   │   └── shared/
│   │       ├── id.go                  # NewID() → uuid.UUID helper
│   │       ├── pagination.go          # Page, PageSize, Cursor value objects
│   │       ├── timeutil.go            # TehranLocation(), TodayTehran()
│   │       └── errors.go              # DomainError base type, ErrNotFound sentinel
│   │
│   ├── application/                   # ── APPLICATION LAYER ─────────────────────
│   │   │                              # Imports: domain only
│   │   ├── auth/
│   │   │   ├── commands.go            # LoginByOTPCmd, RefreshTokenCmd, LogoutCmd
│   │   │   ├── queries.go             # (none — auth is command-heavy)
│   │   │   ├── service.go             # AuthAppService
│   │   │   └── dto.go                 # TokenPairDTO, OTPSentDTO
│   │   │
│   │   ├── user/
│   │   │   ├── commands.go            # CreateClientCmd, UpdateClientCmd, DeactivateCmd
│   │   │   ├── queries.go             # GetClientQuery, ListClientsQuery
│   │   │   ├── service.go             # UserAppService
│   │   │   └── dto.go                 # ClientProfileDTO, ClientListItemDTO
│   │   │
│   │   ├── food/
│   │   │   ├── commands.go            # CreateFoodCmd, UpdateFoodCmd, DeleteFoodCmd
│   │   │   ├── queries.go             # SearchFoodsQuery, GetFoodQuery
│   │   │   ├── service.go             # FoodAppService
│   │   │   └── dto.go                 # FoodDTO, FoodListDTO
│   │   │
│   │   ├── medication/
│   │   │   ├── commands.go            # CreateMedicationCmd, UpdateMedicationCmd
│   │   │   ├── queries.go             # ListMedicationsQuery, GetMedicationQuery
│   │   │   ├── service.go             # MedicationAppService
│   │   │   └── dto.go                 # MedicationDTO
│   │   │
│   │   ├── dietplan/
│   │   │   ├── commands.go            # CreateDietPlanCmd, AddPlanDayCmd, AddMealCmd
│   │   │   │                          # AddMealOptionCmd, AddMealOptionItemCmd
│   │   │   │                          # AddExerciseRecCmd, PrescribeMedCmd
│   │   │   ├── queries.go             # GetActivePlanQuery, GetPlanByIDQuery, ListPlansQuery
│   │   │   ├── service.go             # DietPlanAppService
│   │   │   └── dto.go                 # DietPlanDTO, PlanDayDTO, MealDTO, etc.
│   │   │
│   │   ├── tracking/
│   │   │   ├── commands.go            # LogFoodCmd, LogWaterCmd, LogSleepCmd
│   │   │   │                          # LogExerciseCmd, LogMedicationCmd, RecordMeasurementCmd
│   │   │   ├── queries.go             # GetDailyLogsQuery, GetHistoryQuery
│   │   │   ├── service.go             # TrackingAppService
│   │   │   └── dto.go                 # DailyTrackingSummaryDTO, HistoryDTO
│   │   │
│   │   ├── labresult/
│   │   │   ├── commands.go            # UploadLabResultCmd, DeleteLabResultCmd
│   │   │   ├── queries.go             # ListLabResultsQuery
│   │   │   ├── service.go             # LabResultAppService
│   │   │   └── dto.go                 # LabResultDTO
│   │   │
│   │   ├── message/
│   │   │   ├── commands.go            # SendMessageCmd
│   │   │   ├── queries.go             # GetConversationQuery, GetUnreadCountQuery
│   │   │   ├── service.go             # MessageAppService
│   │   │   └── dto.go                 # MessageDTO, ConversationDTO
│   │   │
│   │   ├── foodrequest/
│   │   │   ├── commands.go            # SubmitFoodRequestCmd, ApproveRequestCmd, RejectRequestCmd
│   │   │   ├── queries.go             # ListPendingRequestsQuery
│   │   │   ├── service.go             # FoodRequestAppService
│   │   │   └── dto.go                 # FoodRequestDTO
│   │   │
│   │   ├── admin/
│   │   │   ├── commands.go            # CreateNutritionistCmd, ToggleNutritionistCmd
│   │   │   ├── queries.go             # ListNutritionistsQuery, PlatformStatsQuery
│   │   │   ├── service.go             # AdminAppService
│   │   │   └── dto.go                 # NutritionistDTO, PlatformStatsDTO
│   │   │
│   │   └── ports/                     # Driven ports (interfaces used by application)
│   │       ├── sms.go                 # SMSSender interface
│   │       ├── push.go                # PushNotifier interface
│   │       ├── storage.go             # FileStorage interface
│   │       ├── token.go               # TokenManager interface (JWT issue/verify)
│   │       └── otp.go                 # OTPStore interface (Redis-backed TTL store)
│   │
│   ├── infrastructure/                # ── INFRASTRUCTURE LAYER ──────────────────
│   │   │                              # Implements domain repo interfaces + app ports
│   │   ├── postgres/
│   │   │   ├── db.go                  # *sql.DB / pgxpool setup, connection string
│   │   │   └── migrations/            # golang-migrate SQL files (also at root /migrations)
│   │   │
│   │   ├── sqlc/                      # Generated code lives here (DO NOT EDIT)
│   │   │   ├── db.go                  # sqlc boilerplate
│   │   │   ├── models.go              # Generated DB row structs
│   │   │   ├── querier.go             # Generated interface
│   │   │   └── *.sql.go               # Generated query implementations
│   │   │
│   │   ├── repository/                # Adapters: sqlc → domain repository interfaces
│   │   │   ├── user_repo.go           # Implements domain/user.UserRepository
│   │   │   ├── food_repo.go           # Implements domain/food.FoodRepository
│   │   │   ├── medication_repo.go
│   │   │   ├── dietplan_repo.go       # Implements domain/dietplan.DietPlanRepository
│   │   │   ├── tracking_repo.go       # Implements all tracking repo interfaces
│   │   │   ├── labresult_repo.go
│   │   │   ├── message_repo.go
│   │   │   └── foodrequest_repo.go
│   │   │
│   │   ├── redis/
│   │   │   ├── client.go              # go-redis client setup
│   │   │   ├── otp_store.go           # Implements application/ports.OTPStore
│   │   │   ├── token_store.go         # Refresh token blacklist / session store
│   │   │   └── rate_limiter.go        # Sliding window rate limit (OTP, API)
│   │   │
│   │   ├── sms/
│   │   │   ├── adapter.go             # Implements application/ports.SMSSender
│   │   │   ├── kavenegar.go           # Kavenegar concrete provider
│   │   │   └── melipayamak.go         # Melipayamak concrete provider (optional)
│   │   │
│   │   ├── push/
│   │   │   ├── adapter.go             # Implements application/ports.PushNotifier
│   │   │   └── webpush.go             # github.com/SherClockHolmes/webpush-go wrapper
│   │   │
│   │   ├── storage/
│   │   │   ├── adapter.go             # Implements application/ports.FileStorage
│   │   │   └── local.go               # Local filesystem: /data/uploads/{type}/{uuid}.ext
│   │   │
│   │   └── jwt/
│   │       ├── manager.go             # Implements application/ports.TokenManager
│   │       └── claims.go              # Custom JWT claims (user_id, role, jti)
│   │
│   └── interface/                     # ── INTERFACE LAYER ────────────────────────
│       │                              # Imports: application layer only
│       ├── http/
│       │   ├── server.go              # gin.Engine setup, route registration
│       │   ├── middleware/
│       │   │   ├── auth.go            # JWT extraction + validation, ctx injection
│       │   │   ├── role.go            # RequireRole(roles...) guard
│       │   │   ├── logging.go         # Structured request/response logging (zerolog)
│       │   │   ├── ratelimit.go       # Per-IP + per-user rate limiting via Redis
│       │   │   ├── cors.go            # CORS for PWA origin
│       │   │   ├── recovery.go        # Panic recovery → 500 with Persian message
│       │   │   └── requestid.go       # X-Request-ID injection
│       │   │
│       │   ├── handler/
│       │   │   ├── auth.go            # POST /auth/otp/send, POST /auth/otp/verify
│       │   │   │                      # POST /auth/login, POST /auth/refresh, POST /auth/logout
│       │   │   ├── user.go            # Nutritionist client management endpoints
│       │   │   ├── food.go            # CRUD + search for foods
│       │   │   ├── medication.go      # CRUD for medications
│       │   │   ├── dietplan.go        # Diet plan builder endpoints
│       │   │   ├── tracking.go        # All client tracking log endpoints
│       │   │   ├── labresult.go       # Lab result upload/list
│       │   │   ├── message.go         # Send/poll messages, unread count
│       │   │   ├── foodrequest.go     # Submit / approve / reject food requests
│       │   │   ├── notification.go    # Push subscription registration
│       │   │   └── admin.go           # Super admin nutritionist management
│       │   │
│       │   └── dto/
│       │       ├── request/           # Inbound JSON binding structs + validation tags
│       │       │   ├── auth.go        # SendOTPRequest, VerifyOTPRequest, LoginRequest
│       │       │   ├── food.go        # CreateFoodRequest, UpdateFoodRequest, SearchFoodRequest
│       │       │   ├── dietplan.go    # CreatePlanRequest, AddDayRequest, …
│       │       │   └── tracking.go    # LogFoodRequest, LogWaterRequest, …
│       │       └── response/          # Outbound JSON structs
│       │           ├── envelope.go    # APIResponse{data, message, errors, meta}
│       │           ├── error.go       # ErrorResponse{code, message(Persian), details}
│       │           └── pagination.go  # PaginatedResponse{items, total, page, page_size}
│       │
│       └── bootstrap/
│           └── wire.go                # Manual DI wiring: build all services, inject deps
│
├── migrations/                        # golang-migrate versioned SQL files
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_foods.up.sql
│   ├── 000002_create_foods.down.sql
│   ├── 000003_create_medications.up.sql
│   ├── 000003_create_medications.down.sql
│   ├── 000004_create_diet_plans.up.sql
│   ├── 000004_create_diet_plans.down.sql
│   ├── 000005_create_tracking_tables.up.sql
│   ├── 000005_create_tracking_tables.down.sql
│   ├── 000006_create_messages.up.sql
│   ├── 000006_create_messages.down.sql
│   ├── 000007_create_lab_results.up.sql
│   ├── 000007_create_lab_results.down.sql
│   ├── 000008_create_food_requests.up.sql
│   ├── 000008_create_food_requests.down.sql
│   └── 000009_create_push_subscriptions.up.sql
│
├── db/
│   └── queries/                       # sqlc source SQL query files (one per domain)
│       ├── users.sql
│       ├── foods.sql
│       ├── medications.sql
│       ├── diet_plans.sql
│       ├── plan_days.sql
│       ├── meals.sql
│       ├── meal_options.sql
│       ├── tracking.sql
│       ├── messages.sql
│       ├── lab_results.sql
│       └── food_requests.sql
│
├── config/
│   ├── config.go                      # Config struct loaded from env vars
│   └── config.example.env             # Template with all required vars documented
│
├── pkg/                               # Reusable, project-agnostic utilities
│   ├── apperror/
│   │   ├── codes.go                   # Error code constants (E4001, E4004, …)
│   │   └── persian.go                 # PersianMessages map[ErrorCode]string
│   ├── logger/
│   │   └── logger.go                  # zerolog setup: JSON output, Tehran timestamp
│   ├── validator/
│   │   └── validator.go               # go-playground/validator + Persian field names
│   └── persian/
│       ├── numbers.go                 # Arabic-Indic ↔ ASCII digit normalisation
│       └── mobile.go                  # Iranian mobile number validation (+98 / 09xx)
│
├── sqlc.yaml                          # sqlc config pointing db/queries/ → internal/infrastructure/sqlc/
├── docker-compose.yml                 # App + Postgres + Redis + Traefik
├── docker-compose.override.yml        # Local dev overrides (hot reload, debug port)
├── Dockerfile
├── Makefile                           # migrate-up, migrate-down, sqlc-gen, lint, test, build
└── .env.example
```

---

## Component Boundaries

### Dependency Graph (arrows = "imports")

```
cmd/api/main.go
    │
    ▼
internal/interface/bootstrap/wire.go   ← assembles everything
    │
    ├──► internal/interface/http/      (handlers + middleware)
    │         │
    │         ▼
    ├──► internal/application/         (use-case services)
    │         │
    │         ├──► internal/domain/    (entities, VOs, repo interfaces)
    │         │         ▲ (implements)
    │         └──► internal/infrastructure/  (concrete adapters)
    │
    └──► pkg/                          (shared utilities — no internal imports)
```

**What is forbidden:**
- `domain` → anything in `internal/` (zero upstream imports)
- `application` → `infrastructure` or `interface` (never touches DB directly)
- `infrastructure` → `application` or `interface` (only implements interfaces)
- `domain` → `pkg/apperror` is acceptable (pkg is project-agnostic)

### Component Responsibilities

| Component | Owns | Communicates With |
|-----------|------|-------------------|
| `domain/user` | User entity, Mobile VO, repo interface | Nothing |
| `domain/dietplan` | DietPlan aggregate, nutrition total computation | `domain/food` (read-only for item data) |
| `domain/tracking` | All tracking log entities | `domain/dietplan` (validates meal/option references) |
| `application/auth` | OTP flow, JWT issuance, token refresh | `domain/user`, ports: OTPStore, SMSSender, TokenManager |
| `application/dietplan` | Plan creation, archiving previous plan | `domain/dietplan`, `domain/food`, `domain/medication` |
| `application/tracking` | Idempotent log writes (local_id dedup) | `domain/tracking`, `domain/dietplan` |
| `infrastructure/repository` | sqlc query execution, DB ↔ domain mapping | `infrastructure/sqlc`, `domain/*/repository` |
| `infrastructure/redis` | OTP TTL, refresh token store, rate limits | `application/ports` |
| `interface/http/handler` | HTTP binding, validation, response shaping | `application/*` services only |
| `interface/bootstrap` | Construct dependency graph, run migrations | All layers |

---

## DDD Aggregate Boundaries

### Aggregate 1: User

```
User (aggregate root)
├── id: UUID
├── role: Role (value object: super_admin | nutritionist | client)
├── fullName: string
├── email: Email (value object — nullable)
├── passwordHash: string (nullable)
├── mobile: Mobile (value object — nullable, Iranian format)
├── dateOfBirth: time.Time (nullable)
├── heightCM: float64 (nullable)
├── gender: Gender (value object)
├── nutritionistID: UUID (nullable FK — not embedded, just a reference)
├── isActive: bool
└── timestamps
```

**Invariants enforced in entity:**
- A client must have a `mobile` and a `nutritionistID`
- A nutritionist/admin must have an `email` and `passwordHash`
- `mobile` must match Iranian mobile format

**Cross-aggregate rule (domain service):** "Only one active diet plan per client" — enforced in `DietPlanDomainService`, not in User entity.

---

### Aggregate 2: Food

```
Food (aggregate root)
├── id: UUID
├── name: string (Persian)
├── categories: []FoodCategory (value object set)
├── nutritionalValues: NutritionalValues (value object)
│   ├── calories: float64
│   ├── protein: float64
│   ├── carbohydrates: float64
│   ├── fat: float64
│   ├── fiber: *float64
│   ├── sugar: *float64
│   └── sodium: *float64
├── measurementUnit: MeasurementUnit (value object)
├── measurementAmount: float64
├── description: *string
├── isActive: bool
└── createdBy: UUID (reference to User — not embedded)
```

**No child entities** — Food is a small aggregate. Categories are a value-object set stored in a junction table.

---

### Aggregate 3: Medication

```
Medication (aggregate root)
├── id: UUID
├── name: string
├── genericName: *string
├── form: MedicationForm (value object)
├── dosageUnit: string
├── description: *string
├── isActive: bool
└── createdBy: UUID
```

---

### Aggregate 4: DietPlan ← Most Complex Aggregate

```
DietPlan (aggregate root)
├── id: UUID
├── clientID: UUID (reference — not embedded)
├── nutritionistID: UUID (reference)
├── period: DateRange (value object: startDate, endDate)
├── dailyWaterTargetML: *int
├── notes: *string
├── status: PlanStatus (value object: active | archived)
│
├── planDays: []PlanDay (child entities, owned by DietPlan)
│   └── PlanDay
│       ├── id: UUID
│       ├── dayNumber: int
│       ├── meals: []Meal (child entities)
│       │   └── Meal
│       │       ├── id: UUID
│       │       ├── title: string
│       │       ├── scheduledTime: time.Time
│       │       ├── order: int
│       │       └── options: []MealOption (child entities)
│       │           └── MealOption
│       │               ├── id: UUID
│       │               ├── optionNumber: int
│       │               └── items: []MealOptionItem (child entities)
│       │                   └── MealOptionItem
│       │                       ├── id: UUID
│       │                       ├── foodID: UUID (reference to Food aggregate)
│       │                       ├── quantity: float64
│       │                       ├── measurementUnit: MeasurementUnit
│       │                       └── notes: *string
│       └── exerciseRecs: []ExerciseRecommendation (child entities)
│
└── prescribedMedications: []PrescribedMedication (child entities)
    └── PrescribedMedication
        ├── id: UUID
        ├── medicationID: UUID (reference to Medication aggregate)
        ├── dosage: string
        ├── frequency: string
        ├── times: []time.Time
        ├── instructions: *string
        └── duration: *DateRange
```

**Key domain methods on DietPlan:**
```go
func (dp *DietPlan) Archive() error           // sets status = archived
func (dp *DietPlan) AddPlanDay(day PlanDay)   // validates day number uniqueness
func (dp *DietPlan) ComputeDailyTotals(dayNumber int, foodLookup FoodLookupFunc) NutritionalRange
func (dp *DietPlan) IsActive() bool
```

**Invariant:** Only the DietPlan aggregate root controls state transitions. The `DietPlanDomainService.CreatePlan()` archives the previous active plan before persisting the new one — within a single DB transaction.

---

### Aggregate 5: TrackingRecord (per-type, thin aggregates)

Each tracking type is a **separate small aggregate** (not one giant tracking aggregate):

```
FoodLog        { id, clientID, date, mealID, selectedOptionID*, loggedAt, notes, localID }
WaterLog       { id, clientID, date, amountML, time*, loggedAt, localID }
SleepLog       { id, clientID, date, sleepTime, wakeTime, quality*, notes, localID }
ExerciseLog    { id, clientID, date, exerciseName, durationMinutes, caloriesBurned*, notes, localID }
MedicationLog  { id, clientID, date, prescribedMedID*, medicationName, dosage, takenAt, notes, localID }
BodyMeasurement{ id, clientID, date, weight*, waist*, hip*, abdomen*, thigh*, chest*, wrist*, recordedBy }
```

`localID` is a client-assigned UUID for **offline sync idempotency** — the server uses `ON CONFLICT (client_id, local_id) DO NOTHING`.

---

### Aggregate 6: Message

```
Message (aggregate root)
├── id: UUID
├── senderID: UUID
├── receiverID: UUID
├── content: *string
├── attachmentType: *AttachmentType (value object: image | file)
├── attachmentPath: *string
├── sentAt: time.Time
└── readAt: *time.Time
```

**Invariant (domain service):** Client can only message their own nutritionist; nutritionist can only message own clients. Checked against User aggregate in `MessageDomainService`.

---

### Aggregate 7: FoodRequest

```
FoodRequest (aggregate root)
├── id: UUID
├── foodName: string
├── description: *string
├── status: RequestStatus (pending | approved | rejected)
├── requestedBy: UUID (clientID)
├── reviewedBy: *UUID (nutritionistID)
├── rejectionReason: *string
└── createdAt: time.Time
```

**Domain methods:**
```go
func (fr *FoodRequest) Approve(reviewerID uuid.UUID) error
func (fr *FoodRequest) Reject(reviewerID uuid.UUID, reason string) error
```

---

## Data Flow for Key Operations

### Flow 1: OTP Login (Client)

```
POST /auth/otp/send { mobile: "09123456789" }
    │
    ▼
interface/http/handler/auth.go
    ├── Bind + validate request DTO
    ├── Normalize mobile (pkg/persian/mobile.go)
    └── Call application/auth.AuthAppService.SendOTP(cmd)
            │
            ├── domain/user.UserRepository.FindByMobile(mobile)     → check client exists & is active
            ├── application/ports.OTPStore.GetAttempts(mobile)      → check rate limit (≤3 per 10min)
            ├── Generate 6-digit OTP
            ├── application/ports.OTPStore.Save(mobile, otp, 2min)  → Redis SET with TTL
            └── application/ports.SMSSender.Send(mobile, otp)       → Kavenegar HTTP call
                    │
                    ▼
            Return OTPSentDTO { expiresInSeconds: 120 }
    │
    ▼
response/envelope.go → 200 { data: { expires_in: 120 }, message: "کد تأیید ارسال شد" }

POST /auth/otp/verify { mobile: "09123456789", otp: "481923" }
    │
    ▼
AuthAppService.VerifyOTP(cmd)
    ├── OTPStore.Verify(mobile, otp)                   → Redis GET + compare + DEL on success
    ├── UserRepository.FindByMobile(mobile)
    ├── TokenManager.IssueAccessToken(userID, role)    → JWT HS256, 15min
    ├── TokenManager.IssueRefreshToken(userID)         → opaque UUID stored in Redis, 30d TTL
    └── Return TokenPairDTO { access_token, refresh_token, expires_in }
```

---

### Flow 2: Diet Plan Creation (Nutritionist)

```
POST /diet-plans { client_id, start_date, end_date, notes, … }
    │
    ▼
handler/dietplan.go
    ├── Extract nutritionistID from JWT claims (middleware/auth.go put it in ctx)
    ├── Bind + validate CreatePlanRequest
    └── Call DietPlanAppService.CreateDietPlan(cmd)
            │
            ├── UserRepository.FindClientByID(clientID)             → verify client belongs to nutritionist
            ├── DietPlanRepository.FindActiveByClientID(clientID)   → get existing active plan (if any)
            ├── DietPlanDomainService.CreateWithArchive(existing, newPlan)
            │       ├── existing.Archive()                          → domain method, sets status=archived
            │       ├── Validate new plan dates (start < end, future dates)
            │       └── Return (planToArchive, newPlan)
            ├── Begin DB transaction
            │   ├── DietPlanRepository.Update(planToArchive)        → set status=archived
            │   └── DietPlanRepository.Save(newPlan)                → insert new active plan
            └── Commit transaction
                    │
                    ▼
            Return DietPlanDTO (with computed nutritional totals = 0, empty days)
```

---

### Flow 3: Offline Sync — Log Water Intake

```
POST /tracking/water { date, amount_ml, time, local_id: "uuid-from-client" }
    │
    ▼
handler/tracking.go → TrackingAppService.LogWater(cmd)
    ├── Validate client ownership (clientID from JWT must match)
    ├── Build WaterLog entity { …, localID: cmd.LocalID }
    └── TrackingRepository.SaveWaterLog(log)
            │
            ▼
        infrastructure/repository/tracking_repo.go
            └── sqlc: INSERT INTO water_logs (…, local_id)
                       ON CONFLICT (client_id, local_id) DO NOTHING
                       RETURNING id
                  ─ If RETURNING returns empty, it was a duplicate → return existing record ID
```

**Idempotency guarantee:** The `(client_id, local_id)` unique constraint handles duplicate sync submissions without error.

---

### Flow 4: Push Notification on New Message

```
MessageAppService.SendMessage(cmd)
    ├── Validate sender/receiver relationship (domain service)
    ├── Create Message entity
    ├── MessageRepository.Save(message)
    └── PushNotifier.Send(receiverID, payload)          ← async, non-blocking
            │
            ▼
    infrastructure/push/webpush.go
        ├── Load receiver's push subscription from DB
        └── webpush.SendNotification(subscription, payload, options)
                  Payload: { title: "پیام جدید", body: senderName, url: "/messages" }
```

Push is fire-and-forget — failure is logged but does not fail the message save.

---

## Configuration Management

```go
// config/config.go — envconfig pattern (no viper: simpler, no YAML needed)
type Config struct {
    // Server
    Port        string `env:"PORT" envDefault:"8080"`
    Environment string `env:"ENVIRONMENT" envDefault:"production"`

    // Database
    DatabaseURL string `env:"DATABASE_URL,required"`

    // Redis
    RedisURL    string `env:"REDIS_URL,required"`

    // JWT
    JWTSecret          string        `env:"JWT_SECRET,required"`
    AccessTokenTTL     time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
    RefreshTokenTTL    time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"720h"`

    // OTP
    OTPLength    int           `env:"OTP_LENGTH" envDefault:"6"`
    OTPTTL       time.Duration `env:"OTP_TTL" envDefault:"2m"`
    OTPMaxTries  int           `env:"OTP_MAX_TRIES" envDefault:"3"`
    OTPRateWindow time.Duration `env:"OTP_RATE_WINDOW" envDefault:"10m"`
    OTPRateMax    int           `env:"OTP_RATE_MAX" envDefault:"3"`

    // SMS
    SMSProvider   string `env:"SMS_PROVIDER" envDefault:"kavenegar"` // kavenegar | melipayamak
    SMSAPIKey     string `env:"SMS_API_KEY,required"`
    SMSSenderLine string `env:"SMS_SENDER_LINE"`

    // Push
    VAPIDPublicKey  string `env:"VAPID_PUBLIC_KEY,required"`
    VAPIDPrivateKey string `env:"VAPID_PRIVATE_KEY,required"`
    VAPIDSubject    string `env:"VAPID_SUBJECT" envDefault:"mailto:admin@nutritrack.ir"`

    // File Storage
    UploadDir    string `env:"UPLOAD_DIR" envDefault:"/data/uploads"`
    MaxFileSizeMB int   `env:"MAX_FILE_SIZE_MB" envDefault:"10"`

    // Timezone (also set TZ=Asia/Tehran in container env)
    Timezone string `env:"TIMEZONE" envDefault:"Asia/Tehran"`

    // Logging
    LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}
```

**Library choice:** `github.com/kelseyhightower/envconfig` over viper — no config files needed for a 12-factor app, env vars only, simpler struct tags.

---

## Cross-Cutting Concerns

### Persian Error Messages

```go
// pkg/apperror/codes.go
type ErrorCode string

const (
    ErrCodeNotFound         ErrorCode = "NOT_FOUND"
    ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
    ErrCodeForbidden        ErrorCode = "FORBIDDEN"
    ErrCodeValidation       ErrorCode = "VALIDATION_ERROR"
    ErrCodeOTPInvalid       ErrorCode = "OTP_INVALID"
    ErrCodeOTPExpired       ErrorCode = "OTP_EXPIRED"
    ErrCodeOTPRateLimit     ErrorCode = "OTP_RATE_LIMIT"
    ErrCodeDuplicateMobile  ErrorCode = "DUPLICATE_MOBILE"
    ErrCodeInactiveUser     ErrorCode = "INACTIVE_USER"
    ErrCodeActivePlanExists ErrorCode = "ACTIVE_PLAN_EXISTS"
    // …
)

// pkg/apperror/persian.go
var PersianMessages = map[ErrorCode]string{
    ErrCodeNotFound:         "مورد درخواستی یافت نشد",
    ErrCodeUnauthorized:     "لطفاً وارد حساب کاربری خود شوید",
    ErrCodeForbidden:        "دسترسی به این بخش مجاز نیست",
    ErrCodeOTPInvalid:       "کد تأیید وارد شده صحیح نیست",
    ErrCodeOTPExpired:       "کد تأیید منقضی شده است. لطفاً کد جدید دریافت کنید",
    ErrCodeOTPRateLimit:     "تعداد درخواست کد تأیید بیش از حد مجاز است. ۱۰ دقیقه دیگر تلاش کنید",
    ErrCodeDuplicateMobile:  "این شماره موبایل قبلاً ثبت شده است",
    ErrCodeActivePlanExists: "این مراجع از قبل دارای برنامه‌ی فعال است",
    // …
}
```

**Response envelope (always Persian `message`):**
```json
{
  "success": false,
  "message": "کد تأیید وارد شده صحیح نیست",
  "error_code": "OTP_INVALID",
  "data": null
}
```

---

### Asia/Tehran Timezone

```go
// pkg/logger/logger.go
var TehranLocation *time.Location

func init() {
    var err error
    TehranLocation, err = time.LoadLocation("Asia/Tehran")
    if err != nil {
        panic("cannot load Asia/Tehran timezone: " + err.Error())
    }
}

// All timestamp responses format in Tehran time:
func FormatTehran(t time.Time) string {
    return t.In(TehranLocation).Format(time.RFC3339)
}
```

**Container requirement:** `TZ=Asia/Tehran` in docker-compose `environment:` for all services.
**DB storage:** All timestamps stored as `TIMESTAMPTZ` (UTC internally, displayed in Tehran tz via app layer).
**"Today" computation:** Use `TodayTehran()` from `domain/shared/timeutil.go` — never `time.Now().UTC()` for date-only logic (3:30 AM UTC = 7 AM Tehran — wrong date otherwise).

---

### Structured Logging

```go
// pkg/logger/logger.go — zerolog
func New(level string) zerolog.Logger {
    lvl, _ := zerolog.ParseLevel(level)
    return zerolog.New(os.Stdout).
        Level(lvl).
        With().
        Timestamp().
        Str("service", "nutritrack-api").
        Logger()
}

// interface/http/middleware/logging.go
func RequestLogger(log zerolog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        log.Info().
            Str("method", c.Request.Method).
            Str("path", c.Request.URL.Path).
            Int("status", c.Writer.Status()).
            Str("request_id", c.GetString("X-Request-ID")).
            Str("user_id", c.GetString("user_id")).   // set by auth middleware
            Dur("latency_ms", time.Since(start)).
            Msg("request")
    }
}
```

---

## Dependency Injection Pattern

Manual DI wired in `internal/interface/bootstrap/wire.go` — no code generation tool needed at this scale:

```go
// internal/interface/bootstrap/wire.go
func BuildApp(cfg *config.Config) (*http.Server, error) {
    // Infrastructure
    db, err := postgres.Connect(cfg.DatabaseURL)
    queries := sqlc.New(db)

    redisClient := redis.NewClient(cfg.RedisURL)

    // Repositories (infrastructure → domain interface)
    userRepo     := repository.NewUserRepo(queries)
    foodRepo     := repository.NewFoodRepo(queries)
    dietPlanRepo := repository.NewDietPlanRepo(queries, db) // db for tx
    trackingRepo := repository.NewTrackingRepo(queries)
    messageRepo  := repository.NewMessageRepo(queries)
    // …

    // Ports
    otpStore     := redis.NewOTPStore(redisClient)
    tokenStore   := redis.NewTokenStore(redisClient)
    rateLimiter  := redis.NewRateLimiter(redisClient)
    smsAdapter   := sms.NewAdapter(cfg)
    pushAdapter  := push.NewAdapter(cfg)
    storageAdapter := storage.NewLocalAdapter(cfg.UploadDir)
    jwtManager   := jwt.NewManager(cfg.JWTSecret, cfg.AccessTokenTTL)

    // Application services
    authSvc     := appauth.NewAuthService(userRepo, otpStore, tokenStore, smsAdapter, jwtManager)
    userSvc     := appuser.NewUserService(userRepo)
    foodSvc     := appfood.NewFoodService(foodRepo)
    dietPlanSvc := appdietplan.NewDietPlanService(dietPlanRepo, userRepo, foodRepo, medicationRepo)
    trackingSvc := apptracking.NewTrackingService(trackingRepo, dietPlanRepo)
    messageSvc  := appmessage.NewMessageService(messageRepo, userRepo, pushAdapter)
    // …

    // HTTP layer
    logger := pkglogger.New(cfg.LogLevel)
    router := httpserver.NewRouter(cfg, logger, rateLimiter,
        authSvc, userSvc, foodSvc, dietPlanSvc, trackingSvc, messageSvc, …)

    return &http.Server{Addr: ":" + cfg.Port, Handler: router}, nil
}
```

---

## Docker Compose Multi-Service Setup

```yaml
# docker-compose.yml
version: "3.9"

services:
  app:
    build: .
    restart: unless-stopped
    environment:
      TZ: Asia/Tehran
      PORT: 8080
      DATABASE_URL: postgres://nutritrack:${DB_PASSWORD}@postgres:5432/nutritrack?sslmode=disable
      REDIS_URL: redis://redis:6379
      JWT_SECRET: ${JWT_SECRET}
      SMS_PROVIDER: kavenegar
      SMS_API_KEY: ${SMS_API_KEY}
      VAPID_PUBLIC_KEY: ${VAPID_PUBLIC_KEY}
      VAPID_PRIVATE_KEY: ${VAPID_PRIVATE_KEY}
      UPLOAD_DIR: /data/uploads
      LOG_LEVEL: info
      ENVIRONMENT: production
    volumes:
      - uploads:/data/uploads
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.api.rule=Host(`api.nutritrack.ir`)"
      - "traefik.http.routers.api.tls.certresolver=letsencrypt"
      - "traefik.http.services.api.loadbalancer.server.port=8080"

  postgres:
    image: postgres:16-alpine
    restart: unless-stopped
    environment:
      TZ: Asia/Tehran
      POSTGRES_DB: nutritrack
      POSTGRES_USER: nutritrack
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U nutritrack"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    environment:
      TZ: Asia/Tehran
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
    volumes:
      - redisdata:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  traefik:
    image: traefik:v3.0
    restart: unless-stopped
    command:
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.letsencrypt.acme.email=${ACME_EMAIL}"
      - "--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json"
      - "--certificatesresolvers.letsencrypt.acme.tlschallenge=true"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - letsencrypt:/letsencrypt

volumes:
  pgdata:
  redisdata:
  uploads:
  letsencrypt:
```

---

## Database Schema Organization

```
migrations/           ← golang-migrate versioned SQL (run at startup or via Makefile)
  000001_users        ← foundation: users table with roles, soft delete
  000002_foods        ← shared food database + food_categories junction
  000003_medications  ← shared medication database
  000004_diet_plans   ← diet_plans, plan_days, meals, meal_options, meal_option_items
                         prescribed_medications, exercise_recommendations
  000005_tracking     ← food_logs, water_logs, sleep_logs, exercise_logs,
                         medication_logs, body_measurements
                         (all with local_id for offline sync)
  000006_messages     ← messages + attachments
  000007_lab_results  ← lab_results with file_path
  000008_food_requests← food_requests with status workflow
  000009_push         ← push_subscriptions (endpoint, keys, user_id)

db/queries/           ← sqlc input SQL (one file per domain table group)
  users.sql           ← GetByID, GetByMobile, GetByEmail, ListClientsByNutritionist
  foods.sql           ← Insert, Update, GetByID, SearchByName (pg_trgm), ListByCategory
  diet_plans.sql      ← Insert, GetActiveByClient, GetByID, Archive, ListByClient
  tracking.sql        ← InsertFoodLog (ON CONFLICT local_id), GetDailyLogs, etc.
  messages.sql        ← Insert, GetConversation, CountUnread, MarkRead

internal/infrastructure/sqlc/  ← sqlc OUTPUT (generated, do not edit)
```

**pg_trgm for Persian search** — migration must include:
```sql
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX foods_name_trgm_idx ON foods USING gin (name gin_trgm_ops);
```

---

## Build Order (Phase Dependencies)

The build order follows domain dependency direction:

```
Phase 1 — Foundation
  1. Project scaffold (go mod, folder structure, Dockerfile, docker-compose.yml)
  2. Config loading (config.go + envconfig)
  3. Database connection + redis connection
  4. Logger setup (zerolog, Tehran tz)
  5. Migration runner (golang-migrate) + first migration (users table)
  6. pkg/apperror + Persian messages map
  7. pkg/persian (mobile validation, number normalization)

Phase 2 — Domain Core + Auth
  8. domain/shared (ID, pagination, timeutil)
  9. domain/user (entity, value objects, repository interface)
  10. infrastructure/repository/user_repo.go + sqlc queries/users.sql
  11. infrastructure/redis (OTP store, token store, rate limiter)
  12. infrastructure/jwt/manager.go
  13. infrastructure/sms/kavenegar.go
  14. application/auth (SendOTP, VerifyOTP, Refresh, Logout)
  15. interface/http (server, middleware: auth, logging, cors, recovery)
  16. interface/http/handler/auth.go
  → Milestone: OTP login works end-to-end

Phase 3 — Shared Databases (Food + Medication)
  17. domain/food + domain/medication
  18. Migrations: foods, medications tables + pg_trgm
  19. infrastructure/repository/food_repo + medication_repo
  20. application/food + application/medication services
  21. interface/http/handler/food.go + medication.go
  → Milestone: Admin/nutritionist can CRUD foods and medications

Phase 4 — Diet Plan Builder
  22. domain/dietplan (all entities + domain service)
  23. Migration: diet_plans + all nested tables
  24. infrastructure/repository/dietplan_repo (with transaction support)
  25. application/dietplan service
  26. interface/http/handler/dietplan.go
  → Milestone: Nutritionist can create and assign diet plans

Phase 5 — Client Tracking
  27. domain/tracking (all entity types)
  28. Migration: all tracking tables (with local_id unique constraints)
  29. infrastructure/repository/tracking_repo
  30. application/tracking service (idempotent writes)
  31. interface/http/handler/tracking.go
  → Milestone: Client can log all daily tracking data (offline sync ready)

Phase 6 — Messaging + Lab Results + Food Requests
  32. domain/message + domain/labresult + domain/foodrequest
  33. Migrations: messages, lab_results, food_requests
  34. infrastructure/storage/local.go (file upload handling)
  35. infrastructure/push/webpush.go
  36. Application services + handlers for messages, lab results, food requests
  → Milestone: Messaging, lab uploads, food request workflow complete

Phase 7 — Notifications + Admin Panel
  37. Push subscription management (DB + handler)
  38. Notification triggers in messaging, diet plan, food request services
  39. application/admin service + interface/http/handler/admin.go
  40. Super admin nutritionist management endpoints
  → Milestone: Push notifications + super admin panel complete

Phase 8 — Hardening
  41. Rate limiting middleware (Redis sliding window)
  42. Row-level authorization audit (every handler checked)
  43. Structured logging audit (all sensitive fields masked)
  44. Integration test suite
  45. Performance indexes review
```

---

## Anti-Patterns to Avoid

### Anti-Pattern 1: Anemic Domain Models
**What:** Entities are just data bags; all logic in app services.
**Why bad:** Invariants can be violated from anywhere; domain becomes meaningless.
**Instead:** Put `Archive()`, `AddPlanDay()`, `ComputeTotals()` on the DietPlan entity itself.

### Anti-Pattern 2: Domain Importing Infrastructure
**What:** `domain/dietplan/entity.go` imports `database/sql` or `github.com/go-redis/redis`.
**Why bad:** Destroys testability; domain can't be unit-tested without DB.
**Instead:** Domain defines interfaces (`DietPlanRepository`); infrastructure implements them.

### Anti-Pattern 3: Fat Handlers
**What:** Gin handler directly builds SQL queries or calls Redis.
**Why bad:** Business logic leaks into HTTP layer; untestable; hard to change transport.
**Instead:** Handler → App Service → Domain → Repository. Handler only: bind, validate, call service, shape response.

### Anti-Pattern 4: Single Giant Tracking Aggregate
**What:** One `TrackingRecord` aggregate with all 6 log types.
**Why bad:** Unnecessary coupling; offline sync for one type shouldn't touch others.
**Instead:** Separate small aggregates per tracking type — each has its own repository and table.

### Anti-Pattern 5: UTC-only Date Logic
**What:** `time.Now().UTC().Truncate(24*time.Hour)` for "today" comparisons.
**Why bad:** At midnight Tehran (20:30 UTC), this gives the wrong date for all Tehran-timezone logic.
**Instead:** Always use `TodayTehran()` which converts to `Asia/Tehran` before date extraction.

### Anti-Pattern 6: Loading Full DietPlan Aggregate on Every Request
**What:** Loading the entire plan (all days, meals, options, items) just to check if a plan is active.
**Why bad:** Huge DB join for a simple status check.
**Instead:** Separate query methods — `FindActiveByClientID()` returns a lightweight summary; `FindByIDFull()` loads the complete aggregate only when needed (plan builder view).

---

## Scalability Considerations

| Concern | At 50 nutritionists / 500 concurrent | At 10K clients |
|---------|--------------------------------------|----------------|
| DB connections | pgxpool (max 20 connections) sufficient | Add PgBouncer |
| Redis | Single node, 256MB sufficient | Sentinel for HA |
| File storage | Local Hetzner disk, path in DB | Move to object store (S3-compatible) |
| Full-text search | pg_trgm on foods table | Already sufficient at this scale |
| Diet plan load | Full aggregate load manageable | Add Redis caching for active plan |
| Push notifications | Synchronous webpush call | Move to async job queue (Asynq) |
| Polling chat | 10s interval × 500 users = 50 req/s | Redis pub/sub or long-poll upgrade |

Current scale targets (50 nutritionists, ~10,000 clients, ~500 concurrent) are comfortably served by the single-binary architecture with a single Postgres instance and single Redis node.

---

## Sources

- Go DDD community patterns: [https://github.com/marcusolsson/goddd](https://github.com/marcusolsson/goddd) — MEDIUM confidence (reference implementation)
- Gin framework docs: [https://gin-gonic.com/docs/](https://gin-gonic.com/docs/) — HIGH confidence
- sqlc docs: [https://docs.sqlc.dev/](https://docs.sqlc.dev/) — HIGH confidence
- golang-migrate: [https://github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) — HIGH confidence
- webpush-go: [https://github.com/SherClockHolmes/webpush-go](https://github.com/SherClockHolmes/webpush-go) — HIGH confidence
- zerolog: [https://github.com/rs/zerolog](https://github.com/rs/zerolog) — HIGH confidence
- envconfig: [https://github.com/kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) — HIGH confidence
- Persian pg_trgm: PostgreSQL docs — HIGH confidence (pg_trgm works on Unicode/Persian text)
- DDD aggregate design: "Implementing Domain-Driven Design" (Vernon) — HIGH confidence (canonical source)
