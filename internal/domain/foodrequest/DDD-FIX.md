# DDD Fix Report: internal/domain/foodrequest
Layer: domain
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | FoodRequest aggregate exposes all fields as exported | CRITICAL | entity/food_request.go + 4 consumers | SAFE | FIXED |
| 2 | Missing `NewFoodRequest()` factory function | HIGH | entity/food_request.go | SAFE | FIXED |
| 3 | No domain-behaviour mutation methods (`Approve`, `Reject`) | MEDIUM | entity/food_request.go | DEFERRED | DEFERRED: out of scope for this run |
| 4 | No domain error variables defined | LOW | entity/food_request.go | DEFERRED | DEFERRED: partial (ErrFoodRequestInvalidInput added; others low priority) |

---

## Changes Applied

### Fix 1 + 2: Unexported fields, getters, Hydrate helper, NewFoodRequest factory

**File:** `internal/domain/foodrequest/entity/food_request.go`

**Change:**
- All 9 struct fields renamed to unexported (`id`, `clientID`, `nutritionistID`, `foodName`, `status`, `rejectionReason`, `createdFoodID`, `createdAt`, `updatedAt`)
- Added getter methods for all 9 fields (`GetID`, `GetClientID`, `GetNutritionistID`, `GetFoodName`, `GetStatus`, `GetRejectionReason`, `GetCreatedFoodID`, `GetCreatedAt`, `GetUpdatedAt`)
- Added `NewFoodRequest(clientID, nutritionistID uuid.UUID, foodName string) (*FoodRequest, error)` factory with required-field validation (returns `ErrFoodRequestInvalidInput` if any arg is zero/empty)
- Added `FromPersistence(...)` package-level function for the infrastructure mapper to reconstruct aggregates from DB rows without bypassing the domain model
- Added `Hydrate(id, status, createdAt, updatedAt)` method for the repository's `Create` method to write back DB-generated fields after insert
- Added `ErrFoodRequestInvalidInput` sentinel error variable
- State predicates (`IsPending`, `IsApproved`, `IsRejected`) updated to read the unexported `status` field

**Before:**
```go
type FoodRequest struct {
    ID              uuid.UUID
    ClientID        uuid.UUID
    NutritionistID  uuid.UUID
    FoodName        string
    Status          FoodRequestStatus
    RejectionReason *string
    CreatedFoodID   *uuid.UUID
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

**After:**
```go
type FoodRequest struct {
    id              uuid.UUID
    clientID        uuid.UUID
    nutritionistID  uuid.UUID
    foodName        string
    status          FoodRequestStatus
    rejectionReason *string
    createdFoodID   *uuid.UUID
    createdAt       time.Time
    updatedAt       time.Time
}
func (r *FoodRequest) GetID() uuid.UUID             { return r.id }
// … (8 more getters)

func NewFoodRequest(clientID, nutritionistID uuid.UUID, foodName string) (*FoodRequest, error) { … }
func FromPersistence(…) *FoodRequest { … }
func (r *FoodRequest) Hydrate(id uuid.UUID, status FoodRequestStatus, createdAt, updatedAt time.Time) { … }
```

**Build:** PASS

---

### Fix 1a: Infrastructure mapper — use FromPersistence

**File:** `internal/infrastructure/persistence/foodrequest/mapper.go`

**Change:** Replaced direct struct literal (`&entity.FoodRequest{ID: row.ID, …}`) with `entity.FromPersistence(…)` call so the infrastructure never constructs the aggregate with exported fields.

**Before:**
```go
func toDomain(row db.FoodRequest) *entity.FoodRequest {
    return &entity.FoodRequest{
        ID:             row.ID,
        ClientID:       row.ClientID,
        …
    }
}
```

**After:**
```go
func toDomain(row db.FoodRequest) *entity.FoodRequest {
    return entity.FromPersistence(
        row.ID, row.ClientID, row.NutritionistID, row.FoodName,
        entity.FoodRequestStatus(row.Status), row.RejectionReason,
        row.CreatedFoodID, row.CreatedAt, row.UpdatedAt,
    )
}
```

**Build:** PASS

---

### Fix 1b: Infrastructure repository — use getters + Hydrate

**File:** `internal/infrastructure/persistence/foodrequest/pg_food_request_repository.go`

**Change:** `Create` method replaced direct field reads (`req.ClientID`, `req.NutritionistID`, `req.FoodName`) with getters, and replaced four individual write-back assignments with a single `req.Hydrate(…)` call.

**Before:**
```go
ClientID:       req.ClientID,
NutritionistID: req.NutritionistID,
FoodName:       req.FoodName,
// …
req.ID = row.ID
req.Status = entity.FoodRequestStatus(row.Status)
req.CreatedAt = row.CreatedAt
req.UpdatedAt = row.UpdatedAt
```

**After:**
```go
ClientID:       req.GetClientID(),
NutritionistID: req.GetNutritionistID(),
FoodName:       req.GetFoodName(),
// …
req.Hydrate(row.ID, entity.FoodRequestStatus(row.Status), row.CreatedAt, row.UpdatedAt)
```

**Build:** PASS

---

### Fix 1c + 2a: Application service — use NewFoodRequest and getters

**File:** `internal/application/foodrequest/food_request_service.go`

**Change:**
- `Submit`: replaced raw struct literal with `frEntity.NewFoodRequest(clientID, *client.NutritionistID, foodName)` factory call; validation errors surface naturally
- `Approve` and `Reject`: replaced `foodReq.NutritionistID` with `foodReq.GetNutritionistID()`

**Before:**
```go
req := &frEntity.FoodRequest{
    ClientID:       clientID,
    NutritionistID: *client.NutritionistID,
    FoodName:       foodName,
    Status:         frEntity.FoodRequestStatusPending,
}
// …
if foodReq.NutritionistID != nutritionistID {
```

**After:**
```go
req, err := frEntity.NewFoodRequest(clientID, *client.NutritionistID, foodName)
if err != nil {
    return nil, err
}
// …
if foodReq.GetNutritionistID() != nutritionistID {
```

**Build:** PASS

---

### Fix 1d: HTTP handler — use getters in foodRequestToMap and push goroutines

**File:** `internal/interfaces/http/handler/food_request_handler.go`

**Change:**
- `foodRequestToMap`: replaced all 9 direct field reads with getter calls
- `Approve` and `Reject` handlers: replaced `result.ClientID` with `result.GetClientID()` in the push notification goroutine

**Before:**
```go
"id":        r.ID,
"client_id": r.ClientID,
// …
clientID := result.ClientID
```

**After:**
```go
"id":        r.GetID(),
"client_id": r.GetClientID(),
// …
clientID := result.GetClientID()
```

**Build:** PASS

---

## Deferred Items

- **[MEDIUM] Add `Approve()` and `Reject()` transition methods** — would require adding state-machine logic to the entity; safe to add independently without touching consumers since the application service currently calls `UpdateStatus` on the repo directly. Out of scope for this fix run.
- **[LOW] Additional sentinel errors** (`ErrFoodRequestAlreadyReviewed`, `ErrFoodRequestNotFound`) — `ErrFoodRequestNotFound` already lives in `internal/domain/shared`; adding a duplicate in the domain sub-package would require updating all call sites. Deferred.

---

## Final Build Status
PASS — `go build ./...` after all fixes  
PASS — `go build ./internal/domain/foodrequest/...`  
PASS — `go build ./internal/infrastructure/persistence/foodrequest/...`  
PASS — `go build ./internal/application/foodrequest/...`  
PASS — `go build ./internal/interfaces/...`

## Remaining Violations
None at CRITICAL or HIGH severity. All MEDIUM and LOW items are documented as deferred above.
