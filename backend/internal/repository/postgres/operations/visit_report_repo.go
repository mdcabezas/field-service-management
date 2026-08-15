package postgresoperations

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type VisitReportRepo struct {
	pool *pgxpool.Pool
}

func NewVisitReportRepo(pool *pgxpool.Pool) *VisitReportRepo {
	return &VisitReportRepo{pool: pool}
}

func (r *VisitReportRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitReport, error) {
	var rep operations.VisitReport
	err := r.pool.QueryRow(ctx,
		`SELECT r.id, r.visit_id, r.report_template_id, r.source, r.recorded_at, r.created_at,
		        rt.name AS report_template_name,
		        CONCAT(COALESCE(vt.name, 'Visita'), ' - ', TO_CHAR(v.scheduled_at, 'DD/MM/YYYY')) AS visit_label
		 FROM operations.visit_reports r
		 LEFT JOIN shared.report_templates rt ON rt.id = r.report_template_id
		 LEFT JOIN operations.visits v ON v.id = r.visit_id
		 LEFT JOIN operations.visit_types vt ON vt.id = v.type
		 WHERE r.id = $1`, id,
	).Scan(&rep.ID, &rep.VisitID, &rep.ReportTemplateID, &rep.Source, &rep.RecordedAt, &rep.CreatedAt,
		&rep.ReportTemplateName, &rep.VisitLabel)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_report", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_report by id: %w", err)
	}
	return &rep, nil
}

func (r *VisitReportRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitReport, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.id, r.visit_id, r.report_template_id, r.source, r.recorded_at, r.created_at,
		        rt.name AS report_template_name,
		        CONCAT(COALESCE(vt.name, 'Visita'), ' - ', TO_CHAR(v.scheduled_at, 'DD/MM/YYYY')) AS visit_label
		 FROM operations.visit_reports r
		 LEFT JOIN shared.report_templates rt ON rt.id = r.report_template_id
		 LEFT JOIN operations.visits v ON v.id = r.visit_id
		 LEFT JOIN operations.visit_types vt ON vt.id = v.type
		 WHERE r.visit_id = $1
		 ORDER BY r.recorded_at`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit reports by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitReport, 0)
	for rows.Next() {
		var rep operations.VisitReport
		if err := rows.Scan(&rep.ID, &rep.VisitID, &rep.ReportTemplateID, &rep.Source, &rep.RecordedAt, &rep.CreatedAt,
			&rep.ReportTemplateName, &rep.VisitLabel); err != nil {
			return nil, fmt.Errorf("scan visit report: %w", err)
		}
		items = append(items, rep)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitReportRepo) Create(ctx context.Context, rep *operations.VisitReport) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_reports (id, visit_id, report_template_id, source, recorded_at, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		rep.ID, rep.VisitID, rep.ReportTemplateID, rep.Source, rep.RecordedAt, rep.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit report: %w", err)
	}
	return nil
}

func (r *VisitReportRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_reports WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit report: %w", err)
	}
	return nil
}

func (r *VisitReportRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[operations.VisitReport], error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.id, r.visit_id, r.report_template_id, r.source, r.recorded_at, r.created_at,
		        rt.name AS report_template_name,
		        CONCAT(COALESCE(vt.name, 'Visita'), ' - ', TO_CHAR(v.scheduled_at, 'DD/MM/YYYY')) AS visit_label
		 FROM operations.visit_reports r
		 LEFT JOIN shared.report_templates rt ON rt.id = r.report_template_id
		 LEFT JOIN operations.visits v ON v.id = r.visit_id
		 LEFT JOIN operations.visit_types vt ON vt.id = v.type
		 ORDER BY r.recorded_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit reports: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitReport, 0)
	for rows.Next() {
		var rep operations.VisitReport
		if err := rows.Scan(&rep.ID, &rep.VisitID, &rep.ReportTemplateID, &rep.Source, &rep.RecordedAt, &rep.CreatedAt,
			&rep.ReportTemplateName, &rep.VisitLabel); err != nil {
			return nil, fmt.Errorf("scan visit report: %w", err)
		}
		items = append(items, rep)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	var total int
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM operations.visit_reports`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count visit reports: %w", err)
	}

	return &repository.ListResult[operations.VisitReport]{Items: items, Total: total}, nil
}
