package operations

import (
	"time"

	"github.com/google/uuid"
)

type VisitAssignment struct {
	ID        uuid.UUID `json:"id" db:"id" `
	VisitID   uuid.UUID `json:"visit_id" db:"visit_id"`
	TechID    uuid.UUID `json:"tech_id" db:"tech_id"`
	RoleID    uuid.UUID `json:"role_id" db:"role_id"`
	TechName  string    `json:"tech_name" db:"tech_name"`
	RoleName  string    `json:"role_name" db:"role_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
