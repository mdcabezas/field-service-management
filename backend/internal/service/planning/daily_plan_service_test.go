package planning

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/customers"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	coreservice "localis-backend/internal/service/core"
	"localis-backend/internal/testutil/mocks"
)

type dailyPlanSvcMocks struct {
	planRepo       *mocks.DailyPlanRepository
	assignmentRepo *mocks.DailyPlanAssignmentRepository
	techRepo       *mocks.TechnicianRepository
	auditRepo      *mocks.CorePlanAuditLogRepository
	svc            *DailyPlanService
}

func newDailyPlanService(t *testing.T) *dailyPlanSvcMocks {
	t.Helper()
	planRepo := mocks.NewDailyPlanRepository(t)
	assignmentRepo := mocks.NewDailyPlanAssignmentRepository(t)
	techRepo := mocks.NewTechnicianRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	svc := NewDailyPlanService(planRepo, assignmentRepo, techRepo, coreservice.NewAuditService(auditRepo))
	return &dailyPlanSvcMocks{planRepo: planRepo, assignmentRepo: assignmentRepo, techRepo: techRepo, auditRepo: auditRepo, svc: svc}
}

func TestCreatePlan_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	plan := &planning.DailyPlan{Name: "Plan A"}
	m.planRepo.On("Create", ctx, plan).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.CreatePlan(ctx, plan)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, plan.ID)
}

func TestCreatePlan_RepoError(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	m.planRepo.On("Create", ctx, mock.Anything).Return(errors.New("db down")).Once()

	err := m.svc.CreatePlan(ctx, &planning.DailyPlan{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "create daily plan")
}

func TestAssignTechnician_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	assignment := &planning.DailyPlanAssignment{DailyPlanID: uuid.New(), TechID: uuid.New(), RoleID: uuid.New()}
	tech := &customers.Technician{ID: assignment.TechID, IsActive: true}

	m.planRepo.On("GetByID", ctx, assignment.DailyPlanID).Return(&planning.DailyPlan{ID: assignment.DailyPlanID}, nil).Once()
	m.techRepo.On("GetByID", ctx, assignment.TechID).Return(tech, nil).Once()
	m.assignmentRepo.On("Create", ctx, assignment).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.AssignTechnician(ctx, assignment)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, assignment.ID)
}

func TestAssignTechnician_Inactive(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	assignment := &planning.DailyPlanAssignment{DailyPlanID: uuid.New(), TechID: uuid.New()}
	tech := &customers.Technician{ID: assignment.TechID, IsActive: false}

	m.planRepo.On("GetByID", ctx, assignment.DailyPlanID).Return(&planning.DailyPlan{ID: assignment.DailyPlanID}, nil).Once()
	m.techRepo.On("GetByID", ctx, assignment.TechID).Return(tech, nil).Once()

	err := m.svc.AssignTechnician(ctx, assignment)
	var ve *service.ValidationError
	require.ErrorAs(t, err, &ve)
	require.Equal(t, "tech_id", ve.Field)
}

func TestAssignTechnician_PlanNotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	assignment := &planning.DailyPlanAssignment{DailyPlanID: uuid.New(), TechID: uuid.New()}
	m.planRepo.On("GetByID", ctx, assignment.DailyPlanID).Return(nil, &service.NotFoundError{Resource: "daily_plan", ID: assignment.DailyPlanID.String()}).Once()

	err := m.svc.AssignTechnician(ctx, assignment)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
	require.Equal(t, "daily_plan", nf.Resource)
}

func TestAssignTechnician_TechNotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	assignment := &planning.DailyPlanAssignment{DailyPlanID: uuid.New(), TechID: uuid.New()}
	m.planRepo.On("GetByID", ctx, assignment.DailyPlanID).Return(&planning.DailyPlan{ID: assignment.DailyPlanID}, nil).Once()
	m.techRepo.On("GetByID", ctx, assignment.TechID).Return(nil, &service.NotFoundError{Resource: "technician", ID: assignment.TechID.String()}).Once()

	err := m.svc.AssignTechnician(ctx, assignment)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
	require.Equal(t, "technician", nf.Resource)
}

func TestCloseDailyPlan_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	plan := &planning.DailyPlan{ID: id, Name: "Plan A"}
	m.planRepo.On("GetByID", ctx, id).Return(plan, nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.CloseDailyPlan(ctx, id)
	require.NoError(t, err)
}

func TestCloseDailyPlan_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	m.planRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "daily_plan", ID: id.String()}).Once()

	err := m.svc.CloseDailyPlan(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestGetByID_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	expected := &planning.DailyPlan{ID: id}
	m.planRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

	got, err := m.svc.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestList_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	expected := &repository.ListResult[planning.DailyPlan]{Items: []planning.DailyPlan{{ID: uuid.New()}}, Total: 1}
	m.planRepo.On("List", ctx, 10, 0).Return(expected, nil).Once()

	got, err := m.svc.List(ctx, 10, 0)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestUpdate_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	updated := &planning.DailyPlan{ID: id, Name: "Plan B"}
	m.planRepo.On("GetByID", ctx, id).Return(&planning.DailyPlan{ID: id, Name: "Plan A"}, nil).Once()
	m.planRepo.On("Update", ctx, id, updated).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.Update(ctx, id, updated)
	require.NoError(t, err)
}

func TestUpdate_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	m.planRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "daily_plan", ID: id.String()}).Once()

	err := m.svc.Update(ctx, id, &planning.DailyPlan{})
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestDelete_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	m.planRepo.On("GetByID", ctx, id).Return(&planning.DailyPlan{ID: id}, nil).Once()
	m.planRepo.On("Delete", ctx, id).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.Delete(ctx, id)
	require.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	m.planRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "daily_plan", ID: id.String()}).Once()

	err := m.svc.Delete(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestListAssignments_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	planID := uuid.New()
	expected := []planning.DailyPlanAssignment{{ID: uuid.New(), DailyPlanID: planID}}
	m.assignmentRepo.On("ListByPlan", ctx, planID).Return(expected, nil).Once()

	got, err := m.svc.ListAssignments(ctx, planID)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestDeleteAssignment_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	m.assignmentRepo.On("GetByID", ctx, id).Return(&planning.DailyPlanAssignment{ID: id}, nil).Once()
	m.assignmentRepo.On("Delete", ctx, id).Return(nil).Once()
	m.auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := m.svc.DeleteAssignment(ctx, id)
	require.NoError(t, err)
}

func TestDeleteAssignment_NotFound(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	m.assignmentRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "daily_plan_assignment", ID: id.String()}).Once()

	err := m.svc.DeleteAssignment(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestGetAssignmentByID_Success(t *testing.T) {
	ctx := context.Background()
	m := newDailyPlanService(t)

	id := uuid.New()
	expected := &planning.DailyPlanAssignment{ID: id}
	m.assignmentRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

	got, err := m.svc.GetAssignmentByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}
