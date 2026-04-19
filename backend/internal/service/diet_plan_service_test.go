//go:build integration

package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestActivatePlanArchivesPrevious — Dim 1: one-active-plan constraint
// TODO: implemented by Plan 03-03 once DietPlanService exists
func TestActivatePlanArchivesPrevious(t *testing.T) {
	t.Skip("stub — Plan 03-03 implements")
	_ = context.Background()
	_ = assert.New(t)
}

// TestOneActivePlanConstraintAtDBLevel — Dim 1: DB partial unique index
func TestOneActivePlanConstraintAtDBLevel(t *testing.T) {
	t.Skip("stub — Plan 03-03 implements")
}

// TestActivationValidation — Dim 3: activation blocked for incomplete plans
func TestActivationValidation(t *testing.T) {
	t.Skip("stub — Plan 03-03 implements")
}

// TestNutritionistCannotAccessOtherClientPlan — Dim 4: row-level auth
func TestNutritionistCannotAccessOtherClientPlan(t *testing.T) {
	t.Skip("stub — Plan 03-03 implements")
}

// TestClientCannotAccessOtherClientPlan — Dim 4: client isolation
func TestClientCannotAccessOtherClientPlan(t *testing.T) {
	t.Skip("stub — Plan 03-03 implements")
}
