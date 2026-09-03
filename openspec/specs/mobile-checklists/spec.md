# Mobile Checklists Specification

## Purpose

The checklists by configurable stages system enables templates to define ordered stages (e.g., "Pre-departure", "On-site", "Post-work"). When creating a visit, the planner selects which stages to include. Required items in a stage block status transitions.

## Requirements

### Requirement: Checklist Stages
New table `inventory.checklist_stages` (id, template_id, name, sort_order, required, created_at). New field on templates: `stages_enabled` (boolean, default false). If stages_enabled = false: flat checklist. If stages_enabled = true: items grouped by stage.

#### Scenario: Stages are enabled on template
- **WHEN** template has stages_enabled = true
- **THEN** checklist items are grouped by stage
- **AND** each stage has a name and sort_order

### Requirement: Stage Management
Create/update/delete stages on template. Reorder stages via sort_order. Each stage has a name and can be marked as required (blocks status transition if incomplete).

#### Scenario: Stage is created
- **WHEN** planner creates a stage on template
- **THEN** stage is created with name, sort_order, and required flag

### Requirement: Stage Selection on Visit Creation
When template has stages_enabled = true, planner selects which stages to include. Default: all stages selected. Emergency visits can skip pre-departure stage.

#### Scenario: Planner selects stages
- **WHEN** planner creates a visit with stages_enabled template
- **THEN** planner selects which stages to include
- **AND** default is all stages selected

### Requirement: Checklist Population
`PopulateChecklist` filtered by selected stages. If no stages selected: populate all (backward compatible). If stages selected: only populate items belonging to selected stages.

#### Scenario: Checklist is populated with stages
- **WHEN** checklist is populated
- **AND** stages are selected
- **THEN** only items belonging to selected stages are populated

### Requirement: Status Transition Validation
When visiting status (e.g., "en_curso" → "completada"), check required stages for completeness. If any required stage has unchecked required items → block transition. Return specific error with blocking stages.

#### Scenario: Status transition is blocked
- **WHEN** technician tries to transition to "completada"
- **AND** required stage "Post-work" has 3 incomplete required items
- **THEN** transition is blocked
- **AND** error message: "Stage 'Post-work' has 3 incomplete required items"

### Requirement: Checklist Display by Stage
Mobile shows checklist grouped by stage. Stage header with name and progress (X/Y completed). Collapsible stages. Required indicator on items that block transitions. Color coding: green (complete), yellow (in progress), grey (not started).

#### Scenario: Technician views checklist by stage
- **WHEN** technician opens checklist
- **THEN** items are grouped by stage
- **AND** stage header shows name and progress
- **AND** stages are collapsible
