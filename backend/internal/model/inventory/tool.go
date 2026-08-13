package inventory

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type Tool struct {
	ID        uuid.UUID         `json:"id" db:"id" `
	Code      *string           `json:"code,omitempty" db:"code"`
	Name      string            `json:"name" db:"name" binding:"required"`
	Status    shared.ToolStatus `json:"status" db:"status" binding:"required"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
}
