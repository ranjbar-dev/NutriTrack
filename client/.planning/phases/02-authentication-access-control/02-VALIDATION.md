# Phase 02 Validation Strategy

**Phase:** 02 - Authentication & Access Control
**Date:** 2026-04-23
**Status:** Active

## Validation Objectives

- Validate AUTH-01 through AUTH-04 with deterministic automated tests.
- Ensure role namespace boundaries and redirect behavior are enforced across direct entry, refresh, and navigation.
- Ensure auth error handling remains secure and user-safe in Persian UI copy.

## Requirement Coverage Matrix

| Requirement | Validation Paths | Evidence Artifact |
|-------------|------------------|-------------------|
| AUTH-01 | `tests/auth/client-otp-validation.spec.ts`, `tests/auth/client-otp-flow.spec.ts` | 02-02 summary + unit test output |
| AUTH-02 | `tests/auth/auth-gateway-neutrality.spec.ts`, `tests/auth/role-credential-login.spec.ts` | 02-03 summary + unit test output |
| AUTH-03 | `tests/auth/session-refresh-singleflight.spec.ts`, `tests/auth/logout-state-cleanup.spec.ts`, `tests/auth/session-expiry-redirect.spec.ts` | 02-01 and 02-04 summaries + unit test output |
| AUTH-04 | `tests/auth/route-access-control.spec.ts`, `tests/platform/shell-role-isolation.spec.ts` | 02-04 summary + unit test output |

## Verification Gates

1. `npm run lint`
2. `npm run typecheck`
3. `npm run test:unit -- tests/auth/session-refresh-singleflight.spec.ts tests/auth/logout-state-cleanup.spec.ts`
4. `npm run test:unit -- tests/auth/client-otp-validation.spec.ts tests/auth/client-otp-flow.spec.ts`
5. `npm run test:unit -- tests/auth/auth-gateway-neutrality.spec.ts tests/auth/role-credential-login.spec.ts`
6. `npm run test:unit -- tests/auth/route-access-control.spec.ts tests/auth/session-expiry-redirect.spec.ts tests/platform/shell-role-isolation.spec.ts`

All gates must pass before phase verification is marked complete.

## Risk Checks

- Assert refresh remains single-flight and does not create retry storms.
- Assert INVALID_TOKEN, TOKEN_REVOKED, and UNAUTHORIZED enforce safe logout and route handoff.
- Assert role namespace guards prevent cross-role route access and navigation leakage.
- Assert auth error rendering uses mapped Persian-safe messages and not raw backend content.
