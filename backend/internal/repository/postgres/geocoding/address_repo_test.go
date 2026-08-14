package postgresgeocoding

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/geocoding"
)

func TestGeocodingAddressRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewGeocodingAddressRepo(pool)

	addr := &geocoding.Address{
		ID:                 uuid.New(),
		Street:             "Test Street",
		Number:             strPtr("123"),
		Apartment:          strPtr("A"),
		Neighborhood:       strPtr("Center"),
		City:               strPtr("Santiago"),
		Region:             strPtr("RM"),
		LocationReferences: strPtr("Near plaza"),
		PostalCode:         strPtr("1234567"),
		Geom:               orb.Point{-70.65, -33.45},
		CreatedAt:          time.Now(),
	}
	err := repo.Create(ctx, addr)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, addr.ID)

	got, err := repo.GetByID(ctx, addr.ID)
	require.NoError(t, err)
	require.Equal(t, addr.Street, got.Street)
	require.Equal(t, addr.Number, got.Number)
	require.Equal(t, addr.Geom, got.Geom)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, a := range result.Items {
		if a.ID == addr.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	addr.Street = "Updated Street"
	err = repo.Update(ctx, addr.ID, addr)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, addr.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated Street", updated.Street)

	err = repo.Delete(ctx, addr.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, addr.ID)
	require.Error(t, err)
}

func TestGeocodingAddressRepo_FindNearby(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewGeocodingAddressRepo(pool)

	// Santiago center
	items, err := repo.FindNearby(ctx, -33.45, -70.65, 5000)
	require.NoError(t, err)
	// Just verify it runs, results depend on data
	require.IsType(t, []geocoding.Address{}, items)
}

func TestGeocodingAddressRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewGeocodingAddressRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
