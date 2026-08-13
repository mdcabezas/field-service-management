## Purpose

CRUD para gestionar clientes, sus direcciones, propiedades, asignaciones de socios, técnicos y certificaciones.

## ADDED Requirements

### Requirement: Customers list view
The system SHALL display a paginated list of customers.

#### Scenario: View customers
- **WHEN** user navigates to /dashboard/customers
- **THEN** system displays table with columns: Name, RUT, Phone, Email, Status
- **AND** supports search by name or RUT

### Requirement: Customer detail view
The system SHALL display customer details with nested sub-resources.

#### Scenario: View customer detail
- **WHEN** user clicks a customer row
- **THEN** system displays customer info at top
- **AND** displays tabs: Addresses, Technicians
- **AND** each tab shows respective list with CRUD actions

### Requirement: Customer Addresses management
The system SHALL allow CRUD for customer addresses with partner assignments and properties.

#### Scenario: Create address
- **WHEN** user clicks "Crear Dirección" in customer detail
- **THEN** system displays form with fields: street, number, commune, city, region, lat, lng
- **AND** on submit sends POST to /api/customer-addresses with customer_id

#### Scenario: Manage address partners
- **WHEN** user opens an address detail
- **THEN** system displays Partner Assignments tab
- **AND** allows linking partners via POST to /api/customer-address-partners

#### Scenario: Manage address properties
- **WHEN** user opens an address detail
- **THEN** system displays Properties tab with list
- **AND** allows CRUD via /api/properties endpoints

### Requirement: Technicians management
The system SHALL allow CRUD for technicians and their certifications.

#### Scenario: Create technician
- **WHEN** user clicks "Crear Técnico" in customer detail
- **THEN** system displays form with employee_number, name, phone, tech_role_id
- **AND** on submit sends POST to /api/technicians with customer_id

#### Scenario: Manage technician certifications
- **WHEN** user opens a technician detail
- **THEN** system displays Certifications tab
- **AND** allows adding certifications via POST to /api/tech-certifications
