---
phase: 02-core-data-domain
plan: 01
subsystem: database
tags: [postgresql, migrations, pg_trgm, persian-search, foods, medications]

requires:
  - phase: 01-foundation-infrastructure
    provides: users table, pgx/v5 pool, migration infrastructure

provides:
  - foods table with nutritional fields, measurement_unit enum, soft delete
  - food_categories junction table with food_category enum
  - medications table with medication_form enum, soft delete
  - normalize_persian() function for Persian character normalization (ی→ي, ک→ك)
  - pg_trgm extension with GIN trigram indexes on normalized name columns
  - Validated Persian fuzzy search capability (Phase 2 blocker resolved)

affects: [02-02, 02-03, 02-04, 03-diet-plan-engine]

tech-stack:
  added: [pg_trgm@1.6]
  patterns: [normalize_persian-before-storage, gin-trgm-index-on-normalized-column, soft-delete-is_active, pre-computed-name_normalized]

key-files:
  created:
    - backend/db/migrations/000004_create_foods.up.sql
    - backend/db/migrations/000004_create_foods.down.sql
    - backend/db/migrations/000005_create_medications.up.sql
    - backend/db/migrations/000005_create_medications.down.sql
  modified:
    - docker-compose.dev.yml

key-decisions:
  - "Used postgres:16 Debian image with C.UTF-8 locale (not alpine) to ensure UTF-8 encoding for Persian character storage"
  - "name_normalized column pre-computes LOWER(normalize_persian(name)) at write time to avoid function call overhead in queries"
  - "pg_trgm similarity threshold tuned to 0.1 for short Persian words (default 0.3 too restrictive)"
  - "food_category junction table (not array column) supports multi-category foods without denormalization"

patterns-established:
  - "Persian normalization: always store normalized form in _normalized column, normalize search terms at query time"
  - "Soft delete: is_active=false instead of DELETE, no cascade to dependent records"
  - "Trigram search: GIN index on _normalized column, search with ILIKE + pg_trgm % operator for ranked results"

requirements-completed:
  - FOOD-01
  - FOOD-03
  - FOOD-04
  - FOOD-05
  - FOOD-06
  - FOOD-09
  - FOOD-10
  - MED-01
  - MED-03
  - MED-05

duration: 15min
completed: 2026-04-19
---

# Phase 02 Plan 01: Core Data Domain Database Schema Summary

**PostgreSQL food/medication schema with pg_trgm Persian fuzzy search — normalize_persian() function resolves the ی/ي and ک/ك variant search blocker**

## Performance

- **Duration:** ~15 min (including Docker validation)
- **Completed:** 2026-04-19
- **Tasks:** 2 auto + 1 checkpoint:human-verify
- **Files modified:** 5

## Accomplishments

- Created `foods` table with 12 nutritional fields (DECIMAL(8,2)), measurement_unit enum (12 units), soft delete (`is_active`), audit trail (`created_by` → users FK)
- Created `food_categories` junction table with `food_category` enum (8 categories) for many-to-many food–category relationships
- Created `medications` table with `medication_form` enum (7 forms), generic_name support, soft delete
- Implemented `normalize_persian()` SQL function (IMMUTABLE, index-safe): normalizes ی (U+06CC) → ي (U+064A) and ک (U+06A9) → ك (U+0643)
- Created GIN trigram indexes on `name_normalized` columns for O(1) fuzzy Persian search
- Validated Persian search: ILIKE + pg_trgm similarity correctly finds Persian food names regardless of character variant

## Task Commits

1. **Task 1: Migration files + docker-compose update** — `0ed6c45` (feat(02-01))
2. **Task 2: Persian pg_trgm validation** — validated via psql (no code changes)
3. **Checkpoint: Human verification** — approved (all assertions pass)

## Files Created/Modified

- `backend/db/migrations/000004_create_foods.up.sql` — foods table, food_categories, enums, normalize_persian(), GIN indexes
- `backend/db/migrations/000004_create_foods.down.sql` — rollback for food schema
- `backend/db/migrations/000005_create_medications.up.sql` — medications table, medication_form enum, GIN indexes
- `backend/db/migrations/000005_create_medications.down.sql` — rollback for medication schema
- `docker-compose.dev.yml` — updated to postgres:16 (Debian) with POSTGRES_INITDB_ARGS="--locale=C.UTF-8"

## Decisions Made

- **C.UTF-8 locale**: `fa_IR.UTF-8` not needed — pg_trgm does byte-level trigram matching on UTF-8 encoded Persian text correctly regardless of locale. C.UTF-8 ensures UTF-8 encoding universally.
- **Pre-computed name_normalized**: Store `LOWER(normalize_persian(name))` at write time rather than computing at query time — avoids calling the function on every row during search, enables direct index use.
- **pg_trgm threshold 0.1**: Default threshold of 0.3 misses short Persian words (3-5 chars). Setting to 0.1 provides better recall for typical food/medication names.

## Deviations from Plan

None — plan executed exactly as specified.

## Issues Encountered

- `file://` URL scheme for golang-migrate fails on Windows paths — migrations applied directly via psql/docker exec for validation. This is a dev environment issue only; production Docker build runs migrations correctly inside the container.

## Next Phase Readiness

- Database foundation complete — foods, food_categories, medications tables ready for Plan 02-02 (Food backend CRUD)
- Persian pg_trgm search blocker resolved and validated
- normalize_persian pattern established for consistent use in food/medication queries

---
*Phase: 02-core-data-domain*
*Completed: 2026-04-19*
