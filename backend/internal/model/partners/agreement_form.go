package partners

import (
	"time"

	"github.com/google/uuid"
)

type PartnerAgreementForm struct {
	ID             uuid.UUID  `json:"id" db:"id" `
	AgreementID    uuid.UUID  `json:"agreement_id" db:"agreement_id" binding:"required"`
	WorkType       string     `json:"work_type" db:"work_type" binding:"required"`
	FormTemplateID *uuid.UUID `json:"form_template_id,omitempty" db:"form_template_id"`
	Quantity       int        `json:"quantity" db:"quantity" binding:"required"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}
