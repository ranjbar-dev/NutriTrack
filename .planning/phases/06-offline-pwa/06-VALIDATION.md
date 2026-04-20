# Phase 6: Offline & PWA — Validation Specification

**Generated:** 2026-04-20
**Phase:** 06-offline-pwa
**Framework:** Vitest 3.x (frontend) + Go test (backend)

---

## Test Framework

| Property | Value |
|----------|-------|
| Frontend framework | Vitest 3.x (already in package.json) |
| Config file | `frontend/vitest.config.ts` |
| Frontend run command | `cd frontend && npm run test` |
| Backend run command | `cd backend && go test ./...` |

---

## Wave-by-Wave Automated Test Map

| Wave | Plan | Req IDs | Test File | Command | Status |
|------|------|---------|-----------|---------|--------|
| 2 | 06-02 | OFFL-05, OFFL-12 | `frontend/tests/useDb.test.ts` | `cd frontend && npm run test -- tests/useDb.test.ts` | ❌ Created in Wave 2 |
| 3 | 06-03 | OFFL-06, OFFL-07, OFFL-08 | `frontend/tests/useSyncManager.test.ts` | `cd frontend && npm run test -- tests/useSyncManager.test.ts` | ❌ Created in Wave 3 |
| 4 | 06-05 | OFFL-02, OFFL-10, OFFL-11 | `frontend/tests/clientPlan.offline.test.ts` | `cd frontend && npm run test -- tests/clientPlan.offline.test.ts` | ❌ Created in Wave 4 |
| 5 | 06-06 | NOTIF-07 | `backend/internal/service/push_service_test.go` | `cd backend && go test ./internal/service/... -run "TestSendToClient|TestReminderAlready" -v` | ❌ Created in Wave 5 |

### Sampling Continuity Check (per Dimension 8c)

| Wave | Tasks with automated verify | Total tasks | Meets ≥2/3 rule |
|------|-----------------------------|-------------|-----------------|
| 2 | 2 (useDb.test.ts) | 3 | ✅ |
| 3 | 2 (useSyncManager.test.ts) | 2 | ✅ |
| 4 | 1 (clientPlan.offline.test.ts) | 2 | ✅ (only 2 tasks) |
| 5 | 1 (push_service_test.go) | 2 | ✅ (only 2 tasks) |

---

## Per-Requirement Automated Verification

| Req ID | Description | Test File | Test Name | Plan |
|--------|-------------|-----------|-----------|------|
| OFFL-05 | Dexie schema initializes with all 6 tables | `tests/useDb.test.ts` | `schema initializes with correct tables` | 06-02 |
| OFFL-12 | iOS eviction detection on app open | `tests/useDb.test.ts` | `iOS eviction detected when table empty after recent fetch` | 06-02 |
| OFFL-06 | Sync queue FIFO processing | `tests/useSyncManager.test.ts` | `OFFL-06: enqueue adds entry with status=pending` | 06-03 |
| OFFL-07 | local_id deduplication (server-side) | `tests/useSyncManager.test.ts` | `OFFL-07: same local_id queued twice yields two entries` | 06-03 |
| OFFL-08 | Exponential backoff delays | `tests/useSyncManager.test.ts` | `OFFL-08: backoffMs returns correct delays` | 06-03 |
| OFFL-02 | Active plan cached and retrievable from IndexedDB | `tests/clientPlan.offline.test.ts` | `OFFL-02: activePlan cached and retrievable` | 06-05 |
| OFFL-11 | Last 50 messages per conversation cached | `tests/clientPlan.offline.test.ts` | `OFFL-11: messages table stores and retrieves last 50 per partner` | 06-05 |
| NOTIF-07 | Preference opt-out blocks push send | `backend/internal/service/push_service_test.go` | `TestSendToClient_SkipsWhenPreferenceDisabled` | 06-06 |

---

## Manual Validation Checklist

### OFFL-01 / UI-06 / UI-07: PWA Shell & Install

- [ ] `npm run build` completes without errors
- [ ] `frontend/.output/public/sw.js` exists
- [ ] `frontend/.output/public/manifest.webmanifest` contains `نوتری‌ترک` and `"display":"standalone"`
- [ ] Chrome DevTools → Application → Manifest shows Persian name, RTL direction, 3 shortcuts
- [ ] Chrome DevTools → Application → Service Workers shows SW registered and active
- [ ] Android Chrome: install prompt appears; app launches standalone from home screen

### OFFL-02 / OFFL-10: Offline Diet Plan View

- [ ] Open app online, load client plan page → observe plan content
- [ ] Enable airplane mode, reload/reopen → plan still renders from IndexedDB cache
- [ ] Reconnect, open app again → plan refreshes from API

### OFFL-03: Tracking Writes Offline

- [ ] Enable airplane mode, log water intake → `ClientSyncStatus` shows pending count ≥ 1
- [ ] `همگام‌سازی` / `در انتظار` labels appear in Persian
- [ ] Reconnect → pending count clears within 30 seconds; server record exists
- [ ] Force 3 failures (block server) → item shows `خطا` with `تلاش مجدد` button
- [ ] Click `تلاش مجدد` → item re-queued and synced

### OFFL-04 / OFFL-11: Offline Messages

- [ ] Open conversation, see messages; enable airplane mode → cached messages still display
- [ ] Type message offline → local echo immediately visible (no error)
- [ ] Reconnect → local echo replaced with server response; server has the message

### OFFL-12: iOS Storage Eviction

- [ ] On iOS Safari (≥16.4): clear website data in Settings → reopen app
- [ ] Persian notice "اطلاعات آفلاین توسط دستگاه حذف شد" appears
- [ ] App re-fetches plan on reconnect (no blank screen or crash)

### NOTIF-01 / NOTIF-02: Push Subscription

- [ ] `POST /api/push/subscribe` with valid client JWT → HTTP 201; row in `push_subscriptions`
- [ ] Notification permission prompt appears in browser after plan loads (D-16)
- [ ] iOS non-standalone: banner "برای دریافت اعلان ابتدا اپ را نصب کنید" shown instead of permission prompt

### NOTIF-03: Event-Driven Push Notifications

- [ ] Nutritionist sends message to client → client device receives push "پیام جدید"
- [ ] Nutritionist activates plan for client → client receives push "برنامه جدید فعال شد"
- [ ] Nutritionist approves food request → client receives push "درخواست تأیید شد"
- [ ] Nutritionist rejects food request → client receives push "درخواست رد شد"

### NOTIF-04: Meal Reminders

- [ ] Set a meal time 3 minutes in the future; wait → push "وقت وعده غذایی" fires within 1 minute
- [ ] Same reminder does NOT fire again in the next scheduler tick (dedup working)

### NOTIF-05: Medication Reminders

- [ ] Set a medication time 3 minutes in the future; wait → push "وقت دارو" fires within 1 minute

### NOTIF-06: Water Reminders

- [ ] Enable water_reminder preference for a client; wait until :00 or :02 of an even hour (08, 10, 12, 14, 16, 18, 20)
- [ ] Push "یادآوری آب" fires within 5 minutes of the hour boundary
- [ ] Same reminder does NOT fire again for same hour (dedup working)

### NOTIF-07 / NOTIF-08: Preferences & Payload

- [ ] Notification settings page at `/client/settings/notifications` shows 6 Persian toggles
- [ ] Disable "new_message" toggle; have nutritionist send message → push NOT received
- [ ] Push payload contains `title`, `body`, and valid `url` for deep-link navigation
- [ ] Notification click navigates to the correct in-app screen

---

## Service Worker Manual Tests (Chrome DevTools)

Service worker behavior cannot be unit-tested with Vitest; validate manually:

1. Application → Service Workers → SW status: "activated and is running"
2. Application → Cache Storage → precached entries appear (JS, CSS, fonts)
3. Network tab → Offline → reload → static assets served from cache (200 from SW)
4. Application → Background Sync → trigger sync tag `nutritrack-sync` → verify `ClientSyncStatus` processes pending items
5. Application → Push → send test push → verify Persian notification appears with RTL layout

---

## Full Suite Command (run before phase sign-off)

```bash
# Frontend unit tests
cd frontend && npm run test

# Backend unit tests (push service)
cd backend && go test ./internal/service/... -run "TestSendToClient|TestReminderAlready" -v

# Build verification
cd frontend && npm run build && ls .output/public/sw.js .output/public/manifest.webmanifest
cd backend && go build ./...
```
