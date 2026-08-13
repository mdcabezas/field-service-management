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

type VisitRepo struct {
	pool *pgxpool.Pool
}

func NewVisitRepo(pool *pgxpool.Pool) *VisitRepo {
	return &VisitRepo{pool: pool}
}

func (r *VisitRepo) Pool() *pgxpool.Pool {
	return r.pool
}

func (r *VisitRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.Visit, error) {
	var v operations.Visit
	err := r.pool.QueryRow(ctx,
		`SELECT id, property_id, partner_id, partner_order_id, parent_visit_id, route_id, daily_plan_id,
		        type, status, priority, source, source_reference, billing_to, partner_supervisor,
		        result, rejection_reason_id, rejection_reason_detail,
		        scheduled_at, started_at, completed_at, alternative_location, notes,
		        created_at, updated_at
		 FROM operations.visits
		 WHERE id = $1`, id,
	).Scan(&v.ID, &v.PropertyID, &v.PartnerID, &v.PartnerOrderID, &v.ParentVisitID, &v.RouteID, &v.DailyPlanID,
		&v.Type, &v.Status, &v.Priority, &v.Source, &v.SourceReference, &v.BillingTo, &v.PartnerSupervisor,
		&v.Result, &v.RejectionReasonID, &v.RejectionReasonDetail,
		&v.ScheduledAt, &v.StartedAt, &v.CompletedAt, &v.AlternativeLocation, &v.Notes,
		&v.CreatedAt, &v.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit by id: %w", err)
	}
	return &v, nil
}

func (r *VisitRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[operations.Visit], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM operations.visits`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count visits: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, property_id, partner_id, partner_order_id, parent_visit_id, route_id, daily_plan_id,
		        type, status, priority, source, source_reference, billing_to, partner_supervisor,
		        result, rejection_reason_id, rejection_reason_detail,
		        scheduled_at, started_at, completed_at, alternative_location, notes,
		        created_at, updated_at
		 FROM operations.visits
		 ORDER BY scheduled_at DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list visits: %w", err)
	}
	defer rows.Close()

	items := make([]operations.Visit, 0)
	for rows.Next() {
		var v operations.Visit
		if err := rows.Scan(&v.ID, &v.PropertyID, &v.PartnerID, &v.PartnerOrderID, &v.ParentVisitID, &v.RouteID, &v.DailyPlanID,
			&v.Type, &v.Status, &v.Priority, &v.Source, &v.SourceReference, &v.BillingTo, &v.PartnerSupervisor,
			&v.Result, &v.RejectionReasonID, &v.RejectionReasonDetail,
			&v.ScheduledAt, &v.StartedAt, &v.CompletedAt, &v.AlternativeLocation, &v.Notes,
			&v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan visit: %w", err)
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[operations.Visit]{Items: items, Total: total}, nil
}

func (r *VisitRepo) Create(ctx context.Context, visit *operations.Visit) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visits (id, property_id, partner_id, partner_order_id, parent_visit_id, route_id, daily_plan_id,
		        type, status, priority, source, source_reference, billing_to, partner_supervisor,
		        result, rejection_reason_id, rejection_reason_detail,
		        scheduled_at, started_at, completed_at, alternative_location, notes,
		        created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)`,
		visit.ID, visit.PropertyID, visit.PartnerID, visit.PartnerOrderID, visit.ParentVisitID, visit.RouteID, visit.DailyPlanID,
		visit.Type, visit.Status, visit.Priority, visit.Source, visit.SourceReference, visit.BillingTo, visit.PartnerSupervisor,
		visit.Result, visit.RejectionReasonID, visit.RejectionReasonDetail,
		visit.ScheduledAt, visit.StartedAt, visit.CompletedAt, visit.AlternativeLocation, visit.Notes,
		visit.CreatedAt, visit.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit: %w", err)
	}
	return nil
}

func (r *VisitRepo) Update(ctx context.Context, id uuid.UUID, visit *operations.Visit) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE operations.visits
		 SET property_id = $2, partner_id = $3, partner_order_id = $4, parent_visit_id = $5, route_id = $6, daily_plan_id = $7,
		     type = $8, status = $9, priority = $10, source = $11, source_reference = $12, billing_to = $13,
		     partner_supervisor = $14, result = $15, rejection_reason_id = $16, rejection_reason_detail = $17,
		     scheduled_at = $18, started_at = $19, completed_at = $20, alternative_location = $21, notes = $22,
		     updated_at = $23
		 WHERE id = $1`,
		id, visit.PropertyID, visit.PartnerID, visit.PartnerOrderID, visit.ParentVisitID, visit.RouteID, visit.DailyPlanID,
		visit.Type, visit.Status, visit.Priority, visit.Source, visit.SourceReference, visit.BillingTo, visit.PartnerSupervisor,
		visit.Result, visit.RejectionReasonID, visit.RejectionReasonDetail,
		visit.ScheduledAt, visit.StartedAt, visit.CompletedAt, visit.AlternativeLocation, visit.Notes,
		visit.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update visit: %w", err)
	}
	return nil
}

func (r *VisitRepo) UpdateInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, visit *operations.Visit) error {
	_, err := tx.Exec(ctx,
		`UPDATE operations.visits
		 SET property_id = $2, partner_id = $3, partner_order_id = $4, parent_visit_id = $5, route_id = $6, daily_plan_id = $7,
		     type = $8, status = $9, priority = $10, source = $11, source_reference = $12, billing_to = $13,
		     partner_supervisor = $14, result = $15, rejection_reason_id = $16, rejection_reason_detail = $17,
		     scheduled_at = $18, started_at = $19, completed_at = $20, alternative_location = $21, notes = $22,
		     updated_at = $23
		 WHERE id = $1`,
		id, visit.PropertyID, visit.PartnerID, visit.PartnerOrderID, visit.ParentVisitID, visit.RouteID, visit.DailyPlanID,
		visit.Type, visit.Status, visit.Priority, visit.Source, visit.SourceReference, visit.BillingTo, visit.PartnerSupervisor,
		visit.Result, visit.RejectionReasonID, visit.RejectionReasonDetail,
		visit.ScheduledAt, visit.StartedAt, visit.CompletedAt, visit.AlternativeLocation, visit.Notes,
		visit.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update visit in tx: %w", err)
	}
	return nil
}

func (r *VisitRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visits WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit: %w", err)
	}
	return nil
}
