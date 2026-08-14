package postgresinventory

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/model/shared"
)

func TestMaintenanceRecordRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewMaintenanceRecordRepo(pool)

	vehicleID := uuid.MustParse("53000000-0000-0000-0000-000000000001")

	record := &inventory.MaintenanceRecord{
		ID:          uuid.New(),
		Type:        shared.MaintenanceTypeVehicle,
		ReferenceID: vehicleID,
		Date:        time.Now(),
		Cost:        floatPtr(100.50),
		Supplier:    strPtr("Test Supplier"),
		Description: strPtr("Test maintenance"),
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, record)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, record.ID)

	got, err := repo.GetByID(ctx, record.ID)
	require.NoError(t, err)
	require.Equal(t, record.Type, got.Type)
	require.Equal(t, record.ReferenceID, got.ReferenceID)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, r := range result.Items {
		if r.ID == record.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, record.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, record.ID)
	require.Error(t, err)
}

func TestMaintenanceRecordRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewMaintenanceRecordRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestMaintenanceScheduleRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewMaintenanceScheduleRepo(pool)

	vehicleID := uuid.MustParse("53000000-0000-0000-0000-000000000001")

	sched := &inventory.MaintenanceSchedule{
		ID:              uuid.New(),
		Type:            shared.MaintenanceTypeVehicle,
		ReferenceID:     vehicleID,
		FrequencyDays:   intPtr(30),
		FrequencyKM:     intPtr(5000),
		LastServiceDate: timePtr(time.Now().AddDate(0, -1, 0)),
		Active:          true,
		CreatedAt:       time.Now(),
	}
	err := repo.Create(ctx, sched)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, sched.ID)

	got, err := repo.GetByID(ctx, sched.ID)
	require.NoError(t, err)
	require.Equal(t, sched.Type, got.Type)
	require.Equal(t, sched.ReferenceID, got.ReferenceID)
	require.Equal(t, sched.Active, got.Active)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, s := range result.Items {
		if s.ID == sched.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	sched.Active = false
	err = repo.Update(ctx, sched.ID, sched)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, sched.ID)
	require.NoError(t, err)
	require.False(t, updated.Active)

	err = repo.Delete(ctx, sched.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, sched.ID)
	require.Error(t, err)
}

func TestMaintenanceScheduleRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewMaintenanceScheduleRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
