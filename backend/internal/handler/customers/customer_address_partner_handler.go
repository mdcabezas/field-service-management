package customershandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/customers"
	"localis-backend/internal/repository"
)

type CustomerAddressPartnerHandler struct {
	repo repository.CustomerAddressPartnerRepository
}

func NewCustomerAddressPartnerHandler(repo repository.CustomerAddressPartnerRepository) *CustomerAddressPartnerHandler {
	return &CustomerAddressPartnerHandler{repo: repo}
}

func (h *CustomerAddressPartnerHandler) GetByID(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *CustomerAddressPartnerHandler) ListByCustomerAddress(c *gin.Context) {
	caID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	items, err := h.repo.ListByCustomerAddress(c.Request.Context(), caID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *CustomerAddressPartnerHandler) Create(c *gin.Context) {
	var cap customers.CustomerAddressPartner
	if err := c.ShouldBindJSON(&cap); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	cap.ID = uuid.New()
	cap.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &cap); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, cap)
}

func (h *CustomerAddressPartnerHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusNoContent, nil)
}
