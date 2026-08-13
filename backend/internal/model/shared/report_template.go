package shared

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ReportTemplate struct {
	ID         uuid.UUID       `json:"id" db:"id" `
	Name       string          `json:"name" db:"name" binding:"required"`
	FieldsJSON json.RawMessage `json:"fields_json" db:"fields_json" binding:"required"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
}
