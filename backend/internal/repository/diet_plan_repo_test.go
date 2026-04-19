//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPlanAggregateLoadTime — Dim 2: ≤500ms SLA for 7×5×3×4 plan
// TODO: implemented by Plan 03-02 once DietPlanRepository exists
func TestPlanAggregateLoadTime(t *testing.T) {
	t.Skip("stub — Plan 03-02 implements")
	_ = context.Background()
	_ = time.Millisecond
	_ = require.New(t)
	_ = assert.New(t)
}
