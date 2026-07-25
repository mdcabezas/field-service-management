#!/bin/bash
# ============================================================================
# reload-postgrest-cache.sh — Recarga el schema cache de PostgREST
# Ejecutar después de cambiar FKs o estructura de tablas
# ============================================================================

set -e

echo "Recargando schema cache de PostgREST..."

# Enviar SIGUSR1 al proceso PostgREST para recargar el cache
docker compose exec -T postgrest kill -SIGUSR1 1 2>/dev/null || {
    echo "Error: No se pudo enviar SIGUSR1 a PostgREST"
    echo "Intentando restart del container..."
    docker compose restart postgrest
}

echo "Schema cache recargado correctamente"
echo "Verificar con: curl -s http://localhost:3000/ | head -20"
