package customershandler

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
	"localis-backend/internal/model/customers"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/testutil/mocks"
)

func mockCtx() context.Context { return context.Background() }

func newCustomerRouter(t *testing.T, h *CustomerHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/customers", h.List)
	r.GET("/customers/:id", h.GetByID)
	r.POST("/customers", h.Create)
	r.PUT("/customers/:id", h.Update)
	r.DELETE("/customers/:id", h.Delete)
	return r
}

func TestCustomerHandler_GetByID(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	id := uuid.New()
	c := &customers.Customer{ID: id, Name: "ACME"}
	repo.EXPECT().GetByID(mockCtx(), id).Return(c, nil)

	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerHandler_GetByID_InvalidUUID(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerHandler_List(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[customers.Customer]{Total: 0}, nil)

	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerHandler_Create(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*customers.Customer")).Return(nil)

	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(`{"name":"ACME"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerHandler_Create_InvalidBody(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBufferString(`{`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerHandler_Update(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*customers.Customer")).Return(nil)

	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/customers/"+id.String(), bytes.NewBufferString(`{"name":"ACME 2"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerHandler_Update_InvalidUUID(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/customers/bad", bytes.NewBufferString(`{"name":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerHandler_Update_Error(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*customers.Customer")).Return(&repoErr{})

	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/customers/"+id.String(), bytes.NewBufferString(`{"name":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerHandler_Delete(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/customers/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerHandler_Delete_InvalidUUID(t *testing.T) {
	repo := mocks.NewCustomerRepository(t)
	r := newCustomerRouter(t, NewCustomerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/customers/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func newPropertyRouter(t *testing.T, h *PropertyHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/properties/:id", h.GetByID)
	r.GET("/customer-addresses/:id/properties", h.ListByCustomerAddress)
	r.POST("/properties", h.Create)
	r.PUT("/properties/:id", h.Update)
	r.DELETE("/properties/:id", h.Delete)
	return r
}

func TestPropertyHandler_GetByID(t *testing.T) {
	repo := mocks.NewPropertyRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&customers.Property{ID: id}, nil)

	r := newPropertyRouter(t, NewPropertyHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/properties/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPropertyHandler_GetByID_InvalidUUID(t *testing.T) {
	repo := mocks.NewPropertyRepository(t)
	r := newPropertyRouter(t, NewPropertyHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/properties/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPropertyHandler_ListByCustomerAddress(t *testing.T) {
	repo := mocks.NewPropertyRepository(t)
	caID := uuid.New()
	repo.EXPECT().ListByCustomerAddress(mockCtx(), caID).Return([]customers.Property{}, nil)

	r := newPropertyRouter(t, NewPropertyHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customer-addresses/"+caID.String()+"/properties", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPropertyHandler_ListByCustomerAddress_InvalidUUID(t *testing.T) {
	repo := mocks.NewPropertyRepository(t)
	r := newPropertyRouter(t, NewPropertyHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customer-addresses/bad/properties", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPropertyHandler_Create(t *testing.T) {
	repo := mocks.NewPropertyRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*customers.Property")).Return(nil)

	caID := uuid.New()
	r := newPropertyRouter(t, NewPropertyHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/properties", bytes.NewBufferString(`{"customer_address_id":"`+caID.String()+`"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPropertyHandler_Update(t *testing.T) {
	repo := mocks.NewPropertyRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*customers.Property")).Return(nil)

	caID := uuid.New()
	r := newPropertyRouter(t, NewPropertyHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/properties/"+id.String(), bytes.NewBufferString(`{"customer_address_id":"`+caID.String()+`"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPropertyHandler_Delete(t *testing.T) {
	repo := mocks.NewPropertyRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newPropertyRouter(t, NewPropertyHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/properties/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newTechnicianRouter(t *testing.T, h *TechnicianHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/technicians", h.List)
	r.GET("/technicians/:id", h.GetByID)
	r.POST("/technicians", h.Create)
	r.PUT("/technicians/:id", h.Update)
	r.DELETE("/technicians/:id", h.Delete)
	return r
}

func TestTechnicianHandler_GetByID(t *testing.T) {
	repo := mocks.NewTechnicianRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&customers.Technician{ID: id}, nil)

	r := newTechnicianRouter(t, NewTechnicianHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/technicians/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechnicianHandler_List(t *testing.T) {
	repo := mocks.NewTechnicianRepository(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[customers.Technician]{Total: 0}, nil)

	r := newTechnicianRouter(t, NewTechnicianHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/technicians", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechnicianHandler_Create(t *testing.T) {
	repo := mocks.NewTechnicianRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*customers.Technician")).Return(nil)

	r := newTechnicianRouter(t, NewTechnicianHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/technicians", bytes.NewBufferString(`{"name":"Juan","is_active":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestTechnicianHandler_Update(t *testing.T) {
	repo := mocks.NewTechnicianRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*customers.Technician")).Return(nil)

	r := newTechnicianRouter(t, NewTechnicianHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/technicians/"+id.String(), bytes.NewBufferString(`{"name":"Juan","is_active":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechnicianHandler_Delete(t *testing.T) {
	repo := mocks.NewTechnicianRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newTechnicianRouter(t, NewTechnicianHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/technicians/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newAddressRouter(t *testing.T, h *CustomerAddressHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/customer-addresses/:id", h.GetByID)
	r.GET("/customers/:id/addresses", h.ListByCustomer)
	r.POST("/customer-addresses", h.Create)
	r.PUT("/customer-addresses/:id", h.Update)
	r.DELETE("/customer-addresses/:id", h.Delete)
	return r
}

func TestCustomerAddressHandler_GetByID(t *testing.T) {
	repo := mocks.NewCustomerAddressRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&customers.CustomerAddress{ID: id}, nil)

	r := newAddressRouter(t, NewCustomerAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customer-addresses/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerAddressHandler_ListByCustomer(t *testing.T) {
	repo := mocks.NewCustomerAddressRepository(t)
	cid := uuid.New()
	repo.EXPECT().ListByCustomer(mockCtx(), cid).Return([]customers.CustomerAddress{}, nil)

	r := newAddressRouter(t, NewCustomerAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customers/"+cid.String()+"/addresses", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerAddressHandler_Create(t *testing.T) {
	repo := mocks.NewCustomerAddressRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*customers.CustomerAddress")).Return(nil)

	cid := uuid.New()
	aid := uuid.New()
	body := `{"customer_id":"` + cid.String() + `","address_id":"` + aid.String() + `","type":"home"}`
	r := newAddressRouter(t, NewCustomerAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/customer-addresses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCustomerAddressHandler_Update(t *testing.T) {
	repo := mocks.NewCustomerAddressRepository(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*customers.CustomerAddress")).Return(nil)

	cid := uuid.New()
	aid := uuid.New()
	body := `{"customer_id":"` + cid.String() + `","address_id":"` + aid.String() + `","type":"home"}`
	r := newAddressRouter(t, NewCustomerAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/customer-addresses/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCustomerAddressHandler_Delete(t *testing.T) {
	repo := mocks.NewCustomerAddressRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newAddressRouter(t, NewCustomerAddressHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/customer-addresses/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newAddressPartnerRouter(t *testing.T, h *CustomerAddressPartnerHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/customer-address-partners/:id", h.GetByID)
	r.GET("/customer-addresses/:id/partners", h.ListByCustomerAddress)
	r.POST("/customer-address-partners", h.Create)
	r.DELETE("/customer-address-partners/:id", h.Delete)
	return r
}

func TestCustomerAddressPartnerHandler_GetByID(t *testing.T) {
	repo := mocks.NewCustomerAddressPartnerRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&customers.CustomerAddressPartner{ID: id}, nil)

	r := newAddressPartnerRouter(t, NewCustomerAddressPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customer-address-partners/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerAddressPartnerHandler_ListByCustomerAddress(t *testing.T) {
	repo := mocks.NewCustomerAddressPartnerRepository(t)
	caID := uuid.New()
	repo.EXPECT().ListByCustomerAddress(mockCtx(), caID).Return([]customers.CustomerAddressPartner{}, nil)

	r := newAddressPartnerRouter(t, NewCustomerAddressPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/customer-addresses/"+caID.String()+"/partners", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestCustomerAddressPartnerHandler_Create(t *testing.T) {
	repo := mocks.NewCustomerAddressPartnerRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*customers.CustomerAddressPartner")).Return(nil)

	caID := uuid.New()
	pid := uuid.New()
	body := `{"customer_address_id":"` + caID.String() + `","partner_id":"` + pid.String() + `","from_date":"2026-01-01T00:00:00Z"}`
	r := newAddressPartnerRouter(t, NewCustomerAddressPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/customer-address-partners", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCustomerAddressPartnerHandler_Delete(t *testing.T) {
	repo := mocks.NewCustomerAddressPartnerRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newAddressPartnerRouter(t, NewCustomerAddressPartnerHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/customer-address-partners/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func newTechCertRouter(t *testing.T, h *TechCertificationHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/tech-certifications/:id", h.GetByID)
	r.GET("/technicians/:id/certifications", h.ListByTech)
	r.POST("/tech-certifications", h.Create)
	r.DELETE("/tech-certifications/:id", h.Delete)
	return r
}

func TestTechCertificationHandler_GetByID(t *testing.T) {
	repo := mocks.NewTechCertificationRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&customers.TechCertification{ID: id}, nil)

	r := newTechCertRouter(t, NewTechCertificationHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tech-certifications/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechCertificationHandler_ListByTech(t *testing.T) {
	repo := mocks.NewTechCertificationRepository(t)
	tid := uuid.New()
	repo.EXPECT().ListByTech(mockCtx(), tid).Return([]customers.TechCertification{}, nil)

	r := newTechCertRouter(t, NewTechCertificationHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/technicians/"+tid.String()+"/certifications", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestTechCertificationHandler_Create(t *testing.T) {
	repo := mocks.NewTechCertificationRepository(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*customers.TechCertification")).Return(nil)

	tid := uuid.New()
	body := `{"tech_id":"` + tid.String() + `"}`
	r := newTechCertRouter(t, NewTechCertificationHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tech-certifications", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestTechCertificationHandler_Delete(t *testing.T) {
	repo := mocks.NewTechCertificationRepository(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newTechCertRouter(t, NewTechCertificationHandler(repo))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/tech-certifications/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

type repoErr struct{}

func (*repoErr) Error() string { return "boom" }

var _ = shared.CustomerAddressType("home")
var _ = time.Time{}
