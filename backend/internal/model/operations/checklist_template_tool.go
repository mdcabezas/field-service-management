package operations

import (
	"time"

	"github.com/google/uuid"
)

type ChecklistTemplateTool struct {
	ID              uuid.UUID `json:"id" db:"id" `
	TemplateID      uuid.UUID `json:"template_id" db:"template_id"`
	ToolID          uuid.UUID `json:"tool_id" db:"tool_id"`
	DefaultQuantity int       `json:"default_quantity" db:"default_quantity" binding:"required"`
	Required        bool      `json:"required" db:"required" binding:"required"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
