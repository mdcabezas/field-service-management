package operationshandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/repository"
)

type VisitCheckpointHandler struct {
	repo repository.VisitCheckpointRepository
}

func NewVisitCheckpointHandler(repo repository.VisitCheckpointRepository) *VisitCheckpointHandler {
	return &VisitCheckpointHandler{repo: repo}
}

func (h *VisitCheckpointHandler) GetByID(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *VisitCheckpointHandler) ListByVisit(c *gin.Context) {
	visitID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	items, err := h.repo.ListByVisit(c.Request.Context(), visitID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *VisitCheckpointHandler) Create(c *gin.Context) {
	var body operations.VisitCheckpoint
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	body.ID = uuid.New()
	body.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, body)
}

func (h *VisitCheckpointHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusNoContent, nil)
}
