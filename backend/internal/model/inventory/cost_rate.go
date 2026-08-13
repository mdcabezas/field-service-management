package inventory

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type CostRate struct {
	ID            uuid.UUID       `json:"id" db:"id" `
	Type          shared.CostType `json:"type" db:"type" binding:"required"`
	ReferenceID   *uuid.UUID      `json:"reference_id,omitempty" db:"reference_id"`
	ReferenceType *string         `json:"reference_type,omitempty" db:"reference_type"`
	Value         float64         `json:"value" db:"value" binding:"required"`
	Unit          shared.CostUnit `json:"unit" db:"unit" binding:"required"`
	ValidFrom     time.Time       `json:"valid_from" db:"valid_from" binding:"required"`
	ValidUntil    *time.Time      `json:"valid_until,omitempty" db:"valid_until"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}
