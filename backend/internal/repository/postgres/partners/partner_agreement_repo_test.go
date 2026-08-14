package postgrespartners

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/partners"
)

func TestPartnerAgreementRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerAgreementRepo(pool)

	partnerID := uuid.MustParse("30000000-0000-0000-0000-000000000001")
	svcType := uuid.MustParse("1d876cdd-8111-4060-9993-553dc4fb7f31")

	agreement := &partners.PartnerAgreement{
		ID:          uuid.New(),
		PartnerID:   partnerID,
		ServiceType: svcType,
		Rate:        floatPtr(100.50),
		Active:      true,
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, agreement)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, agreement.ID)

	got, err := repo.GetByID(ctx, agreement.ID)
	require.NoError(t, err)
	require.Equal(t, partnerID, got.PartnerID)
	require.Equal(t, svcType, got.ServiceType)

	result, err := repo.ListByPartner(ctx, partnerID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, a := range result {
		if a.ID == agreement.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	agreement.Active = false
	err = repo.Update(ctx, agreement.ID, agreement)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, agreement.ID)
	require.NoError(t, err)
	require.False(t, updated.Active)

	err = repo.Delete(ctx, agreement.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, agreement.ID)
	require.Error(t, err)
}

func TestPartnerAgreementRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerAgreementRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestPartnerContactRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerContactRepo(pool)

	partnerID := uuid.MustParse("30000000-0000-0000-0000-000000000001")

	contact := &partners.PartnerContact{
		ID:        uuid.New(),
		PartnerID: partnerID,
		Name:      "Contact " + uuid.New().String(),
		Position:  strPtr("Manager"),
		Phone:     strPtr("+1234567890"),
		Email:     strPtr("contact_" + uuid.New().String() + "@example.com"),
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, contact)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, contact.ID)

	got, err := repo.GetByID(ctx, contact.ID)
	require.NoError(t, err)
	require.Equal(t, partnerID, got.PartnerID)
	require.Equal(t, contact.Name, got.Name)

	result, err := repo.ListByPartner(ctx, partnerID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, c := range result {
		if c.ID == contact.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, contact.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, contact.ID)
	require.Error(t, err)
}

func TestPartnerContactRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerContactRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestPartnerAgreementFormRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerAgreementFormRepo(pool)

	agreementID := uuid.MustParse("32000000-0000-0000-0000-000000000001")

	form := &partners.PartnerAgreementForm{
		ID:             uuid.New(),
		AgreementID:    agreementID,
		WorkType:       "maintenance",
		FormTemplateID: uuidPtr(uuid.MustParse("10000000-0000-0000-0000-000000000001")),
		Quantity:       5,
		CreatedAt:      time.Now(),
	}
	err := repo.Create(ctx, form)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, form.ID)

	got, err := repo.GetByID(ctx, form.ID)
	require.NoError(t, err)
	require.Equal(t, agreementID, got.AgreementID)
	require.Equal(t, form.WorkType, got.WorkType)
	require.Equal(t, form.Quantity, got.Quantity)

	result, err := repo.ListByAgreement(ctx, agreementID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, f := range result {
		if f.ID == form.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, form.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, form.ID)
	require.Error(t, err)
}

func TestPartnerAgreementFormRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerAgreementFormRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestPartnerAgreementDocRepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerAgreementDocRepo(pool)

	agreementID := uuid.MustParse("32000000-0000-0000-0000-000000000001")

	doc := &partners.PartnerAgreementDoc{
		ID:          uuid.New(),
		AgreementID: agreementID,
		WorkType:    "installation",
		DocName:     "Contract " + uuid.New().String(),
		Required:    true,
		CreatedAt:   time.Now(),
	}
	err := repo.Create(ctx, doc)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, doc.ID)

	got, err := repo.GetByID(ctx, doc.ID)
	require.NoError(t, err)
	require.Equal(t, agreementID, got.AgreementID)
	require.Equal(t, doc.WorkType, got.WorkType)

	result, err := repo.ListByAgreement(ctx, agreementID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, d := range result {
		if d.ID == doc.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	err = repo.Delete(ctx, doc.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, doc.ID)
	require.Error(t, err)
}

func TestPartnerAgreementDocRepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewPartnerAgreementDocRepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func TestSLARepo_CRUD(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewSLARepo(pool)

	partnerID := uuid.MustParse("30000000-0000-0000-0000-000000000001")

	sla := &partners.SLA{
		ID:               uuid.New(),
		PartnerID:        partnerID,
		Name:             "SLA " + uuid.New().String(),
		Description:      strPtr("Test SLA"),
		WorkType:         uuidPtr(uuid.MustParse("1b0650f0-1e7c-4eed-a462-27dd182c2c64")),
		ResponseHours:    intPtr(4),
		ResolutionHours:  intPtr(24),
		ComplianceTarget: floatPtr(95.0),
		Active:           true,
		ValidFrom:        time.Now(),
		ValidUntil:       timePtr(time.Now().AddDate(1, 0, 0)),
		CreatedAt:        time.Now(),
	}
	err := repo.Create(ctx, sla)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, sla.ID)

	got, err := repo.GetByID(ctx, sla.ID)
	require.NoError(t, err)
	require.Equal(t, partnerID, got.PartnerID)
	require.Equal(t, sla.Name, got.Name)

	result, err := repo.ListByPartner(ctx, partnerID)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(result), 1)
	found := false
	for _, s := range result {
		if s.ID == sla.ID {
			found = true
			break
		}
	}
	require.True(t, found)

	sla.Active = false
	err = repo.Update(ctx, sla.ID, sla)
	require.NoError(t, err)
	updated, err := repo.GetByID(ctx, sla.ID)
	require.NoError(t, err)
	require.False(t, updated.Active)

	err = repo.Delete(ctx, sla.ID)
	require.NoError(t, err)
	_, err = repo.GetByID(ctx, sla.ID)
	require.Error(t, err)
}

func TestSLARepo_GetByID_NotFound(t *testing.T) {
	ctx := context.Background()
	pool := newPool(t)
	repo := NewSLARepo(pool)

	_, err := repo.GetByID(ctx, uuid.New())
	require.Error(t, err)
}

func intPtr(i int) *int              { return &i }
func floatPtr(f float64) *float64    { return &f }
func uuidPtr(u uuid.UUID) *uuid.UUID { return &u }
func timePtr(t time.Time) *time.Time { return &t }
