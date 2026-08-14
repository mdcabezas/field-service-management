package postgrespartners

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/partners"
	"localis-backend/internal/model/shared"
)

func TestPartnerRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerRepo(pool)

	partner := &partners.Partner{
		ID:        uuid.New(),
		Name:      "Test Partner " + uuid.New().String(),
		TaxID:     strPtr("TAX" + uuid.New().String()),
		Status:    shared.PartnerStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, partner)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, partner.ID)

	got, err := repo.GetByID(ctx, partner.ID)
	require.NoError(t, err)
	require.Equal(t, partner.Name, got.Name)
	require.Equal(t, partner.TaxID, got.TaxID)
	require.Equal(t, partner.Status, got.Status)

	result, err := repo.List(ctx, 100, 0)
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.Total, 1)
	found := false
	for _, p := range result.Items {
		if p.ID == partner.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	newName := "Updated Partner " + uuid.New().String()
	partner.Name = newName
	partner.UpdatedAt = time.Now()
	err = repo.Update(ctx, partner.ID, partner)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, partner.ID)
	require.NoError(t, err)
	require.Equal(t, newName, updated.Name)

	err = repo.Delete(ctx, partner.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, partner.ID)
	require.Error(t, err)
}

func TestPartnerRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}
