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

type VisitEPPUsageRepo struct {
	pool *pgxpool.Pool
}

func NewVisitEPPUsageRepo(pool *pgxpool.Pool) *VisitEPPUsageRepo {
	return &VisitEPPUsageRepo{pool: pool}
}

func (r *VisitEPPUsageRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitEPPUsage, error) {
	var e operations.VisitEPPUsage
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, epp_id, quantity, status, notes, created_at
		 FROM operations.visit_epp_usages
		 WHERE id = $1`, id,
	).Scan(&e.ID, &e.VisitID, &e.EPPID, &e.Quantity, &e.Status, &e.Notes, &e.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_epp_usage", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_epp_usage by id: %w", err)
	}
	return &e, nil
}

func (r *VisitEPPUsageRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitEPPUsage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, epp_id, quantity, status, notes, created_at
		 FROM operations.visit_epp_usages
		 WHERE visit_id = $1
		 ORDER BY created_at`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit epp usages by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitEPPUsage, 0)
	for rows.Next() {
		var e operations.VisitEPPUsage
		if err := rows.Scan(&e.ID, &e.VisitID, &e.EPPID, &e.Quantity, &e.Status, &e.Notes, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit epp usage: %w", err)
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitEPPUsageRepo) Create(ctx context.Context, e *operations.VisitEPPUsage) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_epp_usages (id, visit_id, epp_id, quantity, status, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		e.ID, e.VisitID, e.EPPID, e.Quantity, e.Status, e.Notes, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit epp usage: %w", err)
	}
	return nil
}

func (r *VisitEPPUsageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_epp_usages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit epp usage: %w", err)
	}
	return nil
}
