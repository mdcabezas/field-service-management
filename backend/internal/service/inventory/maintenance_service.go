package inventory

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	inventorymodel "localis-backend/internal/model/inventory"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	coreservice "localis-backend/internal/service/core"
	notifservice "localis-backend/internal/service/notifications"
)

type MaintenanceService struct {
	scheduleRepo repository.MaintenanceScheduleRepository
	recordRepo   repository.MaintenanceRecordRepository
	vehicleRepo  repository.VehicleRepository
	audit        *coreservice.AuditService
	notifService *notifservice.NotificationService
}

func NewMaintenanceService(
	scheduleRepo repository.MaintenanceScheduleRepository,
	recordRepo repository.MaintenanceRecordRepository,
	vehicleRepo repository.VehicleRepository,
	audit *coreservice.AuditService,
	notifService *notifservice.NotificationService,
) *MaintenanceService {
	return &MaintenanceService{
		scheduleRepo: scheduleRepo,
		recordRepo:   recordRepo,
		vehicleRepo:  vehicleRepo,
		audit:        audit,
		notifService: notifService,
	}
}

func (s *MaintenanceService) CheckDueMaintenance(ctx context.Context) ([]inventorymodel.MaintenanceSchedule, error) {
	schedules, err := s.scheduleRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active maintenance schedules: %w", err)
	}

	now := time.Now()
	due := make([]inventorymodel.MaintenanceSchedule, 0)

	for _, sched := range schedules {
		if sched.LastServiceDate == nil || sched.FrequencyDays == nil {
			due = append(due, sched)
			continue
		}

		nextDue := sched.LastServiceDate.AddDate(0, 0, *sched.FrequencyDays)
		if now.After(nextDue) {
			due = append(due, sched)
		}
	}

	return due, nil
}

func (s *MaintenanceService) RecordMaintenance(ctx context.Context, record *inventorymodel.MaintenanceRecord) error {
	schedules, err := s.scheduleRepo.ListActive(ctx)
	if err != nil {
		return fmt.Errorf("list active maintenance schedules: %w", err)
	}

	for i := range schedules {
		sched := &schedules[i]
		if sched.Type == record.Type && sched.ReferenceID == record.ReferenceID {
			sched.LastServiceDate = &record.Date
			if err := s.scheduleRepo.Update(ctx, sched.ID, sched); err != nil {
				return fmt.Errorf("update maintenance schedule: %w", err)
			}
			break
		}
	}

	record.ID = uuid.New()
	record.CreatedAt = time.Now()

	if err := s.recordRepo.Create(ctx, record); err != nil {
		return fmt.Errorf("create maintenance record: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeMaintenanceRecord, record.ID, shared.AuditActionAdd, nil, record, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "maintenance_record", "id", record.ID, "error", err)
	}

	return nil
}

func (s *MaintenanceService) GetScheduleByID(ctx context.Context, id uuid.UUID) (*inventorymodel.MaintenanceSchedule, error) {
	result, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get maintenance schedule by id: %w", err)
	}
	return result, nil
}

func (s *MaintenanceService) ListSchedules(ctx context.Context, limit, offset int) (*repository.ListResult[inventorymodel.MaintenanceSchedule], error) {
	result, err := s.scheduleRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list maintenance schedules: %w", err)
	}
	return result, nil
}

func (s *MaintenanceService) CreateSchedule(ctx context.Context, sched *inventorymodel.MaintenanceSchedule) error {
	sched.ID = uuid.New()
	sched.CreatedAt = time.Now()

	if err := s.scheduleRepo.Create(ctx, sched); err != nil {
		return fmt.Errorf("create maintenance schedule: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeMaintenanceSchedule, sched.ID, shared.AuditActionAdd, nil, sched, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "maintenance_schedule", "id", sched.ID, "error", err)
	}

	return nil
}

func (s *MaintenanceService) UpdateSchedule(ctx context.Context, id uuid.UUID, sched *inventorymodel.MaintenanceSchedule) error {
	existing, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "maintenance_schedule", id.String())
	}

	if err := s.scheduleRepo.Update(ctx, id, sched); err != nil {
		return fmt.Errorf("update maintenance schedule: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeMaintenanceSchedule, id, shared.AuditActionUpdate, existing, sched, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "maintenance_schedule", "id", id, "error", err)
	}

	return nil
}

func (s *MaintenanceService) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	existing, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "maintenance_schedule", id.String())
	}

	if err := s.scheduleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete maintenance schedule: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeMaintenanceSchedule, id, shared.AuditActionDelete, existing, nil, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "maintenance_schedule", "id", id, "error", err)
	}

	return nil
}

func (s *MaintenanceService) GetByID(ctx context.Context, id uuid.UUID) (*inventorymodel.MaintenanceRecord, error) {
	result, err := s.recordRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get maintenance record by id: %w", err)
	}
	return result, nil
}

func (s *MaintenanceService) ListRecords(ctx context.Context, limit, offset int) (*repository.ListResult[inventorymodel.MaintenanceRecord], error) {
	result, err := s.recordRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list maintenance records: %w", err)
	}
	return result, nil
}

func (s *MaintenanceService) DeleteRecord(ctx context.Context, id uuid.UUID) error {
	existing, err := s.recordRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "maintenance_record", id.String())
	}

	if err := s.recordRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete maintenance record: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeMaintenanceRecord, id, shared.AuditActionDelete, existing, nil, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "maintenance_record", "id", id, "error", err)
	}

	return nil
}

func (s *MaintenanceService) UpdateRecord(ctx context.Context, id uuid.UUID, rec *inventorymodel.MaintenanceRecord) error {
	existing, err := s.recordRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "maintenance_record", id.String())
	}

	if err := s.recordRepo.Update(ctx, id, rec); err != nil {
		return fmt.Errorf("update maintenance record: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeMaintenanceRecord, id, shared.AuditActionUpdate, existing, rec, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "maintenance_record", "id", id, "error", err)
	}

	return nil
}
