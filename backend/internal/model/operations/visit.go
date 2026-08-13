package operations

import (
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type Visit struct {
	ID                    uuid.UUID            `json:"id" db:"id" `
	PropertyID            *uuid.UUID           `json:"property_id,omitempty" db:"property_id"`
	PartnerID             *uuid.UUID           `json:"partner_id,omitempty" db:"partner_id"`
	PartnerOrderID        *string              `json:"partner_order_id,omitempty" db:"partner_order_id"`
	ParentVisitID         *uuid.UUID           `json:"parent_visit_id,omitempty" db:"parent_visit_id"`
	RouteID               *uuid.UUID           `json:"route_id,omitempty" db:"route_id"`
	DailyPlanID           *uuid.UUID           `json:"daily_plan_id,omitempty" db:"daily_plan_id"`
	Type                  uuid.UUID            `json:"type" db:"type"`
	Status                shared.VisitStatus   `json:"status" db:"status" binding:"required"`
	Priority              shared.VisitPriority `json:"priority" db:"priority" binding:"required"`
	Source                *shared.VisitSource  `json:"source,omitempty" db:"source"`
	SourceReference       *string              `json:"source_reference,omitempty" db:"source_reference"`
	BillingTo             shared.BillingTo     `json:"billing_to" db:"billing_to" binding:"required"`
	PartnerSupervisor     *string              `json:"partner_supervisor,omitempty" db:"partner_supervisor"`
	Result                *uuid.UUID           `json:"result,omitempty" db:"result"`
	RejectionReasonID     *uuid.UUID           `json:"rejection_reason_id,omitempty" db:"rejection_reason_id"`
	RejectionReasonDetail *string              `json:"rejection_reason_detail,omitempty" db:"rejection_reason_detail"`
	ScheduledAt           time.Time            `json:"scheduled_at" db:"scheduled_at" binding:"required"`
	StartedAt             *time.Time           `json:"started_at,omitempty" db:"started_at"`
	CompletedAt           *time.Time           `json:"completed_at,omitempty" db:"completed_at"`
	AlternativeLocation   *string              `json:"alternative_location,omitempty" db:"alternative_location"`
	Notes                 *string              `json:"notes,omitempty" db:"notes"`
	CreatedAt             time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at" db:"updated_at"`
}
