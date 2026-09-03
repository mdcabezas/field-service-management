package operations

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"localis-backend/internal/model/shared"
)

type VisitPhoto struct {
	ID            uuid.UUID         `json:"id" db:"id" `
	VisitID       uuid.UUID         `json:"visit_id" db:"visit_id"`
	CheckpointID  *uuid.UUID        `json:"checkpoint_id,omitempty" db:"checkpoint_id"`
	URL           *string           `json:"url" db:"url"`
	ThumbnailURL  *string           `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	Geom          orb.Point         `json:"geom,omitempty" db:"geom" binding:"required"`
	Timestamp     time.Time         `json:"timestamp" db:"timestamp" binding:"required"`
	Stage         *shared.PhotoStage `json:"stage,omitempty" db:"stage"`
	FindingType   *uuid.UUID        `json:"finding_type,omitempty" db:"finding_type"`
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
}
