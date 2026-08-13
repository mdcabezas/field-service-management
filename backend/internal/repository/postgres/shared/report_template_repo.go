package postgresshared

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type ReportTemplateRepo struct {
	pool *pgxpool.Pool
}

func NewReportTemplateRepo(pool *pgxpool.Pool) *ReportTemplateRepo {
	return &ReportTemplateRepo{pool: pool}
}

func (r *ReportTemplateRepo) GetByID(ctx context.Context, id uuid.UUID) (*shared.ReportTemplate, error) {
	var t shared.ReportTemplate
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, fields_json, created_at
		 FROM shared.report_templates
		 WHERE id = $1`, id,
	).Scan(&t.ID, &t.Name, &t.FieldsJSON, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "report_template", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get report_template by id: %w", err)
	}
	return &t, nil
}

func (r *ReportTemplateRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[shared.ReportTemplate], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM shared.report_templates`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count report templates: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, fields_json, created_at
		 FROM shared.report_templates
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list report templates: %w", err)
	}
	defer rows.Close()

	items := make([]shared.ReportTemplate, 0)
	for rows.Next() {
		var t shared.ReportTemplate
		if err := rows.Scan(&t.ID, &t.Name, &t.FieldsJSON, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan report template: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[shared.ReportTemplate]{Items: items, Total: total}, nil
}

func (r *ReportTemplateRepo) Create(ctx context.Context, tmpl *shared.ReportTemplate) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO shared.report_templates (id, name, fields_json, created_at)
		 VALUES ($1, $2, $3, $4)`,
		tmpl.ID, tmpl.Name, tmpl.FieldsJSON, tmpl.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create report template: %w", err)
	}
	return nil
}

func (r *ReportTemplateRepo) Update(ctx context.Context, id uuid.UUID, tmpl *shared.ReportTemplate) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE shared.report_templates
		 SET name = $2, fields_json = $3
		 WHERE id = $1`,
		id, tmpl.Name, tmpl.FieldsJSON,
	)
	if err != nil {
		return fmt.Errorf("update report template: %w", err)
	}
	return nil
}

func (r *ReportTemplateRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM shared.report_templates WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete report template: %w", err)
	}
	return nil
}
