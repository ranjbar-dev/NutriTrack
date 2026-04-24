---
phase: 01
slug: platform-foundation
status: draft
shadcn_initialized: false
preset: none
created: 2026-04-22
---

# Phase 01 — UI Design Contract

> Visual and interaction contract for frontend Phase 1 (Platform Foundation). Generated for planner and executor consumption.

---

## Design System

| Property | Value |
|----------|-------|
| Tool | none (custom Nuxt 4 + Tailwind CSS 4 design system) |
| Preset | not applicable |
| Component library | Reka UI primitives (headless base), project-wrapped components |
| Icon library | lucide-vue-next |
| Font | Vazirmatn (base), Persian-first rendering |

### Locked Direction (Phase Context)

- Visual tone is clinical-minimal: calm, low-chroma, high legibility, not decorative.
- Persian-only UI: all user-facing copy and labels must be Persian.
- Default numeral presentation is Persian digits in UI display contexts.
- Date display is Jalali for user-facing views; transport/storage remains API-native Gregorian.

---

## Spacing Scale

Declared values (all multiples of 4):

| Token | Value | Usage |
|-------|-------|-------|
| xs | 4px | Icon gaps, inline micro-spacing |
| sm | 8px | Dense controls, chip padding |
| md | 16px | Default content spacing |
| lg | 24px | Card/internal section padding |
| xl | 32px | Screen section separation |
| 2xl | 48px | Major content breaks |
| 3xl | 64px | Page-level breathing room (rare on mobile) |

Exceptions:
- Minimum tap target: 44px height/width for tappable controls.
- Primary bottom navigation height: 56px minimum (excluding safe-area inset).
- Full-width action buttons: 48px fixed height.

---

## Typography

| Role | Size | Weight | Line Height |
|------|------|--------|-------------|
| Body | 16px | 400 | 1.5 |
| Label | 14px | 600 | 1.4 |
| Heading | 20px | 600 | 1.3 |
| Display | 28px | 600 | 1.2 |

Rules:
- Use only two weights in Phase 1: 400 and 600.
- Do not introduce additional font families in Phase 1.
- Default digits in visible UI text should be Persian digits.
- Numerically dense values (OTP, phone inputs, raw IDs) may use Latin digits only where input correctness requires it.

---

## Color

| Role | Value | Usage |
|------|-------|-------|
| Dominant (60%) | #F4F7F6 | App background, main surfaces |
| Secondary (30%) | #E6ECE9 | Cards, nav containers, grouped sections |
| Accent (10%) | #0F766E | Primary CTA, active nav marker, focus ring, selected controls |
| Destructive | #B42318 | Destructive action text/buttons and critical error emphasis only |

Accent reserved for:
- Primary action button per screen
- Active item in role navigation
- Focus-visible ring and selected input outline
- One high-priority status chip per view (not multiple accents in one block)

Color behavior constraints:
- No purple-led palettes.
- Avoid high-saturation decorative gradients.
- Success/info/warning chips must stay muted and readable in RTL compact layouts.

---

## Copywriting Contract

| Element | Copy |
|---------|------|
| Primary CTA | شروع و ادامه |
| Empty state heading | هنوز داده‌ای ثبت نشده است |
| Empty state body | برای شروع، اولین مورد را ثبت کنید یا صفحه را تازه‌سازی کنید. |
| Error state | ارتباط با سرور برقرار نشد. اتصال اینترنت را بررسی کنید و دوباره تلاش کنید. |
| Destructive confirmation | خروج از حساب: آیا از خروج از حساب کاربری مطمئن هستید؟ |
| Destructive confirmation | پاکسازی داده محلی: داده‌های ذخیره‌شده روی این دستگاه حذف می‌شود. ادامه می‌دهید؟ |

Tone rules:
- کوتاه، روشن، بدون لحن تبلیغاتی.
- عبارت های دستوری باید مستقیم و محترمانه باشند.
- در خطاها همیشه «مسئله + اقدام بعدی» نمایش داده شود.

---

## Registry Safety

| Registry | Blocks Used | Safety Gate |
|----------|-------------|-------------|
| shadcn official | none | not required |
| third-party | none | not applicable |

---

## Layout and Role Boundaries (Phase 1 Critical)

Primary focal point rule:
- On the default Phase 1 landing screen for each role, visual priority must be: page title in app bar -> one primary CTA zone in first viewport -> secondary content cards.

### Route and Layout Partition

- Role shells are isolated by top-level route groups and dedicated layouts:
  - Client: `/client/**` + `client` layout
  - Nutritionist: `/nutritionist/**` + `nutritionist` layout
  - Super Admin: `/admin/**` + `admin` layout
- Shared unauthenticated flows (`/auth/**`, install/update prompts) use a neutral auth layout.
- No cross-role nav items may render in another role shell, even as disabled placeholders.
- Each role shell has independent navigation config and page title map.

### Mobile Shell Composition

- Primary viewport target: 360-430px width.
- Top app bar: 56px, sticky, RTL title alignment.
- Bottom nav is allowed for client shell only in Phase 1 baseline; nutritionist/admin may use compact top tabs + overflow menu.
- Respect safe-area with `env(safe-area-inset-*)` for header/footer padding.
- Max readable content width on larger screens: 480px centered column while preserving mobile rhythm.

---

## Component Patterns (Foundation Set)

Phase 1 must establish reusable base components only:

- `AppShell` (role-aware)
- `AppHeader`
- `BottomNav` (client role)
- `SurfaceCard`
- `PrimaryButton`, `SecondaryButton`, `GhostButton`, `DangerButton`
- `BaseInput`, `BaseTextarea`, `BaseSelect`, `BaseSwitch`
- `InlineNotice` (info/warn/error)
- `LoadingBlock` (skeleton)
- `EmptyState`
- `ErrorState`
- `InstallPromptBanner`
- `UpdateAvailableBanner`

Pattern rules:
- Buttons have one dominant visual priority per viewport.
- Form controls must keep labels visible (no placeholder-only forms).
- Cards are low-elevation and border-led, not shadow-heavy.
- Icon-only interactive controls are allowed only with an explicit accessible label (`aria-label`) and must provide visible text label or tooltip on long-press/focus contexts.

---

## Interaction and Motion Rules

Motion must be subtle, purposeful, and low-cost:

- Duration tokens:
  - fast: 120ms
  - base: 180ms
  - slow: 240ms
- Easing:
  - standard: `cubic-bezier(0.2, 0, 0, 1)`
  - exit: `cubic-bezier(0.4, 0, 1, 1)`
- Allowed animations:
  - page/content fade-up on first mount (max 12px translate)
  - list item stagger only for 3-8 items (20ms step)
  - status badge pulse for sync/update attention (single cycle)
- Disabled animations:
  - parallax effects
  - infinite decorative loops
  - large spring/bounce transitions
- Respect reduced-motion preference: non-essential animation must be removed when `prefers-reduced-motion` is enabled.

---

## State Contracts (Loading, Empty, Error, Offline)

### Loading

- Use skeleton blocks for page-level loading over spinners when layout is known.
- Keep skeleton geometry aligned to final content to avoid layout jump.

### Empty

- Every empty state needs: title, one-line explanation, one primary next action.
- Empty visuals should be line-icon based and monochrome (no illustration-heavy assets in Phase 1).

### Error

- Error blocks must display:
  - short Persian message
  - retry action
  - optional support hint for repeated failures
- API errors should map technical codes to user-safe Persian text.

### Offline and Update (Platform baseline)

- Global connectivity indicator appears as slim top banner, not modal.
- PWA update availability appears as non-blocking in-app banner with actions:
  - `به‌روزرسانی`
  - `بعدا`
- Install prompt appears only after clear user intent moment (first-run completion or explicit install action), never on first paint.

---

## Responsive Behavior (Mobile-First)

Breakpoints and behavior:
- `<640px`: default mobile layout, single-column only.
- `640px-768px`: preserve single-column; increase horizontal padding only.
- `>768px`: still mobile-centric centered canvas; do not switch to desktop dashboard grid in Phase 1.

RTL behavior rules:
- All directional icons must auto-mirror where semantically directional.
- Truncation must preserve readable Persian word boundaries.
- Mixed LTR data (URLs, IDs) must render with bidi-safe wrappers.

---

## Phase 1 Do / Don't Rules

### Do

- Keep interface calm, clinical-minimal, and hierarchy-first.
- Build shell and primitives that can be shared across all role areas.
- Enforce strict role route/layout boundaries from the start.
- Prioritize thumb-reachable controls and safe-area correctness.
- Use explicit Persian microcopy for loading/empty/error/install/update states.

### Don't

- Do not introduce generic dashboard templates with dense KPI cards.
- Do not add desktop-first sidebars or multi-column admin grids in Phase 1.
- Do not use decorative gradients, glassmorphism, or high-chroma accents.
- Do not mix multiple accent colors for interactive hierarchy.
- Do not cache sensitive authenticated payloads broadly in UI-driven PWA behavior.

---

## Checker Sign-Off

- [ ] Dimension 1 Copywriting: PASS
- [ ] Dimension 2 Visuals: PASS
- [ ] Dimension 3 Color: PASS
- [ ] Dimension 4 Typography: PASS
- [ ] Dimension 5 Spacing: PASS
- [ ] Dimension 6 Registry Safety: PASS

**Approval:** approved
