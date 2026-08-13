-- =============================================================================
-- DEV SEED DATA
-- Generic FSM core. No domain-specific data.
-- Run: docker compose exec postgres psql -U fsm_admin -d fsm_fsm -f /docker-entrypoint-initdb.d/07-seed-data-dev.sql
-- =============================================================================

-- ─────────────────────────────────────────────────────────────────────────────
-- SHARED
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO shared.report_templates (id, name, fields_json) VALUES
  ('10000000-0000-0000-0000-000000000001', 'Checklist General', '{"sections":[{"name":"General","fields":["resultado","observaciones"]}]}'),
  ('10000000-0000-0000-0000-000000000002', 'Reporte de Trabajo', '{"sections":[{"name":"Trabajo","fields":["descripcion","horas","materiales"]}]}'),
  ('10000000-0000-0000-0000-000000000003', 'Inspección de Seguridad', '{"sections":[{"name":"Seguridad","fields":["hallazgos","riesgos","acciones"]}]}')
ON CONFLICT (id) DO NOTHING;

INSERT INTO core.users (id, email, role, name) VALUES
  ('10000000-0000-0000-0000-000000000001', 'admin@fsm.cl', 'technician', 'Admin User'),
  ('10000000-0000-0000-0000-000000000002', 'operador@fsm.cl', 'technician', 'Operador User'),
  ('10000000-0000-0000-0000-000000000003', 'carlos@fsm.cl', 'technician', 'Carlos Técnico'),
  ('10000000-0000-0000-0000-000000000004', 'ana@fsm.cl', 'technician', 'Ana Técnica'),
  ('10000000-0000-0000-0000-000000000005', 'pedro@fsm.cl', 'technician', 'Pedro Operario')
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- GEOCODING
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO geocoding.addresses (id, street, number, neighborhood, city, region, postal_code, geom) VALUES
  ('20000000-0000-0000-0000-000000000001', 'Av. Providencia', '1234', 'Providencia', 'Santiago', 'Región Metropolitana', '7500000', ST_SetSRID(ST_MakePoint(-70.6101, -33.4246), 4326)),
  ('20000000-0000-0000-0000-000000000002', 'Calle Los Leones', '567', 'Las Condes', 'Santiago', 'Región Metropolitana', '7550000', ST_SetSRID(ST_MakePoint(-70.5948, -33.4106), 4326)),
  ('20000000-0000-0000-0000-000000000003', 'Av. Recoleta', '890', 'Recoleta', 'Santiago', 'Región Metropolitana', '7510000', ST_SetSRID(ST_MakePoint(-70.6460, -33.4116), 4326)),
  ('20000000-0000-0000-0000-000000000004', 'Calle Maipú', '456', 'Maipú', 'Santiago', 'Región Metropolitana', '9100000', ST_SetSRID(ST_MakePoint(-70.7583, -33.5117), 4326)),
  ('20000000-0000-0000-0000-000000000005', 'Av. El Bosque', '101', 'El Bosque', 'Santiago', 'Región Metropolitana', '8130000', ST_SetSRID(ST_MakePoint(-70.6693, -33.5653), 4326))
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- PARTNERS
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO partners.partners (id, name, tax_id, status) VALUES
  ('30000000-0000-0000-0000-000000000001', 'Servicios Técnicos SpA', '76.123.456-7', 'active'),
  ('30000000-0000-0000-0000-000000000002', 'Mantenimiento Sur Ltda', '76.234.567-8', 'active'),
  ('30000000-0000-0000-0000-000000000003', 'Instalaciones Express SpA', '76.345.678-9', 'active')
ON CONFLICT (id) DO NOTHING;

INSERT INTO partners.partner_contacts (id, partner_id, name, position, phone, email) VALUES
  ('31000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'Juan Pérez', 'Gerente', '+56912345678', 'juan.perez@servtec.cl'),
  ('31000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', 'María González', 'Supervisora', '+56923456789', 'maria.gonzalez@mantsur.cl'),
  ('31000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', 'Carlos Rodríguez', 'Jefe Operaciones', '+56934567890', 'carlos.rodriguez@instexp.cl')
ON CONFLICT (id) DO NOTHING;

INSERT INTO partners.partner_agreements (id, partner_id, service_type, rate, active) VALUES
  ('32000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', NULL, 45000, true),
  ('32000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', NULL, 52000, true),
  ('32000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', NULL, 38000, true)
ON CONFLICT (id) DO NOTHING;

INSERT INTO partners.partner_agreement_docs (id, agreement_id, work_type, doc_name, required) VALUES
  ('33000000-0000-0000-0000-000000000001', '32000000-0000-0000-0000-000000000001', 'Instalación', 'Certificado de habilitación', true),
  ('33000000-0000-0000-0000-000000000002', '32000000-0000-0000-0000-000000000002', 'Mantenimiento', 'Póliza de Seguridad', true)
ON CONFLICT (id) DO NOTHING;

INSERT INTO partners.partner_agreement_forms (id, agreement_id, work_type, form_template_id, quantity) VALUES
  ('34000000-0000-0000-0000-000000000001', '32000000-0000-0000-0000-000000000001', 'Instalación', '10000000-0000-0000-0000-000000000001', 1),
  ('34000000-0000-0000-0000-000000000002', '32000000-0000-0000-0000-000000000002', 'Mantenimiento', '10000000-0000-0000-0000-000000000002', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO partners.slas (id, partner_id, name, description, response_hours, resolution_hours, compliance_target, active, valid_from, valid_until) VALUES
  ('35000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'SLA Prioridad Normal', 'Respuesta en 24h, resolución en 72h', 24, 72, 0.95, true, '2026-01-01', '2026-12-31'),
  ('35000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', 'SLA Urgente', 'Respuesta en 4h, resolución en 24h', 4, 24, 0.98, true, '2026-01-01', '2026-12-31')
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- CUSTOMERS
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO customers.customers (id, name, phone, email, tax_id) VALUES
  ('40000000-0000-0000-0000-000000000001', 'Empresa constructora Los Andes', '+56911111111', 'contacto@losandes.cl', '76.500.001-1'),
  ('40000000-0000-0000-0000-000000000002', 'Comercializadora Central SpA', '+56922222222', 'ventas@central.cl', '76.500.002-2'),
  ('40000000-0000-0000-0000-000000000003', 'Industria Alimentaria Norte', '+56933333333', 'ops@alimnorte.cl', '76.500.003-3'),
  ('40000000-0000-0000-0000-000000000004', 'Minera Sur SpA', '+56944444444', 'admin@minsur.cl', '76.500.004-4'),
  ('40000000-0000-0000-0000-000000000005', 'Retail Express Ltda', '+56955555555', 'tiendas@relexpress.cl', '76.500.005-5')
ON CONFLICT (id) DO NOTHING;

INSERT INTO customers.customer_addresses (id, customer_id, address_id, name, type) VALUES
  ('41000000-0000-0000-0000-000000000001', '40000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 'Sucursal Principal', 'commercial'),
  ('41000000-0000-0000-0000-000000000002', '40000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000002', 'Bodega Las Condes', 'commercial'),
  ('41000000-0000-0000-0000-000000000003', '40000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000003', 'Oficina Central', 'commercial'),
  ('41000000-0000-0000-0000-000000000004', '40000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000004', 'Planta Maipú', 'industrial'),
  ('41000000-0000-0000-0000-000000000005', '40000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000005', 'Tienda El Bosque', 'commercial')
ON CONFLICT (id) DO NOTHING;

INSERT INTO customers.customer_address_partners (id, customer_address_id, partner_id, from_date, to_date) VALUES
  ('42000000-0000-0000-0000-000000000001', '41000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', '2026-01-01', NULL),
  ('42000000-0000-0000-0000-000000000002', '41000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000002', '2026-03-01', NULL),
  ('42000000-0000-0000-0000-000000000003', '41000000-0000-0000-0000-000000000005', '30000000-0000-0000-0000-000000000003', '2026-06-01', NULL)
ON CONFLICT (id) DO NOTHING;

INSERT INTO customers.properties (id, customer_address_id, name, type, notes) VALUES
  ('43000000-0000-0000-0000-000000000001', '41000000-0000-0000-0000-000000000001', 'Equipo zona 3', NULL, 'Mantenimiento preventivo trimestral'),
  ('43000000-0000-0000-0000-000000000002', '41000000-0000-0000-0000-000000000003', 'Sistema eléctrico principal', NULL, 'Inspección anual'),
  ('43000000-0000-0000-0000-000000000003', '41000000-0000-0000-0000-000000000005', 'Refrigeración retail', NULL, 'Monitoreo constante')
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- INVENTORY
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO inventory.materials (id, sku, name, unit, unit_cost) VALUES
  ('50000000-0000-0000-0000-000000000001', 'MAT-001', 'Cable eléctrico 2.5mm', 'metro', 850),
  ('50000000-0000-0000-0000-000000000002', 'MAT-002', 'Tubo PVC 1/2"', 'metro', 1200),
  ('50000000-0000-0000-0000-000000000003', 'MAT-003', 'Cinta aislante', 'unidad', 650),
  ('50000000-0000-0000-0000-000000000004', 'MAT-004', 'Conector RJ45', 'unidad', 350),
  ('50000000-0000-0000-0000-000000000005', 'MAT-005', 'Fusible 20A', 'unidad', 280),
  ('50000000-0000-0000-0000-000000000006', 'MAT-006', 'Soldadura 60/40', 'metro', 420),
  ('50000000-0000-0000-0000-000000000007', 'MAT-007', 'Aceite lubricante', 'litro', 3500),
  ('50000000-0000-0000-0000-000000000008', 'MAT-008', 'Filtro de repuesto', 'unidad', 8900)
ON CONFLICT (id) DO NOTHING;

INSERT INTO inventory.tools (id, code, name, status) VALUES
  ('51000000-0000-0000-0000-000000000001', 'HER-001', 'Taladro percutor Bosch', 'available'),
  ('51000000-0000-0000-0000-000000000002', 'HER-002', 'Multímetro Fluke 87V', 'available'),
  ('51000000-0000-0000-0000-000000000003', 'HER-003', 'Llave dinamometrica 1/2"', 'in_use'),
  ('51000000-0000-0000-0000-000000000004', 'HER-004', 'Osciloscopio digital', 'available'),
  ('51000000-0000-0000-0000-000000000005', 'HER-005', 'Juego de llaves allen', 'available'),
  ('51000000-0000-0000-0000-000000000006', 'HER-006', 'Cortadora de acero', 'maintenance'),
  ('51000000-0000-0000-0000-000000000007', 'HER-007', 'Soldador 60W', 'available'),
  ('51000000-0000-0000-0000-000000000008', 'HER-008', 'Detector multiparámetro', 'available')
ON CONFLICT (id) DO NOTHING;

INSERT INTO inventory.epp_items (id, name, type, lifecycle) VALUES
  ('52000000-0000-0000-0000-000000000001', 'Casco MSA', 'head', 'reusable'),
  ('52000000-0000-0000-0000-000000000002', 'Guantes dieléctricos', 'hands', 'reusable'),
  ('52000000-0000-0000-0000-000000000003', 'Chaleco reflectante', 'body', 'reusable'),
  ('52000000-0000-0000-0000-000000000004', 'Lentes seguridad', 'eyes', 'disposable'),
  ('52000000-0000-0000-0000-000000000005', 'Tapones oídos', 'ears', 'disposable'),
  ('52000000-0000-0000-0000-000000000006', 'Botas seguridad', 'feet', 'reusable'),
  ('52000000-0000-0000-0000-000000000007', 'Mascarilla N95', 'respiratory', 'disposable')
ON CONFLICT (id) DO NOTHING;

INSERT INTO inventory.equipment (id, code, name, status) VALUES
  ('53500000-0000-0000-0000-000000000001', 'EQ-001', 'Generador eléctrico 5kW', 'available'),
  ('53500000-0000-0000-0000-000000000002', 'EQ-002', 'Compresor de taller 25L', 'available'),
  ('53500000-0000-0000-0000-000000000003', 'EQ-003', 'Plataforma elevadora', 'maintenance')
ON CONFLICT (id) DO NOTHING;

INSERT INTO inventory.vehicles (id, type, license_plate, name, brand, model, year, status, capacity) VALUES
  ('53000000-0000-0000-0000-000000000001', NULL, 'AB-1234', 'Camioneta 1', 'Toyota', 'Hilux', 2024, 'available', '1.5 ton'),
  ('53000000-0000-0000-0000-000000000002', NULL, 'CD-5678', 'Furgón 2', 'Hyundai', 'H1', 2023, 'in_use', '800 kg'),
  ('53000000-0000-0000-0000-000000000003', NULL, 'EF-9012', 'Camioneta 3', 'Nissan', 'NP300', 2025, 'available', '1.2 ton')
ON CONFLICT (id) DO NOTHING;

INSERT INTO inventory.rentals (id, type, supplier, item_description, daily_cost, hourly_cost, fixed_cost) VALUES
  ('54000000-0000-0000-0000-000000000001', 'vehicle', 'Rent-a-Car SpA', 'Camioneta 4x4', 35000, 0, 0),
  ('54000000-0000-0000-0000-000000000002', 'tool', 'Herramientas Pro', 'Andamio modular', 0, 2500, 15000),
  ('54000000-0000-0000-0000-000000000003', 'equipment', 'Maquinaria Ltda', 'Plataforma elevadora', 0, 8000, 45000)
ON CONFLICT (id) DO NOTHING;

INSERT INTO inventory.cost_rates (id, type, reference_id, reference_type, value, unit, valid_from, valid_until) VALUES
  ('55000000-0000-0000-0000-000000000001', 'labor', '51000000-0000-0000-0000-000000000001', 'tool', 25000, 'hour', '2026-01-01', '2026-12-31'),
  ('55000000-0000-0000-0000-000000000002', 'vehicle', '53000000-0000-0000-0000-000000000001', 'vehicle', 180, 'km', '2026-01-01', '2026-12-31'),
  ('55000000-0000-0000-0000-000000000003', 'overhead', NULL, NULL, 15000, 'visit', '2026-01-01', '2026-12-31')
ON CONFLICT (id) DO NOTHING;

INSERT INTO inventory.checklist_templates (id, name, work_type) VALUES
  ('56000000-0000-0000-0000-000000000001', 'Checklist Instalación Básica', NULL),
  ('56000000-0000-0000-0000-000000000002', 'Checklist Mantenimiento Preventivo', NULL),
  ('56000000-0000-0000-0000-000000000003', 'Checklist Inspección Seguridad', NULL)
ON CONFLICT (id) DO NOTHING;

-- Checklist template materials
INSERT INTO operations.checklist_template_materials (id, template_id, material_id, default_quantity) VALUES
  ('56100000-0000-0000-0000-000000000001', '56000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 5),
  ('56100000-0000-0000-0000-000000000002', '56000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000003', 2),
  ('56100000-0000-0000-0000-000000000003', '56000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000007', 1),
  ('56100000-0000-0000-0000-000000000004', '56000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000008', 1)
ON CONFLICT (id) DO NOTHING;

-- Checklist template tools
INSERT INTO operations.checklist_template_tools (id, template_id, tool_id) VALUES
  ('56200000-0000-0000-0000-000000000001', '56000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000001'),
  ('56200000-0000-0000-0000-000000000002', '56000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000005'),
  ('56200000-0000-0000-0000-000000000003', '56000000-0000-0000-0000-000000000002', '51000000-0000-0000-0000-000000000002'),
  ('56200000-0000-0000-0000-000000000004', '56000000-0000-0000-0000-000000000003', '51000000-0000-0000-0000-000000000008')
ON CONFLICT (id) DO NOTHING;

-- Checklist template EPPs
INSERT INTO operations.checklist_template_epps (id, template_id, epp_id) VALUES
  ('56300000-0000-0000-0000-000000000001', '56000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000001'),
  ('56300000-0000-0000-0000-000000000002', '56000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000002'),
  ('56300000-0000-0000-0000-000000000003', '56000000-0000-0000-0000-000000000003', '52000000-0000-0000-0000-000000000003'),
  ('56300000-0000-0000-0000-000000000004', '56000000-0000-0000-0000-000000000003', '52000000-0000-0000-0000-000000000004')
ON CONFLICT (id) DO NOTHING;

-- Inventory maintenance
INSERT INTO inventory.maintenance_schedules (id, type, reference_id, frequency_km, frequency_days, active) VALUES
  ('57000000-0000-0000-0000-000000000001', 'vehicle', '53000000-0000-0000-0000-000000000001', 10000, 180, true),
  ('57000000-0000-0000-0000-000000000002', 'tool', '51000000-0000-0000-0000-000000000002', 0, 365, true)
ON CONFLICT (id) DO NOTHING;

INSERT INTO inventory.maintenance_records (id, type, reference_id, date, cost, supplier, description, next_date) VALUES
  ('57100000-0000-0000-0000-000000000001', 'vehicle', '53000000-0000-0000-0000-000000000001', '2026-03-15', 120000, 'Taller Los Andes', 'Cambio de aceite y filtro', '2026-09-15'),
  ('57100000-0000-0000-0000-000000000002', 'tool', '51000000-0000-0000-0000-000000000002', '2026-01-10', 15000, 'Fluke Service', 'Calibración anual', '2027-01-10'),
  ('57100000-0000-0000-0000-000000000003', 'equipment', '53500000-0000-0000-0000-000000000001', '2026-06-20', 85000, 'Maquinaria Ltda', 'Servicio generador 5kW', '2026-12-20')
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- TECHNICIANS (uses core.users)
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO customers.technicians (id, user_id, name, is_active, specialties) VALUES
  ('60000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000003', 'Carlos Técnico', true, ARRAY['electrical','hvac']),
  ('60000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000004', 'Ana Técnica', true, ARRAY['plumbing','general']),
  ('60000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000005', 'Pedro Operario', true, ARRAY['general','painting'])
ON CONFLICT (id) DO NOTHING;

INSERT INTO customers.tech_certifications (id, tech_id, cert_id, number, issuer, issue_date, expiry_date) VALUES
  ('61000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', NULL, 'CERT-2024-001', 'Organismo certificador', '2024-06-01', '2027-06-01'),
  ('61000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000002', NULL, 'CERT-2024-002', 'Organismo certificador', '2024-08-15', '2027-08-15')
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- PLANNING
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO planning.daily_plans (id, name, date, notes) VALUES
  ('70000000-0000-0000-0000-000000000001', 'Plan diario de pruebas', CURRENT_DATE, 'Plan diario de pruebas'),
  ('70000000-0000-0000-0000-000000000002', 'Plan mañana', CURRENT_DATE + 1, 'Plan mañana')
ON CONFLICT (id) DO NOTHING;

INSERT INTO planning.routes (id, type, date, daily_plan_id, notes, status) VALUES
  ('71000000-0000-0000-0000-000000000001', NULL, CURRENT_DATE, '70000000-0000-0000-0000-000000000001', 'Ruta zona norte', 'scheduled'),
  ('71000000-0000-0000-0000-000000000002', NULL, CURRENT_DATE, '70000000-0000-0000-0000-000000000001', 'Ruta zona sur', 'scheduled')
ON CONFLICT (id) DO NOTHING;

INSERT INTO planning.daily_plan_assignments (id, daily_plan_id, tech_id, role_id) VALUES
  ('72000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', (SELECT id FROM core.tech_roles WHERE code = 'technician')),
  ('72000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000002', (SELECT id FROM core.tech_roles WHERE code = 'technician'))
ON CONFLICT (id) DO NOTHING;

INSERT INTO planning.daily_load_materials (id, daily_plan_id, material_id, loaded_quantity) VALUES
  ('73000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 20),
  ('73000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000003', 10)
ON CONFLICT (id) DO NOTHING;

INSERT INTO planning.daily_load_tools (id, daily_plan_id, tool_id, loaded_quantity) VALUES
  ('74000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000001', 1),
  ('74000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000002', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO planning.daily_load_epps (id, daily_plan_id, epp_id, loaded_quantity) VALUES
  ('75000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000001', 2),
  ('75000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000002', 2)
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- OPERATIONS (visits + related)
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO operations.visits (id, property_id, partner_id, route_id, daily_plan_id, type, status, priority, source, billing_to, scheduled_at, notes) VALUES
  ('80000000-0000-0000-0000-000000000001', '43000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', '71000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', NULL, 'scheduled', 'normal', 'portal', 'customer', CURRENT_DATE + INTERVAL '1 hour', 'Mantenimiento preventivo zona 3'),
  ('80000000-0000-0000-0000-000000000002', '43000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', '71000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', NULL, 'en_route', 'urgent', 'phone', 'partner', CURRENT_DATE + INTERVAL '3 hours', 'Diagnóstico sistema eléctrico'),
  ('80000000-0000-0000-0000-000000000003', '43000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', '71000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', NULL, 'completed', 'normal', 'email', 'customer', CURRENT_DATE - INTERVAL '2 hours', 'Instalación completada'),
  ('80000000-0000-0000-0000-000000000004', '43000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', NULL, NULL, NULL, 'cancelled', 'low', 'portal', 'customer', CURRENT_DATE + INTERVAL '5 hours', 'Visit cancelled - rescheduled')
ON CONFLICT (id) DO NOTHING;

INSERT INTO operations.visit_assignments (id, visit_id, tech_id, role_id) VALUES
  ('81000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', (SELECT id FROM core.tech_roles WHERE code = 'technician')),
  ('81000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000002', (SELECT id FROM core.tech_roles WHERE code = 'technician')),
  ('81000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000003', (SELECT id FROM core.tech_roles WHERE code = 'technician'))
ON CONFLICT (id) DO NOTHING;

-- Visit checklist materials
INSERT INTO operations.visit_checklist_materials (id, visit_id, material_id, planned_quantity, confirmed, notes) VALUES
  ('82000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 3, true, 'Usado'),
  ('82000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000003', 1, true, 'Usado')
ON CONFLICT (id) DO NOTHING;

-- Visit checklist tools
INSERT INTO operations.visit_checklist_tools (id, visit_id, tool_id, planned_quantity, confirmed, notes) VALUES
  ('83000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000001', 1, true, 'Devuelto ok'),
  ('83000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000002', 1, true, 'Devuelto ok')
ON CONFLICT (id) DO NOTHING;

-- Visit checklist EPPs
INSERT INTO operations.visit_checklist_epps (id, visit_id, epp_id, planned_quantity, confirmed, notes) VALUES
  ('84000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000001', 1, true, 'Usado'),
  ('84000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000002', 1, true, 'Usado')
ON CONFLICT (id) DO NOTHING;

-- Visit material usages
INSERT INTO operations.visit_material_usages (id, visit_id, material_id, quantity) VALUES
  ('85000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 3),
  ('85000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000003', 1)
ON CONFLICT (id) DO NOTHING;

-- Visit tool usages
INSERT INTO operations.visit_tool_usages (id, visit_id, tool_id) VALUES
  ('86000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000001'),
  ('86000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000002')
ON CONFLICT (id) DO NOTHING;

-- Visit EPP usages
INSERT INTO operations.visit_epp_usages (id, visit_id, epp_id, status) VALUES
  ('87000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000001', 'used'),
  ('87000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000004', 'used')
ON CONFLICT (id) DO NOTHING;

-- Visit checkpoints
INSERT INTO operations.visit_checkpoints (id, visit_id, type, geom, timestamp, device_info) VALUES
  ('88000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000003', 'arrival', ST_SetSRID(ST_MakePoint(-70.6101, -33.4246), 4326), CURRENT_TIMESTAMP - INTERVAL '2 hours', 'iPhone 15'),
  ('88000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000003', 'departure', ST_SetSRID(ST_MakePoint(-70.6101, -33.4246), 4326), CURRENT_TIMESTAMP - INTERVAL '45 minutes', 'iPhone 15')
ON CONFLICT (id) DO NOTHING;

-- Visit measurements
INSERT INTO operations.visit_measurements (id, visit_id, type, value, unit, result, measuring_device, notes) VALUES
  ('89000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000003', NULL, 4.5, 'bar', 'approved', 'Instrumento digital', 'Valor dentro de rango'),
  ('89000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000003', NULL, 0.02, 'ppm', 'approved', 'Detector multiparámetro', 'Sin anomalías detectadas')
ON CONFLICT (id) DO NOTHING;

-- Visit reports
INSERT INTO operations.visit_reports (id, visit_id, report_template_id, source, recorded_at) VALUES
  ('8A000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000002', 'system', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

-- Report entries
INSERT INTO operations.report_entries (id, report_id, data_json) VALUES
  ('8B000000-0000-0000-0000-000000000001', '8A000000-0000-0000-0000-000000000001', '{"section":"General","field":"resultado","value":"aprobado"}'),
  ('8B000000-0000-0000-0000-000000000002', '8A000000-0000-0000-0000-000000000001', '{"section":"General","field":"observaciones","value":"Trabajo realizado sin novedad"}')
ON CONFLICT (id) DO NOTHING;

-- Report images
INSERT INTO operations.report_images (id, report_id, url, pages) VALUES
  ('8C000000-0000-0000-0000-000000000001', '8A000000-0000-0000-0000-000000000001', 'https://storage.example.com/reports/8A000000-0001.jpg', 1),
  ('8C000000-0000-0000-0000-000000000002', '8A000000-0000-0000-0000-000000000001', 'https://storage.example.com/reports/8A000000-0002.jpg', 1)
ON CONFLICT (id) DO NOTHING;

-- Vehicle assignments
INSERT INTO operations.vehicle_assignments (id, vehicle_id, daily_plan_id, visit_id, departure_time, return_time, departure_mileage, return_mileage) VALUES
  ('8D000000-0000-0000-0000-000000000001', '53000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', CURRENT_TIMESTAMP - INTERVAL '3 hours', NULL, 45230, NULL)
ON CONFLICT (id) DO NOTHING;

-- Visit SLA trackings
INSERT INTO operations.visit_sla_trackings (id, visit_id, sla_id, meets_response_sla, meets_resolution_sla, response_time_hours, resolution_time_hours) VALUES
  ('8E000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000003', '35000000-0000-0000-0000-000000000001', true, true, 1.5, 3.2)
ON CONFLICT (id) DO NOTHING;

-- Visit rentals
INSERT INTO operations.visit_rentals (id, visit_id, rental_id, start_time, total_cost, reason) VALUES
  ('8F000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '54000000-0000-0000-0000-000000000001', CURRENT_TIMESTAMP, 35000, 'Traslado a sitio')
ON CONFLICT (id) DO NOTHING;

-- Notifications
INSERT INTO notifications.notifications (id, type, channel, status, recipient, subject, body, entity_type, entity_id) VALUES
  ('90000000-0000-0000-0000-000000000001', 'visit_assigned', 'push', 'read', 'carlos@fsm.cl', 'Visita asignada', 'Se te ha asignado una visita en Av. Providencia 1234', 'visit', '80000000-0000-0000-0000-000000000001'),
  ('90000000-0000-0000-0000-000000000002', 'visit_reminder', 'email', 'sent', 'ana@fsm.cl', 'Recordatorio de visita', 'Tienes una visita programada para hoy a las 15:00', 'visit', '80000000-0000-0000-0000-000000000002'),
  ('90000000-0000-0000-0000-000000000003', 'visit_cancelled', 'whatsapp', 'sent', 'carlos@fsm.cl', 'Visita cancelada', 'La visita ID 80000000-0004 ha sido cancelada', 'visit', '80000000-0000-0000-0000-000000000004')
ON CONFLICT (id) DO NOTHING;

-- Notification templates
INSERT INTO notifications.notification_templates (id, type, channel, subject, body, active) VALUES
  ('91000000-0000-0000-0000-000000000001', 'visit_assigned', 'push', 'Visita asignada', 'Se te ha asignado una visita en {{address}}', true),
  ('91000000-0000-0000-0000-000000000002', 'visit_reminder', 'email', 'Recordatorio', 'Tienes una visita programada para {{date}} a las {{time}}', true),
  ('91000000-0000-0000-0000-000000000003', 'sla_warning', 'push', 'Alerta SLA', 'La visita {{visit_id}} está próxima a vencer el SLA', true)
ON CONFLICT (id) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- AUDIT LOGS
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO core.plan_audit_log (id, entity_type, entity_id, action, new_data, user_id) VALUES
  ('AA000000-0000-0000-0000-000000000001', 'daily_load_material', '73000000-0000-0000-0000-000000000001', 'add', '{"material_id":"50000000-0000-0000-0000-000000000001","quantity":20}', '10000000-0000-0000-0000-000000000001'),
  ('AA000000-0000-0000-0000-000000000002', 'daily_load_tool', '74000000-0000-0000-0000-000000000001', 'add', '{"tool_id":"51000000-0000-0000-0000-000000000001"}', '10000000-0000-0000-0000-000000000001')
ON CONFLICT (id) DO NOTHING;

-- =============================================================================
-- END SEED DATA
-- =============================================================================
