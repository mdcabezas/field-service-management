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

type TechnicianHandler struct {
	repo repository.TechnicianRepository
}

func NewTechnicianHandler(repo repository.TechnicianRepository) *TechnicianHandler {
	return &TechnicianHandler{repo: repo}
}

func (h *TechnicianHandler) GetByID(c *gin.Context) {
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

func (h *TechnicianHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	result, err := h.repo.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, result)
}

func (h *TechnicianHandler) Create(c *gin.Context) {
	var tech customers.Technician
	if err := c.ShouldBindJSON(&tech); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	tech.ID = uuid.New()
	tech.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &tech); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, tech)
}

func (h *TechnicianHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var tech customers.Technician
	if err := c.ShouldBindJSON(&tech); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(c.Request.Context(), id, &tech); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, tech)
}

func (h *TechnicianHandler) Delete(c *gin.Context) {
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
