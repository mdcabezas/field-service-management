package core

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type PlanAuditLog struct {
	ID         uuid.UUID              `json:"id" db:"id"`
	EntityType shared.AuditEntityType `json:"entity_type" db:"entity_type"`
	EntityID   uuid.UUID              `json:"entity_id" db:"entity_id"`
	Action     shared.AuditAction     `json:"action" db:"action"`
	OldData    json.RawMessage        `json:"old_data,omitempty" db:"old_data"`
	NewData    json.RawMessage        `json:"new_data,omitempty" db:"new_data"`
	UserID     *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	Reason     *string                `json:"reason,omitempty" db:"reason"`
	Timestamp  time.Time              `json:"timestamp" db:"timestamp"`
}
