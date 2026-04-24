# NutriTrack Client

## What This Is

NutriTrack Client is a shipped mobile-first Persian RTL Progressive Web App for nutritionists' clients and operational staff to manage diet plans, tracking, messaging, and nutrition workflows against the NutriTrack REST API. v1.0 delivers the complete core workflow: installable shell, role-based auth, client offline daily tracking, messaging and lab exchange, nutritionist workspace and plan authoring, and super-admin catalogue and account governance. This implementation is frontend-only and built on Nuxt 4, Tailwind CSS 4, and Pinia.

## Core Value

Clients and nutritionists can reliably complete the core nutrition workflow on mobile in Persian, even with unstable connectivity for client-side usage.

## Requirements

### Validated

- ✓ Deliver the NutriTrack client as a Nuxt 4 PWA with Persian-only RTL mobile UX — v1.0
- ✓ Implement the frontend against the documented NutriTrack REST API without backend work in this project scope — v1.0
- ✓ Support the core client and nutritionist workflows needed for v1, including auth, plans, tracking, messaging, food requests, and notifications — v1.0
- ✓ PLAT-01: Persian-only RTL mobile app shell with role-aware navigation — v1.0 (Phase 1)
- ✓ PLAT-02: Installable PWA with in-app update prompt — v1.0 (Phase 1)
- ✓ PLAT-03: Persian typography, numerals, Jalali dates, mobile safe-area — v1.0 (Phase 1)
- ✓ AUTH-01: Client OTP authentication flow — v1.0 (Phase 2)
- ✓ AUTH-02: Nutritionist and super-admin email/password login — v1.0 (Phase 2)
- ✓ AUTH-03: Session persistence with token refresh and safe logout — v1.0 (Phase 2)
- ✓ AUTH-04: Role-based route protection — v1.0 (Phase 2)
- ✓ CLNT-01: Today view with active plan, pending actions, water target, sync state — v1.0 (Phase 3)
- ✓ CLNT-02: Full active plan readability across days, meals, exercises, prescriptions — v1.0 (Phase 3)
- ✓ CLNT-03: Archived plan history without losing active plan context — v1.0 (Phase 3)
- ✓ TRCK-01: Food intake logging with mobile-friendly entry — v1.0 (Phase 3)
- ✓ TRCK-02: Water, sleep, exercise, medication, body measurement logging — v1.0 (Phase 3)
- ✓ TRCK-03: Tracking history and progress summaries — v1.0 (Phase 3)
- ✓ OFFL-01: Offline plan and essential data read — v1.0 (Phase 3)
- ✓ OFFL-02: Offline tracking write queue with durable local IDs and sync states — v1.0 (Phase 3)
- ✓ OFFL-03: Sync state visibility with reconnect and manual retry — v1.0 (Phase 3)
- ✓ MSG-01: Conversation history with unread state and polling — v1.0 (Phase 4)
- ✓ MSG-02: Text messages and file attachments from mobile — v1.0 (Phase 4)
- ✓ NOTF-01: Push notification subscribe/unsubscribe — v1.0 (Phase 4)
- ✓ NOTF-02: Notification preference management — v1.0 (Phase 4)
- ✓ LAB-01: Lab result upload and access flows — v1.0 (Phase 4)
- ✓ NUTR-01: Client roster browse, search, and filter — v1.0 (Phase 5)
- ✓ NUTR-02: Client profile with plan summary, tracking history, messages, labs, archived plans — v1.0 (Phase 5)
- ✓ NUTR-03: Diet plan period and metadata management — v1.0 (Phase 5)
- ✓ NUTR-04: Plan days, meals, items, exercises, prescriptions management — v1.0 (Phase 5)
- ✓ CAT-01: Food and category search for plan authoring — v1.0 (Phase 5)
- ✓ CAT-02: Medication search for plan authoring — v1.0 (Phase 5)
- ✓ CAT-03: Food request submission and nutritionist moderation — v1.0 (Phase 5)
- ✓ ADMIN-01: Super-admin platform stats and nutritionist account management — v1.0 (Phase 6)
- ✓ ADMIN-02: Elevated shared catalogue governance — v1.0 (Phase 6)

### Active

*(All v1 requirements shipped. See v2 requirements below for next milestone scope.)*

- [ ] **EXP-01**: Client sync center with queue history and richer recovery controls (v2)
- [ ] **EXP-02**: Richer client progress storytelling and adherence narratives (v2)
- [ ] **EXP-03**: Nutritionist reusable plan defaults, drafts, and authoring accelerators (v2)
- [ ] **ADV-01**: Realtime chat transport (v2)
- [ ] **ADV-02**: Client personal custom food authoring beyond request workflow (v2)
- [ ] **ADV-03**: Advanced analytics dashboards and admin reporting (v2)

### Out of Scope

- Backend services, database design, and infrastructure changes — the request is explicitly client-side only.
- Desktop-first layouts — the PRD defines a mobile-only experience.
- Multi-language support — the product is Persian-only for v1.
- AI diet recommendations, payments, and wearable integrations — explicitly excluded in the PRD.

## Context

- v1.0 shipped 2026-04-24: complete Persian RTL mobile PWA covering all 6 phases, 27 plans, 30 requirements.
- Stack: Nuxt 4, Tailwind CSS 4, Pinia, Vitest — frontend-only against NutriTrack REST API.
- Offline support implemented for client tracking flows using local persistence and sync recovery via offline queue with durable local IDs.
- Admin surfaces (nutritionist, super-admin) are online-first as designed.
- Mobile RTL admin UX walkthrough deferred to v1.1 validation cycle (documented human gate).
- API integration follows `docs/API.md` exactly (JWT auth, OTP login, REST resources, pagination, file uploads, role-based permissions).
- Product behavior follows `docs/PRD.md` (Persian-only RTL UX, client offline support, role boundaries, mobile viewport prioritization).
- v2 scope: realtime chat, richer sync center, advanced analytics, plan authoring accelerators, custom food authoring.

## Constraints

- **Tech stack**: Nuxt 4, Tailwind CSS 4, and Pinia — explicitly requested for the client implementation.
- **Scope**: Frontend only — backend endpoints are consumed as documented, not changed here.
- **Product**: Persian-only RTL mobile PWA — defined by the PRD and not negotiable for v1.
- **Offline**: Client-side offline capture and sync are required — core product behavior depends on it.
- **Integration**: API shapes must match `docs/API.md` exactly — the client is being built against this contract.
- **Design**: UI should use the installed UI Pro Max skill direction — to avoid generic layouts and preserve intentional mobile design.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Build only the client application in this workspace | The user explicitly scoped work to the client side | ✓ Good — clean separation maintained throughout |
| Use Nuxt 4 with Tailwind CSS 4 and Pinia | The user fixed the frontend stack up front | ✓ Good — stack served all 6 phases without friction |
| Treat the API document as the integration contract | The workspace already provides a detailed backend reference | ✓ Good — API contract tests caught shape mismatches early |
| Optimize for Persian RTL mobile UX before desktop concerns | The PRD defines mobile-only Persian usage as the primary product shape | ✓ Good — mobile-first discipline held across all phases |
| Include offline-first client architecture in v1 planning | Offline client tracking is a core product promise, not a later enhancement | ✓ Good — offline queue with durable local IDs proved stable |
| Use UI Pro Max as a design input for frontend phases | The user requested that skill be part of the implementation approach | ✓ Good — prevented generic admin template patterns |
| Keep offline durability bounded to client experience only | Admin and nutritionist surfaces are operational and online-first | ✓ Good — simplified Phase 4+ significantly |
| Sequence roadmap: platform → auth → client offline → messaging → nutritionist → admin | Highest-risk promises first; lower dependencies first | ✓ Good — no phase required rework of an earlier phase |
| Use deny-by-default role namespace guard via global middleware | Role isolation is a hard requirement; guards should fail-safe | ✓ Good — no cross-role data leakage during testing |
| Auth failure tokens (INVALID_TOKEN, TOKEN_REVOKED, UNAUTHORIZED) force logout | Security requirement — invalid sessions must not persist | ✓ Good — consistent auth error handling across all roles |
| Admin read-only client visibility enforced at type level | Compile-time safety for read-only constraint | ✓ Good — caught read/write misuse before runtime |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? -> Move to Out of Scope with reason
2. Requirements validated? -> Move to Validated with phase reference
3. New requirements emerged? -> Add to Active
4. Decisions to log? -> Add to Key Decisions
5. "What This Is" still accurate? -> Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check -> still the right priority?
3. Audit Out of Scope -> reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-04-24 after v1.0 milestone*