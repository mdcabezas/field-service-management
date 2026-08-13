package customers

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type CustomerAddress struct {
	ID         uuid.UUID                  `json:"id" db:"id" `
	CustomerID uuid.UUID                  `json:"customer_id" db:"customer_id" binding:"required"`
	AddressID  uuid.UUID                  `json:"address_id" db:"address_id" binding:"required"`
	Name       *string                    `json:"name,omitempty" db:"name"`
	Type       shared.CustomerAddressType `json:"type" db:"type" binding:"required"`
	CreatedAt  time.Time                  `json:"created_at" db:"created_at"`
}
