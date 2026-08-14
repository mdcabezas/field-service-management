package planning

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/shared"
)

type Route struct {
	ID          uuid.UUID          `json:"id" db:"id" `
	Type        string             `json:"type" db:"type"`
	TypeName    string             `json:"type_name" db:"type_name"`
	Date        time.Time          `json:"date" db:"date" binding:"required"`
	DailyPlanID *uuid.UUID         `json:"daily_plan_id,omitempty" db:"daily_plan_id"`
	Notes       *string            `json:"notes,omitempty" db:"notes"`
	Status      shared.RouteStatus `json:"status" db:"status" binding:"required"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
}

func (r *Route) UnmarshalJSON(data []byte) error {
	type Alias Route
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if aux.Date != "" {
		s := strings.TrimSpace(aux.Date)
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			t, err = time.Parse(time.RFC3339, s)
			if err != nil {
				return err
			}
		}
		r.Date = t
	}
	return nil
}
