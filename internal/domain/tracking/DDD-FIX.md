# DDD Fix Report: internal/domain/tracking
Layer: domain
Fixed: 2025-07
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | 6 aggregate structs have exported fields | HIGH | entity/tracking.go | SAFE | FIXED |
| 2 | No `New*()` factory functions | HIGH | entity/tracking.go | SAFE | FIXED |
| 3 | Callers use direct field access | HIGH | persistence/tracking/mapper.go, pg_tracking_repository.go, application/tracking/tracking_service.go, interfaces/http/handler/tracking_handler.go | SAFE | FIXED |

## Changes Applied

### Fix 1 + 2: Unexported fields, getters, factory functions for all 6 entities
**File:** `internal/domain/tracking/entity/tracking.go`
**Change:** FoodLog, WaterLog, SleepLog, ExerciseLog, MedicationLog, BodyMeasurement — all fields made unexported. Added getter for every field. Added `New*()` and `Reconstitute*()` for each entity.
**Build:** PASS

### Fix 3a: Mapper updated
**File:** `internal/infrastructure/persistence/tracking/mapper.go`
**Change:** All 6 `toDomain*` functions use `entity.Reconstitute*()` instead of struct literals.
**Build:** PASS

### Fix 3b: Repository updated
**File:** `internal/infrastructure/persistence/tracking/pg_tracking_repository.go`
**Change:** All 6 `Upsert*` methods use getter methods for DB params (e.g., `log.ClientID()`, `log.LocalID()`). Struct copy pattern (`*log = *toDomain(row)`) preserved — valid in Go for same-type unexported fields.
**Build:** PASS

### Fix 3c: Service updated
**File:** `internal/application/tracking/tracking_service.go`
**Change:** All LogFood/LogWater/LogSleep/LogExercise/LogMedication/LogBody handlers and BulkSync switch cases use `entity.New*()` factories instead of struct literals.
**Build:** PASS

### Fix 3d: Handler updated
**File:** `internal/interfaces/http/handler/tracking_handler.go`
**Change:** All 6 `*ToMap` helper functions use getter methods (e.g., `log.ID()`, `log.ClientID()`, `log.LoggedAt()`).
**Build:** PASS

## Final Build Status
PASS — `go build ./...` after all fixes
PASS — `go vet ./internal/...` after all fixes

## Remaining Violations
None — all CRITICAL and HIGH findings resolved.
