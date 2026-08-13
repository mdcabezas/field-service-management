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

type ChecklistTemplateMaterialRepo struct {
	pool *pgxpool.Pool
}

func NewChecklistTemplateMaterialRepo(pool *pgxpool.Pool) *ChecklistTemplateMaterialRepo {
	return &ChecklistTemplateMaterialRepo{pool: pool}
}

func (r *ChecklistTemplateMaterialRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.ChecklistTemplateMaterial, error) {
	var m operations.ChecklistTemplateMaterial
	err := r.pool.QueryRow(ctx,
		`SELECT id, template_id, material_id, default_quantity, required, created_at
		 FROM operations.checklist_template_materials
		 WHERE id = $1`, id,
	).Scan(&m.ID, &m.TemplateID, &m.MaterialID, &m.DefaultQuantity, &m.Required, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "checklist_template_material", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get checklist_template_material by id: %w", err)
	}
	return &m, nil
}

func (r *ChecklistTemplateMaterialRepo) ListByTemplate(ctx context.Context, templateID uuid.UUID) ([]operations.ChecklistTemplateMaterial, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, template_id, material_id, default_quantity, required, created_at
		 FROM operations.checklist_template_materials
		 WHERE template_id = $1
		 ORDER BY created_at`, templateID,
	)
	if err != nil {
		return nil, fmt.Errorf("list checklist template materials by template: %w", err)
	}
	defer rows.Close()

	items := make([]operations.ChecklistTemplateMaterial, 0)
	for rows.Next() {
		var m operations.ChecklistTemplateMaterial
		if err := rows.Scan(&m.ID, &m.TemplateID, &m.MaterialID, &m.DefaultQuantity, &m.Required, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan checklist template material: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *ChecklistTemplateMaterialRepo) Create(ctx context.Context, m *operations.ChecklistTemplateMaterial) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.checklist_template_materials (id, template_id, material_id, default_quantity, required, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		m.ID, m.TemplateID, m.MaterialID, m.DefaultQuantity, m.Required, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create checklist template material: %w", err)
	}
	return nil
}

func (r *ChecklistTemplateMaterialRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.checklist_template_materials WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete checklist template material: %w", err)
	}
	return nil
}
