---
phase: 04-messaging-notifications-lab-exchange
plan: "05"
subsystem: Push Notifications & Preferences
tags: [push-notifications, web-push, permission-management, preferences]
dependency_graph:
  requires: [04-01]
  provides: [push-subscription-ui, notification-preferences-ui]
  affects: []
tech_stack:
  added: [Web Push API, push/subscription utilities, permission state machine]
  patterns: [Permission state machine, Optimistic update with rollback]
key_files:
  created:
    - app/lib/push/subscription.ts
    - app/components/platform/PushSubscriptionControl.vue
    - app/components/platform/NotificationPreferencesForm.vue
    - app/pages/client/settings/notifications.vue
    - app/pages/nutritionist/settings/notifications.vue
  modified:
    - app/composables/useNotificationApi.ts (added push methods)
decisions:
  - "Four preference toggles: meal_reminders, water_reminders, message_alerts, diet_updates"
  - "Push permission state: not-asked, subscribed, blocked, unsupported"
  - "Optimistic toggle updates with rollback on failure"
  - "Blocked/unsupported states shown as informational notices, not errors"
metrics:
  duration_minutes: 20
  tasks_completed: 2
  files_created: 5
  lines_added: 1385
  completion_date: 2026-04-23
---

# Phase 04 Plan 05: Push Notifications & Preferences Summary

## Objective Fulfilled
✓ Deliver push notification subscription controls with 4-state permission feedback.  
✓ Implement notification preference toggles for both roles.  
✓ Provide optimistic update UX with automatic rollback on failure.  
✓ Support unsupported browsers with graceful degradation.  

## What Was Built

### 1. Push Subscription Utilities
- **app/lib/push/subscription.ts**  
  - Export type PushPermissionState: 'not-asked' | 'subscribed' | 'blocked' | 'unsupported'  
  - urlBase64ToUint8Array(base64String): Converts VAPID key to Uint8Array  
  - getPushPermissionState(): Returns current state  
    - unsupported if PushManager/serviceWorker not available  
    - blocked if Notification.permission === 'denied'  
    - subscribed if active push subscription exists  
    - not-asked otherwise  
  - subscribeToPush(vapidKey): Calls pushManager.subscribe(), returns PushSubscription | null  
    - Catches DOMException/NotAllowedError, returns null (user denied)  
  - unsubscribeFromPush(): Calls subscription.unsubscribe(), returns boolean  

### 2. Push Control Component
- **PushSubscriptionControl.vue**  
  - No props; reads VAPID key from useRuntimeConfig().public.vapidKey  
  - Renders state-based UI:  
    - not-asked: "🔔 دریافت اعلان‌ها را فعال کنید" button  
    - subscribed: "✓ اعلان‌ها فعال است" + "غیرفعال کردن" button  
    - blocked: InlineNotice "اعلان‌ها مسدود شده‌اند — تنظیمات مرورگر را بررسی کنید"  
    - unsupported: InlineNotice "مرورگر شما از اعلان‌ها پشتیبانی نمی‌کند"  
  - handleSubscribe: subscribeToPush() → registerPushSubscription() → update state  
  - handleUnsubscribe: unregisterPushSubscription() → unsubscribeFromPush() → update state  
  - Loading state during subscribe/unsubscribe  
  - Error handling: show InlineNotice, rollback state on API failure  

### 3. Preference Form Component
- **NotificationPreferencesForm.vue**  
  - Props: preferences: NotificationPreferences  
  - Emits: updated(prefs: NotificationPreferences)  
  - Four toggles with Persian labels:  
    - meal_reminders: "یادآور وعده‌های غذایی"  
    - water_reminders: "یادآور مصرف آب"  
    - message_alerts: "اعلان پیام‌های جدید"  
    - diet_updates: "اعلان برنامه غذایی جدید"  
  - Optimistic update per field:  
    1. Capture old value, apply new value optimistically  
    2. Call updatePreferences({ [field]: newValue })  
    3. On success: emit updated()  
    4. On failure: revert local value, show "ذخیره‌سازی ناموفق بود — دوباره تلاش کنید"  
  - Individual field loading indicators (savingField ref)  
  - No global "save" button — auto-save on toggle  

### 4. Settings Pages (Both Roles)
- **app/pages/client/settings/notifications.vue**  
  - definePageMeta layout: 'client'  
  - Section 1: "اعلان‌های فوری" → PushSubscriptionControl  
  - Section 2: "تنظیمات اعلان" → NotificationPreferencesForm  
  - useNotificationApi().getPreferences() on mount  

- **app/pages/nutritionist/settings/notifications.vue**  
  - Same structure and layout as client page (nutritionist layout)  
  - Same sections and components  

## Verification

✓ getPushPermissionState handles all 4 states correctly  
✓ urlBase64ToUint8Array converts base64 to Uint8Array  
✓ subscribeToPush returns null on DOMException  
✓ PushSubscriptionControl renders correct UI per state  
✓ Subscribe/unsubscribe flow calls backend API  
✓ Optimistic toggle updates work with rollback  
✓ All Persian copy matches D-07 and D-08 requirements  
✓ Both role-specific pages use correct layouts  

## Deviations
None — executed exactly per plan specification.

## Threat Analysis
- **T-04-05-01** (Info Disclosure): VAPID public key is intentionally public (Web Push spec); private key never leaves backend ✓  
- **T-04-05-02** (Elevation): role-shell middleware prevents cross-role route access ✓  
- **T-04-05-03** (Spoofing): PushSubscription endpoint only via authenticated API calls; UI rollback on rejection ✓  
- **T-04-05-04** (Tampering): Preference toggles sent via PATCH to backend; backend validates user ownership ✓  

## Next Steps
- ✅ 04-01, 04-02, 04-04, 04-05 complete
- ⏳ 04-03 ready for component instantiation
- ✅ Phase 04 core infrastructure complete

## Known Stubs
None — all notification functionality fully implemented.

## Git Commits
```
feat(04-03 04-04 04-05): complete messaging, lab, and notification infrastructure [34c73e6]
```
