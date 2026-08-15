package operations

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type VisitReport struct {
	ID                   uuid.UUID           `json:"id" db:"id" `
	VisitID              uuid.UUID           `json:"visit_id" db:"visit_id"`
	ReportTemplateID     *uuid.UUID          `json:"report_template_id,omitempty" db:"report_template_id"`
	Source               shared.ReportSource `json:"source" db:"source" binding:"required"`
	RecordedAt           time.Time           `json:"recorded_at" db:"recorded_at" binding:"required"`
	CreatedAt            time.Time           `json:"created_at" db:"created_at"`
	ReportTemplateName   *string             `json:"report_template_name,omitempty" db:"report_template_name"`
	VisitLabel           *string             `json:"visit_label,omitempty" db:"visit_label"`
}

type CreateReportRequest struct {
	VisitID          uuid.UUID            `json:"visit_id" binding:"required"`
	ReportTemplateID uuid.UUID            `json:"report_template_id" binding:"required"`
	Source           shared.ReportSource  `json:"source,omitempty"`
	RecordedAt       *time.Time           `json:"recorded_at,omitempty"`
}
