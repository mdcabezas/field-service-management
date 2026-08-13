-- ============================================================================
-- 05-comments.sql — COMMENT ON TABLE/COLUMN for OpenAPI (PostgREST)
-- ============================================================================

-- CORE
COMMENT ON SCHEMA core IS 'Core context: users, roles, audit';
COMMENT ON TABLE core.users IS 'System users; JWT sub claim = users.id';
COMMENT ON COLUMN core.users.id IS 'User id, equals JWT sub claim';
COMMENT ON COLUMN core.users.email IS 'Unique login email';
COMMENT ON COLUMN core.users.role IS 'Access role: admin, operator, technician';
COMMENT ON COLUMN core.users.password_hash IS 'Hashed password (NULL when LDAP-authenticated)';

COMMENT ON TABLE core.tech_roles IS 'Functional technician role (Driver, Technician, Helper, Operative)';
COMMENT ON TABLE core.plan_audit_log IS 'Audit trail for checklist/consumption changes';
COMMENT ON COLUMN core.plan_audit_log.entity_type IS 'Audited entity type (daily_load_material, visit_checklist_tool, etc.)';
COMMENT ON COLUMN core.plan_audit_log.old_data IS 'State before change (NULL on INSERT)';
COMMENT ON COLUMN core.plan_audit_log.new_data IS 'State after change (NULL on DELETE)';
COMMENT ON TABLE core.auth_audit_log IS 'Authentication/authorization audit trail';

-- PARTNERS
COMMENT ON SCHEMA partners IS 'Partners context: partner companies and agreements';
COMMENT ON TABLE partners.partners IS 'Partner company (distributor/retailer)';
COMMENT ON TABLE partners.partner_service_types IS 'Partner service type catalog (industry-agnostic)';
COMMENT ON TABLE partners.partner_contacts IS 'Administrative/commercial partner contacts';
COMMENT ON TABLE partners.partner_agreements IS 'Commercial agreement: services contracted by a partner';
COMMENT ON TABLE partners.partner_agreement_docs IS 'Required documentation by work type';
COMMENT ON TABLE partners.partner_agreement_forms IS 'Number of forms required by work type';
COMMENT ON TABLE partners.slas IS 'Service level agreements by partner';
COMMENT ON COLUMN partners.slas.response_hours IS 'Maximum hours to respond';
COMMENT ON COLUMN partners.slas.resolution_hours IS 'Maximum hours to resolve/complete';

-- CUSTOMERS
COMMENT ON SCHEMA customers IS 'Customers context: end clients, addresses, properties, technicians';
COMMENT ON TABLE customers.customers IS 'End client (natural or legal person)';
COMMENT ON TABLE customers.customer_addresses IS 'N:M relationship client ↔ address';
COMMENT ON TABLE customers.customer_address_partners IS 'Active partner(s) for an address';
COMMENT ON TABLE customers.properties IS 'Physical property (type defined by industry pack)';
COMMENT ON TABLE customers.technicians IS 'Field technicians';
COMMENT ON TABLE customers.tech_certifications IS 'Technician certifications (catalog provided by industry pack)';
COMMENT ON COLUMN customers.tech_certifications.cert_id IS 'FK to industry pack certifications catalog';

-- INVENTORY
COMMENT ON SCHEMA inventory IS 'Inventory context: materials, tools, PPE, vehicles, costs, maintenance';
COMMENT ON TABLE inventory.materials IS 'Materials and supplies';
COMMENT ON TABLE inventory.tools IS 'Tools';
COMMENT ON TABLE inventory.epp_items IS 'Personal protective equipment';
COMMENT ON TABLE inventory.equipment IS 'General equipment (compressor, generator, etc.)';
COMMENT ON TABLE inventory.vehicles IS 'Fleet vehicles';
COMMENT ON TABLE inventory.vehicle_types IS 'Vehicle type catalog (industry-agnostic)';
COMMENT ON TABLE inventory.rentals IS 'Rented equipment/vehicles';
COMMENT ON TABLE inventory.checklist_templates IS 'Checklist template by work type';
COMMENT ON TABLE inventory.cost_rates IS 'Cost rates (labor, vehicle, depreciation, PPE, overhead)';
COMMENT ON COLUMN inventory.cost_rates.reference_id IS 'Referenced element ID';
COMMENT ON COLUMN inventory.cost_rates.reference_type IS 'Referenced element type';
COMMENT ON TABLE inventory.maintenance_records IS 'Executed maintenance record';
COMMENT ON TABLE inventory.maintenance_schedules IS 'Preventive maintenance plan';

-- OPERATIONS
COMMENT ON SCHEMA operations IS 'Operations context: visits, assignments, checklists, measurements, checkpoints, photos, reports';
COMMENT ON TABLE operations.visit_types IS 'Visit type catalog (industry-agnostic)';
COMMENT ON TABLE operations.photo_findings IS 'Photo finding type catalog (industry-agnostic)';
COMMENT ON TABLE operations.pre_visit_results IS 'Pre-visit result catalog (industry-agnostic)';
COMMENT ON TABLE operations.rejection_reasons IS 'Rejection reason catalog (industry-agnostic)';
COMMENT ON TABLE operations.visits IS 'Main system entity: a technical visit';
COMMENT ON COLUMN operations.visits.parent_visit_id IS 'Parent visit (if pre_visit or sub-visit)';
COMMENT ON COLUMN operations.visits.partner_order_id IS 'Partner order ID';
COMMENT ON COLUMN operations.visits.billing_to IS 'Who to bill: partner or end customer';

COMMENT ON TABLE operations.visit_assignments IS 'Technician assignment to a visit';
COMMENT ON TABLE operations.checklist_template_materials IS 'Required materials in template';
COMMENT ON TABLE operations.checklist_template_tools IS 'Required tools in template';
COMMENT ON TABLE operations.checklist_template_epps IS 'Required PPE in template';
COMMENT ON TABLE operations.visit_checklist_materials IS 'Planned and confirmed materials for a visit';
COMMENT ON TABLE operations.visit_checklist_tools IS 'Planned and confirmed tools';
COMMENT ON TABLE operations.visit_checklist_epps IS 'Planned and confirmed PPE';
COMMENT ON TABLE operations.visit_material_usages IS 'Actual material consumption';
COMMENT ON TABLE operations.visit_tool_usages IS 'Actual tool usage';
COMMENT ON TABLE operations.visit_epp_usages IS 'Actual PPE usage';
COMMENT ON TABLE operations.visit_measurements IS 'Technical measurements (tightness, pressure, CO, draft, leak)';
COMMENT ON TABLE operations.visit_checkpoints IS 'Mandatory GPS checkpoint (arrival, departure, report_sent)';
COMMENT ON TABLE operations.visit_photos IS 'Georeferenced photos by stage';
COMMENT ON TABLE operations.visit_reports IS 'Visit report (digital or paper)';
COMMENT ON TABLE operations.report_images IS 'Paper report images';
COMMENT ON TABLE operations.report_entries IS 'Structured JSON content of digital report';
COMMENT ON TABLE operations.vehicle_assignments IS 'Vehicle assignment to daily plan or visit';
COMMENT ON TABLE operations.visit_rentals IS 'Rented equipment in a visit';
COMMENT ON TABLE operations.visit_sla_trackings IS 'SLA compliance tracking';

-- PLANNING
COMMENT ON SCHEMA planning IS 'Planning context: daily plan, technician assignment, material loading, routes';
COMMENT ON TABLE planning.daily_plans IS 'Consolidated daily work plan';
COMMENT ON COLUMN planning.daily_plans.name IS 'Human-readable name of the daily plan';
COMMENT ON TABLE planning.daily_plan_assignments IS 'Technician assignment to daily plan';
COMMENT ON TABLE planning.daily_load_materials IS 'Materials loaded/returned for the day';
COMMENT ON TABLE planning.daily_load_tools IS 'Tools loaded/returned';
COMMENT ON TABLE planning.daily_load_epps IS 'PPE loaded/returned';
COMMENT ON TABLE planning.routes IS 'Route or visit container';

-- NOTIFICATIONS
COMMENT ON SCHEMA notifications IS 'Notifications context: templates and notification instances';
COMMENT ON TABLE notifications.notification_templates IS 'Configurable notification template';
COMMENT ON TABLE notifications.notifications IS 'Sent or pending notification instance';
COMMENT ON COLUMN notifications.notifications.entity_type IS 'Associated entity type';
COMMENT ON COLUMN notifications.notifications.entity_id IS 'Associated entity ID (polymorphic)';

-- GEOCODING
COMMENT ON SCHEMA geocoding IS 'Geocoding context: global georeferenced address catalog';
COMMENT ON TABLE geocoding.addresses IS 'Physical address with PostGIS coordinates';

-- SHARED
COMMENT ON SCHEMA shared IS 'Shared global catalogs (rejection_category, report_templates)';
COMMENT ON TABLE shared.report_templates IS 'JSON field template for reports';