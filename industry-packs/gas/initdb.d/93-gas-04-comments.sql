-- ============================================================================
-- 04-domain-gas-comments.sql — Gas Chile industry pack comments for OpenAPI
-- ============================================================================

COMMENT ON SCHEMA domain_gas IS 'Gas Chile specific domain: measurement types, property types, SEC certifications, property assets';
COMMENT ON TABLE domain_gas.measurement_types IS 'Gas domain measurement types (tightness, pressure, CO, etc.)';
COMMENT ON TABLE domain_gas.property_types IS 'Gas domain property types with extra column config';
COMMENT ON TABLE domain_gas.certifications IS 'Official SEC Chile certification catalog';
COMMENT ON COLUMN domain_gas.certifications.code IS 'Code: CL1, CL2, CL3, GLP, TC1, TC2, TC6, GREEN SEAL, DS66';
COMMENT ON COLUMN domain_gas.certifications.class IS 'Classification: Class 1/2/3, TC, Seal, Standard';
COMMENT ON TABLE domain_gas.property_assets IS 'Extra property attributes by type (poles, meter, etc.)';