package notifications

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type Notification struct {
	ID           uuid.UUID                  `json:"id" db:"id" `
	Type         shared.NotificationType    `json:"type" db:"type" binding:"required"`
	Channel      shared.NotificationChannel `json:"channel" db:"channel" binding:"required"`
	Recipient    string                     `json:"recipient" db:"recipient" binding:"required"`
	Subject      *string                    `json:"subject,omitempty" db:"subject"`
	Body         *string                    `json:"body,omitempty" db:"body"`
	Status       shared.NotificationStatus  `json:"status" db:"status" binding:"required"`
	EntityType   *string                    `json:"entity_type,omitempty" db:"entity_type"`
	EntityID     *uuid.UUID                 `json:"entity_id,omitempty" db:"entity_id"`
	ScheduledFor *time.Time                 `json:"scheduled_for,omitempty" db:"scheduled_for"`
	SentAt       *time.Time                 `json:"sent_at,omitempty" db:"sent_at"`
	ReadAt       *time.Time                 `json:"read_at,omitempty" db:"read_at"`
	Attempts     int                        `json:"attempts" db:"attempts"`
	Error        *string                    `json:"error,omitempty" db:"error"`
	CreatedAt    time.Time                  `json:"created_at" db:"created_at"`
}
