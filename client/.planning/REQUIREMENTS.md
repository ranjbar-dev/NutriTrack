# Requirements: NutriTrack Client

**Defined:** 2026-04-22
**Core Value:** Clients and nutritionists can reliably complete the core nutrition workflow on mobile in Persian, even with unstable connectivity for client-side usage.

## v1 Requirements

### Platform

- [x] **PLAT-01**: User sees a Persian-only RTL mobile app shell with role-aware navigation for client, nutritionist, and super admin areas.
- [x] **PLAT-02**: User can install the app as a PWA and receives a clear in-app update prompt when a newer version is available.
- [x] **PLAT-03**: User sees Persian typography, Persian numerals where appropriate, Jalali-aware date presentation, and mobile safe-area handling across core flows.

### Authentication

- [x] **AUTH-01**: Client can request and verify an OTP using the documented mobile-based authentication flow.
- [x] **AUTH-02**: Nutritionist and super admin can log in with email and password.
- [x] **AUTH-03**: Authenticated user keeps a valid session across refreshes through token refresh and is redirected safely when the session expires or logs out.
- [x] **AUTH-04**: Authenticated user is routed only to pages allowed for their role and user identity.

### Client Experience

- [x] **CLNT-01**: Client can open a Today view that surfaces the active plan, pending daily actions, water target, and sync status.
- [x] **CLNT-02**: Client can view the full active diet plan with readable days, meals, options, exercise recommendations, prescriptions, and notes.
- [x] **CLNT-03**: Client can view archived or historical plan information made available by the API without losing context of the active plan.

### Tracking

- [x] **TRCK-01**: Client can log food intake against the documented tracking API with mobile-friendly entry flows.
- [x] **TRCK-02**: Client can log water intake, sleep, exercise, medication intake, and body measurements from the mobile UI.
- [x] **TRCK-03**: Client can review recent tracking history and lightweight progress summaries for the tracked data available in v1.

### Offline and Sync

- [x] **OFFL-01**: Client can read the active plan and essential recent client data while offline.
- [x] **OFFL-02**: Client write actions for supported tracking flows are queued locally with durable local IDs and visible sync states.
- [x] **OFFL-03**: Client can see when queued entries are synced, retrying, or failed, and the app retries sync on reconnect or manual retry.

### Messaging and Notifications

- [x] **MSG-01**: Client and nutritionist can read their conversation history with unread state and stable polling behavior.
- [x] **MSG-02**: Client and nutritionist can send text messages and supported file attachments from mobile screens.
- [x] **NOTF-01**: Authenticated user can subscribe or unsubscribe from push notifications on supported devices.
- [x] **NOTF-02**: Authenticated user can view and update notification preferences for supported reminder and message categories.

### Nutritionist Workspace

- [x] **NUTR-01**: Nutritionist can browse, search, and filter their client roster from a mobile-friendly workspace.
- [x] **NUTR-02**: Nutritionist can open a client profile that shows identity details, current plan summary, tracking history, messages, lab results, and archived plans.
- [x] **NUTR-03**: Nutritionist can create, edit, and manage a client's diet plan period and metadata.
- [x] **NUTR-04**: Nutritionist can manage plan days, meals, meal options, meal items, exercises, and prescriptions from the frontend.

### Catalogue and Requests

- [x] **CAT-01**: Nutritionist and eligible users can search and view foods and categories needed for plan authoring and client workflows.
- [x] **CAT-02**: Nutritionist and eligible users can search and view medications needed for plan authoring and medication workflows.
- [x] **CAT-03**: Client can submit a food request and nutritionist can review, approve, or reject food requests from the frontend.

### Lab Results

- [x] **LAB-01**: Client and nutritionist can upload, view, and access lab results using the documented file or link-based flows.

### Super Admin

- [x] **ADMIN-01**: Super admin can view platform stats and manage nutritionist accounts from mobile-compatible admin screens.
- [x] **ADMIN-02**: Super admin can manage shared food and medication catalogues with the elevated admin endpoints.

## v2 Requirements

### Experience Enhancements

- **EXP-01**: Client can inspect a detailed sync center with queue history and richer recovery controls.
- **EXP-02**: Client sees richer progress storytelling, adherence narratives, and deeper trend analysis.
- **EXP-03**: Nutritionist gets reusable plan defaults, drafts, and more advanced plan-authoring accelerators.

### Advanced Capabilities

- **ADV-01**: Chat uses realtime transport instead of polling.
- **ADV-02**: Client can author personal custom foods beyond the curated request workflow.
- **ADV-03**: Platform offers advanced analytics dashboards and deeper admin reporting.

## Out of Scope

| Feature | Reason |
|---------|--------|
| Backend API implementation | This workspace is scoped to frontend/client-side code only. |
| Desktop-optimized layouts | The PRD defines a mobile-only product shape for v1. |
| Multi-language support | The product is Persian-only and does not need i18n infrastructure in v1. |
| AI recommendations or food-photo recognition | Explicit PRD non-goals and a scope risk before core workflow validation. |
| Wearable integrations | Explicit PRD non-goal with unnecessary integration complexity for v1. |
| Payments or billing | Outside the product scope described in the PRD. |
| Broad offline support for nutritionist or admin surfaces | Offline support is required only for client-side flows in the PRD. |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| PLAT-01 | Phase 1 | Completed (01-04) |
| PLAT-02 | Phase 1 | Completed (01-01) |
| PLAT-03 | Phase 1 | Completed (01-02) |
| AUTH-01 | Phase 2 | Completed (02-02) |
| AUTH-02 | Phase 2 | Completed (02-03) |
| AUTH-03 | Phase 2 | Completed (02-01, 02-04) |
| AUTH-04 | Phase 2 | Completed (02-04) |
| CLNT-01 | Phase 3 | Pending |
| CLNT-02 | Phase 3 | Pending |
| CLNT-03 | Phase 3 | Pending |
| TRCK-01 | Phase 3 | Pending |
| TRCK-02 | Phase 3 | Pending |
| TRCK-03 | Phase 3 | Pending |
| OFFL-01 | Phase 3 | Pending |
| OFFL-02 | Phase 3 | Pending |
| OFFL-03 | Phase 3 | Pending |
| MSG-01 | Phase 4 | Complete |
| MSG-02 | Phase 4 | Complete |
| NOTF-01 | Phase 4 | Complete |
| NOTF-02 | Phase 4 | Complete |
| NUTR-01 | Phase 5 | Complete (05-01, 05-02) |
| NUTR-02 | Phase 5 | Complete (05-01, 05-02) |
| NUTR-03 | Phase 5 | Complete (05-01, 05-03) |
| NUTR-04 | Phase 5 | Complete (05-01, 05-03, 05-04) |
| CAT-01 | Phase 5 | Complete (05-01, 05-04) |
| CAT-02 | Phase 5 | Complete (05-01, 05-04) |
| CAT-03 | Phase 5 | Complete (05-01, 05-05) |
| LAB-01 | Phase 4 | Complete |
| ADMIN-01 | Phase 6 | Complete (06-01, 06-02, 06-03) |
| ADMIN-02 | Phase 6 | Complete (06-01, 06-04) |

**Coverage:**
- v1 requirements: 30 total
- Mapped to phases: 30
- Unmapped: 0 ✓

---
*Requirements defined: 2026-04-22*
*Last updated: 2026-04-24 after Phase 6 completion*