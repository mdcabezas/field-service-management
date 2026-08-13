package planninghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/planning"
	planningservice "localis-backend/internal/service/planning"
)

type DailyLoadEPPHandler struct {
	svc *planningservice.DailyLoadService
}

func NewDailyLoadEPPHandler(svc *planningservice.DailyLoadService) *DailyLoadEPPHandler {
	return &DailyLoadEPPHandler{svc: svc}
}

func (h *DailyLoadEPPHandler) GetByID(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.svc.GetEPPByID(c.Request.Context(), id)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *DailyLoadEPPHandler) ListByPlan(c *gin.Context) {
	planID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid plan id")
		return
	}
	items, err := h.svc.ListEPPsByPlan(c.Request.Context(), planID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *DailyLoadEPPHandler) Create(c *gin.Context) {
	var body planning.DailyLoadEPP
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.LoadEPP(c.Request.Context(), &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, body)
}

func (h *DailyLoadEPPHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var body planning.DailyLoadEPP
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.UpdateEPP(c.Request.Context(), id, &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, body)
}

func (h *DailyLoadEPPHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteEPP(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusNoContent, nil)
}
