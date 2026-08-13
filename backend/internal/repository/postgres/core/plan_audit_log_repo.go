package postgrescore

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/core"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type PlanAuditLogRepo struct {
	pool *pgxpool.Pool
}

func NewPlanAuditLogRepo(pool *pgxpool.Pool) *PlanAuditLogRepo {
	return &PlanAuditLogRepo{pool: pool}
}

func (r *PlanAuditLogRepo) GetByID(ctx context.Context, id uuid.UUID) (*core.PlanAuditLog, error) {
	var e core.PlanAuditLog
	err := r.pool.QueryRow(ctx,
		`SELECT id, entity_type, entity_id, action, old_data, new_data, user_id, reason, timestamp
		 FROM core.plan_audit_log
		 WHERE id = $1`, id,
	).Scan(&e.ID, &e.EntityType, &e.EntityID, &e.Action, &e.OldData, &e.NewData, &e.UserID, &e.Reason, &e.Timestamp)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "audit_log", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get audit_log by id: %w", err)
	}
	return &e, nil
}

func (r *PlanAuditLogRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[core.PlanAuditLog], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM core.plan_audit_log`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count audit logs: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, entity_type, entity_id, action, old_data, new_data, user_id, reason, timestamp
		 FROM core.plan_audit_log
		 ORDER BY timestamp DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	defer rows.Close()

	items := make([]core.PlanAuditLog, 0)
	for rows.Next() {
		var e core.PlanAuditLog
		if err := rows.Scan(&e.ID, &e.EntityType, &e.EntityID, &e.Action, &e.OldData, &e.NewData, &e.UserID, &e.Reason, &e.Timestamp); err != nil {
			return nil, fmt.Errorf("scan audit log: %w", err)
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[core.PlanAuditLog]{Items: items, Total: total}, nil
}

func (r *PlanAuditLogRepo) ListByEntity(ctx context.Context, entityType shared.AuditEntityType, entityID uuid.UUID) ([]core.PlanAuditLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, entity_type, entity_id, action, old_data, new_data, user_id, reason, timestamp
		 FROM core.plan_audit_log
		 WHERE entity_type = $1 AND entity_id = $2
		 ORDER BY timestamp DESC`, entityType, entityID,
	)
	if err != nil {
		return nil, fmt.Errorf("list audit logs by entity: %w", err)
	}
	defer rows.Close()

	items := make([]core.PlanAuditLog, 0)
	for rows.Next() {
		var e core.PlanAuditLog
		if err := rows.Scan(&e.ID, &e.EntityType, &e.EntityID, &e.Action, &e.OldData, &e.NewData, &e.UserID, &e.Reason, &e.Timestamp); err != nil {
			return nil, fmt.Errorf("scan audit log: %w", err)
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *PlanAuditLogRepo) Create(ctx context.Context, entry *core.PlanAuditLog) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO core.plan_audit_log (id, entity_type, entity_id, action, old_data, new_data, user_id, reason, timestamp)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		entry.ID, entry.EntityType, entry.EntityID, entry.Action, entry.OldData, entry.NewData, entry.UserID, entry.Reason, entry.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}
