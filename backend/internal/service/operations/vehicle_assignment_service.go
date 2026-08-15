package operations

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	coreservice "localis-backend/internal/service/core"
)

type VehicleAssignmentService struct {
	vaRepo repository.VehicleAssignmentRepository
	vRepo  repository.VehicleRepository
	vPool  *pgxpool.Pool
	audit  *coreservice.AuditService
}

func NewVehicleAssignmentService(
	vaRepo repository.VehicleAssignmentRepository,
	vRepo repository.VehicleRepository,
	audit *coreservice.AuditService,
) *VehicleAssignmentService {
	s := &VehicleAssignmentService{
		vaRepo: vaRepo,
		vRepo:  vRepo,
		audit:  audit,
	}
	if p, ok := vaRepo.(interface{ Pool() *pgxpool.Pool }); ok {
		s.vPool = p.Pool()
	}
	return s
}

// NOTE: NoTx fallback paths are for development/testing only.
// In production, vPool is always non-nil (extracted from pgx-based repos),
// so the Tx path is always taken. The NoTx paths have TOCTOU race conditions
// because they use plain SELECT without row-level locking (SELECT ... FOR UPDATE).
// Do NOT rely on NoTx paths for concurrent production workloads.

func (s *VehicleAssignmentService) Assign(ctx context.Context, va *operations.VehicleAssignment) error {
	now := time.Now()
	va.DepartureTime = &now
	va.ID = uuid.New()
	va.CreatedAt = now

	if s.vPool != nil {
		tx, err := s.vPool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("start transaction: %w", err)
		}
		defer func() {
			if err := tx.Rollback(ctx); err != nil {
				slog.Warn("transaction rollback failed", "error", err)
			}
		}()

		if err := s.createInTx(ctx, tx, va); err != nil {
			return fmt.Errorf("create vehicle assignment in tx: %w", err)
		}

		vehicle, err := s.getVehicleForUpdate(ctx, tx, va.VehicleID)
		if err != nil {
			return fmt.Errorf("assign vehicle: %w", err)
		}
		if vehicle.Status != shared.VehicleStatusAvailable {
			return &service.ConflictError{Message: "vehicle is not available"}
		}
		vehicle.Status = shared.VehicleStatusInUse
		if err := s.updateVehicleInTx(ctx, tx, va.VehicleID, vehicle); err != nil {
			return fmt.Errorf("update vehicle status in tx: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit transaction: %w", err)
		}
	} else {
		vehicle, err := s.vRepo.GetByID(ctx, va.VehicleID)
		if err != nil {
			return service.HandleRepoGetByIDError(err, "vehicle", va.VehicleID.String())
		}
		if vehicle.Status != shared.VehicleStatusAvailable {
			return &service.ConflictError{Message: "vehicle is not available"}
		}

		if err := s.vaRepo.Create(ctx, va); err != nil {
			return fmt.Errorf("create vehicle assignment: %w", err)
		}

		vehicle.Status = shared.VehicleStatusInUse
		if err := s.vRepo.Update(ctx, va.VehicleID, vehicle); err != nil {
			return fmt.Errorf("update vehicle status: %w", err)
		}
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVehicleAssignment, va.ID, shared.AuditActionAdd, nil, va, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "vehicle_assignment", "id", va.ID, "error", err)
	}

	return nil
}

func (s *VehicleAssignmentService) Return(ctx context.Context, assignmentID uuid.UUID, returnMileage int) error {
	va, err := s.vaRepo.GetByID(ctx, assignmentID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "vehicle_assignment", assignmentID.String())
	}

	vaCopy := *va
	now := time.Now()
	va.ReturnTime = &now
	va.ReturnMileage = &returnMileage

	if s.vPool != nil {
		tx, err := s.vPool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("start transaction: %w", err)
		}
		defer func() {
			if err := tx.Rollback(ctx); err != nil {
				slog.Warn("transaction rollback failed", "error", err)
			}
		}()

		if err := s.updateAssignmentInTx(ctx, tx, assignmentID, va); err != nil {
			return fmt.Errorf("update vehicle assignment in tx: %w", err)
		}

		vehicle, err := s.getVehicleForUpdate(ctx, tx, va.VehicleID)
		if err != nil {
			return fmt.Errorf("return vehicle: %w", err)
		}
		vehicle.Status = shared.VehicleStatusAvailable
		if err := s.updateVehicleInTx(ctx, tx, va.VehicleID, vehicle); err != nil {
			return fmt.Errorf("update vehicle status in tx: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit transaction: %w", err)
		}
	} else {
		if err := s.vaRepo.Update(ctx, assignmentID, va); err != nil {
			return fmt.Errorf("update vehicle assignment: %w", err)
		}
		vehicle, err := s.vRepo.GetByID(ctx, va.VehicleID)
		if err != nil {
			return service.HandleRepoGetByIDError(err, "vehicle", va.VehicleID.String())
		}
		vehicle.Status = shared.VehicleStatusAvailable
		if err := s.vRepo.Update(ctx, va.VehicleID, vehicle); err != nil {
			return fmt.Errorf("update vehicle status: %w", err)
		}
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVehicleAssignment, assignmentID, shared.AuditActionUpdate, &vaCopy, va, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "vehicle_assignment", "id", assignmentID, "error", err)
	}

	return nil
}

func (s *VehicleAssignmentService) GetByID(ctx context.Context, id uuid.UUID) (*operations.VehicleAssignment, error) {
	result, err := s.vaRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get vehicle assignment by id: %w", err)
	}
	return result, nil
}

func (s *VehicleAssignmentService) List(ctx context.Context, limit, offset int) (*repository.ListResult[operations.VehicleAssignment], error) {
	result, err := s.vaRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list vehicle assignments: %w", err)
	}
	return result, nil
}

func (s *VehicleAssignmentService) Update(ctx context.Context, id uuid.UUID, va *operations.VehicleAssignment) error {
	existing, err := s.vaRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "vehicle_assignment", id.String())
	}

	if err := s.vaRepo.Update(ctx, id, va); err != nil {
		return fmt.Errorf("update vehicle assignment: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVehicleAssignment, id, shared.AuditActionUpdate, existing, va, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "vehicle_assignment", "id", id, "error", err)
	}
	return nil
}

func (s *VehicleAssignmentService) Delete(ctx context.Context, id uuid.UUID) error {
	old, err := s.vaRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "vehicle_assignment", id.String())
	}

	vehicleID := old.VehicleID

	if err := s.vaRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete vehicle assignment: %w", err)
	}

	if s.vPool != nil {
		tx, err := s.vPool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("start transaction: %w", err)
		}
		committed := false
		defer func() {
			if !committed {
				if err := tx.Rollback(ctx); err != nil {
					slog.Warn("transaction rollback failed", "error", err)
				}
			}
		}()

		vehicle, err := s.getVehicleForUpdate(ctx, tx, vehicleID)
		if err != nil {
			return fmt.Errorf("get vehicle for update: %w", err)
		}
		vehicle.Status = shared.VehicleStatusAvailable
		if err := s.updateVehicleInTx(ctx, tx, vehicleID, vehicle); err != nil {
			return fmt.Errorf("update vehicle status in tx: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit transaction: %w", err)
		}
		committed = true
	} else {
		vehicle, err := s.vRepo.GetByID(ctx, vehicleID)
		if err != nil {
			return service.HandleRepoGetByIDError(err, "vehicle", vehicleID.String())
		}
		vehicle.Status = shared.VehicleStatusAvailable
		if err := s.vRepo.Update(ctx, vehicleID, vehicle); err != nil {
			return fmt.Errorf("update vehicle status: %w", err)
		}
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVehicleAssignment, id, shared.AuditActionDelete, old, nil, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "vehicle_assignment", "id", id, "error", err)
	}

	return nil
}

func (s *VehicleAssignmentService) createInTx(ctx context.Context, tx pgx.Tx, va *operations.VehicleAssignment) error {
	type inTxAble interface {
		CreateInTx(ctx context.Context, tx pgx.Tx, va *operations.VehicleAssignment) error
	}
	if repo, ok := s.vaRepo.(inTxAble); ok {
		return repo.CreateInTx(ctx, tx, va)
	}
	return s.vaRepo.Create(ctx, va)
}

func (s *VehicleAssignmentService) updateAssignmentInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, va *operations.VehicleAssignment) error {
	type inTxAble interface {
		UpdateInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, va *operations.VehicleAssignment) error
	}
	if repo, ok := s.vaRepo.(inTxAble); ok {
		return repo.UpdateInTx(ctx, tx, id, va)
	}
	return s.vaRepo.Update(ctx, id, va)
}

func (s *VehicleAssignmentService) getVehicleForUpdate(ctx context.Context, tx pgx.Tx, vehicleID uuid.UUID) (*inventory.Vehicle, error) {
	var vehicle inventory.Vehicle
	err := tx.QueryRow(ctx,
		`SELECT id, type, license_plate, name, brand, model, year, status, capacity, created_at
		 FROM inventory.vehicles
		 WHERE id = $1 FOR UPDATE`, vehicleID,
	).Scan(&vehicle.ID, &vehicle.Type, &vehicle.LicensePlate, &vehicle.Name, &vehicle.Brand, &vehicle.Model, &vehicle.Year, &vehicle.Status, &vehicle.Capacity, &vehicle.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &service.NotFoundError{Resource: "vehicle", ID: vehicleID.String()}
		}
		return nil, fmt.Errorf("get vehicle for update: %w", err)
	}
	return &vehicle, nil
}

func (s *VehicleAssignmentService) updateVehicleInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, vehicle *inventory.Vehicle) error {
	type inTxAble interface {
		UpdateInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, vehicle *inventory.Vehicle) error
	}
	if repo, ok := s.vRepo.(inTxAble); ok {
		return repo.UpdateInTx(ctx, tx, id, vehicle)
	}
	return s.vRepo.Update(ctx, id, vehicle)
}
