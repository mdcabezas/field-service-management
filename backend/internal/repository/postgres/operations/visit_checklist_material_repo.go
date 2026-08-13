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

type VisitChecklistMaterialRepo struct {
	pool *pgxpool.Pool
}

func NewVisitChecklistMaterialRepo(pool *pgxpool.Pool) *VisitChecklistMaterialRepo {
	return &VisitChecklistMaterialRepo{pool: pool}
}

func (r *VisitChecklistMaterialRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitChecklistMaterial, error) {
	var m operations.VisitChecklistMaterial
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, material_id, planned_quantity, confirmed, notes, created_at
		 FROM operations.visit_checklist_materials
		 WHERE id = $1`, id,
	).Scan(&m.ID, &m.VisitID, &m.MaterialID, &m.PlannedQuantity, &m.Confirmed, &m.Notes, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_checklist_material", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_checklist_material by id: %w", err)
	}
	return &m, nil
}

func (r *VisitChecklistMaterialRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitChecklistMaterial, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, material_id, planned_quantity, confirmed, notes, created_at
		 FROM operations.visit_checklist_materials
		 WHERE visit_id = $1
		 ORDER BY created_at
		 LIMIT 100`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit checklist materials by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitChecklistMaterial, 0)
	for rows.Next() {
		var m operations.VisitChecklistMaterial
		if err := rows.Scan(&m.ID, &m.VisitID, &m.MaterialID, &m.PlannedQuantity, &m.Confirmed, &m.Notes, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit checklist material: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitChecklistMaterialRepo) Create(ctx context.Context, m *operations.VisitChecklistMaterial) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_checklist_materials (id, visit_id, material_id, planned_quantity, confirmed, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		m.ID, m.VisitID, m.MaterialID, m.PlannedQuantity, m.Confirmed, m.Notes, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit checklist material: %w", err)
	}
	return nil
}

func (r *VisitChecklistMaterialRepo) CreateBatch(ctx context.Context, items []operations.VisitChecklistMaterial) error {
	const maxBatchSize = 1000

	if len(items) == 0 {
		return nil
	}
	if len(items) > maxBatchSize {
		return fmt.Errorf("batch size %d exceeds maximum %d", len(items), maxBatchSize)
	}
	batch := &pgx.Batch{}
	for _, m := range items {
		batch.Queue(
			`INSERT INTO operations.visit_checklist_materials (id, visit_id, material_id, planned_quantity, confirmed, notes, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			m.ID, m.VisitID, m.MaterialID, m.PlannedQuantity, m.Confirmed, m.Notes, m.CreatedAt,
		)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	for range items {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("batch create visit checklist materials: %w", err)
		}
	}
	return nil
}

func (r *VisitChecklistMaterialRepo) Update(ctx context.Context, id uuid.UUID, m *operations.VisitChecklistMaterial) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE operations.visit_checklist_materials
		 SET visit_id = $2, material_id = $3, planned_quantity = $4, confirmed = $5, notes = $6
		 WHERE id = $1`,
		id, m.VisitID, m.MaterialID, m.PlannedQuantity, m.Confirmed, m.Notes,
	)
	if err != nil {
		return fmt.Errorf("update visit checklist material: %w", err)
	}
	return nil
}

func (r *VisitChecklistMaterialRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_checklist_materials WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit checklist material: %w", err)
	}
	return nil
}
