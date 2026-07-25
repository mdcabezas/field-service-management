-- ============================================================================
-- 01-domain-gas-schema.sql — Gas Chile industry pack: lookup tables, SEC certifications, rejection reasons
-- Runs after core schema (prefix 90- ensures ordering via filename)
-- Idempotent: uses IF NOT EXISTS where possible
-- ============================================================================

SET search_path TO shared, public;

CREATE SCHEMA IF NOT EXISTS domain_gas;

-- ============================================================================
-- LOOKUP TABLES (replace enums)
-- ============================================================================

CREATE TABLE IF NOT EXISTS domain_gas.visit_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS domain_gas.measurement_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  default_unit TEXT,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS domain_gas.photo_findings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS domain_gas.vehicle_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS domain_gas.partner_service_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS domain_gas.pre_visit_results (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS domain_gas.route_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS domain_gas.property_types (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  columns_config JSONB DEFAULT '{}',
  active BOOLEAN DEFAULT true
);

-- ============================================================================
-- SEC CHILE CERTIFICATIONS
-- ============================================================================

CREATE TABLE IF NOT EXISTS domain_gas.certifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  class TEXT,
  description TEXT,
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- ============================================================================
-- REJECTION REASONS
-- ============================================================================

CREATE TABLE IF NOT EXISTS domain_gas.rejection_reasons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  category shared.rejection_category NOT NULL,
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- ============================================================================
-- PROPERTY ASSETS (poles, meter, etc.)
-- ============================================================================

CREATE TABLE IF NOT EXISTS domain_gas.property_assets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  property_id UUID NOT NULL REFERENCES customers.properties(id),
  type TEXT NOT NULL,
  value JSONB,
  created_at TIMESTAMPTZ DEFAULT now(),
  UNIQUE(property_id, type)
);

-- ============================================================================
-- ADD FOREIGN KEYS TO CORE TABLES (idempotent)
-- ============================================================================

-- partners.partner_agreements.service_type → domain_gas.partner_service_types
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_partner_agreements_service_type'
  ) THEN
    ALTER TABLE partners.partner_agreements
      ADD CONSTRAINT fk_partner_agreements_service_type
      FOREIGN KEY (service_type) REFERENCES domain_gas.partner_service_types(id);
  END IF;
END $$;

-- partners.slas.work_type → domain_gas.visit_types
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_slas_work_type'
  ) THEN
    ALTER TABLE partners.slas
      ADD CONSTRAINT fk_slas_work_type
      FOREIGN KEY (work_type) REFERENCES domain_gas.visit_types(id);
  END IF;
END $$;

-- customers.properties.type → domain_gas.property_types
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_properties_type'
  ) THEN
    ALTER TABLE customers.properties
      ADD CONSTRAINT fk_properties_type
      FOREIGN KEY (type) REFERENCES domain_gas.property_types(id);
  END IF;
END $$;

-- customers.tech_certifications.cert_id → domain_gas.certifications
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_tech_certifications_cert_id'
  ) THEN
    ALTER TABLE customers.tech_certifications
      ADD CONSTRAINT fk_tech_certifications_cert_id
      FOREIGN KEY (cert_id) REFERENCES domain_gas.certifications(id);
  END IF;
END $$;

-- inventory.vehicles.type → domain_gas.vehicle_types
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_vehicles_type'
  ) THEN
    ALTER TABLE inventory.vehicles
      ADD CONSTRAINT fk_vehicles_type
      FOREIGN KEY (type) REFERENCES domain_gas.vehicle_types(id);
  END IF;
END $$;

-- inventory.checklist_templates.work_type → domain_gas.visit_types
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_checklist_templates_work_type'
  ) THEN
    ALTER TABLE inventory.checklist_templates
      ADD CONSTRAINT fk_checklist_templates_work_type
      FOREIGN KEY (work_type) REFERENCES domain_gas.visit_types(id);
  END IF;
END $$;

-- planning.routes.type → domain_gas.route_types
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_routes_type'
  ) THEN
    ALTER TABLE planning.routes
      ADD CONSTRAINT fk_routes_type
      FOREIGN KEY (type) REFERENCES domain_gas.route_types(id);
  END IF;
END $$;

-- operations.visits.type → domain_gas.visit_types
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_visits_type'
  ) THEN
    ALTER TABLE operations.visits
      ADD CONSTRAINT fk_visits_type
      FOREIGN KEY (type) REFERENCES domain_gas.visit_types(id);
  END IF;
END $$;

-- operations.visits.result → domain_gas.pre_visit_results
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_visits_result'
  ) THEN
    ALTER TABLE operations.visits
      ADD CONSTRAINT fk_visits_result
      FOREIGN KEY (result) REFERENCES domain_gas.pre_visit_results(id);
  END IF;
END $$;

-- operations.visits.rejection_reason_id → domain_gas.rejection_reasons
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_visits_rejection_reason_id'
  ) THEN
    ALTER TABLE operations.visits
      ADD CONSTRAINT fk_visits_rejection_reason_id
      FOREIGN KEY (rejection_reason_id) REFERENCES domain_gas.rejection_reasons(id);
  END IF;
END $$;

-- operations.visit_measurements.type → domain_gas.measurement_types
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_visit_measurements_type'
  ) THEN
    ALTER TABLE operations.visit_measurements
      ADD CONSTRAINT fk_visit_measurements_type
      FOREIGN KEY (type) REFERENCES domain_gas.measurement_types(id);
  END IF;
END $$;

-- operations.visit_photos.finding_type → domain_gas.photo_findings
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.table_constraints
    WHERE constraint_name = 'fk_visit_photos_finding_type'
  ) THEN
    ALTER TABLE operations.visit_photos
      ADD CONSTRAINT fk_visit_photos_finding_type
      FOREIGN KEY (finding_type) REFERENCES domain_gas.photo_findings(id);
  END IF;
END $$;