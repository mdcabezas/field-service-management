package operations

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/service"
	"localis-backend/internal/service/core"
	"localis-backend/internal/testutil/mocks"
)

func newAssignmentService(t *testing.T) (*VisitAssignmentService, *mocks.VisitAssignmentRepository, *mocks.VisitRepository, *mocks.TechnicianRepository, *mocks.CorePlanAuditLogRepository) {
	t.Helper()
	assignmentRepo := mocks.NewVisitAssignmentRepository(t)
	visitRepo := mocks.NewVisitRepository(t)
	techRepo := mocks.NewTechnicianRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	svc := NewVisitAssignmentService(assignmentRepo, visitRepo, techRepo, core.NewAuditService(auditRepo))
	return svc, assignmentRepo, visitRepo, techRepo, auditRepo
}

func TestAssignTechnician_Success(t *testing.T) {
	ctx := context.Background()
	svc, assignmentRepo, visitRepo, techRepo, auditRepo := newAssignmentService(t)

	assignment := &operations.VisitAssignment{ID: uuid.New(), VisitID: uuid.New(), TechID: uuid.New()}
	tech := newTechnicianFixture(true)
	visit := newVisitFixture(assignment.VisitID, shared.VisitStatusScheduled)

	visitRepo.On("GetByID", ctx, assignment.VisitID).Return(visit, nil).Once()
	techRepo.On("GetByID", ctx, assignment.TechID).Return(tech, nil).Once()
	assignmentRepo.On("Create", ctx, assignment).Return(nil).Once()
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := svc.AssignTechnician(ctx, assignment)
	require.NoError(t, err)
}

func TestAssignTechnician_InactiveTechnician(t *testing.T) {
	ctx := context.Background()
	svc, _, visitRepo, techRepo, _ := newAssignmentService(t)

	assignment := &operations.VisitAssignment{ID: uuid.New(), VisitID: uuid.New(), TechID: uuid.New()}
	visit := newVisitFixture(assignment.VisitID, shared.VisitStatusScheduled)
	tech := newTechnicianFixture(false)

	visitRepo.On("GetByID", ctx, assignment.VisitID).Return(visit, nil).Once()
	techRepo.On("GetByID", ctx, assignment.TechID).Return(tech, nil).Once()

	err := svc.AssignTechnician(ctx, assignment)
	var ve *service.ValidationError
	require.ErrorAs(t, err, &ve)
	require.Equal(t, "tech_id", ve.Field)
}

func TestAssignTechnician_VisitNotFound(t *testing.T) {
	ctx := context.Background()
	svc, _, visitRepo, _, _ := newAssignmentService(t)

	assignment := &operations.VisitAssignment{ID: uuid.New(), VisitID: uuid.New(), TechID: uuid.New()}
	visitRepo.On("GetByID", ctx, assignment.VisitID).Return(nil, &service.NotFoundError{Resource: "visit", ID: assignment.VisitID.String()}).Once()

	err := svc.AssignTechnician(ctx, assignment)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestRemoveAssignment_Success(t *testing.T) {
	ctx := context.Background()
	svc, assignmentRepo, _, _, auditRepo := newAssignmentService(t)

	id := uuid.New()
	assignment := &operations.VisitAssignment{ID: id}
	assignmentRepo.On("GetByID", ctx, id).Return(assignment, nil).Once()
	assignmentRepo.On("Delete", ctx, id).Return(nil).Once()
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := svc.RemoveAssignment(ctx, id)
	require.NoError(t, err)
}

func TestRemoveAssignment_NotFound(t *testing.T) {
	ctx := context.Background()
	svc, assignmentRepo, _, _, _ := newAssignmentService(t)

	id := uuid.New()
	assignmentRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "visit_assignment", ID: id.String()}).Once()

	err := svc.RemoveAssignment(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestVisitAssignment_ListByVisit(t *testing.T) {
	ctx := context.Background()
	svc, assignmentRepo, _, _, _ := newAssignmentService(t)

	visitID := uuid.New()
	expected := []operations.VisitAssignment{{ID: uuid.New(), VisitID: visitID}}
	assignmentRepo.On("ListByVisit", ctx, visitID).Return(expected, nil).Once()

	got, err := svc.ListByVisit(ctx, visitID)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestVisitAssignmentService_GetByID(t *testing.T) {
	ctx := context.Background()
	svc, assignmentRepo, _, _, _ := newAssignmentService(t)

	id := uuid.New()
	expected := &operations.VisitAssignment{ID: id, VisitID: uuid.New(), TechID: uuid.New()}
	assignmentRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

	got, err := svc.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}
