# Mobile GPS Specification

## Purpose

The mobile GPS system provides manual GPS capture with distance validation and a waiver system for cases where GPS fix is unavailable. Technicians prove they were on-site by capturing GPS at arrival and departure checkpoints, with distance validation against property coordinates.

## Requirements

### Requirement: Manual GPS Capture
GPS SHALL be captured manually (not continuous tracking) at capture points: arrival, departure, photos, measurements. Use `expo-location` with `Accuracy.Balanced` and 15-second timeout.

#### Scenario: Technician captures GPS at arrival
- **WHEN** the technician taps "Arrival" checkpoint
- **THEN** the app requests GPS with `Accuracy.Balanced`
- **AND** GPS coordinates are captured within 15 seconds
- **AND** coordinates are stored with lat, lng, accuracy, and timestamp

#### Scenario: GPS request times out
- **WHEN** GPS request times out after 15 seconds
- **THEN** the app shows an error message
- **AND** allows retry or skip

### Requirement: Distance Validation
The system SHALL validate GPS distance against property coordinates using Haversine formula. Configurable threshold (default: 500m). Warning does NOT block operation.

#### Scenario: Technician is near property
- **WHEN** technician captures GPS
- **AND** distance to property is < 500m
- **THEN** no warning is shown

#### Scenario: Technician is far from property
- **WHEN** technician captures GPS
- **AND** distance to property is > 500m
- **THEN** a warning is shown: "You appear to be X meters from the property"
- **AND** the technician can proceed (warning does not block)

### Requirement: GPS Waiver
When GPS fails (timeout/no fix), the technician CAN accept a waiver. Waiver stored with visit_id, type, timestamp, accepted_by. Waiver types: "gps_timeout", "no_fix", "low_accuracy".

#### Scenario: Technician accepts GPS waiver
- **WHEN** GPS fails (timeout or no fix)
- **AND** technician taps "Accept Waiver"
- **THEN** waiver is stored via `POST /api/mobile/gps/waiver`
- **AND** waiver includes visit_id, type, timestamp, and accepted_by
- **AND** waiver is shown in audit trail

### Requirement: GPS on Checkpoints
Arrival and departure checkpoints SHALL capture GPS with distance validation. Both stored with lat, lng, accuracy, timestamp, and distance from property.

#### Scenario: Arrival checkpoint with GPS
- **WHEN** technician taps "Arrival"
- **AND** GPS is available
- **THEN** GPS coordinates are captured
- **AND** distance from property is calculated and stored

### Requirement: GPS on Photos
Optional GPS capture when taking photo. If GPS available, attach to photo metadata. If GPS unavailable, photo stored without GPS (not blocked).

#### Scenario: Photo with GPS
- **WHEN** technician takes a photo
- **AND** GPS is available
- **THEN** GPS coordinates are attached to photo metadata

### Requirement: GPS on Measurements
Optional GPS capture when recording measurement. If GPS available, attach to measurement metadata. If GPS unavailable, measurement stored without GPS (not blocked).

#### Scenario: Measurement with GPS
- **WHEN** technician records a measurement
- **AND** GPS is available
- **THEN** GPS coordinates are attached to measurement metadata

### Requirement: GPS Status Indicator
The system SHALL display GPS status with color coding: green (accuracy < 100m), yellow (accuracy 100-500m), red (accuracy > 500m), grey (not acquired).

#### Scenario: GPS status indicator
- **WHEN** technician views GPS status
- **THEN** indicator shows current GPS accuracy
- **AND** green means accuracy < 100m
- **AND** yellow means accuracy 100-500m
- **AND** red means accuracy > 500m
- **AND** grey means GPS not acquired
