// ============================================================================
// API Types — Derived from backend Go struct JSON tags
// ============================================================================

// Auth
export interface LoginRequest {
  employee_number: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
}

export interface RefreshRequest {
  refresh_token: string;
}

export interface RefreshResponse {
  access_token: string;
  refresh_token: string;
}

export interface UserInfo {
  employee_number: string;
  role: UserRole;
}

// Enums
export type UserRole = "admin" | "manager" | "supervisor" | "technician";

// Pagination
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
}

// ============================================================================
// Core
// ============================================================================

export interface User {
  id: string;
  email: string;
  role: UserRole;
  name: string;
  created_at: string;
}

export interface TechRole {
  id: string;
  code: string;
  name: string;
  active: boolean;
  created_at: string;
}

// ============================================================================
// Partners
// ============================================================================

export interface Partner {
  id: string;
  name: string;
  tax_id?: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface PartnerContact {
  id: string;
  partner_id: string;
  name: string;
  position?: string;
  phone?: string;
  email?: string;
  created_at: string;
}

export interface PartnerAgreement {
  id: string;
  partner_id: string;
  service_type: string;
  rate: number;
  active: boolean;
  created_at: string;
}

export interface PartnerAgreementDoc {
  id: string;
  agreement_id: string;
  work_type: string;
  doc_name: string;
  required: boolean;
  created_at: string;
}

export interface PartnerAgreementForm {
  id: string;
  agreement_id: string;
  work_type: string;
  form_template_id: string;
  quantity: number;
  created_at: string;
}

export interface SLA {
  id: string;
  partner_id: string;
  name: string;
  description?: string;
  work_type: string;
  response_hours: number;
  resolution_hours: number;
  compliance_target: number;
  active: boolean;
  valid_from?: string;
  valid_until?: string;
  created_at: string;
}

// ============================================================================
// Customers
// ============================================================================

export interface Customer {
  id: string;
  name: string;
  tax_id?: string;
  phone?: string;
  email?: string;
  created_at: string;
  updated_at: string;
}

export interface CustomerAddress {
  id: string;
  customer_id: string;
  address_id: string;
  name?: string;
  type?: string;
  created_at: string;
}

export interface CustomerAddressPartner {
  id: string;
  customer_address_id: string;
  partner_id: string;
  from_date?: string;
  to_date?: string;
  notes?: string;
  created_at: string;
}

export interface Property {
  id: string;
  customer_address_id: string;
  name?: string;
  type?: string;
  notes?: string;
  created_at: string;
}

export interface Technician {
  id: string;
  user_id: string;
  name: string;
  is_active: boolean;
  specialties?: string[];
  created_at: string;
}

export interface TechCertification {
  id: string;
  tech_id: string;
  cert_id: string;
  number?: string;
  issuer?: string;
  issue_date?: string;
  expiry_date?: string;
  created_at: string;
}

// ============================================================================
// Geocoding
// ============================================================================

export interface GeocodedAddress {
  id: string;
  street: string;
  number: string;
  apartment?: string;
  neighborhood?: string;
  city: string;
  region: string;
  location_references?: string;
  postal_code?: string;
  created_at: string;
}

// ============================================================================
// Shared
// ============================================================================

export interface ReportTemplate {
  id: string;
  name: string;
  fields_json: Record<string, unknown>;
  created_at: string;
}

// ============================================================================
// Inventory
// ============================================================================

export interface Material {
  id: string;
  sku: string;
  name: string;
  unit: string;
  unit_cost: number;
  created_at: string;
}

export interface Tool {
  id: string;
  code: string;
  name: string;
  status: string;
  created_at: string;
}

export interface EPPItem {
  id: string;
  name: string;
  type: string;
  lifecycle: string;
  created_at: string;
}

export interface VehicleType {
  id: string;
  code: string;
  name: string;
  active: boolean;
}

export interface Vehicle {
  id: string;
  type?: string;
  license_plate: string;
  name: string;
  brand?: string;
  model?: string;
  year?: number;
  status: string;
  capacity?: number;
  created_at: string;
}

export interface Rental {
  id: string;
  type: string;
  supplier?: string;
  item_description?: string;
  daily_cost?: number;
  hourly_cost?: number;
  fixed_cost?: number;
  notes?: string;
  created_at: string;
}

export interface CostRate {
  id: string;
  type: string;
  reference_id?: string;
  reference_type?: string;
  value: number;
  unit: string;
  valid_from?: string;
  valid_until?: string;
  created_at: string;
}

export interface MaintenanceRecord {
  id: string;
  type: string;
  reference_id: string;
  date: string;
  cost?: number;
  supplier?: string;
  description?: string;
  next_date?: string;
  created_at: string;
}

export interface MaintenanceSchedule {
  id: string;
  type: string;
  reference_id: string;
  frequency_km?: number;
  frequency_days?: number;
  last_service_date?: string;
  last_service_km?: number;
  active: boolean;
  created_at: string;
}

export interface ChecklistTemplate {
  id: string;
  name: string;
  description?: string;
  visit_type?: string;
  work_type: string;
  created_at: string;
}

export interface ChecklistTemplateMaterial {
  id: string;
  template_id: string;
  material_id: string;
  default_quantity: number;
  required: boolean;
  created_at: string;
}

export interface ChecklistTemplateTool {
  id: string;
  template_id: string;
  tool_id: string;
  default_quantity: number;
  required: boolean;
  created_at: string;
}

export interface ChecklistTemplateEPP {
  id: string;
  template_id: string;
  epp_id: string;
  default_quantity: number;
  required: boolean;
  created_at: string;
}

// ============================================================================
// Operations
// ============================================================================

export interface Visit {
  id: string;
  property_id?: string;
  partner_id?: string;
  partner_order_id?: string;
  parent_visit_id?: string;
  route_id?: string;
  daily_plan_id?: string;
  type: string;
  status: string;
  priority: string;
  source?: string;
  source_reference?: string;
  billing_to: string;
  partner_supervisor?: string;
  result?: string;
  rejection_reason_id?: string;
  rejection_reason_detail?: string;
  scheduled_at: string;
  started_at?: string;
  completed_at?: string;
  alternative_location?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface VisitAssignment {
  id: string;
  visit_id: string;
  tech_id: string;
  role_id: string;
  created_at: string;
}

export interface VisitChecklistMaterial {
  id: string;
  visit_id: string;
  material_id: string;
  planned_quantity: number;
  confirmed: boolean;
  notes?: string;
  created_at: string;
}

export interface VisitChecklistTool {
  id: string;
  visit_id: string;
  tool_id: string;
  planned_quantity: number;
  confirmed: boolean;
  notes?: string;
  created_at: string;
}

export interface VisitChecklistEPP {
  id: string;
  visit_id: string;
  epp_id: string;
  planned_quantity: number;
  confirmed: boolean;
  notes?: string;
  created_at: string;
}

export interface VisitMeasurement {
  id: string;
  visit_id: string;
  type: string;
  value?: number;
  unit?: string;
  result?: string;
  measuring_device?: string;
  notes?: string;
  photo_url?: string;
  created_at: string;
}

export interface VisitCheckpoint {
  id: string;
  visit_id: string;
  type: string;
  geom?: string;
  timestamp?: string;
  device_info?: string;
  created_at: string;
}

export interface VisitPhoto {
  id: string;
  visit_id: string;
  checkpoint_id?: string;
  url: string;
  geom?: string;
  timestamp?: string;
  stage?: string;
  finding_type?: string;
  created_at: string;
}

export interface VisitReport {
  id: string;
  visit_id: string;
  report_template_id: string;
  source?: string;
  recorded_at?: string;
  created_at: string;
}

export interface ReportEntry {
  id: string;
  report_id: string;
  data_json: Record<string, unknown>;
  created_at: string;
}

export interface ReportImage {
  id: string;
  report_id: string;
  url: string;
  pages?: number;
  created_at: string;
}

export interface VehicleAssignment {
  id: string;
  vehicle_id: string;
  daily_plan_id?: string;
  visit_id?: string;
  departure_time?: string;
  return_time?: string;
  departure_mileage?: number;
  return_mileage?: number;
  notes?: string;
  created_at: string;
}

export interface VisitRental {
  id: string;
  visit_id: string;
  rental_id: string;
  start_time?: string;
  end_time?: string;
  hours?: number;
  total_cost?: number;
  reason?: string;
  created_at: string;
}

export interface VisitSLATracking {
  id: string;
  visit_id: string;
  sla_id: string;
  requested_at?: string;
  responded_at?: string;
  resolved_at?: string;
  response_time_hours?: number;
  resolution_time_hours?: number;
  meets_response_sla?: boolean;
  meets_resolution_sla?: boolean;
  created_at: string;
}

// ============================================================================
// Planning
// ============================================================================

export interface DailyPlan {
  id: string;
  name: string;
  date: string;
  notes?: string;
  created_at: string;
}

export interface DailyPlanAssignment {
  id: string;
  daily_plan_id: string;
  tech_id: string;
  role_id: string;
  created_at: string;
}

export interface DailyLoadMaterial {
  id: string;
  daily_plan_id: string;
  material_id: string;
  loaded_quantity: number;
  returned_quantity: number;
  created_at: string;
}

export interface DailyLoadTool {
  id: string;
  daily_plan_id: string;
  tool_id: string;
  loaded_quantity: number;
  returned_quantity: number;
  created_at: string;
}

export interface DailyLoadEPP {
  id: string;
  daily_plan_id: string;
  epp_id: string;
  loaded_quantity: number;
  returned_quantity: number;
  created_at: string;
}

export interface Route {
  id: string;
  type: string;
  type_name?: string;
  date: string;
  daily_plan_id: string;
  notes?: string;
  status: string;
  created_at: string;
}

export interface RouteType {
  id: string;
  code: string;
  name: string;
  active: boolean;
}

// ============================================================================
// Notifications
// ============================================================================

export interface NotificationTemplate {
  id: string;
  type: string;
  channel: string;
  subject?: string;
  body?: string;
  active: boolean;
  created_at: string;
}

export interface Notification {
  id: string;
  type: string;
  channel: string;
  recipient: string;
  subject?: string;
  body?: string;
  status: string;
  entity_type?: string;
  entity_id?: string;
  scheduled_for?: string;
  sent_at?: string;
  read_at?: string;
  attempts?: number;
  error?: string;
  created_at: string;
}

// API Error
export interface ApiError {
  error: string;
  message: string;
  statusCode: number;
}
