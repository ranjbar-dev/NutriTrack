# DDD Fix Report: internal/application/medication
Layer: application
Fixed: 2026-04-22
Based on: DDD-AUDIT.md

## Baseline Build Status
PASS — `go build ./...` before fixes

## Fix Plan

| # | Finding | Severity | Files | Strategy | Status |
|---|---------|----------|-------|----------|--------|
| 1 | Struct literal construction bypasses factory | HIGH | medication_service.go | N/A | PRE-FIXED |
| 2 | Direct mutation of exported aggregate fields | HIGH | medication_service.go | N/A | PRE-FIXED |
| 3 | Deactivate bypasses aggregate lifecycle | MEDIUM | medication_service.go | DEFERRED | DEFERRED |
| 4 | Magic role string literals | MEDIUM | medication_service.go | DEFERRED | DEFERRED |

## Changes Applied

No changes required — all HIGH violations were already resolved prior to this fix pass.

## Pre-Fixed Findings (verified)

- **Factory call**: `CreateMedication` already calls `entity.NewMedication(req.Name, normalized, req.Description, req.Unit, &req.CallerID)` — no struct literal. ✓
- **Setter methods**: `UpdateMedication` already calls `med.SetName(...)`, `med.SetNameNormalized(...)`, `med.SetDescription(...)`, `med.SetUnit(...)` — no direct field assignment. ✓

## Deferred Items

- **[MEDIUM]** `Deactivate()` aggregate method — load-transition-persist pattern; safe to implement but out of scope for HIGH-only pass.
- **[MEDIUM]** Magic role strings — same as dietplan, deferred.

## Final Build Status
PASS — `go build ./...`
PASS — `go vet ./internal/...`

## Remaining Violations
None at HIGH severity.
