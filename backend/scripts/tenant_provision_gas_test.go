package scripts

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const defaultGasTestDB = "postgres://fsm_admin:change_me_in_prod@localhost:5432/fsm_test_gas?sslmode=disable"

// gasPackProvisioned reports whether the domain_gas schema exists in the
// connected database. Gas-pack tests skip when the pack is not provisioned so
// the core-only database (fsm_test) does not fail them.
func gasPackProvisioned(t *testing.T, pool *pgxpool.Pool) bool {
	t.Helper()
	var count int
	err := pool.QueryRow(context.Background(), "SELECT count(*) FROM information_schema.schemata WHERE schema_name = 'domain_gas'").Scan(&count)
	require.NoError(t, err)
	return count == 1
}

func TestTenantProvision_GasPack(t *testing.T) {
	ctx := context.Background()
	pool := testPool(t, defaultGasTestDB)
	if !gasPackProvisioned(t, pool) {
		t.Skip("domain_gas schema not provisioned; skipping gas pack verification")
	}

	// Gas pack domain schema (universal catalogs live in core)
	domainTables := []string{
		"domain_gas.measurement_types",
		"domain_gas.property_types",
		"domain_gas.certifications",
		"domain_gas.property_assets",
	}

	for _, table := range domainTables {
		var count int
		err := pool.QueryRow(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema||'.'||table_name = $1", table).Scan(&count)
		require.NoError(t, err, "checking table %s", table)
		require.Equal(t, 1, count, "table %s should exist", table)
	}

	// Gas dev-data assigns vehicle types from the core catalog
	var typedVehicles int
	err := pool.QueryRow(ctx, "SELECT count(*) FROM inventory.vehicles WHERE type IS NOT NULL").Scan(&typedVehicles)
	require.NoError(t, err)
	require.Greater(t, typedVehicles, 0, "gas dev-data should populate vehicle types")
}

func TestTenantProvision_GasPackFKs(t *testing.T) {
	ctx := context.Background()
	pool := testPool(t, defaultGasTestDB)
	if !gasPackProvisioned(t, pool) {
		t.Skip("domain_gas schema not provisioned; skipping gas pack FK verification")
	}

	// Gas pack FKs should exist; the moved-to-core FKs must point to core catalogs.
	expectedFKs := map[string]string{
		"fk_partner_agreements_service_type": "partners.partner_service_types",
		"fk_slas_work_type":                  "operations.visit_types",
		"fk_properties_type":                 "domain_gas.property_types",
		"fk_tech_certifications_cert_id":     "domain_gas.certifications",
		"fk_checklist_templates_work_type":   "operations.visit_types",
		"fk_routes_type":                     "planning.route_types",
		"fk_visits_type":                     "operations.visit_types",
		"fk_visits_result":                   "operations.pre_visit_results",
		"fk_visits_rejection_reason_id":      "operations.rejection_reasons",
		"fk_visit_measurements_type":         "domain_gas.measurement_types",
		"fk_visit_photos_finding_type":       "operations.photo_findings",
	}

	for fkName, expectedTable := range expectedFKs {
		var target string
		err := pool.QueryRow(ctx, `
			SELECT ccu.table_schema||'.'||ccu.table_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
			WHERE tc.constraint_name = $1
		`, fkName).Scan(&target)
		require.NoError(t, err, "FK %s should exist", fkName)
		require.Equal(t, expectedTable, target, "FK %s should point to %s", fkName, expectedTable)
	}
}
