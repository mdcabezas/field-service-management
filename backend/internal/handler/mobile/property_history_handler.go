package mobilehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/mobile"
)

type PropertyHistoryHandler struct{}

func NewPropertyHistoryHandler() *PropertyHistoryHandler {
	return &PropertyHistoryHandler{}
}

func (h *PropertyHistoryHandler) GetByProperty(c *gin.Context) {
	propertyID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	_ = propertyID
	response := mobile.PropertyHistory{
		Visits: []mobile.PropertyHistoryVisit{},
		Total:  0,
	}

	handler.RespondJSON(c, http.StatusOK, response)
}
