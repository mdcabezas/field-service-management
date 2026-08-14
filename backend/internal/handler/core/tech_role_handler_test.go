package corehandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"localis-backend/internal/model/core"
	"localis-backend/internal/repository"
	"localis-backend/internal/testutil/mocks"
)

func newRoleRouter(t *testing.T, h *TechRoleHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/tech-roles", h.List)
	r.GET("/tech-roles/:id", h.GetByID)
	r.POST("/tech-roles", h.Create)
	r.PUT("/tech-roles/:id", h.Update)
	r.DELETE("/tech-roles/:id", h.Delete)
	return r
}

func TestTechRoleHandler_GetByID(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	role := &core.TechRole{ID: uuid.New(), Code: "tech", Name: "Tech", Active: true}
	repo.EXPECT().GetByID(mockCtx(), role.ID).Return(role, nil)

	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tech-roles/"+role.ID.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechRoleHandler_GetByID_InvalidUUID(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tech-roles/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechRoleHandler_List(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[core.TechRole]{Total: 0}, nil)

	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tech-roles", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechRoleHandler_Create(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.TechRole")).Return(nil)

	body := `{"code":"tech","name":"Tech","active":true}`
	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tech-roles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
	var got core.TechRole
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ID == uuid.Nil {
		t.Fatal("expected ID to be set")
	}
}

func TestTechRoleHandler_Create_InvalidBody(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tech-roles", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechRoleHandler_Update(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*core.TechRole")).Return(nil)

	body := `{"code":"tech","name":"Tech","active":true}`
	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/tech-roles/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechRoleHandler_Update_InvalidUUID(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/tech-roles/bad", bytes.NewBufferString(`{"code":"x","name":"y","active":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechRoleHandler_Update_InvalidBody(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/tech-roles/"+uuid.New().String(), bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechRoleHandler_Delete(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/tech-roles/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechRoleHandler_Delete_InvalidUUID(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/tech-roles/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechRoleHandler_Delete_Error(t *testing.T) {
	repo := mocks.NewCoreTechRoleRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(&repositoryErr{})

	r := newRoleRouter(t, NewTechRoleHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/tech-roles/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

type repositoryErr struct{}

func (repositoryErr) Error() string { return "boom" }
