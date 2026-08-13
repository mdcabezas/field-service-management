## Purpose

CRUD para gestionar socios comerciales, sus contactos, acuerdos, documentos, formularios y SLAs.

## ADDED Requirements

### Requirement: Partners list view
The system SHALL display a paginated list of partner companies.

#### Scenario: View partners
- **WHEN** user navigates to /dashboard/partners
- **THEN** system displays table with columns: Name, RUT, Status, Created
- **AND** supports search by name or RUT
- **AND** each row has expand icon to view contacts, agreements, SLAs

### Requirement: Partner detail view
The system SHALL display partner details with nested sub-resources.

#### Scenario: View partner detail
- **WHEN** user clicks a partner row
- **THEN** system displays partner info at top
- **AND** displays tabs: Contacts, Agreements, SLAs
- **AND** each tab shows respective list with CRUD actions

### Requirement: Partner Contacts management
The system SHALL allow CRUD for partner contacts.

#### Scenario: Create contact
- **WHEN** user clicks "Crear Contacto" in partner detail
- **THEN** system displays form with fields: name, email, phone, position
- **AND** on submit sends POST to /api/partner-contacts with partner_id

### Requirement: Partner Agreements management
The system SHALL allow CRUD for partner agreements with documents and forms.

#### Scenario: Create agreement
- **WHEN** user clicks "Crear Acuerdo" in partner detail
- **THEN** system displays form with fields: type, start_date, end_date, status
- **AND** on submit sends POST to /api/partner-agreements with partner_id

#### Scenario: Manage agreement documents
- **WHEN** user opens an agreement detail
- **THEN** system displays Documents tab with list
- **AND** allows upload of new documents via POST to /api/partner-agreement-docs

#### Scenario: Manage agreement forms
- **WHEN** user opens an agreement detail
- **THEN** system displays Forms tab with list
- **AND** allows creation of new forms via POST to /api/partner-agreement-forms

### Requirement: SLAs management
The system SHALL allow CRUD for partner SLAs.

#### Scenario: Create SLA
- **WHEN** user clicks "Crear SLA" in partner detail
- **THEN** system displays form with fields: type, target_value, unit, period
- **AND** on submit sends POST to /api/slas with partner_id
