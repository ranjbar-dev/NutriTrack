---
phase: "03"
plan: "03"
subsystem: "user-profile"
tags: [avatar, file-upload, magic-bytes, mime-validation, local-storage, ddd]
dependency_graph:
  requires: [03-01, 03-02]
  provides: [avatar-upload-endpoint]
  affects: [user-entity, bootstrap-container, router]
tech_stack:
  added: []
  patterns: [magic-byte-mime-validation, multipart-upload, local-filesystem-storage, role-based-access-control]
key_files:
  created:
    - internal/domain/shared/file_validator.go
    - internal/infrastructure/storage/local_storage.go
    - internal/application/user/avatar_service.go
    - internal/interfaces/http/handler/avatar_handler.go
  modified:
    - bootstrap/wire.go
    - internal/interfaces/http/router/router.go
decisions:
  - "Magic byte validation (not Content-Type header) for MIME detection — prevents spoofed uploads"
  - "bytes.NewReader used instead of custom bytesReader to reconstruct full file stream after reading magic bytes"
  - "dto.Abort used (not dto.Error) for consistency with existing handler pattern"
  - "Access control enforced in service layer (AvatarService.UploadAvatar) not handler — DDD boundary"
metrics:
  duration: "~10 min"
  completed: "2026-04-21"
  tasks_completed: 1
  files_changed: 6
---

# Phase 3 Plan 03: Profile Picture Upload Summary

**One-liner:** Avatar upload endpoint with magic-byte MIME validation (JPEG/PNG/WEBP), local filesystem storage, and role-based access control in service layer.

## What Was Built

Implemented `PUT /users/:id/avatar` — a multipart file upload endpoint for profile pictures with the following layers:

### Domain Layer
- **`internal/domain/shared/file_validator.go`** — `ValidateImageMagicBytes()` reads the first 12 bytes to detect JPEG (`FF D8 FF`), PNG (`89 50 4E 47...`), or WEBP (`52 49 46 46...WEBP`) — never trusts `Content-Type` header.

### Infrastructure Layer
- **`internal/infrastructure/storage/local_storage.go`** — `LocalStorage.SaveAvatar()` persists files to `uploads/avatars/<uuid>.<ext>`, creates directories if needed, returns URL path.

### Application Layer
- **`internal/application/user/avatar_service.go`** — `AvatarService.UploadAvatar()` orchestrates: validate MIME → fetch user → RBAC check → save file → update DB.
  - Superadmin: update anyone
  - Nutritionist: update own avatar or own client's avatar
  - Client: update own avatar only

### Interface Layer
- **`internal/interfaces/http/handler/avatar_handler.go`** — 5 MB size limit via `http.MaxBytesReader`, reads 12 magic bytes, reconstructs full reader with `io.MultiReader + bytes.NewReader`.
- **`internal/interfaces/http/router/router.go`** — `PUT /users/:id/avatar` under protected group; `r.Static("/uploads", "./uploads")` for serving stored avatars.

### DI Container
- **`bootstrap/wire.go`** — Wired `LocalStorage` and `AvatarService` into `Container`.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Replaced custom `bytesReader` with `bytes.NewReader`**
- **Found during:** Implementation of `avatar_handler.go`
- **Issue:** The plan provided a custom `bytesReader` struct implementing `io.Reader` over a byte slice — this is identical to `bytes.NewReader` from the standard library.
- **Fix:** Used `bytes.NewReader(magicBuf)` instead of the custom struct — cleaner, no maintenance burden.
- **Files modified:** `internal/interfaces/http/handler/avatar_handler.go`
- **Commit:** cbe3bb0

**2. [Rule 1 - Bug] Used `dto.Abort` instead of `dto.Error`**
- **Found during:** Implementation of `avatar_handler.go`
- **Issue:** `dto.Error(c, err)` where `err` is of type `error` would not compile — `dto.Error` takes `*shared.AppError`. Other handlers in the codebase use `dto.Abort` with explicit type assertion.
- **Fix:** Used `dto.Abort(c, appErr)` with `appErr, ok := err.(*shared.AppError)` type assertion — consistent with `nutritionist_handler.go` and `auth_handler.go` pattern.
- **Files modified:** `internal/interfaces/http/handler/avatar_handler.go`
- **Commit:** cbe3bb0

## Known Stubs

None — all data flows are wired. Avatar URL is persisted to DB via `userRepo.Update()` and returned in response.

## Threat Flags

| Flag | File | Description |
|------|------|-------------|
| threat_flag: path-traversal | internal/infrastructure/storage/local_storage.go | File extension from `ValidateImageMagicBytes` is trusted (only "jpg"/"png"/"webp") — UUID filename prevents traversal. Low risk. |
| threat_flag: disk-exhaustion | internal/infrastructure/storage/local_storage.go | No cleanup of old avatar on update — each upload creates a new file. Old files accumulate. Future plan should add orphan cleanup. |

## Self-Check: PASSED

- `internal/domain/shared/file_validator.go` ✅
- `internal/infrastructure/storage/local_storage.go` ✅
- `internal/application/user/avatar_service.go` ✅
- `internal/interfaces/http/handler/avatar_handler.go` ✅
- `bootstrap/wire.go` ✅
- `internal/interfaces/http/router/router.go` ✅
- Commit cbe3bb0 ✅
- `go build ./...` passes ✅
