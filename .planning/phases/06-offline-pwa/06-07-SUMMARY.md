---
phase: "06"
plan: "07"
subsystem: frontend/push-preferences
tags: [pwa, service-worker, push, notifications, preferences]
dependency_graph:
  requires: [06-04, 06-05, 06-06]
  provides: [push-click-handler, client-notification-subscribe, notification-preferences-ui]
  affects: [frontend/app/service-worker/sw.ts, frontend/app/pages/client/plan.vue]
tech_stack:
  added: []
  patterns: [post-plan subscription sync, cache-first preferences, rtl notification payload]
key_files:
  created:
    - frontend/app/composables/useNotificationPermission.ts
    - frontend/app/stores/notificationPrefs.ts
    - frontend/app/pages/client/settings/notifications.vue
    - frontend/app/pages/client/profile.vue
    - frontend/.env.example
  modified:
    - frontend/app/service-worker/sw.ts
    - frontend/app/pages/client/plan.vue
    - frontend/nuxt.config.ts
decisions:
  - "Push subscription registration auto-syncs only after the client diet plan loads; permission prompting stays user-initiated from settings"
  - "Phase 6 fixed the injectManifest service-worker path bug by switching Nuxt PWA srcDir to service-worker"
  - "A minimal client profile page was created because no existing /client/profile page existed to host the notification settings entry point"
metrics:
  duration: "~25 minutes"
  completed: "2026-04-20"
  tasks_completed: 3
  files_changed: 8
---

# Phase 06 Plan 07: Frontend Push Notification Flow Summary

**One-liner:** Completed the user-facing push loop with a real service-worker notification handler, subscription permission composable, cached preference store, and Persian notification settings UI.

## What Was Built

### 1. Service worker push + click handling
- Replaced the Wave 1 push stub in `frontend/app/service-worker/sw.ts` with JSON payload parsing and Persian RTL notifications
- Added `tag`, `dir: 'rtl'`, `lang: 'fa'`, and safe fallback payload handling
- Added `notificationclick` deep-link navigation that focuses an existing app window when possible and otherwise opens the target route

### 2. Push permission + subscription registration
- Created `frontend/app/composables/useNotificationPermission.ts`
- Handles iOS non-installed guidance with the exact Persian message `برای دریافت اعلان ابتدا اپ را نصب کنید`
- Requests notification permission only from the settings CTA, then registers the push subscription with `POST /api/client/push/subscribe`
- Re-syncs an existing granted subscription only after the client active plan loads in `client/plan.vue`

### 3. Notification preferences cache + settings screen
- Created `frontend/app/stores/notificationPrefs.ts` with Dexie cache-first reads and backend sync for `GET/PATCH /api/client/push/preferences`
- Created `frontend/app/pages/client/settings/notifications.vue` with 6 Persian toggle switches for:
  - پیام جدید
  - برنامه غذایی جدید
  - نتیجه درخواست غذا
  - یادآور وعده غذایی
  - یادآور دارو
  - یادآور آب
- Added a new `frontend/app/pages/client/profile.vue` page so the existing client bottom-nav profile route is real and links to notification settings

## Deviations / Auto-fixed Issues

### 1. Fixed the pre-existing PWA build blocker
- `frontend/nuxt.config.ts` previously pointed injectManifest at `app/service-worker`, which made the build look for `app/app/service-worker/sw.ts`
- Updated the PWA config to `srcDir: 'service-worker'`, which allowed the custom worker to bundle successfully and unblocked Phase 6 verification

### 2. Actual backend routes use `/client/push/*`
- The plan examples referenced `/push/subscribe` and `/push/preferences`
- The implemented frontend uses the real backend routes from `backend/cmd/api/main.go`: `/api/client/push/subscribe` and `/api/client/push/preferences`

## Requested Output Notes

- **TypeScript lib additions in `sw.ts`:** none required during this plan; the existing `/// <reference lib="webworker" />` and `ServiceWorkerGlobalScope` declaration were already sufficient
- **`VITE_VAPID_PUBLIC_KEY` in service worker context:** not needed by `sw.ts`; the VAPID public key is consumed by `useNotificationPermission.ts`, and the injectManifest service-worker build now succeeds after the path fix. `frontend/.env.example` documents `NUXT_PUBLIC_VAPID_PUBLIC_KEY`, which must match backend `VAPID_PUBLIC_KEY`
- **`profile.vue` settings section:** there was no existing `frontend/app/pages/client/profile.vue`; the page and notification-settings entry point were created in this plan
- **Physical Android end-to-end push test:** not executed in this CLI environment
- **iOS-specific issues encountered:** no runtime iOS device testing was possible here; non-installed iOS guidance is implemented in code

## Validation

- ✅ `go test ./...` in `backend`
- ✅ `npm run test` in `frontend`
- ✅ `npm run build` in `frontend`
- ✅ Targeted ESLint on changed frontend files

## Self-Check: PASSED

- `frontend/app/service-worker/sw.ts` contains `showNotification` and `notificationclick`
- `frontend/app/composables/useNotificationPermission.ts` exists and posts subscriptions to `/api/client/push/subscribe`
- `frontend/app/stores/notificationPrefs.ts` exists and caches preferences in Dexie
- `frontend/app/pages/client/settings/notifications.vue` exists with 6 toggles
- `frontend/app/pages/client/profile.vue` exists and links to `/client/settings/notifications`
