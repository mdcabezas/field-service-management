package scripts

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const defaultCoreTestDB = "postgres://fsm_admin:change_me_in_prod@localhost:5432/fsm_test?sslmode=disable"

func testPool(t *testing.T, defaultURL string) *pgxpool.Pool {
	t.Helper()
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = defaultURL
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })
	return pool
}

func TestTenantProvision_CoreSchema(t *testing.T) {
	ctx := context.Background()
	pool := testPool(t, defaultCoreTestDB)

	// Core tables should exist
	coreTables := []string{
		"core.users",
		"core.tech_roles",
		"core.plan_audit_log",
		"partners.partners",
		"partners.partner_agreements",
		"customers.customers",
		"customers.properties",
		"customers.technicians",
		"inventory.materials",
		"inventory.tools",
		"inventory.epp_items",
		"inventory.vehicles",
		"inventory.vehicle_types",
		"inventory.checklist_templates",
		"planning.daily_plans",
		"planning.routes",
		"planning.route_types",
		"operations.visit_types",
		"operations.photo_findings",
		"operations.pre_visit_results",
		"operations.rejection_reasons",
		"partners.partner_service_types",
		"operations.visits",
		"operations.visit_assignments",
		"notifications.notifications",
		"geocoding.addresses",
	}

	for _, table := range coreTables {
		var count int
		err := pool.QueryRow(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema||'.'||table_name = $1", table).Scan(&count)
		require.NoError(t, err, "checking table %s", table)
		require.Equal(t, 1, count, "table %s should exist", table)
	}
}

func TestTenantProvision_CoreCatalogData(t *testing.T) {
	ctx := context.Background()
	pool := testPool(t, defaultCoreTestDB)

	// Core vehicle_types should exist and have data
	var vtCoreCount int
	err := pool.QueryRow(ctx, "SELECT count(*) FROM inventory.vehicle_types WHERE active = true").Scan(&vtCoreCount)
	require.NoError(t, err)
	require.GreaterOrEqual(t, vtCoreCount, 5, "should have at least 5 core vehicle types")

	// Core route_types should exist and have data
	var rtCoreCount int
	err = pool.QueryRow(ctx, "SELECT count(*) FROM planning.route_types WHERE active = true").Scan(&rtCoreCount)
	require.NoError(t, err)
	require.GreaterOrEqual(t, rtCoreCount, 5, "should have at least 5 core route types")

	// Core seeds should be present
	var count int
	err = pool.QueryRow(ctx, "SELECT count(*) FROM inventory.vehicle_types WHERE code IN ('truck','crane','van','crane_truck','other') AND active = true").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 5, count, "should have 5 core vehicle types seeded")
}

func TestTenantProvision_VehicleTypesFK(t *testing.T) {
	ctx := context.Background()
	pool := testPool(t, defaultCoreTestDB)

	// FK should point to inventory.vehicle_types
	var fkTarget string
	err := pool.QueryRow(ctx, `
		SELECT ccu.table_schema||'.'||ccu.table_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.constraint_column_usage ccu ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_name = 'vehicles_type_fkey'
		  AND tc.table_schema = 'inventory'
		  AND tc.table_name = 'vehicles'
	`).Scan(&fkTarget)
	require.NoError(t, err)
	require.Equal(t, "inventory.vehicle_types", fkTarget, "FK should point to core inventory.vehicle_types")

	// domain_gas.vehicle_types should NOT exist (moved to core)
	var vtCount int
	err = pool.QueryRow(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema='domain_gas' AND table_name='vehicle_types'").Scan(&vtCount)
	require.NoError(t, err)
	require.Equal(t, 0, vtCount, "domain_gas.vehicle_types should not exist (moved to core)")

	// domain_gas.route_types should NOT exist (moved to core)
	var rtCount int
	err = pool.QueryRow(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema='domain_gas' AND table_name='route_types'").Scan(&rtCount)
	require.NoError(t, err)
	require.Equal(t, 0, rtCount, "domain_gas.route_types should not exist (moved to core)")

	// The 5 universal catalogs should NOT exist in domain_gas (moved to core)
	universalCatalogs := []string{"visit_types", "photo_findings", "pre_visit_results", "rejection_reasons", "partner_service_types"}
	for _, table := range universalCatalogs {
		var count int
		err = pool.QueryRow(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema='domain_gas' AND table_name=$1", table).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count, "domain_gas.%s should not exist (moved to core)", table)
	}
}
