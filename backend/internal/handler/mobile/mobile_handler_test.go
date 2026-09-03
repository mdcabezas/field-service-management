package mobilehandler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/auth"
)

func mockCtx() context.Context { return context.Background() }

func TestConfigHandler_GetConfig(t *testing.T) {
	r := gin.New()
	r.GET("/api/mobile/config", NewConfigHandler().GetConfig)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/config", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestConfigHandler_UpdateConfig(t *testing.T) {
	h := NewConfigHandler()
	body := `{"photo_retention_days": 14}`
	r := gin.New()
	r.PUT("/api/mobile/config", h.UpdateConfig)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/mobile/config", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestConfigHandler_UpdateConfig_InvalidBody(t *testing.T) {
	r := gin.New()
	r.PUT("/api/mobile/config", NewConfigHandler().UpdateConfig)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/mobile/config", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestGPSHandler_CreateWaiver(t *testing.T) {
	h := NewGPSHandler()
	visitID := uuid.New()
	body := `{"visit_id":"` + visitID.String() + `","type":"gps_timeout","timestamp":"2026-01-01T00:00:00Z"}`
	r := gin.New()
	r.POST("/api/mobile/gps/waivers", h.CreateWaiver)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/gps/waivers", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGPSHandler_CreateWaiver_InvalidType(t *testing.T) {
	h := NewGPSHandler()
	visitID := uuid.New()
	body := `{"visit_id":"` + visitID.String() + `","type":"invalid_type","timestamp":"2026-01-01T00:00:00Z"}`
	r := gin.New()
	r.POST("/api/mobile/gps/waivers", h.CreateWaiver)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/gps/waivers", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestGPSHandler_CreateWaiver_InvalidBody(t *testing.T) {
	r := gin.New()
	r.POST("/api/mobile/gps/waivers", NewGPSHandler().CreateWaiver)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/gps/waivers", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestGPSHandler_ListWaivers(t *testing.T) {
	h := NewGPSHandler()
	visitID := uuid.New()
	r := gin.New()
	r.GET("/api/mobile/gps/waivers", h.ListWaivers)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/gps/waivers?visit_id="+visitID.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestGPSHandler_ListWaivers_MissingVisitID(t *testing.T) {
	r := gin.New()
	r.GET("/api/mobile/gps/waivers", NewGPSHandler().ListWaivers)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/gps/waivers", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestGPSHandler_DeleteWaiver_NotFound(t *testing.T) {
	id := uuid.New()
	r := gin.New()
	r.DELETE("/api/mobile/gps/waivers/:id", NewGPSHandler().DeleteWaiver)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/mobile/gps/waivers/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestSyncHandler_Push(t *testing.T) {
	r := gin.New()
	r.POST("/api/mobile/sync/push", NewTestSyncHandler().Push)
	w := httptest.NewRecorder()
	body := `{"checkpoints":[],"photos":[],"usages":[],"measurements":[],"reports":[],"status_changes":[]}`
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/push", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSyncHandler_Push_InvalidBody(t *testing.T) {
	r := gin.New()
	r.POST("/api/mobile/sync/push", NewTestSyncHandler().Push)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/push", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestSyncHandler_Pull(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		claims := &auth.Claims{
			UserID: "1001",
			Role:           "admin",
		}
		c.Set(auth.ClaimsKey, claims)
		c.Next()
	})
	r.GET("/api/mobile/sync/pull", NewTestSyncHandler().Pull)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/sync/pull", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestSyncHandler_Retry(t *testing.T) {
	r := gin.New()
	r.POST("/api/mobile/sync/retry", NewTestSyncHandler().Retry)
	w := httptest.NewRecorder()
	body := `{"item_ids":["test-1","test-2"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/retry", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSyncHandler_Retry_InvalidBody(t *testing.T) {
	r := gin.New()
	r.POST("/api/mobile/sync/retry", NewTestSyncHandler().Retry)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/retry", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_GetContacts(t *testing.T) {
	id := uuid.New()
	r := gin.New()
	r.GET("/api/mobile/partners/:id/contacts", NewPartnerHandler().GetContacts)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/partners/"+id.String()+"/contacts", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_GetContacts_InvalidID(t *testing.T) {
	r := gin.New()
	r.GET("/api/mobile/partners/:id/contacts", NewPartnerHandler().GetContacts)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/partners/invalid/contacts", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_GetByProperty(t *testing.T) {
	id := uuid.New()
	r := gin.New()
	r.GET("/api/mobile/properties/:id/partners", NewPartnerHandler().GetByProperty)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/properties/"+id.String()+"/partners", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_GetByProperty_InvalidID(t *testing.T) {
	r := gin.New()
	r.GET("/api/mobile/properties/:id/partners", NewPartnerHandler().GetByProperty)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/properties/invalid/partners", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_ListTemplates(t *testing.T) {
	r := gin.New()
	r.GET("/api/mobile/templates", NewPartnerHandler().ListTemplates)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/templates", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPropertyHistoryHandler_GetByProperty(t *testing.T) {
	id := uuid.New()
	r := gin.New()
	r.GET("/api/mobile/properties/:id/history", NewPropertyHistoryHandler().GetByProperty)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/properties/"+id.String()+"/history", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPropertyHistoryHandler_GetByProperty_InvalidID(t *testing.T) {
	r := gin.New()
	r.GET("/api/mobile/properties/:id/history", NewPropertyHistoryHandler().GetByProperty)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/properties/invalid/history", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAuditHandler_GetByVisit(t *testing.T) {
	id := uuid.New()
	r := gin.New()
	r.GET("/api/mobile/visits/:id/audit", NewAuditHandler().GetByVisit)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/visits/"+id.String()+"/audit", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAuditHandler_GetByVisit_InvalidID(t *testing.T) {
	r := gin.New()
	r.GET("/api/mobile/visits/:id/audit", NewAuditHandler().GetByVisit)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/visits/invalid/audit", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAuditHandler_GetByVisit_WithPagination(t *testing.T) {
	id := uuid.New()
	r := gin.New()
	r.GET("/api/mobile/visits/:id/audit", NewAuditHandler().GetByVisit)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/visits/"+id.String()+"/audit?page=2&limit=10", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}
