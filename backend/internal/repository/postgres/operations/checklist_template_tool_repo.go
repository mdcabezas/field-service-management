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

type ChecklistTemplateToolRepo struct {
	pool *pgxpool.Pool
}

func NewChecklistTemplateToolRepo(pool *pgxpool.Pool) *ChecklistTemplateToolRepo {
	return &ChecklistTemplateToolRepo{pool: pool}
}

func (r *ChecklistTemplateToolRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.ChecklistTemplateTool, error) {
	var t operations.ChecklistTemplateTool
	err := r.pool.QueryRow(ctx,
		`SELECT id, template_id, tool_id, default_quantity, required, created_at
		 FROM operations.checklist_template_tools
		 WHERE id = $1`, id,
	).Scan(&t.ID, &t.TemplateID, &t.ToolID, &t.DefaultQuantity, &t.Required, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "checklist_template_tool", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get checklist_template_tool by id: %w", err)
	}
	return &t, nil
}

func (r *ChecklistTemplateToolRepo) ListByTemplate(ctx context.Context, templateID uuid.UUID) ([]operations.ChecklistTemplateTool, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, template_id, tool_id, default_quantity, required, created_at
		 FROM operations.checklist_template_tools
		 WHERE template_id = $1
		 ORDER BY created_at`, templateID,
	)
	if err != nil {
		return nil, fmt.Errorf("list checklist template tools by template: %w", err)
	}
	defer rows.Close()

	items := make([]operations.ChecklistTemplateTool, 0)
	for rows.Next() {
		var t operations.ChecklistTemplateTool
		if err := rows.Scan(&t.ID, &t.TemplateID, &t.ToolID, &t.DefaultQuantity, &t.Required, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan checklist template tool: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *ChecklistTemplateToolRepo) Create(ctx context.Context, t *operations.ChecklistTemplateTool) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.checklist_template_tools (id, template_id, tool_id, default_quantity, required, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		t.ID, t.TemplateID, t.ToolID, t.DefaultQuantity, t.Required, t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create checklist template tool: %w", err)
	}
	return nil
}

func (r *ChecklistTemplateToolRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.checklist_template_tools WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete checklist template tool: %w", err)
	}
	return nil
}
