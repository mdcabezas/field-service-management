package planning

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	coreservice "localis-backend/internal/service/core"
)

type DailyPlanService struct {
	planRepo       repository.DailyPlanRepository
	assignmentRepo repository.DailyPlanAssignmentRepository
	techRepo       repository.TechnicianRepository
	audit          *coreservice.AuditService
}

func NewDailyPlanService(
	planRepo repository.DailyPlanRepository,
	assignmentRepo repository.DailyPlanAssignmentRepository,
	techRepo repository.TechnicianRepository,
	audit *coreservice.AuditService,
) *DailyPlanService {
	return &DailyPlanService{
		planRepo:       planRepo,
		assignmentRepo: assignmentRepo,
		techRepo:       techRepo,
		audit:          audit,
	}
}

func (s *DailyPlanService) CreatePlan(ctx context.Context, plan *planning.DailyPlan) error {
	plan.ID = uuid.New()
	plan.CreatedAt = time.Now()

	if err := s.planRepo.Create(ctx, plan); err != nil {
		return fmt.Errorf("create daily plan: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyPlan, plan.ID, shared.AuditActionAdd, nil, plan, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_plan", "id", plan.ID, "error", err)
	}

	return nil
}

func (s *DailyPlanService) AssignTechnician(ctx context.Context, a *planning.DailyPlanAssignment) error {
	_, err := s.planRepo.GetByID(ctx, a.DailyPlanID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "daily_plan", a.DailyPlanID.String())
	}

	tech, err := s.techRepo.GetByID(ctx, a.TechID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "technician", a.TechID.String())
	}
	if !tech.IsActive {
		return &service.ValidationError{Field: "tech_id", Message: "technician is not active"}
	}

	a.ID = uuid.New()
	a.CreatedAt = time.Now()

	if err := s.assignmentRepo.Create(ctx, a); err != nil {
		return fmt.Errorf("create daily plan assignment: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyPlanAssignment, a.ID, shared.AuditActionAdd, nil, a, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_plan_assignment", "id", a.ID, "error", err)
	}

	return nil
}

// CloseDailyPlan is deferred — requires a DB migration to add a status column
// on the daily_plans table. Currently only records an audit event without
// modifying the plan state.
func (s *DailyPlanService) CloseDailyPlan(ctx context.Context, planID uuid.UUID) error {
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "daily_plan", planID.String())
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyPlan, planID, shared.AuditActionUpdate, plan, plan, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_plan", "id", planID, "error", err)
	}

	return nil
}

func (s *DailyPlanService) GetByID(ctx context.Context, id uuid.UUID) (*planning.DailyPlan, error) {
	result, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get daily plan by id: %w", err)
	}
	return result, nil
}

func (s *DailyPlanService) List(ctx context.Context, limit, offset int) (*repository.ListResult[planning.DailyPlan], error) {
	result, err := s.planRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list daily plans: %w", err)
	}
	return result, nil
}

func (s *DailyPlanService) Update(ctx context.Context, id uuid.UUID, plan *planning.DailyPlan) error {
	old, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "daily_plan", id.String())
	}

	if err := s.planRepo.Update(ctx, id, plan); err != nil {
		return fmt.Errorf("update daily plan: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyPlan, id, shared.AuditActionUpdate, old, plan, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_plan", "id", id, "error", err)
	}

	return nil
}

func (s *DailyPlanService) Delete(ctx context.Context, id uuid.UUID) error {
	old, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "daily_plan", id.String())
	}

	if err := s.planRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete daily plan: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyPlan, id, shared.AuditActionDelete, old, nil, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_plan", "id", id, "error", err)
	}

	return nil
}

func (s *DailyPlanService) ListAssignments(ctx context.Context, planID uuid.UUID) ([]planning.DailyPlanAssignment, error) {
	result, err := s.assignmentRepo.ListByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("list daily plan assignments: %w", err)
	}
	return result, nil
}

func (s *DailyPlanService) DeleteAssignment(ctx context.Context, id uuid.UUID) error {
	old, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "daily_plan_assignment", id.String())
	}

	if err := s.assignmentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete daily plan assignment: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyPlanAssignment, id, shared.AuditActionDelete, old, nil, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_plan_assignment", "id", id, "error", err)
	}

	return nil
}

func (s *DailyPlanService) GetAssignmentByID(ctx context.Context, id uuid.UUID) (*planning.DailyPlanAssignment, error) {
	result, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get daily plan assignment by id: %w", err)
	}
	return result, nil
}
