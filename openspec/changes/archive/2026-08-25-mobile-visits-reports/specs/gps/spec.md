# GPS Spec

## Delta: New capability — manual GPS capture, distance validation, waiver system

## Context
Technicians need to prove they were on-site. We add manual GPS capture (not automatic, to save battery) with distance validation against property coordinates. A waiver system handles cases where GPS fix is unavailable.

## Functional Requirements

### FR-1: Manual GPS Capture
- GPS captured manually (not continuous tracking)
- Capture points: arrival, departure, photos, measurements
- Use `expo-location` with `Accuracy.Balanced` (not high accuracy, saves battery)
- 15-second timeout on GPS request
- If timeout: show error, allow retry or skip

### FR-2: Distance Validation
- Property GPS resolved via 3 JOINs: visits → properties → customer_addresses → geocoding.addresses.geom
- Haversine formula for distance calculation
- Configurable threshold (default: 500m)
- If distance > threshold: show warning "You appear to be X meters from the property"
- Warning does NOT block operation (technician can proceed)

### FR-3: GPS Waiver
- When GPS fails (timeout/no fix), technician can accept a waiver
- Waiver stored with: visit_id, type, timestamp, accepted_by
- Waiver types: "gps_timeout", "no_fix", "low_accuracy"
- Waiver shown in audit trail
- Waiver acceptance is timestamped

### FR-4: GPS on Checkpoints
- Arrival checkpoint: GPS capture + distance validation
- Departure checkpoint: GPS capture + distance validation
- Both stored with: lat, lng, accuracy, timestamp
- Distance from property calculated and stored

### FR-5: GPS on Photos
- Optional GPS capture when taking photo
- If GPS available: attach to photo metadata
- If GPS unavailable: photo stored without GPS (not blocked)

### FR-6: GPS on Measurements
- Optional GPS capture when recording measurement
- If GPS available: attach to measurement metadata
- If GPS unavailable: measurement stored without GPS (not blocked)

### FR-7: GPS Status Indicator
- Green dot: GPS acquired (accuracy < 100m)
- Yellow dot: GPS acquired (accuracy 100-500m)
- Red dot: GPS acquired (accuracy > 500m)
- Grey dot: GPS not acquired

## Non-Functional Requirements

### NFR-1: Battery Optimization
- Manual capture only (no continuous tracking)
- `Accuracy.Balanced` (not high accuracy)
- 15-second timeout (no infinite waits)
- GPS off when app in background

### NFR-2: Data Cost
- GPS coordinates: ~40 bytes per point
- Negligible data cost

### NFR-3: Accuracy
- Haversine formula: ~0.3% error at earth surface (sufficient for 500m threshold)
- Device GPS accuracy: typically 10-30m outdoors

## API Contract

### POST /api/mobile/gps/waiver
```json
// Request
{
  "visit_id": "...",
  "type": "gps_timeout",
  "timestamp": "2026-01-15T14:30:00Z"
}

// Response
{
  "id": "...",
  "visit_id": "...",
  "type": "gps_timeout",
  "accepted_by": "Juan Pérez",
  "timestamp": "2026-01-15T14:30:00Z"
}
```

### GET /api/mobile/gps/waivers?visit_id=X
```json
// Response
{
  "waivers": [{
    "id": "...",
    "type": "gps_timeout",
    "accepted_by": "Juan Pérez",
    "timestamp": "2026-01-15T14:30:00Z"
  }]
}
```

### DELETE /api/mobile/gps/waivers/:id
```json
// Response 204 No Content
```
