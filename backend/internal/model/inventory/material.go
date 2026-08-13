package inventory

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID        uuid.UUID `json:"id" db:"id" `
	SKU       *string   `json:"sku,omitempty" db:"sku"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Unit      string    `json:"unit" db:"unit" binding:"required"`
	UnitCost  *float64  `json:"unit_cost,omitempty" db:"unit_cost"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
