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
| `domain_gas.measurement_types` | 8 measurement types (tightness, pressure, CO, draft, leak, pH, temperature, other) |
| `domain_gas.property_types` | 5 property types with extra column config (tank, meter, indoor_piping, appliance, other) |
| `domain_gas.certifications` | 9 SEC Chile certifications (CL1, CL2, CL3, GLP, TC1, TC2, TC6, GREEN SEAL, DS66) |
| `domain_gas.property_assets` | Extra property attributes (poles, meter, etc.) |

Industry-agnostic catalogs live in the core and are seeded there (see `docker/init/04-seeds.sql`); this pack adds its specific codes into those core tables via `91-gas-02-seeds.sql`:

| Core catalog | Gas codes added |
|--------------|-----------------|
| `operations.visit_types` | pre_visit, installation, maintenance, emergency, certification, repair, diagnosis |
| `operations.photo_findings` | normal, leak, damage, emergency, incomplete, other |
| `operations.pre_visit_results` | approved, rejected, conditional |
| `operations.rejection_reasons` | 17 reasons in 4 categories (logistics, regulatory, customer, operational) |
| `partners.partner_service_types` | installation, maintenance, certification, emergency, other |
| `planning.route_types` | meter_reading, letter_delivery, tank_collection, tank_delivery, mass_inspection, mixed |

## Foreign Keys Added to Core

Universal catalogs were moved to core (`sql/20260813_catalogs_core.sql`); the pack keeps the FK wiring pointing to those core tables so every tenant, gas or not, ships the catalogs:

| Core Table | Column | References |
|------------|--------|------------|
| `partners.partner_agreements` | `service_type` | `partners.partner_service_types(id)` *(core)* |
| `partners.slas` | `work_type` | `operations.visit_types(id)` *(core)* |
| `customers.properties` | `type` | `domain_gas.property_types(id)` |
| `customers.tech_certifications` | `cert_id` | `domain_gas.certifications(id)` |
| `inventory.vehicles` | `type` | `inventory.vehicle_types(id)` *(core, seeded by 04-seeds.sql)* |
| `inventory.checklist_templates` | `work_type` | `operations.visit_types(id)` *(core)* |
| `planning.routes` | `type` | `planning.route_types(id)` *(core catalog; gas seeds its codes via 91-gas-02-seeds.sql)* |
| `operations.visits` | `type` | `operations.visit_types(id)` *(core)* |
| `operations.visits` | `result` | `operations.pre_visit_results(id)` *(core)* |
| `operations.visits` | `rejection_reason_id` | `operations.rejection_reasons(id)` *(core)* |
| `operations.visit_measurements` | `type` | `domain_gas.measurement_types(id)` |
| `operations.visit_photos` | `finding_type` | `operations.photo_findings(id)` *(core)* |

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
