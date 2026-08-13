package operations

import (
	"time"

	"github.com/google/uuid"
)

type VisitSLATracking struct {
	ID                  uuid.UUID  `json:"id" db:"id" `
	VisitID             uuid.UUID  `json:"visit_id" db:"visit_id"`
	SLAID               uuid.UUID  `json:"sla_id" db:"sla_id"`
	RequestedAt         *time.Time `json:"requested_at,omitempty" db:"requested_at"`
	RespondedAt         *time.Time `json:"responded_at,omitempty" db:"responded_at"`
	ResolvedAt          *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
	ResponseTimeHours   *float64   `json:"response_time_hours,omitempty" db:"response_time_hours"`
	ResolutionTimeHours *float64   `json:"resolution_time_hours,omitempty" db:"resolution_time_hours"`
	MeetsResponseSLA    *bool      `json:"meets_response_sla,omitempty" db:"meets_response_sla"`
	MeetsResolutionSLA  *bool      `json:"meets_resolution_sla,omitempty" db:"meets_resolution_sla"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
}
