## UI-SPEC COMPLETE

---

# Phase 3 UI Design Contract: Diet Plan Engine

**Status:** draft  
**Phase:** 03 – Diet Plan Engine  
**Generated:** auto-chain (no interactive session)  
**Design System:** Tailwind v4 CSS-first (no shadcn) — extends Phase 1/2 patterns  

---

## Design Philosophy

### RTL Mobile-First Drill-Down

All UI is **Persian-only RTL**, rendered inside a `max-width: 430px` centered viewport (set in `main.css`). No desktop breakpoints exist or are needed.

Phase 3 introduces a **drill-down navigation pattern**: the nutritionist descends through Plan → Day → Meal depth levels using forward-navigation with a persistent breadcrumb bar. Each level is a separate Nuxt page route. The client view uses a single page with a horizontal tab bar for day switching.

### Core Principles (inherited from Phases 1 & 2)

1. **Logical RTL properties only** — `ms-`, `me-`, `ps-`, `pe-`, `text-start`, `text-end`, `border-s-`, `float-start` everywhere. No `left`/`right` utility classes.
2. **Card-based layout** — `bg-white rounded-2xl p-4 shadow-sm` for all content groups. No tables.
3. **Touch-first** — minimum 44px touch targets for all interactive elements. No hover-only states.
4. **Persian digits** — `toPersianDigits()` for all numbers shown to users.
5. **Shamsi dates** — `useShamsiDate()` for all date display.
6. **Emerald accent** — `emerald-600` as the single accent color for primary actions, active states, focus rings.
7. **Persian strings hardcoded** — no i18n library; all copy directly in components.

---

## Spacing Scale

8-point grid. All spacing values must be multiples of 4px:

| Token | px | Usage |
|-------|----|-------|
| `p-1` | 4px | tight inline spacing |
| `p-2` | 8px | badge padding, small gaps |
| `p-3` | 12px | card inner row padding |
| `p-4` | 16px | page padding, card padding (standard) |
| `p-6` | 24px | section vertical spacing |
| `p-8` | 32px | large section separators |
| `space-y-3` | 12px | between list items |
| `space-y-4` | 16px | between form sections |
| `gap-2` | 8px | between pills/badges |
| `gap-3` | 12px | between form fields (2-col grid) |
| `gap-4` | 16px | between cards |
| `pb-20` | 80px | bottom padding for fixed nav bar clearance |
| `h-16` | 64px | bottom nav bar height (matches existing) |
| `min-h-[44px]` | 44px | minimum touch target height |

---

## Typography

**Font:** Vazirmatn (already configured in `@theme { --font-family-sans }`)

### Sizes (exactly 4 in use)

| Role | Class | Size | Weight | Line-Height | Usage |
|------|-------|------|--------|-------------|-------|
| Page heading | `text-xl font-bold` | 20px | 700 | 1.4 | Page titles (e.g. "برنامه‌ی جدید") |
| Section heading | `text-sm font-semibold` | 14px | 600 | 1.4 | Card section labels (e.g. "وعده‌ها") |
| Body | `text-base` | 16px | 400 | 1.5 | Main content, food names, labels |
| Caption | `text-xs` | 12px | 400 | 1.4 | Badges, counts, secondary info |

**Weights:** 400 (regular) + 600 (semibold) + 700 (bold heading only).

---

## Color Palette

### Base Colors (inherited)

| Role | Tailwind | Hex | Usage |
|------|----------|-----|-------|
| Page background | `bg-gray-50` | #f9fafb | `min-h-screen bg-gray-50` on layouts |
| Card surface | `bg-white` | #ffffff | All content cards |
| Primary accent | `bg-emerald-600` | #059669 | Primary buttons, active tab indicator, active nav |
| Accent hover | `bg-emerald-700` | #047857 | Button hover state |
| Accent light | `bg-emerald-100` | #d1fae5 | Active status badge background |
| Border | `border-gray-200` | #e5e7eb | Card borders, dividers |
| Input border | `border-gray-300` | #d1d5db | Default input |
| Input focus | `border-emerald-500` | #10b981 | Focus state |
| Heading text | `text-gray-800` | #1f2937 | Page/section headings |
| Body text | `text-gray-700` | #374151 | Labels, body copy |
| Secondary text | `text-gray-600` | #4b5563 | Descriptions, secondary info |
| Muted text | `text-gray-500` | #6b7280 | Captions, placeholders |
| Error | `text-red-600` | #dc2626 | Error messages |
| Error bg | `bg-red-100` | #fee2e2 | Error state backgrounds |

### Status Badge Colors (Phase 3 — new)

| Status | Persian Label | Background | Text | Border |
|--------|--------------|------------|------|--------|
| draft | پیش‌نویس | `bg-orange-100` | `text-orange-700` | none |
| active | فعال | `bg-emerald-100` | `text-emerald-700` | none |
| archived | آرشیو | `bg-gray-100` | `text-gray-600` | none |

**60/30/10 distribution:**
- 60% dominant: `gray-50` page + `white` cards
- 30% secondary: `gray-100` (filters, inactive pills, archived badges), `gray-200` (borders, dividers)
- 10% accent: `emerald-600` — reserved exclusively for: primary CTA buttons, active tab underline, active bottom nav icon, focus rings, active plan badge text, nutritional total highlight

---

## Component Inventory

### Reused from Phase 1/2 (no changes)

| Component | Auto-import alias | Notes |
|-----------|------------------|-------|
| `AppButton.vue` | `<UiAppButton>` | primary/secondary/danger × sm/md/lg |
| `AppInput.vue` | `<UiAppInput>` | label, error, inputDir support |
| `LoadingSpinner.vue` | `<UiLoadingSpinner>` | sm/md/lg sizes |
| `BottomNav.vue` | `<UiBottomNav>` | fixed bottom, exists in both layouts |

### New Components (Phase 3)

All placed in `frontend/app/components/plan/` (auto-imported as `<PlanXxx>`).

| Component | File | Description |
|-----------|------|-------------|
| `PlanStatusBadge` | `plan/StatusBadge.vue` | Shows draft/active/archived pill |
| `PlanBreadcrumb` | `plan/Breadcrumb.vue` | Drill-down breadcrumb + back button bar |
| `PlanDayCard` | `plan/DayCard.vue` | Day row in plan overview (day number, label, meal count) |
| `PlanMealCard` | `plan/MealCard.vue` | Meal row in day view (title, time, option count, reorder buttons) |
| `PlanOptionCard` | `plan/OptionCard.vue` | Option accordion in meal view (expandable, nutrition badges) |
| `PlanFoodItemRow` | `plan/FoodItemRow.vue` | Single food item inside an option (name, qty, unit, delete) |
| `PlanExerciseCard` | `plan/ExerciseCard.vue` | Exercise recommendation card (name, duration, calories) |
| `PlanMedicationCard` | `plan/MedicationCard.vue` | Medication prescription card (name, dosage, times) |
| `PlanNutritionBadges` | `plan/NutritionBadges.vue` | Inline row of computed nutrition values |
| `PlanDayTabBar` | `plan/DayTabBar.vue` | Horizontal scrollable day tabs (client view) |
| `PlanFoodPickerSheet` | `plan/FoodPickerSheet.vue` | Bottom-sheet food search + add item |
| `PlanActivateModal` | `plan/ActivateModal.vue` | Confirmation modal for plan activation |
| `PlanWaterBadge` | `plan/WaterBadge.vue` | Water target display chip in client plan header |

---

## Screen Specifications

---

### Screen 1: Plan Header Form

**Route:** `/nutritionist/clients/:clientId/plans/new`  
**Purpose:** Create a new diet plan (DIET-01). Plan saved in 'draft' status.  
**Layout:** `layout: 'nutritionist'` + `middleware: ['role-guard']` + `roles: ['nutritionist']`

#### Layout Structure

```
┌─────────────────────────────┐
│ [← برگشت]  برنامه‌ی جدید    │  ← Page header (p-4, text-xl font-bold)
├─────────────────────────────┤
│                             │
│  ┌───────────────────────┐  │
│  │ اطلاعات پایه           │  │  ← Section card (bg-white rounded-2xl p-4 shadow-sm)
│  │ [تاریخ شروع]          │  │
│  │ [تاریخ پایان]         │  │
│  │ [هدف آب روزانه (ml)]  │  │
│  └───────────────────────┘  │
│                             │
│  ┌───────────────────────┐  │
│  │ یادداشت               │  │
│  │ [textarea]            │  │
│  └───────────────────────┘  │
│                             │
│  [ذخیره و ادامه →]          │  ← Primary button (full width)
│                             │
└─────────────────────────────┘
```

#### Component Composition

```vue
<div class="p-4 space-y-4">
  <!-- Back + Title -->
  <header class="flex items-center gap-3">
    <button @click="navigateTo(`/nutritionist/clients/${clientId}`)"
            class="text-gray-500 min-h-[44px] min-w-[44px] flex items-center justify-center">
      ← <!-- RTL: right arrow visually = back -->
    </button>
    <h1 class="text-xl font-bold text-gray-800">برنامه‌ی جدید</h1>
  </header>

  <!-- Basic Info Card -->
  <section class="bg-white rounded-2xl p-4 shadow-sm space-y-4">
    <h2 class="text-sm font-semibold text-gray-700">اطلاعات پایه</h2>
    <UiAppInput label="تاریخ شروع (شمسی)" v-model="form.startDate" ... />
    <UiAppInput label="تاریخ پایان (شمسی)" v-model="form.endDate" ... />
    <UiAppInput label="هدف آب روزانه (میلی‌لیتر)" v-model="form.waterTarget" type="number" inputDir="ltr" ... />
  </section>

  <!-- Notes Card -->
  <section class="bg-white rounded-2xl p-4 shadow-sm space-y-4">
    <h2 class="text-sm font-semibold text-gray-700">یادداشت</h2>
    <textarea rows="4" class="w-full rounded-lg border border-gray-300 px-3 py-2.5 ..." />
  </section>

  <UiAppButton type="submit" :loading="isSubmitting">ذخیره و ادامه</UiAppButton>
</div>
```

#### Date Input Pattern

Dates entered as Shamsi (e.g. `۱۴۰۳/۰۱/۱۵`). Use `<UiAppInput type="text" inputDir="ltr" placeholder="مثال: ۱۴۰۳/۰۱/۰۱">`. On submit, convert to Gregorian via `useShamsiDate().toGregorian()` before sending to API. Display validation error if invalid Shamsi format.

#### Form Validation Rules

| Field | Rule | Error Message |
|-------|------|---------------|
| تاریخ شروع | Required, valid Shamsi date | «تاریخ شروع الزامی است» / «تاریخ نامعتبر است» |
| تاریخ پایان | Required, valid Shamsi, ≥ start date | «تاریخ پایان باید بعد از شروع باشد» |
| هدف آب | Optional, integer > 0 if provided | «هدف آب باید عدد مثبت باشد» |
| یادداشت | Optional, max 1000 chars | «یادداشت نباید بیشتر از ۱۰۰۰ کاراکتر باشد» |

#### State

- `isSubmitting: boolean`
- `form: { startDate, endDate, waterTargetMl, notes }`
- `errors: { startDate, endDate, waterTargetMl, notes }`
- On success: `navigateTo('/nutritionist/clients/:clientId/plans/:planId')` (server returns plan ID)

#### Empty / Error States

- Network error: red toast (top, fixed, `z-50`) "ذخیره برنامه با خطا مواجه شد"
- Validation errors: inline below each field (red-600, text-sm)

---

### Screen 2: Plan Overview (Nutritionist)

**Route:** `/nutritionist/clients/:clientId/plans/:planId`  
**Purpose:** Hub page for plan builder. Shows days list, medication prescriptions, and activation CTA.  
**Layout:** `layout: 'nutritionist'`

#### Layout Structure

```
┌─────────────────────────────┐
│ PlanBreadcrumb              │  ← fixed top breadcrumb bar (bg-white border-b)
├─────────────────────────────┤
│                             │
│  ┌──── Plan Header Card ──┐ │
│  │ برنامه‌ی علی رضایی      │ │
│  │ ۱۴۰۳/۰۱/۰۱ تا ۰۷/۰۱  │ │
│  │ 🚰 ۲۵۰۰ میلی‌لیتر آب  │ │
│  │ [پیش‌نویس] badge       │ │
│  └───────────────────────┘ │
│                             │
│  ┌──── روزها ─────────────┐ │  ← Days section card
│  │ روز ۱ — ۳ وعده    ←  │ │  ← PlanDayCard rows, tappable
│  │ روز ۲ — ۲ وعده    ←  │ │
│  │ [+ افزودن روز]         │ │
│  └───────────────────────┘ │
│                             │
│  ┌──── داروها ─────────────┐│  ← Medications section card
│  │ PlanMedicationCard      ││
│  │ [+ افزودن دارو]         ││
│  └───────────────────────┘ │
│                             │
│  [فعال‌سازی برنامه]         │  ← Primary CTA (only if draft, has ≥1 day)
│  [حذف برنامه]              │  ← Danger button (only if draft)
│                             │
└─────────────────────────────┘
```

#### PlanBreadcrumb Component

Fixed bar at top (below device status bar, `sticky top-0 z-40`):

```
bg-white border-b border-gray-200 px-4 py-3

← [کلاینت‌ها] / [نام کلاینت] / برنامه‌ی جدید
```

- Each crumb tappable → navigates back
- Back arrow `←` on far start side (RTL: visually on right)
- Font: `text-sm text-gray-500`, active crumb: `text-gray-800 font-medium`
- Min height: 44px

#### Plan Header Card

```
bg-white rounded-2xl p-4 shadow-sm
```

- Plan date range: `{shamsiStart} تا {shamsiEnd}` — `text-sm text-gray-600`
- Water target: `🚰 {toPersianDigits(waterMl)} میلی‌لیتر` — `text-sm text-gray-600`
- Status: `<PlanStatusBadge :status="plan.status" />`
- Edit icon tap → navigate to `/nutritionist/clients/:clientId/plans/:planId/edit`

#### PlanDayCard Component

```
bg-white rounded-xl px-4 py-3 border-b border-gray-100 last:border-0
flex items-center justify-between
```

- Label start: `روز {toPersianDigits(dayNumber)}` — `text-base font-medium text-gray-800`
- Secondary: `{toPersianDigits(mealCount)} وعده` — `text-xs text-gray-500 ms-2`
- Chevron end: `›` — `text-gray-400 text-xl` (RTL: points left, tap navigates to day view)
- Min height: 56px

"+ افزودن روز" button: `variant="secondary"` full-width at bottom of days list.

#### PlanMedicationCard Component

```
bg-gray-50 rounded-xl px-4 py-3 border border-gray-100
```

- Medication name: `text-base font-medium text-gray-800`
- Dosage row: `text-sm text-gray-600` — `{dosage} — {frequency}`
- Times chips: `text-xs bg-emerald-50 text-emerald-700 rounded-full px-2 py-0.5` for each time
- Edit/Delete: icon buttons at end, `min-h-[44px] min-w-[44px]`

#### Activation Flow

1. Tap "فعال‌سازی برنامه" → `<PlanActivateModal>` appears
2. Modal shows warning: "فعال کردن این برنامه، برنامه‌ی فعال قبلی را آرشیو می‌کند"
3. Two buttons: "فعال‌سازی" (primary) + "انصراف" (secondary)
4. On confirm → PATCH `/api/diet-plans/:id/activate` → plan status updates → activation button disappears
5. If plan is incomplete: API returns Persian error displayed as red toast

#### States

| State | UI |
|-------|----|
| Loading | `<UiLoadingSpinner size="lg">` centered, `py-16` |
| No days | Empty state card: "هنوز روزی اضافه نشده" + "افزودن روز" button |
| No medications | Empty state: "دارویی تجویز نشده" + "افزودن دارو" |
| Draft plan | Activation button visible + delete button visible |
| Active/Archived plan | Activation button hidden, delete button hidden, all cards read-only |

---

### Screen 3: Day View

**Route:** `/nutritionist/clients/:clientId/plans/:planId/days/:dayId`  
**Purpose:** Manage meals and exercise recommendations for a single plan day.

#### Layout Structure

```
┌─────────────────────────────┐
│ PlanBreadcrumb              │  ← 3-level: کلاینت‌ها / نام کلاینت / روز ۱
├─────────────────────────────┤
│                             │
│  Day header                 │  ← "روز ۱ — ۱۴۰۳/۰۱/۰۱" (Shamsi date)
│                             │
│  ┌──── وعده‌ها ─────────────┐│
│  │ PlanMealCard (reorder ↑↓)││
│  │ PlanMealCard             ││
│  │ [+ افزودن وعده]          ││
│  └───────────────────────┘ │
│                             │
│  ┌──── ورزش‌ها ─────────────┐│
│  │ PlanExerciseCard         ││
│  │ [+ افزودن ورزش]          ││
│  └───────────────────────┘ │
│                             │
└─────────────────────────────┘
```

#### PlanMealCard Component

```
bg-white rounded-xl p-4 shadow-sm
```

- Row layout: `flex items-center justify-between`
- Start group:
  - Title: `text-base font-medium text-gray-800`
  - Time (if set): `text-xs text-gray-500 mt-0.5`
  - Option count: `text-xs text-gray-500` — `{toPersianDigits(n)} گزینه`
- End group: reorder buttons + chevron
  - `↑` button: `text-gray-400 min-h-[44px] min-w-[44px]` — disabled if first item
  - `↓` button: same — disabled if last item
  - `›` chevron: tapping navigates to meal view
- Tap on card body (not buttons): navigate to meal view

#### PlanExerciseCard Component

```
bg-gray-50 rounded-xl px-4 py-3 border border-gray-100
```

- Name: `text-base font-medium text-gray-800`
- Duration: `text-sm text-gray-600` — `{toPersianDigits(minutes)} دقیقه`
- Calories burn (if set): `text-xs text-gray-500` — `تخمین: {toPersianDigits(cal)} کالری`
- Description (if set): `text-sm text-gray-600 mt-1`
- Edit/Delete icon buttons at end

#### Add Meal Form (inline sheet or navigate-to-new page)

Implementation: navigate to `/nutritionist/clients/:clientId/plans/:planId/days/:dayId/meals/new`  
Form fields: title (required), scheduled_time (optional, time picker), display_order (auto-assigned).

#### Add Exercise Form

Implementation: inline bottom sheet (slide-up `Transition` — see Animation Specs).  
Form fields: exercise_name, duration_minutes, description, calories_burn_estimate (all in AppInput).

#### Day Label Display

```
"روز {toPersianDigits(day.day_number)}" + Shamsi date computed as:
  plan.start_date + (day.day_number - 1) days → converted via useShamsiDate()
```

---

### Screen 4: Meal View (with Food Picker)

**Route:** `/nutritionist/clients/:clientId/plans/:planId/days/:dayId/meals/:mealId`  
**Purpose:** Manage meal options and food items. This is the deepest level of the drill-down. Most complex screen.

#### Layout Structure

```
┌─────────────────────────────┐
│ PlanBreadcrumb (4-level)    │
├─────────────────────────────┤
│                             │
│  Meal header                │  ← title + time + nutrition totals for whole meal
│  PlanNutritionBadges        │  ← live computed across all options
│                             │
│  ┌──── گزینه ۱ ────────────┐│  ← PlanOptionCard (accordion)
│  │ ▼ [open]                 ││
│  │  PlanFoodItemRow × N    ││
│  │  PlanNutritionBadges    ││
│  │  [+ افزودن ماده غذایی]  ││
│  └───────────────────────┘ │
│                             │
│  ┌──── گزینه ۲ ────────────┐│  ← collapsed
│  │ ▶ [closed]               ││
│  └───────────────────────┘ │
│                             │
│  [+ افزودن گزینه]           │  ← adds new option
│                             │
└─────────────────────────────┘

  ┌─────────────────────────────┐  ← FoodPickerSheet (slide-up, fixed)
  │ [جستجوی غذا...]             │
  │ FoodResultRow × N          │
  │ [افزودن] quantity + unit   │
  └─────────────────────────────┘
```

#### PlanOptionCard Component (Accordion)

```vue
<!-- Collapsed state -->
<div class="bg-white rounded-2xl p-4 shadow-sm">
  <button @click="toggle" class="flex items-center justify-between w-full min-h-[44px]">
    <div class="flex items-center gap-2">
      <span class="text-base font-medium text-gray-800">گزینه {{ toPersianDigits(option.option_number) }}</span>
      <span v-if="!isOpen" class="text-xs text-gray-500">{{ toPersianDigits(itemCount) }} ماده</span>
    </div>
    <span class="text-gray-400">{{ isOpen ? '▼' : '▶' }}</span>
  </button>

  <!-- Expanded state -->
  <div v-show="isOpen" class="mt-3 space-y-2">
    <PlanFoodItemRow v-for="item in option.items" :key="item.id" :item="item" @delete="deleteItem" />
    <PlanNutritionBadges :totals="computedOptionTotals" class="mt-2" />
    <button @click="openFoodPicker" class="...text-emerald-600 text-sm font-medium mt-2">
      + افزودن ماده غذایی
    </button>
  </div>
</div>
```

- Default state: first option open, others closed
- Reorder: `↑`/`↓` buttons on option header (option_number swap)
- Delete option: trash icon on header, confirmation: `confirm()` browser dialog

#### PlanFoodItemRow Component

```
bg-gray-50 rounded-xl px-3 py-2.5 flex items-center justify-between
```

- Start: food name (`text-sm font-medium text-gray-800`)
- Middle: `{toPersianDigits(qty)} {unitLabel}` — `text-xs text-gray-500`
- End: delete button (🗑 or `×`, `text-red-500`, `min-h-[44px] min-w-[44px]`)
- Notes (if set): `text-xs text-gray-400 mt-0.5`

#### PlanNutritionBadges Component

Horizontal scrollable row of computed badges:

```vue
<div class="flex gap-2 overflow-x-auto pb-1 no-scrollbar">
  <span class="badge">کالری: {{ toPersianDigits(Math.round(totals.calories)) }}</span>
  <span class="badge">پروتئین: {{ toPersianDigits(Math.round(totals.protein)) }} گ</span>
  <span class="badge">کربو: {{ toPersianDigits(Math.round(totals.carbs)) }} گ</span>
  <span class="badge">چربی: {{ toPersianDigits(Math.round(totals.fat)) }} گ</span>
  <span class="badge">فیبر: {{ toPersianDigits(Math.round(totals.fiber)) }} گ</span>
</div>
```

Badge style: `bg-emerald-50 text-emerald-700 rounded-full px-2.5 py-1 text-xs font-medium whitespace-nowrap`

#### PlanFoodPickerSheet Component (Bottom Sheet)

```
fixed inset-x-0 bottom-0 z-50 bg-white rounded-t-3xl shadow-2xl
max-height: 80vh (or `max-h-[80vh]`)
padding: p-4 space-y-4
```

Structure:
1. **Handle bar**: `w-10 h-1 bg-gray-300 rounded-full mx-auto mb-3`
2. **Search input**: `<UiAppInput placeholder="جستجوی غذا..." />` — debounced 300ms → calls `GET /api/foods?search=...&limit=20`
3. **Results list**: scrollable, `max-h-[40vh]`, each item is a tappable row
4. **Selected food form** (appears when food tapped):
   - Food name shown as heading
   - Quantity: `<UiAppInput label="مقدار" type="number" inputDir="ltr" />`
   - Unit: `<select>` matching `measurement_unit` enum labels
   - Notes: `<UiAppInput label="یادداشت" />` (optional)
   - Buttons: "افزودن" (primary) + "انصراف" (secondary)

Food Result Row style: `flex items-center justify-between px-3 py-3 border-b border-gray-100 last:border-0 min-h-[56px]`
- Name: `text-base text-gray-800`
- Nutrition summary: `text-xs text-gray-500` — `{calories} کالری / {measurementAmount} {unit}`

Overlay: `fixed inset-0 bg-black/40 z-40` (tap to dismiss sheet)

#### Nutritional Computation (D-14, D-16)

```typescript
// composable: useNutritionComputed.ts
function computeItemNutrition(item: MealOptionItem) {
  const ratio = item.quantity / item.food.measurement_amount
  return {
    calories: item.food.calories * ratio,
    protein:  item.food.protein_g * ratio,
    carbs:    item.food.carbs_g * ratio,
    fat:      item.food.fat_g * ratio,
    fiber:    item.food.fiber_g * ratio,
  }
}

// Option total = sum of all items
const optionTotals = computed(() =>
  option.items.reduce((acc, item) => {
    const n = computeItemNutrition(item)
    acc.calories += n.calories
    acc.protein  += n.protein
    acc.carbs    += n.carbs
    acc.fat      += n.fat
    acc.fiber    += n.fiber
    return acc
  }, { calories: 0, protein: 0, carbs: 0, fat: 0, fiber: 0 })
)

// Meal total = sum of ALL options (nutritionist sees aggregate; client picks one)
// Day total  = sum of ALL meals' totals
```

All values are Vue `computed()` — reactive to `option.items` array changes. No debounce on computation (D-16).

#### Breadcrumb (4 levels)

`کلاینت‌ها ← نام کلاینت ← روز ۱ ← وعده‌ی صبحانه`

---

### Screen 5: Client Plan View

**Route:** `/client/plan`  
**Layout:** `layout: 'client'`  
**Middleware:** `['auth', 'role-guard']` with `roles: ['client']`

#### Layout Structure

```
┌─────────────────────────────┐
│  Plan Header Bar (sticky)   │  ← plan dates + water target
├─────────────────────────────┤
│  PlanDayTabBar (scroll-x)   │  ← "روز ۱" "روز ۲" etc.
├─────────────────────────────┤
│                             │
│  [Active Day Content]       │
│                             │
│  ┌─ وعده‌ی صبحانه ─────────┐ │
│  │ گزینه ۱ [accordion]     │ │
│  │ گزینه ۲ [accordion]     │ │
│  └───────────────────────┘ │
│                             │
│  ┌─ وعده‌ی ناهار ──────────┐ │
│  │ ...                     │ │
│  └───────────────────────┘ │
│                             │
│  ┌─ ورزش‌ها ───────────────┐ │
│  │ PlanExerciseCard × N    │ │
│  └───────────────────────┘ │
│                             │
│  ┌─ داروها ────────────────┐ │  ← Medications (plan-level, shown on last day or fixed)
│  │ PlanMedicationCard × N  │ │
│  └───────────────────────┘ │
│                             │
└─────────────────────────────┘
│         BottomNav           │  ← fixed
└─────────────────────────────┘
```

#### Plan Header Bar (sticky)

```
sticky top-0 z-30 bg-white border-b border-gray-200 px-4 py-3
```

- Plan date range: `text-xs text-gray-500` — `{shamsiStart} تا {shamsiEnd}`
- Water target: `<PlanWaterBadge>` — `🚰 {toPersianDigits(waterMl)} میلی‌لیتر`

#### PlanDayTabBar Component

```
sticky top-[57px] z-20 bg-white border-b border-gray-100
overflow-x-auto whitespace-nowrap px-4 py-0 flex gap-1 no-scrollbar
```

Each tab:
```
inline-flex items-center justify-center px-4 py-3 text-sm min-h-[44px]
border-b-2 transition-colors duration-200
```

- Active: `border-emerald-600 text-emerald-700 font-semibold`
- Inactive: `border-transparent text-gray-500`
- Label: `روز {toPersianDigits(day.day_number)}`
- Auto-scroll to active tab on mount via `scrollIntoView({ behavior: 'smooth', inline: 'center' })`

**Default tab:** computed as `today's day_number` based on `start_date + offset`. Falls back to Day 1 if before plan start or after plan end.

#### Client Meal Card

```
bg-white rounded-2xl shadow-sm overflow-hidden mb-4
```

- Meal header: `px-4 py-3 flex items-center justify-between`
  - Title: `text-base font-semibold text-gray-800`
  - Time (if set): `text-xs text-gray-500`
- Options stacked below header (not inside accordion on desktop, but as accordion on mobile)

#### Client Option Accordion

```vue
<div class="border-t border-gray-100">
  <button @click="toggle" class="w-full flex items-center justify-between px-4 py-3 min-h-[44px]">
    <span class="text-sm font-medium text-gray-700">
      گزینه {{ toPersianDigits(option.option_number) }}
      <span v-if="option.label" class="text-gray-400"> — {{ option.label }}</span>
    </span>
    <span class="text-gray-400">{{ isOpen ? '▼' : '▶' }}</span>
  </button>

  <div v-show="isOpen" class="px-4 pb-3 space-y-2">
    <!-- Food item rows (read-only) -->
    <div v-for="item in option.items" :key="item.id"
         class="flex items-center justify-between py-1.5">
      <span class="text-sm text-gray-800">{{ item.food.name }}</span>
      <span class="text-xs text-gray-500">{{ toPersianDigits(item.quantity) }} {{ unitLabel(item.measurement_unit) }}</span>
    </div>
    <!-- Nutrition totals -->
    <PlanNutritionBadges :totals="computedOptionTotals(option)" class="pt-1" />
  </div>
</div>
```

- Default: all options collapsed. Client taps to expand any option.
- Read-only in Phase 3 (no tracking interaction until Phase 4)

#### Exercise Section (Client)

```
bg-white rounded-2xl shadow-sm p-4 mb-4
```

Section header: `text-sm font-semibold text-gray-700 mb-3` — "ورزش‌های امروز"

Each exercise:
```
flex items-start gap-3 py-2 border-b border-gray-100 last:border-0
```
- Icon: `🏃` (emoji, `text-xl`)
- Name: `text-base font-medium text-gray-800`
- Duration: `text-sm text-gray-600` — `{toPersianDigits(minutes)} دقیقه`
- Calories (if set): `text-xs text-gray-500` — `تخمین مصرف: {toPersianDigits(cal)} کالری`

#### Medications Section (Client, plan-level)

Rendered once at bottom of the day's scroll, regardless of which day is active:

```
bg-white rounded-2xl shadow-sm p-4 mb-4
```

Section header: "داروها"

Each medication:
```
py-3 border-b border-gray-100 last:border-0
```
- Name: `text-base font-medium text-gray-800`
- Dosage + frequency: `text-sm text-gray-600`
- Times: flex-wrap row of time chips — `bg-blue-50 text-blue-700 rounded-full px-2 py-0.5 text-xs`

#### Empty State (No Active Plan)

```vue
<div class="flex flex-col items-center justify-center min-h-[60vh] px-8 text-center space-y-4">
  <span class="text-6xl">📋</span>
  <h2 class="text-lg font-bold text-gray-800">برنامه‌ای فعال ندارید</h2>
  <p class="text-sm text-gray-500">برای دریافت برنامه‌ی غذایی با متخصص تغذیه‌ی خود تماس بگیرید</p>
</div>
```

#### Loading State (Client Plan)

Skeleton loading cards while `GET /api/clients/me/active-plan` is in flight:

```vue
<!-- Skeleton: Day tab bar -->
<div class="flex gap-2 px-4 py-3 overflow-hidden">
  <div v-for="i in 5" class="h-8 w-14 bg-gray-200 rounded-full animate-pulse" />
</div>
<!-- Skeleton: Meal cards -->
<div v-for="i in 3" class="bg-white rounded-2xl shadow-sm p-4 mb-4 space-y-3">
  <div class="h-5 bg-gray-200 rounded w-1/2 animate-pulse" />
  <div class="h-4 bg-gray-100 rounded w-3/4 animate-pulse" />
  <div class="h-4 bg-gray-100 rounded w-2/3 animate-pulse" />
</div>
```

---

## State Management (Pinia)

### Store 1: `usePlanBuilderStore`

**File:** `frontend/app/stores/planBuilder.ts`

```typescript
// State
interface PlanBuilderState {
  // Current plan being built
  plan: DietPlan | null
  days: PlanDay[]
  
  // Loading / error
  loading: boolean
  saving: boolean
  error: string | null
  
  // Food picker
  foodPickerOpen: boolean
  foodPickerTargetOptionId: string | null
  foodSearchQuery: string
  foodSearchResults: Food[]
  foodSearchLoading: boolean
  
  // Inline form state
  activeMealId: string | null
  activeOptionId: string | null
}

// Key Actions
actions: {
  fetchPlan(planId: string): Promise<void>        // GET /api/diet-plans/:id
  createPlan(clientId, payload): Promise<string>  // POST /api/diet-plans → returns planId
  updatePlan(planId, payload): Promise<void>       // PATCH /api/diet-plans/:id
  activatePlan(planId): Promise<void>              // PATCH /api/diet-plans/:id/activate

  addDay(planId): Promise<void>                    // POST /api/diet-plans/:id/days
  reorderMeals(dayId, mealIds): Promise<void>      // PUT (reorder)
  addMeal(dayId, payload): Promise<void>
  updateMeal(mealId, payload): Promise<void>
  deleteMeal(mealId): Promise<void>

  addOption(mealId): Promise<void>
  deleteOption(optionId): Promise<void>

  openFoodPicker(optionId: string): void
  closeFoodPicker(): void
  searchFoods(query: string): Promise<void>        // GET /api/foods?search=...
  addFoodItem(optionId, foodId, qty, unit, notes): Promise<void>
  deleteFoodItem(itemId): Promise<void>

  addExercise(dayId, payload): Promise<void>
  deleteExercise(exerciseId): Promise<void>

  addMedication(planId, payload): Promise<void>
  deleteMedication(medicationId): Promise<void>
}
```

**Important:** The plan builder store holds the **full plan aggregate** in memory after `fetchPlan()`. Individual CRUD actions (addMeal, addFoodItem etc.) optimistically update local state and sync to API. On error, re-fetch to restore consistency.

### Store 2: `useClientPlanStore`

**File:** `frontend/app/stores/clientPlan.ts`

```typescript
interface ClientPlanState {
  plan: DietPlanAggregate | null  // Full nested plan with food nutritional data
  loading: boolean
  error: string | null
  activeDayNumber: number         // Currently selected day tab
}

// Computed (defined as store getters)
getters: {
  activeDay: (state) => state.plan?.days.find(d => d.day_number === state.activeDayNumber)
  computedOptionTotals: (state) => (optionId: string) => NutritionTotals  // via useNutritionComputed
  computedMealTotals: (state) => (mealId: string) => NutritionTotals
  computedDayTotals: (state) => NutritionTotals
}

actions: {
  fetchActivePlan(): Promise<void>    // GET /api/clients/me/active-plan
  setActiveDay(dayNumber: number): void
  initActiveDay(): void               // Called on mount: auto-detect today's day
}
```

**Auto-detect today's day:**
```typescript
function initActiveDay() {
  if (!plan.value) return
  const startDate = new Date(plan.value.start_date)
  const today = new Date()
  const diffDays = Math.floor((today.getTime() - startDate.getTime()) / 86400000)
  const dayNumber = diffDays + 1
  const totalDays = plan.value.days.length
  activeDayNumber.value = Math.max(1, Math.min(dayNumber, totalDays))
}
```

### Composable: `useNutritionComputed`

**File:** `frontend/app/composables/useNutritionComputed.ts`

```typescript
export function useNutritionComputed() {
  function itemTotals(item: MealOptionItem): NutritionTotals {
    const ratio = item.quantity / item.food.measurement_amount
    return {
      calories: item.food.calories * ratio,
      protein:  item.food.protein_g * ratio,
      carbs:    item.food.carbs_g * ratio,
      fat:      item.food.fat_g * ratio,
      fiber:    item.food.fiber_g * ratio,
    }
  }

  function optionTotals(option: MealOption): NutritionTotals {
    return option.items.reduce((acc, item) => {
      const n = itemTotals(item)
      return {
        calories: acc.calories + n.calories,
        protein:  acc.protein + n.protein,
        carbs:    acc.carbs + n.carbs,
        fat:      acc.fat + n.fat,
        fiber:    acc.fiber + n.fiber,
      }
    }, zeroTotals())
  }

  function mealTotals(meal: Meal): NutritionTotals {
    // Sum ALL options (nutritionist sees aggregate; client picks one per tracking)
    return meal.options.reduce((acc, opt) => sumTotals(acc, optionTotals(opt)), zeroTotals())
  }

  function dayTotals(day: PlanDay): NutritionTotals {
    return day.meals.reduce((acc, meal) => sumTotals(acc, mealTotals(meal)), zeroTotals())
  }

  return { itemTotals, optionTotals, mealTotals, dayTotals }
}
```

---

## Navigation Architecture

### Nutritionist Plan Builder Routes

```
/nutritionist/clients/:clientId                     ← Client profile (Phase 5)
  └── [tap "برنامه‌ی جدید" button]
      /nutritionist/clients/:clientId/plans/new      ← Screen 1: Plan Header Form
          └── [on save → navigate to plan overview]
      /nutritionist/clients/:clientId/plans/:planId  ← Screen 2: Plan Overview
          └── /days/:dayId                           ← Screen 3: Day View
              └── /meals/new                         ← New meal form (or inline)
              └── /meals/:mealId                     ← Screen 4: Meal View
      /nutritionist/clients/:clientId/plans          ← Plan list (all plans for client)
```

### Back Navigation

Each drill-down level shows a `PlanBreadcrumb` bar. Back button navigates `useRouter().back()`. On direct URL entry, breadcrumb constructs links from route params.

### Client Routes

```
/client/plan    ← Screen 5: Client Plan View (single page, tab-switched)
```

### Plan Builder Accessed From Client Profile

The nutritionist accesses plan builder via client profile page (`/nutritionist/clients/:clientId`), not from bottom nav. Bottom nav does NOT add a new "برنامه‌ها" tab.

### Route Guards

All nutritionist plan routes:
```typescript
definePageMeta({
  layout: 'nutritionist',
  middleware: ['auth', 'role-guard'],
  roles: ['nutritionist'],
})
```

Client plan route:
```typescript
definePageMeta({
  layout: 'client',
  middleware: ['auth', 'role-guard'],
  roles: ['client'],
})
```

---

## Persian String Constants

### Plan Builder (Nutritionist)

```typescript
const STRINGS = {
  // Plan Header Form
  pageTitle: 'برنامه‌ی جدید',
  fieldStartDate: 'تاریخ شروع (شمسی)',
  fieldEndDate: 'تاریخ پایان (شمسی)',
  fieldWaterTarget: 'هدف آب روزانه (میلی‌لیتر)',
  fieldNotes: 'یادداشت',
  datePlaceholder: 'مثال: ۱۴۰۳/۰۱/۰۱',
  btnSaveAndContinue: 'ذخیره و ادامه',
  
  // Plan Overview
  btnAddDay: '+ افزودن روز',
  btnAddMedication: '+ افزودن دارو',
  btnActivatePlan: 'فعال‌سازی برنامه',
  btnDeletePlan: 'حذف برنامه',
  activateConfirmTitle: 'فعال‌سازی برنامه',
  activateConfirmBody: 'فعال کردن این برنامه، برنامه‌ی فعال قبلی را آرشیو می‌کند. آیا مطمئنید؟',
  activateConfirmBtn: 'فعال‌سازی',
  cancelBtn: 'انصراف',
  emptyDays: 'هنوز روزی اضافه نشده',
  emptyMedications: 'دارویی تجویز نشده',
  sectionDays: 'روزها',
  sectionMedications: 'داروها',
  
  // Day View
  btnAddMeal: '+ افزودن وعده',
  btnAddExercise: '+ افزودن ورزش',
  sectionMeals: 'وعده‌ها',
  sectionExercises: 'ورزش‌ها',
  emptyMeals: 'هنوز وعده‌ای اضافه نشده',
  emptyExercises: 'ورزشی برای این روز ثبت نشده',
  fieldMealTitle: 'عنوان وعده',
  fieldMealTime: 'زمان وعده (اختیاری)',
  
  // Meal View
  btnAddOption: '+ افزودن گزینه',
  btnAddFoodItem: '+ افزودن ماده غذایی',
  optionLabel: (n: number) => `گزینه ${toPersianDigits(n)}`,
  emptyItems: 'هنوز ماده غذایی اضافه نشده',
  
  // Food Picker Sheet
  sheetTitle: 'انتخاب ماده غذایی',
  searchPlaceholder: 'جستجوی غذا...',
  fieldQuantity: 'مقدار',
  fieldUnit: 'واحد',
  fieldItemNotes: 'یادداشت (اختیاری)',
  btnAddItem: 'افزودن',
  searchEmpty: 'نتیجه‌ای یافت نشد',
  searchEmptyHint: 'نام غذا را به فارسی وارد کنید',
  
  // Status Badges
  statusDraft: 'پیش‌نویس',
  statusActive: 'فعال',
  statusArchived: 'آرشیو',
  
  // Nutrition labels
  labelCalories: 'کالری',
  labelProtein: 'پروتئین',
  labelCarbs: 'کربو',
  labelFat: 'چربی',
  labelFiber: 'فیبر',
  unitGram: 'گ',
  
  // Reorder
  btnMoveUp: 'بالاتر',
  btnMoveDown: 'پایین‌تر',
  
  // Errors
  errIncomplete: 'برنامه ناقص است — حداقل یک روز با یک وعده و یک گزینه الزامی است',
  errSaveFailed: 'ذخیره با خطا مواجه شد',
  errLoadFailed: 'بارگذاری برنامه با خطا مواجه شد',
  errDeleteFailed: 'حذف با خطا مواجه شد',
  errStartDate: 'تاریخ شروع الزامی است',
  errEndDate: 'تاریخ پایان باید بعد از شروع باشد',
  errInvalidDate: 'تاریخ نامعتبر است',
  errWaterTarget: 'هدف آب باید عدد مثبت باشد',
  errMealTitle: 'عنوان وعده الزامی است',
  errQuantity: 'مقدار باید بزرگ‌تر از صفر باشد',
  errFormErrors: 'لطفاً خطاهای فرم را اصلاح کنید',
}
```

### Client Plan View

```typescript
const CLIENT_STRINGS = {
  noPlan: 'برنامه‌ای فعال ندارید',
  noPlanHint: 'برای دریافت برنامه‌ی غذایی با متخصص تغذیه‌ی خود تماس بگیرید',
  waterTarget: (ml: number) => `🚰 ${toPersianDigits(ml)} میلی‌لیتر`,
  dayTab: (n: number) => `روز ${toPersianDigits(n)}`,
  sectionExercises: 'ورزش‌های امروز',
  sectionMedications: 'داروها',
  noExercises: 'ورزشی برای امروز تعریف نشده',
  loadError: 'خطا در بارگذاری برنامه',
  btnRetry: 'تلاش مجدد',
  optionLabel: (n: number) => `گزینه ${toPersianDigits(n)}`,
  minuteUnit: 'دقیقه',
  caloriesBurn: (cal: number) => `تخمین مصرف: ${toPersianDigits(cal)} کالری`,
}
```

### Measurement Unit Labels

```typescript
const UNIT_LABELS: Record<string, string> = {
  gram: 'گرم',
  kg: 'کیلوگرم',
  tablespoon: 'قاشق غذاخوری',
  teaspoon: 'قاشق چایخوری',
  cup: 'لیوان',
  piece: 'عدد',
  slice: 'برش',
  palm: 'کف دست',
  matchbox: 'قوطی کبریت',
  bowl: 'کاسه',
  ml: 'میلی‌لیتر',
  liter: 'لیتر',
}
```

---

## Color & Status Palette

### Status Badges (PlanStatusBadge)

```typescript
const STATUS_CLASSES = {
  draft:    'bg-orange-100 text-orange-700',
  active:   'bg-emerald-100 text-emerald-700',
  archived: 'bg-gray-100 text-gray-600',
}
```

Badge element: `inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium`

### Nutritional Badge Colors

All nutritional badges use a single style: `bg-emerald-50 text-emerald-700` — consistent, no per-macro colors.

### Time Chips (Medications)

`bg-blue-50 text-blue-700 rounded-full px-2 py-0.5 text-xs`

### Destructive Actions

- Delete buttons: icon-only `text-red-500` (no background), switches to `text-red-700` on tap
- Danger full button (delete plan): `<UiAppButton variant="danger">` — `bg-red-600 text-white`

---

## Animation Specs

### Page Transitions (Drill-Down)

Use Nuxt `<NuxtPage>` with Vue `<Transition>` for slide-in/slide-out:

```css
/* In main.css or per-page definePageMeta */
.slide-left-enter-active,
.slide-left-leave-active {
  transition: transform 250ms ease, opacity 250ms ease;
}
.slide-left-enter-from { transform: translateX(-20px); opacity: 0; }
.slide-left-leave-to   { transform: translateX(20px);  opacity: 0; }
```

- **Forward navigation** (going deeper): new page slides in from end (RTL: left)
- **Back navigation**: previous page slides in from start (RTL: right)
- Duration: 250ms ease
- No complex animations in v1

### Bottom Sheet (FoodPickerSheet)

```css
.sheet-enter-active, .sheet-leave-active {
  transition: transform 300ms cubic-bezier(0.32, 0.72, 0, 1);
}
.sheet-enter-from, .sheet-leave-to {
  transform: translateY(100%);
}
```

Overlay fade:
```css
.overlay-enter-active, .overlay-leave-active { transition: opacity 300ms; }
.overlay-enter-from, .overlay-leave-to       { opacity: 0; }
```

Trigger: `v-if` on sheet + `<Transition name="sheet">`. Use `<Teleport to="body">` to escape stacking context.

### PlanActivateModal

```css
.modal-enter-active, .modal-leave-active { transition: opacity 200ms, transform 200ms; }
.modal-enter-from, .modal-leave-to       { opacity: 0; transform: scale(0.95); }
```

Centered modal: `fixed inset-0 z-50 flex items-center justify-center px-6`  
Modal card: `bg-white rounded-2xl p-6 shadow-xl w-full max-w-sm`

### Accordion (Options)

```css
/* Use v-show + max-height transition for smooth expand */
.accordion-enter-active, .accordion-leave-active {
  transition: max-height 200ms ease, opacity 200ms;
  overflow: hidden;
}
.accordion-enter-from, .accordion-leave-to { max-height: 0; opacity: 0; }
.accordion-enter-to,   .accordion-leave-from { max-height: 600px; opacity: 1; }
```

Prefer `v-show` over `v-if` for accordion content to preserve reactive computed state.

### Loading Skeleton

`animate-pulse` (Tailwind built-in) on placeholder divs. No third-party skeleton library.

---

## Accessibility Considerations

### Touch Targets

All interactive elements: minimum `min-h-[44px]` and `min-w-[44px]` per Apple HIG. This applies to:
- Breadcrumb back button
- Up/Down reorder buttons
- Accordion toggle buttons
- Tab bar items
- Food item delete buttons

### Focus Management

- When FoodPickerSheet opens: focus moves to search input (`el.focus()` in `onMounted` of sheet)
- When FoodPickerSheet closes: focus returns to "افزودن ماده غذایی" trigger button
- When ActivateModal opens: focus moves to first button (Confirm or Cancel)
- `tabindex` ordering: RTL-correct (end → start in visual layout, but tab order follows DOM order which matches reading order)

### ARIA

- `<PlanBreadcrumb>`: wrap in `<nav aria-label="مسیر">`, each crumb in `<ol>` + `<li>`
- `<PlanDayTabBar>`: `role="tablist"`, each tab `role="tab"` with `aria-selected`
- `<PlanOptionCard>` accordion: `aria-expanded` on button, `aria-controls` on content panel
- `<PlanFoodPickerSheet>`: `role="dialog"`, `aria-modal="true"`, `aria-label="انتخاب ماده غذایی"`
- Status badges: no role needed (presentational) — parent context provides semantic meaning

### RTL-Safe Interactions

- Chevron direction: use `›` (U+203A) in RTL — it visually points left (towards deeper content) without CSS transforms
- Back arrow: use `←` (U+2190) — always points left, appropriate for back in RTL
- Horizontal swipe gestures: NOT implemented in Phase 3 (no touch-swipe for day navigation). Tab bar tap only.

### Color Contrast

- `text-emerald-700` on `bg-emerald-50`: ratio ≈ 5.2:1 ✓ (WCAG AA)
- `text-orange-700` on `bg-orange-100`: ratio ≈ 4.7:1 ✓ (WCAG AA)
- `text-gray-600` on `bg-gray-100`: ratio ≈ 4.6:1 ✓ (WCAG AA)
- `text-gray-800` on `bg-white`: ratio ≈ 14:1 ✓
- `text-white` on `bg-emerald-600`: ratio ≈ 4.5:1 ✓ (WCAG AA)

---

## Component File Map

```
frontend/app/
├── pages/
│   ├── nutritionist/
│   │   └── clients/
│   │       └── [clientId]/
│   │           └── plans/
│   │               ├── index.vue          ← Plan list for client (DIET-11, D-28)
│   │               ├── new.vue            ← Screen 1: Plan Header Form
│   │               └── [planId]/
│   │                   ├── index.vue      ← Screen 2: Plan Overview
│   │                   ├── edit.vue       ← Plan Header Edit
│   │                   └── days/
│   │                       └── [dayId]/
│   │                           ├── index.vue    ← Screen 3: Day View
│   │                           └── meals/
│   │                               ├── new.vue  ← New Meal Form
│   │                               └── [mealId]/
│   │                                   └── index.vue  ← Screen 4: Meal View
│   └── client/
│       └── plan.vue                       ← Screen 5: Client Plan View
│
├── components/
│   ├── plan/
│   │   ├── StatusBadge.vue
│   │   ├── Breadcrumb.vue
│   │   ├── DayCard.vue
│   │   ├── MealCard.vue
│   │   ├── OptionCard.vue
│   │   ├── FoodItemRow.vue
│   │   ├── ExerciseCard.vue
│   │   ├── MedicationCard.vue
│   │   ├── NutritionBadges.vue
│   │   ├── DayTabBar.vue
│   │   ├── FoodPickerSheet.vue
│   │   ├── ActivateModal.vue
│   │   └── WaterBadge.vue
│   └── ui/                               ← existing (no changes)
│
├── stores/
│   ├── planBuilder.ts                    ← new
│   └── clientPlan.ts                     ← new
│
└── composables/
    └── useNutritionComputed.ts           ← new
```

---

## Registry

**No shadcn. No third-party component registries.** All components hand-rolled in the project's established Tailwind v4 pattern. Safety gate: not applicable.

---

*Phase: 03-diet-plan-engine*  
*UI-SPEC generated: auto-chain, 2026-04-19*  
*Upstream sources: 03-CONTEXT.md (D-17..D-28), REQUIREMENTS.md (DIET-01..12), Phase 1/2 established patterns*
