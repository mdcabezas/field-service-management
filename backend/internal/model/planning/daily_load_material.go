package planning

import (
	"time"

	"github.com/google/uuid"
)

type DailyLoadMaterial struct {
	ID               uuid.UUID `json:"id" db:"id" `
	DailyPlanID      uuid.UUID `json:"daily_plan_id" db:"daily_plan_id" binding:"required"`
	MaterialID       uuid.UUID `json:"material_id" db:"material_id" binding:"required"`
	LoadedQuantity   float64   `json:"loaded_quantity" db:"loaded_quantity" binding:"required"`
	ReturnedQuantity float64   `json:"returned_quantity" db:"returned_quantity" binding:"required"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
