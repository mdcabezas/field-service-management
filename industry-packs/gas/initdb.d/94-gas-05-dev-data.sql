-- ============================================================================
-- 94-gas-05-dev-data.sql — Gas Chile industry pack: dev data enrichment
-- Reuses rows seeded by the generic dev seed (07-seed-data-dev.sql).
-- Fills industry-specific type/classification columns via UPDATE by id.
-- Run: docker compose exec postgres psql -U fsm_admin -d fsm_fsm -f /docker-entrypoint-initdb.d/94-gas-05-dev-data.sql
-- ============================================================================

-- Partner agreements: service type (core partners.partner_service_types)
UPDATE partners.partner_agreements SET service_type = (SELECT id FROM partners.partner_service_types WHERE code = 'installation')
  WHERE id = '32000000-0000-0000-0000-000000000001';
UPDATE partners.partner_agreements SET service_type = (SELECT id FROM partners.partner_service_types WHERE code = 'maintenance')
  WHERE id = '32000000-0000-0000-0000-000000000002';
UPDATE partners.partner_agreements SET service_type = (SELECT id FROM partners.partner_service_types WHERE code = 'certification')
  WHERE id = '32000000-0000-0000-0000-000000000003';

-- Properties: type
UPDATE customers.properties SET type = (SELECT id FROM domain_gas.property_types WHERE code = 'other')
  WHERE id = '43000000-0000-0000-0000-000000000001';
UPDATE customers.properties SET type = (SELECT id FROM domain_gas.property_types WHERE code = 'other')
  WHERE id = '43000000-0000-0000-0000-000000000002';
UPDATE customers.properties SET type = (SELECT id FROM domain_gas.property_types WHERE code = 'appliance')
  WHERE id = '43000000-0000-0000-0000-000000000003';

-- Vehicles: type (catalog in core inventory.vehicle_types)
UPDATE inventory.vehicles SET type = (SELECT id FROM inventory.vehicle_types WHERE code = 'truck')
  WHERE id = '53000000-0000-0000-0000-000000000001';
UPDATE inventory.vehicles SET type = (SELECT id FROM inventory.vehicle_types WHERE code = 'van')
  WHERE id = '53000000-0000-0000-0000-000000000002';
UPDATE inventory.vehicles SET type = (SELECT id FROM inventory.vehicle_types WHERE code = 'truck')
  WHERE id = '53000000-0000-0000-0000-000000000003';

-- Checklist templates: work type (core operations.visit_types)
UPDATE inventory.checklist_templates SET work_type = (SELECT id FROM operations.visit_types WHERE code = 'installation')
  WHERE id = '56000000-0000-0000-0000-000000000001';
UPDATE inventory.checklist_templates SET work_type = (SELECT id FROM operations.visit_types WHERE code = 'maintenance')
  WHERE id = '56000000-0000-0000-0000-000000000002';
UPDATE inventory.checklist_templates SET work_type = (SELECT id FROM operations.visit_types WHERE code = 'certification')
  WHERE id = '56000000-0000-0000-0000-000000000003';

-- Technician certifications: SEC Chile certs
UPDATE customers.tech_certifications SET cert_id = (SELECT id FROM domain_gas.certifications WHERE code = 'CL2'),
  number = 'SEC-2024-001', issuer = 'SEC Chile'
  WHERE id = '61000000-0000-0000-0000-000000000001';
UPDATE customers.tech_certifications SET cert_id = (SELECT id FROM domain_gas.certifications WHERE code = 'CL3'),
  number = 'SEC-2024-002', issuer = 'SEC Chile'
  WHERE id = '61000000-0000-0000-0000-000000000002';

-- Routes: type
UPDATE planning.routes SET type = (SELECT id FROM planning.route_types WHERE code = 'meter_reading')
  WHERE id = '71000000-0000-0000-0000-000000000001';
UPDATE planning.routes SET type = (SELECT id FROM planning.route_types WHERE code = 'tank_delivery')
  WHERE id = '71000000-0000-0000-0000-000000000002';

-- Visits: type (core operations.visit_types)
UPDATE operations.visits SET type = (SELECT id FROM operations.visit_types WHERE code = 'maintenance')
  WHERE id = '80000000-0000-0000-0000-000000000001';
UPDATE operations.visits SET type = (SELECT id FROM operations.visit_types WHERE code = 'diagnosis')
  WHERE id = '80000000-0000-0000-0000-000000000002';
UPDATE operations.visits SET type = (SELECT id FROM operations.visit_types WHERE code = 'installation')
  WHERE id = '80000000-0000-0000-0000-000000000003';
UPDATE operations.visits SET type = (SELECT id FROM operations.visit_types WHERE code = 'maintenance')
  WHERE id = '80000000-0000-0000-0000-000000000004';

-- Visit measurements: type
UPDATE operations.visit_measurements SET type = (SELECT id FROM domain_gas.measurement_types WHERE code = 'pressure'),
  measuring_device = 'Manómetro digital', notes = 'Presión dentro de rango'
  WHERE id = '89000000-0000-0000-0000-000000000001';
UPDATE operations.visit_measurements SET type = (SELECT id FROM domain_gas.measurement_types WHERE code = 'co'),
  measuring_device = 'Detector de gases', notes = 'Sin fugas detectadas'
  WHERE id = '89000000-0000-0000-0000-000000000002';

-- ============================================================================
-- END GAS DEV DATA
-- ============================================================================
