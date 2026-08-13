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

type PartnerContactHandler struct {
	repo repository.PartnerContactRepository
}

func NewPartnerContactHandler(repo repository.PartnerContactRepository) *PartnerContactHandler {
	return &PartnerContactHandler{repo: repo}
}

func (h *PartnerContactHandler) GetByID(c *gin.Context) {
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

func (h *PartnerContactHandler) ListByPartner(c *gin.Context) {
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

func (h *PartnerContactHandler) Create(c *gin.Context) {
	var contact partners.PartnerContact
	if err := c.ShouldBindJSON(&contact); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	contact.ID = uuid.New()
	contact.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &contact); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, contact)
}

func (h *PartnerContactHandler) Delete(c *gin.Context) {
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
