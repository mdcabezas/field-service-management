package operationshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	opsmodel "localis-backend/internal/model/operations"
	opssvc "localis-backend/internal/service/operations"
)

type VehicleAssignmentHandler struct {
	svc *opssvc.VehicleAssignmentService
}

func NewVehicleAssignmentHandler(svc *opssvc.VehicleAssignmentService) *VehicleAssignmentHandler {
	return &VehicleAssignmentHandler{svc: svc}
}

func (h *VehicleAssignmentHandler) GetByID(c *gin.Context) {
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

func (h *VehicleAssignmentHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	items, err := h.svc.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *VehicleAssignmentHandler) Create(c *gin.Context) {
	var body opsmodel.VehicleAssignment
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Assign(c.Request.Context(), &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, body)
}

func (h *VehicleAssignmentHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var body opsmodel.VehicleAssignment
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.Update(c.Request.Context(), id, &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, body)
}

func (h *VehicleAssignmentHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusNoContent, nil)
}
