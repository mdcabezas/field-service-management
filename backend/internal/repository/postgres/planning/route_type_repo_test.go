package postgresplanning

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRouteTypeRepo_List(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
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
