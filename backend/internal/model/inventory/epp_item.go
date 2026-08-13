package inventory

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type EPPItem struct {
	ID        uuid.UUID           `json:"id" db:"id" `
	Name      string              `json:"name" db:"name" binding:"required"`
	Type      shared.EPPType      `json:"type" db:"type" binding:"required"`
	Lifecycle shared.EPPLifecycle `json:"lifecycle" db:"lifecycle" binding:"required"`
	CreatedAt time.Time           `json:"created_at" db:"created_at"`
}
