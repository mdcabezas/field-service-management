package partners

import (
	"time"

	"github.com/google/uuid"
)

type SLA struct {
	ID               uuid.UUID  `json:"id" db:"id" `
	PartnerID        uuid.UUID  `json:"partner_id" db:"partner_id" binding:"required"`
	Name             string     `json:"name" db:"name" binding:"required"`
	Description      *string    `json:"description,omitempty" db:"description"`
	WorkType         *uuid.UUID `json:"work_type,omitempty" db:"work_type"`
	ResponseHours    *int       `json:"response_hours,omitempty" db:"response_hours"`
	ResolutionHours  *int       `json:"resolution_hours,omitempty" db:"resolution_hours"`
	ComplianceTarget *float64   `json:"compliance_target,omitempty" db:"compliance_target"`
	Active           bool       `json:"active" db:"active" binding:"required"`
	ValidFrom        time.Time  `json:"valid_from" db:"valid_from" binding:"required"`
	ValidUntil       *time.Time `json:"valid_until,omitempty" db:"valid_until"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}
