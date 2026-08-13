package operations

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type VisitEPPUsage struct {
	ID        uuid.UUID             `json:"id" db:"id" `
	VisitID   uuid.UUID             `json:"visit_id" db:"visit_id"`
	EPPID     uuid.UUID             `json:"epp_id" db:"epp_id"`
	Quantity  int                   `json:"quantity" db:"quantity" binding:"required"`
	Status    shared.EPPUsageStatus `json:"status" db:"status" binding:"required"`
	Notes     *string               `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time             `json:"created_at" db:"created_at"`
}
