# DDD Fix Report: internal/application/tracking
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Struct literal construction bypasses factory | HIGH | tracking_service.go | N/A | PRE-FIXED |
| 2 | No input validation in Log* methods | HIGH | tracking_service.go | SAFE | FIXED |
| 3 | json: tags on app-layer types (SyncEntry, BulkSyncResult) | MEDIUM | tracking_service.go | DEFERRED | DEFERRED |
| 4 | encoding/json in application service | MEDIUM | tracking_service.go | DEFERRED | DEFERRED |
| 5 | GetTracking returns any | MEDIUM | tracking_service.go | DEFERRED | DEFERRED |
| 6 | BulkSync duplicates Log* logic | MEDIUM | tracking_service.go | DEFERRED | DEFERRED |
| 7 | Silent error swallowing in BulkSync | LOW | tracking_service.go | DEFERRED | DEFERRED |

## Changes Applied

### Fix 2: Input validation guard clauses in all six Log* methods

**File:** `internal/application/tracking/tracking_service.go`

Added guard clauses at the start of each `Log*` method returning `shared.ErrValidation`:

| Method | Guard |
|--------|-------|
| `LogFood` | `FoodName == ""` or `Quantity <= 0` |
| `LogWater` | `AmountMl <= 0` |
| `LogSleep` | `DurationMinutes <= 0` |
| `LogExercise` | `ExerciseName == ""` or `DurationMinutes <= 0` |
| `LogMedication` | `MedicationName == ""` |
| `LogBody` | all measurement fields nil |

**Build:** PASS

## Pre-Fixed Findings (verified)

- **Factory functions**: All six entity factories (`NewFoodLog`, `NewWaterLog`, `NewSleepLog`, `NewExerciseLog`, `NewMedicationLog`, `NewBodyMeasurement`) already existed with unexported fields and were already called. ✓
- `BulkSync` also already used factory calls (not struct literals). ✓

## Deferred Items

- **[MEDIUM]** `SyncEntry`/`BulkSyncResult` with json tags — moving to interfaces layer requires HTTP handler changes across multiple files; deferred.
- **[MEDIUM]** `encoding/json` in `BulkSync` — coupled to `SyncEntry` move; deferred together.
- **[MEDIUM]** `GetTracking(any)` typed return — requires updating all callers; deferred.
- **[MEDIUM]** `BulkSync` delegation refactor — deferred.
- **[LOW]** Per-entry failure tracking — deferred.

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None at HIGH severity.
