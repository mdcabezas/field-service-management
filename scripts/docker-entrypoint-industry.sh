#!/bin/bash
# ============================================================================
# docker-entrypoint-industry.sh — Wrapper to run industry pack init scripts
# Runs after standard postgres entrypoint completes
# ============================================================================

set -e

# Run the original postgres entrypoint
docker-entrypoint.sh postgres "$@" &
PID=$!

# Wait for postgres to be ready
until pg_isready -U "${POSTGRES_USER:-postgres}" -d "${POSTGRES_DB:-postgres}" > /dev/null 2>&1; do
  sleep 1
done

# Run industry pack init scripts
if [ -d "/docker-entrypoint-industry.d" ]; then
  echo "Running industry pack initialization scripts..."
  for f in /docker-entrypoint-industry.d/*.sql; do
    [ -f "$f" ] || continue
    echo "Running $f"
    psql -v ON_ERROR_STOP=1 -U "${POSTGRES_USER:-postgres}" -d "${POSTGRES_DB:-postgres}" -f "$f"
  done
  echo "Industry pack initialization complete."
fi

# Wait for postgres process
wait $PID