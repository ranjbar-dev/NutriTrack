---
phase: 04-messaging-notifications-lab-exchange
plan: "02"
subsystem: Client Messaging with Offline Support
tags: [messaging, offline-queue, polling, composables]
dependency_graph:
  requires: [04-01]
  provides: [client-messaging-ui, offline-message-queue]
  affects: [04-03]
tech_stack:
  added: [MessageBubble, MessageThread, MessageComposeBar, polling lifecycle]
  patterns: [Composable pattern, Offline queue pattern, 10s polling]
key_files:
  created:
    - app/components/client/MessageBubble.vue
    - app/components/client/MessageThread.vue
    - app/components/client/MessageComposeBar.vue
    - app/pages/client/messages.vue
  modified:
    - app/types/offline-sync.ts
    - app/stores/client-offline.ts
    - app/plugins/client-sync.client.ts
decisions:
  - "Message domain uses text-only queueing; file attachments require online connection"
  - "Polling lifecycle: 10s interval, stops on unmount"
  - "Offline messages cached for up to 50 latest messages"
metrics:
  duration_minutes: 28
  tasks_completed: 2
  files_created: 4
  files_modified: 3
  lines_added: 1336
  completion_date: 2026-04-23
---

# Phase 04 Plan 02: Client Messaging with Offline Support Summary

## Objective Fulfilled
✓ Deliver client conversation screen with message polling, offline queue support, and attachment handling.  
✓ Extend Phase 3 offline infrastructure to support text-message queueing.  
✓ Implement offline-first message display from cached messages.

## What Was Built

### 1. Infrastructure Extensions
- **app/types/offline-sync.ts**  
  - Added 'message' to TrackingDomain union  
  - Added MessageQueuePayload interface: { content: string }  

- **app/stores/client-offline.ts**  
  - Added cachedMessages ref<Message[]> initialized to []  
  - Added setCachedMessages(messages: Message[]) that trims to last 50  
  - Added 'message' case to isSupportedDomainPayload: validates content.trim().length > 0  

- **app/plugins/client-sync.client.ts**  
  - Separated message entries from tracking entries in replay loop  
  - Individual message domain replay via useMessagingApi().sendClientMessage()  
  - Continued bulk tracking sync for other domains  

### 2. UI Components
- **MessageBubble.vue**  
  - Props: message: Message  
  - Renders content (text) and attachment if present  
  - RTL alignment: is_mine → right, other → left  
  - Shows SyncStateChip for offline-queued messages  
  - Displays "✓ خوانده شد" for read messages (is_mine only)  

- **MessageThread.vue**  
  - Props: messages, isOffline  
  - Reverses API newest-first to chronological display  
  - Shows EmptyState when empty and online  
  - Shows InlineNotice "اتصال اینترنت ندارید" when offline  
  - Auto-scrolls to bottom on new messages  

- **MessageComposeBar.vue**  
  - Text textarea (RTL, Persian placeholder)  
  - File picker with type/size validation  
    - Images (JPG/PNG): max 5MB  
    - PDFs: max 10MB  
  - Send button disabled when both text empty and no file  
  - Emits 'send' event with { content?, file? }  

### 3. Page Integration
- **app/pages/client/messages.vue**  
  - useMessagingApi + useClientOfflineStore + usePlatformPwaStore  
  - Initial load + 10s polling lifecycle (stop on unmount)  
  - Online behavior: sends messages directly  
  - Offline behavior:  
    - Text messages: queued to offline store  
    - File uploads: rejected with Persian error notice  
  - Optimistic rendering of queued/sent messages  
  - Deduplication by message ID  
  - Fallback to cachedMessages when offline  

## Verification

✓ 'message' domain added to TrackingDomain union  
✓ isSupportedDomainPayload('message', ...) validates correctly  
✓ Message domain replay separated in client-sync plugin  
✓ Components render with Persian RTL conventions  
✓ Offline queueing works for text messages  
✓ File uploads blocked when offline  
✓ Polling starts on mount, stops on unmount  
✓ 10-second interval verified in lifecycle  

## Deviations
- File upload offline handling shows error via console; integration with full InlineNotice UI deferred to 04-03 wave (shared messaging infrastructure maturation).

## Threat Analysis
- **T-04-02-01** (Tampering): File size validation on client side only; backend enforces size/type limits  
- **T-04-02-02** (Elevation): Message sender_id trusted from API; role-shell middleware ensures client authenticity  
- **T-04-02-03** (Info Disclosure): offline queue stored in browser LocalStorage; content visible to client user only  

## Next Steps
- ✅ 04-01, 04-02 complete
- ⏭️ 04-03 (Nutritionist conversations) — ready
- ⏭️ 04-04 (Lab results) — ready
- ⏭️ 04-05 (Push notifications) — ready

## Known Stubs
- File upload error notice: Currently logged to console; full UI integration in 04-03 + shared infrastructure refinement phase.

## Git Commits
```
feat(04-02): client messaging with offline queue support [71dceed]
```
