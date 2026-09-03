# Photos Spec

## Delta: New capability — photo capture, compression, upload, retention, sharing

## Context
Technicians need to document their work with photos. We add a complete photo pipeline: capture → compress → local store → upload → thumbnail generation → serve. Photos are the heaviest data in the system and require careful management on low-end devices.

## Functional Requirements

### FR-1: Photo Capture
- Use expo-camera for capture (if available) or expo-image-picker for gallery
- Capture resolution: max 1280x720 (sufficient for documentation)
- Auto-rotate based on device orientation
- Add timestamp watermark on capture

### FR-2: Photo Compression
- Use expo-image-manipulator
- Resize to max 800px width (maintain aspect ratio)
- JPEG quality 70%
- Target output: ~200KB per photo
- Compression happens locally before upload

### FR-3: Photo Storage (Local)
- Store compressed photos in app's document directory
- SQLite metadata: id, visit_id, local_path, thumbnail_path, uploaded, created_at
- Thumbnail: 200px width, 60% quality (~20KB)
- Generated locally after compression

### FR-4: Photo Upload
- Multipart upload to `POST /api/mobile/photos`
- Fields: visit_id, photo (file), timestamp, lat, lng
- Upload queued in outbox (not blocking UI)
- Retry on failure (exponential backoff)
- Progress indicator per photo

### FR-5: Photo Serving
- Backend generates thumbnail on upload (200px, 60% quality)
- `/api/mobile/photos/:id/thumb` — thumbnail (~20KB)
- `/api/mobile/photos/:id/file` — full image (~200KB)
- Original retained for 30 days, then deleted

### FR-6: Photo Retention (Configurable)
- New endpoint: `GET/PUT /api/mobile/config`
- Configurable fields:
  - `photo_retention_days`: days to keep photos locally (default: 7)
  - `photo_max_storage_mb`: max MB for photo storage (default: 100)
  - `photo_warning_threshold_mb`: warning threshold (default: 50)
- Cleanup runs on app startup and after successful sync
- Warning shown when approaching storage limit

### FR-7: Photo Sharing
- Use expo-sharing to share photo file
- Share includes: photo file only (no metadata)
- Opens native share sheet
- Available in visit detail photo gallery

### FR-8: Photo Gallery (Visit Detail)
- Show thumbnails in a grid layout
- Tap to view full image
- Long press to share
- Delete option (with confirmation)
- Photo count badge on visit card

## Non-Functional Requirements

### NFR-1: Storage Efficiency
- Compressed photos: ~200KB each
- Thumbnails: ~20KB each
- 100MB storage cap (configurable)
- ~500 photos at 200KB = 100MB

### NFR-2: Upload Performance
- Upload in background (non-blocking)
- Queue max 50 items
- Exponential backoff on failure
- No duplicate uploads (idempotent)

### NFR-3: Quality
- Minimum 720p resolution for documentation
- JPEG compression artifacts acceptable for field documentation
- Timestamp watermark visible but not obstructive

## API Contract

### POST /api/mobile/photos
```
Content-Type: multipart/form-data

Fields:
- visit_id: UUID
- photo: File (JPEG, max 5MB)
- timestamp: ISO 8601
- lat: float (optional)
- lng: float (optional)
```

```json
// Response
{
  "id": "...",
  "visit_id": "...",
  "url": "/api/mobile/photos/.../file",
  "thumbnail_url": "/api/mobile/photos/.../thumb",
  "created_at": "..."
}
```

### GET /api/mobile/photos/:id/thumb
```
Response: image/jpeg
Cache-Control: max-age=86400
```

### GET /api/mobile/photos/:id/file
```
Response: image/jpeg
Cache-Control: max-age=3600
```

### GET /api/mobile/config
```json
// Response
{
  "photo_retention_days": 7,
  "photo_max_storage_mb": 100,
  "photo_warning_threshold_mb": 50
}
```

### PUT /api/mobile/config
```json
// Request
{ "photo_retention_days": 14 }

// Response
{ "photo_retention_days": 14, "photo_max_storage_mb": 100, "photo_warning_threshold_mb": 50 }
```
