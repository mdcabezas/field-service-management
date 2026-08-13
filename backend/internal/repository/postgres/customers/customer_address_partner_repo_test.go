package postgrescustomers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/customers"
	"localis-backend/internal/model/geocoding"
	"localis-backend/internal/model/partners"
	"localis-backend/internal/model/shared"
	postgresgeocoding "localis-backend/internal/repository/postgres/geocoding"
	postgrespartners "localis-backend/internal/repository/postgres/partners"
)

func TestCustomerAddressPartnerRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCustomerAddressPartnerRepo(pool)

	// Need a customer
	customerRepo := NewCustomerRepo(pool)
	customer := &customers.Customer{
		ID:        uuid.New(),
		Name:      "Customer for Partner " + uuid.New().String(),
		Phone:     strPtr("+1234567890"),
		Email:     strPtr("cust_partner_" + uuid.New().String() + "@example.com"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, customerRepo.Create(ctx, customer))

	// Need an address
	addrRepo := postgresgeocoding.NewGeocodingAddressRepo(pool)
	address := &geocoding.Address{
		ID:         uuid.New(),
		Street:     "Partner Street",
		Number:     strPtr("789"),
		City:       strPtr("Partner City"),
		Region:     strPtr("PC"),
		PostalCode: strPtr("99999"),
		CreatedAt:  time.Now(),
	}
	require.NoError(t, addrRepo.Create(ctx, address))

	// Need a customer_address
	caRepo := NewCustomerAddressRepo(pool)
	customerAddress := &customers.CustomerAddress{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		AddressID:  address.ID,
		Name:       strPtr("Partner Location"),
		Type:       shared.CustomerAddressTypeCommercial,
		CreatedAt:  time.Now(),
	}
	require.NoError(t, caRepo.Create(ctx, customerAddress))

	// Need a partner
	partnerRepo := postgrespartners.NewPartnerRepo(pool)
	partner := &partners.Partner{
		ID:        uuid.New(),
		Name:      "Test Partner " + uuid.New().String(),
		TaxID:     strPtr("TAX999999"),
		Status:    shared.PartnerStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, partnerRepo.Create(ctx, partner))

	// Create CustomerAddressPartner
	cap := &customers.CustomerAddressPartner{
		ID:                uuid.New(),
		CustomerAddressID: customerAddress.ID,
		PartnerID:         partner.ID,
		FromDate:          time.Now(),
		ToDate:            nil,
		Notes:             strPtr("Active partnership"),
		CreatedAt:         time.Now(),
	}
	err := repo.Create(ctx, cap)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, cap.ID)

	// GetByID
	got, err := repo.GetByID(ctx, cap.ID)
	require.NoError(t, err)
	require.Equal(t, cap.CustomerAddressID, got.CustomerAddressID)
	require.Equal(t, cap.PartnerID, got.PartnerID)
	require.Equal(t, cap.Notes, got.Notes)

	// ListByCustomerAddress
	caps, err := repo.ListByCustomerAddress(ctx, customerAddress.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(caps), 1)
	found := false
	for _, p := range caps {
		if p.ID == cap.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	// Delete
	err = repo.Delete(ctx, cap.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, cap.ID)
	require.Error(t, err)

	// Cleanup
	partnerRepo.Delete(ctx, partner.ID)
	caRepo.Delete(ctx, customerAddress.ID)
	addrRepo.Delete(ctx, address.ID)
	customerRepo.Delete(ctx, customer.ID)
}

func TestCustomerAddressPartnerRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCustomerAddressPartnerRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestCustomerAddressPartnerRepo_ListByCustomerAddress_Empty(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewCustomerAddressPartnerRepo(pool)

	// Create customer and address but no partners
	customerRepo := NewCustomerRepo(pool)
	customer := &customers.Customer{
		ID:        uuid.New(),
		Name:      "Empty Partner Customer " + uuid.New().String(),
		Phone:     strPtr("+1000000000"),
		Email:     strPtr("emptypartner_" + uuid.New().String() + "@example.com"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, customerRepo.Create(ctx, customer))

	addrRepo := postgresgeocoding.NewGeocodingAddressRepo(pool)
	address := &geocoding.Address{
		ID:         uuid.New(),
		Street:     "Empty Partner Street",
		Number:     strPtr("111"),
		City:       strPtr("Empty City"),
		Region:     strPtr("EC"),
		PostalCode: strPtr("11111"),
		CreatedAt:  time.Now(),
	}
	require.NoError(t, addrRepo.Create(ctx, address))

	caRepo := NewCustomerAddressRepo(pool)
	customerAddress := &customers.CustomerAddress{
		ID:         uuid.New(),
		CustomerID: customer.ID,
		AddressID:  address.ID,
		Name:       strPtr("No Partners"),
		Type:       shared.CustomerAddressTypeResidential,
		CreatedAt:  time.Now(),
	}
	require.NoError(t, caRepo.Create(ctx, customerAddress))

	caps, err := repo.ListByCustomerAddress(ctx, customerAddress.ID)
	require.NoError(t, err)
	require.Len(t, caps, 0)

	// Cleanup
	caRepo.Delete(ctx, customerAddress.ID)
	addrRepo.Delete(ctx, address.ID)
	customerRepo.Delete(ctx, customer.ID)
}
