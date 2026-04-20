# Phase 6: Offline & PWA — Research

**Researched:** 2026-04-20
**Domain:** Progressive Web App, Service Workers, IndexedDB (Dexie.js), Web Push (VAPID), Background Sync, Go goroutine scheduling
**Confidence:** HIGH (stack verified; libraries confirmed against npm registry and golang proxy)

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

- **D-01** `injectManifest` strategy (custom SW), not zero-config preset — Phase 6 needs explicit runtime caching, notification click handling, and Background Sync registration hooks.
- **D-02** PWA applies to the whole build; offline behavior is client-role scoped. Static assets cached globally; diet-plan, tracking, messaging, and sync features only activate inside client routes and client stores.
- **D-03** Static assets use cache-first with versioned precache. Client API reads (active plan, message history) use network-first with IndexedDB fallback.
- **D-04** Manifest is Persian-only: `name: نوتری‌ترک`, `display: standalone`, app shortcuts for برنامه, ثبت روزانه, پیام‌ها.
- **D-05** Dexie.js as the single offline database. Tables: `activePlan`, `messages`, `syncQueue`, `syncMeta`, `notificationPreferences`, `uiState`. No per-domain tracking tables — offline write source is normalized queue payload.
- **D-06** One active plan snapshot per client (including all nested data + `fetched_at`, `plan_id`, `updated_hint`). Refresh on app open when online.
- **D-07** Last 50 messages per conversation cached in IndexedDB. Immediate cached reads, then online merge by message ID.
- **D-08** Queued messages may include attachment blobs stored in IndexedDB.
- **D-09** Wrap existing `useApi()` with an offline-aware client API layer — not duplicating fetch logic per store.
- **D-10** Queue all Phase 4 + 5 client POST requests that carry `local_id`: food logs, water logs, sleep logs, exercise logs, medication logs, body measurements, lab-result metadata, outgoing messages.
- **D-11** `syncQueue` entry: `id`, `entity_type`, `request_path`, `method`, serialized payload, optional attachment blob, `local_id`, `created_at`, `status`, `retry_count`, `last_error`, `next_attempt_at`.
- **D-12** Retry: exponential backoff 1s → 2s → 4s, max 3 attempts, then `failed` for manual retry.
- **D-13** Reconnect detection: browser `online` events + app-open sync sweep. Background Sync where supported; foreground polling fallback.
- **D-14** Sync success updates in-memory stores immediately. Server-side `local_id` dedup from Phase 4 is the canonical duplicate guard.
- **D-15** Backend: `push_subscriptions` + `notification_preferences` tables.
- **D-16** Client subscribe/unsubscribe endpoints, idempotent by endpoint URL.
- **D-17** `webpush-go` with VAPID keys from env. One JSON envelope: Persian `title`, `body`, `action_url`, `icon`, `type`, optional entity IDs.
- **D-18** Event-driven push from Phase 3/5 service layer + goroutine ticker every minute for reminders.
- **D-19** Preference-aware scheduling; dedup key per reminder window to prevent duplicate sends.
- **D-20** Sync-status surface in client UI: syncing / synced / X pending. Failed items expose manual retry.
- **D-21** Render cached data first; offline empty state if no cache.
- **D-22** Notification preferences under client profile/settings, not a separate admin page.
- **D-23** iOS eviction: on app open detect missing tables, clear stale UI flags, show Persian notice, re-fetch on reconnect.

### Agent's Discretion

- Exact install-prompt copy and iconography (must be Persian, mobile-first)
- Sync-status indicator placement: header or floating pill, visible from plan/tracking/messages
- Failed sync item display: inline per screen or centralized retry sheet
- Exact Workbox cache names and Dexie schema version number

### Deferred Ideas (OUT OF SCOPE)

- Background push delivery for nutritionist/admin roles
- Rich message search, adaptive polling, inbox summarization
- Native-device reminder scheduling outside Web Push/browser capabilities

</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| OFFL-01 | Service Worker caches static assets + API responses for client role only | `@vite-pwa/nuxt` injectManifest + Workbox strategies — §Standard Stack, §Architecture Patterns |
| OFFL-02 | Active diet plan fully viewable offline | Dexie `activePlan` table + network-first-with-IDB-fallback — §Dexie Schema, §Plan Cache Pattern |
| OFFL-03 | All tracking types work offline with queue | `useClientApi` offline wrapper + `syncQueue` Dexie table — §Offline-Aware API Layer |
| OFFL-04 | Queued messages sendable offline | Message offline send path with FormData + Blob serialization — §Message Offline Pattern |
| OFFL-05 | IndexedDB via Dexie.js: activePlan, pending logs, cached messages, outgoing queue | Dexie v4 schema definition — §Dexie Schema |
| OFFL-06 | Sync manager: FIFO, reconnect, last-write-wins | `useSyncManager` composable — §Sync Manager Pattern |
| OFFL-07 | Each queued entry has local_id + synced_at; server deduplicates via ON CONFLICT | Phase 4 backend already implements this; planner must ensure all synced entries report `synced_at` — §Integration Points |
| OFFL-08 | Exponential backoff max 3 retries, then manual retry flag | Retry logic in `useSyncManager` — §Sync Manager Pattern |
| OFFL-09 | Background Sync API where supported, fallback to polling | SW Background Sync registration + foreground `online` listener — §Background Sync |
| OFFL-10 | Diet plan cached on first load, refreshed on app open if online | `clientPlan` store extension + `syncMeta` table — §Plan Cache Pattern |
| OFFL-11 | Cached last 50 messages per conversation | `messages` Dexie table, 50-row window per `partner_id` — §Message Cache Pattern |
| OFFL-12 | iOS PWA storage eviction handled gracefully | Detection + re-fetch logic, Persian notice — §iOS Pitfalls |
| NOTIF-01 | Web Push via VAPID keys using webpush-go | `webpush-go` v1.4.0, VAPID config in env — §Standard Stack |
| NOTIF-02 | Client subscribes to push on first login | `useNotificationPermission` composable + `POST /api/push/subscribe` — §Push Frontend |
| NOTIF-03 | Triggers: new message, new diet plan, food request result | Event-driven hooks in `communication_service.go` + `diet_plan_service.go` — §Push Backend |
| NOTIF-04 | Meal time reminders based on diet plan scheduled times | Goroutine scheduler querying meal times — §Reminder Scheduler |
| NOTIF-05 | Medication reminders based on prescribed medication times | Same scheduler, `plan_medications.times` JSONB field — §Reminder Scheduler |
| NOTIF-06 | Water intake reminders | Fixed-time daily reminder in scheduler — §Reminder Scheduler |
| NOTIF-07 | Client can enable/disable each reminder type in notification preferences | `notification_preferences` table + prefs page — §Push Backend, §Push Frontend |
| NOTIF-08 | Notification payload: title, body, action_url, icon, type | Single push envelope struct — §Push Backend |
| UI-06 | PWA manifest with install prompt | `@vite-pwa/nuxt` manifest config — §Standard Stack |
| UI-07 | Service worker with registerType: autoUpdate | `@vite-pwa/nuxt` `registerType: 'autoUpdate'` — §Standard Stack |

</phase_requirements>

---

## Summary

Phase 6 transforms an already-complete online PWA into a fully offline-capable application by layering three complementary systems: a Workbox-backed service worker for static and API caching, a Dexie.js IndexedDB store for offline data persistence and sync queue management, and a Go-side Web Push notification pipeline using VAPID and a goroutine-based reminder scheduler.

The most important architectural constraint is the integration seam: `useApi.ts` is the correct and only place to intercept outgoing client writes and redirect them to the Dexie sync queue when offline. Touching individual store files minimally (to read from cache first) and routing writes through a single `useClientApi` composable keeps the offline logic contained. The backend already has all the `local_id` deduplication infrastructure in place from Phase 4 — Phase 6 simply consumes it.

Push notifications split cleanly: event-driven pushes (new message, new plan, food request outcome) fire synchronously from within existing service methods; scheduled reminders (meals, medications, water) run in a separate goroutine that ticks every minute against the active plan database. Both paths share a single `NotificationService` that reads preferences, checks the dedup key, and calls `webpush-go`.

**Primary recommendation:** Execute in 6 waves — PWA shell → plan cache → sync queue + tracking → offline messaging → backend push → frontend push preferences. Each wave is independently testable. The most complex wave is Wave 3 (sync queue), and the most platform-risky is Wave 5 (iOS push).

---

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Static asset caching | Service Worker | CDN/Traefik | SW precaches all Nuxt build artifacts; Traefik serves source files |
| Runtime API caching (plan, messages) | Service Worker | IndexedDB | SW intercepts fetch, falls through to IDB on miss |
| Offline data persistence | Browser (IndexedDB/Dexie) | — | Client-side only; no server persistence of offline queue state |
| Sync queue processing | Frontend (composable) | Service Worker (Background Sync) | Foreground primary; SW Background Sync is enhancement |
| Push subscription management | API / Backend | Browser (SW) | Backend stores subscriptions; SW receives push events |
| Reminder scheduling | API / Backend | — | Goroutine ticker in Go process; no browser scheduling |
| Notification preferences | API / Backend | Frontend Store | Backend is source of truth; frontend reflects state |
| PWA manifest + install prompt | Frontend (Nuxt config) | Browser | nuxt.config.ts manifest fed to browser |
| iOS eviction recovery | Frontend (composable) | — | Detected at app-open via Dexie table existence checks |

---

## Standard Stack

### Core (new additions to install)

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `@vite-pwa/nuxt` | 1.1.1 | PWA module for Nuxt 4 — manifest, SW registration, `$pwa` composable | Official Nuxt PWA module wrapping vite-plugin-pwa + Workbox; `injectManifest` strategy gives full custom SW control |
| `dexie` | 4.4.2 | IndexedDB wrapper for offline storage | Versioned schema, TypeScript-first in v4, compound indexes, simpler than raw IDB; explicitly chosen in STACK.md |
| `webpush-go` | v1.4.0 | VAPID-based Web Push delivery from Go backend | Only maintained Go Web Push library; handles VAPID key generation, payload encryption, subscription management |

> [VERIFIED: npm registry] `@vite-pwa/nuxt@1.1.1`, `dexie@4.4.2`, `vite-plugin-pwa@1.2.0`
> [VERIFIED: golang proxy] `webpush-go@v1.4.0` (2025-01-02)

### Already in project (used by Phase 6)

| Library | Current Version | Phase 6 Usage |
|---------|----------------|---------------|
| `@pinia/nuxt` | 0.11.3 | All offline state surfaces through existing Pinia stores |
| `nuxt` | ^4.0.0 | `@vite-pwa/nuxt` integrates via modules array |
| `dayjs` | 1.11.20 | Already in package.json (despite STACK.md recommendation against it — keep for continuity, not worth removing now) |
| `github.com/google/uuid` | 1.6.0 | Used in notification service for dedup keys |
| `github.com/rs/zerolog` | 1.35.0 | Logging in notification service and scheduler |
| `golang.org/x/crypto` | 0.50.0 | Used internally by webpush-go |

### Installation

```bash
# Frontend
cd frontend
npm install dexie@4.4.2
npx nuxi@latest module add @vite-pwa/nuxt

# Backend
cd backend
go get github.com/SherClockHolmes/webpush-go@v1.4.0
go mod tidy
```

**VAPID key generation (one-time setup):**
```bash
# Using webpush-go CLI or any VAPID generator — do this once, store in env
go run -mod=mod github.com/SherClockHolmes/webpush-go/cmd/vapid-keygen@v1.4.0
# Outputs: VAPID_PUBLIC_KEY and VAPID_PRIVATE_KEY (base64url encoded)
```

---

## Architecture Patterns

### System Architecture Diagram

```
CLIENT (Browser / PWA)
┌─────────────────────────────────────────────────────────────────┐
│  Nuxt App (client.vue layout, plan/tracking/messages pages)     │
│       │  reads cached first       │  write action               │
│       ▼                           ▼                             │
│  useClientApi (offline wrapper)──→ navigator.onLine?            │
│       │ online: apiFetch()        │ offline: Dexie syncQueue    │
│       ▼                           ▼                             │
│  Pinia Stores ◄─── merge ─── useSyncManager                    │
│       │                       (on 'online' event + app-open)   │
│       ▼                                                         │
│  useDb() / Dexie Database                                       │
│  ┌──────────┬──────────┬──────────┬───────────┬──────────────┐ │
│  │activePlan│ messages │syncQueue │ syncMeta  │notifPrefs   │ │
│  └──────────┴──────────┴──────────┴───────────┴──────────────┘ │
└─────────────────────────────────────────────────────────────────┘
          ▲ fetch / sync              ▲ push events
          │                           │
Service Worker (sw.ts / Workbox)       │
  ├── precache: JS/CSS/fonts/icons     │
  ├── runtime cache: /api/clients/me/active-plan (network-first)  │
  ├── runtime cache: /api/messages/* (network-first)              │
  ├── Background Sync: 'sync-queue-tag' event                     │
  └── push handler → showNotification() + click → navigate       │
          │                           │
          ▼                           │
Go API (Gin)                          │ webpush-go
  ├── /api/push/subscribe (POST)      │
  ├── /api/push/preferences (GET/PATCH)│
  ├── /api/client/* (tracking, plan, messages)
  │       │
  ├── NotificationService.Send() ◄─── triggered by CommunicationService.SendMessageTo()
  │                                                  DietPlanService.ActivatePlan()
  │                                                  CommunicationService.ApproveFoodRequest()
  │                                                  CommunicationService.RejectFoodRequest()
  └── ReminderScheduler (goroutine ticker, every 60s)
          └── queries active plans → meal times / medication times / water windows
                  └── NotificationService.Send() with dedup key check
```

### Recommended Project Structure (new files only)

```
frontend/
├── public/
│   ├── sw.ts                        # Custom service worker (injectManifest)
│   └── icons/                       # PWA icons (72, 96, 128, 144, 152, 192, 384, 512)
├── app/
│   ├── composables/
│   │   ├── useDb.ts                 # Dexie database singleton + schema
│   │   ├── useClientApi.ts          # Offline-aware API wrapper (wraps useApi)
│   │   ├── useSyncManager.ts        # Queue processor composable
│   │   └── useNotificationPermission.ts  # Push subscription flow
│   ├── stores/
│   │   └── notificationPrefs.ts     # Notification preferences Pinia store
│   ├── components/client/
│   │   ├── SyncStatusBadge.vue      # Syncing / synced / X pending indicator
│   │   └── OfflineBanner.vue        # Offline mode notice banner
│   ├── pages/client/
│   │   └── settings/
│   │       └── notifications.vue    # Preference toggles page
│   └── layouts/
│       └── client.vue               # MODIFIED: add SyncStatusBadge + OfflineBanner

backend/
├── db/migrations/
│   └── 000010_push_notifications.up.sql   # push_subscriptions + notification_preferences
├── internal/
│   ├── repository/
│   │   └── notification_repo.go     # Push subscription + preferences CRUD
│   ├── service/
│   │   ├── notification_service.go  # Push delivery + dedup logic
│   │   └── reminder_scheduler.go    # Goroutine ticker, reminder logic
│   └── handler/
│       └── notification_handler.go  # Subscribe, unsubscribe, preferences endpoints
└── cmd/api/
    └── main.go                      # MODIFIED: wire notification service + scheduler
```

---

### Pattern 1: `@vite-pwa/nuxt` with `injectManifest` Strategy

**What:** Registers a custom service worker (`public/sw.ts`) that imports the Workbox precache manifest and applies explicit runtime caching rules. `registerType: 'autoUpdate'` handles SW lifecycle automatically (satisfies UI-07).

**When to use:** Any time you need push notification handling, Background Sync registration, or per-route caching logic beyond Workbox presets.

**nuxt.config.ts addition:**
```typescript
// Source: @vite-pwa/nuxt docs (injectManifest strategy)
import { defineNuxtConfig } from 'nuxt/config'

export default defineNuxtConfig({
  compatibilityDate: '2025-07-18',
  future: { compatibilityVersion: 4 },

  modules: ['@pinia/nuxt', '@nuxt/eslint', '@vite-pwa/nuxt'],

  pwa: {
    strategies: 'injectManifest',
    srcDir: 'public',
    filename: 'sw.ts',
    registerType: 'autoUpdate',          // UI-07: autoUpdate strategy
    injectManifest: {
      injectionPoint: 'self.__WB_MANIFEST',
    },
    manifest: {
      name: 'نوتری‌ترک',
      short_name: 'نوتری‌ترک',
      description: 'مدیریت برنامه غذایی',
      lang: 'fa',
      dir: 'rtl',
      display: 'standalone',
      background_color: '#ffffff',
      theme_color: '#16a34a',           // Tailwind green-600 — matches app accent
      icons: [
        { src: '/icons/icon-192.png', sizes: '192x192', type: 'image/png' },
        { src: '/icons/icon-512.png', sizes: '512x512', type: 'image/png', purpose: 'any maskable' },
      ],
      shortcuts: [
        { name: 'برنامه', url: '/client/plan', icons: [{ src: '/icons/shortcut-plan.png', sizes: '96x96' }] },
        { name: 'ثبت روزانه', url: '/client/tracking', icons: [{ src: '/icons/shortcut-track.png', sizes: '96x96' }] },
        { name: 'پیام‌ها', url: '/client/messages', icons: [{ src: '/icons/shortcut-msg.png', sizes: '96x96' }] },
      ],
      start_url: '/client/plan',
    },
    devOptions: {
      enabled: true,                    // Enable SW in dev for testing
      type: 'module',
    },
  },
  // ... rest of config
})
```

---

### Pattern 2: Custom Service Worker (`public/sw.ts`)

**What:** The injectManifest SW imports `self.__WB_MANIFEST` for precaching, sets up Workbox runtime caching strategies, registers Background Sync, and handles push events.

```typescript
// Source: vite-plugin-pwa injectManifest docs + Workbox docs
import { cleanupOutdatedCaches, precacheAndRoute } from 'workbox-precaching'
import { registerRoute, Route } from 'workbox-routing'
import { NetworkFirst, CacheFirst } from 'workbox-strategies'
import { BackgroundSyncPlugin } from 'workbox-background-sync'
import { ExpirationPlugin } from 'workbox-expiration'

declare let self: ServiceWorkerGlobalScope

// Precache all versioned build assets (JS, CSS, fonts) — cache-first implicitly
precacheAndRoute(self.__WB_MANIFEST)
cleanupOutdatedCaches()

// Runtime: API plan endpoint — network-first, IDB fallback handled in store
registerRoute(
  new Route(
    ({ url }) => url.pathname === '/api/clients/me/active-plan',
    new NetworkFirst({
      cacheName: 'nutritrack-plan-v1',
      networkTimeoutSeconds: 5,
      plugins: [new ExpirationPlugin({ maxEntries: 1, maxAgeSeconds: 86400 })],
    }),
  ),
)

// Runtime: Messages — network-first with 50-entry cache
registerRoute(
  new Route(
    ({ url }) => url.pathname.startsWith('/api/messages'),
    new NetworkFirst({
      cacheName: 'nutritrack-messages-v1',
      networkTimeoutSeconds: 5,
    }),
  ),
)

// Background Sync — fallback queue for tracking writes when SW intercepts
const bgSyncPlugin = new BackgroundSyncPlugin('nutritrack-sync-queue', {
  maxRetentionTime: 24 * 60, // 24 hours in minutes
})

// Push event handler
self.addEventListener('push', (event) => {
  if (!event.data) return
  const data = event.data.json() as {
    title: string; body: string; action_url: string; icon: string; type: string
  }
  event.waitUntil(
    self.registration.showNotification(data.title, {
      body: data.body,
      icon: data.icon || '/icons/icon-192.png',
      badge: '/icons/badge-72.png',
      dir: 'rtl',
      lang: 'fa',
      data: { action_url: data.action_url },
    }),
  )
})

// Notification click handler — navigate to action_url
self.addEventListener('notificationclick', (event) => {
  event.notification.close()
  const actionUrl = event.notification.data?.action_url || '/client/plan'
  event.waitUntil(
    self.clients.matchAll({ type: 'window' }).then((clients) => {
      const existing = clients.find(c => c.url.includes(self.location.origin))
      if (existing) return existing.navigate(actionUrl)
      return self.clients.openWindow(actionUrl)
    }),
  )
})

// Background sync event (where supported)
self.addEventListener('sync', (event) => {
  if (event.tag === 'nutritrack-sync') {
    // Tell the app to process its sync queue
    event.waitUntil(
      self.clients.matchAll().then(clients =>
        Promise.all(clients.map(c => c.postMessage({ type: 'TRIGGER_SYNC' }))),
      ),
    )
  }
})
```

> **Critical note on `public/sw.ts`:** Nuxt 4 copies `public/` to `.output/public/`. With `injectManifest`, vite-plugin-pwa compiles `sw.ts` → `sw.js` and injects the precache manifest. Do NOT use `app/` directory for the SW file — it must be in `public/`. [ASSUMED]

---

### Pattern 3: Dexie v4 Schema (`composables/useDb.ts`)

**What:** Singleton Dexie database with versioned tables per D-05. All Phase 6 composables import `useDb()`.

```typescript
// Source: Dexie.js v4 docs (TypeScript, versionSchema)
import Dexie, { type EntityTable } from 'dexie'

export interface ActivePlanRecord {
  id: number           // always 1 — singleton row
  plan_id: string
  data: object         // full DietPlanResponse JSON blob
  fetched_at: string   // ISO timestamp
  updated_hint: string // plan updated_at from server
}

export interface MessageRecord {
  id: string           // message UUID from server (or temp local_id for queued)
  partner_id: string   // conversation partner UUID
  data: object         // MessageResponse JSON blob
  is_local: boolean    // true = not yet synced
  cached_at: string
}

export interface SyncQueueEntry {
  id?: number          // auto-increment
  entity_type: string  // 'food_log' | 'water_log' | 'sleep_log' | 'exercise_log' | 'medication_log' | 'body_measurement' | 'lab_result_meta' | 'message'
  request_path: string // e.g. '/client/food-logs'
  method: string       // 'POST'
  payload: string      // JSON.stringify(body) — no blobs here; blobs separate
  attachment_blob?: Blob
  attachment_filename?: string
  attachment_mime?: string
  local_id: string     // UUID (crypto.randomUUID())
  created_at: string
  status: 'pending' | 'syncing' | 'failed' | 'synced'
  retry_count: number
  last_error?: string
  next_attempt_at: string // ISO timestamp
}

export interface SyncMetaRecord {
  key: string          // 'plan_last_fetch', 'messages_last_fetch_{partnerId}'
  value: string        // ISO timestamp or any string
}

export interface NotificationPrefRecord {
  id: number           // always 1 — singleton
  new_message: boolean
  new_plan: boolean
  food_request_result: boolean
  meal_reminders: boolean
  medication_reminders: boolean
  water_reminders: boolean
}

export interface UiStateRecord {
  key: string
  value: string
}

class NutriTrackDB extends Dexie {
  activePlan!: EntityTable<ActivePlanRecord, 'id'>
  messages!: EntityTable<MessageRecord, 'id'>
  syncQueue!: EntityTable<SyncQueueEntry, 'id'>
  syncMeta!: EntityTable<SyncMetaRecord, 'key'>
  notificationPreferences!: EntityTable<NotificationPrefRecord, 'id'>
  uiState!: EntityTable<UiStateRecord, 'key'>

  constructor() {
    super('nutritrack')
    this.version(1).stores({
      activePlan: 'id',
      messages: 'id, partner_id, cached_at',  // compound index for partner queries
      syncQueue: '++id, status, entity_type, next_attempt_at',
      syncMeta: 'key',
      notificationPreferences: 'id',
      uiState: 'key',
    })
  }
}

let _db: NutriTrackDB | null = null

export function useDb(): NutriTrackDB {
  if (!_db) {
    _db = new NutriTrackDB()
  }
  return _db
}
```

> **iOS eviction check (D-23):** On app open, call `useDb().activePlan.count()`. If IndexedDB throws `DatabaseClosedError` or returns 0 unexpectedly, set a `uiState` key `eviction_detected = 'true'` and show the recovery banner.

---

### Pattern 4: Offline-Aware Client API Layer (`composables/useClientApi.ts`)

**What:** Thin wrapper around `useApi().apiFetch` that intercepts POST errors due to network failure and routes writes to the Dexie sync queue instead. GET requests are never queued — they fall through to cached reads in individual stores. This is the integration seam (D-09).

```typescript
// [ASSUMED pattern — aligned with D-09, D-10, D-11]
export function useClientApi() {
  const { apiFetch } = useApi()
  const db = useDb()

  async function clientPost<T>(
    path: string,
    body: Record<string, unknown>,
    options?: { entityType: string; localId: string; attachmentBlob?: Blob; attachmentFilename?: string; attachmentMime?: string },
  ): Promise<T | { queued: true; local_id: string }> {
    if (!navigator.onLine) {
      // Immediately queue without attempting network
      return enqueue(path, body, options)
    }
    try {
      return await apiFetch<T>(path, { method: 'POST', body: JSON.stringify(body) })
    }
    catch (err) {
      // Network error (not 4xx/5xx — those are real failures, not offline)
      if (isNetworkError(err)) {
        return enqueue(path, body, options)
      }
      throw err
    }
  }

  async function enqueue(
    path: string,
    body: Record<string, unknown>,
    options?: { entityType: string; localId: string; attachmentBlob?: Blob; attachmentFilename?: string; attachmentMime?: string },
  ) {
    const localId = options?.localId ?? (body.local_id as string) ?? crypto.randomUUID()
    await db.syncQueue.add({
      entity_type: options?.entityType ?? 'unknown',
      request_path: path,
      method: 'POST',
      payload: JSON.stringify(body),
      attachment_blob: options?.attachmentBlob,
      attachment_filename: options?.attachmentFilename,
      attachment_mime: options?.attachmentMime,
      local_id: localId,
      created_at: new Date().toISOString(),
      status: 'pending',
      retry_count: 0,
      next_attempt_at: new Date().toISOString(),
    })
    return { queued: true as const, local_id: localId }
  }

  return { clientPost }
}

function isNetworkError(err: unknown): boolean {
  // fetch throws TypeError for network failures; createError sets statusCode
  return err instanceof TypeError || (err as { statusCode?: number }).statusCode === undefined
}
```

> **Attachment handling for messages (D-08):** Messages with `File` objects need special handling. Convert the `File` to a `Blob` before storing (already a Blob in JS). Store `attachment_blob`, `attachment_filename`, `attachment_mime` in the syncQueue entry. On sync, reconstruct the `FormData` from the stored Blob.

---

### Pattern 5: Sync Manager (`composables/useSyncManager.ts`)

**What:** Processes `syncQueue` FIFO, single-flight (one item at a time). Called on `online` event and app-open sweep. Registers Background Sync where available (D-13).

```typescript
// [ASSUMED pattern — aligned with D-11..D-14]
export function useSyncManager() {
  const db = useDb()
  const { apiFetch } = useApi()
  const isSyncing = ref(false)
  const pendingCount = ref(0)
  const failedCount = ref(0)

  async function refreshCounts() {
    pendingCount.value = await db.syncQueue.where('status').anyOf(['pending', 'syncing']).count()
    failedCount.value = await db.syncQueue.where('status').equals('failed').count()
  }

  async function processQueue() {
    if (isSyncing.value || !navigator.onLine) return
    isSyncing.value = true
    try {
      // Process FIFO: oldest pending first
      const items = await db.syncQueue
        .where('status').anyOf(['pending'])
        .and(item => item.next_attempt_at <= new Date().toISOString())
        .sortBy('created_at')
      
      for (const item of items) {
        await db.syncQueue.update(item.id!, { status: 'syncing' })
        try {
          let responseBody: unknown
          if (item.attachment_blob) {
            const fd = new FormData()
            const payload = JSON.parse(item.payload) as Record<string, string>
            Object.entries(payload).forEach(([k, v]) => fd.append(k, v))
            fd.append('attachment', item.attachment_blob, item.attachment_filename ?? 'file')
            responseBody = await apiFetch(item.request_path, { method: 'POST', body: fd })
          } else {
            responseBody = await apiFetch(item.request_path, {
              method: 'POST',
              body: item.payload,
            })
          }
          await db.syncQueue.update(item.id!, { status: 'synced' })
          // Emit to stores so they can update optimistic data
          useNuxtApp().callHook('sync:itemSynced', { entity_type: item.entity_type, local_id: item.local_id, response: responseBody })
        }
        catch (err) {
          const retries = item.retry_count + 1
          const backoff = Math.pow(2, Math.min(retries - 1, 2)) * 1000 // 1s, 2s, 4s
          const nextAttempt = new Date(Date.now() + backoff).toISOString()
          if (retries >= 3) {
            await db.syncQueue.update(item.id!, { status: 'failed', retry_count: retries, last_error: String(err) })
          } else {
            await db.syncQueue.update(item.id!, {
              status: 'pending', retry_count: retries,
              last_error: String(err), next_attempt_at: nextAttempt,
            })
          }
        }
      }
    } finally {
      isSyncing.value = false
      await refreshCounts()
    }
  }

  function registerBackgroundSync() {
    if ('serviceWorker' in navigator && 'SyncManager' in window) {
      navigator.serviceWorker.ready.then(reg => reg.sync.register('nutritrack-sync'))
        .catch(() => { /* Background Sync not available — foreground polling handles it */ })
    }
  }

  function setupOnlineListener() {
    window.addEventListener('online', () => processQueue())
    // Also receive message from SW Background Sync trigger
    navigator.serviceWorker?.addEventListener('message', (e) => {
      if (e.data?.type === 'TRIGGER_SYNC') processQueue()
    })
  }

  return { isSyncing, pendingCount, failedCount, processQueue, refreshCounts, registerBackgroundSync, setupOnlineListener }
}
```

---

### Pattern 6: Plan Cache Integration (extend `stores/clientPlan.ts`)

**What:** `fetchActivePlan` reads from Dexie first (instant), then fetches online and updates cache. Satisfies OFFL-02, OFFL-10.

```typescript
// Extend fetchActivePlan in clientPlan.ts
async function fetchActivePlan() {
  loading.value = true
  error.value = null
  const db = useDb()
  
  // 1. Read from IndexedDB cache first — instant render
  const cached = await db.activePlan.get(1)
  if (cached) {
    activePlan.value = cached.data as DietPlanResponse
    initActiveDay()
    loading.value = false  // Unblock UI immediately
  }

  // 2. Only fetch from API if online
  if (!navigator.onLine) return
  
  try {
    const { apiFetch } = useApi()
    const data = await apiFetch<DietPlanResponse>('/clients/me/active-plan')
    activePlan.value = data
    initActiveDay()
    // 3. Write back to cache
    await db.activePlan.put({ id: 1, plan_id: data.id, data, fetched_at: new Date().toISOString(), updated_hint: data.updated_at ?? '' })
    await db.syncMeta.put({ key: 'plan_last_fetch', value: new Date().toISOString() })
  } catch (e: unknown) {
    const err = e as { statusCode?: number; data?: { error?: string } }
    if (err.statusCode === 404) activePlan.value = null
    else if (!cached) error.value = (err.data?.error) ?? 'خطا در بارگذاری برنامه'
    // If we already have cached data, silently ignore online fetch failure
  } finally {
    loading.value = false
  }
}
```

---

### Pattern 7: Backend — Push Notifications Schema

**Migration 000010:**
```sql
-- push_subscriptions table
CREATE TABLE push_subscriptions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    endpoint     TEXT NOT NULL,
    p256dh       TEXT NOT NULL,
    auth_key     TEXT NOT NULL,
    user_agent   TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (client_id, endpoint)  -- idempotent by endpoint (D-16)
);

-- notification_preferences table
CREATE TABLE notification_preferences (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id             UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    new_message           BOOLEAN NOT NULL DEFAULT TRUE,
    new_plan              BOOLEAN NOT NULL DEFAULT TRUE,
    food_request_result   BOOLEAN NOT NULL DEFAULT TRUE,
    meal_reminders        BOOLEAN NOT NULL DEFAULT TRUE,
    medication_reminders  BOOLEAN NOT NULL DEFAULT TRUE,
    water_reminders       BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Reminder dedup: prevents sending same reminder twice per minute window
CREATE TABLE sent_reminders (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    dedup_key   TEXT NOT NULL,     -- e.g. 'meal:plan_id:meal_id:2026-04-20T08:00'
    sent_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (dedup_key)
);

CREATE INDEX idx_sent_reminders_sent_at ON sent_reminders(sent_at);
-- Cleanup job: delete sent_reminders older than 24h (or add TTL via pg_cron in Phase 7)
```

---

### Pattern 8: Backend — `notification_service.go`

```go
// [ASSUMED structure — aligned with D-17, D-18, D-19]
type PushPayload struct {
    Title     string `json:"title"`
    Body      string `json:"body"`
    ActionURL string `json:"action_url"`
    Icon      string `json:"icon"`
    Type      string `json:"type"`        // 'message' | 'new_plan' | 'food_request' | 'meal_reminder' | 'medication_reminder' | 'water_reminder'
    EntityID  string `json:"entity_id,omitempty"`
}

type NotificationService struct {
    repo       repository.NotificationRepository
    vapidPublic  string
    vapidPrivate string
    vapidContact string   // mailto:admin@nutritrack.ir
    logger     zerolog.Logger
}

func (s *NotificationService) SendToClient(ctx context.Context, clientID uuid.UUID, notifType string, payload PushPayload) error {
    // 1. Check preference for this type
    prefs, err := s.repo.GetPreferences(ctx, clientID)
    if err != nil || !prefEnabled(prefs, notifType) {
        return nil  // Silently skip if preference off
    }
    // 2. Get all subscriptions for this client
    subs, err := s.repo.GetSubscriptions(ctx, clientID)
    if err != nil || len(subs) == 0 {
        return nil
    }
    // 3. Send to all devices (best-effort; remove expired subscriptions on 410 Gone)
    jsonBytes, _ := json.Marshal(payload)
    for _, sub := range subs {
        resp, err := webpush.SendNotification(jsonBytes, &webpush.Subscription{
            Endpoint: sub.Endpoint,
            Keys: webpush.Keys{ P256dh: sub.P256dh, Auth: sub.AuthKey },
        }, &webpush.Options{
            VAPIDPublicKey:  s.vapidPublic,
            VAPIDPrivateKey: s.vapidPrivate,
            Subscriber:      s.vapidContact,
            TTL:             30,
        })
        if err != nil {
            s.logger.Warn().Err(err).Str("client_id", clientID.String()).Msg("push send failed")
            continue
        }
        defer resp.Body.Close()
        if resp.StatusCode == 410 {  // Subscription expired
            _ = s.repo.DeleteSubscription(ctx, sub.ID)
        }
    }
    return nil
}
```

---

### Pattern 9: Backend — Reminder Scheduler Goroutine

```go
// backend/internal/service/reminder_scheduler.go
// [ASSUMED pattern — aligned with D-18, D-19]
type ReminderScheduler struct {
    notifSvc *NotificationService
    planRepo repository.DietPlanRepository
    logger   zerolog.Logger
}

func (s *ReminderScheduler) Start(ctx context.Context) {
    ticker := time.NewTicker(60 * time.Second)
    defer ticker.Stop()
    s.logger.Info().Msg("reminder scheduler started")
    for {
        select {
        case <-ticker.C:
            s.runReminderCycle(ctx)
        case <-ctx.Done():
            s.logger.Info().Msg("reminder scheduler stopped")
            return
        }
    }
}

func (s *ReminderScheduler) runReminderCycle(ctx context.Context) {
    now := time.Now()
    windowKey := now.Format("2006-01-02T15:04") // minute-level dedup key

    // Query active plans with meals/medications firing in this minute window
    // (implementation queries plan_days.meals joined against current time's hour:minute)
    // For each: check dedup key, check preference, send push, insert dedup record
}

// main.go wiring (within main() after service initialization):
// scheduler := service.NewReminderScheduler(notifService, planRepo, logger)
// go scheduler.Start(ctx)  // ctx is already the signal-cancellable root context
```

---

### Pattern 10: Config Extension (`config.go`)

```go
// Add to Config struct:
VAPIDPublicKey  string  // base64url VAPID public key (required for push)
VAPIDPrivateKey string  // base64url VAPID private key (required for push)
VAPIDContact    string  // mailto:admin@nutritrack.ir (default)

// Add to Load():
VAPIDPublicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
VAPIDPrivateKey: os.Getenv("VAPID_PRIVATE_KEY"),
VAPIDContact:    getEnv("VAPID_CONTACT", "mailto:admin@nutritrack.ir"),

// Validation: if not empty, both must be present together
```

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Service worker precaching | Manual cache.addAll() list | `workbox-precaching` via `@vite-pwa/nuxt` injectManifest | Workbox handles cache versioning, cleanup of old caches, and URL revision hashing automatically. Manual lists miss chunks on every build. |
| Push payload encryption | Custom VAPID/ECDH encryption | `webpush-go` | Web Push uses ECDH + HKDF + AES-128-GCM. Getting it wrong silently fails. `webpush-go` is battle-tested. |
| IndexedDB raw API | `indexedDB.open()` boilerplate | `dexie` v4 | Raw IDB error handling (blocked upgrades, version conflicts, transaction aborts) is notoriously complex. Dexie handles all of it. |
| Background Sync queuing | Custom SW fetch queue | `workbox-background-sync` BackgroundSyncPlugin | Handles retry, storage, and lifecycle across SW restarts. |
| VAPID key generation | Custom crypto code | `webpush-go`'s CLI / `vapid-keygen` | One-time operation; wrong key format breaks all push. |
| Duplicate reminder prevention | In-memory set | PostgreSQL `sent_reminders` table with UNIQUE dedup key | In-memory set is lost on restart; PostgreSQL survives process restarts. |

---

## Common Pitfalls

### Pitfall 1: iOS Push Requires Home Screen Install + iOS 16.4+
**What goes wrong:** Push notifications silently fail on iOS without showing any permission dialog or error.
**Why it happens:** iOS Safari only supports Web Push for PWA installed to the home screen (not opened in browser). iOS 16.4+ is required. Push subscription call (`PushManager.subscribe`) returns a successful `PushSubscription` object but pushes are never delivered if the app was not launched from the home screen icon.
**How to avoid:** Show install prompt first; only show notification permission prompt after the user has added to home screen. Detect with `navigator.standalone` on iOS. Gracefully handle cases where push permission is denied or not available.
**Warning signs:** `PushSubscription` created successfully but no push events received during testing on iOS.
> [CITED: https://webkit.org/blog/13878/web-push-for-web-apps-on-ios-and-ipados/]

### Pitfall 2: `injectManifest` SW Must Be in `public/`, Not `app/`
**What goes wrong:** `nuxt build` does not process `app/` directory files as static assets. The SW file must be at a path Vite can pick up as a static file entry point, not as a Vue component.
**Why it happens:** vite-plugin-pwa with `injectManifest` uses Rollup to build `srcDir/filename` separately from the main bundle. Nuxt's `app/` directory is for Vue component resolution.
**How to avoid:** Set `pwa.srcDir: 'public'` and `pwa.filename: 'sw.ts'` in `nuxt.config.ts`. The file must be `frontend/public/sw.ts`. [ASSUMED — verify on first build]
**Warning signs:** `[vite-pwa] failed to find sw file` build error.

### Pitfall 3: Dexie Throws in SSR Context
**What goes wrong:** Nuxt 4 runs modules server-side (SSR). `IndexedDB` does not exist in Node.js. Calling `new Dexie()` during SSR throws `ReferenceError: indexedDB is not defined`.
**Why it happens:** `useDb()` is called in store setup code which runs both client and server.
**How to avoid:** Guard `useDb()` with `if (process.client)` or use `onMounted` for Dexie calls. In `plugins/sync.client.ts` (note `.client.ts` suffix — Nuxt auto-skips on server). Alternatively, the `useDb()` singleton initializes lazily, but only call it from composables that run client-side.
**Warning signs:** `indexedDB is not defined` in server console; SSR hydration mismatch.

### Pitfall 4: Background Sync API Has Limited Support
**What goes wrong:** Background Sync API (`SyncManager`) is only reliably supported in Chrome/Chromium on Android. Safari (iOS/macOS), Firefox, and Chrome on iOS do not support it as of 2026.
**Why it happens:** The spec has not been implemented broadly.
**How to avoid:** Background Sync is an enhancement (D-13 explicitly says foreground polling is primary). Always implement the foreground `online` event listener + app-open sweep. Background Sync adds value only on Android Chrome.
**Warning signs:** `SyncManager is not defined` in console on iOS/Firefox — this is expected, not a bug.
> [CITED: https://caniuse.com/background-sync — ~73% global support as of 2026]

### Pitfall 5: FormData + Blobs Cannot Be JSON.stringify'd in syncQueue
**What goes wrong:** `message.ts` sends `new FormData()` with attachment `File` objects. Files/Blobs are not serializable with `JSON.stringify`. Storing the FormData in a string `payload` column loses the attachment.
**Why it happens:** Offline message queue entries with attachments need the binary data preserved.
**How to avoid:** Store the Blob separately in `syncQueue.attachment_blob` (Dexie handles Blob storage natively in IndexedDB). Store only JSON-serializable fields in `payload`. Reconstruct `FormData` from stored Blob + payload fields when processing the queue. The `useClientApi.clientPost` interface must accept `attachmentBlob` separately.
**Warning signs:** Synced messages arriving at server with no attachment when they had one before going offline.

### Pitfall 6: SW `registerType: 'autoUpdate'` Skips Waiting — Test for Update Loops
**What goes wrong:** `autoUpdate` calls `skipWaiting()` immediately and reloads. If the new SW also has a build error, users get stuck in a reload loop.
**Why it happens:** `autoUpdate` combined with a faulty SW build creates an infinite activate → error → update → activate cycle.
**How to avoid:** Test SW in dev with `devOptions.enabled: true`. Build and smoke-test before deploying. Monitor SW registration errors in production via zerolog/Loki.
**Warning signs:** Chrome DevTools → Application → Service Workers shows repeated registration failures.

### Pitfall 7: iOS PWA Storage Eviction (D-23)
**What goes wrong:** iOS can delete all IDB data for a PWA that hasn't been used for ~7 days (low-storage scenario). The app opens, Dexie initializes with empty tables, and shows stale/empty UI with no explanation.
**Why it happens:** iOS uses an evictable storage bucket by default for PWA origins.
**How to avoid:** On every app open, check `db.activePlan.count()`. If 0 and `syncMeta.plan_last_fetch` was recent (within 1 hour), treat as eviction event. Display: `داده‌های آفلاین توسط دستگاه حذف شدند، پس از اتصال به اینترنت بازیابی می‌شوند`. Also call `navigator.storage.persist()` on first visit (iOS 16.4+ respects this).
**Warning signs:** Client reports "empty app" after device ran low on storage.

### Pitfall 8: VAPID Keys Must Not Change After First Deployment
**What goes wrong:** Changing `VAPID_PUBLIC_KEY` invalidates all existing push subscriptions. Clients will not receive push until they re-subscribe. There's no migration path for existing subscriptions.
**Why it happens:** The browser links the subscription cryptographically to the VAPID public key used at subscribe time.
**How to avoid:** Generate VAPID keys once before first production deployment. Store in a password manager and in the secrets manager. Never rotate unless absolutely necessary.
**Warning signs:** All push deliveries suddenly return 401 Unauthorized after a key rotation.

---

## Integration Points

### Frontend → Backend new endpoints

| Endpoint | Method | Purpose | Auth |
|----------|--------|---------|------|
| `/api/push/subscribe` | POST | Save push subscription (idempotent by endpoint) | client |
| `/api/push/unsubscribe` | DELETE | Remove push subscription by endpoint | client |
| `/api/push/preferences` | GET | Fetch notification preferences for current client | client |
| `/api/push/preferences` | PATCH | Update notification preferences | client |

These 4 endpoints are all the new backend surface area. Everything else is triggered internally.

### Backend service hooks (trigger points for event-driven push, D-18)

| Service Method | Push Type | Payload Target |
|----------------|-----------|----------------|
| `CommunicationService.SendMessageTo()` | `message` | receiver's client subscriptions |
| `DietPlanService.ActivatePlan()` | `new_plan` | the plan's assigned client |
| `CommunicationService.ApproveFoodRequest()` | `food_request` | requesting client |
| `CommunicationService.RejectFoodRequest()` | `food_request` | requesting client |

**Implementation pattern:** `NotificationService` is injected into `CommunicationService` and `DietPlanService` constructors. After a successful write, call `notifSvc.SendToClient(ctx, targetClientID, type, payload)`. This must be non-blocking (use `go` goroutine or fire-and-forget) so push delivery failures don't affect the main response.

### `main.go` wiring additions

```go
// After existing service initialization:
notifRepo := repository.NewNotificationRepository(pool)
notifService := service.NewNotificationService(notifRepo, cfg.VAPIDPublicKey, cfg.VAPIDPrivateKey, cfg.VAPIDContact, logger)
notifHandler := handler.NewNotificationHandler(notifService)

// Inject notifService into communication and plan services:
commService := service.NewCommunicationService(commRepo, userRepo, cfg.UploadsDir, notifService, logger)
planService := service.NewDietPlanService(planRepo, notifService, logger)

// Start reminder scheduler (uses the root ctx for graceful shutdown):
scheduler := service.NewReminderScheduler(notifService, planRepo, logger)
go scheduler.Start(ctx)

// Register push routes:
pushRoutes := r.Group("/api/push")
pushRoutes.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("client"))
{
    pushRoutes.POST("/subscribe", notifHandler.Subscribe)
    pushRoutes.DELETE("/unsubscribe", notifHandler.Unsubscribe)
    pushRoutes.GET("/preferences", notifHandler.GetPreferences)
    pushRoutes.PATCH("/preferences", notifHandler.UpdatePreferences)
}
```

### `client.vue` layout additions

```vue
<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <OfflineBanner v-if="!isOnline" />          <!-- new -->
    <slot />
    <SyncStatusBadge />                           <!-- new: sits above bottom nav -->
    <UiBottomNav :items="navItems" />
  </div>
</template>

<script setup lang="ts">
// ... existing navItems
const isOnline = useOnline()   // Vue composable or window.navigator.onLine reactive ref
const { processQueue, registerBackgroundSync, setupOnlineListener } = useSyncManager()

onMounted(() => {
  setupOnlineListener()
  registerBackgroundSync()
  processQueue()               // Sweep on every app open
})
</script>
```

---

## Phased Execution Recommendations

Execute in 6 waves. Each wave is independently testable and deployable.

### Wave 1: PWA Shell + Manifest + Service Worker Skeleton (2 plans)
**Goal:** App installs on Android Chrome + iOS Safari. Static assets cached. No offline data yet.
- Install `@vite-pwa/nuxt`, configure `nuxt.config.ts` with manifest + icons
- Create `public/sw.ts` skeleton (precache only, no runtime routes yet)
- Add PWA icons (at minimum 192×192 and 512×512 PNG)
- Verify: install prompt appears on Chrome Android; app launches standalone
- **Risk:** Nuxt 4 + vite-plugin-pwa `injectManifest` compatibility — test build first

### Wave 2: Dexie Schema + Plan Cache (2 plans)
**Goal:** Active plan loads from cache when offline. One plan snapshot stored in IndexedDB.
- Install `dexie`, create `useDb.ts` with schema version 1
- Extend `clientPlan.ts` `fetchActivePlan()` with cache-read-first pattern
- Add iOS eviction detection + `OfflineBanner`
- Add runtime cache rule to `sw.ts` for `/api/clients/me/active-plan`
- Verify: plan visible after airplane mode (Chrome DevTools)

### Wave 3: Sync Queue + Tracking Offline (3 plans)
**Goal:** All 6 tracking types work offline. Sync queue processes on reconnect.
- Create `useClientApi.ts` offline-aware wrapper
- Create `useSyncManager.ts` with FIFO processing + exponential backoff
- Wrap `foodLog`, `waterLog`, `sleepLog`, `exerciseLog`, `medicationLog`, `bodyMeasurement` store writes with `useClientApi.clientPost`
- Add `SyncStatusBadge` component, wire into `client.vue` layout
- Register Background Sync in `sw.ts`
- Verify: log food offline → see queue → go online → verify server received entry (check `local_id` dedup)
- **This is the highest complexity wave — allocate most time here**

### Wave 4: Offline Messaging (2 plans)
**Goal:** Last 50 messages cached. Messages sendable offline (text + attachments).
- Extend `message.ts` `fetchMessages()` with Dexie `messages` table cache
- Wrap `message.ts` `sendMessage()` with offline queue (incl. Blob storage for attachments)
- Add cached message merge logic on poll
- Verify: load messages → airplane mode → messages still visible; send offline → reconnect → message delivered

### Wave 5: Backend Push Notifications (3 plans)
**Goal:** Event-driven push for messages, plans, food requests. Reminder scheduler running.
- DB migration 000010: `push_subscriptions`, `notification_preferences`, `sent_reminders`
- `notification_repo.go`, `notification_service.go`, `reminder_scheduler.go`
- `notification_handler.go` + route wiring in `main.go`
- Inject `NotificationService` into `CommunicationService` and `DietPlanService`
- Configure VAPID keys in env
- Verify: send message → push received on subscribed device

### Wave 6: Frontend Push + Preferences + Final Polish (2 plans)
**Goal:** Push subscription prompt on first login. Preference toggles in settings.
- Add push event handler to `sw.ts`
- Create `useNotificationPermission.ts` composable + first-login prompt (after plan loads, not on startup)
- Create `notificationPrefs.ts` Pinia store
- Create `pages/client/settings/notifications.vue` preference page
- iOS: test install-first-then-push flow on real device
- Verify: disable meal reminders → no reminder fires; re-enable → resumes

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `workbox-webpack-plugin` | `vite-plugin-pwa` + `@vite-pwa/nuxt` | Webpack → Vite era (~2022) | Simpler config, Vite-native, auto-integrates with Nuxt |
| `idb` raw wrapper | `Dexie.js` v4 with TypeScript EntityTable | Dexie v4 (2024) | TypeScript generics mean no manual type casting on queries |
| Manual service worker for push | `BackgroundSyncPlugin` in Workbox | Workbox 6+ | Plugin handles retry and SW lifecycle automatically |
| `localForage` | `Dexie.js` | Ongoing | Dexie has better schema versioning and query API for complex offline models |

**Deprecated/outdated:**
- `@nuxtjs/pwa` (Nuxt 2 era): replaced by `@vite-pwa/nuxt` for Nuxt 3/4
- `dexie-observable` / `dexie-syncable`: paid/complex cloud sync addons — not needed here

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Vitest 3.x (already in package.json) |
| Config file | vitest.config.ts (or package.json scripts) |
| Quick run command | `npm run test` (maps to `vitest run`) |
| Full suite command | `npm run test` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| OFFL-05 | Dexie schema initializes with correct tables | unit | `vitest run tests/useDb.test.ts` | ❌ Wave 1 |
| OFFL-06 | useSyncManager processes queue FIFO | unit | `vitest run tests/useSyncManager.test.ts` | ❌ Wave 3 |
| OFFL-07 | local_id dedup: same local_id queued twice yields single item | unit | `vitest run tests/useSyncManager.test.ts` | ❌ Wave 3 |
| OFFL-08 | Exponential backoff: retry_count 0/1/2 → next_attempt at 1s/2s/4s | unit | `vitest run tests/useSyncManager.test.ts` | ❌ Wave 3 |
| OFFL-12 | iOS eviction: empty activePlan after recent fetch → eviction banner shown | unit | `vitest run tests/useDb.test.ts` | ❌ Wave 2 |
| NOTIF-07 | Preferences: disabling type blocks send in notification service | unit (Go) | `go test ./internal/service/... -run TestNotificationService` | ❌ Wave 5 |
| OFFL-02, OFFL-10 | Plan cache: stale cache rendered while fresh fetch pending | unit | `vitest run tests/clientPlan.test.ts` | ❌ Wave 2 |

### Wave 0 Gaps
- [ ] `frontend/tests/useDb.test.ts` — covers OFFL-05, OFFL-12 schema init and eviction detection
- [ ] `frontend/tests/useSyncManager.test.ts` — covers OFFL-06, OFFL-07, OFFL-08 queue logic
- [ ] `frontend/tests/clientPlan.offline.test.ts` — covers OFFL-02, OFFL-10 cache-first reads
- [ ] `backend/internal/service/notification_service_test.go` — covers NOTIF-07 preference gating
- [ ] `vitest.config.ts` — ensure `happy-dom` environment is configured (already in devDependencies)

**Note:** Service Worker behavior (precaching, push events, Background Sync) requires end-to-end / manual testing in real browser. Vitest cannot simulate `ServiceWorkerGlobalScope`. Use Chrome DevTools → Application → Service Workers for manual verification.

---

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | yes | VAPID subscription requires authenticated client JWT; push subscription endpoints use existing `middleware.Auth` |
| V3 Session Management | no | No new session concerns — existing JWT + cookie pattern unchanged |
| V4 Access Control | yes | Push subscription endpoints guarded with `RoleGuard("client")` only; no nutritionist access |
| V5 Input Validation | yes | Push subscription body validated (endpoint, p256dh, auth as non-empty strings); preferences body validated with go-playground/validator |
| V6 Cryptography | yes | VAPID keys never exposed to frontend; only public key used client-side; private key stays in env |

### Known Threat Patterns

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Push subscription endpoint spoofing (client A subscribes for client B) | Spoofing | `client_id` always derived from authenticated JWT, never from request body |
| VAPID private key exposure | Information Disclosure | Store only in env var (`VAPID_PRIVATE_KEY`); never log; never include in config struct serialization |
| Sync queue replay (replay the same local_id to duplicate tracking entry) | Tampering | Backend `ON CONFLICT (local_id) DO NOTHING` — already in Phase 4; Phase 6 relies on it |
| Offline message attachment size bypass | Tampering | Client-side blob size check in `useClientApi`; server-side MIME + size validation in `CommunicationService.SendMessageTo()` still enforced on sync |
| Reminder spam (scheduler sends same reminder 60× per hour) | DoS | `sent_reminders` table with UNIQUE `dedup_key` prevents re-send within the same minute window |

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Node.js | npm install, Nuxt build | ✓ | In path (Nuxt already runs) | — |
| Go 1.25 | `go get webpush-go` | ✓ | 1.25.0 per go.mod | — |
| Chrome (Android or DevTools) | SW + push testing | ✓ | Any modern | — |
| iOS device (16.4+) | iOS PWA push testing | ✗ (unknown) | — | iOS emulator does NOT support push; real device required |
| VAPID keys | Push delivery | ✗ (not generated yet) | — | Generate once before Wave 5 |
| `@vite-pwa/nuxt` npm package | PWA | ✗ (not installed) | 1.1.1 | — |
| `dexie` npm package | IndexedDB | ✗ (not installed) | 4.4.2 | — |
| `webpush-go` Go module | Push backend | ✗ (not in go.mod) | v1.4.0 | — |

**Missing dependencies with no fallback:**
- iOS 16.4+ physical device for final push notification validation (Chrome DevTools cannot simulate iOS push)

**Missing dependencies with fallback:**
- VAPID keys: generate before Wave 5; can defer push feature until keys are generated

---

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | `public/sw.ts` is the correct `srcDir`/`filename` for `injectManifest` with Nuxt 4 | Pattern 2 | Build fails; SW not generated. Fix: check vite-plugin-pwa docs for Nuxt 4 path convention |
| A2 | `useNuxtApp().callHook()` is available inside a plain composable to emit sync events to stores | Pattern 5 | Store update post-sync breaks; fix: use a shared Pinia `syncEvent` store or Vue `provide/inject` |
| A3 | Dexie v4 stores `Blob` objects natively in IndexedDB without serialization | Pattern 4 (attachment blobs) | Attachments silently dropped; test with a real attachment in Wave 4 |
| A4 | Nuxt 4 server-side rendering does not execute `useDb()` if it's guarded with `process.client` | Pattern 3 | SSR crash `indexedDB is not defined`; fix: use `.client.ts` plugin suffix |
| A5 | `workbox-background-sync` BackgroundSyncPlugin can be used inside a custom `injectManifest` SW with Workbox 7 | Pattern 2 | Background Sync plugin import fails; fix: check vite-plugin-pwa Workbox version bundled with 1.1.1 |
| A6 | Reminder scheduler goroutine has access to active plans via `DietPlanRepository.ListActivePlansWithSchedule()` — this query does not yet exist | Pattern 9 | Scheduler cannot query meals; fix: add dedicated query in Wave 5 plan |
| A7 | `navigator.standalone` reliably detects iOS PWA home-screen launch | iOS push guidance | Push subscription created in browser context instead; fix: always allow subscription attempt but show install banner |

---

## Open Questions

1. **Workbox version bundled with `@vite-pwa/nuxt@1.1.1`**
   - What we know: `vite-plugin-pwa@1.2.0` is the underlying package; it bundles Workbox 7
   - What's unclear: Exact Workbox minor version affects BackgroundSyncPlugin API
   - Recommendation: Run `npm ls workbox-core` after install to confirm version; assume Workbox 7

2. **SW scope in Nuxt SSR output**
   - What we know: Nuxt 4 outputs to `.output/public/` for static assets
   - What's unclear: Whether the SW registered at `/sw.js` can control requests to `/api/*` when the API is on the same origin
   - Recommendation: In production, Traefik routes `/api/*` to Go backend and `/*` to Nuxt. The SW must be scoped to the Nuxt origin. The fetch handler for API routes in the SW only applies to the network-first cache — actual routing still goes through Traefik. Verify SW registration scope on first build.

3. **Lab result offline upload (OFFL-03 includes lab results)**
   - What we know: `trackingHandler.UploadLabResult` uses `multipart/form-data` with file binary
   - What's unclear: D-10 includes "lab-result metadata" — does this mean only the metadata (title, type, link) or also the file?
   - Recommendation: Queue only the metadata + external link variant offline. File binary uploads (up to 10MB) are too large for IndexedDB blob storage in most browsers. Show "Lab result upload requires internet connection" for file-based uploads when offline. This is consistent with the requirement saying "metadata/file uploads" — treat file as deferred.

---

## Sources

### Primary (HIGH confidence)
- [VERIFIED: npm registry] `@vite-pwa/nuxt@1.1.1` — confirmed current version
- [VERIFIED: npm registry] `dexie@4.4.2` — confirmed current version
- [VERIFIED: golang proxy] `github.com/SherClockHolmes/webpush-go@v1.4.0` — confirmed latest release (2025-01-02)
- [VERIFIED: codebase] `frontend/package.json` — confirmed `@vite-pwa/nuxt` and `dexie` not yet installed
- [VERIFIED: codebase] `backend/go.mod` — confirmed `webpush-go` not yet in module
- [VERIFIED: codebase] `frontend/app/composables/useApi.ts` — integration seam confirmed, uses native `fetch` with `credentials: 'include'`
- [VERIFIED: codebase] `frontend/app/stores/foodLog.ts` — `crypto.randomUUID()` for `local_id` confirmed in all tracking stores
- [VERIFIED: codebase] `backend/internal/config/config.go` — VAPID keys not yet in Config struct
- [VERIFIED: codebase] `backend/cmd/api/main.go` — service wiring pattern confirmed; scheduler goroutine pattern is new

### Secondary (MEDIUM confidence)
- [CITED: https://webkit.org/blog/13878/web-push-for-web-apps-on-ios-and-ipados/] — iOS 16.4+ Web Push for PWA requirement
- [CITED: https://caniuse.com/background-sync] — Background Sync API support ~73%, no iOS/Firefox

### Tertiary (LOW confidence — see Assumptions Log)
- [ASSUMED] Custom SW at `public/sw.ts` is correct path for Nuxt 4 injectManifest
- [ASSUMED] `useNuxtApp().callHook()` available inside composable for sync event emission
- [ASSUMED] Dexie v4 Blob storage works without serialization in all target browsers

---

## Metadata

**Confidence breakdown:**
- Standard stack (libraries + versions): HIGH — verified via npm registry and golang proxy
- Architecture patterns (Dexie schema, SW config, sync queue): HIGH — derived from CONTEXT.md decisions and verified existing code
- Backend push integration: HIGH — existing service patterns are clear; webpush-go API is documented
- iOS push specifics: MEDIUM — based on WebKit blog post; real device testing required
- SW injectManifest path conventions in Nuxt 4: MEDIUM — ASSUMED, verify on first build

**Research date:** 2026-04-20
**Valid until:** 2026-05-20 (30 days; `@vite-pwa/nuxt` minor versions move fast — re-verify before coding Wave 1)
