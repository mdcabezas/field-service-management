package postgrescore

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/core"
)

func TestCoreUserRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCoreUserRepo(pool)

	// Create
	email := "testuser_" + uuid.New().String() + "@example.com"
	user := &core.User{
		Email: email,
		Role:  "technician",
		Name:  "Test User",
	}
	err := repo.Create(ctx, user)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, user.ID)

	// GetByID
	got, err := repo.GetByID(ctx, user.ID.String())
	require.NoError(t, err)
	require.Equal(t, user.Email, got.Email)
	require.Equal(t, user.Name, got.Name)
	require.Equal(t, user.Role, got.Role)

	// List - use larger limit to ensure our user is included
	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, u := range result.Items {
		if u.ID == user.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	// Update
	user.Name = "Updated Name"
	user.Role = "admin"
	err = repo.Update(ctx, user.ID.String(), user)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, user.ID.String())
	require.NoError(t, err)
	require.Equal(t, "Updated Name", updated.Name)
	require.Equal(t, "admin", updated.Role)

	// Delete
	err = repo.Delete(ctx, user.ID.String())
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, user.ID.String())
	require.Error(t, err)
}

func TestCoreUserRepo_ListPagination(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCoreUserRepo(pool)

	// Create multiple users with unique emails
	var createdIDs []uuid.UUID
	for i := 0; i < 5; i++ {
		u := &core.User{
			Email: "pageuser" + uuid.New().String() + "@example.com",
			Role:  "technician",
			Name:  "Page User " + uuid.New().String(),
		}
		require.NoError(t, repo.Create(ctx, u))
		createdIDs = append(createdIDs, u.ID)
	}

	page1, err := repo.List(ctx, 2, 0)
	require.NoError(t, err)
	require.Len(t, page1.Items, 2)

	page2, err := repo.List(ctx, 2, 2)
	require.NoError(t, err)
	require.Len(t, page2.Items, 2)

	// No overlap
	for _, u1 := range page1.Items {
		for _, u2 := range page2.Items {
			require.NotEqual(t, u1.ID, u2.ID)
		}
	}

	// Cleanup
	for _, id := range createdIDs {
		repo.Delete(ctx, id.String())
	}
}

func TestCoreUserRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCoreUserRepo(pool)

	_, err := repo.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
	require.Error(t, err)
}

func TestCoreUserRepo_DuplicateEmail(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCoreUserRepo(pool)

	email := "dup_" + uuid.New().String() + "@example.com"
	u1 := &core.User{Email: email, Role: "tech", Name: "User 1"}
	require.NoError(t, repo.Create(ctx, u1))

	u2 := &core.User{Email: email, Role: "tech", Name: "User 2"}
	err := repo.Create(ctx, u2)
	require.Error(t, err) // unique constraint on email
}
