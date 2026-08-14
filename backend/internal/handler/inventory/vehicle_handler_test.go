package inventoryhandler

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
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	"localis-backend/internal/testutil/mocks"
)

func newRouter(t *testing.T, h *VehicleHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/vehicles", h.List)
	r.GET("/vehicles/:id", h.GetByID)
	r.POST("/vehicles", h.Create)
	r.PUT("/vehicles/:id", h.Update)
	r.DELETE("/vehicles/:id", h.Delete)
	return r
}

func sampleVehicle() *inventory.Vehicle {
	id := uuid.MustParse("53000000-0000-0000-0000-000000000001")
	plate := "AB-123-CD"
	status := shared.VehicleStatusAvailable
	return &inventory.Vehicle{
		ID:           id,
		LicensePlate: &plate,
		Name:         "Camioneta",
		Status:       status,
	}
}

func TestVehicleHandler_GetByID(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	item := sampleVehicle()
	repo.EXPECT().GetByID(mockCtx(), item.ID).Return(item, nil)

	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/vehicles/"+item.ID.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var got inventory.Vehicle
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Name != item.Name {
		t.Errorf("name = %q, want %q", got.Name, item.Name)
	}
}

func TestVehicleHandler_GetByID_InvalidUUID(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/vehicles/not-a-uuid", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestVehicleHandler_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(nil, &service.NotFoundError{Resource: "vehicle", ID: id.String()})

	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/vehicles/"+id.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestVehicleHandler_GetByID_InternalError(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(nil, errors.New("db down"))

	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/vehicles/"+id.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestVehicleHandler_List(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	item := sampleVehicle()
	result := &repository.ListResult[inventory.Vehicle]{
		Items: []inventory.Vehicle{*item},
		Total: 1,
	}
	repo.EXPECT().List(mockCtx(), 20, 0).Return(result, nil)

	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/vehicles", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var got repository.ListResult[inventory.Vehicle]
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Total != 1 || len(got.Items) != 1 {
		t.Errorf("total = %d items = %d, want 1/1", got.Total, len(got.Items))
	}
}

func TestVehicleHandler_List_WithPagination(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	repo.EXPECT().List(mockCtx(), 50, 10).Return(&repository.ListResult[inventory.Vehicle]{}, nil)

	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/vehicles?limit=50&offset=10", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestVehicleHandler_Create(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*inventory.Vehicle")).Return(nil)

	body := `{"name":"Camioneta","status":"active"}`
	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/vehicles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	var got inventory.Vehicle
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.ID == uuid.Nil {
		t.Error("expected ID to be set")
	}
}

func TestVehicleHandler_Create_InvalidBody(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/vehicles", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestVehicleHandler_Update(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*inventory.Vehicle")).Return(nil)

	body := `{"name":"Camioneta 2","status":"active"}`
	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/vehicles/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestVehicleHandler_Update_InvalidUUID(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/vehicles/nope", bytes.NewBufferString(`{"name":"x","status":"active"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestVehicleHandler_Update_InvalidBody(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/vehicles/"+uuid.New().String(), bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestVehicleHandler_Delete(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/vehicles/"+id.String(), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestVehicleHandler_Delete_InvalidUUID(t *testing.T) {
	repo := mocks.NewVehicleRepository(t)
	r := newRouter(t, NewVehicleHandler(repo))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/vehicles/bad", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func mockCtx() context.Context { return context.Background() }
