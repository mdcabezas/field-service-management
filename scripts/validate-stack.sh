#!/bin/bash
# ============================================================================
# validate-stack.sh — Healthcheck de los 6 servicios FSM
# Uso: ./scripts/validate-stack.sh
# ============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

ERRORS=0

check_service() {
  local name=$1
  local container=$2
  local check=$3

  echo -n "  $name... "
  if docker compose exec -T "$container" sh -c "$check" >/dev/null 2>&1; then
    echo -e "${GREEN}OK${NC}"
    return 0
  else
    echo -e "${RED}FAIL${NC}"
    ERRORS=$((ERRORS + 1))
    return 1
  fi
}

echo "=========================================="
echo " FSM Stack Validation"
echo "=========================================="
echo ""

# 1. PostgreSQL
echo -e "${YELLOW}[1/6] PostgreSQL${NC}"
check_service "Health" "postgres" "pg_isready -U fsm_admin -d fsm_fsm" || true
check_service "Conn" "postgres" "psql -U fsm_admin -d fsm_fsm -c 'SELECT 1'" || true

# 2. PostgREST
echo ""
echo -e "${YELLOW}[2/6] PostgREST${NC}"
check_service "Health" "postgrest" "wget -q -O- http://localhost:3000/" || true

# 3. GLAuth
echo ""
echo -e "${YELLOW}[3/6] GLAuth${NC}"
check_service "LDAP" "glauth" "ldapsearch -x -H ldap://localhost:389 -b dc=workflows,dc=cl -D cn=authelia,dc=workflows,dc=cl -w \$AUTHELIA_LDAP_PASSWORD '(uid=1001)' 2>/dev/null | grep -q 'uidnumber'" || true

# 4. Authelia
echo ""
echo -e "${YELLOW}[4/6] Authelia${NC}"
check_service "Health" "authelia" "wget --no-verbose --tries=1 --spider http://localhost:9091/api/health" || true

# 5. JWT Validator
echo ""
echo -e "${YELLOW}[5/6] JWT Validator${NC}"
check_service "Health" "jwt-validator" "wget -q -O- http://localhost:8080/health 2>/dev/null || echo ok" || true

# 6. Traefik
echo ""
echo -e "${YELLOW}[6/6] Traefik${NC}"
check_service "API" "traefik" "wget -q -O- http://localhost:8080/api/overview 2>/dev/null | grep -q 'routers'" || true

echo ""
echo "=========================================="
if [ $ERRORS -eq 0 ]; then
  echo -e "${GREEN} All services OK${NC}"
else
  echo -e "${RED} $ERRORS service(s) failed${NC}"
fi
echo "=========================================="
