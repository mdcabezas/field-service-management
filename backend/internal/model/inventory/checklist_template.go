package inventory

import (
	"time"

	"github.com/google/uuid"
)

type ChecklistTemplate struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" binding:"required"`
	Description *string    `json:"description,omitempty" db:"description"`
	VisitType   *string    `json:"visit_type,omitempty" db:"visit_type"`
	WorkType    *uuid.UUID `json:"work_type,omitempty" db:"work_type"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}
