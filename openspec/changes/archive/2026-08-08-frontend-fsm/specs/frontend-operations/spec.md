## Purpose

CRUD para gestionar visitas con todo el sub-árbol: asignaciones, checklists, mediciones, checkpoints, fotos, reportes, vehículos, arriendos y SLAs.

## ADDED Requirements

### Requirement: Visits list view
The system SHALL display a paginated list of visits.

#### Scenario: View visits
- **WHEN** user navigates to /dashboard/operations/visits
- **THEN** system displays table with columns: Date, Customer, Address, Technician, Status
- **AND** supports filter by date range, technician, status
- **AND** supports search by customer name

### Requirement: Visit detail view
The system SHALL display visit details with nested sub-resources.

#### Scenario: View visit detail
- **WHEN** user clicks a visit row
- **THEN** system displays visit info at top
- **AND** displays tabs: Assignments, Checklist Materials, Checklist Tools, Checklist EPPs, Measurements, Checkpoints, Photos, Reports, Vehicle Assignments, Rentals, SLA Trackings
- **AND** each tab shows respective list with CRUD actions

### Requirement: Visit Assignments management
The system SHALL allow CRUD for technician assignments to visits.

#### Scenario: Create assignment
- **WHEN** user clicks "Crear Asignación" in visit detail
- **THEN** system displays form with fields: technician_id, role, status
- **AND** on submit sends POST to /api/visit-assignments with visit_id

### Requirement: Visit Checklist management
The system SHALL allow CRUD for materials, tools, and EPPs used in visit checklists.

#### Scenario: Add checklist material
- **WHEN** user clicks "Agregar Material" in visit checklist tab
- **THEN** system displays form with fields: material_id, quantity, status
- **AND** on submit sends POST to /api/visit-checklist-materials with visit_id

### Requirement: Visit Measurements management
The system SHALL allow CRUD for measurements taken during visits.

#### Scenario: Create measurement
- **WHEN** user clicks "Agregar Medición" in visit detail
- **THEN** system displays form with fields: type, value, unit, notes
- **AND** on submit sends POST to /api/visit-measurements with visit_id

### Requirement: Visit Checkpoints management
The system SHALL allow CRUD for checkpoints within visits.

#### Scenario: Create checkpoint
- **WHEN** user clicks "Agregar Checkpoint" in visit detail
- **THEN** system displays form with fields: name, status, notes, photo_url
- **AND** on submit sends POST to /api/visit-checkpoints with visit_id

### Requirement: Visit Photos management
The system SHALL allow CRUD for photos attached to visits.

#### Scenario: Add photo
- **WHEN** user clicks "Agregar Foto" in visit detail
- **THEN** system displays form with fields: url, description, finding_type
- **AND** on submit sends POST to /api/visit-photos with visit_id

### Requirement: Visit Reports management
The system SHALL allow CRUD for reports with entries and images.

#### Scenario: Create report
- **WHEN** user clicks "Crear Reporte" in visit detail
- **THEN** system displays form with fields: template_id, title, status
- **AND** on submit sends POST to /api/visit-reports with visit_id

#### Scenario: Add report entry
- **WHEN** user opens a report detail
- **THEN** system displays Entries tab
- **AND** allows adding entries via POST to /api/report-entries

#### Scenario: Add report image
- **WHEN** user opens a report detail
- **THEN** system displays Images tab
- **AND** allows uploading images via POST to /api/report-images

### Requirement: Vehicle Assignments management
The system SHALL allow CRUD for vehicle assignments to visits.

#### Scenario: Create vehicle assignment
- **WHEN** user clicks "Asignar Vehículo" in visit detail
- **THEN** system displays form with fields: vehicle_id, start_date, end_date
- **AND** on submit sends POST to /api/vehicle-assignments

### Requirement: Visit Rentals management
The system SHALL allow CRUD for rentals associated with visits.

#### Scenario: Add rental to visit
- **WHEN** user clicks "Agregar Arriendo" in visit detail
- **THEN** system displays form with fields: rental_id
- **AND** on submit sends POST to /api/visit-rentals with visit_id

### Requirement: SLA Trackings management
The system SHALL allow CRUD for SLA tracking records on visits.

#### Scenario: Create SLA tracking
- **WHEN** user clicks "Agregar Tracking SLA" in visit detail
- **THEN** system displays form with fields: sla_id, status, value, notes
- **AND** on submit sends POST to /api/visit-sla-trackings with visit_id
