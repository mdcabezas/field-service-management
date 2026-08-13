-- ============================================================================
-- 02-domain-gas-seeds.sql — Gas Chile industry pack: lookup table seeds
-- Industry-agnostic catalogs are seeded in core (04-seeds.sql). This file adds
-- gas-specific codes into the core catalogs; ON CONFLICT keeps re-runs safe.
-- ============================================================================

-- Visit types (core operations.visit_types)
INSERT INTO operations.visit_types (code, name) VALUES
  ('pre_visit', 'Pre-visit'),
  ('installation', 'Installation'),
  ('maintenance', 'Maintenance'),
  ('emergency', 'Emergency'),
  ('certification', 'Certification'),
  ('repair', 'Repair'),
  ('diagnosis', 'Diagnosis')
ON CONFLICT (code) DO NOTHING;

-- Measurement types
INSERT INTO domain_gas.measurement_types (code, name, default_unit) VALUES
  ('tightness', 'Tightness', 'm3/h'),
  ('pressure', 'Pressure', 'mbar'),
  ('co', 'Carbon Monoxide', 'ppm'),
  ('draft', 'Draft', 'mmH2O'),
  ('leak', 'Leak', 'l/h'),
  ('ph', 'pH', 'pH'),
  ('temperature', 'Temperature', '°C'),
  ('other', 'Other', NULL);

-- Photo findings (core operations.photo_findings)
INSERT INTO operations.photo_findings (code, name) VALUES
  ('normal', 'Normal'),
  ('leak', 'Leak'),
  ('damage', 'Damage'),
  ('emergency', 'Emergency'),
  ('incomplete', 'Incomplete'),
  ('other', 'Other')
ON CONFLICT (code) DO NOTHING;

-- Vehicle types (core-agnostic catalog, also seeded in 04-seeds.sql core)

-- Partner service types (core partners.partner_service_types)
INSERT INTO partners.partner_service_types (code, name) VALUES
  ('installation', 'Installation'),
  ('maintenance', 'Maintenance'),
  ('certification', 'Certification'),
  ('emergency', 'Emergency'),
  ('other', 'Other')
ON CONFLICT (code) DO NOTHING;

-- Pre-visit results (core operations.pre_visit_results)
INSERT INTO operations.pre_visit_results (code, name) VALUES
  ('approved', 'Approved'),
  ('rejected', 'Rejected'),
  ('conditional', 'Conditional')
ON CONFLICT (code) DO NOTHING;

-- Route types (gas-specific codes into the core catalog; ON CONFLICT for re-run safety)
INSERT INTO planning.route_types (code, name) VALUES
  ('meter_reading', 'Meter Reading'),
  ('letter_delivery', 'Letter Delivery'),
  ('tank_collection', 'Tank Collection'),
  ('tank_delivery', 'Tank Delivery'),
  ('mass_inspection', 'Mass Inspection'),
  ('mixed', 'Mixed')
ON CONFLICT (code) DO NOTHING;

-- Property types
INSERT INTO domain_gas.property_types (code, name, columns_config) VALUES
  ('tank', 'Tank', '{"extra_columns": ["capacity_liters", "brand", "model"]}'),
  ('meter', 'Meter', '{"extra_columns": ["serial_number", "brand", "class"]}'),
  ('indoor_piping', 'Indoor Piping', '{"extra_columns": ["material", "diameter"]}'),
  ('appliance', 'Appliance', '{"extra_columns": ["brand", "model", "gas_type"]}'),
  ('other', 'Other', '{}');

-- SEC Chile certifications
INSERT INTO domain_gas.certifications (code, name, class, description) VALUES
  ('CL1', 'Class 1 Installer', 'Class 1', 'SEC authorized installer for all types of gas installations'),
  ('CL2', 'Class 2 Installer', 'Class 2', 'SEC authorized installer for gas installations up to 10 m³/h'),
  ('CL3', 'Class 3 Installer', 'Class 3', 'SEC authorized installer for gas installations up to 5 m³/h'),
  ('GLP', 'Bulk LPG Installer', 'LPG', 'Specialized installer for bulk LPG systems'),
  ('TC1', 'TC1 Declaration', 'TC', 'Declaration for individual installations'),
  ('TC2', 'TC2 Declaration', 'TC', 'Declaration for collective/building installations'),
  ('TC6', 'TC6 Declaration', 'TC', 'Declaration for industrial/commercial installations'),
  ('GREEN SEAL', 'Green Seal Inspector', 'Seal', 'Inspector authorized to issue SEC Green Seal certifications'),
  ('DS66', 'DS-66 Interior Installations', 'Standard', 'Certification under Supreme Decree 66');

-- Rejection reasons (core operations.rejection_reasons)
INSERT INTO operations.rejection_reasons (code, name, category) VALUES
  ('no_crane_access', 'No crane vehicle access', 'logistics'),
  ('insufficient_space', 'Insufficient space for tank', 'logistics'),
  ('impassable_road', 'Impassable road for truck', 'logistics'),
  ('damaged_vehicle', 'Damaged vehicle/truck', 'logistics'),
  ('no_materials', 'No materials available', 'logistics'),
  ('sec_noncompliance', 'Does not meet current SEC regulations', 'regulatory'),
  ('missing_permit', 'Missing municipal permit', 'regulatory'),
  ('minimum_distance', 'Does not meet minimum distance', 'regulatory'),
  ('incorrect_ventilation', 'Incorrect ventilation system', 'regulatory'),
  ('customer_unavailable', 'Customer not available on site', 'customer'),
  ('customer_cancelled', 'Customer cancelled the visit', 'customer'),
  ('customer_rejected_quote', 'Customer rejected quote', 'customer'),
  ('customer_no_authorization', 'Customer does not authorize work', 'customer'),
  ('pre_visit_rejected', 'Pre-visit not approved', 'operational'),
  ('missing_documentation', 'Missing partner documentation', 'operational'),
  ('partner_change', 'Provider change by customer', 'operational'),
  ('other', 'Other reason (specify)', 'operational')
ON CONFLICT (code) DO NOTHING;