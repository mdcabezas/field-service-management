# Mobile Visits Specification

## Purpose

Visits are the core entity in the mobile app. The system extends the existing visit model with partner-specific data, property history, reopen capability, and audit trail. Checklists are validated before status transitions.

## Requirements

### Requirement: Visit List (Mobile)
Show technician's assigned visits for today. Display: property name, partner name, scheduled time, status, priority. Pull-to-refresh triggers delta sync. Paginated: load 10 visits at a time. Filter by status.

#### Scenario: Technician views visit list
- **WHEN** technician opens visit list
- **THEN** assigned visits for today are shown
- **AND** each visit shows property name, partner name, scheduled time, status, priority
- **AND** pull-to-refresh triggers delta sync

### Requirement: Visit Detail (Lazy Loaded)
Visit details fetched when user opens visit (not on list load). Includes: property info, partner info, assignments, reports, property history, checklist by stages, audit trail.

#### Scenario: Technician opens visit detail
- **WHEN** technician taps a visit
- **THEN** visit details are fetched
- **AND** includes property info, partner info, assignments, reports
- **AND** includes property history, checklist by stages, audit trail

### Requirement: Visit Creation (with partner resolution)
Planner selects property → partner auto-resolved. Planner selects template → filtered by partner. Planner selects stages → from template's stages. SLA auto-linked if partner + work_type match.

#### Scenario: Visit is created with partner resolution
- **WHEN** planner selects a property
- **THEN** partner is auto-resolved
- **AND** templates are filtered by partner
- **AND** stages are selected from template

### Requirement: Status Transitions
Validated against checklist stages. Required stages must be complete before transitioning to "completada". Emergency visits can skip pre-departure. All transitions audited.

#### Scenario: Status transition is validated
- **WHEN** technician tries to transition status
- **THEN** required stages are checked for completeness
- **AND** if incomplete, transition is blocked with specific error

### Requirement: Property History
`GET /api/properties/:id/history` returns previous visits. Shows: date, type, status, technician name. Limited to last 20 visits.

#### Scenario: Property history is shown
- **WHEN** technician opens visit detail
- **THEN** property history is shown
- **AND** history includes date, type, status, technician name

### Requirement: Reopen Visit
`POST /api/visits/:id/reopen` with reason. Within configurable window (default 24h from completion). Requires reason text. Full audit trail preserved.

#### Scenario: Technician reopens visit
- **WHEN** technician taps "Reopen" on completed visit
- **AND** within 24h window
- **AND** provides reason
- **THEN** visit status reverts to "en_curso"
- **AND** audit trail is preserved

### Requirement: Audit Trail
`GET /api/visits/:id/audit` returns all changes. Shows: timestamp, user, action, old_value, new_value. Paginated (50 per page).

#### Scenario: Audit trail is shown
- **WHEN** technician opens audit trail
- **THEN** all changes are shown
- **AND** each entry shows timestamp, user, action, old_value, new_value

### Requirement: Partner Data
`partner_order_id` (external order reference) and `partner_supervisor` (partner's supervisor name). Fields visible/editable only when partner is assigned.

#### Scenario: Partner data is shown
- **WHEN** visit has partner assigned
- **THEN** partner_order_id and partner_supervisor are shown
- **AND** fields are editable
