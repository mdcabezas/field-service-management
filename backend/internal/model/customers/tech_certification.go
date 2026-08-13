package customers

import (
	"time"

	"github.com/google/uuid"
)

type TechCertification struct {
	ID         uuid.UUID  `json:"id" db:"id" `
	TechID     uuid.UUID  `json:"tech_id" db:"tech_id" binding:"required"`
	CertID     *uuid.UUID `json:"cert_id,omitempty" db:"cert_id"`
	Number     *string    `json:"number,omitempty" db:"number"`
	Issuer     *string    `json:"issuer,omitempty" db:"issuer"`
	IssueDate  *time.Time `json:"issue_date,omitempty" db:"issue_date"`
	ExpiryDate *time.Time `json:"expiry_date,omitempty" db:"expiry_date"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}
