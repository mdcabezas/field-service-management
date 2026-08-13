package partners

import (
	"time"

	"github.com/google/uuid"
)

type PartnerContact struct {
	ID        uuid.UUID `json:"id" db:"id" `
	PartnerID uuid.UUID `json:"partner_id" db:"partner_id" binding:"required"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Position  *string   `json:"position,omitempty" db:"position"`
	Phone     *string   `json:"phone,omitempty" db:"phone"`
	Email     *string   `json:"email,omitempty" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
