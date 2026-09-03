package operationshandler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestChecklistStageHandler_ListByTemplate(t *testing.T) {
	id := uuid.New()
	r := gin.New()
	r.GET("/api/operations/checklist-templates/:id/stages", NewChecklistStageHandler().ListByTemplate)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/operations/checklist-templates/"+id.String()+"/stages", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestChecklistStageHandler_ListByTemplate_InvalidID(t *testing.T) {
	r := gin.New()
	r.GET("/api/operations/checklist-templates/:id/stages", NewChecklistStageHandler().ListByTemplate)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/operations/checklist-templates/invalid/stages", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestChecklistStageHandler_Create(t *testing.T) {
	h := NewChecklistStageHandler()
	templateID := uuid.New()
	body := `{"template_id":"` + templateID.String() + `","name":"Electrical","sort_order":1,"required":true}`
	r := gin.New()
	r.POST("/api/operations/checklist-stages", h.Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/operations/checklist-stages", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestChecklistStageHandler_Create_InvalidBody(t *testing.T) {
	r := gin.New()
	r.POST("/api/operations/checklist-stages", NewChecklistStageHandler().Create)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/operations/checklist-stages", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestChecklistStageHandler_Delete(t *testing.T) {
	h := NewChecklistStageHandler()
	templateID := uuid.New()
	createBody := `{"template_id":"` + templateID.String() + `","name":"Test Stage","sort_order":1}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/operations/checklist-stages", bytes.NewBufferString(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/operations/checklist-stages", h.Create)
	r.ServeHTTP(createRec, createReq)

	stageID := uuid.New()
	r.DELETE("/api/operations/checklist-stages/:id", h.Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/operations/checklist-stages/"+stageID.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestChecklistStageHandler_Delete_InvalidID(t *testing.T) {
	r := gin.New()
	r.DELETE("/api/operations/checklist-stages/:id", NewChecklistStageHandler().Delete)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/operations/checklist-stages/invalid", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}
