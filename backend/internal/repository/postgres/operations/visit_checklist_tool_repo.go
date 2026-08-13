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

type VisitChecklistToolRepo struct {
	pool *pgxpool.Pool
}

func NewVisitChecklistToolRepo(pool *pgxpool.Pool) *VisitChecklistToolRepo {
	return &VisitChecklistToolRepo{pool: pool}
}

func (r *VisitChecklistToolRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitChecklistTool, error) {
	var t operations.VisitChecklistTool
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, tool_id, planned_quantity, confirmed, notes, created_at
		 FROM operations.visit_checklist_tools
		 WHERE id = $1`, id,
	).Scan(&t.ID, &t.VisitID, &t.ToolID, &t.PlannedQuantity, &t.Confirmed, &t.Notes, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_checklist_tool", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_checklist_tool by id: %w", err)
	}
	return &t, nil
}

func (r *VisitChecklistToolRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitChecklistTool, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, tool_id, planned_quantity, confirmed, notes, created_at
		 FROM operations.visit_checklist_tools
		 WHERE visit_id = $1
		 ORDER BY created_at
		 LIMIT 100`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit checklist tools by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitChecklistTool, 0)
	for rows.Next() {
		var t operations.VisitChecklistTool
		if err := rows.Scan(&t.ID, &t.VisitID, &t.ToolID, &t.PlannedQuantity, &t.Confirmed, &t.Notes, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit checklist tool: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitChecklistToolRepo) Create(ctx context.Context, t *operations.VisitChecklistTool) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_checklist_tools (id, visit_id, tool_id, planned_quantity, confirmed, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		t.ID, t.VisitID, t.ToolID, t.PlannedQuantity, t.Confirmed, t.Notes, t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit checklist tool: %w", err)
	}
	return nil
}

func (r *VisitChecklistToolRepo) CreateBatch(ctx context.Context, items []operations.VisitChecklistTool) error {
	const maxBatchSize = 1000

	if len(items) == 0 {
		return nil
	}
	if len(items) > maxBatchSize {
		return fmt.Errorf("batch size %d exceeds maximum %d", len(items), maxBatchSize)
	}
	batch := &pgx.Batch{}
	for _, t := range items {
		batch.Queue(
			`INSERT INTO operations.visit_checklist_tools (id, visit_id, tool_id, planned_quantity, confirmed, notes, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			t.ID, t.VisitID, t.ToolID, t.PlannedQuantity, t.Confirmed, t.Notes, t.CreatedAt,
		)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range items {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("batch create visit checklist tools: %w", err)
		}
	}
	return nil
}

func (r *VisitChecklistToolRepo) Update(ctx context.Context, id uuid.UUID, t *operations.VisitChecklistTool) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE operations.visit_checklist_tools
		 SET visit_id = $2, tool_id = $3, planned_quantity = $4, confirmed = $5, notes = $6
		 WHERE id = $1`,
		id, t.VisitID, t.ToolID, t.PlannedQuantity, t.Confirmed, t.Notes,
	)
	if err != nil {
		return fmt.Errorf("update visit checklist tool: %w", err)
	}
	return nil
}

func (r *VisitChecklistToolRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_checklist_tools WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit checklist tool: %w", err)
	}
	return nil
}
