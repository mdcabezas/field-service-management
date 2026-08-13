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

type VisitMeasurementRepo struct {
	pool *pgxpool.Pool
}

func NewVisitMeasurementRepo(pool *pgxpool.Pool) *VisitMeasurementRepo {
	return &VisitMeasurementRepo{pool: pool}
}

func (r *VisitMeasurementRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitMeasurement, error) {
	var m operations.VisitMeasurement
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, type, value, unit, result, measuring_device, notes, photo_url, created_at
		 FROM operations.visit_measurements
		 WHERE id = $1`, id,
	).Scan(&m.ID, &m.VisitID, &m.Type, &m.Value, &m.Unit, &m.Result, &m.MeasuringDevice, &m.Notes, &m.PhotoURL, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_measurement", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_measurement by id: %w", err)
	}
	return &m, nil
}

func (r *VisitMeasurementRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitMeasurement, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, type, value, unit, result, measuring_device, notes, photo_url, created_at
		 FROM operations.visit_measurements
		 WHERE visit_id = $1
		 ORDER BY created_at`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit measurements by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitMeasurement, 0)
	for rows.Next() {
		var m operations.VisitMeasurement
		if err := rows.Scan(&m.ID, &m.VisitID, &m.Type, &m.Value, &m.Unit, &m.Result, &m.MeasuringDevice, &m.Notes, &m.PhotoURL, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit measurement: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitMeasurementRepo) Create(ctx context.Context, m *operations.VisitMeasurement) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_measurements (id, visit_id, type, value, unit, result, measuring_device, notes, photo_url, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		m.ID, m.VisitID, m.Type, m.Value, m.Unit, m.Result, m.MeasuringDevice, m.Notes, m.PhotoURL, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit measurement: %w", err)
	}
	return nil
}

func (r *VisitMeasurementRepo) Update(ctx context.Context, id uuid.UUID, m *operations.VisitMeasurement) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE operations.visit_measurements
		 SET visit_id = $2, type = $3, value = $4, unit = $5, result = $6, measuring_device = $7, notes = $8, photo_url = $9
		 WHERE id = $1`,
		id, m.VisitID, m.Type, m.Value, m.Unit, m.Result, m.MeasuringDevice, m.Notes, m.PhotoURL,
	)
	if err != nil {
		return fmt.Errorf("update visit measurement: %w", err)
	}
	return nil
}

func (r *VisitMeasurementRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_measurements WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit measurement: %w", err)
	}
	return nil
}
