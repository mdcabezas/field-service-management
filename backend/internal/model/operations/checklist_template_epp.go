package operations

import (
	"time"

	"github.com/google/uuid"
)

type ChecklistTemplateEPP struct {
	ID              uuid.UUID `json:"id" db:"id" `
	TemplateID      uuid.UUID `json:"template_id" db:"template_id"`
	EPPID           uuid.UUID `json:"epp_id" db:"epp_id"`
	DefaultQuantity int       `json:"default_quantity" db:"default_quantity" binding:"required"`
	Required        bool      `json:"required" db:"required" binding:"required"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
