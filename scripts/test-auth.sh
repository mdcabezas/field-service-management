#!/bin/bash
# ============================================================================
# test-auth.sh — Test LDAP bind (GLAuth) + token generation (Authelia)
# Uso: ./scripts/test-auth.sh <user> <password>
# ============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

if [ $# -lt 2 ]; then
  echo -e "${RED}Error: user and password are required${NC}"
  echo "Usage: $0 <user> <password>"
  echo "Example: $0 admin <password>"
  exit 1
fi

USER=$1
PASSWORD=$2
BASE_DN="dc=workflows,dc=cl"
LDAP_URL="ldap://localhost:389"
GLAUTH_CONTAINER="fsm-glauth"

echo "=========================================="
echo " Auth Test"
echo "=========================================="
echo ""

# 1. Test LDAP bind
echo -e "${YELLOW}[1/2] LDAP Bind - GLAuth${NC}"
echo "  User: cn=${USER},${BASE_DN}"

if docker compose exec -T "$GLAUTH_CONTAINER" ldapsearch -x \
  -H "$LDAP_URL" \
  -b "$BASE_DN" \
  -D "cn=${USER},${BASE_DN}" \
  -w "$PASSWORD" \
  "(uid=*)" 2>/dev/null | grep -q "uidnumber"; then
  echo -e "  ${GREEN}Bind OK${NC}"
else
  echo -e "  ${RED}Bind FAILED${NC}"
  exit 1
fi

# 2. Get employee number
echo ""
echo -e "${YELLOW}[2/2] Employee Number${NC}"

EMPLOYEE_NUM=$(docker compose exec -T "$GLAUTH_CONTAINER" ldapsearch -x \
  -H "$LDAP_URL" \
  -b "$BASE_DN" \
  -D "cn=${USER},${BASE_DN}" \
  -w "$PASSWORD" \
  "(uid=*)" 2>/dev/null | grep "uidnumber:" | awk '{print $2}')

if [ -n "$EMPLOYEE_NUM" ]; then
  echo -e "  Employee number: ${GREEN}${EMPLOYEE_NUM}${NC}"
  echo ""
  echo "=========================================="
  echo -e " ${GREEN}Auth OK${NC}"
  echo "=========================================="
  echo ""
  echo "Use in Postman:"
  echo "  employee_number = ${EMPLOYEE_NUM}"
  echo "  user_role = admin|operator|technician"
else
  echo -e "  ${RED}Could not extract employee number${NC}"
  exit 1
fi
