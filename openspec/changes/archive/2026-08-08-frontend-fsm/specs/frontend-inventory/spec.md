## Purpose

CRUD para gestionar materiales, herramientas, EPP, vehículos, arriendos, costos, mantenimiento y plantillas de checklist.

## ADDED Requirements

### Requirement: Materials list view
The system SHALL display a paginated list of materials.

#### Scenario: View materials
- **WHEN** user navigates to /dashboard/inventory/materials
- **THEN** system displays table with columns: Name, Code, Unit, Stock, Min Stock
- **AND** supports search by name or code
- **AND** highlights rows where stock < min_stock

### Requirement: Tools list view
The system SHALL display a paginated list of tools.

#### Scenario: View tools
- **WHEN** user navigates to /dashboard/inventory/tools
- **THEN** system displays table with columns: Name, Code, Status, Assigned To
- **AND** supports search and filter by status

### Requirement: EPP Items list view
The system SHALL display a paginated list of personal protective equipment.

#### Scenario: View EPP items
- **WHEN** user navigates to /dashboard/inventory/epp
- **THEN** system displays table with columns: Name, Code, Category, Stock, Expiry Date
- **AND** highlights items nearing expiry

### Requirement: Vehicles list view
The system SHALL display a paginated list of vehicles.

#### Scenario: View vehicles
- **WHEN** user navigates to /dashboard/inventory/vehicles
- **THEN** system displays table with columns: Plate, Type, Status, Assigned To
- **AND** supports CRUD operations

### Requirement: Rentals management
The system SHALL allow CRUD for vehicle and equipment rentals.

#### Scenario: Create rental
- **WHEN** user clicks "Crear Arriendo"
- **THEN** system displays form with fields: item_type, item_id, start_date, end_date, cost
- **AND** on submit sends POST to /api/rentals

### Requirement: Cost Rates management
The system SHALL allow CRUD for cost rates per hour/day.

#### Scenario: View cost rates
- **WHEN** user navigates to /dashboard/inventory/cost-rates
- **THEN** system displays table with columns: Name, Rate, Unit, Effective Date
- **AND** supports CRUD operations

### Requirement: Maintenance Records and Schedules
The system SHALL allow viewing maintenance history and scheduling future maintenance.

#### Scenario: View maintenance records
- **WHEN** user navigates to /dashboard/inventory/maintenance
- **THEN** system displays tabs: Records, Schedules
- **AND** Records tab shows completed maintenance with date, type, cost, vehicle/tool
- **AND** Schedules tab shows upcoming maintenance with status

### Requirement: Checklist Templates management
The system SHALL allow CRUD for checklist templates with associated materials, tools, and EPPs.

#### Scenario: Create checklist template
- **WHEN** user clicks "Crear Plantilla de Checklist"
- **THEN** system displays form with name, description, visit_type
- **AND** allows adding materials, tools, EPPs via sub-forms
- **AND** on submit sends POST to /api/checklist-templates with nested resources
