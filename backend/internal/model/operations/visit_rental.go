package operations

import (
	"time"

	"github.com/google/uuid"
)

type VisitRental struct {
	ID        uuid.UUID  `json:"id" db:"id" `
	VisitID   uuid.UUID  `json:"visit_id" db:"visit_id"`
	RentalID  uuid.UUID  `json:"rental_id" db:"rental_id"`
	StartTime time.Time  `json:"start_time" db:"start_time" binding:"required"`
	EndTime   *time.Time `json:"end_time,omitempty" db:"end_time"`
	Hours     *float64   `json:"hours,omitempty" db:"hours"`
	TotalCost float64    `json:"total_cost" db:"total_cost" binding:"required"`
	Reason    *string    `json:"reason,omitempty" db:"reason"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}
