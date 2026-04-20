# 05-05 Summary: Communication Handler & Route Wiring

## What Was Built

### Communication Handler (`backend/internal/handler/communication_handler.go`)

Messaging endpoints:
- `SendMessage` — POST /api/messages/:partnerId — multipart/form-data, reads `content` + optional `attachment` file + `receiver_id` form field
- `ListMessages` — GET /api/messages/:partnerId — full conversation history
- `PollNewMessages` — GET /api/messages/:partnerId/new?since=<ISO timestamp> — polling endpoint
- `MarkRead` — PUT /api/messages/:partnerId/read — marks conversation messages as read
- `GetUnreadCount` — GET /api/messages/unread-count — badge count
- `DownloadAttachment` — GET /api/messages/attachment/:messageId — streams file with Content-Disposition: attachment

Food request endpoints:
- `ClientCreateFoodRequest` — POST /api/clients/me/food-requests
- `ClientListFoodRequests` — GET /api/clients/me/food-requests
- `NutriListFoodRequests` — GET /api/nutritionist/food-requests?status=pending
- `NutriApproveFoodRequest` — PUT /api/nutritionist/food-requests/:id/approve
- `NutriRejectFoodRequest` — PUT /api/nutritionist/food-requests/:id/reject

### Route Wiring (`backend/cmd/api/main.go`)

Added initialization:
```go
commRepo := repository.NewCommunicationRepository(pool, queries)
commService := service.NewCommunicationService(commRepo, userRepo, cfg)
commHandler := handler.NewCommunicationHandler(commService)
```

Client management routes added to `nutri` group:
- GET/PUT endpoints for list, profile, activate, deactivate, update

Food request routes added to `client` group and `nutri` group.

Messaging group `/api/messages` with **correct ordering**:
- `/unread-count` registered BEFORE `/:partnerId` wildcard
- `/attachment/:messageId` registered BEFORE `/:partnerId` wildcard
- `/:partnerId`, `/:partnerId/new`, `/:partnerId/read` registered after

## Key Technical Details

- Gin route ordering is critical — specific paths must be registered before wildcards in same group
- `c.Request.FormFile("attachment")` gracefully skips if not present (text-only messages valid)
- `receiver_id` passed as form field (not JSON body) because SendMessage uses multipart

## Outcomes

- All routes wired and compiling
- `go build ./...` exit code 0
