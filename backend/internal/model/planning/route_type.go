package planning

import (
	"github.com/google/uuid"
)

type RouteType struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Code   string    `json:"code" db:"code"`
	Name   string    `json:"name" db:"name"`
	Active bool      `json:"active" db:"active"`
}
