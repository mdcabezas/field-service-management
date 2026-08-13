## Purpose

Define el layout autenticado con sidebar, topbar, navegación por módulos y el estilo visual brutalista del sistema.

## ADDED Requirements

### Requirement: Authenticated layout
The system SHALL display a sidebar, topbar, and main content area for authenticated users.

#### Scenario: Layout structure
- **WHEN** user is authenticated and on any /dashboard/* route
- **THEN** system displays fixed sidebar on left (240px width)
- **AND** displays topbar at top (64px height)
- **AND** displays main content area filling remaining space

### Requirement: Sidebar navigation
The system SHALL display navigation links for modules the user has access to.

#### Scenario: Admin sees all modules
- **WHEN** user role is admin
- **THEN** sidebar displays all modules: Core, Partners, Customers, Inventory, Planning, Operations, Notifications, Geocoding, Shared

#### Scenario: Non-admin sees permitted modules only
- **WHEN** user role is manager or supervisor
- **THEN** sidebar displays only modules the role has access to
- **AND** unauthorized modules are hidden, not just disabled

### Requirement: Topbar
The system SHALL display a topbar with user info and logout action.

#### Scenario: Topbar content
- **WHEN** user is authenticated
- **THEN** topbar displays system name "Localis FSM"
- **AND** displays current user's employee number
- **AND** displays logout button

### Requirement: Active route highlighting
The system SHALL highlight the current module in the sidebar.

#### Scenario: Active module
- **WHEN** user is on /dashboard/customers/*
- **THEN** Customers link in sidebar has distinct visual style (background color change)
- **AND** parent module is expanded if nested

### Requirement: Breadcrumbs
The system SHALL display breadcrumbs showing current location.

#### Scenario: Nested route breadcrumbs
- **WHEN** user is on /dashboard/customers/123/addresses
- **THEN** breadcrumb shows: Dashboard > Customers > 123 > Addresses

### Requirement: Brutalist visual style
The system SHALL use a brutalist design aesthetic throughout.

#### Scenario: Visual elements
- **WHEN** any UI component is rendered
- **THEN** borders are 2px solid black (border-2 border-black)
- **AND** no border-radius (rounded-none)
- **AND** high contrast colors (black/white/primary)
- **AND** monospace font for labels and headings
- **AND** no shadows or gradients
