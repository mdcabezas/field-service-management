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

type VisitSLATrackingRepo struct {
	pool *pgxpool.Pool
}

func NewVisitSLATrackingRepo(pool *pgxpool.Pool) *VisitSLATrackingRepo {
	return &VisitSLATrackingRepo{pool: pool}
}

func (r *VisitSLATrackingRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitSLATracking, error) {
	var t operations.VisitSLATracking
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, sla_id, requested_at, responded_at, resolved_at,
		        response_time_hours, resolution_time_hours, meets_response_sla, meets_resolution_sla, created_at
		 FROM operations.visit_sla_trackings
		 WHERE id = $1`, id,
	).Scan(&t.ID, &t.VisitID, &t.SLAID, &t.RequestedAt, &t.RespondedAt, &t.ResolvedAt,
		&t.ResponseTimeHours, &t.ResolutionTimeHours, &t.MeetsResponseSLA, &t.MeetsResolutionSLA, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_sla_tracking", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_sla_tracking by id: %w", err)
	}
	return &t, nil
}

func (r *VisitSLATrackingRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitSLATracking, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, sla_id, requested_at, responded_at, resolved_at,
		        response_time_hours, resolution_time_hours, meets_response_sla, meets_resolution_sla, created_at
		 FROM operations.visit_sla_trackings
		 WHERE visit_id = $1
		 ORDER BY created_at`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit sla trackings by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitSLATracking, 0)
	for rows.Next() {
		var t operations.VisitSLATracking
		if err := rows.Scan(&t.ID, &t.VisitID, &t.SLAID, &t.RequestedAt, &t.RespondedAt, &t.ResolvedAt,
			&t.ResponseTimeHours, &t.ResolutionTimeHours, &t.MeetsResponseSLA, &t.MeetsResolutionSLA, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit sla tracking: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitSLATrackingRepo) Create(ctx context.Context, t *operations.VisitSLATracking) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_sla_trackings (id, visit_id, sla_id, requested_at, responded_at, resolved_at,
		        response_time_hours, resolution_time_hours, meets_response_sla, meets_resolution_sla, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		t.ID, t.VisitID, t.SLAID, t.RequestedAt, t.RespondedAt, t.ResolvedAt,
		t.ResponseTimeHours, t.ResolutionTimeHours, t.MeetsResponseSLA, t.MeetsResolutionSLA, t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit sla tracking: %w", err)
	}
	return nil
}

func (r *VisitSLATrackingRepo) Update(ctx context.Context, id uuid.UUID, t *operations.VisitSLATracking) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE operations.visit_sla_trackings
		 SET visit_id = $2, sla_id = $3, requested_at = $4, responded_at = $5, resolved_at = $6,
		     response_time_hours = $7, resolution_time_hours = $8, meets_response_sla = $9, meets_resolution_sla = $10
		 WHERE id = $1`,
		id, t.VisitID, t.SLAID, t.RequestedAt, t.RespondedAt, t.ResolvedAt,
		t.ResponseTimeHours, t.ResolutionTimeHours, t.MeetsResponseSLA, t.MeetsResolutionSLA,
	)
	if err != nil {
		return fmt.Errorf("update visit sla tracking: %w", err)
	}
	return nil
}

func (r *VisitSLATrackingRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_sla_trackings WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit sla tracking: %w", err)
	}
	return nil
}

func (r *VisitSLATrackingRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[operations.VisitSLATracking], error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, sla_id, requested_at, responded_at, resolved_at,
		        response_time_hours, resolution_time_hours, meets_response_sla, meets_resolution_sla, created_at
		 FROM operations.visit_sla_trackings
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit sla trackings: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitSLATracking, 0)
	for rows.Next() {
		var t operations.VisitSLATracking
		if err := rows.Scan(&t.ID, &t.VisitID, &t.SLAID, &t.RequestedAt, &t.RespondedAt, &t.ResolvedAt,
			&t.ResponseTimeHours, &t.ResolutionTimeHours, &t.MeetsResponseSLA, &t.MeetsResolutionSLA, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit sla tracking: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	var total int
	err = r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM operations.visit_sla_trackings`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count visit sla trackings: %w", err)
	}

	return &repository.ListResult[operations.VisitSLATracking]{Items: items, Total: total}, nil
}
