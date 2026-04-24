---
phase: 03
slug: client-offline-daily-loop
status: draft
shadcn_initialized: false
preset: none
created: 2026-04-23
---

# Phase 03 - UI Design Contract

> Visual and interaction contract for Phase 3 (Client Offline Daily Loop). Generated for planner and executor consumption.

---

## Design System

| Property | Value |
|----------|-------|
| Tool | none (existing custom Nuxt 4 + Tailwind 4 token system) |
| Preset | not applicable |
| Component library | project-wrapped Vue components (`AppShell`, `InlineNotice`, `EmptyState`, `ErrorState`, `ConnectivityBanner`) |
| Icon library | lucide-vue-next |
| Font | Vazirmatn (Persian-first) |

### Locked Direction (Phase 1 + Phase 2 + Phase 3 context)

- Tone remains clinical-minimal, calm, and trust-first.
- Persian-only copy across client daily loop, tracking, sync, and recovery states.
- Mobile-first RTL remains default; desktop keeps centered mobile canvas.
- Offline support is client-only and limited to plan + tracking domains defined in requirements.
- Sync state is a first-class UX element and must be visible without entering settings.

---

## Spacing Scale

Declared values (multiples of 4 only):

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Inline icon-label gap, status dot spacing |
| sm | 8px | Dense controls, quick action chips, segmented toggles |
| md | 16px | Default spacing between cards/inputs |
| lg | 24px | Card padding, section spacing |
| xl | 32px | Major section separation |
| 2xl | 48px | Hero-to-content separation on Today view |
| 3xl | 64px | Large page breathing room on tall screens |

Exceptions:
- Minimum interactive target: 44px.
- Primary tracking CTA buttons: 48px height.
- Quick-add pill controls (water/meal actions): minimum 44px height.
- Sticky sync strip (when visible): 40px minimum height, non-modal.

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
- Use Persian numerals for display values and summaries.
- Keep Latin digits only where exact input correctness is required.
- Sync timestamps shown to users must be Persian-formatted time/date.

---

## Color

| Role | Value | Usage |
|------|-------|-------|
| Dominant (60%) | #F4F7F8 | App background and base surfaces |
| Secondary (30%) | #FFFFFF | Cards, grouped sections, sheets |
| Accent (10%) | #0F6B7A | Primary CTA, selected tracking tab, manual sync trigger |
| Destructive | #B63D3D | Failed sync state emphasis and destructive confirmations |

Accent reserved for:
- Primary action per viewport (for example: "ثبت وعده", "ثبت آب", "همگام سازی دوباره").
- One active filter/segment in tracking entry flows.
- Focus-visible ring for active controls.
- Single high-priority sync action in failure state (not all status chips).

Sync semantic tones:
- Queued: warning surface (amber tint), non-destructive.
- Syncing: info surface (muted blue/teal tint), animated only once per state transition.
- Synced: success surface (muted green tint), no persistent pulse.
- Failed: destructive tone using `--color-danger` plus explicit recovery action.

---

## Copywriting Contract

| Element | Copy |
|---------|------|
| Primary CTA (today) | ثبت پیگیری امروز |
| Primary CTA (sync retry) | همگام سازی دوباره |
| Empty state heading (today plan) | برنامه فعالی برای امروز پیدا نشد |
| Empty state body (today plan) | برای دریافت برنامه، اینترنت را وصل کنید و صفحه را تازه سازی کنید. |
| Empty state heading (tracking history) | هنوز ثبتی ندارید |
| Empty state body (tracking history) | اولین مورد را از بخش امروز ثبت کنید تا اینجا نمایش داده شود. |
| Error state (sync generic) | همگام سازی کامل نشد. اتصال اینترنت را بررسی کنید و دوباره تلاش کنید. |
| Error state (entry failed) | این ثبت به سرور نرسید. می توانید دوباره تلاش کنید یا بعدا ارسال کنید. |
| Destructive confirmation | حذف ثبت آفلاین: این مورد از صف ارسال حذف می شود. مطمئن هستید؟ |

Trust/recovery tone rules:
- Always use "problem + next step" in one short block.
- Never imply user fault; explain system state clearly.
- For queued items, reassure durability before asking for action.
- For failed items, always provide one primary recovery action and one safe fallback.
- Avoid technical terms like payload, conflict, endpoint in user-facing Persian copy.

---

## Registry Safety

| Registry | Blocks Used | Safety Gate |
|----------|-------------|-------------|
| shadcn official | none | not required |
| third-party | none | not applicable |

---

## Phase 3 Surface Inventory

Required client surfaces:

1. Today home (`/client`)
- Active plan snapshot with day/meal grouping.
- Pending actions summary (food, water, sleep, exercise, medication, body).
- Water target progress and quick-add actions.
- Global sync snapshot (queued/syncing/synced/failed counts).

2. Active plan detail (`/client/plan`)
- Full readable plan: days, meals, options, exercises, prescriptions, notes.
- Clear active-vs-history context badge.

3. Plan history (`/client/history/plans`)
- Archived plans list with period and status.
- Entry to view plan details without confusing active plan state.

4. Tracking entry flows (`/client/tracking/*` or sheet-based entry)
- Food intake entry.
- Water log quick entry.
- Sleep entry.
- Exercise entry.
- Medication entry.
- Body measurement entry.

5. Tracking history and lightweight progress (`/client/history/tracking`)
- Recent logs with status chips.
- Simple progress blocks (water target completion, streak-like summaries where available).

6. Sync/recovery touchpoints
- Always-visible lightweight sync indicator in client shell.
- Failed-item list access point with retry controls.
- Manual sync trigger visible when offline queue exists.

---

## Offline State UX Contract

### State visibility model

Every client write must map to one visible sync state:
- `queued`
- `syncing`
- `synced`
- `failed`

State visibility rules:
- Show per-entry state chip in history list rows.
- Show aggregate state in Today header strip.
- Do not hide failures behind a generic offline banner.
- Maintain visibility until state changes; do not auto-dismiss failed state.

### Global state matrix

| Connectivity | Queue | Header Strip | Entry Chips | Primary Guidance |
|--------------|-------|--------------|-------------|------------------|
| Online | 0 pending | Last sync time + "همگام" status | Synced only | No interruption |
| Online | pending > 0 | "در حال همگام سازی" with counter | queued/syncing | Keep user flow non-blocking |
| Offline | pending >= 0 | "آفلاین" + pending count | queued persists | Reassure data is saved locally |
| Reconnected + failures | failed > 0 | "نیاز به بررسی" + retry CTA | failed visible | Route user to retry list |

### Offline read contract

While offline, client can still access:
- Active plan snapshot and essential plan sections.
- Recent tracking history cached locally.
- Water target and daily progress derived from cached + queued writes.

When stale:
- Show "آخرین به روزرسانی" timestamp.
- If data is older than expected and online is available, prompt refresh non-modally.

### Retry and recovery contract

- Auto-retry on reconnect and app-open for queued/failed entries.
- Exponential backoff up to 3 attempts for failed sync (from PRD).
- After max attempts, keep state `failed` with explicit manual retry CTA.
- Manual retry supports:
  - Retry single entry.
  - Retry all failed entries.
- Optional destructive action: remove failed queued item with confirmation copy.

---

## Sync Indicator Contract

### Indicator hierarchy

1. Global strip (highest visibility)
- Location: below `ConnectivityBanner` in client layout, above main content.
- Content: state label + pending/failed counters + optional CTA.
- Height: compact, single-row on mobile.

2. Section-level summary
- Location: top of tracking history and Today pending block.
- Content: counts by state and last sync time.

3. Row-level chip
- Location: each tracking item in recent list.
- Content: one-word state (`در صف`, `در حال ارسال`, `همگام شد`, `نیاز به ارسال دوباره`).

### Indicator behavior rules

- `queued`: shown immediately after local write success.
- `syncing`: shown only during active network attempt.
- `synced`: shown after API success; may downgrade to subtle style after 10 seconds.
- `failed`: persistent until retry success or explicit remove.

### Motion rules for sync

- Use a single transition (120-180ms) for status change.
- No infinite spinner loops for row-level chips.
- Only global syncing strip may show subtle progress animation, disabled under reduced-motion.

---

## Tracking Entry Patterns (Persian Mobile)

### Cross-flow principles

- One-thumb friendly: primary action reachable in lower viewport.
- Fast path first: prefilled date/time defaults to now and editable.
- Keep form density low: progressive disclosure for optional notes.
- Confirm local save instantly before network result.

### Entry pattern by domain

1. Food log
- Pattern: meal-context entry from Today meal card.
- Primary fields: meal context, selected option or manual name, quantity/unit.
- Secondary fields: macros/calories and notes.
- Submission feedback: immediate queued/synced chip on created item.

2. Water log
- Pattern: quick-add pills (for example 250ml, 500ml) plus custom amount.
- Daily progress ring/bar on same screen.
- After add: remain on screen for repeated one-tap entries.

3. Sleep log
- Pattern: paired time inputs (sleep start, wake end) with computed duration preview.
- Validation: wake must be after sleep with overnight support.

4. Exercise log
- Pattern: exercise name + duration as required core; calories optional.
- Optional notes collapsed by default.

5. Medication log
- Pattern: choose prescribed medication when available; manual name fallback required by API.
- Dosage optional but encouraged when prescribed context exists.

6. Body measurement log
- Pattern: card grid inputs grouped by body zones.
- Requirement: allow partial submission because API fields are optional.
- Highlight last recorded value for trust and continuity.

---

## Data Trust and Recovery Copy Rules

### Required Persian microcopy intents

- Durability reassurance after local write:
  - "ثبت شد و پس از اتصال اینترنت همگام می شود."
- Active sync in progress:
  - "در حال ارسال اطلاعات ثبت شده..."
- Partial sync failure:
  - "بخشی از اطلاعات ارسال نشد. لطفا دوباره تلاش کنید."
- Manual retry success:
  - "همگام سازی با موفقیت انجام شد."
- Persistent failure handoff:
  - "ارسال این مورد انجام نشد. می توانید دوباره تلاش کنید یا بعدا اقدام کنید."

### Prohibited copy patterns

- Blame language (for example: "شما اشتباه وارد کردید") for sync/network states.
- Technical backend jargon (for example: 422, payload, timeout code).
- Ambiguous state words without action (for example: "خطا").

### Conflict and overwrite communication

Given last-write-wins strategy:
- Do not expose conflict internals by default.
- If a value changed after sync reconciliation, show concise informational notice:
  - "آخرین تغییر ثبت شده نمایش داده شد."
- Offer user path to review latest value in history list.

---

## API and State Mapping Contract

Tracking endpoint mapping:
- Food -> `POST /tracking/food`
- Water -> `POST /tracking/water`
- Sleep -> `POST /tracking/sleep`
- Exercise -> `POST /tracking/exercise`
- Medication -> `POST /tracking/medication`
- Body -> `POST /tracking/body`
- Bulk replay -> `POST /tracking/sync`

Mapping rules:
- Every write request must include `local_id` generated client-side UUID.
- UI state key is local entry identity (`local_id`) not server ID.
- `synced_at` user-visible derivation should come from local sync metadata.
- Bulk sync error list must map back to row-level failed entries.

---

## Component Contract (Phase 3 additions)

Reuse required from existing components:
- `AppShell`
- `ConnectivityBanner`
- `InlineNotice`
- `EmptyState`
- `ErrorState`

New components to establish in Phase 3 (names can vary, behavior is fixed):
- `SyncStatusStrip` (global client sync summary)
- `SyncStateChip` (queued/syncing/synced/failed)
- `TodayPlanSnapshotCard`
- `PendingActionsCard`
- `WaterQuickAdd`
- `TrackingEntrySheet` per domain
- `FailedSyncList` with single and bulk retry actions

---

## Accessibility and RTL Contract

- All state updates announce via polite live regions where needed.
- Failed sync notices use `role="alert"`.
- Touch targets >=44px for every quick action.
- Ensure RTL layout for labels and summaries; allow bidi-safe LTR for IDs/timestamps when necessary.
- Status color cannot be sole indicator; always pair with text label.
- Respect reduced-motion for all sync state transitions.

---

## Phase 3 Do / Don't

### Do

- Keep Today view action-oriented with clear next steps.
- Surface sync status persistently but lightly.
- Preserve user trust with explicit save/retry messaging.
- Keep entry flows fast and repeatable for daily usage.
- Enforce client-only offline boundary.

### Don't

- Do not gate all content behind heavy dashboard analytics.
- Do not hide failed sync states in toasts only.
- Do not add offline queue logic to nutritionist/admin routes.
- Do not introduce new ad-hoc token scales or extra font weights.
- Do not block logging when offline if domain is supported.

---

## Checker Sign-Off

- [ ] Dimension 1 Copywriting: PASS
- [ ] Dimension 2 Visuals: PASS
- [ ] Dimension 3 Color: PASS
- [ ] Dimension 4 Typography: PASS
- [ ] Dimension 5 Spacing: PASS
- [ ] Dimension 6 Registry Safety: PASS

**Approval:** pending
