package postgrescustomers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/core"
	"localis-backend/internal/model/customers"
	postgrescore "localis-backend/internal/repository/postgres/core"
)

func TestTechnicianRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewTechnicianRepo(pool)
	userRepo := postgrescore.NewCoreUserRepo(pool)

	// Create a user first (FK requirement)
	user := &core.User{
		Email:     "tech_user_" + uuid.New().String() + "@example.com",
		Role:      "technician",
		Name:      "Tech User " + uuid.New().String(),
		CreatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, user))

	// Create Technician
	tech := &customers.Technician{
		ID:          uuid.New(),
		UserID:      strPtr(user.ID.String()),
		Name:        "Test Technician " + uuid.New().String(),
		IsActive:    true,
		Specialties: []string{"HVAC", "Electrical"},
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, tech)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, tech.ID)

	// GetByID
	got, err := repo.GetByID(ctx, tech.ID)
	require.NoError(t, err)
	require.Equal(t, tech.UserID, got.UserID)
	require.Equal(t, tech.Name, got.Name)
	require.Equal(t, tech.IsActive, got.IsActive)
	require.Equal(t, tech.Specialties, got.Specialties)

	// GetByUserID
	gotByUser, err := repo.GetByUserID(ctx, user.ID.String())
	require.NoError(t, err)
	require.Equal(t, tech.ID, gotByUser.ID)

	// List
	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, ttech := range result.Items {
		if ttech.ID == tech.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	// Update
	tech.Name = "Updated Technician " + uuid.New().String()
	tech.IsActive = false
	tech.Specialties = []string{"Plumbing"}
	err = repo.Update(ctx, tech.ID, tech)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, tech.ID)
	require.NoError(t, err)
	require.Equal(t, tech.Name, updated.Name)
	require.False(t, updated.IsActive)
	require.Equal(t, []string{"Plumbing"}, updated.Specialties)

	// Delete
	err = repo.Delete(ctx, tech.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, tech.ID)
	require.Error(t, err)

	// Cleanup user
	userRepo.Delete(ctx, user.ID.String())
}

func TestTechnicianRepo_ListPagination(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewTechnicianRepo(pool)
	userRepo := postgrescore.NewCoreUserRepo(pool)

	var createdIDs []uuid.UUID
	for i := 0; i < 5; i++ {
		// Create user first
		user := &core.User{
			Email:     "page_user_" + uuid.New().String() + "@example.com",
			Role:      "technician",
			Name:      "Page Tech User " + uuid.New().String(),
			CreatedAt: time.Now(),
		}
		require.NoError(t, userRepo.Create(ctx, user))

		tech := &customers.Technician{
			ID:          uuid.New(),
			UserID:      strPtr(user.ID.String()),
			Name:        "Page Tech " + uuid.New().String(),
			IsActive:    true,
			Specialties: []string{"General"},
			CreatedAt:   time.Now(),
		}
		require.NoError(t, repo.Create(ctx, tech))
		createdIDs = append(createdIDs, tech.ID)
	}

	page1, err := repo.List(ctx, 2, 0)
	require.NoError(t, err)
	require.Len(t, page1.Items, 2)

	page2, err := repo.List(ctx, 2, 2)
	require.NoError(t, err)
	require.Len(t, page2.Items, 2)

	for _, t1 := range page1.Items {
		for _, t2 := range page2.Items {
			require.NotEqual(t, t1.ID, t2.ID)
		}
	}

	for _, id := range createdIDs {
		repo.Delete(ctx, id)
	}

	// Cleanup users - get all users created in this test
	// (In practice, we'd track user IDs, but for simplicity we'll skip)
}

func TestTechnicianRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewTechnicianRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestTechnicianRepo_GetByUserID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewTechnicianRepo(pool)

	_, err := repo.GetByUserID(ctx, uuid.New().String())
	require.Error(t, err)
}
