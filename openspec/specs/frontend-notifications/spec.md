## Purpose

CRUD para gestionar plantillas de notificaciones y notificaciones enviadas.

## ADDED Requirements

### Requirement: Notification Templates list view
The system SHALL display a paginated list of notification templates.

#### Scenario: View templates
- **WHEN** user navigates to /dashboard/notifications/templates
- **THEN** system displays table with columns: Name, Type, Channel, Status
- **AND** supports search by name

### Requirement: Notification Templates CRUD
The system SHALL allow full CRUD for notification templates.

#### Scenario: Create template
- **WHEN** user clicks "Crear Plantilla"
- **THEN** system displays form with fields: name, type, channel, subject, body
- **AND** on submit sends POST to /api/notification-templates

#### Scenario: Edit template
- **WHEN** user clicks edit on a template row
- **THEN** system displays pre-filled form
- **AND** on submit sends PUT to /api/notification-templates/:id

#### Scenario: Delete template
- **WHEN** user clicks delete on a template row
- **THEN** system displays confirmation dialog
- **AND** on confirm sends DELETE to /api/notification-templates/:id

### Requirement: Notifications list view
The system SHALL display a paginated list of sent notifications.

#### Scenario: View notifications
- **WHEN** user navigates to /dashboard/notifications
- **THEN** system displays table with columns: Recipient, Template, Channel, Status, Sent At
- **AND** supports filter by status and date range

### Requirement: Notifications CRUD
The system SHALL allow CRUD for notifications.

#### Scenario: Create notification
- **WHEN** user clicks "Crear Notificación"
- **THEN** system displays form with fields: recipient_id, template_id, channel, data
- **AND** on submit sends POST to /api/notifications

#### Scenario: Edit notification
- **WHEN** user clicks edit on a notification row
- **THEN** system displays pre-filled form
- **AND** on submit sends PUT to /api/notifications/:id

#### Scenario: Delete notification
- **WHEN** user clicks delete on a notification row
- **THEN** system displays confirmation dialog
- **AND** on confirm sends DELETE to /api/notifications/:id
