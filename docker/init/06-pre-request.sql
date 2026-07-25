-- ============================================================================
-- 06-pre-request.sql — Pre-request function for Traefik authentication
-- Sets user context based on headers injected by Traefik
-- ============================================================================

-- Main function: set_user_context
-- Called before each request via db-pre-request
-- Reads headers injected by Traefik: X-User-ID, X-User-Role
-- PostgREST stores headers as JSON: current_setting('request.headers', true)::json
CREATE OR REPLACE FUNCTION public.set_user_context()
RETURNS void AS $$
DECLARE
  v_headers json;
  v_role text;
  v_user_id text;
BEGIN
  -- Read headers as JSON from PostgREST transaction-scoped settings
  v_headers := current_setting('request.headers', true)::json;

  -- Extract X-User-ID and X-User-Role (lowercased by PostgREST)
  IF v_headers IS NOT NULL THEN
    v_role := v_headers->>'x-user-role';
    v_user_id := v_headers->>'x-user-id';
  END IF;

  -- If no headers (anonymous access via fsm_api), do nothing
  IF v_role IS NULL AND v_user_id IS NULL THEN
    RETURN;
  END IF;

  -- Propagate user_id for audit trail (app_user_id() reads this)
  IF v_user_id IS NOT NULL THEN
    PERFORM set_config('request.jwt.claim.sub', v_user_id, true);
  END IF;

  -- Propagate role for app-level access control
  IF v_role IS NOT NULL THEN
    PERFORM set_config('request.jwt.claim.role', v_role, true);
  END IF;

END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Grant execution to API roles
GRANT EXECUTE ON FUNCTION public.set_user_context() TO fsm_api, fsm_backend;
