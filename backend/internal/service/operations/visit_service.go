package operations

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	"localis-backend/internal/service/core"
)

type VisitServiceDeps struct {
	VisitRepo      repository.VisitRepository
	AssignmentRepo repository.VisitAssignmentRepository
	TechnicianRepo repository.TechnicianRepository
	CheckpointRepo repository.VisitCheckpointRepository
	CTMRepo        repository.ChecklistTemplateMaterialRepository
	CTTRepo        repository.ChecklistTemplateToolRepository
	CTARepo        repository.ChecklistTemplateEPPRepository
	VCMRepo        repository.VisitChecklistMaterialRepository
	VCTRepo        repository.VisitChecklistToolRepository
	VCERepo        repository.VisitChecklistEPPRepository
	AuditService   *core.AuditService
}

type VisitService struct {
	visitRepo         repository.VisitRepository
	assignmentRepo    repository.VisitAssignmentRepository
	techRepo          repository.TechnicianRepository
	checkpointRepo    repository.VisitCheckpointRepository
	checklistTmplMat  repository.ChecklistTemplateMaterialRepository
	checklistTmplTool repository.ChecklistTemplateToolRepository
	checklistTmplEPP  repository.ChecklistTemplateEPPRepository
	visitCheckMat     repository.VisitChecklistMaterialRepository
	visitCheckTool    repository.VisitChecklistToolRepository
	visitCheckEPP     repository.VisitChecklistEPPRepository
	audit             *core.AuditService
	vPool             *pgxpool.Pool
}

func NewVisitService(deps VisitServiceDeps) *VisitService {
	s := &VisitService{
		visitRepo:         deps.VisitRepo,
		assignmentRepo:    deps.AssignmentRepo,
		techRepo:          deps.TechnicianRepo,
		checkpointRepo:    deps.CheckpointRepo,
		checklistTmplMat:  deps.CTMRepo,
		checklistTmplTool: deps.CTTRepo,
		checklistTmplEPP:  deps.CTARepo,
		visitCheckMat:     deps.VCMRepo,
		visitCheckTool:    deps.VCTRepo,
		visitCheckEPP:     deps.VCERepo,
		audit:             deps.AuditService,
	}
	if p, ok := deps.VisitRepo.(interface{ Pool() *pgxpool.Pool }); ok {
		s.vPool = p.Pool()
	}
	return s
}

var validTransitions = map[shared.VisitStatus][]shared.VisitStatus{
	shared.VisitStatusScheduled:   {shared.VisitStatusEnRoute, shared.VisitStatusCancelled, shared.VisitStatusRescheduled},
	shared.VisitStatusEnRoute:     {shared.VisitStatusInProgress, shared.VisitStatusCancelled},
	shared.VisitStatusInProgress:  {shared.VisitStatusCompleted, shared.VisitStatusCancelled},
	shared.VisitStatusRescheduled: {shared.VisitStatusScheduled, shared.VisitStatusCancelled},
}

func canTransition(from, to shared.VisitStatus) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func (s *VisitService) CreateVisit(ctx context.Context, visit *operations.Visit) error {
	visit.ID = uuid.New()
	visit.Status = shared.VisitStatusScheduled
	visit.Priority = shared.VisitPriorityNormal
	now := time.Now()
	visit.CreatedAt = now
	visit.UpdatedAt = now

	if err := s.visitRepo.Create(ctx, visit); err != nil {
		return fmt.Errorf("create visit: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visit.ID, shared.AuditActionAdd, nil, visit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visit.ID, "error", err)
	}

	return nil
}

// NOTE: NoTx fallback paths are for development/testing only.
// In production, vPool is always non-nil (extracted from pgx-based repos),
// so the Tx path is always taken. The NoTx paths have TOCTOU race conditions
// because they use plain SELECT without row-level locking (SELECT ... FOR UPDATE).
// Do NOT rely on NoTx paths for concurrent production workloads.

func (s *VisitService) StartVisit(ctx context.Context, visitID uuid.UUID, checkpoint *operations.VisitCheckpoint) error {
	if s.vPool != nil {
		return s.startVisitTx(ctx, visitID, checkpoint)
	}
	return s.startVisitNoTx(ctx, visitID, checkpoint)
}

func (s *VisitService) startVisitTx(ctx context.Context, visitID uuid.UUID, checkpoint *operations.VisitCheckpoint) error {
	tx, err := s.vPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.Warn("transaction rollback failed", "error", err)
		}
	}()

	visit, err := s.getVisitForUpdate(ctx, tx, visitID)
	if err != nil {
		return fmt.Errorf("start visit: %w", err)
	}

	if !canTransition(visit.Status, shared.VisitStatusEnRoute) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusEnRoute)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusEnRoute
	now := time.Now()
	visit.StartedAt = &now
	visit.UpdatedAt = now

	if err := s.visitUpdateInTx(ctx, tx, visitID, visit); err != nil {
		return fmt.Errorf("update visit in tx: %w", err)
	}

	if checkpoint != nil {
		checkpoint.VisitID = visitID
		if checkpoint.Timestamp.IsZero() {
			checkpoint.Timestamp = now
		}
		if err := s.checkpointCreateInTx(ctx, tx, checkpoint); err != nil {
			return fmt.Errorf("create checkpoint in tx: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	return nil
}

func (s *VisitService) startVisitNoTx(ctx context.Context, visitID uuid.UUID, checkpoint *operations.VisitCheckpoint) error {
	visit, err := s.visitRepo.GetByID(ctx, visitID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit", visitID.String())
	}

	if !canTransition(visit.Status, shared.VisitStatusEnRoute) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusEnRoute)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusEnRoute
	now := time.Now()
	visit.StartedAt = &now
	visit.UpdatedAt = now

	if err := s.visitRepo.Update(ctx, visitID, visit); err != nil {
		return fmt.Errorf("update visit: %w", err)
	}

	if checkpoint != nil {
		checkpoint.VisitID = visitID
		if checkpoint.Timestamp.IsZero() {
			checkpoint.Timestamp = now
		}
		if err := s.checkpointRepo.Create(ctx, checkpoint); err != nil {
			return fmt.Errorf("create checkpoint: %w", err)
		}
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	return nil
}

func (s *VisitService) BeginWork(ctx context.Context, visitID uuid.UUID) error {
	if s.vPool != nil {
		return s.beginWorkTx(ctx, visitID)
	}
	return s.beginWorkNoTx(ctx, visitID)
}

func (s *VisitService) beginWorkTx(ctx context.Context, visitID uuid.UUID) error {
	tx, err := s.vPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.Warn("transaction rollback failed", "error", err)
		}
	}()

	visit, err := s.getVisitForUpdate(ctx, tx, visitID)
	if err != nil {
		return fmt.Errorf("begin work: %w", err)
	}

	if !canTransition(visit.Status, shared.VisitStatusInProgress) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusInProgress)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusInProgress
	now := time.Now()
	visit.StartedAt = &now

	if err := s.visitUpdateInTx(ctx, tx, visitID, visit); err != nil {
		return fmt.Errorf("update visit in tx: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	return nil
}

func (s *VisitService) beginWorkNoTx(ctx context.Context, visitID uuid.UUID) error {
	visit, err := s.visitRepo.GetByID(ctx, visitID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit", visitID.String())
	}

	if !canTransition(visit.Status, shared.VisitStatusInProgress) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusInProgress)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusInProgress
	now := time.Now()
	visit.StartedAt = &now

	if err := s.visitRepo.Update(ctx, visitID, visit); err != nil {
		return fmt.Errorf("update visit: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	return nil
}

func (s *VisitService) CompleteVisit(ctx context.Context, visitID uuid.UUID) error {
	if s.vPool != nil {
		return s.completeVisitTx(ctx, visitID)
	}
	return s.completeVisitNoTx(ctx, visitID)
}

func (s *VisitService) completeVisitTx(ctx context.Context, visitID uuid.UUID) error {
	tx, err := s.vPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.Warn("transaction rollback failed", "error", err)
		}
	}()

	visit, err := s.getVisitForUpdate(ctx, tx, visitID)
	if err != nil {
		return fmt.Errorf("complete visit: %w", err)
	}

	if !canTransition(visit.Status, shared.VisitStatusCompleted) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusCompleted)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusCompleted
	now := time.Now()
	visit.CompletedAt = &now
	visit.UpdatedAt = now

	if err := s.visitUpdateInTx(ctx, tx, visitID, visit); err != nil {
		return fmt.Errorf("update visit in tx: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	return nil
}

func (s *VisitService) completeVisitNoTx(ctx context.Context, visitID uuid.UUID) error {
	visit, err := s.visitRepo.GetByID(ctx, visitID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit", visitID.String())
	}

	if !canTransition(visit.Status, shared.VisitStatusCompleted) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusCompleted)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusCompleted
	now := time.Now()
	visit.CompletedAt = &now
	visit.UpdatedAt = now

	if err := s.visitRepo.Update(ctx, visitID, visit); err != nil {
		return fmt.Errorf("update visit: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	return nil
}

func (s *VisitService) CancelVisit(ctx context.Context, visitID uuid.UUID, reason *string) error {
	if s.vPool != nil {
		return s.cancelVisitTx(ctx, visitID, reason)
	}
	return s.cancelVisitNoTx(ctx, visitID, reason)
}

func (s *VisitService) cancelVisitTx(ctx context.Context, visitID uuid.UUID, reason *string) error {
	tx, err := s.vPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.Warn("transaction rollback failed", "error", err)
		}
	}()

	visit, err := s.getVisitForUpdate(ctx, tx, visitID)
	if err != nil {
		return fmt.Errorf("cancel visit: %w", err)
	}

	if !canTransition(visit.Status, shared.VisitStatusCancelled) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusCancelled)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusCancelled
	visit.UpdatedAt = time.Now()

	if err := s.visitUpdateInTx(ctx, tx, visitID, visit); err != nil {
		return fmt.Errorf("update visit in tx: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, reason); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	return nil
}

func (s *VisitService) cancelVisitNoTx(ctx context.Context, visitID uuid.UUID, reason *string) error {
	visit, err := s.visitRepo.GetByID(ctx, visitID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit", visitID.String())
	}

	if !canTransition(visit.Status, shared.VisitStatusCancelled) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusCancelled)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusCancelled
	visit.UpdatedAt = time.Now()

	if err := s.visitRepo.Update(ctx, visitID, visit); err != nil {
		return fmt.Errorf("update visit: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, reason); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	return nil
}

func (s *VisitService) RescheduleVisit(ctx context.Context, visitID uuid.UUID, newScheduledAt time.Time) error {
	if s.vPool != nil {
		return s.rescheduleVisitTx(ctx, visitID, newScheduledAt)
	}
	return s.rescheduleVisitNoTx(ctx, visitID, newScheduledAt)
}

func (s *VisitService) rescheduleVisitTx(ctx context.Context, visitID uuid.UUID, newScheduledAt time.Time) error {
	tx, err := s.vPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			slog.Warn("transaction rollback failed", "error", err)
		}
	}()

	visit, err := s.getVisitForUpdate(ctx, tx, visitID)
	if err != nil {
		return fmt.Errorf("reschedule visit: %w", err)
	}

	if !canTransition(visit.Status, shared.VisitStatusRescheduled) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusRescheduled)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusRescheduled
	visit.UpdatedAt = time.Now()

	if err := s.visitUpdateInTx(ctx, tx, visitID, visit); err != nil {
		return fmt.Errorf("update visit in tx: %w", err)
	}

	newVisit := &operations.Visit{
		ID:              uuid.New(),
		PropertyID:      visit.PropertyID,
		PartnerID:       visit.PartnerID,
		PartnerOrderID:  visit.PartnerOrderID,
		ParentVisitID:   &visitID,
		RouteID:         visit.RouteID,
		DailyPlanID:     visit.DailyPlanID,
		Type:            visit.Type,
		Status:          shared.VisitStatusScheduled,
		Priority:        visit.Priority,
		Source:          visit.Source,
		SourceReference: visit.SourceReference,
		BillingTo:       visit.BillingTo,
		ScheduledAt:     newScheduledAt,
		Notes:           visit.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.visitCreateInTx(ctx, tx, newVisit); err != nil {
		return fmt.Errorf("create rescheduled visit in tx: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, newVisit.ID, shared.AuditActionAdd, nil, newVisit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", newVisit.ID, "error", err)
	}

	return nil
}

func (s *VisitService) rescheduleVisitNoTx(ctx context.Context, visitID uuid.UUID, newScheduledAt time.Time) error {
	visit, err := s.visitRepo.GetByID(ctx, visitID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit", visitID.String())
	}

	if !canTransition(visit.Status, shared.VisitStatusRescheduled) {
		return &service.TransitionError{Entity: "visit", From: string(visit.Status), To: string(shared.VisitStatusRescheduled)}
	}

	oldVisit := *visit
	visit.Status = shared.VisitStatusRescheduled
	visit.UpdatedAt = time.Now()

	if err := s.visitRepo.Update(ctx, visitID, visit); err != nil {
		return fmt.Errorf("update visit: %w", err)
	}

	newVisit := &operations.Visit{
		ID:              uuid.New(),
		PropertyID:      visit.PropertyID,
		PartnerID:       visit.PartnerID,
		PartnerOrderID:  visit.PartnerOrderID,
		ParentVisitID:   &visitID,
		RouteID:         visit.RouteID,
		DailyPlanID:     visit.DailyPlanID,
		Type:            visit.Type,
		Status:          shared.VisitStatusScheduled,
		Priority:        visit.Priority,
		Source:          visit.Source,
		SourceReference: visit.SourceReference,
		BillingTo:       visit.BillingTo,
		ScheduledAt:     newScheduledAt,
		Notes:           visit.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.visitRepo.Create(ctx, newVisit); err != nil {
		return fmt.Errorf("create rescheduled visit: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, visitID, shared.AuditActionUpdate, oldVisit, visit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", visitID, "error", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, newVisit.ID, shared.AuditActionAdd, nil, newVisit, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", newVisit.ID, "error", err)
	}

	return nil
}

func (s *VisitService) PopulateChecklist(ctx context.Context, visitID uuid.UUID, templateID uuid.UUID) error {
	if _, err := s.visitRepo.GetByID(ctx, visitID); err != nil {
		return service.HandleRepoGetByIDError(err, "visit", visitID.String())
	}

	tmplMats, err := s.checklistTmplMat.ListByTemplate(ctx, templateID)
	if err != nil {
		return fmt.Errorf("list template materials: %w", err)
	}
	if len(tmplMats) > 0 {
		items := make([]operations.VisitChecklistMaterial, 0, len(tmplMats))
		now := time.Now()
		for _, tm := range tmplMats {
			items = append(items, operations.VisitChecklistMaterial{
				ID:              uuid.New(),
				VisitID:         visitID,
				MaterialID:      tm.MaterialID,
				PlannedQuantity: tm.DefaultQuantity,
				Confirmed:       false,
				CreatedAt:       now,
			})
		}
		if err := s.visitCheckMat.CreateBatch(ctx, items); err != nil {
			return fmt.Errorf("batch create visit checklist materials: %w", err)
		}
	}

	tmplTools, err := s.checklistTmplTool.ListByTemplate(ctx, templateID)
	if err != nil {
		return fmt.Errorf("list template tools: %w", err)
	}
	if len(tmplTools) > 0 {
		items := make([]operations.VisitChecklistTool, 0, len(tmplTools))
		now := time.Now()
		for _, tt := range tmplTools {
			items = append(items, operations.VisitChecklistTool{
				ID:              uuid.New(),
				VisitID:         visitID,
				ToolID:          tt.ToolID,
				PlannedQuantity: tt.DefaultQuantity,
				Confirmed:       false,
				CreatedAt:       now,
			})
		}
		if err := s.visitCheckTool.CreateBatch(ctx, items); err != nil {
			return fmt.Errorf("batch create visit checklist tools: %w", err)
		}
	}

	tmplEPPs, err := s.checklistTmplEPP.ListByTemplate(ctx, templateID)
	if err != nil {
		return fmt.Errorf("list template EPPs: %w", err)
	}
	if len(tmplEPPs) > 0 {
		items := make([]operations.VisitChecklistEPP, 0, len(tmplEPPs))
		now := time.Now()
		for _, te := range tmplEPPs {
			items = append(items, operations.VisitChecklistEPP{
				ID:              uuid.New(),
				VisitID:         visitID,
				EPPID:           te.EPPID,
				PlannedQuantity: te.DefaultQuantity,
				Confirmed:       false,
				CreatedAt:       now,
			})
		}
		if err := s.visitCheckEPP.CreateBatch(ctx, items); err != nil {
			return fmt.Errorf("batch create visit checklist epps: %w", err)
		}
	}

	return nil
}

func (s *VisitService) GetByID(ctx context.Context, id uuid.UUID) (*operations.Visit, error) {
	item, err := s.visitRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get visit: %w", err)
	}
	return item, nil
}

func (s *VisitService) List(ctx context.Context, limit, offset int) (*repository.ListResult[operations.Visit], error) {
	items, err := s.visitRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list visits: %w", err)
	}
	return items, nil
}

func (s *VisitService) Delete(ctx context.Context, id uuid.UUID) error {
	visit, err := s.visitRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "visit", id.String())
	}

	if err := s.visitRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete visit: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeVisit, id, shared.AuditActionDelete, visit, nil, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "visit", "id", id, "error", err)
	}

	return nil
}

func (s *VisitService) getVisitForUpdate(ctx context.Context, tx pgx.Tx, visitID uuid.UUID) (*operations.Visit, error) {
	var v operations.Visit
	err := tx.QueryRow(ctx,
		`SELECT id, property_id, partner_id, partner_order_id, parent_visit_id, route_id, daily_plan_id,
		        type, status, priority, source, source_reference, billing_to, partner_supervisor,
		        result, rejection_reason_id, rejection_reason_detail,
		        scheduled_at, started_at, completed_at, alternative_location, notes,
		        created_at, updated_at
		 FROM operations.visits
		 WHERE id = $1 FOR UPDATE`, visitID,
	).Scan(&v.ID, &v.PropertyID, &v.PartnerID, &v.PartnerOrderID, &v.ParentVisitID, &v.RouteID, &v.DailyPlanID,
		&v.Type, &v.Status, &v.Priority, &v.Source, &v.SourceReference, &v.BillingTo, &v.PartnerSupervisor,
		&v.Result, &v.RejectionReasonID, &v.RejectionReasonDetail,
		&v.ScheduledAt, &v.StartedAt, &v.CompletedAt, &v.AlternativeLocation, &v.Notes,
		&v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &service.NotFoundError{Resource: "visit", ID: visitID.String()}
		}
		return nil, fmt.Errorf("get visit for update: %w", err)
	}
	return &v, nil
}

func (s *VisitService) visitUpdateInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, visit *operations.Visit) error {
	type updateInTxAble interface {
		UpdateInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, visit *operations.Visit) error
	}
	if repo, ok := s.visitRepo.(updateInTxAble); ok {
		return repo.UpdateInTx(ctx, tx, id, visit)
	}
	_, err := tx.Exec(ctx,
		`UPDATE operations.visits
		 SET property_id = $2, partner_id = $3, partner_order_id = $4, parent_visit_id = $5, route_id = $6, daily_plan_id = $7,
		     type = $8, status = $9, priority = $10, source = $11, source_reference = $12, billing_to = $13,
		     partner_supervisor = $14, result = $15, rejection_reason_id = $16, rejection_reason_detail = $17,
		     scheduled_at = $18, started_at = $19, completed_at = $20, alternative_location = $21, notes = $22,
		     updated_at = $23
		 WHERE id = $1`,
		id, visit.PropertyID, visit.PartnerID, visit.PartnerOrderID, visit.ParentVisitID, visit.RouteID, visit.DailyPlanID,
		visit.Type, visit.Status, visit.Priority, visit.Source, visit.SourceReference, visit.BillingTo, visit.PartnerSupervisor,
		visit.Result, visit.RejectionReasonID, visit.RejectionReasonDetail,
		visit.ScheduledAt, visit.StartedAt, visit.CompletedAt, visit.AlternativeLocation, visit.Notes,
		visit.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update visit in tx: %w", err)
	}
	return nil
}

func (s *VisitService) visitCreateInTx(ctx context.Context, tx pgx.Tx, visit *operations.Visit) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO operations.visits (id, property_id, partner_id, partner_order_id, parent_visit_id, route_id, daily_plan_id,
		        type, status, priority, source, source_reference, billing_to, partner_supervisor,
		        result, rejection_reason_id, rejection_reason_detail,
		        scheduled_at, started_at, completed_at, alternative_location, notes,
		        created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)`,
		visit.ID, visit.PropertyID, visit.PartnerID, visit.PartnerOrderID, visit.ParentVisitID, visit.RouteID, visit.DailyPlanID,
		visit.Type, visit.Status, visit.Priority, visit.Source, visit.SourceReference, visit.BillingTo, visit.PartnerSupervisor,
		visit.Result, visit.RejectionReasonID, visit.RejectionReasonDetail,
		visit.ScheduledAt, visit.StartedAt, visit.CompletedAt, visit.AlternativeLocation, visit.Notes,
		visit.CreatedAt, visit.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit in tx: %w", err)
	}
	return nil
}

func (s *VisitService) checkpointCreateInTx(ctx context.Context, tx pgx.Tx, cp *operations.VisitCheckpoint) error {
	type createInTxAble interface {
		CreateInTx(ctx context.Context, tx pgx.Tx, cp *operations.VisitCheckpoint) error
	}
	if repo, ok := s.checkpointRepo.(createInTxAble); ok {
		return repo.CreateInTx(ctx, tx, cp)
	}
	return s.checkpointRepo.Create(ctx, cp)
}
