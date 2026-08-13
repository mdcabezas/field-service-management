package partners

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type Partner struct {
	ID        uuid.UUID            `json:"id" db:"id" `
	Name      string               `json:"name" db:"name" binding:"required"`
	TaxID     *string              `json:"tax_id,omitempty" db:"tax_id"`
	Status    shared.PartnerStatus `json:"status" db:"status" binding:"required"`
	CreatedAt time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" db:"updated_at"`
}
