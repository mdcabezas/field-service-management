# Mobile Sync Specification

## Purpose

The mobile sync engine enables offline-first bidirectional synchronization between the mobile app and the backend. Technicians work in the field with unreliable connectivity, so all data is stored locally in SQLite and synced when connectivity returns. The sync protocol supports full sync (reference data + visits), delta sync (changes only), outbox pattern (ordered push), and conflict detection.

## Requirements

### Requirement: Full Sync on Login
The system SHALL perform a full sync when the technician logs in, fetching reference data (materials, tools, EPP, templates, partners) in a single response (~50KB), then assigned visits paginated (10 at a time).

#### Scenario: Technician logs in and syncs
- **WHEN** the technician authenticates successfully
- **THEN** the mobile app fetches reference data from `GET /api/mobile/sync`
- **AND** reference data includes materials, tools, EPP, templates, and partners
- **AND** assigned visits are fetched paginated with basic info, property GPS, partner info, and checklist template IDs

### Requirement: Delta Sync on Reconnect
The system SHALL perform delta sync when connectivity returns, pushing outbox items first, then pulling changes since last sync.

#### Scenario: Technician reconnects after offline period
- **WHEN** the mobile app detects connectivity
- **AND** there are pending outbox items
- **THEN** the app pushes outbox items via `POST /api/mobile/sync`
- **AND** then pulls changes via `GET /api/mobile/sync?since=<timestamp>`
- **AND** delta sync payload is small (<10KB typical)

### Requirement: Outbox Pattern
All local changes SHALL be stored in SQLite `outbox` table with dependency ordering. Push order: checkpoints → photos → usages → measurements → reports → status changes.

#### Scenario: Local change is queued
- **WHEN** the technician makes a local change (e.g., confirms checklist item)
- **THEN** the change is stored in the outbox table with entity_type, entity_id, operation, payload, status, and created_at
- **AND** failed items are retried on next sync attempt

### Requirement: Conflict Detection
The system SHALL detect conflicts during sync. Status conflicts are resolved with server wins (conservative approach). Data conflicts are resolved with client wins (technician's local edits take precedence).

#### Scenario: Status conflict during sync
- **WHEN** the mobile app pushes a status change
- **AND** the server has a newer status
- **THEN** the server wins (conservative approach for safety-critical transitions)
- **AND** the conflict is logged for audit trail

### Requirement: Manual Retry
Failed sync items SHALL be shown in a retry UI. The user can manually trigger retry for individual items or all failures.

#### Scenario: Technician retries failed sync items
- **WHEN** there are failed sync items
- **AND** the technician opens the retry UI
- **THEN** the technician can select individual items or all failures
- **AND** retry clears error state and requeues the item

### Requirement: Sync Status Indicators
The system SHALL display sync status with color coding: green (fully synced), yellow (sync in progress), red (sync failed with error count), grey (offline).

#### Scenario: Technician sees sync status
- **WHEN** the technician views the sync status indicator
- **THEN** the indicator shows the current sync state
- **AND** green means fully synced
- **AND** yellow means sync in progress
- **AND** red means sync failed with error count
- **AND** grey means offline

### Requirement: Technician Identification
The backend SHALL identify the technician via `user_id` from the JWT `sub` claim. The query MUST use `WHERE user_id = $1`, NOT `employee_number`.

#### Scenario: Backend identifies technician during sync
- **WHEN** the mobile app calls `GET /api/mobile/sync`
- **AND** the JWT contains `sub: "10000000-0000-0000-0000-000000000002"`
- **THEN** the backend queries `SELECT id FROM customers.technicians WHERE user_id = $1`
- **AND** returns the technician's assigned visits

#### Scenario: No technician found for user
- **WHEN** the backend cannot find a technician for the given user_id
- **THEN** the backend returns 404 with message "technician not found for employee_number: <id>"
