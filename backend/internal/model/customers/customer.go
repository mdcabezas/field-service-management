package customers

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `json:"id" db:"id" `
	Name      string    `json:"name" db:"name" binding:"required"`
	Phone     *string   `json:"phone,omitempty" db:"phone"`
	Email     *string   `json:"email,omitempty" db:"email"`
	TaxID     *string   `json:"tax_id,omitempty" db:"tax_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
