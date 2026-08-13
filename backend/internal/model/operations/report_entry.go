package operations

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ReportEntry struct {
	ID        uuid.UUID       `json:"id" db:"id" `
	ReportID  uuid.UUID       `json:"report_id" db:"report_id"`
	DataJSON  json.RawMessage `json:"data_json" db:"data_json"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}
