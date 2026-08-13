## Purpose

CRUD para gestionar direcciones geocodificadas con búsqueda de direcciones cercanas.

## ADDED Requirements

### Requirement: Addresses list view
The system SHALL display a paginated list of geocoded addresses.

#### Scenario: View addresses
- **WHEN** user navigates to /dashboard/geocoding/addresses
- **THEN** system displays table with columns: Street, Number, Commune, City, Lat, Lng
- **AND** supports search by street name or commune
- **AND** supports pagination

### Requirement: Address detail view
The system SHALL display address details with coordinates.

#### Scenario: View address detail
- **WHEN** user clicks an address row
- **THEN** system displays full address info
- **AND** displays map with marker at lat/lng coordinates
- **AND** shows nearby addresses within configurable radius

### Requirement: Nearby addresses search
The system SHALL allow searching for addresses near a given location.

#### Scenario: Search nearby addresses
- **WHEN** user enters lat, lng, and radius in search form
- **THEN** system sends GET to /api/addresses/nearby with query params
- **AND** displays results in table sorted by distance

### Requirement: Addresses CRUD
The system SHALL allow full CRUD for geocoded addresses.

#### Scenario: Create address
- **WHEN** user clicks "Crear Dirección"
- **THEN** system displays form with fields: street, number, commune, city, region, lat, lng
- **AND** on submit sends POST to /api/addresses

#### Scenario: Edit address
- **WHEN** user clicks edit on an address row
- **THEN** system displays pre-filled form
- **AND** on submit sends PUT to /api/addresses/:id

#### Scenario: Delete address
- **WHEN** user clicks delete on an address row
- **THEN** system displays confirmation dialog
- **AND** on confirm sends DELETE to /api/addresses/:id
