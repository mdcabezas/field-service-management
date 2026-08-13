package inventory

import (
	"time"

	"github.com/google/uuid"
)

type VehicleType struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}
