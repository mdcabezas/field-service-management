package operationshandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	operationsmodel "localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	operationssvc "localis-backend/internal/service/operations"
)

type VisitReportHandler struct {
	svc *operationssvc.ReportService
}

func NewVisitReportHandler(svc *operationssvc.ReportService) *VisitReportHandler {
	return &VisitReportHandler{svc: svc}
}

func (h *VisitReportHandler) GetByID(c *gin.Context) {
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

func (h *VisitReportHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	items, err := h.svc.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *VisitReportHandler) ListByVisit(c *gin.Context) {
	visitID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	items, err := h.svc.ListByVisit(c.Request.Context(), visitID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *VisitReportHandler) Create(c *gin.Context) {
	var body operationsmodel.CreateReportRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	source := body.Source
	if source == "" {
		source = shared.ReportSourceSystem
	}
	recordedAt := time.Now()
	if body.RecordedAt != nil {
		recordedAt = *body.RecordedAt
	}
	report, err := h.svc.GenerateReport(c.Request.Context(), body.VisitID, body.ReportTemplateID, source, recordedAt)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, report)
}

func (h *VisitReportHandler) Delete(c *gin.Context) {
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
