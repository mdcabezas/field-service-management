package operations

import (
	"time"

	"github.com/google/uuid"
)

type VisitMaterialUsage struct {
	ID           uuid.UUID `json:"id" db:"id"`
	VisitID      uuid.UUID `json:"visit_id" db:"visit_id"`
	MaterialID   uuid.UUID `json:"material_id" db:"material_id"`
	MaterialName *string   `json:"material_name,omitempty" db:"material_name"`
	Quantity     float64   `json:"quantity" db:"quantity" binding:"required"`
	Notes        *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
