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

func TestMaterialRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewMaterialRepo(pool)

	mat := &inventory.Material{
		ID:        uuid.New(),
		Name:      "Test Material " + uuid.New().String(),
		Unit:      "m",
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, mat)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, mat.ID)

	got, err := repo.GetByID(ctx, mat.ID)
	require.NoError(t, err)
	require.Equal(t, mat.Name, got.Name)
	require.Equal(t, mat.Unit, got.Unit)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, m := range result.Items {
		if m.ID == mat.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	mat.Name = "Updated " + mat.Name
	err = repo.Update(ctx, mat.ID, mat)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, mat.ID)
	require.NoError(t, err)
	require.Equal(t, mat.Name, updated.Name)

	err = repo.Delete(ctx, mat.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, mat.ID)
	require.Error(t, err)
}

func TestMaterialRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewMaterialRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestRentalRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewRentalRepo(pool)

	rental := &inventory.Rental{
		ID:              uuid.New(),
		Type:            shared.RentalTypeVehicle,
		ItemDescription: "Test Rental " + uuid.New().String(),
		CreatedAt:       time.Now(),
	}
	err := repo.Create(ctx, rental)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, rental.ID)

	got, err := repo.GetByID(ctx, rental.ID)
	require.NoError(t, err)
	require.Equal(t, rental.Type, got.Type)
	require.Equal(t, rental.ItemDescription, got.ItemDescription)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, r := range result.Items {
		if r.ID == rental.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, rental.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, rental.ID)
	require.Error(t, err)
}

func TestRentalRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewRentalRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestToolRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewToolRepo(pool)

	tool := &inventory.Tool{
		ID:        uuid.New(),
		Name:      "Test Tool " + uuid.New().String(),
		Status:    shared.ToolStatusAvailable,
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, tool)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, tool.ID)

	got, err := repo.GetByID(ctx, tool.ID)
	require.NoError(t, err)
	require.Equal(t, tool.Name, got.Name)
	require.Equal(t, tool.Status, got.Status)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, t := range result.Items {
		if t.ID == tool.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	tool.Name = "Updated " + tool.Name
	err = repo.Update(ctx, tool.ID, tool)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, tool.ID)
	require.NoError(t, err)
	require.Equal(t, tool.Name, updated.Name)

	err = repo.Delete(ctx, tool.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, tool.ID)
	require.Error(t, err)
}

func TestToolRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewToolRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
