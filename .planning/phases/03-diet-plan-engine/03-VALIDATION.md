# Phase 3: Diet Plan Engine — Validation Architecture

**Source:** Extracted from `03-RESEARCH.md` §Validation Architecture  
**Generated:** 2026-04-19 (revision 2026-04-19)  
**Config:** `nyquist_validation: true`

---

## Overview

Phase 3 has no pre-existing test infrastructure. Wave 0 creates all test stub files first — as empty-but-compiling shells — so the automated verify commands in every subsequent plan have a runnable harness. Executors fill stubs with real assertions as they implement each dimension.

**Test framework setup required in Wave 0 (Plan 03-09):**
- Go: standard `testing` package + `github.com/stretchr/testify` (already a common Go dep; add if missing)
- Frontend: Vitest + `@vue/test-utils` (install via `npm install -D vitest @vue/test-utils`)
- Config file: `frontend/vitest.config.ts`

---

## Quick-Run Commands

| Scope | Command |
|-------|---------|
| Backend service tests | `cd backend && go test ./internal/service/... -run TestDietPlan -v -timeout 30s` |
| Backend repo tests | `cd backend && go test ./internal/repository/... -run TestDietPlan -v -timeout 60s` |
| Full backend | `cd backend && go test ./... -timeout 60s` |
| Frontend unit tests | `cd frontend && npx vitest run tests/` |
| Single frontend test | `cd frontend && npx vitest run tests/useNutritionCalc.test.ts` |

---

## Dimension 1: One-Active-Plan Constraint

**Requirement:** DIET-02, D-02  
**Risk if wrong:** Two clients end up with two active plans — data corruption, both display to the client.

**Test file:** `backend/internal/service/diet_plan_service_test.go`

### Test Strategy

```go
// TestActivatePlanArchivesPrevious
// Covers: service-layer ActivatePlan correctly archives prior active plan
func TestActivatePlanArchivesPrevious(t *testing.T) {
    // 1. Create plan A for clientX → directly set status='active'
    // 2. Create plan B for clientX in 'draft' status
    // 3. Call svc.ActivatePlan(ctx, planB.ID, nutritionistID)
    // 4. Re-fetch plan A → assert status == 'archived'
    // 5. Re-fetch plan B → assert status == 'active'
    // 6. Assert no other active plan exists for clientX
}

// TestOneActivePlanConstraintAtDBLevel
// Covers: PostgreSQL partial unique index prevents bypass via raw SQL
func TestOneActivePlanConstraintAtDBLevel(t *testing.T) {
    // 1. Create and activate plan A for clientX
    // 2. Attempt direct SQL INSERT of a second active plan for clientX
    //    INSERT INTO diet_plans (..., status) VALUES (..., 'active')
    // 3. Assert PostgreSQL returns unique violation error
    //    (error message contains "idx_diet_plans_one_active_per_client" or similar)
}
```

### Manual Verification

```sql
-- Should return 0 rows (no client has two active plans)
SELECT client_id, COUNT(*) 
FROM diet_plans 
WHERE status = 'active' 
GROUP BY client_id 
HAVING COUNT(*) > 1;
```

**Automated command:** `cd backend && go test ./internal/service/... -run TestActivatePlan -v -timeout 30s`

---

## Dimension 2: Batch Load Performance (≤500ms SLA)

**Requirement:** DIET-12, INFRA-10, D-36  
**Risk if wrong:** Client plan view loads slowly — UX failure; also potential N+1 query regression.

**Test file:** `backend/internal/repository/diet_plan_repo_test.go`

### Test Strategy

```go
// TestPlanAggregateLoadTime
// Integration test — requires real PostgreSQL connection
func TestPlanAggregateLoadTime(t *testing.T) {
    // SETUP: seed a realistic plan: 7 days × 5 meals × 3 options × 4 items = 420 items
    //        + 35 exercises + 5 medications
    // Use test helpers to insert data (see testutil package)
    
    // MEASURE:
    start := time.Now()
    plan, err := repo.GetFullPlanAggregate(ctx, testPlanID)
    elapsed := time.Since(start)
    
    require.NoError(t, err)
    require.NotNil(t, plan)
    
    // SLA assertion
    assert.Less(t, elapsed, 500*time.Millisecond,
        "plan aggregate must load in ≤500ms, got %v", elapsed)
    
    // Structure completeness assertions
    assert.Len(t, plan.Days, 7)
    assert.Len(t, plan.Days[0].Meals, 5)
    assert.Len(t, plan.Days[0].Meals[0].Options, 3)
    assert.Len(t, plan.Days[0].Meals[0].Options[0].Items, 4)
    assert.Len(t, plan.Medications, 5)
}
```

### Index Verification

```sql
-- Verify index scan (not sequential scan) for batch ANY() queries
EXPLAIN ANALYZE
SELECT * FROM meals WHERE day_id = ANY(ARRAY['<uuid1>', '<uuid2>', '<uuid3>']::uuid[]);
-- Expected: "Index Scan using idx_meals_day_id"

EXPLAIN ANALYZE
SELECT * FROM meal_options WHERE meal_id = ANY(ARRAY['<uuid>']::uuid[]);
-- Expected: "Index Scan using idx_meal_options_meal_id"
```

### Manual Timing

```bash
time curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/diet-plans/$PLAN_ID > /dev/null
# Expected: real < 0m0.500s
```

**Automated command:** `cd backend && go test ./internal/repository/... -run TestPlanAggregate -v -timeout 60s`

---

## Dimension 3: Activation Validation (Incomplete Plan)

**Requirement:** DIET-01, D-33  
**Risk if wrong:** Nutritionist activates an empty plan — client sees blank diet plan.

**Test file:** `backend/internal/service/diet_plan_service_test.go`

### Test Strategy

```go
// TestActivationBlockedForIncompletePlan
func TestActivationValidation(t *testing.T) {
    cases := []struct {
        name    string
        setup   func(ctx context.Context, planID string) // mutates plan to be incomplete
        wantErr string
    }{
        {
            name:    "no days",
            setup:   func(ctx context.Context, planID string) { /* create plan with 0 days */ },
            wantErr: "برنامه ناقص است",
        },
        {
            name: "day with no meals",
            setup: func(ctx context.Context, planID string) {
                // create one day, no meals
            },
            wantErr: "برنامه ناقص است",
        },
        {
            name: "meal with no options",
            setup: func(ctx context.Context, planID string) {
                // create day + meal, no options
            },
            wantErr: "برنامه ناقص است",
        },
        {
            name: "option with no items",
            setup: func(ctx context.Context, planID string) {
                // create day + meal + option, no items
            },
            wantErr: "برنامه ناقص است",
        },
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // create fresh draft plan for each case
            planID := createDraftPlan(ctx, t)
            tc.setup(ctx, planID)
            err := svc.ActivatePlan(ctx, planID, testNutritionistID)
            require.Error(t, err)
            assert.Contains(t, err.Error(), tc.wantErr)
        })
    }
}
```

**Automated command:** `cd backend && go test ./internal/service/... -run TestActivationValidation -v`

---

## Dimension 4: Row-Level Authorization

**Requirement:** AUTH-11, D-10  
**Risk if wrong:** Nutritionist B can read/modify Nutritionist A's client plans — IDOR vulnerability.

**Test file:** `backend/internal/service/diet_plan_service_test.go`

### Test Strategy

```go
// TestRowLevelAuth_CrossNutritionistAccess
func TestNutritionistCannotAccessOtherClientPlan(t *testing.T) {
    // nutritionistA created planX for clientA
    // nutritionistB attempts to read planX
    _, err := svc.GetPlanAggregate(ctx, planX.ID, nutritionistB.ID)
    require.Error(t, err)
    // Treated as not found (don't reveal existence)
    assert.ErrorIs(t, err, ErrDietPlanNotFound)
}

// TestRowLevelAuth_ClientSelfAccess
func TestClientCannotAccessOtherClientPlan(t *testing.T) {
    // client A has planX active
    // client B attempts to get their "active plan" — should see own plan or 404
    plan, err := svc.GetActivePlanForClient(ctx, clientB.ID)
    // Either nil (clientB has no active plan) or clientB's own plan — never planX
    if err == nil {
        assert.Equal(t, clientB.ID, plan.ClientID)
    }
}
```

### Manual Test

```bash
# Get nutritionist B's token
NUTRITIONIST_B_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -d '{"email":"b@test.com","password":"test"}' | jq -r '.access_token')

# Attempt to access nutritionist A's plan
curl -H "Authorization: Bearer $NUTRITIONIST_B_TOKEN" \
  http://localhost:8080/api/diet-plans/$NUTRITIONIST_A_PLAN_ID
# Expected: 404 Not Found
```

**Automated command:** `cd backend && go test ./internal/service/... -run TestRowLevelAuth -v`

---

## Dimension 5: Nutritional Computation Accuracy (Client-Side)

**Requirement:** DIET-08, D-14, D-16  
**Risk if wrong:** Nutritional totals displayed are mathematically wrong — undermines the core feature.

**Test file:** `frontend/tests/useNutritionCalc.test.ts`

### Formula

```
itemCalories = food.calories × (item.quantity / food.measurement_amount)
optionCalories = Σ itemCalories across all items in option
mealCalories = Σ optionCalories (for display — all options shown, not one)
dayCalories = Σ mealCalories
```

### Test Strategy

```typescript
// frontend/tests/useNutritionCalc.test.ts
import { describe, it, expect } from 'vitest'
import { useNutritionComputed } from '../app/composables/useNutritionComputed'

describe('useNutritionComputed', () => {
  // helpers to build minimal mock objects
  const mockItem = (overrides = {}) => ({
    id: '1', quantity: 100, measurement_unit: 'gram',
    food: {
      calories: 200, protein_g: 10, carbs_g: 25, fat_g: 5, fiber_g: 2,
      measurement_amount: 100, measurement_unit: 'gram',
    },
    ...overrides,
  })

  it('computes item nutrition proportionally to quantity', () => {
    const { itemNutrition } = useNutritionComputed()
    const n = itemNutrition(mockItem({ quantity: 50 }))
    expect(n.calories).toBeCloseTo(100)   // 200 * 50/100
    expect(n.protein_g).toBeCloseTo(5)    // 10 * 50/100
    expect(n.carbs_g).toBeCloseTo(12.5)   // 25 * 50/100
    expect(n.fat_g).toBeCloseTo(2.5)
    expect(n.fiber_g).toBeCloseTo(1)
  })

  it('handles measurement_amount = 1 (piece-based foods)', () => {
    const { itemNutrition } = useNutritionComputed()
    const item = mockItem({
      quantity: 3,
      food: { calories: 80, protein_g: 2, carbs_g: 10, fat_g: 1, fiber_g: 0.5,
              measurement_amount: 1, measurement_unit: 'piece' },
    })
    const n = itemNutrition(item)
    expect(n.calories).toBeCloseTo(240)   // 80 * 3/1
    expect(n.protein_g).toBeCloseTo(6)
  })

  it('sums option items correctly', () => {
    const { optionTotals } = useNutritionComputed()
    const option = {
      id: 'opt1', option_number: 1, label: null,
      items: [mockItem({ quantity: 100 }), mockItem({ quantity: 50 })],
    }
    const n = optionTotals(option)
    // item1: 200 cal, item2: 100 cal → 300 total
    expect(n.calories).toBeCloseTo(300)
    expect(n.protein_g).toBeCloseTo(15)
  })

  it('returns zero totals for empty option', () => {
    const { optionTotals } = useNutritionComputed()
    const n = optionTotals({ id: 'x', option_number: 1, label: null, items: [] })
    expect(n.calories).toBe(0)
    expect(n.protein_g).toBe(0)
  })
})
```

**Automated command:** `cd frontend && npx vitest run tests/useNutritionCalc.test.ts`

---

## Wave 0 Gaps (Test Stubs to Create — Plan 03-09)

The following files must exist BEFORE any implementation plan runs. Plan 03-09 creates them as compilable stubs:

| File | Purpose | Status |
|------|---------|--------|
| `backend/internal/service/diet_plan_service_test.go` | Dims 1, 3, 4 — service-layer tests | ⬜ Wave 0 creates |
| `backend/internal/repository/diet_plan_repo_test.go` | Dim 2 — batch load integration test | ⬜ Wave 0 creates |
| `frontend/tests/useNutritionCalc.test.ts` | Dim 5 — nutrition formula unit tests | ⬜ Wave 0 creates |
| `frontend/vitest.config.ts` | Vitest configuration (if not present) | ⬜ Wave 0 creates |

**Stub requirement:** Files must compile (`go build ./...` + `npx tsc --noEmit`) with empty test bodies. Actual assertions are added when Plan 03-03 (service) and Plan 03-02 (repo) implement the production code.

---

## Vitest Configuration

If `frontend/vitest.config.ts` does not exist, Plan 03-09 creates:

```typescript
// frontend/vitest.config.ts
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'happy-dom',
    globals: true,
  },
  resolve: {
    alias: {
      '~': resolve(__dirname, 'app'),
    },
  },
})
```

Install command: `cd frontend && npm install -D vitest @vue/test-utils happy-dom @vitejs/plugin-vue`

---

*Validation Architecture: Phase 03-diet-plan-engine*  
*Source: 03-RESEARCH.md §Validation Architecture*  
*Created: 2026-04-19*
