package postgresoperations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
)

func TestVisitEPPUsageRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitEPPUsageRepo(pool)

	visitID := uuid.MustParse("80000000-0000-0000-0000-000000000001")
	eppID := uuid.MustParse("52000000-0000-0000-0000-000000000001")

	item := &operations.VisitEPPUsage{
		ID:        uuid.New(),
		VisitID:   visitID,
		EPPID:     eppID,
		Quantity:  2,
		Status:    shared.EPPUsageStatusUsed,
		Notes:     strPtr("Test notes"),
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.Quantity, got.Quantity)
	require.Equal(t, item.Status, got.Status)

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

func TestVisitEPPUsageRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitEPPUsageRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestVisitMaterialUsageRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitMaterialUsageRepo(pool)

	visitID := uuid.MustParse("80000000-0000-0000-0000-000000000001")
	matID := uuid.MustParse("50000000-0000-0000-0000-000000000001")

	item := &operations.VisitMaterialUsage{
		ID:         uuid.New(),
		VisitID:    visitID,
		MaterialID: matID,
		Quantity:   5.5,
		Notes:      strPtr("Test notes"),
		CreatedAt:  time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.Quantity, got.Quantity)

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

func TestVisitMaterialUsageRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitMaterialUsageRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestVisitMeasurementRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitMeasurementRepo(pool)

	visitID := uuid.MustParse("80000000-0000-0000-0000-000000000001")
	typeID := uuid.MustParse("35fd50f2-3af5-4a1f-b99c-1420f31395bc")

	item := &operations.VisitMeasurement{
		ID:              uuid.New(),
		VisitID:         visitID,
		Type:            typeID,
		Value:           floatPtr(10.5),
		Unit:            strPtr("meters"),
		Result:          shared.MeasurementResultApproved,
		MeasuringDevice: strPtr("Laser"),
		Notes:           strPtr("Test notes"),
		CreatedAt:       time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.Result, got.Result)

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

	item.Result = shared.MeasurementResultRejected
	err = repo.Update(ctx, item.ID, item)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, shared.MeasurementResultRejected, updated.Result)

	err = repo.Delete(ctx, item.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, item.ID)
	require.Error(t, err)
}

func TestVisitMeasurementRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVisitMeasurementRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
