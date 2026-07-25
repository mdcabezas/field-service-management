-- ============================================================================
-- 03-domain-gas-grants.sql — Gas Chile industry pack grants
-- ============================================================================

GRANT USAGE ON SCHEMA domain_gas TO fsm_api;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA domain_gas TO fsm_api;

-- Extend search_path for fsm_api to include domain_gas
ALTER ROLE fsm_api SET search_path TO core, partners, customers, inventory, operations, planning, notifications, geocoding, shared, domain_gas;
ALTER ROLE fsm_backend SET search_path TO core, partners, customers, inventory, operations, planning, notifications, geocoding, shared, domain_gas;