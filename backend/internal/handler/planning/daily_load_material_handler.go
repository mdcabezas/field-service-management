package planninghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/planning"
	planningservice "localis-backend/internal/service/planning"
)

type DailyLoadMaterialHandler struct {
	svc *planningservice.DailyLoadService
}

func NewDailyLoadMaterialHandler(svc *planningservice.DailyLoadService) *DailyLoadMaterialHandler {
	return &DailyLoadMaterialHandler{svc: svc}
}

func (h *DailyLoadMaterialHandler) GetByID(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.svc.GetMaterialByID(c.Request.Context(), id)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *DailyLoadMaterialHandler) ListByPlan(c *gin.Context) {
	planID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid plan id")
		return
	}
	items, err := h.svc.ListMaterialsByPlan(c.Request.Context(), planID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *DailyLoadMaterialHandler) Create(c *gin.Context) {
	var body planning.DailyLoadMaterial
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.LoadMaterial(c.Request.Context(), &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, body)
}

func (h *DailyLoadMaterialHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var body planning.DailyLoadMaterial
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.UpdateMaterial(c.Request.Context(), id, &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, body)
}

func (h *DailyLoadMaterialHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteMaterial(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusNoContent, nil)
}
