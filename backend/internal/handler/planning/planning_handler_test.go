package planninghandler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	inventorymodel "localis-backend/internal/model/inventory"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	coreservice "localis-backend/internal/service/core"
	planningservice "localis-backend/internal/service/planning"
	"localis-backend/internal/testutil/mocks"
)

func mockCtx() context.Context { return context.Background() }

func newDailyLoadSvc(t *testing.T) (*mocks.DailyLoadMaterialRepository, *mocks.DailyLoadToolRepository, *mocks.DailyLoadEPPRepository, *mocks.MaterialRepository, *mocks.ToolRepository, *mocks.EPPItemRepository, *planningservice.DailyLoadService) {
	t.Helper()
	loadMatRepo := mocks.NewDailyLoadMaterialRepository(t)
	loadToolRepo := mocks.NewDailyLoadToolRepository(t)
	loadEPPRepo := mocks.NewDailyLoadEPPRepository(t)
	matRepo := mocks.NewMaterialRepository(t)
	toolRepo := mocks.NewToolRepository(t)
	eppRepo := mocks.NewEPPItemRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	auditRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.PlanAuditLog")).Return(nil).Maybe()
	audit := coreservice.NewAuditService(auditRepo)
	svc := planningservice.NewDailyLoadService(loadMatRepo, loadToolRepo, loadEPPRepo, matRepo, toolRepo, eppRepo, audit)
	return loadMatRepo, loadToolRepo, loadEPPRepo, matRepo, toolRepo, eppRepo, svc
}

func TestDailyLoadMaterialHandler_GetByID(t *testing.T) {
	loadMatRepo, _, _, _, _, _, svc := newDailyLoadSvc(t)
	id := uuid.New()
	loadMatRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyLoadMaterial{ID: id}, nil)

	r := gin.New()
	r.GET("/daily-load-materials/:id", NewDailyLoadMaterialHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-load-materials/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyLoadMaterialHandler_ListByPlan(t *testing.T) {
	loadMatRepo, _, _, _, _, _, svc := newDailyLoadSvc(t)
	pid := uuid.New()
	loadMatRepo.EXPECT().ListByPlan(mockCtx(), pid).Return([]planning.DailyLoadMaterial{}, nil)

	r := gin.New()
	r.GET("/daily-plans/:id/materials", NewDailyLoadMaterialHandler(svc).ListByPlan)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-plans/"+pid.String()+"/materials", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyLoadMaterialHandler_Create(t *testing.T) {
	loadMatRepo, _, _, matRepo, _, _, svc := newDailyLoadSvc(t)
	mid := uuid.New()
	pid := uuid.New()
	matRepo.EXPECT().GetByID(mockCtx(), mid).Return(&inventorymodel.Material{ID: mid}, nil)
	loadMatRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*planning.DailyLoadMaterial")).Return(nil)

	body := `{"daily_plan_id":"` + pid.String() + `","material_id":"` + mid.String() + `","loaded_quantity":5,"returned_quantity":1}`
	r := gin.New()
	r.POST("/daily-load-materials", NewDailyLoadMaterialHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/daily-load-materials", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyLoadMaterialHandler_Update(t *testing.T) {
	loadMatRepo, _, _, _, _, _, svc := newDailyLoadSvc(t)
	id := uuid.New()
	loadMatRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyLoadMaterial{ID: id}, nil)
	loadMatRepo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*planning.DailyLoadMaterial")).Return(nil)

	mid := uuid.New()
	pid := uuid.New()
	body := `{"daily_plan_id":"` + pid.String() + `","material_id":"` + mid.String() + `","loaded_quantity":5,"returned_quantity":1}`
	r := gin.New()
	r.PUT("/daily-load-materials/:id", NewDailyLoadMaterialHandler(svc).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/daily-load-materials/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyLoadMaterialHandler_Delete(t *testing.T) {
	loadMatRepo, _, _, _, _, _, svc := newDailyLoadSvc(t)
	id := uuid.New()
	loadMatRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyLoadMaterial{ID: id}, nil)
	loadMatRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/daily-load-materials/:id", NewDailyLoadMaterialHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/daily-load-materials/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyLoadToolHandler_GetByID(t *testing.T) {
	_, loadToolRepo, _, _, _, _, svc := newDailyLoadSvc(t)
	id := uuid.New()
	loadToolRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyLoadTool{ID: id}, nil)

	r := gin.New()
	r.GET("/daily-load-tools/:id", NewDailyLoadToolHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-load-tools/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyLoadToolHandler_ListByPlan(t *testing.T) {
	_, loadToolRepo, _, _, _, _, svc := newDailyLoadSvc(t)
	pid := uuid.New()
	loadToolRepo.EXPECT().ListByPlan(mockCtx(), pid).Return([]planning.DailyLoadTool{}, nil)

	r := gin.New()
	r.GET("/daily-plans/:id/tools", NewDailyLoadToolHandler(svc).ListByPlan)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-plans/"+pid.String()+"/tools", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyLoadToolHandler_Create(t *testing.T) {
	_, loadToolRepo, _, _, toolRepo, _, svc := newDailyLoadSvc(t)
	tid := uuid.New()
	pid := uuid.New()
	toolRepo.EXPECT().GetByID(mockCtx(), tid).Return(&inventorymodel.Tool{ID: tid}, nil)
	loadToolRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*planning.DailyLoadTool")).Return(nil)

	body := `{"daily_plan_id":"` + pid.String() + `","tool_id":"` + tid.String() + `","loaded_quantity":2,"returned_quantity":1}`
	r := gin.New()
	r.POST("/daily-load-tools", NewDailyLoadToolHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/daily-load-tools", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyLoadToolHandler_Update(t *testing.T) {
	_, loadToolRepo, _, _, _, _, svc := newDailyLoadSvc(t)
	id := uuid.New()
	loadToolRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyLoadTool{ID: id}, nil)
	loadToolRepo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*planning.DailyLoadTool")).Return(nil)

	tid := uuid.New()
	pid := uuid.New()
	body := `{"daily_plan_id":"` + pid.String() + `","tool_id":"` + tid.String() + `","loaded_quantity":2,"returned_quantity":1}`
	r := gin.New()
	r.PUT("/daily-load-tools/:id", NewDailyLoadToolHandler(svc).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/daily-load-tools/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyLoadToolHandler_Delete(t *testing.T) {
	_, loadToolRepo, _, _, _, _, svc := newDailyLoadSvc(t)
	id := uuid.New()
	loadToolRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyLoadTool{ID: id}, nil)
	loadToolRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/daily-load-tools/:id", NewDailyLoadToolHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/daily-load-tools/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyLoadEPPHandler_GetByID(t *testing.T) {
	_, _, loadEPPRepo, _, _, _, svc := newDailyLoadSvc(t)
	id := uuid.New()
	loadEPPRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyLoadEPP{ID: id}, nil)

	r := gin.New()
	r.GET("/daily-load-epps/:id", NewDailyLoadEPPHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-load-epps/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyLoadEPPHandler_ListByPlan(t *testing.T) {
	_, _, loadEPPRepo, _, _, _, svc := newDailyLoadSvc(t)
	pid := uuid.New()
	loadEPPRepo.EXPECT().ListByPlan(mockCtx(), pid).Return([]planning.DailyLoadEPP{}, nil)

	r := gin.New()
	r.GET("/daily-plans/:id/epps", NewDailyLoadEPPHandler(svc).ListByPlan)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-plans/"+pid.String()+"/epps", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyLoadEPPHandler_Create(t *testing.T) {
	_, _, loadEPPRepo, _, _, eppRepo, svc := newDailyLoadSvc(t)
	eid := uuid.New()
	pid := uuid.New()
	eppRepo.EXPECT().GetByID(mockCtx(), eid).Return(&inventorymodel.EPPItem{ID: eid}, nil)
	loadEPPRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*planning.DailyLoadEPP")).Return(nil)

	body := `{"daily_plan_id":"` + pid.String() + `","epp_id":"` + eid.String() + `","loaded_quantity":4,"returned_quantity":1}`
	r := gin.New()
	r.POST("/daily-load-epps", NewDailyLoadEPPHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/daily-load-epps", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyLoadEPPHandler_Update(t *testing.T) {
	_, _, loadEPPRepo, _, _, _, svc := newDailyLoadSvc(t)
	id := uuid.New()
	loadEPPRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyLoadEPP{ID: id}, nil)
	loadEPPRepo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*planning.DailyLoadEPP")).Return(nil)

	eid := uuid.New()
	pid := uuid.New()
	body := `{"daily_plan_id":"` + pid.String() + `","epp_id":"` + eid.String() + `","loaded_quantity":4,"returned_quantity":1}`
	r := gin.New()
	r.PUT("/daily-load-epps/:id", NewDailyLoadEPPHandler(svc).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/daily-load-epps/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyLoadEPPHandler_Delete(t *testing.T) {
	_, _, loadEPPRepo, _, _, _, svc := newDailyLoadSvc(t)
	id := uuid.New()
	loadEPPRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyLoadEPP{ID: id}, nil)
	loadEPPRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/daily-load-epps/:id", NewDailyLoadEPPHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/daily-load-epps/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newDailyPlanSvc(t *testing.T) (*mocks.DailyPlanRepository, *mocks.TechnicianRepository, *planningservice.DailyPlanService) {
	t.Helper()
	planRepo := mocks.NewDailyPlanRepository(t)
	assignmentRepo := mocks.NewDailyPlanAssignmentRepository(t)
	techRepo := mocks.NewTechnicianRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	auditRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.PlanAuditLog")).Return(nil).Maybe()
	audit := coreservice.NewAuditService(auditRepo)
	svc := planningservice.NewDailyPlanService(planRepo, assignmentRepo, techRepo, audit)
	return planRepo, techRepo, svc
}

func TestDailyPlanHandler_GetByID(t *testing.T) {
	planRepo, _, svc := newDailyPlanSvc(t)
	id := uuid.New()
	planRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyPlan{ID: id}, nil)

	r := gin.New()
	r.GET("/daily-plans/:id", NewDailyPlanHandler(svc).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-plans/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyPlanHandler_List(t *testing.T) {
	planRepo, _, svc := newDailyPlanSvc(t)
	planRepo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[planning.DailyPlan]{Total: 0}, nil)

	r := gin.New()
	r.GET("/daily-plans", NewDailyPlanHandler(svc).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-plans", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyPlanHandler_Create(t *testing.T) {
	planRepo, _, svc := newDailyPlanSvc(t)
	planRepo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*planning.DailyPlan")).Return(nil)

	body := `{"name":"Plan A","date":"2026-01-01T00:00:00Z"}`
	r := gin.New()
	r.POST("/daily-plans", NewDailyPlanHandler(svc).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/daily-plans", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyPlanHandler_Update(t *testing.T) {
	planRepo, _, svc := newDailyPlanSvc(t)
	id := uuid.New()
	planRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyPlan{ID: id}, nil)
	planRepo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*planning.DailyPlan")).Return(nil)

	body := `{"name":"Plan A","date":"2026-01-01T00:00:00Z"}`
	r := gin.New()
	r.PUT("/daily-plans/:id", NewDailyPlanHandler(svc).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/daily-plans/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyPlanHandler_Delete(t *testing.T) {
	planRepo, _, svc := newDailyPlanSvc(t)
	id := uuid.New()
	planRepo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyPlan{ID: id}, nil)
	planRepo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/daily-plans/:id", NewDailyPlanHandler(svc).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/daily-plans/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyPlanAssignmentHandler_GetByID(t *testing.T) {
	repo := mocks.NewDailyPlanAssignmentRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyPlanAssignment{ID: id}, nil)

	r := gin.New()
	r.GET("/daily-plan-assignments/:id", NewDailyPlanAssignmentHandler(repo).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-plan-assignments/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyPlanAssignmentHandler_ListByPlan(t *testing.T) {
	repo := mocks.NewDailyPlanAssignmentRepository(t)
	pid := uuid.New()
	repo.EXPECT().ListByPlan(mockCtx(), pid).Return([]planning.DailyPlanAssignment{}, nil)

	r := gin.New()
	r.GET("/daily-plans/:id/assignments", NewDailyPlanAssignmentHandler(repo).ListByPlan)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-plans/"+pid.String()+"/assignments", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyPlanAssignmentHandler_Create(t *testing.T) {
	repo := mocks.NewDailyPlanAssignmentRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*planning.DailyPlanAssignment")).Return(nil)

	pid := uuid.New()
	tid := uuid.New()
	rid := uuid.New()
	body := `{"daily_plan_id":"` + pid.String() + `","tech_id":"` + tid.String() + `","role_id":"` + rid.String() + `"}`
	r := gin.New()
	r.POST("/daily-plan-assignments", NewDailyPlanAssignmentHandler(repo).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/daily-plan-assignments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyPlanAssignmentHandler_List(t *testing.T) {
	repo := mocks.NewDailyPlanAssignmentRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[planning.DailyPlanAssignment]{Total: 0}, nil)

	r := gin.New()
	r.GET("/daily-plan-assignments", NewDailyPlanAssignmentHandler(repo).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/daily-plan-assignments", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestDailyPlanAssignmentHandler_Update(t *testing.T) {
	repo := mocks.NewDailyPlanAssignmentRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*planning.DailyPlanAssignment")).Return(nil)
	repo.EXPECT().GetByID(mockCtx(), id).Return(&planning.DailyPlanAssignment{ID: id}, nil)

	pid := uuid.New()
	tid := uuid.New()
	rid := uuid.New()
	body := `{"daily_plan_id":"` + pid.String() + `","tech_id":"` + tid.String() + `","role_id":"` + rid.String() + `"}`
	r := gin.New()
	r.PUT("/daily-plan-assignments/:id", NewDailyPlanAssignmentHandler(repo).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/daily-plan-assignments/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestDailyPlanAssignmentHandler_Delete(t *testing.T) {
	repo := mocks.NewDailyPlanAssignmentRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/daily-plan-assignments/:id", NewDailyPlanAssignmentHandler(repo).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/daily-plan-assignments/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRouteHandler_GetByID(t *testing.T) {
	repo := mocks.NewRouteRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&planning.Route{ID: id}, nil)

	r := gin.New()
	r.GET("/routes/:id", NewRouteHandler(repo).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/routes/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRouteHandler_List(t *testing.T) {
	repo := mocks.NewRouteRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[planning.Route]{Total: 0}, nil)

	r := gin.New()
	r.GET("/routes", NewRouteHandler(repo).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/routes", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRouteHandler_Create(t *testing.T) {
	repo := mocks.NewRouteRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*planning.Route")).Return(nil)

	body := `{"type":"gas","date":"2026-01-01T00:00:00Z","status":"scheduled"}`
	r := gin.New()
	r.POST("/routes", NewRouteHandler(repo).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/routes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRouteHandler_Update(t *testing.T) {
	repo := mocks.NewRouteRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*planning.Route")).Return(nil)

	body := `{"type":"gas","date":"2026-01-01T00:00:00Z","status":"scheduled"}`
	r := gin.New()
	r.PUT("/routes/:id", NewRouteHandler(repo).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/routes/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRouteHandler_Delete(t *testing.T) {
	repo := mocks.NewRouteRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/routes/:id", NewRouteHandler(repo).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/routes/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

var _ = shared.RouteStatus("scheduled")
