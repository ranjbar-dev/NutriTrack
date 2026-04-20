---
phase: "06"
plan: "05"
subsystem: frontend/stores
tags: [offline, pwa, cache-first, local-echo, messages, dexie, indexeddb]
dependency_graph:
  requires: ["06-03", "06-04"]
  provides: ["D-07-messages", "D-08-messages"]
  affects: ["frontend/app/stores/message.ts", "frontend/app/pages/client/messages.vue", "frontend/app/pages/nutritionist/messages/[partnerId].vue"]
tech_stack:
  added: []
  patterns: ["cache-first read", "local echo optimistic UI", "offline-then-sync", "ID-based merge"]
key_files:
  created:
    - frontend/tests/clientPlan.offline.test.ts
  modified:
    - frontend/app/stores/message.ts
    - frontend/app/pages/client/messages.vue
    - frontend/app/pages/nutritionist/messages/[partnerId].vue
decisions:
  - "sendMessage return type changed to Promise<void>; local echo added to messages array directly for optimistic UI — components updated accordingly"
  - "hook callback cast `as never` consistent with syncQueue.ts callHook pattern to suppress Nuxt strict hook typing"
  - "Message.attachment_type restricted to 'image'|'pdf' matching existing types.ts — not 'file'"
metrics:
  duration: "~20 min"
  completed: "2026-04-20"
  tasks_completed: 2
  files_changed: 4
---

# Phase 06 Plan 05: Message Store Cache-First Reads + Offline Send Queue Summary

**One-liner:** Cache-first message fetch using Dexie (D-07) + local-echo offline send queue (D-08) with ID-based dedup merge in useMessageStore.

## What Was Built

### Task 1: Extended `frontend/app/stores/message.ts`

**D-07 — Cache-first reads:**
- `fetchMessages()` now reads Dexie `messages` table indexed by `partner_id` before any network call
- If offline, returns cached data (or Persian error if cache empty)
- On successful fetch, merges server data with cached data using `Map<id, Message>` dedup — preserves local echoes
- `persistMessages()` helper saves last 50 messages per `partner_id` to IndexedDB

**D-08 — Offline send queue:**
- `sendMessage()` return type changed to `Promise<void>`
- Creates `localEcho: Message` immediately and pushes to `messages` array (instant UI update, no network wait)
- Uses `clientPost()` from `useOfflineApi` — if offline, enqueued; if online, sent directly
- On successful server response, replaces local echo (`local_${uuid}`) with server message
- Handles file attachments via `attachmentBlob` option for sync queue storage

**D-14 — Hook integration:**
- Registers `sync:itemSynced` hook to replace local echo when sync manager processes a queued message
- Uses `as never` double-cast consistent with syncQueue.ts pattern for Nuxt hook strict typing

### Task 2: Offline tests (`frontend/tests/clientPlan.offline.test.ts`)

- **OFFL-11:** Verifies messages table stores 55 entries, last-50 trim logic returns 50 starting from `msg-5`
- **OFFL-02:** Verifies `activePlan` singleton cached and retrievable from Dexie
- Both tests use isolated `TestDB` with fresh Dexie instance (`Date.now()` suffix) to avoid cross-test contamination
- All 20 tests pass (14 passed + 6 todo)

### Updated Vue components

Both `client/messages.vue` and `nutritionist/messages/[partnerId].vue`:
- Removed `const msg = await sendMessage(...)` pattern
- Updated to `await sendMessage(...)` then check `messageStore.error` for form reset

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Adapted localEcho to match actual Message interface**
- **Found during:** Task 1 implementation
- **Issue:** Plan specified `content: content ?? null`, `attachment_type: 'image' | 'file' | null`, `attachment_path`, `read_at: null` — but the actual `Message` type has `content?: string`, `attachment_type?: 'image' | 'pdf'`, no `attachment_path`, `read_at?: string`
- **Fix:** Used conditional spread `...(content ? { content } : {})`, `attachment_type: file ? (...'image' : 'pdf') : undefined`, removed `attachment_path`
- **Files modified:** `frontend/app/stores/message.ts`
- **Commit:** cf12d83

**2. [Rule 1 - Bug] Updated Vue components for void return type**
- **Found during:** Task 1 — both `messages.vue` files check `if (msg)` on return value
- **Issue:** Changing `sendMessage` to `Promise<void>` breaks the `if (msg)` form-reset pattern
- **Fix:** Updated both components to check `!messageStore.error` instead
- **Files modified:** `frontend/app/pages/client/messages.vue`, `frontend/app/pages/nutritionist/messages/[partnerId].vue`
- **Commit:** cf12d83

**3. [Rule 1 - Bug] Fixed hook callback TypeScript error**
- **Found during:** TypeScript typecheck
- **Issue:** `nuxtApp.hook('sync:itemSynced' as never, handler)` fails because with name cast to `never`, TypeScript expects handler type `never` too
- **Fix:** Cast handler `as never` as well — consistent with syncQueue.ts `callHook` pattern
- **Files modified:** `frontend/app/stores/message.ts`
- **Commit:** cf12d83

## Known Stubs

None — all functionality is fully wired:
- `fetchMessages` reads and writes Dexie
- `sendMessage` uses `clientPost` from `useOfflineApi`
- `pollNewMessages` updates Dexie cache on success

## Threat Flags

None — no new network endpoints or auth paths introduced. The message store only consumes existing `/messages/*` endpoints.

## Self-Check: PASSED

- `frontend/app/stores/message.ts` — ✓ exists, modified
- `frontend/tests/clientPlan.offline.test.ts` — ✓ created
- `frontend/app/pages/client/messages.vue` — ✓ updated
- `frontend/app/pages/nutritionist/messages/[partnerId].vue` — ✓ updated
- Commit `cf12d83` — ✓ exists in git log
- All 14 tests pass (+ 6 todo) — ✓ confirmed
- `stores/message.ts` has zero TypeScript errors — ✓ confirmed via nuxi typecheck filter
