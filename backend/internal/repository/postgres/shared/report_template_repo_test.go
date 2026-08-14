package postgresshared

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/shared"
)

func TestReportTemplateRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewReportTemplateRepo(pool)

	tmpl := &shared.ReportTemplate{
		ID:         uuid.New(),
		Name:       "Test Template " + uuid.New().String(),
		FieldsJSON: json.RawMessage(`{"field": "value"}`),
		CreatedAt:  time.Now(),
	}
	err := repo.Create(ctx, tmpl)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, tmpl.ID)

	got, err := repo.GetByID(ctx, tmpl.ID)
	require.NoError(t, err)
	require.Equal(t, tmpl.Name, got.Name)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, t := range result.Items {
		if t.ID == tmpl.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	tmpl.Name = "Updated " + tmpl.Name
	err = repo.Update(ctx, tmpl.ID, tmpl)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, tmpl.ID)
	require.NoError(t, err)
	require.Equal(t, tmpl.Name, updated.Name)

	err = repo.Delete(ctx, tmpl.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, tmpl.ID)
	require.Error(t, err)
}

func TestReportTemplateRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewReportTemplateRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
