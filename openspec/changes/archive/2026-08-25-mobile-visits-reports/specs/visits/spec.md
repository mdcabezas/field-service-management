# Visits Spec

## Delta: Extended — partner context, property history, reopen, audit, checklist validation

## Context
Visits are the core entity in the mobile app. We extend the existing visit model with partner-specific data, property history, reopen capability, and audit trail. Checklists are now validated before status transitions.

## Functional Requirements

### FR-1: Visit List (Mobile)
- Show technician's assigned visits for today
- Display: property name, partner name, scheduled time, status, priority
- Pull-to-refresh triggers delta sync
- Paginated: load 10 visits at a time
- Filter by status (all, assigned, in_progress, completed)

### FR-2: Visit Detail (Lazy Loaded)
- Visit details fetched when user opens visit (not on list load)
- Includes: property info, partner info, assignments, reports
- Partner info: name, order ID, supervisor, contacts
- Property coordinates for GPS validation
- Property history (previous visits)
- Checklist by stages
- Audit trail

### FR-3: Visit Creation (with partner resolution)
- Planner selects property → partner auto-resolved from `customer_address_partners`
- Planner selects template → filtered by partner (partner-specific first, then globals)
- Planner selects stages → from template's stages
- SLA auto-linked if partner + work_type match
- partner_order_id and partner_supervisor editable

### FR-4: Status Transitions
- Validated against checklist stages
- Required stages must be complete before transitioning to "completada"
- Emergency visits can skip pre-departure
- All transitions audited

### FR-5: Property History
- `GET /api/properties/:id/history` returns previous visits
- Shows: date, type, status, technician name
- Limited to last 20 visits

### FR-6: Reopen Visit
- `POST /api/visits/:id/reopen` with reason
- Within configurable window (default 24h from completion)
- Requires reason text
- Full audit trail preserved
- Allows adding photos, modifying data
- Status reverts to "en_curso"

### FR-7: Audit Trail
- `GET /api/visits/:id/audit` returns all changes
- Shows: timestamp, user, action, old_value, new_value
- Paginated (50 per page)

### FR-8: Partner Data
- partner_order_id: external order reference (shown prominently if partner assigned)
- partner_supervisor: partner's supervisor name
- These fields visible/editable only when partner is assigned

## Non-Functional Requirements

### NFR-1: Offline Support
- Visit list cached in SQLite
- Visit detail cached on first load
- Changes queued in outbox
- Conflict detection on status changes

### NFR-2: Performance
- Visit list loads in <500ms from SQLite
- Visit detail loads in <1s (lazy loaded)
- Property history loads in <500ms

## API Contract

### GET /api/mobile/visits?status=X&page=N&limit=10
```json
// Response
{
  "visits": [{
    "id": "...",
    "property": { "id": "...", "name": "...", "address": "...", "lat": ..., "lng": ... },
    "partner": { "id": "...", "name": "...", "order_id": "...", "supervisor": "..." },
    "type": "installation",
    "status": "assigned",
    "priority": "normal",
    "scheduled_at": "2026-01-15T14:00:00Z",
    "assignments": [{ "technician_name": "..." }]
  }],
  "total": 25,
  "page": 1,
  "limit": 10
}
```

### GET /api/mobile/visits/:id
```json
// Response
{
  "id": "...",
  "property": { "id": "...", "name": "...", "address": "...", "lat": ..., "lng": ... },
  "partner": { "id": "...", "name": "...", "order_id": "...", "supervisor": "..." },
  "partner_contacts": [{ "name": "...", "role": "...", "phone": "..." }],
  "type": "installation",
  "status": "assigned",
  "priority": "normal",
  "scheduled_at": "2026-01-15T14:00:00Z",
  "checklist_template_id": "...",
  "stages_enabled": true,
  "selected_stages": ["stage-id-1", "stage-id-2"],
  "sla": { "response_hours": 4, "resolution_hours": 24 },
  "reports": [],
  "audit": []
}
```

### POST /api/visits/:id/reopen
```json
// Request
{ "reason": "Need to add additional photos" }

// Response
{ "id": "...", "status": "en_curso", "reopened_at": "...", "reopen_reason": "..." }
```

### GET /api/properties/:id/history
```json
// Response
{
  "visits": [{
    "id": "...",
    "type": "installation",
    "status": "completada",
    "technician": "Juan Pérez",
    "completed_at": "2025-12-10T16:30:00Z"
  }],
  "total": 5
}
```

### GET /api/visits/:id/audit?page=N&limit=50
```json
// Response
{
  "entries": [{
    "id": "...",
    "timestamp": "2026-01-15T14:30:00Z",
    "user": "Juan Pérez",
    "action": "status_change",
    "field": "status",
    "old_value": "assigned",
    "new_value": "en_curso"
  }],
  "total": 15,
  "page": 1
}
```
