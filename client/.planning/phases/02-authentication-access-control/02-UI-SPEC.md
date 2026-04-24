---
phase: 02
slug: authentication-access-control
status: draft
shadcn_initialized: false
preset: none
created: 2026-04-23
---

# Phase 02 - UI Design Contract

> Visual and interaction contract for Phase 2 (Authentication and Access Control). Generated for planner and executor consumption.

---

## Design System

| Property | Value |
|----------|-------|
| Tool | none (existing custom Nuxt 4 + Tailwind 4 token system) |
| Preset | not applicable |
| Component library | project-wrapped Vue components (`AppShell`, `InlineNotice`, `EmptyState`, `ErrorState`) |
| Icon library | lucide-vue-next |
| Font | Vazirmatn (Persian-first) |

### Locked Direction (from Phase 1 + Phase 2 context)

- Tone remains clinical-minimal, calm, and high-legibility.
- Persian-only copy across all auth and guard states.
- Mobile-first RTL layout remains the default; desktop stays centered mobile canvas.
- Auth methods are role-locked:
  - Client: OTP only.
  - Nutritionist and super admin: email/password only.
- Shared auth surfaces remain neutral and must not expose role-private navigation.

---

## Spacing Scale

Declared values (multiples of 4 only):

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Inline icon-label gap, helper text separation |
| sm | 8px | Field subtext spacing, dense row gaps |
| md | 16px | Default vertical rhythm between controls |
| lg | 24px | Form section spacing and card padding |
| xl | 32px | Major block separation inside auth screens |
| 2xl | 48px | Distance between hero/title block and first actionable block |
| 3xl | 64px | Page-level top breathing room on tall screens |

Exceptions:
- Minimum interactive target: 44px.
- Primary auth buttons: fixed 48px height.
- OTP cells: 44px min width and height, with 8px inter-cell gap.

---

## Typography

| Role | Size | Weight | Line Height |
|------|------|--------|-------------|
| Body | 16px | 400 | 1.5 |
| Label | 14px | 600 | 1.4 |
| Heading | 20px | 600 | 1.3 |
| Display | 28px | 600 | 1.2 |

Rules:
- Keep exactly two weights: 400 and 600.
- Use Persian numerals for display text.
- Keep Latin digits for strict input correctness contexts only: mobile number, OTP code, and machine identifiers.

---

## Color

| Role | Value | Usage |
|------|-------|-------|
| Dominant (60%) | #F4F7F8 | App/auth background and large surfaces |
| Secondary (30%) | #FFFFFF | Form cards, grouped field containers, notices |
| Accent (10%) | #0F6B7A | Primary CTA, active role selector, focus ring, active OTP cell |
| Destructive | #B63D3D | Account/session failure emphasis, destructive confirmations |

Accent reserved for:
- Primary action button per screen.
- Selected role entry card/chip in auth gateway.
- Focus-visible ring around active field.
- Active OTP slot and single high-priority inline status.

---

## Copywriting Contract

| Element | Copy |
|---------|------|
| Primary CTA (client) | دریافت کد تایید |
| Primary CTA (client verify) | تایید و ورود |
| Primary CTA (nutritionist/admin) | ورود به حساب |
| Empty state heading | اطلاعاتی برای نمایش وجود ندارد |
| Empty state body | صفحه را دوباره بارگذاری کنید یا به ورود بازگردید. |
| Error state (generic) | ورود انجام نشد. اطلاعات را بررسی کنید و دوباره تلاش کنید. |
| Destructive confirmation | خروج از حساب: آیا مطمئن هستید می خواهید خارج شوید؟ |

Error messaging rules:
- Always show "problem + next step".
- Never reveal whether a specific account exists.
- Never expose backend/internal code names in UI text.
- Keep auth failure messages short and recovery-oriented.

---

## Registry Safety

| Registry | Blocks Used | Safety Gate |
|----------|-------------|-------------|
| shadcn official | none | not required |
| third-party | none | not applicable |

---

## Auth Surface Inventory

Phase 2 must deliver these UI surfaces:

1. Auth gateway (`/auth`)
- Neutral role-entry screen with exactly 3 explicit entry points:
  - ورود کاربر (OTP)
  - ورود متخصص تغذیه
  - ورود ادمین
- No private data, no role dashboard preview, no cross-role hints.

2. Client OTP request (`/auth/client`)
- Mobile number input.
- Primary CTA: دریافت کد تایید.
- Inline notice for cooldown and rate limits.

3. Client OTP verify (`/auth/client/verify`)
- 6-digit OTP input with auto-advance and backspace support.
- Resend control with cooldown.
- Editable mobile number shortcut (back to previous step).

4. Nutritionist credential sign-in (`/auth/nutritionist`)
- Email + password.
- Password visibility toggle.
- Submit disabled until valid form state.

5. Admin credential sign-in (`/auth/admin`)
- Same control pattern as nutritionist.
- Distinct heading and recovery copy for admin context.

6. Session-expired handoff notice (rendered in auth route after forced logout)
- One short inline banner: جلسه شما منقضی شد. دوباره وارد شوید.

---

## Form and Input Contract

### Client OTP request

| Field | Type | Validation | Input Behavior |
|------|------|------------|----------------|
| Mobile | tel | Iranian mobile format, required | Keep Latin digits; strip spaces; preserve leading zero |

Interactions:
- Disable CTA while request is pending.
- On success, route to verify step and start 60-second resend cooldown.
- Show normalized mobile in verify step in Persian sentence context.

### Client OTP verify

| Field | Type | Validation | Input Behavior |
|------|------|------------|----------------|
| Code | text (6 digits) | required, exactly 6 digits | Auto-focus next cell, paste support for full code, keep Latin digits |

Interactions:
- Verify button enabled only at 6 digits.
- Resend disabled during countdown; countdown visible in Persian copy.
- After 3 invalid attempts from backend response, lock verify action until resend.

### Nutritionist/Admin credential

| Field | Type | Validation | Input Behavior |
|------|------|------------|----------------|
| Email | email | required, valid email format | Keep LTR rendering inside field |
| Password | password | required, min 6 | Show/hide toggle; clear validation helper |

Interactions:
- Submit disabled while pending.
- Enter key submits when form valid.
- On repeated failure, keep password field value but move focus to top error notice.

---

## API Error to UI Message Mapping

| API Code | Surface | UI Copy (Persian) | UX Action |
|---------|---------|-------------------|-----------|
| INVALID_CREDENTIALS | Credential login | ایمیل یا رمز عبور نادرست است. دوباره تلاش کنید. | Keep user on form, focus error notice |
| OTP_INVALID | OTP verify | کد تایید صحیح نیست. دوباره وارد کنید. | Clear OTP cells, focus first cell |
| OTP_EXPIRED | OTP verify | اعتبار کد تایید تمام شده است. کد جدید دریافت کنید. | Disable verify, promote resend |
| OTP_MAX_ATTEMPTS | OTP verify | تعداد تلاش بیش از حد مجاز بود. کد جدید دریافت کنید. | Lock verify until resend |
| OTP_RATE_LIMIT | OTP send/verify | درخواست های شما بیش از حد مجاز است. کمی بعد دوباره تلاش کنید. | Disable send action temporarily |
| INVALID_TOKEN / TOKEN_REVOKED / UNAUTHORIZED | Any protected screen | جلسه شما منقضی شد. دوباره وارد شوید. | Force logout, clear state, redirect to role auth |
| RATE_LIMIT_EXCEEDED | Any auth call | درخواست ها زیاد است. چند لحظه دیگر دوباره تلاش کنید. | Keep form values, show retry guidance |
| INTERNAL_ERROR | Any auth call | خطایی رخ داد. لطفا دوباره تلاش کنید. | Keep context, allow retry |

Security rule:
- Do not render raw backend `message` directly; map code to controlled Persian copy.

---

## Guard and Transition Contract

### Route namespace rules

- Role namespaces remain strict:
  - `/client/**` for client.
  - `/nutritionist/**` for nutritionist.
  - `/admin/**` for super admin.
- Deny-by-default on namespace mismatch.
- Authenticated users must be redirected from `/auth/**` to their role root unless explicitly in logout/session-expired flow.

### Transition matrix

| From | Condition | To | UI Feedback |
|------|-----------|----|-------------|
| Protected route | No valid session | Matching auth entry route | Inline notice: برای ادامه وارد شوید |
| Protected route | Access token expired, refresh success | Original route | No disruptive toast; keep context |
| Protected route | Access token expired, refresh failed | Matching auth entry route | Session-expired banner on auth page |
| Any route | Role mismatch | Current role root route | Short notice: دسترسی به این بخش مجاز نیست |
| Auth route | Already authenticated | Role root route | Silent redirect after brief loading state |
| Logout action | API success or safe fallback | Matching auth entry route | Confirmation notice: با موفقیت خارج شدید |

### Single-flight refresh UX

- During refresh, block duplicate refresh requests globally.
- Keep current screen skeleton overlay if refresh is in-flight for more than 300ms.
- If refresh finishes under 300ms, do not flash loading UI.

---

## State Contract (Loading, Empty, Error, Offline)

### Loading states
- Button-level loading for auth submits.
- Page-level skeleton only for session recovery or guarded route bootstrap.

### Empty states
- Minimal usage in auth phase; only for unsupported deep-link fallback pages.
- Must use `EmptyState` component with one return action.

### Error states
- Field-level validation under the field.
- Form-level/auth-level error with `InlineNotice` (tone `error`).
- Fatal route-level auth failure may use `ErrorState` with retry + go-to-auth action.

### Offline behavior in auth
- Because authentication is online-required, show clear offline blocker:
  - اتصال اینترنت برقرار نیست. برای ورود، اینترنت را وصل کنید.
- Disable submit controls while offline.
- Keep typed values in-memory while offline banner is visible.

---

## Components and Reuse Requirements

Required reuse from Phase 1 primitives:
- `AppShell` with auth role shell.
- `ConnectivityBanner` and `UpdateAvailableBanner` remain active in auth layout.
- `InlineNotice` for non-fatal errors and recovery guidance.
- `ErrorState` / `EmptyState` only for route-level exceptional conditions.

New Phase 2 UI primitives to add:
- `AuthRolePicker`
- `OtpInput`
- `AuthFormCard`
- `SessionExpiredNotice`

All new components must preserve Phase 1 token usage and not introduce new ad-hoc colors or font scales.

---

## Accessibility and RTL Contract

- Every input has visible Persian label (no placeholder-only forms).
- Error notices use `role="alert"` for screen-reader announcements.
- OTP cells announce index context (digit 1 of 6, etc.).
- Keep touch targets >=44px.
- Respect `prefers-reduced-motion` baseline already defined in global CSS.
- Preserve RTL layout while allowing LTR rendering for email/mobile/OTP fields where readability requires it.

---

## Security UX Boundaries

- Never reveal whether mobile/email exists in system through differentiated copy.
- Never show token values or technical diagnostics in UI.
- Clear all role-scoped client state and persisted cache immediately on logout and forced session invalidation.
- Redirect targets after auth must be validated to app-internal allowed namespaces only.

---

## Checker Sign-Off

- [ ] Dimension 1 Copywriting: PASS
- [ ] Dimension 2 Visuals: PASS
- [ ] Dimension 3 Color: PASS
- [ ] Dimension 4 Typography: PASS
- [ ] Dimension 5 Spacing: PASS
- [ ] Dimension 6 Registry Safety: PASS

**Approval:** pending
