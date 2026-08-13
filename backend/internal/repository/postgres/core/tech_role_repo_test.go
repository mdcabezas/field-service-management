package postgrescore

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/core"
)

func TestCoreTechRoleRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCoreTechRoleRepo(pool)

	// Create
	role := &core.TechRole{
		Code:   "crud_" + uuid.New().String(),
		Name:   "CRUD Technician",
		Active: true,
	}
	err := repo.Create(ctx, role)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, role.ID)

	// GetByID
	got, err := repo.GetByID(ctx, role.ID)
	require.NoError(t, err)
	require.Equal(t, role.Code, got.Code)
	require.Equal(t, role.Name, got.Name)
	require.Equal(t, role.Active, got.Active)

	// List
	result, err := repo.List(ctx, 10, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, r := range result.Items {
		if r.ID == role.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	// Update
	role.Name = "Senior Tech Updated"
	role.Active = false
	err = repo.Update(ctx, role.ID, role)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, role.ID)
	require.NoError(t, err)
	require.Equal(t, "Senior Tech Updated", updated.Name)
	require.False(t, updated.Active)

	// Delete
	err = repo.Delete(ctx, role.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, role.ID)
	require.Error(t, err)
}

func TestCoreTechRoleRepo_ListPagination(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCoreTechRoleRepo(pool)

	// Create additional roles
	var createdIDs []uuid.UUID
	for i := 0; i < 3; i++ {
		r := &core.TechRole{
			Code:   "page_" + uuid.New().String(),
			Name:   "Page Role " + uuid.New().String(),
			Active: true,
		}
		require.NoError(t, repo.Create(ctx, r))
		createdIDs = append(createdIDs, r.ID)
	}

	page1, err := repo.List(ctx, 2, 0)
	require.NoError(t, err)
	require.Len(t, page1.Items, 2)

	page2, err := repo.List(ctx, 2, 2)
	require.NoError(t, err)
	require.Len(t, page2.Items, 2)

	for _, r1 := range page1.Items {
		for _, r2 := range page2.Items {
			require.NotEqual(t, r1.ID, r2.ID)
		}
	}

	// Cleanup
	for _, id := range createdIDs {
		repo.Delete(ctx, id)
	}
}

func TestCoreTechRoleRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCoreTechRoleRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestCoreTechRoleRepo_DuplicateCode(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCoreTechRoleRepo(pool)

	code := "dup_" + uuid.New().String()
	r1 := &core.TechRole{Code: code, Name: "Role 1", Active: true}
	require.NoError(t, repo.Create(ctx, r1))

	r2 := &core.TechRole{Code: code, Name: "Role 2", Active: true}
	err := repo.Create(ctx, r2)
	require.Error(t, err) // unique constraint on code
}

func TestCoreTechRoleRepo_ActiveFilter(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCoreTechRoleRepo(pool)

	// Create inactive role
	inactive := &core.TechRole{Code: "inactive_" + uuid.New().String(), Name: "Inactive", Active: false}
	require.NoError(t, repo.Create(ctx, inactive))

	// List all (should include inactive)
	all, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)

	// Create active role
	active := &core.TechRole{Code: "active_" + uuid.New().String(), Name: "Active", Active: true}
	require.NoError(t, repo.Create(ctx, active))

	// List again
	all2, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.Greater(t, all2.Total, all.Total)
}
