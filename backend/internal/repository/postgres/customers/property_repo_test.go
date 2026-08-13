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

func TestPropertyRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPropertyRepo(pool)

	// Need a customer
	customerRepo := NewCustomerRepo(pool)
	customer := &customers.Customer{
		ID:        uuid.New(),
		Name:      "Customer for Property " + uuid.New().String(),
		Phone:     strPtr("+1234567890"),
		Email:     strPtr("cust_prop_" + uuid.New().String() + "@example.com"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, customerRepo.Create(ctx, customer))

	// Need an address
	addrRepo := postgresgeocoding.NewGeocodingAddressRepo(pool)
	address := &geocoding.Address{
		ID:         uuid.New(),
		Street:     "Property Street",
		Number:     strPtr("456"),
		City:       strPtr("Property City"),
		Region:     strPtr("PC"),
		PostalCode: strPtr("67890"),
		CreatedAt:  time.Now(),
	}
	require.NoError(t, addrRepo.Create(ctx, address))

	// Need a customer_address
	caRepo := NewCustomerAddressRepo(pool)
	customerAddress := &customers.CustomerAddress{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		AddressID:  address.ID,
		Name:       strPtr("Property Location"),
		Type:       shared.CustomerAddressTypeResidential,
		CreatedAt:  time.Now(),
	}
	require.NoError(t, caRepo.Create(ctx, customerAddress))

	// Create Property
	prop := &customers.Property{
		ID:                uuid.New(),
		CustomerAddressID: customerAddress.ID,
		Name:              strPtr("Main Building"),
		Type:              nil,
		Notes:             strPtr("Main office building"),
		CreatedAt:         time.Now(),
	}
	err := repo.Create(ctx, prop)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, prop.ID)

	// GetByID
	got, err := repo.GetByID(ctx, prop.ID)
	require.NoError(t, err)
	require.Equal(t, prop.CustomerAddressID, got.CustomerAddressID)
	require.Equal(t, prop.Name, got.Name)
	require.Equal(t, prop.Notes, got.Notes)

	// ListByCustomerAddress
	properties, err := repo.ListByCustomerAddress(ctx, customerAddress.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(properties), 1)
	found := false
	for _, p := range properties {
		if p.ID == prop.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	// Update
	newName := "Updated Building"
	prop.Name = strPtr(newName)
	prop.Notes = strPtr("Updated notes")
	err = repo.Update(ctx, prop.ID, prop)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, prop.ID)
	require.NoError(t, err)
	require.Equal(t, newName, *updated.Name)
	require.Equal(t, "Updated notes", *updated.Notes)

	// Delete
	err = repo.Delete(ctx, prop.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, prop.ID)
	require.Error(t, err)

	// Cleanup
	caRepo.Delete(ctx, customerAddress.ID)
	addrRepo.Delete(ctx, address.ID)
	customerRepo.Delete(ctx, customer.ID)
}

func TestPropertyRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPropertyRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestPropertyRepo_ListByCustomerAddress_Empty(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPropertyRepo(pool)

	// Create customer and address but no property
	customerRepo := NewCustomerRepo(pool)
	customer := &customers.Customer{
		ID:        uuid.New(),
		Name:      "Empty Prop Customer " + uuid.New().String(),
		Phone:     strPtr("+1000000000"),
		Email:     strPtr("emptyprop_" + uuid.New().String() + "@example.com"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, customerRepo.Create(ctx, customer))

	addrRepo := postgresgeocoding.NewGeocodingAddressRepo(pool)
	address := &geocoding.Address{
		ID:         uuid.New(),
		Street:     "Empty Prop Street",
		Number:     strPtr("999"),
		City:       strPtr("Empty City"),
		Region:     strPtr("EC"),
		PostalCode: strPtr("00000"),
		CreatedAt:  time.Now(),
	}
	require.NoError(t, addrRepo.Create(ctx, address))

	caRepo := NewCustomerAddressRepo(pool)
	customerAddress := &customers.CustomerAddress{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		AddressID:  address.ID,
		Name:       strPtr("No Properties"),
		Type:       shared.CustomerAddressTypeResidential,
		CreatedAt:  time.Now(),
	}
	require.NoError(t, caRepo.Create(ctx, customerAddress))

	properties, err := repo.ListByCustomerAddress(ctx, customerAddress.ID)
	require.NoError(t, err)
	require.Len(t, properties, 0)

	// Cleanup
	caRepo.Delete(ctx, customerAddress.ID)
	addrRepo.Delete(ctx, address.ID)
	customerRepo.Delete(ctx, customer.ID)
}
