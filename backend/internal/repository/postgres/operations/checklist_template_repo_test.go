package postgresoperations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/operations"
)

func TestChecklistTemplateEPPRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewChecklistTemplateEPPRepo(pool)

	tmplID := uuid.MustParse("56000000-0000-0000-0000-000000000001")
	eppID := uuid.MustParse("52000000-0000-0000-0000-000000000001")

	item := &operations.ChecklistTemplateEPP{
		ID:              uuid.New(),
		TemplateID:      tmplID,
		EPPID:           eppID,
		DefaultQuantity: 2,
		Required:        true,
		CreatedAt:       time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.DefaultQuantity, got.DefaultQuantity)
	require.Equal(t, item.Required, got.Required)

	result, err := repo.ListByTemplate(ctx, tmplID)
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

func TestChecklistTemplateEPPRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewChecklistTemplateEPPRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestChecklistTemplateMaterialRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewChecklistTemplateMaterialRepo(pool)

	tmplID := uuid.MustParse("56000000-0000-0000-0000-000000000001")
	matID := uuid.MustParse("50000000-0000-0000-0000-000000000001")

	item := &operations.ChecklistTemplateMaterial{
		ID:              uuid.New(),
		TemplateID:      tmplID,
		MaterialID:      matID,
		DefaultQuantity: 5,
		Required:        true,
		CreatedAt:       time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.DefaultQuantity, got.DefaultQuantity)
	require.Equal(t, item.Required, got.Required)

	result, err := repo.ListByTemplate(ctx, tmplID)
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

func TestChecklistTemplateMaterialRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewChecklistTemplateMaterialRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestChecklistTemplateToolRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewChecklistTemplateToolRepo(pool)

	tmplID := uuid.MustParse("56000000-0000-0000-0000-000000000001")
	toolID := uuid.MustParse("51000000-0000-0000-0000-000000000001")

	item := &operations.ChecklistTemplateTool{
		ID:              uuid.New(),
		TemplateID:      tmplID,
		ToolID:          toolID,
		DefaultQuantity: 3,
		Required:        true,
		CreatedAt:       time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.DefaultQuantity, got.DefaultQuantity)
	require.Equal(t, item.Required, got.Required)

	result, err := repo.ListByTemplate(ctx, tmplID)
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

func TestChecklistTemplateToolRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewChecklistTemplateToolRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
