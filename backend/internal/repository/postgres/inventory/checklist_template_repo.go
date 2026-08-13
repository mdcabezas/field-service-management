package postgresinventory

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type ChecklistTemplateRepo struct {
	pool *pgxpool.Pool
}

func NewChecklistTemplateRepo(pool *pgxpool.Pool) *ChecklistTemplateRepo {
	return &ChecklistTemplateRepo{pool: pool}
}

func (r *ChecklistTemplateRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.ChecklistTemplate, error) {
	var c inventory.ChecklistTemplate
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, description, visit_type, work_type, created_at
		 FROM inventory.checklist_templates
		 WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Description, &c.VisitType, &c.WorkType, &c.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "checklist_template", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get checklist_template by id: %w", err)
	}
	return &c, nil
}

func (r *ChecklistTemplateRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[inventory.ChecklistTemplate], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.checklist_templates`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count checklist templates: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, description, visit_type, work_type, created_at
		 FROM inventory.checklist_templates
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list checklist templates: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.ChecklistTemplate, 0)
	for rows.Next() {
		var c inventory.ChecklistTemplate
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.VisitType, &c.WorkType, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan checklist template: %w", err)
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[inventory.ChecklistTemplate]{Items: items, Total: total}, nil
}

func (r *ChecklistTemplateRepo) Create(ctx context.Context, tmpl *inventory.ChecklistTemplate) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.checklist_templates (id, name, description, visit_type, work_type, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		tmpl.ID, tmpl.Name, tmpl.Description, tmpl.VisitType, tmpl.WorkType, tmpl.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create checklist template: %w", err)
	}
	return nil
}

func (r *ChecklistTemplateRepo) Update(ctx context.Context, id uuid.UUID, tmpl *inventory.ChecklistTemplate) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inventory.checklist_templates
		 SET name = $2, description = $3, visit_type = $4, work_type = $5
		 WHERE id = $1`,
		id, tmpl.Name, tmpl.Description, tmpl.VisitType, tmpl.WorkType,
	)
	if err != nil {
		return fmt.Errorf("update checklist template: %w", err)
	}
	return nil
}

func (r *ChecklistTemplateRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inventory.checklist_templates WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete checklist template: %w", err)
	}
	return nil
}
