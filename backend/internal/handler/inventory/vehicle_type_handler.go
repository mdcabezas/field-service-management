package inventoryhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	"localis-backend/internal/repository/postgres/inventory"
)

type VehicleTypeHandler struct {
	repo *postgresinventory.VehicleTypeRepo
}

func NewVehicleTypeHandler(repo *postgresinventory.VehicleTypeRepo) *VehicleTypeHandler {
	return &VehicleTypeHandler{repo: repo}
}

func (h *VehicleTypeHandler) List(c *gin.Context) {
	result, err := h.repo.List(c.Request.Context())
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, result)
}
