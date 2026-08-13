package planning

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"localis-backend/internal/model/planning"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
	coreservice "localis-backend/internal/service/core"
)

type DailyLoadService struct {
	loadMatRepo  repository.DailyLoadMaterialRepository
	loadToolRepo repository.DailyLoadToolRepository
	loadEPPRepo  repository.DailyLoadEPPRepository
	matRepo      repository.MaterialRepository
	toolRepo     repository.ToolRepository
	eppRepo      repository.EPPItemRepository
	audit        *coreservice.AuditService
}

func NewDailyLoadService(
	loadMatRepo repository.DailyLoadMaterialRepository,
	loadToolRepo repository.DailyLoadToolRepository,
	loadEPPRepo repository.DailyLoadEPPRepository,
	matRepo repository.MaterialRepository,
	toolRepo repository.ToolRepository,
	eppRepo repository.EPPItemRepository,
	audit *coreservice.AuditService,
) *DailyLoadService {
	return &DailyLoadService{
		loadMatRepo:  loadMatRepo,
		loadToolRepo: loadToolRepo,
		loadEPPRepo:  loadEPPRepo,
		matRepo:      matRepo,
		toolRepo:     toolRepo,
		eppRepo:      eppRepo,
		audit:        audit,
	}
}

func (s *DailyLoadService) LoadMaterial(ctx context.Context, m *planning.DailyLoadMaterial) error {
	_, err := s.matRepo.GetByID(ctx, m.MaterialID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "material", m.MaterialID.String())
	}

	m.ID = uuid.New()
	m.ReturnedQuantity = 0
	m.CreatedAt = time.Now()

	if err := s.loadMatRepo.Create(ctx, m); err != nil {
		return fmt.Errorf("create daily load material: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyLoadMaterial, m.ID, shared.AuditActionAdd, nil, m, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_load_material", "id", m.ID, "error", err)
	}

	return nil
}

func (s *DailyLoadService) LoadTool(ctx context.Context, t *planning.DailyLoadTool) error {
	_, err := s.toolRepo.GetByID(ctx, t.ToolID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "tool", t.ToolID.String())
	}

	t.ID = uuid.New()
	t.ReturnedQuantity = 0
	t.CreatedAt = time.Now()

	if err := s.loadToolRepo.Create(ctx, t); err != nil {
		return fmt.Errorf("create daily load tool: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyLoadTool, t.ID, shared.AuditActionAdd, nil, t, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_load_tool", "id", t.ID, "error", err)
	}

	return nil
}

func (s *DailyLoadService) LoadEPP(ctx context.Context, e *planning.DailyLoadEPP) error {
	_, err := s.eppRepo.GetByID(ctx, e.EPPID)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "epp_item", e.EPPID.String())
	}

	e.ID = uuid.New()
	e.ReturnedQuantity = 0
	e.CreatedAt = time.Now()

	if err := s.loadEPPRepo.Create(ctx, e); err != nil {
		return fmt.Errorf("create daily load epp: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyLoadEPP, e.ID, shared.AuditActionAdd, nil, e, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_load_epp", "id", e.ID, "error", err)
	}

	return nil
}

func (s *DailyLoadService) ReturnMaterial(ctx context.Context, id uuid.UUID, returnedQuantity float64) error {
	if returnedQuantity < 0 {
		return &service.ValidationError{Field: "returned_quantity", Message: "must not be negative"}
	}

	old, err := s.loadMatRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "daily_load_material", id.String())
	}

	oldCopy := *old
	old.ReturnedQuantity = returnedQuantity
	if err := s.loadMatRepo.Update(ctx, id, old); err != nil {
		return fmt.Errorf("update daily load material: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyLoadMaterial, id, shared.AuditActionUpdate, &oldCopy, old, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_load_material", "id", id, "error", err)
	}

	return nil
}

func (s *DailyLoadService) ReturnTool(ctx context.Context, id uuid.UUID, returnedQuantity float64) error {
	if returnedQuantity < 0 {
		return &service.ValidationError{Field: "returned_quantity", Message: "must not be negative"}
	}

	old, err := s.loadToolRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "daily_load_tool", id.String())
	}

	oldCopy := *old
	old.ReturnedQuantity = returnedQuantity
	if err := s.loadToolRepo.Update(ctx, id, old); err != nil {
		return fmt.Errorf("update daily load tool: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyLoadTool, id, shared.AuditActionUpdate, &oldCopy, old, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_load_tool", "id", id, "error", err)
	}

	return nil
}

func (s *DailyLoadService) ReturnEPP(ctx context.Context, id uuid.UUID, returnedQuantity float64) error {
	if returnedQuantity < 0 {
		return &service.ValidationError{Field: "returned_quantity", Message: "must not be negative"}
	}

	old, err := s.loadEPPRepo.GetByID(ctx, id)
	if err != nil {
		return service.HandleRepoGetByIDError(err, "daily_load_epp", id.String())
	}

	oldCopy := *old
	old.ReturnedQuantity = returnedQuantity
	if err := s.loadEPPRepo.Update(ctx, id, old); err != nil {
		return fmt.Errorf("update daily load epp: %w", err)
	}

	if err := s.audit.RecordChange(ctx, shared.AuditEntityTypeDailyLoadEPP, id, shared.AuditActionUpdate, &oldCopy, old, nil); err != nil {
		slog.Warn("failed to record audit change", "entity", "daily_load_epp", "id", id, "error", err)
	}

	return nil
}

func (s *DailyLoadService) ListMaterialsByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyLoadMaterial, error) {
	result, err := s.loadMatRepo.ListByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("list daily load materials: %w", err)
	}
	return result, nil
}

func (s *DailyLoadService) ListToolsByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyLoadTool, error) {
	result, err := s.loadToolRepo.ListByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("list daily load tools: %w", err)
	}
	return result, nil
}

func (s *DailyLoadService) ListEPPsByPlan(ctx context.Context, planID uuid.UUID) ([]planning.DailyLoadEPP, error) {
	result, err := s.loadEPPRepo.ListByPlan(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("list daily load epps: %w", err)
	}
	return result, nil
}

func (s *DailyLoadService) GetMaterialByID(ctx context.Context, id uuid.UUID) (*planning.DailyLoadMaterial, error) {
	result, err := s.loadMatRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get daily load material by id: %w", err)
	}
	return result, nil
}

func (s *DailyLoadService) GetToolByID(ctx context.Context, id uuid.UUID) (*planning.DailyLoadTool, error) {
	result, err := s.loadToolRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get daily load tool by id: %w", err)
	}
	return result, nil
}

func (s *DailyLoadService) GetEPPByID(ctx context.Context, id uuid.UUID) (*planning.DailyLoadEPP, error) {
	result, err := s.loadEPPRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get daily load epp by id: %w", err)
	}
	return result, nil
}

func (s *DailyLoadService) UpdateMaterial(ctx context.Context, id uuid.UUID, m *planning.DailyLoadMaterial) error {
	if _, err := s.loadMatRepo.GetByID(ctx, id); err != nil {
		return service.HandleRepoGetByIDError(err, "daily_load_material", id.String())
	}
	if err := s.loadMatRepo.Update(ctx, id, m); err != nil {
		return fmt.Errorf("update daily load material: %w", err)
	}
	return nil
}

func (s *DailyLoadService) UpdateTool(ctx context.Context, id uuid.UUID, t *planning.DailyLoadTool) error {
	if _, err := s.loadToolRepo.GetByID(ctx, id); err != nil {
		return service.HandleRepoGetByIDError(err, "daily_load_tool", id.String())
	}
	if err := s.loadToolRepo.Update(ctx, id, t); err != nil {
		return fmt.Errorf("update daily load tool: %w", err)
	}
	return nil
}

func (s *DailyLoadService) UpdateEPP(ctx context.Context, id uuid.UUID, e *planning.DailyLoadEPP) error {
	if _, err := s.loadEPPRepo.GetByID(ctx, id); err != nil {
		return service.HandleRepoGetByIDError(err, "daily_load_epp", id.String())
	}
	if err := s.loadEPPRepo.Update(ctx, id, e); err != nil {
		return fmt.Errorf("update daily load epp: %w", err)
	}
	return nil
}

func (s *DailyLoadService) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	if _, err := s.loadMatRepo.GetByID(ctx, id); err != nil {
		return service.HandleRepoGetByIDError(err, "daily_load_material", id.String())
	}
	if err := s.loadMatRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete daily load material: %w", err)
	}
	return nil
}

func (s *DailyLoadService) DeleteTool(ctx context.Context, id uuid.UUID) error {
	if _, err := s.loadToolRepo.GetByID(ctx, id); err != nil {
		return service.HandleRepoGetByIDError(err, "daily_load_tool", id.String())
	}
	if err := s.loadToolRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete daily load tool: %w", err)
	}
	return nil
}

func (s *DailyLoadService) DeleteEPP(ctx context.Context, id uuid.UUID) error {
	if _, err := s.loadEPPRepo.GetByID(ctx, id); err != nil {
		return service.HandleRepoGetByIDError(err, "daily_load_epp", id.String())
	}
	if err := s.loadEPPRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete daily load epp: %w", err)
	}
	return nil
}
