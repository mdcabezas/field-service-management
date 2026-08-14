package inventoryhandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/testutil/mocks"
)

func TestChecklistTemplateHandler_GetByID(t *testing.T) {
	repo := mocks.NewChecklistTemplateRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&inventory.ChecklistTemplate{ID: id}, nil)

	r := gin.New()
	r.GET("/checklist-templates/:id", NewChecklistTemplateHandler(repo).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/checklist-templates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestChecklistTemplateHandler_List(t *testing.T) {
	repo := mocks.NewChecklistTemplateRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[inventory.ChecklistTemplate]{Total: 0}, nil)

	r := gin.New()
	r.GET("/checklist-templates", NewChecklistTemplateHandler(repo).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/checklist-templates", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestChecklistTemplateHandler_Create(t *testing.T) {
	repo := mocks.NewChecklistTemplateRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*inventory.ChecklistTemplate")).Return(nil)

	r := gin.New()
	r.POST("/checklist-templates", NewChecklistTemplateHandler(repo).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/checklist-templates", bytes.NewBufferString(`{"name":"Checklist"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestChecklistTemplateHandler_Update(t *testing.T) {
	repo := mocks.NewChecklistTemplateRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*inventory.ChecklistTemplate")).Return(nil)

	r := gin.New()
	r.PUT("/checklist-templates/:id", NewChecklistTemplateHandler(repo).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/checklist-templates/"+id.String(), bytes.NewBufferString(`{"name":"Checklist 2"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestChecklistTemplateHandler_Delete(t *testing.T) {
	repo := mocks.NewChecklistTemplateRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/checklist-templates/:id", NewChecklistTemplateHandler(repo).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/checklist-templates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCostRateHandler_GetByID(t *testing.T) {
	repo := mocks.NewCostRateRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&inventory.CostRate{ID: id}, nil)

	r := gin.New()
	r.GET("/cost-rates/:id", NewCostRateHandler(repo).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cost-rates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCostRateHandler_List(t *testing.T) {
	repo := mocks.NewCostRateRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[inventory.CostRate]{Total: 0}, nil)

	r := gin.New()
	r.GET("/cost-rates", NewCostRateHandler(repo).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/cost-rates", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCostRateHandler_Create(t *testing.T) {
	repo := mocks.NewCostRateRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*inventory.CostRate")).Return(nil)

	r := gin.New()
	r.POST("/cost-rates", NewCostRateHandler(repo).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/cost-rates", bytes.NewBufferString(`{"type":"labor","value":100,"unit":"hour","valid_from":"2026-01-01T00:00:00Z"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCostRateHandler_Update(t *testing.T) {
	repo := mocks.NewCostRateRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*inventory.CostRate")).Return(nil)

	r := gin.New()
	r.PUT("/cost-rates/:id", NewCostRateHandler(repo).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/cost-rates/"+id.String(), bytes.NewBufferString(`{"type":"labor","value":120,"unit":"hour","valid_from":"2026-01-01T00:00:00Z"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCostRateHandler_Delete(t *testing.T) {
	repo := mocks.NewCostRateRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/cost-rates/:id", NewCostRateHandler(repo).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/cost-rates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestEPPItemHandler_GetByID(t *testing.T) {
	repo := mocks.NewEPPItemRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&inventory.EPPItem{ID: id}, nil)

	r := gin.New()
	r.GET("/epp-items/:id", NewEPPItemHandler(repo).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/epp-items/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestEPPItemHandler_List(t *testing.T) {
	repo := mocks.NewEPPItemRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[inventory.EPPItem]{Total: 0}, nil)

	r := gin.New()
	r.GET("/epp-items", NewEPPItemHandler(repo).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/epp-items", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestEPPItemHandler_Create(t *testing.T) {
	repo := mocks.NewEPPItemRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*inventory.EPPItem")).Return(nil)

	r := gin.New()
	r.POST("/epp-items", NewEPPItemHandler(repo).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/epp-items", bytes.NewBufferString(`{"name":"Casco","type":"head","lifecycle":"reusable"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestEPPItemHandler_Update(t *testing.T) {
	repo := mocks.NewEPPItemRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*inventory.EPPItem")).Return(nil)

	r := gin.New()
	r.PUT("/epp-items/:id", NewEPPItemHandler(repo).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/epp-items/"+id.String(), bytes.NewBufferString(`{"name":"Casco","type":"head","lifecycle":"reusable"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestEPPItemHandler_Delete(t *testing.T) {
	repo := mocks.NewEPPItemRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/epp-items/:id", NewEPPItemHandler(repo).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/epp-items/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMaterialHandler_GetByID(t *testing.T) {
	repo := mocks.NewMaterialRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&inventory.Material{ID: id}, nil)

	r := gin.New()
	r.GET("/materials/:id", NewMaterialHandler(repo).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/materials/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMaterialHandler_List(t *testing.T) {
	repo := mocks.NewMaterialRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[inventory.Material]{Total: 0}, nil)

	r := gin.New()
	r.GET("/materials", NewMaterialHandler(repo).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/materials", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMaterialHandler_Create(t *testing.T) {
	repo := mocks.NewMaterialRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*inventory.Material")).Return(nil)

	r := gin.New()
	r.POST("/materials", NewMaterialHandler(repo).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/materials", bytes.NewBufferString(`{"name":"Cable","unit":"m"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestMaterialHandler_Update(t *testing.T) {
	repo := mocks.NewMaterialRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*inventory.Material")).Return(nil)

	r := gin.New()
	r.PUT("/materials/:id", NewMaterialHandler(repo).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/materials/"+id.String(), bytes.NewBufferString(`{"name":"Cable","unit":"m"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMaterialHandler_Delete(t *testing.T) {
	repo := mocks.NewMaterialRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/materials/:id", NewMaterialHandler(repo).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/materials/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRentalHandler_GetByID(t *testing.T) {
	repo := mocks.NewRentalRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&inventory.Rental{ID: id}, nil)

	r := gin.New()
	r.GET("/rentals/:id", NewRentalHandler(repo).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/rentals/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRentalHandler_List(t *testing.T) {
	repo := mocks.NewRentalRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[inventory.Rental]{Total: 0}, nil)

	r := gin.New()
	r.GET("/rentals", NewRentalHandler(repo).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/rentals", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRentalHandler_Create(t *testing.T) {
	repo := mocks.NewRentalRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*inventory.Rental")).Return(nil)

	r := gin.New()
	r.POST("/rentals", NewRentalHandler(repo).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/rentals", bytes.NewBufferString(`{"type":"vehicle","item_description":"Retroexcavadora"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRentalHandler_Update(t *testing.T) {
	repo := mocks.NewRentalRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*inventory.Rental")).Return(nil)

	r := gin.New()
	r.PUT("/rentals/:id", NewRentalHandler(repo).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/rentals/"+id.String(), bytes.NewBufferString(`{"type":"vehicle","item_description":"Retroexcavadora"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRentalHandler_Delete(t *testing.T) {
	repo := mocks.NewRentalRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/rentals/:id", NewRentalHandler(repo).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/rentals/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestToolHandler_GetByID(t *testing.T) {
	repo := mocks.NewToolRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&inventory.Tool{ID: id}, nil)

	r := gin.New()
	r.GET("/tools/:id", NewToolHandler(repo).GetByID)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tools/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestToolHandler_List(t *testing.T) {
	repo := mocks.NewToolRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[inventory.Tool]{Total: 0}, nil)

	r := gin.New()
	r.GET("/tools", NewToolHandler(repo).List)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tools", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestToolHandler_Create(t *testing.T) {
	repo := mocks.NewToolRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*inventory.Tool")).Return(nil)

	r := gin.New()
	r.POST("/tools", NewToolHandler(repo).Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tools", bytes.NewBufferString(`{"name":"Taladro","status":"available"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestToolHandler_Update(t *testing.T) {
	repo := mocks.NewToolRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*inventory.Tool")).Return(nil)

	r := gin.New()
	r.PUT("/tools/:id", NewToolHandler(repo).Update)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/tools/"+id.String(), bytes.NewBufferString(`{"name":"Taladro","status":"available"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestToolHandler_Delete(t *testing.T) {
	repo := mocks.NewToolRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := gin.New()
	r.DELETE("/tools/:id", NewToolHandler(repo).Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/tools/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

var _ = shared.CostType("labor")
var _ = shared.CostUnit("hour")
var _ = shared.EPPType("head")
var _ = shared.EPPLifecycle("reusable")
var _ = shared.RentalType("vehicle")
var _ = shared.ToolStatus("available")
