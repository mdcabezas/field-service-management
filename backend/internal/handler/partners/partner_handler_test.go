package partnershandler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"localis-backend/internal/model/partners"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/testutil/mocks"
)

func mockCtx() context.Context { return context.Background() }

func newPartnerRouter(t *testing.T, h *PartnerHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/partners", h.List)
	r.GET("/partners/:id", h.GetByID)
	r.POST("/partners", h.Create)
	r.PUT("/partners/:id", h.Update)
	r.DELETE("/partners/:id", h.Delete)
	return r
}

func TestPartnerHandler_GetByID(t *testing.T) {
	repo := mocks.NewPartnerRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&partners.Partner{ID: id}, nil)

	r := newPartnerRouter(t, NewPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partners/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_GetByID_InvalidUUID(t *testing.T) {
	repo := mocks.NewPartnerRepository(t)
	r := newPartnerRouter(t, NewPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partners/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_List(t *testing.T) {
	repo := mocks.NewPartnerRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[partners.Partner]{Total: 0}, nil)

	r := newPartnerRouter(t, NewPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partners", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_Create(t *testing.T) {
	repo := mocks.NewPartnerRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*partners.Partner")).Return(nil)

	r := newPartnerRouter(t, NewPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/partners", bytes.NewBufferString(`{"name":"ACME","status":"active"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPartnerHandler_Update(t *testing.T) {
	repo := mocks.NewPartnerRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*partners.Partner")).Return(nil)

	r := newPartnerRouter(t, NewPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/partners/"+id.String(), bytes.NewBufferString(`{"name":"ACME 2","status":"active"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_Update_Error(t *testing.T) {
	repo := mocks.NewPartnerRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*partners.Partner")).Return(&repoErr{})

	r := newPartnerRouter(t, NewPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/partners/"+id.String(), bytes.NewBufferString(`{"name":"x","status":"active"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerHandler_Delete(t *testing.T) {
	repo := mocks.NewPartnerRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newPartnerRouter(t, NewPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/partners/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newContactRouter(t *testing.T, h *PartnerContactHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/partner-contacts/:id", h.GetByID)
	r.GET("/partners/:id/contacts", h.ListByPartner)
	r.POST("/partner-contacts", h.Create)
	r.DELETE("/partner-contacts/:id", h.Delete)
	return r
}

func TestPartnerContactHandler_GetByID(t *testing.T) {
	repo := mocks.NewPartnerContactRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&partners.PartnerContact{ID: id}, nil)

	r := newContactRouter(t, NewPartnerContactHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partner-contacts/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerContactHandler_ListByPartner(t *testing.T) {
	repo := mocks.NewPartnerContactRepository(t)
	pid := uuid.New()
	repo.EXPECT().ListByPartner(mockCtx(), pid).Return([]partners.PartnerContact{}, nil)

	r := newContactRouter(t, NewPartnerContactHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partners/"+pid.String()+"/contacts", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerContactHandler_Create(t *testing.T) {
	repo := mocks.NewPartnerContactRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*partners.PartnerContact")).Return(nil)

	pid := uuid.New()
	body := `{"partner_id":"` + pid.String() + `","name":"Juan"}`
	r := newContactRouter(t, NewPartnerContactHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/partner-contacts", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPartnerContactHandler_Delete(t *testing.T) {
	repo := mocks.NewPartnerContactRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newContactRouter(t, NewPartnerContactHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/partner-contacts/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newAgreementRouter(t *testing.T, h *PartnerAgreementHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/partner-agreements/:id", h.GetByID)
	r.GET("/partners/:id/agreements", h.ListByPartner)
	r.POST("/partner-agreements", h.Create)
	r.PUT("/partner-agreements/:id", h.Update)
	r.DELETE("/partner-agreements/:id", h.Delete)
	return r
}

func TestPartnerAgreementHandler_GetByID(t *testing.T) {
	repo := mocks.NewPartnerAgreementRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&partners.PartnerAgreement{ID: id}, nil)

	r := newAgreementRouter(t, NewPartnerAgreementHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partner-agreements/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerAgreementHandler_ListByPartner(t *testing.T) {
	repo := mocks.NewPartnerAgreementRepository(t)
	pid := uuid.New()
	repo.EXPECT().ListByPartner(mockCtx(), pid).Return([]partners.PartnerAgreement{}, nil)

	r := newAgreementRouter(t, NewPartnerAgreementHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partners/"+pid.String()+"/agreements", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerAgreementHandler_Create(t *testing.T) {
	repo := mocks.NewPartnerAgreementRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*partners.PartnerAgreement")).Return(nil)

	pid := uuid.New()
	stid := uuid.New()
	body := `{"partner_id":"` + pid.String() + `","service_type":"` + stid.String() + `","active":true}`
	r := newAgreementRouter(t, NewPartnerAgreementHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/partner-agreements", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPartnerAgreementHandler_Update(t *testing.T) {
	repo := mocks.NewPartnerAgreementRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*partners.PartnerAgreement")).Return(nil)

	pid := uuid.New()
	stid := uuid.New()
	body := `{"partner_id":"` + pid.String() + `","service_type":"` + stid.String() + `","active":true}`
	r := newAgreementRouter(t, NewPartnerAgreementHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/partner-agreements/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerAgreementHandler_Delete(t *testing.T) {
	repo := mocks.NewPartnerAgreementRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newAgreementRouter(t, NewPartnerAgreementHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/partner-agreements/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newAgreementDocRouter(t *testing.T, h *PartnerAgreementDocHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/partner-agreement-docs/:id", h.GetByID)
	r.GET("/partner-agreements/:id/docs", h.ListByAgreement)
	r.POST("/partner-agreement-docs", h.Create)
	r.DELETE("/partner-agreement-docs/:id", h.Delete)
	return r
}

func TestPartnerAgreementDocHandler_GetByID(t *testing.T) {
	repo := mocks.NewPartnerAgreementDocRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&partners.PartnerAgreementDoc{ID: id}, nil)

	r := newAgreementDocRouter(t, NewPartnerAgreementDocHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partner-agreement-docs/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerAgreementDocHandler_ListByAgreement(t *testing.T) {
	repo := mocks.NewPartnerAgreementDocRepository(t)
	aid := uuid.New()
	repo.EXPECT().ListByAgreement(mockCtx(), aid).Return([]partners.PartnerAgreementDoc{}, nil)

	r := newAgreementDocRouter(t, NewPartnerAgreementDocHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partner-agreements/"+aid.String()+"/docs", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerAgreementDocHandler_Create(t *testing.T) {
	repo := mocks.NewPartnerAgreementDocRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*partners.PartnerAgreementDoc")).Return(nil)

	aid := uuid.New()
	body := `{"agreement_id":"` + aid.String() + `","work_type":"install","doc_name":"contrato","required":true}`
	r := newAgreementDocRouter(t, NewPartnerAgreementDocHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/partner-agreement-docs", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPartnerAgreementDocHandler_Delete(t *testing.T) {
	repo := mocks.NewPartnerAgreementDocRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newAgreementDocRouter(t, NewPartnerAgreementDocHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/partner-agreement-docs/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newAgreementFormRouter(t *testing.T, h *PartnerAgreementFormHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/partner-agreement-forms/:id", h.GetByID)
	r.GET("/partner-agreements/:id/forms", h.ListByAgreement)
	r.POST("/partner-agreement-forms", h.Create)
	r.DELETE("/partner-agreement-forms/:id", h.Delete)
	return r
}

func TestPartnerAgreementFormHandler_GetByID(t *testing.T) {
	repo := mocks.NewPartnerAgreementFormRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&partners.PartnerAgreementForm{ID: id}, nil)

	r := newAgreementFormRouter(t, NewPartnerAgreementFormHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partner-agreement-forms/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerAgreementFormHandler_ListByAgreement(t *testing.T) {
	repo := mocks.NewPartnerAgreementFormRepository(t)
	aid := uuid.New()
	repo.EXPECT().ListByAgreement(mockCtx(), aid).Return([]partners.PartnerAgreementForm{}, nil)

	r := newAgreementFormRouter(t, NewPartnerAgreementFormHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partner-agreements/"+aid.String()+"/forms", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPartnerAgreementFormHandler_Create(t *testing.T) {
	repo := mocks.NewPartnerAgreementFormRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*partners.PartnerAgreementForm")).Return(nil)

	aid := uuid.New()
	body := `{"agreement_id":"` + aid.String() + `","work_type":"install","quantity":3}`
	r := newAgreementFormRouter(t, NewPartnerAgreementFormHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/partner-agreement-forms", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPartnerAgreementFormHandler_Delete(t *testing.T) {
	repo := mocks.NewPartnerAgreementFormRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newAgreementFormRouter(t, NewPartnerAgreementFormHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/partner-agreement-forms/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newSLARouter(t *testing.T, h *SLAHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/slas/:id", h.GetByID)
	r.GET("/partners/:id/slas", h.ListByPartner)
	r.POST("/slas", h.Create)
	r.PUT("/slas/:id", h.Update)
	r.DELETE("/slas/:id", h.Delete)
	return r
}

func TestSLAHandler_GetByID(t *testing.T) {
	repo := mocks.NewSLARepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&partners.SLA{ID: id}, nil)

	r := newSLARouter(t, NewSLAHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/slas/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestSLAHandler_ListByPartner(t *testing.T) {
	repo := mocks.NewSLARepository(t)
	pid := uuid.New()
	repo.EXPECT().ListByPartner(mockCtx(), pid).Return([]partners.SLA{}, nil)

	r := newSLARouter(t, NewSLAHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/partners/"+pid.String()+"/slas", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestSLAHandler_Create(t *testing.T) {
	repo := mocks.NewSLARepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*partners.SLA")).Return(nil)

	pid := uuid.New()
	body := `{"partner_id":"` + pid.String() + `","name":"SLA 1","active":true,"valid_from":"2026-01-01T00:00:00Z"}`
	r := newSLARouter(t, NewSLAHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/slas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSLAHandler_Update(t *testing.T) {
	repo := mocks.NewSLARepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*partners.SLA")).Return(nil)

	pid := uuid.New()
	body := `{"partner_id":"` + pid.String() + `","name":"SLA 1","active":true,"valid_from":"2026-01-01T00:00:00Z"}`
	r := newSLARouter(t, NewSLAHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/slas/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestSLAHandler_Delete(t *testing.T) {
	repo := mocks.NewSLARepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newSLARouter(t, NewSLAHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/slas/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

type repoErr struct{}

func (*repoErr) Error() string { return "boom" }

var _ = shared.PartnerStatus("active")
var _ = time.Time{}
