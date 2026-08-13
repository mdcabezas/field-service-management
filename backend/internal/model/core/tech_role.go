package core

import (
	"time"

	"github.com/google/uuid"
)

type TechRole struct {
	ID        uuid.UUID `json:"id" db:"id" `
	Code      string    `json:"code" db:"code" binding:"required"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Active    bool      `json:"active" db:"active" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
