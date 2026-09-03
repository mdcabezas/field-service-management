package inventory

import (
	"time"

	"github.com/google/uuid"
)

type ChecklistStage struct {
	ID         uuid.UUID `json:"id" db:"id"`
	TemplateID uuid.UUID `json:"template_id" db:"template_id"`
	Name       string    `json:"name" db:"name" binding:"required"`
	SortOrder  int       `json:"sort_order" db:"sort_order"`
	Required   bool      `json:"required" db:"required"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
