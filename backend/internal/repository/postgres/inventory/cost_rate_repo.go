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

type CostRateRepo struct {
	pool *pgxpool.Pool
}

func NewCostRateRepo(pool *pgxpool.Pool) *CostRateRepo {
	return &CostRateRepo{pool: pool}
}

func (r *CostRateRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.CostRate, error) {
	var c inventory.CostRate
	err := r.pool.QueryRow(ctx,
		`SELECT id, type, reference_id, reference_type, value, unit, valid_from, valid_until, created_at
		 FROM inventory.cost_rates
		 WHERE id = $1`, id,
	).Scan(&c.ID, &c.Type, &c.ReferenceID, &c.ReferenceType, &c.Value, &c.Unit, &c.ValidFrom, &c.ValidUntil, &c.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "cost_rate", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get cost_rate by id: %w", err)
	}
	return &c, nil
}

func (r *CostRateRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[inventory.CostRate], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.cost_rates`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count cost rates: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, type, reference_id, reference_type, value, unit, valid_from, valid_until, created_at
		 FROM inventory.cost_rates
		 ORDER BY valid_from DESC
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list cost rates: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.CostRate, 0)
	for rows.Next() {
		var c inventory.CostRate
		if err := rows.Scan(&c.ID, &c.Type, &c.ReferenceID, &c.ReferenceType, &c.Value, &c.Unit, &c.ValidFrom, &c.ValidUntil, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan cost rate: %w", err)
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[inventory.CostRate]{Items: items, Total: total}, nil
}

func (r *CostRateRepo) Create(ctx context.Context, rate *inventory.CostRate) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.cost_rates (id, type, reference_id, reference_type, value, unit, valid_from, valid_until, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		rate.ID, rate.Type, rate.ReferenceID, rate.ReferenceType, rate.Value, rate.Unit, rate.ValidFrom, rate.ValidUntil, rate.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create cost rate: %w", err)
	}
	return nil
}

func (r *CostRateRepo) Update(ctx context.Context, id uuid.UUID, rate *inventory.CostRate) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inventory.cost_rates
		 SET type = $2, reference_id = $3, reference_type = $4, value = $5, unit = $6, valid_from = $7, valid_until = $8
		 WHERE id = $1`,
		id, rate.Type, rate.ReferenceID, rate.ReferenceType, rate.Value, rate.Unit, rate.ValidFrom, rate.ValidUntil,
	)
	if err != nil {
		return fmt.Errorf("update cost rate: %w", err)
	}
	return nil
}

func (r *CostRateRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inventory.cost_rates WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete cost rate: %w", err)
	}
	return nil
}
