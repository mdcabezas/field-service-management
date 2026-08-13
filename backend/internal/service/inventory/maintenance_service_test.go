package inventory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	inventorymodel "localis-backend/internal/model/inventory"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	coreservice "localis-backend/internal/service/core"
	notifservice "localis-backend/internal/service/notifications"
	"localis-backend/internal/testutil/mocks"
)

type maintenanceSvcMocks struct {
	scheduleRepo *mocks.MaintenanceScheduleRepository
	recordRepo   *mocks.MaintenanceRecordRepository
	vehicleRepo  *mocks.VehicleRepository
	auditRepo    *mocks.CorePlanAuditLogRepository
	notifRepo    *mocks.NotificationRepository
	svc          *MaintenanceService
}

func newMaintenanceService(t *testing.T) *maintenanceSvcMocks {
	t.Helper()
	scheduleRepo := mocks.NewMaintenanceScheduleRepository(t)
	recordRepo := mocks.NewMaintenanceRecordRepository(t)
	vehicleRepo := mocks.NewVehicleRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	notifRepo := mocks.NewNotificationRepository(t)
	audit := coreservice.NewAuditService(auditRepo)
	notifSvc := notifservice.NewNotificationService(notifRepo)
	svc := NewMaintenanceService(scheduleRepo, recordRepo, vehicleRepo, audit, notifSvc)
	return &maintenanceSvcMocks{scheduleRepo: scheduleRepo, recordRepo: recordRepo, vehicleRepo: vehicleRepo, auditRepo: auditRepo, notifRepo: notifRepo, svc: svc}
}

func TestCheckDueMaintenance_NoDates(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	schedules := []inventorymodel.MaintenanceSchedule{
		{ID: uuid.New(), Type: shared.MaintenanceTypeVehicle, ReferenceID: uuid.New()},
		{ID: uuid.New(), Type: shared.MaintenanceTypeTool, ReferenceID: uuid.New(), FrequencyDays: intPtr(30)},
	}
	m.scheduleRepo.On("ListActive", ctx).Return(schedules, nil).Once()

	due, err := m.svc.CheckDueMaintenance(ctx)
	require.NoError(t, err)
	require.Len(t, due, 2)
}

func TestCheckDueMaintenance_DueAndNotDue(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	past := time.Now().AddDate(0, 0, -40)
	future := time.Now().AddDate(0, 0, -10)
	days := 30
	schedules := []inventorymodel.MaintenanceSchedule{
		{ID: uuid.New(), ReferenceID: uuid.New(), FrequencyDays: &days, LastServiceDate: &past},
		{ID: uuid.New(), ReferenceID: uuid.New(), FrequencyDays: &days, LastServiceDate: &future},
	}
	m.scheduleRepo.On("ListActive", ctx).Return(schedules, nil).Once()

	due, err := m.svc.CheckDueMaintenance(ctx)
	require.NoError(t, err)
	require.Len(t, due, 1)
	require.Equal(t, schedules[0].ID, due[0].ID)
}

func TestCheckDueMaintenance_RepoError(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	m.scheduleRepo.On("ListActive", ctx).Return(nil, errors.New("db down")).Once()

	_, err := m.svc.CheckDueMaintenance(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list active maintenance schedules")
}

func TestRecordMaintenance_Success_UpdatesSchedule(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	refID := uuid.New()
	record := &inventorymodel.MaintenanceRecord{Type: shared.MaintenanceTypeVehicle, ReferenceID: refID, Date: time.Now()}
	sched := inventorymodel.MaintenanceSchedule{ID: uuid.New(), Type: shared.MaintenanceTypeVehicle, ReferenceID: refID, Active: true}

	m.scheduleRepo.On("ListActive", ctx).Return([]inventorymodel.MaintenanceSchedule{sched}, nil).Once()
	var updatedSched *inventorymodel.MaintenanceSchedule
	m.scheduleRepo.On("Update", ctx, sched.ID, mock.MatchedBy(func(s *inventorymodel.MaintenanceSchedule) bool {
		updatedSched = s
		return s.ID == sched.ID && s.LastServiceDate != nil
	})).Return(nil).Once()
	m.recordRepo.On("Create", ctx, record).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.RecordMaintenance(ctx, record)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, record.ID)
	require.NotZero(t, record.CreatedAt)
	require.NotNil(t, updatedSched)
	require.NotNil(t, updatedSched.LastServiceDate)
	require.True(t, updatedSched.LastServiceDate.Equal(record.Date))
}

func TestRecordMaintenance_NoMatchingSchedule(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	record := &inventorymodel.MaintenanceRecord{Type: shared.MaintenanceTypeVehicle, ReferenceID: uuid.New(), Date: time.Now()}
	m.scheduleRepo.On("ListActive", ctx).Return([]inventorymodel.MaintenanceSchedule{}, nil).Once()
	m.recordRepo.On("Create", ctx, record).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.RecordMaintenance(ctx, record)
	require.NoError(t, err)
}

func TestRecordMaintenance_ListActiveError(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	m.scheduleRepo.On("ListActive", ctx).Return(nil, errors.New("db down")).Once()

	err := m.svc.RecordMaintenance(ctx, &inventorymodel.MaintenanceRecord{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "list active maintenance schedules")
}

func TestRecordMaintenance_ScheduleUpdateError(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	refID := uuid.New()
	record := &inventorymodel.MaintenanceRecord{Type: shared.MaintenanceTypeVehicle, ReferenceID: refID, Date: time.Now()}
	sched := inventorymodel.MaintenanceSchedule{ID: uuid.New(), Type: shared.MaintenanceTypeVehicle, ReferenceID: refID}

	m.scheduleRepo.On("ListActive", ctx).Return([]inventorymodel.MaintenanceSchedule{sched}, nil).Once()
	m.scheduleRepo.On("Update", ctx, sched.ID, mock.MatchedBy(func(s *inventorymodel.MaintenanceSchedule) bool {
		return s.ID == sched.ID && s.LastServiceDate != nil
	})).Return(errors.New("db down")).Once()

	err := m.svc.RecordMaintenance(ctx, record)
	require.Error(t, err)
	require.Contains(t, err.Error(), "update maintenance schedule")
}

func TestRecordMaintenance_CreateError(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	record := &inventorymodel.MaintenanceRecord{Type: shared.MaintenanceTypeVehicle, ReferenceID: uuid.New(), Date: time.Now()}
	m.scheduleRepo.On("ListActive", ctx).Return([]inventorymodel.MaintenanceSchedule{}, nil).Once()
	m.recordRepo.On("Create", ctx, record).Return(errors.New("db down")).Once()

	err := m.svc.RecordMaintenance(ctx, record)
	require.Error(t, err)
	require.Contains(t, err.Error(), "create maintenance record")
}

func TestCreateSchedule_Success(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	sched := &inventorymodel.MaintenanceSchedule{Type: shared.MaintenanceTypeVehicle, ReferenceID: uuid.New(), Active: true}
	m.scheduleRepo.On("Create", ctx, sched).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.CreateSchedule(ctx, sched)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, sched.ID)
}

func TestUpdateSchedule_Success(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	id := uuid.New()
	existing := &inventorymodel.MaintenanceSchedule{ID: id, Active: true}
	updated := &inventorymodel.MaintenanceSchedule{ID: id, Active: false}

	m.scheduleRepo.On("GetByID", ctx, id).Return(existing, nil).Once()
	m.scheduleRepo.On("Update", ctx, id, updated).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.UpdateSchedule(ctx, id, updated)
	require.NoError(t, err)
}

func TestUpdateSchedule_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	id := uuid.New()
	m.scheduleRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "maintenance_schedule", ID: id.String()}).Once()

	err := m.svc.UpdateSchedule(ctx, id, &inventorymodel.MaintenanceSchedule{})
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestDeleteSchedule_Success(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	id := uuid.New()
	m.scheduleRepo.On("GetByID", ctx, id).Return(&inventorymodel.MaintenanceSchedule{ID: id}, nil).Once()
	m.scheduleRepo.On("Delete", ctx, id).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.DeleteSchedule(ctx, id)
	require.NoError(t, err)
}

func TestGetScheduleByID_Success(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	id := uuid.New()
	expected := &inventorymodel.MaintenanceSchedule{ID: id}
	m.scheduleRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

	got, err := m.svc.GetScheduleByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestListSchedules_Success(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	expected := &repository.ListResult[inventorymodel.MaintenanceSchedule]{Items: []inventorymodel.MaintenanceSchedule{{ID: uuid.New()}}, Total: 1}
	m.scheduleRepo.On("List", ctx, 10, 0).Return(expected, nil).Once()

	got, err := m.svc.ListSchedules(ctx, 10, 0)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestUpdateRecord_Success(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	id := uuid.New()
	m.recordRepo.On("GetByID", ctx, id).Return(&inventorymodel.MaintenanceRecord{ID: id}, nil).Once()
	m.recordRepo.On("Update", ctx, id, mock.Anything).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.UpdateRecord(ctx, id, &inventorymodel.MaintenanceRecord{ID: id})
	require.NoError(t, err)
}

func TestDeleteRecord_Success(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	id := uuid.New()
	m.recordRepo.On("GetByID", ctx, id).Return(&inventorymodel.MaintenanceRecord{ID: id}, nil).Once()
	m.recordRepo.On("Delete", ctx, id).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.DeleteRecord(ctx, id)
	require.NoError(t, err)
}

func TestGetByID_Success(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	id := uuid.New()
	expected := &inventorymodel.MaintenanceRecord{ID: id}
	m.recordRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

	got, err := m.svc.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestListRecords_Success(t *testing.T) {
	ctx := context.Background()
	m := newMaintenanceService(t)

	expected := &repository.ListResult[inventorymodel.MaintenanceRecord]{Items: []inventorymodel.MaintenanceRecord{{ID: uuid.New()}}, Total: 1}
	m.recordRepo.On("List", ctx, 10, 0).Return(expected, nil).Once()

	got, err := m.svc.ListRecords(ctx, 10, 0)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func intPtr(v int) *int {
	return &v
}
