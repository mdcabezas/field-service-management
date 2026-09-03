package inventory

import (
	"time"

	"github.com/google/uuid"
)

type ChecklistTemplate struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	Name            string     `json:"name" db:"name" binding:"required"`
	Description     *string    `json:"description,omitempty" db:"description"`
	VisitType        *uuid.UUID `json:"visit_type,omitempty" db:"visit_type"`
	VisitTypeDisplay *string    `json:"visit_type_display,omitempty" db:"visit_type_display"`
	WorkType         *uuid.UUID `json:"work_type,omitempty" db:"work_type"`
	WorkTypeDisplay  *string    `json:"work_type_display,omitempty" db:"work_type_display"`
	PartnerID       *uuid.UUID `json:"partner_id,omitempty" db:"partner_id"`
	StagesEnabled   bool       `json:"stages_enabled" db:"stages_enabled"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}
