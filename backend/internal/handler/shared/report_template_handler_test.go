package sharedhandler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/testutil/mocks"
)

func newReportTemplateRouter(t *testing.T, h *ReportTemplateHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/report-templates", h.List)
	r.GET("/report-templates/:id", h.GetByID)
	r.POST("/report-templates", h.Create)
	r.PUT("/report-templates/:id", h.Update)
	r.DELETE("/report-templates/:id", h.Delete)
	return r
}

func mockCtx() context.Context { return context.Background() }

func TestReportTemplateHandler_GetByID(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&shared.ReportTemplate{ID: id, Name: "t"}, nil)

	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/report-templates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_GetByID_InvalidUUID(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/report-templates/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_GetByID_Error(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(nil, &repoErr{})

	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/report-templates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_List(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[shared.ReportTemplate]{Total: 0}, nil)

	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/report-templates", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_List_Error(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(nil, &repoErr{})

	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/report-templates", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_Create(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*shared.ReportTemplate")).Return(nil)

	body := `{"name":"t","fields_json":[{"k":"v"}]}`
	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/report-templates", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_Create_InvalidBody(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/report-templates", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_Update(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*shared.ReportTemplate")).Return(nil)

	body := `{"name":"t","fields_json":[]}`
	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/report-templates/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_Update_InvalidUUID(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/report-templates/bad", bytes.NewBufferString(`{"name":"x","fields_json":[]}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_Update_Error(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*shared.ReportTemplate")).Return(&repoErr{})

	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/report-templates/"+id.String(), bytes.NewBufferString(`{"name":"x","fields_json":[]}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_Delete(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/report-templates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_Delete_InvalidUUID(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/report-templates/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestReportTemplateHandler_Delete_Error(t *testing.T) {
	repo := mocks.NewReportTemplateRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(&repoErr{})

	r := newReportTemplateRouter(t, NewReportTemplateHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/report-templates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

type repoErr struct{}

func (*repoErr) Error() string { return "boom" }
