package operations

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type VisitMeasurement struct {
	ID              uuid.UUID                `json:"id" db:"id" `
	VisitID         uuid.UUID                `json:"visit_id" db:"visit_id"`
	Type            uuid.UUID                `json:"type" db:"type"`
	Value           *float64                 `json:"value,omitempty" db:"value"`
	Unit            *string                  `json:"unit,omitempty" db:"unit"`
	Result          *shared.MeasurementResult `json:"result,omitempty" db:"result"`
	MeasuringDevice *string                  `json:"measuring_device,omitempty" db:"measuring_device"`
	Notes           *string                  `json:"notes,omitempty" db:"notes"`
	PhotoURL        *string                  `json:"photo_url,omitempty" db:"photo_url"`
	Lat             *float64                 `json:"lat,omitempty" db:"lat"`
	Lng             *float64                 `json:"lng,omitempty" db:"lng"`
	CreatedAt       time.Time                `json:"created_at" db:"created_at"`
}
