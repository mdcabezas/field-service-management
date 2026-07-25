-- ============================================================================
-- 03-schema.sql — DDL FSM Core (single-tenant)
-- Bounded context schemas, polymorphic split Option A,
-- updated_at + audit_change triggers, CHECK route↔plan via trigger
-- ============================================================================
-- Execute on PostgreSQL 16 with PostGIS 3.5

BEGIN;

-- ============================================================================
-- EXTENSIONS
-- ============================================================================
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ============================================================================
-- ENUMS (shared schema)
-- ============================================================================

SET search_path TO shared, public;

CREATE TYPE visit_status AS ENUM (
  'scheduled', 'en_route', 'in_progress',
  'completed', 'cancelled', 'rescheduled'
);

CREATE TYPE visit_priority AS ENUM (
  'low', 'normal', 'urgent', 'emergency'
);

CREATE TYPE visit_source AS ENUM (
  'email', 'phone', 'portal', 'whatsapp', 'other'
);

CREATE TYPE billing_to AS ENUM (
  'partner', 'customer'
);

CREATE TYPE photo_stage AS ENUM (
  'diagnosis', 'execution', 'review', 'report', 'other'
);

CREATE TYPE checkpoint_type AS ENUM (
  'arrival', 'departure', 'report_sent'
);

CREATE TYPE report_source AS ENUM (
  'system', 'paper'
);

CREATE TYPE measurement_result AS ENUM (
  'approved', 'rejected', 'pending', 'observation'
);

CREATE TYPE vehicle_status AS ENUM (
  'available', 'in_use', 'maintenance', 'retired'
);

CREATE TYPE tool_status AS ENUM (
  'available', 'in_use', 'maintenance', 'retired'
);

CREATE TYPE tech_status AS ENUM (
  'active', 'inactive', 'vacation', 'leave'
);

CREATE TYPE rental_type AS ENUM (
  'vehicle', 'tool', 'equipment'
);

CREATE TYPE route_status AS ENUM (
  'scheduled', 'in_progress', 'completed', 'cancelled'
);

CREATE TYPE epp_usage_status AS ENUM (
  'used', 'not_needed', 'missing'
);

CREATE TYPE rejection_category AS ENUM (
  'logistics', 'regulatory', 'customer', 'operational'
);

CREATE TYPE cost_type AS ENUM (
  'labor', 'vehicle', 'tool_depreciation', 'epp', 'overhead'
);

CREATE TYPE cost_unit AS ENUM (
  'hour', 'km', 'day', 'visit', 'unit'
);

CREATE TYPE customer_address_type AS ENUM (
  'residential', 'commercial', 'industrial', 'institutional'
);

CREATE TYPE partner_status AS ENUM (
  'active', 'inactive'
);

CREATE TYPE epp_type AS ENUM (
  'head', 'hands', 'body', 'eyes', 'ears', 'feet', 'respiratory'
);

CREATE TYPE epp_lifecycle AS ENUM (
  'disposable', 'reusable'
);

CREATE TYPE notification_type AS ENUM (
  'visit_assigned', 'visit_reminder', 'visit_cancelled',
  'visit_rescheduled', 'checklist_pending', 'report_pending',
  'sla_warning', 'sla_expired', 'maintenance_scheduled',
  'certification_expired', 'other'
);

CREATE TYPE notification_channel AS ENUM (
  'email', 'sms', 'whatsapp', 'push', 'in_app'
);

CREATE TYPE notification_status AS ENUM (
  'pending', 'sent', 'failed', 'read'
);

CREATE TYPE audit_entity_type AS ENUM (
  'daily_load_material', 'daily_load_tool', 'daily_load_epp',
  'visit_checklist_material', 'visit_checklist_tool', 'visit_checklist_epp',
  'visit_material_usage', 'visit_tool_usage', 'visit_epp_usage'
);

CREATE TYPE audit_action AS ENUM (
  'add', 'update', 'delete'
);

CREATE TYPE maintenance_type AS ENUM (
  'vehicle', 'tool', 'equipment'
);

-- ============================================================================
-- HELPER FUNCTIONS
-- ============================================================================

SET search_path TO core, public;

-- app_user_id: retrieves employee_number from JWT claim sub (signed by Go backend)
CREATE OR REPLACE FUNCTION app_user_id()
RETURNS TEXT AS $$
  SELECT NULLIF(current_setting('request.jwt.claim.sub', true), '')::TEXT;
$$ LANGUAGE sql STABLE;

-- touch_updated_at: trigger function for automatic updated_at
CREATE OR REPLACE FUNCTION touch_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- audit_change: trigger function for plan_audit_log
CREATE OR REPLACE FUNCTION audit_change()
RETURNS TRIGGER AS $$
DECLARE
  v_entity_type shared.audit_entity_type;
  v_action shared.audit_action;
  v_old JSONB;
  v_new JSONB;
BEGIN
  IF TG_OP = 'INSERT' THEN
    v_action := 'add';
    v_old := NULL;
    v_new := to_jsonb(NEW);
  ELSIF TG_OP = 'UPDATE' THEN
    v_action := 'update';
    v_old := to_jsonb(OLD);
    v_new := to_jsonb(NEW);
  ELSIF TG_OP = 'DELETE' THEN
    v_action := 'delete';
    v_old := to_jsonb(OLD);
    v_new := NULL;
  END IF;

  CASE TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME
    WHEN 'planning.daily_load_materials'  THEN v_entity_type := 'daily_load_material';
    WHEN 'planning.daily_load_tools'      THEN v_entity_type := 'daily_load_tool';
    WHEN 'planning.daily_load_epps'       THEN v_entity_type := 'daily_load_epp';
    WHEN 'operations.visit_checklist_materials' THEN v_entity_type := 'visit_checklist_material';
    WHEN 'operations.visit_checklist_tools'     THEN v_entity_type := 'visit_checklist_tool';
    WHEN 'operations.visit_checklist_epps'      THEN v_entity_type := 'visit_checklist_epp';
    WHEN 'operations.visit_material_usages'     THEN v_entity_type := 'visit_material_usage';
    WHEN 'operations.visit_tool_usages'         THEN v_entity_type := 'visit_tool_usage';
    WHEN 'operations.visit_epp_usages'          THEN v_entity_type := 'visit_epp_usage';
    ELSE v_entity_type := NULL;
  END CASE;

  IF v_entity_type IS NOT NULL THEN
    INSERT INTO plan_audit_log (entity_type, entity_id, action, old_data, new_data, user_id)
    VALUES (
      v_entity_type,
      CASE WHEN TG_OP = 'DELETE' THEN OLD.id ELSE NEW.id END,
      v_action,
      v_old,
      v_new,
      core.app_user_id()
    );
  END IF;

  IF TG_OP = 'DELETE' THEN RETURN OLD; END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- validate_visit_route_plan: D2 trigger — route and daily_plan must be consistent
CREATE OR REPLACE FUNCTION validate_visit_route_plan()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.route_id IS NOT NULL AND NEW.daily_plan_id IS NOT NULL THEN
    IF NEW.daily_plan_id != (SELECT daily_plan_id FROM planning.routes WHERE id = NEW.route_id) THEN
      RAISE EXCEPTION 'daily_plan_id does not match the specified route (D2)';
    END IF;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- SHARED (catalogs — before other bounded contexts due to dependencies)
-- Only rejection_category lives here (it is generic)
-- ============================================================================

CREATE TABLE shared.report_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  fields_json JSONB NOT NULL DEFAULT '[]',
  created_at TIMESTAMPTZ DEFAULT now()
);

-- ============================================================================
-- GEOCODING (global catalog — before customers due to FK)
-- ============================================================================

CREATE TABLE geocoding.addresses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  street TEXT NOT NULL,
  number TEXT,
  apartment TEXT,
  neighborhood TEXT,
  city TEXT,
  region TEXT,
  location_references TEXT,
  postal_code TEXT,
  geom GEOMETRY(Point, 4326),
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_addresses_geom ON geocoding.addresses USING GIST(geom);

-- ============================================================================
-- CORE (no companies, no company_id)
-- ============================================================================

CREATE TABLE core.users (
  employee_number TEXT PRIMARY KEY,
  email TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'technician',
  name TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE UNIQUE INDEX idx_users_email ON core.users(email);

CREATE TABLE core.tech_roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE core.plan_audit_log (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  entity_type shared.audit_entity_type NOT NULL,
  entity_id UUID NOT NULL,
  action shared.audit_action NOT NULL,
  old_data JSONB,
  new_data JSONB,
  user_id TEXT,
  reason TEXT,
  timestamp TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_plan_audit_log_entity ON core.plan_audit_log(entity_type, entity_id);
CREATE INDEX idx_plan_audit_log_user ON core.plan_audit_log(user_id);
CREATE INDEX idx_plan_audit_log_timestamp ON core.plan_audit_log(timestamp);

-- ============================================================================
-- PARTNERS
-- ============================================================================

CREATE TABLE partners.partners (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  tax_id TEXT,
  status shared.partner_status DEFAULT 'active',
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE partners.partner_contacts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  partner_id UUID NOT NULL REFERENCES partners.partners(id),
  name TEXT NOT NULL,
  position TEXT,
  phone TEXT,
  email TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_partner_contacts_partner ON partners.partner_contacts(partner_id);

CREATE TABLE partners.partner_agreements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  partner_id UUID NOT NULL REFERENCES partners.partners(id),
  service_type UUID NOT NULL,  -- FK added by industry pack
  rate NUMERIC,
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_partner_agreements_partner ON partners.partner_agreements(partner_id);

CREATE TABLE partners.partner_agreement_docs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  agreement_id UUID NOT NULL REFERENCES partners.partner_agreements(id),
  work_type TEXT NOT NULL,
  doc_name TEXT NOT NULL,
  required BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_partner_agreement_docs_agreement ON partners.partner_agreement_docs(agreement_id);

CREATE TABLE partners.partner_agreement_forms (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  agreement_id UUID NOT NULL REFERENCES partners.partner_agreements(id),
  work_type TEXT NOT NULL,
  form_template_id UUID REFERENCES shared.report_templates(id),
  quantity INTEGER DEFAULT 1,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_partner_agreement_forms_agreement ON partners.partner_agreement_forms(agreement_id);

CREATE TABLE partners.slas (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  partner_id UUID NOT NULL REFERENCES partners.partners(id),
  name TEXT NOT NULL,
  description TEXT,
  work_type UUID,  -- FK added by industry pack
  response_hours INTEGER,
  resolution_hours INTEGER,
  compliance_target NUMERIC,
  active BOOLEAN DEFAULT true,
  valid_from DATE DEFAULT CURRENT_DATE,
  valid_until DATE,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_slas_partner ON partners.slas(partner_id);

-- ============================================================================
-- CUSTOMERS
-- ============================================================================

CREATE TABLE customers.customers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  phone TEXT,
  email TEXT,
  tax_id TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE customers.customer_addresses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id UUID NOT NULL REFERENCES customers.customers(id),
  address_id UUID NOT NULL REFERENCES geocoding.addresses(id),
  name TEXT,
  type shared.customer_address_type DEFAULT 'residential',
  created_at TIMESTAMPTZ DEFAULT now(),
  UNIQUE(customer_id, address_id)
);

CREATE INDEX idx_customer_addresses_customer ON customers.customer_addresses(customer_id);
CREATE INDEX idx_customer_addresses_address ON customers.customer_addresses(address_id);

CREATE TABLE customers.customer_address_partners (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_address_id UUID NOT NULL REFERENCES customers.customer_addresses(id),
  partner_id UUID NOT NULL REFERENCES partners.partners(id),
  from_date DATE NOT NULL DEFAULT CURRENT_DATE,
  to_date DATE,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_customer_address_partners_ca ON customers.customer_address_partners(customer_address_id);
CREATE INDEX idx_customer_address_partners_partner ON customers.customer_address_partners(partner_id);

CREATE TABLE customers.properties (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_address_id UUID NOT NULL REFERENCES customers.customer_addresses(id),
  name TEXT,
  type UUID,  -- FK added by industry pack
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_properties_ca ON customers.properties(customer_address_id);

CREATE TABLE customers.technicians (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id TEXT REFERENCES core.users(employee_number),
  name TEXT NOT NULL,
  is_active BOOLEAN DEFAULT true,
  specialties TEXT[],
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE customers.tech_certifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tech_id UUID NOT NULL REFERENCES customers.technicians(id),
  cert_id UUID,  -- FK added by industry pack
  number TEXT,
  issuer TEXT,
  issue_date DATE,
  expiry_date DATE,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_tech_certifications_tech ON customers.tech_certifications(tech_id);

-- ============================================================================
-- INVENTORY
-- ============================================================================

CREATE TABLE inventory.materials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sku TEXT,
  name TEXT NOT NULL,
  unit TEXT NOT NULL,
  unit_cost NUMERIC DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE inventory.tools (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code TEXT,
  name TEXT NOT NULL,
  status shared.tool_status DEFAULT 'available',
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE inventory.epp_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  type shared.epp_type NOT NULL,
  lifecycle shared.epp_lifecycle DEFAULT 'reusable',
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE inventory.vehicles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type UUID,  -- FK added by industry pack
  license_plate TEXT,
  name TEXT NOT NULL,
  brand TEXT,
  model TEXT,
  year INTEGER,
  status shared.vehicle_status DEFAULT 'available',
  capacity TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE inventory.rentals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type shared.rental_type NOT NULL,
  supplier TEXT,
  item_description TEXT NOT NULL,
  daily_cost NUMERIC,
  hourly_cost NUMERIC,
  fixed_cost NUMERIC,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE inventory.cost_rates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type shared.cost_type NOT NULL,
  reference_id UUID,
  reference_type TEXT,
  value NUMERIC NOT NULL,
  unit shared.cost_unit NOT NULL,
  valid_from DATE DEFAULT CURRENT_DATE,
  valid_until DATE,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_cost_rates_type ON inventory.cost_rates(type);

CREATE TABLE inventory.maintenance_records (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type shared.maintenance_type NOT NULL,
  reference_id UUID NOT NULL,
  date DATE NOT NULL,
  cost NUMERIC,
  supplier TEXT,
  description TEXT,
  next_date DATE,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_maintenance_records_ref ON inventory.maintenance_records(type, reference_id);

CREATE TABLE inventory.maintenance_schedules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type shared.maintenance_type NOT NULL,
  reference_id UUID NOT NULL,
  frequency_km INTEGER,
  frequency_days INTEGER,
  last_service_date DATE,
  last_service_km INTEGER,
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_maintenance_schedules_ref ON inventory.maintenance_schedules(type, reference_id);

CREATE TABLE inventory.checklist_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  work_type UUID,  -- FK added by industry pack
  created_at TIMESTAMPTZ DEFAULT now()
);

-- ============================================================================
-- PLANNING (before operations due to visits↔routes↔daily_plans dependencies)
-- ============================================================================

CREATE TABLE planning.daily_plans (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  date DATE NOT NULL,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_daily_plans_date ON planning.daily_plans(date);

CREATE TABLE planning.daily_plan_assignments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  daily_plan_id UUID NOT NULL REFERENCES planning.daily_plans(id),
  tech_id UUID NOT NULL REFERENCES customers.technicians(id),
  role_id UUID NOT NULL REFERENCES core.tech_roles(id),
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_daily_plan_assignments_plan ON planning.daily_plan_assignments(daily_plan_id);
CREATE INDEX idx_daily_plan_assignments_tech ON planning.daily_plan_assignments(tech_id);

-- Option A split: daily_load_items → 3 child tables
CREATE TABLE planning.daily_load_materials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  daily_plan_id UUID NOT NULL REFERENCES planning.daily_plans(id),
  material_id UUID NOT NULL REFERENCES inventory.materials(id),
  loaded_quantity NUMERIC NOT NULL,
  returned_quantity NUMERIC DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_dlm_plan ON planning.daily_load_materials(daily_plan_id);

CREATE TABLE planning.daily_load_tools (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  daily_plan_id UUID NOT NULL REFERENCES planning.daily_plans(id),
  tool_id UUID NOT NULL REFERENCES inventory.tools(id),
  loaded_quantity NUMERIC NOT NULL,
  returned_quantity NUMERIC DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_dlt_plan ON planning.daily_load_tools(daily_plan_id);

CREATE TABLE planning.daily_load_epps (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  daily_plan_id UUID NOT NULL REFERENCES planning.daily_plans(id),
  epp_id UUID NOT NULL REFERENCES inventory.epp_items(id),
  loaded_quantity NUMERIC NOT NULL,
  returned_quantity NUMERIC DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_dle_plan ON planning.daily_load_epps(daily_plan_id);

CREATE TABLE planning.routes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type UUID,  -- FK added by industry pack
  date DATE NOT NULL,
  daily_plan_id UUID REFERENCES planning.daily_plans(id),
  notes TEXT,
  status shared.route_status DEFAULT 'scheduled',
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_routes_daily_plan ON planning.routes(daily_plan_id);

-- ============================================================================
-- OPERATIONS
-- ============================================================================

CREATE TABLE operations.visits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  property_id UUID REFERENCES customers.properties(id),
  partner_id UUID REFERENCES partners.partners(id),
  partner_order_id TEXT,
  parent_visit_id UUID REFERENCES operations.visits(id),
  route_id UUID REFERENCES planning.routes(id),
  daily_plan_id UUID REFERENCES planning.daily_plans(id),
  type UUID NOT NULL,  -- FK added by industry pack
  status shared.visit_status NOT NULL DEFAULT 'scheduled',
  priority shared.visit_priority DEFAULT 'normal',
  source shared.visit_source,
  source_reference TEXT,
  billing_to shared.billing_to DEFAULT 'customer',
  partner_supervisor TEXT,
  result UUID,  -- FK added by industry pack
  rejection_reason_id UUID,  -- FK added by industry pack
  rejection_reason_detail TEXT,
  scheduled_at TIMESTAMPTZ NOT NULL,
  started_at TIMESTAMPTZ,
  completed_at TIMESTAMPTZ,
  alternative_location TEXT,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_visits_property ON operations.visits(property_id);
CREATE INDEX idx_visits_partner ON operations.visits(partner_id);
CREATE INDEX idx_visits_route ON operations.visits(route_id);
CREATE INDEX idx_visits_daily_plan ON operations.visits(daily_plan_id);
CREATE INDEX idx_visits_status ON operations.visits(status);
CREATE INDEX idx_visits_scheduled_at ON operations.visits(scheduled_at);
CREATE INDEX idx_visits_parent ON operations.visits(parent_visit_id);

CREATE TABLE operations.visit_assignments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  tech_id UUID NOT NULL REFERENCES customers.technicians(id),
  role_id UUID NOT NULL REFERENCES core.tech_roles(id),
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_visit_assignments_visit ON operations.visit_assignments(visit_id);
CREATE INDEX idx_visit_assignments_tech ON operations.visit_assignments(tech_id);

-- Option A split: checklist_template_items → 3 child tables
CREATE TABLE operations.checklist_template_materials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  template_id UUID NOT NULL REFERENCES inventory.checklist_templates(id),
  material_id UUID NOT NULL REFERENCES inventory.materials(id),
  default_quantity INTEGER DEFAULT 1,
  required BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_ctm_template ON operations.checklist_template_materials(template_id);

CREATE TABLE operations.checklist_template_tools (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  template_id UUID NOT NULL REFERENCES inventory.checklist_templates(id),
  tool_id UUID NOT NULL REFERENCES inventory.tools(id),
  default_quantity INTEGER DEFAULT 1,
  required BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_ctt_template ON operations.checklist_template_tools(template_id);

CREATE TABLE operations.checklist_template_epps (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  template_id UUID NOT NULL REFERENCES inventory.checklist_templates(id),
  epp_id UUID NOT NULL REFERENCES inventory.epp_items(id),
  default_quantity INTEGER DEFAULT 1,
  required BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_cte_template ON operations.checklist_template_epps(template_id);

-- Option A split: visit_checklist → 3 child tables
CREATE TABLE operations.visit_checklist_materials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  material_id UUID NOT NULL REFERENCES inventory.materials(id),
  planned_quantity INTEGER DEFAULT 1,
  confirmed BOOLEAN DEFAULT false,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_vcm_visit ON operations.visit_checklist_materials(visit_id);

CREATE TABLE operations.visit_checklist_tools (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  tool_id UUID NOT NULL REFERENCES inventory.tools(id),
  planned_quantity INTEGER DEFAULT 1,
  confirmed BOOLEAN DEFAULT false,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_vct_visit ON operations.visit_checklist_tools(visit_id);

CREATE TABLE operations.visit_checklist_epps (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  epp_id UUID NOT NULL REFERENCES inventory.epp_items(id),
  planned_quantity INTEGER DEFAULT 1,
  confirmed BOOLEAN DEFAULT false,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_vce_visit ON operations.visit_checklist_epps(visit_id);

CREATE TABLE operations.visit_material_usages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  material_id UUID NOT NULL REFERENCES inventory.materials(id),
  quantity NUMERIC NOT NULL,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_vmu_visit ON operations.visit_material_usages(visit_id);

CREATE TABLE operations.visit_tool_usages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  tool_id UUID NOT NULL REFERENCES inventory.tools(id),
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_vtu_visit ON operations.visit_tool_usages(visit_id);

CREATE TABLE operations.visit_epp_usages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  epp_id UUID NOT NULL REFERENCES inventory.epp_items(id),
  quantity INTEGER DEFAULT 1,
  status shared.epp_usage_status NOT NULL,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_veu_visit ON operations.visit_epp_usages(visit_id);

CREATE TABLE operations.visit_measurements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  type UUID NOT NULL,  -- FK added by industry pack
  value NUMERIC,
  unit TEXT,
  result shared.measurement_result,
  measuring_device TEXT,
  notes TEXT,
  photo_url TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_visit_measurements_visit ON operations.visit_measurements(visit_id);
CREATE INDEX idx_visit_measurements_type ON operations.visit_measurements(type);

CREATE TABLE operations.visit_checkpoints (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  type shared.checkpoint_type NOT NULL,
  geom GEOMETRY(Point, 4326),
  timestamp TIMESTAMPTZ NOT NULL DEFAULT now(),
  device_info TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_visit_checkpoints_visit ON operations.visit_checkpoints(visit_id);
CREATE INDEX idx_visit_checkpoints_geom ON operations.visit_checkpoints USING GIST(geom);

CREATE TABLE operations.visit_photos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  checkpoint_id UUID REFERENCES operations.visit_checkpoints(id),
  url TEXT NOT NULL,
  geom GEOMETRY(Point, 4326),
  timestamp TIMESTAMPTZ DEFAULT now(),
  stage shared.photo_stage,
  finding_type UUID,  -- FK added by industry pack
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_visit_photos_visit ON operations.visit_photos(visit_id);
CREATE INDEX idx_visit_photos_checkpoint ON operations.visit_photos(checkpoint_id);
CREATE INDEX idx_visit_photos_geom ON operations.visit_photos USING GIST(geom);

CREATE TABLE operations.visit_reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  report_template_id UUID REFERENCES shared.report_templates(id),
  source shared.report_source NOT NULL,
  recorded_at TIMESTAMPTZ DEFAULT now(),
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_visit_reports_visit ON operations.visit_reports(visit_id);

CREATE TABLE operations.report_images (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  report_id UUID NOT NULL REFERENCES operations.visit_reports(id),
  url TEXT NOT NULL,
  pages INTEGER DEFAULT 1,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_report_images_report ON operations.report_images(report_id);

CREATE TABLE operations.report_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  report_id UUID NOT NULL REFERENCES operations.visit_reports(id),
  data_json JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_report_entries_report ON operations.report_entries(report_id);

-- Vehicle assignments (single table, D3)
CREATE TABLE operations.vehicle_assignments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  vehicle_id UUID NOT NULL REFERENCES inventory.vehicles(id),
  daily_plan_id UUID REFERENCES planning.daily_plans(id),
  visit_id UUID REFERENCES operations.visits(id),
  departure_time TIMESTAMPTZ,
  return_time TIMESTAMPTZ,
  departure_mileage INTEGER,
  return_mileage INTEGER,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_vehicle_assignments_vehicle ON operations.vehicle_assignments(vehicle_id);
CREATE INDEX idx_vehicle_assignments_daily_plan ON operations.vehicle_assignments(daily_plan_id);
CREATE INDEX idx_vehicle_assignments_visit ON operations.vehicle_assignments(visit_id);

CREATE TABLE operations.visit_rentals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  rental_id UUID NOT NULL REFERENCES inventory.rentals(id),
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ,
  hours NUMERIC,
  total_cost NUMERIC NOT NULL,
  reason TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_visit_rentals_visit ON operations.visit_rentals(visit_id);
CREATE INDEX idx_visit_rentals_rental ON operations.visit_rentals(rental_id);

CREATE TABLE operations.visit_sla_trackings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  visit_id UUID NOT NULL REFERENCES operations.visits(id),
  sla_id UUID NOT NULL REFERENCES partners.slas(id),
  requested_at TIMESTAMPTZ,
  responded_at TIMESTAMPTZ,
  resolved_at TIMESTAMPTZ,
  response_time_hours NUMERIC,
  resolution_time_hours NUMERIC,
  meets_response_sla BOOLEAN,
  meets_resolution_sla BOOLEAN,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_visit_sla_trackings_visit ON operations.visit_sla_trackings(visit_id);
CREATE INDEX idx_visit_sla_trackings_sla ON operations.visit_sla_trackings(sla_id);

-- ============================================================================
-- NOTIFICATIONS
-- ============================================================================

CREATE TABLE notifications.notification_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type shared.notification_type NOT NULL,
  channel shared.notification_channel NOT NULL,
  subject TEXT,
  body TEXT,
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE notifications.notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type shared.notification_type NOT NULL,
  channel shared.notification_channel NOT NULL,
  recipient TEXT NOT NULL,
  subject TEXT,
  body TEXT,
  status shared.notification_status DEFAULT 'pending',
  entity_type TEXT,
  entity_id UUID,
  scheduled_for TIMESTAMPTZ,
  sent_at TIMESTAMPTZ,
  read_at TIMESTAMPTZ,
  attempts INTEGER DEFAULT 0,
  error TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_notifications_status ON notifications.notifications(status);
CREATE INDEX idx_notifications_entity ON notifications.notifications(entity_type, entity_id);

-- ============================================================================
-- TRIGGERS — updated_at
-- ============================================================================

CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON core.users
  FOR EACH ROW EXECUTE FUNCTION touch_updated_at();

CREATE TRIGGER trg_partners_updated_at BEFORE UPDATE ON partners.partners
  FOR EACH ROW EXECUTE FUNCTION touch_updated_at();

CREATE TRIGGER trg_customers_updated_at BEFORE UPDATE ON customers.customers
  FOR EACH ROW EXECUTE FUNCTION touch_updated_at();

CREATE TRIGGER trg_visits_updated_at BEFORE UPDATE ON operations.visits
  FOR EACH ROW EXECUTE FUNCTION touch_updated_at();

-- ============================================================================
-- TRIGGERS — D2 (route ↔ daily_plan consistency)
-- ============================================================================

CREATE TRIGGER trg_validate_visit_route_plan BEFORE INSERT OR UPDATE OF route_id, daily_plan_id
  ON operations.visits FOR EACH ROW EXECUTE FUNCTION validate_visit_route_plan();

-- ============================================================================
-- TRIGGERS — audit
-- ============================================================================

CREATE TRIGGER trg_audit_daily_load_materials AFTER INSERT OR UPDATE OR DELETE ON planning.daily_load_materials
  FOR EACH ROW EXECUTE FUNCTION audit_change();
CREATE TRIGGER trg_audit_daily_load_tools AFTER INSERT OR UPDATE OR DELETE ON planning.daily_load_tools
  FOR EACH ROW EXECUTE FUNCTION audit_change();
CREATE TRIGGER trg_audit_daily_load_epps AFTER INSERT OR UPDATE OR DELETE ON planning.daily_load_epps
  FOR EACH ROW EXECUTE FUNCTION audit_change();
CREATE TRIGGER trg_audit_visit_checklist_materials AFTER INSERT OR UPDATE OR DELETE ON operations.visit_checklist_materials
  FOR EACH ROW EXECUTE FUNCTION audit_change();
CREATE TRIGGER trg_audit_visit_checklist_tools AFTER INSERT OR UPDATE OR DELETE ON operations.visit_checklist_tools
  FOR EACH ROW EXECUTE FUNCTION audit_change();
CREATE TRIGGER trg_audit_visit_checklist_epps AFTER INSERT OR UPDATE OR DELETE ON operations.visit_checklist_epps
  FOR EACH ROW EXECUTE FUNCTION audit_change();
CREATE TRIGGER trg_audit_visit_material_usages AFTER INSERT OR UPDATE OR DELETE ON operations.visit_material_usages
  FOR EACH ROW EXECUTE FUNCTION audit_change();
CREATE TRIGGER trg_audit_visit_tool_usages AFTER INSERT OR UPDATE OR DELETE ON operations.visit_tool_usages
  FOR EACH ROW EXECUTE FUNCTION audit_change();
CREATE TRIGGER trg_audit_visit_epp_usages AFTER INSERT OR UPDATE OR DELETE ON operations.visit_epp_usages
  FOR EACH ROW EXECUTE FUNCTION audit_change();

-- ============================================================================
-- GRANTS
-- ============================================================================

GRANT USAGE ON SCHEMA core, partners, customers, inventory, operations,
  planning, notifications, geocoding, shared TO fsm_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA core TO fsm_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA partners TO fsm_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA customers TO fsm_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA inventory TO fsm_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA operations TO fsm_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA planning TO fsm_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA notifications TO fsm_api;
GRANT SELECT ON ALL TABLES IN SCHEMA geocoding TO fsm_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA shared TO fsm_api;

GRANT fsm_api TO fsm_backend;

COMMIT;