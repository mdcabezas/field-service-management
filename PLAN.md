# PLAN — Generic FSM System (Single-tenant)

## 0. Estado general

- **Stack**: PostgreSQL 16 + PostGIS 3.5 + PostgREST v12.2.3 + Docker Compose dev
- **Modo**: Single-tenant (sin RLS, sin company_id)
- **Arquitectura**: Core FSM genérico + Industry Packs (gas, electricity, HVAC, plumbing, etc.)
- **Esquemas Core**: 9 bounded contexts, ~49 tablas, 27 enums nativos PostgreSQL
- **Industry Packs**: Cargados dinámicamente, añaden lookup tables, FKs, seeds específicos
- **Arquitectura PostgREST**: JWT firmado por backend Go; fsm_api consume vía PostgREST.
- **Auth**: GLAuth (LDAP) → Authelia (MFA) + Backend Go (JWT) → Traefik → PostgREST
- **Status actual**: Core genérico completo. Gas Industry Pack creado. GLAuth + Authelia + Traefik + JWT auth funcionando. Pendiente: F6+ (backend).

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
| `docker-compose.yml` | postgres + postgrest + traefik + jwt-validator + authelia + glauth, red interna `fsm-internal`, healthcheck | ✅ |
| `postgrest.conf` | db-schemas (9 core), db-anon-role=fsm_api, db-pre-request=public.set_user_context | ✅ |
| `traefik/traefik.yml` | Traefik v3 static config (entrypoint :8088, dashboard :8080, file provider) | ✅ |
| `traefik/dynamic.yml` | JWT forwardAuth middleware + authelia middleware + strip-auth-header + PostgREST router | ✅ |
| `traefik/jwt-validator/` | Minimal Go JWT validator (HMAC-SHA256, stdlib only, ~10MB image) | ✅ |
| `glauth/config.toml` | GLAuth LDAP config (4 users, 4 groups, TOML) | ✅ |
| `authelia/config/configuration.yml` | Authelia config (LDAP backend via GLAuth, TOTP, session, access control) | ✅ |
| `.env.example` | POSTGRES_PASSWORD, FSM_API_PASSWORD, PGRST_JWT_SECRET, AUTHELIA_* | ✅ |
| `.gitignore` | Excluye `.env` | ✅ |
| `scripts/reload-postgrest-cache.sh` | Script para recargar schema cache después de cambiar FKs | ✅ |

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
5. `postgrest.conf` usaba key obsoleta `role-claim-key` → migrada a `jwt-role-claim-key = "$$.role"`
6. `inventory.tools.estado` era TEXT → migrado a enum `shared.tool_status`
7. `04-seeds.sql` INSERT tech_roles tenía 3 valores en 2 columnas → corregido
8. Archivo deprecado multi-tenant `sql/schema.sql` → eliminado
9. **F5.1 Industry Pack Architecture** — Core genérico + Industry Packs cargados dinámicamente
10. **F5.2 English Naming Refactor** — Nombres de enums, columnas, seeds y comentarios en inglés (estándar industria FSM)
11. **F5.3 Traefik + JWT Authentication** — API Gateway con validación JWT via forwardAuth middleware

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
| `visit_types` | pre_visit, installation, maintenance, emergency, certification, repair, diagnosis |
| `measurement_types` | tightness, pressure, co, draft, leak, ph, temperature, other |
| `photo_findings` | normal, leak, damage, emergency, incomplete, other |
| `vehicle_types` | truck, crane, van, crane_truck, other |
| `partner_service_types` | installation, maintenance, certification, emergency, other |
| `pre_visit_results` | approved, rejected, conditional |
| `route_types` | meter_reading, letter_delivery, tank_collection, tank_delivery, mass_inspection, mixed |
| `property_types` | tank, meter, indoor_piping, appliance, other (con columns_config JSONB) |
| `certifications` | 9 certs SEC Chile (CL1, CL2, CL3, GLP, TC1, TC2, TC6, GREEN SEAL, DS66) |
| `rejection_reasons` | 17 razones en 4 categorías (logistics, regulatory, customer, operational) |
| `property_assets` | Atributos extra por propiedad (poles, meter, etc.) |

### FKs añadidas a Core (12)

| Core Table | Column | Referencia |
|------------|--------|------------|
| `partners.partner_agreements` | `service_type` | `domain_gas.partner_service_types` |
| `partners.slas` | `work_type` | `domain_gas.visit_types` |
| `customers.properties` | `type` | `domain_gas.property_types` |
| `customers.tech_certifications` | `cert_id` | `domain_gas.certifications` |
| `inventory.vehicles` | `type` | `domain_gas.vehicle_types` |
| `inventory.checklist_templates` | `work_type` | `domain_gas.visit_types` |
| `planning.routes` | `type` | `domain_gas.route_types` |
| `operations.visits` | `type` | `domain_gas.visit_types` |
| `operations.visits` | `result` | `domain_gas.pre_visit_results` |
| `operations.visits` | `rejection_reason_id` | `domain_gas.rejection_reasons` |
| `operations.visit_measurements` | `type` | `domain_gas.measurement_types` |
| `operations.visit_photos` | `finding_type` | `domain_gas.photo_findings` |

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
| JWT válido → app_user_id() | employee_number correcto |
| JWT válido → tech_roles query | 4 rows |
| JWT válido → OpenAPI schema | 5 paths, 3 definitions |
| Authelia login → GLAuth LDAP bind | 200 OK |
| Authelia TOTP → MFA challenge | 200 OK |

### Architecture: JWT Auth Flow

```
Client → Traefik :8088
  → JWT Validator (forwardAuth) validates signature
  → Injects X-User-ID (employee_number), X-User-Role headers
  → Strip Authorization header
  → PostgREST :3000
    → db-pre-request: set_user_context()
      → Reads headers from request.headers JSON
      → Propagates employee_number + role via set_config
    → Main query executes as fsm_api
```

### Architecture: Auth Flow (GLAuth + Authelia)

```
┌─────────────────────────────────────────────────────────────────┐
│                        TRAEFIK :8088                            │
│  ┌─────────────────────┐  ┌──────────────────────────────────┐ │
│  │ /api/*              │  │ /* (browser)                     │ │
│  │ → jwt-auth          │  │ → authelia forwardAuth           │ │
│  │ → PostgREST :3000   │  │ → Authelia :9091 → GLAuth :389   │ │
│  └─────────────────────┘  └──────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘

GLAuth (LDAP)
├── uid=1000, cn=authelia (service user)
├── uid=1001, cn=admin (admins)
├── uid=1002, cn=operador (operators)
└── uid=1003, cn=tecnico (technicians)

Authelia → LDAP bind → Session cookie
Backend Go → LDAP search → Generate JWT {sub: employee_number, role}
PostgREST → app_user_id() returns employee_number (TEXT)
```

---

## 5. Fases pendientes

### F6 Backend Go [SIGUIENTE]

Opciones de diseño de API:

**Opción A — REST por esquema (bounded context mapping)**:
```
/api/core/users
/api/operations/visits
/api/planning/daily_plans
```
- A favor: mapeo 1:1 con bounded contexts, fácil split futuro a microservicios
- Contra: URLs verbosas

**Opción B — Recurso-driven**:
```
/api/users
/api/visits
/api/daily-plans
```
- A favor: REST standard
- Contra: requiere routing interno por schema

**Tareas F6**:
1. Scaffolding módulo Go (go.mod, cmd/, internal/)
2. LDAP client (go-ldap/ldap) para buscar usuarios en GLAuth
3. JWT signing (HS256 with PGRST_JWT_SECRET — shared with Traefik JWT validator)
4. Endpoint `/auth/token` (LDAP bind → JWT generation)
5. Capa de servicios por bounded context
6. Tipos shared Go ↔ enums PG
7. Tests unitarios + integration
8. Dockerfile

### F7 Definiciones pendientes

| Tema | Decisión |
|------|----------|
| API style (REST por esquema vs recurso) | pendiente |
| JWT algorithm (HS256 vs RS256) | HS256 (ya implementado en JWT validator) |
| Migration tooling (goose / atlas / manual) | pendiente |
| Logging (slog / zap / zerolog) | pendiente |
| ORM (sqlc / pgx direct / bun / gorm) | pendiente |
| Tests framework (testify / stdlib) | pendiente |
| Auth flow (login endpoint / IdP externo) | **GLAuth (LDAP) + Authelia** |
| Background workers (asynq / river / stdlib) | pendiente |

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
├── postgrest.conf
├── opencode.json
├── sql/
│   └── schema.md
├── traefik/
│   ├── traefik.yml              # Traefik v3 static config
│   ├── dynamic.yml              # JWT forwardAuth + authelia + PostgREST router
│   └── jwt-validator/
│       ├── Dockerfile           # Multi-stage Go build (~10MB)
│       ├── main.go              # Minimal JWT validator (stdlib only)
│       └── go.mod
├── glauth/
│   └── config.toml              # GLAuth LDAP config (users, groups)
├── authelia/
│   └── config/
│       └── configuration.yml    # Authelia config (LDAP backend, TOTP, session)
├── docker/
│   └── init/
│       ├── 01-roles.sql
│       ├── 02-schemas.sql
│       ├── 03-schema.sql
│       ├── 04-seeds.sql
│       ├── 05-comments.sql
│       └── 06-pre-request.sql   # set_user_context() for Traefik auth
├── scripts/
│   └── reload-postgrest-cache.sh
├── industry-packs/
│   └── gas/
│       ├── init/
│       │   ├── 01-schemas.sql
│       │   ├── 02-tables.sql
│       │   ├── 03-seeds.sql
│       │   ├── 04-grants.sql
│       │   └── 05-comments.sql
│       ├── docker-compose.override.yml
│       └── README.md
├── postman/
│   ├── GAS-FSM.postman_collection.json
│   ├── scripts/
│   │   ├── jwt-generator.js
│   │   └── generate-collection.js
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
| Auth | GLAuth (LDAP) → Authelia (MFA) + Backend Go (JWT) → Traefik → PostgREST |
| Naming convention | English (industry standard: snake_case) |