## 1. Fix TypeScript Interfaces

- [x] 1.1 Rewrite `types/api.ts` — Core domain (User, TechRole)
- [x] 1.2 Rewrite `types/api.ts` — Partners domain (Partner, PartnerContact, PartnerAgreement, PartnerAgreementDoc, PartnerAgreementForm, SLA)
- [x] 1.3 Rewrite `types/api.ts` — Customers domain (Customer, CustomerAddress, CustomerAddressPartner, Property, Technician, TechCertification)
- [x] 1.4 Rewrite `types/api.ts` — Geocoding domain (GeocodedAddress)
- [x] 1.5 Rewrite `types/api.ts` — Shared domain (ReportTemplate)
- [x] 1.6 Rewrite `types/api.ts` — Inventory domain (Material, Tool, EPPItem, Vehicle, Rental, CostRate, MaintenanceRecord, MaintenanceSchedule, ChecklistTemplate)
- [x] 1.7 Rewrite `types/api.ts` — Operations domain (Visit, VisitAssignment, VisitChecklistMaterial, VisitChecklistTool, VisitChecklistEPP, VisitMeasurement, VisitCheckpoint, VisitPhoto, VisitReport, ReportEntry, ReportImage, VehicleAssignment, VisitRental, VisitSLATracking, ChecklistTemplateMaterial, ChecklistTemplateTool, ChecklistTemplateEPP)
- [x] 1.8 Rewrite `types/api.ts` — Planning domain (DailyPlan, DailyPlanAssignment, DailyLoadMaterial, DailyLoadTool, DailyLoadEPP, Route)
- [x] 1.9 Rewrite `types/api.ts` — Notifications domain (NotificationTemplate, Notification)

## 2. Fix List Pages (column accessorKey + cell formatters)

- [x] 2.1 Fix operations/visits/page.tsx — `date` → `scheduled_at`
- [x] 2.2 Fix operations/vehicle-assignments/page.tsx — `start_date`/`end_date` → `departure_time`/`return_time`
- [x] 2.3 Fix operations/reports/page.tsx — `template_id` → `report_template_id`, remove phantom fields
- [x] 2.4 Fix operations/sla-trackings/page.tsx — update to match new fields
- [x] 2.5 Fix planning/daily-plans/page.tsx — remove phantom `technician_id`/`status`
- [x] 2.6 Fix planning/assignments/page.tsx — `visit_id` → `tech_id`
- [x] 2.7 Fix planning/routes/page.tsx — update to match new fields
- [x] 2.8 Fix partners/page.tsx — `rut` → `tax_id`
- [x] 2.9 Fix customers/page.tsx — `rut` → `tax_id`, remove phantom `status`
- [x] 2.10 Fix inventory pages — `code` → `sku` (materials), `plate` → `license_plate` (vehicles), `item_type` → `type` (rentals), `name`/`rate` → `value` (cost-rates)
- [x] 2.11 Fix shared/report-templates/page.tsx — `fields` → `fields_json`, remove phantom fields

## 3. Fix Detail Pages

- [x] 3.1 Fix operations/visits/[id]/page.tsx — update all field references
- [x] 3.2 Fix operations/vehicle-assignments/[id]/page.tsx — update field references
- [x] 3.3 Fix other operations detail pages as needed
- [x] 3.4 Fix planning detail pages as needed
- [x] 3.5 Fix partners detail pages — `rut` → `tax_id`
- [x] 3.6 Fix customers detail pages — `rut` → `tax_id`
- [x] 3.7 Fix inventory detail pages as needed
- [x] 3.8 Fix core detail pages as needed

## 4. Fix Forms

- [x] 4.1 Fix user-form.tsx — `username` → `name`, `employee_number` → `email`
- [x] 4.2 Fix partner-form.tsx — `rut` → `tax_id`
- [x] 4.3 Fix customer-form.tsx — `rut` → `tax_id`
- [x] 4.4 Fix other forms as needed based on field renames

## 5. Verify

- [x] 5.1 Run `npm run build` — confirm zero errors
- [ ] 5.2 Spot-check key pages render data correctly
