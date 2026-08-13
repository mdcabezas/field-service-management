package partners

import (
	"time"

	"github.com/google/uuid"
)

type PartnerAgreementDoc struct {
	ID          uuid.UUID `json:"id" db:"id"`
	AgreementID uuid.UUID `json:"agreement_id" db:"agreement_id"`
	WorkType    string    `json:"work_type" db:"work_type"`
	DocName     string    `json:"doc_name" db:"doc_name"`
	Required    bool      `json:"required" db:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
