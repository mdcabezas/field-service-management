# Sync Engine Spec

## Delta: New capability — offline-first bidirectional sync

## Context
The mobile app must work fully offline and sync when connectivity returns. The backend has existing endpoints but no sync protocol. We add a new sync protocol that supports full/delta sync, outbox pattern, and conflict detection.

## Functional Requirements

### FR-1: Full Sync
- On login, mobile fetches reference data (materials, tools, EPP, templates, partners, etc.) in a single response (~50KB)
- Then fetches assigned visits paginated (10 at a time, `?page=N&limit=10&status=assigned`)
- Each visit includes: basic info, property GPS, partner info, checklist template IDs
- Does NOT include checklist items, photos, or audit (lazy loaded)

### FR-2: Delta Sync
- On reconnect, mobile pushes outbox items first, then pulls changes
- Push: `POST /api/mobile/sync` with body `{ checkpoints: [], photos: [], usages: [], measurements: [], reports: [], status_changes: [] }`
- Pull: `GET /api/mobile/sync?since=<timestamp>` returns changes since last sync
- Delta sync payload is small (<10KB typical)

### FR-3: Outbox Pattern
- All local changes stored in SQLite `outbox` table with dependency ordering
- Push order: checkpoints → photos → usages → measurements → reports → status changes
- Each outbox item: `id, entity_type, entity_id, operation, payload, status, created_at`
- Failed items retry on next sync attempt

### FR-4: Conflict Detection
- Status conflicts: server wins (conservative approach for safety-critical transitions)
- Data conflicts: client wins (technician's local edits take precedence)
- Conflict log stored for audit trail
- No merge — simple resolution strategy

### FR-5: Manual Retry
- Failed sync items shown in a retry UI
- User can manually trigger retry for individual items or all failures
- Retry clears error state and requeues the item

### FR-6: Sync Status Indicators
- Green: fully synced
- Yellow: sync in progress
- Red: sync failed (with error count)
- Grey: offline (no connection)

### FR-7: Technician Identification
- Backend identifies the technician via `user_id` from the JWT `sub` claim
- Query: `SELECT id FROM customers.technicians WHERE user_id = $1`
- MUST NOT use `employee_number` (which is `EMP-XXX` format, not a UUID)
- If no technician found for the user_id, return 404 with message "technician not found"

## Non-Functional Requirements

### NFR-1: Performance
- Full sync completes in <10 seconds on 3G
- Delta sync completes in <5 seconds on 3G
- Push payload compressed (gzip) to minimize data transfer

### NFR-2: Reliability
- Outbox items survive app restart
- Sync resumes from last successful point (not from beginning)
- No data loss on crash (WAL journal mode)

### NFR-3: Data Cost
- Full sync: ~50KB (reference) + ~20KB per 10 visits
- Delta sync: typically <10KB
- Photos uploaded separately (not in sync payload)

## API Contract

### POST /api/mobile/sync
```json
// Request
{
  "checkpoints": [{ "visit_id": "...", "type": "arrival", "timestamp": "...", "lat": ..., "lng": ..., "accuracy": ... }],
  "photos": [{ "visit_id": "...", "timestamp": "...", "local_id": "..." }],
  "usages": [{ "visit_id": "...", "material_id": "...", "quantity": ..., "notes": "..." }],
  "measurements": [{ "visit_id": "...", "key": "...", "value": "...", "unit": "..." }],
  "reports": [{ "visit_id": "...", "report_template_id": "...", "payload": {...} }],
  "status_changes": [{ "visit_id": "...", "status": "...", "timestamp": "...", "notes": "..." }]
}

// Response
{
  "synced_at": "2026-01-15T10:30:00Z",
  "results": {
    "checkpoints": [{ "local_id": "...", "server_id": "...", "status": "ok" }],
    "photos": [{ "local_id": "...", "server_id": "...", "status": "ok" }],
    "usages": [{ "local_id": "...", "server_id": "...", "status": "ok" }],
    "measurements": [{ "local_id": "...", "server_id": "...", "status": "ok" }],
    "reports": [{ "local_id": "...", "server_id": "...", "status": "ok" }],
    "status_changes": [{ "local_id": "...", "server_id": "...", "status": "ok", "conflict": false }]
  }
}
```

### GET /api/mobile/sync?since=<timestamp>
```json
// Response
{
  "reference_data": { /* materials, tools, EPP, templates, partners — only on full sync */ },
  "visits": [{ "id": "...", "property": {...}, "partner": {...}, "status": "...", ... }],
  "checklist_items": [{ "visit_id": "...", "stage": "...", "confirmed": true, ... }],
  "photos": [{ "visit_id": "...", "url": "...", "thumbnail_url": "..." }],
  "synced_at": "2026-01-15T10:30:00Z"
}
```

### POST /api/mobile/sync/retry
```json
// Request
{ "item_ids": ["..."] }

// Response
{ "retried": 3, "failed": 0 }
```
