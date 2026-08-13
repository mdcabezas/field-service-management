# FSM Schema Documentation

Field Service Management system for multi-industry field service operations.

**Stack**: PostgreSQL 16 + PostGIS 3.5 + PostgREST 12 + Go backend

**Mode**: Single-tenant

**Core tables**: ~49

**Industry packs**: Load additional lookup tables, certifications, and FKs per industry

## Conventions

| Element | Convention |
|---------|-----------|
| Names | `snake_case` |
| Primary keys | `UUID DEFAULT gen_random_uuid()` |
| Timestamps | `TIMESTAMPTZ DEFAULT now()` |
| Enums | Native PostgreSQL types in `shared` schema (27) |
| Industry lookup tables | Tables in `domain_<industry>` schema with UUID PK |
| PostGIS | `GEOMETRY(Point, 4326)` for locations |

## Architecture Decisions

| ID | Decision | Value |
|----|----------|-------|
| D1 | `addresses` | **Global catalog** in `geocoding` schema |
| D2 | `route_id`↔`daily_plan_id` | **TRIGGER constraint** on `visits` |
| D3 | `vehicle_assignments` | **Single table** with optional `daily_plan_id` and `visit_id` |
| D4 | Native enums | **27 PostgreSQL enums** in `shared` |
| D5 | Industry packs | **`domain_<industry>` schemas** with lookup tables, certifications, FKs |
| A | Item polymorphism | **Option A**: split into `_materials`, `_tools`, `_epps` |
| C | Domain schema | **Industry packs** isolate industry-specific elements (enums→tables, certifications, seeds) |

## Authentication

JWT signed by Go backend, consumed by PostgREST via `set_user_context()`.

**JWT payload**:
```json
{
  "sub": "1001",
  "role": "admin|operator|technician",
  "iat": 1700000000,
  "exp": 1700086400
}
```

- `sub`: employee_number (GLAuth uidnumber)
- `role`: single role per user (mapped from GLAuth group)

**Auth flow**:
```
Client → Go Backend → LDAP bind (GLAuth) → JWT
API → Traefik → Go Backend → PostgREST → set_user_context()
```

## Schemas (bounded contexts)

### Core (9 schemas)

| Schema | Tables | Possible future microservice |
|--------|--------|-------------------------------|
| `core` | users, tech_roles, plan_audit_log | IAM/Audit |
| `partners` | partners, partner_contacts, partner_agreements, partner_agreement_forms, partner_agreement_docs, slas | Partners |
| `customers` | customers, customer_addresses, customer_address_partners, properties, technicians, tech_certifications | CRM |
| `inventory` | materials, tools, epp_items, vehicles, rentals, cost_rates, maintenance_records, maintenance_schedules, checklist_templates | Inventory |
| `operations` | visits, visit_assignments, visit_checklist_materials/_tools/_epps, visit_material_usages, visit_tool_usages, visit_epp_usages, visit_measurements, visit_checkpoints, visit_photos, visit_reports, report_images, report_entries, vehicle_assignments, visit_rentals, visit_sla_trackings, checklist_template_materials/_tools/_epps | Operations |
| `planning` | daily_plans, daily_plan_assignments, daily_load_materials/_tools/_epps, routes | Planner |
| `notifications` | notification_templates, notifications | Notifier |
| `geocoding` | addresses | Geo |
| `shared` | report_templates, rejection_category (enum) | — (generic catalog) |

### Industry Pack: Gas Chile (`domain_gas`)

See [industry-packs/gas/README.md](../industry-packs/gas/README.md) for full documentation.

---

## Core Schemas

### `core`

#### `core.users`

| Field | Type | Description |
|-------|------|-------------|
| employee_number | TEXT PK | GLAuth uidnumber, JWT sub claim |
| email | TEXT | Unique login email |
| role | TEXT | admin, operator, technician |
| name | TEXT | Full name |
| is_active | BOOLEAN | User active status |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

#### `core.tech_roles`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| code | TEXT UNIQUE | |
| name | TEXT | |
| active | BOOLEAN | |
| created_at | TIMESTAMPTZ | |

Seed: driver, technician, helper, operative.

#### `core.plan_audit_log`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| entity_type | shared.audit_entity_type | |
| entity_id | UUID | |
| action | shared.audit_action | add/update/delete |
| old_data | JSONB | Previous state (NULL on INSERT) |
| new_data | JSONB | New state (NULL on DELETE) |
| user_id | TEXT | employee_number from JWT sub |
| reason | TEXT | |
| timestamp | TIMESTAMPTZ | |

---

### `partners`

#### `partners.partners`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| name | TEXT | Partner name |
| tax_id | TEXT | |
| status | shared.partner_status | active/inactive |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | (trigger) |

#### `partners.partner_contacts`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| partner_id | UUID FK → partners | |
| name | TEXT | |
| position | TEXT | |
| phone | TEXT | |
| email | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `partners.partner_agreements`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| partner_id | UUID FK → partners | |
| service_type | UUID | FK added by industry pack |
| rate | NUMERIC | |
| active | BOOLEAN | |
| created_at | TIMESTAMPTZ | |

#### `partners.partner_agreement_docs`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| agreement_id | UUID FK → partner_agreements | |
| work_type | TEXT | |
| doc_name | TEXT | |
| required | BOOLEAN | |
| created_at | TIMESTAMPTZ | |

#### `partners.partner_agreement_forms`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| agreement_id | UUID FK → partner_agreements | |
| work_type | TEXT | |
| form_template_id | UUID FK → shared.report_templates | |
| quantity | INTEGER | |
| created_at | TIMESTAMPTZ | |

#### `partners.slas`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| partner_id | UUID FK → partners | |
| name | TEXT | |
| description | TEXT | |
| work_type | UUID | FK added by industry pack |
| response_hours | INTEGER | |
| resolution_hours | INTEGER | |
| compliance_target | NUMERIC | |
| active | BOOLEAN | |
| valid_from | DATE | |
| valid_until | DATE | |
| created_at | TIMESTAMPTZ | |

---

### `customers`

#### `customers.customers`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| name | TEXT | |
| phone | TEXT | |
| email | TEXT | |
| tax_id | TEXT | |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | (trigger) |

#### `customers.customer_addresses`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| customer_id | UUID FK → customers | |
| address_id | UUID FK → geocoding.addresses | |
| name | TEXT | |
| type | shared.customer_address_type | |
| created_at | TIMESTAMPTZ | |

#### `customers.customer_address_partners`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| customer_address_id | UUID FK → customer_addresses | |
| partner_id | UUID FK → partners | |
| from_date | DATE | |
| to_date | DATE | NULL = active |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `customers.properties`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| customer_address_id | UUID FK → customer_addresses | |
| name | TEXT | |
| type | UUID | FK added by industry pack |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

Extra attributes (poles, meter, etc.) in `domain_<industry>.property_assets`.

#### `customers.technicians`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| user_id | TEXT FK → core.users(employee_number) | |
| name | TEXT | |
| is_active | BOOLEAN | |
| specialties | TEXT[] | |
| created_at | TIMESTAMPTZ | |

#### `customers.tech_certifications`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| tech_id | UUID FK → technicians | |
| cert_id | UUID | FK added by industry pack |
| number | TEXT | |
| issuer | TEXT | |
| issue_date | DATE | |
| expiry_date | DATE | |
| created_at | TIMESTAMPTZ | |

---

### `inventory`

#### `inventory.materials`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| sku | TEXT | |
| name | TEXT | |
| unit | TEXT | |
| unit_cost | NUMERIC | |
| created_at | TIMESTAMPTZ | |

#### `inventory.tools`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| code | TEXT | |
| name | TEXT | |
| status | shared.tool_status | available/in_use/maintenance/retired |
| created_at | TIMESTAMPTZ | |

#### `inventory.epp_items`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| name | TEXT | |
| type | shared.epp_type | |
| lifecycle | shared.epp_lifecycle | |
| created_at | TIMESTAMPTZ | |

#### `inventory.vehicles`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| type | UUID | FK added by industry pack |
| license_plate | TEXT | |
| name | TEXT | |
| brand | TEXT | |
| model | TEXT | |
| year | INTEGER | |
| status | shared.vehicle_status | |
| capacity | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `inventory.rentals`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| type | shared.rental_type | |
| supplier | TEXT | |
| item_description | TEXT | |
| daily_cost | NUMERIC | |
| hourly_cost | NUMERIC | |
| fixed_cost | NUMERIC | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `inventory.cost_rates`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| type | shared.cost_type | |
| reference_id | UUID | |
| reference_type | TEXT | |
| value | NUMERIC | |
| unit | shared.cost_unit | |
| valid_from | DATE | |
| valid_until | DATE | |
| created_at | TIMESTAMPTZ | |

#### `inventory.maintenance_records`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| type | shared.maintenance_type | |
| reference_id | UUID | |
| date | DATE | |
| cost | NUMERIC | |
| supplier | TEXT | |
| description | TEXT | |
| next_date | DATE | |
| created_at | TIMESTAMPTZ | |

#### `inventory.maintenance_schedules`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| type | shared.maintenance_type | |
| reference_id | UUID | |
| frequency_km | INTEGER | |
| frequency_days | INTEGER | |
| last_service_date | DATE | |
| last_service_km | INTEGER | |
| active | BOOLEAN | |
| created_at | TIMESTAMPTZ | |

#### `inventory.checklist_templates`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| name | TEXT | |
| work_type | UUID | FK added by industry pack |
| created_at | TIMESTAMPTZ | |

---

### `operations`

#### `operations.visits`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| property_id | UUID FK → properties | |
| partner_id | UUID FK → partners | |
| partner_order_id | TEXT | |
| parent_visit_id | UUID FK → visits | Parent visit (pre_visit or sub-visit) |
| route_id | UUID FK → routes | D2 CHECK |
| daily_plan_id | UUID FK → daily_plans | D2 CHECK |
| type | UUID NOT NULL | FK added by industry pack |
| status | shared.visit_status | |
| priority | shared.visit_priority | |
| source | shared.visit_source | |
| source_reference | TEXT | |
| billing_to | shared.billing_to | partner/customer |
| partner_supervisor | TEXT | |
| result | UUID | FK added by industry pack |
| rejection_reason_id | UUID | FK added by industry pack |
| rejection_reason_detail | TEXT | |
| scheduled_at | TIMESTAMPTZ NOT NULL | |
| started_at | TIMESTAMPTZ | |
| completed_at | TIMESTAMPTZ | |
| alternative_location | TEXT | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | (trigger) |

#### `operations.visit_assignments`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| tech_id | UUID FK → technicians | |
| role_id | UUID FK → core.tech_roles | |
| created_at | TIMESTAMPTZ | |

#### `operations.checklist_template_materials` (Option A)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| template_id | UUID FK → checklist_templates | |
| material_id | UUID FK → materials | |
| default_quantity | INTEGER | |
| required | BOOLEAN | |
| created_at | TIMESTAMPTZ | |

#### `operations.checklist_template_tools` (Option A)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| template_id | UUID FK → checklist_templates | |
| tool_id | UUID FK → tools | |
| default_quantity | INTEGER | |
| required | BOOLEAN | |
| created_at | TIMESTAMPTZ | |

#### `operations.checklist_template_epps` (Option A)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| template_id | UUID FK → checklist_templates | |
| epp_id | UUID FK → epp_items | |
| default_quantity | INTEGER | |
| required | BOOLEAN | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_checklist_materials` (Option A)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| material_id | UUID FK → materials | |
| planned_quantity | INTEGER | |
| confirmed | BOOLEAN | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_checklist_tools` (Option A)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| tool_id | UUID FK → tools | |
| planned_quantity | INTEGER | |
| confirmed | BOOLEAN | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_checklist_epps` (Option A)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| epp_id | UUID FK → epp_items | |
| planned_quantity | INTEGER | |
| confirmed | BOOLEAN | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_material_usages`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| material_id | UUID FK → materials | |
| quantity | NUMERIC | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_tool_usages`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| tool_id | UUID FK → tools | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_epp_usages`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| epp_id | UUID FK → epp_items | |
| quantity | INTEGER | |
| status | shared.epp_usage_status | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_measurements`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| type | UUID NOT NULL | FK added by industry pack |
| value | NUMERIC | |
| unit | TEXT | |
| result | shared.measurement_result | |
| measuring_device | TEXT | |
| notes | TEXT | |
| photo_url | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_checkpoints`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| type | shared.checkpoint_type | arrival/departure/report_sent |
| geom | GEOMETRY(Point, 4326) | GPS |
| timestamp | TIMESTAMPTZ NOT NULL DEFAULT now() | |
| device_info | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_photos`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| checkpoint_id | UUID FK → visit_checkpoints | |
| url | TEXT NOT NULL | |
| geom | GEOMETRY(Point, 4326) | |
| timestamp | TIMESTAMPTZ DEFAULT now() | |
| stage | shared.photo_stage | |
| finding_type | UUID | FK added by industry pack |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_reports`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| report_template_id | UUID FK → shared.report_templates | |
| source | shared.report_source | system/paper |
| recorded_at | TIMESTAMPTZ DEFAULT now() | |
| created_at | TIMESTAMPTZ | |

#### `operations.report_images`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| report_id | UUID FK → visit_reports | |
| url | TEXT NOT NULL | |
| pages | INTEGER DEFAULT 1 | |
| created_at | TIMESTAMPTZ | |

#### `operations.report_entries`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| report_id | UUID FK → visit_reports | |
| data_json | JSONB NOT NULL DEFAULT '{}' | |
| created_at | TIMESTAMPTZ | |

#### `operations.vehicle_assignments` (single table, D3)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| vehicle_id | UUID FK → vehicles | |
| daily_plan_id | UUID FK → daily_plans | |
| visit_id | UUID FK → visits | |
| departure_time | TIMESTAMPTZ | |
| return_time | TIMESTAMPTZ | |
| departure_mileage | INTEGER | |
| return_mileage | INTEGER | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_rentals`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| rental_id | UUID FK → rentals | |
| start_time | TIMESTAMPTZ NOT NULL | |
| end_time | TIMESTAMPTZ | |
| hours | NUMERIC | |
| total_cost | NUMERIC NOT NULL | |
| reason | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `operations.visit_sla_trackings`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| visit_id | UUID FK → visits | |
| sla_id | UUID FK → slas | |
| requested_at | TIMESTAMPTZ | |
| responded_at | TIMESTAMPTZ | |
| resolved_at | TIMESTAMPTZ | |
| response_time_hours | NUMERIC | |
| resolution_time_hours | NUMERIC | |
| meets_response_sla | BOOLEAN | |
| meets_resolution_sla | BOOLEAN | |
| created_at | TIMESTAMPTZ | |

---

### `planning`

#### `planning.daily_plans`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| date | DATE NOT NULL | |
| notes | TEXT | |
| created_at | TIMESTAMPTZ | |

#### `planning.daily_plan_assignments`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| daily_plan_id | UUID FK → daily_plans | |
| tech_id | UUID FK → technicians | |
| role_id | UUID FK → core.tech_roles | |
| created_at | TIMESTAMPTZ | |

#### `planning.daily_load_materials` (Option A)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| daily_plan_id | UUID FK → daily_plans | |
| material_id | UUID FK → materials | |
| loaded_quantity | NUMERIC NOT NULL | |
| returned_quantity | NUMERIC DEFAULT 0 | |
| created_at | TIMESTAMPTZ | |

#### `planning.daily_load_tools` (Option A)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| daily_plan_id | UUID FK → daily_plans | |
| tool_id | UUID FK → tools | |
| loaded_quantity | NUMERIC NOT NULL | |
| returned_quantity | NUMERIC DEFAULT 0 | |
| created_at | TIMESTAMPTZ | |

#### `planning.daily_load_epps` (Option A)

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| daily_plan_id | UUID FK → daily_plans | |
| epp_id | UUID FK → epp_items | |
| loaded_quantity | NUMERIC NOT NULL | |
| returned_quantity | NUMERIC DEFAULT 0 | |
| created_at | TIMESTAMPTZ | |

#### `planning.routes`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| type | UUID | FK added by industry pack |
| date | DATE NOT NULL | |
| daily_plan_id | UUID FK → daily_plans | |
| notes | TEXT | |
| status | shared.route_status | |
| created_at | TIMESTAMPTZ | |

---

### `notifications`

#### `notifications.notification_templates`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| type | shared.notification_type | |
| channel | shared.notification_channel | |
| subject | TEXT | |
| body | TEXT | |
| active | BOOLEAN | |
| created_at | TIMESTAMPTZ | |

#### `notifications.notifications`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| type | shared.notification_type | |
| channel | shared.notification_channel | |
| recipient | TEXT | |
| subject | TEXT | |
| body | TEXT | |
| status | shared.notification_status | |
| entity_type | TEXT | (polymorphic) |
| entity_id | UUID | (polymorphic) |
| scheduled_for | TIMESTAMPTZ | |
| sent_at | TIMESTAMPTZ | |
| read_at | TIMESTAMPTZ | |
| attempts | INTEGER DEFAULT 0 | |
| error | TEXT | |
| created_at | TIMESTAMPTZ | |

---

### `geocoding`

#### `geocoding.addresses`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| street | TEXT NOT NULL | |
| number | TEXT | |
| apartment | TEXT | |
| neighborhood | TEXT | |
| city | TEXT | |
| region | TEXT | |
| location_references | TEXT | |
| postal_code | TEXT | |
| geom | GEOMETRY(Point, 4326) | PostGIS |
| created_at | TIMESTAMPTZ | |

---

### `shared`

#### `shared.report_templates`

| Field | Type | Description |
|-------|------|-------------|
| id | UUID PK | |
| name | TEXT | |
| fields_json | JSONB | |
| created_at | TIMESTAMPTZ | |

Enum `shared.rejection_category` remains here (generic).

---

### `domain_gas` (Gas Chile Industry Pack)

See [industry-packs/gas/README.md](../industry-packs/gas/README.md) for table definitions, seeds, and FK constraints.

---

## Helper Functions

| Function | Type | Purpose |
|----------|------|---------|
| `app_user_id()` | SQL STABLE | Returns `employee_number` (TEXT) from JWT `sub` claim |
| `touch_updated_at()` | TRIGGER | Sets `NEW.updated_at = now()` |
| `audit_change()` | TRIGGER | Writes changes to `core.plan_audit_log` |
| `validate_visit_route_plan()` | TRIGGER | D2: validates route ↔ daily_plan consistency |

## Triggers

### Automatic updated_at
- `partners.partners`, `customers.customers`, `operations.visits`, `core.users`

### Audit (8 tables)
- `planning.daily_load_materials`, `daily_load_tools`, `daily_load_epps`
- `operations.visit_checklist_materials`, `visit_checklist_tools`, `visit_checklist_epps`
- `operations.visit_material_usages`, `visit_tool_usages`, `visit_epp_usages`

## Enums (27)

| Enum | Values |
|------|--------|
| visit_status | scheduled, en_route, in_progress, completed, cancelled, rescheduled |
| visit_priority | low, normal, urgent, emergency |
| visit_source | email, phone, portal, whatsapp, other |
| billing_to | partner, customer |
| photo_stage | diagnosis, execution, review, report, other |
| checkpoint_type | arrival, departure, report_sent |
| report_source | system, paper |
| measurement_result | approved, rejected, pending, observation |
| vehicle_status | available, in_use, maintenance, retired |
| tool_status | available, in_use, maintenance, retired |
| tech_status | active, inactive, vacation, leave |
| rental_type | vehicle, tool, equipment |
| route_status | scheduled, in_progress, completed, cancelled |
| epp_usage_status | used, not_needed, missing |
| rejection_category | logistics, regulatory, customer, operational |
| cost_type | labor, vehicle, tool_depreciation, epp, overhead |
| cost_unit | hour, km, day, visit, unit |
| customer_address_type | residential, commercial, industrial, institutional |
| partner_status | active, inactive |
| epp_type | head, hands, body, eyes, ears, feet, respiratory |
| epp_lifecycle | disposable, reusable |
| notification_type | visit_assigned, visit_reminder, visit_cancelled, visit_rescheduled, checklist_pending, report_pending, sla_warning, sla_expired, maintenance_scheduled, certification_expired, other |
| notification_channel | email, sms, whatsapp, push, in_app |
| notification_status | pending, sent, failed, read |
| audit_entity_type | daily_load_material, daily_load_tool, daily_load_epp, visit_checklist_material, visit_checklist_tool, visit_checklist_epp, visit_material_usage, visit_tool_usage, visit_epp_usage |
| audit_action | add, update, delete |
| maintenance_type | vehicle, tool, equipment |

## Local dev

```bash
cp .env.example .env
docker compose up -d
docker compose exec postgres psql -U fsm_admin -d fsm_fsm
```

## Load Industry Pack (Gas)

```bash
docker compose -f docker-compose.yml -f industry-packs/gas/docker-compose.override.yml up -d
```