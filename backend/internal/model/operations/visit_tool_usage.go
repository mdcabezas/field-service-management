package operations

import (
	"time"

	"github.com/google/uuid"
)

type VisitToolUsage struct {
	ID        uuid.UUID `json:"id" db:"id" `
	VisitID   uuid.UUID `json:"visit_id" db:"visit_id"`
	ToolID    uuid.UUID `json:"tool_id" db:"tool_id"`
	Notes     *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
