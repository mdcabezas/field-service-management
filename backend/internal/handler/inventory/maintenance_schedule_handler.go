package inventoryhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	invmodel "localis-backend/internal/model/inventory"
	invcsv "localis-backend/internal/service/inventory"
)

type MaintenanceScheduleHandler struct {
	svc *invcsv.MaintenanceService
}

func NewMaintenanceScheduleHandler(svc *invcsv.MaintenanceService) *MaintenanceScheduleHandler {
	return &MaintenanceScheduleHandler{svc: svc}
}

func (h *MaintenanceScheduleHandler) GetByID(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.svc.GetScheduleByID(c.Request.Context(), id)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *MaintenanceScheduleHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	result, err := h.svc.ListSchedules(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, result)
}

func (h *MaintenanceScheduleHandler) Create(c *gin.Context) {
	var item invmodel.MaintenanceSchedule
	if err := c.ShouldBindJSON(&item); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.CreateSchedule(c.Request.Context(), &item); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, item)
}

func (h *MaintenanceScheduleHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var item invmodel.MaintenanceSchedule
	if err := c.ShouldBindJSON(&item); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.UpdateSchedule(c.Request.Context(), id, &item); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *MaintenanceScheduleHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteSchedule(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusNoContent, nil)
}
