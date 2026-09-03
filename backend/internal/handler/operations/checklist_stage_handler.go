package operationshandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/inventory"
)

type ChecklistStageHandler struct {
	stages map[string][]inventory.ChecklistStage
}

func NewChecklistStageHandler() *ChecklistStageHandler {
	return &ChecklistStageHandler{
		stages: make(map[string][]inventory.ChecklistStage),
	}
}

func (h *ChecklistStageHandler) ListByTemplate(c *gin.Context) {
	templateID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid template id")
		return
	}

	stages := h.stages[templateID.String()]
	if stages == nil {
		stages = []inventory.ChecklistStage{}
	}

	handler.RespondJSON(c, http.StatusOK, stages)
}

func (h *ChecklistStageHandler) Create(c *gin.Context) {
	var body struct {
		TemplateID uuid.UUID `json:"template_id" binding:"required"`
		Name       string    `json:"name" binding:"required"`
		SortOrder  int       `json:"sort_order"`
		Required   bool      `json:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	stage := inventory.ChecklistStage{
		ID:         uuid.New(),
		TemplateID: body.TemplateID,
		Name:       body.Name,
		SortOrder:  body.SortOrder,
		Required:   body.Required,
		CreatedAt:  time.Now(),
	}

	h.stages[body.TemplateID.String()] = append(h.stages[body.TemplateID.String()], stage)

	handler.RespondJSON(c, http.StatusCreated, stage)
}

func (h *ChecklistStageHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Name      *string `json:"name"`
		SortOrder *int    `json:"sort_order"`
		Required  *bool   `json:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	for templateID, stages := range h.stages {
		for i, s := range stages {
			if s.ID == id {
				if body.Name != nil {
					h.stages[templateID][i].Name = *body.Name
				}
				if body.SortOrder != nil {
					h.stages[templateID][i].SortOrder = *body.SortOrder
				}
				if body.Required != nil {
					h.stages[templateID][i].Required = *body.Required
				}
				handler.RespondJSON(c, http.StatusOK, h.stages[templateID][i])
				return
			}
		}
	}

	handler.RespondError(c, http.StatusNotFound, "stage not found")
}

func (h *ChecklistStageHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	for templateID, stages := range h.stages {
		for i, s := range stages {
			if s.ID == id {
				h.stages[templateID] = append(stages[:i], stages[i+1:]...)
				c.Status(http.StatusNoContent)
				return
			}
		}
	}

	handler.RespondError(c, http.StatusNotFound, "stage not found")
}
