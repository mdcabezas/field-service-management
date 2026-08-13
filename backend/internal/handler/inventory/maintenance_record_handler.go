package inventoryhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	invmodel "localis-backend/internal/model/inventory"
	invcsv "localis-backend/internal/service/inventory"
)

type MaintenanceRecordHandler struct {
	svc *invcsv.MaintenanceService
}

func NewMaintenanceRecordHandler(svc *invcsv.MaintenanceService) *MaintenanceRecordHandler {
	return &MaintenanceRecordHandler{svc: svc}
}

func (h *MaintenanceRecordHandler) GetByID(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *MaintenanceRecordHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	result, err := h.svc.ListRecords(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, result)
}

func (h *MaintenanceRecordHandler) Create(c *gin.Context) {
	var item invmodel.MaintenanceRecord
	if err := c.ShouldBindJSON(&item); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.RecordMaintenance(c.Request.Context(), &item); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, item)
}

func (h *MaintenanceRecordHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var item invmodel.MaintenanceRecord
	if err := c.ShouldBindJSON(&item); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.UpdateRecord(c.Request.Context(), id, &item); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *MaintenanceRecordHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteRecord(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusNoContent, nil)
}
