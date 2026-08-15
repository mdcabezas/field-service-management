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

type ReportImageRepo struct {
	pool *pgxpool.Pool
}

func NewReportImageRepo(pool *pgxpool.Pool) *ReportImageRepo {
	return &ReportImageRepo{pool: pool}
}

func (r *ReportImageRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.ReportImage, error) {
	var img operations.ReportImage
	err := r.pool.QueryRow(ctx,
		`SELECT id, report_id, url, pages, created_at
		 FROM operations.report_images
		 WHERE id = $1`, id,
	).Scan(&img.ID, &img.ReportID, &img.URL, &img.Pages, &img.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "report_image", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get report_image by id: %w", err)
	}
	return &img, nil
}

func (r *ReportImageRepo) ListByReport(ctx context.Context, reportID uuid.UUID) ([]operations.ReportImage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, report_id, url, pages, created_at
		 FROM operations.report_images
		 WHERE report_id = $1
		 ORDER BY created_at`, reportID,
	)
	if err != nil {
		return nil, fmt.Errorf("list report images by report: %w", err)
	}
	defer rows.Close()

	items := make([]operations.ReportImage, 0)
	for rows.Next() {
		var img operations.ReportImage
		if err := rows.Scan(&img.ID, &img.ReportID, &img.URL, &img.Pages, &img.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan report image: %w", err)
		}
		items = append(items, img)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *ReportImageRepo) Create(ctx context.Context, img *operations.ReportImage) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.report_images (id, report_id, url, pages, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		img.ID, img.ReportID, img.URL, img.Pages, img.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create report image: %w", err)
	}
	return nil
}

func (r *ReportImageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.report_images WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete report image: %w", err)
	}
	return nil
}

func (r *ReportImageRepo) DeleteByReport(ctx context.Context, reportID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.report_images WHERE report_id = $1`, reportID)
	if err != nil {
		return fmt.Errorf("delete report images by report: %w", err)
	}
	return nil
}
