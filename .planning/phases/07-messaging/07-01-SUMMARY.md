---
phase: "07"
plan: "01"
subsystem: messaging
tags: [messaging, chat, attachments, unread-count, ddd, go]
dependency_graph:
  requires: [06-03]
  provides: [message-entity, message-service, message-endpoints]
  affects: [bootstrap/wire.go, router.go]
tech_stack:
  added: []
  patterns: [magic-byte-validation, multipart-upload, io.Reader-nil-check, domain-repository-interface]
key_files:
  created:
    - migrations/000009_messages.up.sql
    - migrations/000009_messages.down.sql
    - db/queries/messages.sql
    - internal/infrastructure/persistence/sqlc/messages.sql.go
    - internal/domain/message/entity/message.go
    - internal/domain/message/repository/message_repository.go
    - internal/infrastructure/persistence/message/mapper.go
    - internal/infrastructure/persistence/message/pg_message_repository.go
    - internal/application/message/message_service.go
    - internal/interfaces/http/handler/message_handler.go
  modified:
    - internal/infrastructure/persistence/sqlc/models.go
    - internal/domain/shared/apperror.go
    - internal/infrastructure/storage/local_storage.go
    - bootstrap/wire.go
    - internal/interfaces/http/router/router.go
decisions:
  - "io.Reader nil-check pattern: pass io.Reader as nil when no attachment (not a zero-value interface)"
  - "Magic byte detection: reused hasMagic/detectMIME pattern from lab_result_service.go"
  - "SaveAttachment returns URL path (same as SaveAvatar), not FS path (unlike SaveLabResult)"
  - "Conversation ordering: ASC by created_at — clients see oldest first"
  - "MarkRead called after listing — marks opposite party's messages as read on fetch"
metrics:
  duration: "15m"
  completed: "2026-04-21"
  tasks_completed: 14
  files_created: 15
---

# Phase 7 Plan 01: Message Domain, Conversation Endpoints, Attachment Upload, Unread Count — Summary

**One-liner:** Chat messaging between clients and nutritionists with magic-byte-validated file attachments, read tracking, and paginated conversation history.

## What Was Built

Complete DDD messaging system implementing the Plan 07-01 specification:

1. **Migration 000009** — `messages` table with full conversation structure: sender_id, receiver_id, content, nullable attachment fields (path, type, size, name), read_at, created_at. Two indexes: receiver_id for unread queries, LEAST/GREATEST composite for bidirectional conversation lookups.

2. **sqlc query file** (`db/queries/messages.sql`) — 5 named queries: CreateMessage, ListConversationMessages, CountConversationMessages, MarkConversationRead, CountUnreadMessages.

3. **Hand-written sqlc Go implementation** (`messages.sql.go`) — All 5 queries implemented with proper Scan patterns; `Message` struct appended to `models.go`.

4. **Domain layer** — `entity.Message` with `HasAttachment()` method, `MessageRepository` interface (Create, ListConversation, MarkRead, CountUnread).

5. **AppError** — Added `ErrMessageNotFound` to Persian error catalog.

6. **Infrastructure** — `PgMessageRepository` implementing the domain interface, mapper from db.Message → entity.Message, `isNotFound` helper.

7. **LocalStorage extension** — `SaveAttachment` method (returns URL path to `/uploads/attachments/<uuid>.<ext>`).

8. **Application service** (`MessageService`) — SendAsClient, SendAsNutritionist (both validate ownership), GetClientConversation, GetNutritionistConversation (both auto-mark-read on fetch), GetUnreadCount. Magic-byte MIME detection for PDF/JPEG/PNG attachments with per-type size limits (5 MB image, 10 MB PDF).

9. **HTTP handler** — 5 endpoints, multipart/form-data with optional `file` field, `io.Reader` nil pattern for no-attachment case.

10. **Wire + Router** — MessageService wired into Container; 5 routes registered under protected group.

## Routes Registered

| Method | Path | Handler | Role |
|--------|------|---------|------|
| GET | /api/v1/messages/unread-count | GetUnreadCount | any |
| GET | /api/v1/messages | GetClientMessages | client |
| POST | /api/v1/messages | SendAsClient | client |
| GET | /api/v1/clients/:id/messages | GetNutritionistMessages | nutritionist |
| POST | /api/v1/clients/:id/messages | SendAsNutritionist | nutritionist |

## Deviations from Plan

None — plan executed exactly as written.

## Known Stubs

None. All endpoints are wired to real service → repository → database path.

## Threat Flags

| Flag | File | Description |
|------|------|-------------|
| threat_flag: unauthenticated-attachment-access | router.go | `/uploads/attachments/*` served publicly via `r.Static("/uploads", "./uploads")`. Any user with the URL can download attachments without auth. This is existing behavior (same as avatars). Future mitigation in Phase 8 security audit. |

## Self-Check: PASSED

- `migrations/000009_messages.up.sql` ✅
- `migrations/000009_messages.down.sql` ✅
- `db/queries/messages.sql` ✅
- `internal/infrastructure/persistence/sqlc/messages.sql.go` ✅
- `internal/domain/message/entity/message.go` ✅
- `internal/domain/message/repository/message_repository.go` ✅
- `internal/infrastructure/persistence/message/mapper.go` ✅
- `internal/infrastructure/persistence/message/pg_message_repository.go` ✅
- `internal/application/message/message_service.go` ✅
- `internal/interfaces/http/handler/message_handler.go` ✅
- Commit `54d834a` ✅
- `go build ./...` exit 0 ✅
