# Technology Stack

**Project:** NutriTrack Client  
**Researched:** 2026-04-22  
**Scope:** Frontend-only, mobile-first Persian RTL PWA  
**Overall recommendation confidence:** HIGH

## Recommendation Summary

For this project, the right 2026 stack is:

- Nuxt 4 for the app shell, routing, SSR-aware data flow, and PWA host framework.
- Tailwind CSS 4 for styling, with custom tokens and a bespoke mobile UI rather than a prebuilt component framework.
- Pinia for client state only: session, UI state, offline queue metadata, and a few workflow stores.
- TanStack Query for server state and cache lifecycle.
- Zod plus vee-validate for forms and runtime validation.
- Dexie on IndexedDB for offline-first persistence and an explicit outbox/sync model.
- @vite-pwa/nuxt with a custom service worker path for installability, caching, push, and selective offline behavior.
- Native browser Push API, Notification API, and multipart uploads, with thin utilities around them instead of heavyweight SaaS abstractions.

This is the best fit for the documented NutriTrack API because the backend contract is custom REST with JWT access and refresh tokens, OTP login for clients, multipart uploads, and explicit push subscription endpoints. That makes a Nuxt-native data layer plus a custom auth implementation a better fit than OAuth-focused auth frameworks.

## Recommended Stack

### Core Framework

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| Nuxt | 4.4.x | Application framework | Current Nuxt 4 docs show Nuxt 4.4.2, SSR-aware data fetching, auto-imports, file routing, and Vite-by-default. This is the correct host for a mobile PWA that still benefits from SSR and route-level hydration control. | HIGH |
| Vue | 3.x via Nuxt | Reactive UI runtime | Comes with Nuxt 4; no reason to deviate. | HIGH |
| TypeScript | 5.5+ strict mode | Application typing | Required for Zod-first contracts, Pinia store safety, and API typing discipline. | HIGH |

### Styling and UI Layer

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| Tailwind CSS | 4.2.x | Utility-first styling | Tailwind's current Nuxt guide installs Tailwind 4 with the Vite plugin directly. Best fit for a custom, intentional mobile UI instead of a generic dashboard kit. | HIGH |
| @tailwindcss/vite | Tailwind 4 companion | Tailwind integration for Nuxt/Vite | Official Tailwind Nuxt guidance uses the Vite plugin path, not the old Nuxt module path. | HIGH |
| Reka UI | Latest stable | Accessible unstyled primitives | Strong fit for custom UI Pro Max direction: unstyled, accessible, RTL-aware, and Vue-native. Use for dialogs, drawers, tabs, selects, tooltips, switches, and menu primitives without surrendering visual identity. | MEDIUM |
| @fontsource/vazirmatn | Latest stable | Persian UI typography | Vazirmatn is a strong Persian-first font choice and is readily consumable through Fontsource. Better than default Latin-biased stacks for a Persian-only product. | HIGH |
| lucide-vue-next | v1.x | Icon library | Tree-shakable Vue icon components with good accessibility guidance. Good default until a product-specific icon language emerges. | HIGH |

### State, Data, and API Layer

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| @pinia/nuxt + pinia | Pinia 3.x docs track | Client state management | Official Pinia docs support Nuxt 3 and 4 and provide the right integration path. Use Pinia for auth session state, optimistic local workflow state, install-prompt state, and sync orchestration metadata. | HIGH |
| @tanstack/vue-query | 5.x | Server state, caching, invalidation | In 2026 this remains the standard answer for async server state in Vue apps. It handles cache freshness, pagination, retries, optimistic updates, and request deduping better than hand-rolled Pinia fetch stores. | HIGH |
| Nuxt useFetch and $fetch | Built-in | HTTP client and SSR-safe data fetching | Official Nuxt 4 guidance recommends useFetch or useAsyncData plus $fetch. This should be the base HTTP layer, wrapped once in a project-specific useApiFetch composable. | HIGH |
| Zod | 4.x | Runtime validation and inferred types | Zod 4 is stable and remains the cleanest way to validate request and response boundaries against an API doc that is detailed but not machine-generated. | HIGH |
| VueUse | 14.x | Browser and app utilities | Use for network state, permissions, storage fallbacks, event listeners, and PWA-adjacent utilities. Small, SSR-friendly, and already idiomatic in Vue/Nuxt projects. | HIGH |

### Forms and Validation

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| vee-validate | 4.x | Form state and submission flow | Still the most established Vue form library, especially for composition API apps. Good fit for OTP login, profile forms, measurement forms, diet tracking forms, and preferences. | HIGH |
| @vee-validate/zod | Latest stable | Zod integration for forms | Lets the project define one schema per domain and reuse it for both form validation and API payload validation. | HIGH |

### Offline, Storage, and PWA

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| Dexie | Latest stable | IndexedDB abstraction | Best fit for offline-first NutriTrack. IndexedDB is required for structured offline data and outbox queues; Dexie is the cleanest mature wrapper for that job. | HIGH |
| @vite-pwa/nuxt | Latest stable compatible with Nuxt 4 | PWA integration | This is the standard Vite/Nuxt PWA path. It provides manifest wiring, service worker registration, install prompts, update prompts, and offline-ready hooks. | HIGH |
| Workbox via @vite-pwa/nuxt injectManifest | Via module | Custom service worker behavior | NutriTrack needs more than static precache. Push handling, fine-grained runtime caching, and offline sync hooks justify injectManifest instead of a fully generated worker. | MEDIUM |

### Uploads and Media

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| Native FormData + $fetch | Built-in | File upload transport | The API contract is plain multipart/form-data. Use the browser's native file input and FormData directly. It is simpler and lighter than a generic uploader framework. | HIGH |
| browser-image-compression | 2.x | Client-side image compression | Good fit for avatar uploads and chat images on mobile networks. Supports web workers and reduces payload size before upload. | MEDIUM |

### Notifications and Push

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| Native Push API + PushManager | Browser API | Web push subscription | The backend API already exposes push subscription endpoints. The client should subscribe using the standard Push API and post subscriptions to the backend. | HIGH |
| Native Notification API | Browser API | Permission and display model | Required for a standards-based PWA push flow. Must be requested from explicit user gestures. | HIGH |

### Testing

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| Vitest | 4.x | Unit and component testing | Best fit with Vite/Nuxt projects. Fast local loop and current ecosystem standard. | HIGH |
| @nuxt/test-utils | Latest stable | Nuxt-aware test harness | Needed to test composables, plugins, and Nuxt-specific runtime boundaries without fighting the framework. | HIGH |
| @vue/test-utils | Latest stable | Vue component testing | Standard Vue test utility layer under Vitest. | HIGH |
| MSW | 2.x | API mocking in tests and local development | Excellent for contract-faithful REST mocks without patching the app internals. Reusable across local dev, unit tests, and integration-style tests. | HIGH |
| Playwright | Latest stable | End-to-end and mobile PWA verification | Best fit for mobile viewport flows, installability checks, offline smoke tests, uploads, and push-permission UX verification. | HIGH |

### Date, Calendar, and Locale Handling

| Technology | Version | Purpose | Why | Confidence |
|------------|---------|---------|-----|------------|
| date-fns | 4.x | Core date arithmetic | Current date-fns highlights v4 with time zone support and remains modular, immutable, and TypeScript-friendly. | HIGH |
| date-fns-jalali | 4.1.0-0 lineage | Jalali calendar operations | Still the cleanest fit for Jalali-aware formatting and arithmetic while keeping the date-fns programming model. Use it for Persian-facing date UI. | MEDIUM |
| Intl.DateTimeFormat with fa-IR-u-ca-persian | Browser built-in | Lightweight localized formatting | Use native Intl for simple display formatting and Persian digits where arithmetic is not needed. | HIGH |

## Prescriptive Architecture Choices

### 1. State split

Use this split from day one:

- Pinia: auth session, current user, notification preference flags, install prompt visibility, sync engine state, temporary UI state.
- TanStack Query: anything fetched from the API and expected to be revalidated or invalidated.
- Dexie: offline cache and outbox persistence.

Do not use Pinia as the primary cache for server resources like foods, medications, diet plans, messages, or tracking history. That is exactly the class of problem TanStack Query solves better.

### 2. Auth implementation

Implement auth as a project-local module, not a generic auth framework.

Recommended approach:

- `useAuthStore` in Pinia for session state.
- Access token kept in memory plus short-lived persistence policy based on product security decisions.
- Refresh token stored in the safest browser-accessible mechanism the backend contract allows, usually a secure cookie if the backend is configured for it, otherwise a clearly isolated persisted token store.
- One `useApiFetch` or `apiFetch` wrapper that injects auth headers, retries one refresh path on 401, and logs out cleanly on refresh failure.
- Route middleware for role-gated pages: `client`, `nutritionist`, `super_admin`.

Why not use a generic auth framework:

- The NutriTrack API is not OAuth-first.
- Client login is OTP plus JWT, not provider-based login.
- The refresh flow is custom.
- Auth.js style integrations add complexity without solving the hard part here.

### 3. Offline model

Offline support is a core product promise, so do not treat it as cache-only.

Recommended offline tables in Dexie:

- `dietPlans`
- `planDays`
- `trackingEntries`
- `trackingDrafts`
- `messages`
- `messageDrafts`
- `foodRequestsDrafts`
- `labResultDrafts`
- `syncQueue`
- `syncRuns`
- `notificationPrefs`

Recommended sync model:

- Every offline write creates an outbox item with deterministic operation metadata.
- Reconnect triggers a sync runner from app resume, explicit user action, and online events.
- Background Sync is an enhancement only, not the primary guarantee, because Safari and iOS support remains weak.

### 4. PWA strategy

Use `@vite-pwa/nuxt` with `injectManifest`.

Why:

- Push notification handling needs custom service worker code.
- Offline queue replay benefits from custom worker hooks.
- NutriTrack needs deliberate runtime caching decisions, not just a generic precache.

Cache strategy guidance:

- Precache app shell, fonts, icons, and static assets.
- Cache diet-plan reads and low-churn reference data with stale-while-revalidate.
- Avoid long-lived opaque caching for authenticated JSON responses carrying sensitive health-adjacent user data.
- Never assume message send or tracking POST success from a cache layer; those must go through the explicit outbox.

### 5. Dates and calendar rules

Store and exchange dates exactly as the API defines them:

- Dates: Gregorian `YYYY-MM-DD`
- Timestamps: RFC 3339

Render to users in Persian and Jalali.

Practical rule:

- Use API-native Gregorian strings in storage and transport.
- Convert only at the display and user-input boundary.
- Use `date-fns-jalali` when arithmetic or parsing must respect Jalali semantics.
- Use `Intl.DateTimeFormat('fa-IR-u-ca-persian')` for lightweight display-only formatting.

## What To Use For Each Concern

| Concern | Recommendation | Why |
|---------|---------------|-----|
| Auth | Custom Nuxt composables + Pinia + API fetch wrapper | Best fit for OTP plus JWT plus refresh flow |
| API client | `useFetch` and `$fetch` wrapped in `useApiFetch` | Native Nuxt path, SSR-aware, less dependency weight |
| Response validation | Zod schemas at the boundary | Protects the app from API drift and bad offline replay data |
| Forms | vee-validate + Zod | Best Vue DX for complex mobile forms |
| Server cache | TanStack Query | Better than Pinia for async cache lifecycle |
| Offline persistence | Dexie | Structured IndexedDB with mature API |
| PWA | `@vite-pwa/nuxt` + custom SW | Best Nuxt-native PWA path |
| Image upload | Native file input + FormData + browser-image-compression | Lightest stack for current API |
| Push | Push API + Notification API + service worker | Matches backend subscription endpoints |
| Test mocks | MSW | Reusable network mocking layer |
| E2E | Playwright | Strong mobile and PWA coverage |
| Persian dates | date-fns + date-fns-jalali + Intl | Correct split between arithmetic and display |

## What Not To Use

| Category | Do Not Use | Why Not |
|----------|------------|---------|
| Auth | Auth.js, sidebase/nuxt-auth, or OAuth-centric auth stacks | Wrong abstraction for OTP plus custom JWT refresh API |
| HTTP client | Axios by default | Nuxt already standardizes on `$fetch` and `useFetch`; Axios adds duplicate concepts and bundle weight |
| Server state | Pinia-only fetch stores for all API data | Leads to homegrown cache invalidation and stale data bugs |
| Offline storage | localStorage for core offline workflows | Too small, string-only, and unsafe for queueable structured data |
| PWA | Auto-generated service worker only | Too rigid for push plus offline outbox needs |
| Uploads | Uppy as the default initial uploader | Strong product, but unnecessary weight until the backend supports resumable uploads, remote sources, or complex dashboard flows |
| UI kit | Heavy prebuilt component suites such as Vuetify or Quasar for the main app shell | They will fight the UI Pro Max custom direction and produce a generic app feel |
| Dates | Moment, moment-jalaali, or ad hoc string math | Outdated, heavier, and easier to get wrong than date-fns-based tooling |
| Sync | Background Sync as the only offline replay path | Browser support is still too inconsistent, especially on Safari/iOS |

## Suggested Package Set

### Production

```bash
npm install nuxt @pinia/nuxt pinia @tanstack/vue-query zod vee-validate @vee-validate/zod @vite-pwa/nuxt dexie @vueuse/core reka-ui @fontsource/vazirmatn lucide-vue-next browser-image-compression date-fns date-fns-jalali
```

### Development

```bash
npm install -D tailwindcss @tailwindcss/vite vitest @nuxt/test-utils @vue/test-utils msw playwright
```

## Implementation Notes For This Repo

### Nuxt module and plugin posture

Recommended Nuxt modules/plugins at project start:

- `@pinia/nuxt`
- `@vite-pwa/nuxt`
- Tailwind Vite plugin via `vite.plugins`

Recommended app-level wrappers to create immediately:

- `composables/useApiFetch.ts`
- `composables/useAuth.ts`
- `plugins/vue-query.ts`
- `plugins/dexie.client.ts`
- `plugins/pwa.client.ts`
- `middleware/auth.global.ts` or role-specific named middleware

### Offline boundary decisions

Offline-first should cover at least:

- viewing active diet plan
- recording meal choices
- water intake logging
- sleep logging
- exercise logging
- medication intake logging
- weight and measurement entry
- composing messages and drafts

Uploads should be split:

- Small image uploads can queue metadata and prompt retry when online.
- Large files such as lab-result PDFs should degrade gracefully to pending-upload state rather than pretending to complete offline.

## Confidence Assessment

| Area | Level | Notes |
|------|-------|-------|
| Nuxt 4 + Tailwind 4 + Pinia integration | HIGH | Verified against current official Nuxt, Tailwind, and Pinia docs |
| Server-state recommendation with TanStack Query | HIGH | Strong current ecosystem fit and official Vue docs remain active |
| Forms and validation with vee-validate + Zod | HIGH | Current docs are clear and ecosystem support remains strong |
| Offline storage with Dexie | HIGH | Mature IndexedDB abstraction and strong fit for explicit offline queues |
| PWA with @vite-pwa/nuxt | HIGH | Current docs confirm Nuxt integration path and install/update hooks |
| Custom auth over auth frameworks | HIGH | Driven by product API contract, not just ecosystem preference |
| Reka UI as primitive layer | MEDIUM | Good fit for the design constraints and RTL support, though still a product choice rather than a hard requirement |
| date-fns-jalali recommendation | MEDIUM | Good fit and active enough to use, but exact long-term maintenance strength is weaker than core date-fns |
| browser-image-compression recommendation | MEDIUM | Good lightweight fit for mobile uploads, but should stay narrowly scoped to images |

## Alternatives Considered

| Category | Recommended | Alternative | Why Not |
|----------|-------------|-------------|---------|
| Auth | Custom auth module | Auth.js or sidebase/nuxt-auth | Wrong center of gravity for this API |
| Server state | TanStack Query | Pinia-only async stores | Too much manual cache logic |
| Offline DB | Dexie | localForage or idb-keyval | Too limited for relational-like offline queues |
| UI primitives | Reka UI | Full UI framework | Too prescriptive for a bespoke mobile design system |
| Uploads | Native FormData | Uppy | Overkill until resumable or remote-source uploads matter |
| Dates | date-fns plus date-fns-jalali | Day.js plus Jalali plugin | Weaker fit for explicit functional date arithmetic and shared utility style |

## Sources

- Nuxt 4 introduction and data-fetching docs: https://nuxt.com/docs/4.x/getting-started/introduction
- Nuxt `useFetch`: https://nuxt.com/docs/4.x/api/composables/use-fetch
- Nuxt `$fetch`: https://nuxt.com/docs/4.x/api/utils/dollarfetch
- Pinia Nuxt integration: https://pinia.vuejs.org/ssr/nuxt.html
- Tailwind CSS Nuxt guide: https://tailwindcss.com/docs/installation/framework-guides/nuxt
- Vite PWA Nuxt integration: https://vite-pwa-org.netlify.app/frameworks/nuxt.html
- TanStack Query Vue overview: https://tanstack.com/query/latest/docs/framework/vue/overview
- VueUse: https://vueuse.org/
- Zod: https://zod.dev/
- vee-validate: https://vee-validate.logaretm.com/v4/
- Dexie docs: https://dexie.org/docs/Tutorial/Getting-started
- Reka UI: https://reka-ui.com/
- Lucide Vue: https://lucide.dev/guide/packages/lucide-vue-next
- Fontsource Vazirmatn: https://fontsource.org/fonts/vazirmatn
- Vitest: https://vitest.dev/guide/
- Playwright: https://playwright.dev/docs/intro
- MSW: https://mswjs.io/docs/
- Notification API: https://developer.mozilla.org/en-US/docs/Web/API/Notification
- PushManager: https://developer.mozilla.org/en-US/docs/Web/API/PushManager
- Background Synchronization API: https://developer.mozilla.org/en-US/docs/Web/API/Background_Synchronization_API
- date-fns: https://date-fns.org/
- date-fns-jalali: https://github.com/date-fns-jalali/date-fns-jalali
- browser-image-compression: https://minimalfolk.github.io/browser-image-compression/
- Uppy: https://github.com/transloadit/uppy
