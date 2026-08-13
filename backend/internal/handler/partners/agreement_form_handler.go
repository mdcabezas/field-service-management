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

type PartnerAgreementFormHandler struct {
	repo repository.PartnerAgreementFormRepository
}

func NewPartnerAgreementFormHandler(repo repository.PartnerAgreementFormRepository) *PartnerAgreementFormHandler {
	return &PartnerAgreementFormHandler{repo: repo}
}

func (h *PartnerAgreementFormHandler) GetByID(c *gin.Context) {
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

func (h *PartnerAgreementFormHandler) ListByAgreement(c *gin.Context) {
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

func (h *PartnerAgreementFormHandler) Create(c *gin.Context) {
	var form partners.PartnerAgreementForm
	if err := c.ShouldBindJSON(&form); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	form.ID = uuid.New()
	form.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &form); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, form)
}

func (h *PartnerAgreementFormHandler) Delete(c *gin.Context) {
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
