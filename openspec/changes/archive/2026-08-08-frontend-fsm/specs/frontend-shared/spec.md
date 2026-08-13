## Purpose

CRUD para gestionar plantillas de reportes compartidos.

## ADDED Requirements

### Requirement: Report Templates list view
The system SHALL display a paginated list of report templates.

#### Scenario: View templates
- **WHEN** user navigates to /dashboard/shared/report-templates
- **THEN** system displays table with columns: Name, Description, Fields Count, Status
- **AND** supports search by name

### Requirement: Report Templates CRUD
The system SHALL allow full CRUD for report templates.

#### Scenario: Create template
- **WHEN** user clicks "Crear Plantilla"
- **THEN** system displays form with fields: name, description, fields (JSON structure)
- **AND** fields editor allows adding/removing field definitions
- **AND** on submit sends POST to /api/report-templates

#### Scenario: Edit template
- **WHEN** user clicks edit on a template row
- **THEN** system displays pre-filled form with existing fields
- **AND** on submit sends PUT to /api/report-templates/:id

#### Scenario: Delete template
- **WHEN** user clicks delete on a template row
- **THEN** system displays confirmation dialog
- **AND** on confirm sends DELETE to /api/report-templates/:id

#### Scenario: View template detail
- **WHEN** user clicks a template row
- **THEN** system displays full template info
- **AND** displays field definitions in readable format
- **AND** shows usage count (how many reports use this template)
