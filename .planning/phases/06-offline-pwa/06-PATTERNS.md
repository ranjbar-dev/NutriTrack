# Phase 6: Offline & PWA — Pattern Map

**Mapped:** 2026-04-20
**Files analyzed:** 21 new/modified files
**Analogs found:** 19 / 21

---

## File Classification

| New / Modified File | Role | Data Flow | Closest Analog | Match Quality |
|---|---|---|---|---|
| `frontend/nuxt.config.ts` | config | — | itself | exact (extend) |
| `frontend/app/service-worker/sw.ts` | service-worker | event-driven | _(none)_ | no analog |
| `frontend/app/composables/usePWA.client.ts` | composable/hook | event-driven | `useMessagePolling.ts` | role-match |
| `frontend/app/db/index.ts` | utility | file-I/O | _(none)_ | no analog |
| `frontend/app/composables/useOfflineApi.ts` | composable/utility | request-response | `useApi.ts` | exact |
| `frontend/app/stores/syncQueue.ts` | store | CRUD + event-driven | `message.ts` store | role-match |
| `frontend/app/composables/useSyncManager.ts` | composable/hook | event-driven | `useMessagePolling.ts` | exact |
| `frontend/app/stores/message.ts` _(extend)_ | store | CRUD + streaming | itself | exact (extend) |
| `frontend/app/stores/clientPlan.ts` _(extend)_ | store | CRUD | itself | exact (extend) |
| `frontend/app/stores/foodLog.ts` _(extend)_ | store | CRUD | itself | exact (extend) |
| `frontend/app/stores/waterLog.ts` _(extend)_ | store | CRUD | itself | exact (extend) |
| `frontend/app/stores/medicationLog.ts` _(extend)_ | store | CRUD | itself | exact (extend) |
| `frontend/app/stores/sleepLog.ts` _(extend)_ | store | CRUD | `foodLog.ts` | exact |
| `frontend/app/stores/exerciseLog.ts` _(extend)_ | store | CRUD | `foodLog.ts` | exact |
| `frontend/app/stores/bodyMeasurement.ts` _(extend)_ | store | CRUD | itself | exact (extend) |
| `frontend/app/stores/labResult.ts` _(extend)_ | store | file-I/O | itself | exact (extend) |
| `frontend/app/stores/notificationPreferences.ts` | store | CRUD | `auth.ts` store | role-match |
| `frontend/app/layouts/client.vue` _(extend)_ | layout | — | itself | exact (extend) |
| `frontend/app/components/ClientSyncStatus.vue` | component | — | `client.vue` layout | partial |
| `backend/db/migrations/000010_create_push.up.sql` | migration | — | `000009_create_communication.up.sql` | exact |
| `backend/db/queries/push_subscriptions.sql` | query | CRUD | `messages.sql` | exact |
| `backend/db/queries/notification_preferences.sql` | query | CRUD | `messages.sql` | exact |
| `backend/internal/repository/push_repo.go` | repository | CRUD | `communication_repo.go` | exact |
| `backend/internal/service/push_service.go` | service | event-driven | `communication_service.go` | role-match |
| `backend/internal/service/reminder_scheduler.go` | service | event-driven | _(none — first goroutine ticker)_ | no analog |
| `backend/internal/handler/push_handler.go` | handler | request-response | `communication_handler.go` | exact |
| `backend/internal/model/dto/push_dto.go` | model | — | existing dto files | exact |
| `backend/internal/config/config.go` _(extend)_ | config | — | itself | exact (extend) |
| `backend/cmd/api/main.go` _(extend)_ | bootstrap | — | itself | exact (extend) |

---

## Pattern Assignments

---

### `frontend/nuxt.config.ts` (config, extend)

**Analog:** itself (`frontend/nuxt.config.ts`)

**Current modules block** (line 26):
```ts
modules: ['@pinia/nuxt', '@nuxt/eslint'],
```

**Extension pattern — add `@vite-pwa/nuxt` as third module and a `pwa` key at root level:**
```ts
modules: ['@pinia/nuxt', '@nuxt/eslint', '@vite-pwa/nuxt'],

pwa: {
  // D-01: injectManifest strategy — explicit service worker, not zero-config preset
  strategies: 'injectManifest',
  srcDir: 'app/service-worker',
  filename: 'sw.ts',
  registerType: 'autoUpdate',

  manifest: {
    // D-04: Persian-only manifest
    name: 'نوتری‌ترک',
    short_name: 'نوتری‌ترک',
    lang: 'fa',
    dir: 'rtl',
    display: 'standalone',
    background_color: '#f9fafb',
    theme_color: '#16a34a',
    shortcuts: [
      { name: 'برنامه', url: '/client/plan' },
      { name: 'ثبت روزانه', url: '/client/tracking' },
      { name: 'پیام‌ها', url: '/client/messages' },
    ],
  },

  workbox: {
    // D-03: injectManifest hands control to sw.ts for all caching strategies
    globPatterns: ['**/*.{js,css,woff2,png,svg,ico}'],
  },
},
```

**Constraint:** `@vite-pwa/nuxt` and `vite-plugin-pwa` must be added to `frontend/package.json` dependencies.

---

### `frontend/app/service-worker/sw.ts` (service-worker, event-driven)

**Analog:** _(no existing service worker — no analog in repo)_

**Pattern basis:** `@vite-pwa/nuxt` injectManifest skeleton + Workbox runtime caching.

**Skeleton to implement:**
```ts
import { cleanupOutdatedCaches, precacheAndRoute } from 'workbox-precaching'
import { registerRoute } from 'workbox-routing'
import { NetworkFirst, CacheFirst } from 'workbox-strategies'

declare let self: ServiceWorkerGlobalScope

// D-03: precache static assets injected by vite-plugin-pwa
precacheAndRoute(self.__WB_MANIFEST)
cleanupOutdatedCaches()

// D-03: network-first for client API reads (plan, messages)
registerRoute(
  ({ url }) => url.pathname.startsWith('/api/clients/me/active-plan') ||
               url.pathname.startsWith('/api/messages'),
  new NetworkFirst({ cacheName: 'api-client-reads', networkTimeoutSeconds: 5 }),
)

// D-13: Background Sync registration hook
self.addEventListener('sync', (event: SyncEvent) => {
  if (event.tag === 'nutritrack-sync') {
    event.waitUntil(triggerForegroundSync())
  }
})

// D-18: Push notification click — navigate to action_url
self.addEventListener('notificationclick', (event) => {
  event.notification.close()
  const data = event.notification.data
  if (data?.action_url) {
    event.waitUntil(clients.openWindow(data.action_url))
  }
})
```

**Constraint:** `triggerForegroundSync()` must be a thin postMessage bridge to `useSyncManager` — the actual queue processing stays in the Vue layer to keep auth-cookie access inside the browser context.

---

### `frontend/app/composables/usePWA.client.ts` (composable, event-driven)

**Analog:** `frontend/app/composables/useMessagePolling.ts` (lines 1–31)

**Lifecycle pattern to copy from `useMessagePolling.ts`:**
```ts
// useMessagePolling.ts — onMounted/onUnmounted lifecycle binding
onMounted(start)
onUnmounted(stop)
```

**Module-level singleton pattern to copy from `useApi.ts` (lines 11–13):**
```ts
// useApi.ts — singleton state outside the composable function
let isRefreshing = false
let refreshQueue: Array<...> = []
```

**Implement as:**
```ts
// Module-level singleton — shared across components
let deferredPrompt: BeforeInstallPromptEvent | null = null
let registration: ServiceWorkerRegistration | null = null

export function usePWA() {
  const canInstall = ref(false)
  const needsUpdate = ref(false)

  function start() {
    window.addEventListener('beforeinstallprompt', (e) => {
      e.preventDefault()
      deferredPrompt = e as BeforeInstallPromptEvent
      canInstall.value = true
    })
    navigator.serviceWorker?.ready.then((reg) => {
      registration = reg
      reg.addEventListener('updatefound', () => { needsUpdate.value = true })
    })
  }

  async function promptInstall() {
    if (!deferredPrompt) return
    await deferredPrompt.prompt()
    deferredPrompt = null
    canInstall.value = false
  }

  function applyUpdate() {
    registration?.waiting?.postMessage({ type: 'SKIP_WAITING' })
    window.location.reload()
  }

  onMounted(start)

  return { canInstall, needsUpdate, promptInstall, applyUpdate }
}
```

**Constraint:** File name must end in `.client.ts` — this composable uses `window`/`navigator.serviceWorker` which are browser-only. Nuxt 4 auto-excludes `.client.ts` from SSR.

---

### `frontend/app/db/index.ts` (utility, file-I/O / IndexedDB)

**Analog:** _(no existing IndexedDB abstraction in repo — no analog)_

**D-05 schema to implement:**
```ts
import Dexie, { type Table } from 'dexie'

export interface CachedPlan {
  id: 1 // singleton row
  plan_id: string
  fetched_at: string
  updated_hint: string | null
  data: object // full DietPlanResponse JSON
}

export interface CachedMessage {
  id: string      // server message UUID
  partner_id: string
  sent_at: string
  payload: object // full MessageResponse JSON
}

export interface SyncQueueEntry {
  // D-11 fields
  id?: number     // autoIncrement
  entity_type: string
  request_path: string
  method: string
  payload: string           // JSON string
  attachment_blob?: Blob    // D-08: optional for file messages/lab uploads
  local_id: string
  created_at: string
  status: 'pending' | 'processing' | 'failed'
  retry_count: number
  last_error: string | null
  next_attempt_at: string
}

export interface NotificationPreference {
  id: 1 // singleton row
  new_message: boolean
  new_plan: boolean
  food_request_result: boolean
  meal_reminder: boolean
  medication_reminder: boolean
  water_reminder: boolean
}

class NutriTrackDB extends Dexie {
  activePlan!: Table<CachedPlan>
  messages!: Table<CachedMessage>
  syncQueue!: Table<SyncQueueEntry>
  notificationPreferences!: Table<NotificationPreference>
  uiState!: Table<{ key: string; value: unknown }>

  constructor() {
    super('NutriTrackDB')
    this.version(1).stores({
      activePlan: 'id',
      messages: 'id, partner_id, sent_at',
      syncQueue: '++id, status, created_at, entity_type',
      notificationPreferences: 'id',
      uiState: 'key',
    })
  }
}

export const db = new NutriTrackDB()
```

**Constraint:** Import `db` only in client-side composables/stores. Never import in pages/components directly — always go through a store or composable so iOS eviction (D-23) is handled in one place.

---

### `frontend/app/composables/useOfflineApi.ts` (composable, request-response)

**Analog:** `frontend/app/composables/useApi.ts` (full file, lines 1–102)

**Core seam to copy and extend — `apiFetch` signature and error escalation from `useApi.ts` (lines 26–59):**
```ts
// useApi.ts — the pattern to wrap
async function apiFetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${config.public.apiBase}${endpoint}`
  // ... sets Content-Type, credentials: 'include' ...
  const response = await fetch(url, { ...options, credentials: 'include', headers })
  if (response.status === 401 && !endpoint.includes('/auth/refresh')) {
    return handleRefreshAndRetry<T>(endpoint, options)
  }
  if (!response.ok) {
    const body = await response.json().catch(() => ({}))
    throw createError({ statusCode: response.status, message: body.error || 'خطایی رخ داد' })
  }
  if (response.status === 204) return {} as T
  return response.json()
}
```

**D-09 extension pattern:**
```ts
// useOfflineApi.ts — offline-aware layer on top of useApi
const QUEUEABLE_PATHS: RegExp[] = [
  /^\/client\/(food|water|sleep|exercise|medication)-logs$/,
  /^\/client\/body-measurements$/,
  /^\/client\/lab-results$/,
  /^\/messages$/,
]

export function useOfflineApi() {
  const { apiFetch } = useApi()
  const syncStore = useSyncQueueStore()

  async function clientApiFetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const isWrite = options.method === 'POST' || options.method === 'PUT' || options.method === 'PATCH'
    const isQueueable = isWrite && QUEUEABLE_PATHS.some(r => r.test(endpoint))

    if (!isQueueable) {
      return apiFetch<T>(endpoint, options) // fall through for non-queueable or reads
    }

    if (!navigator.onLine) {
      // D-09: queue immediately without attempting network
      await syncStore.enqueue(endpoint, options)
      return buildOptimisticResponse(endpoint, options) as T
    }

    try {
      return await apiFetch<T>(endpoint, options)
    } catch (err: any) {
      const isNetworkFailure = err?.statusCode == null // no status = transport error
      if (isNetworkFailure) {
        await syncStore.enqueue(endpoint, options)
        return buildOptimisticResponse(endpoint, options) as T
      }
      throw err // server-side error (4xx/5xx) — bubble up normally
    }
  }

  return { clientApiFetch }
}
```

**Constraint:** `apiFetch` from `useApi` already owns auth-refresh and cookie handling — never bypass it. `useOfflineApi` must call `apiFetch` underneath, not raw `fetch`.

---

### `frontend/app/stores/syncQueue.ts` (store, CRUD + event-driven)

**Analog:** `frontend/app/stores/message.ts` (full file, lines 1–80) + `frontend/app/stores/foodLog.ts` (lines 1–75)

**Store structure to copy from `message.ts` (lines 3–8):**
```ts
// message.ts — setup store pattern
export const useMessageStore = defineStore('message', () => {
  const messages = ref<Message[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
```

**`local_id` generation pattern from `foodLog.ts` (lines 28–30):**
```ts
// foodLog.ts — crypto.randomUUID() for local_id
const payload: LogFoodPayload = {
  local_id: crypto.randomUUID(),
  date: new Date().toISOString().slice(0, 10),
```

**D-11/D-12 queue store to implement:**
```ts
export const useSyncQueueStore = defineStore('syncQueue', () => {
  const pendingItems = ref<SyncQueueEntry[]>([])
  const isProcessing = ref(false)
  const failedItems = computed(() => pendingItems.value.filter(i => i.status === 'failed'))
  const pendingCount = computed(() => pendingItems.value.filter(i => i.status === 'pending').length)

  async function enqueue(path: string, options: RequestInit) {
    // Persist to Dexie, push to pendingItems ref
  }

  // D-12: exponential backoff 1s → 2s → 4s, max 3 retries
  async function processQueue() {
    if (isProcessing.value) return // single-flight guard
    isProcessing.value = true
    try {
      const items = pendingItems.value
        .filter(i => i.status === 'pending' && new Date(i.next_attempt_at) <= new Date())
        .sort((a, b) => a.created_at.localeCompare(b.created_at)) // FIFO

      for (const item of items) {
        await processItem(item)
      }
    } finally {
      isProcessing.value = false
    }
  }

  return { pendingItems, isProcessing, failedItems, pendingCount, enqueue, processQueue }
})
```

**Constraint:** After a successful sync item replay (D-14), call the appropriate in-memory store update directly (e.g., `useFoodLogStore().upsertLocal(serverResponse)`) — do not re-fetch the full list.

---

### `frontend/app/composables/useSyncManager.ts` (composable, event-driven)

**Analog:** `frontend/app/composables/useMessagePolling.ts` (full file, lines 1–31)

**Entire pattern to copy verbatim and adapt:**
```ts
// useMessagePolling.ts — interval + lifecycle pattern
export function useMessagePolling(partnerId: string) {
  let intervalId: ReturnType<typeof setInterval> | null = null
  let lastSeen = new Date().toISOString()

  function start() {
    if (intervalId) return
    intervalId = setInterval(async () => {
      const newMsgs = await messageStore.pollNewMessages(partnerId, lastSeen)
      if (newMsgs.length > 0) lastSeen = newMsgs[newMsgs.length - 1].sent_at
    }, 10_000)
  }

  function stop() {
    if (intervalId) { clearInterval(intervalId); intervalId = null }
  }

  onMounted(start)
  onUnmounted(stop)

  return { start, stop }
}
```

**D-13 adaptation:**
```ts
export function useSyncManager() {
  const syncStore = useSyncQueueStore()
  let intervalId: ReturnType<typeof setInterval> | null = null

  function start() {
    if (intervalId) return
    // Foreground timer fallback (D-13) — runs even without Background Sync support
    intervalId = setInterval(() => {
      if (navigator.onLine) syncStore.processQueue()
    }, 30_000)

    // Background Sync registration (D-13) — best-effort, not required
    navigator.serviceWorker?.ready.then((reg) => {
      if ('sync' in reg) reg.sync.register('nutritrack-sync').catch(() => {})
    })
  }

  window.addEventListener('online', () => syncStore.processQueue()) // reconnect sweep

  function stop() {
    if (intervalId) { clearInterval(intervalId); intervalId = null }
  }

  onMounted(start)
  onUnmounted(stop)

  return { start, stop }
}
```

---

### `frontend/app/stores/message.ts` _(extend)_ (store, CRUD + streaming)

**Analog:** itself (`frontend/app/stores/message.ts`, full file, lines 1–80)

**`fetchMessages` to extend (lines 9–21) — add Dexie read-first, then merge by ID:**
```ts
// Current pattern to keep as-is:
async function fetchMessages(partnerId: string, limit = 50, offset = 0) {
  const { apiFetch } = useApi()
  loading.value = true
  error.value = null
  try {
    const data = await apiFetch<Message[]>(`/messages/${partnerId}?...`)
    messages.value = data
  } catch (e: any) {
    error.value = e?.data?.error ?? 'خطا در دریافت پیام‌ها'
  } finally {
    loading.value = false
  }
}
```

**D-07 extension pattern:**
```ts
async function fetchMessages(partnerId: string, limit = 50, offset = 0) {
  // 1. Read cache immediately (D-07)
  const cached = await db.messages.where('partner_id').equals(partnerId)
    .sortBy('sent_at')
  if (cached.length > 0) messages.value = cached.map(c => c.payload as Message)

  loading.value = true
  if (!navigator.onLine) { loading.value = false; return }

  try {
    const { apiFetch } = useApi()
    const data = await apiFetch<Message[]>(`/messages/${partnerId}?limit=${limit}&offset=${offset}`)
    // Merge by ID (D-07)
    const byId = new Map(messages.value.map(m => [m.id, m]))
    data.forEach(m => byId.set(m.id, m))
    messages.value = [...byId.values()].sort((a, b) => a.sent_at.localeCompare(b.sent_at))
    // Persist last 50 to Dexie
    await persistMessages(partnerId, messages.value.slice(-50))
  } catch (e: any) {
    if (messages.value.length === 0)
      error.value = 'خطا در دریافت پیام‌ها'
  } finally {
    loading.value = false
  }
}
```

**`sendMessage` to extend (lines 38–53) — add optimistic echo + queue fallback (D-08):**
```ts
// Replace direct apiFetch POST with useOfflineApi().clientApiFetch
// Add optimistic echo before the call:
const localEcho: Message = {
  id: `local_${crypto.randomUUID()}`, // temp local ID, replaced on sync success
  sender_id: authStore.user!.id,
  receiver_id: receiverId,
  content: content ?? null,
  sent_at: new Date().toISOString(),
  read_at: null,
  attachment_type: null, attachment_path: null, attachment_name: null,
}
messages.value.push(localEcho)
```

---

### `frontend/app/stores/clientPlan.ts` _(extend)_ (store, CRUD)

**Analog:** itself (`frontend/app/stores/clientPlan.ts`, full file, lines 1–123)

**`fetchActivePlan` to extend (lines 12–33) — add Dexie read-first (D-06):**
```ts
// Current pattern (lines 12–33) — preserve error.value = null, loading pattern exactly
async function fetchActivePlan() {
  loading.value = true
  error.value = null
  try {
    const { apiFetch } = useApi()
    const data = await apiFetch<DietPlanResponse>('/clients/me/active-plan')
    activePlan.value = data
    initActiveDay()
  } catch (e: unknown) {
    const err = e as { statusCode?: number; data?: { error?: string } }
    if (err.statusCode === 404) activePlan.value = null
    else error.value = (err.data?.error) ?? 'خطا در بارگذاری برنامه'
  } finally {
    loading.value = false
  }
}
```

**D-06 extension — insert cache read before network call:**
```ts
async function fetchActivePlan() {
  // 1. Try Dexie first (D-06)
  const cached = await db.activePlan.get(1)
  if (cached) {
    activePlan.value = cached.data as DietPlanResponse
    initActiveDay()
  }

  if (!navigator.onLine) {
    if (!cached) error.value = 'داده‌ای در حافظه موجود نیست. پس از اتصال دوباره تلاش کنید.'
    loading.value = false
    return
  }

  loading.value = true
  error.value = null
  try {
    const { apiFetch } = useApi()
    const data = await apiFetch<DietPlanResponse>('/clients/me/active-plan')
    activePlan.value = data
    initActiveDay()
    // Persist snapshot (D-06)
    await db.activePlan.put({ id: 1, plan_id: data.id, fetched_at: new Date().toISOString(), updated_hint: null, data })
  } catch (e: unknown) {
    // existing error handling unchanged
  } finally {
    loading.value = false
  }
}
```

---

### Tracking Stores: `foodLog.ts`, `waterLog.ts`, `medicationLog.ts`, `sleepLog.ts`, `exerciseLog.ts`, `bodyMeasurement.ts`, `labResult.ts` _(extend)_

**Analog per store:** itself (each already has `local_id` + POST via `apiFetch`)

**Canonical `local_id` + POST pattern from `foodLog.ts` (lines 26–41):**
```ts
// foodLog.ts lines 26–41 — this is THE pattern all tracking write methods follow
async function logFood(mealId: string, selectedOptionId: string, notes?: string) {
  const { apiFetch } = useApi()
  const payload: LogFoodPayload = {
    local_id: crypto.randomUUID(),   // ← preserve exactly; Phase 6 queues this payload
    date: new Date().toISOString().slice(0, 10),
    meal_id: mealId,
    selected_option_id: selectedOptionId,
    is_skipped: false,
    ...(notes ? { notes } : {}),
  }
  const entry = await apiFetch<FoodLogEntry>('/client/food-logs', {
    method: 'POST',
    body: JSON.stringify(payload),
  })
  upsertLocal(entry)  // ← after sync replay, call this same function with server response
}
```

**D-09/D-10 extension — swap `apiFetch` call site for `clientApiFetch`:**
```ts
// Replace in every tracking write method:
// const { apiFetch } = useApi()
// const entry = await apiFetch<FoodLogEntry>('/client/food-logs', { method: 'POST', body: ... })

// With:
const { clientApiFetch } = useOfflineApi()
const entry = await clientApiFetch<FoodLogEntry>('/client/food-logs', {
  method: 'POST',
  body: JSON.stringify(payload),
})
// upsertLocal(entry) call is unchanged — works with both real and optimistic responses
```

**`waterLog.ts` pattern (lines 29–42) — same shape, copy for remaining stores:**
```ts
async function addWater(amountMl: number, loggedTime?: string) {
  const { apiFetch } = useApi()
  const payload: LogWaterPayload = {
    local_id: crypto.randomUUID(),
    date: new Date().toISOString().slice(0, 10),
    amount_ml: amountMl,
    ...(loggedTime ? { logged_time: loggedTime } : {}),
  }
  const entry = await apiFetch<WaterLogEntry>('/client/water-logs', { method: 'POST', body: JSON.stringify(payload) })
  logs.value.push(entry)
}
```

**`labResult.ts` special case (lines 26–31) — FormData + file blob for D-08:**
```ts
// labResult.ts — FormData POST without JSON.stringify
async function uploadLabResult(formData: FormData) {
  const { apiFetch } = useApi()
  const result = await apiFetch<LabResultResponse>('/client/lab-results', { method: 'POST', body: formData })
  labResults.value.unshift(result)
}
```
When offline, the attachment Blob must be extracted from FormData and stored in `SyncQueueEntry.attachment_blob` (D-08). `useOfflineApi` must handle this special case: detect `body instanceof FormData`, extract the file field, store the blob, then reassemble on replay.

---

### `frontend/app/stores/notificationPreferences.ts` (store, CRUD)

**Analog:** `frontend/app/stores/auth.ts` (lines 11–97) for profile-data store shape; `frontend/app/stores/message.ts` for fetch/update pattern.

**Store skeleton — copy `auth.ts` setup store structure (lines 11–16):**
```ts
// auth.ts — profile store template
export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => !!user.value)
  // ...
  async function checkAuth() {
    try {
      const { apiFetch } = useApi()
      const data = await apiFetch<{ user: User }>('/auth/me')
      user.value = data.user
    } catch { user.value = null }
  }
```

**Implement as:**
```ts
export const useNotificationPrefsStore = defineStore('notificationPrefs', () => {
  const prefs = ref<NotificationPrefsResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const pushSubscribed = ref(false)

  async function fetchPrefs() {
    const { apiFetch } = useApi()
    loading.value = true
    try {
      prefs.value = await apiFetch<NotificationPrefsResponse>('/client/notification-preferences')
      // Sync to Dexie for offline display
      await db.notificationPreferences.put({ id: 1, ...prefs.value })
    } catch (e: any) {
      error.value = e?.data?.error ?? 'خطا در دریافت تنظیمات اعلان'
    } finally { loading.value = false }
  }

  async function updatePrefs(patch: Partial<NotificationPrefsResponse>) {
    const { apiFetch } = useApi()
    await apiFetch('/client/notification-preferences', { method: 'PATCH', body: JSON.stringify(patch) })
    prefs.value = { ...prefs.value!, ...patch }
  }

  async function subscribePush(subscription: PushSubscriptionJSON) {
    const { apiFetch } = useApi()
    await apiFetch('/client/push-subscriptions', { method: 'POST', body: JSON.stringify(subscription) })
    pushSubscribed.value = true
  }

  async function unsubscribePush(endpoint: string) {
    const { apiFetch } = useApi()
    await apiFetch('/client/push-subscriptions', { method: 'DELETE', body: JSON.stringify({ endpoint }) })
    pushSubscribed.value = false
  }

  return { prefs, loading, error, pushSubscribed, fetchPrefs, updatePrefs, subscribePush, unsubscribePush }
})
```

---

### `frontend/app/layouts/client.vue` _(extend)_ (layout)

**Analog:** itself (`frontend/app/layouts/client.vue`, lines 1–15)

**Current layout (full file):**
```vue
<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <slot />
    <UiBottomNav :items="navItems" />
  </div>
</template>
```

**D-20 extension — add `<ClientSyncStatus>` and PWA prompts inside the wrapper div:**
```vue
<template>
  <div class="min-h-screen bg-gray-50 pb-20">
    <ClientSyncStatus />  <!-- D-20: visible from all client screens -->
    <slot />
    <UiBottomNav :items="navItems" />
    <ClientPwaInstallBanner v-if="canInstall" @install="promptInstall" />
    <ClientPwaUpdateBanner v-if="needsUpdate" @update="applyUpdate" />
  </div>
</template>

<script setup lang="ts">
const { canInstall, needsUpdate, promptInstall, applyUpdate } = usePWA()
// navItems unchanged
</script>
```

**Constraint:** Only this layout is modified (D-02). `nutritionist.vue` and `admin.vue` layouts must remain untouched.

---

### `frontend/app/components/ClientSyncStatus.vue` (component)

**Analog:** `frontend/app/layouts/client.vue` (lines 1–15) for mobile-first, Persian-only, card-based style convention.

**D-20 Persian copy strings:**
- `همگام‌سازی در حال انجام` — when `isProcessing`
- `همه‌چیز همگام است` — when `pendingCount === 0`
- `X مورد در انتظار` — when `pendingCount > 0`
- `تلاش مجدد` — manual retry CTA for failed items

**Component pattern:**
```vue
<script setup lang="ts">
const syncStore = useSyncQueueStore()
const { pendingCount, isProcessing, failedItems } = storeToRefs(syncStore)
</script>
<template>
  <!-- compact pill, shown only when pendingCount > 0 or isProcessing -->
  <div v-if="isProcessing || pendingCount > 0" class="...">
    <span v-if="isProcessing">همگام‌سازی در حال انجام</span>
    <span v-else>{{ pendingCount }} مورد در انتظار</span>
  </div>
  <div v-if="failedItems.length > 0" class="...">
    <button @click="syncStore.retryFailed()">تلاش مجدد</button>
  </div>
</template>
```

---

## Backend Pattern Assignments

---

### `backend/db/migrations/000010_create_push.up.sql` (migration)

**Analog:** `backend/db/migrations/000009_create_communication.up.sql` (full file)

**Exact structural pattern to copy:**
```sql
-- 000009: CREATE TYPE + CREATE TABLE + FK to users(id) + indexes
CREATE TYPE food_request_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE food_requests (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requested_by UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_food_requests_requested_by ON food_requests (requested_by, created_at DESC);
```

**D-15 tables to create:**
```sql
-- Migration 000010: Web Push subscriptions and notification preferences

CREATE TABLE push_subscriptions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    endpoint   TEXT NOT NULL,
    p256dh     TEXT NOT NULL,
    auth_key   TEXT NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_push_subscriptions_endpoint UNIQUE (endpoint)  -- D-16: idempotent by endpoint
);
CREATE INDEX idx_push_subscriptions_client ON push_subscriptions (client_id);

CREATE TABLE notification_preferences (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id            UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    new_message          BOOLEAN NOT NULL DEFAULT TRUE,
    new_plan             BOOLEAN NOT NULL DEFAULT TRUE,
    food_request_result  BOOLEAN NOT NULL DEFAULT TRUE,
    meal_reminder        BOOLEAN NOT NULL DEFAULT TRUE,
    medication_reminder  BOOLEAN NOT NULL DEFAULT TRUE,
    water_reminder       BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_notification_preferences_client UNIQUE (client_id)
);
```

**Constraint:** Migration file number must follow current sequence (next after `000009`). Create matching `.down.sql` with `DROP TABLE` + `DROP INDEX` in reverse order.

---

### `backend/db/queries/push_subscriptions.sql` (query)

**Analog:** `backend/db/queries/messages.sql` (full file)

**SQLC comment format to copy exactly:**
```sql
-- messages.sql
-- name: CreateMessage :one
INSERT INTO messages (sender_id, ...) VALUES ($1, ...) RETURNING ...;

-- name: ListMessages :many
SELECT ... FROM messages WHERE ... LIMIT $3 OFFSET $4;
```

**Queries to implement:**
```sql
-- name: UpsertPushSubscription :one
INSERT INTO push_subscriptions (client_id, endpoint, p256dh, auth_key, user_agent)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (endpoint) DO UPDATE
  SET client_id = EXCLUDED.client_id, p256dh = EXCLUDED.p256dh,
      auth_key = EXCLUDED.auth_key, user_agent = EXCLUDED.user_agent
RETURNING *;

-- name: DeletePushSubscriptionByEndpoint :exec
DELETE FROM push_subscriptions WHERE endpoint = $1 AND client_id = $2;

-- name: ListPushSubscriptionsByClient :many
SELECT * FROM push_subscriptions WHERE client_id = $1;

-- name: ListAllActiveSubscriptions :many
-- Used by reminder scheduler to fan out to all active client subscriptions
SELECT ps.* FROM push_subscriptions ps
JOIN users u ON u.id = ps.client_id
WHERE u.status = 'active';
```

---

### `backend/db/queries/notification_preferences.sql` (query)

**Analog:** `backend/db/queries/messages.sql` — same SQLC comment format.

```sql
-- name: GetNotificationPreferences :one
SELECT * FROM notification_preferences WHERE client_id = $1;

-- name: UpsertNotificationPreferences :one
INSERT INTO notification_preferences (client_id)
VALUES ($1)
ON CONFLICT (client_id) DO NOTHING
RETURNING *;

-- name: UpdateNotificationPreferences :one
UPDATE notification_preferences
SET new_message = $2, new_plan = $3, food_request_result = $4,
    meal_reminder = $5, medication_reminder = $6, water_reminder = $7,
    updated_at = NOW()
WHERE client_id = $1
RETURNING *;
```

---

### `backend/internal/repository/push_repo.go` (repository, CRUD)

**Analog:** `backend/internal/repository/communication_repo.go` (full file)

**Interface + struct + constructor pattern to copy (lines 14–41):**
```go
// communication_repo.go — canonical repo pattern
type CommunicationRepository interface {
    CreateMessage(ctx context.Context, ...) (*dto.MessageResponse, error)
    // ...
}

type communicationRepository struct {
    pool *pgxpool.Pool
    q    *sqlc.Queries
}

func NewCommunicationRepository(pool *pgxpool.Pool) CommunicationRepository {
    return &communicationRepository{pool: pool, q: sqlc.New(pool)}
}
```

**pgtype.UUID conversion to copy (lines 43–56):**
```go
// communication_repo.go lines 43–56 — pgtype.UUID wrapping pattern
func (r *communicationRepository) CreateMessage(...) (*dto.MessageResponse, error) {
    msg, err := r.q.CreateMessage(ctx, sqlc.CreateMessageParams{
        SenderID:   pgtype.UUID{Bytes: senderID, Valid: true},
        ReceiverID: pgtype.UUID{Bytes: receiverID, Valid: true},
        Content:    textOrNull(content),
    })
    if err != nil { return nil, err }
    return messageToDTO(msg), nil
}
```

**Implement as:**
```go
type PushRepository interface {
    UpsertSubscription(ctx context.Context, clientID uuid.UUID, endpoint, p256dh, authKey string, userAgent *string) (*dto.PushSubscriptionResponse, error)
    DeleteSubscription(ctx context.Context, clientID uuid.UUID, endpoint string) error
    ListSubscriptionsByClient(ctx context.Context, clientID uuid.UUID) ([]dto.PushSubscriptionResponse, error)
    ListAllActiveSubscriptions(ctx context.Context) ([]dto.PushSubscriptionResponse, error)
    GetNotificationPreferences(ctx context.Context, clientID uuid.UUID) (*dto.NotificationPrefsResponse, error)
    UpsertNotificationPreferences(ctx context.Context, clientID uuid.UUID) error
    UpdateNotificationPreferences(ctx context.Context, clientID uuid.UUID, prefs dto.UpdateNotificationPrefsRequest) (*dto.NotificationPrefsResponse, error)
}
```

---

### `backend/internal/service/push_service.go` (service, event-driven)

**Analog:** `backend/internal/service/communication_service.go` (full file, lines 1–263)

**Package-level sentinel errors to copy (lines 21–28):**
```go
// communication_service.go — Persian sentinel errors
var (
    ErrCommNotFound            = errors.New("رکورد یافت نشد")
    ErrCommUnauthorized        = errors.New("دسترسی غیرمجاز")
    ErrMsgAttachmentTooLarge   = errors.New("حجم فایل بیش از حد مجاز است")
)
```

**Struct + constructor pattern (lines 30–39):**
```go
// communication_service.go — struct with repo + config + logger
type CommunicationService struct {
    repo       repository.CommunicationRepository
    userRepo   repository.UserRepository
    uploadsDir string
    logger     zerolog.Logger
}

func NewCommunicationService(repo repository.CommunicationRepository, userRepo repository.UserRepository, uploadsDir string, logger zerolog.Logger) *CommunicationService {
    return &CommunicationService{repo: repo, userRepo: userRepo, uploadsDir: uploadsDir, logger: logger}
}
```

**D-17 service struct:**
```go
var (
    ErrPushSubscriptionNotFound = errors.New("اشتراک پوش یافت نشد")
    ErrPushSendFailed           = errors.New("ارسال اعلان ناموفق بود")
)

type PushService struct {
    repo      repository.PushRepository
    userRepo  repository.UserRepository
    vapidSub  string  // "mailto:admin@example.com"
    vapidPub  string
    vapidPriv string
    logger    zerolog.Logger
}

func NewPushService(repo repository.PushRepository, userRepo repository.UserRepository, vapidPub, vapidPriv, vapidSub string, logger zerolog.Logger) *PushService {
    return &PushService{repo: repo, userRepo: userRepo, vapidPub: vapidPub, vapidPriv: vapidPriv, vapidSub: vapidSub, logger: logger}
}

// D-17: shared push envelope
type PushPayload struct {
    Title    string `json:"title"`
    Body     string `json:"body"`
    ActionURL string `json:"action_url"`
    Icon     string `json:"icon"`
    Type     string `json:"type"`
    EntityID string `json:"entity_id,omitempty"`
}

func (s *PushService) SendToClient(ctx context.Context, clientID uuid.UUID, payload PushPayload) error {
    subs, err := s.repo.ListSubscriptionsByClient(ctx, clientID)
    if err != nil { return err }
    for _, sub := range subs {
        // webpush-go call — errors logged but not fatal (stale sub = delete)
        if err := s.send(sub, payload); err != nil {
            s.logger.Warn().Err(err).Str("endpoint", sub.Endpoint).Msg("push send failed")
        }
    }
    return nil
}
```

**Constraint:** `webpush-go` must be added to `backend/go.mod`. VAPID keys are loaded via `config.go` env vars, not hardcoded.

---

### `backend/internal/service/reminder_scheduler.go` (service, event-driven)

**Analog:** _(no goroutine ticker exists in this repo — no analog)_

**`main.go` goroutine pattern to copy (lines 311–316) — only existing goroutine in the codebase:**
```go
// main.go lines 311–316 — goroutine + channel for graceful shutdown
go func() {
    logger.Info().Str("addr", srv.Addr).Msg("HTTP server listening")
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        logger.Fatal().Err(err).Msg("HTTP server error")
    }
}()
<-ctx.Done()
```

**D-18/D-19 scheduler pattern:**
```go
type ReminderScheduler struct {
    pushSvc  *PushService
    planRepo repository.DietPlanRepository
    logger   zerolog.Logger
    dedup    sync.Map // key: "clientID:type:minuteBucket" → struct{}{}
}

func NewReminderScheduler(pushSvc *PushService, planRepo repository.DietPlanRepository, logger zerolog.Logger) *ReminderScheduler {
    return &ReminderScheduler{pushSvc: pushSvc, planRepo: planRepo, logger: logger}
}

// Start launches the ticker goroutine — called from main.go before <-ctx.Done()
func (s *ReminderScheduler) Start(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Minute)
    go func() {
        defer ticker.Stop()
        for {
            select {
            case t := <-ticker.C:
                s.tick(ctx, t)
            case <-ctx.Done():
                return
            }
        }
    }()
}

func (s *ReminderScheduler) tick(ctx context.Context, now time.Time) {
    bucket := now.Format("2006-01-02T15:04") // minute-level dedup key (D-19)
    // Query active plans, check meal/medication/water times, send push if not deduped
}
```

**Constraint:** `s.dedup` is in-memory only (per D-19). It resets on server restart — acceptable because the scheduler only targets the current minute bucket.

---

### `backend/internal/handler/push_handler.go` (handler, request-response)

**Analog:** `backend/internal/handler/communication_handler.go` (full file, lines 1–316)

**Handler struct + constructor to copy (lines 17–25):**
```go
// communication_handler.go — struct + constructor
type CommunicationHandler struct {
    commService *service.CommunicationService
    uploadsDir  string
}

func NewCommunicationHandler(commService *service.CommunicationService, uploadsDir string) *CommunicationHandler {
    return &CommunicationHandler{commService: commService, uploadsDir: uploadsDir}
}
```

**`authUUID` + `ShouldBindJSON` + error dispatch pattern (lines 31–77):**
```go
// communication_handler.go — canonical client endpoint pattern
func (h *CommunicationHandler) SendMessage(c *gin.Context) {
    senderID, ok := authUUID(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
        return
    }
    var req dto.FoodRequestCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "اطلاعات ورودی نامعتبر است"})
        return
    }
    result, err := h.commService.DoSomething(c.Request.Context(), senderID, req)
    if err != nil {
        h.handleCommError(c, err)
        return
    }
    c.JSON(http.StatusCreated, result)
}
```

**Error switch pattern (lines 299–316):**
```go
// communication_handler.go — error type switch
func (h *CommunicationHandler) handleCommError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, service.ErrCommNotFound):
        c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
    case errors.Is(err, service.ErrCommUnauthorized):
        c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
    default:
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "خطای داخلی سرور"})
    }
}
```

**Push handler endpoints (D-16):**
```go
type PushHandler struct {
    pushService *service.PushService
}

// POST /api/client/push-subscriptions   — subscribe (idempotent by endpoint)
// DELETE /api/client/push-subscriptions — unsubscribe
// GET  /api/client/notification-preferences
// PATCH /api/client/notification-preferences
```

---

### `backend/internal/config/config.go` _(extend)_ (config)

**Analog:** itself (`backend/internal/config/config.go`, full file, lines 1–67)

**Field + env-read pattern to copy (lines 9–31):**
```go
// config.go — existing pattern
type Config struct {
    Port        string
    DatabaseURL string
    JWTSecret   string
    Environment string
    FrontendURL string
    UploadsDir  string
    SMSAPIKey   string
    SMSTemplate string
}

func Load() (*Config, error) {
    cfg := &Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: os.Getenv("DATABASE_URL"),  // required fields use os.Getenv
        UploadsDir:  getEnv("UPLOADS_DIR", "./uploads"),  // optional fields use getEnv with default
    }
```

**D-17 extension — add three VAPID fields:**
```go
type Config struct {
    // ... existing fields unchanged ...
    VAPIDPublicKey  string // VAPID_PUBLIC_KEY (required for push)
    VAPIDPrivateKey string // VAPID_PRIVATE_KEY (required for push)
    VAPIDEmail      string // VAPID_EMAIL e.g. "mailto:admin@nutritrack.ir"
}

// In Load():
cfg.VAPIDPublicKey  = os.Getenv("VAPID_PUBLIC_KEY")
cfg.VAPIDPrivateKey = os.Getenv("VAPID_PRIVATE_KEY")
cfg.VAPIDEmail      = getEnv("VAPID_EMAIL", "mailto:admin@nutritrack.ir")

// In validate(): warn but don't fatal if VAPID keys are empty in development
if cfg.Environment == "production" && (cfg.VAPIDPublicKey == "" || cfg.VAPIDPrivateKey == "") {
    return errors.New("VAPID_PUBLIC_KEY and VAPID_PRIVATE_KEY are required in production")
}
```

---

### `backend/cmd/api/main.go` _(extend)_ (bootstrap)

**Analog:** itself (`backend/cmd/api/main.go`, full file, lines 1–407)

**Repo → service → handler wiring pattern to copy (lines 99–127):**
```go
// main.go — canonical wiring pattern
commRepo := repository.NewCommunicationRepository(pool)
commService := service.NewCommunicationService(commRepo, userRepo, cfg.UploadsDir, logger)
commHandler := handler.NewCommunicationHandler(commService, cfg.UploadsDir)
```

**Route registration pattern to copy (lines 264–286):**
```go
// main.go — client group registration
client := r.Group("/api/client")
client.Use(middleware.Auth(jwtSecret), middleware.RoleGuard("client"))
{
    client.POST("/food-logs", trackingHandler.LogFood)
    // ...
}
```

**D-18 extension — add push wiring and scheduler startup:**
```go
// After existing repo/service init:
pushRepo := repository.NewPushRepository(pool)
pushService := service.NewPushService(pushRepo, userRepo, cfg.VAPIDPublicKey, cfg.VAPIDPrivateKey, cfg.VAPIDEmail, logger)
pushHandler := handler.NewPushHandler(pushService)

// Start reminder scheduler BEFORE blocking on <-ctx.Done()
scheduler := service.NewReminderScheduler(pushService, planRepo, logger)
scheduler.Start(ctx)

// Register client push routes inside existing 'client' group:
client.POST("/push-subscriptions", pushHandler.Subscribe)
client.DELETE("/push-subscriptions", pushHandler.Unsubscribe)
client.GET("/notification-preferences", pushHandler.GetPreferences)
client.PATCH("/notification-preferences", pushHandler.UpdatePreferences)
```

**Constraint:** `planRepo` already exists as `planRepo := repository.NewDietPlanRepository(pool)` (line 105). Pass it to `NewReminderScheduler` — do not create a second instance.

---

## Shared Patterns

### Authentication — `authUUID` helper in handlers

**Source:** `backend/internal/handler/communication_handler.go` lines 32–37
**Apply to:** `push_handler.go` — all four push endpoints must call `authUUID(c)` first.

```go
// Every client handler endpoint starts with:
clientID, ok := authUUID(c)
if !ok {
    c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "احراز هویت الزامی است"})
    return
}
```

`authUUID` is a package-private helper already defined in the `handler` package — do not redefine it.

---

### Frontend auth threading — `useApi()` must not be bypassed

**Source:** `frontend/app/composables/useApi.ts` lines 23–101
**Apply to:** `useOfflineApi.ts`, `syncQueue.ts` (replay path), `notificationPreferences.ts`

All network calls must go through `apiFetch` from `useApi()`, which owns token refresh (lines 61–98) and `credentials: 'include'` (line 39). The offline layer wraps this; it never calls raw `fetch` for API requests.

---

### Persian error strings — frontend

**Source:** `frontend/app/stores/foodLog.ts` line 18, `waterLog.ts` line 22, `message.ts` line 16

```ts
// All stores use this pattern:
error.value = err.data?.error ?? 'خطا در بارگذاری ...'
// For offline states, use:
error.value = 'داده‌ای در حافظه موجود نیست. پس از اتصال دوباره تلاش کنید.'
// For queued offline writes:
// Show toast, not error.value, because the operation succeeded locally
```

---

### Persian error strings — backend

**Source:** `backend/internal/service/communication_service.go` lines 21–28, `tracking_service.go` lines 22–28

All service-layer sentinel errors are Persian strings. New push service errors must follow the same style:
```go
var (
    ErrPushSubscriptionNotFound = errors.New("اشتراک پوش یافت نشد")
    ErrPushDeliveryFailed       = errors.New("ارسال اعلان ناموفق بود")
    ErrPushPrefsNotFound        = errors.New("تنظیمات اعلان یافت نشد")
)
```

---

### Store `$reset` convention

**Source:** `frontend/app/stores/foodLog.ts` lines 68–72, `clientPlan.ts` lines 101–107

Every Pinia setup store exposes `$reset()`. New stores (`syncQueue.ts`, `notificationPreferences.ts`) must include:
```ts
function $reset() {
  // reset all refs to initial values
  // clear Dexie tables owned by this store if applicable
}
```

---

### Migration file naming convention

**Source:** `backend/db/migrations/` directory listing

Format: `NNNNNN_<snake_case_description>.<up|down>.sql`
Next migration must be `000010_create_push_notifications.up.sql` / `.down.sql`.
Run order is enforced by `golang-migrate` sequence number.

---

### SQLC query generation

**Source:** `backend/db/queries/messages.sql` (full file), `backend/internal/repository/sqlc/` (generated)

After adding new `.sql` query files, run `sqlc generate` from `backend/` to regenerate `internal/repository/sqlc/`. The generated types are what the repository implementations use — never write raw SQL in Go service files.

---

## No Analog Found

| File | Role | Data Flow | Reason |
|------|------|-----------|--------|
| `frontend/app/service-worker/sw.ts` | service-worker | event-driven | First service worker in project; no existing SW file |
| `frontend/app/db/index.ts` | utility | IndexedDB / file-I/O | No client-side DB abstraction exists; Dexie is a new dependency |
| `backend/internal/service/reminder_scheduler.go` | service | event-driven | First goroutine ticker in the codebase; no background job pattern exists |

For these three files, use the patterns from `RESEARCH.md` + the `@vite-pwa/nuxt` / Dexie.js / `webpush-go` documentation respectively. The closest structural reference for the scheduler goroutine is the HTTP server goroutine in `main.go` lines 311–316.

---

## Metadata

**Analog search scope:** `frontend/app/stores/`, `frontend/app/composables/`, `frontend/app/layouts/`, `backend/internal/service/`, `backend/internal/handler/`, `backend/internal/repository/`, `backend/db/migrations/`, `backend/db/queries/`, `backend/internal/config/`, `backend/cmd/api/`
**Files read:** 24
**Pattern extraction date:** 2026-04-20
