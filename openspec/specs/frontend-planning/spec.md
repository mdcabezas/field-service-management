## Purpose

CRUD para gestionar planes diarios, asignaciones, carga de materiales/herramientas/EPP y rutas.

## ADDED Requirements

### Requirement: Daily Plans list view
The system SHALL display a paginated list of daily plans.

#### Scenario: View daily plans
- **WHEN** user navigates to /dashboard/planning/daily-plans
- **THEN** system displays table with columns: Date, Technician, Status, Created
- **AND** supports filter by date range and technician
- **AND** supports search by technician name

### Requirement: Daily Plan detail view
The system SHALL display daily plan details with nested sub-resources.

#### Scenario: View daily plan detail
- **WHEN** user clicks a daily plan row
- **THEN** system displays plan info at top
- **AND** displays tabs: Assignments, Materials, Tools, EPPs
- **AND** each tab shows respective list with CRUD actions

### Requirement: Daily Plan Assignments management
The system SHALL allow CRUD for visit assignments within a daily plan.

#### Scenario: Create assignment
- **WHEN** user clicks "Crear Asignación" in daily plan detail
- **THEN** system displays form with fields: visit_id, time_slot, priority
- **AND** on submit sends POST to /api/daily-plan-assignments with daily_plan_id

### Requirement: Daily Load Materials management
The system SHALL allow CRUD for materials loaded for a daily plan.

#### Scenario: Add material to plan
- **WHEN** user clicks "Agregar Material" in daily plan detail
- **THEN** system displays form with fields: material_id, quantity
- **AND** on submit sends POST to /api/daily-load-materials with daily_plan_id

### Requirement: Daily Load Tools management
The system SHALL allow CRUD for tools loaded for a daily plan.

#### Scenario: Add tool to plan
- **WHEN** user clicks "Agregar Herramienta" in daily plan detail
- **THEN** system displays form with fields: tool_id, quantity
- **AND** on submit sends POST to /api/daily-load-tools with daily_plan_id

### Requirement: Daily Load EPPs management
The system SHALL allow CRUD for EPPs loaded for a daily plan.

#### Scenario: Add EPP to plan
- **WHEN** user clicks "Agregar EPP" in daily plan detail
- **THEN** system displays form with fields: epp_item_id, quantity
- **AND** on submit sends POST to /api/daily-load-epps with daily_plan_id

### Requirement: Routes management
The system SHALL allow CRUD for routes used in daily planning.

#### Scenario: View routes
- **WHEN** user navigates to /dashboard/planning/routes
- **THEN** system displays table with columns: Name, Description, Stops Count
- **AND** supports CRUD operations
