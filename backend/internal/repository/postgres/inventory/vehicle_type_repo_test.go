package postgresinventory

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestVehicleTypeRepo_List(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVehicleTypeRepo(pool)

	result, err := repo.List(ctx)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Should have the 5 core vehicle types
	require.GreaterOrEqual(t, result.Total, 5)
	require.Len(t, result.Items, result.Total)

	codes := make(map[string]bool)
	for _, vt := range result.Items {
		require.NotEmpty(t, vt.ID)
		require.NotEmpty(t, vt.Code)
		require.NotEmpty(t, vt.Name)
		codes[vt.Code] = true
	}

	// Core vehicle types should exist
	expectedCodes := []string{"truck", "crane", "van", "crane_truck", "other"}
	for _, c := range expectedCodes {
		require.True(t, codes[c], "missing expected vehicle type code: %s", c)
	}
}

func TestVehicleTypeRepo_GetByID(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVehicleTypeRepo(pool)

	// First list to get a valid ID
	list, err := repo.List(ctx)
	require.NoError(t, err)
	require.Greater(t, len(list.Items), 0)

	first := list.Items[0]
	vt, err := repo.GetByID(ctx, first.ID)
	require.NoError(t, err)
	require.NotNil(t, vt)
	require.Equal(t, first.ID, vt.ID)
	require.Equal(t, first.Code, vt.Code)
	require.Equal(t, first.Name, vt.Name)
	require.Equal(t, first.Active, vt.Active)
}

func TestVehicleTypeRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVehicleTypeRepo(pool)

	nonExistentID := uuid.New()
	vt, err := repo.GetByID(ctx, nonExistentID)
	require.NoError(t, err)
	require.Nil(t, vt)
}
