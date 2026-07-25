# GAS-FSM Postman Collection

Postman collection for testing all PostgREST endpoints of the GAS-FSM API.

## Files

```
postman/
├── GAS-FSM.postman_collection.json          # Main collection (370 requests, 65 tables)
├── GAS-FSM-dev.postman_environment.json     # Dev environment (localhost:3000)
├── GAS-FSM-staging.postman_environment.json # Staging environment
├── GAS-FSM-prod.postman_environment.json    # Production environment
└── scripts/
    ├── generate-collection.js               # Collection generator script
    └── jwt-generator.js                     # Pre-request JWT generator
```

## Quick Start

### 1. Import Collection
- Open Postman
- Click **Import** → **File** → Select `GAS-FSM.postman_collection.json`

### 2. Import Environment
- Click **Environments** → **Import** → Select `GAS-FSM-dev.postman_environment.json`
- Set `PGRST_JWT_SECRET` to your actual secret

### 3. Select Environment
- Top-right dropdown → Select **GAS-FSM-dev**

### 4. Generate JWT
- Expand **Auth** folder → Click **Generate JWT (dev)**
- JWT is now cached in `jwt_token` variable

### 5. Start Querying
- Expand any schema folder (e.g., **operations**)
- Click any request to load it
- Click **Send**

## Authentication

JWT is generated automatically via pre-request script:
- Uses `PGRST_JWT_SECRET` from environment
- Generates HS256 JWT with `sub` (employee_number) and `role` (user_role)
- Token is cached in `jwt_token` and auto-refreshed on expiry

**Auth flow**:
```
GLAuth (LDAP) → Authelia (MFA) → Traefik → Backend Go → JWT
Postman → Pre-request script → JWT → Authorization header → PostgREST
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `base_url` | PostgREST URL | `http://localhost:3000` |
| `PGRST_JWT_SECRET` | JWT signing secret | `dev-jwt-secret-change-in-prod` |
| `employee_number` | Employee number (GLAuth uidnumber) for JWT sub | `1001` |
| `user_role` | User role for JWT | `admin` |
| `jwt_ttl_minutes` | JWT expiry in minutes | `60` |
| `jwt_token` | Generated JWT (auto-populated) | — |

## Collection Structure

```
GAS-FSM PostgREST API
├── Auth (1 request)
│   └── Generate JWT (dev)
├── core (3 tables × 5 requests)
│   ├── users
│   ├── tech_roles
│   └── plan_audit_log (read-only)
├── partners (6 tables × 7 requests)
│   ├── partners (+ embed contacts/agreements/slas)
│   ├── partner_contacts
│   ├── partner_agreements (+ embed docs/forms)
│   ├── partner_agreement_docs
│   ├── partner_agreement_forms
│   └── slas
├── customers (6 tables × 7 requests)
│   ├── customers (+ embed addresses)
│   ├── customer_addresses (+ embed partner_assignments)
│   ├── customer_address_partners
│   ├── properties (+ embed property_assets)
│   ├── technicians (+ embed certifications)
│   └── tech_certifications
├── inventory (9 tables × 5 requests)
│   ├── materials, tools, epp_items
│   ├── vehicles (+ embed vehicle_types)
│   ├── rentals, cost_rates
│   ├── maintenance_records, maintenance_schedules
│   └── checklist_templates (+ embed materials/tools/epps)
├── operations (21 tables × 7 requests)
│   ├── visits (+ embed 14 relationships)
│   ├── visit_assignments
│   ├── checklist_template_materials/tools/epps
│   ├── visit_checklist_materials/tools/epps
│   ├── visit_material_usages/tools/epps
│   ├── visit_measurements
│   ├── visit_checkpoints (PostGIS)
│   ├── visit_photos (PostGIS)
│   ├── visit_reports (+ embed images/entries)
│   ├── report_images/entries
│   ├── vehicle_assignments
│   ├── visit_rentals
│   ├── visit_sla_trackings
│   └── SPECIAL: Full Visit Embed (14 relationships)
├── planning (6 tables × 5 requests)
│   ├── daily_plans (+ embed assignments/loads/routes)
│   ├── daily_plan_assignments
│   ├── daily_load_materials/tools/epps
│   └── routes
├── notifications (2 tables × 5 requests)
│   ├── notification_templates
│   └── notifications
├── geocoding (1 table × 6 requests)
│   ├── addresses
│   └── SPECIAL: Geo Query (PostGIS st_dwithin, st_distance)
├── shared (1 table × 5 requests)
│   └── report_templates
└── domain_gas (11 tables × 5 requests)
    ├── visit_types, measurement_types, photo_findings
    ├── vehicle_types, partner_service_types
    ├── pre_visit_results, route_types, property_types
    ├── certifications, rejection_reasons
    └── property_assets
```

## Request Types

Each table includes:

| Request | Method | Description |
|---------|--------|-------------|
| List | GET | `/{schema}/{table}?limit=20&offset=0` |
| List (select *) | GET | `/{schema}/{table}?select=*` |
| Get by ID | GET | `/{schema}/{table}?id=eq.{{id}}` |
| Embed | GET | `/{schema}/{table}?select=*,<related>(*)` |
| Create | POST | `/{schema}/{table}` (body: JSON) |
| Update | PATCH | `/{schema}/{table}?id=eq.{{id}}` (body: JSON) |
| Delete | DELETE | `/{schema}/{table}?id=eq.{{id}}` |

## PostGIS Geo Queries

### Within 1km of Santiago

```http
GET /geocoding/addresses?select=*,st_distance(geom,ST_SetSRID(ST_MakePoint(-70.6483,-33.4569),4326)::geography)&geom=st_dwithin(geom,ST_SetSRID(ST_MakePoint(-70.6483,-33.4569),4326)::geography,1000)&order=geom asc
```

### Bounding Box

```http
GET /geocoding/address?geom=st_intersects(geom,ST_MakeEnvelope(-70.7,-33.5,-70.6,-33.4,4326))
```

## Special Embeds

### Visits with All Relationships

```http
GET /operations/visits?select=*,
  visit_assignments(*,technician:tech_id(name,specialties)),
  visit_checklist_materials(*,material:material_id(name,unit,unit_cost)),
  visit_checklist_tools(*,tool:tool_id(name,code)),
  visit_checklist_epps(*,epp:epp_id(name,type,lifecycle)),
  visit_material_usages(*,material:material_id(name)),
  visit_tool_usages(*,tool:tool_id(name)),
  visit_epp_usages(*,epp:epp_id(name)),
  visit_measurements(*,mtype:type(code,name,default_unit)),
  visit_checkpoints(*),
  visit_photos(*,finding:finding_type(code,name)),
  visit_reports(*,report_images(*),report_entries(*)),
  vehicle_assignments(*,vehicle:vehicle_id(name,license_plate,vtype:type(code,name))),
  visit_rentals(*,rental:rental_id(item_description,daily_cost)),
  visit_sla_trackings(*,sla:sla_id(name,response_hours,resolution_hours))
```

## Regenerating Collection

If you modify the schema, regenerate the collection:

```bash
cd postman
node scripts/generate-collection.js GAS-FSM.postman_collection.json
```

## PostgREST Filtering

| Operator | Example | Description |
|----------|---------|-------------|
| `eq` | `?status=eq.active` | Equals |
| `neq` | `?status=neq.cancelled` | Not equals |
| `gt` | `?created_at=gt.2024-01-01` | Greater than |
| `gte` | `?year=gte.2023` | Greater or equal |
| `lt` | `?cost=lt.100000` | Less than |
| `lte` | `?cost=lte.50000` | Less or equal |
| `like` | `?name=like.*gas*` | Pattern match |
| `in` | `?status=in.(active,inactive)` | In list |
| `is` | `?notes=is.null` | Is null |

## PostgREST Pagination

```http
GET /operations/visits?limit=20&offset=0
Range: 0-19
Prefer: count=exact
```

Response headers:
- `Content-Range: 0-19/150`
- `X-Total-Count: 150`
