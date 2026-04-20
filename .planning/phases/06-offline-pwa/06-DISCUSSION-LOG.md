# Phase 6: Offline & PWA - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-20
**Phase:** 06-offline-pwa
**Areas discussed:** PWA shell strategy, offline cache boundaries, sync queue behavior, push notification architecture, client offline UX

---

## PWA shell strategy

| Option | Description | Selected |
|--------|-------------|----------|
| InjectManifest custom service worker | Explicit control over runtime caching, push handling, and background sync hooks | ✓ |
| Zero-config PWA preset | Faster setup but weaker control for custom sync/notification flows | |
| Let the agent decide | Defer selection to implementation | |

**User's choice:** Auto-selected recommended option: custom `injectManifest` service worker.
**Notes:** Phase 6 requires runtime caching, push click handling, and sync orchestration, so explicit service-worker control is safer than a minimal preset.

---

## Offline cache boundaries

| Option | Description | Selected |
|--------|-------------|----------|
| Cache active plan + last 50 messages + normalized sync queue | Covers the required offline read/write loop with bounded storage | ✓ |
| Cache every client API response aggressively | Larger offline surface but higher invalidation/storage complexity | |
| Let the agent decide | Defer selection to implementation | |

**User's choice:** Auto-selected recommended option: bounded cache focused on plan, messages, and queued writes.
**Notes:** This matches roadmap success criteria without expanding into unnecessary offline copies of every endpoint.

---

## Sync queue behavior

| Option | Description | Selected |
|--------|-------------|----------|
| Central FIFO queue with single-flight sync and exponential backoff | Predictable ordering, simpler reconciliation, clearer status UI | ✓ |
| Per-store independent sync loops | More parallelism but fragmented state and harder retry handling | |
| Let the agent decide | Defer selection to implementation | |

**User's choice:** Auto-selected recommended option: one shared FIFO queue.
**Notes:** Existing `local_id` support from Phase 4 makes a single normalized queue the most coherent design.

---

## Push notification architecture

| Option | Description | Selected |
|--------|-------------|----------|
| `webpush-go` + persisted subscriptions/preferences + server ticker for reminders | Matches approved stack and supports both event-driven and scheduled pushes | ✓ |
| Client-side local reminders only | Simpler but does not satisfy server-driven message/plan notifications | |
| Let the agent decide | Defer selection to implementation | |

**User's choice:** Auto-selected recommended option: server-driven Web Push with VAPID.
**Notes:** Scheduled meal/medication/water reminders and cross-device delivery require backend participation.

---

## Client offline UX

| Option | Description | Selected |
|--------|-------------|----------|
| Visible sync-status indicator + offline saved toast + profile-based notification settings | Clear operational feedback without redesigning core screens | ✓ |
| Silent background sync with minimal user feedback | Cleaner UI but poor recoverability for failed queue items | |
| Let the agent decide | Defer selection to implementation | |

**User's choice:** Auto-selected recommended option: visible sync/install status with manual retry path.
**Notes:** Phase 6 explicitly requires pending counts, sync states, and manual retry for failed items.

---

## the agent's Discretion

- Exact visual treatment of install/update prompt
- Exact iconography for sync states and notification preference toggles
- Concrete runtime cache naming and Dexie schema versioning details

## Deferred Ideas

- Native app background scheduling beyond browser/Web Push capabilities
- Offline support for nutritionist/admin surfaces
- Message search or richer inbox features
