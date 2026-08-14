package operations

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/partners"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	"localis-backend/internal/service/core"
	notifsvc "localis-backend/internal/service/notifications"
	"localis-backend/internal/testutil/mocks"
)

func newSLAService(t *testing.T) (*SLATrackingService, *mocks.VisitSLATrackingRepository, *mocks.SLARepository, *mocks.VisitRepository, *mocks.CorePlanAuditLogRepository) {
	t.Helper()
	slaTrackingRepo := mocks.NewVisitSLATrackingRepository(t)
	slaRepo := mocks.NewSLARepository(t)
	visitRepo := mocks.NewVisitRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	notifRepo := mocks.NewNotificationRepository(t)
	notifRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Maybe()
	notifRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
	notifSvc := notifsvc.NewNotificationService(notifRepo)
	svc := NewSLATrackingService(slaTrackingRepo, slaRepo, visitRepo, core.NewAuditService(auditRepo), notifSvc)
	return svc, slaTrackingRepo, slaRepo, visitRepo, auditRepo
}

func TestSLACreateTracking_Success(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, _, visitRepo, auditRepo := newSLAService(t)

	tracking := &operations.VisitSLATracking{VisitID: uuid.New(), SLAID: uuid.New()}
	visit := newVisitFixture(tracking.VisitID, shared.VisitStatusScheduled)

	visitRepo.On("GetByID", ctx, tracking.VisitID).Return(visit, nil).Once()
	slaTrackingRepo.On("Create", ctx, tracking).Return(nil).Once()
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := svc.CreateTracking(ctx, tracking)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, tracking.ID)
	require.NotNil(t, tracking.RequestedAt)
}

func TestSLACreateTracking_VisitNotFound(t *testing.T) {
	ctx := context.Background()
	svc, _, _, visitRepo, _ := newSLAService(t)

	tracking := &operations.VisitSLATracking{VisitID: uuid.New(), SLAID: uuid.New()}
	visitRepo.On("GetByID", ctx, tracking.VisitID).Return(nil, &service.NotFoundError{Resource: "visit", ID: tracking.VisitID.String()}).Once()

	err := svc.CreateTracking(ctx, tracking)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestSLAUpdateResponseTime_Meets(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, slaRepo, _, auditRepo := newSLAService(t)

	id := uuid.New()
	requested := time.Now().Add(-2 * time.Hour)
	tracking := &operations.VisitSLATracking{ID: id, VisitID: uuid.New(), SLAID: uuid.New(), RequestedAt: &requested}
	sla := &partners.SLA{ID: tracking.SLAID, Name: "Std SLA"}
	responseHours := 24
	sla.ResponseHours = &responseHours

	slaTrackingRepo.On("GetByID", ctx, id).Return(tracking, nil).Once()
	slaRepo.On("GetByID", ctx, tracking.SLAID).Return(sla, nil).Once()
	slaTrackingRepo.On("Update", ctx, id, tracking).Return(nil).Once()
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := svc.UpdateResponseTime(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, tracking.RespondedAt)
	require.NotNil(t, tracking.ResponseTimeHours)
	require.NotNil(t, tracking.MeetsResponseSLA)
	require.True(t, *tracking.MeetsResponseSLA)
}

func TestSLAUpdateResponseTime_NotFound(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, _, _, _ := newSLAService(t)

	id := uuid.New()
	slaTrackingRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "visit_sla_tracking", ID: id.String()}).Once()

	err := svc.UpdateResponseTime(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestSLADelete_Success(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, _, _, auditRepo := newSLAService(t)

	id := uuid.New()
	tracking := &operations.VisitSLATracking{ID: id}
	slaTrackingRepo.On("GetByID", ctx, id).Return(tracking, nil).Once()
	slaTrackingRepo.On("Delete", ctx, id).Return(nil).Once()
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
}

func TestSLAUpdateResolutionTime_Meets(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, slaRepo, _, auditRepo := newSLAService(t)

	id := uuid.New()
	requested := time.Now().Add(-3 * time.Hour)
	resolved := time.Now().Add(-1 * time.Hour)
	tracking := &operations.VisitSLATracking{ID: id, VisitID: uuid.New(), SLAID: uuid.New(), RequestedAt: &requested, ResolvedAt: &resolved}
	sla := &partners.SLA{ID: tracking.SLAID, Name: "Std SLA"}
	resolutionHours := 4
	sla.ResolutionHours = &resolutionHours

	slaTrackingRepo.On("GetByID", ctx, id).Return(tracking, nil).Once()
	slaRepo.On("GetByID", ctx, tracking.SLAID).Return(sla, nil).Once()
	slaTrackingRepo.On("Update", ctx, id, tracking).Return(nil).Once()
	auditRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Maybe()

	err := svc.UpdateResolutionTime(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, tracking.ResolutionTimeHours)
	require.NotNil(t, tracking.MeetsResolutionSLA)
	require.True(t, *tracking.MeetsResolutionSLA)
}

func TestSLAUpdateResolutionTime_NotFound(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, _, _, _ := newSLAService(t)

	id := uuid.New()
	slaTrackingRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "visit_sla_tracking", ID: id.String()}).Once()

	err := svc.UpdateResolutionTime(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestSLAUpdate(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, _, _, auditRepo := newSLAService(t)

	id := uuid.New()
	existing := &operations.VisitSLATracking{ID: id}
	updated := &operations.VisitSLATracking{ID: id, VisitID: uuid.New()}

	slaTrackingRepo.On("GetByID", ctx, id).Return(existing, nil).Once()
	slaTrackingRepo.On("Update", ctx, id, updated).Return(nil).Once()
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := svc.Update(ctx, id, updated)
	require.NoError(t, err)
}

func TestSLACheckSLAWarnings(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, _ := newSLAService(t)

	// CheckSLAWarnings is not fully implemented (returns error)
	_, err := svc.CheckSLAWarnings(ctx)
	require.Error(t, err)
}

func TestSLAService_GetByID(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, _, _, _ := newSLAService(t)

	id := uuid.New()
	expected := &operations.VisitSLATracking{ID: id, VisitID: uuid.New(), SLAID: uuid.New()}
	slaTrackingRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

	got, err := svc.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestSLAService_List(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, _, _, _ := newSLAService(t)

	expected := &repository.ListResult[operations.VisitSLATracking]{
		Items: []operations.VisitSLATracking{{ID: uuid.New(), VisitID: uuid.New()}},
		Total: 1,
	}
	slaTrackingRepo.On("List", ctx, 20, 0).Return(expected, nil).Once()

	got, err := svc.List(ctx, 20, 0)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestSLAService_ListByVisit(t *testing.T) {
	ctx := context.Background()
	svc, slaTrackingRepo, _, _, _ := newSLAService(t)

	visitID := uuid.New()
	expected := []operations.VisitSLATracking{{ID: uuid.New(), VisitID: visitID}}
	slaTrackingRepo.On("ListByVisit", ctx, visitID).Return(expected, nil).Once()

	got, err := svc.ListByVisit(ctx, visitID)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}
