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

type CustomerAddressHandler struct {
	repo repository.CustomerAddressRepository
}

func NewCustomerAddressHandler(repo repository.CustomerAddressRepository) *CustomerAddressHandler {
	return &CustomerAddressHandler{repo: repo}
}

func (h *CustomerAddressHandler) GetByID(c *gin.Context) {
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

func (h *CustomerAddressHandler) ListByCustomer(c *gin.Context) {
	customerID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	items, err := h.repo.ListByCustomer(c.Request.Context(), customerID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, gin.H{"items": items, "total": len(items)})
}

func (h *CustomerAddressHandler) Create(c *gin.Context) {
	var addr customers.CustomerAddress
	if err := c.ShouldBindJSON(&addr); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	addr.ID = uuid.New()
	addr.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &addr); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, addr)
}

func (h *CustomerAddressHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var addr customers.CustomerAddress
	if err := c.ShouldBindJSON(&addr); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(c.Request.Context(), id, &addr); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, addr)
}

func (h *CustomerAddressHandler) Delete(c *gin.Context) {
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
