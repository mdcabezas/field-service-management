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

type PropertyHandler struct {
	repo repository.PropertyRepository
}

func NewPropertyHandler(repo repository.PropertyRepository) *PropertyHandler {
	return &PropertyHandler{repo: repo}
}

func (h *PropertyHandler) GetByID(c *gin.Context) {
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

func (h *PropertyHandler) ListByCustomerAddress(c *gin.Context) {
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

func (h *PropertyHandler) Create(c *gin.Context) {
	var prop customers.Property
	if err := c.ShouldBindJSON(&prop); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	prop.ID = uuid.New()
	prop.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &prop); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, prop)
}

func (h *PropertyHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var prop customers.Property
	if err := c.ShouldBindJSON(&prop); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(c.Request.Context(), id, &prop); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, prop)
}

func (h *PropertyHandler) Delete(c *gin.Context) {
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
