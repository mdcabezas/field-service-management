-- ============================================================================
-- 20260813_route_types_core.sql — Move route_types catalog from domain_gas to core planning
-- Feature: generic FSM core owns the route_types catalog; gas pack seeds its codes
-- Preserves existing IDs so planning.routes.type references stay valid.
-- Idempotent: safe to run multiple times.
-- Compatible with PostgreSQL 16.
-- ============================================================================

BEGIN;

-- 1. Create the core catalog table (matches core DDL in docker/init/03-schema.sql)
CREATE TABLE IF NOT EXISTS planning.route_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

-- 2. Migrate existing rows from the gas pack, preserving IDs
INSERT INTO planning.route_types (id, code, name, active)
SELECT id, code, name, active FROM domain_gas.route_types
ON CONFLICT (id) DO NOTHING;

-- 3. Seed generic industry-agnostic route types (packs may add more)
INSERT INTO planning.route_types (code, name) VALUES
  ('maintenance', 'Maintenance'),
  ('installation', 'Installation'),
  ('repair', 'Repair'),
  ('inspection', 'Inspection'),
  ('delivery', 'Delivery'),
  ('collection', 'Collection'),
  ('emergency', 'Emergency'),
  ('other', 'Other')
ON CONFLICT (code) DO NOTHING;

-- 4. Repoint the FK on planning.routes to the core catalog
ALTER TABLE planning.routes DROP CONSTRAINT IF EXISTS fk_routes_type;

DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_routes_type'
  ) THEN
    ALTER TABLE planning.routes
      ADD CONSTRAINT fk_routes_type
      FOREIGN KEY (type) REFERENCES planning.route_types(id);
  END IF;
END $$;

-- 5. Drop the now-unused gas-pack table
DROP TABLE IF EXISTS domain_gas.route_types;

COMMIT;
