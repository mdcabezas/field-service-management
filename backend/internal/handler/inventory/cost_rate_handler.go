package inventoryhandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/repository"
)

type CostRateHandler struct {
	repo repository.CostRateRepository
}

func NewCostRateHandler(repo repository.CostRateRepository) *CostRateHandler {
	return &CostRateHandler{repo: repo}
}

func (h *CostRateHandler) GetByID(c *gin.Context) {
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

func (h *CostRateHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	result, err := h.repo.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, result)
}

func (h *CostRateHandler) Create(c *gin.Context) {
	var item inventory.CostRate
	if err := c.ShouldBindJSON(&item); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	item.ID = uuid.New()
	item.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &item); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, item)
}

func (h *CostRateHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var item inventory.CostRate
	if err := c.ShouldBindJSON(&item); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(c.Request.Context(), id, &item); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *CostRateHandler) Delete(c *gin.Context) {
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
