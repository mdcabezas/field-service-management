## Purpose

CRUD para gestionar usuarios del sistema y roles de técnicos.

## ADDED Requirements

### Requirement: Users list view
The system SHALL display a paginated list of users.

#### Scenario: View users
- **WHEN** admin navigates to /dashboard/core/users
- **THEN** system displays table with columns: Employee Number, Username, Role, Status
- **AND** supports pagination with limit/offset
- **AND** supports search by username or employee number

#### Scenario: Create user
- **WHEN** admin clicks "Crear Usuario" button
- **THEN** system displays form with fields: employee_number, username, password, role
- **AND** role field is a select with options: admin, manager, supervisor, technician
- **AND** on submit sends POST to /api/users

#### Scenario: Edit user
- **WHEN** admin clicks edit action on a user row
- **THEN** system displays pre-filled form with user data
- **AND** on submit sends PUT to /api/users/:id

#### Scenario: Delete user
- **WHEN** admin clicks delete action on a user row
- **THEN** system displays confirmation dialog
- **AND** on confirm sends DELETE to /api/users/:id

### Requirement: Tech Roles list view
The system SHALL display a paginated list of technician roles.

#### Scenario: View tech roles
- **WHEN** admin navigates to /dashboard/core/tech-roles
- **THEN** system displays table with columns: Name, Description
- **AND** supports CRUD operations similar to Users
