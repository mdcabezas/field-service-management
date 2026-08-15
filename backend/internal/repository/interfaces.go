package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"localis-backend/internal/model/core"
	"localis-backend/internal/model/customers"
	"localis-backend/internal/model/geocoding"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/model/notifications"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/partners"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/model/shared"
)

type ListResult[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

type CoreUserRepository interface {
	GetByID(ctx context.Context, id string) (*core.User, error)
	List(ctx context.Context, limit, offset int) (*ListResult[core.User], error)
	Create(ctx context.Context, user *core.User) error
	Update(ctx context.Context, id string, user *core.User) error
	Delete(ctx context.Context, id string) error
}

type CoreTechRoleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*core.TechRole, error)
	List(ctx context.Context, limit, offset int) (*ListResult[core.TechRole], error)
	Create(ctx context.Context, role *core.TechRole) error
	Update(ctx context.Context, id uuid.UUID, role *core.TechRole) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CorePlanAuditLogRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*core.PlanAuditLog, error)
	List(ctx context.Context, limit, offset int) (*ListResult[core.PlanAuditLog], error)
	ListByEntity(ctx context.Context, entityType shared.AuditEntityType, entityID uuid.UUID) ([]core.PlanAuditLog, error)
	Create(ctx context.Context, entry *core.PlanAuditLog) error
}

type PartnerRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*partners.Partner, error)
	List(ctx context.Context, limit, offset int) (*ListResult[partners.Partner], error)
	Create(ctx context.Context, partner *partners.Partner) error
	Update(ctx context.Context, id uuid.UUID, partner *partners.Partner) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PartnerContactRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*partners.PartnerContact, error)
	ListByPartner(ctx context.Context, partnerID uuid.UUID) ([]partners.PartnerContact, error)
	Create(ctx context.Context, contact *partners.PartnerContact) error
	Update(ctx context.Context, id uuid.UUID, contact *partners.PartnerContact) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PartnerAgreementRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*partners.PartnerAgreement, error)
	ListByPartner(ctx context.Context, partnerID uuid.UUID) ([]partners.PartnerAgreement, error)
	Create(ctx context.Context, agreement *partners.PartnerAgreement) error
	Update(ctx context.Context, id uuid.UUID, agreement *partners.PartnerAgreement) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PartnerAgreementDocRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*partners.PartnerAgreementDoc, error)
	ListByAgreement(ctx context.Context, agreementID uuid.UUID) ([]partners.PartnerAgreementDoc, error)
	Create(ctx context.Context, doc *partners.PartnerAgreementDoc) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PartnerAgreementFormRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*partners.PartnerAgreementForm, error)
	ListByAgreement(ctx context.Context, agreementID uuid.UUID) ([]partners.PartnerAgreementForm, error)
	Create(ctx context.Context, form *partners.PartnerAgreementForm) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SLARepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*partners.SLA, error)
	ListByPartner(ctx context.Context, partnerID uuid.UUID) ([]partners.SLA, error)
	Create(ctx context.Context, sla *partners.SLA) error
	Update(ctx context.Context, id uuid.UUID, sla *partners.SLA) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CustomerRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*customers.Customer, error)
	List(ctx context.Context, limit, offset int) (*ListResult[customers.Customer], error)
	Create(ctx context.Context, customer *customers.Customer) error
	Update(ctx context.Context, id uuid.UUID, customer *customers.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CustomerAddressRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*customers.CustomerAddress, error)
	ListByCustomer(ctx context.Context, customerID uuid.UUID) ([]customers.CustomerAddress, error)
	Create(ctx context.Context, addr *customers.CustomerAddress) error
	Update(ctx context.Context, id uuid.UUID, addr *customers.CustomerAddress) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CustomerAddressPartnerRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*customers.CustomerAddressPartner, error)
	ListByCustomerAddress(ctx context.Context, caID uuid.UUID) ([]customers.CustomerAddressPartner, error)
	Create(ctx context.Context, cap *customers.CustomerAddressPartner) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PropertyRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*customers.Property, error)
	ListByCustomerAddress(ctx context.Context, caID uuid.UUID) ([]customers.Property, error)
	Create(ctx context.Context, prop *customers.Property) error
	Update(ctx context.Context, id uuid.UUID, prop *customers.Property) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TechnicianRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*customers.Technician, error)
	GetByUserID(ctx context.Context, userID string) (*customers.Technician, error)
	List(ctx context.Context, limit, offset int) (*ListResult[customers.Technician], error)
	Create(ctx context.Context, tech *customers.Technician) error
	Update(ctx context.Context, id uuid.UUID, tech *customers.Technician) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TechCertificationRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*customers.TechCertification, error)
	ListByTech(ctx context.Context, techID uuid.UUID) ([]customers.TechCertification, error)
	Create(ctx context.Context, cert *customers.TechCertification) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MaterialRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.Material, error)
	List(ctx context.Context, limit, offset int) (*ListResult[inventory.Material], error)
	Create(ctx context.Context, mat *inventory.Material) error
	Update(ctx context.Context, id uuid.UUID, mat *inventory.Material) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ToolRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.Tool, error)
	List(ctx context.Context, limit, offset int) (*ListResult[inventory.Tool], error)
	Create(ctx context.Context, tool *inventory.Tool) error
	Update(ctx context.Context, id uuid.UUID, tool *inventory.Tool) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EPPItemRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.EPPItem, error)
	List(ctx context.Context, limit, offset int) (*ListResult[inventory.EPPItem], error)
	Create(ctx context.Context, item *inventory.EPPItem) error
	Update(ctx context.Context, id uuid.UUID, item *inventory.EPPItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VehicleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.Vehicle, error)
	List(ctx context.Context, limit, offset int) (*ListResult[inventory.Vehicle], error)
	Create(ctx context.Context, vehicle *inventory.Vehicle) error
	Update(ctx context.Context, id uuid.UUID, vehicle *inventory.Vehicle) error
	UpdateInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, vehicle *inventory.Vehicle) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RentalRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.Rental, error)
	List(ctx context.Context, limit, offset int) (*ListResult[inventory.Rental], error)
	Create(ctx context.Context, rental *inventory.Rental) error
	Update(ctx context.Context, id uuid.UUID, rental *inventory.Rental) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CostRateRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.CostRate, error)
	List(ctx context.Context, limit, offset int) (*ListResult[inventory.CostRate], error)
	Create(ctx context.Context, rate *inventory.CostRate) error
	Update(ctx context.Context, id uuid.UUID, rate *inventory.CostRate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MaintenanceRecordRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.MaintenanceRecord, error)
	List(ctx context.Context, limit, offset int) (*ListResult[inventory.MaintenanceRecord], error)
	Create(ctx context.Context, rec *inventory.MaintenanceRecord) error
	Update(ctx context.Context, id uuid.UUID, rec *inventory.MaintenanceRecord) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type MaintenanceScheduleRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.MaintenanceSchedule, error)
	List(ctx context.Context, limit, offset int) (*ListResult[inventory.MaintenanceSchedule], error)
	ListActive(ctx context.Context) ([]inventory.MaintenanceSchedule, error)
	Create(ctx context.Context, sched *inventory.MaintenanceSchedule) error
	Update(ctx context.Context, id uuid.UUID, sched *inventory.MaintenanceSchedule) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ChecklistTemplateRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.ChecklistTemplate, error)
	List(ctx context.Context, limit, offset int) (*ListResult[inventory.ChecklistTemplate], error)
	Create(ctx context.Context, tmpl *inventory.ChecklistTemplate) error
	Update(ctx context.Context, id uuid.UUID, tmpl *inventory.ChecklistTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type DailyPlanRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyPlan, error)
	List(ctx context.Context, limit, offset int) (*ListResult[planning.DailyPlan], error)
	Create(ctx context.Context, plan *planning.DailyPlan) error
	Update(ctx context.Context, id uuid.UUID, plan *planning.DailyPlan) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type DailyPlanAssignmentRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyPlanAssignment, error)
	List(ctx context.Context, limit, offset int) (*ListResult[planning.DailyPlanAssignment], error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyPlanAssignment, error)
	Create(ctx context.Context, a *planning.DailyPlanAssignment) error
	Update(ctx context.Context, id uuid.UUID, a *planning.DailyPlanAssignment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type DailyLoadMaterialRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyLoadMaterial, error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyLoadMaterial, error)
	Create(ctx context.Context, m *planning.DailyLoadMaterial) error
	Update(ctx context.Context, id uuid.UUID, m *planning.DailyLoadMaterial) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type DailyLoadToolRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyLoadTool, error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyLoadTool, error)
	Create(ctx context.Context, t *planning.DailyLoadTool) error
	Update(ctx context.Context, id uuid.UUID, t *planning.DailyLoadTool) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type DailyLoadEPPRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyLoadEPP, error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyLoadEPP, error)
	Create(ctx context.Context, e *planning.DailyLoadEPP) error
	Update(ctx context.Context, id uuid.UUID, e *planning.DailyLoadEPP) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RouteRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*planning.Route, error)
	List(ctx context.Context, limit, offset int) (*ListResult[planning.Route], error)
	ListByPlan(ctx context.Context, planID uuid.UUID) ([]planning.Route, error)
	Create(ctx context.Context, route *planning.Route) error
	Update(ctx context.Context, id uuid.UUID, route *planning.Route) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.Visit, error)
	List(ctx context.Context, limit, offset int) (*ListResult[operations.Visit], error)
	Create(ctx context.Context, visit *operations.Visit) error
	Update(ctx context.Context, id uuid.UUID, visit *operations.Visit) error
	UpdateInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, visit *operations.Visit) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitAssignmentRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitAssignment, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitAssignment, error)
	Create(ctx context.Context, a *operations.VisitAssignment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ChecklistTemplateMaterialRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.ChecklistTemplateMaterial, error)
	ListByTemplate(ctx context.Context, templateID uuid.UUID) ([]operations.ChecklistTemplateMaterial, error)
	Create(ctx context.Context, item *operations.ChecklistTemplateMaterial) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ChecklistTemplateToolRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.ChecklistTemplateTool, error)
	ListByTemplate(ctx context.Context, templateID uuid.UUID) ([]operations.ChecklistTemplateTool, error)
	Create(ctx context.Context, item *operations.ChecklistTemplateTool) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ChecklistTemplateEPPRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.ChecklistTemplateEPP, error)
	ListByTemplate(ctx context.Context, templateID uuid.UUID) ([]operations.ChecklistTemplateEPP, error)
	Create(ctx context.Context, item *operations.ChecklistTemplateEPP) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitChecklistMaterialRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitChecklistMaterial, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitChecklistMaterial, error)
	Create(ctx context.Context, item *operations.VisitChecklistMaterial) error
	CreateBatch(ctx context.Context, items []operations.VisitChecklistMaterial) error
	Update(ctx context.Context, id uuid.UUID, item *operations.VisitChecklistMaterial) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitChecklistToolRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitChecklistTool, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitChecklistTool, error)
	Create(ctx context.Context, item *operations.VisitChecklistTool) error
	CreateBatch(ctx context.Context, items []operations.VisitChecklistTool) error
	Update(ctx context.Context, id uuid.UUID, item *operations.VisitChecklistTool) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitChecklistEPPRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitChecklistEPP, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitChecklistEPP, error)
	Create(ctx context.Context, item *operations.VisitChecklistEPP) error
	CreateBatch(ctx context.Context, items []operations.VisitChecklistEPP) error
	Update(ctx context.Context, id uuid.UUID, item *operations.VisitChecklistEPP) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitMaterialUsageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitMaterialUsage, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitMaterialUsage, error)
	Create(ctx context.Context, usage *operations.VisitMaterialUsage) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitToolUsageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitToolUsage, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitToolUsage, error)
	Create(ctx context.Context, usage *operations.VisitToolUsage) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitEPPUsageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitEPPUsage, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitEPPUsage, error)
	Create(ctx context.Context, usage *operations.VisitEPPUsage) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitMeasurementRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitMeasurement, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitMeasurement, error)
	Create(ctx context.Context, m *operations.VisitMeasurement) error
	Update(ctx context.Context, id uuid.UUID, m *operations.VisitMeasurement) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitCheckpointRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitCheckpoint, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitCheckpoint, error)
	Create(ctx context.Context, cp *operations.VisitCheckpoint) error
	CreateInTx(ctx context.Context, tx pgx.Tx, cp *operations.VisitCheckpoint) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitPhotoRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitPhoto, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitPhoto, error)
	Create(ctx context.Context, photo *operations.VisitPhoto) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitReportRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitReport, error)
	List(ctx context.Context, limit, offset int) (*ListResult[operations.VisitReport], error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitReport, error)
	Create(ctx context.Context, report *operations.VisitReport) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReportImageRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.ReportImage, error)
	ListByReport(ctx context.Context, reportID uuid.UUID) ([]operations.ReportImage, error)
	Create(ctx context.Context, img *operations.ReportImage) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByReport(ctx context.Context, reportID uuid.UUID) error
}

type ReportEntryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.ReportEntry, error)
	ListByReport(ctx context.Context, reportID uuid.UUID) ([]operations.ReportEntry, error)
	Create(ctx context.Context, entry *operations.ReportEntry) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByReport(ctx context.Context, reportID uuid.UUID) error
}

type VehicleAssignmentRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VehicleAssignment, error)
	List(ctx context.Context, limit, offset int) (*ListResult[operations.VehicleAssignment], error)
	Create(ctx context.Context, va *operations.VehicleAssignment) error
	Update(ctx context.Context, id uuid.UUID, va *operations.VehicleAssignment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitRentalRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitRental, error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitRental, error)
	Create(ctx context.Context, r *operations.VisitRental) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitSLATrackingRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitSLATracking, error)
	List(ctx context.Context, limit, offset int) (*ListResult[operations.VisitSLATracking], error)
	ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitSLATracking, error)
	Create(ctx context.Context, t *operations.VisitSLATracking) error
	Update(ctx context.Context, id uuid.UUID, t *operations.VisitSLATracking) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type NotificationTemplateRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*notifications.NotificationTemplate, error)
	List(ctx context.Context, limit, offset int) (*ListResult[notifications.NotificationTemplate], error)
	Create(ctx context.Context, tmpl *notifications.NotificationTemplate) error
	Update(ctx context.Context, id uuid.UUID, tmpl *notifications.NotificationTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type NotificationRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*notifications.Notification, error)
	List(ctx context.Context, limit, offset int) (*ListResult[notifications.Notification], error)
	Create(ctx context.Context, notif *notifications.Notification) error
	Update(ctx context.Context, id uuid.UUID, notif *notifications.Notification) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type GeocodingAddressRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*geocoding.Address, error)
	List(ctx context.Context, limit, offset int) (*ListResult[geocoding.Address], error)
	Create(ctx context.Context, addr *geocoding.Address) error
	Update(ctx context.Context, id uuid.UUID, addr *geocoding.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindNearby(ctx context.Context, lat, lng, radiusMeters float64) ([]geocoding.Address, error)
}

type ReportTemplateRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*shared.ReportTemplate, error)
	List(ctx context.Context, limit, offset int) (*ListResult[shared.ReportTemplate], error)
	Create(ctx context.Context, tmpl *shared.ReportTemplate) error
	Update(ctx context.Context, id uuid.UUID, tmpl *shared.ReportTemplate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type VisitTypeRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitType, error)
	List(ctx context.Context, limit, offset int) (*ListResult[operations.VisitType], error)
	Create(ctx context.Context, vt *operations.VisitType) error
	Update(ctx context.Context, id uuid.UUID, vt *operations.VisitType) error
	Delete(ctx context.Context, id uuid.UUID) error
}
