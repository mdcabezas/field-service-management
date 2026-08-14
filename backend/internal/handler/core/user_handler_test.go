package corehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"localis-backend/internal/model/core"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	"localis-backend/internal/testutil/mocks"
)

func newUserRouter(t *testing.T, h *UserHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/users", h.List)
	r.GET("/users/:id", h.GetByID)
	r.POST("/users", h.Create)
	r.PUT("/users/:id", h.Update)
	r.DELETE("/users/:id", h.Delete)
	return r
}

func sampleUser() *core.User {
	return &core.User{
		ID:    uuid.New(),
		Email: "u@example.com",
		Role:  "admin",
		Name:  "Admin",
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	u := sampleUser()
	repo.EXPECT().GetByID(mockCtx(), u.ID.String()).Return(u, nil)

	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/"+u.ID.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id.String()).Return(nil, &service.NotFoundError{Resource: "user", ID: id.String()})

	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestUserHandler_GetByID_InvalidUUID(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestUserHandler_GetByID_InternalError(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id.String()).Return(nil, errors.New("db down"))

	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestUserHandler_List(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	u := sampleUser()
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[core.User]{Items: []core.User{*u}, Total: 1}, nil)

	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var got repository.ListResult[core.User]
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Total != 1 {
		t.Fatalf("total = %d", got.Total)
	}
}

func TestUserHandler_Create(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*core.User")).Return(nil)

	body := `{"email":"u@example.com","role":"admin","name":"Admin"}`
	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestUserHandler_Create_InvalidBody(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestUserHandler_Update(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id.String(), mock.AnythingOfType("*core.User")).Return(nil)

	body := `{"email":"u@example.com","role":"manager","name":"Admin"}`
	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/users/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestUserHandler_Update_InvalidUUID(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/users/bad", bytes.NewBufferString(`{"email":"a@b.c","role":"x","name":"y"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestUserHandler_Delete(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id.String()).Return(nil)

	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/users/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestUserHandler_Delete_InvalidUUID(t *testing.T) {
	repo := mocks.NewCoreUserRepository(t)
	r := newUserRouter(t, NewUserHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/users/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func mockCtx() context.Context { return context.Background() }
