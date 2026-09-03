# Sync Protocol

Offline-first synchronization between mobile app and backend server.

---

## Overview

The sync system uses an **outbox pattern** with ordered processing and supports both full and delta synchronization.

---

## Outbox Pattern

### Structure
```typescript
interface OutboxItem {
  id: string;
  entity_type: 'checkpoint' | 'photo' | 'usage' | 'measurement' | 'report' | 'status';
  entity_id: string;
  operation: 'create' | 'update' | 'delete';
  payload: string; // JSON serialized
  status: 'pending' | 'syncing' | 'synced' | 'failed';
  error?: string;
  retry_count: number;
  created_at: string;
}
```

### Dependency Order
Items processed in this order:
1. **checkpoint** - GPS arrival/departure
2. **photo** - Photo uploads
3. **usage** - Material usage
4. **measurement** - Key/value measurements
5. **report** - Completed reports
6. **status** - Visit status transitions

### Processing Logic
```typescript
function processOutbox() {
  // 1. Get pending items ordered by dependency + created_at
  const items = orderOutboxItems(getPendingItems());
  
  // 2. Process each item
  for (const item of items) {
    const result = await pushItem(item);
    if (result.status === 'ok') {
      markSynced(item.id, result.server_id);
    } else if (result.status === 'conflict') {
      handleConflict(item, result);
    } else {
      incrementRetry(item.id);
    }
  }
}
```

---

## Full Sync (On Login)

Triggered when:
- First app launch
- After logout/login
- Manual "Full Sync" button

### Flow
```
Mobile                          Server
  |                               |
  |--- GET /sync/pull?since=null -->|
  |<-- 200 + ReferenceData + ------|
  |      Visits + Items            |
  |                               |
  |--- Store in SQLite             |
  |                               |
```

### Pull Response
```json
{
  "reference_data": {
    "materials": [...],
    "tools": [...],
    "epp": [...],
    "templates": [...],
    "partners": [...]
  },
  "visits": [...],
  "checklist_items": [...],
  "photos": [...],
  "synced_at": "2026-01-01T10:00:00Z"
}
```

### Local Storage
- Reference data → `materials`, `tools`, `epp`, `checklist_templates`, `partners`
- Visits → `visits` table
- Checklist items → `checklist_items` table
- Photos → `photos` table (metadata only, files already on device)

---

## Delta Sync (On Reconnect)

Triggered when:
- App comes online after being offline
- Periodic background sync (configurable)
- Manual sync button

### Flow
```
Mobile                          Server
  |                               |
  |--- GET /sync/pull?since=ts -->|
  |<-- 200 + Changes only -------|
  |                               |
  |--- Apply changes to SQLite ---|
  |                               |
```

### Since Parameter
- `last_sync` from `sync_meta` table
- Only returns changes after that timestamp

### Change Types
- New visits assigned
- Visit status updates
- Checklist template changes
- Reference data updates
- Partner changes

---

## Push Sync (Continuous)

Triggered when:
- Any local change created
- After delta sync completes
- Manual retry

### Flow
```
Mobile                          Server
  |                               |
  |--- POST /sync/push ---------->|
  |      { items... }            |
  |<-- 200 + Results ------------|
  |      { local_id, server_id,  |
  |        status, conflict? }   |
  |                               |
  |--- Update local IDs ---------|
  |--- Mark synced --------------|
  |                               |
```

### Request Format
```json
{
  "checkpoints": [...],
  "photos": [...],
  "usages": [...],
  "measurements": [...],
  "reports": [...],
  "status_changes": [...]
}
```

### Response Format
```json
{
  "synced_at": "2026-01-01T10:00:00Z",
  "results": {
    "checkpoints": [
      {"local_id": "local-1", "server_id": "server-1", "status": "ok"}
    ],
    "photos": [
      {"local_id": "local-2", "server_id": "server-2", "status": "ok"}
    ],
    ...
  }
}
```

### Conflict Resolution
| Entity | Strategy |
|--------|----------|
| Status | Server wins |
| Data (photos, measurements, etc.) | Client wins |
| Checklist items | Merge (last write wins) |

---

## Retry Logic

### Failed Items
- Stored with `status: 'failed'` and `error` message
- `retry_count` incremented
- Max 3 retries before manual intervention

### Manual Retry
```
Mobile                          Server
  |                               |
  |--- POST /sync/retry --------->|
  |      { item_ids: [...] }     |
  |<-- 200 { retried: 2 } -------|
  |                               |
```

### Exponential Backoff
```
Attempt 1: immediate
Attempt 2: 30 seconds
Attempt 3: 2 minutes
```

---

## Photo Upload Flow

Separate from JSON sync (multipart upload):

```
Mobile                          Server
  |                               |
  |--- POST /photos (multipart) -->|
  |      visit_id, photo, ...    |
  |<-- 201 { id, url, ... } -----|
  |                               |
  |--- Add to outbox (photo) --->|
  |      (with server_id)        |
  |                               |
```

### Compression
- Original: max 1280x720
- Compressed: 800px width, 70% JPEG (~200KB)
- Thumbnail: 200px width, 60% JPEG (~20KB)

---

## GPS Validation

### Haversine Distance
```typescript
function haversineDistance(lat1, lng1, lat2, lng2): number {
  // Returns distance in meters
}
```

### Checkpoint Validation
- Required: manual capture (no auto-tracking)
- Accuracy: `Balanced` (~10-50m)
- Timeout: 15 seconds
- Threshold: 500m from property (configurable)

### Waiver Flow
```
Mobile                          Server
  |                               |
  |--- POST /gps/waivers -------->|
  |      { visit_id, type, ... } |
  |<-- 201 { waiver } -----------|
  |                               |
```

---

## Status Transitions

### Valid Transitions
```
assigned    -> en_route, cancelled
en_route    -> in_progress, cancelled
in_progress -> completed, cancelled
completed   -> assigned (reopen, 24h window)
cancelled   -> assigned
rescheduled -> assigned
```

### Reopen
- Only within 24h of completion
- Requires reason
- Sets status to `in_progress`
- Creates audit entry

---

## Offline Behavior

### Queue All Changes
- All mutations go to outbox immediately
- UI updates optimistically
- Background sync when online

### Storage Limits
- Photos: 100MB max (configurable)
- Retention: 7 days (configurable)
- Cleanup on app start + after sync

---

## Error Handling

### Network Errors
- Retry with exponential backoff
- Queue for next sync attempt
- Show pending count in UI

### Server Errors
- 4xx: Don't retry (invalid data)
- 5xx: Retry with backoff
- Conflict: Resolve per strategy

---

## Sync State

```typescript
interface SyncState {
  lastSync: string | null;      // ISO timestamp
  isOnline: boolean;
  isSyncing: boolean;
  pendingItems: number;         // Outbox count
  error: string | null;
}
```

### UI Indicators
- Green: Synced, no pending
- Yellow: Pending items
- Blue: Syncing
- Red: Error (with retry button)

---

## Configuration

```typescript
interface MobileConfig {
  photo_retention_days: 7;
  photo_max_storage_mb: 100;
  photo_warning_threshold_mb: 50;
}
```

### Config Sync
- Downloaded on full sync
- Stored in `config` table
- Can be overridden locally (PUT /config)

---

## Testing Endpoints

```bash
# Full sync
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8081/api/mobile/sync/pull

# Delta sync
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8081/api/mobile/sync/pull?since=2026-01-01T10:00:00Z"

# Push
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"checkpoints":[],"photos":[],"usages":[],"measurements":[],"reports":[],"status_changes":[]}' \
  http://localhost:8081/api/mobile/sync/push

# Retry
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"item_ids":["local-1","local-2"]}' \
  http://localhost:8081/api/mobile/sync/retry
```