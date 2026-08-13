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

type ChecklistTemplateEPPRepo struct {
	pool *pgxpool.Pool
}

func NewChecklistTemplateEPPRepo(pool *pgxpool.Pool) *ChecklistTemplateEPPRepo {
	return &ChecklistTemplateEPPRepo{pool: pool}
}

func (r *ChecklistTemplateEPPRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.ChecklistTemplateEPP, error) {
	var e operations.ChecklistTemplateEPP
	err := r.pool.QueryRow(ctx,
		`SELECT id, template_id, epp_id, default_quantity, required, created_at
		 FROM operations.checklist_template_epps
		 WHERE id = $1`, id,
	).Scan(&e.ID, &e.TemplateID, &e.EPPID, &e.DefaultQuantity, &e.Required, &e.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "checklist_template_epp", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get checklist_template_epp by id: %w", err)
	}
	return &e, nil
}

func (r *ChecklistTemplateEPPRepo) ListByTemplate(ctx context.Context, templateID uuid.UUID) ([]operations.ChecklistTemplateEPP, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, template_id, epp_id, default_quantity, required, created_at
		 FROM operations.checklist_template_epps
		 WHERE template_id = $1
		 ORDER BY created_at`, templateID,
	)
	if err != nil {
		return nil, fmt.Errorf("list checklist template epps by template: %w", err)
	}
	defer rows.Close()

	items := make([]operations.ChecklistTemplateEPP, 0)
	for rows.Next() {
		var e operations.ChecklistTemplateEPP
		if err := rows.Scan(&e.ID, &e.TemplateID, &e.EPPID, &e.DefaultQuantity, &e.Required, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan checklist template epp: %w", err)
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *ChecklistTemplateEPPRepo) Create(ctx context.Context, e *operations.ChecklistTemplateEPP) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.checklist_template_epps (id, template_id, epp_id, default_quantity, required, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		e.ID, e.TemplateID, e.EPPID, e.DefaultQuantity, e.Required, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create checklist template epp: %w", err)
	}
	return nil
}

func (r *ChecklistTemplateEPPRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.checklist_template_epps WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete checklist template epp: %w", err)
	}
	return nil
}
