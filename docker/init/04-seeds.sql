-- ============================================================================
-- 04-seeds.sql — Initial data (core only)
-- Industry-specific seeds go in industry-packs/{industry}/init/
-- ============================================================================

-- Global tech roles
INSERT INTO core.tech_roles (code, name) VALUES
  ('driver', 'Driver'),
  ('technician', 'Technician'),
  ('helper', 'Helper'),
  ('operative', 'Operative');

-- Global vehicle types (industry-agnostic; packs may add more)
INSERT INTO inventory.vehicle_types (code, name) VALUES
  ('truck', 'Truck'),
  ('crane', 'Crane'),
  ('van', 'Van'),
  ('crane_truck', 'Crane Truck'),
  ('other', 'Other');

-- Global route types (industry-agnostic; packs may add more)
INSERT INTO planning.route_types (code, name) VALUES
  ('maintenance', 'Maintenance'),
  ('installation', 'Installation'),
  ('repair', 'Repair'),
  ('inspection', 'Inspection'),
  ('delivery', 'Delivery'),
  ('collection', 'Collection'),
  ('emergency', 'Emergency'),
  ('other', 'Other');

-- Global visit types (industry-agnostic; packs may add more)
INSERT INTO operations.visit_types (code, name) VALUES
  ('pre_visit', 'Pre-visit'),
  ('installation', 'Installation'),
  ('maintenance', 'Maintenance'),
  ('repair', 'Repair'),
  ('diagnosis', 'Diagnosis'),
  ('inspection', 'Inspection'),
  ('emergency', 'Emergency'),
  ('certification', 'Certification'),
  ('other', 'Other');

-- Global photo findings (industry-agnostic; packs may add more)
INSERT INTO operations.photo_findings (code, name) VALUES
  ('normal', 'Normal'),
  ('damage', 'Damage'),
  ('incomplete', 'Incomplete'),
  ('other', 'Other');

-- Global pre-visit results (industry-agnostic; packs may add more)
INSERT INTO operations.pre_visit_results (code, name) VALUES
  ('approved', 'Approved'),
  ('rejected', 'Rejected'),
  ('conditional', 'Conditional');

-- Global partner service types (industry-agnostic; packs may add more)
INSERT INTO partners.partner_service_types (code, name) VALUES
  ('installation', 'Installation'),
  ('maintenance', 'Maintenance'),
  ('repair', 'Repair'),
  ('certification', 'Certification'),
  ('emergency', 'Emergency'),
  ('other', 'Other');

-- Global rejection reasons (industry-agnostic; packs may add more)
INSERT INTO operations.rejection_reasons (code, name, category) VALUES
  ('customer_unavailable', 'Customer not available on site', 'customer'),
  ('customer_cancelled', 'Customer cancelled the visit', 'customer'),
  ('missing_documentation', 'Missing documentation', 'operational'),
  ('other', 'Other reason (specify)', 'operational');