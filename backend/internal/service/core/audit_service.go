package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	coremodel "localis-backend/internal/model/core"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type AuditService struct {
	auditRepo repository.CorePlanAuditLogRepository
}

func NewAuditService(auditRepo repository.CorePlanAuditLogRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

func (s *AuditService) RecordChange(ctx context.Context, entityType shared.AuditEntityType, entityID uuid.UUID, action shared.AuditAction, oldData, newData any, reason *string) error {
	var oldRaw, newRaw json.RawMessage
	if oldData != nil {
		b, err := json.Marshal(oldData)
		if err != nil {
			return fmt.Errorf("marshal old audit data: %w", err)
		}
		oldRaw = b
	}
	if newData != nil {
		b, err := json.Marshal(newData)
		if err != nil {
			return fmt.Errorf("marshal new audit data: %w", err)
		}
		newRaw = b
	}

	userID := service.UserIDFromCtx(ctx)
	var userIDPtr *uuid.UUID
	if userID != "" {
		if uid, err := uuid.Parse(userID); err == nil {
			userIDPtr = &uid
		}
	}

	entry := &coremodel.PlanAuditLog{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		OldData:    oldRaw,
		NewData:    newRaw,
		UserID:     userIDPtr,
		Reason:     reason,
		Timestamp:  time.Now(),
	}

	if err := s.auditRepo.Create(ctx, entry); err != nil {
		return fmt.Errorf("create audit log entry: %w", err)
	}
	return nil
}

func (s *AuditService) ListByEntity(ctx context.Context, entityType shared.AuditEntityType, entityID uuid.UUID) ([]coremodel.PlanAuditLog, error) {
	items, err := s.auditRepo.ListByEntity(ctx, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("list audit logs by entity: %w", err)
	}
	return items, nil
}
