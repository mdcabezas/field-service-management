package postgresoperations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/operations"
)

func TestVisitChecklistEPPRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitChecklistEPPRepo(pool)

	visitID := uuid.MustParse("80000000-0000-0000-0000-000000000001")
	eppID := uuid.MustParse("52000000-0000-0000-0000-000000000001")

	item := &operations.VisitChecklistEPP{
		ID:              uuid.New(),
		VisitID:         visitID,
		EPPID:           eppID,
		PlannedQuantity: 2,
		Confirmed:       true,
		Notes:           strPtr("Test notes"),
		CreatedAt:       time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.PlannedQuantity, got.PlannedQuantity)
	require.Equal(t, item.Confirmed, got.Confirmed)

	result, err := repo.ListByVisit(ctx, visitID)
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

	item.Confirmed = false
	err = repo.Update(ctx, item.ID, item)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.False(t, updated.Confirmed)

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestVisitChecklistEPPRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitChecklistEPPRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestVisitChecklistMaterialRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitChecklistMaterialRepo(pool)

	visitID := uuid.MustParse("80000000-0000-0000-0000-000000000001")
	matID := uuid.MustParse("50000000-0000-0000-0000-000000000001")

	item := &operations.VisitChecklistMaterial{
		ID:              uuid.New(),
		VisitID:         visitID,
		MaterialID:      matID,
		PlannedQuantity: 5,
		Confirmed:       true,
		Notes:           strPtr("Test notes"),
		CreatedAt:       time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.PlannedQuantity, got.PlannedQuantity)
	require.Equal(t, item.Confirmed, got.Confirmed)

	result, err := repo.ListByVisit(ctx, visitID)
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

	item.Confirmed = false
	err = repo.Update(ctx, item.ID, item)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.False(t, updated.Confirmed)

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestVisitChecklistMaterialRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitChecklistMaterialRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestVisitChecklistToolRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitChecklistToolRepo(pool)

	visitID := uuid.MustParse("80000000-0000-0000-0000-000000000001")
	toolID := uuid.MustParse("51000000-0000-0000-0000-000000000001")

	item := &operations.VisitChecklistTool{
		ID:              uuid.New(),
		VisitID:         visitID,
		ToolID:          toolID,
		PlannedQuantity: 3,
		Confirmed:       true,
		Notes:           strPtr("Test notes"),
		CreatedAt:       time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.PlannedQuantity, got.PlannedQuantity)
	require.Equal(t, item.Confirmed, got.Confirmed)

	result, err := repo.ListByVisit(ctx, visitID)
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

	item.Confirmed = false
	err = repo.Update(ctx, item.ID, item)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.False(t, updated.Confirmed)

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestVisitChecklistToolRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitChecklistToolRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
