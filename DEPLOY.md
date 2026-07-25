# Despliegue Producción — workflows.cl

## Prerrequisitos

- Docker + Docker Compose v2
- Dominio `workflows.cl` apuntando al servidor (A record)
- Subdominios configurados:
  - `app.workflows.cl` → servidor
  - `auth.workflows.cl` → servidor
- Puerto 80 y 443 abiertos
- Git instalado

## Pasos de despliegue

### 1. Clonar repo

```bash
git clone <repo-url>
cd localis
```

### 2. Configurar variables de producción

```bash
cp .env.prod.example .env.prod
```

### 3. Generar contraseñas seguras

```bash
# Generar cada contraseña
openssl rand -base64 32
```

Editar `.env.prod` con las contraseñas generadas:

```bash
POSTGRES_PASSWORD=<contraseña_generada>
FSM_API_PASSWORD=<contraseña_generada>
PGRST_JWT_SECRET=<contraseña_generada>
AUTHELIA_SESSION_SECRET=<contraseña_generada>
AUTHELIA_STORAGE_ENCRYPTION_KEY=<contraseña_generada>
AUTHELIA_JWT_SECRET=<contraseña_generada>
```

### 4. Configurar GLAuth

Editar `glauth/config.toml` con passwords hasheados para producción:

```bash
# Generar hash SHA256 de una contraseña
echo -n "tu_password" | sha256sum | awk '{print $1}'
```

Actualizar `passsha256` para cada usuario:
- `authelia` (service user)
- `admin`
- `operador`
- `tecnico`

### 5. Levantar stack

```bash
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
```

### 6. Verificar servicios

```bash
# Verificar que todos los contenedores estén corriendo
docker compose -f docker-compose.prod.yml ps

# Verificar logs
docker compose -f docker-compose.prod.yml logs -f
```

### 7. Verificar endpoints

```bash
# Health checks
curl -k https://app.workflows.cl/users
curl -k https://auth.workflows.cl/api/health
```

### 8. Configurar TOTP (MFA)

1. Acceder a `https://auth.workflows.cl`
2. Login con `admin` / `<password>`
3. Seguir instrucciones para configurar Authenticator app
4. Repetir para cada usuario

## Comandos útiles

```bash
# Ver logs de un servicio
docker compose -f docker-compose.prod.yml logs -f authelia
docker compose -f docker-compose.prod.yml logs -f traefik

# Reiniciar un servicio
docker compose -f docker-compose.prod.yml restart authelia

# Parar todo
docker compose -f docker-compose.prod.yml down

# Parar y eliminar volumes (CUIDADO: borra datos)
docker compose -f docker-compose.prod.yml down -v
```

## Variables de entorno

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `POSTGRES_PASSWORD` | Password de PostgreSQL | `openssl rand -base64 32` |
| `FSM_API_PASSWORD` | Password del usuario `fsm_api` | `openssl rand -base64 32` |
| `PGRST_JWT_SECRET` | Secret para JWT signing | `openssl rand -base64 32` |
| `AUTHELIA_SESSION_SECRET` | Secret para sesiones Authelia | 32+ chars |
| `AUTHELIA_STORAGE_ENCRYPTION_KEY` | Key para storage Authelia | 32+ chars |
| `AUTHELIA_JWT_SECRET` | Secret para JWT de Authelia | 32+ chars |
| `AUTHELIA_SESSION_DOMAIN` | Dominio para cookies | `workflows.cl` |
| `AUTHELIA_AUTHELIA_URL` | URL de Authelia | `https://auth.workflows.cl` |
| `AUTHELIA_DEFAULT_REDIRECT` | URL de redirect post-login | `https://app.workflows.cl` |

## Troubleshooting

### Authelia no inicia
```bash
docker compose -f docker-compose.prod.yml logs authelia
# Verificar que GLAuth esté corriendo
docker compose -f docker-compose.prod.yml ps glauth
```

### Traefik no obtiene certificados TLS
```bash
# Verificar que los puertos 80/443 estén abiertos
# Verificar DNS apunta al servidor
nslookup app.workflows.cl
nslookup auth.workflows.cl
```

### PostgREST no responde
```bash
# Verificar conexión a PostgreSQL
docker compose -f docker-compose.prod.yml exec postgres psql -U fsm_admin -d fsm_fsm -c "\dt"
```
