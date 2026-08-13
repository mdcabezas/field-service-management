package notifications

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type NotificationTemplate struct {
	ID        uuid.UUID                  `json:"id" db:"id" `
	Type      shared.NotificationType    `json:"type" db:"type" binding:"required"`
	Channel   shared.NotificationChannel `json:"channel" db:"channel" binding:"required"`
	Subject   *string                    `json:"subject,omitempty" db:"subject"`
	Body      *string                    `json:"body,omitempty" db:"body"`
	Active    bool                       `json:"active" db:"active" binding:"required"`
	CreatedAt time.Time                  `json:"created_at" db:"created_at"`
}
