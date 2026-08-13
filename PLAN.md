# PLAN — Generic FSM System (Single-tenant)

## 0. Estado general

- **Stack**: PostgreSQL 16 + PostGIS 3.5 + Go Backend (Gin + pgx) + Docker Compose dev
- **Modo**: Single-tenant (sin RLS, sin company_id)
- **Arquitectura**: Core FSM generico + Industry Packs (gas, electricity, HVAC, plumbing, etc.)
- **Esquemas Core**: 9 bounded contexts, ~49 tablas, 27 enums nativos PostgreSQL
- **Industry Packs**: Cargados dinamicamente, anaden lookup tables, FKs, seeds especificos
- **Auth**: GLAuth (LDAP) → Go Backend (JWT) → PostgreSQL directo via pgx
- **Status actual**: Core generico completo. Gas Industry Pack creado. Go Backend completo con handlers, services, repos. Auth, RBAC, rate limiting, security headers implementados.

---

## 1. Ejecutado (F0-F5)

### ✅ F0 Decisiones de arquitectura

| ID | Decisión | Estado |
|----|----------|--------|
| D1 | `addresses` = catálogo global `geocoding` | ✅ |
| D2 | `route_id` ↔ `daily_plan_id` = trigger function `validate_visit_route_plan()` | ✅ |
| D3 | `vehicle_assignments` = tabla única con daily_plan_id y visit_id opcionales | ✅ |
| D4 | 27 enums nativos PostgreSQL en `shared` | ✅ |
| D5 | Certificaciones/licencias en Industry Pack (no en core) | ✅ |
| A | Polimorfismo items → split en `_materials`/`_tools`/`_epps` (Opción A) | ✅ |
| — | Single-tenant (sin company_id, sin RLS, sin app_tenant_id) | ✅ |
| — | JWT payload: `{sub, role}` — sin company_id | ✅ |

### ✅ F1 Init scripts Core (`docker/init/`)

| Archivo | Contenido | Estado |
|---------|-----------|--------|
| `01-roles.sql` | Roles: fsm_api (login), fsm_backend (login), fsm_migrator (nologin). search_path multi-schema | ✅ |
| `02-schemas.sql` | 9 esquemas core: core, partners, customers, inventory, operations, planning, notifications, geocoding, shared | ✅ |
| `03-schema.sql` | DDL completo (~1000 líneas): 27 enums, 3 funciones helper, ~49 tablas, 13 triggers, grants | ✅ |
| `04-seeds.sql` | tech_roles (4) | ✅ |
| `05-comments.sql` | COMMENT ON para OpenAPI (schema/table/column) — todo en inglés | ✅ |

### ✅ F2 Dev stack

| Archivo | Contenido | Estado |
|---------|-----------|--------|
| `docker-compose.yml` | postgres + traefik + go-backend + glauth, red interna `fsm-internal`, healthcheck | ✅ |
| `traefik/traefik.yml` | Traefik v3 static config (entrypoint :8088, dashboard :8080, file provider) | ✅ |
| `traefik/dynamic.yml` | Go Backend routes (api, auth, health) | ✅ |
| `glauth/config.toml` | GLAuth LDAP config (4 users, 4 groups, TOML) | ✅ |
| `backend/` | Go Backend (Gin + pgx): auth + CRUD + business logic | ✅ |
| `.env.example` | POSTGRES_PASSWORD, JWT_SECRET, LDAP_SERVICE_PASSWORD | ✅ |
| `.gitignore` | Excluye `.env`, `backend/server` | ✅ |

### ✅ F3 Documentación

| Archivo | Estado |
|---------|--------|
| `sql/schema.md` | Actualizado: single-tenant, 27 enums, 9 esquemas core, ~49 tablas | ✅ |
| `sql/schema.sql` | 🗑️ Eliminado (contenía multi-tenant obsoleto) | ✅ |

### ✅ F4 Validación de invariantes (pendiente levantar stack)

| Invariante | Esperado |
|------------|----------|
| `company_id` (non-comment) | 0 |
| `ENABLE ROW LEVEL SECURITY` | 0 |
| `CREATE POLICY` | 0 |
| `app_tenant_id` | 0 |
| `app_user_id` | ≥1 |
| `tool_status` (enum) | ≥1 |
| `validate_visit_route_plan` | ≥1 |
| `CREATE TYPE` (enums) | 27 |
| `CREATE SCHEMA` (core) | 9 |
| Order: `geocoding.addresses` < `customers.*` | ✓ |
| Order: `shared.report_templates` < `partners.partner_agreement_forms` | ✓ |

### ✅ F5 Bugfixes aplicados durante refactor

1. `customers.customer_addresses` referenciaba `geocoding.addresses` antes de su creación → orden de DDL reordenado
2. `partner_agreement_forms` referenciaba `shared.report_templates` antes de su creación → schemas shared y geocoding movidos antes de core/partners/customers
3. D2 CHECK constraint usaba subquery (inválido en PostgreSQL) → reemplazado por trigger function
4. `04-seeds.sql` INSERTaba `es_global` en `tech_roles` → columna eliminada
5. `inventory.tools.estado` era TEXT → migrado a enum `shared.tool_status`
6. `04-seeds.sql` INSERT tech_roles tenía 3 valores en 2 columnas → corregido
7. **F5.1 Industry Pack Architecture** — Core genérico + Industry Packs cargados dinámicamente
8. **F5.2 English Naming Refactor** — Nombres de enums, columnas, seeds y comentarios en inglés (estándar industria FSM)
9. **F5.3 Traefik + JWT Authentication** — API Gateway con Go Backend (auth + proxy)

---

## 2. Arquitectura: Core Genérico + Industry Packs

### Core FSM (Genérico)
- 9 esquemas bounded context
- 27 enums universales (status, priority, type, etc.)
- ~49 tablas base
- Sin datos específicos de industria
- Extensible vía Industry Packs

### Industry Packs (Ej: Gas Chile)

```
industry-packs/gas/
├── init/
│   ├── 01-schemas.sql      # CREATE SCHEMA domain_gas
│   ├── 02-tables.sql       # 11 lookup tables + 12 FKs a core
│   ├── 03-seeds.sql        # Datos Chile gas (SEC, tipos visita, etc.)
│   ├── 04-grants.sql       # GRANTs para fsm_api + search_path
│   └── 05-comments.sql     # COMMENTS para OpenAPI
├── docker-compose.override.yml  # Monta init en postgres
└── README.md
```

**Carga dinámica**: `docker compose -f docker-compose.yml -f industry-packs/gas/docker-compose.override.yml up -d`

### Patrones de Industry Pack

| Elemento | Core | Industry Pack |
|----------|------|---------------|
| Enums universales | ✅ (27 en `shared`) | — |
| Lookup tables específicas | — | ✅ (visit_types, measurement_types, etc.) |
| Certificaciones/licencias | — | ✅ (SEC, HVAC licenses, etc.) |
| Rejection reasons genéricos | ✅ `rejection_category` enum | ✅ `rejection_reasons` tabla |
| FKs a lookups | Columnas UUID sin FK | Industry pack añade FKs |
| Seeds | Solo tech_roles | Datos específicos industria |

---

## 3. Gas Industry Pack (Creado)

### Tablas añadidas (`domain_gas`)

| Tabla | Propósito |
|-------|-----------|
| `measurement_types` | tightness, pressure, co, draft, leak, ph, temperature, other |
| `property_types` | tank, meter, indoor_piping, appliance, other (con columns_config JSONB) |
| `certifications` | 9 certs SEC Chile (CL1, CL2, CL3, GLP, TC1, TC2, TC6, GREEN SEAL, DS66) |
| `property_assets` | Atributos extra por propiedad (poles, meter, etc.) |

Los catálogos universales (`visit_types`, `photo_findings`, `pre_visit_results`, `rejection_reasons`, `partner_service_types`) se movieron al core (`operations.*` / `partners.*`, migración `sql/20260813_catalogs_core.sql`); el pack gas siembra sus códigos en esas tablas core.

### FKs añadidas a Core (12)

| Core Table | Column | Referencia |
|------------|--------|------------|
| `partners.partner_agreements` | `service_type` | `partners.partner_service_types` *(core catalog)* |
| `partners.slas` | `work_type` | `operations.visit_types` *(core catalog)* |
| `customers.properties` | `type` | `domain_gas.property_types` |
| `customers.tech_certifications` | `cert_id` | `domain_gas.certifications` |
| `inventory.vehicles` | `type` | `inventory.vehicle_types` *(moved to core for multi-tenant — seeded in 04-seeds.sql)* |
| `inventory.checklist_templates` | `work_type` | `operations.visit_types` *(core catalog)* |
| `planning.routes` | `type` | `planning.route_types` *(moved to core catalog; gas seeds its codes)* |
| `operations.visits` | `type` | `operations.visit_types` *(core catalog)* |
| `operations.visits` | `result` | `operations.pre_visit_results` *(core catalog)* |
| `operations.visits` | `rejection_reason_id` | `operations.rejection_reasons` *(core catalog)* |
| `operations.visit_measurements` | `type` | `domain_gas.measurement_types` |
| `operations.visit_photos` | `finding_type` | `operations.photo_findings` *(core catalog)* |

---

## 4. Stack Running — Inspección

Estado: Stack levantado y validado. Traefik + JWT auth funcionando.

### Puertos expuestos

| Servicio | Puerto | Propósito |
|----------|--------|-----------|
| Traefik (HTTP) | :8088 | API Gateway (JWT validation) |
| Traefik (Dashboard) | :8080 | Traefik admin dashboard |
| PostgREST (direct) | :3000 | Direct access (sin JWT, para debugging) |
| PostgreSQL | :5432 | Database direct access |

### Validación completa

| Test | Resultado |
|------|-----------|
| Sin JWT → Traefik | 401 Unauthorized |
| JWT inválido → Traefik | 401 Unauthorized |
| LDAP bind (svc-localis) → Go Backend login | 200 OK |
| JWT válido → /api/users | 200 OK |

### Architecture: JWT Auth Flow

```
Client → Traefik :8088
  → Go Backend :8080
    → Validates JWT (HS256)
    → Extracts employee_number, role from claims
    → Queries PostgreSQL directly via pgx
```

### Architecture: Auth Flow (GLAuth + Go Backend)

```
┌─────────────────────────────────────────────────────────────────┐
│                        TRAEFIK :8088                            │
│  ┌─────────────────────┐  ┌──────────────────────────────────┐ │
│  │ /api/* /auth/*      │  │ /* (browser)                     │ │
│  │ → Go Backend :8080  │  │ → Go Backend :8080               │ │
│  └─────────────────────┘  └──────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘

GLAuth (LDAP)
├── uid=1000, cn=svc-localis (service user)
├── uid=1001, cn=admin (admins)
├── uid=1002, cn=operador (operators)
└── uid=1003, cn=tecnico (technicians)

Go Backend → LDAP search → Generate JWT {sub: employee_number, role}
Go Backend → PostgreSQL via pgx (direct queries)
```

---

## 5. Fases pendientes

### F6 Backend Go ✅ COMPLETADO

Go Backend implementado con la siguiente arquitectura:

- **Framework**: Gin v1.12.0 + pgx/v5
- **Auth**: LDAP bind → JWT (HS256, jti, issuer validation) → RBAC middleware
- **Handler/Service/Repo**: Thin controller pattern, ~55 tablas, interfaces en `repository/interfaces.go`
- **Seguridad**: Rate limiting (login 10/min, API 100/min), CORS allowlist, security headers, body limit 10MB
- **Operational**: Request ID, structured logging (slog), graceful shutdown
- **DB**: Direct PostgreSQL via pgx (no ORM), PostGIS geometry handling, transactions on critical paths

### F7 Pendiente

| Tema | Estado |
|------|--------|
| API style (REST por esquema) | ✅ Recurso-driven implementado |
| JWT algorithm | ✅ HS256 (hardcoded, secret ≥32 bytes) |
| Migration tooling | Pendiente |
| Logging | ✅ slog (structured) |
| ORM | ✅ pgx directo (hand-written repos) |
| Tests framework | Pendiente (sin tests unitarios) |
| Auth flow | ✅ GLAuth (LDAP) + Go Backend |
| Background workers | Pendiente |

### F8 UI / Frontend
No definido.

### F9+ Integraciones externas
Identificadas, no en scope:
- SEC: validaciones de normas (gas)
- Abastible/Gasco/Metrogas: APIs de partners
- Google Maps: geocoding chileno
- Twilio/SendGrid: notificaciones
- SII (DTE): facturación electrónica chilena
- Firma digital (FirmaChile): informes firmados

---

## 6. Inventario de artefactos

```
localis/
├── AGENTS.md
├── PLAN.md
├── .env.example
├── .gitignore
├── docker-compose.yml
├── opencode.json
├── sql/
│   └── schema.md
├── traefik/
│   ├── traefik.yml
│   ├── dynamic.yml
│   ├── traefik.prod.yml
│   └── dynamic.prod.yml
├── glauth/
│   └── config.toml
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   └── Dockerfile
├── docker/init/
│   ├── 01-roles.sql
│   ├── 02-schemas.sql
│   ├── 03-schema.sql
│   ├── 04-seeds.sql
│   ├── 05-comments.sql
│   └── 06-pre-request.sql
├── scripts/
│   ├── validate-stack.sh
│   └── test-auth.sh
├── industry-packs/
│   └── gas/
│       ├── init/
│       ├── docker-compose.override.yml
│       └── README.md
├── postman/
│   ├── FSM.postman_collection.json
│   ├── scripts/
│   └── README.md
└── examples/
    ├── altogasspa/
    ├── cegas/
    ├── evj/
    ├── ibssagroup/
    ├── ipiex/
    └── rigas/
```

### Métricas finales

| Métrica | Valor |
|---------|-------|
| Tablas Core | ~49 |
| Enums Core | 27 nativos |
| Tablas Industry Pack (Gas) | 11 |
| FKs Industry Pack (Gas) | 12 |
| Triggers | 13 (3 updated_at + 1 D2 + 9 audit) |
| Funciones | 4 (app_user_id, touch_updated_at, audit_change, validate_visit_route_plan) |
| Esquemas Core | 9 bounded contexts |
| Esquemas + Industry Pack | 10 |
| GRANT statements | 9 core + industry pack |
| Init files Core | 6 |
| Init files Industry Pack (Gas) | 5 |
| Auth | GLAuth (LDAP) → Go Backend (JWT) → Traefik → PostgreSQL via pgx |
| Backend handlers | ~50 handler files |
| Backend services | ~15 service files |
| Backend repos | ~50 repository files |
| Naming convention | English (industry standard: snake_case) |