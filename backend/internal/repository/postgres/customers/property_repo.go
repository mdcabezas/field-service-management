package postgrescustomers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/model/customers"
	"localis-backend/internal/service"
)

type PropertyRepo struct {
	pool *pgxpool.Pool
}

func NewPropertyRepo(pool *pgxpool.Pool) *PropertyRepo {
	return &PropertyRepo{pool: pool}
}

func (r *PropertyRepo) GetByID(ctx context.Context, id uuid.UUID) (*customers.Property, error) {
	var p customers.Property
	err := r.pool.QueryRow(ctx,
		`SELECT id, customer_address_id, name, type, notes, created_at
		 FROM customers.properties
		 WHERE id = $1`, id,
	).Scan(&p.ID, &p.CustomerAddressID, &p.Name, &p.Type, &p.Notes, &p.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "property", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get property by id: %w", err)
	}
	return &p, nil
}

func (r *PropertyRepo) ListByCustomerAddress(ctx context.Context, caID uuid.UUID) ([]customers.Property, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, customer_address_id, name, type, notes, created_at
		 FROM customers.properties
		 WHERE customer_address_id = $1
		 ORDER BY created_at`, caID,
	)
	if err != nil {
		return nil, fmt.Errorf("list properties: %w", err)
	}
	defer rows.Close()

	items := make([]customers.Property, 0)
	for rows.Next() {
		var p customers.Property
		if err := rows.Scan(&p.ID, &p.CustomerAddressID, &p.Name, &p.Type, &p.Notes, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan property: %w", err)
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *PropertyRepo) Create(ctx context.Context, prop *customers.Property) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO customers.properties (id, customer_address_id, name, type, notes, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		prop.ID, prop.CustomerAddressID, prop.Name, prop.Type, prop.Notes, prop.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create property: %w", err)
	}
	return nil
}

func (r *PropertyRepo) Update(ctx context.Context, id uuid.UUID, prop *customers.Property) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE customers.properties
		 SET name = $2, type = $3, notes = $4
		 WHERE id = $1`,
		id, prop.Name, prop.Type, prop.Notes,
	)
	if err != nil {
		return fmt.Errorf("update property: %w", err)
	}
	return nil
}

func (r *PropertyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM customers.properties WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete property: %w", err)
	}
	return nil
}
