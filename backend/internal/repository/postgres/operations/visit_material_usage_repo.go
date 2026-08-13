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

type VisitMaterialUsageRepo struct {
	pool *pgxpool.Pool
}

func NewVisitMaterialUsageRepo(pool *pgxpool.Pool) *VisitMaterialUsageRepo {
	return &VisitMaterialUsageRepo{pool: pool}
}

func (r *VisitMaterialUsageRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitMaterialUsage, error) {
	var m operations.VisitMaterialUsage
	err := r.pool.QueryRow(ctx,
		`SELECT id, visit_id, material_id, quantity, notes, created_at
		 FROM operations.visit_material_usages
		 WHERE id = $1`, id,
	).Scan(&m.ID, &m.VisitID, &m.MaterialID, &m.Quantity, &m.Notes, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_material_usage", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_material_usage by id: %w", err)
	}
	return &m, nil
}

func (r *VisitMaterialUsageRepo) ListByVisit(ctx context.Context, visitID uuid.UUID) ([]operations.VisitMaterialUsage, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, visit_id, material_id, quantity, notes, created_at
		 FROM operations.visit_material_usages
		 WHERE visit_id = $1
		 ORDER BY created_at`, visitID,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit material usages by visit: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitMaterialUsage, 0)
	for rows.Next() {
		var m operations.VisitMaterialUsage
		if err := rows.Scan(&m.ID, &m.VisitID, &m.MaterialID, &m.Quantity, &m.Notes, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan visit material usage: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *VisitMaterialUsageRepo) Create(ctx context.Context, m *operations.VisitMaterialUsage) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_material_usages (id, visit_id, material_id, quantity, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		m.ID, m.VisitID, m.MaterialID, m.Quantity, m.Notes, m.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create visit material usage: %w", err)
	}
	return nil
}

func (r *VisitMaterialUsageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_material_usages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit material usage: %w", err)
	}
	return nil
}
