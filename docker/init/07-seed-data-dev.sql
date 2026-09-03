-- =============================================================================
-- DEV SEED DATA — GasConecta SpA
-- Chilean gas maintenance company: installations, pre-visits, PH tests, maintenance
-- =============================================================================

-- PHASE 1: CORE TECH ROLES
INSERT INTO core.tech_roles (id, code, name, active) VALUES
  ('eadc5f0e-99f5-429d-b0f0-592cb6c0c24d', 'driver',     'Driver',     true),
  ('224b0832-b719-433f-8098-886282e4521d', 'technician', 'Technician', true),
  ('9811bfc2-6644-493c-9cdf-48924a48308a', 'helper',     'Helper',     true),
  ('343cc4d8-4af3-4245-b32d-7bceb6985150', 'operative',  'Operative',  true),
  ('f1a2b3c4-d5e6-4f7a-8b9c-0d1e2f3a4b5c', 'supervisor', 'Supervisor', true),
  ('a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d', 'foreman',    'Foreman',    true)
ON CONFLICT (id) DO NOTHING;

-- PHASE 2: USERS (20) with bcrypt password hashes
INSERT INTO core.users (id, email, role, name, password_hash) VALUES
  ('10000000-0000-0000-0000-000000000001', 'admin@fsm.cl',    'admin',     'Admin GasConecta',  '$2a$10$pNSweEcjk/sQjkX23oyLT.GvtvcuPcATN9pI.SL5hiZzima58OMm6'),
  ('10000000-0000-0000-0000-000000000002', 'carlos@fsm.cl',   'technician','Carlos Tecnico',    '$2a$10$woc0D9hCFGmBEoVhPH/OMeGxAzOAyfax0BAOSW/k.pJanJSxvnyPS'),
  ('10000000-0000-0000-0000-000000000003', 'ana@fsm.cl',      'technician','Ana Tecnica',       '$2a$10$QvlUqHRYMx6T/jVF.bHlR.Kqqfjyrev5UCWsf4zV31rPbe8B1/Unq'),
  ('10000000-0000-0000-0000-000000000004', 'pedro@fsm.cl',    'technician','Pedro Operario',    '$2a$10$Lx6hvzBFxIs7Do8uDdrWSe2HdmUNM8Jns/JkrW5WSKu/xwL.11oi.'),
  ('10000000-0000-0000-0000-000000000005', 'operador@fsm.cl', 'operator',  'Operador Central',  '$2a$10$USC8YMgZp/.PkrkFHFFBfui/Rshci.90Q9t2yNrZxcuHsNxbMjJz6'),
  ('10000000-0000-0000-0000-000000000006', 'rodrigo@fsm.cl',  'technician','Rodrigo Silva',     '$2a$10$8rxZfwDfhv.7bXmvFo0nm.wvykBdndzS2QQlGAQzYdCSwlYZ9I47q'),
  ('10000000-0000-0000-0000-000000000007', 'lucia@fsm.cl',    'technician','Lucia Fernandez',   '$2a$10$KB9Zb4MkKxzEjM3YTfy4hO28F8JHrwBYTn6i8Bsl9DFj/V1rL3u0W'),
  ('10000000-0000-0000-0000-000000000008', 'javier@fsm.cl',   'technician','Javier Morales',    '$2a$10$Mxb3.6lFa1YsW10x7928feLMSo2Bl8Fsfrocyv7L7zcPrUu.iKgb2'),
  ('10000000-0000-0000-0000-000000000009', 'marcos@fsm.cl',   'technician','Marcos Reyes',      '$2a$10$zRDPa3dhlPzsB/9kmYERVuR9BOwgHBtTqL/6EaeYaiVxOfdA1sx/G'),
  ('10000000-0000-0000-0000-000000000010', 'paula@fsm.cl',    'technician','Paula Vargas',      '$2a$10$oaKH34wYnIoFySd2q5bI9Oi5x2FeKleSwoSjEEIvaaxzEM2eojtyq'),
  ('10000000-0000-0000-0000-000000000011', 'fernanda@fsm.cl', 'technician','Fernanda Lopez',    '$2a$10$PLOdOYNNB7nT9mmRWVUR4O4lIxEuXCqoxSaAVTDQIl9JiC91gzvPK'),
  ('10000000-0000-0000-0000-000000000012', 'diego@fsm.cl',    'technician','Diego Castillo',    '$2a$10$WlIJnrBwcHFEx0l5704IyeUMcUunyGV6JSzSuvWQRhIBJ7AEEU6Di'),
  ('10000000-0000-0000-0000-000000000013', 'camila@fsm.cl',   'technician','Camila Herrera',    '$2a$10$dTIS6geboYhNfqOMhS2XZOrsmyYKhMbV1QikAyMBkgzZn9yf9ueUu'),
  ('10000000-0000-0000-0000-000000000014', 'sebastian@fsm.cl','technician','Sebastian Munoz',   '$2a$10$w41sOKA0GDpejYIeSJspe.yGtKQNJbkTDrveUa6fUdqG9MV2KHS/W'),
  ('10000000-0000-0000-0000-000000000015', 'patricio@fsm.cl', 'technician','Patricio Torres',   '$2a$10$/2i7.8XX.m5q3Fh4Tfq/PuaXjEVTOzEsKCukEE7bTE4XSCdn6mxH2'),
  ('10000000-0000-0000-0000-000000000016', 'valentina@fsm.cl','technician','Valentina Rojas',   '$2a$10$Po5MgAp/2EABnvgQOGFjPeNl32BXUP5PCZek/whN5XmpP5.kw1yre'),
  ('10000000-0000-0000-0000-000000000017', 'andres@fsm.cl',   'technician','Andres Guzman',     '$2a$10$890oPP/OiK4ZblqwINZFieAmcQ1jOOX55GdpBjT29imMv5D1W65N.'),
  ('10000000-0000-0000-0000-000000000018', 'claudio@fsm.cl',  'technician','Claudio Diaz',      '$2a$10$PKWOA3tGY9WBED4RMZv4ROHDrsEpmFVaQpeoybYRs0b8pZpNDYTuC'),
  ('10000000-0000-0000-0000-000000000019', 'isabel@fsm.cl',   'technician','Isabel Castro',     '$2a$10$ZwKsymmjebDAZfZ17.acCuttm4HNeiYlim5VpnwGBdWqa3Rc399rm'),
  ('10000000-0000-0000-0000-000000000020', 'renato@fsm.cl',   'technician','Renato Figueroa',   '$2a$10$Z7eebgxocJ70eVZwNYUVzuJiK31aEBifnn2kicGoQGdp7X1G8TksC')
ON CONFLICT (id) DO NOTHING;

-- PHASE 3: SHARED TEMPLATES
INSERT INTO shared.report_templates (id, name, fields_json) VALUES
  ('10000000-0000-0000-0000-000000000001', 'Checklist General',       '{"sections":[{"name":"General","fields":["resultado","observaciones"]}]}'),
  ('10000000-0000-0000-0000-000000000002', 'Reporte de Trabajo',      '{"sections":[{"name":"Trabajo","fields":["descripcion","horas","materiales"]}]}'),
  ('10000000-0000-0000-0000-000000000003', 'Inspeccion de Seguridad', '{"sections":[{"name":"Seguridad","fields":["hallazgos","riesgos","acciones"]}]}')
ON CONFLICT (id) DO NOTHING;

-- PHASE 4: GEOCODING ADDRESSES (15)
INSERT INTO geocoding.addresses (id, street, number, neighborhood, city, region, postal_code, geom) VALUES
  ('20000000-0000-0000-0000-000000000001', 'Av. Providencia',  '1234', 'Providencia',  'Santiago',     'Region Metropolitana', '7500000', ST_SetSRID(ST_MakePoint(-70.6123, -33.4256), 4326)),
  ('20000000-0000-0000-0000-000000000002', 'Av. Apoquindo',    '5678', 'Las Condes',   'Santiago',     'Region Metropolitana', '7550000', ST_SetSRID(ST_MakePoint(-70.5950, -33.4100), 4326)),
  ('20000000-0000-0000-0000-000000000003', 'Los Abedules',     '90',   'Nunoa',        'Santiago',     'Region Metropolitana', '7510000', ST_SetSRID(ST_MakePoint(-70.6000, -33.4500), 4326)),
  ('20000000-0000-0000-0000-000000000004', 'Av. La Estrella',  '456',  'Quilicura',    'Santiago',     'Region Metropolitana', '8700000', ST_SetSRID(ST_MakePoint(-70.7200, -33.3500), 4326)),
  ('20000000-0000-0000-0000-000000000005', 'Av. Vitacura',     '3210', 'Vitacura',     'Santiago',     'Region Metropolitana', '7630000', ST_SetSRID(ST_MakePoint(-70.5700, -33.3900), 4326)),
  ('20000000-0000-0000-0000-000000000006', 'Av. Circunvalar',  '1155', NULL,           'Talca',        'Region del Maule',    '3460000', ST_SetSRID(ST_MakePoint(-71.6554, -35.4264), 4326)),
  ('20000000-0000-0000-0000-000000000007', 'Av. Recoleta',     '1234', NULL,           'Talca',        'Region del Maule',    '3460000', ST_SetSRID(ST_MakePoint(-71.6600, -35.4300), 4326)),
  ('20000000-0000-0000-0000-000000000008', 'Calle 1 Norte',    '850',  NULL,           'Talca',        'Region del Maule',    '3460000', ST_SetSRID(ST_MakePoint(-71.6520, -35.4230), 4326)),
  ('20000000-0000-0000-0000-000000000009', 'Av. San Martin',   '890',  NULL,           'Curico',       'Region del Maule',    '3300000', ST_SetSRID(ST_MakePoint(-71.2335, -34.9828), 4326)),
  ('20000000-0000-0000-0000-000000000010', 'Av. Freire',       '567',  NULL,           'Curico',       'Region del Maule',    '3300000', ST_SetSRID(ST_MakePoint(-71.2400, -34.9800), 4326)),
  ('20000000-0000-0000-0000-000000000011', 'Av. Puerto',       's/n',  NULL,           'Curico',       'Region del Maule',    '3300000', ST_SetSRID(ST_MakePoint(-71.2300, -34.9750), 4326)),
  ('20000000-0000-0000-0000-000000000012', 'Av. Confederacion','456',  NULL,           'Linares',      'Region del Maule',    '3730000', ST_SetSRID(ST_MakePoint(-71.5900, -35.8500), 4326)),
  ('20000000-0000-0000-0000-000000000013', 'Av. San Martin',   '1200', NULL,           'Linares',      'Region del Maule',    '3730000', ST_SetSRID(ST_MakePoint(-71.5850, -35.8470), 4326)),
  ('20000000-0000-0000-0000-000000000014', 'Ruta 115',         'km 5', NULL,           'San Clemente', 'Region del Maule',    '3530000', ST_SetSRID(ST_MakePoint(-71.4833, -35.5333), 4326)),
  ('20000000-0000-0000-0000-000000000015', 'Av. Comercio',     '320',  NULL,           'Cauquenes',    'Region del Maule',    '3690000', ST_SetSRID(ST_MakePoint(-72.3167, -35.9667), 4326))
ON CONFLICT (id) DO NOTHING;

-- PHASE 5: PARTNERS (4)
INSERT INTO partners.partners (id, name, tax_id, status) VALUES
  ('30000000-0000-0000-0000-000000000001', 'SecurGas SpA',   '76.123.456-7', 'active'),
  ('30000000-0000-0000-0000-000000000002', 'PetroGas Chile', '76.234.567-8', 'active'),
  ('30000000-0000-0000-0000-000000000003', 'GasMaule SpA',   '76.800.001-1', 'active'),
  ('30000000-0000-0000-0000-000000000004', 'SecurGas Sur',   '76.800.002-2', 'active')
ON CONFLICT (id) DO NOTHING;

INSERT INTO partners.partner_contacts (id, partner_id, name, position, phone, email) VALUES
  ('31000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'Juan Perez',     'Gerente',         '+56912345678', 'juan.perez@securgas.cl'),
  ('31000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', 'Maria Gonzalez', 'Supervisora',     '+56923456789', 'maria.gonzalez@petrogas.cl'),
  ('31000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', 'Pedro Soto',     'Jefe Operaciones','+56934567890', 'pedro.soto@gasmaule.cl'),
  ('31000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000004', 'Claudia Rios',   'Gerenta',         '+56945678901', 'claudia.rios@securgassur.cl')
ON CONFLICT (id) DO NOTHING;

INSERT INTO partners.partner_agreements (id, partner_id, service_type, rate, active) VALUES
  ('32000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', NULL, 45000, true),
  ('32000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', NULL, 52000, true),
  ('32000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', NULL, 42000, true),
  ('32000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000004', NULL, 48000, true)
ON CONFLICT (id) DO NOTHING;

INSERT INTO partners.partner_agreement_docs (id, agreement_id, work_type, doc_name, required) VALUES
  ('33000000-0000-0000-0000-000000000001', '32000000-0000-0000-0000-000000000001', 'Instalacion',   'Certificado de habilitacion', true),
  ('33000000-0000-0000-0000-000000000002', '32000000-0000-0000-0000-000000000002', 'Mantenimiento', 'Poliza de Seguridad',         true),
  ('33000000-0000-0000-0000-000000000003', '32000000-0000-0000-0000-000000000003', 'Instalacion',   'Certificado SEC',             true),
  ('33000000-0000-0000-0000-000000000004', '32000000-0000-0000-0000-000000000004', 'Emergencia',    'Poliza de Responsabilidad',   true)
ON CONFLICT (id) DO NOTHING;

INSERT INTO partners.partner_agreement_forms (id, agreement_id, work_type, form_template_id, quantity) VALUES
  ('34000000-0000-0000-0000-000000000001', '32000000-0000-0000-0000-000000000001', 'Instalacion',   '10000000-0000-0000-0000-000000000001', 1),
  ('34000000-0000-0000-0000-000000000002', '32000000-0000-0000-0000-000000000002', 'Mantenimiento', '10000000-0000-0000-0000-000000000002', 1)
ON CONFLICT (id) DO NOTHING;

-- PHASE 6: SLAs (4)
INSERT INTO partners.slas (id, partner_id, name, description, response_hours, resolution_hours, compliance_target, active, valid_from, valid_until) VALUES
  ('35000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', 'SLA Gas Normal Santiago',  'Respuesta 24h, resolucion 72h',  24, 72, 0.95, true, '2026-01-01', '2026-12-31'),
  ('35000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', 'SLA Gas Urgente Santiago', 'Respuesta 4h, resolucion 24h',    4,  24, 0.98, true, '2026-01-01', '2026-12-31'),
  ('35000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', 'SLA Gas Normal Maule',     'Respuesta 24h, resolucion 72h',  24, 72, 0.95, true, '2026-01-01', '2026-12-31'),
  ('35000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000004', 'SLA Gas Urgente Maule',    'Respuesta 4h, resolucion 24h',    4,  24, 0.98, true, '2026-01-01', '2026-12-31')
ON CONFLICT (id) DO NOTHING;

-- PHASE 7: CUSTOMERS (15)
INSERT INTO customers.customers (id, name, phone, email, tax_id) VALUES
  ('40000000-0000-0000-0000-000000000001', 'Restaurante La Cava',           '+56912345678', 'contactos@lacava.cl',          '76.500.001-1'),
  ('40000000-0000-0000-0000-000000000002', 'Hotel Los Andes',                '+56923456789', 'facilities@losandes.cl',       '76.500.002-2'),
  ('40000000-0000-0000-0000-000000000003', 'Conjunto Residencial El Bosque', '+56934567890', 'admin@elbosque.cl',            '76.500.003-3'),
  ('40000000-0000-0000-0000-000000000004', 'Industrial Panaderia Norte',     '+56945678901', 'mantencion@panaderianorte.cl', '76.500.004-4'),
  ('40000000-0000-0000-0000-000000000005', 'Colegio San Patricio',           '+56956789012', 'admin@sanpatricio.cl',         '76.500.005-5'),
  ('40000000-0000-0000-0000-000000000006', 'Hotel Plaza Talca',              '+56967890123', 'reservas@hotelplazatalca.cl',  '76.800.006-6'),
  ('40000000-0000-0000-0000-000000000007', 'Clinica Maule',                  '+56978901234', 'admision@clinicamaule.cl',     '76.800.007-7'),
  ('40000000-0000-0000-0000-000000000008', 'Restaurante Donde Juancho',      '+56989012345', 'reservas@dondejuancho.cl',     '76.800.008-8'),
  ('40000000-0000-0000-0000-000000000009', 'Restaurante Rincon Curicano',    '+56990123456', 'contacto@rinconcuricano.cl',   '76.800.009-9'),
  ('40000000-0000-0000-0000-000000000010', 'Supermercado Express Curico',    '+56901234567', 'admin@expresscurico.cl',       '76.800.010-0'),
  ('40000000-0000-0000-0000-000000000011', 'Hotel Central Curico',           '+56912340001', 'recepcion@hotelcentral.cl',    '76.800.011-1'),
  ('40000000-0000-0000-0000-000000000012', 'Bodega Linares SpA',             '+56923400002', 'bodega@linaresspa.cl',         '76.800.012-2'),
  ('40000000-0000-0000-0000-000000000013', 'Farmacia Cruz Verde Linares',    '+56934500003', 'farmacia@cruzverde.cl',        '76.800.013-3'),
  ('40000000-0000-0000-0000-000000000014', 'Agricola San Clemente',          '+56945600004', 'contacto@agrosancl.cl',        '76.800.014-4'),
  ('40000000-0000-0000-0000-000000000015', 'Molino Cauquenes',               '+56956700005', 'admin@molinocauquenes.cl',     '76.800.015-5')
ON CONFLICT (id) DO NOTHING;

-- PHASE 8: CUSTOMER ADDRESSES (15)
INSERT INTO customers.customer_addresses (id, customer_id, address_id, name, type) VALUES
  ('41000000-0000-0000-0000-000000000001', '40000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 'Restaurante La Cava',           'commercial'),
  ('41000000-0000-0000-0000-000000000002', '40000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000002', 'Hotel Los Andes',                'commercial'),
  ('41000000-0000-0000-0000-000000000003', '40000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000003', 'Conjunto Residencial El Bosque', 'residential'),
  ('41000000-0000-0000-0000-000000000004', '40000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000004', 'Industrial Panaderia Norte',     'industrial'),
  ('41000000-0000-0000-0000-000000000005', '40000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000005', 'Colegio San Patricio',           'commercial'),
  ('41000000-0000-0000-0000-000000000006', '40000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000006', 'Hotel Plaza Talca',              'commercial'),
  ('41000000-0000-0000-0000-000000000007', '40000000-0000-0000-0000-000000000007', '20000000-0000-0000-0000-000000000007', 'Clinica Maule',                  'institutional'),
  ('41000000-0000-0000-0000-000000000008', '40000000-0000-0000-0000-000000000008', '20000000-0000-0000-0000-000000000008', 'Restaurante Donde Juancho',      'commercial'),
  ('41000000-0000-0000-0000-000000000009', '40000000-0000-0000-0000-000000000009', '20000000-0000-0000-0000-000000000009', 'Restaurante Rincon Curicano',    'commercial'),
  ('41000000-0000-0000-0000-000000000010', '40000000-0000-0000-0000-000000000010', '20000000-0000-0000-0000-000000000010', 'Supermercado Express Curico',    'commercial'),
  ('41000000-0000-0000-0000-000000000011', '40000000-0000-0000-0000-000000000011', '20000000-0000-0000-0000-000000000011', 'Hotel Central Curico',           'commercial'),
  ('41000000-0000-0000-0000-000000000012', '40000000-0000-0000-0000-000000000012', '20000000-0000-0000-0000-000000000012', 'Bodega Linares SpA',             'industrial'),
  ('41000000-0000-0000-0000-000000000013', '40000000-0000-0000-0000-000000000013', '20000000-0000-0000-0000-000000000013', 'Farmacia Cruz Verde Linares',    'commercial'),
  ('41000000-0000-0000-0000-000000000014', '40000000-0000-0000-0000-000000000014', '20000000-0000-0000-0000-000000000014', 'Agricola San Clemente',          'industrial'),
  ('41000000-0000-0000-0000-000000000015', '40000000-0000-0000-0000-000000000015', '20000000-0000-0000-0000-000000000015', 'Molino Cauquenes',               'industrial')
ON CONFLICT (id) DO NOTHING;

INSERT INTO customers.customer_address_partners (id, customer_address_id, partner_id, from_date, to_date) VALUES
  ('42000000-0000-0000-0000-000000000001', '41000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', '2026-01-01', NULL),
  ('42000000-0000-0000-0000-000000000002', '41000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001', '2026-01-01', NULL),
  ('42000000-0000-0000-0000-000000000003', '41000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000003', '2026-01-01', NULL),
  ('42000000-0000-0000-0000-000000000004', '41000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000003', '2026-01-01', NULL),
  ('42000000-0000-0000-0000-000000000005', '41000000-0000-0000-0000-000000000009', '30000000-0000-0000-0000-000000000004', '2026-01-01', NULL),
  ('42000000-0000-0000-0000-000000000006', '41000000-0000-0000-0000-000000000012', '30000000-0000-0000-0000-000000000003', '2026-01-01', NULL),
  ('42000000-0000-0000-0000-000000000007', '41000000-0000-0000-0000-000000000015', '30000000-0000-0000-0000-000000000004', '2026-01-01', NULL)
ON CONFLICT (id) DO NOTHING;

-- PHASE 9: PROPERTIES (15)
INSERT INTO customers.properties (id, customer_address_id, name, notes) VALUES
  ('43000000-0000-0000-0000-000000000001', '41000000-0000-0000-0000-000000000001', 'Restaurante La Cava',           'Estanque GLP 45kg'),
  ('43000000-0000-0000-0000-000000000002', '41000000-0000-0000-0000-000000000002', 'Hotel Los Andes',                'Estanque GLP 120kg + medidor'),
  ('43000000-0000-0000-0000-000000000003', '41000000-0000-0000-0000-000000000003', 'Conjunto Residencial El Bosque', 'Estanque GLP 200kg'),
  ('43000000-0000-0000-0000-000000000004', '41000000-0000-0000-0000-000000000004', 'Industrial Panaderia Norte',     'Estanque GLP 450kg'),
  ('43000000-0000-0000-0000-000000000005', '41000000-0000-0000-0000-000000000005', 'Colegio San Patricio',           'Estanque GLP 100kg'),
  ('43000000-0000-0000-0000-000000000006', '41000000-0000-0000-0000-000000000006', 'Hotel Plaza Talca',              'Estanque GLP 200kg + medidor'),
  ('43000000-0000-0000-0000-000000000007', '41000000-0000-0000-0000-000000000007', 'Clinica Maule',                  'Estanque GLP 120kg + tuberia'),
  ('43000000-0000-0000-0000-000000000008', '41000000-0000-0000-0000-000000000008', 'Restaurante Donde Juancho',      'Estanque GLP 45kg'),
  ('43000000-0000-0000-0000-000000000009', '41000000-0000-0000-0000-000000000009', 'Restaurante Rincon Curicano',    'Estanque GLP 45kg'),
  ('43000000-0000-0000-0000-000000000010', '41000000-0000-0000-0000-000000000010', 'Supermercado Express Curico',    'Estanque GLP 100kg'),
  ('43000000-0000-0000-0000-000000000011', '41000000-0000-0000-0000-000000000011', 'Hotel Central Curico',           'Estanque GLP 120kg'),
  ('43000000-0000-0000-0000-000000000012', '41000000-0000-0000-0000-000000000012', 'Bodega Linares SpA',             'Estanque GLP 450kg'),
  ('43000000-0000-0000-0000-000000000013', '41000000-0000-0000-0000-000000000013', 'Farmacia Cruz Verde Linares',    'Estanque GLP 45kg'),
  ('43000000-0000-0000-0000-000000000014', '41000000-0000-0000-0000-000000000014', 'Agricola San Clemente',          'Estanque GLP 200kg'),
  ('43000000-0000-0000-0000-000000000015', '41000000-0000-0000-0000-000000000015', 'Molino Cauquenes',               'Estanque GLP 450kg + tuberia')
ON CONFLICT (id) DO NOTHING;

-- PHASE 10: INVENTORY — MATERIALS (10)
INSERT INTO inventory.materials (id, sku, name, unit, unit_cost) VALUES
  ('50000000-0000-0000-0000-000000000001', 'GAS-001', 'Tubo acero galvanizado 1/2"',   'ml',     3500),
  ('50000000-0000-0000-0000-000000000002', 'GAS-002', 'Tubo acero galvanizado 3/4"',   'ml',     5200),
  ('50000000-0000-0000-0000-000000000003', 'GAS-003', 'Conexion flexible gas 1/2"',    'unidad',  8500),
  ('50000000-0000-0000-0000-000000000004', 'GAS-004', 'Conexion flexible gas 3/4"',    'unidad', 12000),
  ('50000000-0000-0000-0000-000000000005', 'GAS-005', 'Valvula de corte 1/2"',         'unidad', 15000),
  ('50000000-0000-0000-0000-000000000006', 'GAS-006', 'Valvula de corte 3/4"',         'unidad', 22000),
  ('50000000-0000-0000-0000-000000000007', 'GAS-007', 'Cinta teflon gas',              'rollo',   1200),
  ('50000000-0000-0000-0000-000000000008', 'GAS-008', 'Abrazaderas acero',             'unidad',   800),
  ('50000000-0000-0000-0000-000000000009', 'GAS-009', 'Pasta selladora gas',           'tubo',    4500),
  ('50000000-0000-0000-0000-000000000010', 'GAS-010', 'Manometro calibrado',           'unidad', 85000)
ON CONFLICT (id) DO NOTHING;

-- PHASE 11: INVENTORY — TOOLS (8)
INSERT INTO inventory.tools (id, code, name, status) VALUES
  ('51000000-0000-0000-0000-000000000001', 'HG-001', 'Detector de gas multigas',      'available'),
  ('51000000-0000-0000-0000-000000000002', 'HG-002', 'Manometro calibrado 0-10 bar',  'available'),
  ('51000000-0000-0000-0000-000000000003', 'HG-003', 'Llave inglesa 12"',             'available'),
  ('51000000-0000-0000-0000-000000000004', 'HG-004', 'Llave inglesa 18"',             'available'),
  ('51000000-0000-0000-0000-000000000005', 'HG-005', 'Llave Stillson 14"',            'available'),
  ('51000000-0000-0000-0000-000000000006', 'HG-006', 'Cortatubos industrial',         'available'),
  ('51000000-0000-0000-0000-000000000007', 'HG-007', 'Enroscadora manual 1/2"',       'available'),
  ('51000000-0000-0000-0000-000000000008', 'HG-008', 'Detector de fugas ultrasonico', 'available')
ON CONFLICT (id) DO NOTHING;

-- PHASE 12: INVENTORY — EPP (10)
INSERT INTO inventory.epp_items (id, name, type, lifecycle) VALUES
  ('52000000-0000-0000-0000-000000000001', 'Casco MSA',              'head',        'reusable'),
  ('52000000-0000-0000-0000-000000000002', 'Guantes nitrilo gas',    'hands',       'disposable'),
  ('52000000-0000-0000-0000-000000000003', 'Chaleco reflectante',    'body',        'reusable'),
  ('52000000-0000-0000-0000-000000000004', 'Lentes seguridad',       'eyes',        'disposable'),
  ('52000000-0000-0000-0000-000000000005', 'Tapones oidos',          'ears',        'disposable'),
  ('52000000-0000-0000-0000-000000000006', 'Botas seguridad',        'feet',        'reusable'),
  ('52000000-0000-0000-0000-000000000007', 'Mascarilla N95',         'respiratory', 'disposable'),
  ('52000000-0000-0000-0000-000000000008', 'Mascarilla gas GLP',     'respiratory', 'disposable'),
  ('52000000-0000-0000-0000-000000000009', 'Guantes cuero soldador', 'hands',       'reusable'),
  ('52000000-0000-0000-0000-000000000010', 'Arnes seguridad',        'body',        'reusable')
ON CONFLICT (id) DO NOTHING;

-- PHASE 13: INVENTORY — VEHICLES (8)
INSERT INTO inventory.vehicles (id, type, license_plate, name, brand, model, year, status, capacity) VALUES
  ('53000000-0000-0000-0000-000000000001', NULL, 'FG-3456', 'Camioneta Tecnico 1', 'Toyota',  'Hilux',  2023, 'available', '5 personas'),
  ('53000000-0000-0000-0000-000000000002', NULL, 'GH-7890', 'Camioneta Tecnico 2', 'Nissan',  'NP300',  2024, 'available', '5 personas'),
  ('53000000-0000-0000-0000-000000000003', NULL, 'IJ-1234', 'Furgon Materiales',   'Hyundai', 'H1 Van', 2023, 'available', '1200 kg'),
  ('53000000-0000-0000-0000-000000000004', NULL, 'KL-5678', 'Camion Liviano',      'Toyota',  'Dyna',   2022, 'available', '3500 kg'),
  ('53000000-0000-0000-0000-000000000005', NULL, 'MN-9012', 'Camioneta Talca 1',   'Toyota',  'Hilux',  2023, 'available', '5 personas'),
  ('53000000-0000-0000-0000-000000000006', NULL, 'OP-3456', 'Camioneta Talca 2',   'Nissan',  'NP300',  2024, 'available', '5 personas'),
  ('53000000-0000-0000-0000-000000000007', NULL, 'QR-7890', 'Furgon Curico',       'Toyota',  'Hilux',  2022, 'available', '5 personas'),
  ('53000000-0000-0000-0000-000000000008', NULL, 'ST-1234', 'Camioneta Linares',   'Nissan',  'NP300',  2023, 'available', '5 personas')
ON CONFLICT (id) DO NOTHING;

-- PHASE 14: TECHNICIANS (17)
INSERT INTO customers.technicians (id, user_id, name, is_active, specialties, employee_number) VALUES
  ('60000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000002', 'Carlos Tecnico',    true, ARRAY['gas_installation','gas_maintenance','ph_testing'], 'EMP-001'),
  ('60000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000003', 'Ana Tecnica',       true, ARRAY['gas_installation','inspection','emergency'],     'EMP-002'),
  ('60000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000004', 'Pedro Operario',    true, ARRAY['gas_installation','general'],                    'EMP-003'),
  ('60000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000006', 'Rodrigo Silva',     true, ARRAY['gas_installation','gas_maintenance'],            'EMP-004'),
  ('60000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000007', 'Lucia Fernandez',   true, ARRAY['inspection','ph_testing'],                       'EMP-005'),
  ('60000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000008', 'Javier Morales',    true, ARRAY['gas_installation','general'],                    'EMP-006'),
  ('60000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000009', 'Marcos Reyes',      true, ARRAY['driving','logistics'],                           'EMP-007'),
  ('60000000-0000-0000-0000-000000000008', '10000000-0000-0000-0000-000000000010', 'Paula Vargas',      true, ARRAY['gas_installation','inspection'],                'EMP-008'),
  ('60000000-0000-0000-0000-000000000009', '10000000-0000-0000-0000-000000000011', 'Fernanda Lopez',    true, ARRAY['gas_installation','emergency'],                 'EMP-009'),
  ('60000000-0000-0000-0000-000000000010', '10000000-0000-0000-0000-000000000012', 'Diego Castillo',    true, ARRAY['gas_installation','general'],                   'EMP-010'),
  ('60000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000013', 'Camila Herrera',    true, ARRAY['general','logistics'],                          'EMP-011'),
  ('60000000-0000-0000-0000-000000000012', '10000000-0000-0000-0000-000000000014', 'Sebastian Munoz',   true, ARRAY['gas_installation','inspection'],                'EMP-012'),
  ('60000000-0000-0000-0000-000000000013', '10000000-0000-0000-0000-000000000015', 'Patricio Torres',   true, ARRAY['gas_maintenance','inspection'],                 'EMP-013'),
  ('60000000-0000-0000-0000-000000000014', '10000000-0000-0000-0000-000000000016', 'Valentina Rojas',   true, ARRAY['ph_testing','gas_installation'],                'EMP-014'),
  ('60000000-0000-0000-0000-000000000015', '10000000-0000-0000-0000-000000000017', 'Andres Guzman',     true, ARRAY['general'],                                      'EMP-015'),
  ('60000000-0000-0000-0000-000000000016', '10000000-0000-0000-0000-000000000018', 'Claudio Diaz',      true, ARRAY['gas_installation','gas_maintenance'],           'EMP-016'),
  ('60000000-0000-0000-0000-000000000017', '10000000-0000-0000-0000-000000000019', 'Isabel Castro',     true, ARRAY['gas_installation','general'],                   'EMP-017')
ON CONFLICT (id) DO NOTHING;

INSERT INTO customers.tech_certifications (id, tech_id, cert_id, number, issuer, issue_date, expiry_date) VALUES
  ('61000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', NULL, 'CERT-2024-001', 'SEC Chile',         '2024-06-01', '2027-06-01'),
  ('61000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000001', NULL, 'CERT-2024-003', 'SEC Chile',         '2024-06-01', '2027-06-01'),
  ('61000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000002', NULL, 'CERT-2024-002', 'SEC Chile',         '2024-08-15', '2027-08-15'),
  ('61000000-0000-0000-0000-000000000004', '60000000-0000-0000-0000-000000000002', NULL, 'CERT-2024-004', 'SEC Chile',         '2024-08-15', '2027-08-15'),
  ('61000000-0000-0000-0000-000000000005', '60000000-0000-0000-0000-000000000003', NULL, 'CERT-2024-005', 'SEC Chile',         '2024-10-01', '2027-10-01'),
  ('61000000-0000-0000-0000-000000000006', '60000000-0000-0000-0000-000000000004', NULL, 'CERT-2024-006', 'SEC Chile',         '2024-03-01', '2027-03-01'),
  ('61000000-0000-0000-0000-000000000007', '60000000-0000-0000-0000-000000000005', NULL, 'CERT-2024-007', 'SEC Chile',         '2024-04-01', '2027-04-01'),
  ('61000000-0000-0000-0000-000000000008', '60000000-0000-0000-0000-000000000008', NULL, 'CERT-2024-008', 'SEC Chile',         '2024-05-01', '2027-05-01'),
  ('61000000-0000-0000-0000-000000000009', '60000000-0000-0000-0000-000000000009', NULL, 'CERT-2024-009', 'SEC Chile',         '2024-02-01', '2027-02-01'),
  ('61000000-0000-0000-0000-000000000010', '60000000-0000-0000-0000-000000000012', NULL, 'CERT-2024-010', 'SEC Chile',         '2024-01-01', '2027-01-01'),
  ('61000000-0000-0000-0000-000000000011', '60000000-0000-0000-0000-000000000013', NULL, 'CERT-2024-011', 'SEC Chile',         '2024-07-01', '2027-07-01'),
  ('61000000-0000-0000-0000-000000000012', '60000000-0000-0000-0000-000000000014', NULL, 'CERT-2024-012', 'SEC Chile',         '2024-09-01', '2027-09-01'),
  ('61000000-0000-0000-0000-000000000013', '60000000-0000-0000-0000-000000000016', NULL, 'CERT-2024-013', 'SEC Chile',         '2024-11-01', '2027-11-01')
ON CONFLICT (id) DO NOTHING;

-- PHASE 15: CHECKLIST TEMPLATES (3)
INSERT INTO inventory.checklist_templates (id, name, description, visit_type, work_type) VALUES
  ('56000000-0000-0000-0000-000000000001', 'Previsita Instalacion', 'Checklist para previsitgas', (SELECT id FROM operations.visit_types WHERE code = 'pre_visit'), (SELECT id FROM partners.partner_service_types WHERE code = 'installation')),
  ('56000000-0000-0000-0000-000000000002', 'Instalacion Estanque',  'Checklist instalacion gas', (SELECT id FROM operations.visit_types WHERE code = 'installation'), (SELECT id FROM partners.partner_service_types WHERE code = 'installation')),
  ('56000000-0000-0000-0000-000000000003', 'Prueba Hermeticidad',   'Checklist PH',             (SELECT id FROM operations.visit_types WHERE code = 'inspection'), (SELECT id FROM partners.partner_service_types WHERE code = 'inspection'))
ON CONFLICT (id) DO UPDATE SET description = EXCLUDED.description, visit_type = EXCLUDED.visit_type, work_type = EXCLUDED.work_type;

-- Checklist template materials
INSERT INTO operations.checklist_template_materials (id, template_id, material_id, default_quantity) VALUES
  ('56100000-0000-0000-0000-000000000001', '56000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 2),
  ('56100000-0000-0000-0000-000000000002', '56000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000003', 1),
  ('56100000-0000-0000-0000-000000000003', '56000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000005', 1),
  ('56100000-0000-0000-0000-000000000004', '56000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000002', 5),
  ('56100000-0000-0000-0000-000000000005', '56000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000004', 2),
  ('56100000-0000-0000-0000-000000000006', '56000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000006', 2),
  ('56100000-0000-0000-0000-000000000007', '56000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000007', 2),
  ('56100000-0000-0000-0000-000000000008', '56000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000008', 10),
  ('56100000-0000-0000-0000-000000000009', '56000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000009', 1),
  ('56100000-0000-0000-0000-000000000010', '56000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000010', 1)
ON CONFLICT (id) DO NOTHING;

-- Checklist template tools
INSERT INTO operations.checklist_template_tools (id, template_id, tool_id) VALUES
  ('56200000-0000-0000-0000-000000000001', '56000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000001'),
  ('56200000-0000-0000-0000-000000000002', '56000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000002'),
  ('56200000-0000-0000-0000-000000000003', '56000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000003'),
  ('56200000-0000-0000-0000-000000000004', '56000000-0000-0000-0000-000000000002', '51000000-0000-0000-0000-000000000006'),
  ('56200000-0000-0000-0000-000000000005', '56000000-0000-0000-0000-000000000002', '51000000-0000-0000-0000-000000000007'),
  ('56200000-0000-0000-0000-000000000006', '56000000-0000-0000-0000-000000000002', '51000000-0000-0000-0000-000000000005'),
  ('56200000-0000-0000-0000-000000000007', '56000000-0000-0000-0000-000000000002', '51000000-0000-0000-0000-000000000002'),
  ('56200000-0000-0000-0000-000000000008', '56000000-0000-0000-0000-000000000003', '51000000-0000-0000-0000-000000000002'),
  ('56200000-0000-0000-0000-000000000009', '56000000-0000-0000-0000-000000000003', '51000000-0000-0000-0000-000000000008')
ON CONFLICT (id) DO NOTHING;

-- Checklist template EPPs
INSERT INTO operations.checklist_template_epps (id, template_id, epp_id) VALUES
  ('56300000-0000-0000-0000-000000000001', '56000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000001'),
  ('56300000-0000-0000-0000-000000000002', '56000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000004'),
  ('56300000-0000-0000-0000-000000000003', '56000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000006'),
  ('56300000-0000-0000-0000-000000000004', '56000000-0000-0000-0000-000000000002', '52000000-0000-0000-0000-000000000001'),
  ('56300000-0000-0000-0000-000000000005', '56000000-0000-0000-0000-000000000002', '52000000-0000-0000-0000-000000000002'),
  ('56300000-0000-0000-0000-000000000006', '56000000-0000-0000-0000-000000000002', '52000000-0000-0000-0000-000000000004'),
  ('56300000-0000-0000-0000-000000000007', '56000000-0000-0000-0000-000000000002', '52000000-0000-0000-0000-000000000006'),
  ('56300000-0000-0000-0000-000000000008', '56000000-0000-0000-0000-000000000002', '52000000-0000-0000-0000-000000000008'),
  ('56300000-0000-0000-0000-000000000009', '56000000-0000-0000-0000-000000000003', '52000000-0000-0000-0000-000000000001'),
  ('56300000-0000-0000-0000-000000000010', '56000000-0000-0000-0000-000000000003', '52000000-0000-0000-0000-000000000004'),
  ('56300000-0000-0000-0000-000000000011', '56000000-0000-0000-0000-000000000003', '52000000-0000-0000-0000-000000000002'),
  ('56300000-0000-0000-0000-000000000012', '56000000-0000-0000-0000-000000000003', '52000000-0000-0000-0000-000000000008')
ON CONFLICT (id) DO NOTHING;

-- PHASE 16: PLANNING (4 plans, 4 routes)
INSERT INTO planning.daily_plans (id, name, date, notes) VALUES
  ('70000000-0000-0000-0000-000000000001', 'Ruta Santiago 21/08',  CURRENT_DATE,      'Ruta instalaciones Santiago'),
  ('70000000-0000-0000-0000-000000000002', 'Ruta Talca 21/08',    CURRENT_DATE,      'Ruta instalaciones Talca'),
  ('70000000-0000-0000-0000-000000000003', 'Ruta Curico 21/08',   CURRENT_DATE,      'Ruta mixta Curico'),
  ('70000000-0000-0000-0000-000000000004', 'Emergencia Maule 21/08', CURRENT_DATE,  'Emergencias Maule')
ON CONFLICT (id) DO NOTHING;

INSERT INTO planning.routes (id, type, date, daily_plan_id, notes, status) VALUES
  ('71000000-0000-0000-0000-000000000001', NULL, CURRENT_DATE, '70000000-0000-0000-0000-000000000001', 'Ruta instalaciones Santiago', 'scheduled'),
  ('71000000-0000-0000-0000-000000000002', NULL, CURRENT_DATE, '70000000-0000-0000-0000-000000000002', 'Ruta instalaciones Talca',   'scheduled'),
  ('71000000-0000-0000-0000-000000000003', NULL, CURRENT_DATE, '70000000-0000-0000-0000-000000000003', 'Ruta mixta Curico',           'scheduled'),
  ('71000000-0000-0000-0000-000000000004', NULL, CURRENT_DATE, '70000000-0000-0000-0000-000000000004', 'Ruta emergencias Maule',      'in_progress')
ON CONFLICT (id) DO NOTHING;

INSERT INTO planning.daily_plan_assignments (id, daily_plan_id, tech_id, role_id) VALUES
  ('72000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', (SELECT id FROM core.tech_roles WHERE code = 'technician')),
  ('72000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000006', (SELECT id FROM core.tech_roles WHERE code = 'helper')),
  ('72000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000004', (SELECT id FROM core.tech_roles WHERE code = 'technician')),
  ('72000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000006', (SELECT id FROM core.tech_roles WHERE code = 'helper')),
  ('72000000-0000-0000-0000-000000000005', '70000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000009', (SELECT id FROM core.tech_roles WHERE code = 'technician')),
  ('72000000-0000-0000-0000-000000000006', '70000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000010', (SELECT id FROM core.tech_roles WHERE code = 'operative')),
  ('72000000-0000-0000-0000-000000000007', '70000000-0000-0000-0000-000000000004', '60000000-0000-0000-0000-000000000001', (SELECT id FROM core.tech_roles WHERE code = 'technician')),
  ('72000000-0000-0000-0000-000000000008', '70000000-0000-0000-0000-000000000004', '60000000-0000-0000-0000-000000000013', (SELECT id FROM core.tech_roles WHERE code = 'technician'))
ON CONFLICT (id) DO NOTHING;

-- Daily loads
INSERT INTO planning.daily_load_materials (id, daily_plan_id, material_id, loaded_quantity) VALUES
  ('73000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 20),
  ('73000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000003', 10),
  ('73000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000005', 5),
  ('73000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000002', 20),
  ('73000000-0000-0000-0000-000000000005', '70000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000004', 8)
ON CONFLICT (id) DO NOTHING;

INSERT INTO planning.daily_load_tools (id, daily_plan_id, tool_id, loaded_quantity) VALUES
  ('74000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000001', 1),
  ('74000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000002', 1),
  ('74000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000003', 2),
  ('74000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000006', 1),
  ('74000000-0000-0000-0000-000000000005', '70000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000007', 1)
ON CONFLICT (id) DO NOTHING;

INSERT INTO planning.daily_load_epps (id, daily_plan_id, epp_id, loaded_quantity) VALUES
  ('75000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000001', 3),
  ('75000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000004', 3),
  ('75000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000006', 3),
  ('75000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000002', 5),
  ('75000000-0000-0000-0000-000000000005', '70000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000008', 10)
ON CONFLICT (id) DO NOTHING;
-- PHASE 17: OPERATIONS.VISITS (19 visits)

INSERT INTO operations.visits (
  id, property_id, partner_id, route_id, daily_plan_id,
  type, status, priority, source, billing_to,
  scheduled_at, started_at, completed_at, notes, result, rejection_reason_id
) VALUES
  ('80000000-0000-0000-0000-000000000001', '43000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000001', '71000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001',
   '7c542e71-8be4-4a5f-91ee-a1fb013a4063', 'scheduled', 'normal', 'portal', 'customer',
   (CURRENT_DATE + interval '0 hours'), NULL, NULL, 'Mantenimiento preventivo estanque principal', NULL, NULL),
  ('80000000-0000-0000-0000-000000000002', '43000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000002', '71000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001',
   'facd0f93-c680-437d-aa98-42df62814701', 'in_progress', 'normal', 'phone', 'partner',
   (CURRENT_DATE + interval '2 hours'), NULL, NULL, 'Inspeccion de hermeticidad gas', NULL, NULL),
  ('80000000-0000-0000-0000-000000000003', '43000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000003', '71000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000002',
   '40e1f2b8-4dda-4af1-92ff-1aa07e802010', 'completed', 'normal', 'email', 'customer',
   (CURRENT_DATE + interval '4 hours'), NULL, NULL, 'Instalacion nuevo estanque', 'd38f7943-ef5d-4ad3-82b0-529c3e48ae91', NULL),
  ('80000000-0000-0000-0000-000000000004', '43000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000004', '71000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001',
   '40e1f2b8-4dda-4af1-92ff-1aa07e802010', 'cancelled', 'normal', 'portal', 'customer',
   (CURRENT_DATE + interval '8 hours'), NULL, NULL, 'Instalacion cancelada por cliente', NULL, '6965b638-740a-40ee-80f5-7c41867cfa53'),
  ('80000000-0000-0000-0000-000000000005', '43000000-0000-0000-0000-000000000005', '30000000-0000-0000-0000-000000000001', '71000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000002',
   '29097af9-1645-4ce5-b54f-d042b431487a', 'completed', 'normal', 'portal', 'customer',
   (CURRENT_DATE + interval '0 hours'), NULL, NULL, 'Previsita aprobada', 'd38f7943-ef5d-4ad3-82b0-529c3e48ae91', NULL),
  ('80000000-0000-0000-0000-000000000006', '43000000-0000-0000-0000-000000000006', '30000000-0000-0000-0000-000000000002', '71000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000003',
   '7c542e71-8be4-4a5f-91ee-a1fb013a4063', 'scheduled', 'urgent', 'phone', 'partner',
   (CURRENT_DATE + interval '2 hours'), NULL, NULL, 'Mantenimiento correctivo urgente', NULL, NULL),
  ('80000000-0000-0000-0000-000000000007', '43000000-0000-0000-0000-000000000007', '30000000-0000-0000-0000-000000000003', '71000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000003',
   'facd0f93-c680-437d-aa98-42df62814701', 'completed', 'normal', 'email', 'customer',
   (CURRENT_DATE + interval '4 hours'), NULL, NULL, 'Inspeccion semestral', 'd38f7943-ef5d-4ad3-82b0-529c3e48ae91', NULL),
  ('80000000-0000-0000-0000-000000000008', '43000000-0000-0000-0000-000000000008', '30000000-0000-0000-0000-000000000004', '71000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000003',
   '45dc0986-5f43-40fb-a8a3-a26a41144aab', 'in_progress', 'emergency', 'phone', 'customer',
   (CURRENT_DATE + interval '0 hours'), NULL, NULL, 'Fuga de gas emergencia', NULL, NULL),
  ('80000000-0000-0000-0000-000000000009', '43000000-0000-0000-0000-000000000009', '30000000-0000-0000-0000-000000000001', '71000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000004',
   '40e1f2b8-4dda-4af1-92ff-1aa07e802010', 'scheduled', 'normal', 'portal', 'customer',
   (CURRENT_DATE + interval '24 hours'), NULL, NULL, 'Instalacion tanque nuevo', NULL, NULL),
  ('80000000-0000-0000-0000-000000000010', '43000000-0000-0000-0000-000000000010', '30000000-0000-0000-0000-000000000002', '71000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001',
   '7c542e71-8be4-4a5f-91ee-a1fb013a4063', 'completed', 'normal', 'portal', 'customer',
   (CURRENT_DATE + interval '6 hours'), NULL, NULL, 'Mantenimiento anual', 'd38f7943-ef5d-4ad3-82b0-529c3e48ae91', NULL),
  ('80000000-0000-0000-0000-000000000011', '43000000-0000-0000-0000-000000000011', '30000000-0000-0000-0000-000000000003', '71000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000002',
   '29097af9-1645-4ce5-b54f-d042b431487a', 'completed', 'normal', 'whatsapp', 'partner',
   (CURRENT_DATE + interval '1 hours'), NULL, NULL, 'Previsita condicional', '635296cb-38d9-48b0-8816-31aa2a85705d', NULL),
  ('80000000-0000-0000-0000-000000000012', '43000000-0000-0000-0000-000000000012', '30000000-0000-0000-0000-000000000004', '71000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000003',
   'facd0f93-c680-437d-aa98-42df62814701', 'scheduled', 'urgent', 'email', 'customer',
   (CURRENT_DATE + interval '12 hours'), NULL, NULL, 'Inspeccion post-siniestro', NULL, NULL),
  ('80000000-0000-0000-0000-000000000013', '43000000-0000-0000-0000-000000000013', '30000000-0000-0000-0000-000000000001', '71000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000004',
   '40e1f2b8-4dda-4af1-92ff-1aa07e802010', 'en_route', 'normal', 'portal', 'customer',
   (CURRENT_DATE + interval '1 hours'), NULL, NULL, 'Instalacion en transito', NULL, NULL),
  ('80000000-0000-0000-0000-000000000014', '43000000-0000-0000-0000-000000000014', '30000000-0000-0000-0000-000000000002', '71000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001',
   '7c542e71-8be4-4a5f-91ee-a1fb013a4063', 'scheduled', 'low', 'portal', 'customer',
   (CURRENT_DATE + interval '48 hours'), NULL, NULL, 'Mantenimiento preventivo menor', NULL, NULL),
  ('80000000-0000-0000-0000-000000000015', '43000000-0000-0000-0000-000000000015', '30000000-0000-0000-0000-000000000003', '71000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000002',
   '45dc0986-5f43-40fb-a8a3-a26a41144aab', 'completed', 'emergency', 'phone', 'customer',
   (CURRENT_DATE + interval '0 hours'), NULL, NULL, 'Emergencia resuelta', 'd38f7943-ef5d-4ad3-82b0-529c3e48ae91', NULL),
  ('80000000-0000-0000-0000-000000000016', '43000000-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000004', '71000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000003',
   '29097af9-1645-4ce5-b54f-d042b431487a', 'rescheduled', 'normal', 'portal', 'customer',
   (CURRENT_DATE + interval '24 hours'), NULL, NULL, 'Previsita reprogramada', NULL, NULL),
  ('80000000-0000-0000-0000-000000000017', '43000000-0000-0000-0000-000000000002', '30000000-0000-0000-0000-000000000001', '71000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000004',
   'facd0f93-c680-437d-aa98-42df62814701', 'scheduled', 'normal', 'email', 'customer',
   (CURRENT_DATE + interval '36 hours'), NULL, NULL, 'Inspeccion trimestral', NULL, NULL),
  ('80000000-0000-0000-0000-000000000018', '43000000-0000-0000-0000-000000000003', '30000000-0000-0000-0000-000000000002', '71000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001',
   '40e1f2b8-4dda-4af1-92ff-1aa07e802010', 'completed', 'normal', 'portal', 'customer',
   (CURRENT_DATE + interval '10 hours'), NULL, NULL, 'Instalacion verificada', 'd38f7943-ef5d-4ad3-82b0-529c3e48ae91', NULL),
  ('80000000-0000-0000-0000-000000000019', '43000000-0000-0000-0000-000000000004', '30000000-0000-0000-0000-000000000003', '71000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000002',
   '7c542e71-8be4-4a5f-91ee-a1fb013a4063', 'scheduled', 'normal', 'whatsapp', 'partner',
   (CURRENT_DATE + interval '6 hours'), NULL, NULL, 'Mantenimiento programado', NULL, NULL)
ON CONFLICT (id) DO UPDATE SET
  property_id = EXCLUDED.property_id, partner_id = EXCLUDED.partner_id,
  route_id = EXCLUDED.route_id, daily_plan_id = EXCLUDED.daily_plan_id,
  type = EXCLUDED.type, status = EXCLUDED.status, priority = EXCLUDED.priority,
  source = EXCLUDED.source, billing_to = EXCLUDED.billing_to,
  scheduled_at = EXCLUDED.scheduled_at, notes = EXCLUDED.notes,
  result = EXCLUDED.result, rejection_reason_id = EXCLUDED.rejection_reason_id;

-- PHASE 18: VISIT DATA

-- Visit assignments
INSERT INTO operations.visit_assignments (id, visit_id, tech_id, role_id) VALUES
  ('81000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000001', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '60000000-0000-0000-0000-000000000006', '9811bfc2-6644-493c-9cdf-48924a48308a'),
  ('81000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000004', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000002', '60000000-0000-0000-0000-000000000007', '9811bfc2-6644-493c-9cdf-48924a48308a'),
  ('81000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000009', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000003', '60000000-0000-0000-0000-000000000010', 'eadc5f0e-99f5-429d-b0f0-592cb6c0c24d'),
  ('81000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000004', '60000000-0000-0000-0000-000000000001', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000008', '80000000-0000-0000-0000-000000000005', '60000000-0000-0000-0000-000000000013', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000009', '80000000-0000-0000-0000-000000000006', '60000000-0000-0000-0000-000000000002', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000010', '80000000-0000-0000-0000-000000000006', '60000000-0000-0000-0000-000000000006', '9811bfc2-6644-493c-9cdf-48924a48308a'),
  ('81000000-0000-0000-0000-000000000011', '80000000-0000-0000-0000-000000000007', '60000000-0000-0000-0000-000000000005', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000012', '80000000-0000-0000-0000-000000000008', '60000000-0000-0000-0000-000000000001', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000013', '80000000-0000-0000-0000-000000000008', '60000000-0000-0000-0000-000000000014', '9811bfc2-6644-493c-9cdf-48924a48308a'),
  ('81000000-0000-0000-0000-000000000014', '80000000-0000-0000-0000-000000000008', '60000000-0000-0000-0000-000000000017', 'eadc5f0e-99f5-429d-b0f0-592cb6c0c24d'),
  ('81000000-0000-0000-0000-000000000015', '80000000-0000-0000-0000-000000000009', '60000000-0000-0000-0000-000000000008', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000016', '80000000-0000-0000-0000-000000000010', '60000000-0000-0000-0000-000000000011', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000017', '80000000-0000-0000-0000-000000000010', '60000000-0000-0000-0000-000000000015', '9811bfc2-6644-493c-9cdf-48924a48308a'),
  ('81000000-0000-0000-0000-000000000018', '80000000-0000-0000-0000-000000000011', '60000000-0000-0000-0000-000000000012', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000019', '80000000-0000-0000-0000-000000000012', '60000000-0000-0000-0000-000000000004', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000020', '80000000-0000-0000-0000-000000000012', '60000000-0000-0000-0000-000000000007', 'eadc5f0e-99f5-429d-b0f0-592cb6c0c24d'),
  ('81000000-0000-0000-0000-000000000021', '80000000-0000-0000-0000-000000000013', '60000000-0000-0000-0000-000000000016', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000022', '80000000-0000-0000-0000-000000000014', '60000000-0000-0000-0000-000000000010', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000023', '80000000-0000-0000-0000-000000000015', '60000000-0000-0000-0000-000000000001', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000024', '80000000-0000-0000-0000-000000000015', '60000000-0000-0000-0000-000000000006', '9811bfc2-6644-493c-9cdf-48924a48308a'),
  ('81000000-0000-0000-0000-000000000025', '80000000-0000-0000-0000-000000000015', '60000000-0000-0000-0000-000000000017', 'eadc5f0e-99f5-429d-b0f0-592cb6c0c24d'),
  ('81000000-0000-0000-0000-000000000026', '80000000-0000-0000-0000-000000000016', '60000000-0000-0000-0000-000000000003', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000027', '80000000-0000-0000-0000-000000000017', '60000000-0000-0000-0000-000000000009', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000028', '80000000-0000-0000-0000-000000000017', '60000000-0000-0000-0000-000000000011', '9811bfc2-6644-493c-9cdf-48924a48308a'),
  ('81000000-0000-0000-0000-000000000029', '80000000-0000-0000-0000-000000000018', '60000000-0000-0000-0000-000000000002', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000030', '80000000-0000-0000-0000-000000000018', '60000000-0000-0000-0000-000000000007', '9811bfc2-6644-493c-9cdf-48924a48308a'),
  ('81000000-0000-0000-0000-000000000031', '80000000-0000-0000-0000-000000000019', '60000000-0000-0000-0000-000000000005', '224b0832-b719-433f-8098-886282e4521d'),
  ('81000000-0000-0000-0000-000000000032', '80000000-0000-0000-0000-000000000019', '60000000-0000-0000-0000-000000000015', 'eadc5f0e-99f5-429d-b0f0-592cb6c0c24d')
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, tech_id = EXCLUDED.tech_id, role_id = EXCLUDED.role_id;

-- Visit checklist materials
INSERT INTO operations.visit_checklist_materials (id, visit_id, material_id, planned_quantity, confirmed) VALUES
  ('82000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000001', 1, true),
  ('82000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '50000000-0000-0000-0000-000000000003', 1, true),
  ('82000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000007', 1, true),
  ('82000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000002', '50000000-0000-0000-0000-000000000008', 1, true),
  ('82000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000001', 1, true),
  ('82000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000003', 1, true),
  ('82000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000007', 1, true),
  ('82000000-0000-0000-0000-000000000008', '80000000-0000-0000-0000-000000000005', '50000000-0000-0000-0000-000000000001', 1, true),
  ('82000000-0000-0000-0000-000000000009', '80000000-0000-0000-0000-000000000006', '50000000-0000-0000-0000-000000000007', 1, true),
  ('82000000-0000-0000-0000-000000000010', '80000000-0000-0000-0000-000000000006', '50000000-0000-0000-0000-000000000008', 1, true),
  ('82000000-0000-0000-0000-000000000011', '80000000-0000-0000-0000-000000000007', '50000000-0000-0000-0000-000000000001', 1, true),
  ('82000000-0000-0000-0000-000000000012', '80000000-0000-0000-0000-000000000007', '50000000-0000-0000-0000-000000000003', 1, true),
  ('82000000-0000-0000-0000-000000000013', '80000000-0000-0000-0000-000000000008', '50000000-0000-0000-0000-000000000007', 1, true),
  ('82000000-0000-0000-0000-000000000014', '80000000-0000-0000-0000-000000000009', '50000000-0000-0000-0000-000000000001', 1, true),
  ('82000000-0000-0000-0000-000000000015', '80000000-0000-0000-0000-000000000009', '50000000-0000-0000-0000-000000000003', 1, true),
  ('82000000-0000-0000-0000-000000000016', '80000000-0000-0000-0000-000000000009', '50000000-0000-0000-0000-000000000007', 1, true),
  ('82000000-0000-0000-0000-000000000017', '80000000-0000-0000-0000-000000000010', '50000000-0000-0000-0000-000000000007', 1, true),
  ('82000000-0000-0000-0000-000000000018', '80000000-0000-0000-0000-000000000010', '50000000-0000-0000-0000-000000000008', 1, true),
  ('82000000-0000-0000-0000-000000000019', '80000000-0000-0000-0000-000000000011', '50000000-0000-0000-0000-000000000001', 1, false),
  ('82000000-0000-0000-0000-000000000020', '80000000-0000-0000-0000-000000000013', '50000000-0000-0000-0000-000000000001', 1, false),
  ('82000000-0000-0000-0000-000000000021', '80000000-0000-0000-0000-000000000013', '50000000-0000-0000-0000-000000000003', 1, false),
  ('82000000-0000-0000-0000-000000000022', '80000000-0000-0000-0000-000000000015', '50000000-0000-0000-0000-000000000007', 1, false),
  ('82000000-0000-0000-0000-000000000023', '80000000-0000-0000-0000-000000000015', '50000000-0000-0000-0000-000000000008', 1, false),
  ('82000000-0000-0000-0000-000000000024', '80000000-0000-0000-0000-000000000017', '50000000-0000-0000-0000-000000000001', 1, false),
  ('82000000-0000-0000-0000-000000000025', '80000000-0000-0000-0000-000000000017', '50000000-0000-0000-0000-000000000003', 1, false),
  ('82000000-0000-0000-0000-000000000026', '80000000-0000-0000-0000-000000000018', '50000000-0000-0000-0000-000000000007', 1, false),
  ('82000000-0000-0000-0000-000000000027', '80000000-0000-0000-0000-000000000018', '50000000-0000-0000-0000-000000000008', 1, false),
  ('82000000-0000-0000-0000-000000000028', '80000000-0000-0000-0000-000000000019', '50000000-0000-0000-0000-000000000001', 1, false),
  ('82000000-0000-0000-0000-000000000029', '80000000-0000-0000-0000-000000000019', '50000000-0000-0000-0000-000000000003', 1, false)
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, material_id = EXCLUDED.material_id, planned_quantity = EXCLUDED.planned_quantity, confirmed = EXCLUDED.confirmed;

-- Visit checklist tools
INSERT INTO operations.visit_checklist_tools (id, visit_id, tool_id, planned_quantity, confirmed) VALUES
  ('83000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000001', 1, true),
  ('83000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '51000000-0000-0000-0000-000000000005', 1, true),
  ('83000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000002', '51000000-0000-0000-0000-000000000002', 1, true),
  ('83000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000003', '51000000-0000-0000-0000-000000000001', 1, true),
  ('83000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000003', '51000000-0000-0000-0000-000000000005', 1, true),
  ('83000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000006', '51000000-0000-0000-0000-000000000002', 1, true),
  ('83000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000007', '51000000-0000-0000-0000-000000000001', 1, true),
  ('83000000-0000-0000-0000-000000000008', '80000000-0000-0000-0000-000000000008', '51000000-0000-0000-0000-000000000002', 1, true),
  ('83000000-0000-0000-0000-000000000009', '80000000-0000-0000-0000-000000000009', '51000000-0000-0000-0000-000000000001', 1, true),
  ('83000000-0000-0000-0000-000000000010', '80000000-0000-0000-0000-000000000010', '51000000-0000-0000-0000-000000000005', 1, true),
  ('83000000-0000-0000-0000-000000000011', '80000000-0000-0000-0000-000000000010', '51000000-0000-0000-0000-000000000002', 1, true),
  ('83000000-0000-0000-0000-000000000012', '80000000-0000-0000-0000-000000000013', '51000000-0000-0000-0000-000000000001', 1, false),
  ('83000000-0000-0000-0000-000000000013', '80000000-0000-0000-0000-000000000015', '51000000-0000-0000-0000-000000000002', 1, false),
  ('83000000-0000-0000-0000-000000000014', '80000000-0000-0000-0000-000000000018', '51000000-0000-0000-0000-000000000001', 1, false),
  ('83000000-0000-0000-0000-000000000015', '80000000-0000-0000-0000-000000000018', '51000000-0000-0000-0000-000000000005', 1, false)
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, tool_id = EXCLUDED.tool_id, planned_quantity = EXCLUDED.planned_quantity, confirmed = EXCLUDED.confirmed;

-- Visit checklist epps
INSERT INTO operations.visit_checklist_epps (id, visit_id, epp_id, planned_quantity, confirmed) VALUES
  ('84000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000001', 1, true),
  ('84000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000002', 1, true),
  ('84000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000002', '52000000-0000-0000-0000-000000000003', 1, true),
  ('84000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000002', '52000000-0000-0000-0000-000000000004', 1, true),
  ('84000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000003', '52000000-0000-0000-0000-000000000001', 1, true),
  ('84000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000006', '52000000-0000-0000-0000-000000000003', 1, true),
  ('84000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000008', '52000000-0000-0000-0000-000000000001', 1, true),
  ('84000000-0000-0000-0000-000000000008', '80000000-0000-0000-0000-000000000008', '52000000-0000-0000-0000-000000000002', 1, true),
  ('84000000-0000-0000-0000-000000000009', '80000000-0000-0000-0000-000000000008', '52000000-0000-0000-0000-000000000003', 1, true),
  ('84000000-0000-0000-0000-000000000010', '80000000-0000-0000-0000-000000000008', '52000000-0000-0000-0000-000000000004', 1, true),
  ('84000000-0000-0000-0000-000000000011', '80000000-0000-0000-0000-000000000009', '52000000-0000-0000-0000-000000000001', 1, true),
  ('84000000-0000-0000-0000-000000000012', '80000000-0000-0000-0000-000000000010', '52000000-0000-0000-0000-000000000003', 1, true),
  ('84000000-0000-0000-0000-000000000013', '80000000-0000-0000-0000-000000000013', '52000000-0000-0000-0000-000000000001', 1, false),
  ('84000000-0000-0000-0000-000000000014', '80000000-0000-0000-0000-000000000015', '52000000-0000-0000-0000-000000000001', 1, false),
  ('84000000-0000-0000-0000-000000000015', '80000000-0000-0000-0000-000000000015', '52000000-0000-0000-0000-000000000002', 1, false),
  ('84000000-0000-0000-0000-000000000016', '80000000-0000-0000-0000-000000000018', '52000000-0000-0000-0000-000000000001', 1, false)
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, epp_id = EXCLUDED.epp_id, planned_quantity = EXCLUDED.planned_quantity, confirmed = EXCLUDED.confirmed;

-- Visit material usages
INSERT INTO operations.visit_material_usages (id, visit_id, material_id, quantity, notes) VALUES
  ('85000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000001', 3, 'Tubos PVC instalacion'),
  ('85000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000003', '50000000-0000-0000-0000-000000000002', 2, 'Abrazaderas'),
  ('85000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000007', '50000000-0000-0000-0000-000000000001', 1, 'Sellador'),
  ('85000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000010', '50000000-0000-0000-0000-000000000003', 5, 'Fittings varios'),
  ('85000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000015', '50000000-0000-0000-0000-000000000004', 4, 'Reparacion emergencia'),
  ('85000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000018', '50000000-0000-0000-0000-000000000001', 2, 'Material instalacion')
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, material_id = EXCLUDED.material_id, quantity = EXCLUDED.quantity;

-- Visit tool usages
INSERT INTO operations.visit_tool_usages (id, visit_id, tool_id, notes) VALUES
  ('86000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000003', '51000000-0000-0000-0000-000000000001', 'Llave inglesa para uniones'),
  ('86000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000007', '51000000-0000-0000-0000-000000000003', 'Detector de fugas'),
  ('86000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000010', '51000000-0000-0000-0000-000000000002', 'Llave especial'),
  ('86000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000015', '51000000-0000-0000-0000-000000000001', 'Herramienta emergencia'),
  ('86000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000018', '51000000-0000-0000-0000-000000000004', 'Juego de llaves')
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, tool_id = EXCLUDED.tool_id, notes = EXCLUDED.notes;

-- Visit EPP usages
INSERT INTO operations.visit_epp_usages (id, visit_id, epp_id, quantity, status, notes) VALUES
  ('87000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000001', 1, 'used', 'Casco'),
  ('87000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', '52000000-0000-0000-0000-000000000002', 1, 'used', 'Guantes'),
  ('87000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000002', '52000000-0000-0000-0000-000000000003', 1, 'used', 'Detector gas'),
  ('87000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000003', '52000000-0000-0000-0000-000000000001', 1, 'used', 'Casco'),
  ('87000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000006', '52000000-0000-0000-0000-000000000004', 1, 'not_needed', 'N/A'),
  ('87000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000008', '52000000-0000-0000-0000-000000000001', 2, 'used', 'Casco emergencia'),
  ('87000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000008', '52000000-0000-0000-0000-000000000002', 2, 'used', 'Guantes emergencia'),
  ('87000000-0000-0000-0000-000000000008', '80000000-0000-0000-0000-000000000010', '52000000-0000-0000-0000-000000000003', 1, 'used', 'Detector'),
  ('87000000-0000-0000-0000-000000000009', '80000000-0000-0000-0000-000000000015', '52000000-0000-0000-0000-000000000001', 1, 'missing', 'Casco perdido'),
  ('87000000-0000-0000-0000-000000000010', '80000000-0000-0000-0000-000000000018', '52000000-0000-0000-0000-000000000002', 1, 'used', 'Guantes')
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, epp_id = EXCLUDED.epp_id, quantity = EXCLUDED.quantity, status = EXCLUDED.status;

-- Visit checkpoints
INSERT INTO operations.visit_checkpoints (id, visit_id, type, geom, timestamp, device_info) VALUES
  ('88000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', 'arrival', ST_SetSRID(ST_MakePoint(-70.65, -33.45), 4326), (CURRENT_DATE + interval '8 hours'), 'Android 13'),
  ('88000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000001', 'departure', ST_SetSRID(ST_MakePoint(-70.65, -33.45), 4326), (CURRENT_DATE + interval '10 hours'), 'Android 13'),
  ('88000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000002', 'arrival', ST_SetSRID(ST_MakePoint(-70.66, -33.44), 4326), (CURRENT_DATE + interval '9 hours'), 'iOS 17'),
  ('88000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000003', 'arrival', ST_SetSRID(ST_MakePoint(-71.23, -33.02), 4326), (CURRENT_DATE + interval '14 hours'), 'Android 14'),
  ('88000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000003', 'departure', ST_SetSRID(ST_MakePoint(-71.23, -33.02), 4326), (CURRENT_DATE + interval '17 hours'), 'Android 14'),
  ('88000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000005', 'arrival', ST_SetSRID(ST_MakePoint(-70.64, -33.46), 4326), (CURRENT_DATE + interval '10 hours'), 'Android 13'),
  ('88000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000005', 'report_sent', ST_SetSRID(ST_MakePoint(-70.64, -33.46), 4326), (CURRENT_DATE + interval '11 hours'), 'Android 13'),
  ('88000000-0000-0000-0000-000000000008', '80000000-0000-0000-0000-000000000007', 'arrival', ST_SetSRID(ST_MakePoint(-71.25, -33.03), 4326), (CURRENT_DATE + interval '8 hours'), 'iOS 17'),
  ('88000000-0000-0000-0000-000000000009', '80000000-0000-0000-0000-000000000008', 'arrival', ST_SetSRID(ST_MakePoint(-70.67, -33.43), 4326), (CURRENT_DATE + interval '2 hours'), 'Android 13'),
  ('88000000-0000-0000-0000-000000000010', '80000000-0000-0000-0000-000000000010', 'arrival', ST_SetSRID(ST_MakePoint(-70.63, -33.47), 4326), (CURRENT_DATE + interval '9 hours'), 'Android 14'),
  ('88000000-0000-0000-0000-000000000011', '80000000-0000-0000-0000-000000000010', 'departure', ST_SetSRID(ST_MakePoint(-70.63, -33.47), 4326), (CURRENT_DATE + interval '11 hours'), 'Android 14'),
  ('88000000-0000-0000-0000-000000000012', '80000000-0000-0000-0000-000000000015', 'arrival', ST_SetSRID(ST_MakePoint(-70.65, -33.45), 4326), (CURRENT_DATE + interval '3 hours'), 'Android 13'),
  ('88000000-0000-0000-0000-000000000013', '80000000-0000-0000-0000-000000000015', 'report_sent', ST_SetSRID(ST_MakePoint(-70.65, -33.45), 4326), (CURRENT_DATE + interval '5 hours'), 'Android 13'),
  ('88000000-0000-0000-0000-000000000014', '80000000-0000-0000-0000-000000000018', 'arrival', ST_SetSRID(ST_MakePoint(-70.64, -33.46), 4326), (CURRENT_DATE + interval '14 hours'), 'iOS 17')
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, type = EXCLUDED.type, timestamp = EXCLUDED.timestamp, device_info = EXCLUDED.device_info;

-- Visit measurements
INSERT INTO operations.visit_measurements (id, visit_id, type, value, unit, result, measuring_device, notes) VALUES
  ('89000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000002', NULL, 95.5, '%', 'approved', 'Manometro digital', 'Presion OK'),
  ('89000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000003', NULL, 99.1, '%', 'approved', 'Manometro digital', 'Hermeticidad OK'),
  ('89000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000007', NULL, 97.8, '%', 'approved', 'Manometro digital', 'Inspeccion OK'),
  ('89000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000008', NULL, 72.3, '%', 'observation', 'Manometro digital', 'Presion baja detectada'),
  ('89000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000010', NULL, 98.5, '%', 'approved', 'Manometro digital', 'Mantenimiento OK'),
  ('89000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000015', NULL, 88.0, '%', 'pending', 'Manometro digital', 'Requiere seguimiento'),
  ('89000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000018', NULL, 99.3, '%', 'approved', 'Manometro digital', 'Instalacion verificada')
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, value = EXCLUDED.value, unit = EXCLUDED.unit, result = EXCLUDED.result;

-- Visit photos
INSERT INTO operations.visit_photos (id, visit_id, checkpoint_id, url, geom, timestamp, stage, finding_type) VALUES
  ('90000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000003', NULL, '/photos/v3_estanque.jpg', ST_SetSRID(ST_MakePoint(-71.23, -33.02), 4326), (CURRENT_DATE + interval '14 hours'), 'execution', '2bff71a3-6de9-4643-afb9-3d8f168f02ba'),
  ('90000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000003', NULL, '/photos/v3_soldadura.jpg', ST_SetSRID(ST_MakePoint(-71.23, -33.02), 4326), (CURRENT_DATE + interval '15 hours'), 'execution', '2bff71a3-6de9-4643-afb9-3d8f168f02ba'),
  ('90000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000007', NULL, '/photos/v7_inspeccion.jpg', ST_SetSRID(ST_MakePoint(-71.25, -33.03), 4326), (CURRENT_DATE + interval '9 hours'), 'diagnosis', '2bff71a3-6de9-4643-afb9-3d8f168f02ba'),
  ('90000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000008', NULL, '/photos/v8_fuga.jpg', ST_SetSRID(ST_MakePoint(-70.67, -33.43), 4326), (CURRENT_DATE + interval '3 hours'), 'diagnosis', 'c280c476-386b-48f2-ad21-a3bd63067189'),
  ('90000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000010', NULL, '/photos/v10_mantencion.jpg', ST_SetSRID(ST_MakePoint(-70.63, -33.47), 4326), (CURRENT_DATE + interval '10 hours'), 'execution', '2bff71a3-6de9-4643-afb9-3d8f168f02ba'),
  ('90000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000015', NULL, '/photos/v15_emergencia.jpg', ST_SetSRID(ST_MakePoint(-70.65, -33.45), 4326), (CURRENT_DATE + interval '4 hours'), 'execution', 'c280c476-386b-48f2-ad21-a3bd63067189'),
  ('90000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000018', NULL, '/photos/v18_verificacion.jpg', ST_SetSRID(ST_MakePoint(-70.64, -33.46), 4326), (CURRENT_DATE + interval '15 hours'), 'review', '2bff71a3-6de9-4643-afb9-3d8f168f02ba')
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, url = EXCLUDED.url, stage = EXCLUDED.stage;

-- Visit reports
INSERT INTO operations.visit_reports (id, visit_id, report_template_id, source, recorded_at) VALUES
  ('8a000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001', 'system', (CURRENT_DATE + interval '17 hours')),
  ('8a000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000002', 'paper', (CURRENT_DATE + interval '10 hours')),
  ('8a000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000010', '10000000-0000-0000-0000-000000000001', 'system', (CURRENT_DATE + interval '11 hours')),
  ('8a000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000015', '10000000-0000-0000-0000-000000000003', 'system', (CURRENT_DATE + interval '5 hours')),
  ('8a000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000018', '10000000-0000-0000-0000-000000000001', 'system', (CURRENT_DATE + interval '16 hours'))
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, report_template_id = EXCLUDED.report_template_id, source = EXCLUDED.source;

-- Report entries
INSERT INTO operations.report_entries (id, report_id, data_json) VALUES
  ('8b000000-0000-0000-0000-000000000001', '8a000000-0000-0000-0000-000000000001', '{"status":"completed","checklist":true,"signature":true}'::jsonb),
  ('8b000000-0000-0000-0000-000000000002', '8a000000-0000-0000-0000-000000000002', '{"status":"completed","notes":"Inspeccion satisfactoria"}'::jsonb),
  ('8b000000-0000-0000-0000-000000000003', '8a000000-0000-0000-0000-000000000003', '{"status":"completed","checklist":true,"signature":true}'::jsonb),
  ('8b000000-0000-0000-0000-000000000004', '8a000000-0000-0000-0000-000000000004', '{"status":"completed","incident":true,"photos":2}'::jsonb),
  ('8b000000-0000-0000-0000-000000000005', '8a000000-0000-0000-0000-000000000005', '{"status":"completed","checklist":true}'::jsonb)
ON CONFLICT (id) DO UPDATE SET report_id = EXCLUDED.report_id, data_json = EXCLUDED.data_json;

-- Vehicle assignments
INSERT INTO operations.vehicle_assignments (id, vehicle_id, daily_plan_id, visit_id, departure_time, return_time, departure_mileage, return_mileage, notes) VALUES
  ('8c000000-0000-0000-0000-000000000001', '53000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '2 hours'), 12345, 12365, 'Vehiculo Santiago sur'),
  ('8c000000-0000-0000-0000-000000000002', '53000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000002', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '3 hours'), 23456, 23480, 'Vehiculo Santiago norte'),
  ('8c000000-0000-0000-0000-000000000003', '53000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000003', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '4 hours'), 34567, 34600, 'Vehiculo Talca'),
  ('8c000000-0000-0000-0000-000000000004', '53000000-0000-0000-0000-000000000004', '70000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000006', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '3 hours'), 45678, 45710, 'Vehiculo Curico'),
  ('8c000000-0000-0000-0000-000000000005', '53000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000008', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '2 hours'), 12370, 12400, 'Vehiculo emergencia'),
  ('8c000000-0000-0000-0000-000000000006', '53000000-0000-0000-0000-000000000005', '70000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000010', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '3 hours'), 56789, 56815, 'Vehiculo Santiago centro'),
  ('8c000000-0000-0000-0000-000000000007', '53000000-0000-0000-0000-000000000006', '70000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000013', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '2 hours'), 67890, 67910, 'Vehiculo en ruta'),
  ('8c000000-0000-0000-0000-000000000008', '53000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000015', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '3 hours'), 12410, 12440, 'Vehiculo emergencia 2'),
  ('8c000000-0000-0000-0000-000000000009', '53000000-0000-0000-0000-000000000007', '70000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000018', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '3 hours'), 78901, 78930, 'Vehiculo verificacion')
ON CONFLICT (id) DO UPDATE SET vehicle_id = EXCLUDED.vehicle_id, visit_id = EXCLUDED.visit_id;

-- Visit SLA trackings
INSERT INTO operations.visit_sla_trackings (id, visit_id, sla_id, requested_at, responded_at, resolved_at, response_time_hours, resolution_time_hours, meets_response_sla, meets_resolution_sla) VALUES
  ('8d000000-0000-0000-0000-000000000001', '80000000-0000-0000-0000-000000000001', '35000000-0000-0000-0000-000000000001', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '2.5 hours'), (CURRENT_DATE + interval '8.0 hours'), 2.5, 8.0, TRUE, TRUE),
  ('8d000000-0000-0000-0000-000000000002', '80000000-0000-0000-0000-000000000002', '35000000-0000-0000-0000-000000000001', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '1.0 hours'), (CURRENT_DATE + interval '6.0 hours'), 1.0, 6.0, TRUE, TRUE),
  ('8d000000-0000-0000-0000-000000000003', '80000000-0000-0000-0000-000000000003', '35000000-0000-0000-0000-000000000001', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '3.0 hours'), (CURRENT_DATE + interval '12.0 hours'), 3.0, 12.0, TRUE, TRUE),
  ('8d000000-0000-0000-0000-000000000004', '80000000-0000-0000-0000-000000000006', '35000000-0000-0000-0000-000000000002', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '0.5 hours'), (CURRENT_DATE + interval '2.0 hours'), 0.5, 2.0, TRUE, TRUE),
  ('8d000000-0000-0000-0000-000000000005', '80000000-0000-0000-0000-000000000008', '35000000-0000-0000-0000-000000000002', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '0.25 hours'), (CURRENT_DATE + interval '1.5 hours'), 0.25, 1.5, TRUE, TRUE),
  ('8d000000-0000-0000-0000-000000000006', '80000000-0000-0000-0000-000000000010', '35000000-0000-0000-0000-000000000001', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '4.0 hours'), (CURRENT_DATE + interval '24.0 hours'), 4.0, 24.0, TRUE, TRUE),
  ('8d000000-0000-0000-0000-000000000007', '80000000-0000-0000-0000-000000000012', '35000000-0000-0000-0000-000000000002', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '1.0 hours'), NULL, 1.0, NULL, TRUE, NULL),
  ('8d000000-0000-0000-0000-000000000008', '80000000-0000-0000-0000-000000000015', '35000000-0000-0000-0000-000000000002', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '0.5 hours'), (CURRENT_DATE + interval '3.0 hours'), 0.5, 3.0, TRUE, TRUE),
  ('8d000000-0000-0000-0000-000000000009', '80000000-0000-0000-0000-000000000018', '35000000-0000-0000-0000-000000000001', (CURRENT_DATE + interval '0 hours'), (CURRENT_DATE + interval '2.0 hours'), (CURRENT_DATE + interval '10.0 hours'), 2.0, 10.0, TRUE, TRUE)
ON CONFLICT (id) DO UPDATE SET visit_id = EXCLUDED.visit_id, sla_id = EXCLUDED.sla_id;

-- PHASE 19: AUDIT LOGS

-- Auth audit log
INSERT INTO core.auth_audit_log (id, user_id, role, ip, endpoint, method, status_code, timestamp) VALUES
  ('aa000000-0000-0000-0000-000000000001', NULL, 'admin', '192.168.1.10', '/api/auth/login', 'POST', 200, (CURRENT_DATE + interval '0 hours')),
  ('aa000000-0000-0000-0000-000000000002', NULL, 'admin', '192.168.1.10', '/api/visits', 'GET', 200, (CURRENT_DATE + interval '1 hours')),
  ('aa000000-0000-0000-0000-000000000003', NULL, 'tech', '10.0.0.5', '/api/visits/80000000-0000-0000-0000-000000000003/checkpoints', 'POST', 201, (CURRENT_DATE + interval '2 hours')),
  ('aa000000-0000-0000-0000-000000000004', NULL, 'tech', '10.0.0.5', '/api/visits/80000000-0000-0000-0000-000000000003/measurements', 'POST', 201, (CURRENT_DATE + interval '4 hours')),
  ('aa000000-0000-0000-0000-000000000005', NULL, 'admin', '192.168.1.10', '/api/auth/logout', 'POST', 200, (CURRENT_DATE + interval '8 hours'))
ON CONFLICT (id) DO NOTHING;

-- Plan audit log
INSERT INTO core.plan_audit_log (id, entity_type, entity_id, action, old_data, new_data, user_id, reason, timestamp) VALUES
  ('aa000000-0000-0000-0000-000000000001', 'visit_checklist_material', '82000000-0000-0000-0000-000000000001', 'update', '{"confirmed":false}'::jsonb, '{"confirmed":true}'::jsonb, NULL, 'Visita completada', (CURRENT_DATE + interval '18 hours')),
  ('aa000000-0000-0000-0000-000000000002', 'daily_load_material', '70000000-0000-0000-0000-000000000001', 'update', '{"status":"scheduled"}'::jsonb, '{"status":"in_progress"}'::jsonb, NULL, 'Plan iniciado', (CURRENT_DATE + interval '0 hours')),
  ('aa000000-0000-0000-0000-000000000003', 'visit_checklist_material', '80000000-0000-0000-0000-000000000008', 'update', '{"status":"scheduled"}'::jsonb, '{"status":"in_progress"}'::jsonb, NULL, 'Emergencia atendida', (CURRENT_DATE + interval '2 hours')),
  ('aa000000-0000-0000-0000-000000000004', 'visit_material_usage', '85000000-0000-0000-0000-000000000001', 'add', NULL, '{"quantity":3}'::jsonb, NULL, 'Visita cancelada por cliente', (CURRENT_DATE + interval '8 hours')),
  ('aa000000-0000-0000-0000-000000000005', 'visit_epp_usage', '87000000-0000-0000-0000-000000000005', 'update', '{"status":"not_needed"}'::jsonb, '{"status":"used"}'::jsonb, NULL, 'Mantenimiento finalizado', (CURRENT_DATE + interval '6 hours'))
ON CONFLICT (id) DO NOTHING;

-- END OF SEED DATA

