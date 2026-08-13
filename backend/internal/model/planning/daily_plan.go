package planning

import (
	"time"

	"github.com/google/uuid"
)

type DailyPlan struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Date      time.Time `json:"date" db:"date" binding:"required"`
	Notes     *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
