# Gas Chile Industry Pack

Industry-specific extensions for FSM core, providing lookup tables, SEC certifications, and rejection reasons for the Chilean gas sector.

## Structure

```
industry-packs/gas/
├── docker-compose.override.yml    # Dev: mounts initdb.d/ to postgres
├── docker-compose.prod.yml        # Prod: mounts initdb.d/ to postgres
├── initdb.d/                      # Auto-executed by PostgreSQL entrypoint
│   ├── 90-gas-01-schemas.sql      # Schema + tables + FKs (idempotent)
│   ├── 91-gas-02-seeds.sql        # Lookup data
│   ├── 92-gas-03-grants.sql       # Grants + search_path
│   └── 93-gas-04-comments.sql     # OpenAPI comments
└── README.md
```

## Tables Added

| Table | Purpose |
|-------|---------|
| `domain_gas.visit_types` | 7 visit types (pre_visit, installation, maintenance, emergency, certification, repair, diagnosis) |
| `domain_gas.measurement_types` | 8 measurement types (tightness, pressure, CO, draft, leak, pH, temperature, other) |
| `domain_gas.photo_findings` | 6 photo finding types (normal, leak, damage, emergency, incomplete, other) |
| `domain_gas.vehicle_types` | 5 vehicle types (truck, crane, van, crane_truck, other) |
| `domain_gas.partner_service_types` | 5 service types (installation, maintenance, certification, emergency, other) |
| `domain_gas.pre_visit_results` | 3 results (approved, rejected, conditional) |
| `domain_gas.route_types` | 6 route types (meter_reading, letter_delivery, tank_collection, tank_delivery, mass_inspection, mixed) |
| `domain_gas.property_types` | 5 property types with extra column config (tank, meter, indoor_piping, appliance, other) |
| `domain_gas.certifications` | 9 SEC Chile certifications (CL1, CL2, CL3, GLP, TC1, TC2, TC6, GREEN SEAL, DS66) |
| `domain_gas.rejection_reasons` | 17 reasons in 4 categories (logistics, regulatory, customer, operational) |
| `domain_gas.property_assets` | Extra property attributes (poles, meter, etc.) |

## Foreign Keys Added to Core

The industry pack adds FK constraints from core tables to `domain_gas` lookup tables:

| Core Table | Column | References |
|------------|--------|------------|
| `partners.partner_agreements` | `service_type` | `domain_gas.partner_service_types(id)` |
| `partners.slas` | `work_type` | `domain_gas.visit_types(id)` |
| `customers.properties` | `type` | `domain_gas.property_types(id)` |
| `customers.tech_certifications` | `cert_id` | `domain_gas.certifications(id)` |
| `inventory.vehicles` | `type` | `domain_gas.vehicle_types(id)` |
| `inventory.checklist_templates` | `work_type` | `domain_gas.visit_types(id)` |
| `planning.routes` | `type` | `domain_gas.route_types(id)` |
| `operations.visits` | `type` | `domain_gas.visit_types(id)` |
| `operations.visits` | `result` | `domain_gas.pre_visit_results(id)` |
| `operations.visits` | `rejection_reason_id` | `domain_gas.rejection_reasons(id)` |
| `operations.visit_measurements` | `type` | `domain_gas.measurement_types(id)` |
| `operations.visit_photos` | `finding_type` | `domain_gas.photo_findings(id)` |

## Usage

### Development

```bash
# From project root
docker compose -f docker-compose.yml -f industry-packs/gas/docker-compose.override.yml up -d
```

### Production

```bash
# From project root
docker compose -f docker-compose.prod.yml -f industry-packs/gas/docker-compose.prod.yml --env-file .env.prod up -d
```

### Manual Execution (if needed)

```bash
# Copy files to container
docker cp industry-packs/gas/initdb.d/*.sql fsm-postgres:/tmp/

# Execute each file
docker compose exec postgres psql -U fsm_admin -d fsm_fsm -f /tmp/90-gas-01-schemas.sql
docker compose exec postgres psql -U fsm_admin -d fsm_fsm -f /tmp/91-gas-02-seeds.sql
docker compose exec postgres psql -U fsm_admin -d fsm_fsm -f /tmp/92-gas-03-grants.sql
docker compose exec postgres psql -U fsm_admin -d fsm_fsm -f /tmp/93-gas-04-comments.sql
```

## Creating New Industry Packs

1. Copy `industry-packs/gas/` to `industry-packs/<your-industry>/`
2. Modify lookup tables in `90-<industry>-01-schemas.sql`
3. Replace seeds in `91-<industry>-02-seeds.sql`
4. Update grants in `92-<industry>-03-grants.sql`
5. Update comments in `93-<industry>-04-comments.sql`
6. Update `docker-compose.override.yml` and `docker-compose.prod.yml`

## Core vs Industry Pack Philosophy

- **Core (fsm):** Generic FSM tables, enums for universal concepts (status, priority, etc.), no industry-specific data
- **Industry Pack:** All industry-specific lookups, certifications, regulations, and FKs linking to core

This allows the same core deployment to serve gas, electricity, HVAC, plumbing, etc. by simply loading different industry packs.
