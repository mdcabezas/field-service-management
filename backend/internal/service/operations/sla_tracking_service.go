package operations

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	coreservice "localis-backend/internal/service/core"
	notifservice "localis-backend/internal/service/notifications"
)

type SLATrackingService struct {
	slaTrackingRepo repository.VisitSLATrackingRepository
	slaRepo         repository.SLARepository
	visitRepo       repository.VisitRepository
	audit           *coreservice.AuditService
	notifService    *notifservice.NotificationService
}

func NewSLATrackingService(
	slaTrackingRepo repository.VisitSLATrackingRepository,
	slaRepo repository.SLARepository,
	visitRepo repository.VisitRepository,
	audit *coreservice.AuditService,
	notifService *notifservice.NotificationService,
) *SLATrackingService {
	return &SLATrackingService{
		slaTrackingRepo: slaTrackingRepo,
		slaRepo:         slaRepo,
		visitRepo:       visitRepo,
		audit:           audit,
		notifService:    notifService,
	}
}

func (s *SLATrackingService) CreateTracking(ctx context.Context, tracking *operations.VisitSLATracking) error {
	if _, err := s.visitRepo.GetByID(ctx, tracking.VisitID); err != nil {
		return service.HandleRepoGetByIDError(err, "visit", tracking.VisitID.String())
	}

	now := time.Now()
	tracking.ID = uuid.New()
	tracking.RequestedAt = &now
	tracking.CreatedAt = now

	if err := s.slaTrackingRepo.Create(ctx, tracking); err != nil {
		return fmt.Errorf("create sla tracking: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitSLATracking, tracking.ID, shared.AuditActionAdd, nil, tracking, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_sla_tracking", "id", tracking.ID, "error", err)
	}

	return nil
}

func (s *SLATrackingService) UpdateResponseTime(ctx context.Context, id uuid.UUID) error {
	tracking, err := s.slaTrackingRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit_sla_tracking", id.String())
	}

	now := time.Now()
	tracking.RespondedAt = &now

	if tracking.RequestedAt != nil {
		hours := now.Sub(*tracking.RequestedAt).Hours()
		tracking.ResponseTimeHours = &hours
	}

	sla, err := s.slaRepo.GetByID(ctx, tracking.SLAID)
	if err != nil {
		return fmt.Errorf("get sla: %w", err)
	}

	if sla.ResponseHours != nil && tracking.ResponseTimeHours != nil {
		meets := *tracking.ResponseTimeHours <= float64(*sla.ResponseHours)
		tracking.MeetsResponseSLA = &meets

		if !meets {
			if err := s.notifService.NotifySLAWarning(ctx, tracking.VisitID, sla.ID, sla.Name, now); err != nil {
				slog.Warn("failed to send SLA warning", "visit_id", tracking.VisitID, "error", err)
			}
		}
	}

	if err := s.slaTrackingRepo.Update(ctx, id, tracking); err != nil {
		return fmt.Errorf("update sla tracking response time: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitSLATracking, id, shared.AuditActionUpdate, nil, tracking, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_sla_tracking", "id", id, "error", err)
	}

	return nil
}

func (s *SLATrackingService) UpdateResolutionTime(ctx context.Context, id uuid.UUID) error {
	tracking, err := s.slaTrackingRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit_sla_tracking", id.String())
	}

	now := time.Now()
	tracking.ResolvedAt = &now

	if tracking.RequestedAt != nil {
		hours := now.Sub(*tracking.RequestedAt).Hours()
		tracking.ResolutionTimeHours = &hours
	}

	sla, err := s.slaRepo.GetByID(ctx, tracking.SLAID)
	if err != nil {
		return fmt.Errorf("get sla: %w", err)
	}

	if sla.ResolutionHours != nil && tracking.ResolutionTimeHours != nil {
		meets := *tracking.ResolutionTimeHours <= float64(*sla.ResolutionHours)
		tracking.MeetsResolutionSLA = &meets

		if !meets {
			if err := s.notifService.NotifySLAWarning(ctx, tracking.VisitID, sla.ID, sla.Name, now); err != nil {
				slog.Warn("failed to send SLA warning", "visit_id", tracking.VisitID, "error", err)
			}
		}
	}

	if err := s.slaTrackingRepo.Update(ctx, id, tracking); err != nil {
		return fmt.Errorf("update sla tracking resolution time: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitSLATracking, id, shared.AuditActionUpdate, nil, tracking, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_sla_tracking", "id", id, "error", err)
	}

	return nil
}

func (s *SLATrackingService) Update(ctx context.Context, id uuid.UUID, tracking *operations.VisitSLATracking) error {
	existing, err := s.slaTrackingRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit_sla_tracking", id.String())
	}

	if tracking.VisitID != uuid.Nil {
		existing.VisitID = tracking.VisitID
	}
	if tracking.SLAID != uuid.Nil {
		existing.SLAID = tracking.SLAID
	}
	if tracking.RequestedAt != nil {
		existing.RequestedAt = tracking.RequestedAt
	}
	if tracking.RespondedAt != nil {
		existing.RespondedAt = tracking.RespondedAt
	}
	if tracking.ResolvedAt != nil {
		existing.ResolvedAt = tracking.ResolvedAt
	}
	if tracking.ResponseTimeHours != nil {
		existing.ResponseTimeHours = tracking.ResponseTimeHours
	}
	if tracking.ResolutionTimeHours != nil {
		existing.ResolutionTimeHours = tracking.ResolutionTimeHours
	}
	if tracking.MeetsResponseSLA != nil {
		existing.MeetsResponseSLA = tracking.MeetsResponseSLA
	}
	if tracking.MeetsResolutionSLA != nil {
		existing.MeetsResolutionSLA = tracking.MeetsResolutionSLA
	}

	if err := s.slaTrackingRepo.Update(ctx, id, existing); err != nil {
		return fmt.Errorf("update sla tracking: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitSLATracking, id, shared.AuditActionUpdate, existing, tracking, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_sla_tracking", "id", id, "error", err)
	}
	return nil
}

func (s *SLATrackingService) CheckSLAWarnings(ctx context.Context) (int, error) {
	return 0, fmt.Errorf("CheckSLAWarnings not implemented: requires ListAll method on VisitSLATrackingRepository")
}

func (s *SLATrackingService) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitSLATracking, error) {
	item, err := s.slaTrackingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get sla tracking: %w", err)
	}
	return item, nil
}

func (s *SLATrackingService) List(ctx context.Context, limit, offset int) (*repository.ListResult[operations.VisitSLATracking], error) {
	result, err := s.slaTrackingRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list sla trackings: %w", err)
	}
	return result, nil
}

func (s *SLATrackingService) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitSLATracking, error) {
	items, err := s.slaTrackingRepo.ListByVisit(ctx, visitID)
	if err != nil {
		return nil, fmt.Errorf("list sla tracking by visit: %w", err)
	}
	return items, nil
}

func (s *SLATrackingService) Delete(ctx context.Context, id uuid.UUID) error {
	tracking, err := s.slaTrackingRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit_sla_tracking", id.String())
	}

	if err := s.slaTrackingRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete sla tracking: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisitSLATracking, id, shared.AuditActionDelete, tracking, nil, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit_sla_tracking", "id", id, "error", err)
	}

	return nil
}
