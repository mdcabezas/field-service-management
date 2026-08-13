package core

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	coremodel "localis-backend/internal/model/core"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/service"
	"localis-backend/internal/testutil/mocks"
)

func newAuditService(t *testing.T) (*AuditService, *mocks.CorePlanAuditLogRepository) {
	t.Helper()
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	return NewAuditService(auditRepo), auditRepo
}

func TestRecordChange_Success(t *testing.T) {
	ctx := service.WithUserID(context.Background(), uuid.New().String())
	svc, auditRepo := newAuditService(t)

	entityID := uuid.New()
	oldData := map[string]any{"name": "old"}
	newData := map[string]any{"name": "new"}
	reason := "updated by test"

	var captured *coremodel.PlanAuditLog
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once().
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(*coremodel.PlanAuditLog)
		})

	err := svc.RecordChange(ctx, shared.AuditEntityTypeVisit, entityID, shared.AuditActionUpdate, oldData, newData, &reason)
	require.NoError(t, err)
	require.NotNil(t, captured)
	require.Equal(t, shared.AuditEntityTypeVisit, captured.EntityType)
	require.Equal(t, entityID, captured.EntityID)
	require.Equal(t, shared.AuditActionUpdate, captured.Action)
	require.JSONEq(t, `{"name":"old"}`, string(captured.OldData))
	require.JSONEq(t, `{"name":"new"}`, string(captured.NewData))
	require.NotNil(t, captured.UserID)
	require.Equal(t, service.UserIDFromCtx(ctx), captured.UserID.String())
	require.NotNil(t, captured.Reason)
	require.Equal(t, reason, *captured.Reason)
	require.NotZero(t, captured.Timestamp)
}

func TestRecordChange_NoUserInCtx(t *testing.T) {
	ctx := context.Background()
	svc, auditRepo := newAuditService(t)

	var captured *coremodel.PlanAuditLog
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once().
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(*coremodel.PlanAuditLog)
		})

	err := svc.RecordChange(ctx, shared.AuditEntityTypeVisit, uuid.New(), shared.AuditActionAdd, nil, nil, nil)
	require.NoError(t, err)
	require.Nil(t, captured.UserID)
	require.Nil(t, captured.Reason)
	require.Empty(t, captured.OldData)
	require.Empty(t, captured.NewData)
}

func TestRecordChange_InvalidUserIDInCtx(t *testing.T) {
	ctx := service.WithUserID(context.Background(), "not-a-uuid")
	svc, auditRepo := newAuditService(t)

	var captured *coremodel.PlanAuditLog
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once().
		Run(func(args mock.Arguments) {
			captured = args.Get(1).(*coremodel.PlanAuditLog)
		})

	err := svc.RecordChange(ctx, shared.AuditEntityTypeVisit, uuid.New(), shared.AuditActionDelete, nil, nil, nil)
	require.NoError(t, err)
	require.Nil(t, captured.UserID)
}

func TestRecordChange_MarshalError(t *testing.T) {
	svc, _ := newAuditService(t)

	err := svc.RecordChange(context.Background(), shared.AuditEntityTypeVisit, uuid.New(), shared.AuditActionAdd, make(chan int), nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "marshal old audit data")
}

func TestRecordChange_RepoError(t *testing.T) {
	ctx := context.Background()
	svc, auditRepo := newAuditService(t)

	auditRepo.On("Create", ctx, mock.Anything).Return(errors.New("db down")).Once()

	err := svc.RecordChange(ctx, shared.AuditEntityTypeVisit, uuid.New(), shared.AuditActionAdd, nil, nil, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "create audit log entry")
}

func TestListByEntity_Success(t *testing.T) {
	ctx := context.Background()
	svc, auditRepo := newAuditService(t)

	entityID := uuid.New()
	expected := []coremodel.PlanAuditLog{
		{ID: uuid.New(), EntityType: shared.AuditEntityTypeVisit, EntityID: entityID, Action: shared.AuditActionAdd},
	}
	auditRepo.On("ListByEntity", ctx, shared.AuditEntityTypeVisit, entityID).Return(expected, nil).Once()

	got, err := svc.ListByEntity(ctx, shared.AuditEntityTypeVisit, entityID)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestListByEntity_RepoError(t *testing.T) {
	ctx := context.Background()
	svc, auditRepo := newAuditService(t)

	auditRepo.On("ListByEntity", ctx, shared.AuditEntityTypeVisit, uuid.Nil).Return(nil, errors.New("db down")).Once()

	_, err := svc.ListByEntity(ctx, shared.AuditEntityTypeVisit, uuid.Nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "list audit logs by entity")
}
