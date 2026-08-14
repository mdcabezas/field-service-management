package operations

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	"localis-backend/internal/service/core"
	"localis-backend/internal/testutil/mocks"
)

type visitServiceMocks struct {
	visitRepo  *mocks.VisitRepository
	checkpoint *mocks.VisitCheckpointRepository
	ctmRepo    *mocks.ChecklistTemplateMaterialRepository
	cttRepo    *mocks.ChecklistTemplateToolRepository
	ctaRepo    *mocks.ChecklistTemplateEPPRepository
	vcmRepo    *mocks.VisitChecklistMaterialRepository
	vctRepo    *mocks.VisitChecklistToolRepository
	vceRepo    *mocks.VisitChecklistEPPRepository
	auditRepo  *mocks.CorePlanAuditLogRepository
	service    *VisitService
}

func newVisitFixture(id uuid.UUID, status shared.VisitStatus) *operations.Visit {
	now := time.Now()
	return &operations.Visit{
		ID:          id,
		Type:        uuid.New(),
		Status:      status,
		Priority:    shared.VisitPriorityNormal,
		BillingTo:   shared.BillingToPartner,
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func newCheckpointFixture(visitID uuid.UUID) *operations.VisitCheckpoint {
	return &operations.VisitCheckpoint{
		ID:        uuid.New(),
		VisitID:   visitID,
		Type:      shared.CheckpointTypeArrival,
		Timestamp: time.Now(),
		CreatedAt: time.Now(),
	}
}

// newVisitService builds a VisitService with all repos mocked. The mocks do not
// implement Pool(), so vPool stays nil and the NoTx code paths are exercised.
func newVisitService(t *testing.T) *visitServiceMocks {
	t.Helper()
	visitRepo := mocks.NewVisitRepository(t)
	assignmentRepo := mocks.NewVisitAssignmentRepository(t)
	techRepo := mocks.NewTechnicianRepository(t)
	checkpointRepo := mocks.NewVisitCheckpointRepository(t)
	ctmRepo := mocks.NewChecklistTemplateMaterialRepository(t)
	cttRepo := mocks.NewChecklistTemplateToolRepository(t)
	ctaRepo := mocks.NewChecklistTemplateEPPRepository(t)
	vcmRepo := mocks.NewVisitChecklistMaterialRepository(t)
	vctRepo := mocks.NewVisitChecklistToolRepository(t)
	vceRepo := mocks.NewVisitChecklistEPPRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)

	svc := NewVisitService(VisitServiceDeps{
		VisitRepo:      visitRepo,
		AssignmentRepo: assignmentRepo,
		TechnicianRepo: techRepo,
		CheckpointRepo: checkpointRepo,
		CTMRepo:        ctmRepo,
		CTTRepo:        cttRepo,
		CTARepo:        ctaRepo,
		VCMRepo:        vcmRepo,
		VCTRepo:        vctRepo,
		VCERepo:        vceRepo,
		AuditService:   core.NewAuditService(auditRepo),
	})
	return &visitServiceMocks{
		visitRepo:  visitRepo,
		checkpoint: checkpointRepo,
		ctmRepo:    ctmRepo,
		cttRepo:    cttRepo,
		ctaRepo:    ctaRepo,
		vcmRepo:    vcmRepo,
		vctRepo:    vctRepo,
		vceRepo:    vceRepo,
		auditRepo:  auditRepo,
		service:    svc,
	}
}

func expectAudit(m *visitServiceMocks, times int) {
	if times <= 0 {
		m.auditRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
		return
	}
	m.auditRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(times)
}

func TestCanTransition(t *testing.T) {
	tests := []struct {
		name string
		from shared.VisitStatus
		to   shared.VisitStatus
		want bool
	}{
		{"scheduled to en_route", shared.VisitStatusScheduled, shared.VisitStatusEnRoute, true},
		{"scheduled to cancelled", shared.VisitStatusScheduled, shared.VisitStatusCancelled, true},
		{"scheduled to rescheduled", shared.VisitStatusScheduled, shared.VisitStatusRescheduled, true},
		{"scheduled to completed", shared.VisitStatusScheduled, shared.VisitStatusCompleted, false},
		{"en_route to in_progress", shared.VisitStatusEnRoute, shared.VisitStatusInProgress, true},
		{"en_route to cancelled", shared.VisitStatusEnRoute, shared.VisitStatusCancelled, true},
		{"en_route to scheduled", shared.VisitStatusEnRoute, shared.VisitStatusScheduled, false},
		{"in_progress to completed", shared.VisitStatusInProgress, shared.VisitStatusCompleted, true},
		{"in_progress to cancelled", shared.VisitStatusInProgress, shared.VisitStatusCancelled, true},
		{"in_progress to en_route", shared.VisitStatusInProgress, shared.VisitStatusEnRoute, false},
		{"rescheduled to scheduled", shared.VisitStatusRescheduled, shared.VisitStatusScheduled, true},
		{"completed to scheduled", shared.VisitStatusCompleted, shared.VisitStatusScheduled, false},
		{"cancelled to scheduled", shared.VisitStatusCancelled, shared.VisitStatusScheduled, false},
		{"unknown to scheduled", shared.VisitStatus("bogus"), shared.VisitStatusScheduled, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canTransition(tt.from, tt.to)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCreateVisit(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	input := &operations.Visit{
		Type:        uuid.New(),
		BillingTo:   shared.BillingToPartner,
		ScheduledAt: time.Now(),
	}
	m.visitRepo.On("Create", ctx, input).Return(nil).Once()
	expectAudit(m, 1)

	err := m.service.CreateVisit(ctx, input)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, input.ID)
	require.Equal(t, shared.VisitStatusScheduled, input.Status)
	require.Equal(t, shared.VisitPriorityNormal, input.Priority)
	require.False(t, input.CreatedAt.IsZero())
	require.False(t, input.UpdatedAt.IsZero())
}

func TestCreateVisit_RepoError(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	input := &operations.Visit{Type: uuid.New(), BillingTo: shared.BillingToPartner, ScheduledAt: time.Now()}
	repoErr := errors.New("db down")
	m.visitRepo.On("Create", ctx, input).Return(repoErr).Once()

	err := m.service.CreateVisit(ctx, input)
	require.Error(t, err)
	require.ErrorIs(t, err, repoErr)
}

func TestStartVisit_Success(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusScheduled)
	cp := newCheckpointFixture(id)

	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()
	m.visitRepo.On("Update", ctx, id, visit).Return(nil).Once()
	m.checkpoint.On("Create", ctx, cp).Return(nil).Once()
	expectAudit(m, 1)

	err := m.service.StartVisit(ctx, id, cp)
	require.NoError(t, err)
	require.Equal(t, shared.VisitStatusEnRoute, visit.Status)
	require.NotNil(t, visit.StartedAt)
	require.Equal(t, id, cp.VisitID)
}

func TestStartVisit_InvalidTransition(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusCompleted)
	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()

	err := m.service.StartVisit(ctx, id, nil)
	var te *service.TransitionError
	require.ErrorAs(t, err, &te)
	require.Equal(t, "visit", te.Entity)
	require.Equal(t, string(shared.VisitStatusCompleted), te.From)
	require.Equal(t, string(shared.VisitStatusEnRoute), te.To)
}

func TestStartVisit_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	m.visitRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "visit", ID: id.String()}).Once()

	err := m.service.StartVisit(ctx, id, nil)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestBeginWork_Success(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusEnRoute)
	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()
	m.visitRepo.On("Update", ctx, id, visit).Return(nil).Once()
	expectAudit(m, 1)

	err := m.service.BeginWork(ctx, id)
	require.NoError(t, err)
	require.Equal(t, shared.VisitStatusInProgress, visit.Status)
}

func TestBeginWork_InvalidTransition(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusScheduled)
	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()

	err := m.service.BeginWork(ctx, id)
	var te *service.TransitionError
	require.ErrorAs(t, err, &te)
}

func TestCompleteVisit_Success(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusInProgress)
	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()
	m.visitRepo.On("Update", ctx, id, visit).Return(nil).Once()
	expectAudit(m, 1)

	err := m.service.CompleteVisit(ctx, id)
	require.NoError(t, err)
	require.Equal(t, shared.VisitStatusCompleted, visit.Status)
	require.NotNil(t, visit.CompletedAt)
}

func TestCompleteVisit_InvalidTransition(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusCancelled)
	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()

	err := m.service.CompleteVisit(ctx, id)
	var te *service.TransitionError
	require.ErrorAs(t, err, &te)
}

func TestCancelVisit_Success(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusScheduled)
	reason := "customer cancelled"

	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()
	m.visitRepo.On("Update", ctx, id, visit).Return(nil).Once()
	expectAudit(m, 1)

	err := m.service.CancelVisit(ctx, id, &reason)
	require.NoError(t, err)
	require.Equal(t, shared.VisitStatusCancelled, visit.Status)
}

func TestCancelVisit_InvalidTransition(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusCompleted)
	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()

	err := m.service.CancelVisit(ctx, id, nil)
	var te *service.TransitionError
	require.ErrorAs(t, err, &te)
}

func TestRescheduleVisit_Success(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusScheduled)
	newScheduledAt := time.Now().Add(24 * time.Hour)

	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()
	m.visitRepo.On("Update", ctx, id, visit).Return(nil).Once()
	m.visitRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
	expectAudit(m, 2)

	err := m.service.RescheduleVisit(ctx, id, newScheduledAt)
	require.NoError(t, err)
	require.Equal(t, shared.VisitStatusRescheduled, visit.Status)
}

func TestRescheduleVisit_InvalidTransition(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusCompleted)
	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()

	err := m.service.RescheduleVisit(ctx, id, time.Now())
	var te *service.TransitionError
	require.ErrorAs(t, err, &te)
}

func TestPopulateChecklist(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	visitID := uuid.New()
	templateID := uuid.New()
	m.visitRepo.On("GetByID", ctx, visitID).Return(newVisitFixture(visitID, shared.VisitStatusScheduled), nil).Once()

	mat := operations.ChecklistTemplateMaterial{ID: uuid.New(), TemplateID: templateID, MaterialID: uuid.New(), DefaultQuantity: 2}
	tool := operations.ChecklistTemplateTool{ID: uuid.New(), TemplateID: templateID, ToolID: uuid.New(), DefaultQuantity: 1}
	epp := operations.ChecklistTemplateEPP{ID: uuid.New(), TemplateID: templateID, EPPID: uuid.New(), DefaultQuantity: 3}

	m.ctmRepo.On("ListByTemplate", ctx, templateID).Return([]operations.ChecklistTemplateMaterial{mat}, nil).Once()
	m.vcmRepo.On("CreateBatch", ctx, mock.Anything).Return(nil).Once()
	m.cttRepo.On("ListByTemplate", ctx, templateID).Return([]operations.ChecklistTemplateTool{tool}, nil).Once()
	m.vctRepo.On("CreateBatch", ctx, mock.Anything).Return(nil).Once()
	m.ctaRepo.On("ListByTemplate", ctx, templateID).Return([]operations.ChecklistTemplateEPP{epp}, nil).Once()
	m.vceRepo.On("CreateBatch", ctx, mock.Anything).Return(nil).Once()

	err := m.service.PopulateChecklist(ctx, visitID, templateID)
	require.NoError(t, err)
}

func TestPopulateChecklist_EmptyTemplate(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	visitID := uuid.New()
	templateID := uuid.New()
	m.visitRepo.On("GetByID", ctx, visitID).Return(newVisitFixture(visitID, shared.VisitStatusScheduled), nil).Once()

	m.ctmRepo.On("ListByTemplate", ctx, templateID).Return([]operations.ChecklistTemplateMaterial{}, nil).Once()
	m.cttRepo.On("ListByTemplate", ctx, templateID).Return([]operations.ChecklistTemplateTool{}, nil).Once()
	m.ctaRepo.On("ListByTemplate", ctx, templateID).Return([]operations.ChecklistTemplateEPP{}, nil).Once()

	err := m.service.PopulateChecklist(ctx, visitID, templateID)
	require.NoError(t, err)
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	expected := newVisitFixture(id, shared.VisitStatusScheduled)
	m.visitRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

	got, err := m.service.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestGetByID_Error(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	repoErr := errors.New("db down")
	m.visitRepo.On("GetByID", ctx, id).Return(nil, repoErr).Once()

	_, err := m.service.GetByID(ctx, id)
	require.Error(t, err)
	require.ErrorIs(t, err, repoErr)
}

func TestDelete_Success(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	visit := newVisitFixture(id, shared.VisitStatusScheduled)
	m.visitRepo.On("GetByID", ctx, id).Return(visit, nil).Once()
	m.visitRepo.On("Delete", ctx, id).Return(nil).Once()
	expectAudit(m, 1)

	err := m.service.Delete(ctx, id)
	require.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	id := uuid.New()
	m.visitRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "visit", ID: id.String()}).Once()

	err := m.service.Delete(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestVisitService_List(t *testing.T) {
	ctx := context.Background()
	m := newVisitService(t)

	expected := &repository.ListResult[operations.Visit]{
		Items: []operations.Visit{*newVisitFixture(uuid.New(), shared.VisitStatusScheduled)},
		Total: 1,
	}
	m.visitRepo.On("List", ctx, 20, 0).Return(expected, nil).Once()

	got, err := m.service.List(ctx, 20, 0)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}
