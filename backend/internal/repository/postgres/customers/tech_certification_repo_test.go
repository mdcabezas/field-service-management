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

func TestTechCertificationRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewTechCertificationRepo(pool)

	// Need a technician first (which needs a user)
	techRepo := NewTechnicianRepo(pool)
	userRepo := postgrescore.NewCoreUserRepo(pool)

	user := &core.User{
		Email:     "cert_user_" + uuid.New().String() + "@example.com",
		Role:      "technician",
		Name:      "Cert Tech User " + uuid.New().String(),
		CreatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, user))

	tech := &customers.Technician{
		ID:          uuid.New(),
		UserID:      strPtr(user.ID.String()),
		Name:        "Cert Tech " + uuid.New().String(),
		IsActive:    true,
		Specialties: []string{"HVAC"},
		CreatedAt:   time.Now(),
	}
	require.NoError(t, techRepo.Create(ctx, tech))

	// Create Certification
	cert := &customers.TechCertification{
		ID:         uuid.New(),
		TechID:     tech.ID,
		CertID:     nil,
		Number:     strPtr("CERT-12345"),
		Issuer:     strPtr("HVAC Authority"),
		IssueDate:  timePtr(time.Now()),
		ExpiryDate: timePtr(time.Now().AddDate(1, 0, 0)),
		CreatedAt:  time.Now(),
	}
	err := repo.Create(ctx, cert)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, cert.ID)

	// GetByID
	got, err := repo.GetByID(ctx, cert.ID)
	require.NoError(t, err)
	require.Equal(t, cert.TechID, got.TechID)
	require.Equal(t, cert.Number, got.Number)
	require.Equal(t, cert.Issuer, got.Issuer)

	// ListByTech
	certs, err := repo.ListByTech(ctx, tech.ID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(certs), 1)
	found := false
	for _, c := range certs {
		if c.ID == cert.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	// Delete
	err = repo.Delete(ctx, cert.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, cert.ID)
	require.Error(t, err)

	// Cleanup
	techRepo.Delete(ctx, tech.ID)
	userRepo.Delete(ctx, user.ID.String())
}

func TestTechCertificationRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewTechCertificationRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestTechCertificationRepo_ListByTech_Empty(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewTechCertificationRepo(pool)

	// Create a technician without certifications
	techRepo := NewTechnicianRepo(pool)
	userRepo := postgrescore.NewCoreUserRepo(pool)

	user := &core.User{
		Email:     "empty_cert_user_" + uuid.New().String() + "@example.com",
		Role:      "technician",
		Name:      "No Cert Tech User " + uuid.New().String(),
		CreatedAt: time.Now(),
	}
	require.NoError(t, userRepo.Create(ctx, user))

	tech := &customers.Technician{
		ID:          uuid.New(),
		UserID:      strPtr(user.ID.String()),
		Name:        "No Cert Tech " + uuid.New().String(),
		IsActive:    true,
		Specialties: []string{"General"},
		CreatedAt:   time.Now(),
	}
	require.NoError(t, techRepo.Create(ctx, tech))

	certs, err := repo.ListByTech(ctx, tech.ID)
	require.NoError(t, err)
	require.Len(t, certs, 0)

	techRepo.Delete(ctx, tech.ID)
	userRepo.Delete(ctx, user.ID.String())
}

func timePtr(t time.Time) *time.Time {
	return &t
}
