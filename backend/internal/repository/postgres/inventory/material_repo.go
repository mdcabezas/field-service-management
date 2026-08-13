package postgresinventory

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/inventory"
	"localis-backend/internal/repository"
	"localis-backend/internal/service"
)

type MaterialRepo struct {
	pool *pgxpool.Pool
}

func NewMaterialRepo(pool *pgxpool.Pool) *MaterialRepo {
	return &MaterialRepo{pool: pool}
}

func (r *MaterialRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.Material, error) {
	var m inventory.Material
	err := r.pool.QueryRow(ctx,
		`SELECT id, sku, name, unit, unit_cost, created_at
		 FROM inventory.materials
		 WHERE id = $1`, id,
	).Scan(&m.ID, &m.SKU, &m.Name, &m.Unit, &m.UnitCost, &m.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "material", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get material by id: %w", err)
	}
	return &m, nil
}

func (r *MaterialRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[inventory.Material], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.materials`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count materials: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, sku, name, unit, unit_cost, created_at
		 FROM inventory.materials
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list materials: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.Material, 0)
	for rows.Next() {
		var m inventory.Material
		if err := rows.Scan(&m.ID, &m.SKU, &m.Name, &m.Unit, &m.UnitCost, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan material: %w", err)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[inventory.Material]{Items: items, Total: total}, nil
}

func (r *MaterialRepo) Create(ctx context.Context, mat *inventory.Material) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.materials (id, sku, name, unit, unit_cost, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		mat.ID, mat.SKU, mat.Name, mat.Unit, mat.UnitCost, mat.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create material: %w", err)
	}
	return nil
}

func (r *MaterialRepo) Update(ctx context.Context, id uuid.UUID, mat *inventory.Material) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inventory.materials
		 SET sku = $2, name = $3, unit = $4, unit_cost = $5
		 WHERE id = $1`,
		id, mat.SKU, mat.Name, mat.Unit, mat.UnitCost,
	)
	if err != nil {
		return fmt.Errorf("update material: %w", err)
	}
	return nil
}

func (r *MaterialRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inventory.materials WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete material: %w", err)
	}
	return nil
}
