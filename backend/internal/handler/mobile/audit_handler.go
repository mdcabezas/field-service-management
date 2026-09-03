package mobilehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/mobile"
)

type AuditHandler struct{}

func NewAuditHandler() *AuditHandler {
	return &AuditHandler{}
}

func (h *AuditHandler) GetByVisit(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	_ = id
	limit, offset := handler.ParsePagination(c)

	response := mobile.AuditResponse{
		Entries: []mobile.AuditEntry{},
		Total:   0,
		Page:    offset/limit + 1,
		Limit:   limit,
	}

	handler.RespondJSON(c, http.StatusOK, response)
}
