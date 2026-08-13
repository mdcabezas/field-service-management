-- ============================================================================
-- 20260813_catalogs_core.sql — Move 5 universal catalogs from domain_gas to core
--   operations.visit_types          (was domain_gas.visit_types)
--   operations.photo_findings       (was domain_gas.photo_findings)
--   operations.pre_visit_results    (was domain_gas.pre_visit_results)
--   operations.rejection_reasons    (was domain_gas.rejection_reasons)
--   partners.partner_service_types  (was domain_gas.partner_service_types)
-- Feature: generic FSM core owns these industry-agnostic catalogs; gas pack
--          seeds its codes into the core tables.
-- Preserves existing IDs so existing references stay valid.
-- Idempotent: safe to run multiple times. Compatible with PostgreSQL 16.
-- ============================================================================

BEGIN;

-- 1. Create the core catalog tables (match core DDL in docker/init/03-schema.sql)
CREATE TABLE IF NOT EXISTS operations.visit_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS operations.photo_findings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS operations.pre_visit_results (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS operations.rejection_reasons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  category shared.rejection_category NOT NULL,
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS partners.partner_service_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

-- 2. Migrate existing rows from the gas pack, preserving IDs
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'domain_gas' AND table_name = 'visit_types') THEN
    INSERT INTO operations.visit_types (id, code, name, active)
    SELECT id, code, name, active FROM domain_gas.visit_types
    ON CONFLICT (id) DO NOTHING;
  END IF;

  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'domain_gas' AND table_name = 'photo_findings') THEN
    INSERT INTO operations.photo_findings (id, code, name, active)
    SELECT id, code, name, active FROM domain_gas.photo_findings
    ON CONFLICT (id) DO NOTHING;
  END IF;

  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'domain_gas' AND table_name = 'pre_visit_results') THEN
    INSERT INTO operations.pre_visit_results (id, code, name, active)
    SELECT id, code, name, active FROM domain_gas.pre_visit_results
    ON CONFLICT (id) DO NOTHING;
  END IF;

  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'domain_gas' AND table_name = 'rejection_reasons') THEN
    INSERT INTO operations.rejection_reasons (id, code, name, category, active, created_at)
    SELECT id, code, name, category, active, created_at FROM domain_gas.rejection_reasons
    ON CONFLICT (id) DO NOTHING;
  END IF;

  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'domain_gas' AND table_name = 'partner_service_types') THEN
    INSERT INTO partners.partner_service_types (id, code, name, active)
    SELECT id, code, name, active FROM domain_gas.partner_service_types
    ON CONFLICT (id) DO NOTHING;
  END IF;
END $$;

-- 3. Seed generic industry-agnostic catalog codes (packs may add more)
INSERT INTO operations.visit_types (code, name) VALUES
  ('pre_visit', 'Pre-visit'),
  ('installation', 'Installation'),
  ('maintenance', 'Maintenance'),
  ('repair', 'Repair'),
  ('diagnosis', 'Diagnosis'),
  ('inspection', 'Inspection'),
  ('emergency', 'Emergency'),
  ('certification', 'Certification'),
  ('other', 'Other')
ON CONFLICT (code) DO NOTHING;

INSERT INTO operations.photo_findings (code, name) VALUES
  ('normal', 'Normal'),
  ('damage', 'Damage'),
  ('incomplete', 'Incomplete'),
  ('other', 'Other')
ON CONFLICT (code) DO NOTHING;

INSERT INTO operations.pre_visit_results (code, name) VALUES
  ('approved', 'Approved'),
  ('rejected', 'Rejected'),
  ('conditional', 'Conditional')
ON CONFLICT (code) DO NOTHING;

INSERT INTO partners.partner_service_types (code, name) VALUES
  ('installation', 'Installation'),
  ('maintenance', 'Maintenance'),
  ('repair', 'Repair'),
  ('certification', 'Certification'),
  ('emergency', 'Emergency'),
  ('other', 'Other')
ON CONFLICT (code) DO NOTHING;

INSERT INTO operations.rejection_reasons (code, name, category) VALUES
  ('customer_unavailable', 'Customer not available on site', 'customer'),
  ('customer_cancelled', 'Customer cancelled the visit', 'customer'),
  ('missing_documentation', 'Missing documentation', 'operational'),
  ('other', 'Other reason (specify)', 'operational')
ON CONFLICT (code) DO NOTHING;

-- 4. Repoint the FKs from the gas tables to the core catalogs
ALTER TABLE partners.partner_agreements DROP CONSTRAINT IF EXISTS fk_partner_agreements_service_type;
ALTER TABLE partners.slas DROP CONSTRAINT IF EXISTS fk_slas_work_type;
ALTER TABLE inventory.checklist_templates DROP CONSTRAINT IF EXISTS fk_checklist_templates_work_type;
ALTER TABLE operations.visits DROP CONSTRAINT IF EXISTS fk_visits_type;
ALTER TABLE operations.visits DROP CONSTRAINT IF EXISTS fk_visits_result;
ALTER TABLE operations.visits DROP CONSTRAINT IF EXISTS fk_visits_rejection_reason_id;
ALTER TABLE operations.visit_photos DROP CONSTRAINT IF EXISTS fk_visit_photos_finding_type;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_partner_agreements_service_type') THEN
    ALTER TABLE partners.partner_agreements
      ADD CONSTRAINT fk_partner_agreements_service_type
      FOREIGN KEY (service_type) REFERENCES partners.partner_service_types(id);
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_slas_work_type') THEN
    ALTER TABLE partners.slas
      ADD CONSTRAINT fk_slas_work_type
      FOREIGN KEY (work_type) REFERENCES operations.visit_types(id);
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_checklist_templates_work_type') THEN
    ALTER TABLE inventory.checklist_templates
      ADD CONSTRAINT fk_checklist_templates_work_type
      FOREIGN KEY (work_type) REFERENCES operations.visit_types(id);
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_visits_type') THEN
    ALTER TABLE operations.visits
      ADD CONSTRAINT fk_visits_type
      FOREIGN KEY (type) REFERENCES operations.visit_types(id);
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_visits_result') THEN
    ALTER TABLE operations.visits
      ADD CONSTRAINT fk_visits_result
      FOREIGN KEY (result) REFERENCES operations.pre_visit_results(id);
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_visits_rejection_reason_id') THEN
    ALTER TABLE operations.visits
      ADD CONSTRAINT fk_visits_rejection_reason_id
      FOREIGN KEY (rejection_reason_id) REFERENCES operations.rejection_reasons(id);
  END IF;
END $$;

DO $$ BEGIN
  IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_visit_photos_finding_type') THEN
    ALTER TABLE operations.visit_photos
      ADD CONSTRAINT fk_visit_photos_finding_type
      FOREIGN KEY (finding_type) REFERENCES operations.photo_findings(id);
  END IF;
END $$;

-- 5. Drop the now-unused gas-pack tables
DROP TABLE IF EXISTS domain_gas.visit_types;
DROP TABLE IF EXISTS domain_gas.photo_findings;
DROP TABLE IF EXISTS domain_gas.pre_visit_results;
DROP TABLE IF EXISTS domain_gas.rejection_reasons;
DROP TABLE IF EXISTS domain_gas.partner_service_types;

-- 6. Drop leftover domain_gas.vehicle_types (moved to core inventory.vehicle_types earlier; orphaned residue)
DROP TABLE IF EXISTS domain_gas.vehicle_types;

COMMIT;