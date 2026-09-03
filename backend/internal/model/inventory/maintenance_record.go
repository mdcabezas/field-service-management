package inventory

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type MaintenanceRecord struct {
	ID              uuid.UUID              `json:"id" db:"id" `
	Type            shared.MaintenanceType `json:"type" db:"type" binding:"required"`
	ReferenceID     uuid.UUID              `json:"reference_id" db:"reference_id" binding:"required"`
	ReferenceDisplay string                `json:"reference_display,omitempty" db:"reference_display"`
	Date            time.Time              `json:"date" db:"date" binding:"required"`
	Cost            *float64               `json:"cost,omitempty" db:"cost"`
	Supplier        *string                `json:"supplier,omitempty" db:"supplier"`
	Description     *string                `json:"description,omitempty" db:"description"`
	NextDate        *time.Time             `json:"next_date,omitempty" db:"next_date"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}
