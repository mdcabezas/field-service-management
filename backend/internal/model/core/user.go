package core

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email" binding:"required,email,max=254"`
	Role      string    `json:"role" db:"role" binding:"required,max=50"`
	Name      string    `json:"name" db:"name" binding:"required,max=200"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserWithPassword struct {
	User
	PasswordHash *string `db:"password_hash"`
}
