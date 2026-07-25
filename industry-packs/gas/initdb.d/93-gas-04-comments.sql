-- ============================================================================
-- 04-domain-gas-comments.sql — Gas Chile industry pack comments for OpenAPI
-- ============================================================================

COMMENT ON SCHEMA domain_gas IS 'Gas Chile specific domain: lookup tables, SEC certifications, rejection reasons';
COMMENT ON TABLE domain_gas.visit_types IS 'Gas domain visit types (pre_visit, installation, etc.)';
COMMENT ON TABLE domain_gas.measurement_types IS 'Gas domain measurement types (tightness, pressure, CO, etc.)';
COMMENT ON TABLE domain_gas.photo_findings IS 'Gas domain photo finding types';
COMMENT ON TABLE domain_gas.vehicle_types IS 'Gas domain vehicle types';
COMMENT ON TABLE domain_gas.partner_service_types IS 'Gas domain partner service types';
COMMENT ON TABLE domain_gas.pre_visit_results IS 'Gas domain pre-visit results';
COMMENT ON TABLE domain_gas.route_types IS 'Gas domain route types (meter_reading, etc.)';
COMMENT ON TABLE domain_gas.property_types IS 'Gas domain property types with extra column config';
COMMENT ON TABLE domain_gas.certifications IS 'Official SEC Chile certification catalog';
COMMENT ON COLUMN domain_gas.certifications.code IS 'Code: CL1, CL2, CL3, GLP, TC1, TC2, TC6, GREEN SEAL, DS66';
COMMENT ON COLUMN domain_gas.certifications.class IS 'Classification: Class 1/2/3, TC, Seal, Standard';
COMMENT ON TABLE domain_gas.rejection_reasons IS 'Gas domain specific rejection reasons';
COMMENT ON COLUMN domain_gas.rejection_reasons.category IS 'Category: logistics, regulatory, customer, operational';
COMMENT ON TABLE domain_gas.property_assets IS 'Extra property attributes by type (poles, meter, etc.)';