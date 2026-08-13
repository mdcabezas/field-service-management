package partners

import (
	"time"

	"github.com/google/uuid"
)

type PartnerAgreement struct {
	ID          uuid.UUID `json:"id" db:"id" `
	PartnerID   uuid.UUID `json:"partner_id" db:"partner_id" binding:"required"`
	ServiceType uuid.UUID `json:"service_type" db:"service_type" binding:"required"`
	Rate        *float64  `json:"rate,omitempty" db:"rate"`
	Active      bool      `json:"active" db:"active" binding:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
