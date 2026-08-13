# FSM Postman Collection

Postman collection for testing all endpoints of the Go Backend API.

## Files

```
postman/
├── FSM.postman_collection.json          # Main collection (auto-generated)
├── FSM-dev.postman_environment.json     # Dev environment (localhost:8081)
├── FSM-staging.postman_environment.json # Staging environment
├── FSM-prod.postman_environment.json    # Production environment
└── scripts/
    ├── generate-collection.js           # Collection generator script
    └── jwt-generator.js                 # Pre-request JWT generator
```

## Quick Start

### 1. Import Collection
- Open Postman
- Click **Import** → **File** → Select `FSM.postman_collection.json`

### 2. Import Environment
- Click **Environments** → **Import** → Select `FSM-dev.postman_environment.json`
- Set `JWT_SECRET` to your actual secret

### 3. Select Environment
- Top-right dropdown → Select **FSM-dev**

### 4. Login
- Expand **Auth** folder → Click **Login**
- JWT is cached in `jwt_token` variable and auto-refreshed

### 5. Start Querying
- Expand any schema folder (e.g., **operations**)
- Click any request to load it
- Click **Send**

## Authentication

JWT is generated client-side via pre-request script:
- Uses `JWT_SECRET` from environment
- Generates HS256 JWT with `sub` (employee_number) and `role` (user_role)
- Token is cached in `jwt_token` and auto-refreshed on expiry

**Auth flow**:
```
Postman → Pre-request script → JWT → Go Backend (/api/*) → validates JWT
Go Backend (direct login) → POST /auth/login → LDAP bind → returns JWT
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `base_url` | Go Backend URL | `http://localhost:8081` |
| `JWT_SECRET` | JWT signing secret | `dev-jwt-secret-change-in-prod` |
| `employee_number` | Employee number for JWT sub | `1001` |
| `user_role` | User role for JWT | `admin` |
| `jwt_ttl_minutes` | JWT expiry in minutes | `60` |
| `jwt_token` | Generated JWT (auto-populated) | — |

## Collection Structure

```
FSM Go Backend API
├── Auth (3 requests)
│   ├── Login
│   ├── Refresh Token
│   └── Me
├── core (2 resources)
│   ├── users
│   └── tech_roles
├── partners (6 resources)
│   ├── partners
│   ├── partner_contacts
│   ├── partner_agreements
│   ├── partner_agreement_docs
│   ├── partner_agreement_forms
│   └── slas
├── customers (6 resources)
│   ├── customers
│   ├── customer_addresses
│   ├── customer_address_partners
│   ├── properties
│   ├── technicians
│   └── tech_certifications
├── inventory (9 resources)
│   ├── materials, tools, epp_items
│   ├── vehicles, rentals, cost_rates
│   ├── maintenance_records, maintenance_schedules
│   └── checklist_templates
├── operations (21 resources)
│   ├── visits
│   ├── visit_assignments
│   ├── checklist_template_materials/tools/epps
│   ├── visit_checklist_materials/tools/epps
│   ├── visit_material_usages/tool_usages/epp_usages
│   ├── visit_measurements, visit_checkpoints, visit_photos
│   ├── visit_reports, report_images, report_entries
│   ├── vehicle_assignments
│   ├── visit_rentals
│   └── visit_sla_trackings
├── planning (6 resources)
│   ├── daily_plans, daily_plan_assignments
│   ├── daily_load_materials/tools/epps
│   └── routes
├── notifications (2 resources)
│   ├── notification_templates
│   └── notifications
├── geocoding (1 resource)
│   └── addresses
└── shared (1 resource)
    └── report_templates
```

## Request Types

Each resource includes:

| Request | Method | Description |
|---------|--------|-------------|
| List | GET | `GET /api/{resource}?limit=20&offset=0` |
| Get by ID | GET | `GET /api/{resource}/:id` |
| Create | POST | `POST /api/{resource}` (body: JSON) |
| Update | PUT | `PUT /api/{resource}/:id` (body: JSON) |
| Delete | DELETE | `DELETE /api/{resource}/:id` |

## Regenerating Collection

If you modify the schema, regenerate the collection:

```bash
cd postman
node scripts/generate-collection.js FSM.postman_collection.json
```
