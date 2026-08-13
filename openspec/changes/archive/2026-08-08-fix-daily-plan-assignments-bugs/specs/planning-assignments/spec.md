## Purpose

Daily-plan technician assignments let an admin or manager attach technicians (Users with `role=technician`) and TechRoles to a DailyPlan for routing work. This capability covers the full editing flow: list all assignments, create a new assignment, edit an existing one, delete it, and surface technician/role choices as dropdowns rather than raw UUIDs.

## ADDED Requirements

### Requirement: List all assignments

The system SHALL display a list of all daily-plan assignments across all daily plans, regardless of which daily plan is selected. The list MUST NOT be silently empty when no daily plan is selected.

#### Scenario: User opens the assignments page

- **WHEN** an admin or manager navigates to `/planning/assignments`
- **THEN** the system fetches every assignment via `GET /api/daily-plan-assignments` and renders them in a table
- **AND** the table is NOT empty when assignments exist in the database

#### Scenario: Backend has no assignments

- **WHEN** the user opens `/planning/assignments`
- **AND** the database contains zero assignments
- **THEN** the system renders an empty-state message (e.g. "No hay asignaciones") instead of indefinitely showing "Cargando..."

### Requirement: Create assignment preserves user-entered daily_plan_id

The system SHALL send the user-supplied `daily_plan_id` to the backend when creating an assignment. The create flow MUST NOT overwrite `daily_plan_id` with an empty string or any value not chosen by the user.

#### Scenario: User creates an assignment

- **WHEN** the user fills the create form selecting a daily plan, technician, and TechRole
- **AND** submits the form
- **THEN** the system POSTs `{ daily_plan_id, tech_id, role_id }` to `/api/daily-plan-assignments` with the user-chosen `daily_plan_id`
- **AND** the backend returns 201 Created
- **AND** the user is redirected to `/planning/assignments`

#### Scenario: User selects an invalid daily plan

- **WHEN** the user selects a daily plan that does not exist
- **AND** submits the form
- **THEN** the system surfaces the backend error as a toast saying "Error al crear asignación"
- **AND** the user remains on the form

### Requirement: Technician selector dropdown

The technician field on the create and edit forms SHALL be a dropdown populated with all Users whose `role === "technician"`. The dropdown MUST NOT accept raw UUID text input.

#### Scenario: User opens the create form

- **WHEN** the user navigates to `/planning/assignments/new`
- **THEN** the "Técnico" field is a `<select>` populated from `GET /api/users` filtered to `role=technician`
- **AND** each option's label is the user's name and its value is the user's id

#### Scenario: No technicians exist

- **WHEN** the user opens the create form
- **AND** the technicians query returns an empty list
- **THEN** the dropdown renders a single disabled option "No hay técnicos disponibles"
- **AND** the submit button is disabled

### Requirement: TechRole selector dropdown

The TechRole field on the create and edit forms SHALL be a dropdown populated with all active TechRoles. The dropdown MUST NOT accept raw UUID text input.

#### Scenario: User opens the create form

- **WHEN** the user navigates to `/planning/assignments/new`
- **THEN** the "Rol" field is a `<select>` populated from `GET /api/tech-roles` filtered to `active=true`
- **AND** each option's label is the role's name and its value is the role's id

#### Scenario: No active TechRoles exist

- **WHEN** the user opens the create form
- **AND** the active TechRoles query returns an empty list
- **THEN** the dropdown renders a single disabled option "No hay roles disponibles"
- **AND** the submit button is disabled

### Requirement: DailyPlan selector dropdown

The daily-plan field on the create form SHALL be a dropdown populated with all DailyPlans. The dropdown MUST NOT accept raw UUID text input.

#### Scenario: User opens the create form

- **WHEN** the user navigates to `/planning/assignments/new`
- **THEN** the "Plan Diario" field is a `<select>` populated from `GET /api/daily-plans`
- **AND** each option's label is the plan's date (and notes if present) and its value is the plan's id

#### Scenario: No daily plans exist

- **WHEN** the user opens the create form
- **AND** the daily plans query returns an empty list
- **THEN** the dropdown renders a single disabled option "No hay planes diarios"
- **AND** the submit button is disabled
- **AND** the form shows a hint linking to `/planning/daily-plans/new`

### Requirement: Edit existing assignment

The system SHALL allow editing the technician or TechRole of an existing assignment via `PUT /api/daily-plan-assignments/:id`. The `daily_plan_id` is immutable after creation (editing it would be a move between plans; out of scope).

#### Scenario: User opens an existing assignment

- **WHEN** an admin or manager navigates to `/planning/assignments/{id}` for an existing assignment
- **THEN** the form is populated with the assignment's current `daily_plan_id` (disabled), `tech_id`, and `role_id`
- **AND** the form shows an "Actualizar" button (not "Crear")

#### Scenario: User edits and saves

- **WHEN** the user changes the technician or TechRole dropdown selection
- **AND** clicks "Actualizar"
- **THEN** the system PUTs `{ tech_id, role_id }` to `/api/daily-plan-assignments/{id}`
- **AND** the backend returns 200 OK with the updated assignment
- **AND** the user is redirected to `/planning/assignments`

#### Scenario: Backend has no PUT endpoint

- **WHEN** the backend is missing `PUT /api/daily-plan-assignments/:id`
- **THEN** an admin or manager attempting to edit receives a 404/405 error surfaced as a toast
- **AND** the assignment is not modified

### Requirement: Delete existing assignment

The system SHALL allow deleting an existing assignment via `DELETE /api/daily-plan-assignments/:id`. This requirement is unchanged from current behavior and is restated for completeness.

#### Scenario: User deletes an assignment

- **WHEN** an admin clicks "Eliminar" on `/planning/assignments/{id}`
- **AND** confirms the deletion dialog
- **THEN** the system DELETEs `/api/daily-plan-assignments/{id}`
- **AND** the user is redirected to `/planning/assignments`

### Requirement: Backend exposes list-all endpoint

The backend SHALL expose `GET /api/daily-plan-assignments` returning all assignments (paginated). This endpoint is required for the list-all UI.

#### Scenario: Frontend lists all assignments

- **WHEN** the frontend calls `GET /api/daily-plan-assignments`
- **THEN** the backend returns `{ items: DailyPlanAssignment[], total: number }`
- **AND** the response is paginated (default limit 50, offset 0)

### Requirement: Backend exposes update endpoint

The backend SHALL expose `PUT /api/daily-plan-assignments/:id` accepting `{ tech_id, role_id }` and returning the updated assignment. The `daily_plan_id` field is not accepted on update.

#### Scenario: Frontend updates an assignment

- **WHEN** the frontend calls `PUT /api/daily-plan-assignments/{id}` with `{ tech_id, role_id }`
- **THEN** the backend updates the assignment's `tech_id` and `role_id`
- **AND** returns 200 OK with the updated assignment

#### Scenario: Assignment not found

- **WHEN** the frontend calls `PUT /api/daily-plan-assignments/{id}` for a non-existent id
- **THEN** the backend returns 404 Not Found
