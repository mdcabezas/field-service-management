#!/bin/bash
# ============================================================================
# provision-tenant.sh — Create a new tenant database for FSM SaaS (Option A)
#
# One PostgreSQL cluster, one database per tenant, each seeded fresh:
#   core init (02-07) + one industry pack (90-94).
# Each tenant gets its own backend instance pointing at its DATABASE_URL.
#
# Usage:
#   ./scripts/provision-tenant.sh <tenant_name> [industry_pack]
#
# Examples:
#   ./scripts/provision-tenant.sh altogasspa    # core + default pack
#   ./scripts/provision-tenant.sh municipal_a  municipal  # core + municipal pack (future)
#
# Requirements: docker, running fsm-postgres, superuser fsm_admin
# ============================================================================
set -euo pipefail

TENANT="${1:?Usage: provision-tenant.sh <tenant_name> [industry_pack]}"
INDUSTRY="${2:-gas}"
CONTAINER="${POSTGRES_CONTAINER:-fsm-postgres}"
DB_USER="${POSTGRES_USER:-fsm_admin}"
DB_NAME="fsm_${TENANT}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CORE_DIR="$ROOT/docker/init"
PACK_DIR="$ROOT/industry-packs/$INDUSTRY/initdb.d"

run_sql() {
  docker exec -i "$CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 "$@" 2>&1
}

if [ ! -d "$PACK_DIR" ]; then
  echo "ERROR: industry pack not found: $PACK_DIR"
  exit 1
fi

echo "==> Creating database: $DB_NAME (industry: $INDUSTRY)"
docker exec "$CONTAINER" psql -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;" >/dev/null

echo "==> Applying core schema+seeds"
for f in 02-schemas.sql 03-schema.sql 04-seeds.sql 05-comments.sql 06-pre-request.sql 07-seed-data-dev.sql; do
  file="$CORE_DIR/$f"
  [ -f "$file" ] || { echo "  !! missing core file: $file"; exit 1; }
  run_sql < "$file" | grep -iE "ERROR|FATAL" && exit 1
  echo "  OK $f"
done

echo "==> Applying industry pack: $INDUSTRY"
for f in "$PACK_DIR"/*.sql; do
  run_sql < "$f" | grep -iE "ERROR|FATAL" && exit 1
  echo "  OK $(basename "$f")"
done

echo "==> Done. Tenant '$TENANT' provisioned."
echo ""
echo "Backend instance: set DATABASE_URL=postgres://$DB_USER:...@postgres:5432/$DB_NAME?sslmode=disable"