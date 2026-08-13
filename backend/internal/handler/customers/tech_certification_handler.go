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

type TechCertificationHandler struct {
	repo repository.TechCertificationRepository
}

func NewTechCertificationHandler(repo repository.TechCertificationRepository) *TechCertificationHandler {
	return &TechCertificationHandler{repo: repo}
}

func (h *TechCertificationHandler) GetByID(c *gin.Context) {
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

func (h *TechCertificationHandler) ListByTech(c *gin.Context) {
	techID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	items, err := h.repo.ListByTech(c.Request.Context(), techID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *TechCertificationHandler) Create(c *gin.Context) {
	var cert customers.TechCertification
	if err := c.ShouldBindJSON(&cert); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	cert.ID = uuid.New()
	cert.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &cert); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, cert)
}

func (h *TechCertificationHandler) Delete(c *gin.Context) {
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
