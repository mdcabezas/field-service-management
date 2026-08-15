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

type VisitTypeRepo struct {
	pool *pgxpool.Pool
}

func NewVisitTypeRepo(pool *pgxpool.Pool) *VisitTypeRepo {
	return &VisitTypeRepo{pool: pool}
}

func (r *VisitTypeRepo) GetByID(ctx context.Context, id uuid.UUID) (*operations.VisitType, error) {
	var vt operations.VisitType
	err := r.pool.QueryRow(ctx,
		`SELECT id, code, name, active
		 FROM operations.visit_types
		 WHERE id = $1`, id,
	).Scan(&vt.ID, &vt.Code, &vt.Name, &vt.Active)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "visit_type", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get visit_type by id: %w", err)
	}
	return &vt, nil
}

func (r *VisitTypeRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[operations.VisitType], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM operations.visit_types`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count visit types: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, code, name, active
		 FROM operations.visit_types
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list visit types: %w", err)
	}
	defer rows.Close()

	items := make([]operations.VisitType, 0)
	for rows.Next() {
		var vt operations.VisitType
		if err := rows.Scan(&vt.ID, &vt.Code, &vt.Name, &vt.Active); err != nil {
			return nil, fmt.Errorf("scan visit type: %w", err)
		}
		items = append(items, vt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[operations.VisitType]{Items: items, Total: total}, nil
}

func (r *VisitTypeRepo) Create(ctx context.Context, vt *operations.VisitType) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO operations.visit_types (id, code, name, active)
		 VALUES ($1, $2, $3, $4)`,
		vt.ID, vt.Code, vt.Name, vt.Active,
	)
	if err != nil {
		return fmt.Errorf("create visit type: %w", err)
	}
	return nil
}

func (r *VisitTypeRepo) Update(ctx context.Context, id uuid.UUID, vt *operations.VisitType) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE operations.visit_types
		 SET code = $2, name = $3, active = $4
		 WHERE id = $1`,
		id, vt.Code, vt.Name, vt.Active,
	)
	if err != nil {
		return fmt.Errorf("update visit type: %w", err)
	}
	return nil
}

func (r *VisitTypeRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM operations.visit_types WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete visit type: %w", err)
	}
	return nil
}