package operationshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	operationssvc "localis-backend/internal/service/operations"
)

type VisitHandler struct {
	svc *operationssvc.VisitService
}

func NewVisitHandler(svc *operationssvc.VisitService) *VisitHandler {
	return &VisitHandler{svc: svc}
}

func (h *VisitHandler) GetByID(c *gin.Context) {
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

func (h *VisitHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	items, err := h.svc.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *VisitHandler) Create(c *gin.Context) {
	var body operations.Visit
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.CreateVisit(c.Request.Context(), &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, body)
}

func (h *VisitHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var body operations.Visit
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := c.Request.Context()

	switch body.Status {
	case shared.VisitStatusEnRoute:
		if err := h.svc.StartVisit(ctx, id, nil); err != nil {
			handler.HandleServiceError(c, err)
			return
		}
	case shared.VisitStatusInProgress:
		if err := h.svc.BeginWork(ctx, id); err != nil {
			handler.HandleServiceError(c, err)
			return
		}
	case shared.VisitStatusCompleted:
		if err := h.svc.CompleteVisit(ctx, id); err != nil {
			handler.HandleServiceError(c, err)
			return
		}
	case shared.VisitStatusCancelled:
		if err := h.svc.CancelVisit(ctx, id, body.Notes); err != nil {
			handler.HandleServiceError(c, err)
			return
		}
	default:
		handler.RespondError(c, http.StatusBadRequest, "unsupported status transition")
		return
	}

	handler.RespondJSON(c, http.StatusOK, gin.H{"id": id, "status": body.Status})
}

func (h *VisitHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *VisitHandler) Reopen(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "reason is required")
		return
	}

	if err := h.svc.ReopenVisit(c.Request.Context(), id, body.Reason); err != nil {
		handler.HandleServiceError(c, err)
		return
	}

	handler.RespondJSON(c, http.StatusOK, gin.H{"id": id, "status": "in_progress"})
}
