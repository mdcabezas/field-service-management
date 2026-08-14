package inventoryhandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	inventorymodel "localis-backend/internal/model/inventory"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	coreservice "localis-backend/internal/service/core"
	invcsv "localis-backend/internal/service/inventory"
	notifservice "localis-backend/internal/service/notifications"
	"localis-backend/internal/testutil/mocks"
)

func newMaintSvc(t *testing.T) (*mocks.MaintenanceScheduleRepository, *mocks.MaintenanceRecordRepository, *invcsv.MaintenanceService) {
	t.Helper()
	scheduleRepo := mocks.NewMaintenanceScheduleRepository(t)
	recordRepo := mocks.NewMaintenanceRecordRepository(t)
	vehicleRepo := mocks.NewVehicleRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	notifRepo := mocks.NewNotificationRepository(t)
	audit := coreservice.NewAuditService(auditRepo)
	notifSvc := notifservice.NewNotificationService(notifRepo)
	svc := invcsv.NewMaintenanceService(scheduleRepo, recordRepo, vehicleRepo, audit, notifSvc)
	auditRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.PlanAuditLog")).Return(nil).Maybe()
	return scheduleRepo, recordRepo, svc
}

func TestMaintenanceRecordHandler_GetByID(t *testing.T) {
	_, recordRepo, svc := newMaintSvc(t)
	id := uuid.New()
	recordRepo.EXPECT().GetByID(mockCtx(), id).Return(&inventorymodel.MaintenanceRecord{ID: id}, nil)

	r := gin.New()
	r.GET("/maintenance-records/:id", NewMaintenanceRecordHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/maintenance-records/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMaintenanceRecordHandler_List(t *testing.T) {
	_, recordRepo, svc := newMaintSvc(t)
	recordRepo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[inventorymodel.MaintenanceRecord]{Total: 0}, nil)

	r := gin.New()
	r.GET("/maintenance-records", NewMaintenanceRecordHandler(svc).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/maintenance-records", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMaintenanceRecordHandler_Create(t *testing.T) {
	scheduleRepo, recordRepo, svc := newMaintSvc(t)
	refID := uuid.New()
	scheduleRepo.EXPECT().ListActive(mockCtx()).Return([]inventorymodel.MaintenanceSchedule{}, nil)
	recordRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*inventory.MaintenanceRecord")).Return(nil)

	body := `{"type":"vehicle","reference_id":"` + refID.String() + `","date":"2026-01-01T00:00:00Z"}`
	r := gin.New()
	r.POST("/maintenance-records", NewMaintenanceRecordHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/maintenance-records", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestMaintenanceRecordHandler_Update(t *testing.T) {
	_, recordRepo, svc := newMaintSvc(t)
	id := uuid.New()
	recordRepo.EXPECT().GetByID(mockCtx(), id).Return(&inventorymodel.MaintenanceRecord{ID: id}, nil)
	recordRepo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*inventory.MaintenanceRecord")).Return(nil)

	refID := uuid.New()
	body := `{"type":"vehicle","reference_id":"` + refID.String() + `","date":"2026-01-01T00:00:00Z"}`
	r := gin.New()
	r.PUT("/maintenance-records/:id", NewMaintenanceRecordHandler(svc).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/maintenance-records/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestMaintenanceRecordHandler_Delete(t *testing.T) {
	_, recordRepo, svc := newMaintSvc(t)
	id := uuid.New()
	recordRepo.EXPECT().GetByID(mockCtx(), id).Return(&inventorymodel.MaintenanceRecord{ID: id}, nil)
	recordRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/maintenance-records/:id", NewMaintenanceRecordHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/maintenance-records/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMaintenanceScheduleHandler_GetByID(t *testing.T) {
	scheduleRepo, _, svc := newMaintSvc(t)
	id := uuid.New()
	scheduleRepo.EXPECT().GetByID(mockCtx(), id).Return(&inventorymodel.MaintenanceSchedule{ID: id}, nil)

	r := gin.New()
	r.GET("/maintenance-schedules/:id", NewMaintenanceScheduleHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/maintenance-schedules/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMaintenanceScheduleHandler_List(t *testing.T) {
	scheduleRepo, _, svc := newMaintSvc(t)
	scheduleRepo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[inventorymodel.MaintenanceSchedule]{Total: 0}, nil)

	r := gin.New()
	r.GET("/maintenance-schedules", NewMaintenanceScheduleHandler(svc).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/maintenance-schedules", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMaintenanceScheduleHandler_Create(t *testing.T) {
	scheduleRepo, _, svc := newMaintSvc(t)
	refID := uuid.New()
	scheduleRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*inventory.MaintenanceSchedule")).Return(nil)

	body := `{"type":"vehicle","reference_id":"` + refID.String() + `","active":true}`
	r := gin.New()
	r.POST("/maintenance-schedules", NewMaintenanceScheduleHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/maintenance-schedules", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestMaintenanceScheduleHandler_Update(t *testing.T) {
	scheduleRepo, _, svc := newMaintSvc(t)
	id := uuid.New()
	refID := uuid.New()
	scheduleRepo.EXPECT().GetByID(mockCtx(), id).Return(&inventorymodel.MaintenanceSchedule{ID: id}, nil)
	scheduleRepo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*inventory.MaintenanceSchedule")).Return(nil)

	body := `{"type":"vehicle","reference_id":"` + refID.String() + `","active":true}`
	r := gin.New()
	r.PUT("/maintenance-schedules/:id", NewMaintenanceScheduleHandler(svc).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/maintenance-schedules/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestMaintenanceScheduleHandler_Delete(t *testing.T) {
	scheduleRepo, _, svc := newMaintSvc(t)
	id := uuid.New()
	scheduleRepo.EXPECT().GetByID(mockCtx(), id).Return(&inventorymodel.MaintenanceSchedule{ID: id}, nil)
	scheduleRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/maintenance-schedules/:id", NewMaintenanceScheduleHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/maintenance-schedules/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

var _ = shared.MaintenanceTypeVehicle
