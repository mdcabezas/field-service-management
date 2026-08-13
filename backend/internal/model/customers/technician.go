package customers

import (
	"time"

	"github.com/google/uuid"
)

type Technician struct {
	ID          uuid.UUID `json:"id" db:"id" `
	UserID      *string   `json:"user_id,omitempty" db:"user_id"`
	Name        string    `json:"name" db:"name" binding:"required"`
	IsActive    bool      `json:"is_active" db:"is_active" binding:"required"`
	Specialties []string  `json:"specialties,omitempty" db:"specialties"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
