package operations

import (
	"github.com/google/uuid"
	"localis-backend/internal/model/customers"
)

func newTechnicianFixture(active bool) *customers.Technician {
	return &customers.Technician{
		ID:       uuid.New(),
		Name:     "Tech Test",
		IsActive: active,
	}
}
