# Mobile Visits & Reports — Technical Design

## Architecture Overview

```
┌─────────────────────────────────────────────────────┐
│                  MOBILE APP                         │
│  Expo SDK 57 · React Native 0.79                   │
├─────────────────────────────────────────────────────┤
│  ┌─────────┐  ┌──────────┐  ┌──────────────────┐  │
│  │ Zustand │  │ Expo     │  │ op-sqlite        │  │
│  │ Stores  │  │ Router   │  │ (local DB)       │  │
│  └────┬────┘  └────┬─────┘  └────────┬─────────┘  │
│       │            │                  │             │
│  ┌────┴────────────┴──────────────────┴─────────┐  │
│  │              lib/database.ts                  │  │
│  │  SQLite + PRAGMAs + migrations + queries      │  │
│  └──────────────────┬───────────────────────────┘  │
│                     │                              │
│  ┌──────────────────┴───────────────────────────┐  │
│  │              lib/sync-engine.ts               │  │
│  │  Full sync · Delta sync · Outbox · Retry      │  │
│  └──────────────────┬───────────────────────────┘  │
│                     │                              │
│  ┌──────────────────┴───────────────────────────┐  │
│  │              lib/auth.ts                      │  │
│  │  Login · Token refresh · SecureStore          │  │
│  └──────────────────┬───────────────────────────┘  │
│                     │                              │
│  ┌──────────────────┴───────────────────────────┐  │
│  │              lib/photo-upload.ts              │  │
│  │  Compress · Queue · Upload · Retry            │  │
│  └──────────────────────────────────────────────┘  │
└─────────────────────┬───────────────────────────────┘
                      │ HTTPS
┌─────────────────────┴───────────────────────────────┐
│                  BACKEND API                        │
│  Go · Gin · PostgreSQL                              │
├─────────────────────────────────────────────────────┤
│  /api/mobile/sync     Sync protocol                │
│  /api/mobile/photos   Photo upload/serve            │
│  /api/mobile/config   Configuration                 │
│  /api/mobile/gps/*    GPS waivers                   │
│  /api/visits/*        Visit CRUD + status           │
│  /api/checklist-*     Checklist stages              │
│  /api/properties/*    Property history              │
│  /api/partners/*      Partner contacts              │
│  /api/templates       Template filtering            │
└─────────────────────────────────────────────────────┘
```

## SQLite Schema (Mobile)

```sql
-- PRAGMA optimizations for low-end devices
PRAGMA journal_mode = MEMORY;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 2000;
PRAGMA temp_store = MEMORY;
PRAGMA mmap_size = 67108864;

-- Reference data (full sync)
CREATE TABLE materials (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  sku TEXT,
  unit TEXT,
  unit_cost REAL
);

CREATE TABLE tools (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  sku TEXT,
  description TEXT
);

CREATE TABLE epp (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  sku TEXT,
  description TEXT
);

CREATE TABLE checklist_templates (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  visit_type TEXT,
  work_type TEXT,
  partner_id TEXT,
  stages_enabled INTEGER DEFAULT 0
);

CREATE TABLE checklist_stages (
  id TEXT PRIMARY KEY,
  template_id TEXT NOT NULL,
  name TEXT NOT NULL,
  sort_order INTEGER DEFAULT 0,
  required INTEGER DEFAULT 0
);

CREATE TABLE partners (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE partner_contacts (
  id TEXT PRIMARY KEY,
  partner_id TEXT NOT NULL,
  name TEXT NOT NULL,
  role TEXT,
  phone TEXT,
  email TEXT
);

-- Visits (delta sync)
CREATE TABLE visits (
  id TEXT PRIMARY KEY,
  property_id TEXT,
  property_name TEXT,
  property_address TEXT,
  property_lat REAL,
  property_lng REAL,
  partner_id TEXT,
  partner_name TEXT,
  partner_order_id TEXT,
  partner_supervisor TEXT,
  type TEXT,
  status TEXT DEFAULT 'assigned',
  priority TEXT DEFAULT 'normal',
  scheduled_at TEXT,
  checklist_template_id TEXT,
  stages_enabled INTEGER DEFAULT 0,
  selected_stages TEXT, -- JSON array of stage IDs
  sla_response_hours INTEGER,
  sla_resolution_hours INTEGER,
  created_at TEXT DEFAULT (datetime('now')),
  updated_at TEXT DEFAULT (datetime('now'))
);

-- Checklist items
CREATE TABLE checklist_items (
  id TEXT PRIMARY KEY,
  visit_id TEXT NOT NULL,
  stage_id TEXT,
  stage_name TEXT,
  type TEXT NOT NULL, -- material, tool, epp
  ref_id TEXT NOT NULL, -- material/tool/epp ID
  name TEXT NOT NULL,
  quantity REAL DEFAULT 1,
  unit TEXT,
  confirmed INTEGER DEFAULT 0,
  required INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now'))
);

-- GPS checkpoints
CREATE TABLE checkpoints (
  id TEXT PRIMARY KEY,
  visit_id TEXT NOT NULL,
  type TEXT NOT NULL, -- arrival, departure
  timestamp TEXT NOT NULL,
  lat REAL,
  lng REAL,
  accuracy REAL,
  distance_from_property REAL,
  synced INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now'))
);

-- GPS waivers
CREATE TABLE gps_waivers (
  id TEXT PRIMARY KEY,
  visit_id TEXT NOT NULL,
  type TEXT NOT NULL, -- gps_timeout, no_fix, low_accuracy
  timestamp TEXT NOT NULL,
  accepted_by TEXT,
  synced INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now'))
);

-- Photos (local metadata)
CREATE TABLE photos (
  id TEXT PRIMARY KEY,
  visit_id TEXT NOT NULL,
  local_path TEXT NOT NULL,
  thumbnail_path TEXT,
  server_id TEXT,
  uploaded INTEGER DEFAULT 0,
  lat REAL,
  lng REAL,
  timestamp TEXT,
  created_at TEXT DEFAULT (datetime('now'))
);

-- Measurements
CREATE TABLE measurements (
  id TEXT PRIMARY KEY,
  visit_id TEXT NOT NULL,
  key TEXT NOT NULL,
  value TEXT,
  unit TEXT,
  lat REAL,
  lng REAL,
  synced INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now'))
);

-- Material usage
CREATE TABLE usages (
  id TEXT PRIMARY KEY,
  visit_id TEXT NOT NULL,
  material_id TEXT NOT NULL,
  quantity REAL NOT NULL,
  notes TEXT,
  synced INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now'))
);

-- Audit trail (local cache)
CREATE TABLE audit_entries (
  id TEXT PRIMARY KEY,
  visit_id TEXT NOT NULL,
  timestamp TEXT NOT NULL,
  user_name TEXT,
  action TEXT,
  field_name TEXT,
  old_value TEXT,
  new_value TEXT,
  created_at TEXT DEFAULT (datetime('now'))
);

-- Outbox (sync queue)
CREATE TABLE outbox (
  id TEXT PRIMARY KEY,
  entity_type TEXT NOT NULL, -- checkpoint, photo, usage, measurement, report, status
  entity_id TEXT NOT NULL,
  operation TEXT NOT NULL, -- create, update, delete
  payload TEXT NOT NULL, -- JSON
  status TEXT DEFAULT 'pending', -- pending, sending, sent, failed
  error TEXT,
  retry_count INTEGER DEFAULT 0,
  created_at TEXT DEFAULT (datetime('now'))
);

-- Configuration
CREATE TABLE config (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

-- Sync metadata
CREATE TABLE sync_meta (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);
```

## Sync Protocol

```
FULL SYNC (on login):
═══════════════════════════════════════════════════

Mobile                          Backend
  │                               │
  │── GET /api/mobile/sync ──────>│
  │   (no since param)            │
  │                               │
  │<── { reference_data, visits }─│
  │   (~50KB reference)           │
  │   (~20KB per 10 visits)       │
  │                               │
  │ Store in SQLite               │
  │ Mark: last_sync = synced_at   │

DELTA SYNC (on reconnect):
═══════════════════════════════════════════════════

Mobile                          Backend
  │                               │
  │ 1. Push outbox items          │
  │── POST /api/mobile/sync ─────>│
  │   { checkpoints, photos,      │
  │     usages, measurements,     │
  │     reports, status_changes } │
  │                               │
  │<── { results, synced_at } ────│
  │                               │
  │ 2. Pull changes               │
  │── GET /api/mobile/sync? ─────>│
  │   since=<last_sync>           │
  │                               │
  │<── { visits, checklist_items,─│
  │      photos, audit }          │
  │                               │
  │ Merge into SQLite             │
  │ Mark: last_sync = synced_at   │
```

## Outbox Ordering

```
PUSH ORDER (dependency graph):
═══════════════════════════════════════════════════

1. checkpoints     (no dependencies)
2. photos          (depends on: checkpoints for GPS)
3. usages          (depends on: nothing, but logical order)
4. measurements    (depends on: nothing, but logical order)
5. reports         (depends on: photos, usages, measurements)
6. status_changes  (depends on: everything above)

Each outbox item:
├── id: UUID
├── entity_type: checkpoint|photo|usage|measurement|report|status
├── entity_id: UUID (local)
├── operation: create|update|delete
├── payload: JSON (full entity)
├── status: pending|sending|sent|failed
├── error: error message (if failed)
├── retry_count: integer
└── created_at: ISO timestamp
```

## Photo Pipeline

```
CAPTURE → COMPRESS → STORE → UPLOAD → SERVE
═══════════════════════════════════════════════════

1. CAPTURE (expo-camera or expo-image-picker)
   └── Max resolution: 1280x720

2. COMPRESS (expo-image-manipulator)
   ├── Resize: max 800px width
   ├── Quality: 70%
   └── Output: ~200KB JPEG

3. STORE (local filesystem)
   ├── Full: documents/photos/<uuid>.jpg
   ├── Thumbnail: documents/thumbs/<uuid>.jpg (200px, 60%)
   └── Metadata: SQLite photos table

4. UPLOAD (outbox → POST /api/mobile/photos)
   ├── Multipart form data
   ├── Queue in outbox
   ├── Retry on failure
   └── Progress indicator

5. SERVE (backend)
   ├── Thumbnail: /api/mobile/photos/:id/thumb (20KB)
   ├── Full: /api/mobile/photos/:id/file (200KB)
   └── Cache: 24h for thumb, 1h for file
```

## Photo Retention

```
CLEANUP STRATEGY:
═══════════════════════════════════════════════════

TRIGGERS:
├── App startup
└── After successful sync

PROCESS:
├── 1. Calculate total photo storage size
├── 2. If > warning_threshold_mb: show warning
├── 3. If > max_storage_mb: delete oldest photos
├── 4. Delete photos older than retention_days
├── 5. Clean up orphaned files (no DB record)
└── 6. Vacuum SQLite

CONFIG:
├── photo_retention_days: 7 (default)
├── photo_max_storage_mb: 100 (default)
└── photo_warning_threshold_mb: 50 (default)
```

## Distance Validation

```
HAVERSINE FORMULA:
═══════════════════════════════════════════════════

Distance = 2R × arcsin(√(
  sin²((φ₂-φ₁)/2) + cos(φ₁)cos(φ₂)sin²((λ₂-λ₁)/2)
))

Where:
├── R = 6,371,000m (Earth radius)
├── φ = latitude (radians)
└── λ = longitude (radians)

PROPERTY GPS RESOLUTION:
═══════════════════════════════════════════════════

visits.property_id
  → properties.customer_address_id
    → customer_addresses.address_id
      → geocoding.addresses.geom (PostGIS POINT)

MOBILE:
├── Receives property_lat, property_lng in sync payload
├── Calculates distance on-device
└── No need for PostGIS on mobile

THRESHOLD:
├── Default: 500m
├── Configurable per visit type
└── Warning only (does not block)
```

## Partner Auto-Resolution

```
RESOLUTION CHAIN:
═══════════════════════════════════════════════════

1. Visit creation with property_id
2. Query: customer_address_partners
   WHERE customer_address_id = property's address
   AND is_active = true
3. If found: set partner_id on visit
4. If multiple: return list, let user choose
5. If none: leave null

SLA AUTO-LINK:
═══════════════════════════════════════════════════

1. Visit has partner_id + type (work_type)
2. Query: partners.slas
   WHERE partner_id = X
   AND work_type = Y
   AND active = true
3. If found: create VisitSLATracking
4. If not found: skip
5. If multiple: most specific wins

TEMPLATE FILTERING:
═══════════════════════════════════════════════════

1. Get templates where partner_id = X (partner-specific)
2. Get templates where partner_id IS NULL (global)
3. Merge: partner-specific first
4. If same name in both: partner wins (override)
```

## Checklist Stage Validation

```
VALIDATION RULES:
═══════════════════════════════════════════════════

WHEN: status transition to "completada"
CHECK:
├── For each selected stage where required = true:
│   ├── Get all checklist items in that stage
│   ├── Filter items where required = true
│   ├── Check all are confirmed
│   └── If any unconfirmed: BLOCK transition
├── If all required stages complete: ALLOW transition
└── Emergency visits: all stages become optional

ERROR RESPONSE:
{
  "error": "Checklist incomplete",
  "blocking_stages": [
    { "stage_id": "...", "name": "Post-work", "incomplete_required": 3 }
  ]
}
```

## Testing Structure

```
BACKEND:
═══════════════════════════════════════════════════

backend/internal/
├── service/operations/
│   ├── visit_service_test.go          (expand, +15 tests)
│   ├── report_service_test.go         (expand, +5 tests)
│   ├── checklist_stage_test.go        (NEW, +12 tests)
│   └── sla_tracking_service_test.go   (expand, +3 tests)
├── service/mobile/
│   ├── sync_service_test.go           (NEW, +10 tests)
│   ├── photo_upload_service_test.go   (NEW, +7 tests)
│   └── partner_resolve_test.go        (NEW, +8 tests)
├── handler/mobile/
│   ├── sync_handler_test.go           (NEW, +5 tests)
│   ├── photo_handler_test.go          (NEW, +4 tests)
│   ├── config_handler_test.go         (NEW, +2 tests)
│   ├── audit_handler_test.go          (NEW, +3 tests)
│   └── property_history_handler_test.go (NEW, +3 tests)
└── handler/operations/
    └── checklist_stage_handler_test.go (NEW, +5 tests)

MOBILE:
═══════════════════════════════════════════════════

mobile/__tests__/
├── lib/
│   ├── database.test.ts               (NEW, 10 tests)
│   ├── sync-engine.test.ts            (NEW, 15 tests)
│   ├── auth.test.ts                   (NEW, 5 tests)
│   └── photo-upload.test.ts           (NEW, 8 tests)
├── utils/
│   ├── distance.test.ts               (NEW, 5 tests)
│   ├── ordered-sync.test.ts           (NEW, 7 tests)
│   ├── uuid.test.ts                   (NEW, 3 tests)
│   └── date.test.ts                   (NEW, 3 tests)
├── stores/
│   ├── auth-store.test.ts             (NEW, 2 tests)
│   ├── sync-store.test.ts             (NEW, 2 tests)
│   └── visit-store.test.ts            (NEW, 2 tests)
└── integration/
    ├── sync-flow.test.ts              (NEW, 5 tests)
    ├── auth-flow.test.ts              (NEW, 3 tests)
    └── photo-flow.test.ts             (NEW, 3 tests)

FRONTEND (Playwright):
═══════════════════════════════════════════════════

frontend/tests/
├── checklist-stages.spec.ts           (NEW, 10 scenarios)
├── property-history.spec.ts           (NEW, 3 scenarios)
├── partner-customization.spec.ts      (NEW, 5 scenarios)
├── reopen-visit.spec.ts               (NEW, 4 scenarios)
├── audit-trail.spec.ts                (NEW, 3 scenarios)
└── visit-checklist-full.spec.ts       (NEW, 5 scenarios)
```
