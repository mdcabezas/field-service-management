package shared

type VisitStatus string

const (
	VisitStatusScheduled   VisitStatus = "scheduled"
	VisitStatusEnRoute     VisitStatus = "en_route"
	VisitStatusInProgress  VisitStatus = "in_progress"
	VisitStatusCompleted   VisitStatus = "completed"
	VisitStatusCancelled   VisitStatus = "cancelled"
	VisitStatusRescheduled VisitStatus = "rescheduled"
)

type VisitPriority string

const (
	VisitPriorityLow       VisitPriority = "low"
	VisitPriorityNormal    VisitPriority = "normal"
	VisitPriorityUrgent    VisitPriority = "urgent"
	VisitPriorityEmergency VisitPriority = "emergency"
)

type VisitSource string

const (
	VisitSourceEmail    VisitSource = "email"
	VisitSourcePhone    VisitSource = "phone"
	VisitSourcePortal   VisitSource = "portal"
	VisitSourceWhatsapp VisitSource = "whatsapp"
	VisitSourceOther    VisitSource = "other"
)

type BillingTo string

const (
	BillingToPartner  BillingTo = "partner"
	BillingToCustomer BillingTo = "customer"
)

type PhotoStage string

const (
	PhotoStageDiagnosis PhotoStage = "diagnosis"
	PhotoStageExecution PhotoStage = "execution"
	PhotoStageReview    PhotoStage = "review"
	PhotoStageReport    PhotoStage = "report"
	PhotoStageOther     PhotoStage = "other"
)

type CheckpointType string

const (
	CheckpointTypeArrival    CheckpointType = "arrival"
	CheckpointTypeDeparture  CheckpointType = "departure"
	CheckpointTypeReportSent CheckpointType = "report_sent"
)

type ReportSource string

const (
	ReportSourceSystem ReportSource = "system"
	ReportSourcePaper  ReportSource = "paper"
)

type MeasurementResult string

const (
	MeasurementResultApproved    MeasurementResult = "approved"
	MeasurementResultRejected    MeasurementResult = "rejected"
	MeasurementResultPending     MeasurementResult = "pending"
	MeasurementResultObservation MeasurementResult = "observation"
)

type VehicleStatus string

const (
	VehicleStatusAvailable   VehicleStatus = "available"
	VehicleStatusInUse       VehicleStatus = "in_use"
	VehicleStatusMaintenance VehicleStatus = "maintenance"
	VehicleStatusRetired     VehicleStatus = "retired"
)

type ToolStatus string

const (
	ToolStatusAvailable   ToolStatus = "available"
	ToolStatusInUse       ToolStatus = "in_use"
	ToolStatusMaintenance ToolStatus = "maintenance"
	ToolStatusRetired     ToolStatus = "retired"
)

type TechStatus string

const (
	TechStatusActive   TechStatus = "active"
	TechStatusInactive TechStatus = "inactive"
	TechStatusVacation TechStatus = "vacation"
	TechStatusLeave    TechStatus = "leave"
)

type RentalType string

const (
	RentalTypeVehicle   RentalType = "vehicle"
	RentalTypeTool      RentalType = "tool"
	RentalTypeEquipment RentalType = "equipment"
)

type RouteStatus string

const (
	RouteStatusScheduled  RouteStatus = "scheduled"
	RouteStatusInProgress RouteStatus = "in_progress"
	RouteStatusCompleted  RouteStatus = "completed"
	RouteStatusCancelled  RouteStatus = "cancelled"
)

type EPPUsageStatus string

const (
	EPPUsageStatusUsed      EPPUsageStatus = "used"
	EPPUsageStatusNotNeeded EPPUsageStatus = "not_needed"
	EPPUsageStatusMissing   EPPUsageStatus = "missing"
)

type RejectionCategory string

const (
	RejectionCategoryLogistics   RejectionCategory = "logistics"
	RejectionCategoryRegulatory  RejectionCategory = "regulatory"
	RejectionCategoryCustomer    RejectionCategory = "customer"
	RejectionCategoryOperational RejectionCategory = "operational"
)

type CostType string

const (
	CostTypeLabor            CostType = "labor"
	CostTypeVehicle          CostType = "vehicle"
	CostTypeToolDepreciation CostType = "tool_depreciation"
	CostTypeEPP              CostType = "epp"
	CostTypeOverhead         CostType = "overhead"
)

type CostUnit string

const (
	CostUnitHour  CostUnit = "hour"
	CostUnitKm    CostUnit = "km"
	CostUnitDay   CostUnit = "day"
	CostUnitVisit CostUnit = "visit"
	CostUnitUnit  CostUnit = "unit"
)

type CustomerAddressType string

const (
	CustomerAddressTypeResidential   CustomerAddressType = "residential"
	CustomerAddressTypeCommercial    CustomerAddressType = "commercial"
	CustomerAddressTypeIndustrial    CustomerAddressType = "industrial"
	CustomerAddressTypeInstitutional CustomerAddressType = "institutional"
)

type PartnerStatus string

const (
	PartnerStatusActive   PartnerStatus = "active"
	PartnerStatusInactive PartnerStatus = "inactive"
)

type EPPType string

const (
	EPPTypeHead        EPPType = "head"
	EPPTypeHands       EPPType = "hands"
	EPPTypeBody        EPPType = "body"
	EPPTypeEyes        EPPType = "eyes"
	EPPTypeEars        EPPType = "ears"
	EPPTypeFeet        EPPType = "feet"
	EPPTypeRespiratory EPPType = "respiratory"
)

type EPPLifecycle string

const (
	EPPLifecycleDisposable EPPLifecycle = "disposable"
	EPPLifecycleReusable   EPPLifecycle = "reusable"
)

type NotificationType string

const (
	NotificationTypeVisitAssigned        NotificationType = "visit_assigned"
	NotificationTypeVisitReminder        NotificationType = "visit_reminder"
	NotificationTypeVisitCancelled       NotificationType = "visit_cancelled"
	NotificationTypeVisitRescheduled     NotificationType = "visit_rescheduled"
	NotificationTypeChecklistPending     NotificationType = "checklist_pending"
	NotificationTypeReportPending        NotificationType = "report_pending"
	NotificationTypeSLAWarning           NotificationType = "sla_warning"
	NotificationTypeSLAExpired           NotificationType = "sla_expired"
	NotificationTypeMaintenanceScheduled NotificationType = "maintenance_scheduled"
	NotificationTypeCertificationExpired NotificationType = "certification_expired"
	NotificationTypeOther                NotificationType = "other"
)

type NotificationChannel string

const (
	NotificationChannelEmail    NotificationChannel = "email"
	NotificationChannelSMS      NotificationChannel = "sms"
	NotificationChannelWhatsapp NotificationChannel = "whatsapp"
	NotificationChannelPush     NotificationChannel = "push"
	NotificationChannelInApp    NotificationChannel = "in_app"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
	NotificationStatusRead    NotificationStatus = "read"
)

type AuditEntityType string

const (
	AuditEntityTypeDailyLoadMaterial      AuditEntityType = "daily_load_material"
	AuditEntityTypeDailyLoadTool          AuditEntityType = "daily_load_tool"
	AuditEntityTypeDailyLoadEPP           AuditEntityType = "daily_load_epp"
	AuditEntityTypeVisitChecklistMaterial AuditEntityType = "visit_checklist_material"
	AuditEntityTypeVisitChecklistTool     AuditEntityType = "visit_checklist_tool"
	AuditEntityTypeVisitChecklistEPP      AuditEntityType = "visit_checklist_epp"
	AuditEntityTypeVisitMaterialUsage     AuditEntityType = "visit_material_usage"
	AuditEntityTypeVisitToolUsage         AuditEntityType = "visit_tool_usage"
	AuditEntityTypeVisitEPPUsage          AuditEntityType = "visit_epp_usage"
	AuditEntityTypeVisit                  AuditEntityType = "visit"
	AuditEntityTypeVisitAssignment        AuditEntityType = "visit_assignment"
	AuditEntityTypeDailyPlan              AuditEntityType = "daily_plan"
	AuditEntityTypeDailyPlanAssignment    AuditEntityType = "daily_plan_assignment"
	AuditEntityTypeVehicleAssignment      AuditEntityType = "vehicle_assignment"
	AuditEntityTypeVisitSLATracking       AuditEntityType = "visit_sla_tracking"
	AuditEntityTypeMaintenanceRecord      AuditEntityType = "maintenance_record"
	AuditEntityTypeMaintenanceSchedule    AuditEntityType = "maintenance_schedule"
	AuditEntityTypeVisitReport            AuditEntityType = "visit_report"
	AuditEntityTypeNotification           AuditEntityType = "notification"
	AuditEntityTypePartner                AuditEntityType = "partner"
	AuditEntityTypeCustomer               AuditEntityType = "customer"
	AuditEntityTypeTechnician             AuditEntityType = "technician"
	AuditEntityTypeVehicle                AuditEntityType = "vehicle"
	AuditEntityTypeMaterial               AuditEntityType = "material"
	AuditEntityTypeTool                   AuditEntityType = "tool"
	AuditEntityTypeEPPItem                AuditEntityType = "epp_item"
	AuditEntityTypeUser                   AuditEntityType = "user"
)

type AuditAction string

const (
	AuditActionAdd    AuditAction = "add"
	AuditActionUpdate AuditAction = "update"
	AuditActionDelete AuditAction = "delete"
)

type MaintenanceType string

const (
	MaintenanceTypeVehicle   MaintenanceType = "vehicle"
	MaintenanceTypeTool      MaintenanceType = "tool"
	MaintenanceTypeEquipment MaintenanceType = "equipment"
)
