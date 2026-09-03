# Mobile Partners Specification

## Purpose

The partner-specific customization system enables partner-scoped templates, auto-resolution of partners on visit creation, and custom data fields. Different partners require different checklists, reports, SLAs, and data fields.

## Requirements

### Requirement: Partner-Scoped Checklist Templates
New field `ChecklistTemplate.PartnerID UUID` (nullable). NULL = global template (applies to all partners). UUID = template specific to that partner. Partner-specific templates override globals.

#### Scenario: Partner-specific templates are shown
- **WHEN** planner selects a property with partner assigned
- **THEN** partner-specific templates are shown first
- **AND** global templates are shown after
- **AND** partner-specific templates override globals (same name = partner wins)

### Requirement: Partner-Scoped Report Templates
New fields `ReportTemplate.PartnerID UUID` (nullable) and `ReportTemplate.WorkType TEXT` (nullable). Filtering by partner + work_type when creating report.

#### Scenario: Report templates are filtered by partner
- **WHEN** technician creates a report
- **THEN** report templates are filtered by partner + work_type
- **AND** partner-specific reports appear first in selection

### Requirement: Auto-Resolve Partner on Visit Creation
When planner selects a property, backend resolves partner via `customer_address_partners`. `partner_id` auto-set on visit. If multiple partners: show selector. If no partner: leave null.

#### Scenario: Partner is auto-resolved
- **WHEN** planner selects a property
- **THEN** backend resolves partner via `customer_address_partners`
- **AND** `partner_id` is auto-set on visit
- **AND** if multiple partners, selector is shown

### Requirement: Auto-Resolve SLA
When visit has `partner_id` + visit `type`, backend searches `partners.slas` for active SLA matching partner + work_type. Auto-creates `VisitSLATracking` record.

#### Scenario: SLA is auto-linked
- **WHEN** visit has partner_id and type
- **THEN** backend searches for active SLA matching partner + work_type
- **AND** auto-creates VisitSLATracking record

### Requirement: Partner Data Fields on Visit
`partner_order_id` (external order reference) and `partner_supervisor` (partner's supervisor name). Fields visible/editable only when partner is assigned. Mobile shows partner info at top of visit detail.

#### Scenario: Partner data fields are shown
- **WHEN** visit has partner assigned
- **THEN** partner_order_id and partner_supervisor are shown
- **AND** fields are editable

### Requirement: Partner Contacts Reference
When visit has partner_id, show partner contacts as reference. Endpoint: `GET /api/partners/:id/contacts`. Mobile: collapsible section "Contactos del partner".

#### Scenario: Partner contacts are shown
- **WHEN** visit has partner assigned
- **THEN** partner contacts are shown in collapsible section
- **AND** contacts include name, role, phone, email
