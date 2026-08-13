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

type ReportEntryRepo struct {
	pool *pgxpool.Pool
}

func NewReportEntryRepo(pool *pgxpool.Pool) *ReportEntryRepo {
	return &ReportEntryRepo{pool: pool}
}

func (r *ReportEntryRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.ReportEntry, error) {
	var e operations.ReportEntry
	err := r.pool.QueryRow(ctx,
		`SELECT id, report_id, data_json, created_at
		 FROM operations.report_entries
		 WHERE id = $1`, id,
	).Scan(&e.ID, &e.ReportID, &e.DataJSON, &e.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "report_entry", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get report_entry by id: %w", err)
	}
	return &e, nil
}

func (r *ReportEntryRepo) ListByReport(ctx context.Context, reportID uuid.UUID) ([]operations.ReportEntry, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, report_id, data_json, created_at
		 FROM operations.report_entries
		 WHERE report_id = $1
		 ORDER BY created_at`, reportID,
	)
	if err != nil {
		return nil, fmt.Errorf("list report entries by report: %w", err)
	}
	defer rows.Close()

	items := make([]operations.ReportEntry, 0)
	for rows.Next() {
		var e operations.ReportEntry
		if err := rows.Scan(&e.ID, &e.ReportID, &e.DataJSON, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan report entry: %w", err)
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *ReportEntryRepo) Create(ctx context.Context, e *operations.ReportEntry) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.report_entries (id, report_id, data_json, created_at)
		 VALUES ($1, $2, $3, $4)`,
		e.ID, e.ReportID, e.DataJSON, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create report entry: %w", err)
	}
	return nil
}

func (r *ReportEntryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.report_entries WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete report entry: %w", err)
	}
	return nil
}
