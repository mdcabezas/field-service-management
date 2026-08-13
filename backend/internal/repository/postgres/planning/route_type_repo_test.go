package postgresplanning

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://fsm_admin:change_me_in_prod@localhost:5432/fsm_test?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })
	return pool
}

func TestRouteTypeRepo_List(t *testing.T) {
	ctx := context.Background()
	pool := testPool(t)
	repo := NewRouteTypeRepo(pool)

	result, err := repo.List(ctx)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Should have the 8 core route types
	require.GreaterOrEqual(t, result.Total, 8)
	require.Len(t, result.Items, result.Total)

	codes := make(map[string]bool)
	for _, rt := range result.Items {
		require.NotEmpty(t, rt.ID)
		require.NotEmpty(t, rt.Code)
		require.NotEmpty(t, rt.Name)
		require.True(t, rt.Active)
		codes[rt.Code] = true
	}

	// Core route types should exist
	expectedCodes := []string{
		"maintenance", "installation", "repair", "inspection",
		"delivery", "collection", "emergency", "other",
	}
	for _, c := range expectedCodes {
		require.True(t, codes[c], "missing expected route type code: %s", c)
	}
}
