package customers

import (
	"time"

	"github.com/google/uuid"
)

type CustomerAddressPartner struct {
	ID                uuid.UUID  `json:"id" db:"id" `
	CustomerAddressID uuid.UUID  `json:"customer_address_id" db:"customer_address_id" binding:"required"`
	PartnerID         uuid.UUID  `json:"partner_id" db:"partner_id" binding:"required"`
	FromDate          time.Time  `json:"from_date" db:"from_date" binding:"required"`
	ToDate            *time.Time `json:"to_date,omitempty" db:"to_date"`
	Notes             *string    `json:"notes,omitempty" db:"notes"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}
