# Feature Landscape

**Domain:** Persian RTL nutrition management PWA for nutritionists and clients  
**Project:** NutriTrack client  
**Researched:** 2026-04-22  
**Scope:** Frontend-only, against the documented REST API, with offline support required for client flows only.

## Executive View

For this product category, the frontend is credible only if it behaves like a daily-use mobile tool rather than a CRUD shell. The table stakes are not just "show the data". They are: a plan-first home screen, fast daily logging, chat, notifications, clear role-based navigation, and an offline model that users can understand and trust.

The strongest differentiators for v1 are not flashy features. They are product qualities that directly reduce friction for Persian-speaking mobile users: excellent RTL information hierarchy, Jalali-aware scheduling, low-friction one-thumb tracking, and explicit offline/sync feedback. Those will matter more than adding broad admin surfaces or speculative intelligence features.

For v1, the product should deliberately exclude anything that expands scope without improving the core nutrition workflow: desktop optimization, real-time chat, AI recommendations, wearable integrations, advanced analytics, and client-side authoring workflows that bypass nutritionist curation.

## Table Stakes

Features users will expect in a serious nutrition management PWA. Missing these will make the product feel incomplete even if the backend exists.

| Feature | Why Expected | Complexity | Dependencies | Notes |
|---------|--------------|------------|--------------|-------|
| Role-based mobile app shell | Users need immediate role-appropriate navigation after login instead of one shared generic UI. | Medium | Auth endpoints, profile endpoints, route guards | Separate client, nutritionist, and admin information architecture. Client home should be plan-centric; nutritionist home should be client/workload-centric. |
| Authentication flows for both audiences | Client OTP login and staff email/password login are baseline trust and access requirements. | Medium | `POST /auth/login`, `POST /auth/otp/send`, `POST /auth/otp/verify`, token refresh/logout | OTP must feel lightweight on mobile and clearly communicate rate limits, expiry, and retry states. |
| Plan-first client dashboard | Clients expect to open the app and immediately understand today's meals, water target, medication times, and pending actions. | Medium | `GET /plans/active`, tracking read endpoints, notification preferences | This is the primary client surface. It should privilege “what do I do now?” over generic profile content. |
| Full active-plan viewing with clear hierarchy | A nutrition client must be able to inspect plan period, day structure, meals, options, exercise recommendations, medications, and notes without confusion. | Medium | `GET /plans/active`, `GET /plans/:id` | The nested plan structure is complex; the UI must flatten it into readable cards, timelines, or grouped sections. |
| Fast daily tracking flows | Logging food, water, sleep, exercise, medication intake, and body measurements is core daily product value. | High | `POST /tracking/food`, `/water`, `/sleep`, `/exercise`, `/medication`, `/body` | Tracking must optimize for repeat entry, not form completeness. Use quick actions, presets, and obvious completion states. |
| Client offline read/write for core flows | Offline support is a product promise, not a nice-to-have. Without it, the client app breaks the PRD. | High | IndexedDB/service worker architecture, `POST /tracking/sync`, idempotent `local_id` handling | Needs cached active plan, cached recent messages, queued writes, reconnect sync, retry states, and visible sync status. |
| Sync visibility and recovery UX | Users must know whether entries are saved locally, synced, retrying, or failed. Hidden offline state creates mistrust. | High | Offline queue, bulk sync endpoint, local persistence | This is functionally table stakes because offline is mandatory. A silent queue is not sufficient. |
| Messaging with attachments and unread state | Direct client-nutritionist chat is a core care workflow. Users will expect readable conversation history and attachment support. | Medium | Message list/send endpoints, unread count endpoint, push subscriptions | Since the API uses polling, the frontend must make polling feel stable and battery-conscious. |
| Push notifications and preferences | Diet updates, messages, and reminders are expected for adherence in a mobile care product. | Medium | Push subscribe/unsubscribe, notification preferences | Permission prompts should be delayed until contextual value is clear, not shown blindly at first paint. |
| Food discovery plus food request fallback | Clients need a way to understand foods in their plan and request missing foods without dead ends. | Medium | Foods search/read endpoints, food requests endpoints | For v1, the request flow matters more than exposing broad manual food authoring to clients. |
| Nutritionist client roster and client profile workspace | Nutritionists need a mobile workspace to find clients, inspect progress, message them, and jump into plan actions. | High | Client management endpoints, tracking read endpoint, plans endpoints, messages endpoints | The nutritionist profile should aggregate tracking tabs, current plan summary, lab results, and message access. |
| Diet plan authoring and editing UI | Nutritionists cannot deliver value without building plans, meals, options, exercise recommendations, and prescriptions from mobile. | High | Full plan CRUD tree, foods search, medications search | This is the most complex v1 authoring surface. Progressive disclosure is mandatory to avoid a hostile mobile form. |
| Shared food and medication catalogue management | Nutritionists and admins need searchable catalogue screens to support plan authoring and curation. | Medium | Foods, categories, and medications endpoints | Search, pagination, and create/edit states are more important than dense desktop-style tables. |
| Lab result upload and viewing | Clients expect to submit test results, and nutritionists expect a simple review/download path. | Medium | Lab results endpoints, file upload handling | Explicitly online-only. The UX should say that before the user invests effort. |
| Basic admin operational screens | Super admin needs lightweight surfaces for nutritionist management and platform database control. | Medium | Admin stats, nutritionist management, admin foods/medications endpoints | This is table stakes for the role, but should remain operationally simple in v1. |
| Profile, avatar, and session management | Users expect to manage identity basics and safely recover sessions. | Low | `/auth/me`, avatar upload, refresh/logout | This should stay lightweight and subordinate to the core workflow. |

## Differentiators

Features and UX qualities that are not strictly required to ship, but would make this frontend noticeably better than a generic care portal.

| Feature | Value Proposition | Complexity | Dependencies | Notes |
|---------|-------------------|------------|--------------|-------|
| Explicit sync center for clients | Builds trust by showing queued items, last successful sync, retry actions, and per-entry status. | High | Offline queue model, bulk sync endpoint, local metadata | This turns offline support from hidden plumbing into a user-visible reliability feature. |
| Persian-native nutrition UX | Jalali-aware dates, Persian numerals where appropriate, RTL-first spacing, and familiar portion units make the app feel local instead of translated. | Medium | UI layer, date/formatting utilities, plan and food metadata | This is a stronger differentiator than visual decoration. It directly affects comprehension speed. |
| “Today” timeline instead of generic dashboards | A chronological daily agenda for meals, water, meds, exercise, and reminders reduces cognitive load for clients. | Medium | Active plan, tracking state, reminder preferences | Strong fit for mobile PWA use and one-handed interaction. |
| One-thumb quick logging patterns | Chips, steppers, repeat-last-entry, and tap-to-complete interactions can make tracking dramatically faster than standard forms. | Medium | Tracking endpoints, local draft state | Particularly valuable for water, medication intake, and meal adherence. |
| Plan comprehension aids | Nutritional rollups, meal option comparison, and clear “what changed since last plan” summaries help clients and nutritionists understand the plan faster. | Medium | Full plan object, archived plans, computed totals from API | This is more useful than adding richer analytics early. |
| Conversation-aware care UX | Message composer shortcuts tied to plan context, unread badge discipline, and attachment previews make chat feel like part of care, not a bolted-on inbox. | Medium | Messages endpoints, unread count, uploads | Can be shipped incrementally without needing real-time transport. |
| Compact progress storytelling | Simple weight and measurement trends, streak-like adherence summaries, and day completion markers provide motivation without full analytics bloat. | Medium | Tracking history reads, archived plans | Keep this lightweight and narrative, not dashboard-heavy. |
| Guided nutritionist mobile authoring | Stepwise plan-building flows, saved drafts, and reusable defaults reduce the pain of authoring complex plans on mobile. | High | Plan CRUD tree, foods/medications search, local transient state | This can become a true product differentiator because mobile plan authoring is otherwise cumbersome. |

## Anti-Features / Deferred Items

Items that should be deliberately excluded from v1 or clearly deferred because they expand scope faster than they improve the core workflow.

| Anti-Feature / Deferred Item | Why Avoid in v1 | What to Do Instead |
|------------------------------|-----------------|-------------------|
| Desktop-first or tablet-optimized layouts | Conflicts with the documented mobile-only product shape and adds large design/test surface. | Optimize exclusively for narrow mobile viewports and installable PWA behavior. |
| Real-time chat via WebSocket | The API contract is polling-based, and realtime transport adds complexity without changing the care workflow meaningfully. | Ship strong polling UX with unread counts, attachment support, and push notifications. |
| Client-created custom foods as a primary flow | The PRD points clients toward food requests, and open-ended client food creation can undermine catalogue quality and nutritionist oversight. | Expose request-food flow first; revisit client personal foods only after adoption evidence. |
| Broad offline support for nutritionists or admins | Mandatory offline support is explicitly client-only. Extending it multiplies complexity across authoring and management flows. | Keep staff experiences online-first and reliable on good mobile networks. |
| Offline lab-result upload | File upload conflicts with the documented online-only constraint and complicates storage/state handling. | Make upload clearly internet-required and preserve draft metadata only if needed later. |
| Advanced analytics dashboards | Dense charts and multi-metric dashboards add design and interpretation cost before the core loop is proven. | Ship a few simple trend and summary views tied to concrete actions. |
| Complex automation builders for reminders | Custom rule builders are expensive and unnecessary for initial adherence support. | Start with notification preferences plus schedule-driven reminders from the plan. |
| AI meal recommendations or coaching | Explicit PRD non-goal and a distraction from validated workflow value. | Invest in clearer plans, tracking UX, and messaging. |
| Food-photo logging or calorie recognition | Explicit PRD non-goal and high implementation risk for weak early payoff. | Keep logging structured and tied to plan options and curated foods. |
| Wearables or health-device integrations | Explicit PRD non-goal with major data and support implications. | Keep manual tracking focused and low friction. |
| Payment, billing, or subscription management | Outside product scope and unrelated to the care workflow being validated. | Omit from v1. |
| Large admin audit/reporting suite | The admin role only needs operational control in v1, not a full BI environment. | Provide basic stats, account management, and database curation only. |
| Gamification layer | Streak systems, badges, and rewards can dilute a medically grounded workflow if added too early. | Use simple completion and progress cues instead of game mechanics. |

## Information Architecture Implications

The frontend should be organized around tasks, not resource nouns.

### Client IA

1. Home / Today
2. Active Plan
3. Track
4. Messages
5. Results / Progress
6. Settings

Recommended behavior:
- Home should merge schedule awareness, reminders, and pending offline sync state.
- Track should prioritize fast entry over historical browsing.
- Results / Progress should combine measurements, archived plans, and lightweight trends.

### Nutritionist IA

1. Clients
2. Client Profile
3. Plans
4. Messages
5. Foods
6. Medications
7. Settings

Recommended behavior:
- Client Profile should be the central workspace, not a passive details page.
- Plans should be reachable from both the client list and the client profile.
- Catalogue screens should support plan authoring, not exist as isolated admin CRUD.

### Super Admin IA

1. Nutritionists
2. Foods
3. Medications
4. Platform Stats

Recommended behavior:
- Keep admin shallow and operational.
- Do not let the admin surface dominate product effort early.

## Feature Dependencies

```text
Authentication -> Role-based app shell
Role-based app shell -> Client home / Nutritionist workspace / Admin operations
Foods + Medications + Client management -> Diet plan authoring
Diet plan authoring -> Active plan viewing -> Daily tracking context
Active plan caching + Local queue + Sync endpoint -> Reliable client offline experience
Messages + Push subscriptions + Notification preferences -> Care communication loop
Tracking history -> Progress views -> Lightweight adherence summaries
Food requests -> Catalogue quality improvement without exposing broad client authoring
```

## MVP Recommendation

Prioritize:

1. Client auth, role-aware shell, and active-plan viewing
2. Client offline-first tracking with visible sync state for food, water, sleep, exercise, medication, and measurements
3. Messaging, push notifications, and basic notification preferences
4. Nutritionist client list, client profile, and basic plan authoring/editing
5. Food search/request flow and minimal lab-result upload/view

Defer:

- Advanced analytics: simple trends are enough for v1.
- Admin depth: keep to core operational needs.
- Client-authored custom foods, realtime chat, and broader automation: they complicate the product before the main loop is proven.

## Suggested Phase Framing

1. **Foundations and access**
   - Auth, session handling, role-based shell, Persian RTL mobile scaffolding
2. **Client core loop**
   - Active plan, today view, tracking, offline queue, sync visibility
3. **Care communication**
   - Messages, attachments, unread state, push subscriptions, preferences
4. **Nutritionist workspace**
   - Client list, profile, history tabs, basic plan authoring
5. **Catalogue and supporting workflows**
   - Foods, medications, food requests, lab results, lightweight admin tools

## Confidence Notes

| Area | Confidence | Notes |
|------|------------|-------|
| Client table stakes | HIGH | Directly supported by PRD offline strategy, role definitions, and tracking/message API surface. |
| Nutritionist table stakes | HIGH | Plan management, client management, and messaging are explicitly defined in both PRD and API. |
| Differentiators | MEDIUM | Based on documented constraints plus common mobile care-product UX expectations rather than explicit product mandates. |
| Anti-features / deferrals | HIGH | Strongly supported by PRD non-goals, mobile-only scope, and API constraints such as polling chat and client-only offline support. |

## Sources

- `.planning/PROJECT.md`
- `docs/PRD.md`
- `docs/API.md`