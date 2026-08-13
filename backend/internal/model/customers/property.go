package customers

import (
	"time"

	"github.com/google/uuid"
)

type Property struct {
	ID                uuid.UUID  `json:"id" db:"id" `
	CustomerAddressID uuid.UUID  `json:"customer_address_id" db:"customer_address_id" binding:"required"`
	Name              *string    `json:"name,omitempty" db:"name"`
	Type              *uuid.UUID `json:"type,omitempty" db:"type"`
	Notes             *string    `json:"notes,omitempty" db:"notes"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
}
