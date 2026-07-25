#!/bin/bash
# ============================================================================
# generate-jwt.sh — Generador standalone JWT para testing
# Uso: ./scripts/generate-jwt.sh [employee_number] [role] [ttl_minutes]
# Requiere: PGRST_JWT_SECRET环境变量
# ============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

if [ -z "$PGRST_JWT_SECRET" ]; then
  echo -e "${RED}Error: PGRST_JWT_SECRET environment variable is required${NC}"
  echo "Usage: PGRST_JWT_SECRET=<secret> $0 [employee_number] [role] [ttl_minutes]"
  exit 1
fi

EMPLOYEE_NUMBER=${1:-1001}
ROLE=${2:-admin}
TTL_MINUTES=${3:-60}
SECRET=$PGRST_JWT_SECRET

echo "=========================================="
echo " JWT Generator (standalone)"
echo "=========================================="
echo ""

# Check if node is available
if ! command -v node &>/dev/null; then
  echo -e "${RED}Error: node not found${NC}"
  echo "Install Node.js or use Postman pre-request script instead"
  exit 1
fi

# Generate JWT using Node.js
JWT=$(node -e "
const crypto = require('crypto');

function base64url(str) {
  return Buffer.from(str)
    .toString('base64')
    .replace(/=/g, '')
    .replace(/\+/g, '-')
    .replace(/\//g, '_');
}

function sign(header, payload, secret) {
  const h = base64url(JSON.stringify(header));
  const p = base64url(JSON.stringify(payload));
  const data = h + '.' + p;
  const sig = crypto
    .createHmac('sha256', secret)
    .update(data)
    .digest('base64')
    .replace(/=/g, '')
    .replace(/\+/g, '-')
    .replace(/\//g, '_');
  return data + '.' + sig;
}

const now = Math.floor(Date.now() / 1000);
const header = { alg: 'HS256', typ: 'JWT' };
const payload = {
  sub: '${EMPLOYEE_NUMBER}',
  role: '${ROLE}',
  iat: now,
  exp: now + (${TTL_MINUTES} * 60)
};

console.log(sign(header, payload, '${SECRET}'));
")

echo -e "${YELLOW}Payload:${NC}"
echo "  sub:    ${EMPLOYEE_NUMBER}"
echo "  role:   ${ROLE}"
echo "  ttl:    ${TTL_MINUTES} minutes"
echo ""
echo -e "${YELLOW}Token:${NC}"
echo "$JWT"
echo ""
echo "=========================================="
echo -e " ${GREEN}Generated OK${NC}"
echo "=========================================="
echo ""
echo "Use in Postman:"
echo "  Set jwt_token = $JWT"
