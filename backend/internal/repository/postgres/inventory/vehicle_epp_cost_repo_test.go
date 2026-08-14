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

func TestVehicleRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVehicleRepo(pool)

	vehicle := &inventory.Vehicle{
		ID:           uuid.New(),
		Name:         "Test Vehicle " + uuid.New().String(),
		Status:       shared.VehicleStatusAvailable,
		LicensePlate: strPtr("TEST-" + uuid.New().String()[:4]),
		CreatedAt:    time.Now(),
	}
	err := repo.Create(ctx, vehicle)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, vehicle.ID)

	got, err := repo.GetByID(ctx, vehicle.ID)
	require.NoError(t, err)
	require.Equal(t, vehicle.Name, got.Name)
	require.Equal(t, vehicle.Status, got.Status)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, v := range result.Items {
		if v.ID == vehicle.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	vehicle.Name = "Updated " + vehicle.Name
	err = repo.Update(ctx, vehicle.ID, vehicle)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, vehicle.ID)
	require.NoError(t, err)
	require.Equal(t, vehicle.Name, updated.Name)

	err = repo.Delete(ctx, vehicle.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, vehicle.ID)
	require.Error(t, err)
}

func TestVehicleRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewVehicleRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestEPPItemRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewEPPItemRepo(pool)

	item := &inventory.EPPItem{
		ID:        uuid.New(),
		Name:      "Test EPP " + uuid.New().String(),
		Type:      shared.EPPTypeHead,
		Lifecycle: shared.EPPLifecycleReusable,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, item)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, item.ID)

	got, err := repo.GetByID(ctx, item.ID)
	require.NoError(t, err)
	require.Equal(t, item.Name, got.Name)
	require.Equal(t, item.Type, got.Type)
	require.Equal(t, item.Lifecycle, got.Lifecycle)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, i := range result.Items {
		if i.ID == item.ID {
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

func TestEPPItemRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewEPPItemRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestCostRateRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCostRateRepo(pool)

	rate := &inventory.CostRate{
		ID:        uuid.New(),
		Type:      shared.CostTypeLabor,
		Value:     50.0,
		Unit:      shared.CostUnitHour,
		ValidFrom: time.Now(),
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, rate)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, rate.ID)

	got, err := repo.GetByID(ctx, rate.ID)
	require.NoError(t, err)
	require.Equal(t, rate.Type, got.Type)
	require.Equal(t, rate.Value, got.Value)
	require.Equal(t, rate.Unit, got.Unit)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, r := range result.Items {
		if r.ID == rate.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	rate.Value = 75.0
	err = repo.Update(ctx, rate.ID, rate)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, rate.ID)
	require.NoError(t, err)
	require.Equal(t, 75.0, updated.Value)

	err = repo.Delete(ctx, rate.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, rate.ID)
	require.Error(t, err)
}

func TestCostRateRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCostRateRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
