---
phase: 01-foundation-infrastructure
plan: 05
title: "Frontend Auth Integration: Store, Pages & Route Guards"
subsystem: frontend
tags: [nuxt4, auth, pinia, jwt, otp, middleware, rtl, persian, cookies]

# Dependency graph
requires:
  - phase: 01-02
    provides: Nuxt 4 app with RTL/Persian foundation, toPersianDigits/toLatinDigits utils, constants
  - phase: 01-04
    provides: Backend auth API endpoints (login, OTP, refresh, logout, me)
provides:
  - useApi composable with transparent 401 refresh and mutex/queue pattern
  - Pinia auth store with login, OTP, logout, checkAuth, clearUser actions
  - Login page for admin/nutritionist email+password auth
  - OTP page for client mobile+code auth with two-step flow
  - Global auth middleware blocking unauthenticated access
  - Role guard middleware for per-page role validation
  - Unauthorized page with Persian message
affects: [02-all-plans, 03-all-plans, 04-all-plans]

# Tech tracking
tech-stack:
  added: []
  patterns: [mutex-refresh-queue, cookie-based-auth, global-middleware, role-guard-meta, two-step-otp-flow]

key-files:
  created:
    - frontend/app/composables/useApi.ts
    - frontend/app/stores/auth.ts
    - frontend/app/pages/auth/login.vue
    - frontend/app/pages/auth/otp.vue
    - frontend/app/middleware/auth.global.ts
    - frontend/app/middleware/role-guard.ts
    - frontend/app/pages/unauthorized.vue
  modified:
    - frontend/nuxt.config.ts

key-decisions:
  - "Updated apiBase to http://localhost:8080/api (was http://localhost:8080) to match backend /api prefix"
  - "Auth store normalizes both mobile AND code via toLatinDigits() for API calls"
  - "Unauthorized page uses auth layout with middleware:[] to avoid redirect loop"
  - "getDefaultRoute() helper in auth store uses ROLE_DEFAULT_ROUTES from constants.ts"

patterns-established:
  - "useApi() composable for all API calls — never use $fetch or useFetch for auth requests"
  - "credentials:'include' on every fetch (httpOnly cookies, never Authorization header)"
  - "auth.global.ts runs on ALL routes, public paths opt out via path check"
  - "role-guard middleware reads to.meta.roles for per-page role validation"
  - "definePageMeta({ middleware: [] }) to skip auth on login/otp/unauthorized pages"

requirements-completed: [AUTH-09]

# Metrics
duration: 7min
completed: 2026-04-19
---

# Phase 1 Plan 5: Frontend Auth Integration: Store, Pages & Route Guards Summary

**useApi composable with mutex/queue refresh pattern, Pinia auth store, Persian RTL login/OTP pages, and global auth + role-guard middleware**

## Performance

- **Duration:** 7 min
- **Started:** 2026-04-19T19:48:11Z
- **Completed:** 2026-04-19T19:55:04Z
- **Tasks:** 3/3
- **Files modified:** 8

## Accomplishments

- `useApi` composable with transparent 401 → refresh → retry flow using mutex/queue pattern (AUTH-09, D-02)
- Concurrent 401s share a single refresh request — no thundering herd
- All fetch calls use `credentials:'include'` for httpOnly cookie auth (D-01)
- Failed refresh clears auth state and redirects to /auth/login
- Pinia auth store manages user state with login, requestOTP, verifyOTP, logout, checkAuth, clearUser actions
- Login page renders Persian RTL form with email/password for admin/nutritionist roles
- OTP page renders two-step flow: mobile number entry → 6-digit code verification for clients
- Persian digit normalization: `toLatinDigits()` on mobile and OTP code inputs before API calls
- Post-login redirect per role: super_admin → /admin, nutritionist → /nutritionist/clients, client → /client/plan (D-18)
- Global auth middleware checks auth on ALL routes, redirects unauthenticated users to /auth/login
- checkAuth() restores session from cookie on hard refresh (calls GET /api/auth/me)
- Role guard middleware validates per-page `meta.roles` array
- Unauthorized page shows Persian access-denied message
- Updated `apiBase` to include `/api` prefix matching backend route structure
- Build passes successfully with all new files

## Task Commits

Each task was committed atomically:

1. **Task 1: useApi composable with 401 refresh queue** - `de4e8ac` (feat)
2. **Task 2: Pinia auth store, login page, and OTP page** - `09ff49f` (feat)
3. **Task 3: Global auth middleware and role guard** - `73c419c` (feat)

## Files Created/Modified

- `frontend/app/composables/useApi.ts` - API fetch wrapper with auto-refresh on 401
- `frontend/app/stores/auth.ts` - Pinia auth store with user state and auth actions
- `frontend/app/pages/auth/login.vue` - Admin/nutritionist email+password login page
- `frontend/app/pages/auth/otp.vue` - Client OTP request + verify page
- `frontend/app/middleware/auth.global.ts` - Global auth guard blocking unauthenticated access
- `frontend/app/middleware/role-guard.ts` - Per-page role validation middleware
- `frontend/app/pages/unauthorized.vue` - Access denied page with Persian message
- `frontend/nuxt.config.ts` - Updated apiBase to include /api prefix

## Decisions Made

- **apiBase update:** Changed default from `http://localhost:8080` to `http://localhost:8080/api` because all backend routes are under `/api/` prefix. Without this, every endpoint would need manual `/api` prefixing.
- **Mobile normalization:** Added `toLatinDigits()` for both mobile number and OTP code in the auth store, since users may type with Persian keyboard (e.g., ۰۹۱۲ instead of 0912).
- **Unauthorized page layout:** Uses `auth` layout with `middleware: []` to prevent redirect loops (global auth middleware would otherwise redirect to login).
- **getDefaultRoute helper:** Added to auth store for DRY role-based redirect logic, reusing `ROLE_DEFAULT_ROUTES` from constants.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] apiBase missing /api prefix**
- **Found during:** Task 1
- **Issue:** nuxt.config apiBase was `http://localhost:8080` but backend routes are all under `/api/*`. The useApi composable constructs URLs as `${apiBase}${endpoint}`, so endpoints like `/auth/login` would result in `http://localhost:8080/auth/login` instead of the correct `http://localhost:8080/api/auth/login`.
- **Fix:** Updated apiBase default to `http://localhost:8080/api`
- **Files modified:** `frontend/nuxt.config.ts`
- **Commit:** `de4e8ac`

**2. [Rule 2 - Missing functionality] Mobile number normalization in requestOTP**
- **Found during:** Task 2
- **Issue:** Plan only showed `toLatinDigits()` on OTP code, but mobile number input can also contain Persian digits (user types ۰۹۱۲ on Persian keyboard). Backend expects Latin digits for mobile validation.
- **Fix:** Added `toLatinDigits(mobile)` in both `requestOTP()` and `verifyOTP()` store actions
- **Files modified:** `frontend/app/stores/auth.ts`
- **Commit:** `09ff49f`

---

**Total deviations:** 2 auto-fixed (1 bug, 1 missing functionality)
**Impact on plan:** Both fixes are correctness requirements. No scope creep.

## Known Stubs

None — all auth flows are fully wired with proper API calls, error handling, and redirects.

## Self-Check: PASSED

All 7 created files and 1 modified file verified present. All 3 task commits (`de4e8ac`, `09ff49f`, `73c419c`) verified in git log.
