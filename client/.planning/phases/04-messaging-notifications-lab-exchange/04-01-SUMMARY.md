---
phase: 04-messaging-notifications-lab-exchange
plan: "01"
subsystem: Types & API Contracts
tags: [types, contracts, composables, api]
dependency_graph:
  requires: []
  provides: [messaging-types, lab-types, notification-types, messaging-api, lab-api, notification-api]
  affects: [04-02, 04-03, 04-04, 04-05]
tech_stack:
  added: [TypeScript interfaces, Nuxt composables, useAsyncData, FormData]
  patterns: [Composable pattern, API wrapper pattern]
key_files:
  created:
    - app/types/messaging.ts
    - app/types/lab.ts
    - app/types/notifications.ts
    - app/composables/useMessagingApi.ts
    - app/composables/useLabApi.ts
    - app/composables/useNotificationApi.ts
  modified: []
decisions: []
metrics:
  duration_minutes: 12
  tasks_completed: 2
  files_created: 6
  lines_added: 224
  completion_date: 2026-04-23
---

# Phase 04 Plan 01: Type Contracts & API Composables Summary

## Objective Fulfilled
✓ Define typed API contracts for messaging, lab results, and notification preferences as Wave 1 foundation.  
✓ All Wave 2 and Wave 3 plans depend on these contracts.  
✓ Separated contracts from UI implementation to ensure early type safety.

## What Was Built

### 1. Type Definitions (3 files)
- **app/types/messaging.ts**  
  - `Message` interface with sender/receiver IDs, content, attachment, read state  
  - `MessageAttachment` with URL, MIME type, size, filename  
  - `SendMessageRequest` with optional content and file  
  - `UnreadCountResponse` wrapper  
  - `PaginatedMessages` for list responses  

- **app/types/lab.ts**  
  - `LabResult` interface with file metadata (original_name, file_type, file_size) and optional link  
  - `UploadLabResultRequest` for client/nutritionist uploads  
  - `LabResourceType` discriminator type ('file' | 'link')  
  - **Helper function**: `getLabResourceType(result: LabResult): LabResourceType`  
    Returns 'link' if result.link is not null, otherwise 'file'  

- **app/types/notifications.ts**  
  - `NotificationPreferences` with four boolean toggles  
    - meal_reminders, water_reminders, message_alerts, diet_updates  
  - `UpdateNotificationPreferencesRequest` as Partial<NotificationPreferences>  

### 2. Composable APIs (3 files)
- **app/composables/useMessagingApi.ts**  
  - `getClientConversation(page?, pageSize?)`: useAsyncData GET /api/v1/messages  
  - `sendClientMessage(req)`: $fetch POST with FormData (content + file support)  
  - `getNutritionistConversation(clientId, page?, pageSize?)`: useAsyncData GET  
  - `sendNutritionistMessage(clientId, req)`: $fetch POST /api/v1/clients/:id/messages  
  - `getUnreadCount()`: $fetch GET /api/v1/messages/unread-count  

- **app/composables/useLabApi.ts**  
  - `listLabResults(clientId, page?, pageSize?)`: useAsyncData GET /api/v1/clients/:id/lab-results  
  - `uploadLabResult(clientId, req)`: $fetch POST with FormData (multipart fields)  
  - `getDownloadUrl(labId): string`: Returns "/api/v1/lab-results/:id/download" (pure function)  

- **app/composables/useNotificationApi.ts**  
  - `getPreferences()`: useAsyncData GET /api/v1/notifications/preferences  
  - `updatePreferences(req)`: $fetch PATCH /api/v1/notifications/preferences  
  - `registerPushSubscription(subscription)`: $fetch POST /api/v1/push/subscribe  
  - `unregisterPushSubscription()`: $fetch DELETE /api/v1/push/subscribe  

## Verification

✓ All 6 files created at specified paths  
✓ All type exports verified (5 + 4 + 2 types)  
✓ All composable functions exported  
✓ FormData patterns follow existing tracking composable conventions  
✓ useAsyncData + $fetch patterns consistent with useTrackingApi  

## Deviations
None — executed exactly per plan specification.

## Threat Analysis
- **T-04-01-01** (Info Disclosure): No auth headers hardcoded; handled by auth-fetch.client.ts plugin ✓  
- **T-04-01-02** (Tampering): Backend enforces file size/type limits; client validation in later plans  
- **T-04-01-03** (Info Disclosure): getDownloadUrl returns URL only; auth-fetch injects Bearer token  

## Next Steps
- ✅ 04-01 complete
- ⏭️ 04-02 (Client messaging page, polling, offline queue) — ready to start
- ⏭️ 04-03 (Nutritionist conversations) — ready to start
- ⏭️ 04-04 (Lab upload & results) — ready to start
- ⏭️ 04-05 (Push notifications & preferences) — ready to start

## Known Stubs
None — all types fully defined with no placeholder values.

## Git Commit
```
feat(04-01): add type contracts for messaging, lab, and notifications [6d6a6cb]
```
