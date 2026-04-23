# Phase 4 Research: Messaging, Notifications & Lab Exchange

**Generated:** 2026-04-23
**Phase:** 04-messaging-notifications-lab-exchange
**Sources:** docs/API.md §16–17, §20 · docs/PRD.md §5.10–5.11, §6, §7 · existing composable/store patterns

---

## 1. Messaging Strategy

### Polling Approach (D-01)

PRD §5.11 specifies polling every **10 seconds** while the chat screen is **open**. No global background poll.

```typescript
// Pattern: onMounted → setInterval, onUnmounted → clearInterval
// Refresh call: getClientConversation() or getNutritionistConversation(clientId)
// Unread badge: GET /messages/unread-count — on nav arrival, NOT per poll tick
```

**API contracts (from docs/API.md §17):**

| Action | Client | Nutritionist |
|--------|--------|--------------|
| Read conversation | `GET /messages?page=1&page_size=20` | `GET /clients/:id/messages?page=1&page_size=20` |
| Send message | `POST /messages` (multipart) | `POST /clients/:id/messages` (multipart) |
| Unread count | `GET /messages/unread-count` | `GET /messages/unread-count` |

Response is paginated, **newest-first**. UI reverses to render oldest-first (chronological thread order). D-02 prohibits client-side re-sorting beyond this reversal.

### Attachment Constraints (D-03)

| Attachment type | Accepted formats | Max size |
|-----------------|-----------------|---------|
| Image | JPG, PNG | 5 MB |
| File | PDF | 10 MB |

Validate `file.size` and `file.type` client-side before FormData submission. Show inline Persian error copy; do not submit. Backend enforces the same limits — frontend validation is UX guard only.

### Offline Message Queue (PRD §6 offline matrix)

PRD specifies "Send messages → Queue for sync". Text-only messages can be queued offline. File attachments cannot be queued (D-05 rationale: file binary data is unsuitable for IndexedDB queue; consistent with lab-upload online-only decision).

**Required offline-sync type extension:**

```typescript
// Extend TrackingDomain union in app/types/offline-sync.ts:
export type TrackingDomain = 'food' | 'water' | 'sleep' | 'exercise' | 'medication' | 'body' | 'message'

// New payload type:
export interface MessageQueuePayload {
  content: string  // text-only; file_data not queued
}
```

Extend `client-offline` store:
- Add `'message'` case to `isSupportedDomainPayload()` — valid if `content` is a non-empty string
- Add `cachedMessages ref<Message[]>` (max 50 entries per PRD §6 data freshness)
- On reconnect, replay orchestration calls `POST /messages` for each queued 'message' domain entry

Extend `plugins/client-sync.client.ts` to handle `domain === 'message'` → call `useMessagingApi().sendClientMessage({ content: payload.content })`.

### Message Object Shape (from docs/API.md §17)

```typescript
interface MessageAttachment {
  url: string
  type: string   // MIME type, e.g. "application/pdf"
  size: number   // bytes
  name: string
}

interface Message {
  id: string
  sender_id: string
  receiver_id: string
  content: string | null       // null when attachment-only
  is_mine: boolean             // true when caller is sender
  read_at: string | null
  attachment: MessageAttachment | null
  created_at: string
}
```

---

## 2. Lab Results Strategy

### File vs Link Branching (D-06)

`GET /lab-results/:id/download` returns:
- **File-backed**: file as attachment (`Content-Disposition: attachment`) — use `window.open(downloadUrl)` or `<a href download>`
- **Link-backed**: `302` redirect to external URL — use `window.open(url, '_blank')`

Determine branch from `LabResult` object: `result.link !== null` → link type; `result.original_name !== null` → file type.

### Online-Only Upload (D-05)

Do NOT enqueue lab uploads. On upload attempt while offline: show `ConnectivityBanner` or `InlineNotice` with Persian copy. Disable the upload CTA when `platform-pwa.offline === true`.

### Upload State Surface (D-04)

States: `'idle' | 'uploading' | 'success' | 'failure'`

- `idle`: upload form enabled
- `uploading`: spinner, CTA disabled, no cancel (simple v1 UX)
- `success`: Persian success notice, list refreshes
- `failure`: Persian recovery copy + retry button

Use native `$fetch` with `FormData` body (NOT `useFetch`) for multipart upload — `useFetch` does not support streaming progress and adding FormData requires manual header management.

```typescript
// Multipart pattern:
const form = new FormData()
form.append('title', req.title)
form.append('result_type', req.result_type)
if (req.file) form.append('file', req.file)
if (req.link) form.append('link', req.link)
const result = await $fetch<LabResult>(`/api/v1/clients/${clientId}/lab-results`, {
  method: 'POST',
  body: form,
})
```

### Lab Result Object Shape (from docs/API.md §16)

```typescript
interface LabResult {
  id: string
  client_id: string
  nutritionist_id: string
  title: string
  result_type: string            // e.g. "blood_sugar", "lipid_panel"
  test_date: string | null       // "YYYY-MM-DD"
  original_name: string | null   // null for link-type
  file_type: string | null       // MIME, null for link-type
  file_size: number              // 0 for link-type
  link: string | null            // null for file-type
  notes: string
  created_at: string
}
```

---

## 3. Push Subscription Strategy (D-08)

### Web Push API Pattern

```typescript
// Support check:
const isSupported = 'PushManager' in window && 'serviceWorker' in navigator

// Subscribe:
const registration = await navigator.serviceWorker.ready
const subscription = await registration.pushManager.subscribe({
  userVisibleOnly: true,
  applicationServerKey: urlBase64ToUint8Array(config.public.vapidKey),
})

// Unsubscribe:
const existing = await registration.pushManager.getSubscription()
if (existing) await existing.unsubscribe()
```

`urlBase64ToUint8Array` is a standard ~8-line utility (base64 → Uint8Array); implement inline in `app/lib/push/subscription.ts` — no library needed.

**VAPID public key:** `NUXT_PUBLIC_VAPID_KEY` runtime config. Exposed as `useRuntimeConfig().public.vapidKey`.

**Backend push registration endpoint:** `POST /push/subscribe` is **not listed** in Phase 4 API contract scope. Implement the full client-side subscribe/unsubscribe flow and subscription state feedback. Hold the `PushSubscription` object in component state ready for future wiring to the registration endpoint (Phase 5 scope).

### Permission State Mapping

| `Notification.permission` | UI state | Persian label |
|--------------------------|----------|---------------|
| `'default'` | `not-asked` | "دریافت اعلان‌ها را فعال کنید" |
| `'granted'` (subscribed) | `subscribed` | "اعلان‌ها فعال است" + "غیرفعال کردن" |
| `'denied'` | `blocked` | "اعلان‌ها مسدود شده‌اند — تنظیمات مرورگر را بررسی کنید" |
| PushManager unavailable | `unsupported` | "مرورگر شما از اعلان‌ها پشتیبانی نمی‌کند" |

**Safari/iOS caveat (from PITFALLS.md):** Web Push supported only on iOS 16.4+ installed PWA. `unsupported` state is the graceful path — never show an error, show a soft notice.

---

## 4. Notification Preferences Strategy (D-07)

**API (from docs/API.md §20):**
- `GET /notifications/preferences` — any auth, returns preferences object
- `PATCH /notifications/preferences` — any auth, upsert (first call creates record)

**Preference fields:**

```typescript
interface NotificationPreferences {
  id: string
  user_id: string
  meal_reminders: boolean
  water_reminders: boolean
  message_alerts: boolean
  diet_updates: boolean
}
```

Both client and nutritionist expose all four toggles (contract is shared per D-07). Persian label copy per toggle:

| Field | Persian label |
|-------|---------------|
| `meal_reminders` | "یادآور وعده‌های غذایی" |
| `water_reminders` | "یادآور مصرف آب" |
| `message_alerts` | "اعلان پیام‌های جدید" |
| `diet_updates` | "اعلان برنامه غذایی جدید" |

On toggle: immediately call `PATCH` with the changed field. Optimistic UI — revert on error. Show `InlineNotice` on failure.

---

## 5. Established Patterns to Reuse

| Pattern | Source file | Reuse in Phase 4 |
|---------|------------|-----------------|
| Typed composable + `useFetch` | `useTrackingApi.ts` | `useMessagingApi`, `useLabApi`, `useNotificationApi` |
| `$fetch` for mutations | `useAuthApi.ts` | `sendClientMessage`, `uploadLabResult` |
| Upload state surface (idle/uploading/success/failure) | `TrackingEntrySheet.vue` | `LabUploadSheet`, `MessageComposeBar` |
| `EmptyState`, `ErrorState` | Platform primitives | Message thread empty, lab list empty |
| `InlineNotice` | `InlineNotice.vue` | Push permission denied notice, connectivity notice |
| `SyncStateChip` | `SyncStateChip.vue` | Queued-offline message indicator in thread |
| `ConnectivityBanner` | `ConnectivityBanner.vue` | Lab upload blocked while offline |
| Role layouts + middleware | Phase 1–2 | New pages use existing `client.vue` / `nutritionist.vue` layouts |

---

## 6. Standard Stack

| Concern | Solution |
|---------|---------|
| Polling | `setInterval` / `clearInterval` in component lifecycle |
| Multipart upload | `$fetch` + `FormData` |
| Push subscription | Native Web Push API — no library |
| File download | `window.open(url)` or `<a href download>` |
| Message cache (offline) | Pinia ref in `client-offline` store, max 50 entries |
| Notification preferences | `useNotificationApi` composable |
| Auth continuity (D-10) | Existing `auth-fetch.client.ts` plugin wraps all new endpoints automatically |

---

*Research generated for: 04-messaging-notifications-lab-exchange*
*Sources: docs/API.md §16–17 §20 · docs/PRD.md §5.10–5.11 §6 §7 · .planning/phases/03-*/03-CONTEXT.md*
