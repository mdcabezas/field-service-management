package operations

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/service"
	"localis-backend/internal/service/core"
	"localis-backend/internal/testutil/mocks"
)

type vehicleSvcMocks struct {
	vaRepo    *mocks.VehicleAssignmentRepository
	vRepo     *mocks.VehicleRepository
	auditRepo *mocks.CorePlanAuditLogRepository
	svc       *VehicleAssignmentService
}

func newVehicleService(t *testing.T) *vehicleSvcMocks {
	t.Helper()
	vaRepo := mocks.NewVehicleAssignmentRepository(t)
	vRepo := mocks.NewVehicleRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	svc := NewVehicleAssignmentService(vaRepo, vRepo, core.NewAuditService(auditRepo))
	return &vehicleSvcMocks{vaRepo: vaRepo, vRepo: vRepo, auditRepo: auditRepo, svc: svc}
}

func TestVehicleAssign_Success(t *testing.T) {
	ctx := context.Background()
	m := newVehicleService(t)

	va := &operations.VehicleAssignment{VehicleID: uuid.New()}
	vehicle := &inventory.Vehicle{ID: va.VehicleID, Status: shared.VehicleStatusAvailable, Name: "Truck 1"}

	m.vRepo.On("GetByID", ctx, va.VehicleID).Return(vehicle, nil).Once()
	m.vaRepo.On("Create", ctx, va).Return(nil).Once()
	m.vRepo.On("Update", ctx, va.VehicleID, vehicle).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.Assign(ctx, va)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, va.ID)
	require.NotNil(t, va.DepartureTime)
	require.Equal(t, shared.VehicleStatusInUse, vehicle.Status)
}

func TestVehicleAssign_NotAvailable(t *testing.T) {
	ctx := context.Background()
	m := newVehicleService(t)

	va := &operations.VehicleAssignment{VehicleID: uuid.New()}
	vehicle := &inventory.Vehicle{ID: va.VehicleID, Status: shared.VehicleStatusInUse, Name: "Truck 1"}
	m.vRepo.On("GetByID", ctx, va.VehicleID).Return(vehicle, nil).Once()

	err := m.svc.Assign(ctx, va)
	var ce *service.ConflictError
	require.ErrorAs(t, err, &ce)
}

func TestVehicleAssign_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newVehicleService(t)

	va := &operations.VehicleAssignment{VehicleID: uuid.New()}
	m.vRepo.On("GetByID", ctx, va.VehicleID).Return(nil, &service.NotFoundError{Resource: "vehicle", ID: va.VehicleID.String()}).Once()

	err := m.svc.Assign(ctx, va)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestVehicleReturn_Success(t *testing.T) {
	ctx := context.Background()
	m := newVehicleService(t)

	id := uuid.New()
	va := &operations.VehicleAssignment{ID: id, VehicleID: uuid.New()}
	vehicle := &inventory.Vehicle{ID: va.VehicleID, Status: shared.VehicleStatusInUse, Name: "Truck 1"}

	m.vaRepo.On("GetByID", ctx, id).Return(va, nil).Once()
	m.vaRepo.On("Update", ctx, id, va).Return(nil).Once()
	m.vRepo.On("GetByID", ctx, va.VehicleID).Return(vehicle, nil).Once()
	m.vRepo.On("Update", ctx, va.VehicleID, vehicle).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.Return(ctx, id, 5000)
	require.NoError(t, err)
	require.NotNil(t, va.ReturnTime)
	require.Equal(t, 5000, *va.ReturnMileage)
	require.Equal(t, shared.VehicleStatusAvailable, vehicle.Status)
}

func TestVehicleReturn_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newVehicleService(t)

	id := uuid.New()
	m.vaRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "vehicle_assignment", ID: id.String()}).Once()

	err := m.svc.Return(ctx, id, 100)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}
