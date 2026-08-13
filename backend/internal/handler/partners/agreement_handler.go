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

type PartnerAgreementHandler struct {
	repo repository.PartnerAgreementRepository
}

func NewPartnerAgreementHandler(repo repository.PartnerAgreementRepository) *PartnerAgreementHandler {
	return &PartnerAgreementHandler{repo: repo}
}

func (h *PartnerAgreementHandler) GetByID(c *gin.Context) {
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

func (h *PartnerAgreementHandler) ListByPartner(c *gin.Context) {
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

func (h *PartnerAgreementHandler) Create(c *gin.Context) {
	var agreement partners.PartnerAgreement
	if err := c.ShouldBindJSON(&agreement); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	agreement.ID = uuid.New()
	agreement.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &agreement); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, agreement)
}

func (h *PartnerAgreementHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var agreement partners.PartnerAgreement
	if err := c.ShouldBindJSON(&agreement); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(c.Request.Context(), id, &agreement); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, agreement)
}

func (h *PartnerAgreementHandler) Delete(c *gin.Context) {
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
