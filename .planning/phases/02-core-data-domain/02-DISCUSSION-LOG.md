# Phase 2: Discussion Log

**Date:** 2025-01-20
**Mode:** Auto-approved (--auto flag / yolo mode)
**Participants:** System (autonomous decision selection)

## Session Summary

Phase 2 context gathered autonomously with recommended defaults for all implementation decisions. All gray areas were auto-selected and resolved using best practices from Phase 1 patterns, requirements analysis, and established Persian/mobile-first conventions.

## Gray Areas Discussed

### 1. Database Schema & Data Integrity
**Selected approach:** Soft delete with audit trail
- Rationale: Preserves data integrity for existing diet plan references, enables audit compliance
- Food/medication items never hard deleted, marked `is_active = FALSE`
- `created_by` tracking for all items, visible to Super Admin

### 2. Persian Search Implementation
**Selected approach:** pg_trgm with dual normalization (storage + query)
- Rationale: Handles Persian character variants (ی/ي, ک/ك) correctly, provides fuzzy search
- Docker PostgreSQL uses fa_IR.UTF-8 locale (blocker noted in STATE.md, to be validated during execution)
- Trigram GIN index on normalized name column for performance

### 3. Food Management UI Pattern
**Selected approach:** Card-based mobile list with debounced search
- Rationale: Consistent with Phase 1 mobile-first pattern, better UX than tables on mobile
- Horizontal scrollable category filters (pills)
- 300ms search debounce to reduce API load
- "Load More" infinite scroll pattern (20 items per page)

### 4. Medication Management UI Pattern
**Selected approach:** Same card pattern as food, simplified fields
- Rationale: Consistency across data management UIs, less complex than food (no nutrition data)
- Search covers both commercial name and generic name
- Form type as dropdown (tablet/capsule/syrup/etc.)

### 5. Super Admin Panel Design
**Selected approach:** Dashboard with stat cards + nutritionist management table
- Rationale: Provides operational visibility without overwhelming complexity
- 4 stat cards: nutritionists, clients, foods, active plans
- Nutritionist list with activate/deactivate toggle
- Read-only view of nutritionist clients (oversight, not direct management)

### 6. API Structure & Authorization
**Selected approach:** RESTful routes with role-based repository filtering
- Rationale: Follows Phase 1 conventions, clear separation of concerns
- `/api/foods`, `/api/medications`, `/api/admin/*` routes
- Row-level authorization at repository layer (nutritionists can only delete own items, Super Admin can delete any)
- Standard CRUD operations with soft delete semantics

### 7. Form Validation & Error Handling
**Selected approach:** Persian validator tags + client-side normalization
- Rationale: Consistent with Phase 1 validation pattern, clear error messages in Persian
- Duplicate name check (case-insensitive + Persian-normalized)
- Nutritional values validated for range and precision (2 decimal places)
- Field-level inline errors + toast notifications for form-level errors

### 8. Loading & Performance Strategy
**Selected approach:** Indexed search with pagination, 300ms debounce
- Rationale: Balances responsiveness with server load
- Target <300ms for 20-item list load
- Debounced search to prevent excessive API calls
- Loading states on buttons during save operations

## Decisions Carried Forward from Phase 1

- Handler→Service→Repository layered architecture
- sqlc for all database queries (no ORM)
- Mobile-first, no desktop breakpoints
- Persian error messages throughout
- RTL layout with Tailwind v4 logical properties
- Pinia stores for state management
- JWT auth cookies with role-based middleware
- Structured JSON logging

## Canonical References Identified

- `.planning/REQUIREMENTS.md` (FOOD-01 through FOOD-10, MED-01 through MED-05, ADMIN-01 through ADMIN-08)
- `.planning/ROADMAP.md` §Phase 2
- `.planning/research/ARCHITECTURE.md`
- `.planning/research/STACK.md`
- `.planning/phases/01-foundation-infrastructure/01-CONTEXT.md`

## Deferred Ideas (Out of Scope)

- Food photos/images → Backlog
- Bulk import for foods from CSV → Backlog
- Medication interaction warnings → Backlog
- Food edit history/versioning → Backlog
- Advanced search filters (calorie range, macro ratios) → Backlog

## Next Steps

Context document created at `.planning/phases/02-core-data-domain/02-CONTEXT.md` with 30 implementation decisions across 8 categories.

**Auto-advancing to:** `/gsd-plan-phase 2` (per --auto workflow)

---

*Discussion completed: 2025-01-20*
*Mode: Autonomous (--auto)*
*Total decisions: 37 (30 implementation + 7 discretion)*
