# Despliegue Produccion — workflows.cl

## Prerrequisitos

- Docker + Docker Compose v2
- Dominio `workflows.cl` apuntando al servidor (A record)
- Subdominio configurado:
  - `app.workflows.cl` → servidor
- Puerto 80 y 443 abiertos
- Git instalado

## Pasos de despliegue

### 1. Clonar repo

```bash
git clone <repo-url>
cd localis
```

### 2. Configurar variables de produccion

```bash
cp .env.prod.example .env.prod
```

### 3. Generar contrasenas seguras

```bash
# Generar cada contrasena
openssl rand -base64 32
```

Editar `.env.prod` con las contrasenas generadas:

```bash
POSTGRES_PASSWORD=<contrasena_generada>
JWT_SECRET=<contrasena_generada>
LDAP_SERVICE_PASSWORD=<contrasena_generada>
```

### 4. Configurar GLAuth

Editar `glauth/config.toml` con passwords hasheados para produccion:

```bash
# Generar hash SHA256 de una contrasena
echo -n "tu_password" | sha256sum | awk '{print $1}'
```

Actualizar `passsha256` para cada usuario:
- `svc-localis` (service user)
- `admin`
- `operador`
- `tecnico`

### 5. Levantar stack

```bash
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d
```

### 6. Verificar servicios

```bash
# Verificar que todos los contenedores esten corriendo
docker compose -f docker-compose.prod.yml ps

# Verificar logs
docker compose -f docker-compose.prod.yml logs -f
```

### 7. Verificar endpoints

```bash
# Health checks
curl -k https://app.workflows.cl/health
curl -k -X POST https://app.workflows.cl/auth/login \
  -H "Content-Type: application/json" \
  -d '{"employee_number": "1001", "password": "<password>"}'
```

## Comandos utiles

```bash
# Ver logs de un servicio
docker compose -f docker-compose.prod.yml logs -f go-backend
docker compose -f docker-compose.prod.yml logs -f traefik

# Reiniciar un servicio
docker compose -f docker-compose.prod.yml restart go-backend

# Parar todo
docker compose -f docker-compose.prod.yml down

# Parar y eliminar volumes (CUIDADO: borra datos)
docker compose -f docker-compose.prod.yml down -v
```

## Variables de entorno

| Variable | Descripcion | Ejemplo |
|----------|-------------|---------|
| `POSTGRES_PASSWORD` | Password de PostgreSQL | `openssl rand -base64 32` |
| `JWT_SECRET` | Secret para JWT signing (HS256) | `openssl rand -base64 32` |
| `LDAP_SERVICE_PASSWORD` | Password del service user LDAP | `openssl rand -base64 32` |

## Troubleshooting

### Go Backend no inicia
```bash
docker compose -f docker-compose.prod.yml logs go-backend
# Verificar que GLAuth este corriendo
docker compose -f docker-compose.prod.yml ps glauth
# Verificar que PostgreSQL este corriendo
docker compose -f docker-compose.prod.yml ps postgres
```

### Traefik no obtiene certificados TLS
```bash
# Verificar que los puertos 80/443 esten abiertos
# Verificar DNS apunta al servidor
nslookup app.workflows.cl
```
