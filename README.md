# Field Service Management (FSM)

Generic FSM system for multi-industry field service operations. Single-tenant architecture with PostgreSQL, Go Backend (Gin), and JWT authentication.

## Stack

| Component | Version | Purpose |
|-----------|---------|---------|
| PostgreSQL | 16 + PostGIS 3.5 | Database |
| Go Backend | 1.26 | REST API (Gin + pgx) |
| Next.js Frontend | 15 | Web UI (React) |
| Traefik | 3.0 | API Gateway |
| GLAuth | latest | LDAP server |
| Docker Compose | v2 | Container orchestration |

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        TRAEFIK :8088                            │
│  ┌─────────────────────┐  ┌──────────────────────────────────┐ │
│  │ /api/* /auth/*      │  │ /* (browser)                     │ │
│  │ → Go Backend :8081  │  │ → Next.js Frontend :3000        │ │
│  └─────────────────────┘  └──────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘

Go Backend (Gin + pgx)
├── LDAP bind (GLAuth :389) → validates credentials
├── JWT generation (HS256)
├── JWT middleware (validates tokens)
└── Direct PostgreSQL queries via pgx

Next.js Frontend
└── Client-side fetch to Go Backend via /api/* and /auth/*
```

- **Single-tenant**: Sin RLS, sin `company_id`
- **Core + Industry Packs**: FSM generico + packs especificos por industria
- **JWT Auth**: Go Backend validates JWT, queries PostgreSQL directly

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
├── docker-compose.yml          # Dev stack
├── docker-compose.prod.yml     # Production stack
├── .env.example                # Environment template
├── .env.prod.example           # Production template
│
├── docker/init/                # PostgreSQL init scripts
│   ├── 01-roles.sh             # Database roles
│   ├── 02-schemas.sql          # 9 bounded contexts
│   ├── 03-schema.sql           # DDL (~1000 lines)
│   ├── 04-seeds.sql            # Initial data
│   ├── 05-comments.sql         # OpenAPI comments
│   ├── 06-pre-request.sql      # JWT auth function
│   └── 07-seed-data-dev.sql    # Dev seed data
│
├── traefik/                    # API Gateway
│   ├── traefik.yml             # Static config (dev)
│   ├── traefik.prod.yml        # Static config (prod)
│   ├── dynamic.yml             # Routes (dev)
│   └── dynamic.prod.yml        # Routes (prod)
│
├── glauth/                     # LDAP server
│   └── config.toml             # Users and groups
│
├── backend/                    # Go Backend (Gin + pgx)
│   ├── cmd/server/main.go      # Entry point
│   ├── internal/               # Handlers, services, repositories
│   ├── scripts/                # Tenant provisioning tests
│   └── Dockerfile
│
├── frontend/                   # Next.js Frontend (React)
│   ├── src/                    # App routes and components
│   ├── tests/                  # Playwright tests
│   └── Dockerfile
│
├── industry-packs/             # Industry-specific extensions
│   └── gas/                    # Gas industry pack
│       ├── initdb.d/           # Auto-executed scripts
│       └── README.md
│
├── scripts/                    # Operational scripts
│   ├── validate-stack.sh       # Healthcheck
│   ├── test-auth.sh            # LDAP test
│   └── provision-tenant.sh     # Tenant DB provisioning
│
├── postman/                    # API collection
│   ├── FSM.postman_collection.json
│   └── scripts/
│
├── sql/
│   ├── schema.md               # Schema documentation
│   └── 20260813_*.sql          # Standalone migrations
│
├── DEPLOY.md                   # Production deployment
├── PLAN.md                     # Project plan
└── AGENTS.md                   # AI agent instructions
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `POSTGRES_PASSWORD` | PostgreSQL superuser password | Yes |
| `JWT_SECRET` | JWT signing secret (HS256) | Yes |
| `LDAP_SERVICE_PASSWORD` | LDAP service bind password | Yes |
| `DATABASE_URL` | Backend PostgreSQL connection URL | Yes |

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
2. Traefik → Go Backend (:8081)
3. Go Backend validates JWT (HS256)
4. Go Backend extracts employee_number, role from claims
5. Go Backend queries PostgreSQL directly via pgx
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
| `domain_gas.measurement_types` | 8 measurement types |
| `domain_gas.certifications` | 9 SEC Chile certifications |
| `domain_gas.property_types` | 5 property types |
| `domain_gas.property_assets` | Extra property attributes |

Industry-agnostic catalogs live in the core (see Core Catalogs below): `operations.visit_types`, `operations.photo_findings`, `operations.pre_visit_results`, `operations.rejection_reasons`, and `partners.partner_service_types`. The gas pack seeds its specific codes into these core tables.

### Core Catalogs

The following lookup catalogs are seeded in the core (`docker/init/04-seeds.sql`) so they work for every tenant regardless of industry pack:

| Table | Generic seeds |
|-------|---------------|
| `inventory.vehicle_types` | truck, crane, van, crane_truck, other |
| `planning.route_types` | maintenance, installation, repair, inspection, delivery, collection, emergency, other |
| `operations.visit_types` | pre_visit, installation, maintenance, repair, diagnosis, inspection, emergency, certification, other |
| `operations.photo_findings` | normal, damage, incomplete, other |
| `operations.pre_visit_results` | approved, rejected, conditional |
| `operations.rejection_reasons` | customer_unavailable, customer_cancelled, missing_documentation, other |
| `partners.partner_service_types` | installation, maintenance, repair, certification, emergency, other |

Industry packs may add their own codes to these catalogs (see `industry-packs/gas/initdb.d/91-gas-02-seeds.sql`).

See `industry-packs/gas/README.md` for details.

### Using Industry Packs

```bash
# Dev
docker compose -f docker-compose.yml -f industry-packs/gas/docker-compose.override.yml up -d

# Prod
docker compose -f docker-compose.prod.yml -f industry-packs/gas/docker-compose.prod.yml --env-file .env.prod up -d
```

### Multi-tenant (SaaS, Option A: database-per-tenant)

Core + one industry pack per tenant database. One backend instance per tenant, pointed at its own `DATABASE_URL`.

```bash
# Provision a new tenant (core init + gas pack)
./scripts/provision-tenant.sh altogasspa gas

# Point a backend instance at the tenant
DATABASE_URL="postgres://fsm_admin:...@postgres:5432/fsm_altogasspa?sslmode=disable"
```

`inventory.vehicle_types`, `planning.route_types`, `operations.visit_types`, `operations.photo_findings`, `operations.pre_visit_results`, `operations.rejection_reasons`, and `partners.partner_service_types` live in the core catalog (seeded by `04-seeds.sql`), so these lookups work for every tenant regardless of industry pack.

## Scripts

| Script | Purpose | Usage |
|--------|---------|-------|
| `validate-stack.sh` | Healthcheck all services | `./scripts/validate-stack.sh` |
| `test-auth.sh` | Test LDAP bind | `./scripts/test-auth.sh <user> <password>` |
| `provision-tenant.sh` | Create a tenant DB (core + industry pack) | `./scripts/provision-tenant.sh <name> [pack]` |

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
| Core tables | 63 |
| Core enums | 27 |
| Industry pack tables | 4 |
| Industry pack FKs | 11 |
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
