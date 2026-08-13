## Purpose

Controla el acceso a módulos y acciones según el rol del usuario (admin, manager, supervisor).

## ADDED Requirements

### Requirement: Role-based module access
The system SHALL restrict module visibility based on user role.

#### Scenario: Admin access
- **WHEN** user role is admin
- **THEN** user can access all modules and all actions within each module

#### Scenario: Manager access
- **WHEN** user role is manager
- **THEN** user can access Partners, Customers, Inventory, Planning, Operations, Notifications modules
- **AND** user cannot access Core (Users, Tech Roles) or system configuration

#### Scenario: Supervisor access
- **WHEN** user role is supervisor
- **THEN** user can access Operations module with limited actions (read, create visit assignments)
- **AND** user cannot access other modules

### Requirement: Action-level permissions
The system SHALL restrict actions within modules based on user role.

#### Scenario: Create/edit/delete restricted
- **WHEN** user role does not have create permission for a resource
- **THEN** "Create" button is hidden from the UI
- **AND** edit and delete actions are hidden from table rows

#### Scenario: Read-only access
- **WHEN** user role has read-only access to a module
- **THEN** user can view list and detail views
- **AND** all mutation actions (create, edit, delete) are hidden

### Requirement: Unauthorized action handling
The system SHALL prevent unauthorized actions via API.

#### Scenario: API returns 403
- **WHEN** user attempts an action their role doesn't permit
- **THEN** API returns 403 Forbidden
- **AND** frontend displays "No tiene permisos para realizar esta acción"
- **AND** user remains on current page

### Requirement: Role-based route protection
The system SHALL redirect users away from unauthorized routes.

#### Scenario: Accessing restricted route
- **WHEN** user navigates to a route their role doesn't have access to
- **THEN** system redirects to /dashboard
- **AND** displays toast notification "Acceso no autorizado"
