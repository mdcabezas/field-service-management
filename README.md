# Field Service Management (FSM)

Generic FSM system for multi-industry field service operations. Single-tenant architecture with PostgreSQL, PostgREST, and JWT authentication via Traefik.

## Stack

| Component | Version | Purpose |
|-----------|---------|---------|
| PostgreSQL | 16 + PostGIS 3.5 | Database |
| PostgREST | 12.2.3 | REST API |
| Traefik | 3.0 | API Gateway |
| GLAuth | latest | LDAP server |
| Authelia | 4.38 | MFA authentication |
| Docker Compose | v2 | Container orchestration |

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        TRAEFIK :8088                            │
│  ┌─────────────────────┐  ┌──────────────────────────────────┐ │
│  │ /api/*              │  │ /* (browser)                     │ │
│  │ → jwt-auth          │  │ → authelia forwardAuth           │ │
│  │ → PostgREST :3000   │  │ → Authelia :9091 → GLAuth :389   │ │
│  └─────────────────────┘  └──────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

- **Single-tenant**: Sin RLS, sin `company_id`
- **Core + Industry Packs**: FSM genérico + packs específicos por industria
- **JWT Auth**: Traefik valida JWT, inyecta headers → PostgREST

## Quick Start

### Development

```bash
# 1. Clone repo
git clone git@github.com:mdcabezas/field-service-management.git
cd field-service-management

# 2. Configure environment
cp .env.example .env

# 3. Generate secrets
sed -i 's/<REPLACE_WITH_GENERATED_SECRET>/$(openssl rand -base64 32)/g' .env

# 4. Start stack
docker compose --profile dev up -d

# 5. Verify
./scripts/validate-stack.sh
```

### Production

```bash
# 1. Configure production environment
cp .env.prod.example .env.prod

# 2. Generate secure passwords
openssl rand -base64 32  # Run for each variable

# 3. Start stack
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d

# 4. See DEPLOY.md for detailed instructions
```

## Project Structure

```
field-service-management/
├── docker-compose.yml          # Dev stack (6 services)
├── docker-compose.prod.yml     # Production stack
├── postgrest.conf              # PostgREST configuration
├── .env.example                # Environment template
├── .env.prod.example           # Production template
│
├── docker/init/                # PostgreSQL init scripts
│   ├── 01-roles.sh             # Database roles
│   ├── 02-schemas.sql          # 9 bounded contexts
│   ├── 03-schema.sql           # DDL (~1000 lines)
│   ├── 04-seeds.sql            # Initial data
│   ├── 05-comments.sql         # OpenAPI comments
│   └── 06-pre-request.sql      # JWT auth function
│
├── traefik/                    # API Gateway
│   ├── traefik.yml             # Static config (dev)
│   ├── traefik.prod.yml        # Static config (prod)
│   ├── dynamic.yml             # Routes (dev)
│   ├── dynamic.prod.yml        # Routes (prod)
│   └── jwt-validator/          # JWT validation service
│
├── glauth/                     # LDAP server
│   └── config.toml             # Users and groups
│
├── authelia/                   # MFA authentication
│   └── config/
│       └── configuration.yml   # Authelia config
│
├── industry-packs/             # Industry-specific extensions
│   └── gas/                    # Gas industry pack
│       ├── init/               # SQL scripts
│       ├── initdb.d/           # Auto-executed scripts
│       └── README.md
│
├── scripts/                    # Operational scripts
│   ├── validate-stack.sh       # Healthcheck
│   ├── test-auth.sh            # LDAP test
│   ├── generate-jwt.sh         # JWT generator
│   └── reload-postgrest-cache.sh
│
├── postman/                    # API collection
│   ├── GAS-FSM.postman_collection.json
│   └── scripts/
│
├── sql/
│   └── schema.md               # Schema documentation
│
├── DEPLOY.md                   # Production deployment
├── PLAN.md                     # Project plan
└── AGENTS.md                   # AI agent instructions
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `POSTGRES_PASSWORD` | PostgreSQL superuser password | Yes |
| `FSM_API_PASSWORD` | PostgREST API role password | Yes |
| `PGRST_JWT_SECRET` | JWT signing secret (shared) | Yes |
| `AUTHELIA_SESSION_SECRET` | Authelia session secret (32+ chars) | Yes |
| `AUTHELIA_STORAGE_ENCRYPTION_KEY` | Authelia storage key (32+ chars) | Yes |
| `AUTHELIA_JWT_SECRET` | Authelia JWT secret (32+ chars) | Yes |
| `AUTHELIA_LDAP_PASSWORD` | LDAP bind password for Authelia | Yes |
| `AUTHELIA_SESSION_DOMAIN` | Cookie domain (dev: localhost) | No |
| `AUTHELIA_AUTHELIA_URL` | Authelia URL (dev: http://localhost:9091) | No |
| `AUTHELIA_DEFAULT_REDIRECT` | Post-login redirect URL | No |

Generate secrets with:
```bash
openssl rand -base64 32
```

## Authentication

### Users (GLAuth)

| Employee Number | Username | Role | Password |
|-----------------|----------|------|----------|
| 1001 | admin | admin | dev only |
| 1002 | operador | operator | dev only |
| 1003 | tecnico | technician | dev only |

### JWT Flow

```
1. Client → Traefik (:8088)
2. Traefik → JWT Validator (forwardAuth)
3. JWT Validator validates HMAC-SHA256 signature
4. JWT Validator returns X-User-ID, X-User-Role headers
5. Traefik strips Authorization header
6. Traefik → PostgREST (:3000)
7. PostgREST calls set_user_context()
8. set_user_context() reads headers → sets JWT claims
9. Query executes as fsm_api with user context
```

### JWT Payload

```json
{
  "sub": "1001",
  "role": "admin",
  "iat": 1700000000,
  "exp": 1700086400
}
```

## Industry Packs

### Gas Industry Pack

Adds Chilean gas industry-specific tables and constraints:

| Table | Purpose |
|-------|---------|
| `domain_gas.visit_types` | 7 visit types |
| `domain_gas.measurement_types` | 8 measurement types |
| `domain_gas.certifications` | 9 SEC Chile certifications |
| `domain_gas.rejection_reasons` | 17 rejection reasons |
| `domain_gas.property_types` | 5 property types |
| `domain_gas.vehicle_types` | 5 vehicle types |
| `domain_gas.route_types` | 6 route types |
| `domain_gas.photo_findings` | 6 photo finding types |
| `domain_gas.partner_service_types` | 5 service types |
| `domain_gas.pre_visit_results` | 3 results |
| `domain_gas.property_assets` | Extra property attributes |

See `industry-packs/gas/README.md` for details.

### Using Industry Packs

```bash
# Dev
docker compose -f docker-compose.yml -f industry-packs/gas/docker-compose.override.yml up -d

# Prod
docker compose -f docker-compose.prod.yml -f industry-packs/gas/docker-compose.prod.yml --env-file .env.prod up -d
```

## Scripts

| Script | Purpose | Usage |
|--------|---------|-------|
| `validate-stack.sh` | Healthcheck all 6 services | `./scripts/validate-stack.sh` |
| `test-auth.sh` | Test LDAP bind | `./scripts/test-auth.sh <user> <password>` |
| `generate-jwt.sh` | Generate JWT token | `PGRST_JWT_SECRET=<secret> ./scripts/generate-jwt.sh` |
| `reload-postgrest-cache.sh` | Reload PostgREST schema cache | `./scripts/reload-postgrest-cache.sh` |

## Database

### Schemas (9 bounded contexts)

| Schema | Purpose |
|--------|---------|
| `core` | Users, roles, audit log |
| `partners` | Partner companies, agreements, SLAs |
| `customers` | Customers, addresses, properties, technicians |
| `inventory` | Materials, tools, PPE, vehicles, costs |
| `operations` | Visits, assignments, checklists, measurements |
| `planning` | Daily plans, routes, material loading |
| `notifications` | Templates and notification instances |
| `geocoding` | Global address catalog with PostGIS |
| `shared` | Universal enums and catalogs |

### Key Metrics

| Metric | Value |
|--------|-------|
| Core tables | ~49 |
| Core enums | 27 |
| Industry pack tables | 11 |
| Industry pack FKs | 12 |
| Triggers | 13 |
| Functions | 4 |

## Contributing

1. Create feature branch from `develop`
2. Make changes
3. Run `./scripts/validate-stack.sh`
4. Commit with conventional commits
5. Create PR to `develop`

## License

Private - Workflows.cl
