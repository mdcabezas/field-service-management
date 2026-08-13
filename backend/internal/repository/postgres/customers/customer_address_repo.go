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

type CustomerAddressRepo struct {
	pool *pgxpool.Pool
}

func NewCustomerAddressRepo(pool *pgxpool.Pool) *CustomerAddressRepo {
	return &CustomerAddressRepo{pool: pool}
}

func (r *CustomerAddressRepo) GetByID(ctx context.Context, id uuid.UUID) (*customers.CustomerAddress, error) {
	var ca customers.CustomerAddress
	err := r.pool.QueryRow(ctx,
		`SELECT id, customer_id, address_id, name, type, created_at
		 FROM customers.customer_addresses
		 WHERE id = $1`, id,
	).Scan(&ca.ID, &ca.CustomerID, &ca.AddressID, &ca.Name, &ca.Type, &ca.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &service.NotFoundError{Resource: "customer_address", ID: id.String()}
	}
	if err != nil {
		return nil, fmt.Errorf("get customer_address by id: %w", err)
	}
	return &ca, nil
}

func (r *CustomerAddressRepo) ListByCustomer(ctx context.Context, customerID uuid.UUID) ([]customers.CustomerAddress, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, customer_id, address_id, name, type, created_at
		 FROM customers.customer_addresses
		 WHERE customer_id = $1
		 ORDER BY created_at`, customerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list customer addresses: %w", err)
	}
	defer rows.Close()

	items := make([]customers.CustomerAddress, 0)
	for rows.Next() {
		var ca customers.CustomerAddress
		if err := rows.Scan(&ca.ID, &ca.CustomerID, &ca.AddressID, &ca.Name, &ca.Type, &ca.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan customer address: %w", err)
		}
		items = append(items, ca)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return items, nil
}

func (r *CustomerAddressRepo) Create(ctx context.Context, addr *customers.CustomerAddress) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO customers.customer_addresses (id, customer_id, address_id, name, type, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		addr.ID, addr.CustomerID, addr.AddressID, addr.Name, addr.Type, addr.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create customer address: %w", err)
	}
	return nil
}

func (r *CustomerAddressRepo) Update(ctx context.Context, id uuid.UUID, addr *customers.CustomerAddress) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE customers.customer_addresses
		 SET name = $2, type = $3
		 WHERE id = $1`,
		id, addr.Name, addr.Type,
	)
	if err != nil {
		return fmt.Errorf("update customer address: %w", err)
	}
	return nil
}

func (r *CustomerAddressRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM customers.customer_addresses WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete customer address: %w", err)
	}
	return nil
}
