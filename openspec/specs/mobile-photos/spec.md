# Mobile Photos Specification

## Purpose

The mobile photo system provides a complete photo pipeline: capture → compress → local store → upload → thumbnail generation → serve. Photos are the heaviest data in the system and require careful management on low-end devices with configurable retention and storage limits.

## Requirements

### Requirement: Photo Capture
Use expo-camera for capture (if available) or expo-image-picker for gallery. Capture resolution max 1280x720. Auto-rotate based on device orientation. Add timestamp watermark on capture.

#### Scenario: Technician captures photo
- **WHEN** technician taps "Take Photo"
- **THEN** camera opens with max 1280x720 resolution
- **AND** photo is auto-rotated based on device orientation
- **AND** timestamp watermark is applied

### Requirement: Photo Compression
Use expo-image-manipulator to resize to max 800px width (maintain aspect ratio) at JPEG quality 70%. Target output ~200KB per photo. Compression happens locally before upload.

#### Scenario: Photo is compressed
- **WHEN** photo is captured
- **THEN** photo is resized to max 800px width
- **AND** JPEG quality is set to 70%
- **AND** output is ~200KB

### Requirement: Photo Storage (Local)
Store compressed photos in app's document directory. SQLite metadata: id, visit_id, local_path, thumbnail_path, uploaded, created_at. Thumbnail: 200px width, 60% quality (~20KB).

#### Scenario: Photo is stored locally
- **WHEN** photo is compressed
- **THEN** photo is stored in app's document directory
- **AND** metadata is saved in SQLite
- **AND** thumbnail is generated (200px, 60% quality)

### Requirement: Photo Upload
Multipart upload to `POST /api/mobile/photos`. Upload queued in outbox (not blocking UI). Retry on failure with exponential backoff. Progress indicator per photo.

#### Scenario: Photo is uploaded
- **WHEN** photo is queued for upload
- **THEN** photo is uploaded via multipart to `POST /api/mobile/photos`
- **AND** upload is non-blocking (queued in outbox)
- **AND** progress indicator is shown

### Requirement: Photo Serving
Backend generates thumbnail on upload (200px, 60% quality). Endpoints: `/api/mobile/photos/:id/thumb` (~20KB) and `/api/mobile/photos/:id/file` (~200KB). Original retained for 30 days.

#### Scenario: Photo is served
- **WHEN** mobile app requests photo
- **THEN** thumbnail is served from `/api/mobile/photos/:id/thumb`
- **AND** full image is served from `/api/mobile/photos/:id/file`

### Requirement: Photo Retention (Configurable)
New endpoint `GET/PUT /api/mobile/config` with configurable fields: photo_retention_days (default: 7), photo_max_storage_mb (default: 100), photo_warning_threshold_mb (default: 50). Cleanup runs on app startup and after successful sync.

#### Scenario: Photo retention cleanup
- **WHEN** app starts or sync completes
- **THEN** photos older than retention days are deleted
- **AND** warning is shown when approaching storage limit

### Requirement: Photo Sharing
Use expo-sharing to share photo file. Share includes photo file only (no metadata). Opens native share sheet. Available in visit detail photo gallery.

#### Scenario: Technician shares photo
- **WHEN** technician long-presses a photo
- **AND** taps "Share"
- **THEN** native share sheet opens with photo file

### Requirement: Photo Gallery (Visit Detail)
Show thumbnails in a grid layout. Tap to view full image. Long press to share. Delete option with confirmation. Photo count badge on visit card.

#### Scenario: Technician views photo gallery
- **WHEN** technician opens visit detail
- **THEN** photo thumbnails are shown in grid layout
- **AND** tap opens full image
- **AND** long press shows share option
- **AND** delete option requires confirmation
