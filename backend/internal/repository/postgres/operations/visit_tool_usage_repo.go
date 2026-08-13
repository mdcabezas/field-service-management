package postgresoperations

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/service"
)

type VisitToolUsageRepo struct {
	pool *pgxpool.Pool
}

func NewVisitToolUsageRepo(pool *pgxpool.Pool) *VisitToolUsageRepo {
	return &VisitToolUsageRepo{pool: pool}
}

func (r *VisitToolUsageRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitToolUsage, error) {
	var t operations.VisitToolUsage
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, tool_id, notes, created_at
		 FROM operations.visit_tool_usages
		 WHERE id = $1`, id,
	).Scan(&t.ID, &t.VisitID, &t.ToolID, &t.Notes, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_tool_usage", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_tool_usage by id: %w", err)
	}
	return &t, nil
}

func (r *VisitToolUsageRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitToolUsage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, tool_id, notes, created_at
		 FROM operations.visit_tool_usages
		 WHERE visit_id = $1
		 ORDER BY created_at`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit tool usages by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitToolUsage, 0)
	for rows.Next() {
		var t operations.VisitToolUsage
		if err := rows.Scan(&t.ID, &t.VisitID, &t.ToolID, &t.Notes, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit tool usage: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitToolUsageRepo) Create(ctx context.Context, t *operations.VisitToolUsage) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_tool_usages (id, visit_id, tool_id, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.VisitID, t.ToolID, t.Notes, t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit tool usage: %w", err)
	}
	return nil
}

func (r *VisitToolUsageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_tool_usages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit tool usage: %w", err)
	}
	return nil
}
