package operations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	"localis-backend/internal/service/core"
)

type VisitAssignmentService struct {
	assignmentRepo repository.VisitAssignmentRepository
	visitRepo      repository.VisitRepository
	techRepo       repository.TechnicianRepository
	audit          *core.AuditService
}

func NewVisitAssignmentService(
	assignmentRepo repository.VisitAssignmentRepository,
	visitRepo repository.VisitRepository,
	techRepo repository.TechnicianRepository,
	audit *core.AuditService,
) *VisitAssignmentService {
	return &VisitAssignmentService{
		assignmentRepo: assignmentRepo,
		visitRepo:      visitRepo,
		techRepo:       techRepo,
		audit:          audit,
	}
}

func (s *VisitAssignmentService) AssignTechnician(ctx context.Context, assignment *operations.VisitAssignment) error {
	_, err := s.visitRepo.GetByID(ctx, assignment.VisitID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit", assignment.VisitID.String())
	}

	tech, err := s.techRepo.GetByID(ctx, assignment.TechID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "technician", assignment.TechID.String())
	}

	if !tech.IsActive {
		return &service.ValidationError{Field: "tech_id", Message: "technician is not active"}
	}

	if err := s.assignmentRepo.Create(ctx, assignment); err != nil {
		return fmt.Errorf("create visit assignment: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitAssignment, assignment.ID, shared.AuditActionAdd, nil, assignment, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_assignment", "id", assignment.ID, "error", err)
	}

	return nil
}

func (s *VisitAssignmentService) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitAssignment, error) {
	result, err := s.assignmentRepo.ListByVisit(ctx, visitID)
	if err != nil {
		return nil, fmt.Errorf("list visit assignments: %w", err)
	}
	return result, nil
}

func (s *VisitAssignmentService) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitAssignment, error) {
	result, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get visit assignment by id: %w", err)
	}
	return result, nil
}

func (s *VisitAssignmentService) RemoveAssignment(ctx context.Context, id uuid.UUID) error {
	assignment, err := s.assignmentRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit_assignment", id.String())
	}

	if err := s.assignmentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete visit assignment: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitAssignment, id, shared.AuditActionDelete, assignment, nil, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_assignment", "id", id, "error", err)
	}

	return nil
}
