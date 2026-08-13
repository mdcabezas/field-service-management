package partnershandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/partners"
	"localis-backend/internal/repository"
)

type SLAHandler struct {
	repo repository.SLARepository
}

func NewSLAHandler(repo repository.SLARepository) *SLAHandler {
	return &SLAHandler{repo: repo}
}

func (h *SLAHandler) GetByID(c *gin.Context) {
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

func (h *SLAHandler) ListByPartner(c *gin.Context) {
	partnerID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	items, err := h.repo.ListByPartner(c.Request.Context(), partnerID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *SLAHandler) Create(c *gin.Context) {
	var sla partners.SLA
	if err := c.ShouldBindJSON(&sla); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	sla.ID = uuid.New()
	sla.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &sla); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, sla)
}

func (h *SLAHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var sla partners.SLA
	if err := c.ShouldBindJSON(&sla); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(c.Request.Context(), id, &sla); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, sla)
}

func (h *SLAHandler) Delete(c *gin.Context) {
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
