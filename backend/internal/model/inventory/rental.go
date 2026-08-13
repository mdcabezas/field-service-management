package inventory

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type Rental struct {
	ID              uuid.UUID         `json:"id" db:"id" `
	Type            shared.RentalType `json:"type" db:"type" binding:"required"`
	Supplier        *string           `json:"supplier,omitempty" db:"supplier"`
	ItemDescription string            `json:"item_description" db:"item_description" binding:"required"`
	DailyCost       *float64          `json:"daily_cost,omitempty" db:"daily_cost"`
	HourlyCost      *float64          `json:"hourly_cost,omitempty" db:"hourly_cost"`
	FixedCost       *float64          `json:"fixed_cost,omitempty" db:"fixed_cost"`
	Notes           *string           `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
}
