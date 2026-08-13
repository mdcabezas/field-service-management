package planning

import (
	"time"

	"github.com/google/uuid"
)

type DailyPlanAssignment struct {
	ID          uuid.UUID `json:"id" db:"id" `
	DailyPlanID uuid.UUID `json:"daily_plan_id" db:"daily_plan_id" binding:"required"`
	TechID      uuid.UUID `json:"tech_id" db:"tech_id" binding:"required"`
	RoleID      uuid.UUID `json:"role_id" db:"role_id" binding:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
