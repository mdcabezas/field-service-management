package postgresplanning

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/planning"
)

func TestDailyLoadEPPRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyLoadEPPRepo(pool)

	planID := uuid.MustParse("70000000-0000-0000-0000-000000000001")
	eppID := uuid.MustParse("52000000-0000-0000-0000-000000000001")

	item := &planning.DailyLoadEPP{
		ID:               uuid.New(),
		DailyPlanID:      planID,
		EPPID:            eppID,
		LoadedQuantity:   2,
		ReturnedQuantity: 1,
		CreatedAt:        time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.LoadedQuantity, got.LoadedQuantity)
	require.Equal(t, item.ReturnedQuantity, got.ReturnedQuantity)

	result, err := repo.ListByPlan(ctx, planID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, r := range result {
		if r.ID == item.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestDailyLoadEPPRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyLoadEPPRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestDailyLoadMaterialRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyLoadMaterialRepo(pool)

	planID := uuid.MustParse("70000000-0000-0000-0000-000000000001")
	matID := uuid.MustParse("50000000-0000-0000-0000-000000000001")

	item := &planning.DailyLoadMaterial{
		ID:               uuid.New(),
		DailyPlanID:      planID,
		MaterialID:       matID,
		LoadedQuantity:   10,
		ReturnedQuantity: 5,
		CreatedAt:        time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.LoadedQuantity, got.LoadedQuantity)
	require.Equal(t, item.ReturnedQuantity, got.ReturnedQuantity)

	result, err := repo.ListByPlan(ctx, planID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, r := range result {
		if r.ID == item.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestDailyLoadMaterialRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyLoadMaterialRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestDailyLoadToolRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyLoadToolRepo(pool)

	planID := uuid.MustParse("70000000-0000-0000-0000-000000000001")
	toolID := uuid.MustParse("51000000-0000-0000-0000-000000000001")

	item := &planning.DailyLoadTool{
		ID:               uuid.New(),
		DailyPlanID:      planID,
		ToolID:           toolID,
		LoadedQuantity:   5,
		ReturnedQuantity: 3,
		CreatedAt:        time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.LoadedQuantity, got.LoadedQuantity)
	require.Equal(t, item.ReturnedQuantity, got.ReturnedQuantity)

	result, err := repo.ListByPlan(ctx, planID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, r := range result {
		if r.ID == item.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestDailyLoadToolRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewDailyLoadToolRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
