# 05-03 Summary: DTOs, Repository & Service

## What Was Built

### DTOs (`backend/internal/model/dto/communication_dto.go`)

Request types:
- `SendMessageRequest` — receiver_id (UUID), content (text), attachment_path/name/mime (optional)
- `CreateFoodRequestRequest` — food_name, description (optional)
- `ApproveFoodRequestRequest` — (empty, food created from stored food_name)
- `RejectFoodRequestRequest` — reject_reason (optional)
- `UpdateClientProfileRequest` — height_cm, date_of_birth, gender (all optional)

Response types:
- `MessageResponse` — id, sender_id, receiver_id, content, attachment fields, created_at, read_at
- `FoodRequestResponse` — id, client_id, nutritionist_id, food_name, description, status, reject_reason, reviewed_at, created_at
- `ClientProfileResponse` — full client info: id, name, mobile, height_cm, date_of_birth, gender, is_active, nutritionist_id, created_at
- `ClientListItemResponse` — condensed for list views

### Repository (`backend/internal/repository/communication_repo.go`)

`CommunicationRepository` interface + `pgxCommunicationRepository` implementation wrapping sqlc `*Queries`:
- All 12 methods delegating to sqlc: SendMessage, GetConversation, GetNewMessages, MarkRead, GetUnreadCount, GetMessageByID, CreateFoodRequest, GetFoodRequestByID, ListClientFoodRequests, ListNutriPendingFoodRequests, ApproveFoodRequest, RejectFoodRequest

### User Repository Update (`backend/internal/repository/user_repo.go`)

- Changed from `sqlc.DBTX` to `*pgxpool.Pool` (needed for dynamic SQL in SearchClients)
- `SearchClientsParams` struct added
- `SearchClients` — raw SQL with dynamic ORDER BY (`name` or `last_activity`) + ILIKE filter (safe: only 2 allowed ORDER BY values)
- `GetClientByIDForNutritionist` — verifies ownership before returning
- `UpdateClientProfile` — height_cm, date_of_birth, gender update

### Service (`backend/internal/service/communication_service.go`)

Full `CommunicationService` with:
- `verifyConversationPair` — validates nutritionist-client relationship before any message op
- `SendMessageTo` — multipart file upload with mimetype detection, size limits (5MB images, 10MB PDFs), UUID filename storage at `uploads/messages/{senderID}/{uuid}.{ext}`
- `GetMessages`, `GetNewMessages`, `MarkRead`, `GetUnreadCount`
- `GetMessageAttachment` — returns file path + mime for streaming
- `CreateFoodRequest`, `ListClientFoodRequests`, `ListNutriPendingFoodRequests`, `ApproveFoodRequest`, `RejectFoodRequest`

## Decisions Made

- File upload uses `mimetype.DetectReader` which reads from the reader — must seek back to 0 after detection
- Food request approval does NOT auto-create food; returns `food_name` so frontend redirects to `/nutritionist/foods/create?name=<food_name>`
- `verifyConversationPair` checks `nutritionist_id` field on client user record

## Outcomes

- Backend service layer complete and compiling
