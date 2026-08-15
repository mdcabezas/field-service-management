package operations

import (
	"github.com/google/uuid"
)

type VisitType struct {
	ID     uuid.UUID `json:"id" db:"id"`
	Code   string    `json:"code" db:"code" binding:"required"`
	Name   string    `json:"name" db:"name" binding:"required"`
	Active bool      `json:"active" db:"active"`
}