# Phase 4 Validation: Messaging, Notifications & Lab Exchange

**Generated:** 2026-04-23
**Phase:** 04-messaging-notifications-lab-exchange

---

## UAT Criteria

These criteria map 1:1 to Phase 4 success criteria from ROADMAP.md. Each item must be true before the phase is marked complete.

---

### SC-1 — Conversation History (MSG-01)

> Client and nutritionist can read conversation history with unread state and stable refresh behavior.

- [ ] Client opens `/client/messages` and sees conversation thread with assigned nutritionist, ordered oldest-first
- [ ] Messages refresh without manual action within ~10 seconds of new message arrival
- [ ] Unread badge is visible on the messages nav item when unread messages exist; clears after opening conversation
- [ ] While offline: last-cached messages are visible; a connectivity notice is shown instead of a spinner loop
- [ ] Nutritionist opens `/nutritionist/messages` and sees list of client conversations with unread indicators
- [ ] Nutritionist selects a client and sees the conversation thread for that client

---

### SC-2 — Send Messages with Attachments (MSG-02)

> Client and nutritionist can send Persian text messages and supported file attachments.

- [ ] Client types a text message and sends; message appears in thread with `is_mine: true` alignment
- [ ] Client attaches a JPG/PNG ≤5 MB and sends; attachment thumbnail/name visible in thread
- [ ] Client attaches a PDF ≤10 MB and sends; attachment name and download link visible
- [ ] Client attaches an image >5 MB; inline Persian error shows before upload (message not sent)
- [ ] Client attaches a non-supported file type (e.g., .zip); inline Persian validation error shows
- [ ] Client attempts to send while offline; text-only message is queued with sync state chip ("در صف ارسال")
- [ ] Client attempts attachment send while offline; file attachment error shown ("فایل‌ها نیاز به اتصال دارند")
- [ ] Queued messages sync after reconnect; thread updates with server-assigned IDs
- [ ] Nutritionist sends text message to client; message appears in thread
- [ ] Nutritionist sends file attachment; attachment visible in conversation

---

### SC-3 — Push Notification Subscription (NOTF-01)

> Authenticated user can subscribe or unsubscribe from push notifications on supported devices.

- [ ] Notification settings page shows subscribe CTA with label "دریافت اعلان‌ها را فعال کنید"
- [ ] Tapping subscribe CTA triggers browser permission dialog
- [ ] After granting permission: UI shows "اعلان‌ها فعال است" + "غیرفعال کردن" button
- [ ] Tapping deactivate: unsubscribes from push; UI returns to `not-asked` state
- [ ] After denying permission: UI shows "اعلان‌ها مسدود شده‌اند — تنظیمات مرورگر را بررسی کنید" (no CTA)
- [ ] On browser without PushManager (unsupported): UI shows soft notice "مرورگر شما از اعلان‌ها پشتیبانی نمی‌کند" — no error state

---

### SC-4 — Notification Preferences (NOTF-02)

> Authenticated user can view and update notification preferences for reminder and message categories.

- [ ] User sees four Persian-labeled preference toggles on the notification settings page
- [ ] Toggles reflect current server state on page load (`GET /notifications/preferences`)
- [ ] Toggling any preference immediately calls `PATCH /notifications/preferences` and reflects the change
- [ ] First-time toggle (no existing record) succeeds via upsert — no 404 error
- [ ] On preference update failure: `InlineNotice` shown, toggle reverts to previous state (optimistic rollback)
- [ ] Both client and nutritionist see the same four toggles (role-aware copy, shared endpoint)

---

### SC-5 — Lab Results (LAB-01)

> Client and nutritionist can upload, view, and access lab results using file or link-based flows.

- [ ] Client opens `/client/labs`; sees uploaded lab result list with title, type, and date
- [ ] Lab list shows `EmptyState` when no results exist
- [ ] Client taps "افزودن نتیجه آزمایش"; upload sheet opens with title, type, date, notes, file/link fields
- [ ] Client attaches a PDF ≤10 MB and submits; result appears in list (upload state: uploading → success)
- [ ] Client provides an external link and submits; result appears in list
- [ ] Client leaves both file and link empty and submits; Persian validation error shown inline ("فایل یا لینک الزامی است")
- [ ] File-backed result: "دانلود" action triggers file download (browser download behavior)
- [ ] Link-backed result: "مشاهده" action opens external URL in new tab
- [ ] Upload attempted while offline: CTA disabled or `ConnectivityBanner` shown; no stuck-spinner
- [ ] Nutritionist on `/nutritionist/clients/:id/labs` sees all lab results for that client
- [ ] Nutritionist can trigger download/view for file-backed and link-backed results

---

## Regression Guards

| Guard | Verification |
|-------|-------------|
| Auth interceptor intact (D-10) | 401 on any new endpoint → redirect to auth login (not a white page) |
| Role isolation | Client cannot access `/nutritionist/messages/*` routes; nutritionist cannot access `/client/messages` directly |
| Phase 3 offline queue unaffected | Existing tracking domain queue entries (`food`, `water`, etc.) still replay correctly after message domain extension |
| Persian RTL layout | All new screens render RTL with correct text directionality in mobile viewport |

---

## Test Coverage Requirements

| Spec file | Scope |
|-----------|-------|
| `tests/client/messaging-poll.spec.ts` | Polling lifecycle (mount, unmount, refresh interval, dedup) |
| `tests/client/messaging-offline-queue.spec.ts` | Message domain enqueue, payload validation, text-only constraint |
| `tests/nutritionist/messaging-conversations.spec.ts` | Nutritionist unread list, per-client conversation render |
| `tests/client/lab-upload-state.spec.ts` | Upload state machine transitions, offline guard, file validation |
| `tests/platform/push-subscription.spec.ts` | Permission state machine, subscribe/unsubscribe actions, unsupported fallback |

---

*Validation generated for: 04-messaging-notifications-lab-exchange*
*Maps to ROADMAP.md Phase 4 success criteria SC-1 through SC-5*
