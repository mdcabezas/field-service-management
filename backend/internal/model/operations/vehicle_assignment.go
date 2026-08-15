package operations

import (
	"time"

	"github.com/google/uuid"
)

type VehicleAssignment struct {
	ID               uuid.UUID  `json:"id" db:"id" `
	VehicleID        uuid.UUID  `json:"vehicle_id" db:"vehicle_id"`
	DailyPlanID      *uuid.UUID `json:"daily_plan_id,omitempty" db:"daily_plan_id"`
	VisitID          *uuid.UUID `json:"visit_id,omitempty" db:"visit_id"`
	VehicleName      string     `json:"vehicle_name" db:"vehicle_name"`
	DailyPlanName    string     `json:"daily_plan_name" db:"daily_plan_name"`
	VisitLabel       string     `json:"visit_label" db:"visit_label"`
	DepartureTime    *time.Time `json:"departure_time,omitempty" db:"departure_time"`
	ReturnTime       *time.Time `json:"return_time,omitempty" db:"return_time"`
	DepartureMileage *int       `json:"departure_mileage,omitempty" db:"departure_mileage"`
	ReturnMileage    *int       `json:"return_mileage,omitempty" db:"return_mileage"`
	Notes            *string    `json:"notes,omitempty" db:"notes"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}
