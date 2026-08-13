package planninghandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	"localis-backend/internal/repository/postgres/planning"
)

type RouteTypeHandler struct {
	repo *postgresplanning.RouteTypeRepo
}

func NewRouteTypeHandler(repo *postgresplanning.RouteTypeRepo) *RouteTypeHandler {
	return &RouteTypeHandler{repo: repo}
}

func (h *RouteTypeHandler) List(c *gin.Context) {
	result, err := h.repo.List(c.Request.Context())
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, result)
}
