package postgrescustomers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/customers"
)

func TestCustomerRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCustomerRepo(pool)

	// Create
	customer := &customers.Customer{
		ID:        uuid.New(),
		Name:      "Test Customer " + uuid.New().String(),
		Phone:     strPtr("+1234567890"),
		Email:     strPtr("test_" + uuid.New().String() + "@example.com"),
		TaxID:     strPtr("TAX123456"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, customer)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, customer.ID)

	// GetByID
	got, err := repo.GetByID(ctx, customer.ID)
	require.NoError(t, err)
	require.Equal(t, customer.Name, got.Name)
	require.Equal(t, customer.Phone, got.Phone)
	require.Equal(t, customer.Email, got.Email)
	require.Equal(t, customer.TaxID, got.TaxID)

	// List
	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, c := range result.Items {
		if c.ID == customer.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	// Update
	newName := "Updated Customer " + uuid.New().String()
	customer.Name = newName
	customer.UpdatedAt = time.Now()
	err = repo.Update(ctx, customer.ID, customer)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, customer.ID)
	require.NoError(t, err)
	require.Equal(t, newName, updated.Name)

	// Delete
	err = repo.Delete(ctx, customer.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, customer.ID)
	require.Error(t, err)
}

func TestCustomerRepo_ListPagination(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCustomerRepo(pool)

	var createdIDs []uuid.UUID
	for i := 0; i < 5; i++ {
		c := &customers.Customer{
			ID:        uuid.New(),
			Name:      "Page Customer " + uuid.New().String(),
			Phone:     strPtr("+100000000" + string(rune('0'+i))),
			Email:     strPtr("pagecust" + uuid.New().String() + "@example.com"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		require.NoError(t, repo.Create(ctx, c))
		createdIDs = append(createdIDs, c.ID)
	}

	page1, err := repo.List(ctx, 2, 0)
	require.NoError(t, err)
	require.Len(t, page1.Items, 2)

	page2, err := repo.List(ctx, 2, 2)
	require.NoError(t, err)
	require.Len(t, page2.Items, 2)

	for _, c1 := range page1.Items {
		for _, c2 := range page2.Items {
			require.NotEqual(t, c1.ID, c2.ID)
		}
	}

	for _, id := range createdIDs {
		repo.Delete(ctx, id)
	}
}

func TestCustomerRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCustomerRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
