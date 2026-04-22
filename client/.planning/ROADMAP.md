# Roadmap: NutriTrack Client

## Overview

NutriTrack Client ships as a dependency-first Persian RTL mobile PWA: first establish the Nuxt 4 app shell, design system, and role boundaries; then stabilize authentication and role isolation; then deliver the client's offline-capable daily care loop before layering messaging, nutritionist operations, and super-admin governance. This sequence keeps the highest-risk product promises, especially offline client usage and mobile-first Persian UX, from becoming expensive rewrites later.

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: Platform Foundation** - Establish the Persian RTL mobile app shell, design primitives, and installable PWA baseline.
- [x] **Phase 2: Authentication & Access Control** - Deliver role-specific sign-in, session refresh, and route protection for client, nutritionist, and super admin users.
- [ ] **Phase 3: Client Offline Daily Loop** - Deliver the client's today view, plan access, tracking flows, and offline sync visibility.
- [ ] **Phase 4: Messaging, Notifications & Lab Exchange** - Deliver mobile communication flows, attachments, lab-result exchange, and notification controls.
- [ ] **Phase 5: Nutritionist Workspace & Plan Authoring** - Deliver the nutritionist's client workspace, plan management, catalogue access, and food-request handling.
- [ ] **Phase 6: Admin Governance** - Deliver super-admin stats, nutritionist management, and elevated shared catalogue control.

## Phase Details

### Phase 1: Platform Foundation
**Goal**: Users can open an installable Persian RTL mobile shell that establishes the visual, navigational, and PWA foundation for every role.
**Depends on**: Nothing (first phase)
**Requirements**: PLAT-01, PLAT-02, PLAT-03
**Success Criteria** (what must be TRUE):
  1. User can open a Persian-only RTL mobile app shell with role-aware navigation structure for client, nutritionist, and super admin areas.
  2. User sees consistent Persian typography, Persian numerals where appropriate, Jalali-aware date presentation, and safe-area mobile layout behavior across shared UI primitives.
  3. User can install the app as a PWA and sees a clear in-app update prompt when a newer client version is available.
**Plans**: 4 plans
Plans:
- [x] 01-01-PLAN.md - Baseline platform wiring, conservative PWA config, and test harness
- [x] 01-02-PLAN.md - Persian RTL design tokens, locale formatting, and shared shell primitives
- [x] 01-03-PLAN.md - Install/update/connectivity banners and cache-boundary regression guards
- [x] 01-04-PLAN.md - Role shell isolation routes, layouts, and middleware enforcement
**UI hint**: yes

### Phase 2: Authentication & Access Control
**Goal**: Users can sign in through the correct role-specific flow, keep stable sessions, and access only the surfaces allowed to them.
**Depends on**: Phase 1
**Requirements**: AUTH-01, AUTH-02, AUTH-03, AUTH-04
**Success Criteria** (what must be TRUE):
  1. Client can request and verify an OTP with the documented mobile-based flow and land in the client experience.
  2. Nutritionist and super admin can log in with email and password and land in their correct workspaces.
  3. Authenticated user remains signed in across refreshes through token refresh and is redirected safely when the session expires or they log out.
  4. User cannot access pages or data outside the routes allowed for their role and user identity.
**Plans**: 4 plans
Plans:
- [x] 02-01-PLAN.md - Auth core infrastructure, typed session store, and refresh orchestration
- [x] 02-02-PLAN.md - Client OTP flow screens and OTP verification lifecycle
- [x] 02-03-PLAN.md - Nutritionist and admin credential auth flows with secure error handling
- [x] 02-04-PLAN.md - Role route guards, session bootstrap, and logout/reset enforcement
**UI hint**: yes

### Phase 3: Client Offline Daily Loop
**Goal**: Client users can understand today's diet work, review plans, and record daily adherence even with unstable connectivity.
**Depends on**: Phase 2
**Requirements**: CLNT-01, CLNT-02, CLNT-03, TRCK-01, TRCK-02, TRCK-03, OFFL-01, OFFL-02, OFFL-03
**Success Criteria** (what must be TRUE):
  1. Client can open a Today view showing the active plan, pending daily actions, water target, and current sync state.
  2. Client can read the full active plan and available archived plan history without losing the context of which plan is current.
  3. Client can log food intake, water, sleep, exercise, medication intake, and body measurements from mobile-friendly Persian flows.
  4. Client can review recent tracking history and lightweight progress summaries from the data available in v1.
  5. While offline, client can read essential recent plan data and queue supported tracking writes with visible synced, retrying, or failed states plus reconnect or manual retry behavior.
**Plans**: TBD
**UI hint**: yes

### Phase 4: Messaging, Notifications & Lab Exchange
**Goal**: Client and nutritionist users can exchange messages, files, reminders, and lab results through coherent mobile communication flows.
**Depends on**: Phase 3
**Requirements**: MSG-01, MSG-02, NOTF-01, NOTF-02, LAB-01
**Success Criteria** (what must be TRUE):
  1. Client and nutritionist can read conversation history with unread state and stable refresh behavior suitable for polling-based chat.
  2. Client and nutritionist can send Persian text messages and supported file attachments from mobile conversation screens.
  3. Authenticated user can subscribe or unsubscribe from push notifications on supported devices and manage notification preferences for reminder and message categories.
  4. Client and nutritionist can upload, view, and access lab results using the documented file or link-based flows.
**Plans**: TBD
**UI hint**: yes

### Phase 5: Nutritionist Workspace & Plan Authoring
**Goal**: Nutritionists can manage clients and author care plans from a mobile-first operational workspace without inheriting client offline complexity.
**Depends on**: Phase 4
**Requirements**: NUTR-01, NUTR-02, NUTR-03, NUTR-04, CAT-01, CAT-02, CAT-03
**Success Criteria** (what must be TRUE):
  1. Nutritionist can browse, search, and filter their client roster, then open a client profile with identity details, current plan summary, tracking history, messages, lab results, and archived plans.
  2. Nutritionist can create, edit, and manage a client's diet-plan period, metadata, days, meals, meal options, meal items, exercises, and prescriptions from the frontend.
  3. Nutritionist and eligible users can search and view shared foods and medications needed for plan authoring and medication workflows.
  4. Client can submit a food request and nutritionist can review, approve, or reject that request from the frontend.
**Plans**: TBD
**UI hint**: yes

### Phase 6: Admin Governance
**Goal**: Super admins can oversee platform health, nutritionist accounts, and shared catalogue governance from mobile-compatible admin screens.
**Depends on**: Phase 5
**Requirements**: ADMIN-01, ADMIN-02
**Success Criteria** (what must be TRUE):
  1. Super admin can view platform stats from mobile-compatible admin screens.
  2. Super admin can create, update, activate, and deactivate nutritionist accounts from the frontend.
  3. Super admin can manage shared food and medication catalogues using elevated admin endpoints.
**Plans**: TBD
**UI hint**: yes

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Platform Foundation | 4/4 | Completed | 2026-04-22 (01-01, 01-02, 01-04, 01-03) |
| 2. Authentication & Access Control | 4/4 | Completed | 2026-04-23 (02-01, 02-02, 02-03, 02-04) |
| 3. Client Offline Daily Loop | 0/TBD | Not started | - |
| 4. Messaging, Notifications & Lab Exchange | 0/TBD | Not started | - |
| 5. Nutritionist Workspace & Plan Authoring | 0/TBD | Not started | - |
| 6. Admin Governance | 0/TBD | Not started | - |