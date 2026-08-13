package operations

import (
	"time"

	"github.com/google/uuid"
)

type VisitChecklistTool struct {
	ID              uuid.UUID `json:"id" db:"id" `
	VisitID         uuid.UUID `json:"visit_id" db:"visit_id"`
	ToolID          uuid.UUID `json:"tool_id" db:"tool_id"`
	PlannedQuantity int       `json:"planned_quantity" db:"planned_quantity" binding:"required"`
	Confirmed       bool      `json:"confirmed" db:"confirmed" binding:"required"`
	Notes           *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
