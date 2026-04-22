# Phase 2: Authentication & Access Control - Research

**Researched:** 2026-04-23
**Domain:** Nuxt 4 client authentication UX, token lifecycle, and role-bound route access control
**Confidence:** HIGH

## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Client authentication uses OTP flow only through documented endpoints (`/auth/otp/send`, `/auth/otp/verify`) and lands in client route namespace.
- **D-02:** Nutritionist and super admin authentication use email/password through `/auth/login` and land in role-specific route namespaces.
- **D-03:** The frontend must not introduce alternate auth methods (social, magic links, self-registration) in this phase.
- **D-04:** Access token + refresh token lifecycle is implemented via API contract, with refresh handled through a single-flight guard to avoid parallel refresh storms.
- **D-05:** On auth failure (`INVALID_TOKEN`, `TOKEN_REVOKED`, `UNAUTHORIZED`), user is safely logged out and redirected to role-appropriate auth entry.
- **D-06:** Logout uses `/auth/logout` and clears all user-scoped runtime state, persisted state, and role caches.
- **D-07:** Route middleware enforces namespace access by role (`/client/**`, `/nutritionist/**`, `/admin/**`) with deny-by-default behavior.
- **D-08:** Shared auth pages remain neutral and do not expose role-private navigation.
- **D-09:** Access guard checks apply consistently on direct URL entry, refresh, and route transitions.
- **D-10:** All auth copy stays Persian-only, short, and recovery-oriented.
- **D-11:** OTP UX must expose cooldown/retry states and clear validation feedback without leaking sensitive auth internals.
- **D-12:** Numeric entry formatting allows Latin digits where input correctness requires it (OTP/mobile), while display contexts remain Persian-first.

### the agent's Discretion
- Exact form layout composition per role entry screen.
- Internal store/composable file boundaries, as long as session and guard responsibilities stay explicit.
- Retry/backoff timing details for refresh and transient network failures within reasonable UX constraints.

### Deferred Ideas (OUT OF SCOPE)
- MFA beyond current role model, social login, or magic links.
- Server-driven RBAC matrix editor or permission-management UI.
- Advanced account recovery beyond current API contract.

## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| AUTH-01 | Client can request and verify an OTP using the documented mobile-based authentication flow. | OTP form contract, cooldown/attempt UX, and API-code mapping are specified in Standard Stack, Architecture Patterns, and Common Pitfalls. |
| AUTH-02 | Nutritionist and super admin can log in with email and password. | Credential form validation, secure error mapping, and role landing rules are defined in Architecture Patterns and Security Domain. |
| AUTH-03 | Session persists with refresh, with safe expiry/logout redirects. | Token lifecycle model, single-flight refresh guard, and logout cleanup contract are defined in Standard Stack, Architecture Patterns, and Don't Hand-Roll. |
| AUTH-04 | Authenticated user can access only role-allowed routes/identity scope. | Namespace middleware strategy and deny-by-default route protection are defined in Architectural Responsibility Map and Architecture Patterns. |

## Project Constraints (from copilot-instructions.md)

- Frontend-only scope; no backend or infrastructure changes in this phase. [VERIFIED: repository instruction file]
- Stack is fixed to Nuxt 4, Tailwind CSS 4, and Pinia with TypeScript-first composable patterns. [VERIFIED: repository instruction file]
- `docs/API.md` is the backend contract; `docs/PRD.md` is product behavior contract. [VERIFIED: repository instruction file]
- Persian-only RTL mobile PWA quality is primary, desktop is secondary. [VERIFIED: repository instruction file]
- Role boundaries must remain explicit in routing, state, cache, and persistence. [VERIFIED: repository instruction file]

## Summary

Phase 2 should be implemented as a strict separation of concerns: role-specific auth UI in auth routes, centralized session/token logic in one auth client layer, and namespace authorization in global route middleware. This aligns directly with existing Phase 1 route shell isolation and Nuxt 4 middleware behavior. [VERIFIED: workspace code + CITED: https://nuxt.com/docs/4.x/guide/directory-structure/middleware]

The API contract already defines OTP, credential login, refresh, logout, and auth error codes. The frontend should treat backend `code` as the canonical discriminator and map each code to a fixed Persian recovery message rather than surfacing backend `message`. This preserves UX clarity and avoids user/account enumeration patterns. [VERIFIED: docs/API.md + CITED: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#authentication-and-error-messages]

Nuxt 4 already provides `$fetch` (ofetch-backed) and route middleware primitives needed to implement auth without additional auth SDKs. The recommended implementation is a single shared API client with request/response interceptors and a process-wide single-flight refresh promise to serialize refresh attempts under concurrent 401 responses. [CITED: https://nuxt.com/docs/4.x/getting-started/data-fetching + CITED: https://github.com/unjs/ofetch]

**Primary recommendation:** Implement one centralized auth transport layer (token attach + single-flight refresh + error-to-logout), then keep page-level auth UI thin and role-specific.

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| OTP send/verify UI validation and cooldown feedback | Browser / Client | API / Backend | Input shape, timers, and UX lock states are client concerns; OTP truth and attempt counters are backend authority. [VERIFIED: docs/API.md + ASSUMED] |
| Credential login form validation and submission | Browser / Client | API / Backend | Client validates format and submits; backend validates credentials and issues tokens. [VERIFIED: docs/API.md] |
| Token issue/refresh/revoke protocol | API / Backend | Browser / Client | Token minting/revocation is server-side; client stores and rotates tokens per contract. [VERIFIED: docs/API.md] |
| Single-flight refresh guard | Browser / Client | — | Concurrent request coordination is a client transport concern. [CITED: https://github.com/unjs/ofetch + ASSUMED] |
| Role namespace route access control | Frontend Server (SSR) | Browser / Client | Nuxt route middleware runs during SSR and on client transitions; enforce deny-by-default there. [CITED: https://nuxt.com/docs/4.x/guide/directory-structure/middleware] |
| Secure error mapping | Browser / Client | — | UI-safe Persian messaging is frontend responsibility; do not expose backend internals. [VERIFIED: docs/API.md + CITED: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#authentication-and-error-messages] |
| Logout cleanup (state + persistence + caches) | Browser / Client | API / Backend | API revokes session; client must clear local scoped state and persisted artifacts. [VERIFIED: docs/API.md + ASSUMED] |

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Nuxt | 4.4.2 | App framework, route middleware, SSR-aware navigation | Native middleware model fits role-bound route guards and redirect flow. [VERIFIED: npm registry + CITED: https://nuxt.com/docs/4.x/guide/directory-structure/middleware] |
| Pinia | 3.0.4 | Auth/session state store (tokens, role, user id, status) | First-party store ecosystem for Nuxt and existing project pattern continuity. [VERIFIED: npm registry + CITED: https://pinia.vuejs.org/core-concepts/] |
| Tailwind CSS | 4.2.4 | Auth surface styling tokens and utility composition | Existing project stack lock and visual consistency with Phase 1 design tokens. [VERIFIED: npm registry + VERIFIED: repository instruction file] |
| Nuxt `$fetch` (ofetch-backed) | bundled with Nuxt 4.4.2 | API calls with interceptors for auth header, refresh handling, and error mapping | Supports `onRequest`/`onResponseError` interception and client reuse via shared fetch instance patterns. [CITED: https://nuxt.com/docs/4.x/getting-started/data-fetching + CITED: https://github.com/unjs/ofetch] |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| @pinia/nuxt | 0.11.2 | Nuxt integration for Pinia stores | Required for store auto-registration in Nuxt app runtime. [VERIFIED: package.json] |
| zod | 4.3.6 | Schema validation for OTP/mobile/email/password inputs and API response guards | Use at user-input boundaries and response parsing boundaries for deterministic errors. [VERIFIED: npm registry + ASSUMED] |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Nuxt middleware + custom auth transport | Third-party auth framework (e.g., sidebase/nuxt-auth) | Adds abstraction mismatch with OTP+role split custom contract; less direct control over code-to-message mapping. [ASSUMED] |
| Pinia auth store | useState-only scattered session refs | Lower structural clarity and harder centralized logout cleanup across roles. [ASSUMED] |
| zod runtime validation | Manual regex-only validation | Lower maintainability for evolving payload/error contracts and weaker type inference. [ASSUMED] |

**Installation:**
```bash
npm install zod
```

**Version verification:**
- `nuxt` -> `4.4.2`, modified `2026-03-12T11:59:41.082Z`. [VERIFIED: npm registry]
- `pinia` -> `3.0.4`, modified `2025-11-05T09:25:14.059Z`. [VERIFIED: npm registry]
- `tailwindcss` -> `4.2.4`, modified `2026-04-22T15:33:29.666Z`. [VERIFIED: npm registry]
- `zod` -> `4.3.6`, modified `2026-01-25T21:51:57.252Z`. [VERIFIED: npm registry]

## Architecture Patterns

### System Architecture Diagram

```text
[Auth Route UI (/auth/*)]
   | submit form (OTP send / OTP verify / credentials)
   v
[Auth Composable Layer]
   | validates input + calls api client
   v
[Shared API Client ($fetch/ofetch instance)]
   | onRequest: attach access token
   | onResponseError 401: trigger single-flight refresh
   |   |- refresh success -> replay original request
   |   |- refresh fail -> hard logout cleanup
   v
[NutriTrack API (/auth/login, /auth/otp/*, /auth/refresh, /auth/logout)]

[Nuxt Global Route Middleware]
   | reads auth store role/session
   | deny-by-default namespace check
   |- allow -> target page
   |- deny/no session -> role auth entry
```

### Recommended Project Structure

```text
app/
  composables/
    useAuthApi.ts              # login/sendOtp/verifyOtp/logout wrappers
    useAuthSession.ts          # token load/save/clear + refresh single-flight
    useAuthErrorMap.ts         # api code -> Persian safe message mapping
  stores/
    auth-session.ts            # role, user_id, token status, auth bootstrap state
  middleware/
    role-shell.global.ts       # namespace guard (extend existing)
    auth-redirect.global.ts    # auth<->protected redirect rules
  pages/
    auth/
      index.vue                # role picker
      client.vue               # OTP request
      client/verify.vue        # OTP verify
      nutritionist.vue         # email/password
      admin.vue                # email/password
```

### Pattern 1: Return `navigateTo()` from route middleware
**What:** Nuxt route middleware must return redirect/abort values directly.
**When to use:** Role mismatch, unauthenticated access to protected route, authenticated access to auth routes.
**Example:**
```typescript
// Source: https://nuxt.com/docs/4.x/api/utils/navigate-to
export default defineNuxtRouteMiddleware((to) => {
  if (to.path.startsWith('/admin') && !isAdmin()) {
    return navigateTo('/auth/admin')
  }
})
```

### Pattern 2: Centralized fetch instance with interceptors
**What:** Create one API client instance with auth header injection + response error interception.
**When to use:** Any call requiring consistent token attach and refresh behavior.
**Example:**
```typescript
// Source: https://github.com/unjs/ofetch
const api = $fetch.create({
  baseURL: '/api/v1',
  onRequest({ options }) {
    if (accessToken.value) {
      options.headers = { ...(options.headers || {}), Authorization: `Bearer ${accessToken.value}` }
    }
  },
  async onResponseError(ctx) {
    if (ctx.response.status === 401) {
      await refreshSingleFlight()
    }
  }
})
```

### Pattern 3: Single-flight refresh guard
**What:** Store one in-memory `refreshPromise`; reuse it for all simultaneous 401 failures.
**When to use:** Concurrent protected requests after access token expiry.
**Example:**
```typescript
// Source: project pattern recommendation; Promise-based fan-in [ASSUMED]
let refreshPromise: Promise<void> | null = null

async function refreshSingleFlight() {
  if (!refreshPromise) {
    refreshPromise = doRefresh().finally(() => {
      refreshPromise = null
    })
  }
  await refreshPromise
}
```

### Anti-Patterns to Avoid

- **Per-page auth logic duplication:** leads to inconsistent redirects and stale token behavior; keep transport/session logic centralized. [ASSUMED]
- **Showing backend raw `message` directly:** can leak internals and create discrepancy factors; map by `code` to fixed Persian messages. [VERIFIED: docs/API.md + CITED: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#authentication-and-error-messages]
- **Using `useRoute()` inside middleware guards:** Nuxt warns this is inaccurate in middleware context; use `to` and `from` args. [CITED: https://nuxt.com/docs/4.x/guide/directory-structure/middleware]
- **Non-returned `navigateTo()` in middleware:** redirect may behave unexpectedly; always return it. [CITED: https://nuxt.com/docs/4.x/api/utils/navigate-to]

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Generic network layer with ad-hoc retry branches in every call | Per-request token attach/refresh logic | One shared `$fetch.create()` client with interceptors | Prevents inconsistent token behavior and duplicated edge-case code. [CITED: https://github.com/unjs/ofetch] |
| Field validation with scattered inline checks | Repeated regex snippets per component | Schema-driven validation (zod) | Centralizes constraints and error mapping for OTP/credentials. [ASSUMED] |
| Role access checks in every page component | Manual checks in page setup | Global middleware + named page middleware where needed | Nuxt route middleware is purpose-built and runs in known order. [CITED: https://nuxt.com/docs/4.x/guide/directory-structure/middleware] |
| Displaying backend auth messages directly | Raw `message` passthrough | Controlled `code -> Persian message` map | Reduces user enumeration and internal leakage risk. [VERIFIED: docs/API.md + CITED: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#authentication-and-error-messages] |

**Key insight:** Custom auth UX is required, but auth plumbing should be standardized and centralized; most production auth defects come from duplicated edge handling, not from missing UI screens. [ASSUMED]

## Common Pitfalls

### Pitfall 1: Refresh storm after token expiry
**What goes wrong:** multiple concurrent 401 responses trigger multiple refresh calls.
**Why it happens:** no shared refresh promise/lock around refresh workflow.
**How to avoid:** single-flight guard around refresh + replay queue for pending requests.
**Warning signs:** repeated `/auth/refresh` bursts in network logs; out-of-order token overwrites.

### Pitfall 2: Redirect loops between auth and protected routes
**What goes wrong:** middleware redirects authenticated users into `/auth/*` or unauth users into protected paths repeatedly.
**Why it happens:** missing `to.path` guards and missing return semantics from middleware redirects.
**How to avoid:** explicit path guards and `return navigateTo(...)` in each redirect branch.
**Warning signs:** browser history flood, repeated same-path transitions.
[CITED: https://nuxt.com/docs/4.x/guide/directory-structure/middleware + CITED: https://nuxt.com/docs/4.x/api/utils/navigate-to]

### Pitfall 3: Leaky auth error UX
**What goes wrong:** different messages/status pathways reveal account existence or auth internals.
**Why it happens:** direct backend message rendering and inconsistent mapping.
**How to avoid:** deterministic code mapping to generic recovery-oriented Persian messages.
**Warning signs:** different copy for unknown email vs wrong password; role/account hints in UI.
[CITED: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#authentication-and-error-messages]

### Pitfall 4: Incomplete logout cleanup
**What goes wrong:** old role data persists after logout or token revoke.
**Why it happens:** only tokens are cleared; stores/caches/persisted payloads remain.
**How to avoid:** central logout routine clears session store, role cookies, user-scoped caches, and pending request state.
**Warning signs:** user sees previous role content after logout/login switch.
[ASSUMED]

## Code Examples

Verified patterns from official sources:

### Middleware redirect correctness
```typescript
// Source: https://nuxt.com/docs/4.x/api/utils/navigate-to
export default defineNuxtRouteMiddleware((to) => {
  if (to.path !== '/auth') {
    return navigateTo('/auth')
  }
})
```

### Nuxt fetch split: useFetch for setup data, $fetch for event actions
```typescript
// Source: https://nuxt.com/docs/4.x/getting-started/data-fetching
const { data } = await useFetch('/api/me')

async function submitLogin(payload: { email: string; password: string }) {
  return await $fetch('/api/v1/auth/login', {
    method: 'POST',
    body: payload
  })
}
```

### ofetch response error access for deterministic code mapping
```typescript
// Source: https://github.com/unjs/ofetch
try {
  await $fetch('/api/v1/auth/login', { method: 'POST', body: payload })
} catch (error: any) {
  const apiCode = error?.data?.code
  showSafePersianMessage(mapAuthCode(apiCode))
}
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Component-local auth checks per page | Centralized middleware + shared auth client | Nuxt 3/4 composable and middleware-first app patterns [ASSUMED] | Lower auth drift and easier verification of route boundaries. |
| Ad-hoc fetch wrappers per feature | Interceptor-capable fetch client instances | Mature ofetch interceptor model (documented in current README) [CITED: https://github.com/unjs/ofetch] | Consistent token attach/refresh/error behavior. |
| ASVS 4.x references only | ASVS 5.0.0 stable available | ASVS 5.0.0 published and listed as latest stable [CITED: https://owasp.org/www-project-application-security-verification-standard/] | Security mapping should cite versioned controls explicitly. |

**Deprecated/outdated:**
- Using raw backend auth text in UI is outdated for secure auth UX; controlled error mapping is preferred. [CITED: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#authentication-and-error-messages]

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Promise-based single-flight refresh is sufficient without additional request queue metadata. | Architecture Patterns | Could miss corner cases for non-idempotent retried requests. |
| A2 | zod should be introduced for auth validation in this phase. | Standard Stack / Don't Hand-Roll | Adds dependency overhead if team prefers lightweight custom validation. |
| A3 | Existing app persistence/caching surfaces can be fully cleaned by a centralized logout routine without extra plugin hooks. | Common Pitfalls | Potential residual auth artifacts if hidden persistence exists. |
| A4 | Nuxt-auth external module would be less suitable for this custom OTP/role split contract. | Alternatives Considered | Could overlook module capabilities if evaluated deeply later. |

## Open Questions

1. Should refresh retry behavior be zero retry or bounded retry (e.g., 1 retry with short backoff) for transient network failures?
- What we know: API contract defines refresh endpoint and auth failure codes, but not client retry policy. [VERIFIED: docs/API.md]
- What's unclear: desired UX for flaky networks during token refresh.
- Recommendation: lock as product decision during planning; default to one bounded retry with strict timeout and then logout. [ASSUMED]

2. Should tokens be persisted in cookies, localStorage, or memory+cookie split in this repo?
- What we know: current middleware reads `nt_role` cookie and Phase 1 already uses cookie-based route context. [VERIFIED: app/middleware/role-shell.global.ts]
- What's unclear: security/persistence preference for access/refresh tokens in this frontend-only deployment.
- Recommendation: keep access token in memory, store refresh token in secure persistence boundary compatible with current API contract, and document threat tradeoff explicitly. [ASSUMED]

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| node | Nuxt build/dev/test pipeline | yes | v24.14.1 | — |
| npm | package/version management | yes | 11.11.0 | — |
| git | workflow + commit_docs automation | yes | available | — |
| docker | not required for this frontend phase execution | yes | available | not needed |

**Missing dependencies with no fallback:**
- None. [VERIFIED: environment probe]

**Missing dependencies with fallback:**
- None. [VERIFIED: environment probe]

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | Vitest 3.2.2 + Playwright 1.54.1 [VERIFIED: package.json] |
| Config file | `vitest.config.ts`, `playwright.config.ts` [VERIFIED: workspace files] |
| Quick run command | `npm run test:unit` [VERIFIED: package.json] |
| Full suite command | `npm run test:unit && npm run test:e2e` [VERIFIED: package.json] |

### Phase Requirements -> Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| AUTH-01 | OTP request/verify input validation, cooldown, retry lock behavior | unit + e2e | `npm run test:unit -- tests/auth/otp-flow.spec.ts` and `npm run test:e2e -- tests/e2e/auth-otp.spec.ts` | no (Wave 0) |
| AUTH-02 | Nutritionist/admin credential login and role landing | unit + e2e | `npm run test:unit -- tests/auth/credential-auth.spec.ts` and `npm run test:e2e -- tests/e2e/auth-role-login.spec.ts` | no (Wave 0) |
| AUTH-03 | token refresh single-flight, session expiry redirect, logout cleanup | unit + integration | `npm run test:unit -- tests/auth/session-lifecycle.spec.ts` | no (Wave 0) |
| AUTH-04 | namespace guard enforcement on direct entry, refresh, transition | unit + e2e | `npm run test:unit -- tests/auth/role-guard.spec.ts` and `npm run test:e2e -- tests/e2e/role-namespace-guard.spec.ts` | no (Wave 0) |

### Sampling Rate

- **Per task commit:** `npm run test:unit`
- **Per wave merge:** `npm run test:unit && npm run test:e2e`
- **Phase gate:** Full suite green before `/gsd-verify-work`

### Wave 0 Gaps

- [ ] `tests/auth/otp-flow.spec.ts` - covers AUTH-01
- [ ] `tests/auth/credential-auth.spec.ts` - covers AUTH-02
- [ ] `tests/auth/session-lifecycle.spec.ts` - covers AUTH-03
- [ ] `tests/auth/role-guard.spec.ts` - covers AUTH-04
- [ ] `tests/e2e/auth-otp.spec.ts` - OTP journey smoke/e2e
- [ ] `tests/e2e/auth-role-login.spec.ts` - credential role landing e2e
- [ ] `tests/e2e/role-namespace-guard.spec.ts` - direct URL + transition access checks

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | yes | Role-specific auth flow split (OTP for client, credential for nutritionist/admin) + generic auth failures in UI. [VERIFIED: docs/API.md + CITED: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html] |
| V3 Session Management | yes | Access/refresh lifecycle, refresh endpoint usage, forced logout on invalid/revoked token. [VERIFIED: docs/API.md] |
| V4 Access Control | yes | Global namespace middleware deny-by-default and role-root redirects. [VERIFIED: app/middleware/role-shell.global.ts + CITED: https://nuxt.com/docs/4.x/guide/directory-structure/middleware] |
| V5 Input Validation | yes | OTP/mobile/email/password validation at client boundary with deterministic user-safe feedback. [VERIFIED: docs/API.md + ASSUMED] |
| V6 Cryptography | yes | JWT handling only; no client-side cryptographic primitive implementation in this phase. [VERIFIED: docs/API.md + ASSUMED] |

### Known Threat Patterns for this stack

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| User enumeration via detailed login errors | Information Disclosure | Code-based generic Persian error mapping for login/OTP failures. [CITED: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html#authentication-and-error-messages] |
| Refresh token abuse after expiry/revoke | Elevation of Privilege | Refresh failure triggers hard logout and local state purge. [VERIFIED: docs/API.md + ASSUMED] |
| Role namespace probing (`/admin`, `/nutritionist`) | Elevation of Privilege | Global middleware deny-by-default with role prefix checks. [VERIFIED: app/middleware/role-shell.global.ts] |
| Concurrent refresh race causing inconsistent auth state | Tampering | Single-flight refresh guard and atomic token updates. [ASSUMED] |
| Open redirect misuse in auth transitions | Spoofing | Use internal-only route targets with `navigateTo` and never trust external redirect params. [CITED: https://nuxt.com/docs/4.x/api/utils/navigate-to + ASSUMED] |

## Sources

### Primary (HIGH confidence)
- `docs/API.md` - auth endpoints, token lifecycle, logout semantics, error codes, role model.
- `docs/PRD.md` - role auth model, OTP and session expectations.
- `app/middleware/role-shell.global.ts` - existing role namespace enforcement baseline.
- https://nuxt.com/docs/4.x/guide/directory-structure/middleware - middleware behavior, order, server/client run semantics.
- https://nuxt.com/docs/4.x/api/utils/navigate-to - return semantics and redirect correctness in middleware.
- https://nuxt.com/docs/4.x/getting-started/data-fetching - `$fetch`/`useFetch` guidance in Nuxt 4.
- https://github.com/unjs/ofetch - interceptor lifecycle and shared client creation patterns.
- https://pinia.vuejs.org/core-concepts/ - store design and SSR caveats.
- npm registry package metadata for version verification (`nuxt`, `pinia`, `tailwindcss`, `zod`).

### Secondary (MEDIUM confidence)
- https://owasp.org/www-project-application-security-verification-standard/ - ASVS project scope and latest stable version metadata.
- https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html - authentication error handling and anti-enumeration guidance.

### Tertiary (LOW confidence)
- None.

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH - project lock + npm verification + framework docs.
- Architecture: HIGH - aligns with Nuxt middleware/fetch docs and existing codebase boundary.
- Pitfalls: MEDIUM - some implementation-specific race/logout details rely on engineering assumptions.

**Graph context:** `.planning/graphs/graph.json` not present; graph enrichment skipped. [VERIFIED: filesystem check]

**Research date:** 2026-04-23
**Valid until:** 2026-05-23
