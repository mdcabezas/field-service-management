package inventory

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type Vehicle struct {
	ID           uuid.UUID            `json:"id" db:"id" `
	Type         *uuid.UUID           `json:"type,omitempty" db:"type"`
	LicensePlate *string              `json:"license_plate,omitempty" db:"license_plate"`
	Name         string               `json:"name" db:"name" binding:"required"`
	Brand        *string              `json:"brand,omitempty" db:"brand"`
	Model        *string              `json:"model,omitempty" db:"model"`
	Year         *int                 `json:"year,omitempty" db:"year"`
	Status       shared.VehicleStatus `json:"status" db:"status" binding:"required"`
	Capacity     *string              `json:"capacity,omitempty" db:"capacity"`
	CreatedAt    time.Time            `json:"created_at" db:"created_at"`
}
