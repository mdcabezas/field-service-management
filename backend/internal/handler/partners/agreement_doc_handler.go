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

type PartnerAgreementDocHandler struct {
	repo repository.PartnerAgreementDocRepository
}

func NewPartnerAgreementDocHandler(repo repository.PartnerAgreementDocRepository) *PartnerAgreementDocHandler {
	return &PartnerAgreementDocHandler{repo: repo}
}

func (h *PartnerAgreementDocHandler) GetByID(c *gin.Context) {
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

func (h *PartnerAgreementDocHandler) ListByAgreement(c *gin.Context) {
	agreementID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	items, err := h.repo.ListByAgreement(c.Request.Context(), agreementID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *PartnerAgreementDocHandler) Create(c *gin.Context) {
	var doc partners.PartnerAgreementDoc
	if err := c.ShouldBindJSON(&doc); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	doc.ID = uuid.New()
	doc.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &doc); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, doc)
}

func (h *PartnerAgreementDocHandler) Delete(c *gin.Context) {
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
