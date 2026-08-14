package postgresoperations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/operations"
)

func TestReportEntryRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewReportEntryRepo(pool)

	reportID := uuid.MustParse("8a000000-0000-0000-0000-000000000001")

	entry := &operations.ReportEntry{
		ID:        uuid.New(),
		ReportID:  reportID,
		DataJSON:  []byte(`{"key": "value"}`),
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, entry)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, entry.ID)

	got, err := repo.GetByID(ctx, entry.ID)
	require.NoError(t, err)
	require.Equal(t, entry.ReportID, got.ReportID)

	result, err := repo.ListByReport(ctx, reportID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, r := range result {
		if r.ID == entry.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, entry.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, entry.ID)
	require.Error(t, err)
}

func TestReportEntryRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewReportEntryRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestReportImageRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewReportImageRepo(pool)

	reportID := uuid.MustParse("8a000000-0000-0000-0000-000000000001")

	img := &operations.ReportImage{
		ID:        uuid.New(),
		ReportID:  reportID,
		URL:       "http://example.com/image.jpg",
		Pages:     2,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, img)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, img.ID)

	got, err := repo.GetByID(ctx, img.ID)
	require.NoError(t, err)
	require.Equal(t, img.URL, got.URL)
	require.Equal(t, img.Pages, got.Pages)

	result, err := repo.ListByReport(ctx, reportID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, r := range result {
		if r.ID == img.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, img.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, img.ID)
	require.Error(t, err)
}

func TestReportImageRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewReportImageRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
