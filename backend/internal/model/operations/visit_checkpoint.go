package operations

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"localis-backend/internal/model/shared"
)

type VisitCheckpoint struct {
	ID         uuid.UUID             `json:"id" db:"id" `
	VisitID    uuid.UUID             `json:"visit_id" db:"visit_id"`
	Type       shared.CheckpointType `json:"type" db:"type" binding:"required"`
	Geom       orb.Point             `json:"geom,omitempty" db:"geom" binding:"required"`
	Timestamp  time.Time             `json:"timestamp" db:"timestamp" binding:"required"`
	DeviceInfo *string               `json:"device_info,omitempty" db:"device_info"`
	CreatedAt  time.Time             `json:"created_at" db:"created_at"`
}
