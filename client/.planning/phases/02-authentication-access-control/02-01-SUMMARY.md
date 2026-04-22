---
phase: 02-authentication-access-control
plan: 01
subsystem: auth-core
tags: [auth, session, refresh, logout]
requires: []
provides:
  - Typed authentication contracts and envelope parsing
  - Single-flight refresh orchestration for expired-token replay
  - Deterministic forced logout cleanup path
affects: [client-auth, nutritionist-auth, admin-auth]
tech-stack:
  added: []
  patterns: [typed-api-client, single-flight-refresh, role-aware-logout]
key-files:
  created: [app/types/auth.ts, app/lib/auth/error-map.ts, app/composables/useAuthApi.ts, app/composables/useSessionRefresh.ts, app/stores/auth-session.ts, app/plugins/auth-fetch.client.ts, tests/auth/session-refresh-singleflight.spec.ts, tests/auth/logout-state-cleanup.spec.ts]
  modified: []
key-decisions:
  - "Auth API errors are mapped to fixed Persian recovery messages before reaching UI state."
  - "Refresh retries are coalesced to one in-flight promise across concurrent 401 responses."
requirements-completed: [AUTH-03]
duration: 34min
completed: 2026-04-23
---

# Phase 2 Plan 01: Auth Core Infrastructure Summary

**Typed token lifecycle foundation with single-flight refresh and deterministic forced logout cleanup**

## Task Commits
1. `d74ee83` - feat(02-01): add typed auth contracts and safe api mapping
2. `a5a2b5c` - feat(02-01): implement single-flight refresh and logout cleanup

## Accomplishments
- Added typed request/response contracts for `/auth/login`, `/auth/otp/send`, `/auth/otp/verify`, `/auth/refresh`, and `/auth/logout`.
- Added controlled error-code mapping layer to prevent raw backend messages from leaking into auth UI.
- Implemented session refresh controller with single-flight coalescing and auth-fetch retry integration.
- Added role-safe logout cleanup orchestration and route-target resolution helpers.

## Deviations from Plan
None - plan executed as specified.

## Verification Results
- `npm run lint` -> PASS
- `npm run test:unit -- tests/auth/session-refresh-singleflight.spec.ts tests/auth/logout-state-cleanup.spec.ts` -> PASS (6/6)

## Known Stubs
None.

## Self-Check: PASSED
- Summary file created.
- Commits `d74ee83` and `a5a2b5c` present in git history.
