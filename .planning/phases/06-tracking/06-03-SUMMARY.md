---
phase: "06"
plan: "03"
subsystem: lab-results
tags: [lab-results, file-upload, magic-bytes, download, DDD, postgresql]
dependency_graph:
  requires: [06-01, 06-02]
  provides: [lab-result-upload, lab-result-download, lab-result-list]
  affects: [bootstrap/wire.go, router.go]
tech_stack:
  added: []
  patterns: [magic-byte-mime-detection, multipart-upload, file-attachment-download, paginated-list]
key_files:
  created:
    - migrations/000008_lab_results.up.sql
    - migrations/000008_lab_results.down.sql
    - db/queries/lab_results.sql
    - internal/infrastructure/persistence/sqlc/lab_results.sql.go
    - internal/domain/labresult/entity/lab_result.go
    - internal/domain/labresult/repository/lab_result_repository.go
    - internal/infrastructure/persistence/labresult/mapper.go
    - internal/infrastructure/persistence/labresult/pg_lab_result_repository.go
    - internal/application/labresult/lab_result_service.go
    - internal/interfaces/http/handler/lab_result_handler.go
  modified:
    - internal/infrastructure/persistence/sqlc/models.go
    - internal/domain/shared/apperror.go
    - internal/infrastructure/storage/local_storage.go
    - bootstrap/wire.go
    - internal/interfaces/http/router/router.go
decisions:
  - Magic-byte MIME detection (not Content-Type header) prevents spoofing for lab result uploads — same pattern as avatar uploads
  - SaveLabResult returns filesystem path (not URL) because download endpoint serves via c.FileAttachment directly
  - ErrLabResultNotFound added to shared error catalog for consistent Persian error responses
metrics:
  duration: "~20m"
  completed: "2026-04-21"
  tasks_completed: 14
  files_changed: 15
---

# Phase 6 Plan 03: Lab Results Upload & Download Summary

**One-liner:** Lab result file upload (PDF/JPEG/PNG) with magic-byte MIME validation, paginated listing, and authenticated file download via `c.FileAttachment`.

## What Was Built

Complete lab results feature across all DDD layers:

1. **Migration 000008** — `lab_results` table with FK to `users` (client + nutritionist), `file_path`, `original_name`, `file_type`, `file_size`, `notes`, `created_at`. Index on `client_id`.
2. **sqlc layer** — hand-written `lab_results.sql.go` with `CreateLabResult`, `GetLabResultByID`, `ListLabResultsByClientID` (paginated), `CountLabResultsByClientID`.
3. **Domain** — `entity.LabResult` struct + `LabResultRepository` interface.
4. **AppError** — `ErrLabResultNotFound` added to Persian error catalog.
5. **Infrastructure** — `PgLabResultRepository` + `mapper.go` (toDomain helper + `isNotFound`).
6. **Storage** — `LocalStorage.SaveLabResult` saves to `<basePath>/lab-results/<uuid>.<ext>`, returns filesystem path.
7. **Application Service** — `LabResultService` with:
   - Magic-byte MIME detection (PDF: `%PDF`, JPEG: `FF D8 FF`, PNG: `89 50 4E 47...`)
   - 10 MB size limit
   - Role-based access control (superadmin/nutritionist/client)
   - `UploadLabResult`, `ListClientLabResults`, `GetLabResultForDownload`
8. **HTTP Handler** — `LabResultHandler` with `Upload`, `List`, `Download` endpoints.
9. **Wire** — `LabResultService` wired in `bootstrap/wire.go`.
10. **Routes** — `POST /clients/:id/lab-results`, `GET /clients/:id/lab-results`, `GET /lab-results/:id/download`.

## API Routes

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/clients/:id/lab-results` | Upload lab result file (multipart/form-data, field: `file`) |
| GET | `/api/v1/clients/:id/lab-results` | List paginated lab results for a client |
| GET | `/api/v1/lab-results/:id/download` | Download a lab result file as attachment |

## Deviations from Plan

None — plan executed exactly as written. The incomplete `checkClientAccess` stub mentioned in the plan prompt was not included; only the final clean methods were implemented.

## Known Stubs

None — all endpoints are fully wired with real repository and storage.

## Self-Check: PASSED

- `migrations/000008_lab_results.up.sql` ✅
- `migrations/000008_lab_results.down.sql` ✅
- `db/queries/lab_results.sql` ✅
- `internal/infrastructure/persistence/sqlc/lab_results.sql.go` ✅
- `internal/domain/labresult/entity/lab_result.go` ✅
- `internal/domain/labresult/repository/lab_result_repository.go` ✅
- `internal/infrastructure/persistence/labresult/mapper.go` ✅
- `internal/infrastructure/persistence/labresult/pg_lab_result_repository.go` ✅
- `internal/application/labresult/lab_result_service.go` ✅
- `internal/interfaces/http/handler/lab_result_handler.go` ✅
- Commit `4e91f98` ✅
- `go build ./...` exit 0 ✅
