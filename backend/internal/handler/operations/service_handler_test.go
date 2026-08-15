package operationshandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"localis-backend/internal/model/customers"
	inventorymodel "localis-backend/internal/model/inventory"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/partners"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	coreservice "localis-backend/internal/service/core"
	notifservice "localis-backend/internal/service/notifications"
	operationssvc "localis-backend/internal/service/operations"
	"localis-backend/internal/testutil/mocks"
)

func newVehicleAssignmentSvc(t *testing.T) (*mocks.VehicleAssignmentRepository, *mocks.VehicleRepository, *operationssvc.VehicleAssignmentService) {
	t.Helper()
	vaRepo := mocks.NewVehicleAssignmentRepository(t)
	vRepo := mocks.NewVehicleRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	auditRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.PlanAuditLog")).Return(nil).Maybe()
	audit := coreservice.NewAuditService(auditRepo)
	svc := operationssvc.NewVehicleAssignmentService(vaRepo, vRepo, audit)
	return vaRepo, vRepo, svc
}

func TestVehicleAssignmentHandler_GetByID(t *testing.T) {
	vaRepo, _, svc := newVehicleAssignmentSvc(t)
	id := uuid.New()
	vaRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VehicleAssignment{ID: id}, nil)

	r := gin.New()
	r.GET("/vehicle-assignments/:id", NewVehicleAssignmentHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/vehicle-assignments/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVehicleAssignmentHandler_List(t *testing.T) {
	vaRepo, _, svc := newVehicleAssignmentSvc(t)
	vaRepo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[operations.VehicleAssignment]{Total: 0}, nil)

	r := gin.New()
	r.GET("/vehicle-assignments", NewVehicleAssignmentHandler(svc).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/vehicle-assignments", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVehicleAssignmentHandler_Create(t *testing.T) {
	vaRepo, vRepo, svc := newVehicleAssignmentSvc(t)
	vid := uuid.New()
	vRepo.EXPECT().GetByID(mockCtx(), vid).Return(&inventorymodel.Vehicle{ID: vid, Status: shared.VehicleStatusAvailable}, nil)
	vaRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VehicleAssignment")).Return(nil)
	vRepo.EXPECT().Update(mockCtx(), vid, mock.AnythingOfType("*inventory.Vehicle")).Return(nil)

	body := `{"vehicle_id":"` + vid.String() + `"}`
	r := gin.New()
	r.POST("/vehicle-assignments", NewVehicleAssignmentHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/vehicle-assignments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVehicleAssignmentHandler_Update(t *testing.T) {
	vaRepo, _, svc := newVehicleAssignmentSvc(t)
	id := uuid.New()
	vaRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VehicleAssignment{ID: id}, nil)
	vaRepo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*operations.VehicleAssignment")).Return(nil)

	vid := uuid.New()
	body := `{"vehicle_id":"` + vid.String() + `"}`
	r := gin.New()
	r.PUT("/vehicle-assignments/:id", NewVehicleAssignmentHandler(svc).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/vehicle-assignments/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVehicleAssignmentHandler_Delete(t *testing.T) {
	vaRepo, vRepo, svc := newVehicleAssignmentSvc(t)
	id := uuid.New()
	vehicleID := uuid.New()
	vaRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VehicleAssignment{ID: id, VehicleID: vehicleID}, nil)
	vaRepo.EXPECT().Delete(mockCtx(), id).Return(nil)
	vRepo.EXPECT().GetByID(mockCtx(), vehicleID).Return(&inventorymodel.Vehicle{ID: vehicleID, Status: shared.VehicleStatusInUse}, nil)
	vRepo.EXPECT().Update(mockCtx(), vehicleID, mock.Anything).Return(nil)

	r := gin.New()
	r.DELETE("/vehicle-assignments/:id", NewVehicleAssignmentHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/vehicle-assignments/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newVisitAssignmentSvc(t *testing.T) (*mocks.VisitAssignmentRepository, *mocks.VisitRepository, *mocks.TechnicianRepository, *operationssvc.VisitAssignmentService) {
	t.Helper()
	assignmentRepo := mocks.NewVisitAssignmentRepository(t)
	visitRepo := mocks.NewVisitRepository(t)
	techRepo := mocks.NewTechnicianRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	auditRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.PlanAuditLog")).Return(nil).Maybe()
	audit := coreservice.NewAuditService(auditRepo)
	svc := operationssvc.NewVisitAssignmentService(assignmentRepo, visitRepo, techRepo, audit)
	return assignmentRepo, visitRepo, techRepo, svc
}

func TestVisitAssignmentHandler_GetByID(t *testing.T) {
	assignmentRepo, _, _, svc := newVisitAssignmentSvc(t)
	id := uuid.New()
	assignmentRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitAssignment{ID: id}, nil)

	r := gin.New()
	r.GET("/visit-assignments/:id", NewVisitAssignmentHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visit-assignments/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitAssignmentHandler_ListByVisit(t *testing.T) {
	assignmentRepo, _, _, svc := newVisitAssignmentSvc(t)
	vid := uuid.New()
	assignmentRepo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitAssignment{}, nil)

	r := gin.New()
	r.GET("/visits/:id/assignments", NewVisitAssignmentHandler(svc).ListByVisit)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/assignments", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitAssignmentHandler_Create(t *testing.T) {
	assignmentRepo, visitRepo, techRepo, svc := newVisitAssignmentSvc(t)
	vid := uuid.New()
	tid := uuid.New()
	rid := uuid.New()
	visitRepo.EXPECT().GetByID(mockCtx(), vid).Return(&operations.Visit{ID: vid}, nil)
	techRepo.EXPECT().GetByID(mockCtx(), tid).Return(&customers.Technician{ID: tid, IsActive: true}, nil)
	assignmentRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitAssignment")).Return(nil)

	body := `{"visit_id":"` + vid.String() + `","tech_id":"` + tid.String() + `","role_id":"` + rid.String() + `"}`
	r := gin.New()
	r.POST("/visit-assignments", NewVisitAssignmentHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/visit-assignments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVisitAssignmentHandler_Delete(t *testing.T) {
	assignmentRepo, _, _, svc := newVisitAssignmentSvc(t)
	id := uuid.New()
	assignmentRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitAssignment{ID: id}, nil)
	assignmentRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/visit-assignments/:id", NewVisitAssignmentHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/visit-assignments/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newVisitSvc(t *testing.T) (*mocks.VisitRepository, *operationssvc.VisitService) {
	t.Helper()
	visitRepo := mocks.NewVisitRepository(t)
	assignmentRepo := mocks.NewVisitAssignmentRepository(t)
	techRepo := mocks.NewTechnicianRepository(t)
	checkpointRepo := mocks.NewVisitCheckpointRepository(t)
	ctmRepo := mocks.NewChecklistTemplateMaterialRepository(t)
	cttRepo := mocks.NewChecklistTemplateToolRepository(t)
	ctaRepo := mocks.NewChecklistTemplateEPPRepository(t)
	vcmRepo := mocks.NewVisitChecklistMaterialRepository(t)
	vctRepo := mocks.NewVisitChecklistToolRepository(t)
	vceRepo := mocks.NewVisitChecklistEPPRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	auditRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.PlanAuditLog")).Return(nil).Maybe()
	audit := coreservice.NewAuditService(auditRepo)
	svc := operationssvc.NewVisitService(operationssvc.VisitServiceDeps{
		VisitRepo:      visitRepo,
		AssignmentRepo: assignmentRepo,
		TechnicianRepo: techRepo,
		CheckpointRepo: checkpointRepo,
		CTMRepo:        ctmRepo,
		CTTRepo:        cttRepo,
		CTARepo:        ctaRepo,
		VCMRepo:        vcmRepo,
		VCTRepo:        vctRepo,
		VCERepo:        vceRepo,
		AuditService:   audit,
	})
	return visitRepo, svc
}

func TestVisitHandler_GetByID(t *testing.T) {
	visitRepo, svc := newVisitSvc(t)
	id := uuid.New()
	visitRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.Visit{ID: id}, nil)

	r := gin.New()
	r.GET("/visits/:id", NewVisitHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visits/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitHandler_List(t *testing.T) {
	visitRepo, svc := newVisitSvc(t)
	visitRepo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[operations.Visit]{Total: 0}, nil)

	r := gin.New()
	r.GET("/visits", NewVisitHandler(svc).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visits", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitHandler_Create(t *testing.T) {
	visitRepo, svc := newVisitSvc(t)
	visitRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.Visit")).Return(nil)

	body := `{"type":"` + uuid.New().String() + `","status":"scheduled","priority":"normal","billing_to":"partner","scheduled_at":"2026-01-01T00:00:00Z"}`
	r := gin.New()
	r.POST("/visits", NewVisitHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/visits", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVisitHandler_Update_UnsupportedStatus(t *testing.T) {
	_, svc := newVisitSvc(t)
	id := uuid.New()
	body := `{"type":"` + uuid.New().String() + `","status":"rescheduled","priority":"normal","billing_to":"partner","scheduled_at":"2026-01-01T00:00:00Z"}`
	r := gin.New()
	r.PUT("/visits/:id", NewVisitHandler(svc).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/visits/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitHandler_Delete(t *testing.T) {
	visitRepo, svc := newVisitSvc(t)
	id := uuid.New()
	visitRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.Visit{ID: id}, nil)
	visitRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/visits/:id", NewVisitHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/visits/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newReportSvc(t *testing.T) (*mocks.VisitReportRepository, *mocks.ReportEntryRepository, *mocks.ReportImageRepository, *operationssvc.ReportService) {
	t.Helper()
	visitReportRepo := mocks.NewVisitReportRepository(t)
	reportEntryRepo := mocks.NewReportEntryRepository(t)
	reportImageRepo := mocks.NewReportImageRepository(t)
	visitPhotoRepo := mocks.NewVisitPhotoRepository(t)
	visitMeasureRepo := mocks.NewVisitMeasurementRepository(t)
	reportTmplRepo := mocks.NewReportTemplateRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	auditRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.PlanAuditLog")).Return(nil).Maybe()
	audit := coreservice.NewAuditService(auditRepo)

	svc := operationssvc.NewReportService(visitReportRepo, reportEntryRepo, reportImageRepo, visitPhotoRepo, visitMeasureRepo, reportTmplRepo, audit)
	return visitReportRepo, reportEntryRepo, reportImageRepo, svc
}

func TestVisitReportHandler_GetByID(t *testing.T) {
	reportRepo, _, _, svc := newReportSvc(t)
	id := uuid.New()
	reportRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitReport{ID: id}, nil)

	r := gin.New()
	r.GET("/visit-reports/:id", NewVisitReportHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visit-reports/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitReportHandler_List(t *testing.T) {
	reportRepo, _, _, svc := newReportSvc(t)
	reportRepo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[operations.VisitReport]{Total: 0}, nil)

	r := gin.New()
	r.GET("/visit-reports", NewVisitReportHandler(svc).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visit-reports", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitReportHandler_ListByVisit(t *testing.T) {
	reportRepo, _, _, svc := newReportSvc(t)
	vid := uuid.New()
	reportRepo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitReport{}, nil)

	r := gin.New()
	r.GET("/visits/:id/reports", NewVisitReportHandler(svc).ListByVisit)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/reports", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitReportHandler_Create(t *testing.T) {
	reportRepo, entryRepo, _, svc := newReportSvc(t)
	vid := uuid.New()
	tid := uuid.New()
	reportRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitReport")).Return(nil)
	entryRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.ReportEntry")).Return(nil)

	body := `{"visit_id":"` + vid.String() + `","report_template_id":"` + tid.String() + `","source":"system","recorded_at":"2026-01-01T00:00:00Z"}`
	r := gin.New()
	r.POST("/visit-reports", NewVisitReportHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/visit-reports", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVisitReportHandler_Create_MissingTemplate(t *testing.T) {
	_, _, _, svc := newReportSvc(t)
	vid := uuid.New()
	body := `{"visit_id":"` + vid.String() + `"}`
	r := gin.New()
	r.POST("/visit-reports", NewVisitReportHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/visit-reports", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitReportHandler_Create_WithSourceAndRecordedAt(t *testing.T) {
	reportRepo, entryRepo, _, svc := newReportSvc(t)
	vid := uuid.New()
	tid := uuid.New()
	reportRepo.EXPECT().Create(mockCtx(), mock.MatchedBy(func(r *operations.VisitReport) bool {
		return r.Source == shared.ReportSourcePaper && r.RecordedAt.Year() == 2026
	})).Return(nil)
	entryRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.ReportEntry")).Return(nil)

	body := `{"visit_id":"` + vid.String() + `","report_template_id":"` + tid.String() + `","source":"paper","recorded_at":"2026-01-15T10:30:00Z"}`
	r := gin.New()
	r.POST("/visit-reports", NewVisitReportHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/visit-reports", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVisitReportHandler_Delete(t *testing.T) {
	reportRepo, entryRepo, imageRepo, svc := newReportSvc(t)
	id := uuid.New()
	reportRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitReport{ID: id}, nil)
	entryRepo.EXPECT().DeleteByReport(mockCtx(), id).Return(nil)
	imageRepo.EXPECT().DeleteByReport(mockCtx(), id).Return(nil)
	reportRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/visit-reports/:id", NewVisitReportHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/visit-reports/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newSLASvc(t *testing.T) (*mocks.VisitSLATrackingRepository, *mocks.VisitRepository, *operationssvc.SLATrackingService) {
	t.Helper()
	slaTrackingRepo := mocks.NewVisitSLATrackingRepository(t)
	slaRepo := mocks.NewSLARepository(t)
	visitRepo := mocks.NewVisitRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	auditRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.PlanAuditLog")).Return(nil).Maybe()
	notifRepo := mocks.NewNotificationRepository(t)
	audit := coreservice.NewAuditService(auditRepo)
	notifSvc := notifservice.NewNotificationService(notifRepo)
	svc := operationssvc.NewSLATrackingService(slaTrackingRepo, slaRepo, visitRepo, audit, notifSvc)
	return slaTrackingRepo, visitRepo, svc
}

func TestVisitSLATrackingHandler_GetByID(t *testing.T) {
	slaRepo, _, svc := newSLASvc(t)
	id := uuid.New()
	slaRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitSLATracking{ID: id}, nil)

	r := gin.New()
	r.GET("/visit-sla-trackings/:id", NewVisitSLATrackingHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visit-sla-trackings/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitSLATrackingHandler_List(t *testing.T) {
	slaRepo, _, svc := newSLASvc(t)
	slaRepo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[operations.VisitSLATracking]{Total: 0}, nil)

	r := gin.New()
	r.GET("/visit-sla-trackings", NewVisitSLATrackingHandler(svc).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visit-sla-trackings", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitSLATrackingHandler_ListByVisit(t *testing.T) {
	slaRepo, _, svc := newSLASvc(t)
	vid := uuid.New()
	slaRepo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitSLATracking{}, nil)

	r := gin.New()
	r.GET("/visits/:id/sla-trackings", NewVisitSLATrackingHandler(svc).ListByVisit)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/sla-trackings", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestVisitSLATrackingHandler_Create(t *testing.T) {
	slaRepo, visitRepo, svc := newSLASvc(t)
	vid := uuid.New()
	sid := uuid.New()
	visitRepo.EXPECT().GetByID(mockCtx(), vid).Return(&operations.Visit{ID: vid}, nil)
	slaRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitSLATracking")).Return(nil)
	body := `{"visit_id":"` + vid.String() + `","sla_id":"` + sid.String() + `"}`
	r := gin.New()
	r.POST("/visit-sla-trackings", NewVisitSLATrackingHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/visit-sla-trackings", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVisitSLATrackingHandler_Update(t *testing.T) {
	slaRepo, _, svc := newSLASvc(t)
	id := uuid.New()
	slaRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitSLATracking{ID: id}, nil)
	slaRepo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*operations.VisitSLATracking")).Return(nil)

	vid := uuid.New()
	sid := uuid.New()
	body := `{"visit_id":"` + vid.String() + `","sla_id":"` + sid.String() + `"}`
	r := gin.New()
	r.PUT("/visit-sla-trackings/:id", NewVisitSLATrackingHandler(svc).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/visit-sla-trackings/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestVisitSLATrackingHandler_Delete(t *testing.T) {
	slaRepo, _, svc := newSLASvc(t)
	id := uuid.New()
	slaRepo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitSLATracking{ID: id}, nil)
	slaRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/visit-sla-trackings/:id", NewVisitSLATrackingHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/visit-sla-trackings/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

var _ = shared.VisitStatus("scheduled")
var _ = shared.VisitPriority("normal")
var _ = shared.VehicleStatusAvailable
var _ = &partners.SLA{}
