# Partner-Specific Customization Spec

## Delta: New capability — partner-scoped templates, auto-resolution, custom data

## Context
Different partners require different checklists, reports, SLAs, and data fields. Currently, templates are global and the partner relationship is barely used in visit/checklist/report flows.

## Functional Requirements

### FR-1: Partner-Scoped Checklist Templates
- New field: `ChecklistTemplate.PartnerID UUID` (nullable)
- NULL = global template (applies to all partners)
- UUID = template specific to that partner
- Filtering: when partner selected on visit, show partner-specific templates first, then globals
- Partner-specific templates override globals (same name = partner wins)

### FR-2: Partner-Scoped Report Templates
- New field: `ReportTemplate.PartnerID UUID` (nullable)
- New field: `ReportTemplate.WorkType TEXT` (nullable)
- Filtering: when creating report, filter by partner + work_type
- Partner-specific reports appear first in selection

### FR-3: Auto-Resolve Partner on Visit Creation
- When planner selects a property, backend resolves partner via `customer_address_partners`
- `partner_id` auto-set on visit
- If multiple partners: show selector
- If no partner: leave null (manual assignment later)

### FR-4: Auto-Resolve SLA
- When visit has `partner_id` + visit `type` (work_type)
- Backend searches `partners.slas` for active SLA matching partner + work_type
- Auto-creates `VisitSLATracking` record
- If no matching SLA: skip (not blocking)
- If multiple matches: use most specific (partner+work_type over partner-only)

### FR-5: Partner Data Fields on Visit
- `partner_order_id` — external order reference (shown prominently if partner assigned)
- `partner_supervisor` — partner's supervisor name
- These fields are visible/editable only when partner is assigned
- Mobile shows partner info at top of visit detail

### FR-6: Partner Contacts Reference
- When visit has partner_id, show partner contacts as reference
- Endpoint: `GET /api/partners/:id/contacts`
- Mobile: collapsible section "Contactos del partner"

## Non-Functional Requirements

### NFR-1: Backward Compatibility
- PartnerID on templates is nullable (NULL = global)
- Existing templates continue to work as globals
- No breaking changes to existing API contracts

### NFR-2: Performance
- Partner resolution happens once at visit creation
- Template filtering is client-side after initial load
- SLA resolution is O(1) lookup (indexed query)

## API Contract

### GET /api/templates?partner_id=X&type=checklist
```json
// Response
{
  "templates": [
    { "id": "...", "name": "Fibra Óptica - Instalación", "partner_id": "...", "stages_enabled": true },
    { "id": "...", "name": "Instalación General", "partner_id": null, "stages_enabled": true }
  ]
}
```

### GET /api/properties/:id/partners
```json
// Response
{
  "partners": [{
    "id": "...",
    "name": "Telefónica",
    "address_role": "primary",
    "is_active": true
  }]
}
```

### GET /api/partners/:id/contacts
```json
// Response
{
  "contacts": [{
    "name": "Carlos López",
    "role": "Supervisor",
    "phone": "+56 9 1234 5678",
    "email": "carlos@telefonica.com"
  }]
}
```
