package operations

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	operationsmodel "localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	coreservice "localis-backend/internal/service/core"
)

type ReportService struct {
	visitReportRepo  repository.VisitReportRepository
	reportEntryRepo  repository.ReportEntryRepository
	reportImageRepo  repository.ReportImageRepository
	visitPhotoRepo   repository.VisitPhotoRepository
	visitMeasureRepo repository.VisitMeasurementRepository
	reportTmplRepo   repository.ReportTemplateRepository
	audit            *coreservice.AuditService
}

func NewReportService(
	visitReportRepo repository.VisitReportRepository,
	reportEntryRepo repository.ReportEntryRepository,
	reportImageRepo repository.ReportImageRepository,
	visitPhotoRepo repository.VisitPhotoRepository,
	visitMeasureRepo repository.VisitMeasurementRepository,
	reportTmplRepo repository.ReportTemplateRepository,
	audit *coreservice.AuditService,
) *ReportService {
	return &ReportService{
		visitReportRepo:  visitReportRepo,
		reportEntryRepo:  reportEntryRepo,
		reportImageRepo:  reportImageRepo,
		visitPhotoRepo:   visitPhotoRepo,
		visitMeasureRepo: visitMeasureRepo,
		reportTmplRepo:   reportTmplRepo,
		audit:            audit,
	}
}

func (s *ReportService) GenerateReport(ctx context.Context, visitID, templateID uuid.UUID, source shared.ReportSource, recordedAt time.Time) (*operationsmodel.VisitReport, error) {
	report := &operationsmodel.VisitReport{
		ID:               uuid.New(),
		VisitID:          visitID,
		ReportTemplateID: &templateID,
		Source:           source,
		RecordedAt:       recordedAt,
		CreatedAt:        time.Now(),
	}

	if err := s.visitReportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("create visit report: %w", err)
	}

	summary := map[string]any{
		"type":     "initial",
		"visit_id": visitID.String(),
	}
	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		return report, fmt.Errorf("marshal report summary: %w", err)
	}

	entry := &operationsmodel.ReportEntry{
		ID:        uuid.New(),
		ReportID:  report.ID,
		DataJSON:  summaryJSON,
		CreatedAt: time.Now(),
	}

	if err := s.reportEntryRepo.Create(ctx, entry); err != nil {
		return report, fmt.Errorf("create report entry: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitReport, report.ID, shared.AuditActionAdd, nil, report, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_report", "id", report.ID, "error", err)
	}

	return report, nil
}

func (s *ReportService) FinalizeReport(ctx context.Context, reportID uuid.UUID) error {
	report, err := s.visitReportRepo.GetByID(ctx, reportID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit_report", reportID.String())
	}

	existingEntries, err := s.reportEntryRepo.ListByReport(ctx, reportID)
	if err != nil {
		return fmt.Errorf("list report entries: %w", err)
	}
	for _, e := range existingEntries {
		var data map[string]any
		if jsonErr := json.Unmarshal(e.DataJSON, &data); jsonErr == nil {
			if t, ok := data["type"].(string); ok && t == "final_summary" {
				return nil
			}
		}
	}

	photos, err := s.visitPhotoRepo.ListByVisit(ctx, report.VisitID)
	if err != nil {
		return fmt.Errorf("list visit photos: %w", err)
	}
	measurements, err := s.visitMeasureRepo.ListByVisit(ctx, report.VisitID)
	if err != nil {
		return fmt.Errorf("list visit measurements: %w", err)
	}

	summary := map[string]any{
		"type":              "final_summary",
		"visit_id":          report.VisitID.String(),
		"photo_count":       len(photos),
		"measurement_count": len(measurements),
	}

	if len(measurements) > 0 {
		var measList []map[string]any
		for _, m := range measurements {
			measList = append(measList, map[string]any{
				"type":   m.Type.String(),
				"value":  m.Value,
				"unit":   m.Unit,
				"result": func() string { if m.Result != nil { return string(*m.Result) }; return "" }(),
			})
		}
		summary["measurements"] = measList
	}

	summaryJSON, err := json.Marshal(summary)
	if err != nil {
		return fmt.Errorf("marshal report summary: %w", err)
	}

	entry := &operationsmodel.ReportEntry{
		ID:        uuid.New(),
		ReportID:  reportID,
		DataJSON:  summaryJSON,
		CreatedAt: time.Now(),
	}

	if err := s.reportEntryRepo.Create(ctx, entry); err != nil {
		return fmt.Errorf("create report entry: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitReport, reportID, shared.AuditActionUpdate, nil, entry, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_report", "id", reportID, "error", err)
	}

	return nil
}

func (s *ReportService) GetByID(ctx context.Context, id uuid.UUID) (*operationsmodel.VisitReport, error) {
	item, err := s.visitReportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get visit report: %w", err)
	}
	return item, nil
}

func (s *ReportService) List(ctx context.Context, limit, offset int) (*repository.ListResult[operationsmodel.VisitReport], error) {
	result, err := s.visitReportRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list visit reports: %w", err)
	}
	return result, nil
}

func (s *ReportService) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operationsmodel.VisitReport, error) {
	items, err := s.visitReportRepo.ListByVisit(ctx, visitID)
	if err != nil {
		return nil, fmt.Errorf("list visit reports: %w", err)
	}
	return items, nil
}

func (s *ReportService) Delete(ctx context.Context, id uuid.UUID) error {
	report, err := s.visitReportRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit_report", id.String())
	}

	if err := s.reportEntryRepo.DeleteByReport(ctx, id); err != nil {
		return fmt.Errorf("delete report entries: %w", err)
	}
	if err := s.reportImageRepo.DeleteByReport(ctx, id); err != nil {
		return fmt.Errorf("delete report images: %w", err)
	}

	if err := s.visitReportRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete visit report: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitReport, id, shared.AuditActionDelete, report, nil, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_report", "id", id, "error", err)
	}

	return nil
}
