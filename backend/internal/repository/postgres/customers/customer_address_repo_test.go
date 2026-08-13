package postgrescustomers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/customers"
	"localis-backend/internal/model/geocoding"
	"localis-backend/internal/model/shared"
	postgresgeocoding "localis-backend/internal/repository/postgres/geocoding"
)

func TestCustomerAddressRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCustomerAddressRepo(pool)

	// Need a customer first
	customerRepo := NewCustomerRepo(pool)
	customer := &customers.Customer{
		ID:        uuid.New(),
		Name:      "Customer for Address " + uuid.New().String(),
		Phone:     strPtr("+1234567890"),
		Email:     strPtr("cust_" + uuid.New().String() + "@example.com"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, customerRepo.Create(ctx, customer))

	// Need an address from geocoding
	addrRepo := postgresgeocoding.NewGeocodingAddressRepo(pool)
	address := &geocoding.Address{
		ID:         uuid.New(),
		Street:     "Test Street",
		Number:     strPtr("123"),
		City:       strPtr("Test City"),
		Region:     strPtr("TS"),
		PostalCode: strPtr("12345"),
		CreatedAt:  time.Now(),
	}
	require.NoError(t, addrRepo.Create(ctx, address))

	// Create CustomerAddress
	ca := &customers.CustomerAddress{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		AddressID:  address.ID,
		Name:       strPtr("Home"),
		Type:       shared.CustomerAddressTypeResidential,
		CreatedAt:  time.Now(),
	}
	err := repo.Create(ctx, ca)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, ca.ID)

	// GetByID
	got, err := repo.GetByID(ctx, ca.ID)
	require.NoError(t, err)
	require.Equal(t, ca.CustomerID, got.CustomerID)
	require.Equal(t, ca.AddressID, got.AddressID)
	require.Equal(t, ca.Name, got.Name)
	require.Equal(t, ca.Type, got.Type)

	// ListByCustomer
	addresses, err := repo.ListByCustomer(ctx, customer.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(addresses), 1)
	found := false
	for _, a := range addresses {
		if a.ID == ca.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	// Update
	newName := "Office Updated"
	ca.Name = strPtr(newName)
	ca.Type = shared.CustomerAddressTypeCommercial
	err = repo.Update(ctx, ca.ID, ca)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, ca.ID)
	require.NoError(t, err)
	require.Equal(t, newName, *updated.Name)
	require.Equal(t, shared.CustomerAddressTypeCommercial, updated.Type)

	// Delete
	err = repo.Delete(ctx, ca.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, ca.ID)
	require.Error(t, err)

	// Cleanup
	addrRepo.Delete(ctx, address.ID)
	customerRepo.Delete(ctx, customer.ID)
}

func TestCustomerAddressRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCustomerAddressRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestCustomerAddressRepo_ListByCustomer_Empty(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCustomerAddressRepo(pool)

	// Create a customer without addresses
	customerRepo := NewCustomerRepo(pool)
	customer := &customers.Customer{
		ID:        uuid.New(),
		Name:      "Empty Customer " + uuid.New().String(),
		Phone:     strPtr("+1000000000"),
		Email:     strPtr("empty_" + uuid.New().String() + "@example.com"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, customerRepo.Create(ctx, customer))

	addresses, err := repo.ListByCustomer(ctx, customer.ID)
	require.NoError(t, err)
	require.Len(t, addresses, 0)

	customerRepo.Delete(ctx, customer.ID)
}
