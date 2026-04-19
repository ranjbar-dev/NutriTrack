# Phase 2: Core Data Domain - Context

**Gathered:** 2025-01-20
**Status:** Ready for planning

<domain>
## Phase Boundary

Nutritionists and Super Admin can populate and manage the shared food and medication databases with full Persian search capability, and Super Admin has complete platform operational control including nutritionist account management and platform statistics dashboard. This phase delivers the data foundation that the diet plan engine (Phase 3) will consume.

</domain>

<decisions>
## Implementation Decisions

### Database Schema & Architecture
- **D-01:** Food and medication tables use soft delete (`is_active` boolean, default true) — never hard delete. Deactivating a food/medication does not cascade-delete existing diet plan references.
- **D-02:** Food categories stored in a separate junction table (`food_categories`) for many-to-many relationship. Category enum in database: `breakfast`, `lunch`, `dinner`, `snack`, `fruit`, `beverage`, `supplement`, `other`.
- **D-03:** Audit trail: `created_by` (UUID FK to users) and `created_at` tracked on all food and medication items. Super Admin can see who created each item; nutritionists only see their own created items unless Super Admin.
- **D-04:** Food measurement units as PostgreSQL enum: `gram`, `kg`, `tablespoon`, `teaspoon`, `cup`, `piece`, `slice`, `palm`, `matchbox`, `bowl`, `ml`, `liter`. Each food item has a default `measurement_unit` and `measurement_amount` (base quantity).
- **D-05:** Nutritional fields stored as DECIMAL(8,2) for precision: `calories`, `protein_g`, `carbs_g`, `fat_g`, `fiber_g`, `sugar_g`, `sodium_mg`. All per the base measurement amount.

### Persian Search Implementation
- **D-06:** pg_trgm extension enabled on PostgreSQL for fuzzy Persian search. Trigram GIN index on `LOWER(name)` column for both food and medication tables.
- **D-07:** Persian character normalization handled at TWO boundaries: (1) storage — before INSERT/UPDATE via a PostgreSQL function `normalize_persian(text)` that converts `ی→ي` and `ک→ك`; (2) query — search terms normalized in the same way before ILIKE or similarity queries.
- **D-08:** Search query pattern: `WHERE LOWER(normalize_persian(name)) ILIKE LOWER(normalize_persian('%' || $search_term || '%'))` combined with pg_trgm similarity scoring for ranking.
- **D-09:** Docker PostgreSQL container uses `fa_IR.UTF-8` locale to ensure correct Persian text collation and sorting. Verified during Phase 2 execution (noted as Phase 2 blocker in STATE.md).

### Food Management UI/UX
- **D-10:** Food list page displays as a scrollable card-based list (not table) — mobile-first, each card shows: name (bold, large), categories (small badges), calories + protein (summary line), edit/delete icons.
- **D-11:** Add/Edit food form is a multi-section form: (1) Basic Info (name, description), (2) Categories (checkbox group), (3) Nutrition (6 fields in 2-column grid), (4) Measurement (unit dropdown + amount input). Auto-save draft to localStorage on input change.
- **D-12:** Category filter: horizontal scrollable pill buttons above food list. "All" (default) + 8 category pills. Active/inactive toggle as a separate filter (default: active only).
- **D-13:** Search bar with 300ms debounce at the top of food list. Searches as-you-type with Persian normalization. Placeholder: "جستجوی غذا..." (Search food...).
- **D-14:** Pagination: 20 items per page with "Load More" button at bottom (infinite scroll pattern) instead of traditional page numbers. Scroll position preserved when navigating back.
- **D-15:** Empty states: (a) No foods exist: "هیچ غذایی ثبت نشده" + "افزودن غذا" button, (b) No search results: "نتیجه‌ای یافت نشد" + "پاک کردن جستجو" button, (c) No foods in selected category: "در این دسته غذایی وجود ندارد".

### Medication Management UI/UX
- **D-16:** Medication list follows same card-based pattern as food list. Each card shows: name, generic name (if provided), form (badge), edit/delete icons.
- **D-17:** Add/Edit medication form: (1) Names (commercial name, generic name), (2) Form (dropdown: tablet/capsule/syrup/injection/drop/powder/other), (3) Dosage unit (text field, e.g., "mg", "ml"), (4) Description (optional).
- **D-18:** Search bar with Persian normalization, same 300ms debounce. Searches both commercial name and generic name fields.
- **D-19:** Pagination: 20 items per page, "Load More" button, same as food list.

### Super Admin Panel
- **D-20:** Super Admin dashboard at `/admin` shows 4 stat cards: Total Nutritionists (active count), Total Clients (all), Total Foods (active), Total Active Diet Plans. Stats refresh on page load, no real-time updates.
- **D-21:** Nutritionist management page at `/admin/nutritionists` — table/card list showing: name, email, status (active/inactive badge), client count, created_at date. Actions: Activate, Deactivate, View Clients (read-only).
- **D-22:** Create Nutritionist form: email, password, full name, phone (optional). Password must be 8+ characters. Email uniqueness validated server-side.
- **D-23:** "View Clients" for a nutritionist opens a read-only modal/page showing that nutritionist's client list (name, mobile, status, plan status). No CRUD actions — Super Admin cannot edit nutritionist clients directly.
- **D-24:** Super Admin can access all food and medication items from any user. Food/medication lists show `created_by` name for audit. Super Admin can edit/delete any item.
- **D-25:** Activate/Deactivate nutritionist is a toggle button. When deactivated, nutritionist cannot log in (validated at auth service layer). Clients under that nutritionist remain in the system.

### API Structure & Routing
- **D-26:** Food CRUD routes under `/api/foods`: `GET /api/foods` (list with search/filter/pagination), `POST /api/foods` (create), `GET /api/foods/:id` (single), `PUT /api/foods/:id` (update), `DELETE /api/foods/:id` (soft delete).
- **D-27:** Medication CRUD routes under `/api/medications`: `GET /api/medications`, `POST /api/medications`, `GET /api/medications/:id`, `PUT /api/medications/:id`, `DELETE /api/medications/:id`.
- **D-28:** Super Admin routes under `/api/admin`: `GET /api/admin/stats` (platform stats), `GET /api/admin/nutritionists` (list), `POST /api/admin/nutritionists` (create), `PATCH /api/admin/nutritionists/:id/status` (activate/deactivate), `GET /api/admin/nutritionists/:id/clients` (read-only client list).
- **D-29:** Row-level authorization enforced at repository layer: Nutritionists can only soft-delete their own created foods/medications (check `created_by = current_user_id`). Super Admin bypasses this check (can delete any).
- **D-30:** Search/filter API accepts query params: `?search=term&category=breakfast&is_active=true&page=1&limit=20`. Returns JSON: `{data: [...], total: N, page: 1, limit: 20, has_more: bool}`.

### Validation & Error Handling
- **D-31:** Persian text fields validated for non-empty after trimming. Food name required (max 200 chars), description optional (max 1000 chars). Medication name required (max 200 chars), generic name optional (max 200 chars).
- **D-32:** Nutritional values validated: must be >= 0, max 2 decimal places. Calories max 9999.99, protein/carbs/fat/fiber/sugar max 999.99, sodium max 9999.99.
- **D-33:** Duplicate food name check: case-insensitive + Persian-normalized. Show error: "غذا با این نام قبلاً ثبت شده است" (Food with this name already exists). Same for medication.
- **D-34:** Form-level errors displayed as Persian toast notifications (top-center). Field-level errors inline below each field in red.

### Loading & Performance
- **D-35:** Food list initial load target: <300ms for 20 items with search/filter. Use `SELECT * FROM foods WHERE ... LIMIT 20 OFFSET 0` with indexed columns.
- **D-36:** Search debounced at 300ms to prevent excessive API calls. Show loading spinner in search bar during search. No skeleton loader for search results (existing list stays visible).
- **D-37:** Create/Update operations show blocking loading state on submit button ("در حال ذخیره..."). Disable form inputs during save.

### Agent's Discretion
- Exact card component styling (shadows, border radius, padding)
- Toast notification library choice (or custom implementation)
- Stat card icon choices for Super Admin dashboard
- Exact color palette for category badges
- Form field order and visual grouping details
- Loading spinner animation style
- Horizontal pill scroll behavior (snap or free scroll)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requirements
- `.planning/REQUIREMENTS.md` §Food Database (FOOD-01 through FOOD-10) — Food CRUD, categories, search, pagination, audit requirements
- `.planning/REQUIREMENTS.md` §Medication Database (MED-01 through MED-05) — Medication CRUD, forms, search requirements
- `.planning/REQUIREMENTS.md` §Super Admin (ADMIN-01 through ADMIN-08) — Super Admin platform management, nutritionist CRUD, statistics requirements
- `.planning/ROADMAP.md` §Phase 2 — Phase goals, success criteria, and dependencies

### Architecture
- `.planning/research/ARCHITECTURE.md` — Layered backend architecture patterns to follow
- `.planning/research/STACK.md` — Library versions and installation commands
- `.planning/phases/01-foundation-infrastructure/01-CONTEXT.md` — Established patterns: handler→service→repository, sqlc queries, Persian utilities, mobile-first UI conventions

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- **Handler→Service→Repository pattern**: Established in Phase 1. `handler/auth_handler.go`, `service/auth_service.go`, `repository/user_repo.go` show the pattern. Food and medication handlers will follow the same structure.
- **Persian error messages**: All error responses in Phase 1 use Persian text (e.g., `"اطلاعات ورودی نامعتبر است"`). Continue this for all food/medication/admin errors.
- **DTO structs with validation tags**: `dto.LoginRequest`, `dto.OTPRequestDTO` use `validator` tags. Food/medication DTOs will use the same approach.
- **UI components**: `AppButton.vue`, `AppInput.vue`, `LoadingSpinner.vue`, `BottomNav.vue` available for reuse. Food forms can use `AppInput` for text fields, `AppButton` for submit.
- **Pinia store pattern**: `stores/auth.ts` shows the pattern (state, getters, actions, API calls). Food and medication stores will follow the same structure.
- **Middleware stack**: `middleware/auth.go`, `middleware/role_guard.go` already handle JWT and role checks. Food/medication routes will use these.

### Established Patterns
- **sqlc for all queries**: Phase 1 uses `db/queries/users.sql`, `db/queries/otp.sql` compiled to type-safe Go code. Food and medication queries go in `db/queries/foods.sql` and `db/queries/medications.sql`.
- **Soft delete**: Not used in Phase 1, but specified for Phase 2. Pattern: add `is_active BOOLEAN DEFAULT TRUE` column, never DELETE, always `UPDATE SET is_active = FALSE`.
- **Mobile-first layouts**: All Phase 1 pages are mobile-only, no desktop breakpoints. Phase 2 continues this. Bottom nav for all roles.
- **Persian date/number display**: `useShamsiDate()` and `toPersianDigits()` used in Phase 1. Use throughout Phase 2 for any dates/numbers.
- **Role-based routing**: Phase 1 uses `middleware: ['auth', 'role']` in Nuxt pages. Super Admin pages use `middleware: ['auth', 'role.admin']`, nutritionist pages use `middleware: ['auth', 'role.nutritionist']`.

### Integration Points
- **Bottom nav**: Super Admin bottom nav (in `layouts/admin.vue`) will need a new tab for Foods or Stats. Nutritionist bottom nav (in `layouts/nutritionist.vue`) will need tabs for Clients, Foods, Medications.
- **API base URL**: Frontend already uses `apiBase` from `useRuntimeConfig()` for API calls. Food/medication API calls use the same base.
- **Auth store**: `stores/auth.ts` exposes `user` and `isAuthenticated`. Food/medication forms can read `user.id` for `created_by` field.

</code_context>

<specifics>
## Specific Ideas

- **Persian pg_trgm validation**: STATE.md notes this as a Phase 2 blocker/concern. Need to verify Persian trigram search works correctly with fa_IR.UTF-8 locale in Docker. Test with sample Persian food names early in execution.
- **Shared database philosophy**: PROJECT.md Key Decision: "Shared food/medication database across all nutritionists — reduces duplication, builds richer database over time." Food and medication are platform-wide, not per-nutritionist.
- **Super Admin as platform operator**: Super Admin manages the platform (nutritionists, global data) but does NOT create diet plans or manage clients directly. Read-only access to nutritionist clients is for oversight only.
- **Mobile-first cards over tables**: Continuing Phase 1 pattern, all lists (food, medication, nutritionists) display as mobile-friendly cards, not desktop tables.
- **RTL logical properties**: Tailwind v4 logical properties (`ms-`, `me-`, `ps-`, `pe-`) used everywhere, established in Phase 1. No physical LTR properties.

</specifics>

<deferred>
## Deferred Ideas

- **Food photos/images**: Not in Phase 2 requirements. If needed, defer to Phase 3 or backlog.
- **Bulk import for foods**: Useful for seeding a large food database from CSV/Excel. Not in Phase 2 success criteria — defer to backlog.
- **Medication interaction warnings**: Complex feature requiring a drug interaction database. Out of scope for v1 — defer to backlog.
- **Food history/versioning**: Tracking edits to food items over time. Not required for Phase 2 — defer to backlog.
- **Advanced search filters**: Filtering by calorie range, macronutrient ratios, etc. Phase 2 only requires category and active/inactive filters — defer advanced filters to backlog.

</deferred>

---

*Phase: 02-core-data-domain*
*Context gathered: 2025-01-20*
