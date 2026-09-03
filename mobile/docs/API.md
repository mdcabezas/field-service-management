# Mobile API Endpoints

Base URL: `http://localhost:8081/api/mobile`

## Authentication

All endpoints require Bearer token in Authorization header:
```
Authorization: Bearer <access_token>
```

Tokens obtained via `/auth/login` and refreshed via `/auth/refresh`.

---

## Config

### GET /config
Get mobile app configuration.

**Response:** 200 OK
```json
{
  "photo_retention_days": 7,
  "photo_max_storage_mb": 100,
  "photo_warning_threshold_mb": 50
}
```

### PUT /config
Update mobile app configuration.

**Request:**
```json
{
  "photo_retention_days": 14,
  "photo_max_storage_mb": 200,
  "photo_warning_threshold_mb": 100
}
```

**Response:** 200 OK (returns updated config)

---

## Sync

### POST /sync/push
Push local changes to server.

**Request:**
```json
{
  "checkpoints": [
    {
      "visit_id": "uuid",
      "type": "arrival",
      "timestamp": "2026-01-01T10:00:00Z",
      "lat": 19.4326,
      "lng": -99.1332,
      "accuracy": 5.0,
      "local_id": "local-uuid-1"
    }
  ],
  "photos": [
    {
      "visit_id": "uuid",
      "local_id": "local-uuid-2",
      "timestamp": "2026-01-01T10:05:00Z",
      "lat": 19.4326,
      "lng": -99.1332
    }
  ],
  "usages": [
    {
      "visit_id": "uuid",
      "material_id": "uuid",
      "quantity": 5,
      "notes": "Used for repair",
      "local_id": "local-uuid-3"
    }
  ],
  "measurements": [
    {
      "visit_id": "uuid",
      "key": "voltage",
      "value": "220",
      "unit": "V",
      "local_id": "local-uuid-4"
    }
  ],
  "reports": [
    {
      "visit_id": "uuid",
      "report_template_id": "uuid",
      "payload": {"field1": "value1"},
      "local_id": "local-uuid-5"
    }
  ],
  "status_changes": [
    {
      "visit_id": "uuid",
      "status": "in_progress",
      "timestamp": "2026-01-01T10:00:00Z",
      "notes": "Started work",
      "local_id": "local-uuid-6"
    }
  ]
}
```

**Response:** 200 OK
```json
{
  "synced_at": "2026-01-01T10:00:00Z",
  "results": {
    "checkpoints": [{"local_id": "local-uuid-1", "server_id": "server-uuid-1", "status": "ok"}],
    "photos": [{"local_id": "local-uuid-2", "server_id": "server-uuid-2", "status": "ok"}],
    "usages": [{"local_id": "local-uuid-3", "server_id": "server-uuid-3", "status": "ok"}],
    "measurements": [{"local_id": "local-uuid-4", "server_id": "server-uuid-4", "status": "ok"}],
    "reports": [{"local_id": "local-uuid-5", "server_id": "server-uuid-5", "status": "ok"}],
    "status_changes": [{"local_id": "local-uuid-6", "server_id": "server-uuid-6", "status": "ok"}]
  }
}
```

### GET /sync/pull
Pull server changes since last sync.

**Query Parameters:**
- `since` (optional): ISO timestamp for delta sync

**Response:** 200 OK
```json
{
  "reference_data": {
    "materials": [{"id": "uuid", "name": "Cable", "sku": "CAB-001", "unit": "m", "unit_cost": 1.5}],
    "tools": [{"id": "uuid", "name": "Multimeter", "sku": "MUL-001", "description": "Digital multimeter"}],
    "epp": [{"id": "uuid", "name": "Helmet", "sku": "HEL-001", "description": "Safety helmet"}],
    "templates": [{"id": "uuid", "name": "Electrical", "visit_type": "maintenance", "work_type": "electrical", "partner_id": "uuid", "stages_enabled": true}],
    "partners": [{"id": "uuid", "name": "Partner Corp"}]
  },
  "visits": [
    {
      "id": "uuid",
      "property_id": "uuid",
      "property_name": "Building A",
      "property_address": "123 Main St",
      "property_lat": 19.4326,
      "property_lng": -99.1332,
      "partner_id": "uuid",
      "partner_name": "Partner Corp",
      "partner_order_id": "ORD-001",
      "partner_supervisor": "John Doe",
      "type": "maintenance",
      "status": "assigned",
      "priority": "normal",
      "scheduled_at": "2026-01-01T10:00:00Z",
      "checklist_template_id": "uuid",
      "stages_enabled": true,
      "selected_stages": ["electrical", "plumbing"],
      "sla_response_hours": 4,
      "sla_resolution_hours": 24
    }
  ],
  "checklist_items": [
    {"id": "uuid", "visit_id": "uuid", "stage_id": "stage-1", "stage_name": "Electrical", "type": "material", "ref_id": "uuid", "name": "Cable", "quantity": 10, "unit": "m", "confirmed": false, "required": true}
  ],
  "photos": [
    {"id": "uuid", "visit_id": "uuid", "url": "/api/mobile/photos/uuid/file", "thumbnail_url": "/api/mobile/photos/uuid/thumb"}
  ],
  "synced_at": "2026-01-01T10:00:00Z"
}
```

### POST /sync/retry
Retry failed sync items.

**Request:**
```json
{
  "item_ids": ["local-uuid-1", "local-uuid-2"]
}
```

**Response:** 200 OK
```json
{
  "retried": 2,
  "failed": 0
}
```

---

## Photos

### POST /photos
Upload photo (multipart/form-data).

**Form Fields:**
- `visit_id` (string, required)
- `photo` (file, required) - JPEG/PNG, max 5MB
- `timestamp` (string, optional) - ISO timestamp
- `lat` (string, optional)
- `lng` (string, optional)

**Response:** 201 Created
```json
{
  "id": "uuid",
  "visit_id": "uuid",
  "url": "/api/mobile/photos/uuid/file",
  "thumbnail_url": "/api/mobile/photos/uuid/thumb",
  "original_path": "/path/to/photo.jpg",
  "thumbnail_path": "/path/to/thumb.jpg",
  "lat": 19.4326,
  "lng": -99.1332,
  "timestamp": "2026-01-01T10:00:00Z",
  "created_at": "2026-01-01T10:00:00Z"
}
```

### GET /photos/:id/file
Serve original photo file.

### GET /photos/:id/thumb
Serve thumbnail photo file.

---

## GPS

### POST /gps/waivers
Create GPS waiver.

**Request:**
```json
{
  "visit_id": "uuid",
  "type": "gps_timeout",
  "timestamp": "2026-01-01T10:00:00Z"
}
```

**Valid types:** `gps_timeout`, `no_fix`, `low_accuracy`

**Response:** 201 Created
```json
{
  "id": "uuid",
  "visit_id": "uuid",
  "type": "gps_timeout",
  "accepted_by": "Technician Name",
  "timestamp": "2026-01-01T10:00:00Z",
  "created_at": "2026-01-01T10:00:00Z"
}
```

### GET /gps/waivers
List waivers for a visit.

**Query Parameters:**
- `visit_id` (required)

**Response:** 200 OK
```json
{
  "waivers": [
    {"id": "uuid", "visit_id": "uuid", "type": "gps_timeout", "accepted_by": "Tech", "timestamp": "2026-01-01T10:00:00Z", "created_at": "2026-01-01T10:00:00Z"}
  ]
}
```

### DELETE /gps/waivers/:id
Delete a waiver.

**Response:** 204 No Content

---

## Partner

### GET /partners/:id/contacts
Get partner contacts.

**Response:** 200 OK
```json
{
  "contacts": [
    {"id": "uuid", "partner_id": "uuid", "name": "John", "role": "Supervisor", "phone": "+52-55-1234-5678", "email": "john@partner.com"}
  ]
}
```

### GET /properties/:id/partners
Get partners for a property.

**Response:** 200 OK
```json
{
  "partners": [
    {"id": "uuid", "name": "Partner Corp", "order_id": "ORD-001", "supervisor": "John Doe"}
  ]
}
```

### GET /templates
List checklist templates filtered by partner.

**Query Parameters:**
- `partner_id` (optional)

**Response:** 200 OK
```json
{
  "templates": [
    {"id": "uuid", "name": "Electrical", "visit_type": "maintenance", "work_type": "electrical", "partner_id": "uuid", "stages_enabled": true}
  ]
}
```

---

## Visits

### GET /visits/:id/audit
Get audit trail for a visit.

**Query Parameters:**
- `page` (optional, default: 1)
- `limit` (optional, default: 20)

**Response:** 200 OK
```json
{
  "entries": [
    {"id": "uuid", "timestamp": "2026-01-01T10:00:00Z", "user_name": "Tech", "action": "update", "field_name": "status", "old_value": "assigned", "new_value": "en_route"}
  ],
  "total": 1,
  "page": 1,
  "limit": 20
}
```

---

## Properties

### GET /properties/:id/history
Get visit history for a property.

**Response:** 200 OK
```json
{
  "visits": [
    {"id": "uuid", "type": "maintenance", "status": "completed", "technician": "John", "completed_at": "2026-01-01T12:00:00Z"}
  ],
  "total": 1
}
```

---

## Status Codes

| Code | Description |
|------|-------------|
| 200 | OK |
| 201 | Created |
| 204 | No Content |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## Error Response Format

```json
{
  "error": "Error message"
}
```