# Phase 5 Verification: Communication & Collaboration

## Goal

> Clients and nutritionists can communicate via messaging, clients can request food additions, and nutritionists have a complete client management workspace

## Requirement Coverage

### Messaging (MSG-01 → MSG-07)

| Req | Description | Implementation | Status |
|-----|-------------|----------------|--------|
| MSG-01 | Chat-style messaging between client and assigned nutritionist | `messages` table, `SendMessageTo`/`GetMessages` service, `SendMessage`/`ListMessages` handler, `/api/messages/:partnerId` routes | ✅ COVERED |
| MSG-02 | Client can only message their assigned nutritionist; nutritionist messages own clients | `verifyConversationPair()` in `communication_service.go` checks `nutritionist_id` FK before all message operations | ✅ COVERED |
| MSG-03 | Text messages with optional attachments (images JPG/PNG ≤5MB, PDFs ≤10MB) | `SendMessage` handler reads multipart form, `communication_service.go` validates magic bytes via `mimetype.DetectReader`, enforces size limits by MIME type, stores at `uploads/messages/{senderID}/{uuid}.{ext}` | ✅ COVERED |
| MSG-04 | Polling-based delivery (every 10 seconds when chat is open) | `useMessagePolling.ts` composable sets `setInterval` at 10s calling `PollNewMessages` endpoint; cleared on `onUnmounted` | ✅ COVERED |
| MSG-05 | Unread message count shown as badge | `GetUnreadCount` endpoint, `fetchUnreadCount` action in `message.ts` store, badge shown in nutritionist messages index | ✅ COVERED |
| MSG-06 | Messages ordered chronologically, no editing or deletion | `GetConversation` SQL uses `ORDER BY created_at ASC`; no UPDATE/DELETE message endpoints exist | ✅ COVERED |
| MSG-07 | Read receipts (read_at timestamp) | `MarkRead` endpoint sets `read_at = NOW()` for messages where `receiver_id = caller AND read_at IS NULL`; `read_at` included in `MessageResponse` | ✅ COVERED |

### Food Requests (FREQ-01 → FREQ-04)

| Req | Description | Implementation | Status |
|-----|-------------|----------------|--------|
| FREQ-01 | Client can submit food addition request with name and optional description | `CreateFoodRequest` handler + `food_requests` table + `ClientCreateFoodRequest` route | ✅ COVERED |
| FREQ-02 | Request goes to client's assigned nutritionist | `CreateFoodRequest` SQL inserts `nutritionist_id` from client's `nutritionist_id` FK column | ✅ COVERED |
| FREQ-03 | Nutritionist can approve (creates food item) or reject (with optional reason) | `ApproveFoodRequest` returns `food_name` → frontend redirects to `/nutritionist/foods/create?name=X`; `RejectFoodRequest` stores `reject_reason` | ✅ COVERED |
| FREQ-04 | Client receives notification of approval/rejection | `ListClientFoodRequests` returns current status + reject_reason; client food-requests page shows result badge | ✅ COVERED |

### Client Management (CLNT-02 → CLNT-07)

| Req | Description | Implementation | Status |
|-----|-------------|----------------|--------|
| CLNT-02 | Client list view: name, mobile, status, plan status, last activity, searchable by name/mobile | `NutriListClients` handler + `SearchClients` raw SQL with ILIKE; `ClientListItemResponse` DTO; `nutritionist/clients.vue` search bar | ✅ COVERED |
| CLNT-03 | Filterable by active/inactive, sortable by name or last activity | `SearchClients` accepts `is_active` filter + dynamic `ORDER BY` (name/last_activity); filter tabs + sort toggle in `nutritionist/clients.vue` | ✅ COVERED |
| CLNT-04 | Client profile view with personal info, current plan summary, and history tabs | `NutriGetClientProfile` handler + `nutritionist/clients/[clientId]/index.vue` with Overview/Tracking/Plans tabs | ✅ COVERED |
| CLNT-05 | Quick actions from client profile: new diet plan, send message, deactivate client | Profile page buttons link to plan creation, message conversation, and call activate/deactivate store actions | ✅ COVERED |
| CLNT-06 | Nutritionist can activate/deactivate clients | `NutriActivateClient` + `NutriDeactivateClient` handlers; `ActivateClient`/`DeactivateClient` service methods toggle `is_active` | ✅ COVERED |
| CLNT-07 | Height and date of birth editable only by nutritionist | `NutriUpdateClientProfile` handler (nutritionist-auth middleware); `height_cm`/`date_of_birth`/`gender` in `UpdateClientProfileRequest`; no client-side edit endpoint for these fields | ✅ COVERED |

### Security (SEC-04, SEC-05, SEC-08)

| Req | Description | Implementation | Status |
|-----|-------------|----------------|--------|
| SEC-04 | File upload validation: type checking, size limits, magic byte verification, UUID filenames | `communication_service.go`: `mimetype.DetectReader` for magic bytes, size limit per MIME type, `uuid.New()` for filename | ✅ COVERED |
| SEC-05 | Content-Disposition: attachment on file downloads to prevent content sniffing | `DownloadAttachment` handler sets `Content-Disposition: attachment; filename="..."` header | ✅ COVERED |
| SEC-08 | Per-client file storage limits | `communication_service.go` enforces per-upload size limits (5MB images, 10MB PDFs) at the service layer | ⚠️ PARTIAL — aggregate per-client quota tracking not implemented (per-file limits enforced) |

## Success Criteria Checklist

| # | Criterion | Evidence | Status |
|---|-----------|----------|--------|
| 1 | Text + optional attachments, 10s polling, unread badge, read receipts, chronological | See MSG-01 → MSG-07 above | ✅ |
| 2 | Food request submit → approve/reject → client sees result | See FREQ-01 → FREQ-04 above | ✅ |
| 3 | Client list searchable/filterable/sortable; client profile with history tabs | See CLNT-02 → CLNT-04 above | ✅ |
| 4 | Quick actions from client profile: plan, message, activate/deactivate | See CLNT-05, CLNT-06, CLNT-07 above | ✅ |
| 5 | File uploads validated (magic bytes, size, UUID names, attachment downloads) | See SEC-04, SEC-05 above | ✅ |

## Known Gaps / Deferrals

- **SEC-08 aggregate quota**: Per-file limits (5MB/10MB) are enforced, but total-per-client storage quota tracking is not implemented. This is noted for Phase 7 hardening.
- **FREQ-04 push notification**: "Client receives notification" is fulfilled via polling (visible on next list refresh), not push notification — push notifications are deferred to Phase 6 (NOTIF-*).

## Phase Goal: ACHIEVED

All 20 Phase 5 requirements are covered (19 fully, 1 partially). The communication and collaboration layer is fully functional with backend API, database schema, and frontend UI complete and building cleanly.
