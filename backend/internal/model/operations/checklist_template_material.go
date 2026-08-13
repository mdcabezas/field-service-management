package operations

import (
	"time"

	"github.com/google/uuid"
)

type ChecklistTemplateMaterial struct {
	ID              uuid.UUID `json:"id" db:"id" `
	TemplateID      uuid.UUID `json:"template_id" db:"template_id"`
	MaterialID      uuid.UUID `json:"material_id" db:"material_id"`
	DefaultQuantity int       `json:"default_quantity" db:"default_quantity" binding:"required"`
	Required        bool      `json:"required" db:"required" binding:"required"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
