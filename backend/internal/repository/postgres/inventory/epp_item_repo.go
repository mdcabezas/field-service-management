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

type EPPItemRepo struct {
	pool *pgxpool.Pool
}

func NewEPPItemRepo(pool *pgxpool.Pool) *EPPItemRepo {
	return &EPPItemRepo{pool: pool}
}

func (r *EPPItemRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.EPPItem, error) {
	var e inventory.EPPItem
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, type, lifecycle, created_at
		 FROM inventory.epp_items
		 WHERE id = $1`, id,
	).Scan(&e.ID, &e.Name, &e.Type, &e.Lifecycle, &e.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "epp_item", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get epp_item by id: %w", err)
	}
	return &e, nil
}

func (r *EPPItemRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[inventory.EPPItem], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.epp_items`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count epp items: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, name, type, lifecycle, created_at
		 FROM inventory.epp_items
		 ORDER BY name
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list epp items: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.EPPItem, 0)
	for rows.Next() {
		var e inventory.EPPItem
		if err := rows.Scan(&e.ID, &e.Name, &e.Type, &e.Lifecycle, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan epp item: %w", err)
		}
		items = append(items, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[inventory.EPPItem]{Items: items, Total: total}, nil
}

func (r *EPPItemRepo) Create(ctx context.Context, item *inventory.EPPItem) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.epp_items (id, name, type, lifecycle, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		item.ID, item.Name, item.Type, item.Lifecycle, item.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create epp item: %w", err)
	}
	return nil
}

func (r *EPPItemRepo) Update(ctx context.Context, id uuid.UUID, item *inventory.EPPItem) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inventory.epp_items
		 SET name = $2, type = $3, lifecycle = $4
		 WHERE id = $1`,
		id, item.Name, item.Type, item.Lifecycle,
	)
	if err != nil {
		return fmt.Errorf("update epp item: %w", err)
	}
	return nil
}

func (r *EPPItemRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inventory.epp_items WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete epp item: %w", err)
	}
	return nil
}
