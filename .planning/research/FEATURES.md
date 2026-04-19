# Feature Landscape

**Domain:** Nutritionist-Client Management PWA (Persian Market)
**Researched:** 2025-07-14
**Confidence:** MEDIUM — Based on competitive analysis of Nutrium, Practice Better, Healthie, Foodzilla, That Clean Life, Nutritics, and consumer apps (MyFitnessPal, Cronometer). No web search available; findings driven by training data + detailed PRD analysis.

---

## Competitive Context

NutriTrack sits in the **nutritionist practice management** category — a B2B2C tool where the nutritionist is the buyer and the client is the end-user. The major Western competitors are:

| Platform | Strength | Weakness for Persian Market |
|----------|----------|----------------------------|
| Nutrium | Full-featured practitioner platform, client portal | No Persian/RTL, USDA food DB, subscription pricing |
| Practice Better | Scheduling, billing, telehealth, protocols | Overly complex, English-only, no offline |
| Healthie | Telehealth-first, EHR integration, billing | Enterprise-oriented, no Persian support |
| That Clean Life | Beautiful meal plans, recipe database | Recipe-focused (not tracking), English-only |
| Foodzilla | AI-assisted, modern UX | No client tracking, no Persian |
| Nutritics | Deep nutritional analysis (micronutrients) | Desktop-oriented, academic focus |

**Key insight:** No credible Persian-language competitor exists in this space. Iranian nutritionists currently use WhatsApp/Telegram for messaging, Excel/Word for diet plans, and paper for tracking. NutriTrack's primary competition is **not other software** but **the existing manual workflow**. This means table stakes are defined by "what replaces the current manual process" not "feature parity with Western SaaS."

---

## Table Stakes

Features users expect. Missing any of these = nutritionists won't switch from their current WhatsApp+Excel workflow.

### 1. Authentication & Role-Based Access

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Three-role system (Admin, Nutritionist, Client) | Fundamental access control; nutritionists must not see each other's clients | Medium | Row-level isolation is the hard part |
| OTP login for clients (SMS) | Iranian clients expect SMS-based auth; no email culture for consumer apps | Medium | Requires Iranian SMS gateway (Kavenegar/Melipayamak) integration |
| JWT with refresh tokens | Standard session management for PWA | Low | Well-understood pattern |
| Nutritionist-creates-client flow | Mirrors real-world onboarding (client doesn't self-register) | Low | Matches how clinics actually work |

### 2. Client Management

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Client list with search/filter | Nutritionists managing 20-100+ clients need fast lookup | Low | Search by name + mobile |
| Client profile with personal info | Height, weight, DOB, gender — minimum clinical data | Low | Static demographic data |
| Client activation/deactivation | Nutritionists need to manage client lifecycle without data loss | Low | Soft delete pattern |
| Client history overview | Must see all tracking data at a glance to prepare for sessions | Medium | Aggregation across multiple tracking tables |

### 3. Diet Plan Creation & Viewing

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Plan with date range (start/end) | Diet plans are time-bounded in clinical practice | Low | Basic metadata |
| Multi-day plans with daily meals | Core of what a nutritionist delivers; replaces Word documents | High | Deeply nested: Plan → Days → Meals → Options → Items |
| Meal options (client picks one per meal) | Standard practice — give 2-3 options per meal for variety | High | Adds a nesting layer; computation complexity |
| Food item selection from database | Must link to nutritional data for computation | Medium | Food picker modal with search |
| Real-time nutritional totals | Nutritionist MUST see calorie/macro totals while building plan | Medium | Pure computation but critical UX |
| One active plan per client | Clinical rule; prevents confusion | Low | DB constraint + application logic |
| Plan archival & history | Clients and nutritionists need to reference past plans | Low | Status enum, history list |
| Client plan view (mobile) | The core client deliverable — "what do I eat today?" | Medium | Day navigation, meal display, RTL layout |

### 4. Shared Food Database

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Food items with full macro data | Calories, protein, carbs, fat per unit — foundation for plan computation | Low | CRUD with nutritional fields |
| Persian food names & categories | Platform is Persian-only; USDA database is useless | Low | 8 meal-type categories |
| Multiple measurement units | Iranian cooking uses cups, spoons, palm-size, matchbox-size etc. | Low | 12 unit types — domain-specific |
| Search & filtering | Nutritionists building plans need fast food lookup | Medium | Persian full-text search with pg_trgm |
| Shared across nutritionists | All nutritionists contribute to and benefit from a growing DB | Low | Platform-wide resource, created_by tracking |

### 5. Client Daily Tracking

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Food intake logging (which option was eaten) | Clients need to report adherence; nutritionists need to see compliance | Low | Select from plan options or mark skipped |
| Weight tracking with history chart | Weight is the #1 metric nutritionists and clients track | Low | Simple line chart with Shamsi dates |
| Body measurements (waist, hip, etc.) | Standard clinical measurements taken at appointments | Low | 7 body sites, both client and nutritionist can record |
| Water intake tracking | Commonly prescribed alongside diet plans | Low | Tap-to-add glasses, daily total vs target |

### 6. Messaging

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Text messaging (client ↔ nutritionist) | Replaces WhatsApp; must exist for platform to be adopted | Medium | Chat UI, chronological messages |
| Image attachments | Clients send food photos; nutritionists send diagrams | Medium | Upload, display, file storage |
| File attachments (PDF) | Sharing lab results, resources | Medium | Upload with size/type validation |
| Unread message badge | Standard messaging UX expectation | Low | Count query, badge display |

### 7. Persian & RTL Native Experience

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Full RTL layout | Non-negotiable for Persian-only app | Medium | Tailwind RTL plugin, consistent across all components |
| Shamsi (Jalali) calendar | Iranian users don't use Gregorian dates for daily planning | Medium | All date displays, pickers, and inputs must use Shamsi |
| Persian numerals | Expected in a native Persian experience | Low | Number formatting utility |
| Persian error messages & labels | All UI text in natural Persian | Low | But must be reviewed by native speaker |

### 8. Push Notifications

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| New message notification | Essential for messaging adoption; without it, users revert to Telegram | Medium | Web Push via VAPID |
| New diet plan notification | Client must know when their nutritionist assigns a new plan | Low | Triggered on plan creation |
| Permission management (enable/disable types) | Standard mobile UX; users expect notification control | Low | Settings toggle per type |

---

## Differentiators

Features that set NutriTrack apart from competitors (and from the manual workflow). Not expected by users, but create significant value.

### 1. Full Offline Support (Client Side)

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| View diet plan offline | Clients often check their plan in areas with poor connectivity (gyms, kitchens) | High | Service Worker + IndexedDB caching of full plan tree |
| Log tracking data offline | Removes friction from logging — no "I'll do it later" excuse | High | Queue to IndexedDB, sync on reconnect |
| Background sync with conflict resolution | Seamless experience — client never worries about connectivity | High | Exponential backoff, local_id deduplication, last-write-wins |
| Sync status indicator | Transparency builds trust in offline-first UX | Medium | Visual indicator for pending/syncing/synced/failed |
| Offline message viewing & queuing | Cached messages + queued outgoing = complete offline chat | High | Complex state management with sync |

**Why differentiating:** No competitor in the nutritionist space offers genuine offline support. Western platforms assume always-on connectivity. For Iranian users on mobile data with spotty coverage, this is a **major** value-add.

### 2. Medication Prescription & Tracking

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Shared medication database | Nutritionists prescribe supplements and sometimes coordinate with medications | Low | Similar pattern to food database |
| Medication prescription as part of diet plan | Supplements (vitamin D, iron, omega-3) are standard in Iranian nutrition practice | Medium | Linked to plan with dosage/frequency/timing |
| Medication intake checklist (client side) | Clients can track supplement compliance alongside diet | Medium | Pre-populated from prescription, tap-to-mark-taken |
| Medication reminders (push notifications) | Improves supplement compliance rates significantly | Medium | Time-based push notifications from prescription schedule |

**Why differentiating:** Most nutrition platforms (Nutrium, That Clean Life, Foodzilla) don't handle medications at all. This bridges the gap between nutrition and clinical practice, which is very relevant in Iran where nutritionists often recommend supplements as standard practice.

### 3. Food Request System

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Client submits food request | Empowers clients when they can't find an item | Low | Simple form: name + description |
| Nutritionist reviews and approves/rejects | Quality-controlled database growth | Medium | Approval → food creation flow |
| Approved items added to shared database | Crowdsourced Persian food database expansion | Low | Bridges the gap of no pre-existing Persian food DB |

**Why differentiating:** Unique workflow not seen in any competitor. Solves the cold-start problem of a Persian food database — clients help identify local and regional foods that nutritionists can then add with proper nutritional data.

### 4. Multi-Dimensional Tracking Suite

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Sleep tracking (time, duration, quality) | Sleep affects weight loss; nutritionists increasingly assess it | Low | Date + sleep/wake times + quality enum |
| Exercise tracking (activity, duration, calories) | Correlates with dietary goals; provides complete picture | Low | Free-text exercise name + metrics |
| Comprehensive body measurements (6 sites) | Goes beyond weight — waist/hip ratio, body composition trends | Low | Multiple measurement fields per entry |
| Nutritionist-recorded measurements | Both parties can contribute data (at-clinic vs at-home) | Low | `recorded_by` field distinguishes source |

**Why differentiating:** While individual tracking features exist elsewhere, the combination of food + water + sleep + exercise + medication + body measurements in a single client view gives nutritionists a uniquely **holistic** picture. Most platforms offer 2-3 of these.

### 5. Lab Results Upload

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Upload lab results (PDF/JPG/PNG) | Nutritionists need blood work to create proper plans (iron, thyroid, vitamin D, etc.) | Medium | File upload + storage + download |
| Categorized by test type | Quick filtering of blood vs thyroid vs hormone tests | Low | Enum-based categorization |
| Link-based results (external URL) | Some Iranian labs provide online result portals | Low | Alternative to file upload |

**Why differentiating:** Bridges the gap between nutrition and medical data. Competitors typically handle this through separate document management or not at all. Having lab results alongside diet history and body measurements creates a comprehensive clinical record.

### 6. Repeating Day Patterns

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| 7-day cycle that repeats throughout plan period | Massive time saver — nutritionist creates 7 days, system applies across 4+ weeks | Medium | Frontend modulo mapping, backend stores only template days |

**Why differentiating:** Most competitors require creating each day individually or duplicating entire plans. Cycle-based planning matches how nutritionists actually think about weekly meal patterns.

### 7. PWA (Progressive Web App) Distribution

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Installable on home screen (Android + iOS) | No app store needed — instant deployment, no review process | Medium | Web manifest, service worker registration |
| Standalone app experience | Feels like native app without address bar | Low | `display: standalone` in manifest |
| Instant updates (no app store delays) | Bug fixes and features deploy immediately | Low | Service worker update detection |

**Why differentiating:** In the Iranian market, Google Play has restrictions and Apple App Store requires a developer account. PWA sidesteps both distribution channels entirely. Updates bypass store review processes.

---

## Anti-Features

Features to explicitly NOT build. Each has a clear rationale.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| **AI-powered diet recommendations** | Adds massive complexity (ML infrastructure, training data, liability); nutritionists want to be the decision-makers, not have software prescribe for them | Let nutritionists use their expertise; provide data, not prescriptions |
| **Calorie auto-detection from food photos** | Requires ML model, unreliable accuracy, massive engineering effort for marginal value | Manual food logging from plan options is faster and more accurate |
| **Real-time video/voice consultation** | WebRTC infrastructure is complex, costly, and a completely different product domain; Telegram/WhatsApp already serve this well | Keep messaging as text+images+files; clients already use external apps for calls |
| **Payment/subscription billing** | Iranian payment landscape (Shetab, Shaparak) is unique and complex; adds regulatory burden; nutritionists handle payments outside the platform | Out of scope; can add later if demand emerges |
| **Wearable device integration** | Requires device-specific APIs (Fitbit, Apple Health, Google Fit), ongoing maintenance, and many users don't have wearables | Manual tracking is sufficient; lower adoption ceiling isn't worth the integration cost |
| **Desktop-optimized UI** | Splits design effort; mobile-first for clients, and nutritionists in Iran primarily use phones too | Single mobile viewport maximizes quality and development speed |
| **Multi-language / i18n** | Persian-only market; adding i18n infrastructure is premature and adds complexity to every component | Hardcode Persian strings; revisit only if international expansion becomes a goal |
| **Recipe database with cooking instructions** | That Clean Life's territory; recipes are content-heavy and require curation/licensing | Focus on food items with nutritional data, not how to cook them |
| **Grocery list generation** | Nice-to-have but adds complexity to plan engine; low priority for initial launch | Can be derived from plan items later as a read-only feature |
| **Appointment scheduling / calendar** | Separate domain (scheduling SaaS); nutritionists use existing tools for this | Integrate with external calendar apps if needed, don't build one |
| **Social features / community** | Client privacy is paramount; nutritionist-client relationship is private | Keep all interactions 1:1 between assigned nutritionist and client |
| **Gamification (badges, streaks, leaderboards)** | Nutritional adherence gamification can create unhealthy relationships with food; inappropriate for clinical context | Focus on objective tracking data and professional guidance |
| **Detailed micronutrient analysis** | Requires extremely comprehensive food database (100+ nutrients per item); accuracy concerns with non-USDA data | Track macros (calories, protein, carbs, fat) + fiber/sugar/sodium; sufficient for practice |

---

## Feature Dependencies

Critical ordering constraints based on data model and UX dependencies:

```
Authentication (users table)
├── Client Management (requires users)
│   ├── Diet Plan Engine (requires clients + food database)
│   │   ├── Food Logging (requires active plan with meals/options)
│   │   ├── Medication Tracking (requires prescribed medications in plan)
│   │   ├── Exercise Tracking (independent, but exercise recommendations live in plan)
│   │   └── Meal/Medication Reminders (requires plan schedule times)
│   │
│   ├── Body Measurements (requires clients, independent of plan)
│   ├── Weight Tracking (requires clients, independent of plan)
│   ├── Water Tracking (requires clients; target from plan is optional)
│   ├── Sleep Tracking (requires clients, fully independent)
│   ├── Lab Results (requires clients, fully independent)
│   └── Messaging (requires client-nutritionist relationship)
│
├── Food Database (shared resource, no user dependency beyond created_by)
│   ├── Diet Plan Engine (food items used in meal options)
│   └── Food Request System (requires food DB + client-nutritionist relationship)
│
├── Medication Database (shared resource, parallel to food DB)
│   └── Medication Prescription (used in diet plan)
│
├── Offline Support (requires ALL tracking features to be built first)
│   ├── Plan Caching (requires plan viewing)
│   ├── Tracking Queue (requires all log endpoints)
│   └── Message Caching (requires messaging)
│
└── Push Notifications (requires push subscription infrastructure)
    ├── Message Notifications (requires messaging)
    ├── Plan Assignment Notification (requires plan creation)
    ├── Meal Reminders (requires plan with scheduled times)
    ├── Medication Reminders (requires medication prescriptions)
    └── Water Reminders (requires water tracking + target)
```

### Critical Path

The longest dependency chain that gates all subsequent work:

```
Auth → Food DB → Diet Plan Engine → Food Logging → Offline Food Logging → PWA
```

This chain spans the entire project. The Diet Plan Engine is the **bottleneck** — it's the most complex feature and blocks both client tracking (food logs depend on plans) and offline support (must cache plans).

### Parallel-Safe Features (after Plan Engine exists)

These features share no data dependencies and can be built concurrently:
- Water tracking, Sleep tracking, Exercise tracking, Body measurements (all independent tracking)
- Messaging system (independent of plans)
- Lab results upload (independent of plans)
- Food request system (depends only on food DB + relationships)

---

## MVP Recommendation

### Prioritize (Phase 1-3 Essentials — must ship before any user touches the product):

1. **Authentication with all 3 roles** — Foundation; nothing works without it
2. **Food database with Persian search** — Plans can't be built without food items; needs initial seeding
3. **Diet plan engine with full nesting** — The core value proposition; what nutritionists actually pay for
4. **Client registration and management** — Nutritionists need to onboard clients before anything else
5. **Client plan viewing (mobile)** — The single screen clients use most; "what do I eat today?"

### Include in First Usable Release (Phase 4-5 — tracking + communication):

6. **Food logging** — Direct plan adherence tracking; most requested by nutritionists
7. **Weight & body measurement tracking** — The primary outcomes nutritionists measure
8. **Water intake tracking** — Simple, high-frequency interaction; builds daily app habit
9. **Messaging (text + attachments)** — Replaces WhatsApp; critical for adoption
10. **Push notifications (messages + new plans)** — Without this, messaging fails (users won't poll manually)

### Defer (Phase 6 — offline; high value but high complexity):

11. **Full offline support** — Massive value but requires ALL tracking features to be stable first
12. **Sleep tracking** — Useful but lower priority than diet adherence metrics
13. **Exercise tracking** — Useful but lower priority than food/weight tracking
14. **Medication prescription & tracking** — Differentiator but not blocking initial adoption
15. **Lab results upload** — Nice-to-have; nutritionists can use messaging to receive lab results initially
16. **Food request system** — Differentiator but edge case; nutritionists can add foods themselves initially

### Rationale for MVP Ordering

The MVP must answer one question: **"Can a nutritionist create a plan and can a client view it?"** If yes, the platform replaces Word documents. Everything else (tracking, messaging, offline) layers on top.

Food logging and weight tracking come immediately after because they create the **feedback loop** that makes the platform sticky. Without tracking, it's a one-way broadcast (nutritionist → client) and provides no advantage over a PDF sent via Telegram.

Messaging must be in the first usable release because if clients need to switch to WhatsApp to ask their nutritionist a question, the platform loses its "single pane of glass" value.

Offline support is deferred not because it's unimportant but because it's **architecturally expensive** and requires all online features to be stable. Building offline sync on top of buggy online features creates compounding problems.

---

## Feature Prioritization Matrix

| Feature | User Value | Business Value | Complexity | Priority |
|---------|-----------|---------------|------------|----------|
| Diet Plan Engine | ★★★★★ | ★★★★★ | High | **P0 — Ship-blocking** |
| Client Management | ★★★★★ | ★★★★★ | Low | **P0 — Ship-blocking** |
| Food Database | ★★★★★ | ★★★★★ | Medium | **P0 — Ship-blocking** |
| Auth (3 roles + OTP) | ★★★★★ | ★★★★★ | Medium | **P0 — Ship-blocking** |
| Client Plan View | ★★★★★ | ★★★★★ | Medium | **P0 — Ship-blocking** |
| Food Logging | ★★★★☆ | ★★★★★ | Low | **P1 — First release** |
| Weight/Body Tracking | ★★★★☆ | ★★★★☆ | Low | **P1 — First release** |
| Water Tracking | ★★★☆☆ | ★★★★☆ | Low | **P1 — First release** |
| Messaging | ★★★★★ | ★★★★★ | Medium | **P1 — First release** |
| Push (messages + plans) | ★★★★☆ | ★★★★★ | Medium | **P1 — First release** |
| Medication DB & Tracking | ★★★☆☆ | ★★★☆☆ | Medium | **P2 — Second release** |
| Sleep Tracking | ★★☆☆☆ | ★★☆☆☆ | Low | **P2 — Second release** |
| Exercise Tracking | ★★☆☆☆ | ★★☆☆☆ | Low | **P2 — Second release** |
| Lab Results Upload | ★★★☆☆ | ★★☆☆☆ | Medium | **P2 — Second release** |
| Food Request System | ★★☆☆☆ | ★★★☆☆ | Medium | **P2 — Second release** |
| Offline Plan Viewing | ★★★★★ | ★★★★☆ | High | **P3 — Polish release** |
| Offline Tracking Sync | ★★★★☆ | ★★★★☆ | High | **P3 — Polish release** |
| Offline Messaging | ★★★☆☆ | ★★★☆☆ | High | **P3 — Polish release** |
| Meal/Med Reminders | ★★★☆☆ | ★★★☆☆ | Medium | **P3 — Polish release** |
| Super Admin Panel | ★★☆☆☆ | ★★★★★ | Low | **P1 — First release** |
| Repeating Day Pattern | ★★★☆☆ | ★★☆☆☆ | Medium | **P2 — Second release** |

---

## Feature Complexity Notes

### High Complexity Features (warrant extra attention)

1. **Diet Plan Engine** — The deeply nested data model (Plan → Days → Meals → Options → Items) creates complexity in CRUD operations, API design (deep fetches), and frontend state management. The plan builder UI is the single most complex frontend component. Expect 3-4 weeks.

2. **Offline Sync Manager** — Managing a queue in IndexedDB, handling reconnection, deduplicating via local_id, retrying failures with backoff, and showing sync state is a full feature unto itself. Each tracking type adds a sync concern. Testing offline→online transitions is notoriously difficult.

3. **Service Worker Caching** — Caching the plan tree (which is deeply nested JSON), managing cache invalidation (plan updates), and handling the "stale cache + background refresh" pattern requires careful architecture. ETag/If-Modified-Since negotiation adds backend work.

### Low Complexity Features (good candidates for parallel work)

1. **Water Tracking** — Simple counter + daily total. Few edge cases.
2. **Sleep Tracking** — One record per day, upsert, simple computation.
3. **Exercise Tracking** — Free-text log, no plan dependency for the basic version.
4. **Body Measurements** — CRUD with date, multiple fields per record.
5. **Lab Results Upload** — Standard file upload + metadata form.

---

## Sources

- **Competitive platforms (training data, MEDIUM confidence):** Nutrium (nutrium.com), Practice Better (practicebetter.io), Healthie (gethealthie.com), Foodzilla (foodzilla.io), That Clean Life (thatcleanlife.com), Nutritics (nutritics.com)
- **Consumer tracking apps (training data, MEDIUM confidence):** MyFitnessPal, Cronometer, Lose It!, FatSecret
- **PRD analysis (HIGH confidence):** Detailed PRD reviewed line-by-line for feature extraction
- **Technical capabilities (HIGH confidence):** Dexie.js docs via Context7 (offline sync patterns), vite-pwa docs via Context7 (service worker strategies)
- **Persian market context (LOW confidence):** Based on general knowledge of Iranian tech ecosystem; no primary research conducted

### Confidence Notes

- Table stakes categorization: **MEDIUM** — based on competitive analysis of Western platforms and PRD's own goals/non-goals. No direct user research or Persian market validation available.
- Complexity estimates: **MEDIUM** — based on architectural analysis of the data model and technology stack. Actual complexity depends on team skill and edge cases discovered during implementation.
- Anti-features list: **HIGH** — directly sourced from PRD Section 2 (Non-Goals) with additional entries based on domain knowledge.
- Feature dependencies: **HIGH** — derived directly from the data model foreign key relationships in PRD Section 9.
