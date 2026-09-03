package mobilehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
)

type PartnerHandler struct{}

func NewPartnerHandler() *PartnerHandler {
	return &PartnerHandler{}
}

func (h *PartnerHandler) GetContacts(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	_ = id
	handler.RespondJSON(c, http.StatusOK, gin.H{"contacts": []interface{}{}})
}

func (h *PartnerHandler) GetByProperty(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	_ = id
	handler.RespondJSON(c, http.StatusOK, gin.H{"partners": []interface{}{}})
}

func (h *PartnerHandler) ListTemplates(c *gin.Context) {
	partnerIDStr := c.Query("partner_id")
	_ = partnerIDStr

	handler.RespondJSON(c, http.StatusOK, gin.H{"templates": []interface{}{}})
}
