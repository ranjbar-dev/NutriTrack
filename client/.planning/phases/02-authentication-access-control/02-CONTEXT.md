# Phase 2: Authentication & Access Control - Context

**Gathered:** 2026-04-23
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver role-specific sign-in, session refresh, and route protection for client, nutritionist, and super admin users on top of the completed Phase 1 platform shell. This phase owns authentication and access boundaries only, not business feature flows.

</domain>

<decisions>
## Implementation Decisions

### Auth flow split by role
- **D-01:** Client authentication uses OTP flow only through documented endpoints (`/auth/otp/send`, `/auth/otp/verify`) and lands in client route namespace.
- **D-02:** Nutritionist and super admin authentication use email/password through `/auth/login` and land in role-specific route namespaces.
- **D-03:** The frontend must not introduce alternate auth methods (social, magic links, self-registration) in this phase.

### Session and token lifecycle
- **D-04:** Access token + refresh token lifecycle is implemented via API contract, with refresh handled through a single-flight guard to avoid parallel refresh storms.
- **D-05:** On auth failure (`INVALID_TOKEN`, `TOKEN_REVOKED`, `UNAUTHORIZED`), user is safely logged out and redirected to role-appropriate auth entry.
- **D-06:** Logout uses `/auth/logout` and clears all user-scoped runtime state, persisted state, and role caches.

### Route and access protection
- **D-07:** Route middleware enforces namespace access by role (`/client/**`, `/nutritionist/**`, `/admin/**`) with deny-by-default behavior.
- **D-08:** Shared auth pages remain neutral and do not expose role-private navigation.
- **D-09:** Access guard checks apply consistently on direct URL entry, refresh, and route transitions.

### UX and localization constraints
- **D-10:** All auth copy stays Persian-only, short, and recovery-oriented.
- **D-11:** OTP UX must expose cooldown/retry states and clear validation feedback without leaking sensitive auth internals.
- **D-12:** Numeric entry formatting allows Latin digits where input correctness requires it (OTP/mobile), while display contexts remain Persian-first.

### the agent's Discretion
- Exact form layout composition per role entry screen.
- Internal store/composable file boundaries, as long as session and guard responsibilities stay explicit.
- Retry/backoff timing details for refresh and transient network failures within reasonable UX constraints.

</decisions>

<specifics>
## Specific Ideas

- Keep role onboarding fast: one clear CTA per auth screen.
- Prefer deterministic error mapping from API codes to Persian user-safe messages.
- Preserve Phase 1 visual tone and shell behavior while adding auth state transitions.

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Product and API contracts
- `docs/PRD.md` — authentication model by role, OTP constraints, and session expectations.
- `docs/API.md` — exact endpoint contracts, response envelope, error codes, and role access matrix.

### Planning and requirements
- `.planning/PROJECT.md` — frontend-only scope and core product constraints.
- `.planning/REQUIREMENTS.md` — AUTH-01..AUTH-04 requirement definitions and status.
- `.planning/ROADMAP.md` — Phase 2 goal, dependencies, and success criteria.
- `.planning/STATE.md` — current project position and continuity notes.

### Prior phase foundation
- `.planning/phases/01-platform-foundation/01-CONTEXT.md` — platform and role shell decisions already locked.
- `.planning/phases/01-platform-foundation/01-UI-SPEC.md` — visual and interaction contract baseline.
- `.planning/phases/01-platform-foundation/01-03-SUMMARY.md` — latest completed platform baseline summary.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `app/layouts/auth.vue`: neutral auth layout scaffold from Phase 1.
- `app/layouts/client.vue`: client shell boundary ready for auth-gated entry.
- `app/layouts/nutritionist.vue`: nutritionist shell boundary ready for auth-gated entry.
- `app/layouts/admin.vue`: admin shell boundary ready for auth-gated entry.
- `app/middleware/role-shell.global.ts`: route-level role namespace enforcement baseline.
- `app/stores/platform-pwa.ts`: existing typed store pattern to mirror for auth/session store shape.

### Established Patterns
- Role namespaces are explicitly isolated.
- Persian RTL shell and shared primitives are established and should be reused.
- PWA and cache boundaries are conservative by default.

### Integration Points
- New auth/session composables and stores must plug into existing role layouts and middleware checks.
- API client/interceptor behavior must integrate with route guards and logout cleanup.

</code_context>

<deferred>
## Deferred Ideas

- MFA beyond current role model, social login, or magic links.
- Server-driven RBAC matrix editor or permission-management UI.
- Advanced account recovery beyond current API contract.

</deferred>

---

*Phase: 02-authentication-access-control*
*Context gathered: 2026-04-23*