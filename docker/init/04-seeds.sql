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