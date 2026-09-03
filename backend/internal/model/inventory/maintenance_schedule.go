package inventory

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type MaintenanceSchedule struct {
	ID              uuid.UUID              `json:"id" db:"id" `
	Type            shared.MaintenanceType `json:"type" db:"type" binding:"required"`
	ReferenceID     uuid.UUID              `json:"reference_id" db:"reference_id" binding:"required"`
	ReferenceDisplay string                `json:"reference_display,omitempty" db:"reference_display"`
	FrequencyKM     *int                   `json:"frequency_km,omitempty" db:"frequency_km"`
	FrequencyDays   *int                   `json:"frequency_days,omitempty" db:"frequency_days"`
	LastServiceDate *time.Time             `json:"last_service_date,omitempty" db:"last_service_date"`
	LastServiceKM   *int                   `json:"last_service_km,omitempty" db:"last_service_km"`
	Active          bool                   `json:"active" db:"active" binding:"required"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
}
