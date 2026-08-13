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

type VisitChecklistEPPRepo struct {
	pool *pgxpool.Pool
}

func NewVisitChecklistEPPRepo(pool *pgxpool.Pool) *VisitChecklistEPPRepo {
	return &VisitChecklistEPPRepo{pool: pool}
}

func (r *VisitChecklistEPPRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitChecklistEPP, error) {
	var e operations.VisitChecklistEPP
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, epp_id, planned_quantity, confirmed, notes, created_at
		 FROM operations.visit_checklist_epps
		 WHERE id = $1`, id,
	).Scan(&e.ID, &e.VisitID, &e.EPPID, &e.PlannedQuantity, &e.Confirmed, &e.Notes, &e.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_checklist_epp", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_checklist_epp by id: %w", err)
	}
	return &e, nil
}

func (r *VisitChecklistEPPRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitChecklistEPP, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, epp_id, planned_quantity, confirmed, notes, created_at
		 FROM operations.visit_checklist_epps
		 WHERE visit_id = $1
		 ORDER BY created_at
		 LIMIT 100`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit checklist epps by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitChecklistEPP, 0)
	for rows.Next() {
		var e operations.VisitChecklistEPP
		if err := rows.Scan(&e.ID, &e.VisitID, &e.EPPID, &e.PlannedQuantity, &e.Confirmed, &e.Notes, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit checklist epp: %w", err)
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitChecklistEPPRepo) Create(ctx context.Context, e *operations.VisitChecklistEPP) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_checklist_epps (id, visit_id, epp_id, planned_quantity, confirmed, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		e.ID, e.VisitID, e.EPPID, e.PlannedQuantity, e.Confirmed, e.Notes, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit checklist epp: %w", err)
	}
	return nil
}

func (r *VisitChecklistEPPRepo) CreateBatch(ctx context.Context, items []operations.VisitChecklistEPP) error {
	const maxBatchSize = 1000

	if len(items) == 0 {
		return nil
	}
	if len(items) > maxBatchSize {
		return fmt.Errorf("batch size %d exceeds maximum %d", len(items), maxBatchSize)
	}
	batch := &pgx.Batch{}
	for _, e := range items {
		batch.Queue(
			`INSERT INTO operations.visit_checklist_epps (id, visit_id, epp_id, planned_quantity, confirmed, notes, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			e.ID, e.VisitID, e.EPPID, e.PlannedQuantity, e.Confirmed, e.Notes, e.CreatedAt,
		)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range items {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("batch create visit checklist epps: %w", err)
		}
	}
	return nil
}

func (r *VisitChecklistEPPRepo) Update(ctx context.Context, id uuid.UUID, e *operations.VisitChecklistEPP) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE operations.visit_checklist_epps
		 SET visit_id = $2, epp_id = $3, planned_quantity = $4, confirmed = $5, notes = $6
		 WHERE id = $1`,
		id, e.VisitID, e.EPPID, e.PlannedQuantity, e.Confirmed, e.Notes,
	)
	if err != nil {
		return fmt.Errorf("update visit checklist epp: %w", err)
	}
	return nil
}

func (r *VisitChecklistEPPRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_checklist_epps WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit checklist epp: %w", err)
	}
	return nil
}
