#!/bin/bash
# ============================================================================
# 01-roles.sh — PostgreSQL roles for FSM (single-tenant)
# Uses environment variables for passwords
# ============================================================================

set -e

if [ -z "$FSM_API_PASSWORD" ]; then
  echo "Error: FSM_API_PASSWORD environment variable is required"
  exit 1
fi

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  -- Application roles (with IF NOT EXISTS via DO block)
  DO \$\$ BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'fsm_api') THEN
      CREATE ROLE fsm_api LOGIN PASSWORD '${FSM_API_PASSWORD}';
    END IF;
  END \$\$;

  DO \$\$ BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'fsm_backend') THEN
      CREATE ROLE fsm_backend LOGIN PASSWORD '${FSM_API_PASSWORD}';
    END IF;
  END \$\$;

  DO \$\$ BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'fsm_migrator') THEN
      CREATE ROLE fsm_migrator NOLOGIN;
    END IF;
  END \$\$;

  -- Role grants
  GRANT fsm_migrator TO fsm_backend;
  GRANT fsm_migrator TO fsm_api;

  -- search_path for application roles (industry-specific schemas added by industry pack)
  ALTER ROLE fsm_api SET search_path TO core, partners, customers, inventory, operations, planning, notifications, geocoding, shared;
  ALTER ROLE fsm_backend SET search_path TO core, partners, customers, inventory, operations, planning, notifications, geocoding, shared;
EOSQL
