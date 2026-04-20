//go:build integration

package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ranjbar-dev/nutritrack/backend/internal/repository"
)

// TestPlanAggregateLoadTime — Dim 2: ≤500ms SLA for 7×5×3×4 plan
func TestPlanAggregateLoadTime(t *testing.T) {
	databaseURL := os.Getenv("NUTRITRACK_TEST_DATABASE_URL")
	planIDRaw := os.Getenv("NUTRITRACK_PERF_PLAN_ID")
	if databaseURL == "" || planIDRaw == "" {
		t.Skip("set NUTRITRACK_TEST_DATABASE_URL and NUTRITRACK_PERF_PLAN_ID to run the performance harness")
	}

	planID, err := uuid.Parse(planIDRaw)
	require.NoError(t, err)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := repository.NewDietPlanRepository(pool)

	start := time.Now()
	plan, err := repo.GetFullPlanAggregate(ctx, planID)
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, plan)
	assert.LessOrEqual(t, elapsed, 500*time.Millisecond, "diet plan aggregate exceeded the 500ms target")
}
