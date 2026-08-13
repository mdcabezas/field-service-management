package operations

import (
	"time"

	"github.com/google/uuid"
)

type ReportImage struct {
	ID        uuid.UUID `json:"id" db:"id" `
	ReportID  uuid.UUID `json:"report_id" db:"report_id"`
	URL       string    `json:"url" db:"url" binding:"required"`
	Pages     int       `json:"pages" db:"pages" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
