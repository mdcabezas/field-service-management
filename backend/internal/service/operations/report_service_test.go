package operations

import (
	"context"
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

func newReportService(t *testing.T) (*ReportService, *mocks.VisitReportRepository, *mocks.ReportEntryRepository, *mocks.VisitPhotoRepository, *mocks.VisitMeasurementRepository, *mocks.CorePlanAuditLogRepository) {
	t.Helper()
	visitReportRepo := mocks.NewVisitReportRepository(t)
	reportEntryRepo := mocks.NewReportEntryRepository(t)
	reportImageRepo := mocks.NewReportImageRepository(t)
	visitPhotoRepo := mocks.NewVisitPhotoRepository(t)
	visitMeasureRepo := mocks.NewVisitMeasurementRepository(t)
	reportTmplRepo := mocks.NewReportTemplateRepository(t)
	auditRepo := mocks.NewCorePlanAuditLogRepository(t)
	svc := NewReportService(visitReportRepo, reportEntryRepo, reportImageRepo, visitPhotoRepo, visitMeasureRepo, reportTmplRepo, core.NewAuditService(auditRepo))
	return svc, visitReportRepo, reportEntryRepo, visitPhotoRepo, visitMeasureRepo, auditRepo
}

func TestGenerateReport(t *testing.T) {
	ctx := context.Background()
	svc, visitReportRepo, reportEntryRepo, _, _, auditRepo := newReportService(t)

	visitID := uuid.New()
	templateID := uuid.New()

	visitReportRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
	reportEntryRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	report, err := svc.GenerateReport(ctx, visitID, templateID)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, report.ID)
	require.Equal(t, visitID, report.VisitID)
}

func TestFinalizeReport(t *testing.T) {
	ctx := context.Background()
	svc, visitReportRepo, reportEntryRepo, visitPhotoRepo, visitMeasureRepo, auditRepo := newReportService(t)

	reportID := uuid.New()
	report := &operations.VisitReport{ID: reportID, VisitID: uuid.New(), Source: shared.ReportSourceSystem, RecordedAt: time.Now()}

	visitReportRepo.On("GetByID", ctx, reportID).Return(report, nil).Once()
	reportEntryRepo.On("ListByReport", ctx, reportID).Return([]operations.ReportEntry{}, nil).Once()
	visitPhotoRepo.On("ListByVisit", ctx, report.VisitID).Return([]operations.VisitPhoto{}, nil).Once()
	visitMeasureRepo.On("ListByVisit", ctx, report.VisitID).Return([]operations.VisitMeasurement{}, nil).Once()
	reportEntryRepo.On("Create", ctx, mock.Anything).Return(nil).Once()
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := svc.FinalizeReport(ctx, reportID)
	require.NoError(t, err)
}

func TestFinalizeReport_AlreadyFinalized(t *testing.T) {
	ctx := context.Background()
	svc, visitReportRepo, reportEntryRepo, _, _, _ := newReportService(t)

	reportID := uuid.New()
	report := &operations.VisitReport{ID: reportID, VisitID: uuid.New(), Source: shared.ReportSourceSystem, RecordedAt: time.Now()}
	entry := &operations.ReportEntry{ID: uuid.New(), ReportID: reportID, DataJSON: []byte(`{"type":"final_summary"}`)}

	visitReportRepo.On("GetByID", ctx, reportID).Return(report, nil).Once()
	reportEntryRepo.On("ListByReport", ctx, reportID).Return([]operations.ReportEntry{*entry}, nil).Once()

	err := svc.FinalizeReport(ctx, reportID)
	require.NoError(t, err)
}

func TestReport_Delete_NotFound(t *testing.T) {
	ctx := context.Background()
	svc, visitReportRepo, _, _, _, _ := newReportService(t)

	id := uuid.New()
	visitReportRepo.On("GetByID", ctx, id).Return(nil, &service.NotFoundError{Resource: "visit_report", ID: id.String()}).Once()

	err := svc.Delete(ctx, id)
	var nf *service.NotFoundError
	require.ErrorAs(t, err, &nf)
}

func TestReportService_GetByID(t *testing.T) {
	ctx := context.Background()
	svc, visitReportRepo, _, _, _, _ := newReportService(t)

	id := uuid.New()
	expected := &operations.VisitReport{ID: id, VisitID: uuid.New(), Source: shared.ReportSourceSystem, RecordedAt: time.Now()}
	visitReportRepo.On("GetByID", ctx, id).Return(expected, nil).Once()

	got, err := svc.GetByID(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestReportService_List(t *testing.T) {
	ctx := context.Background()
	svc, visitReportRepo, _, _, _, _ := newReportService(t)

	expected := &repository.ListResult[operations.VisitReport]{
		Items: []operations.VisitReport{{ID: uuid.New(), VisitID: uuid.New()}},
		Total: 1,
	}
	visitReportRepo.On("List", ctx, 20, 0).Return(expected, nil).Once()

	got, err := svc.List(ctx, 20, 0)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestReportService_ListByVisit(t *testing.T) {
	ctx := context.Background()
	svc, visitReportRepo, _, _, _, _ := newReportService(t)

	visitID := uuid.New()
	expected := []operations.VisitReport{{ID: uuid.New(), VisitID: visitID}}
	visitReportRepo.On("ListByVisit", ctx, visitID).Return(expected, nil).Once()

	got, err := svc.ListByVisit(ctx, visitID)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestReportService_Delete_Success(t *testing.T) {
	ctx := context.Background()
	svc, visitReportRepo, _, _, _, auditRepo := newReportService(t)

	id := uuid.New()
	report := &operations.VisitReport{ID: id}
	visitReportRepo.On("GetByID", ctx, id).Return(report, nil).Once()
	visitReportRepo.On("Delete", ctx, id).Return(nil).Once()
	auditRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
}
