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

type RentalRepo struct {
	pool *pgxpool.Pool
}

func NewRentalRepo(pool *pgxpool.Pool) *RentalRepo {
	return &RentalRepo{pool: pool}
}

func (r *RentalRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory.Rental, error) {
	var rt inventory.Rental
	err := r.pool.QueryRow(ctx,
		`SELECT id, type, supplier, item_description, daily_cost, hourly_cost, fixed_cost, notes, created_at
		 FROM inventory.rentals
		 WHERE id = $1`, id,
	).Scan(&rt.ID, &rt.Type, &rt.Supplier, &rt.ItemDescription, &rt.DailyCost, &rt.HourlyCost, &rt.FixedCost, &rt.Notes, &rt.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "rental", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get rental by id: %w", err)
	}
	return &rt, nil
}

func (r *RentalRepo) List(ctx context.Context, limit, offset int) (*repository.ListResult[inventory.Rental], error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM inventory.rentals`).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count rentals: %w", err)
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, type, supplier, item_description, daily_cost, hourly_cost, fixed_cost, notes, created_at
		 FROM inventory.rentals
		 ORDER BY item_description
		 LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list rentals: %w", err)
	}
	defer rows.Close()

	items := make([]inventory.Rental, 0)
	for rows.Next() {
		var rt inventory.Rental
		if err := rows.Scan(&rt.ID, &rt.Type, &rt.Supplier, &rt.ItemDescription, &rt.DailyCost, &rt.HourlyCost, &rt.FixedCost, &rt.Notes, &rt.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan rental: %w", err)
		}
		items = append(items, rt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return &repository.ListResult[inventory.Rental]{Items: items, Total: total}, nil
}

func (r *RentalRepo) Create(ctx context.Context, rental *inventory.Rental) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.rentals (id, type, supplier, item_description, daily_cost, hourly_cost, fixed_cost, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		rental.ID, rental.Type, rental.Supplier, rental.ItemDescription, rental.DailyCost, rental.HourlyCost, rental.FixedCost, rental.Notes, rental.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create rental: %w", err)
	}
	return nil
}

func (r *RentalRepo) Update(ctx context.Context, id uuid.UUID, rental *inventory.Rental) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE inventory.rentals
		 SET type = $2, supplier = $3, item_description = $4, daily_cost = $5, hourly_cost = $6, fixed_cost = $7, notes = $8
		 WHERE id = $1`,
		id, rental.Type, rental.Supplier, rental.ItemDescription, rental.DailyCost, rental.HourlyCost, rental.FixedCost, rental.Notes,
	)
	if err != nil {
		return fmt.Errorf("update rental: %w", err)
	}
	return nil
}

func (r *RentalRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM inventory.rentals WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete rental: %w", err)
	}
	return nil
}
