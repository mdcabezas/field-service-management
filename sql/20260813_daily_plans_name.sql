-- ============================================================================
-- 20260813_daily_plans_name.sql — Add name column to planning.daily_plans
-- Feature: display daily plan name in route editor (instead of UUID)
-- Idempotent: safe to run multiple times.
-- Compatible with PostgreSQL 16.
-- ============================================================================

BEGIN;

ALTER TABLE planning.daily_plans
  ADD COLUMN IF NOT EXISTS name VARCHAR(255) NOT NULL DEFAULT 'Plan sin nombre';

UPDATE planning.daily_plans
SET name = 'Plan ' || to_char(date, 'YYYY-MM-DD')
WHERE name = 'Plan sin nombre'
  AND date IS NOT NULL;

ALTER TABLE planning.daily_plans
  ALTER COLUMN name DROP DEFAULT;

COMMENT ON COLUMN planning.daily_plans.name IS 'Human-readable name of the daily plan';

COMMIT;
