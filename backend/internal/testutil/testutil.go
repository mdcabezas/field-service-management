package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/auth"
	"localis-backend/internal/testutil/mocks"
)

func TestPool(t *testing.T) *pgxpool.Pool {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://fsm_admin:change_me_in_prod@localhost:5432/fsm_test_gas?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to create test pool: %v", err)
	}
	t.Cleanup(func() { pool.Close() })
	return pool
}

func NewMockLDAPAuth(t *testing.T) *auth.LDAPAuth {
	t.Helper()
	ldapAuth, err := auth.NewLDAPAuth("ldaps://localhost:636", "dc=test,dc=com", "password")
	if err != nil {
		t.Fatalf("failed to create LDAP auth: %v", err)
	}
	return ldapAuth
}

func NewMockJWTAuth(t *testing.T) *auth.JWTAuth {
	t.Helper()
	secret := "test-jwt-secret-" + uuid.New().String()
	jwtAuth, err := auth.NewJWTAuth(secret)
	if err != nil {
		t.Fatalf("failed to create JWT auth: %v", err)
	}
	return jwtAuth
}

func NewMockCoreUserRepo(t *testing.T) *mocks.CoreUserRepository {
	return mocks.NewCoreUserRepository(t)
}

func NewMockCoreTechRoleRepo(t *testing.T) *mocks.CoreTechRoleRepository {
	return mocks.NewCoreTechRoleRepository(t)
}

func NewMockCorePlanAuditLogRepo(t *testing.T) *mocks.CorePlanAuditLogRepository {
	return mocks.NewCorePlanAuditLogRepository(t)
}

func NewMockPartnerRepo(t *testing.T) *mocks.PartnerRepository {
	return mocks.NewPartnerRepository(t)
}

func NewMockPartnerContactRepo(t *testing.T) *mocks.PartnerContactRepository {
	return mocks.NewPartnerContactRepository(t)
}

func NewMockPartnerAgreementRepo(t *testing.T) *mocks.PartnerAgreementRepository {
	return mocks.NewPartnerAgreementRepository(t)
}

func NewMockPartnerAgreementDocRepo(t *testing.T) *mocks.PartnerAgreementDocRepository {
	return mocks.NewPartnerAgreementDocRepository(t)
}

func NewMockPartnerAgreementFormRepo(t *testing.T) *mocks.PartnerAgreementFormRepository {
	return mocks.NewPartnerAgreementFormRepository(t)
}

func NewMockSLARepo(t *testing.T) *mocks.SLARepository {
	return mocks.NewSLARepository(t)
}
