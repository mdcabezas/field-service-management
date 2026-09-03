package postgresoperations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
)

func TestVisitPhotoRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitPhotoRepo(pool)

	visitID := uuid.MustParse("80000000-0000-0000-0000-000000000001")

	item := &operations.VisitPhoto{
		ID:        uuid.New(),
		VisitID:   visitID,
		URL:       strPtr("http://example.com/photo.jpg"),
		Geom:      orb.Point{-70.65, -33.45},
		Timestamp: time.Now(),
		Stage:     photoStagePtr(shared.PhotoStageDiagnosis),
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.URL, got.URL)
	require.Equal(t, item.Geom, got.Geom)
	require.Equal(t, item.Stage, got.Stage)

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

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestVisitPhotoRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitPhotoRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestVisitRentalRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitRentalRepo(pool)

	visitID := uuid.MustParse("80000000-0000-0000-0000-000000000001")
	rentalID := uuid.MustParse("54000000-0000-0000-0000-000000000001")

	item := &operations.VisitRental{
		ID:        uuid.New(),
		VisitID:   visitID,
		RentalID:  rentalID,
		StartTime: time.Now(),
		EndTime:   timePtr(time.Now().Add(2 * time.Hour)),
		Hours:     floatPtr(2.0),
		TotalCost: 150.0,
		Reason:    strPtr("Equipment rental"),
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.RentalID, got.RentalID)
	require.Equal(t, item.TotalCost, got.TotalCost)

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

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestVisitRentalRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitRentalRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestVisitToolUsageRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitToolUsageRepo(pool)

	visitID := uuid.MustParse("80000000-0000-0000-0000-000000000001")
	toolID := uuid.MustParse("51000000-0000-0000-0000-000000000001")

	item := &operations.VisitToolUsage{
		ID:        uuid.New(),
		VisitID:   visitID,
		ToolID:    toolID,
		Notes:     strPtr("Tool usage"),
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.ToolID, got.ToolID)

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

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestVisitToolUsageRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitToolUsageRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
