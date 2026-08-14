package geocodinghandler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"localis-backend/internal/model/geocoding"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	"localis-backend/internal/testutil/mocks"
)

func newAddressRouter(t *testing.T, h *AddressHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/addresses", h.List)
	r.GET("/addresses/:id", h.GetByID)
	r.POST("/addresses", h.Create)
	r.PUT("/addresses/:id", h.Update)
	r.DELETE("/addresses/:id", h.Delete)
	r.GET("/addresses/nearby", h.FindNearby)
	return r
}

func mockCtx() context.Context { return context.Background() }

func TestAddressHandler_GetByID(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	id := uuid.New()
	addr := &geocoding.Address{ID: id, Street: "Av. X"}
	repo.EXPECT().GetByID(mockCtx(), id).Return(addr, nil)

	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_GetByID_InvalidUUID(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(nil, &service.NotFoundError{Resource: "address", ID: id.String()})

	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_List(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[geocoding.Address]{Total: 0}, nil)

	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_Create(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*geocoding.Address")).Return(nil)

	body := `{"street":"Av. X"}`
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_Create_InvalidBody(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_Update(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*geocoding.Address")).Return(nil)

	body := `{"street":"Av. Y"}`
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/addresses/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_Update_InvalidUUID(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/addresses/bad", bytes.NewBufferString(`{"street":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_Update_InvalidBody(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/addresses/"+uuid.New().String(), bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_Delete(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/addresses/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_FindNearby_Valid(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	repo.EXPECT().FindNearby(mockCtx(), -33.45, -70.66, 1000.0).Return([]geocoding.Address{}, nil)

	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses/nearby?lat=-33.45&lng=-70.66&radius=1000", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_FindNearby_MissingParams(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses/nearby", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_FindNearby_InvalidLat(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses/nearby?lat=abc&lng=-70&radius=1000", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_FindNearby_OutOfRange(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses/nearby?lat=200&lng=-70&radius=1000", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_FindNearby_BadRadius(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses/nearby?lat=-33&lng=-70&radius=0", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestAddressHandler_FindNearby_InternalError(t *testing.T) {
	repo := mocks.NewGeocodingAddressRepository(t)
	repo.EXPECT().FindNearby(mockCtx(), -33.0, -70.0, 10.0).Return(nil, errBoom())

	r := newAddressRouter(t, NewAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/addresses/nearby?lat=-33&lng=-70&radius=10", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

func errBoom() error { return &boomError{} }

type boomError struct{}

func (*boomError) Error() string { return "boom" }
