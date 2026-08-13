package geocoding

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type Address struct {
	ID                 uuid.UUID `json:"id" db:"id" `
	Street             string    `json:"street" db:"street" binding:"required"`
	Number             *string   `json:"number,omitempty" db:"number"`
	Apartment          *string   `json:"apartment,omitempty" db:"apartment"`
	Neighborhood       *string   `json:"neighborhood,omitempty" db:"neighborhood"`
	City               *string   `json:"city,omitempty" db:"city"`
	Region             *string   `json:"region,omitempty" db:"region"`
	LocationReferences *string   `json:"location_references,omitempty" db:"location_references"`
	PostalCode         *string   `json:"postal_code,omitempty" db:"postal_code"`
	Geom               orb.Point `json:"geom,omitempty" db:"geom"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}
